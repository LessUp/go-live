# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

> **Navigation**: [Unreleased](#unreleased) | [1.0.0](#100---2025-03-22) | [Version Summary](#version-history-summary)

---

## [Unreleased]

[Empty]

---

## [1.1.0] - 2025-04-16

### Overview

Enhanced documentation release featuring bilingual (English/Chinese) documentation site, professional changelog management, and improved project structure.

Documentation Release / 文档发布：
- Completely restructured docs directory with English and Chinese versions
- Professionalized changelog management system
- Optimized README.md in both languages
- Added GitHub Pages documentation site with language switching

### Added

#### Documentation / 文档
- **Bilingual Documentation Site** - Complete docs in English and Chinese
  - English docs at `/docs/en/` - Usage guide, design docs, API reference
  - 中文文档位于 `/docs/zh/` - 使用指南、设计说明、API 参考
- **GitHub Pages Integration** - Auto-generated documentation site
  - Language switcher for seamless navigation
  - Jekyll Cayman theme with custom styling
- **Professional Changelog System** - Enhanced changelog management
  - `CHANGELOG_GUIDE.md` - Detailed changelog writing guidelines
  - `RELEASE_WORKFLOW.md` - Complete release process documentation
  - Bilingual release note templates

#### Observability
- OpenTelemetry tracing support for distributed observability
  - Configurable via `OTEL_EXPORTER_OTLP_ENDPOINT` and `OTEL_SERVICE_NAME`
  - Supports both stdout and OTLP (grpc/http) exporters
  - HTTP middleware for automatic span creation

#### Code Quality
- `HTTPHandlers.Close()` method for graceful shutdown of rate limiter goroutine
- Unit tests for `internal/uploader` package (23.5% coverage)
- Unit tests for `internal/otel` package (47.4% coverage)
- Named constant `mtuSize` replacing magic number 1500

### Changed

#### Documentation / 文档
- **Restructured README.md** - More professional structure with quick navigation
- **Enhanced Quick Start** - Clearer installation and setup instructions
- **Improved Configuration Tables** - Better organized environment variable reference

#### Improvements
- JWT parser refactored with consolidated options pattern (reduced code duplication)
- Enhanced server shutdown with proper error logging
- Removed unused functions (`getRoom`, `pruneRoom`) from SFU Manager

### Fixed

#### Resource Management
- **Resource leak**: Rate limiter GC goroutine now properly stopped on shutdown via `HTTPHandlers.Close()`
- Silent error ignore on server shutdown - now logs errors properly

---

## [1.0.0] - 2025-03-22

### Overview

Initial stable release featuring a lightweight WebRTC SFU server with WHIP/WHEP protocol support, flexible authentication, recording capabilities, and comprehensive observability.

### Added

#### Core Features

| Feature | Description |
|---------|-------------|
| **WebRTC SFU** | Room-based Selective Forwarding Unit using Pion WebRTC v3 |
| **WHIP Protocol** | HTTP-based WebRTC ingest for publishers (OBS, browsers) |
| **WHEP Protocol** | HTTP-based WebRTC playback for viewers |
| **Room Management** | Single publisher per room, multiple subscribers |

#### Authentication System

- **Global Token**: Single `AUTH_TOKEN` for all rooms
- **Per-room Tokens**: `ROOM_TOKENS` format `room1:tok1;room2:tok2`
- **JWT Authentication**: 
  - HMAC-based (HS256/HS384/HS512)
  - Role-based access (`role=admin` or `admin=true`)
  - Room restriction via `room` claim
  - Audience validation via `JWT_AUDIENCE`
- **Admin Token**: `ADMIN_TOKEN` for management endpoints

#### Media & Recording

- Video recording: VP8/VP9 → IVF format
- Audio recording: Opus → OGG format (48kHz, stereo)
- Automatic S3/MinIO upload with `UPLOAD_RECORDINGS`
- Configurable local file deletion after upload
- Recording file naming: `{room}_{trackID}_{timestamp}.{ext}`

#### Security Features

- **Rate Limiting**: Per-IP token bucket algorithm
  - Configurable RPS via `RATE_LIMIT_RPS`
  - Burst capacity via `RATE_LIMIT_BURST`
  - Automatic garbage collection of stale entries
- **CORS**: Configurable allowed origins via `ALLOWED_ORIGIN`
- **Token Comparison**: Constant-time comparison using `crypto/subtle`
- **Input Validation**: Room name regex `^[A-Za-z0-9_-]{1,64}$`
- **SDP Size Limit**: 1MB maximum request body

#### Observability

**Prometheus Metrics** (`/metrics`):

| Metric | Type | Description |
|--------|------|-------------|
| `live_rooms` | Gauge | Active room count |
| `live_subscribers` | GaugeVec | Subscribers per room |
| `live_rtp_bytes_total` | CounterVec | Total RTP bytes |
| `live_rtp_packets_total` | CounterVec | Total RTP packets |

**Other Endpoints**:
- `GET /healthz` - Health check (returns `ok`)
- `GET /api/rooms` - Room status with publisher/subscriber counts
- `GET /api/records` - Recording file metadata list

#### Deployment

- **Docker**: Multi-stage Dockerfile, minimal image size
- **Docker Compose**: Ready-to-use `docker-compose.yml`
- **Static Binary**: Web assets embedded via `go:embed`
- **Environment Variables**: All configuration via env vars

#### Development Infrastructure

- **Makefile**: Build, test, lint, security, coverage targets
- **Test Suites**:
  - Unit tests with race detection
  - Integration tests (`-tags=integration`)
  - Security tests
  - E2E tests (`-tags=e2e`)
  - Performance benchmarks
- **CI/CD**: GitHub Actions workflows
- **Static Analysis**: golangci-lint, gosec integration

#### Web Interface

| Page | Path | Description |
|------|------|-------------|
| Home | `/` | Redirects to publisher |
| Publisher | `/web/publisher.html` | Browser-based streaming |
| Player | `/web/player.html` | Browser-based playback |
| Records | `/web/records.html` | Recording browser |

**Frontend Features**:
- Shared static resource layer
- Bootstrap configuration API (`/api/bootstrap`)
- Unified WebRTC lifecycle management
- Safe DOM rendering (no innerHTML)
- Loading, empty, and error states

#### Documentation

- `README.md` - English documentation
- `README.zh-CN.md` - Chinese documentation
- `docs/design.md` - Architecture and module design
- `docs/usage.md` - Deployment and usage guide
- `docs/api.md` - Complete API reference
- `CONTRIBUTING.md` - Contribution guidelines
- `CODE_OF_CONDUCT.md` - Community standards
- `SECURITY.md` - Security policy
- `CLAUDE.md` - AI assistant guidance

### API Reference

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `POST` | `/api/whip/publish/{room}` | Token/JWT | Publish stream (SDP Offer → Answer) |
| `POST` | `/api/whep/play/{room}` | Token/JWT | Play stream (SDP Offer → Answer) |
| `GET` | `/api/bootstrap` | None | Frontend runtime configuration |
| `GET` | `/api/rooms` | None | List active rooms |
| `GET` | `/api/records` | None | List recording files |
| `POST` | `/api/admin/rooms/{room}/close` | Admin | Force close a room |
| `GET` | `/healthz` | None | Health check |
| `GET` | `/metrics` | None | Prometheus metrics |
| `GET` | `/debug/pprof/*` | None | pprof endpoints (when enabled) |

### Configuration Reference

#### Server Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `HTTP_ADDR` | `:8080` | HTTP listen address |
| `ALLOWED_ORIGIN` | `*` | CORS allowed origin |
| `TLS_CERT_FILE` | - | TLS certificate path |
| `TLS_KEY_FILE` | - | TLS private key path |

#### Authentication

| Variable | Default | Description |
|----------|---------|-------------|
| `AUTH_TOKEN` | - | Global authentication token |
| `ROOM_TOKENS` | - | Per-room tokens (`room:token;...`) |
| `JWT_SECRET` | - | JWT HMAC signing key |
| `JWT_AUDIENCE` | - | Required JWT audience claim |
| `ADMIN_TOKEN` | - | Admin API token |

#### WebRTC/ICE

| Variable | Default | Description |
|----------|---------|-------------|
| `STUN_URLS` | `stun:stun.l.google.com:19302` | STUN servers (comma-separated) |
| `TURN_URLS` | - | TURN servers (comma-separated) |
| `TURN_USERNAME` | - | TURN username |
| `TURN_PASSWORD` | - | TURN password |

#### Recording

| Variable | Default | Description |
|----------|---------|-------------|
| `RECORD_ENABLED` | `0` | Enable recording (`1` to enable) |
| `RECORD_DIR` | `records` | Recording output directory |
| `UPLOAD_RECORDINGS` | `0` | Enable S3 upload (`1` to enable) |
| `DELETE_RECORDING_AFTER_UPLOAD` | `0` | Delete local file after upload |

#### S3/MinIO

| Variable | Default | Description |
|----------|---------|-------------|
| `S3_ENDPOINT` | - | S3/MinIO endpoint |
| `S3_REGION` | - | S3 region |
| `S3_BUCKET` | - | Target bucket name |
| `S3_ACCESS_KEY` | - | Access Key ID |
| `S3_SECRET_KEY` | - | Secret Access Key |
| `S3_USE_SSL` | `1` | Use HTTPS for S3 connection |
| `S3_PATH_STYLE` | `0` | Use path-style addressing |
| `S3_PREFIX` | - | Object key prefix |

#### Limits & Rate Limiting

| Variable | Default | Description |
|----------|---------|-------------|
| `MAX_SUBS_PER_ROOM` | `0` | Max subscribers per room (0=unlimited) |
| `RATE_LIMIT_RPS` | `0` | Requests per second per IP (0=disabled) |
| `RATE_LIMIT_BURST` | `0` | Burst capacity for rate limiter |

#### Debug

| Variable | Default | Description |
|----------|---------|-------------|
| `PPROF` | `0` | Enable pprof endpoints (`1` to enable) |
| `OTEL_SERVICE_NAME` | `live-webrtc-go` | OpenTelemetry service name |

### Technical Details

#### Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| `github.com/pion/webrtc/v3` | v3.2.24 | WebRTC implementation |
| `github.com/pion/rtp` | - | RTP packet handling |
| `github.com/golang-jwt/jwt/v5` | v5 | JWT authentication |
| `github.com/minio/minio-go/v7` | v7 | S3/MinIO client |
| `github.com/prometheus/client_golang` | - | Prometheus metrics |
| `go.opentelemetry.io/otel` | - | OpenTelemetry tracing |
| `golang.org/x/time/rate` | - | Rate limiting |

#### Project Structure

```
├── cmd/server/           # Entry point
├── internal/
│   ├── api/              # HTTP handlers & middleware
│   ├── config/           # Configuration management
│   ├── metrics/          # Prometheus metrics
│   ├── otel/             # OpenTelemetry tracing
│   ├── sfu/              # WebRTC SFU core logic
│   ├── uploader/         # S3/MinIO upload
│   └── testutil/         # Test utilities
├── test/
│   ├── integration/      # Integration tests
│   ├── e2e/              # End-to-end tests
│   ├── security/         # Security tests
│   ├── performance/      # Performance tests
│   └── load/             # Load testing
├── docs/                 # Documentation (GitHub Pages)
├── web/                  # Frontend static files
└── scripts/              # Development scripts
```

---

## Version History Summary

| Version | Date | Type | Highlights |
|---------|------|------|------------|
| [1.0.0] | 2025-03-22 | Major | Initial stable release with WHIP/WHEP, auth, recording, metrics |

---

## Contributing to this Changelog

### Guidelines

1. **Location**: Add entries under `[Unreleased]` section
2. **Categories**: Use appropriate subsection (Added, Changed, Deprecated, Removed, Fixed, Security)
3. **Style**: 
   - Keep entries concise but descriptive
   - Use imperative mood ("Add feature" not "Added feature")
   - Reference issues/PRs where applicable
4. **Grouping**: Group related changes together
5. **Release**: Move `[Unreleased]` to new version section on release

### Example Entry

```markdown
### Added
- Add WebSocket support for real-time room events (#123)

### Fixed
- Fix memory leak in track fanout when subscriber disconnects (#124)
```

---

## Versioning Policy

This project follows [Semantic Versioning](https://semver.org/):

| Version Component | Increment When | Example |
|-------------------|----------------|---------|
| **MAJOR** (X.0.0) | Incompatible API changes | Breaking endpoint changes |
| **MINOR** (0.X.0) | Backwards-compatible features | New API endpoints |
| **PATCH** (0.0.X) | Backwards-compatible fixes | Bug fixes |

**Pre-release versions**: Use hyphen suffix (e.g., `1.1.0-beta.1`, `2.0.0-rc.2`)

---

[Unreleased]: https://github.com/LessUp/go-live/compare/v1.1.0...HEAD
[1.1.0]: https://github.com/LessUp/go-live/releases/tag/v1.1.0
[1.0.0]: https://github.com/LessUp/go-live/releases/tag/v1.0.0
