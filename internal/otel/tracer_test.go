package otel

import (
	"context"
	"os"
	"testing"

	"go.opentelemetry.io/otel"
)

func TestInitTracer_Stdout(t *testing.T) {
	// Ensure no OTLP endpoint is set
	origEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	_ = os.Unsetenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	defer func() {
		if origEndpoint != "" {
			_ = os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", origEndpoint)
		}
	}()

	shutdown, err := InitTracer("test-service")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if shutdown == nil {
		t.Error("expected shutdown function to be returned")
	}

	// Verify shutdown works
	if err := shutdown(context.Background()); err != nil {
		t.Errorf("expected shutdown to succeed, got %v", err)
	}
}

func TestInitTracer_ServiceName(t *testing.T) {
	origEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	_ = os.Unsetenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	defer func() {
		if origEndpoint != "" {
			_ = os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", origEndpoint)
		}
	}()

	shutdown, err := InitTracer("my-custom-service")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer func() {
		_ = shutdown(context.Background())
	}()

	// Verify tracer provider is set
	tp := otel.GetTracerProvider()
	if tp == nil {
		t.Error("expected tracer provider to be set")
	}
}

func TestInitTracer_MultipleShutdown(t *testing.T) {
	origEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	_ = os.Unsetenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	defer func() {
		if origEndpoint != "" {
			_ = os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", origEndpoint)
		}
	}()

	shutdown, err := InitTracer("test-service")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// First shutdown should work
	if err := shutdown(context.Background()); err != nil {
		t.Errorf("expected first shutdown to succeed, got %v", err)
	}

	// Second shutdown should also work (no-op)
	if err := shutdown(context.Background()); err != nil {
		t.Errorf("expected second shutdown to succeed, got %v", err)
	}
}
