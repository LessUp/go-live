# Delta Spec: WHIP/WHEP Protocol

## MODIFIED Requirements

### Requirement: WHIP Publishing — Success Status Code

The system SHALL return `201 Created` on a successful WHIP publish, with the SDP answer
in the response body and `Content-Type: application/sdp`.

> **Replaces:** the vague "returns SDP answer in response body" wording in the base spec
> and the `200 OK` shown in RFC 0002's flow diagram.

#### Scenario: Successful publish returns 201
- **WHEN** publisher sends a valid SDP offer to `POST /api/whip/publish/{room}` with a valid auth token
- **THEN** system returns `201 Created` with `Content-Type: application/sdp` and the SDP answer in the body

### Requirement: WHIP Publishing — Duplicate Publisher

The system SHALL return `409 Conflict` when a publisher attempts to publish to a room
that already has an active publisher.

> **Replaces:** the existing "403 Forbidden" wording for the duplicate-publisher scenario
> in the base spec. The correct semantic is Conflict, not Forbidden.

#### Scenario: Room already has publisher returns 409
- **WHEN** a publisher sends a valid WHIP publish request to a room that already has an active publisher
- **THEN** system returns `409 Conflict` with a JSON error body `{"error": "publisher already exists in this room"}`

### Requirement: WHEP Playback — Success Status Code

The system SHALL return `201 Created` on a successful WHEP play, with the SDP answer in
the response body and `Content-Type: application/sdp`.

#### Scenario: Successful playback returns 201
- **WHEN** viewer sends a valid SDP offer to `POST /api/whep/play/{room}` with a valid auth token and an active publisher is present
- **THEN** system returns `201 Created` with `Content-Type: application/sdp` and the SDP answer in the body

### Requirement: WHEP Playback — No Active Publisher

The system SHALL return `404 Not Found` when a viewer attempts to play from a room that
has no active publisher, or from a room that does not exist.

> **Replaces:** the existing "returns appropriate error" wording for the no-publisher
> scenario in the base spec.

#### Scenario: No publisher in room returns 404
- **WHEN** viewer sends a valid WHEP play request to a room with no active publisher
- **THEN** system returns `404 Not Found` with a JSON error body `{"error": "no active publisher in room"}`

#### Scenario: Room does not exist returns 404
- **WHEN** viewer sends a valid WHEP play request to a room name that has never been used
- **THEN** system returns `404 Not Found` with a JSON error body `{"error": "no active publisher in room"}`

### Requirement: WHEP Playback — Subscriber Limit

The system SHALL return `403 Forbidden` when `MAX_SUBS_PER_ROOM` is configured and the
room has reached its subscriber limit.

> **Replaces:** "403 Forbidden" wording retained from the base spec, but now linked to
> the sentinel error `ErrSubscriberLimitReached` for precise handler discrimination.

#### Scenario: Room subscriber limit reached returns 403
- **WHEN** `MAX_SUBS_PER_ROOM` is set and the room already has that many subscribers and a new viewer attempts to subscribe
- **THEN** system returns `403 Forbidden` with a JSON error body `{"error": "subscriber limit reached"}`
