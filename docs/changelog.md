---
layout: default
title: Changelog
description: Release notes and version history for Go-Live
nav_order: 5
---

# Changelog

All notable changes to Go-Live are documented in this file.

---

## [v1.1.0](https://github.com/LessUp/go-live/releases/tag/v1.1.0) - 2026-04-15

### Added
- WHIP/WHEP protocol support for streaming and playback
- Room-based SFU architecture with publisher/subscriber model
- Recording functionality (VP8/VP9/Opus codecs)
- S3/MinIO upload integration for recording files
- Prometheus metrics integration (`live_rooms`, `live_subscribers`, `live_rtp_bytes_total`)
- OpenTelemetry tracing support
- Token-based authentication (global + per-room tokens)
- JWT authentication with role-based access control
- Rate limiting middleware
- CORS support with configurable allowed origins
- Frontend web UI with WHIP/WHEP client
- Bootstrap API for frontend configuration

### Security
- Auth middleware for all API endpoints
- Bearer token and X-Auth-Token support
- Room-level access control

### Documentation
- Complete API reference in English and Chinese
- Architecture design documentation
- Usage guide with deployment instructions
- Troubleshooting guide

### Infrastructure
- Docker and Docker Compose support
- GitHub Actions CI/CD pipeline
- Automated testing (unit, integration, security, e2e)
- Code coverage reporting with Codecov

---

## Initial Release

- Core SFU functionality with Pion WebRTC
- Basic room management
- Health check endpoint (`/healthz`)

---

## Version History

| Version | Date | Description |
|---------|------|-------------|
| v1.1.0 | 2026-04-15 | Initial public release with full feature set |

---

For more details, see the [GitHub Releases](https://github.com/LessUp/go-live/releases) page.
