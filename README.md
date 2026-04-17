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

### Run from Source (30 seconds)

```bash
git clone https://github.com/LessUp/go-live.git
cd go-live
go run ./cmd/server
```

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

# Authentication (optional)
AUTH_TOKEN=your-secret-token
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
| `STUN_URLS` | `stun:stun.l.google.com:19302` | STUN servers |
| `TURN_URLS` | - | TURN servers (for NAT environments) |
| `TURN_USERNAME` | - | TURN username |
| `TURN_PASSWORD` | - | TURN password |

### Recording

| Variable | Default | Description |
|----------|---------|-------------|
| `RECORD_ENABLED` | `0` | Enable recording (`1` to enable) |
| `RECORD_DIR` | `records` | Recording output directory |

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
