---
layout: default
title: 使用指南
---

# 使用指南

本文覆盖本地开发、容器部署与常见排障技巧，帮助你快速体验 live-webrtc-go。

## 前置要求

- Go 1.22 及以上
- （可选）Docker / Docker Compose
- 浏览器支持 WebRTC，并允许访问摄像头/麦克风

## 本地开发

### 方式一：使用启动脚本

```bash
./scripts/start.sh
```

脚本会：

- 创建 `records`、`.gocache`、`.gomodcache`
- 把 `GOCACHE` / `GOMODCACHE` 指向仓库目录内
- 加载 `.env.local`（如果存在）

如需在启动前执行 `go mod tidy`：

```bash
RUN_TIDY=1 ./scripts/start.sh
```

### 方式二：手动运行

```bash
go run ./cmd/server
```

启动后常用地址：

- 推流：`http://localhost:8080/web/publisher.html`
- 播放：`http://localhost:8080/web/player.html`
- 房间列表：`http://localhost:8080/api/rooms`
- 录制列表：`http://localhost:8080/api/records`
- 指标：`http://localhost:8080/metrics`
- 健康检查：`http://localhost:8080/healthz`

## Docker / Compose

### 构建镜像

```bash
docker build -t live-webrtc-go:latest .
```

### 直接运行容器

```bash
docker run --rm -p 8080:8080 \
  -e RECORD_ENABLED=1 \
  -e RECORD_DIR=/records \
  -v "$PWD/records:/records" \
  live-webrtc-go:latest
```

### 使用 docker compose

```bash
docker compose up -d
```

## 常用环境变量

| 变量 | 作用 |
|------|------|
| `HTTP_ADDR` | HTTP 监听地址，默认 `:8080` |
| `ALLOWED_ORIGIN` | CORS 白名单 |
| `AUTH_TOKEN` / `ROOM_TOKENS` | 全局 / 房间级鉴权 |
| `JWT_SECRET` / `JWT_AUDIENCE` | JWT 鉴权与 audience 校验 |
| `RECORD_ENABLED` / `RECORD_DIR` | 控制录制与输出目录 |
| `UPLOAD_RECORDINGS` 及 S3 相关变量 | 启用录制上传 |
| `MAX_SUBS_PER_ROOM` | 每房间订阅者上限 |
| `RATE_LIMIT_RPS` / `RATE_LIMIT_BURST` | HTTP 限流阈值 |
| `ADMIN_TOKEN` | 管理接口 Token |

## API 示例

```bash
curl -X POST "http://localhost:8080/api/whip/publish/demo" \
  -H "Content-Type: application/sdp" \
  --data-binary @offer.sdp

curl http://localhost:8080/api/rooms | jq

curl -X POST \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  http://localhost:8080/api/admin/rooms/demo/close
```

## 验证命令

```bash
make test
make test-all
make ci
```

## 常见问题

- **提示 `publisher already exists`**：同一房间只允许一个发布者
- **浏览器无法建立连接**：检查 HTTPS / TURN 配置
- **录制文件缺失**：确认 `RECORD_ENABLED=1` 且进程对 `RECORD_DIR` 有写权限
- **`/api/records` 返回空数组**：如果录制目录还不存在，这是预期行为
- **限流误伤**：可将 `RATE_LIMIT_RPS=0` 关闭限流，或调大 `RATE_LIMIT_BURST`

如需更深入背景，可继续阅读 `docs/design.md`。
