# Changelog 编写指南

本指南详细说明如何为 live-webrtc-go 项目编写高质量的 Changelog 条目。

---

## 目录

- [Changelog 格式](#changelog-格式)
- [分类定义](#分类定义)
- [编写规范](#编写规范)
- [示例](#示例)
- [术语表](#术语表)

---

## Changelog 格式

本项目遵循 [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) 规范，并采用 [Semantic Versioning](https://semver.org/) 进行版本管理。

### 基本结构

```markdown
## [版本号] - YYYY-MM-DD

### Added
- 新功能描述

### Changed
- 变更描述

### Deprecated
- 弃用功能描述

### Removed
- 移除功能描述

### Fixed
- 修复描述

### Security
- 安全修复描述
```

---

## 分类定义

| 分类 | 用途 | 示例 |
|------|------|------|
| **Added** | 新增功能、API、配置选项 | 添加 WebSocket 支持、新增指标 |
| **Changed** | 现有功能的变更、改进 | 优化性能、改进错误处理 |
| **Deprecated** | 即将移除的功能 | 标记废弃的 API 端点 |
| **Removed** | 已移除的功能 | 删除旧的配置选项 |
| **Fixed** | Bug 修复 | 修复内存泄漏、修复竞态条件 |
| **Security** | 安全修复 | 修复漏洞、更新依赖 |

---

## 编写规范

### 1. 使用祈使语气

使用命令式语气，如 "Add" 而非 "Added"、"Fix" 而非 "Fixed"。

```markdown
✅ Add WebSocket support for real-time events
❌ Added WebSocket support for real-time events
```

### 2. 首字母大写

每条条目首字母大写。

```markdown
✅ Add new feature
❌ add new feature
```

### 3. 结尾不加句号

条目结尾不加句号（除非是完整句子）。

```markdown
✅ Add WebSocket support
❌ Add WebSocket support.
```

### 4. 引用相关 Issue/PR

在条目中引用相关的 Issue 或 PR 编号。

```markdown
✅ Fix memory leak in track fanout (#123)
```

### 5. 分组相关变更

使用子标题对相关变更进行分组。

```markdown
### Added

#### Observability
- Add OpenTelemetry tracing support (#50)
- Add distributed trace context propagation (#52)

#### API
- Add `/api/admin/rooms/{room}/stats` endpoint (#51)
```

### 6.  Breaking Changes 标记

对于破坏性变更，在条目前加 `**Breaking:**` 标记。

```markdown
### Changed
- **Breaking:** Change API response format for `/api/rooms` endpoint (#60)
  - Old: `{"rooms": [...]}`
  - New: `[...]`
  - Migration: Update client code to handle array directly
```

---

## 示例

### 示例 1: 新功能发布

```markdown
## [1.2.0] - 2025-04-15

### Added

#### Features
- Add WebSocket support for real-time room events (#100)
- Add support for H.264 codec in addition to VP8/VP9 (#98)

#### API
- Add `/api/admin/rooms/{room}/stats` endpoint for detailed room metrics (#99)

#### Observability
- Add Grafana dashboard templates (#97)
- Add OpenTelemetry metrics exporter (#96)

### Changed

#### Performance
- Improve track fanout performance by 30% using sync.Pool (#95)
- Reduce memory allocation in RTP processing (#94)

### Fixed

#### Bugs
- Fix race condition in subscriber cleanup (#93)
- Fix memory leak when publisher disconnects unexpectedly (#92)
- Fix CORS preflight handling for Safari browsers (#91)

### Security
- Update golang.org/x/crypto to fix CVE-2025-XXXX (#90)
```

### 示例 2: Bug 修复版本

```markdown
## [1.1.2] - 2025-03-28

### Fixed
- Fix ICE connection failure on certain NAT configurations (#85)
- Fix recording file corruption on unexpected shutdown (#84)
- Fix token comparison timing attack vulnerability (#83)
```

### 示例 3: 破坏性变更

```markdown
## [2.0.0] - 2025-06-01

### Added
- Add pluggable authentication provider interface (#150)

### Changed
- **Breaking:** Rename `AUTH_TOKEN` to `API_TOKEN` (#149)
  - Migration: Update environment variable name
- **Breaking:** Change `/api/rooms` response format (#148)
  - Old: `{"rooms": [{"name": "room1"}]}`
  - New: `[{"name": "room1"}]`
  - Migration: Update client to handle array response

### Removed
- **Breaking:** Remove deprecated `webrtc.UseMDNS` configuration option (#147)
```

---

## 术语表

| 术语 | 说明 |
|------|------|
| **SFU** | Selective Forwarding Unit，选择性转发单元 |
| **WHIP** | WebRTC-HTTP Ingestion Protocol，WebRTC HTTP 摄入协议 |
| **WHEP** | WebRTC-HTTP Egress Protocol，WebRTC HTTP 流出协议 |
| **RTP** | Real-time Transport Protocol，实时传输协议 |
| **ICE** | Interactive Connectivity Establishment，交互式连接建立 |
| **SDP** | Session Description Protocol，会话描述协议 |

---

## 相关链接

- [Keep a Changelog](https://keepachangelog.com/)
- [Semantic Versioning](https://semver.org/)
- [CHANGELOG.md](/CHANGELOG.md)
