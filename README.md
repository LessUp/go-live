# live-webrtc-go

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)

English | [简体中文](README.zh-CN.md)

A lightweight live streaming service built with Go + [Pion WebRTC](https://github.com/pion/webrtc). It provides WHIP publishing, WHEP playback, embedded web pages, configurable auth, room status queries, recording, and Prometheus metrics.

## Features

- WebRTC SFU with simple room-based relay
- WHIP / WHEP endpoints for publish and play
- Embedded web UI under `/web/`
- Room list and bootstrap config endpoints
- Optional token / room-token / JWT auth
- Optional local recording (`.ivf` / `.ogg`)
- Prometheus metrics and `/healthz`
- Dockerfile and docker-compose support

## Requirements

- Go 1.22+
- A browser with WebRTC support
- Optional: Docker / Docker Compose

## Quick Start

```bash
git clone https://github.com/LessUp/go-live.git
cd go-live
go run ./cmd/server
```

Or use the local helper script:

```bash
./scripts/start.sh
```

The helper script prepares local cache directories and loads `.env.local` if present.
It does **not** run `go mod tidy` by default anymore. If you want that behavior, use:

```bash
RUN_TIDY=1 ./scripts/start.sh
```

Useful URLs after startup:

- Home: http://localhost:8080/web/index.html
- Publisher: http://localhost:8080/web/publisher.html
- Player: http://localhost:8080/web/player.html
- Records page: http://localhost:8080/web/records.html
- Bootstrap config: http://localhost:8080/api/bootstrap
- Rooms: http://localhost:8080/api/rooms
- Health check: http://localhost:8080/healthz
- Metrics: http://localhost:8080/metrics

## HTTP API

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/whip/publish/{room}` | SDP Offer → Answer, establish publish connection |
| `POST` | `/api/whep/play/{room}` | SDP Offer → Answer, establish play connection |
| `GET` | `/api/bootstrap` | Browser runtime config |
| `GET` | `/api/rooms` | Online room list and stats |
| `GET` | `/api/records` | Recording list metadata |
| `POST` | `/api/admin/rooms/{room}/close` | Close a room with admin auth |
| `GET` | `/healthz` | Health check |
| `GET` | `/metrics` | Prometheus metrics |

Room names are restricted to `A-Z a-z 0-9 _ -` and max length 64.

## Configuration

Important environment variables:

| Variable | Default | Description |
|---|---|---|
| `HTTP_ADDR` | `:8080` | HTTP listen address |
| `ALLOWED_ORIGIN` | `*` | Allowed CORS origin |
| `AUTH_TOKEN` | empty | Global token |
| `ROOM_TOKENS` | empty | Room-specific tokens: `room1:tok1;room2:tok2` |
| `JWT_SECRET` | empty | HMAC JWT secret |
| `JWT_AUDIENCE` | empty | Required JWT audience when set |
| `ADMIN_TOKEN` | empty | Admin token for room close API |
| `STUN_URLS` | Google STUN | Comma-separated STUN servers |
| `TURN_URLS` | empty | Comma-separated TURN servers |
| `TURN_USERNAME` | empty | TURN username |
| `TURN_PASSWORD` | empty | TURN password |
| `RECORD_ENABLED` | `0` | Enable recording when `1` |
| `RECORD_DIR` | `records` | Recording directory |
| `UPLOAD_RECORDINGS` | `0` | Enable object-storage upload |
| `DELETE_RECORDING_AFTER_UPLOAD` | `0` | Delete local file after upload |
| `S3_ENDPOINT` | empty | S3 / MinIO endpoint |
| `S3_REGION` | empty | S3 region |
| `S3_BUCKET` | empty | Target bucket |
| `S3_ACCESS_KEY` | empty | Access key |
| `S3_SECRET_KEY` | empty | Secret key |
| `S3_USE_SSL` | `1` | Use SSL for S3 |
| `S3_PATH_STYLE` | `0` | Use path-style access |
| `S3_PREFIX` | empty | Object key prefix |
| `MAX_SUBS_PER_ROOM` | `0` | Per-room subscriber limit |
| `RATE_LIMIT_RPS` | `0` | Per-IP rate limit |
| `RATE_LIMIT_BURST` | `0` | Rate-limit burst |
| `TLS_CERT_FILE` | empty | TLS cert path |
| `TLS_KEY_FILE` | empty | TLS key path |
| `PPROF` | `0` | Reserved config flag |

## Development and Verification

Common commands:

```bash
make build
make fmt
make lint
make security
make test
make test-all
make coverage
make ci
```

Default local verification (`make test`) runs:

- unit tests
- integration tests
- security tests

`make test-all` additionally runs:

- e2e tests
- performance tests

## Recording behavior

- `/api/records` returns a JSON list of local recordings
- If `RECORD_DIR` does not exist yet, the API returns an empty list
- Static `/records/` file serving is mounted only when recording is enabled

## Docker

```bash
docker build -t live-webrtc-go:latest .
```

```bash
docker run --rm -p 8080:8080 \
  -e RECORD_ENABLED=1 \
  -e RECORD_DIR=/records \
  -v "$PWD/records:/records" \
  live-webrtc-go:latest
```

## License

MIT License
