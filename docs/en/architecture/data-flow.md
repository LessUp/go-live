# Data Flow

Detailed documentation of request and data flow through the system.

## WHIP Publishing Flow

```mermaid
sequenceDiagram
    participant P as Publisher<br/>(OBS/Browser)
    participant H as HTTP Handler
    participant M as Manager
    participant R as Room
    participant PC as PeerConnection
    participant TF as TrackFanout

    P->>H: POST /api/whip/publish/{room}<br/>SDP Offer
    H->>H: CORS Check
    H->>H: Rate Limit Check
    H->>H: Auth Check
    
    H->>M: Publish(roomName, offer)
    M->>M: getOrCreateRoom(roomName)
    M->>R: Publish(offer)
    
    R->>PC: NewPeerConnection(iceConfig)
    PC->>PC: SetRemoteDescription(offer)
    PC->>PC: CreateAnswer()
    PC->>PC: SetLocalDescription(answer)
    
    PC-->>R: OnTrack callback
    R->>TF: attachTrackFeed(track)
    TF->>TF: Start readLoop()
    
    R-->>M: answer SDP
    M-->>H: answer SDP
    H-->>P: 201 Created + SDP Answer
    
    Note over TF: readLoop runs continuously<br/>distributing RTP packets
```

## WHEP Playback Flow

```mermaid
sequenceDiagram
    participant V as Viewer<br/>(Browser)
    participant H as HTTP Handler
    participant M as Manager
    participant R as Room
    participant PC as PeerConnection
    participant TF as TrackFanout

    V->>H: POST /api/whep/play/{room}<br/>SDP Offer
    H->>H: CORS + Rate Limit + Auth Check
    
    H->>M: Subscribe(roomName, offer)
    M->>R: Subscribe(offer)
    
    R->>R: Check subscriber limit
    
    loop For each existing TrackFanout
        R->>TF: attachToSubscriber()
        TF->>PC: AddTrack(localTrack)
    end
    
    R->>PC: NewPeerConnection(iceConfig)
    PC->>PC: SetRemoteDescription(offer)
    PC->>PC: CreateAnswer()
    PC->>PC: SetLocalDescription(answer)
    
    PC-->>R: OnICEStateChange handler
    R-->>M: answer SDP
    M-->>H: answer SDP
    H-->>V: 201 Created + SDP Answer
    
    Note over PC,TF: RTP packets flow from<br/>TrackFanout to Subscriber
```

## RTP Packet Flow

```mermaid
flowchart LR
    subgraph Publisher
        PT[Publisher Track]
    end

    subgraph SFU
        RT[Remote Track] -->|ReadRTP| BUF[RTP Buffer]
        BUF -->|WriteRTP| LT1[Local Track 1]
        BUF -->|WriteRTP| LT2[Local Track 2]
        BUF -->|WriteRTP| LT3[Local Track N]
        BUF -->|WriteRTP| REC[Recorder]
    end

    subgraph Subscribers
        LT1 --> S1[Subscriber 1]
        LT2 --> S2[Subscriber 2]
        LT3 --> SN[Subscriber N]
    end

    PT -->|RTP| RT
```

## Disconnection Flow

### Publisher Disconnect

```mermaid
sequenceDiagram
    participant PC as PeerConnection
    participant R as Room
    participant TF as TrackFanout
    participant U as Uploader
    participant S as Subscribers

    PC->>PC: ICE State Change<br/>(Failed/Disconnected/Closed)
    PC-->>R: OnICEStateChange callback
    
    R->>R: closePublisher()
    
    loop For each TrackFanout
        R->>TF: close()
        TF->>TF: Stop readLoop
        TF-->>R: Return recording file path
    end
    
    R->>S: Close all subscriber connections
    R->>R: Clear subscriber list
    
    alt Upload Enabled
        R->>U: uploadRecording(paths)
        U->>U: PutObject to S3
        opt Delete After Upload
            U->>U: Remove local files
        end
    end
    
    R->>R: pruneIfEmpty()
```

### Subscriber Disconnect

