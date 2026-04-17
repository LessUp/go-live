package api

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"

	"live-webrtc-go/internal/config"
	"live-webrtc-go/internal/sfu"
)

func TestValidRoomName(t *testing.T) {
	tests := []struct {
		name  string
		room  string
		valid bool
	}{
		{"simple", "room1", true},
		{"with underscore", "room_1", true},
		{"with dash", "room-1", true},
		{"uppercase", "ROOM1", true},
		{"mixed case", "RoomOne", true},
		{"single char", "r", true},
		{"max length 64", "1234567890123456789012345678901234567890123456789012345678901234", true},
		{"empty", "", false},
		{"too long 65", "12345678901234567890123456789012345678901234567890123456789012345", false},
		{"with space", "room 1", false},
		{"with special char", "room@1", false},
		{"with dot", "room.1", false},
		{"with slash", "room/1", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validRoomName(tt.room)
			if got != tt.valid {
				t.Errorf("validRoomName(%q) = %v, want %v", tt.room, got, tt.valid)
			}
		})
	}
}

func TestRegisterRoutes_Healthz(t *testing.T) {
	cfg := &config.Config{}
	mgr := sfu.NewManager(cfg)
	h := NewHTTPHandlers(mgr, cfg)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux, fstest{}, "records")

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
	if w.Body.String() != "ok" {
		t.Errorf("expected body 'ok', got %q", w.Body.String())
	}
}

func TestRegisterRoutes_Rooms(t *testing.T) {
	cfg := &config.Config{}
	mgr := sfu.NewManager(cfg)
	h := NewHTTPHandlers(mgr, cfg)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux, fstest{}, "records")

	req := httptest.NewRequest(http.MethodGet, "/api/rooms", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestRegisterRoutes_WHIP_InvalidRoom(t *testing.T) {
	cfg := &config.Config{}
	mgr := sfu.NewManager(cfg)
	h := NewHTTPHandlers(mgr, cfg)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux, fstest{}, "records")

	req := httptest.NewRequest(http.MethodPost, "/api/whip/publish/invalid@room", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestRegisterRoutes_WHEP_InvalidRoom(t *testing.T) {
	cfg := &config.Config{}
	mgr := sfu.NewManager(cfg)
	h := NewHTTPHandlers(mgr, cfg)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux, fstest{}, "records")

	req := httptest.NewRequest(http.MethodPost, "/api/whep/play/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestRegisterRoutes_Records_RequiresAuth(t *testing.T) {
	cfg := &config.Config{
		RecordEnabled: true,
		AdminToken:    "admin-secret",
	}
	mgr := sfu.NewManager(cfg)
	h := NewHTTPHandlers(mgr, cfg)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux, fstest{}, "records")

	t.Run("no auth", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/records", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("with auth", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/records", nil)
		req.Header.Set("X-Auth-Token", "admin-secret")
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		// Should return 200 or 500 (if records dir doesn't exist), not 401
		if w.Code == http.StatusUnauthorized {
			t.Errorf("expected status != %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})
}

func TestRegisterRoutes_AdminCloseRoom_RequiresAuth(t *testing.T) {
	cfg := &config.Config{
		AdminToken: "admin-secret",
	}
	mgr := sfu.NewManager(cfg)
	h := NewHTTPHandlers(mgr, cfg)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux, fstest{}, "records")

	t.Run("no auth", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/admin/rooms/test-room/close", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("invalid room name", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/admin/rooms/invalid@room/close", nil)
		req.Header.Set("X-Auth-Token", "admin-secret")
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

// fstest implements fs.FS for testing
type fstest struct{}

func (f fstest) Open(name string) (fs.File, error) {
	return nil, fs.ErrNotExist
}
