---
layout: default
title: API Reference
nav_order: 4
lang: en
---

# API Reference

This document details all HTTP API endpoints of live-webrtc-go.

{: .no_toc }

## Table of Contents

{: .no_toc .text-delta }

1. TOC
{:toc}

---

## Authentication

The system supports three authentication methods, tried in order of priority:

### 1. Bearer Token

```http
Authorization: Bearer <token>
```

### 2. X-Auth-Token Header

```http
X-Auth-Token: <token>
```

### 3. URL Query Parameter

```http
GET /api/rooms?token=<token>
```

### Authentication Priority

1. Room-level Token (`ROOM_TOKENS`)
2. Global Token (`AUTH_TOKEN`)
3. JWT (`JWT_SECRET`)
4. No authentication (when not configured)

---

## Streaming Endpoints

### WHIP Publish

Publish media stream to a specified room.

```http
POST /api/whip/publish/{room}
Content-Type: application/sdp
Authorization: Bearer <token>
```

**Parameters**

| Parameter | Location | Type | Required | Description |
|-----------|----------|------|----------|-------------|
| `room` | path | string | Yes | Room name, matches `^[A-Za-z0-9_-]{1,64}$` |

**Request Body**

SDP Offer (text/plain)

**Response**

| Status Code | Description |
|-------------|-------------|
| 200 | Success, returns SDP Answer |
| 400 | Invalid SDP or room name |
| 401 | Authentication failed |
| 409 | Room already has a publisher |
| 429 | Rate limit exceeded |

**Example**

```bash
curl -X POST "http://localhost:8080/api/whip/publish/demo" \
  -H "Content-Type: application/sdp" \
  -H "Authorization: Bearer mytoken" \
  --data-binary @offer.sdp
```

---

### WHEP Play

Subscribe to media stream from a specified room.

```http
POST /api/whep/play/{room}
Content-Type: application/sdp
Authorization: Bearer <token>
```

**Parameters**

| Parameter | Location | Type | Required | Description |
|-----------|----------|------|----------|-------------|
| `room` | path | string | Yes | Room name |

**Request Body**

SDP Offer (text/plain)

**Response**

| Status Code | Description |
|-------------|-------------|
| 200 | Success, returns SDP Answer |
| 400 | Invalid SDP or room name |
| 401 | Authentication failed |
| 404 | Room not found |
| 429 | Rate limit exceeded |

**Example**

```bash
curl -X POST "http://localhost:8080/api/whep/play/demo" \
  -H "Content-Type: application/sdp" \
  -H "Authorization: Bearer mytoken" \
  --data-binary @offer.sdp
```

---

## Query Endpoints

### Get Bootstrap Configuration

Returns runtime configuration required by the frontend application.

```http
GET /api/bootstrap
```

**Response**

```json
{
  "authEnabled": true,
  "recordEnabled": true,
  "iceServers": [
    {
      "urls": ["stun:stun.l.google.com:19302"]
    }
  ],
  "features": {
    "rooms": true,
    "records": true
  }
}
```

**Field Descriptions**

| Field | Type | Description |
|-------|------|-------------|
| `authEnabled` | boolean | Whether authentication is enabled |
| `recordEnabled` | boolean | Whether recording is enabled |
| `iceServers` | array | ICE server configuration |
| `features` | object | Feature flags |

---

### Get Room List

Returns all active rooms and their status.

```http
GET /api/rooms
```

**Response**

```json
[
  {
    "name": "demo",
    "hasPublisher": true,
    "subscriberCount": 5
  },
  {
    "name": "test",
    "hasPublisher": false,
    "subscriberCount": 0
  }
]
```

**Field Descriptions**

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Room name |
| `hasPublisher` | boolean | Whether room has a publisher |
| `subscriberCount` | number | Number of subscribers |

---

### Get Recording List

Returns metadata for all recording files.

```http
GET /api/records
```

**Response**

```json
[
  {
    "name": "demo_video0_1710123456.ivf",
    "size": 1048576,
    "modTime": "2024-03-10T12:34:56Z"
  }
]
```

**Field Descriptions**

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | File name |
| `size` | number | File size in bytes |
| `modTime` | string | Modification time (ISO 8601) |

---

## Admin Endpoints

### Close Room

Forcefully close a specified room, disconnecting all connections.

```http
POST /api/admin/rooms/{room}/close
Authorization: Bearer <admin-token>
```

**Parameters**

| Parameter | Location | Type | Required | Description |
|-----------|----------|------|----------|-------------|
| `room` | path | string | Yes | Room name to close |

**Response**

| Status Code | Description |
|-------------|-------------|
| 200 | Successfully closed |
| 401 | Authentication failed (requires Admin Token) |
| 404 | Room not found |

**Example**

```bash
curl -X POST \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  http://localhost:8080/api/admin/rooms/demo/close
```

---

## Health and Metrics

### Health Check

```http
GET /healthz
```

**Response**

```
ok
```

Status code: 200

---

### Prometheus Metrics

```http
GET /metrics
```

**Available Metrics**

| Metric Name | Type | Labels | Description |
|-------------|------|--------|-------------|
| `live_rooms` | Gauge | - | Active room count |
| `live_subscribers` | GaugeVec | `room` | Subscribers per room |
| `live_rtp_bytes_total` | CounterVec | `room` | Total RTP bytes |
| `live_rtp_packets_total` | CounterVec | `room` | Total RTP packets |

**Example**

```bash
curl http://localhost:8080/metrics
```

---

## Error Responses

All error responses follow this format:

```json
{
  "error": "Error description"
}
```

### Common Error Codes

| Status Code | Error Message | Reason |
|-------------|---------------|--------|
| 400 | `invalid room name` | Invalid room name format |
| 400 | `invalid SDP` | Invalid SDP format |
| 401 | `unauthorized` | Authentication failed or not provided |
| 404 | `room not found` | Room does not exist (during WHEP play) |
| 409 | `publisher already exists` | Room already has a publisher |
| 429 | `too many requests` | Rate limit triggered |
| 500 | `internal server error` | Internal server error |

---

## CORS Configuration

All API responses include CORS headers:

```http
Access-Control-Allow-Origin: <ALLOWED_ORIGIN>
Access-Control-Allow-Methods: GET, POST, OPTIONS
Access-Control-Allow-Headers: Content-Type, Authorization, X-Auth-Token
```

Preflight requests (OPTIONS) automatically return 204.

---

## Request Limits

| Limit | Value | Configuration |
|-------|-------|---------------|
| SDP request body size | 1 MB | Hardcoded |
| Room name length | 1-64 characters | Regex `^[A-Za-z0-9_-]{1,64}$` |
| Request rate | Configurable | `RATE_LIMIT_RPS`, `RATE_LIMIT_BURST` |
| Subscribers per room | Configurable | `MAX_SUBS_PER_ROOM` |

---

## Room Name Rules

- **Allowed characters**: `A-Z`, `a-z`, `0-9`, `_`, `-`
- **Maximum length**: 64 characters
- **Pattern**: `^[A-Za-z0-9_-]{1,64}$`

**Valid examples**: `room1`, `my-room`, `live_stream_01`

**Invalid examples**: `my room`, `room@123`, `a` (too short if min>1)
