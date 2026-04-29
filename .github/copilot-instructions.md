# GitHub Copilot Instructions

> Go-Live WebRTC SFU Server - AI Agent Context

## Project Overview

**Go-Live** is a lightweight WebRTC SFU (Selective Forwarding Unit) server built with Go 1.22+ and Pion WebRTC. It implements WHIP/WHEP protocols for real-time streaming with room-based broadcast, recording, and observability.

## Architecture Pattern

```
HTTP Request → CORS → Rate Limiter → Auth → Handler → SFU Room

Manager → Room → trackFanout (three-layer hierarchy)
- Manager: Room lifecycle management
- Room: Single publisher, multiple subscribers
- trackFanout: RTP packet distribution
```

## Code Style

### Error Handling
```go
// ✅ Always wrap errors with context
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}

// ✅ Log errors with slog
if err := conn.Close(); err != nil {
    slog.Error("failed to close connection", "error", err)
}
```

### Concurrency
```go
// ✅ Protect shared state with mutex
type Room struct {
    mu          sync.RWMutex
    subscribers map[string]*PeerConnection
}

// ✅ Use context for cancellation
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
```

## Key Files

| File | Purpose |
|------|---------|
| `internal/sfu/room.go` | SFU core logic, publisher/subscriber management |
| `internal/api/handlers.go` | HTTP handlers for WHIP/WHEP endpoints |
| `internal/api/middleware.go` | Auth, CORS, rate limiting |
| `internal/config/config.go` | Environment configuration |
| `openspec/specs/` | Requirements specifications (source of truth) |

## HTTP Handler Pattern

```go
func (h *HTTPHandlers) HandleSomething(w http.ResponseWriter, r *http.Request) {
    // 1. CORS headers
    h.allowCORS(w, r)
    
    // 2. Rate limiting
    if !h.allowRate(r) {
        http.Error(w, "rate limit exceeded", 429)
        return
    }
    
    // 3. Authentication
    if !h.authOK(r, h.config) {
        http.Error(w, "unauthorized", 401)
        return
    }
    
    // 4. Business logic
}
```

## API Status Codes

| Endpoint | Success | Error Conditions |
|----------|---------|-----------------|
| WHIP Publish | 201 Created | 409 Conflict (duplicate publisher) |
| WHEP Play | 201 Created | 404 Not Found (no publisher), 403 Forbidden (subscriber limit) |

## Naming Conventions

- Room names: `^[A-Za-z0-9_-]{1,64}$`
- Exported functions: must have documentation comments
- Error messages: lowercase, no trailing punctuation

## Security Guidelines

- Use `crypto/subtle.ConstantTimeCompare` for token comparison
- Validate all user inputs (room names, SDP size limit: 1MB)
- Never commit secrets or credentials

## Build & Test Commands

```bash
make build        # Build binary
make test         # Unit + integration + security tests
make lint         # gofmt + go vet + golangci-lint
make security     # gosec scan
```

## OpenSpec Workflow

1. `/opsx:explore` - Investigate requirements
2. `/opsx:propose` - Create change proposal
3. `/opsx:apply` - Implement tasks
4. `/opsx:archive` - Archive and merge specs

Specs location: `openspec/specs/`
