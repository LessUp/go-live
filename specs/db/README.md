# Database Specifications

This directory contains database schema definitions, migration specifications, and data model documentation.

---

## Database Overview

The live-webrtc-go project uses **in-memory data structures** with no persistent database. State is managed through Go data structures with mutex protection.

---

## In-Memory Data Models

### Room State

```go
type Room struct {
    Name        string
    Publisher   *Publisher
    Subscribers map[string]*Subscriber
    mu          sync.RWMutex
}
```

### Configuration

Configuration is loaded from environment variables at startup. See the configuration spec in the main documentation.

---

## Future Database Support

If persistent database support is added in the future (e.g., for room persistence, user management, etc.), the schema files will be placed here in DBML or SQL format.

---

## Naming Convention

Schema files should be named as: `schema-{version}.dbml` or `schema-{version}.sql`
