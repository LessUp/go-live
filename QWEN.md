# Go-Live (live-webrtc-go)

## Project Overview

**Go-Live** is a lightweight, high-performance **WebRTC SFU** (Selective Forwarding Unit) server built with Go and [Pion WebRTC](https://github.com/pion/webrtc). It supports WHIP/WHEP protocols for streaming, room-based broadcast, recording, and comprehensive observability.

### Architecture

The server follows a simple, layered architecture:

```
HTTP Request → CORS → Rate Limiter → Auth → Handler → SFU Room
```

- **SFU Core**: Room-based publisher/subscriber model with efficient RTP fanout
- **HTTP Layer**: WHIP/WHEP endpoints, auth middleware, rate limiting, CORS
- **Recording**: VP8/VP9 → IVF, Opus → OGG with optional S3/MinIO upload
- **Observability**: Prometheus metrics, OpenTelemetry tracing, health checks

### Tech Stack

| Category | Technology |
|----------|------------|
| Language | Go 1.22+ |
| WebRTC | Pion WebRTC v3 |
| Auth | JWT (golang-jwt v5), Token-based |
| Metrics | Prometheus client_golang |
| Tracing | OpenTelemetry (OTLP gRPC/HTTP) |
| Storage | MinIO S3 SDK v7 |
| Frontend | Embedded static HTML/JS |

### Project Structure

```
├── cmd/server/           # Application entry point
│   ├── main.go           # HTTP server initialization, embeds web/
│   └── web/              # Embedded static files (publisher/player UI)
├── internal/
│   ├── api/              # HTTP handlers and routing
│   │   ├── handlers.go   # WHIP/WHEP/Rooms/Admin handlers
│   │   ├── middleware.go # Auth, CORS, rate limiting
│   │   └── routes.go     # URL routing with room name validation
│   ├── config/           # Configuration from environment variables
│   ├── sfu/              # WebRTC SFU core
│   │   ├── manager.go    # Room lifecycle management
│   │   ├── room.go       # PeerConnection, tracks, recording
│   │   └── track.go      # RTP distribution (fanout)
│   ├── metrics/          # Prometheus metrics (rooms, subscribers, RTP)
│   ├── otel/             # OpenTelemetry tracing middleware
│   ├── uploader/         # S3/MinIO upload for recording files
│   └── testutil/         # Test utilities
├── specs/                # Single Source of Truth (Spec-Driven Development)
│   ├── product/          # Product requirements
│   ├── rfc/              # Technical designs
│   ├── api/              # OpenAPI specs
│   ├── db/               # Database schemas
│   └── testing/          # BDD test specs
├── test/                 # Test implementations
│   ├── integration/      # Integration tests
│   ├── e2e/              # End-to-end tests
│   ├── security/         # Security tests
│   ├── performance/      # Benchmarks
│   ├── load/             # Load testing tools
│   └── reports/          # Test reports (generated)
├── docs/                 # Documentation (Jekyll-based GitHub Pages)
│   ├── en/               # English docs
│   ├── zh/               # Chinese docs
│   └── changelog/        # Changelog templates & release notes
├── scripts/              # Development scripts
└── .github/              # GitHub workflows (CI, Pages, dependency review)
```

## Building and Running

### Quick Start

```bash
# Run directly (loads .env.local if exists)
./scripts/start.sh

# Or use go run
go run ./cmd/server
```

### Build from Source

```bash
make build          # Binary at bin/server
./bin/server
```

### Docker

```bash
docker build -t live-webrtc-go:latest .
docker run --rm -p 8080:8080 live-webrtc-go:latest
```

### Configuration

All configuration is via environment variables. Key variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `HTTP_ADDR` | `:8080` | HTTP listen address |
| `AUTH_TOKEN` | - | Global auth token |
| `ADMIN_TOKEN` | - | Admin API token |
| `RECORD_ENABLED` | `0` | Enable recording (`1` to enable) |
| `RECORD_DIR` | `records` | Recording output directory |
| `STUN_URLS` | `stun:stun.l.google.com:19302` | STUN servers |
| `PPROF` | `0` | Enable pprof endpoints |

See `.env.local` example for full configuration.

### API Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `POST` | `/api/whip/publish/{room}` | Token/JWT | Publish stream |
| `POST` | `/api/whep/play/{room}` | Token/JWT | Play stream |
| `GET` | `/api/rooms` | - | List rooms |
| `GET` | `/api/bootstrap` | - | Frontend config |
| `GET` | `/api/records` | Admin | List recordings |
| `POST` | `/api/admin/rooms/{room}/close` | Admin | Close room |
| `GET` | `/healthz` | - | Health check |
| `GET` | `/metrics` | - | Prometheus metrics |

## Development

### Makefile Commands

```bash
make build       # Build binary to bin/
make test        # All tests (unit + integration + security)
make test-unit   # Unit tests only
make test-all    # All tests (incl. e2e, performance)
make lint        # gofmt + go vet + golangci-lint
make security    # gosec security scan
make coverage    # Generate coverage report
make ci          # Full CI pipeline
```

### Testing Conventions

- Unit tests: standard `*_test.go` files alongside source
- Integration tests: `test/integration/` with `-tags=integration`
- E2E tests: `test/e2e/` with `-tags=e2e`
- Security tests: `test/security/` with `-tags=security`
- Performance tests: `test/performance/`

### Code Style

- Go standard formatting (`gofmt -s`)
- Error wrapping with context: `fmt.Errorf("operation failed: %w", err)`
- Structured logging via `log/slog`
- `#nosec` comments for gosec false positives (with explanation)

### Spec-Driven Development

This project follows Spec-Driven Development (SDD). All code changes should be based on specs in `/specs/`:
- Product specs: `/specs/product/`
- Technical RFCs: `/specs/rfc/`
- API specs: `/specs/api/`
- DB specs: `/specs/db/`

See `AGENTS.md` for the full SDD workflow.

## Supported Codecs

| Codec | Type | Format |
|-------|------|--------|
| VP8 | Video | IVF |
| VP9 | Video | IVF |
| Opus | Audio | OGG |

## Links

- **Repository**: https://github.com/LessUp/go-live
- **Documentation**: https://lessup.github.io/go-live/
- **Releases**: https://github.com/LessUp/go-live/releases
