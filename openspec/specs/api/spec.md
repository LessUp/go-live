# API Specification

## Purpose

Reference to the OpenAPI specification defining all HTTP endpoints for the live-webrtc-go server.

## Change History

| Date | Change | Source |
|------|--------|--------|
| 2026-04-30 | API contract alignment (schemas, status codes, OpenAPI) | foundation-contract-alignment |

## Requirements

### Requirement: OpenAPI Specification Reference

The system SHALL maintain an authoritative API specification in OpenAPI 3.0.3 format.

#### Scenario: API spec location
- **WHEN** developer needs to reference API endpoints
- **THEN** OpenAPI spec is available at `openspec/specs/api/openapi.yaml`

### Requirement: Streaming Endpoints

The system SHALL provide WHIP and WHEP endpoints for stream publishing and playback.

#### Scenario: WHIP publish endpoint
- **WHEN** client calls POST /api/whip/publish/{room}
- **THEN** system accepts SDP offer and returns SDP answer with `201 Created`

#### Scenario: WHEP play endpoint
- **WHEN** client calls POST /api/whep/play/{room}
- **THEN** system accepts SDP offer and returns SDP answer with `201 Created`

### Requirement: Query Endpoints

The system SHALL provide query endpoints for room and recording information.

#### Scenario: List rooms
- **WHEN** client calls GET /api/rooms
- **THEN** system returns list of active rooms with metadata

#### Scenario: List recordings
- **WHEN** client calls GET /api/records
- **THEN** system returns list of recording files with metadata

### Requirement: Admin Endpoints

The system SHALL provide admin endpoints requiring admin authentication.

#### Scenario: Close room
- **WHEN** admin calls POST /api/admin/rooms/{room}/close with admin token
- **THEN** system force closes the room

### Requirement: Health Endpoints

The system SHALL provide health and metrics endpoints without authentication.

#### Scenario: Health check
- **WHEN** client calls GET /healthz
- **THEN** system returns health status

#### Scenario: Prometheus metrics
- **WHEN** client calls GET /metrics
- **THEN** system returns Prometheus metrics

## Reference

The authoritative API specification is maintained in OpenAPI 3.0.3 format at:

**File:** `openspec/specs/api/openapi.yaml`

## Authentication

### Token Auth
- Header: `Authorization: Bearer <token>`
- Global token via `AUTH_TOKEN`
- Per-room tokens via `ROOM_TOKENS`

### JWT Auth
- Header: `Authorization: Bearer <jwt>`
- HMAC signing (HS256/HS384/HS512)
- Claims: `role`, `room`, `aud`

### Admin Auth
- Header: `Authorization: Bearer <admin-token>`
- Token via `ADMIN_TOKEN`

## Schemas

### Error
```json
{
  "error": "string"
}
```

### Room
```json
{
  "name": "string",
  "hasPublisher": true,
  "tracks": 0,
  "subscribers": 0
}
```

### Recording
```json
{
  "name": "string",
  "size": 0,
  "modTime": "2025-01-01T00:00:00Z",
  "url": "/records/filename"
}
```
