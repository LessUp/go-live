# Version Requirements Specification

## Purpose
Lock and document minimum and recommended versions for Go runtime, dependencies, and build/CI tooling to ensure reproducible builds and stable deployment.

## Runtime Requirements

### Go
- **Minimum**: 1.22
- **Recommended**: 1.22.x (latest patch)
- **Reason**: Go 1.22 introduced `for range` iterator feature and improved standard library. Locked at 1.22 to ensure compatibility across local dev, CI, and Docker builds.
- **Where**: 
  - `go.mod`: Declared as `go 1.22`
  - `.github/workflows/ci.yml`: `actions/setup-go@v5` with `go-version: 1.22`
  - `Dockerfile`: `FROM golang:1.22-alpine`

### Docker
- **Build Base**: `golang:1.22-alpine` (matches Go version)
- **Runtime Base**: `alpine:latest` (minimal footprint)
- **Reason**: Alpine reduces image size; pinned Go version ensures build consistency.

## Key Dependencies

### WebRTC & Streaming
- **pion/webrtc**: v3.2.24 (latest stable)
  - Reason: Provides RFC-compliant WebRTC implementation; pinned to ensure WHIP/WHEP compatibility
- **pion/rtcp**: v1.2.14
- **pion/rtp**: v1.8.7
- **pion/interceptor**: v0.1.29

### Authentication & Security
- **golang-jwt**: v5.2.0 (JWT token signing)
  - Reason: Stable JWT implementation with Go 1.22 support

### Storage & Observability
- **minio/minio-go**: v7.0.66 (S3-compatible object storage)
- **prometheus/client_golang**: v1.19.0 (metrics export)
- **go.opentelemetry.io/otel**: v1.28.0 (distributed tracing)
- **golang.org/x/time**: v0.5.0 (rate limiting primitives)

**Update Policy**: Dependencies reviewed and updated quarterly or when security advisories are announced.

## Build & CI Tooling

### GitHub Actions
- **checkout**: v4 (fetch repo)
- **setup-go**: v5 (install Go)
- **golangci-lint**: v6 (linting)
- **gosec**: master (security scanning)
- **codecov**: v4 (coverage reporting)
- **docker/setup-buildx**: v3 (Docker multi-arch builds)
- **docker/login**: v3 (Docker registry auth)
- **docker/metadata**: v5 (image tagging)
- **docker/build-push**: v6 (publish images)
- **upload-artifact**: v4 (store build outputs)

### Local Development Tools
- **make**: system default (build orchestration)
- **gofmt**: Go 1.22 (formatting, via `gofmt -s`)
- **go vet**: Go 1.22 (static analysis)
- **golangci-lint**: 1.56.0+ (comprehensive linting)
- **gosec**: 2.19.0+ (security scanning)

### Jekyll & Docs
- Ruby: 3.2+ (required by GitHub Pages)
- Jekyll: 4.3+ (static site generation)

## Constraints

1. **No pre-release versions in production**: All locked versions must be stable releases.
2. **Synchronize Go versions**: Local dev (`go.mod`), Docker (`Dockerfile`), and CI (`ci.yml`) must use the same Go version.
3. **Update only after testing**: Version bumps require full `make test` suite to pass before committing.
4. **Security patches**: Any CVE affecting locked versions triggers immediate patched version deployment.

## Related Documentation
- See `openspec/specs/dependencies/` for detailed dependency rationale and risk assessment.
- See `Makefile` for build/test commands that rely on these versions.
- See `.github/workflows/` for CI tooling version declarations.

## Last Updated
2024-12-20 (Phase 1 version pinning)
