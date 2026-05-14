# Getting Started

This guide covers local development, container deployment, and basic usage of Go-Live.

## Requirements

| Dependency | Version | Notes |
|------------|---------|-------|
| Go | 1.22+ | Required for compilation |
| Docker | 20.10+ | Optional, for containerized deployment |
| Browser | Chrome 90+ / Firefox 88+ | WebRTC support required |

## Quick Start

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

### Docker Compose

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
    volumes:
      - ./records:/records
    restart: unless-stopped
```

## Stream with OBS

1. Open OBS → Settings → Stream
2. Select "WHIP" as the service
3. Server URL: `http://localhost:8080/api/whip/publish/myroom`
4. Bearer Token: `your-token` (if authentication is configured)
5. Start streaming

## Play with JavaScript

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

// Send WHEP request
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

## Next Steps

- [Architecture Overview](/en/architecture/overview) - Understand the system design
- [API Endpoints](/en/api/endpoints) - Complete API reference
- [Configuration](/en/api/configuration) - All configuration options
