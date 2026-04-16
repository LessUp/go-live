# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A lightweight **WebRTC SFU (Selective Forwarding Unit)** server built with Go 1.22+ and [Pion WebRTC](https://github.com/pion/webrtc). It provides real-time live streaming capabilities through WHIP/WHEP protocols with room-based relay, optional recording, authentication, rate limiting, and observability features.

### Key Features
- **WHIP Publishing**: Browser or OBS can push streams via HTTP POST
- **WHEP Playback**: Viewers subscribe to rooms via HTTP POST
- **Room-based SFU**: Single publisher per room, multiple subscribers
- **Recording**: VP8/VP9 → IVF, Opus → OGG with optional S3/MinIO upload
- **Auth**: Token (global/room-level) or JWT with role-based access
- **Observability**: Prometheus metrics, OpenTelemetry tracing, health checks

## Architecture

```
┌─────────────┐   WHIP (POST /api/whip/publish/{room})   ┌──────────────┐
│  Publisher  │ ──────────────────────────────────────▶ │ HTTP Server  │
│  (OBS/Web)  │                                         │   :8080      │
└─────────────┘                                         └──────┬───────┘
                                                               │
                                      Creates PeerConnection   │
                                      & TrackFanout            ▼
                                                        ┌─────────────┐
┌─────────────┐   WHEP (POST /api/whep/play/{room})    │    Room     │
│   Viewer    │ ──────────────────────────────────────▶ │   (SFU)     │
│  (Browser)  │ ◀────────────────────────────────────── │             │
└─────────────┘        RTP Packets (WebRTC)             └──────┬──────┘
                                                               │
                                          Records & Uploads    ▼
                                                        ┌─────────────┐
                                                        │ Object Store│
                                                        │  (S3/MinIO) │
                                                        └─────────────┘
```

## Project Structure

```
├── cmd/server/           # Entry point
│   ├── main.go           # HTTP server setup, graceful shutdown
│   └── web/              # Embedded static files (publisher/player HTML)
├── internal/
│   ├── api/              # HTTP layer
│   │   ├── handlers.go   # WHIP/WHEP/Rooms/Records/Admin endpoints
│   │   ├── middleware.go # CORS, rate limiting, auth (token/JWT)
│   │   └── routes.go     # URL routing, room name validation
│   ├── config/           # Environment variable configuration
│   ├── sfu/              # Core WebRTC SFU logic
│   │   ├── manager.go    # Room lifecycle management
│   │   ├── room.go       # PeerConnection, track fanout, recording
│   │   └── track.go      # RTP packet distribution to subscribers
│   ├── metrics/          # Prometheus gauges/counters
│   ├── otel/             # OpenTelemetry tracer initialization
│   ├── uploader/         # S3/MinIO upload client
│   └── testutil/         # Test helpers
├── test/
│   ├── integration/      # Integration tests (requires -tags=integration)
│   ├── e2e/              # End-to-end tests (requires -tags=e2e)
│   ├── security/         # Security tests
│   ├── performance/      # Performance benchmarks
│   └── load/             # Load testing tools
├── docs/                 # GitHub Pages documentation
└── web/                  # Source static files (embedded into binary)
```

## Build & Run

```bash
# Build
go build -o bin/live-webrtc-go ./cmd/server

# Run directly
go run ./cmd/server

# Development helper (loads .env.local, sets cache dirs)
./scripts/start.sh

# With module tidy
RUN_TIDY=1 ./scripts/start.sh
```

## Test Commands

```bash
make test          # Unit + integration + security (default)
make test-all      # Adds e2e + performance (longer timeout)
make test-unit     # go test -v -race ./internal/...
make test-integration  # Requires -tags=integration
make test-e2e      # Requires -tags=e2e, 10m timeout
make coverage      # HTML + XML coverage reports
```

## Lint & Security

```bash
make lint        # gofmt -s + go vet + golangci-lint
make fmt         # gofmt -s -w .
make security    # gosec ./...
```

