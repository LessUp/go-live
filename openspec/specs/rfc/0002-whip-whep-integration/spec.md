# RFC 0002: WHIP/WHEP Protocol Integration

## Summary

This RFC documents the integration of WHIP (WebRTC-HTTP Ingestion Protocol) and WHEP (WebRTC-HTTP Egress Protocol) for standard HTTP-based WebRTC stream publishing and playback.

## Motivation

WHIP and WHEP are emerging standards for WebRTC signaling over HTTP. Adopting these protocols enables:
- Compatibility with OBS Studio and modern browsers
- Simple HTTP-based signaling (no WebSocket required)
- Easy integration with CDNs and reverse proxies
- Standard REST API patterns

## Design

### Protocol Flow

#### WHIP Publishing

```
Publisher                          Server
    │                                │
    │──── POST /api/whip/publish/{room} ────▶│
    │      (SDP Offer)               │
    │                                │
    │◀──── 200 OK ──────────────────▶│
    │      (SDP Answer)              │
    │                                │
    │◀──────── ICE / WebRTC ────────▶│
    │         RTP Stream             │
```

#### WHEP Playback

```
Viewer                            Server
    │                                │
    │──── POST /api/whep/play/{room} ───────▶│
    │      (SDP Offer)               │
    │                                │
    │◀──── 200 OK ──────────────────▶│
    │      (SDP Answer)              │
    │                                │
    │◀──────── ICE / WebRTC ────────▶│
    │         RTP Stream             │
```

### API Design

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `POST` | `/api/whip/publish/{room}` | Token/JWT | Publish stream (SDP Offer → Answer) |
| `POST` | `/api/whep/play/{room}` | Token/JWT | Play stream (SDP Offer → Answer) |

### SDP Exchange

1. Client creates SDP Offer
2. Client POSTs SDP Offer to server
3. Server creates PeerConnection
4. Server sets remote description (client's offer)
5. Server creates local description (answer)
6. Server returns SDP Answer in HTTP 200 response
7. Client sets remote description (server's answer)
8. ICE negotiation begins

### Integration with SFU

```
HTTP Handler (api/handlers.go)
    │
    ├─ Authenticate request
    ├─ Validate room name
    ├─ Check rate limit
    │
    ▼
SFU Manager (sfu/manager.go)
    │
    ├─ Get or create Room
    ├─ Create Publisher/Subscriber PeerConnection
    └─ Setup track fanout
    │
    ▼
Room (sfu/room.go)
    │
    ├─ Handle SDP exchange
    ├─ Manage ICE candidates
    ├─ Setup track forwarding
    └─ Handle recording
```

## Alternatives Considered

### WebSocket-based Signaling
- **Rejected**: More complex infrastructure requirement
- **WHIP/WHEP Advantage**: Simple HTTP, works with any HTTP client

### Custom Protocol
- **Rejected**: Interoperability concerns
- **WHIP/WHEP Advantage**: Industry standard, growing adoption

## Status

✅ Approved - WHIP/WHEP are well-defined protocols.
