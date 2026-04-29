---
layout: default
title: Usage Guide
nav_order: 2
lang: en
---

# Usage Guide

This guide covers local development, container deployment, API usage, and troubleshooting for live-webrtc-go.

{: .no_toc }

## Table of Contents

{: .no_toc .text-delta }

1. TOC
{:toc}

---

## Requirements

| Dependency | Version | Notes |
|------------|---------|-------|
| Go | 1.22+ | Required for compilation |
| Docker | 20.10+ | Optional, for containerized deployment |
| Docker Compose | 2.0+ | Optional, for multi-container orchestration |
| Browser | Chrome 90+ / Firefox 88+ | WebRTC support required |

### WebRTC Port Requirements

When running behind NAT, you need to:
- Configure STUN/TURN servers
- Ensure UDP ports are not blocked by firewalls

---

## Local Development

### Method 1: Direct Execution

```bash
# Clone the repository
git clone https://github.com/LessUp/go-live.git
cd go-live

# Download dependencies
go mod tidy

# Run the server
go run ./cmd/server
```

### Method 2: Using Development Script (Recommended)

```bash
# Basic startup (loads .env.local if exists)
./scripts/start.sh

# Run go mod tidy before starting
RUN_TIDY=1 ./scripts/start.sh
```

Script features:
- Creates `records`, `.gocache`, `.gomodcache` directories
- Sets `GOCACHE` and `GOMODCACHE` environment variables
- Loads `.env.local` file (if exists)
- Starts the server

### Environment Variables File

Create `.env.local` file:

```bash
# Server configuration
HTTP_ADDR=:8080
ALLOWED_ORIGIN=*

# Authentication
AUTH_TOKEN=your-secret-token
# ROOM_TOKENS=room1:token1;room2:token2
# JWT_SECRET=jwt-signing-secret

# WebRTC configuration
# STUN_URLS=stun:stun.l.google.com:19302
# TURN_URLS=turn:turn.example.com:3478
# TURN_USERNAME=username
# TURN_PASSWORD=password

# Recording
RECORD_ENABLED=1
RECORD_DIR=records

# Rate limiting
RATE_LIMIT_RPS=10
RATE_LIMIT_BURST=20
```

### Verification After Startup

```bash
# Health check
curl http://localhost:8080/healthz
# Output: ok

# List rooms
curl http://localhost:8080/api/rooms
# Output: []

# Get frontend configuration
curl http://localhost:8080/api/bootstrap | jq
```

---

## Configuration Reference

### Core Settings

| Variable | Default | Description |
|----------|---------|-------------|
| `HTTP_ADDR` | `:8080` | HTTP listen address, format `host:port` |
| `ALLOWED_ORIGIN` | `*` | CORS allowed origin, `*` means any |

### Authentication

| Variable | Format | Description |
|----------|--------|-------------|
| `AUTH_TOKEN` | string | Global access token |
| `ROOM_TOKENS` | `room1:tok1;room2:tok2` | Per-room token mapping |
| `JWT_SECRET` | string | JWT HMAC signing key |
| `JWT_AUDIENCE` | string | Required JWT audience claim |
| `ADMIN_TOKEN` | string | Admin API access token |

**Authentication Priority**:
1. Room-specific Token → 2. Global Token → 3. JWT → 4. No authentication

### WebRTC/ICE Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `STUN_URLS` | `stun:stun.l.google.com:19302` | STUN server list |
| `TURN_URLS` | - | TURN server list |
| `TURN_USERNAME` | - | TURN username |
| `TURN_PASSWORD` | - | TURN password |

**Example**:
```bash
STUN_URLS=stun:stun1.l.google.com:19302,stun:stun2.l.google.com:19302
TURN_URLS=turn:turn.example.com:3478,turns:turn.example.com:5349
TURN_USERNAME=myuser
TURN_PASSWORD=mypassword
```

### Recording Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `RECORD_ENABLED` | `0` | Enable recording, set to `1` to enable |
| `RECORD_DIR` | `records` | Recording file storage directory |

### S3 Upload Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `UPLOAD_RECORDINGS` | `0` | Enable upload, set to `1` to enable |
| `DELETE_RECORDING_AFTER_UPLOAD` | `0` | Delete local file after upload |
| `S3_ENDPOINT` | - | S3/MinIO endpoint address |
| `S3_REGION` | - | S3 region |
| `S3_BUCKET` | - | Target bucket name |
| `S3_ACCESS_KEY` | - | Access Key ID |
| `S3_SECRET_KEY` | - | Secret Access Key |
| `S3_USE_SSL` | `1` | Use HTTPS connection |
| `S3_PATH_STYLE` | `0` | Use path-style addressing |
| `S3_PREFIX` | - | Object key prefix |

