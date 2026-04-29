# Design: Foundation Contract Alignment

## Context

See [RFC 0002](../../specs/rfc/0002-whip-whep-integration/spec.md) for the original
WHIP/WHEP integration design. This change makes no architectural additions; it only
closes the gap between that RFC, the specs, and the implementation.

## Decisions

### 1. WHIP/WHEP Success Code Is 201 Created

The implementation already returns `201 Created` for both `ServeWHIPPublish` and
`ServeWHEPPlay`. RFC 0002's flow diagram incorrectly shows `200 OK`. The spec delta
declares `201 Created` as normative. No code change is required for this decision.

### 2. Duplicate Publisher Returns 409 Conflict

**Current behaviour:** `room.Publish` returns `errors.New("publisher already exists in
this room")`; the handler maps all SFU errors to `400 Bad Request`.

**Target behaviour:** the handler detects this specific error and returns `409 Conflict`.

**Implementation note:** introduce a package-level sentinel in `internal/sfu/room.go`:
```go
var ErrPublisherExists = errors.New("publisher already exists in this room")
```
Replace the inline `errors.New(...)` with `ErrPublisherExists`. In
`internal/api/handlers.go` use `errors.Is(err, sfu.ErrPublisherExists)` to map to 409
before falling back to 400 for all other SFU errors.

### 3. Subscriber Limit Returns 403 Forbidden

**Current behaviour:** `room.Subscribe` returns `fmt.Errorf("subscriber limit reached")`;
the handler maps it to 400.

**Target behaviour:** 403 Forbidden.

**Implementation note:** introduce `var ErrSubscriberLimitReached = errors.New("subscriber
limit reached")` in `internal/sfu/room.go`. Use `errors.Is` in the handler to map to 403.

### 4. WHEP Without Active Publisher Returns 404 Not Found

**Current behaviour:** `manager.Subscribe` calls an internal `ensureRoom` helper that
auto-creates a room if one does not yet exist, and then `room.Subscribe` proceeds even
when `r.publisher == nil`. The subscriber therefore connects to an idle peer regardless
of whether the room or its publisher ever existed.

**Target behaviour:** if no publisher is present in the room at subscribe time, return
`404 Not Found`. This covers two distinct cases:

1. **Room does not exist** — `manager.Subscribe` SHALL return `ErrNoPublisher` immediately,
   without creating a new room. This requires changing `manager.Subscribe` so that it looks
   up the room with the existing read-path (the same map lookup used by `GetRoom`) rather
   than calling `ensureRoom`. If the room is absent from the map, it returns `ErrNoPublisher`
   directly. `ensureRoom` continues to be called only on the publish path.

2. **Room exists but publisher is nil** — `room.Subscribe` SHALL return `ErrNoPublisher`
   when `r.publisher == nil` at the moment the subscriber attempts to join. This check
   happens inside the existing `r.mu.Lock()` scope so no additional locking is needed.

**Implementation note:** introduce `var ErrNoPublisher = errors.New("no active publisher
in room")` in `internal/sfu/room.go`. Both call sites above return this same sentinel.
Use `errors.Is(err, sfu.ErrNoPublisher)` in the handler to map to `404 Not Found`.

### 5. RoomInfo Uses Lower-CamelCase JSON Tags

**Current behaviour:** `RoomInfo` struct has no JSON tags; `encoding/json` serialises
field names as-is (PascalCase: `Name`, `HasPublisher`, `Tracks`, `Subscribers`).

**Target behaviour:** JSON output uses lower-camelCase (`name`, `hasPublisher`, `tracks`,
`subscribers`).

**Implementation note:** add struct tags to `RoomInfo` in `internal/sfu/manager.go`:
```go
type RoomInfo struct {
    Name         string `json:"name"`
    HasPublisher bool   `json:"hasPublisher"`
    Tracks       int    `json:"tracks"`
    Subscribers  int    `json:"subscribers"`
}
```

### 6. Recording List Shape: name/size/modTime/url

The implementation in `ServeRecordsList` already emits `{name, size, modTime, url}` using
lower-camelCase tags on a local `rec` struct. The existing `api/spec.md` schema
(`room`, `trackID`, `filename`, `size`, `createdAt`) does not match. The delta spec
updates the authoritative shape to match the implementation. No code change is required.

### 7. OpenAPI Source of Truth

`openspec/specs/api/openapi.yaml` is declared as the canonical machine-readable contract.
The spec delta for `api` records this obligation. Creating the actual YAML file is a
subsequent task; this change only establishes the normative reference.

### 8. No New Auth or Query Mechanisms

This change does not add API key query-string auth, OAuth, or any authentication
mechanism not already present. It does not add filtering or pagination to list endpoints.

## Concurrency Considerations

The sentinel error approach introduces no new shared state. `ErrPublisherExists`,
`ErrSubscriberLimitReached`, and `ErrNoPublisher` are package-level value sentinels
(not variables guarded by a mutex). The check `errors.Is(err, sfu.ErrXxx)` in the HTTP
handler is read-only and goroutine-safe.

The `Subscribe` check for `r.publisher == nil` happens inside the existing `r.mu.Lock()`
scope in `room.Subscribe`, so no additional locking is needed.

## Rollback Strategy

Because this change modifies HTTP status codes, rollback must be atomic across two
artifacts:

1. **Handler mapping** — revert the `errors.Is` branches in `ServeWHIPPublish` and
   `ServeWHEPPlay` to a single `http.Error(w, err.Error(), http.StatusBadRequest)`.
2. **Spec and docs** — revert the delta specs and any documentation changes together so
   the spec and implementation remain consistent.

The change is non-breaking for clients that only check for non-2xx status codes on
failure. Clients that expect a specific 4xx code (e.g. testing for 400 on duplicate
publish) will need to update to 409. This is considered a corrective fix, not a
breaking change, because the prior 400 was always incorrect.
