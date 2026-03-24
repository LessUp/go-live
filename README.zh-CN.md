# live-webrtc-go

[![Docs](https://img.shields.io/badge/Docs-GitHub%20Pages-blue?logo=github)](https://lessup.github.io/go-live/)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

[English](README.md) | 简体中文

使用 Go + [Pion WebRTC](https://github.com/pion/webrtc) 构建的轻量级在线直播服务。提供 WHIP 推流、WHEP 播放、嵌入式 Web 页面、可配置鉴权、房间状态查询、录制与 Prometheus 指标。

## 功能特点

- 基于房间的轻量级 WebRTC SFU
- WHIP / WHEP 推流与播放接口
- 内嵌 Web 页面（`/web/`）
- `/api/bootstrap` 与 `/api/rooms` 运行时查询接口
- 可选全局 Token / 房间 Token / JWT 鉴权
- 可选本地录制（`.ivf` / `.ogg`）
- `/metrics` 与 `/healthz`
- Dockerfile 与 docker-compose 支持

## 运行要求

- Go 1.22+
- 支持 WebRTC 的浏览器
- 可选：Docker / Docker Compose

## 快速开始

```bash
git clone https://github.com/LessUp/go-live.git
cd go-live
go run ./cmd/server
```

也可以直接使用本地启动脚本：

```bash
./scripts/start.sh
```

脚本会准备本地缓存目录，并在存在时加载 `.env.local`。
现在脚本**默认不会**执行 `go mod tidy`。如需执行：

```bash
RUN_TIDY=1 ./scripts/start.sh
```

启动后可访问：

- 首页：http://localhost:8080/web/index.html
- 推流页：http://localhost:8080/web/publisher.html
- 播放页：http://localhost:8080/web/player.html
- 录制列表页：http://localhost:8080/web/records.html
- 运行时配置：http://localhost:8080/api/bootstrap
- 房间列表：http://localhost:8080/api/rooms
- 健康检查：http://localhost:8080/healthz
- 指标：http://localhost:8080/metrics

## HTTP API

| 方法 | 路径 | 说明 |
|------|------|------|
| `POST` | `/api/whip/publish/{room}` | SDP Offer → Answer，建立推流连接 |
| `POST` | `/api/whep/play/{room}` | SDP Offer → Answer，建立播放连接 |
| `GET` | `/api/bootstrap` | 浏览器运行时配置 |
| `GET` | `/api/rooms` | 房间列表与状态 |
| `GET` | `/api/records` | 录制文件元数据 |
| `POST` | `/api/admin/rooms/{room}/close` | 使用管理员鉴权关闭房间 |
| `GET` | `/healthz` | 健康检查 |
| `GET` | `/metrics` | Prometheus 指标 |

房间名限制为 `A-Z a-z 0-9 _ -`，最大长度 64。

## 关键环境变量

| 变量 | 默认值 | 说明 |
|---|---|---|
| `HTTP_ADDR` | `:8080` | HTTP 监听地址 |
| `ALLOWED_ORIGIN` | `*` | 允许的 CORS 来源 |
| `AUTH_TOKEN` | 空 | 全局令牌 |
| `ROOM_TOKENS` | 空 | 房间级令牌，格式 `room1:tok1;room2:tok2` |
| `JWT_SECRET` | 空 | HMAC JWT 密钥 |
| `JWT_AUDIENCE` | 空 | 设置后要求 JWT audience 匹配 |
| `ADMIN_TOKEN` | 空 | 管理接口令牌 |
| `STUN_URLS` | Google STUN | STUN 列表 |
| `TURN_URLS` | 空 | TURN 列表 |
| `TURN_USERNAME` | 空 | TURN 用户名 |
| `TURN_PASSWORD` | 空 | TURN 密码 |
| `RECORD_ENABLED` | `0` | 设为 `1` 启用录制 |
| `RECORD_DIR` | `records` | 录制目录 |
| `UPLOAD_RECORDINGS` | `0` | 启用对象存储上传 |
| `DELETE_RECORDING_AFTER_UPLOAD` | `0` | 上传成功后删除本地文件 |
| `S3_ENDPOINT` | 空 | S3 / MinIO 地址 |
| `S3_REGION` | 空 | S3 区域 |
| `S3_BUCKET` | 空 | 目标桶 |
| `S3_ACCESS_KEY` | 空 | Access Key |
| `S3_SECRET_KEY` | 空 | Secret Key |
| `S3_USE_SSL` | `1` | 是否启用 SSL |
| `S3_PATH_STYLE` | `0` | 是否使用 path-style |
| `S3_PREFIX` | 空 | 上传对象前缀 |
| `MAX_SUBS_PER_ROOM` | `0` | 每房间订阅者上限 |
| `RATE_LIMIT_RPS` | `0` | 每 IP 限流速率 |
| `RATE_LIMIT_BURST` | `0` | 限流突发容量 |
| `TLS_CERT_FILE` | 空 | TLS 证书路径 |
| `TLS_KEY_FILE` | 空 | TLS 私钥路径 |
| `PPROF` | `0` | 预留调试配置 |

## 开发与验证

常用命令：

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

默认本地验证（`make test`）会运行：

- 单元测试
- 集成测试
- 安全测试

`make test-all` 额外运行：

- e2e 测试
- performance 测试

## 录制行为

- `/api/records` 返回本地录制文件 JSON 列表
- 若 `RECORD_DIR` 不存在，则返回空数组，不再返回 500
- 仅当启用录制时才会挂载静态 `/records/` 目录访问

## 容器运行

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

## 许可协议

MIT License