**MinIO Example**:
```bash
S3_ENDPOINT=minio.example.com:9000
S3_ACCESS_KEY=minioadmin
S3_SECRET_KEY=minioadmin
S3_BUCKET=recordings
S3_USE_SSL=0
S3_PATH_STYLE=1
```

### Rate Limiting Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `RATE_LIMIT_RPS` | `0` | Requests per second per IP, `0` disables |
| `RATE_LIMIT_BURST` | `0` | Burst capacity |

### TLS Configuration

| Variable | Description |
|----------|-------------|
| `TLS_CERT_FILE` | TLS certificate file path |
| `TLS_KEY_FILE` | TLS private key file path |

```bash
TLS_CERT_FILE=/etc/ssl/cert.pem
TLS_KEY_FILE=/etc/ssl/key.pem
```

### Debug Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PPROF` | `0` | Enable pprof endpoints |
| `OTEL_SERVICE_NAME` | `live-webrtc-go` | OpenTelemetry service name |

---

## Docker Deployment

### Build Image

```bash
docker build -t live-webrtc-go:latest .
```

### Basic Run

```bash
docker run --rm -p 8080:8080 live-webrtc-go:latest
```

### Enable Recording

```bash
docker run --rm -p 8080:8080 \
  -e RECORD_ENABLED=1 \
  -e RECORD_DIR=/records \
  -v $(pwd)/records:/records \
  live-webrtc-go:latest
```

### Full Configuration Example

```bash
docker run --rm -p 8080:8080 \
  -e HTTP_ADDR=:8080 \
  -e AUTH_TOKEN=mytoken \
  -e RECORD_ENABLED=1 \
  -e RECORD_DIR=/records \
  -e UPLOAD_RECORDINGS=1 \
  -e S3_ENDPOINT=s3.amazonaws.com \
  -e S3_BUCKET=my-bucket \
  -e S3_ACCESS_KEY=$AWS_ACCESS_KEY_ID \
  -e S3_SECRET_KEY=$AWS_SECRET_ACCESS_KEY \
  -v $(pwd)/records:/records \
  live-webrtc-go:latest
```

### Docker Compose

`docker-compose.yml`:

```yaml
version: '3.8'

services:
  live-webrtc:
    build: .
    ports:
      - "8080:8080"
    environment:
      - HTTP_ADDR=:8080
      - AUTH_TOKEN=${AUTH_TOKEN}
      - RECORD_ENABLED=1
      - RECORD_DIR=/records
      - RATE_LIMIT_RPS=10
      - RATE_LIMIT_BURST=20
    volumes:
      - ./records:/records
    restart: unless-stopped
```

Start:
```bash
docker compose up -d
```

---

## Kubernetes Deployment

### Deployment Example

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: live-webrtc
spec:
  replicas: 3
  selector:
    matchLabels:
      app: live-webrtc
  template:
    metadata:
      labels:
        app: live-webrtc
    spec:
      containers:
      - name: live-webrtc
        image: live-webrtc:latest
        ports:
        - containerPort: 8080
        env:
        - name: HTTP_ADDR
          value: ":8080"
        - name: AUTH_TOKEN
          valueFrom:
            secretKeyRef:
              name: live-webrtc-secret
              key: auth-token
        volumeMounts:
        - name: records
          mountPath: /records
      volumes:
      - name: records
        persistentVolumeClaim:
          claimName: records-pvc
```

### Service Example

```yaml
apiVersion: v1
kind: Service
metadata:
  name: live-webrtc
spec:
  selector:
    app: live-webrtc
  ports:
  - port: 80
    targetPort: 8080
  type: LoadBalancer
```

---

## API Usage Examples

### Using curl

```bash
# Set variables
HOST=http://localhost:8080
TOKEN=your-token

# Health check
curl $HOST/healthz

# Get frontend configuration
curl $HOST/api/bootstrap | jq

# List rooms
curl $HOST/api/rooms | jq

# List recording files
curl $HOST/api/records | jq

# Close room (admin API)
curl -X POST \
  -H "Authorization: Bearer $TOKEN" \
  $HOST/api/admin/rooms/myroom/close
```

### Stream with OBS

1. Open OBS → Settings → Stream
2. Select "WHIP" as the service
3. Server URL: `http://localhost:8080/api/whip/publish/myroom`
4. Bearer Token: `your-token` (if authentication is configured)
5. Start streaming

### Play with JavaScript