```mermaid
sequenceDiagram
    participant PC as PeerConnection
    participant R as Room
    participant TF as TrackFanout

    PC->>PC: ICE State Change
    PC-->>R: OnICEStateChange callback
    
    R->>R: removeSubscriber(id)
    
    loop For each TrackFanout
        R->>TF: Remove local track binding
    end
    
    R->>PC: Close PeerConnection
    R->>R: pruneIfEmpty()
```

## Authentication Flow

```mermaid
flowchart TB
    REQ[Request] --> HEADER{Check Headers}
    
    HEADER -->|Authorization: Bearer| TOKEN[Extract Token]
    HEADER -->|X-Auth-Token| TOKEN
    HEADER -->|No Auth Header| NOAUTH[No Auth]
    
    TOKEN --> ROOM{Room Token<br/>ROOM_TOKENS?}
    ROOM -->|Found| ROOMCHECK[Compare Token]
    ROOM -->|Not Found| GLOBAL{Global Token<br/>AUTH_TOKEN?}
    
    ROOMCHECK -->|Match| ALLOW[Allow]
    ROOMCHECK -->|No Match| GLOBAL
    
    GLOBAL -->|Found| GLOBALCHECK[Compare Token]
    GLOBAL -->|Not Found| JWT{JWT<br/>JWT_SECRET?}
    
    GLOBALCHECK -->|Match| ALLOW
    GLOBALCHECK -->|No Match| JWT
    
    JWT -->|Found| JWTCHECK[Verify JWT]
    JWT -->|Not Found| NOAUTH
    
    JWTCHECK -->|Valid| JWTAUTH{Check Claims}
    JWTCHECK -->|Invalid| DENY[401 Unauthorized]
    
    JWTAUTH -->|Room matches| ALLOW
    JWTAUTH -->|Room mismatch| DENY
    
    NOAUTH -->|Auth configured| DENY
    NOAUTH -->|Auth not configured| ALLOW
```

## Recording Flow

```mermaid
sequenceDiagram
    participant TF as TrackFanout
    participant R as Recorder
    participant FS as FileSystem
    participant S3 as S3/MinIO

    TF->>TF: readLoop receives RTP
    
    alt Video Track (VP8/VP9)
        TF->>R: WriteRTP to IVF Writer
    else Audio Track (Opus)
        TF->>R: WriteRTP to OGG Writer
    end
    
    Note over R: File: {room}_{trackID}_{timestamp}.{ext}
    
    R->>FS: Write to RECORD_DIR
    
    opt Publisher Disconnects
        TF->>R: Close()
        R->>FS: Close file
        
        alt UPLOAD_RECORDINGS=1
            TF->>S3: PutObject(bucket, key, file)
            
            opt DELETE_RECORDING_AFTER_UPLOAD=1
                TF->>FS: Remove local file
            end
        end
    end
```

## Metrics Update Flow

```mermaid
flowchart TB
    subgraph RTP["RTP Processing"]
        READ[ReadRTP] --> UPDATE[Update Metrics]
        UPDATE --> DISTRIB[Distribute to Subscribers]
    end

    subgraph Metrics["Prometheus Metrics"]
        UPDATE --> ROOMS[live_rooms Gauge]
        UPDATE --> SUBS[live_subscribers GaugeVec]
        UPDATE --> BYTES[live_rtp_bytes_total CounterVec]
        UPDATE --> PKTS[live_rtp_packets_total CounterVec]
    end

    subgraph Export["Export"]
        ROOMS --> PROM[/metrics endpoint]
        SUBS --> PROM
        BYTES --> PROM
        PKTS --> PROM
    end
```

## Request Rate Limiting

```mermaid
flowchart TB
    REQ[Request] --> IP[Extract Client IP]
    IP --> BUCKET{Token Bucket<br/>Available?}
    
    BUCKET -->|Yes| ALLOW[Allow Request]
    BUCKET -->|No| REJECT[429 Too Many Requests]
    
    ALLOW --> PROCESS[Process Request]
    PROCESS --> UPDATE[Update Bucket<br/>-1 Token]
    
    Note over BUCKET: Refills at RATE_LIMIT_RPS<br/>Burst capacity: RATE_LIMIT_BURST
```
