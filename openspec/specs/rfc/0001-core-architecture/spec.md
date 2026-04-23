# RFC 0001: Core SFU Architecture

## Summary

This RFC documents the core architecture of the WebRTC SFU server, including the Manager → Room → TrackFanout hierarchy and the overall system design.

## Motivation

Document the architectural decisions to ensure consistency in implementation and provide reference for future contributors.

## Design

### System Architecture

```
┌─────────────┐   WHIP (POST /api/whip/publish/{room})   ┌──────────────┐
│  Publisher  │ ──────────────────────────────────────▶ │ HTTP Server  │
│  (OBS/Web)  │                                         │   :8080      │
└─────────────┘                                         └──────┬───────┘
                                                               │
                                      Creates PeerConnection   │
                                      & TrackFanout            ▼
                                                        ┌─────────────┐
┌─────────────┐   WHEP (POST /api/whep/play/{room})    │    Room     │
│   Viewer    │ ──────────────────────────────────────▶ │   (SFU)     │
│  (Browser)  │ ◀────────────────────────────────────── │             │
└───────────┘        RTP Packets (WebRTC)             └──────┬──────┘
                                                               │
                                          Records & Uploads    ▼
                                                        ┌─────────────┐
                                                        │ Object Store│
                                                        │  (S3/MinIO) │
                                                        └─────────────┘
```

### Core Hierarchy

```
Manager
  ├── Room 1
  │     ├── Publisher (PeerConnection)
  │     ├── Subscribers (map of PeerConnections)
  │     └── TrackFanout (per-track RTP distribution)
  ├── Room 2
  └── ...
```

### Package Organization

```
internal/
├── api/              # HTTP layer
│   ├── handlers.go   # WHIP/WHEP/Rooms/Records/Admin endpoints
│   ├── middleware.go # CORS, rate limiting, auth (token/JWT)
│   └── routes.go     # URL routing, room name validation
├── config/           # Environment variable configuration
├── sfu/              # Core WebRTC SFU logic
│   ├── manager.go    # Room lifecycle management
│   ├── room.go       # PeerConnection, track fanout, recording
│   └── track.go      # RTP packet distribution to subscribers
├── metrics/          # Prometheus gauges/counters
├── otel/             # OpenTelemetry tracer initialization
├── uploader/         # S3/MinIO upload client
└── testutil/         # Test helpers
```

### Concurrency Model

1. **Room State Protection**
   - Each Room has `sync.RWMutex`
   - Read operations use `RLock()`
   - Write operations use `Lock()`

2. **Goroutine Lifecycle**
   - `track.go` readLoop must exit on subscriber disconnect
   - Room cleanup must stop all track fanout goroutines
   - Graceful shutdown via context cancellation

3. **Track Fanout Pattern**
   ```go
   // Publisher track receives RTP packets
   go readLoop(downTrack, subscriber)

   // Each subscriber has independent readLoop
   // Exit on: subscriber disconnect, room close, publisher disconnect
   ```

## Alternatives Considered

### MCU (Multipoint Control Unit)
- **Rejected**: MCU decodes and re-encodes all streams, high CPU cost
- **SFU Advantage**: Forward-only, minimal CPU, scales better

### Room Persistence
- **Rejected**: In-memory only for simplicity
- **Future**: Could add database for room/record persistence

## Status

✅ Approved - Architecture is stable and proven in production.
