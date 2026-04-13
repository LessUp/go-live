// Package otel 初始化 OpenTelemetry TracerProvider，支持 stdout 和 OTLP 两种导出模式。
package otel

import (
	"context"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// InitTracer 创建 TracerProvider 并注册为全局 provider。
// 当 OTEL_EXPORTER_OTLP_ENDPOINT 有值时使用 OTLP 导出，否则使用 stdout。
// 返回 shutdown 函数，应在服务退出时调用以刷新 span。
func InitTracer(serviceName string) (func(context.Context) error, error) {
	var exporter sdktrace.SpanExporter
	var err error

	if endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"); endpoint != "" {
		exporter, err = newOTLPExporter(endpoint)
		if err != nil {
			return nil, err
		}
	} else {
		exporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			return nil, err
		}
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter, sdktrace.WithBatchTimeout(time.Second)),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp.Shutdown, nil
}

func newOTLPExporter(endpoint string) (sdktrace.SpanExporter, error) {
	protocol := os.Getenv("OTEL_EXPORTER_OTLP_PROTOCOL")
	if protocol == "grpc" || protocol == "" {
		return otlptracegrpc.New(context.Background(),
			otlptracegrpc.WithEndpoint(endpoint),
		)
	}
	return otlptracehttp.New(context.Background(),
		otlptracehttp.WithEndpoint(endpoint),
	)
}
