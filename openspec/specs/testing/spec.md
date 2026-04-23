# Test Specifications

## Purpose

BDD test specifications defining acceptance test scenarios for the live-webrtc-go server.

## Requirements

### Requirement: WHIP Publishing Tests

The system SHALL be tested for WHIP publishing scenarios.

#### Scenario: Successful publish with valid token
- **GIVEN** a running server at http://localhost:8080
- **AND** valid auth token "secret-token"
- **WHEN** POST to /api/whip/publish/test-room with SDP offer and Authorization header
- **THEN** response is 200 OK with SDP answer

#### Scenario: Publish fails with invalid token
- **GIVEN** a running server
- **WHEN** POST to /api/whip/publish/test-room with wrong token
- **THEN** response is 401 Unauthorized

#### Scenario: Publish fails with missing token
- **GIVEN** a running server
- **WHEN** POST without Authorization header
- **THEN** response is 401 Unauthorized

#### Scenario: Publish fails with invalid room name
- **GIVEN** a running server
- **WHEN** POST to /api/whip/publish/invalid room!
- **THEN** response is 400 Bad Request

### Requirement: WHEP Playback Tests

The system SHALL be tested for WHEP playback scenarios.

#### Scenario: Successful playback with valid token
- **GIVEN** room with active publisher
- **WHEN** POST to /api/whep/play/{room} with valid auth
- **THEN** response is 200 OK with SDP answer

#### Scenario: Playback fails when room full
- **GIVEN** MAX_SUBS_PER_ROOM configured and limit reached
- **WHEN** new subscriber attempts to join
- **THEN** response is 403 Forbidden

### Requirement: Authentication Tests

The system SHALL be tested for authentication scenarios.

#### Scenario: Per-room token works for specific room only
- **GIVEN** ROOM_TOKENS configures room1:tok1
- **WHEN** POST to /api/whip/publish/room1 with tok1
- **THEN** response is 200 OK
- **WHEN** POST to /api/whip/publish/room2 with tok1
- **THEN** response is 401 Unauthorized

#### Scenario: JWT room restriction enforced
- **GIVEN** JWT with room="room1" claim
- **WHEN** POST to /api/whip/publish/room1 with JWT
- **THEN** response is 200 OK
- **WHEN** POST to /api/whip/publish/room2 with same JWT
- **THEN** response is 401 Unauthorized

### Requirement: Room Management Tests

The system SHALL be tested for room management scenarios.

#### Scenario: List active rooms
- **GIVEN** rooms "room1" and "room2" exist with publishers
- **WHEN** GET /api/rooms
- **THEN** response contains both rooms with subscriber counts

#### Scenario: Admin force close room
- **GIVEN** room with active publisher and subscribers
- **WHEN** POST /api/admin/rooms/{room}/close with admin token
- **THEN** room is closed and all connections terminated

### Requirement: Publisher Disconnect Tests

The system SHALL be tested for publisher disconnect handling.

#### Scenario: Publisher disconnect triggers cleanup
- **GIVEN** room with active publisher and 2 subscribers
- **WHEN** publisher disconnects
- **THEN** room is cleaned up
- **AND** all subscribers are notified
- **AND** all subscriber PeerConnections are closed

## Test File Reference

Gherkin feature files are located at: `test/integration/*.feature`
