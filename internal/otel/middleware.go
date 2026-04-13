package otel

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// TraceMiddleware 返回一个 HTTP 中间件，为每个请求创建 OTEL span。
func TraceMiddleware(next http.Handler) http.Handler {
	return otelhttp.NewHandler(next, "live-webrtc-go",
		otelhttp.WithPublicEndpoint(),
		otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
			return r.Method + " " + r.URL.Path
		}),
	)
}
