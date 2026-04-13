# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project overview

Go 1.22+ WebRTC SFU server using [Pion](https://github.com/pion/webrtc). Supports WHIP publish, WHEP playback, room-based relay, recording, and Prometheus metrics.

## Build & run

```bash
go build -o bin/live-webrtc-go ./cmd/server
go run ./cmd/server
```

Dev helper (loads `.env.local`, prepares cache dirs): `./scripts/start.sh`

## Test commands

- `make test` — unit + integration + security tests (default verification)
- `make test-all` — adds e2e + performance suites
- `make test-unit` — unit tests only (`go test -v -race ./internal/...`)
- `make test-integration` — requires `-tags=integration`
- `make test-e2e` — requires `-tags=e2e`, 10m timeout
- `make coverage` — generates HTML + XML coverage reports

## Lint & format

- `make lint` — runs `gofmt -s` check + `go vet` + `golangci-lint` (no project config, uses defaults)
- `make fmt` — formats with `gofmt -s -w .`
- `make security` — runs `gosec`

## Key env vars

See README "Configuration" table. Critical ones: `HTTP_ADDR`, `RECORD_ENABLED`, `RECORD_DIR`, `ALLOWED_ORIGIN`. Auth via `AUTH_TOKEN`, `ROOM_TOKENS`, or `JWT_SECRET`.

## Repo conventions

- Primary branch: `master`
- Entry point: `cmd/server`
- Core packages: `internal/{api,config,metrics,sfu,testutil,uploader}`
- Tagged test suites under `test/{e2e,integration,load,performance,security}`
- Room names: `A-Z a-z 0-9 _ -`, max length 64

## Working style

- Propose a plan before implementing non-trivial changes
- Explain tradeoffs when multiple approaches exist
- Be terse — no trailing summaries
