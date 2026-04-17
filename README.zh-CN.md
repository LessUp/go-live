# Go-Live

[![CI](https://github.com/LessUp/go-live/actions/workflows/ci.yml/badge.svg)](https://github.com/LessUp/go-live/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/LessUp/go-live)](https://goreportcard.com/report/github.com/LessUp/go-live)
![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)
[![Release](https://img.shields.io/github/v/release/LessUp/go-live)](https://github.com/LessUp/go-live/releases)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue?logo=docker)](Dockerfile)
[![Docs](https://img.shields.io/badge/Docs-GitHub%20Pages-blue?logo=github)](https://lessup.github.io/go-live/)

[English](README.md) | **简体中文**

基于 Go 和 [Pion WebRTC](https://github.com/pion/webrtc) 构建的轻量级、高性能 **WebRTC SFU**（选择性转发单元）服务器。支持 WHIP/WHEP 协议推流、房间广播、录制功能和完整可观测性。

---

## 📋 目录

- [特性](#-特性)
- [系统架构](#️-系统架构)
- [快速开始](#-快速开始)
- [安装](#-安装)
- [配置](#️-配置)
- [API 参考](#-api-参考)
- [文档](#-文档)
- [开发](#️-开发)
- [Docker 部署](#-docker-部署)
- [参与贡献](#-参与贡献)
- [许可协议](#-许可协议)

---

## ✨ 特性

| 特性 | 说明 |
|------|------|
| 🎥 **WHIP/WHEP 协议** | 标准 HTTP 协议的 WebRTC 推流和播放，兼容 OBS 和现代浏览器 |
| 🏠 **房间 SFU 架构** | 每房间单发布者、多订阅者，高效的 RTP 转发 |
| 🔐 **灵活认证体系** | 支持全局 Token、房间级 Token、JWT 角色认证 |
| 📹 **录制与上传** | VP8/VP9 → IVF、Opus → OGG，支持 S3/MinIO 自动上传 |
| 📊 **完整可观测性** | Prometheus 指标、OpenTelemetry 追踪、健康检查端点 |
| 🐳 **云原生就绪** | 支持 Docker、Docker Compose，适配 Kubernetes |
| 🌐 **嵌入式 Web 界面** | 内置推流和播放页面，开箱即用 |
| ⚡ **高性能** | 基于 Go 并发模型的低延迟媒体转发 |

---

## 🏗️ 系统架构

```
                         ┌─────────────────────────────────────┐
                         │          HTTP Server :8080          │
                         │                                     │
    ┌──────────┐         │  ┌─────────┐    ┌─────────────┐    │
    │  推流端   │ ──WHIP──▶│  │  认证    │───▶│   SFU       │    │
    │(OBS/网页) │         │  │ 中间件    │    │  管理器      │    │
    └──────────┘         │  └─────────┘    └──────┬──────┘    │
                         │                        │           │
    ┌──────────┐         │        ┌───────────────┤           │
    │  观看端   │ ◀──WHEP──│        │               │           │
    │ (浏览器)  │───────▶  │        ▼               ▼           │
    └──────────┘         │ ┌─────────────┐ ┌───────────┐      │
                         │ │    房间      │ │   录制    │      │
                         │ │  (转发)      │ │  与上传   │      │
                         │ └──────┬──────┘ └─────┬─────┘      │
                         └────────┼──────────────┼────────────┘
                                  │              │
                        ┌─────────▼──────────────▼───────────┐
                        │          对象存储                   │
                        │         (S3/MinIO)                 │
                        └────────────────────────────────────┘
```

### 请求处理链

```
HTTP 请求 → CORS → 限流 → 认证 → 处理器 → SFU 房间
```

---

## 🚀 快速开始

### 从源码运行（30 秒）

```bash
git clone https://github.com/LessUp/go-live.git
cd go-live
go run ./cmd/server
```

### 使用 Docker 运行

```bash
docker run --rm -p 8080:8080 ghcr.io/lessup/go-live:latest
```

### 访问应用

| 页面 | 地址 | 说明 |
|------|------|------|
| 🏠 首页 | `http://localhost:8080/` | 前端控制台 |
| 📤 推流页 | `http://localhost:8080/web/publisher.html` | 浏览器推流 |
| 📥 播放页 | `http://localhost:8080/web/player.html` | 观看直播 |
| 📋 录制列表 | `http://localhost:8080/web/records.html` | 录制文件浏览 |
| 📊 指标监控 | `http://localhost:8080/metrics` | Prometheus 指标 |

---

## 📦 安装

### 二进制下载

从 [GitHub Releases](https://github.com/LessUp/go-live/releases) 下载预编译二进制文件：

```bash
# Linux AMD64
curl -LO https://github.com/LessUp/go-live/releases/latest/download/live-webrtc-go-linux-amd64
chmod +x live-webrtc-go-linux-amd64
./live-webrtc-go-linux-amd64
```

### 从源码构建

```bash
git clone https://github.com/LessUp/go-live.git
cd go-live
make build
./bin/server
```

---

## ⚙️ 配置

通过环境变量进行配置。创建 `.env.local` 文件：

```bash
# 核心配置
HTTP_ADDR=:8080
ALLOWED_ORIGIN=*

# 认证（可选）
AUTH_TOKEN=your-secret-token
ADMIN_TOKEN=your-admin-token

# 录制（可选）
RECORD_ENABLED=1
RECORD_DIR=records

# S3 上传（可选）
UPLOAD_RECORDINGS=1
S3_ENDPOINT=minio.example.com:9000
S3_ACCESS_KEY=minioadmin
S3_SECRET_KEY=minioadmin
S3_BUCKET=recordings
```

### 核心配置

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `HTTP_ADDR` | `:8080` | HTTP 监听地址 |
| `ALLOWED_ORIGIN` | `*` | CORS 允许来源 |

### 身份认证

| 变量 | 说明 |
|------|------|
| `AUTH_TOKEN` | 全局认证令牌 |
| `ROOM_TOKENS` | 房间级令牌：`room1:tok1;room2:tok2` |
| `JWT_SECRET` | JWT HMAC 签名密钥 |
| `ADMIN_TOKEN` | 管理接口令牌 |

### WebRTC/ICE

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `STUN_URLS` | `stun:stun.l.google.com:19302` | STUN 服务器 |
| `TURN_URLS` | - | TURN 服务器（NAT 环境必需） |
| `TURN_USERNAME` | - | TURN 用户名 |
| `TURN_PASSWORD` | - | TURN 密码 |

### 录制功能

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `RECORD_ENABLED` | `0` | 启用录制（设为 `1` 启用） |
| `RECORD_DIR` | `records` | 录制输出目录 |

> 💡 查看[完整配置指南](https://lessup.github.io/go-live/zh/usage.html#配置说明)了解所有选项，包括 S3 上传、限流、TLS 和调试设置。

---

## 🔌 API 参考

### 流媒体接口

| 方法 | 路径 | 认证 | 说明 |
|------|------|------|------|
| `POST` | `/api/whip/publish/{room}` | Token/JWT | 推流到房间 |
| `POST` | `/api/whep/play/{room}` | Token/JWT | 从房间播放 |

### 查询接口

| 方法 | 路径 | 说明 |
|------|------|------|
| `GET` | `/api/bootstrap` | 前端运行时配置 |
| `GET` | `/api/rooms` | 活跃房间列表 |
| `GET` | `/api/records` | 录制文件列表 |

### 管理接口

| 方法 | 路径 | 认证 | 说明 |
|------|------|------|------|
| `POST` | `/api/admin/rooms/{room}/close` | Admin Token | 强制关闭房间 |

### 健康与监控

| 方法 | 路径 | 说明 |
|------|------|------|
| `GET` | `/healthz` | 健康检查（返回 `ok`） |
| `GET` | `/metrics` | Prometheus 指标 |

### 快速测试

```bash
# 健康检查
curl http://localhost:8080/healthz

# 查看房间
curl http://localhost:8080/api/rooms

# 启动配置
curl http://localhost:8080/api/bootstrap | jq

# 关闭房间（管理员）
curl -X POST \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  http://localhost:8080/api/admin/rooms/myroom/close
```

> 💡 查看[完整 API 文档](https://lessup.github.io/go-live/zh/api.html)了解请求/响应格式和示例。

---

## 📚 文档

| 文档 | 链接 |
|------|------|
| English Docs | https://lessup.github.io/go-live/en/ |
| 中文文档 | https://lessup.github.io/go-live/zh/ |
| 使用指南 | https://lessup.github.io/go-live/zh/usage.html |
| 设计说明 | https://lessup.github.io/go-live/zh/design.html |
| API 参考 | https://lessup.github.io/go-live/zh/api.html |
| 更新日志 | https://lessup.github.io/go-live/changelog.html |

---

## 🛠️ 开发

### Makefile 命令

```bash
make build       # 编译二进制到 bin/
make test        # 运行所有测试（单元 + 集成 + 安全）
make test-unit   # 仅运行单元测试
make test-all    # 运行所有测试（含 e2e 和性能）
make lint        # 运行代码检查（gofmt + go vet + golangci-lint）
make security    # 运行 gosec 安全扫描
make coverage    # 生成覆盖率报告
make ci          # 完整 CI 流水线（lint + test + security）
```

### 项目结构

```
├── cmd/server/           # 应用入口
│   ├── main.go           # HTTP 服务初始化
│   └── web/              # 嵌入式静态文件
├── internal/
│   ├── api/              # HTTP 处理器和路由
│   ├── config/           # 配置管理
│   ├── sfu/              # WebRTC SFU 核心
│   ├── metrics/          # Prometheus 指标
│   ├── otel/             # OpenTelemetry 追踪
│   ├── uploader/         # S3/MinIO 上传
│   └── testutil/         # 测试工具
├── specs/                # 单一事实来源（规范）
├── test/                 # 测试实现
└── docs/                 # 文档
```

### 支持的编解码器

| 编解码器 | 类型 | 格式 |
|----------|------|------|
| VP8 | 视频 | IVF |
| VP9 | 视频 | IVF |
| Opus | 音频 | OGG |

---

## 🐳 Docker 部署

### 构建镜像

```bash
docker build -t live-webrtc-go:latest .
```

### 基础运行

```bash
docker run --rm -p 8080:8080 live-webrtc-go:latest
```

### Docker Compose

```yaml
services:
  live-webrtc:
    build: .
    ports:
      - "8080:8080"
    environment:
      - AUTH_TOKEN=${AUTH_TOKEN}
      - RECORD_ENABLED=1
      - RECORD_DIR=/records
    volumes:
      - ./records:/records
    restart: unless-stopped
```

```bash
docker compose up -d
```

### 完整配置

```bash
docker run --rm -p 8080:8080 \
  -e AUTH_TOKEN=mysecret \
  -e RECORD_ENABLED=1 \
  -e RECORD_DIR=/records \
  -e UPLOAD_RECORDINGS=1 \
  -e S3_ENDPOINT=minio:9000 \
  -e S3_ACCESS_KEY=minioadmin \
  -e S3_SECRET_KEY=minioadmin \
  -e S3_BUCKET=recordings \
  -v $(pwd)/records:/records \
  live-webrtc-go:latest
```

> 💡 查看 [Docker 部署指南](https://lessup.github.io/go-live/zh/usage.html#docker-部署)了解 Kubernetes 示例和生产最佳实践。

---

## 🤝 参与贡献

我们欢迎各种形式的贡献！请查看我们的[贡献指南](CONTRIBUTING.md)了解详情。

- [贡献指南](CONTRIBUTING.md)
- [行为准则](CODE_OF_CONDUCT.md)
- [安全策略](SECURITY.md)
- [规范驱动开发](AGENTS.md)

---

## 📄 许可协议

本项目基于 [MIT License](LICENSE) 开源。

---

## 🔗 相关链接

- [GitHub 仓库](https://github.com/LessUp/go-live)
- [问题反馈](https://github.com/LessUp/go-live/issues)
- [版本发布](https://github.com/LessUp/go-live/releases)
- [在线文档](https://lessup.github.io/go-live/)
- [Pion WebRTC](https://github.com/pion/webrtc)

---

<div align="center">

**[⬆ 回到顶部](#go-live)**

Made with ❤️ by [LessUp](https://github.com/LessUp)

</div>
