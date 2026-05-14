# Configuration

Complete reference for all environment variables and configuration options.

## Core Settings

| Variable | Default | Description |
|----------|---------|-------------|
| `HTTP_ADDR` | `:8080` | HTTP listen address, format `host:port` |
| `ALLOWED_ORIGIN` | `*` | CORS allowed origin, `*` means any |

## Authentication

| Variable | Format | Description |
|----------|--------|-------------|
| `AUTH_TOKEN` | string | Global access token |
| `ROOM_TOKENS` | `room1:tok1;room2:tok2` | Per-room token mapping |
| `JWT_SECRET` | string | JWT HMAC signing key |
| `JWT_AUDIENCE` | string | Required JWT audience claim |
| `ADMIN_TOKEN` | string | Admin API access token |

**Authentication Priority**: Room-specific Token → Global Token → JWT → No authentication

## WebRTC / ICE Configuration

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

## Recording Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `RECORD_ENABLED` | `0` | Enable recording, set to `1` to enable |
| `RECORD_DIR` | `records` | Recording file storage directory |

## S3 Upload Configuration

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

## Rate Limiting Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `RATE_LIMIT_RPS` | `0` | Requests per second per IP, `0` disables |
| `RATE_LIMIT_BURST` | `0` | Burst capacity |

## Room Limits

| Variable | Default | Description |
|----------|---------|-------------|
| `MAX_SUBS_PER_ROOM` | `0` | Max subscribers per room, `0` = unlimited |

## TLS Configuration

| Variable | Description |
|----------|-------------|
| `TLS_CERT_FILE` | TLS certificate file path |
| `TLS_KEY_FILE` | TLS private key file path |

```bash
TLS_CERT_FILE=/etc/ssl/cert.pem
TLS_KEY_FILE=/etc/ssl/key.pem
```

## Debug Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PPROF` | `0` | Enable pprof endpoints |
| `OTEL_SERVICE_NAME` | `live-webrtc-go` | OpenTelemetry service name |

## Environment File

Create `.env.local` file for development:

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

# Recording
RECORD_ENABLED=1
RECORD_DIR=records

# Rate limiting
RATE_LIMIT_RPS=10
RATE_LIMIT_BURST=20
```

## Full Docker Example

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

## Kubernetes Deployment

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

## Troubleshooting

| Issue | Possible Cause | Solution |
|-------|---------------|----------|
| `publisher already exists in this room` | Room already has a publisher | Use a different room name or wait for publisher to disconnect |
| `unauthorized` | Authentication failed | Check token or JWT configuration |
| `too many requests` | Rate limit triggered | Increase `RATE_LIMIT_BURST` or wait |
| `no active publisher in room` | Room has no publisher | Ensure publisher is connected |
| `subscriber limit reached` | `MAX_SUBS_PER_ROOM` hit | Increase limit or wait for subscribers to disconnect |
| `ICE connection failed` | NAT traversal issue | Configure TURN server |

### Debugging Steps

1. **Check Service Status**: `curl http://localhost:8080/healthz`
2. **View Metrics**: `curl http://localhost:8080/metrics`
3. **Enable pprof**: `PPROF=1 go run ./cmd/server`
4. **Check Logs**: Focus on `ERROR` level logs
