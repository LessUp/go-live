# Authentication System

## Overview

The project implements a flexible authentication system supporting global tokens, per-room tokens, and JWT-based authentication with role-based access control.

---

## User Stories

### As an Administrator
- I want to set a global auth token for all rooms
- I want to set per-room tokens for granular access control
- I want to use JWT for role-based access (admin vs viewer)
- I want to protect admin endpoints with a separate admin token

### As a Developer
- I want secure token comparison (timing-attack resistant)
- I want flexible auth middleware that can be applied to any endpoint
- I want clear error messages for auth failures

---

## Requirements

### Functional Requirements

1. **Global Token Auth**
   - Single `AUTH_TOKEN` for all rooms
   - Checked via `Authorization: Bearer <token>` header

2. **Per-room Token Auth**
   - `ROOM_TOKENS` format: `room1:tok1;room2:tok2`
   - Room-specific tokens override global token

3. **JWT Authentication**
   - HMAC-based signing (HS256/HS384/HS512)
   - Signed with `JWT_SECRET`
   - Role-based access via `role=admin` or `admin=true` claims
   - Room restriction via `room` claim
   - Audience validation via `JWT_AUDIENCE`

4. **Admin Token**
   - `ADMIN_TOKEN` for admin endpoints
   - Separate from room access tokens

5. **Security**
   - Constant-time token comparison using `crypto/subtle.ConstantTimeCompare`
   - Never log tokens or secrets

---

## Acceptance Criteria

1. ✅ Requests without auth token are rejected (401)
2. ✅ Requests with invalid token are rejected (401)
3. ✅ Global token works for all rooms
4. ✅ Per-room tokens work only for specified rooms
5. ✅ Valid JWT with correct role grants access
6. ✅ JWT with wrong role is rejected
7. ✅ Admin endpoints require admin token
8. ✅ Token comparison is timing-attack resistant

---

## Edge Cases

1. **Token and JWT Both Present**: Token takes precedence if both provided
2. **Empty Token**: Empty string tokens should be rejected
3. **JWT Expiration**: Expired JWT must be rejected
4. **JWT Audience Mismatch**: Wrong audience must be rejected
5. **Room Claim Mismatch**: JWT room claim must match requested room

---

## Authentication Flow

```
Request → Check Auth Token/Room Token → Check JWT → Allow/Deny
```

### Middleware Application

```go
// Auth applied to room endpoints
mux.HandleFunc("/api/whip/publish/{room}", authOKRoom(handler))
mux.HandleFunc("/api/whep/play/{room}", authOKRoom(handler))

// Admin auth for admin endpoints
mux.HandleFunc("/api/admin/rooms/{room}/close", adminOK(handler))
```

---

## Out of Scope

- OAuth2/OIDC integration
- Database-backed user management
- API key rotation