```javascript
// Get configuration
const config = await fetch('/api/bootstrap').then(r => r.json());

// Create PeerConnection
const pc = new RTCPeerConnection({
  iceServers: config.iceServers
});

// Create transceivers (receive only)
pc.addTransceiver('video', { direction: 'recvonly' });
pc.addTransceiver('audio', { direction: 'recvonly' });

// Handle remote tracks
pc.ontrack = (event) => {
  const video = document.getElementById('video');
  video.srcObject = event.streams[0];
};

// Create offer
const offer = await pc.createOffer();
await pc.setLocalDescription(offer);

// Send WHIP request
const response = await fetch('/api/whep/play/myroom', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/sdp',
    'Authorization': 'Bearer your-token'
  },
  body: offer.sdp
});

const answer = await response.text();
await pc.setRemoteDescription({ type: 'answer', sdp: answer });
```

---

## Frontend Integration

### Bootstrap API Response

```json
{
  "authEnabled": true,
  "recordEnabled": true,
  "iceServers": [
    {
      "urls": ["stun:stun.l.google.com:19302"]
    }
  ],
  "features": {
    "rooms": true,
    "records": true
  }
}
```

### Authentication Header Format

Two formats supported:

```
Authorization: Bearer <token>
X-Auth-Token: <token>
```

---

## Troubleshooting

### Common Issues

| Issue | Possible Cause | Solution |
|-------|---------------|----------|
| `publisher already exists in this room` | Room already has a publisher | Use a different room name or wait for publisher to disconnect |
| `unauthorized` | Authentication failed | Check token or JWT configuration |
| `too many requests` | Rate limit triggered | Increase `RATE_LIMIT_BURST` or wait |
| `no active publisher in room` | Room has no publisher | Ensure publisher is connected |
| `subscriber limit reached` | `MAX_SUBS_PER_ROOM` hit | Increase limit or wait for subscribers to disconnect |
| `ICE connection failed` | NAT traversal issue | Configure TURN server |
| No video/audio | Codec mismatch | Check browser-supported codecs |

### Debugging Steps

1. **Check Service Status**
   ```bash
   curl http://localhost:8080/healthz
   curl http://localhost:8080/api/rooms
   ```

2. **View Metrics**
   ```bash
   curl http://localhost:8080/metrics
   ```

3. **Enable pprof**
   ```bash
   PPROF=1 go run ./cmd/server
   # Access http://localhost:8080/debug/pprof/
   ```

4. **Check Logs**
   - View console `slog` output
   - Focus on `ERROR` level logs

### WebRTC Connection Issues

**Symptom**: Browser shows "ICE connection failed"

**Troubleshooting**:
1. Check STUN server reachability
2. Configure TURN server (for NAT environments)
3. Confirm UDP ports are not blocked

**TURN Configuration Example**:
```bash
TURN_URLS=turn:turn.example.com:3478
TURN_USERNAME=username
TURN_PASSWORD=password
```

### Authentication Issues

**Symptom**: API returns `401 Unauthorized`

**Troubleshooting**:
1. Confirm correct Authorization header is sent
2. Check token is correct
3. If using JWT, verify signature and expiration

### Recording Issues

**Symptom**: Recording directory is empty

**Troubleshooting**:
1. Confirm `RECORD_ENABLED=1`
2. Check directory permissions: `ls -la records/`
3. Check logs for write errors
4. Confirm publisher is using VP8/VP9/Opus codecs

### Upload Issues

**Symptom**: Files not uploaded to S3

**Troubleshooting**:
1. Confirm `UPLOAD_RECORDINGS=1`
2. Check all S3_* variables are configured
3. Verify S3 credentials are valid
4. Check network connectivity

```bash
# Test S3 connection
aws s3 ls --endpoint-url http://minio:9000
```

### Performance Issues

**Symptom**: High latency, stuttering

**Troubleshooting**:
1. Check server CPU/memory usage
2. Check subscriber count
3. Consider setting `MAX_SUBS_PER_ROOM` limit
4. Check network bandwidth

```bash
# View metrics
curl http://localhost:8080/metrics | grep live_
```

---

## Development Commands Reference

```bash
# Build
make build

# Format
make fmt

# Lint
make lint

# Security scan
make security

# Run tests
make test          # Unit + Integration + Security
make test-all      # Include e2e + Performance
make test-unit     # Unit tests only

# Coverage
make coverage
open coverage.html
```

---

## Related Documentation

- [Design Documentation](design.md) - Architecture and module details
- [API Reference](api.md) - Complete API documentation
- [GitHub Repository](https://github.com/LessUp/go-live) - Source code
