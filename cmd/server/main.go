// 程序入口：启动一个轻量级 WebRTC 服务，提供 WHIP 推流与 WHEP 播放接口，
// 同时暴露房间/录制查询、Prometheus 指标与健康检查，并内嵌示例网页。
package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"live-webrtc-go/internal/api"
	"live-webrtc-go/internal/config"
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
		log.Fatalf("initialize uploader: %v", err)
	}
	mgr := sfu.NewManager(cfg)
	h := api.NewHTTPHandlers(mgr, cfg)

	mux := http.NewServeMux()
	staticFS, err := fs.Sub(webFS, "web")
	if err != nil {
		log.Fatalf("load embedded web assets: %v", err)
	}
	h.RegisterRoutes(mux, staticFS, cfg.RecordDir)

	addr := cfg.HTTPAddr
	fmt.Printf("Live WebRTC server listening on %s\n", addr)
	fmt.Printf("Open http://localhost%s/web/publisher.html and http://localhost%s/web/player.html\n", addr, addr)

	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
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
			log.Fatal(err)
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
