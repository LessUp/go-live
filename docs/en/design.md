---
layout: default
title: Design Documentation
nav_order: 3
lang: en
---

# Design Documentation

This document describes the system architecture, module breakdown, data flow, and extension points of live-webrtc-go for architectural review and secondary development.

{: .no_toc }

## Table of Contents

{: .no_toc .text-delta }

1. TOC
{:toc}

---

## System Architecture

### Architecture Overview

```
                          ┌─────────────────────────────────────┐
                          │          HTTP Server :8080          │
                          │                                     │
     ┌──────────┐         │  ┌─────────┐    ┌─────────────┐    │
     │ Publisher│ ──WHIP──▶│  │  Auth   │───▶│    SFU      │    │
     │(OBS/Web) │         │  │Middleware│    │  Manager    │    │
     └──────────┘         │  └─────────┘    └──────┬──────┘    │
                          │                        │           │
     ┌──────────┐         │        ┌───────────────┤           │
     │  Viewer  │ ◀──WHEP──│        │               │           │
     │(Browser) │───────▶  │        ▼               ▼           │
     └──────────┘         │ ┌─────────────┐ ┌───────────┐      │
                          │ │    Room     │ │ Recording │      │
                          │ │  (Fanout)   │ │  & Upload │      │
                          │ └──────┬──────┘ └─────┬─────┘      │
                          └────────┼──────────────┼────────────┘
                                   │              │
                         ┌─────────▼──────────────▼───────────┐
                         │          Object Storage            │
                         │           (S3/MinIO)               │
                         └────────────────────────────────────┘
```

### Request Processing Chain

```
HTTP Request
    │
    ▼
┌─────────────┐
│    CORS     │ ← ALLOWED_ORIGIN
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ Rate Limiter│ ← RATE_LIMIT_RPS, RATE_LIMIT_BURST
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Auth      │ ← AUTH_TOKEN / ROOM_TOKENS / JWT_SECRET
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Handler   │ → Business Logic
└─────────────┘
```

---

## Core Concepts

### Room

The Room is the core abstraction of the SFU. Each room:
- Has at most one Publisher
- Can have multiple Subscribers
- Has independent Track Fanout logic
- Can have its own authentication token

### Track Fanout

When a publisher pushes media tracks, the system creates a Track Fanout:
- Read RTP packets from the publisher's PeerConnection
- Copy and distribute to all subscribers
- Optionally write to recording files

### PeerConnection

Each WebRTC connection:
- Publisher: Receives media tracks
- Subscriber: Sends media tracks
- ICE negotiation completed through WHIP/WHEP protocols

---

## Module Details

### cmd/server

**Responsibility**: Application entry point, service initialization

```go
// main.go main flow
1. config.Load()           // Load configuration
2. uploader.Init()         // Initialize uploader
3. sfu.NewManager()        // Create room manager
4. api.NewHTTPHandlers()   // Create HTTP handlers
5. RegisterRoutes()        // Register routes
6. otel.InitTracer()       // Initialize tracing
7. http.Server.Listen()    // Start server
8. Graceful Shutdown       // Graceful exit
```

### internal/config

**Responsibility**: Environment variable parsing and defaults

```
┌─────────────────────────────────────┐
│             Config                  │
├─────────────────────────────────────┤
│ HTTPAddr        string              │
│ AllowedOrigin   string              │
│ AuthToken       string              │
│ RoomTokens      map[string]string   │
│ JWTSecret       string              │
│ RecordEnabled   bool                │
│ RecordDir       string              │
│ S3Endpoint      string              │
│ RateLimitRPS    float64             │
│ STUN/TURN       []string            │
│ ...                                 │
└─────────────────────────────────────┘
```

### internal/api

**Responsibility**: HTTP request handling

| File | Function |
|------|----------|
| `handlers.go` | WHIP/WHEP/Rooms/Records/Admin endpoint handling |
| `middleware.go` | CORS, rate limiting, Token/JWT authentication |
| `routes.go` | URL routing, parameter extraction, room name validation |

**Authentication Priority**:
```
1. Room-specific Token (ROOM_TOKENS)
    ↓ (not found or failed)
2. Global Token (AUTH_TOKEN)
    ↓ (not found or failed)
3. JWT (JWT_SECRET)
    ↓ (not found or failed)
4. Allow (no auth configured)
```

