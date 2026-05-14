# API Endpoints

Complete reference for all HTTP API endpoints.

## Authentication

The system supports two authentication methods, tried in order of priority:

### 1. Bearer Token

```http
Authorization: Bearer <token>
```

### 2. X-Auth-Token Header

```http
X-Auth-Token: <token>
```

### Authentication Priority

1. Room-level Token (`ROOM_TOKENS`)
2. Global Token (`AUTH_TOKEN`)
3. JWT (`JWT_SECRET`)
4. No authentication (when not configured)

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

**Request Body**: SDP Offer (text/plain)

**Response**

| Status Code | Description |
|-------------|-------------|
| 201 | Success, returns SDP Answer |
| 400 | Invalid SDP or room name |
| 401 | Authentication failed |
| 409 | Room already has a publisher |
| 429 | Rate limit exceeded |

```bash
curl -X POST "http://localhost:8080/api/whip/publish/demo" \
  -H "Content-Type: application/sdp" \
  -H "Authorization: Bearer mytoken" \
  --data-binary @offer.sdp
```

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

**Request Body**: SDP Offer (text/plain)

**Response**

| Status Code | Description |
|-------------|-------------|
| 201 | Success, returns SDP Answer |
| 400 | Invalid SDP or room name |
| 401 | Authentication failed |
| 403 | Subscriber limit reached (`MAX_SUBS_PER_ROOM`) |
| 404 | No active publisher in room |
| 429 | Rate limit exceeded |

```bash
curl -X POST "http://localhost:8080/api/whep/play/demo" \
  -H "Content-Type: application/sdp" \
  -H "Authorization: Bearer mytoken" \
  --data-binary @offer.sdp
```

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
    "tracks": 2,
    "subscribers": 5
  }
]
```

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
    "modTime": "2024-03-10T12:34:56Z",
    "url": "/records/demo_video0_1710123456.ivf"
  }
]
```

## Admin Endpoints

### Close Room

Forcefully close a specified room, disconnecting all connections.

```http
POST /api/admin/rooms/{room}/close
Authorization: Bearer <admin-token>
```

**Response**

| Status Code | Description |
|-------------|-------------|
| 200 | Successfully closed |
| 401 | Authentication failed (requires Admin Token) |
| 404 | Room not found |

## Health and Metrics

### Health Check

```http
GET /healthz
```

Returns `ok` with status code 200.

### Prometheus Metrics

```http
GET /metrics
```

| Metric Name | Type | Labels | Description |
|-------------|------|--------|-------------|
| `live_rooms` | Gauge | - | Active room count |
| `live_subscribers` | GaugeVec | `room` | Subscribers per room |
| `live_rtp_bytes_total` | CounterVec | `room` | Total RTP bytes |
| `live_rtp_packets_total` | CounterVec | `room` | Total RTP packets |

## Error Responses

Domain-specific WHIP/WHEP errors use this JSON format:

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
| 403 | `subscriber limit reached` | `MAX_SUBS_PER_ROOM` limit hit |
| 404 | `no active publisher in room` | No publisher in room |
| 409 | `publisher already exists in this room` | Room already has a publisher |
| 429 | `too many requests` | Rate limit triggered |
| 500 | `internal server error` | Internal server error |

## Request Limits

| Limit | Value | Configuration |
|-------|-------|---------------|
| SDP request body size | 1 MB | Hardcoded |
| Room name length | 1-64 characters | Regex `^[A-Za-z0-9_-]{1,64}$` |
| Request rate | Configurable | `RATE_LIMIT_RPS`, `RATE_LIMIT_BURST` |
| Subscribers per room | Configurable | `MAX_SUBS_PER_ROOM` |

## CORS Configuration

All API responses include CORS headers:

```http
Access-Control-Allow-Origin: <ALLOWED_ORIGIN>
Access-Control-Allow-Methods: GET, POST, OPTIONS
Access-Control-Allow-Headers: Content-Type, Authorization, X-Auth-Token
```
