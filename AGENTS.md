# AGENTS.md

> 本文件为 AI Agent（包括 Claude、Copilot、GLM 等）提供项目级上下文和操作指南。

## 项目概述

**Go-Live** 是一个轻量级 WebRTC SFU（Selective Forwarding Unit）服务器，使用 Go 1.22+ 和 Pion WebRTC 构建。支持 WHIP/WHEP 协议，实现浏览器/OBS 推流、观众订阅播放、录制存储等功能。

## 核心架构

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
└─────────────┘        RTP Packets (WebRTC)             └─────────────┘
```

### 三层架构

| 层级 | 文件 | 职责 |
|------|------|------|
| **Manager** | `internal/sfu/manager.go` | 房间生命周期管理、创建/销毁 Room |
| **Room** | `internal/sfu/room.go` | 单发布者、多订阅者、录制管理 |
| **trackFanout** | `internal/sfu/track.go` | RTP 包分发，每个订阅者一个 goroutine |

### 关键约束

1. **单发布者模型**: 每个 Room 只允许一个 publisher，重复发布返回 403 Forbidden
2. **房间名称**: 正则 `^[A-Za-z0-9_-]{1,64}$`，在 `internal/api/routes.go` 校验
3. **SDP 大小限制**: 最大 1MB，防止内存攻击
4. **并发安全**: 所有共享状态使用 `sync.Mutex` 或 `sync.RWMutex` 保护

## 关键代码路径

### 发布流程
```
POST /api/whip/publish/{room}
  → middleware: CORS → RateLimit → Auth
  → handlers.HandleWHIPPublish()
  → manager.GetOrCreateRoom(room)
  → room.SetPublisher(offer)
  → 返回 SDP Answer
```

### 订阅流程
```
POST /api/whep/play/{room}
  → middleware: CORS → RateLimit → Auth
  → handlers.HandleWHEPPlay()
  → manager.GetRoom(room)
  → room.AddSubscriber(offer)
  → 返回 SDP Answer
```

### RTP 分发
```
Publisher → onTrack callback
  → trackFanout.AddTrack(track)
  → goroutine: read RTP packets
  → for each subscriber: subscriber.WriteRTP(packet)
```

## 认证系统

### 双模式认证

| 模式 | 配置 | 适用场景 |
|------|------|----------|
| **Token** | `AUTH_TOKEN` / `ROOM_TOKENS` | 简单场景、内部服务 |
| **JWT** | `JWT_SECRET` + claims | 生产环境、细粒度控制 |

### JWT Claims 支持
- `role=admin` 或 `admin=true`: 管理员权限
- `room=<name>`: 限制特定房间
- `aud=<audience>`: 受众验证（可选）

### 权限检查
```go
// 管理员端点
if !adminOK(r, c.AdminToken) {
    http.Error(w, "unauthorized", 401)
    return
}

// 房间端点
if !authOKRoom(r, c, room) {
    http.Error(w, "unauthorized", 401)
    return
}
```

## 录制系统

### 格式映射
| 媒体类型 | 编解码器 | 容器格式 |
|----------|----------|----------|
| Video | VP8/VP9 | IVF |
| Audio | Opus | OGG (48kHz, stereo) |

### 存储流程
```
Room.startRecording()
  → 检测 track 编解码器
  → 创建 IVF/OGG writer
  → goroutine: write packets to file
  → on close: upload to S3 (if configured)
```

### 文件命名
`{room}_{trackID}_{timestamp}.{ivf|ogg}`

## API 端点参考

| Method | Path | Auth | 描述 |
|--------|------|------|------|
| `POST` | `/api/whip/publish/{room}` | Token/JWT | 发布流（SDP Offer → Answer）|
| `POST` | `/api/whep/play/{room}` | Token/JWT | 播放流（SDP Offer → Answer）|
| `GET` | `/api/bootstrap` | None | 前端运行时配置 |
| `GET` | `/api/rooms` | None | 列出活跃房间 |
| `GET` | `/api/records` | None | 列出录制文件 |
| `POST` | `/api/admin/rooms/{room}/close` | Admin | 强制关闭房间 |
| `GET` | `/healthz` | None | 健康检查 |
| `GET` | `/metrics` | None | Prometheus 指标 |

## 代码风格规范

### 错误处理
```go
// ✅ 正确：包装错误并添加上下文
if err != nil {
    return fmt.Errorf("failed to create room: %w", err)
}

// ❌ 错误：忽略错误
_ = conn.Close()

// ✅ 正确：记录错误
if err := conn.Close(); err != nil {
    slog.Error("failed to close connection", "error", err)
}
```

### 并发模式
```go
// ✅ 正确：使用 mutex 保护共享状态
type Room struct {
    mu          sync.RWMutex
    subscribers map[string]*PeerConnection
}

