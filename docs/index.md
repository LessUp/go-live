---
layout: default
title: 首页
---

# live-webrtc-go

使用 Go + [Pion WebRTC](https://github.com/pion/webrtc) 构建的轻量级在线直播服务示例。实现了 WHIP 推流、WHEP 播放、嵌入式 Web 页面、可配置鉴权与房间状态查询，可作为开源参考或项目脚手架。

## 功能特点

- **WebRTC SFU**：基于 Pion 实现最小可用的房间转发逻辑，支持多观众同时观看。
- **WHIP / WHEP 接口**：HTTP API 兼容现代浏览器或 OBS WHIP 插件推流与播放。
- **可选鉴权**：支持 Token / JWT 三层优先级鉴权体系。
- **房间状态查询**：`GET /api/rooms` 返回在线房间、发布者与订阅者统计。
- **录制能力**：可选将 VP8/VP9 保存为 IVF、Opus 保存为 OGG，并支持 S3/MinIO 上传。
- **监控指标**：`GET /metrics` 暴露 Prometheus 指标（RTP 字节/包、订阅者数、房间数）。
- **容器化**：提供 Dockerfile 与 docker-compose.yml，支持挂载录制目录与配置环境变量。

## 快速开始

```bash
git clone https://github.com/LessUp/go-live.git
cd go-live
go mod tidy
go run ./cmd/server
```

启动后访问：

| 页面 | 地址 |
|------|------|
| 推流页 | `http://localhost:8080/web/publisher.html` |
| 播放页 | `http://localhost:8080/web/player.html` |
| 房间列表 | `http://localhost:8080/api/rooms` |
| 健康检查 | `http://localhost:8080/healthz` |
| 指标监控 | `http://localhost:8080/metrics` |

## HTTP API

| 方法 | 路径 | 说明 |
|------|------|------|
| `POST` | `/api/whip/publish/{room}` | WHIP 推流 |
| `POST` | `/api/whep/play/{room}` | WHEP 播放 |
| `GET` | `/api/rooms` | 房间列表与在线状态 |
| `GET` | `/api/records` | 录制文件列表 |
| `POST` | `/api/admin/rooms/{room}/close` | 关闭指定房间 |
| `GET` | `/healthz` | 健康检查 |
| `GET` | `/metrics` | Prometheus 指标 |

## 项目结构

```
├── cmd/server          # 入口程序（HTTP + WebRTC）
│   └── web             # 嵌入式静态页面
├── internal/api        # HTTP handlers (WHIP/WHEP/Rooms)
├── internal/config     # 配置加载
├── internal/metrics    # Prometheus 指标
├── internal/sfu        # WebRTC SFU 管理逻辑
├── internal/uploader   # S3/MinIO 上传
├── docs/               # 项目文档（本站）
└── test/               # 测试套件
```

## 文档导航

- [**使用指南**](usage) — 完整的启动、部署、API 示例与排障。
- [**设计说明**](design) — 架构背景、模块拆分与数据流。

## 许可协议

本项目以 MIT 协议开源。欢迎 Issue 与 Pull Request！
