# live-webrtc-go

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
- [快速开始](#-快速开始)
- [安装](#-安装)
- [配置](#-配置)
- [API 参考](#-api-参考)
- [文档](#-文档)
- [开发](#-开发)
- [Docker 部署](#-docker-部署)
- [参与贡献](#-参与贡献)
- [许可协议](#-许可协议)

---

## ✨ 特性

| 特性 | 说明 |
|------|------|
| 🎥 **WHIP/WHEP 协议** | 标准 HTTP 协议的 WebRTC 推流和播放，兼容 OBS 和现代浏览器 |
| 🏠 **房间 SFU 架构** | 每房间单发布者、多订阅者，高效的媒体转发 |
| 🔐 **灵活认证体系** | 支持全局 Token、房间级 Token、JWT 角色认证 |
| 📹 **录制与上传** | VP8/VP9 → IVF、Opus → OGG，支持 S3/MinIO 自动上传 |
| 📊 **完整可观测性** | Prometheus 指标、OpenTelemetry 追踪、健康检查端点 |
| 🐳 **云原生就绪** | 支持 Docker、Docker Compose，适配 Kubernetes |
| 🌐 **嵌入式 Web 界面** | 内置推流和播放页面，开箱即用 |

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

---

## 🚀 快速开始

### 前置要求

- Go 1.22+（源码构建）
- 或 Docker（容器部署）

### 从源码运行

```bash
# 克隆代码仓库
git clone https://github.com/LessUp/go-live.git
cd go-live

# 直接运行
go run ./cmd/server

# 或使用开发脚本
./scripts/start.sh
```

### 使用 Docker 运行

```bash
docker run --rm -p 8080:8080 ghcr.io/lessup/go-live:latest
```

### 访问应用

| 页面 | 地址 | 说明 |
|------|------|------|
| 🏠 首页 | http://localhost:8080/ | 重定向到推流页 |
| 📤 推流页 | http://localhost:8080/web/publisher.html | 浏览器推流 |
| 📥 播放页 | http://localhost:8080/web/player.html | 观看直播 |
| 📋 录制列表 | http://localhost:8080/web/records.html | 录制文件浏览 |
| 📊 指标监控 | http://localhost:8080/metrics | Prometheus 指标 |

---

## 📦 安装

### 二进制安装

从 [GitHub Releases](https://github.com/LessUp/go-live/releases) 下载预编译二进制文件：

```bash
# Linux AMD64
curl -LO https://github.com/LessUp/go-live/releases/latest/download/live-webrtc-go-linux-amd64
chmod +x live-webrtc-go-linux-amd64
./live-webrtc-go-linux-amd64
```

### 从源码构建

```bash
# 克隆并构建
git clone https://github.com/LessUp/go-live.git
cd go-live
make build

# 二进制位于 bin/server
./bin/server
```

---

## ⚙️ 配置

通过环境变量进行配置：

### 核心配置

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `HTTP_ADDR` | `:8080` | HTTP 监听地址 |
| `ALLOWED_ORIGIN` | `*` | CORS 允许来源 |

### 身份认证

| 变量 | 说明 |
|------|------|
| `AUTH_TOKEN` | 全局认证令牌 |
| `ROOM_TOKENS` | 房间级令牌，格式：`room1:tok1;room2:tok2` |
| `JWT_SECRET` | JWT HMAC 签名密钥 |
| `ADMIN_TOKEN` | 管理接口令牌 |

### WebRTC/ICE

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `STUN_URLS` | `stun:stun.l.google.com:19302` | STUN 服务器 |
| `TURN_URLS` | - | TURN 服务器 |
| `TURN_USERNAME` | - | TURN 用户名 |
| `TURN_PASSWORD` | - | TURN 密码 |

### 录制功能

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `RECORD_ENABLED` | `0` | 启用录制（设为 `1` 启用） |
| `RECORD_DIR` | `records` | 录制输出目录 |

查看[完整配置指南](https://lessup.github.io/go-live/zh/usage.html#配置说明)了解所有选项。

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
| `GET` | `/healthz` | 健康检查 |
| `GET` | `/metrics` | Prometheus 指标 |

查看[完整 API 文档](https://lessup.github.io/go-live/zh/api.html)了解详情。

---

## 📚 文档

| 文档 | 说明 |
|------|------|
| [English Docs](https://lessup.github.io/go-live/en/) | 英文完整文档 |
| [中文文档](https://lessup.github.io/go-live/zh/) | 中文完整文档 |
| [使用指南](https://lessup.github.io/go-live/zh/usage.html) | 本地开发、Docker 部署、故障排除 |
| [设计说明](https://lessup.github.io/go-live/zh/design.html) | 系统架构和模块详解 |
| [API 参考](https://lessup.github.io/go-live/zh/api.html) | 完整 HTTP API 文档 |

---

## 🛠️ 开发

### Makefile 命令

```bash
make build      # 编译二进制到 bin/
make test       # 运行所有测试
make lint       # 运行代码检查
make security   # 运行安全扫描
make coverage   # 生成覆盖率报告
make ci         # 完整 CI 流水线
```

### 运行测试

```bash
# 单元测试
make test-unit

# 集成测试
make test-integration

# 所有测试
make test-all
```

---

## 🐳 Docker 部署

### 构建镜像

```bash
docker build -t live-webrtc-go:latest .
```

### Docker Compose

```bash
docker compose up -d
```

### 完整配置示例

```bash
docker run --rm -p 8080:8080 \
  -e AUTH_TOKEN=mysecret \
  -e RECORD_ENABLED=1 \
  -e RECORD_DIR=/records \
  -v $(pwd)/records:/records \
  live-webrtc-go:latest
```

查看 [Docker 部署指南](https://lessup.github.io/go-live/zh/usage.html#docker-部署)了解更多详情。

---

## 🤝 参与贡献

我们欢迎各种形式的贡献！请查看我们的[贡献指南](CONTRIBUTING.md)了解详情。

- [贡献指南](CONTRIBUTING.md)
- [行为准则](CODE_OF_CONDUCT.md)
- [安全策略](SECURITY.md)

---

## 📄 许可协议

本项目基于 [MIT License](LICENSE) 开源。

---

## 🔗 相关链接

- [GitHub 仓库](https://github.com/LessUp/go-live)
- [问题反馈](https://github.com/LessUp/go-live/issues)
- [版本发布](https://github.com/LessUp/go-live/releases)
- [Pion WebRTC](https://github.com/pion/webrtc)

---

<div align="center">

**[⬆ 回到顶部](#live-webrtc-go)**

Made with ❤️ by [LessUp](https://github.com/LessUp)

</div>