func (r *Room) AddSubscriber(id string, pc *PeerConnection) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.subscribers[id] = pc
}
```

### 日志规范
```go
// 使用 slog 结构化日志
slog.Info("publisher joined",
    "room", roomName,
    "track_id", trackID,
)

slog.Error("failed to start recording",
    "room", roomName,
    "error", err,
)
```

## 常见任务模板

### 添加新的 API 端点
1. 在 `internal/api/handlers.go` 添加 handler 函数
2. 在 `internal/api/routes.go` 注册路由
3. 添加必要的 auth 检查（`authOKRoom` 或 `adminOK`）
4. 添加 rate limit 检查（`allowRate`）
5. 添加 CORS headers（`allowCORS`）
6. 在 `internal/api/handlers_test.go` 编写测试

### 添加新的配置项
1. 在 `internal/config/config.go` 的 `Config` struct 添加字段
2. 在 `Load()` 函数中解析环境变量
3. 在 `internal/config/config_test.go` 添加测试
4. 更新 README.md 和 docs/ 中的配置说明

### 修改 SFU 行为
1. 理解 `Manager` → `Room` → `trackFanout` 层级
2. 注意 `internal/sfu/track.go` 中 goroutine 的生命周期
3. 修改时考虑并发安全
4. 更新 `internal/metrics/` 中的相关指标

## WebRTC 调试技巧

### ICE 连接问题
```bash
# 检查 STUN/TURN 配置
STUN_URLS=stun:stun.l.google.com:19302
TURN_URLS=turn:turn.example.com:3478
TURN_USERNAME=user
TURN_PASSWORD=pass

# 查看详细日志
GODEBUG=nettrace=1 go run ./cmd/server
```

### SDP 分析
```bash
# 使用 sdp-analyze 工具
go install github.com/pion/sdp/v3/cmd/sdp-analyze@latest

# 打印 SDP 详情
curl -X POST http://localhost:8080/api/whip/publish/test \
  -H "Authorization: Bearer $TOKEN" \
  -d @offer.sdp | sdp-analyze
```

### 常见 WebRTC 错误

| 错误 | 原因 | 解决方案 |
|------|------|----------|
| ICE failed | NAT 穿透失败 | 配置 TURN 服务器 |
| DTLS error | 证书问题 | 检查 TLS 配置 |
| No tracks | 编解码器不匹配 | 确认 VP8/VP9/Opus 支持 |

## OpenSpec 工作流

本项目采用 OpenSpec 规范驱动开发：

```
openspec/
├── specs/           # 规范文件（Single Source of Truth）
│   ├── whip-whep/   # WHIP/WHEP 协议规范
│   ├── room-sfu/    # SFU 行为规范
│   ├── authentication/ # 认证规范
│   └── ...
└── config.yaml      # OpenSpec 配置
```

### 开发流程
1. **探索**: `/opsx:explore` 分析需求和问题
2. **提案**: `/opsx:propose` 创建变更提案
3. **实现**: `/opsx:apply` 执行任务
4. **归档**: `/opsx:archive` 合并规范

## 测试策略

### 测试分类
| 类型 | 目录 | 标签 | 命令 |
|------|------|------|------|
| Unit | `internal/*_test.go` | - | `make test-unit` |
| Integration | `test/integration/` | `integration` | `make test-integration` |
| Security | `test/security/` | `security` | `make test-security` |
| E2E | `test/e2e/` | `e2e` | `make test-e2e` |
| Performance | `test/performance/` | `performance` | `make test-performance` |

### 测试原则
- 使用 table-driven 测试
- 独立测试使用 `t.Parallel()`
- Mock 外部依赖（如 S3 上传）

## 依赖版本

| 依赖 | 版本 | 用途 |
|------|------|------|
| `github.com/pion/webrtc/v3` | v3.2.24 | WebRTC 实现 |
| `github.com/golang-jwt/jwt/v5` | v5.2.0 | JWT 认证 |
| `github.com/minio/minio-go/v7` | v7.0.66 | S3 上传 |
| `github.com/prometheus/client_golang` | v1.19.0 | Prometheus 指标 |
| `go.opentelemetry.io/otel` | v1.28.0 | 分布式追踪 |

## 安全注意事项

1. **永远不要提交**: 密钥、Token、凭证文件
2. **输入验证**: 所有用户输入（房间名、SDP 大小）
3. **Token 比较**: 使用 `crypto/subtle.ConstantTimeCompare`
4. **定期扫描**: `make security` 运行 gosec

## 相关文档

- [README.md](./README.md) - 英文文档
- [README.zh-CN.md](./README.zh-CN.md) - 中文文档
- [CLAUDE.md](./CLAUDE.md) - Claude Code 指南
- [openspec/specs/](./openspec/specs/) - 规范文件
