package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"live-webrtc-go/internal/config"
	"live-webrtc-go/internal/sfu"
)

func TestTokenMatch(t *testing.T) {
	tests := []struct {
		name        string
		headerKey   string
		headerValue string
		expectToken string
		want        bool
	}{
		{
			name:        "X-Auth-Token header match",
			headerKey:   "X-Auth-Token",
			headerValue: "secret",
			expectToken: "secret",
			want:        true,
		},
		{
			name:        "X-Auth-Token header mismatch",
			headerKey:   "X-Auth-Token",
			headerValue: "wrong",
			expectToken: "secret",
			want:        false,
		},
		{
			name:        "Authorization Bearer match",
			headerKey:   "Authorization",
			headerValue: "Bearer secret",
			expectToken: "secret",
			want:        true,
		},
		{
			name:        "Authorization Bearer with spaces",
			headerKey:   "Authorization",
			headerValue: "Bearer  secret  ",
			expectToken: "secret",
			want:        true,
		},
		{
			name:        "missing token",
			headerKey:   "",
			headerValue: "",
			expectToken: "secret",
			want:        false,
		},
		{
			name:        "empty expected token",
			headerKey:   "X-Auth-Token",
			headerValue: "secret",
			expectToken: "",
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.headerKey != "" {
				req.Header.Set(tt.headerKey, tt.headerValue)
			}
			got := tokenMatch(req, tt.expectToken)
			if got != tt.want {
				t.Errorf("tokenMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAllowRate_NoLimiter(t *testing.T) {
	cfg := &config.Config{RateLimitRPS: 0}
	mgr := sfu.NewManager(cfg)
	h := NewHTTPHandlers(mgr, cfg)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "192.168.1.1:12345"

	// Should always allow when no limiter configured
	if !h.allowRate(req) {
		t.Error("expected allowRate to return true when no limiter configured")
	}
}

func TestAuthOKRoom_NoAuth(t *testing.T) {
	cfg := &config.Config{
		AuthToken:  "",
		JWTSecret:  "",
		RoomTokens: map[string]string{},
	}
	mgr := sfu.NewManager(cfg)
	h := NewHTTPHandlers(mgr, cfg)

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	// Should allow when no auth configured
	if !h.authOKRoom(req, "test-room") {
		t.Error("expected authOKRoom to return true when no auth configured")
	}
}

func TestAuthOKRoom_GlobalToken(t *testing.T) {
	cfg := &config.Config{
		AuthToken:  "global-secret",
		JWTSecret:  "",
		RoomTokens: map[string]string{},
	}
	mgr := sfu.NewManager(cfg)
	h := NewHTTPHandlers(mgr, cfg)

	t.Run("valid token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-Auth-Token", "global-secret")
		if !h.authOKRoom(req, "test-room") {
			t.Error("expected authOKRoom to return true with valid token")
		}
	})

	t.Run("invalid token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-Auth-Token", "wrong-token")
		if h.authOKRoom(req, "test-room") {
			t.Error("expected authOKRoom to return false with invalid token")
		}
	})

	t.Run("missing token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		if h.authOKRoom(req, "test-room") {
			t.Error("expected authOKRoom to return false with missing token")
		}
	})
}

func TestAuthOKRoom_RoomToken(t *testing.T) {
	cfg := &config.Config{
		AuthToken: "global-secret",
		JWTSecret: "",
		RoomTokens: map[string]string{
			"special-room": "room-secret",
		},
	}
	mgr := sfu.NewManager(cfg)
	h := NewHTTPHandlers(mgr, cfg)

	t.Run("room-specific token valid", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-Auth-Token", "room-secret")
		if !h.authOKRoom(req, "special-room") {
			t.Error("expected authOKRoom to return true with room-specific token")
		}
	})

	t.Run("room-specific token invalid", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-Auth-Token", "global-secret")
		if h.authOKRoom(req, "special-room") {
			t.Error("expected authOKRoom to return false when room-specific token required but global token provided")
		}
	})

	t.Run("other room uses global token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-Auth-Token", "global-secret")
		if !h.authOKRoom(req, "other-room") {
			t.Error("expected authOKRoom to return true with global token for room without specific token")
		}
	})
}

func TestAdminOK(t *testing.T) {
	cfg := &config.Config{
		AdminToken: "admin-secret",
		JWTSecret:  "",
	}
	mgr := sfu.NewManager(cfg)
	h := NewHTTPHandlers(mgr, cfg)

	t.Run("valid admin token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-Auth-Token", "admin-secret")
		if !h.adminOK(req) {
			t.Error("expected adminOK to return true with valid admin token")
		}
	})

	t.Run("invalid admin token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-Auth-Token", "wrong-token")
		if h.adminOK(req) {
			t.Error("expected adminOK to return false with invalid admin token")
		}
	})

	t.Run("missing admin token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		if h.adminOK(req) {
			t.Error("expected adminOK to return false with missing admin token")
		}
	})
}
