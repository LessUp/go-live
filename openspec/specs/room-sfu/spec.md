# Room-based SFU Relay

## Purpose

A room-based Selective Forwarding Unit (SFU) where each room supports one publisher and multiple subscribers with efficient RTP packet forwarding.

## Requirements

### Requirement: Room Lifecycle

The system SHALL manage room lifecycle automatically without explicit creation.

#### Scenario: Room auto-creation
- **WHEN** first publisher joins a room
- **THEN** system creates the room automatically

#### Scenario: Room cleanup
- **WHEN** publisher disconnects from room
- **THEN** system destroys the room and closes all subscriber connections

#### Scenario: Room name validation
- **WHEN** client references a room name
- **THEN** system validates against pattern ^[A-Za-z0-9_-]{1,64}$

### Requirement: Publisher Management

The system SHALL allow exactly one publisher per room.

#### Scenario: Publisher joins
- **WHEN** publisher successfully establishes connection
- **THEN** publisher's media tracks are distributed to all subscribers

#### Scenario: Duplicate publisher
- **WHEN** second publisher attempts to join existing room
- **THEN** system rejects with appropriate error

#### Scenario: Publisher disconnect
- **WHEN** publisher disconnects
- **THEN** system triggers room cleanup and notifies all subscribers

### Requirement: Subscriber Management

The system SHALL support multiple subscribers per room with optional limits.

#### Scenario: Subscriber joins
- **WHEN** viewer successfully establishes connection
- **THEN** viewer receives all tracks from publisher via RTP

#### Scenario: Subscriber limit
- **WHEN** MAX_SUBS_PER_ROOM is configured and limit reached
- **THEN** new subscriber requests are rejected

#### Scenario: Subscriber disconnect
- **WHEN** individual subscriber disconnects
- **THEN** only that subscriber is affected, others continue normally

### Requirement: RTP Forwarding

The system SHALL forward RTP packets efficiently from publisher to all subscribers.

#### Scenario: Video forwarding
- **WHEN** publisher sends video RTP packets
- **THEN** system forwards packets to all subscriber connections

#### Scenario: Audio forwarding
- **WHEN** publisher sends audio RTP packets
- **THEN** system forwards packets to all subscriber connections

#### Scenario: Multiple codec support
- **WHEN** publisher uses VP8 or VP9 codec
- **THEN** system forwards without transcoding

### Requirement: Room Query

The system SHALL provide API to list active rooms.

#### Scenario: List rooms
- **WHEN** client calls GET /api/rooms
- **THEN** system returns list of rooms with publisher status and subscriber counts

### Requirement: Admin Room Close

The system SHALL provide API to force close a room.

#### Scenario: Force close room
- **WHEN** admin calls POST /api/admin/rooms/{room}/close with valid admin token
- **THEN** system closes room and disconnects all participants
