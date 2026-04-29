# Proposal: Foundation Contract Alignment

## Summary

Align the implementation, documentation, and specifications so that every observable
HTTP contract is stated once, in one authoritative place, without contradiction.

## Problem Statement

Three sources of truth have drifted apart:

| Source | Status |
|--------|--------|
| `internal/api/handlers.go` | Returns 201 for WHIP/WHEP, 400 for all SFU errors, no structured error discrimination |
| `openspec/specs/whip-whep/spec.md` | Vague on success code; says 200 in RFC 0002 flow diagram; says "appropriate error" for no-publisher |
| `openspec/specs/api/spec.md` | Room schema has `subscriberCount`/`createdAt` fields that do not exist in the `RoomInfo` struct; Recording schema has `room`/`trackID`/`filename`/`createdAt` but implementation emits `name`/`size`/`modTime`/`url` |

These mismatches make it impossible to write a test that matches both the spec and the
implementation, and they create confusion for clients integrating against the API.

## Goals

1. **Canonical WHIP/WHEP status codes** — declare 201 Created as the normative success
   code for both publish and play endpoints, matching the existing implementation.

2. **Discriminated error codes** — replace the current blanket 400 with semantically
   correct codes for identifiable error conditions:
   - `409 Conflict` when a second publisher attempts to claim a room that already has one
   - `403 Forbidden` when a subscriber limit is in force and has been reached
   - `404 Not Found` when a subscriber tries to play from a room that has no active publisher

3. **Aligned room response model** — declare the canonical JSON field names for the room
   list response using lower-camelCase (`name`, `hasPublisher`, `tracks`, `subscribers`),
   matching what the implementation must emit after JSON tags are added to `RoomInfo`.

4. **Aligned recording response model** — confirm the canonical JSON shape as
   `name`/`size`/`modTime`/`url`, matching what `ServeRecordsList` already emits.

5. **Authoritative OpenAPI source of truth** — declare `openspec/specs/api/openapi.yaml`
   as the canonical machine-readable contract that must be kept in sync with the spec
   and implementation.

## Out of Scope

The following are explicitly excluded from this change:

- **TURN timed credentials** — no credential rotation or time-limited TURN token endpoint
- **Webhook / event delivery** — no push notification of room or recording lifecycle events
- **Simulcast, SVC, or selective subscription** — no layered encoding or per-subscriber
  layer selection
- **RTMP, SRT, or HLS gateways** — no ingest or egress transcoding pipelines
- **Distributed room placement** — no multi-node cluster routing or room migration

## Motivation

This is a prerequisite for every subsequent feature. Until the contract is unambiguous,
integration tests cannot be both spec-conforming and green, and any new endpoint risks
introducing further drift.
