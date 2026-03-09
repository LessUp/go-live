# live-webrtc-go

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go&logoColor=white)

English | [简体中文](README.zh-CN.md)

A lightweight live streaming service built with Go + [Pion WebRTC](https://github.com/pion/webrtc). Implements WHIP publishing, WHEP playback, embedded web pages, configurable auth, and room status queries.

## Features

- **WebRTC SFU** — Minimal room relay logic via Pion, multi-viewer support
- **WHIP / WHEP** — HTTP API compatible with modern browsers and OBS WHIP plugin
- **Optional Auth** — Bearer token or X-Auth-Token header when `AUTH_TOKEN` is configured
- **Room Status** — `GET /api/rooms` returns online rooms, publisher & subscriber stats
- **Health Check** — `GET /healthz` for liveness probes
- **Embedded Frontend** — Simple publish/play pages with room & token input
- **Recording** — Optional VP8/VP9→IVF, Opus→OGG recording (`RECORD_ENABLED=1`)
- **Prometheus Metrics** — `GET /metrics` exposes RTP bytes/packets, subscribers, rooms
- **Containerized** — Dockerfile + docker-compose.yml with recording volume mount

## Quick Start

```bash
git clone https://github.com/LessUp/go-live.git
cd go-live
go mod tidy
go run ./cmd/server
```

- Publish page: http://localhost:8080/web/publisher.html
- Play page: http://localhost:8080/web/player.html
- Room list: http://localhost:8080/api/rooms
- Health check: http://localhost:8080/healthz

## HTTP API

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/whip/publish/{room}` | SDP Offer → Answer, establish publish connection |
| `POST` | `/api/whep/play/{room}` | SDP Offer → Answer, establish play connection |
| `GET` | `/api/rooms` | Online room list with stats |
| `GET` | `/healthz` | Health check |
| `GET` | `/metrics` | Prometheus metrics |

## Configuration

Environment variables: `ADDR`, `AUTH_TOKEN`, `STUN_SERVERS`, `TURN_SERVERS`, `RECORD_ENABLED`, `RECORD_DIR`, `CORS_ORIGINS`, `MAX_SUBSCRIBERS_PER_ROOM`.

## License

MIT License
