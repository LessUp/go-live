# AI Agent Instructions

This file provides guidance to AI coding assistants (Claude Code, Cursor, Windsurf, etc.) when working with this repository.

---

## Project Philosophy: Spec-Driven Development (SDD)

This project strictly follows the **Spec-Driven Development** paradigm. All code implementation must be based on the specification documents in the `/specs` directory as the **Single Source of Truth**.

---

## Directory Context

| Directory | Purpose |
|-----------|---------|
| `/specs/product/` | Product feature definitions and acceptance criteria (PRDs) |
| `/specs/rfc/` | Technical design documents and architecture decisions (RFCs) |
| `/specs/api/` | API interface definitions (OpenAPI, GraphQL, etc.) |
| `/specs/db/` | Database schema specifications |
| `/specs/testing/` | BDD test specifications (Gherkin `.feature` files) |
| `/docs/` | User guides, tutorials, and developer documentation |

---

## AI Agent Workflow Instructions

When you (the AI) are asked to implement a new feature, modify existing functionality, or fix a bug, **you MUST follow this workflow strictly. Do NOT skip any steps**:

### Step 1: Review Specifications (Review Specs)

**BEFORE writing any code:**

1. Read the relevant documents in `/specs`:
   - Product specs in `/specs/product/`
   - RFCs in `/specs/rfc/`
   - API specs in `/specs/api/`
   - DB specs in `/specs/db/`

2. If the user's request conflicts with existing specs:
   - **STOP** coding immediately
   - Point out the conflict to the user
   - Ask if they want to update the spec first

### Step 2: Spec-First Update

**If this is a new feature or changes existing interfaces/database:**

1. **Propose changes to the relevant spec documents FIRST**:
   - Update or create product specs in `/specs/product/`
   - Update or create RFCs in `/specs/rfc/`
   - Update API specs in `/specs/api/openapi.yaml`
   - Update DB specs in `/specs/db/`

2. **Wait for user confirmation** on the spec changes before proceeding to code

3. Only after specs are confirmed, proceed to implementation

### Step 3: Code Implementation

**When writing code:**

1. **Follow specs 100%**: Variable names, API paths, data types, status codes must match specs exactly
2. **No gold-plating**: Do NOT add features not defined in specs
3. **Use existing patterns**: Follow conventions defined in `/specs/rfc/` for architecture patterns
4. **Error handling**: Wrap errors with context using `fmt.Errorf("operation failed: %w", err)`

### Step 4: Test Against Spec

**Verification:**

1. Write unit and integration tests based on acceptance criteria in specs
2. Test cases MUST cover all boundary conditions described in specs
3. For API changes, verify against `/specs/api/openapi.yaml`
4. Run `make test` and `make lint` before confirming completion

---

## Code Generation Rules

1. **API Changes**: Any external API change MUST update `/specs/api/openapi.yaml` first
2. **Database Changes**: Any schema change MUST update `/specs/db/` first
3. **Architecture Decisions**: Reference `/specs/rfc/` for design patterns - do NOT invent new patterns without discussion
4. **Uncertain Details**: If a technical detail is unclear, consult `/specs/rfc/` - do NOT make up design decisions
5. **No Spec, No Code**: If no spec exists for a feature, create one first before implementing

---

## Build & Run Commands

```bash
# Build
go build -o bin/live-webrtc-go ./cmd/server

# Run
go run ./cmd/server

# Development mode (loads .env.local)
./scripts/start.sh

# With module tidy
RUN_TIDY=1 ./scripts/start.sh
```

---

## Test Commands

```bash
make test              # Unit + integration + security
make test-unit         # Unit tests only
make test-integration  # Integration tests (-tags=integration)
make test-e2e          # E2E tests (-tags=e2e, 10m timeout)
make test-all          # All tests
make coverage          # Generate coverage report
```

---

## Lint & Security

```bash
make lint        # gofmt + go vet + golangci-lint
make fmt         # gofmt -s -w .
make security    # gosec ./...
make ci          # Full CI pipeline
```

---

## Custom Skills & Hooks

### Verify Skill

After making code changes, run verification to catch issues early:

```bash
# Run lint and unit tests
make lint && make test-unit
```

Report any failures clearly. If all pass, confirm briefly.

### Auto-Format Hook

Go files are automatically formatted with `gofmt -s` after Write/Edit operations. This hook runs automatically in Claude Code.

---

## Project Structure Overview

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
├── specs/                # ← Single Source of Truth (Specs)
│   ├── product/          # Product requirements
│   ├── rfc/              # Technical designs
│   ├── api/              # API definitions
│   ├── db/               # Database schemas
│   └── testing/          # BDD test specs
├── test/                 # Test implementations
│   ├── integration/      # Integration tests
│   ├── e2e/              # End-to-end tests
│   ├── security/         # Security tests
│   ├── performance/      # Benchmarks
│   ├── load/             # Load testing tools
│   └── reports/          # Test reports (generated)
├── docs/                 # User & developer documentation
│   ├── en/               # English docs
│   ├── zh/               # Chinese docs
│   └── changelog/        # Changelog templates & release notes
└── .github/              # GitHub workflows and templates
```

---

## Why These Rules Exist

1. **Prevent AI Hallucination**: Forcing spec review anchors the AI's thinking to the project context
2. **Document-Code Synchronization**: Specs-first approach ensures docs and code stay in sync
3. **PR Quality**: Implementation aligned with business logic produces higher quality pull requests

---

## Additional Resources

- [CONTRIBUTING.md](./CONTRIBUTING.md) - Contribution guidelines
- [CLAUDE.md](./CLAUDE.md) - Claude Code specific configuration
- [CHANGELOG.md](./CHANGELOG.md) - Project history
- [specs/README.md](./specs/README.md) - Specifications overview
- [docs/](./docs/) - User and developer documentation
