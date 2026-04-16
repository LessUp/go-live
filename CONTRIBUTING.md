# Contributing to live-webrtc-go

Thank you for your interest in contributing! This document provides guidelines for contributing to the project.

## Quick Links

- [Code of Conduct](CODE_OF_CONDUCT.md)
- [Security Policy](SECURITY.md)
- [Changelog](CHANGELOG.md)
- [GitHub Issues](https://github.com/LessUp/go-live/issues)

## Table of Contents

- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Development Workflow](#development-workflow)
- [Code Standards](#code-standards)
- [Testing Guidelines](#testing-guidelines)
- [Pull Request Process](#pull-request-process)
- [Commit Guidelines](#commit-guidelines)

## Development Setup

### Prerequisites

| Tool | Version | Purpose |
|------|---------|---------|
| Go | 1.22+ | Build and run |
| Make | Any | Build automation |
| Docker | 20.10+ | Container testing |
| golangci-lint | Latest | Linting |

### Quick Start

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/go-live.git
cd go-live

# Download dependencies
go mod download

# Install tools (optional)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/securego/gosec/v2/cmd/gosec@latest

# Verify setup
make test
```

### IDE Setup

**VS Code** (recommended):
```json
// .vscode/settings.json
{
  "go.formatTool": "gofmt",
  "go.lintTool": "golangci-lint",
  "go.lintOnSave": "workspace",
  "editor.formatOnSave": true
}
```

**GoLand**: Enable `gofmt` and `golangci-lint` in preferences.

## Project Structure

```
├── cmd/server/           # Application entry point
│   ├── main.go           # HTTP server initialization
│   └── web/              # Embedded static files
├── internal/
│   ├── api/              # HTTP handlers and routing
│   │   ├── handlers.go   # WHIP/WHEP/Rooms handlers
│   │   ├── middleware.go # Auth, CORS, rate limiting
│   │   └── routes.go     # URL routing
│   ├── config/           # Configuration management
│   ├── sfu/              # WebRTC SFU core
│   │   ├── manager.go    # Room lifecycle
│   │   ├── room.go       # PeerConnection, tracks
│   │   └── track.go      # RTP distribution
│   ├── metrics/          # Prometheus metrics
│   ├── otel/             # OpenTelemetry tracing
│   ├── uploader/         # S3/MinIO upload
│   └── testutil/         # Test utilities
├── test/
│   ├── integration/      # Integration tests
│   ├── e2e/              # End-to-end tests
│   ├── security/         # Security tests
│   └── performance/      # Benchmarks
├── docs/                 # Documentation (GitHub Pages)
├── web/                  # Frontend source files
├── scripts/              # Development scripts
└── Makefile              # Build automation
```

## Development Workflow

### 1. Create a Branch

```bash
git checkout master && git pull
git checkout -b feature/your-feature
```

**Branch naming**:
- `feature/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation
- `refactor/` - Code refactoring
- `test/` - Test updates

### 2. Make Changes

- Keep changes focused
- Follow [Code Standards](#code-standards)
- Add/update tests
- Update documentation

### 3. Verify Locally

```bash
make fmt       # Format code
make lint      # Run linters
make security  # Security scan
make test      # Run tests
make ci        # Full CI pipeline
```

### 4. Commit & Push

```bash
git add .
git commit -m "feat: add new feature"
git push origin feature/your-feature
```

## Code Standards

### Go Style

- Follow [Effective Go](https://golang.org/doc/effective_go)
- Use `gofmt -s` for formatting
- Run `go vet` before committing
- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

### Code Organization

```go
// Package comment
package example

import (
    // Standard library
    "context"
    
    // Third-party
    "github.com/pkg/errors"
    
    // Local packages
    "project/internal/config"
)

// Constants
const maxRetries = 3

// Types
type Service struct { ... }

// Constructor
func NewService(cfg *config.Config) *Service { ... }

// Methods (grouped by receiver)
func (s *Service) Start() error { ... }
func (s *Service) Stop() error { ... }

// Functions
func helper() { ... }
```

### Error Handling

```go
// ✅ Good: Wrap with context
if err := doSomething(); err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}

// ✅ Good: Log with context
slog.Error("operation failed", "error", err, "room", roomID)

// ❌ Bad: Silent ignore
_ = doSomething()

// ❌ Bad: Generic error
return errors.New("error")
```

### Concurrency

```go
// ✅ Good: Proper cleanup
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()

// ✅ Good: Protected access
type SafeMap struct {
    mu sync.RWMutex
    m  map[string]string
}

// ❌ Bad: Goroutine leak
go func() {
    for {
        // No exit condition!
    }
}()
```

### Comments

```go
// Exported functions need documentation
// Start begins the service and returns immediately.
func (s *Service) Start() error { ... }

// Unexported functions can have shorter comments
// validateRoom checks room name validity.
func validateRoom(name string) bool { ... }
```

## Testing Guidelines

### Test Structure

```go
func TestSomething(t *testing.T) {
    t.Parallel()  // When independent
    
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"valid", "input", "output", false},
        {"invalid", "", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Process(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("Process() error = %v, wantErr %v", err, tt.wantErr)
            }
            if got != tt.want {
                t.Errorf("Process() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Test Commands

| Command | Description |
|---------|-------------|
| `make test` | Unit + integration + security |
| `make test-unit` | Unit tests only |
| `make test-integration` | Integration tests |
| `make test-e2e` | End-to-end tests |
| `make coverage` | Generate coverage report |

### Coverage Target

- Core packages: ≥80%
- New code: Must include tests

## Pull Request Process

### Before Submitting

- [ ] `make fmt` - Code formatted
- [ ] `make lint` - No lint errors
- [ ] `make security` - No security issues
- [ ] `make test` - All tests pass
- [ ] Documentation updated

### PR Template

```markdown
## Summary
Brief description of changes.

## Motivation
Why this change is needed.

## Changes
- Change 1
- Change 2

## Testing
- [ ] Unit tests added/updated
- [ ] Manual testing performed

## Checklist
- [ ] `make test` passes
- [ ] `make lint` passes
- [ ] Documentation updated
```

### Review Process

1. Automated checks must pass
2. At least one approval required
3. All conversations resolved
4. Squash and merge to `master`

## Commit Guidelines

### Format

```
<type>: <subject>

<body (optional)>

<footer (optional)>
```

### Types

| Type | When to Use |
|------|-------------|
| `feat` | New feature |
| `fix` | Bug fix |
| `docs` | Documentation |
| `style` | Formatting (no code change) |
| `refactor` | Code refactoring |
| `test` | Adding tests |
| `chore` | Maintenance tasks |

### Examples

```
feat: add WebSocket support for room events

- Add /ws endpoint for real-time subscriptions
- Support join/leave/publish/unpublish events
- Include connection lifecycle management

Closes #42
```

```
fix: resolve memory leak in track fanout

The readLoop was not properly closing when subscribers
disconnected, causing goroutine leaks.

Fixes #44
```

## Getting Help

- [Open an issue](https://github.com/LessUp/go-live/issues) for bugs/features
- [Start a discussion](https://github.com/LessUp/go-live/discussions) for questions
- Check existing issues before creating new ones

## Recognition

Contributors are recognized in release notes. Thank you for improving live-webrtc-go!
