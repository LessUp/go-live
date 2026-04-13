package api

import (
	"io/fs"
	"net/http"
	"net/http/pprof"
	"regexp"
	"strings"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var roomNamePattern = regexp.MustCompile(`^[A-Za-z0-9_-]{1,64}$`)

func validRoomName(room string) bool {
	return roomNamePattern.MatchString(room)
}

func ValidRoomNameForTest(room string) bool {
	return validRoomName(room)
}

// RegisterRoutes 将所有 HTTP 路由注册到给定的 ServeMux。
// webFS 为内嵌的静态页面文件系统（已 Sub 到 "web" 子目录），recordDir 为录制文件目录。
func (h *HTTPHandlers) RegisterRoutes(mux *http.ServeMux, webFS fs.FS, recordDir string) {
	// API：WHIP 推流（POST）
	mux.HandleFunc("/api/whip/publish/", func(w http.ResponseWriter, r *http.Request) {
		room := strings.TrimPrefix(r.URL.Path, "/api/whip/publish/")
		if !validRoomName(room) {
			http.Error(w, "invalid room", http.StatusBadRequest)
			return
		}
		h.ServeWHIPPublish(w, r, room)
	})

	// API：WHEP 播放（POST）
	mux.HandleFunc("/api/whep/play/", func(w http.ResponseWriter, r *http.Request) {
		room := strings.TrimPrefix(r.URL.Path, "/api/whep/play/")
		if !validRoomName(room) {
			http.Error(w, "invalid room", http.StatusBadRequest)
			return
		}
		h.ServeWHEPPlay(w, r, room)
	})

	// API：房间列表、录制文件列表与前端运行时配置（GET）
	mux.HandleFunc("/api/rooms", h.ServeRooms)
	mux.HandleFunc("/api/records", h.ServeRecordsList)
	mux.HandleFunc("/api/bootstrap", h.ServeBootstrap)

	// 管理接口：关闭房间（POST /api/admin/rooms/{room}/close）
	mux.HandleFunc("/api/admin/rooms/", func(w http.ResponseWriter, r *http.Request) {
		p := strings.TrimPrefix(r.URL.Path, "/api/admin/rooms/")
		if strings.HasSuffix(p, "/close") {
			room := strings.TrimSuffix(p, "/close")
			room = strings.TrimSuffix(room, "/")
			if !validRoomName(room) {
				http.Error(w, "invalid room", http.StatusBadRequest)
				return
			}
			h.ServeAdminCloseRoom(w, r, room)
			return
		}
		http.NotFound(w, r)
	})

	// 健康检查：用于存活探测与基础监控
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// Prometheus 指标：采集房间数量、订阅者数、RTP 字节/包等
	mux.Handle("/metrics", promhttp.Handler())

	// 录制文件静态服务：仅在启用录制时暴露 RECORD_DIR 下内容
	if h.cfg.RecordEnabled {
		mux.Handle("/records/", http.StripPrefix("/records/", http.FileServer(http.Dir(recordDir))))
	}

	// pprof 调试端点：仅在 PPROF=1 时启用
	if h.cfg.PprofEnabled {
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}

	// 内嵌静态页面：publisher.html / player.html 等示例
	mux.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.FS(webFS))))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/web/index.html", http.StatusFound)
			return
		}
		http.NotFound(w, r)
	})
}
