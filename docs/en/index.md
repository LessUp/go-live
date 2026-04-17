---
layout: default
title: Home
description: Lightweight WebRTC SFU Server - Built with Go and Pion WebRTC
nav_order: 1
lang: en
---

# Go-Live

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

# Or use the development script (recommended)
./scripts/start.sh
```

After the server starts, visit `http://localhost:8080` to see the frontend interface.

---

## ✨ Features

<div class="features-grid">
  <div class="feature-card">
    <div class="feature-icon">📡</div>
    <h3>WHIP/WHEP Protocols</h3>
    <p>Full support for WHIP publishing and WHEP playback, compatible with OBS and browsers</p>
  </div>
  <div class="feature-card">
    <div class="feature-icon">🏠</div>
    <h3>Room-based Broadcast</h3>
    <p>SFU architecture with one publisher and multiple subscribers per room</p>
  </div>
  <div class="feature-card">
    <div class="feature-icon">🎥</div>
    <h3>Recording & Upload</h3>
    <p>Built-in recording with automatic S3/MinIO object storage upload</p>
  </div>
  <div class="feature-card">
    <div class="feature-icon">📊</div>
    <h3>Observability</h3>
    <p>Prometheus metrics and OpenTelemetry distributed tracing integration</p>
  </div>
  <div class="feature-card">
    <div class="feature-icon">🔒</div>
    <h3>Authentication</h3>
    <p>Token and JWT authentication with configurable per-room access control</p>
  </div>
  <div class="feature-card">
    <div class="feature-icon">⚡</div>
    <h3>High Performance</h3>
    <p>Low-latency, high-throughput media stream distribution in Go</p>
  </div>
</div>

---

## 📖 Documentation Navigation

| Documentation | Description | Link |
|---------------|-------------|------|
| Usage Guide | Deployment, configuration, API examples, troubleshooting | [Usage Guide]({{ site.baseurl }}/en/usage.html) |
| API Reference | Complete REST API documentation | [API Reference]({{ site.baseurl }}/en/api.html) |
| Design | System architecture, modules, data flow | [Design]({{ site.baseurl }}/en/design.html) |

---

## 🛠️ Quick Deployment

### Docker Deployment

```bash
# Build image
docker build -t live-webrtc-go:latest .

# Run
docker run --rm -p 8080:8080 live-webrtc-go:latest

# Enable recording
docker run --rm -p 8080:8080 \
  -e RECORD_ENABLED=1 \
  -e RECORD_DIR=/records \
  -v $(pwd)/records:/records \
  live-webrtc-go:latest
```

### Docker Compose

```yaml
version: '3.8'
services:
  live-webrtc:
    build: .
    ports:
      - "8080:8080"
    environment:
      - RECORD_ENABLED=1
    volumes:
      - ./records:/records
    restart: unless-stopped
```

---

## 📦 Requirements

| Dependency | Version | Notes |
|------------|---------|-------|
| Go | 1.22+ | Required for compilation |
| Docker | 20.10+ | Optional for container deployment |
| Browser | Chrome 90+ / Firefox 88+ | WebRTC support required |

---

## 🔗 Related Links

- **GitHub**: [https://github.com/LessUp/go-live](https://github.com/LessUp/go-live)
- **Releases**: [https://github.com/LessUp/go-live/releases](https://github.com/LessUp/go-live/releases)
- **Issues**: [https://github.com/LessUp/go-live/issues](https://github.com/LessUp/go-live/issues)

<style>
.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
  margin: 2rem 0;
}

.feature-card {
  padding: 1.5rem;
  border: 1px solid #e1e4e8;
  border-radius: 8px;
  background: #f6f8fa;
  transition: all 0.2s ease;
}

.feature-card:hover {
  border-color: #0969da;
  box-shadow: 0 2px 8px rgba(9, 105, 218, 0.1);
  transform: translateY(-2px);
}

.feature-icon {
  font-size: 2rem;
  margin-bottom: 0.5rem;
}

.feature-card h3 {
  margin: 0 0 0.5rem 0;
  font-size: 1.1rem;
  color: #24292f;
}

.feature-card p {
  margin: 0;
  color: #57606a;
  font-size: 0.9rem;
  line-height: 1.5;
}
</style>
