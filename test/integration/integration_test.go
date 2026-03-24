//go:build integration
// +build integration

package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"live-webrtc-go/internal/api"
	"live-webrtc-go/internal/sfu"
	"live-webrtc-go/internal/testutil"
)

func setupIntegrationTest() (*api.HTTPHandlers, *sfu.Manager) {
	cfg := testutil.TestConfig()
	mgr := sfu.NewManager(cfg)
	h := api.NewHTTPHandlers(mgr, cfg)
	return h, mgr
}

func TestIntegrationRoomLifecycleDoesNotLeakOnFailedPublish(t *testing.T) {
	h, mgr := setupIntegrationTest()
	req := httptest.NewRequest("POST", "/api/whip/publish/test-room", bytes.NewReader([]byte("invalid-sdp")))
	w := httptest.NewRecorder()
	h.ServeWHIPPublish(w, req, "test-room")
	if w.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Result().StatusCode)
	}
	if got := len(mgr.ListRooms()); got != 0 {
		t.Fatalf("expected no leaked rooms, got %d", got)
	}
}

func TestIntegrationAuthentication(t *testing.T) {
	cfg := testutil.TestConfig()
	cfg.AuthToken = "test-auth-token"
	mgr := sfu.NewManager(cfg)
	h := api.NewHTTPHandlers(mgr, cfg)
	sdpOffer := "v=0\r\no=- 1234567890 1234567890 IN IP4 127.0.0.1\r\ns=-\r\nt=0 0\r\n"

	req := httptest.NewRequest("POST", "/api/whip/publish/test-room", bytes.NewReader([]byte(sdpOffer)))
	w := httptest.NewRecorder()
	h.ServeWHIPPublish(w, req, "test-room")
	if w.Result().StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 without auth, got %d", w.Result().StatusCode)
	}

	req2 := httptest.NewRequest("POST", "/api/whip/publish/test-room", bytes.NewReader([]byte(sdpOffer)))
	req2.Header.Set("X-Auth-Token", "test-auth-token")
	w2 := httptest.NewRecorder()
	h.ServeWHIPPublish(w2, req2, "test-room")
	if w2.Result().StatusCode == http.StatusUnauthorized {
		t.Fatal("expected auth to succeed")
	}
}

func TestIntegrationAdminCloseRoom(t *testing.T) {
	cfg := testutil.TestConfig()
	cfg.AdminToken = "admin-token"
	mgr := sfu.NewManager(cfg)
	h := api.NewHTTPHandlers(mgr, cfg)
	mgr.EnsureRoom("test-room")

	req := httptest.NewRequest("POST", "/api/admin/rooms/test-room/close", nil)
	req.Header.Set("Authorization", "Bearer admin-token")
	w := httptest.NewRecorder()
	h.ServeAdminCloseRoom(w, req, "test-room")
	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Result().StatusCode)
	}
	if got := len(mgr.ListRooms()); got != 0 {
		t.Fatalf("expected room to be removed, got %d rooms", got)
	}
}

func TestIntegrationRecordsList(t *testing.T) {
	cfg := testutil.TestConfig()
	cfg.RecordDir = t.TempDir()
	mgr := sfu.NewManager(cfg)
	h := api.NewHTTPHandlers(mgr, cfg)

	for _, file := range []struct {
		name string
		data []byte
	}{{"test1.ivf", []byte("a")}, {"test2.ivf", []byte("b")}, {"test.ogg", []byte("c")}, {"ignore.txt", []byte("d")}} {
		if err := os.WriteFile(cfg.RecordDir+"/"+file.name, file.data, 0o644); err != nil {
			t.Fatalf("write %s: %v", file.name, err)
		}
	}

	req := httptest.NewRequest("GET", "/api/records", nil)
	w := httptest.NewRecorder()
	h.ServeRecordsList(w, req)
	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Result().StatusCode)
	}
	var records []map[string]any
	if err := json.NewDecoder(w.Result().Body).Decode(&records); err != nil {
		t.Fatalf("decode records: %v", err)
	}
	if len(records) != 3 {
		t.Fatalf("expected 3 media files, got %d", len(records))
	}
}

func TestIntegrationRecordsListMissingDirReturnsEmptyArray(t *testing.T) {
	cfg := testutil.TestConfig()
	cfg.RecordDir = t.TempDir() + "/missing"
	mgr := sfu.NewManager(cfg)
	h := api.NewHTTPHandlers(mgr, cfg)
	req := httptest.NewRequest("GET", "/api/records", nil)
	w := httptest.NewRecorder()
	h.ServeRecordsList(w, req)
	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Result().StatusCode)
	}
	var records []map[string]any
	if err := json.NewDecoder(w.Result().Body).Decode(&records); err != nil {
		t.Fatalf("decode records: %v", err)
	}
	if len(records) != 0 {
		t.Fatalf("expected empty array, got %d items", len(records))
	}
}

func TestIntegrationCORS(t *testing.T) {
	cfg := testutil.TestConfig()
	cfg.AllowedOrigin = "https://example.com"
	mgr := sfu.NewManager(cfg)
	h := api.NewHTTPHandlers(mgr, cfg)

	req := httptest.NewRequest("OPTIONS", "/api/rooms", nil)
	req.Header.Set("Origin", "https://example.com")
	w := httptest.NewRecorder()
	h.ServeRooms(w, req)
	if w.Result().Header.Get("Access-Control-Allow-Origin") != "https://example.com" {
		t.Fatalf("expected exact allowed origin")
	}

	req2 := httptest.NewRequest("OPTIONS", "/api/rooms", nil)
	req2.Header.Set("Origin", "https://malicious.com")
	w2 := httptest.NewRecorder()
	h.ServeRooms(w2, req2)
	if got := w2.Result().Header.Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("expected no CORS origin for disallowed origin, got %q", got)
	}
}

func TestIntegrationRateLimiting(t *testing.T) {
	cfg := testutil.TestConfig()
	cfg.RateLimitRPS = 1
	cfg.RateLimitBurst = 1
	mgr := sfu.NewManager(cfg)
	h := api.NewHTTPHandlers(mgr, cfg)
	sdpOffer := "v=0\r\no=- 1234567890 1234567890 IN IP4 127.0.0.1\r\ns=-\r\nt=0 0\r\n"

	req := httptest.NewRequest("POST", "/api/whip/publish/test-room", bytes.NewReader([]byte(sdpOffer)))
	req.RemoteAddr = "127.0.0.1:10001"
	w := httptest.NewRecorder()
	h.ServeWHIPPublish(w, req, "test-room")
	if w.Result().StatusCode == http.StatusTooManyRequests {
		t.Fatal("expected first request to pass limiter")
	}

	req2 := httptest.NewRequest("POST", "/api/whip/publish/test-room", bytes.NewReader([]byte(sdpOffer)))
	req2.RemoteAddr = "127.0.0.1:10001"
	w2 := httptest.NewRecorder()
	h.ServeWHIPPublish(w2, req2, "test-room")
	if w2.Result().StatusCode != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d", w2.Result().StatusCode)
	}

	time.Sleep(1100 * time.Millisecond)
	req3 := httptest.NewRequest("POST", "/api/whip/publish/test-room", bytes.NewReader([]byte(sdpOffer)))
	req3.RemoteAddr = "127.0.0.1:10001"
	w3 := httptest.NewRecorder()
	h.ServeWHIPPublish(w3, req3, "test-room")
	if w3.Result().StatusCode == http.StatusTooManyRequests {
		t.Fatal("expected limiter to refill")
	}
}
