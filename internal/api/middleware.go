package api

import (
	"net"
	"net/http"
	"strings"

	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/time/rate"
)

// allowCORS 设置基础跨域响应头，适配示例页面与教学演示。
func (h *HTTPHandlers) allowCORS(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	ao := h.cfg.AllowedOrigin
	if ao == "*" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	} else if origin != "" && (ao == origin || hostMatch(ao, origin)) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Vary", "Origin")
	}
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Auth-Token")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

// allowRate 根据请求 IP 进行限流，避免单个客户端耗尽资源。
func (h *HTTPHandlers) allowRate(r *http.Request) bool {
	if h.limiter == nil || h.cfg.RateLimitRPS <= 0 {
		return true
	}
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	if host == "" {
		host = r.RemoteAddr
	}
	h.mu.Lock()
	limiter, ok := h.limiter[host]
	if !ok {
		burst := h.cfg.RateLimitBurst
		if burst <= 0 {
			burst = 1
		}
		limiter = rate.NewLimiter(rate.Limit(h.cfg.RateLimitRPS), burst)
		h.limiter[host] = limiter
	}
	h.mu.Unlock()
	return limiter.Allow()
}

// authOKRoom 校验访问权限：优先房间级 Token，再回退到全局 Token 或 JWT；
// JWT 可包含 room 声明以限制访问到指定房间。
func (h *HTTPHandlers) authOKRoom(r *http.Request, room string) bool {
	// 优先匹配房间级 Token，再回退到全局 Token 或 JWT。
	// room-specific token overrides global config if set
	if tok, ok := h.cfg.RoomTokens[room]; ok && tok != "" {
		if tokenMatch(r, tok) {
			return true
		}
		if h.cfg.JWTSecret != "" && jwtOKRoom(r, room, h.cfg.JWTSecret) {
			return true
		}
		return false
	}
	if h.cfg.AuthToken != "" {
		if tokenMatch(r, h.cfg.AuthToken) {
			return true
		}
		if h.cfg.JWTSecret != "" && jwtOKRoom(r, room, h.cfg.JWTSecret) {
			return true
		}
		return false
	}
	if h.cfg.JWTSecret != "" {
		if jwtOKRoom(r, room, h.cfg.JWTSecret) {
			return true
		}
		return false
	}
	return true
}

// adminOK 校验管理接口调用方，默认使用 ADMIN_TOKEN，也支持 JWT 指定管理员角色。
func (h *HTTPHandlers) adminOK(r *http.Request) bool {
	if h.cfg.AdminToken != "" && tokenMatch(r, h.cfg.AdminToken) {
		return true
	}
	if h.cfg.JWTSecret != "" && jwtAdmin(r, h.cfg.JWTSecret) {
		return true
	}
	return false
}

// tokenMatch 从 X-Auth-Token 或 Authorization: Bearer 中读取并比对令牌。
func tokenMatch(r *http.Request, expect string) bool {
	if t := r.Header.Get("X-Auth-Token"); t != "" {
		return t == expect
	}
	auth := r.Header.Get("Authorization")
	if strings.HasPrefix(strings.ToLower(auth), "bearer ") {
		return strings.TrimSpace(auth[7:]) == expect
	}
	return false
}

// jwtOKRoom 验证 HMAC JWT 并（可选）校验 claims.room 与目标房间一致。
// 为简化演示，不强制验证 exp/iat/aud。
func jwtOKRoom(r *http.Request, room, secret string) bool {
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(strings.ToLower(auth), "bearer ") {
		return false
	}
	tokenString := strings.TrimSpace(auth[7:])
	parsed, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrInvalidKeyType
		}
		return []byte(secret), nil
	})
	if err != nil || !parsed.Valid {
		return false
	}
	if claims, ok := parsed.Claims.(jwt.MapClaims); ok {
		if v, ok := claims["room"].(string); ok && v != "" && v != room {
			return false
		}
	}
	return true
}

// jwtAdmin 验证 HMAC JWT 并判断是否具备管理员权限（role=admin 或 admin=true/1）。
func jwtAdmin(r *http.Request, secret string) bool {
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(strings.ToLower(auth), "bearer ") {
		return false
	}
	tokenString := strings.TrimSpace(auth[7:])
	parsed, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrInvalidKeyType
		}
		return []byte(secret), nil
	})
	if err != nil || !parsed.Valid {
		return false
	}
	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}
	if role, ok := claims["role"].(string); ok && strings.EqualFold(role, "admin") {
		return true
	}
	if adminBool, ok := claims["admin"].(bool); ok && adminBool {
		return true
	}
	if adminNum, ok := claims["admin"].(float64); ok && adminNum == 1 {
		return true
	}
	return false
}

// hostMatch 简单比对来源主机名是否与配置相符。
func hostMatch(expect, origin string) bool {
	u := origin
	if i := strings.Index(origin, "://"); i >= 0 {
		u = origin[i+3:]
	}
	if j := strings.Index(u, "/"); j >= 0 {
		u = u[:j]
	}
	host, _, err := net.SplitHostPort(u)
	if err != nil {
		host = u
	}
	return host == expect || origin == expect
}
