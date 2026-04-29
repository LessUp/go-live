# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A lightweight **WebRTC SFU (Selective Forwarding Unit)** server built with Go 1.22+ and [Pion WebRTC](https://github.com/pion/webrtc). It provides real-time live streaming capabilities through WHIP/WHEP protocols with room-based relay, optional recording, authentication, rate limiting, and observability features.

### Key Features
- **WHIP Publishing**: Browser or OBS can push streams via HTTP POST
- **WHEP Playback**: Viewers subscribe to rooms via HTTP POST
- **Room-based SFU**: Single publisher per room, multiple subscribers
- **Recording**: VP8/VP9 → IVF, Opus → OGG with optional S3/MinIO upload
- **Auth**: Token (global/room-level) or JWT with role-based access
- **Observability**: Prometheus metrics, OpenTelemetry tracing, health checks

## OpenSpec Workflow

This project uses [OpenSpec](https://github.com/Fission-AI/OpenSpec) for spec-driven development. All specifications are in `openspec/specs/` as the source of truth.

### Commands

| Command | Purpose |
|---------|---------|
| `/opsx:propose <change>` | Start a new change proposal |
| `/opsx:explore` | Investigate problems and clarify requirements |
| `/opsx:apply` | Implement tasks from the proposal |
| `/opsx:archive` | Archive completed change and merge specs |

### Workflow

1. **Review specs** - Check `openspec/specs/` for existing requirements
2. **Propose change** - Use `/opsx:propose` to create:
   - `proposal.md` - Why and what changes
   - `specs/` - Delta specs (ADDED/MODIFIED/REMOVED)
   - `design.md` - Technical approach
   - `tasks.md` - Implementation checklist
3. **Implement** - Use `/opsx:apply` to work through tasks
4. **Archive** - Use `/opsx:archive` to merge delta specs and preserve history

### Spec Format

Requirements use this structure:
```markdown
### Requirement: <name>
Description with SHALL/MUST keywords.

#### Scenario: <name>
- **WHEN** condition
- **THEN** expected outcome
```

### Delta Specs

Changes are tracked using delta operations:
- `## ADDED Requirements` - New capabilities
- `## MODIFIED Requirements` - Changed behavior (include full updated content)
- `## REMOVED Requirements` - Deprecated features (include Reason and Migration)

### Capabilities

| Capability | Location | Description |
|------------|----------|-------------|
| WHIP/WHEP | `openspec/specs/whip-whep/` | Protocol specs |
| Room SFU | `openspec/specs/room-sfu/` | SFU behavior |
| Auth | `openspec/specs/authentication/` | Auth system |
| Recording | `openspec/specs/recording/` | Recording specs |
| Observability | `openspec/specs/observability/` | Metrics/tracing |
| API | `openspec/specs/api/` | API overview |
| Testing | `openspec/specs/testing/` | Test specs |
| RFCs | `openspec/specs/rfc/` | Architecture decisions |

## Architecture

```
┌─────────────┐   WHIP (POST /api/whip/publish/{room})   ┌──────────────┐
│  Publisher  │ ──────────────────────────────────────▶ │ HTTP Server  │
│  (OBS/Web)  │                                         │   :8080      │
└─────────────┘                                         └──────┬───────┘
                                                               │
                                      Creates PeerConnection   │
                                      & TrackFanout            ▼
                                                        ┌─────────────┐
┌─────────────┐   WHEP (POST /api/whep/play/{room})    │    Room     │
│   Viewer    │ ──────────────────────────────────────▶ │   (SFU)     │
│  (Browser)  │ ◀────────────────────────────────────── │             │
└─────────────┘        RTP Packets (WebRTC)             └──────┬──────┘
                                                               │
                                          Records & Uploads    ▼
                                                        ┌─────────────┐
                                                        │ Object Store│
                                                        │  (S3/MinIO) │
                                                        └─────────────┘
```

## Project Structure

```
├── cmd/server/           # Entry point
│   ├── main.go           # HTTP server setup, graceful shutdown
│   └── web/              # Embedded static files (publisher/player HTML)
├── internal/
│   ├── api/              # HTTP layer
│   │   ├── handlers.go   # WHIP/WHEP/Rooms/Records/Admin endpoints
│   │   ├── middleware.go # CORS, rate limiting, auth (token/JWT)
│   │   └── routes.go     # URL routing, room name validation
│   ├── config/           # Environment variable configuration
│   ├── sfu/              # Core WebRTC SFU logic
│   │   ├── manager.go    # Room lifecycle management
│   │   ├── room.go       # PeerConnection, track fanout, recording
│   │   └── track.go      # RTP packet distribution to subscribers
│   ├── metrics/          # Prometheus gauges/counters
│   ├── otel/             # OpenTelemetry tracer initialization
│   ├── uploader/         # S3/MinIO upload client
│   └── testutil/         # Test helpers
├── test/
│   ├── integration/      # Integration tests (requires -tags=integration)
│   ├── e2e/              # End-to-end tests (requires -tags=e2e)
│   ├── security/         # Security tests
│   ├── performance/      # Performance benchmarks
│   └── load/             # Load testing tools
├── docs/                 # GitHub Pages documentation
└── web/                  # Source static files (embedded into binary)
```

## Build & Run

```bash
# Build
go build -o bin/live-webrtc-go ./cmd/server

# Run directly
go run ./cmd/server

# Development helper (loads .env.local, sets cache dirs)
./scripts/start.sh

# With module tidy
RUN_TIDY=1 ./scripts/start.sh
```

## Test Commands

```bash
make test          # Unit + integration + security (default)
make test-all      # Adds e2e + performance (longer timeout)
make test-unit     # go test -v -race ./internal/...
make test-integration  # Requires -tags=integration
make test-e2e      # Requires -tags=e2e, 10m timeout
make coverage      # HTML + XML coverage reports
```

## Lint & Security

```bash
make lint        # gofmt -s + go vet + golangci-lint
make fmt         # gofmt -s -w .
make security    # gosec ./...
```

## Key Environment Variables

| Variable | Default | Purpose |
|----------|---------|---------|
| `HTTP_ADDR` | `:8080` | Listen address |
| `ALLOWED_ORIGIN` | `*` | CORS origin |
| `AUTH_TOKEN` | - | Global auth token |
| `ROOM_TOKENS` | - | Per-room tokens (`room1:tok1;room2:tok2`) |
| `JWT_SECRET` | - | HMAC secret for JWT auth |
| `ADMIN_TOKEN` | - | Admin API token |
| `RECORD_ENABLED` | `0` | Enable recording (`1` to enable) |
| `RECORD_DIR` | `records` | Recording output directory |
| `UPLOAD_RECORDINGS` | `0` | Enable S3 upload |
| `S3_*` | - | S3/MinIO configuration |
| `RATE_LIMIT_RPS` | `0` | Requests per second per IP |
| `MAX_SUBS_PER_ROOM` | `0` | Subscriber limit (0 = unlimited) |
| `STUN_URLS` | Google STUN | Comma-separated STUN servers |
| `TURN_URLS` | - | Comma-separated TURN servers |

## Code Conventions

### Go Style
- Follow standard Go idioms and [Effective Go](https://golang.org/doc/effective_go)
- Use `gofmt -s` for formatting (enforced by `make lint`)
- Run `go vet` and `golangci-lint` before committing

### Error Handling
- Wrap errors with context: `fmt.Errorf("operation failed: %w", err)`
- Log errors with `slog` at appropriate levels
- Never silently ignore errors in production code

### Concurrency
- Protect shared state with `sync.Mutex` or `sync.RWMutex`
- Use `context.Context` for cancellation and timeouts
- Ensure goroutines can exit (avoid leaks)

### Testing
- Write table-driven tests for multiple cases
- Use `t.Parallel()` where tests are independent
- Mock external dependencies (e.g., `uploadRecordingFile` in room.go)

### Naming
- Room names: `^[A-Za-z0-9_-]{1,64}$` (enforced in routes.go)
- Exported functions need documentation comments
- Use descriptive variable names (avoid single letters except loops)

## Repository Conventions

- **Primary branch**: `master`
- **Commit messages**: Present tense, concise ("Add feature" not "Added feature")
- **PRs**: Include summary, motivation, test evidence
- **Versioning**: Semantic Versioning (MAJOR.MINOR.PATCH)

## Working Style

1. **Plan before implementing** non-trivial changes
2. **Explain tradeoffs** when multiple approaches exist
3. **Keep changes focused** - one logical change per PR
4. **Update documentation** when behavior changes
5. **Be terse** - no trailing summaries in code review

## Common Tasks

### Adding a new API endpoint
1. Add handler in `internal/api/handlers.go`
2. Add route in `internal/api/routes.go`
3. Add auth check if needed (use `authOKRoom` or `adminOK`)
4. Add rate limit check with `allowRate`
5. Add CORS headers with `allowCORS`
6. Write tests in `internal/api/handlers_test.go`

### Adding a new configuration option
1. Add field to `Config` struct in `internal/config/config.go`
2. Parse from environment in `Load()` function
3. Add test in `internal/config/config_test.go`
4. Update README.md and docs/usage.md

### Modifying SFU behavior
1. Understand `Manager` → `Room` → `trackFanout` hierarchy
2. Be careful with goroutine lifecycle in `track.go`
3. Update metrics if adding new observables
4. Test with multiple publishers/subscribers

## Debugging

```bash
# Enable pprof endpoints
PPROF=1 go run ./cmd/server

# Access pprof
open http://localhost:8080/debug/pprof/

# View metrics
curl http://localhost:8080/metrics

# Check health
curl http://localhost:8080/healthz
```

### WebRTC 调试技巧

#### ICE 连接问题排查
```bash
# 1. 检查 STUN/TURN 配置
echo $STUN_URLS
echo $TURN_URLS

# 2. 测试 STUN 连通性
stun stun.l.google.com:19302

# 3. 启用 WebRTC 详细日志
GODEBUG=nettrace=1 go run ./cmd/server

# 4. 查看 ICE 候选
# 在浏览器控制台查看 RTCPeerConnection.iceGatheringState
```

#### SDP 分析
```bash
# 保存 SDP offer
curl -X POST http://localhost:8080/api/whip/publish/test \
  -H "Authorization: Bearer $TOKEN" \
  -d @offer.sdp > answer.sdp

# 分析 SDP 内容
cat answer.sdp | grep -E "m=|a=rtpmap|a=fmtp"
```

#### 常见 WebRTC 错误

| 错误 | 原因 | 解决方案 |
|------|------|----------|
| ICE failed | NAT 穿透失败 | 配置 TURN 服务器 `TURN_URLS` |
| DTLS error | 证书问题 | 检查 TLS 配置，确保时间同步 |
| No tracks received | 编解码器不匹配 | 确认 VP8/VP9/Opus 支持 |
| Connection timeout | 防火墙阻断 | 开放 UDP 端口或使用 TURN |
| High latency | 网络拥塞 | 检查带宽、启用 simulcast |

#### 性能调优
```bash
# 查看当前 goroutine 数量
curl http://localhost:8080/debug/pprof/goroutine?debug=1

# 内存分析
go tool pprof http://localhost:8080/debug/pprof/heap

# CPU 分析 (30秒采样)
go tool pprof http://localhost:8080/debug/pprof/profile?seconds=30
```

## Security Considerations

- Never commit secrets, tokens, or credentials
- Validate all user inputs (room names, SDP size limits)
- Use `crypto/subtle.ConstantTimeCompare` for token comparison
- Run `make security` for security scans
- Review auth changes carefully

## 常见错误 FAQ

### 认证错误

| 错误信息 | 原因 | 解决方案 |
|----------|------|----------|
| `401 Unauthorized` | Token 缺失或无效 | 检查 `Authorization: Bearer $TOKEN` |
| `403 Forbidden` | JWT 房间限制 | 确认 JWT `room` claim 匹配目标房间 |
| `403 Forbidden` | 房间已有发布者 | 先关闭现有发布者或使用新房间名 |

### 录制问题

| 问题 | 检查项 |
|------|--------|
| 录制文件未生成 | `RECORD_ENABLED=1`，检查 `RECORD_DIR` 权限 |
| 文件损坏 | 确认编解码器为 VP8/VP9/Opus |
| S3 上传失败 | 检查 `S3_*` 配置和网络连通性 |

### 连接问题

| 症状 | 诊断步骤 |
|------|----------|
| 观众无法连接 | 1. 检查发布者是否在线 2. 验证认证 3. 检查 ICE 候选 |
| 画面卡顿 | 1. 检查带宽 2. 查看服务器 CPU/内存 3. 检查丢包率 |
| 延迟过高 | 1. 确认 TURN 服务器位置 2. 检查网络路由 |

### 配置问题

```bash
# 验证配置加载
go run ./cmd/server -h 2>&1 | head -20

# 检查环境变量
env | grep -E "AUTH|RECORD|STUN|TURN|S3"

# 测试配置
curl -v http://localhost:8080/api/bootstrap
```