### internal/sfu

**Responsibility**: WebRTC SFU core logic

```
┌─────────────────────────────────────────────────────┐
│                     Manager                          │
│  - Manage all Room instances                         │
│  - Create/Delete Room                                │
│  - Count rooms                                       │
└──────────────────────┬──────────────────────────────┘
                       │ 1:N
                       ▼
┌─────────────────────────────────────────────────────┐
│                      Room                            │
│  - Publisher PeerConnection                          │
│  - Subscriber PeerConnections                        │
│  - TrackFeeds (TrackFanout map)                      │
└──────────────────────┬──────────────────────────────┘
                       │ 1:N
                       ▼
┌─────────────────────────────────────────────────────┐
│                  TrackFanout                         │
│  - Remote Track (from publisher)                     │
│  - Local Tracks (to subscribers)                     │
│  - readLoop: RTP distribution                        │
│  - Optional: Recorder (IVF/OGG writer)               │
└─────────────────────────────────────────────────────┘
```

**Key Methods**:

| Method | Purpose |
|--------|---------|
| `Manager.Publish()` | Create room, establish publisher connection |
| `Manager.Subscribe()` | Create subscriber connection, bind existing tracks |
| `Room.attachTrackFeed()` | Distribute new track to all subscribers |
| `TrackFanout.readLoop()` | RTP packet reading and distribution loop |

### internal/metrics

**Responsibility**: Prometheus metrics exposure

| Metric | Type | Description |
|--------|------|-------------|
| `live_rooms` | Gauge | Active room count |
| `live_subscribers` | GaugeVec | Subscribers per room |
| `rtp_bytes_total` | CounterVec | Total RTP bytes |
| `rtp_packets_total` | CounterVec | Total RTP packets |

### internal/uploader

**Responsibility**: S3/MinIO file upload

```
Upload Flow:
1. Check Enabled() → client != nil
2. Open local file
3. Build object key (prefix + filename)
4. client.PutObject()
5. (Optional) Delete local file
```

---

## Data Flow

### Publishing Flow (WHIP)

```
1. Publisher → POST /api/whip/publish/{room}
   │
2. HTTPHandlers.ServeWHIPPublish()
   │  ├─ CORS check
   │  ├─ Rate limit check
   │  └─ Authentication check
   │
3. Manager.Publish(roomName, sdpOffer)
   │  ├─ getOrCreateRoom()
   │  └─ Room.Publish(sdpOffer)
   │
4. Room.Publish()
   │  ├─ Create MediaEngine + Interceptors
   │  ├─ NewPeerConnection(ICEConfig)
   │  ├─ SetRemoteDescription(offer)
   │  ├─ CreateAnswer()
   │  ├─ SetLocalDescription(answer)
   │  └─ OnTrack: attachTrackFeed()
   │
5. Return SDP Answer
   │
6. TrackFanout.readLoop() runs continuously
   │  ├─ Read RTP from Remote Track
   │  ├─ Write to recorder (if enabled)
   │  └─ Distribute to all Local Tracks
```

### Playback Flow (WHEP)

```
1. Viewer → POST /api/whep/play/{room}
   │
2. HTTPHandlers.ServeWHEPPlay()
   │  ├─ CORS/Rate limit/Authentication check
   │  └─ Manager.Subscribe()
   │
3. Manager.Subscribe(roomName, sdpOffer)
   │  └─ Room.Subscribe(sdpOffer)
   │
4. Room.Subscribe()
   │  ├─ Check subscriber limit
   │  ├─ NewPeerConnection()
   │  ├─ Iterate existing TrackFeeds
   │  │   └─ TrackFanout.attachToSubscriber()
   │  ├─ SetRemoteDescription/CreateAnswer
   │  └─ OnICEStateChange: removeSubscriber()
   │
5. Return SDP Answer
```

### Disconnection

```
ICE State Change (Failed/Disconnected/Closed)
    │
    ▼
┌─────────────────────────────────────┐
│ Publisher Disconnect                 │
├─────────────────────────────────────┤
│ 1. closePublisher()                  │
│ 2. Close all TrackFanouts            │
│ 3. Upload recording files            │
│ 4. Clear subscriber list             │
│ 5. pruneIfEmpty()                    │
└─────────────────────────────────────┘

┌─────────────────────────────────────┐
│ Subscriber Disconnect                │
├─────────────────────────────────────┤
│ 1. removeSubscriber()                │
│ 2. Remove binding from TrackFanouts  │
│ 3. Close PeerConnection              │
│ 4. pruneIfEmpty()                    │
└─────────────────────────────────────┘
```

