---
layout: home
hero:
  name: Go-Live
  text: WebRTC SFU 服务器
  tagline: 轻量级、高性能的 WHIP/WHEP 流媒体服务
  actions:
    - theme: brand
      text: 快速开始
      link: /zh/getting-started
    - theme: alt
      text: 系统架构
      link: /zh/architecture/overview
    - theme: alt
      text: GitHub
      link: https://github.com/LessUp/go-live
features:
  - icon: 📡
    title: WHIP/WHEP 协议
    details: 完整支持 WHIP 推流和 WHEP 播放，兼容 OBS Studio 和现代浏览器。
  - icon: 🏠
    title: 房间 SFU 架构
    details: 单发布者、多订阅者的 SFU 架构。高效的 RTP 转发，无需转码。
  - icon: 🎥
    title: 录制与上传
    details: 内置录制功能，VP8/VP9 转 IVF，Opus 转 OGG。支持自动上传到 S3/MinIO。
  - icon: 🔒
    title: 多层认证
    details: Token 认证（全局/房间级），JWT 角色权限控制。
  - icon: 📊
    title: 可观测性
    details: Prometheus 指标、OpenTelemetry 分布式追踪、健康检查。
  - icon: ⚡
    title: 高性能
    details: Go 语言实现的低延迟、高吞吐量媒体分发。高效内存管理。
---
