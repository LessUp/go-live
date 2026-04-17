# Changelog Writing Guide

This guide explains how to write high-quality changelog entries for the live-webrtc-go project.

---

## Table of Contents

- [Changelog Format](#changelog-format)
- [Category Definitions](#category-definitions)
- [Writing Guidelines](#writing-guidelines)
- [Examples](#examples)
- [Glossary](#glossary)

---

## Changelog Format

This project follows the [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) specification and uses [Semantic Versioning](https://semver.org/) for version management.

### Basic Structure

```markdown
## [Version] - YYYY-MM-DD

### Added
- Feature description

### Changed
- Change description

### Deprecated
- Deprecated feature description

### Removed
- Removed feature description

### Fixed
- Fix description

### Security
- Security fix description
```

---

## Category Definitions

| Category | Purpose | Example |
|----------|---------|---------|
| **Added** | New features, APIs, config options | Add WebSocket support, new metrics |
| **Changed** | Existing feature changes, improvements | Optimize performance, improve error handling |
| **Deprecated** | Soon-to-be-removed features | Mark deprecated API endpoint |
| **Removed** | Removed features | Delete old config option |
| **Fixed** | Bug fixes | Fix memory leak, fix race condition |
| **Security** | Security fixes | Fix vulnerability, update dependencies |

---

## Writing Guidelines

### 1. Use Imperative Mood

Use imperative mood, e.g., "Add" not "Added", "Fix" not "Fixed".

```markdown
✅ Add WebSocket support for real-time events
❌ Added WebSocket support for real-time events
```

### 2. Capitalize First Letter

Capitalize the first letter of each entry.

```markdown
✅ Add new feature
❌ add new feature
```

### 3. No Trailing Period

Do not add a period at the end of entries (unless it's a complete sentence).

```markdown
✅ Add WebSocket support
❌ Add WebSocket support.
```

### 4. Reference Related Issues/PRs

Reference related issue or PR numbers in entries.

```markdown
✅ Fix memory leak in track fanout (#123)
```

### 5. Group Related Changes

Use subheadings to group related changes together.

```markdown
### Added

#### Observability
- Add OpenTelemetry tracing support (#50)
- Add distributed trace context propagation (#52)

#### API
- Add `/api/admin/rooms/{room}/stats` endpoint (#51)
```

### 6. Mark Breaking Changes

For breaking changes, prefix entries with `**Breaking:**`.

```markdown
### Changed
- **Breaking:** Change API response format for `/api/rooms` endpoint (#60)
  - Old: `{"rooms": [...]}`
  - New: `[...]`
  - Migration: Update client code to handle array directly
```

---

## Examples

### Example 1: Feature Release

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

### Example 2: Bug Fix Release

```markdown
## [1.1.2] - 2025-03-28

### Fixed
- Fix ICE connection failure on certain NAT configurations (#85)
- Fix recording file corruption on unexpected shutdown (#84)
- Fix token comparison timing attack vulnerability (#83)
```

### Example 3: Breaking Changes

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

## Glossary

| Term | Description |
|------|-------------|
| **SFU** | Selective Forwarding Unit |
| **WHIP** | WebRTC-HTTP Ingestion Protocol |
| **WHEP** | WebRTC-HTTP Egress Protocol |
| **RTP** | Real-time Transport Protocol |
| **ICE** | Interactive Connectivity Establishment |
| **SDP** | Session Description Protocol |

---

## Related Links

- [Keep a Changelog](https://keepachangelog.com/)
- [Semantic Versioning](https://semver.org/)
- [CHANGELOG.md](/CHANGELOG.md)
