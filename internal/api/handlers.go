// Package api 提供 HTTP 层路由与横切逻辑：CORS、限流、鉴权与业务接口。
package api

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
	"live-webrtc-go/internal/config"
	"live-webrtc-go/internal/sfu"
)

const maxSDPBodyBytes int64 = 1 << 20

// HTTPHandlers 聚合了房间管理器与配置，负责对外暴露 WHIP/WHEP/管理等 API。
type HTTPHandlers struct {
	mgr             *sfu.Manager
	cfg             *config.Config
	mu              sync.Mutex
	limiter         map[string]*rate.Limiter // per-IP 限流器
	limiterLastSeen map[string]time.Time     // per-IP 最后访问时间
	limiterDone     chan struct{}            // 停止清理 goroutine
}

type bootstrapResponse struct {
	AuthEnabled   bool             `json:"authEnabled"`
	RecordEnabled bool             `json:"recordEnabled"`
	ICEServers    []map[string]any `json:"iceServers"`
	Features      map[string]bool  `json:"features"`
}

// NewHTTPHandlers 组合房间管理器与配置，并在启用速率限制时初始化每 IP 的限流器。
func NewHTTPHandlers(m *sfu.Manager, c *config.Config) *HTTPHandlers {
	h := &HTTPHandlers{mgr: m, cfg: c}
	if c.RateLimitRPS > 0 {
		h.limiter = make(map[string]*rate.Limiter)
		h.limiterLastSeen = make(map[string]time.Time)
		h.limiterDone = make(chan struct{})
		go h.limiterGC()
	}
	return h
}

// Close stops background goroutines and releases resources.
// Should be called during server shutdown.
func (h *HTTPHandlers) Close() {
	if h.limiterDone != nil {
		close(h.limiterDone)
	}
}

func (h *HTTPHandlers) readSDPBody(w http.ResponseWriter, r *http.Request) (string, bool) {
	defer func() { _ = r.Body.Close() }()
	limited := http.MaxBytesReader(w, r.Body, maxSDPBodyBytes)
	body, err := io.ReadAll(limited)
	if err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)
			return "", false
		}
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return "", false
	}
	return string(body), true
}

// ServeRooms handles GET /api/rooms
func (h *HTTPHandlers) ServeRooms(w http.ResponseWriter, r *http.Request) {
	h.allowCORS(w, r)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !h.allowRate(r) {
		http.Error(w, "too many requests", http.StatusTooManyRequests)
		return
	}
	rooms := h.mgr.ListRooms()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rooms); err != nil {
		slog.Error("encode rooms response", "error", err)
	}
}

// ServeBootstrap 返回浏览器页面需要的公开运行时配置。
func (h *HTTPHandlers) ServeBootstrap(w http.ResponseWriter, r *http.Request) {
	h.allowCORS(w, r)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !h.allowRate(r) {
		http.Error(w, "too many requests", http.StatusTooManyRequests)
		return
	}

	cfg := h.cfg.ICEConfig()
	iceServers := make([]map[string]any, 0, len(cfg.ICEServers))
	for _, server := range cfg.ICEServers {
		item := map[string]any{"urls": server.URLs}
		if server.Username != "" {
			item["username"] = server.Username
		}
		if credential, ok := server.Credential.(string); ok && credential != "" {
			item["credential"] = credential
		}
		iceServers = append(iceServers, item)
	}

	resp := bootstrapResponse{
		AuthEnabled:   h.cfg.AuthToken != "" || len(h.cfg.RoomTokens) > 0 || h.cfg.JWTSecret != "",
		RecordEnabled: h.cfg.RecordEnabled,
		ICEServers:    iceServers,
		Features: map[string]bool{
			"rooms":   true,
			"records": true,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("encode bootstrap response", "error", err)
	}
}

// ServeWHIPPublish 处理 WHIP 推流：POST /api/whip/publish/{room}
// 请求体为 SDP Offer，返回 SDP Answer（201 Created）。
func (h *HTTPHandlers) ServeWHIPPublish(w http.ResponseWriter, r *http.Request, room string) {
	h.allowCORS(w, r)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !h.allowRate(r) {
		http.Error(w, "too many requests", http.StatusTooManyRequests)
		return
	}
	if !h.authOKRoom(r, room) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	offerSDP, ok := h.readSDPBody(w, r)
	if !ok {
		return
	}
	answer, err := h.mgr.Publish(r.Context(), room, offerSDP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/sdp")
	w.WriteHeader(http.StatusCreated)
	// G705: SDP is machine-generated content with fixed Content-Type, not user-rendered HTML
	// #nosec G705
	_, _ = w.Write([]byte(answer))
}

// ServeWHEPPlay 处理 WHEP 播放：POST /api/whep/play/{room}
// 请求体为 SDP Offer，返回 SDP Answer（201 Created）。
func (h *HTTPHandlers) ServeWHEPPlay(w http.ResponseWriter, r *http.Request, room string) {
	h.allowCORS(w, r)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !h.allowRate(r) {
		http.Error(w, "too many requests", http.StatusTooManyRequests)
		return
	}
	if !h.authOKRoom(r, room) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	offerSDP, ok := h.readSDPBody(w, r)
	if !ok {
		return
	}
	answer, err := h.mgr.Subscribe(r.Context(), room, offerSDP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/sdp")
	w.WriteHeader(http.StatusCreated)
	// G705: SDP is machine-generated content with fixed Content-Type, not user-rendered HTML
	// #nosec G705
	_, _ = w.Write([]byte(answer))
}

// ServeRecordsList 列出 RECORD_DIR 下的 ivf/ogg 文件并返回元数据。
func (h *HTTPHandlers) ServeRecordsList(w http.ResponseWriter, r *http.Request) {
	h.allowCORS(w, r)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !h.allowRate(r) {
		http.Error(w, "too many requests", http.StatusTooManyRequests)
		return
	}
	dir := h.cfg.RecordDir
	entries, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode([]any{}); err != nil {
				slog.Error("encode empty records response", "error", err)
			}
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	type rec struct {
		Name    string `json:"name"`
		Size    int64  `json:"size"`
		ModTime string `json:"modTime"`
		URL     string `json:"url"`
	}
	var list []rec
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		ext := strings.ToLower(filepath.Ext(name))
		if ext != ".ivf" && ext != ".ogg" {
			continue
		}
		fi, err := e.Info()
		if err != nil {
			continue
		}
		list = append(list, rec{
			Name:    name,
			Size:    fi.Size(),
			ModTime: fi.ModTime().UTC().Format(time.RFC3339),
			URL:     "/records/" + name,
		})
	}
	sort.Slice(list, func(i, j int) bool {
		if list[i].ModTime == list[j].ModTime {
			return list[i].Name < list[j].Name
		}
		return list[i].ModTime > list[j].ModTime
	})
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(list); err != nil {
		slog.Error("encode records response", "error", err)
	}
}

// ServeAdminCloseRoom 管理接口：关闭指定房间，释放资源并返回 200。
func (h *HTTPHandlers) ServeAdminCloseRoom(w http.ResponseWriter, r *http.Request, room string) {
	h.allowCORS(w, r)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !h.adminOK(r) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	ok := h.mgr.CloseRoom(room)
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
