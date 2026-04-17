# Room-based SFU Relay

## Overview

The project implements a room-based Selective Forwarding Unit (SFU) where each room supports one publisher and multiple subscribers with efficient media forwarding.

---

## User Stories

### As a Room Administrator
- I want to create rooms for different streams
- I want to see active rooms and their statistics
- I want to force close a room if needed

### As a Developer
- I want efficient RTP packet forwarding from publisher to subscribers
- I want proper cleanup when publishers or subscribers disconnect
- I want metrics on room activity

---

## Requirements

### Functional Requirements

1. **Room Lifecycle**
   - Rooms are created automatically when first publisher joins
   - Rooms are destroyed when publisher leaves
   - Room names are unique and validated

2. **Publisher Management**
   - One publisher per room
   - Publisher's media tracks are distributed to all subscribers
   - Publisher disconnect triggers room cleanup

3. **Subscriber Management**
   - Multiple subscribers per room
   - Optional limit via `MAX_SUBS_PER_ROOM`
   - Subscribers receive all tracks from publisher

4. **RTP Forwarding**
   - Publisher's RTP packets are forwarded to all subscribers
   - Efficient fanout using track-level distribution
   - Track both video and audio streams

---

## Acceptance Criteria

1. ✅ Room is automatically created on first publisher join
2. ✅ Multiple subscribers can join the same room
3. ✅ RTP packets are forwarded from publisher to all subscribers
4. ✅ Room is cleaned up when publisher disconnects
5. ✅ Subscriber limits are enforced when configured
6. ✅ Active rooms can be listed via API
7. ✅ Rooms can be force-closed via admin API

---

## Edge Cases

1. **Publisher Reconnect**: Old room must be cleaned up before new publisher joins
2. **Subscriber Crash**: Individual subscriber failure should not affect others
3. **Memory Leaks**: All goroutines must exit when room is closed
4. **Concurrent Operations**: Room state must be protected with mutex

---

## Out of Scope

- Room persistence across server restarts
- Room password protection
- Multi-publisher rooms (e.g., video conferencing)
