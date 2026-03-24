# Contributing to live-webrtc-go

Thank you for contributing.

## Development setup

```bash
git clone https://github.com/YOUR_USERNAME/go-live.git
cd go-live
go mod download
make install-tools
```

## Recommended workflow

1. Create a branch from `master`
2. Make focused changes
3. Run local verification
4. Open a pull request

## Local verification

Use these commands before submitting changes:

```bash
make fmt
make lint
make security
make test
make test-all
make coverage
```

Command meanings:

- `make test` → unit + integration + security tests
- `make test-all` → adds e2e + performance suites
- `make ci` → local CI-style verification pipeline

## Coding expectations

- Follow normal Go idioms
- Add tests for behavior changes
- Update docs when config or runtime behavior changes
- Avoid introducing stale compatibility shims or dead code
- Keep fixes focused and minimal

## Security expectations

- Never commit secrets, credentials, or tokens
- Validate boundary inputs
- Run `make security` for non-trivial changes
- If you modify auth, routing, or request parsing, add tests

## Pull requests

Please include:

- a concise summary
- why the change is needed
- test evidence (`make test`, `make test-all`, etc.)
- any config or docs changes

## Notes

- The primary branch is `master`
- Go baseline is 1.22+
- Tagged test suites live under `test/`
