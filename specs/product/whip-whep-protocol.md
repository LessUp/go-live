# WHIP/WHEP Protocol Support

## Overview

The project implements WHIP (WebRTC-HTTP Ingestion Protocol) and WHEP (WebRTC-HTTP Egress Protocol) for WebRTC stream publishing and playback over HTTP.

---

## User Stories

### As a Stream Publisher (e.g., OBS, Browser)
- I want to publish a WebRTC stream to a named room via HTTP POST
- I want to receive an SDP answer back to establish the WebRTC connection
- I want to be authenticated before publishing

### As a Stream Viewer (Browser)
- I want to play a WebRTC stream from a named room via HTTP POST
- I want to receive an SDP answer back to establish the WebRTC connection
- I want to be authenticated before viewing

---

## Requirements

### Functional Requirements

1. **WHIP Publishing**
   - Endpoint: `POST /api/whip/publish/{room}`
   - Request: SDP offer in request body
   - Response: SDP answer in response body
   - Creates a PeerConnection for the publisher
   - Single publisher per room

2. **WHEP Playback**
   - Endpoint: `POST /api/whep/play/{room}`
   - Request: SDP offer in request body
   - Response: SDP answer in response body
   - Creates a PeerConnection for the viewer
   - Multiple viewers per room

3. **Authentication**
   - Support global auth token (`AUTH_TOKEN`)
   - Support per-room tokens (`ROOM_TOKENS`)
   - Support JWT authentication (`JWT_SECRET`)
   - Return 401 if auth fails

4. **Room Validation**
   - Room names must match: `^[A-Za-z0-9_-]{1,64}$`
   - Return 400 for invalid room names

---

## Acceptance Criteria

1. ✅ Publisher can successfully publish a stream to a room
2. ✅ Multiple viewers can play the same stream simultaneously
3. ✅ Unauthenticated requests are rejected with 401
4. ✅ Invalid room names return 400
5. ✅ Publisher disconnect triggers cleanup of all viewers
6. ✅ CORS headers are properly set for browser access

---

## Edge Cases

1. **Room Full**: If `MAX_SUBS_PER_ROOM` is set and limit is reached, reject new viewers
2. **Publisher Disconnect**: When publisher leaves, all viewers must be notified/cleaned up
3. **SDP Size**: Reject requests with SDP larger than 1MB
4. **Concurrent Access**: Handle race conditions when multiple publishers try same room

---

## Out of Scope

- TURN/STUN server implementation (configuration only)
- Transcoding (pure SFU forwarding)
- Recording (separate spec)
