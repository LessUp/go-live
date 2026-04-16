# live-webrtc-go

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
- [Quick Start](#-quick-start)
- [Installation](#-installation)
- [Configuration](#-configuration)
- [API Reference](#-api-reference)
- [Documentation](#-documentation)
- [Development](#-development)
- [Docker Deployment](#-docker-deployment)
- [Contributing](#-contributing)
- [License](#-license)

---

## ✨ Features

| Feature | Description |
|---------|-------------|
| 🎥 **WHIP/WHEP Protocol** | Standard HTTP-based WebRTC ingest and playback, compatible with OBS and modern browsers |
| 🏠 **Room-based SFU** | Single publisher, multiple subscribers per room with efficient media forwarding |
| 🔐 **Flexible Authentication** | Global token, per-room tokens, or JWT with role-based access control |
| 📹 **Recording & Upload** | VP8/VP9 → IVF, Opus → OGG with automatic S3/MinIO upload |
| 📊 **Full Observability** | Prometheus metrics, OpenTelemetry tracing, health check endpoints |
| 🐳 **Cloud Native** | Docker and Docker Compose support, Kubernetes-ready |
| 🌐 **Embedded Web UI** | Built-in publisher and player pages, ready to use out of the box |

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

---

## 🚀 Quick Start

### Prerequisites

- Go 1.22+ (for building from source)
- Or Docker (for containerized deployment)

### Run from Source

```bash
# Clone the repository
git clone https://github.com/LessUp/go-live.git
cd go-live

# Run directly
go run ./cmd/server

# Or use the development script
./scripts/start.sh
```

### Run with Docker

```bash
docker run --rm -p 8080:8080 ghcr.io/lessup/go-live:latest
```

### Access the Application

| Page | URL | Description |
|------|-----|-------------|
| 🏠 Home | http://localhost:8080/ | Redirects to publisher |
| 📤 Publisher | http://localhost:8080/web/publisher.html | Browser-based streaming |
| 📥 Player | http://localhost:8080/web/player.html | Watch live streams |
| 📋 Records | http://localhost:8080/web/records.html | Recording browser |
| 📊 Metrics | http://localhost:8080/metrics | Prometheus metrics |

---

## 📦 Installation

### Binary Installation

Download pre-built binaries from [GitHub Releases](https://github.com/LessUp/go-live/releases):

```bash
# Linux AMD64
curl -LO https://github.com/LessUp/go-live/releases/latest/download/live-webrtc-go-linux-amd64
chmod +x live-webrtc-go-linux-amd64
./live-webrtc-go-linux-amd64
```

### Build from Source

```bash
# Clone and build
git clone https://github.com/LessUp/go-live.git
cd go-live
make build

# Binary will be at bin/server
./bin/server
```

---

## ⚙️ Configuration

Configuration is via environment variables:

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
| `TURN_URLS` | - | TURN servers |
| `TURN_USERNAME` | - | TURN username |
| `TURN_PASSWORD` | - | TURN password |

### Recording

| Variable | Default | Description |
|----------|---------|-------------|
| `RECORD_ENABLED` | `0` | Enable recording (`1` to enable) |
| `RECORD_DIR` | `records` | Recording output directory |

See [full configuration guide](https://lessup.github.io/go-live/en/usage.html#configuration-reference) for all options.

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
| `GET` | `/healthz` | Health check |
| `GET` | `/metrics` | Prometheus metrics |

See [complete API documentation](https://lessup.github.io/go-live/en/api.html) for details.

---

## 📚 Documentation

| Document | Description |
|----------|-------------|
| [English Docs](https://lessup.github.io/go-live/en/) | Complete documentation in English |
| [中文文档](https://lessup.github.io/go-live/zh/) | 完整中文文档 |
| [Usage Guide](https://lessup.github.io/go-live/en/usage.html) | Local development, Docker deployment, troubleshooting |
| [Design Docs](https://lessup.github.io/go-live/en/design.html) | System architecture and module details |
| [API Reference](https://lessup.github.io/go-live/en/api.html) | Complete HTTP API documentation |

---

## 🛠️ Development

### Makefile Commands

```bash
make build      # Build binary to bin/
make test       # Run all tests
make lint       # Run linters
make security   # Run security scan
make coverage   # Generate coverage report
make ci         # Full CI pipeline
```

### Running Tests

```bash
# Unit tests
make test-unit

# Integration tests
make test-integration

# All tests
make test-all
```

---

## 🐳 Docker Deployment

### Build Image

```bash
docker build -t live-webrtc-go:latest .
```

### Docker Compose

```bash
docker compose up -d
```

### Full Configuration Example

```bash
docker run --rm -p 8080:8080 \
  -e AUTH_TOKEN=mysecret \
  -e RECORD_ENABLED=1 \
  -e RECORD_DIR=/records \
  -v $(pwd)/records:/records \
  live-webrtc-go:latest
```

See [Docker deployment guide](https://lessup.github.io/go-live/en/usage.html#docker-deployment) for more details.

---

## 🤝 Contributing

Contributions are welcome! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

- [Contributing Guidelines](CONTRIBUTING.md)
- [Code of Conduct](CODE_OF_CONDUCT.md)
- [Security Policy](SECURITY.md)

---

## 📄 License

This project is licensed under the [MIT License](LICENSE).

---

## 🔗 Links

- [GitHub Repository](https://github.com/LessUp/go-live)
- [Issue Tracker](https://github.com/LessUp/go-live/issues)
- [GitHub Releases](https://github.com/LessUp/go-live/releases)
- [Pion WebRTC](https://github.com/pion/webrtc)

---

<div align="center">

**[⬆ Back to Top](#live-webrtc-go)**

Made with ❤️ by [LessUp](https://github.com/LessUp)

</div>
