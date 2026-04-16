---
layout: default
title: Home
description: Lightweight WebRTC SFU Server - Built with Go and Pion WebRTC
nav_order: 1
lang: en
---

# live-webrtc-go

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/LessUp/go-live/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/LessUp/go-live)](https://goreportcard.com/report/github.com/LessUp/go-live)
![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)
[![Release](https://img.shields.io/github/v/release/LessUp/go-live)](https://github.com/LessUp/go-live/releases)

[中文]({{ site.baseurl }}/zh/) | **English**

A lightweight, high-performance **WebRTC SFU** (Selective Forwarding Unit) server built with Go and [Pion WebRTC](https://github.com/pion/webrtc). Supports WHIP/WHEP protocols for streaming, room-based broadcast, recording, and comprehensive observability.

---

## 🚀 Quick Start

```bash
# Clone the repository
git clone https://github.com/LessUp/go-live.git
cd go-live

# Run directly
go run ./cmd/server

# Or use the development script
./scripts/start.sh
```

Access the server at `http://localhost:8080`

---

## ✨ Key Features

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

## 📖 Documentation

| Document | Description |
|----------|-------------|
| [Usage Guide](usage.md) | Local development, Docker deployment, configuration reference, troubleshooting |
| [Design Documentation](design.md) | System architecture, module breakdown, data flow diagrams |
| [API Reference](api.md) | Complete HTTP API documentation, request/response formats, error codes |

---

## 🛠️ Development

```bash
# Build
make build

# Run tests
make test

# Full CI pipeline
make ci
```

---

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guidelines](https://github.com/LessUp/go-live/blob/master/CONTRIBUTING.md) for details.

---

## 📄 License

This project is licensed under the [MIT License](https://github.com/LessUp/go-live/blob/master/LICENSE).

---

## 🔗 Links

- [GitHub Repository](https://github.com/LessUp/go-live)
- [Issue Tracker](https://github.com/LessUp/go-live/issues)
- [Changelog](https://github.com/LessUp/go-live/blob/master/CHANGELOG.md)
- [Pion WebRTC](https://github.com/pion/webrtc)
