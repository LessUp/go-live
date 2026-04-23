# Authentication System

## Purpose

Flexible authentication supporting global tokens, per-room tokens, and JWT-based authentication with role-based access control.

## Requirements

### Requirement: Global Token Auth

The system SHALL support a single global auth token for all rooms.

#### Scenario: Valid global token
- **WHEN** client sends Authorization: Bearer <AUTH_TOKEN>
- **THEN** system grants access to any room

#### Scenario: Invalid global token
- **WHEN** client sends wrong token
- **THEN** system returns 401 Unauthorized

### Requirement: Per-room Token Auth

The system SHALL support room-specific tokens that override global token.

#### Scenario: Valid room token
- **WHEN** client sends room-specific token from ROOM_TOKENS config
- **THEN** system grants access only to that specific room

#### Scenario: Room token wrong room
- **WHEN** client uses room1 token to access room2
- **THEN** system returns 401 Unauthorized

### Requirement: JWT Authentication

The system SHALL support JWT-based authentication with HMAC signing.

#### Scenario: Valid JWT with role
- **WHEN** client sends valid JWT signed with JWT_SECRET
- **THEN** system validates and grants access based on claims

#### Scenario: JWT admin role
- **WHEN** JWT contains role=admin or admin=true claim
- **THEN** system grants admin-level access

#### Scenario: JWT room restriction
- **WHEN** JWT contains room claim
- **THEN** system grants access only to that specific room

#### Scenario: JWT expired
- **WHEN** JWT has expired
- **THEN** system returns 401 Unauthorized

#### Scenario: JWT wrong audience
- **WHEN** JWT aud claim does not match JWT_AUDIENCE config
- **THEN** system returns 401 Unauthorized

### Requirement: Admin Token

The system SHALL require separate admin token for admin endpoints.

#### Scenario: Valid admin token
- **WHEN** client sends Authorization: Bearer <ADMIN_TOKEN> to admin endpoint
- **THEN** system grants admin access

#### Scenario: Missing admin token
- **WHEN** client calls admin endpoint without valid admin token
- **THEN** system returns 401 Unauthorized

### Requirement: Security

The system SHALL use timing-attack resistant token comparison.

#### Scenario: Token comparison
- **WHEN** comparing any token
- **THEN** system uses crypto/subtle.ConstantTimeCompare

#### Scenario: Token logging
- **WHEN** logging occurs
- **THEN** tokens and secrets are never logged
