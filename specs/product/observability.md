# Metrics & Tracing (Observability)

## Overview

The project implements comprehensive observability features including Prometheus metrics, OpenTelemetry tracing, and health check endpoints.

---

## User Stories

### As an Operator
- I want to monitor active room counts and subscribers
- I want to track RTP traffic volumes
- I want to view health status via HTTP endpoint
- I want distributed tracing for debugging

### As a Developer
- I want metrics integrated into performance debugging
- I want trace context propagation for request tracking
- I want configurable observability levels

---

## Requirements

### Functional Requirements

1. **Prometheus Metrics**

| Metric | Type | Description |
|--------|------|-------------|
| `live_rooms` | Gauge | Active room count |
| `live_subscribers` | GaugeVec | Subscribers per room |
| `live_rtp_bytes_total` | CounterVec | Total RTP bytes |
| `live_rtp_packets_total` | CounterVec | Total RTP packets |

2. **OpenTelemetry Tracing**
   - Configurable via `OTEL_EXPORTER_OTLP_ENDPOINT`
   - Service name via `OTEL_SERVICE_NAME`
   - Supports stdout and OTLP (grpc/http) exporters
   - HTTP middleware for automatic span creation

3. **Health Check**
   - `GET /healthz` endpoint
   - Returns `ok` when healthy

4. **Room & Recording APIs**
   - `GET /api/rooms` - Room status with publisher/subscriber counts
   - `GET /api/records` - Recording file metadata list

---

## Acceptance Criteria

1. ✅ Metrics are exposed at `/metrics` endpoint
2. ✅ Health check returns 200 when healthy
3. ✅ OpenTelemetry traces are exported when configured
4. ✅ Room metrics reflect real-time state
5. ✅ RTP counters increase with traffic
6. ✅ Gauge values are updated on room/subscriber changes

---

## Edge Cases

1. **High Cardinality**: Room labels should not cause metric explosion
2. **Metric Reset**: Gauges should be updated, not just incremented
3. **Tracing Overhead**: Tracing should not significantly impact performance
4. **Exporter Failure**: Tracing should degrade gracefully if exporter is down

---

## Configuration Reference

| Variable | Default | Description |
|----------|---------|-------------|
| `OTEL_SERVICE_NAME` | `live-webrtc-go` | OpenTelemetry service name |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | - | OTLP exporter endpoint |

---

## Debug Endpoints

When `PPROF=1`:

| Endpoint | Description |
|----------|-------------|
| `/debug/pprof/` | pprof index |
| `/debug/pprof/profile` | CPU profile |
| `/debug/pprof/heap` | Heap profile |
| `/debug/pprof/goroutine` | Goroutine profile |

---

## Out of Scope

- Alertmanager integration
- Custom Grafana dashboards (provided separately)
- Log aggregation
