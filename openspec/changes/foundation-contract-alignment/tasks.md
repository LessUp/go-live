# Tasks: Foundation Contract Alignment

## 1. Spec

- [ ] 1.1 Review and approve `openspec/changes/foundation-contract-alignment/specs/whip-whep/spec.md` delta
- [ ] 1.2 Review and approve `openspec/changes/foundation-contract-alignment/specs/api/spec.md` delta
- [ ] 1.3 Review and approve `openspec/changes/foundation-contract-alignment/specs/recording/spec.md` delta
- [ ] 1.4 Create `openspec/specs/api/openapi.yaml` with OpenAPI 3.0.3 definitions covering all endpoints, request/response schemas, and status codes declared in the delta specs

## 2. Tests

- [ ] 2.1 Add unit test: `ServeWHIPPublish` with duplicate publisher returns 409
- [ ] 2.2 Add unit test: `ServeWHEPPlay` with subscriber limit reached returns 403
- [ ] 2.3 Add unit tests: `ServeWHEPPlay` returns 404 for both no-publisher code paths:
  - 2.3a room is absent from the manager (`manager.Subscribe` returns `ErrNoPublisher` — room was never created)
  - 2.3b room exists but publisher is nil (`room.Subscribe` returns `ErrNoPublisher` — room created, publisher disconnected or not yet connected)
- [ ] 2.4 Add unit test: `ServeWHIPPublish` success returns 201 with `Content-Type: application/sdp`
- [ ] 2.5 Add unit test: `ServeWHEPPlay` success returns 201 with `Content-Type: application/sdp`
- [ ] 2.6 Add unit test: `GET /api/rooms` returns JSON array with lower-camelCase fields (`name`, `hasPublisher`, `tracks`, `subscribers`)
- [ ] 2.7 Add unit test: `GET /api/records` returns JSON array with fields `name`, `size`, `modTime`, `url`
- [ ] 2.8 Verify all existing handler tests still pass after status code changes (`make test-unit`)

## 3. Implementation

- [ ] 3.1 Add sentinel errors to `internal/sfu/room.go`:
  - `var ErrPublisherExists = errors.New("publisher already exists in this room")`
  - `var ErrSubscriberLimitReached = errors.New("subscriber limit reached")`
  - `var ErrNoPublisher = errors.New("no active publisher in room")`
- [ ] 3.2 Replace inline `errors.New("publisher already exists in this room")` with `ErrPublisherExists` in `room.Publish`
- [ ] 3.3 Replace `fmt.Errorf("subscriber limit reached")` with `ErrSubscriberLimitReached` in `room.Subscribe`
- [ ] 3.4 Add check in `room.Subscribe` (inside `r.mu.Lock()` scope): if `r.publisher == nil`, return `ErrNoPublisher` before proceeding. This handles the case where the room exists but no publisher has connected yet.
- [ ] 3.5 Change `manager.Subscribe` to NOT call `ensureRoom` for the subscribe path. Instead, perform a read-locked map lookup (matching `GetRoom`). If the room is absent, return `ErrNoPublisher` directly. `ensureRoom` remains on the publish path only. This eliminates the auto-create-on-subscribe behaviour.
- [ ] 3.6 Update `ServeWHIPPublish` in `internal/api/handlers.go`:
  - map `sfu.ErrPublisherExists` → 409 Conflict
  - retain 400 for all other errors
- [ ] 3.7 Update `ServeWHEPPlay` in `internal/api/handlers.go`:
  - map `sfu.ErrNoPublisher` → 404 Not Found
  - map `sfu.ErrSubscriberLimitReached` → 403 Forbidden
  - retain 400 for all other errors
- [ ] 3.8 Add JSON struct tags to `RoomInfo` in `internal/sfu/manager.go`:
  `json:"name"`, `json:"hasPublisher"`, `json:"tracks"`, `json:"subscribers"`

## 4. Docs

- [ ] 4.1 Update `README.md` API table: WHIP/WHEP success code column shows `201`; error conditions documented (409/403/404)
- [ ] 4.2 Update `RFC 0002` flow diagram comment: change `200 OK` to `201 Created`
- [ ] 4.3 Update `openspec/specs/api/spec.md` Room schema and Recording schema to match delta specs
- [ ] 4.4 Update `openspec/specs/whip-whep/spec.md` base spec to absorb the delta after this change is archived

## 5. Verification

- [ ] 5.1 Run `make lint` — must pass with no new violations
- [ ] 5.2 Run `make test-unit` — all tests green
- [ ] 5.3 Run `make test-integration` — all tests green
- [ ] 5.4 Run `make security` — no new findings
- [ ] 5.5 Manually verify: duplicate publish → curl returns HTTP 409
- [ ] 5.6 Manually verify: subscribe with no publisher → curl returns HTTP 404
- [ ] 5.7 Manually verify: room list JSON keys are lower-camelCase
- [ ] 5.8 Archive this change with `/opsx:archive` after all tasks pass
