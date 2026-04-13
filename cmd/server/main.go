// 程序入口：启动一个轻量级 WebRTC 服务，提供 WHIP 推流与 WHEP 播放接口，
// 同时暴露房间/录制查询、Prometheus 指标与健康检查，并内嵌示例网页。
package main

import (
	"context"
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"live-webrtc-go/internal/api"
	"live-webrtc-go/internal/config"
	liveotel "live-webrtc-go/internal/otel"
	"live-webrtc-go/internal/sfu"
	"live-webrtc-go/internal/uploader"
)

// web 目录下的静态资源打包进二进制，便于教学演示与单文件部署。
//
//go:embed web
var webFS embed.FS

// main 负责：
// 1) 加载配置并初始化上传器与房间管理器
// 2) 注册 HTTP 路由（WHIP/WHEP/房间/录制/管理/指标/健康检查/静态页面）
// 3) 启动 HTTP/HTTPS 服务并实现优雅退出
func main() {
	cfg := config.Load()
	if err := uploader.Init(cfg); err != nil {
		slog.Error("initialize uploader", "error", err)
		os.Exit(1)
	}
	mgr := sfu.NewManager(cfg)
	h := api.NewHTTPHandlers(mgr, cfg)

	mux := http.NewServeMux()
	staticFS, err := fs.Sub(webFS, "web")
	if err != nil {
		slog.Error("load embedded web assets", "error", err)
		os.Exit(1)
	}
	h.RegisterRoutes(mux, staticFS, cfg.RecordDir)

	// 初始化 OpenTelemetry tracing
	otelShutdown, err := liveotel.InitTracer(cfg.OTELServiceName)
	if err != nil {
		slog.Error("initialize tracer", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := otelShutdown(context.Background()); err != nil {
			slog.Error("tracer shutdown", "error", err)
		}
	}()

	addr := cfg.HTTPAddr
	slog.Info("server starting", "addr", addr)
	slog.Info("publisher page", "url", "http://localhost"+addr+"/web/publisher.html")
	slog.Info("player page", "url", "http://localhost"+addr+"/web/player.html")

	srv := &http.Server{
		Addr:              addr,
		Handler:           liveotel.TraceMiddleware(mux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	go func() {
		var err error
		if cfg.TLSCertFile != "" && cfg.TLSKeyFile != "" {
			err = srv.ListenAndServeTLS(cfg.TLSCertFile, cfg.TLSKeyFile)
		} else {
			err = srv.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	mgr.CloseAll()
}