## Key Environment Variables

| Variable | Default | Purpose |
|----------|---------|---------|
| `HTTP_ADDR` | `:8080` | Listen address |
| `ALLOWED_ORIGIN` | `*` | CORS origin |
| `AUTH_TOKEN` | - | Global auth token |
| `ROOM_TOKENS` | - | Per-room tokens (`room1:tok1;room2:tok2`) |
| `JWT_SECRET` | - | HMAC secret for JWT auth |
| `ADMIN_TOKEN` | - | Admin API token |
| `RECORD_ENABLED` | `0` | Enable recording (`1` to enable) |
| `RECORD_DIR` | `records` | Recording output directory |
| `UPLOAD_RECORDINGS` | `0` | Enable S3 upload |
| `S3_*` | - | S3/MinIO configuration |
| `RATE_LIMIT_RPS` | `0` | Requests per second per IP |
| `MAX_SUBS_PER_ROOM` | `0` | Subscriber limit (0 = unlimited) |
| `STUN_URLS` | Google STUN | Comma-separated STUN servers |
| `TURN_URLS` | - | Comma-separated TURN servers |

## Code Conventions

### Go Style
- Follow standard Go idioms and [Effective Go](https://golang.org/doc/effective_go)
- Use `gofmt -s` for formatting (enforced by `make lint`)
- Run `go vet` and `golangci-lint` before committing

### Error Handling
- Wrap errors with context: `fmt.Errorf("operation failed: %w", err)`
- Log errors with `slog` at appropriate levels
- Never silently ignore errors in production code

### Concurrency
- Protect shared state with `sync.Mutex` or `sync.RWMutex`
- Use `context.Context` for cancellation and timeouts
- Ensure goroutines can exit (avoid leaks)

### Testing
- Write table-driven tests for multiple cases
- Use `t.Parallel()` where tests are independent
- Mock external dependencies (e.g., `uploadRecordingFile` in room.go)

### Naming
- Room names: `^[A-Za-z0-9_-]{1,64}$` (enforced in routes.go)
- Exported functions need documentation comments
- Use descriptive variable names (avoid single letters except loops)

## Repository Conventions

- **Primary branch**: `master`
- **Commit messages**: Present tense, concise ("Add feature" not "Added feature")
- **PRs**: Include summary, motivation, test evidence
- **Versioning**: Semantic Versioning (MAJOR.MINOR.PATCH)

## Working Style

1. **Plan before implementing** non-trivial changes
2. **Explain tradeoffs** when multiple approaches exist
3. **Keep changes focused** - one logical change per PR
4. **Update documentation** when behavior changes
5. **Be terse** - no trailing summaries in code review

## Common Tasks

### Adding a new API endpoint
1. Add handler in `internal/api/handlers.go`
2. Add route in `internal/api/routes.go`
3. Add auth check if needed (use `authOKRoom` or `adminOK`)
4. Add rate limit check with `allowRate`
5. Add CORS headers with `allowCORS`
6. Write tests in `internal/api/handlers_test.go`

### Adding a new configuration option
1. Add field to `Config` struct in `internal/config/config.go`
2. Parse from environment in `Load()` function
3. Add test in `internal/config/config_test.go`
4. Update README.md and docs/usage.md

### Modifying SFU behavior
1. Understand `Manager` → `Room` → `trackFanout` hierarchy
2. Be careful with goroutine lifecycle in `track.go`
3. Update metrics if adding new observables
4. Test with multiple publishers/subscribers

## Debugging

```bash
# Enable pprof endpoints
PPROF=1 go run ./cmd/server

# Access pprof
open http://localhost:8080/debug/pprof/

# View metrics
curl http://localhost:8080/metrics

# Check health
curl http://localhost:8080/healthz
```

## Security Considerations

- Never commit secrets, tokens, or credentials
- Validate all user inputs (room names, SDP size limits)
- Use `crypto/subtle.ConstantTimeCompare` for token comparison
- Run `make security` for security scans
- Review auth changes carefully
