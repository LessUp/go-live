---
layout: default
title: Go-Live Documentation
description: WebRTC SFU Server Documentation - 轻量级 WebRTC SFU 服务器文档
nav_order: 0
---

# Go-Live Documentation

## Select Language / 选择语言

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

## About / 关于

**Go-Live** is a lightweight, high-performance WebRTC SFU (Selective Forwarding Unit) server built with Go and Pion WebRTC. It supports WHIP/WHEP protocols for streaming, room-based broadcast, recording, and comprehensive observability.

**Go-Live** 是基于 Go 和 Pion WebRTC 构建的轻量级、高性能 WebRTC SFU（选择性转发单元）服务器。支持 WHIP/WHEP 协议推流、房间广播、录制功能和完整可观测性。

### Quick Links / 快速链接

| Resource / 资源 | Link / 链接 |
|----------------|-------------|
| GitHub Repository | [github.com/LessUp/go-live](https://github.com/LessUp/go-live) |
| Releases / 版本发布 | [GitHub Releases](https://github.com/LessUp/go-live/releases) |
| Issues / 问题反馈 | [GitHub Issues](https://github.com/LessUp/go-live/issues) |
| Changelog / 更新日志 | [CHANGELOG.md](https://github.com/LessUp/go-live/blob/master/CHANGELOG.md) |

<style>
.language-selector {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 2rem;
  margin: 2rem 0;
}

.lang-card {
  display: block;
  padding: 2rem;
  border: 2px solid #e1e4e8;
  border-radius: 12px;
  text-decoration: none;
  color: inherit;
  transition: all 0.3s ease;
}

.lang-card:hover {
  border-color: #0969da;
  box-shadow: 0 4px 12px rgba(9, 105, 218, 0.15);
  text-decoration: none;
}

.lang-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.lang-card h2 {
  margin: 0 0 0.5rem 0;
  color: #24292f;
}

.lang-card p {
  color: #57606a;
  margin-bottom: 1rem;
}

.lang-card ul {
  list-style: none;
  padding: 0;
  margin: 0 0 1.5rem 0;
}

.lang-card li {
  padding: 0.25rem 0;
  color: #57606a;
}

.lang-card li::before {
  content: "✓ ";
  color: #1a7f37;
  font-weight: bold;
}

.lang-button {
  display: inline-block;
  padding: 0.75rem 1.5rem;
  background: #0969da;
  color: white;
  border-radius: 6px;
  font-weight: 600;
  transition: background 0.2s;
}

.lang-card:hover .lang-button {
  background: #0550ae;
}
</style>
