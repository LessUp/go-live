# System Architecture Overview

This document describes the overall architecture of Go-Live, a WebRTC SFU server built with Go and Pion WebRTC.

## Architecture Diagram

```mermaid
flowchart TB
    subgraph Clients
        P[Publisher<br/>OBS / Browser]
        S1[Viewer 1<br/>Browser]
        S2[Viewer 2<br/>Browser]
        S3[Viewer N<br/>Browser]
    end

    subgraph Server["Go-Live Server :8080"]
        subgraph API["HTTP Layer"]
            CORS[CORS Middleware]
            RL[Rate Limiter]
            AUTH[Auth Middleware]
            H[Handlers]
        end

        subgraph SFU["SFU Core"]
            M[Manager]
            R[Room]
            TF[TrackFanout]
        end

        subgraph Storage["Storage Layer"]
            REC[Recording<br/>IVF/OGG]
            S3[S3/MinIO<br/>Upload]
        end

        subgraph Obs["Observability"]
            PROM[Prometheus<br/>Metrics]
            OTEL[OpenTelemetry<br/>Tracing]
        end
    end

    P -->|WHIP| CORS
    S1 -->|WHEP| CORS
    S2 -->|WHEP| CORS
    S3 -->|WHEP| CORS

    CORS --> RL --> AUTH --> H
    H --> M --> R --> TF

    TF --> REC --> S3
    R --> PROM
    H --> OTEL
```

## Component Dependency Graph

```mermaid
flowchart LR
    subgraph cmd["cmd/server"]
        main[main.go]
    end

    subgraph api["internal/api"]
        handlers[handlers.go]
        middleware[middleware.go]
        routes[routes.go]
    end

    subgraph sfu["internal/sfu"]
        manager[manager.go]
        room[room.go]
        track[track.go]
    end

    subgraph infra["Infrastructure"]
        config[config]
        uploader[uploader]
        metrics[metrics]
        otel[otel]
    end

    main --> config
    main --> uploader
    main --> manager
    main --> handlers
    main --> otel

    handlers --> middleware
    handlers --> routes
    handlers --> manager

    manager --> room
    room --> track

    room --> uploader
    room --> metrics
```

## Request Processing Chain

```mermaid
flowchart TB
    REQ[HTTP Request] --> CORS{CORS Check}
    CORS -->|Pass| RL{Rate Limit}
    CORS -->|Fail| ERR1[403 Forbidden]
    RL -->|Pass| AUTH{Authentication}
    RL -->|Fail| ERR2[429 Too Many Requests]
    AUTH -->|Pass| H[Handler]
    AUTH -->|Fail| ERR3[401 Unauthorized]
    H --> RES[Response]
```

## Core Concepts

### Room

The Room is the core abstraction of the SFU. Each room:
- Has at most one Publisher
- Can have multiple Subscribers
- Has independent Track Fanout logic
- Can have its own authentication token

### Track Fanout

When a publisher pushes media tracks, the system creates a Track Fanout:
- Read RTP packets from the publisher's PeerConnection
- Copy and distribute to all subscribers
- Optionally write to recording files

### PeerConnection

Each WebRTC connection:
- **Publisher**: Receives media tracks
- **Subscriber**: Sends media tracks
- ICE negotiation completed through WHIP/WHEP protocols

## Module Responsibilities

| Module | Responsibility |
|--------|----------------|
| `cmd/server` | Application entry point, service initialization |
| `internal/config` | Environment variable parsing and defaults |
| `internal/api` | HTTP request handling, routing, middleware |
| `internal/sfu` | WebRTC SFU core logic |
| `internal/metrics` | Prometheus metrics exposure |
| `internal/uploader` | S3/MinIO file upload |
| `internal/otel` | OpenTelemetry tracer initialization |

## Key Design Decisions

### 1. Single Publisher per Room

Simplifies the SFU logic and ensures predictable stream quality. Multiple publishers would require stream selection or mixing.

### 2. In-Memory Room State

Rooms are stored in memory for simplicity and performance. For multi-instance deployment, external storage (Redis/Database) would be needed.

### 3. RTP Forwarding without Transcoding

The SFU forwards RTP packets directly without decoding/encoding, minimizing latency and CPU usage.

### 4. Recording at SFU Level

Recording happens at the TrackFanout level, capturing the exact RTP packets being distributed to subscribers.

## Performance Characteristics

| Metric | Value |
|--------|-------|
| Latency | < 100ms (same region) |
| Concurrent Subscribers | 1000+ per room |
| Memory (idle) | < 50MB |
| CPU Efficiency | Single core handles 500+ concurrent |

## Next Steps

- [SFU Core](/en/architecture/sfu-core) - Detailed SFU implementation
- [Data Flow](/en/architecture/data-flow) - Request and data flow diagrams
- [Deployment](/en/architecture/deployment) - Deployment patterns and topology
