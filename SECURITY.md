# Security Policy

## Reporting a Vulnerability

If you discover a security vulnerability, please report it privately:

- **Email**: Create an issue with the `security` label (it will be made private)
- **GitHub**: Use the "Report a vulnerability" feature in the Security tab

Please do **not** disclose vulnerabilities publicly until they have been addressed.

## Supported Versions

| Version | Supported |
| ------- | --------- |
| 1.0.x   | ✅ Active |

## Security Features

This project implements the following security measures:

### Authentication
- Token-based authentication (global and per-room)
- JWT authentication with role-based access
- Constant-time token comparison to prevent timing attacks

### Input Validation
- Room name validation: `^[A-Za-z0-9_-]{1,64}$`
- SDP body size limit: 1MB max
- Rate limiting per IP

### Network
- CORS configuration via `ALLOWED_ORIGIN`
- TLS support via `TLS_CERT_FILE` and `TLS_KEY_FILE`
- ICE/TURN server configuration

### Best Practices
- No secrets in logs
- No hardcoded credentials
- Dependencies regularly updated
- Security tests in CI pipeline

## Security Configuration

### Recommended Production Settings

```bash
# Enable authentication
AUTH_TOKEN=<strong-random-token>

# Or use JWT
JWT_SECRET=<strong-random-secret>

# Configure CORS properly
ALLOWED_ORIGIN=https://your-domain.com

# Enable rate limiting
RATE_LIMIT_RPS=10
RATE_LIMIT_BURST=20

# Use TLS
TLS_CERT_FILE=/path/to/cert.pem
TLS_KEY_FILE=/path/to/key.pem

# Configure TURN for NAT traversal
TURN_URLS=turn:turn.example.com:3478
TURN_USERNAME=<username>
TURN_PASSWORD=<password>
```

### Secrets Management

- Never commit secrets to the repository
- Use environment variables or secret management systems
- Rotate tokens and secrets regularly
- Use strong, randomly generated tokens (min 32 characters)

## Security Audit

Run security scans:

```bash
# Run gosec
make security

# Or directly
go install github.com/securego/gosec/v2/cmd/gosec@latest
gosec ./...
```

## Disclosure Policy

1. Report received and acknowledged within 48 hours
2. Vulnerability confirmed and severity assessed
3. Fix developed and tested
4. Patch released and announced
5. Public disclosure after fix is available

Thank you for helping keep this project secure!