---

## Authentication System

### Token Authentication

```
Priority 1: Room Token (ROOM_TOKENS)
┌─────────────────────────────────────┐
│ ROOM_TOKENS="room1:abc;room2:def"   │
│                                      │
│ Access room1 → Check token == "abc" │
│ Access room2 → Check token == "def" │
│ Access room3 → Fallback to Global Token │
└─────────────────────────────────────┘

Priority 2: Global Token (AUTH_TOKEN)
┌─────────────────────────────────────┐
│ AUTH_TOKEN="secret123"              │
│                                      │
│ All rooms use the same token         │
└─────────────────────────────────────┘
```

### JWT Authentication

```go
// JWT Claims structure
type roomClaims struct {
    Room  string `json:"room,omitempty"`   // Restrict to room
    Role  string `json:"role,omitempty"`   // "admin" role
    Admin any    `json:"admin,omitempty"`  // true/1 for admin
    jwt.RegisteredClaims
}

// Use cases
1. Room access: claims.Room == room or claims.Room == ""
2. Admin API: claims.Role == "admin" or claims.Admin == true
```

---

## Recording and Upload

### Recording Format

| Codec | File Format | Writer |
|-------|-------------|--------|
| Opus | .ogg | oggwriter (48kHz, 2ch) |
| VP8 | .ivf | ivfwriter |
| VP9 | .ivf | ivfwriter |

### File Naming

```
{room}_{trackID}_{unixTimestamp}.{ext}

Example: demo_video0_1710123456.ivf
```

### Upload Flow

```
Room.closePublisher()
    │
    ▼
TrackFanout.close() → Return recording file path
    │
    ▼
uploader.Enabled()?
    │ Yes
    ▼
go uploadRecording(path)
    │
    ▼
Upload(ctx, path)
    ├─ PutObject(S3Bucket, objectKey, file)
    └─ (Optional) os.Remove(localFile)
```

---

## Observability

### Prometheus Metrics

```
# Active room count
live_rooms

# Subscribers per room
live_subscribers{room="demo"}

# RTP bytes (cumulative)
live_rtp_bytes_total{room="demo"}

# RTP packets (cumulative)
live_rtp_packets_total{room="demo"}
```

### OpenTelemetry Tracing

```
Environment variables:
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317
OTEL_EXPORTER_OTLP_PROTOCOL=grpc
OTEL_SERVICE_NAME=live-webrtc-go

Traced spans:
- HTTP Handler: {method} {path}
```

### Health Check

```
GET /healthz → "ok" (200 OK)
```

---

## Extension Points

### 1. Multi-instance Deployment

Current room state is in memory. For multi-instance deployment:
- External storage (Redis/Database) for room mappings
- Session affinity (Sticky Session)
- Or client-side redirection

### 2. Media Processing

Can be inserted before `TrackFanout.readLoop()`:
- Transcoding (FFmpeg integration)
- Simulcast / Multiple bitrates
- Screenshots / Watermarks

### 3. Authentication Extensions

Extend in `middleware.go`:
- OAuth2 integration
- Webhook callback validation
- IP whitelist

### 4. Storage Extensions

Implement `rtpWriter` interface:
```go
type rtpWriter interface {
    WriteRTP(*rtp.Packet) error
    Close() error
}
```

Supports:
- Real-time repackaging (MP4)
- Streaming upload (without local storage)
- CDN push

---

## Performance Considerations

### Memory Usage

- Each TrackFanout: ~1-2 MB (RTP buffer)
- Each subscriber: ~1500 bytes (MTU buffer)
- Recording buffer: Depends on write frequency

### CPU Usage

- RTP packet processing: Main loop in `readLoop()`
- Codec negotiation: Only during connection establishment
- Metrics update: Per RTP packet

### Optimization Suggestions

1. Zero-copy RTP forwarding (requires TrackFanout modification)
2. Batch metrics updates
3. Connection pooling (multi-room scenarios)
4. SIMD optimization (large number of subscribers)
