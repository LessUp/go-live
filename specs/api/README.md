# API Specifications

This directory contains API specification documents that define the HTTP APIs, WebSocket protocols, and other interface definitions.

---

## API Specs

| File | Description | Format |
|------|-------------|--------|
| [openapi.yaml](./openapi.yaml) | REST API OpenAPI Definition | YAML |

---

## API Endpoints Overview

### Streaming

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/whip/publish/{room}` | Publish stream to room |
| `POST` | `/api/whep/play/{room}` | Play stream from room |

### Query

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/bootstrap` | Frontend runtime configuration |
| `GET` | `/api/rooms` | List active rooms |
| `GET` | `/api/records` | List recording files |

### Admin

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/admin/rooms/{room}/close` | Force close a room |

### Health & Metrics

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/healthz` | Health check |
| `GET` | `/metrics` | Prometheus metrics |

---

## API Conventions

### Authentication

All endpoints (except `/healthz` and `/api/bootstrap`) require authentication via:
- **Token**: `Authorization: Bearer <token>` header
- **JWT**: `Authorization: Bearer <jwt>` header

### Error Responses

All error responses follow a consistent format:

```json
{
  "error": "Error message describing what went wrong"
}
```

### Response Codes

| Code | Meaning |
|------|---------|
| `200` | Success |
| `400` | Bad Request (invalid input) |
| `401` | Unauthorized (missing/invalid auth) |
| `403` | Forbidden (insufficient permissions) |
| `404` | Not Found |
| `500` | Internal Server Error |
