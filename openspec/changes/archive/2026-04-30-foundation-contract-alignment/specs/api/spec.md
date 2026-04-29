# Delta Spec: API

## MODIFIED Requirements

### Requirement: WHIP/WHEP Success Code

The system SHALL respond with `201 Created` for successful WHIP publish and WHEP play
requests. API documentation and the OpenAPI specification MUST reflect this code.

#### Scenario: WHIP publish endpoint returns 201
- **WHEN** client calls `POST /api/whip/publish/{room}` with a valid auth token and SDP offer and no publisher is already active
- **THEN** system returns `201 Created` with `Content-Type: application/sdp` and the SDP answer

#### Scenario: WHEP play endpoint returns 201
- **WHEN** client calls `POST /api/whep/play/{room}` with a valid auth token and SDP offer and a publisher is active
- **THEN** system returns `201 Created` with `Content-Type: application/sdp` and the SDP answer

### Requirement: Room List Response Model

The system SHALL return a JSON array from `GET /api/rooms`. Each element in the array
MUST conform to the following schema with lower-camelCase field names:

```json
{
  "name":         "string   — room identifier",
  "hasPublisher": "boolean  — true if an active publisher is connected",
  "tracks":       "integer  — number of active media tracks in the room",
  "subscribers":  "integer  — number of active subscriber connections"
}
```

> **Replaces:** the incorrect Room schema in the base `api/spec.md` that listed
> `subscriberCount` and `createdAt` fields which are not present in the implementation.

#### Scenario: Room list fields are lower-camelCase
- **WHEN** client calls `GET /api/rooms`
- **THEN** system returns a JSON array where each object has exactly the keys `name`, `hasPublisher`, `tracks`, `subscribers`

### Requirement: Recording List Response Model

The system SHALL return a JSON array from `GET /api/records`. Each element in the array
MUST conform to the following schema:

```json
{
  "name":    "string  — filename of the recording (e.g. room_trackID_ts.ivf)",
  "size":    "integer — file size in bytes",
  "modTime": "string  — last-modified timestamp in RFC 3339 / UTC (e.g. 2025-01-01T00:00:00Z)",
  "url":     "string  — relative URL to download the recording (e.g. /records/filename)"
}
```

> **Replaces:** the incorrect Recording schema in the base `api/spec.md` that listed
> `room`, `trackID`, `filename`, `size`, `createdAt` — fields that do not match the
> implementation's `ServeRecordsList` output.

#### Scenario: Recording list fields are correct
- **WHEN** client calls `GET /api/records` and recordings exist
- **THEN** system returns a JSON array where each object has exactly the keys `name`, `size`, `modTime`, `url`

> **Sort order:** see the recording delta spec for the authoritative sort-order requirement.

### Requirement: OpenAPI Source of Truth

The system's API contract SHALL be maintained in an OpenAPI 3.0.3 document at
`openspec/specs/api/openapi.yaml`. This file MUST be kept in sync with the narrative
spec and the implementation. It is the authoritative machine-readable reference for
client code generation and contract testing.

#### Scenario: OpenAPI file is present and valid
- **WHEN** a developer references the API contract
- **THEN** `openspec/specs/api/openapi.yaml` exists and validates against the OpenAPI 3.0.3 schema

## ADDED Requirements

### Requirement: Error Response Format

All error responses from the system SHALL use `Content-Type: application/json` and a
JSON body with a single `error` string field:

```json
{ "error": "human-readable error message" }
```

#### Scenario: Error body is JSON
- **WHEN** any API endpoint returns a 4xx or 5xx status code
- **THEN** the response body is a JSON object with an `error` field containing a human-readable message
