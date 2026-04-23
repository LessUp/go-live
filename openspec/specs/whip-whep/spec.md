# WHIP/WHEP Protocol

## Purpose

WebRTC-HTTP Ingestion Protocol (WHIP) and WebRTC-HTTP Egress Protocol (WHEP) for stream publishing and playback over HTTP. Enables compatibility with OBS Studio and modern browsers through simple HTTP-based signaling.

## Requirements

### Requirement: WHIP Publishing

The system SHALL accept WebRTC streams via WHIP protocol at endpoint `POST /api/whip/publish/{room}`.

#### Scenario: Successful publish
- **WHEN** publisher sends SDP offer to POST /api/whip/publish/{room} with valid auth token
- **THEN** system returns SDP answer in response body and creates PeerConnection for publisher

#### Scenario: Publisher authentication failure
- **WHEN** publisher sends request without valid auth token
- **THEN** system returns 401 Unauthorized

#### Scenario: Invalid room name
- **WHEN** publisher sends request with room name not matching ^[A-Za-z0-9_-]{1,64}$
- **THEN** system returns 400 Bad Request

#### Scenario: Room already has publisher
- **WHEN** publisher attempts to publish to room with existing publisher
- **THEN** system returns 403 Forbidden

### Requirement: WHEP Playback

The system SHALL allow viewers to play WebRTC streams via WHEP protocol at endpoint `POST /api/whep/play/{room}`.

#### Scenario: Successful playback
- **WHEN** viewer sends SDP offer to POST /api/whep/play/{room} with valid auth token
- **THEN** system returns SDP answer and creates PeerConnection for viewer

#### Scenario: Viewer authentication failure
- **WHEN** viewer sends request without valid auth token
- **THEN** system returns 401 Unauthorized

#### Scenario: Room full
- **WHEN** MAX_SUBS_PER_ROOM is configured and limit is reached
- **THEN** system returns 403 Forbidden

#### Scenario: No publisher in room
- **WHEN** viewer attempts to play from room with no publisher
- **THEN** system returns appropriate error

### Requirement: SDP Exchange

The system SHALL handle SDP offer/answer exchange following WHIP/WHEP protocol.

#### Scenario: SDP offer accepted
- **WHEN** client sends valid SDP offer in request body
- **THEN** system creates PeerConnection, sets remote description, and returns SDP answer

#### Scenario: SDP size limit
- **WHEN** SDP request body exceeds 1MB
- **THEN** system rejects the request

### Requirement: CORS Support

The system SHALL provide CORS headers for browser-based clients.

#### Scenario: Browser CORS request
- **WHEN** browser sends cross-origin request
- **THEN** system returns appropriate CORS headers based on ALLOWED_ORIGIN config
