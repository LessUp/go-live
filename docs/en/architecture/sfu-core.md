# SFU Core Implementation

Detailed documentation of the SFU (Selective Forwarding Unit) core logic.

## Component Hierarchy

```mermaid
flowchart TB
    M[Manager] -->|1:N| R1[Room 1]
    M -->|1:N| R2[Room 2]
    M -->|1:N| RN[Room N]

    R1 -->|1:N| TF1[TrackFanout 1]
    R1 -->|1:N| TF2[TrackFanout 2]

    TF1 -->|reads| RT1[Remote Track<br/>from Publisher]
    TF1 -->|writes| LT1[Local Tracks<br/>to Subscribers]
    TF1 -->|writes| REC1[Recorder<br/>IVF/OGG]
```

## Manager

The Manager is the top-level component that manages all rooms.

```go
type Manager struct {
    rooms     map[string]*Room    // Room name → Room instance
    roomsMu   sync.RWMutex        // Protects rooms map
    iceConfig ICEConfig           // ICE server configuration
    config    Config              // Server configuration
}
```

### Key Methods

| Method | Purpose |
|--------|---------|
| `Publish(roomName, sdpOffer)` | Create room, establish publisher connection |
| `Subscribe(roomName, sdpOffer)` | Create subscriber connection, bind existing tracks |
| `CloseRoom(roomName)` | Force close a room |
| `ListRooms()` | Return all active rooms |
| `RoomCount()` | Return number of active rooms |

## Room

Each Room represents an isolated streaming session.

```go
type Room struct {
    name         string
    publisher    *webrtc.PeerConnection
    subscribers  map[string]*webrtc.PeerConnection
    trackFeeds   map[uint32]*TrackFanout  // SSRC → TrackFanout
    mu           sync.RWMutex
    config       Config
}
```

### Room Lifecycle

```mermaid
stateDiagram-v2
    [*] --> Empty: Room Created
    Empty --> Active: Publisher Connects
    Active --> Recording: Recording Enabled
    Recording --> Active: Publisher Disconnects
    Active --> Empty: Publisher Disconnects
    Empty --> [*]: Room Pruned
```

### Key Methods

| Method | Purpose |
|--------|---------|
| `Publish(sdpOffer)` | Establish publisher PeerConnection |
| `Subscribe(sdpOffer)` | Add subscriber, bind to existing tracks |
| `attachTrackFeed(track)` | Distribute new track to all subscribers |
| `closePublisher()` | Clean up publisher and all recordings |
| `removeSubscriber(id)` | Remove subscriber connection |
| `pruneIfEmpty()` | Delete room if no publisher/subscribers |

## TrackFanout

TrackFanout handles RTP packet distribution for a single media track.

```go
type TrackFanout struct {
    remoteTrack  *webrtc.TrackRemote
    localTracks  map[string]*webrtc.TrackLocalStaticRTP
    recorder     rtpWriter  // IVF or OGG writer
    stopChan     chan struct{}
}
```

### RTP Distribution Flow

```mermaid
sequenceDiagram
    participant P as Publisher
    participant TF as TrackFanout
    participant S1 as Subscriber 1
    participant S2 as Subscriber 2
    participant R as Recorder

    P->>TF: RTP Packet
    TF->>TF: ReadRTP()
    
    par Distribute
        TF->>S1: WriteRTP()
        TF->>S2: WriteRTP()
        TF->>R: WriteRTP()
    end
```

### readLoop

The core loop that reads RTP packets and distributes them:

```go
func (tf *TrackFanout) readLoop() {
    for {
        rtp, err := tf.remoteTrack.ReadRTP()
        if err != nil {
            return
        }

        // Distribute to all subscribers
        for _, local := range tf.localTracks {
            local.WriteRTP(rtp)
        }

        // Write to recorder if enabled
        if tf.recorder != nil {
            tf.recorder.WriteRTP(rtp)
        }
    }
}
```

## PeerConnection Management

### Publisher Connection

```mermaid
sequenceDiagram
    participant C as Client
    participant H as Handler
    participant M as Manager
    participant R as Room
    participant PC as PeerConnection

    C->>H: POST /whip/publish/{room}
    H->>M: Publish(room, offer)
    M->>R: getOrCreateRoom(room)
    M->>R: Publish(offer)
    R->>PC: NewPeerConnection(iceConfig)
    R->>PC: SetRemoteDescription(offer)
    R->>PC: CreateAnswer()
    R->>PC: SetLocalDescription(answer)
    R->>PC: OnTrack → attachTrackFeed()
    R-->>H: answer SDP
    H-->>C: 201 Created + answer
```

### Subscriber Connection

```mermaid
sequenceDiagram
    participant C as Client
    participant H as Handler
    participant M as Manager
    participant R as Room
    participant PC as PeerConnection

    C->>H: POST /whep/play/{room}
    H->>M: Subscribe(room, offer)
    M->>R: Subscribe(offer)
    R->>PC: NewPeerConnection(iceConfig)
    
    loop For each TrackFanout
        R->>PC: AddTrack(localTrack)
    end
    
    R->>PC: SetRemoteDescription(offer)
    R->>PC: CreateAnswer()
    R->>PC: SetLocalDescription(answer)
    R-->>H: answer SDP
    H-->>C: 201 Created + answer
```

## ICE Configuration

ICE servers are configured via environment variables:

```go
type ICEConfig struct {
    STUNURLs   []string
    TURNURLs   []string
    TURNUser   string
    TURNPass   string
}
```

Default STUN: `stun:stun.l.google.com:19302`

## Memory Layout

```
Manager
├── rooms map[string]*Room
│   ├── "room1" → Room
│   │   ├── publisher *PeerConnection
│   │   ├── subscribers map[string]*PeerConnection
│   │   │   ├── "sub1" → PeerConnection
│   │   │   └── "sub2" → PeerConnection
│   │   └── trackFeeds map[uint32]*TrackFanout
│   │       ├── 12345 → TrackFanout (video)
│   │       │   ├── remoteTrack
│   │       │   ├── localTracks map[string]*TrackLocalStaticRTP
│   │       │   └── recorder
│   │       └── 67890 → TrackFanout (audio)
│   └── "room2" → Room
└── iceConfig
```

## Concurrency Model

- **Manager.roomsMu**: Protects the rooms map
- **Room.mu**: Protects publisher, subscribers, trackFeeds
- **TrackFanout**: Single goroutine (readLoop) per track

### Goroutine Lifecycle

```mermaid
flowchart TB
    subgraph Main["Main Goroutine"]
        HTTP[HTTP Server]
    end

    subgraph PerRoom["Per Room"]
        PCH[Publisher OnTrack Handlers]
        SUB[Subscriber Handlers]
    end

    subgraph PerTrack["Per Track"]
        RL[readLoop]
    end

    HTTP -->|Create Room| PCH
    PCH -->|OnTrack| RL
    RL -->|WriteRTP| SUB
```

## Next Steps

- [Data Flow](/en/architecture/data-flow) - Complete request flow diagrams
- [WHIP Protocol](/en/protocols/whip) - WHIP publishing details
- [WHEP Protocol](/en/protocols/whep) - WHEP playback details
