# Go-Live

[![CI](https://github.com/LessUp/go-live/actions/workflows/ci.yml/badge.svg)](https://github.com/LessUp/go-live/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/LessUp/go-live)](https://goreportcard.com/report/github.com/LessUp/go-live)
![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)
[![Release](https://img.shields.io/github/v/release/LessUp/go-live)](https://github.com/LessUp/go-live/releases)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue?logo=docker)](Dockerfile)
[![Docs](https://img.shields.io/badge/Docs-GitHub%20Pages-blue?logo=github)](https://lessup.github.io/go-live/)

[English](README.md) | [简体中文](README.zh-CN.md)

A lightweight, high-performance **WebRTC SFU** (Selective Forwarding Unit) server built with Go and [Pion WebRTC](https://github.com/pion/webrtc). Supports WHIP/WHEP protocols for streaming, room-based broadcast, recording, and comprehensive observability.

---

## 📋 Table of Contents

- [Features](#-features)
- [Architecture](#️-architecture)
- [Quick Start](#-quick-start)
- [Installation](#-installation)
- [Configuration](#️-configuration)
- [API Reference](#-api-reference)
- [Documentation](#-documentation)
- [Development](#️-development)
- [Docker Deployment](#-docker-deployment)
- [Alternatives](#-alternatives)
- [Troubleshooting](#-troubleshooting)
- [Performance](#-performance)
- [Contributing](#-contributing)
- [License](#-license)

---

## ✨ Features

| Feature | Description |
|---------|-------------|
| 🎥 **WHIP/WHEP Protocol** | Standard HTTP-based WebRTC ingest and playback, compatible with OBS and modern browsers |
| 🏠 **Room-based SFU** | Single publisher, multiple subscribers per room with efficient RTP forwarding |
| 🔐 **Flexible Authentication** | Global token, per-room tokens, or JWT with role-based access control |
| 📹 **Recording & Upload** | VP8/VP9 → IVF, Opus → OGG with automatic S3/MinIO upload |
| 📊 **Full Observability** | Prometheus metrics, OpenTelemetry tracing, health check endpoints |
| 🐳 **Cloud Native** | Docker and Docker Compose support, Kubernetes-ready |
| 🌐 **Embedded Web UI** | Built-in publisher and player pages, ready to use out of the box |
| ⚡ **High Performance** | Low-latency media forwarding with Go's concurrent programming model |

---

## 🏗️ Architecture

```
                         ┌─────────────────────────────────────┐
                         │          HTTP Server :8080          │
                         │                                     │
    ┌──────────┐         │  ┌─────────┐    ┌─────────────┐    │
    │ Publisher│ ──WHIP──▶│  │  Auth   │───▶│    SFU      │    │
    │(OBS/Web) │         │  │Middleware│    │  Manager    │    │
    └──────────┘         │  └─────────┘    └──────┬──────┘    │
                         │                        │           │
    ┌──────────┐         │        ┌───────────────┤           │
    │  Viewer  │ ◀──WHEP──│        │               │           │
    │(Browser) │───────▶  │        ▼               ▼           │
    └──────────┘         │ ┌─────────────┐ ┌───────────┐      │
                         │ │    Room     │ │ Recording │      │
                         │ │  (Fanout)   │ │  & Upload │      │
                         │ └──────┬──────┘ └─────┬─────┘      │
                         └────────┼──────────────┼────────────┘
                                  │              │
                        ┌─────────▼──────────────▼───────────┐
                        │          Object Storage            │
                        │           (S3/MinIO)               │
                        └────────────────────────────────────┘
```

### Request Processing Chain

```
HTTP Request → CORS → Rate Limiter → Auth → Handler → SFU Room
```

---

## 🚀 Quick Start

### Prerequisites

- **Go 1.22+** - [Download Go](https://go.dev/dl/)
- **Git** - For cloning the repository
- **Port 8080** - Available on your machine

### Run from Source

```bash
# Clone and run
git clone https://github.com/LessUp/go-live.git
cd go-live
go run ./cmd/server
```

### Verify Installation

```bash
curl http://localhost:8080/healthz
# Expected response: ok
```

### Stream with OBS

1. Open OBS Studio → **Settings** → **Stream**
2. **Service**: Select `WHIP`
3. **Server**: `http://localhost:8080/api/whip/publish/{room}` (e.g., `myroom`)
4. **Bearer Token**: Your `AUTH_TOKEN` (if configured)
5. Click **Start Streaming**

> 💡 See [OBS WHIP Guide](https://lessup.github.io/go-live/en/usage.html#obs-setup) for detailed instructions.

### Run with Docker

```bash
docker run --rm -p 8080:8080 ghcr.io/lessup/go-live:latest
```

### Access the Application

| Page | URL | Description |
|------|-----|-------------|
| 🏠 Home | `http://localhost:8080/` | Frontend dashboard |
| 📤 Publisher | `http://localhost:8080/web/publisher.html` | Browser-based streaming |
| 📥 Player | `http://localhost:8080/web/player.html` | Watch live streams |
| 📋 Records | `http://localhost:8080/web/records.html` | Recording browser |
| 📊 Metrics | `http://localhost:8080/metrics` | Prometheus metrics |
| ❤️ Health | `http://localhost:8080/healthz` | Health check |

### Quick Test Flow

1. Open **Publisher** page → Start streaming to a room
2. Open **Player** page (new tab) → Enter same room name → Watch stream
3. Check **Metrics** → See active connections

---

## 📦 Installation

### Binary Download

Download pre-built binaries from [GitHub Releases](https://github.com/LessUp/go-live/releases):

```bash
# Linux AMD64
curl -LO https://github.com/LessUp/go-live/releases/latest/download/live-webrtc-go-linux-amd64
chmod +x live-webrtc-go-linux-amd64
./live-webrtc-go-linux-amd64
```

### Build from Source

```bash
git clone https://github.com/LessUp/go-live.git
cd go-live
make build
./bin/server
```

---

## ⚙️ Configuration

Configuration is via environment variables. Create an `.env.local` file:

```bash
# Core
HTTP_ADDR=:8080
ALLOWED_ORIGIN=*

# Authentication (optional but recommended for production)
AUTH_TOKEN=change-me-in-production    # ⚠️ Use a secure random token!
ADMIN_TOKEN=your-admin-token

# Recording (optional)
RECORD_ENABLED=1
RECORD_DIR=records

# S3 Upload (optional)
UPLOAD_RECORDINGS=1
S3_ENDPOINT=minio.example.com:9000
S3_ACCESS_KEY=minioadmin
S3_SECRET_KEY=minioadmin
S3_BUCKET=recordings
```

> ⚠️ **Security Warning**: Never use the example tokens in production. Generate secure random tokens with `openssl rand -hex 32`.

### Core Settings

| Variable | Default | Description |
|----------|---------|-------------|
| `HTTP_ADDR` | `:8080` | HTTP listen address |
| `ALLOWED_ORIGIN` | `*` | CORS allowed origin |

### Authentication

| Variable | Description |
|----------|-------------|
| `AUTH_TOKEN` | Global authentication token |
| `ROOM_TOKENS` | Per-room tokens: `room1:tok1;room2:tok2` |
| `JWT_SECRET` | JWT HMAC signing key |
| `ADMIN_TOKEN` | Token for admin endpoints |

### WebRTC/ICE

| Variable | Default | Description |
|----------|---------|-------------|
| `STUN_URLS` | `stun:stun.l.google.com:19302` | STUN servers (comma-separated) |
| `TURN_URLS` | - | TURN servers (comma-separated) |
| `TURN_USERNAME` | - | TURN username |
| `TURN_PASSWORD` | - | TURN password |

### Rate Limiting & Limits

| Variable | Default | Description |
|----------|---------|-------------|
| `RATE_LIMIT_RPS` | `0` | Requests per second per IP (0 = disabled) |
| `MAX_SUBS_PER_ROOM` | `0` | Max subscribers per room (0 = unlimited) |
| `SDP_MAX_SIZE` | `32768` | Max SDP offer size in bytes |

### Recording

| Variable | Default | Description |
|----------|---------|-------------|
| `RECORD_ENABLED` | `0` | Enable recording (`1` to enable) |
| `RECORD_DIR` | `records` | Recording output directory |
| `UPLOAD_RECORDINGS` | `0` | Enable S3 upload (`1` to enable) |

### S3/MinIO Upload

| Variable | Default | Description |
|----------|---------|-------------|
| `S3_ENDPOINT` | - | S3/MinIO endpoint |
| `S3_ACCESS_KEY` | - | Access key |
| `S3_SECRET_KEY` | - | Secret key |
| `S3_BUCKET` | `recordings` | Bucket name |
| `S3_REGION` | `us-east-1` | Region |
| `S3_USE_SSL` | `0` | Use HTTPS (`1` to enable) |

### Debug

| Variable | Default | Description |
|----------|---------|-------------|
| `PPROF` | `0` | Enable pprof endpoints (`1` to enable) |
| `LOG_LEVEL` | `info` | Log level: `debug`, `info`, `warn`, `error` |

> 💡 See [full configuration guide](https://lessup.github.io/go-live/en/usage.html#configuration-reference) for all options including S3 upload, rate limiting, TLS, and debug settings.

---

## 🔌 API Reference

### Streaming Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `POST` | `/api/whip/publish/{room}` | Token/JWT | Publish stream to room |
| `POST` | `/api/whep/play/{room}` | Token/JWT | Play stream from room |

### Query Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/bootstrap` | Runtime config for frontend |
| `GET` | `/api/rooms` | List active rooms |
| `GET` | `/api/records` | List recording files |

### Admin Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `POST` | `/api/admin/rooms/{room}/close` | Admin Token | Force close a room |

### Health & Metrics

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/healthz` | Health check (returns `ok`) |
| `GET` | `/metrics` | Prometheus metrics |

### Quick Test with curl

```bash
# Health check
curl http://localhost:8080/healthz

# List rooms
curl http://localhost:8080/api/rooms

# Bootstrap config
curl http://localhost:8080/api/bootstrap | jq

# Close room (admin)
curl -X POST \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  http://localhost:8080/api/admin/rooms/myroom/close
```

> 💡 See [complete API documentation](https://lessup.github.io/go-live/en/api.html) for request/response formats and examples.

---

## 📚 Documentation

| Document | Link |
|----------|------|
| English Docs | https://lessup.github.io/go-live/en/ |
| 中文文档 | https://lessup.github.io/go-live/zh/ |
| Usage Guide | https://lessup.github.io/go-live/en/usage.html |
| Design Docs | https://lessup.github.io/go-live/en/design.html |
| API Reference | https://lessup.github.io/go-live/en/api.html |
| Changelog | https://lessup.github.io/go-live/changelog.html |

---

## 🛠️ Development

### Makefile Commands

```bash
make build       # Build binary to bin/
make test        # Run all tests (unit + integration + security)
make test-unit   # Run unit tests only
make test-all    # Run all tests including e2e and performance
make lint        # Run linters (gofmt + go vet + golangci-lint)
make security    # Run gosec security scan
make coverage    # Generate coverage report
make ci          # Full CI pipeline (lint + test + security)
```

### Project Structure

```
├── cmd/server/           # Application entry point
│   ├── main.go           # HTTP server initialization
│   └── web/              # Embedded static files
├── internal/
│   ├── api/              # HTTP handlers and routing
│   ├── config/           # Configuration management
│   ├── sfu/              # WebRTC SFU core
│   ├── metrics/          # Prometheus metrics
│   ├── otel/             # OpenTelemetry tracing
│   ├── uploader/         # S3/MinIO upload
│   └── testutil/         # Test utilities
├── specs/                # Single Source of Truth (Specs)
├── test/                 # Test implementations
└── docs/                 # Documentation
```

### Supported Codecs

| Codec | Type | Format |
|-------|------|--------|
| VP8 | Video | IVF |
| VP9 | Video | IVF |
| Opus | Audio | OGG |

---

## 🐳 Docker Deployment

### Production Checklist

- [ ] Set secure `AUTH_TOKEN` and `ADMIN_TOKEN`
- [ ] Configure HTTPS/TLS (required for WebRTC in production)
- [ ] Set up TURN server for NAT traversal
- [ ] Configure `ALLOWED_ORIGIN` to your domain
- [ ] Set up monitoring (Prometheus + Grafana)

### Build Image

```bash
docker build -t live-webrtc-go:latest .
```

### Basic Run

```bash
docker run --rm -p 8080:8080 live-webrtc-go:latest
```

### Docker Compose

```yaml
services:
  live-webrtc:
    build: .
    ports:
      - "8080:8080"
    environment:
      - AUTH_TOKEN=${AUTH_TOKEN}
      - RECORD_ENABLED=1
      - RECORD_DIR=/records
    volumes:
      - ./records:/records
    restart: unless-stopped
```

```bash
docker compose up -d
```

### Full Configuration

```bash
docker run --rm -p 8080:8080 \
  -e AUTH_TOKEN=mysecret \
  -e RECORD_ENABLED=1 \
  -e RECORD_DIR=/records \
  -e UPLOAD_RECORDINGS=1 \
  -e S3_ENDPOINT=minio:9000 \
  -e S3_ACCESS_KEY=minioadmin \
  -e S3_SECRET_KEY=minioadmin \
  -e S3_BUCKET=recordings \
  -v $(pwd)/records:/records \
  live-webrtc-go:latest
```

> 💡 See [Docker deployment guide](https://lessup.github.io/go-live/en/usage.html#docker-deployment) for Kubernetes examples and production best practices.

---

## 🔄 Alternatives

How does go-live compare to other WebRTC servers?

| Project | Type | Language | Key Difference |
|---------|------|----------|----------------|
| **go-live** | SFU | Go | Lightweight, WHIP/WHEP native, single binary |
| [LiveKit](https://livekit.io) | SFU | Go | Full-featured, SDK ecosystem, scalable |
| [Mediasoup](https://mediasoup.org) | SFU | Node.js/C++ | High performance, Node.js integration |
| [Janus](https://janus.conf.meetecho.com) | MCU/SFU | C | Multi-protocol, plugin architecture |
| [Pionion](https://github.com/pion/ion) | SFU | Go | Pion-based, microservices |

**Choose go-live if you want:**
- Single binary deployment
- WHIP/WHEP protocol (OBS/browser native)
- Minimal dependencies
- Easy containerization

---

## 🤝 Contributing

Contributions are welcome! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

- [Contributing Guidelines](CONTRIBUTING.md)
- [Code of Conduct](CODE_OF_CONDUCT.md)
- [Security Policy](SECURITY.md)
- [Spec-Driven Development](AGENTS.md)

---

## 📄 License

This project is licensed under the [MIT License](LICENSE).

---

## 🔧 Troubleshooting

### Common Issues

| Issue | Cause | Solution |
|-------|-------|----------|
| **ICE connection failed** | STUN/TURN unreachable | Check firewall allows UDP; verify STUN_URLS |
| **No video/audio** | Publisher not connected | Verify WHIP endpoint returns 201; check browser console |
| **Auth rejected (401)** | Missing/invalid token | Ensure `Authorization: Bearer <token>` header is set |
| **Room not found** | Publisher disconnected | Rooms auto-close when publisher leaves |
| **High latency** | Network/ICE issues | Use TURN server for NAT traversal; check network |

### Debug Commands

```bash
# Check health
curl http://localhost:8080/healthz

# View metrics
curl http://localhost:8080/metrics

# List active rooms
curl http://localhost:8080/api/rooms

# Enable pprof (PPROF=1)
curl http://localhost:8080/debug/pprof/
```

### Network Requirements

- **UDP ports**: WebRTC uses dynamic UDP ports for media
- **STUN**: UDP 3478 (outbound)
- **TURN**: UDP/TCP 3478 (if configured)
- **Firewall**: Allow outbound UDP; for inbound, configure TURN

> 💡 See [Troubleshooting Guide](https://lessup.github.io/go-live/en/troubleshooting.html) for detailed diagnostics.

---

## 📈 Performance

### Benchmarks

| Metric | Value | Notes |
|--------|-------|-------|
| Publisher latency | < 50ms | Local network |
| Subscriber fanout | 1000+ per room | Depends on server resources |
| Memory per subscriber | ~2MB | Video + audio tracks |
| CPU usage | ~5% per 100 subs | VP8 passthrough (no transcoding) |

### Optimization Tips

- Use **TURN relay** for challenging NAT environments
- Increase `MAX_SUBS_PER_ROOM` for large audiences
- Enable **Prometheus** to monitor resource usage
- Deploy multiple instances behind a load balancer for horizontal scaling

---

## 🔗 Links

- [GitHub Repository](https://github.com/LessUp/go-live)
- [Issue Tracker](https://github.com/LessUp/go-live/issues)
- [Releases](https://github.com/LessUp/go-live/releases)
- [Documentation](https://lessup.github.io/go-live/)
- [Pion WebRTC](https://github.com/pion/webrtc)

---

<div align="center">

**[⬆ Back to Top](#go-live)**

Made with ❤️ by [LessUp](https://github.com/LessUp)

</div>
