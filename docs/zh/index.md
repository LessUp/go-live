---
layout: default
title: 首页
description: 轻量级 WebRTC SFU 服务器 - 基于 Go 和 Pion WebRTC 构建
nav_order: 1
lang: zh
---

# live-webrtc-go

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/LessUp/go-live/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/LessUp/go-live)](https://goreportcard.com/report/github.com/LessUp/go-live)
![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)
[![Release](https://img.shields.io/github/v/release/LessUp/go-live)](https://github.com/LessUp/go-live/releases)

**中文** | [English]({{ site.baseurl }}/en/)

基于 Go 和 [Pion WebRTC](https://github.com/pion/webrtc) 构建的轻量级、高性能 **WebRTC SFU**（选择性转发单元）服务器。支持 WHIP/WHEP 协议推流、房间广播、录制功能和完整可观测性。

---

## 🚀 快速开始

```bash
# 克隆代码仓库
git clone https://github.com/LessUp/go-live.git
cd go-live

# 直接运行
go run ./cmd/server

# 或使用开发脚本
./scripts/start.sh
```

访问服务器：`http://localhost:8080`

---

## ✨ 核心特性

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
    │  推流端   │ ──WHIP──▶│  │  认证   │───▶│    SFU      │    │
    │(OBS/浏览器)│        │  │ 中间件   │    │  管理器      │    │
    └──────────┘         │  └─────────┘    └──────┬──────┘    │
                         │                        │           │
    ┌──────────┐         │        ┌───────────────┤           │
    │  播放端   │ ◀──WHEP──│        │               │           │
    │ (浏览器)  │───────▶  │        ▼               ▼           │
    └──────────┘         │ ┌─────────────┐ ┌───────────┐      │
                         │ │    房间     │ │   录制    │      │
                         │ │  (转发)     │ │  与上传   │      │
                         │ └──────┬──────┘ └─────┬─────┘      │
                         └────────┼──────────────┼────────────┘
                                  │              │
                        ┌─────────▼──────────────▼───────────┐
                        │          对象存储                   │
                        │         (S3/MinIO)                 │
                        └────────────────────────────────────┘
```

---

## 📖 文档导航

| 文档 | 内容 |
|------|------|
| [使用指南](usage.md) | 本地开发、Docker 部署、配置说明、故障排除 |
| [设计说明](design.md) | 系统架构、模块拆分、数据流向 |
| [API 参考](api.md) | 完整 HTTP API 文档、请求/响应格式、错误码 |

---

## 🛠️ 开发指南

```bash
# 构建
make build

# 运行测试
make test

# 完整 CI 流水线
make ci
```

---

## 🤝 参与贡献

我们欢迎各种形式的贡献！请查看我们的[贡献指南](https://github.com/LessUp/go-live/blob/master/CONTRIBUTING.md)了解详情。

---

## 📄 许可协议

本项目基于 [MIT License](https://github.com/LessUp/go-live/blob/master/LICENSE) 开源。

---

## 🔗 相关链接

- [GitHub 仓库](https://github.com/LessUp/go-live)
- [问题反馈](https://github.com/LessUp/go-live/issues)
- [更新日志](https://github.com/LessUp/go-live/blob/master/CHANGELOG.md)
- [Pion WebRTC](https://github.com/pion/webrtc)
