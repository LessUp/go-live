---
layout: home
hero:
  name: Go-Live
  text: WebRTC SFU Server
  tagline: Lightweight, high-performance streaming with WHIP/WHEP protocols
  actions:
    - theme: brand
      text: Getting Started
      link: /en/getting-started
    - theme: alt
      text: Architecture
      link: /en/architecture/overview
    - theme: alt
      text: GitHub
      link: https://github.com/LessUp/go-live
features:
  - icon: 📡
    title: WHIP/WHEP Protocols
    details: Full support for WHIP publishing and WHEP playback, compatible with OBS Studio and modern browsers.
  - icon: 🏠
    title: Room-based SFU
    details: SFU architecture with one publisher and multiple subscribers per room. Efficient RTP forwarding without transcoding.
  - icon: 🎥
    title: Recording & Upload
    details: Built-in recording with VP8/VP9 to IVF, Opus to OGG. Automatic S3/MinIO upload support.
  - icon: 🔒
    title: Multi-layer Auth
    details: Token authentication (global/per-room), JWT with role-based access control.
  - icon: 📊
    title: Observability
    details: Prometheus metrics, OpenTelemetry distributed tracing, health checks.
  - icon: ⚡
    title: High Performance
    details: Low-latency, high-throughput media distribution in Go. Efficient memory management.
---
