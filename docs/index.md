---
layout: default
title: Go-Live Documentation
description: Lightweight WebRTC SFU Server - WHIP/WHEP streaming, recording, and observability with Go and Pion WebRTC
nav_order: 0
---

# Go-Live

<div class="hero-section">
  <p class="hero-tagline">轻量级 · 高性能 · 生产就绪</p>
  <p class="hero-tagline-en">Lightweight · High Performance · Production Ready</p>
</div>

**Go-Live** 是基于 Go 和 [Pion WebRTC](https://github.com/pion/webrtc) 构建的轻量级 WebRTC SFU（选择性转发单元）服务器。专为实时音视频流分发而设计，支持 WHIP/WHEP 协议、房间广播、录制存储和完整的可观测性体系。

<div class="language-selector">
  <a href="{{ site.baseurl }}/en/" class="lang-card">
    <div class="lang-icon">🇺🇸</div>
    <h2>English</h2>
    <p>Lightweight WebRTC SFU Server Documentation</p>
    <ul>
      <li>Quick Start Guide</li>
      <li>Deployment Instructions</li>
      <li>API Reference</li>
      <li>Architecture Design</li>
    </ul>
    <span class="lang-button">Enter Documentation →</span>
  </a>

  <a href="{{ site.baseurl }}/zh/" class="lang-card">
    <div class="lang-icon">🇨🇳</div>
    <h2>中文</h2>
    <p>轻量级 WebRTC SFU 服务器文档</p>
    <ul>
      <li>快速开始指南</li>
      <li>部署说明</li>
      <li>API 参考</li>
      <li>架构设计</li>
    </ul>
    <span class="lang-button">进入文档 →</span>
  </a>
</div>

---

## 🚀 30秒快速开始 / Quick Start in 30s

```bash
# 克隆仓库 / Clone repository
git clone https://github.com/LessUp/go-live.git && cd go-live

# 启动服务 / Start server
go run ./cmd/server

# 访问 http://localhost:8080 即可开始推流
```

---

## ✨ 核心能力 / Core Capabilities

<div class="features-grid">
  <div class="feature-card">
    <div class="feature-icon">📡</div>
    <h3>WHIP/WHEP Support</h3>
    <p>Full support for WHIP publishing and WHEP playback protocols</p>
  </div>
  <div class="feature-card">
    <div class="feature-icon">🏠</div>
    <h3>Room-based Broadcast</h3>
    <p>One publisher, multiple subscribers per room with efficient fanout</p>
  </div>
  <div class="feature-card">
    <div class="feature-icon">🎥</div>
    <h3>Recording & Upload</h3>
    <p>Built-in recording with automatic S3/MinIO upload support</p>
  </div>
  <div class="feature-card">
    <div class="feature-icon">📊</div>
    <h3>Observability</h3>
    <p>Prometheus metrics and OpenTelemetry tracing integration</p>
  </div>
  <div class="feature-card">
    <div class="feature-icon">🔒</div>
    <h3>Authentication</h3>
    <p>Token-based and JWT authentication with per-room access control</p>
  </div>
  <div class="feature-card">
    <div class="feature-icon">⚡</div>
    <h3>High Performance</h3>
    <p>Built with Go for low-latency, high-throughput media distribution</p>
  </div>
</div>

---

## 📊 性能指标 / Performance Metrics

| 指标 Metric | 数值 Value |
|------------|-----------|
| 延迟 Latency | < 100ms (同区域) |
| 并发订阅者 | 1000+ / Room |
| 内存占用 | < 50MB (空载) |
| CPU 效率 | 单核支持 500+ 并发 |

---

## 🔗 快速链接 / Quick Links

| 资源 Resource | 链接 Link |
|--------------|-----------|
| GitHub 仓库 | [github.com/LessUp/go-live](https://github.com/LessUp/go-live) |
| 版本发布 Releases | [GitHub Releases](https://github.com/LessUp/go-live/releases) |
| 问题反馈 Issues | [GitHub Issues](https://github.com/LessUp/go-live/issues) |
| 更新日志 Changelog | [CHANGELOG.md](https://github.com/LessUp/go-live/blob/master/CHANGELOG.md) |
