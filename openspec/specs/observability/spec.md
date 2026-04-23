# Metrics & Tracing (Observability)

## Purpose

Comprehensive observability including Prometheus metrics, OpenTelemetry tracing, and health check endpoints for monitoring and debugging.

## Requirements

### Requirement: Prometheus Metrics

The system SHALL expose Prometheus metrics at /metrics endpoint.

#### Scenario: Room metrics
- **WHEN** rooms are created or destroyed
- **THEN** live_rooms gauge reflects current active room count

#### Scenario: Subscriber metrics
- **WHEN** subscribers join or leave rooms
- **THEN** live_subscribers GaugeVec reflects counts per room

#### Scenario: RTP metrics
- **WHEN** RTP packets are forwarded
- **THEN** live_rtp_bytes_total and live_rtp_packets_total counters increment

### Requirement: OpenTelemetry Tracing

The system SHALL support distributed tracing via OpenTelemetry.

#### Scenario: OTLP exporter configured
- **WHEN** OTEL_EXPORTER_OTLP_ENDPOINT is set
- **THEN** system exports traces to configured OTLP endpoint

#### Scenario: Service name
- **WHEN** OTEL_SERVICE_NAME is set
- **THEN** traces use configured service name

#### Scenario: HTTP middleware tracing
- **WHEN** HTTP request is received
- **THEN** system creates span for request tracing

### Requirement: Health Check

The system SHALL provide health check endpoint.

#### Scenario: Health check
- **WHEN** client calls GET /healthz
- **THEN** system returns 200 OK with body "ok"

### Requirement: Debug Endpoints

The system SHALL provide pprof endpoints when enabled.

#### Scenario: Pprof enabled
- **WHEN** PPROF=1 environment variable is set
- **THEN** system exposes /debug/pprof/* endpoints

#### Scenario: CPU profile
- **WHEN** client requests /debug/pprof/profile
- **THEN** system returns CPU profile

#### Scenario: Heap profile
- **WHEN** client requests /debug/pprof/heap
- **THEN** system returns heap profile
