---
layout: default
title: Go-Live 首页
description: 轻量级 WebRTC SFU 服务器 - 基于 Go 和 Pion WebRTC 构建
nav_order: 1
lang: zh
---

# Go-Live

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

# 或者使用开发脚本（推荐）
./scripts/start.sh
```

服务启动后，访问 `http://localhost:8080` 查看前端界面。

---

## ✨ 核心特性

<div class="features-grid">
  <div class="feature-card">
    <div class="feature-icon">📡</div>
    <h3>WHIP/WHEP 协议</h3>
    <p>完整支持 WHIP 推流和 WHEP 播放协议，兼容 OBS、浏览器等主流工具</p>
  </div>
  <div class="feature-card">
    <div class="feature-icon">🏠</div>
    <h3>房间广播</h3>
    <p>基于房间的 SFU 架构，一个发布者，多个订阅者，高效 RTP 分发</p>
  </div>
  <div class="feature-card">
    <div class="feature-icon">🎥</div>
    <h3>录制与上传</h3>
    <p>内置录制功能，支持自动上传到 S3/MinIO 对象存储</p>
  </div>
  <div class="feature-card">
    <div class="feature-icon">📊</div>
    <h3>可观测性</h3>
    <p>Prometheus 指标采集和 OpenTelemetry 分布式追踪集成</p>
  </div>
  <div class="feature-card">
    <div class="feature-icon">🔒</div>
    <h3>认证体系</h3>
    <p>支持 Token 和 JWT 认证，可配置房间级访问控制</p>
  </div>
  <div class="feature-card">
    <div class="feature-icon">⚡</div>
    <h3>高性能</h3>
    <p>基于 Go 语言开发，低延迟、高并发的媒体流分发能力</p>
  </div>
</div>

---

## 📖 文档导航

| 文档 | 描述 | 链接 |
|------|------|------|
| 使用指南 | 部署、配置、API 示例、故障排除 | [使用指南]({{ site.baseurl }}/zh/usage.html) |
| API 参考 | 完整的 REST API 文档 | [API 参考]({{ site.baseurl }}/zh/api.html) |
| 架构设计 | 系统架构、模块说明、数据流 | [设计文档]({{ site.baseurl }}/zh/design.html) |

---

## 🛠️ 快速部署

### Docker 部署

```bash
# 构建镜像
docker build -t live-webrtc-go:latest .

# 运行
docker run --rm -p 8080:8080 live-webrtc-go:latest

# 启用录制
docker run --rm -p 8080:8080 \
  -e RECORD_ENABLED=1 \
  -e RECORD_DIR=/records \
  -v $(pwd)/records:/records \
  live-webrtc-go:latest
```

### Docker Compose

```yaml
version: '3.8'
services:
  live-webrtc:
    build: .
    ports:
      - "8080:8080"
    environment:
      - RECORD_ENABLED=1
    volumes:
      - ./records:/records
    restart: unless-stopped
```

---

## 📦 系统要求

| 依赖 | 版本 | 说明 |
|------|------|------|
| Go | 1.22+ | 编译所需 |
| Docker | 20.10+ | 容器部署可选 |
| 浏览器 | Chrome 90+ / Firefox 88+ | WebRTC 支持 |

---

## 🔗 相关链接

- **GitHub**: [https://github.com/LessUp/go-live](https://github.com/LessUp/go-live)
- **Releases**: [https://github.com/LessUp/go-live/releases](https://github.com/LessUp/go-live/releases)
- **Issues**: [https://github.com/LessUp/go-live/issues](https://github.com/LessUp/go-live/issues)

<style>
.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
  margin: 2rem 0;
}

.feature-card {
  padding: 1.5rem;
  border: 1px solid #e1e4e8;
  border-radius: 8px;
  background: #f6f8fa;
  transition: all 0.2s ease;
}

.feature-card:hover {
  border-color: #0969da;
  box-shadow: 0 2px 8px rgba(9, 105, 218, 0.1);
  transform: translateY(-2px);
}

.feature-icon {
  font-size: 2rem;
  margin-bottom: 0.5rem;
}

.feature-card h3 {
  margin: 0 0 0.5rem 0;
  font-size: 1.1rem;
  color: #24292f;
}

.feature-card p {
  margin: 0;
  color: #57606a;
  font-size: 0.9rem;
  line-height: 1.5;
}
</style>
