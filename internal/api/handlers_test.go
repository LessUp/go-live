package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"live-webrtc-go/internal/sfu"
	"live-webrtc-go/internal/testutil"
)

func setupTestHandlers() (*HTTPHandlers, *sfu.Manager) {
	cfg := testutil.TestConfig()
	mgr := sfu.NewManager(cfg)
	h := NewHTTPHandlers(mgr, cfg)
	return h, mgr
}

func signedJWT(t *testing.T, secret string, claims jwt.MapClaims) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("sign jwt: %v", err)
	}
	return s
}

func TestServeRoomsSuccess(t *testing.T) {
	h, _ := setupTestHandlers()
	req := httptest.NewRequest("GET", "/api/rooms", nil)
	w := httptest.NewRecorder()
	h.ServeRooms(w, req)
	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Result().StatusCode)
	}
	var rooms []sfu.RoomInfo
	if err := json.NewDecoder(w.Result().Body).Decode(&rooms); err != nil {
		t.Fatalf("decode rooms: %v", err)
	}
	if rooms == nil {
		t.Fatal("expected array")
	}
}

func TestServeBootstrapSuccess(t *testing.T) {
	cfg := testutil.TestConfig()
	cfg.RecordEnabled = true
	cfg.AuthToken = "bootstrap-token"
	cfg.TURN = []string{"turn:turn.example.com:3478"}
	cfg.TURNUsername = "turn-user"
	cfg.TURNPassword = "turn-pass"
	mgr := sfu.NewManager(cfg)
	h := NewHTTPHandlers(mgr, cfg)
	req := httptest.NewRequest("GET", "/api/bootstrap", nil)
	w := httptest.NewRecorder()
	h.ServeBootstrap(w, req)
	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Result().StatusCode)
	}
}

func TestServeWHIPPublishRequiresAuth(t *testing.T) {
	cfg := testutil.TestConfig()
	cfg.AuthToken = "required-token"
	mgr := sfu.NewManager(cfg)
	h := NewHTTPHandlers(mgr, cfg)
	req := httptest.NewRequest("POST", "/api/whip/publish/test-room", strings.NewReader("v=0\r\n"))
	w := httptest.NewRecorder()
	h.ServeWHIPPublish(w, req, "test-room")
	if w.Result().StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Result().StatusCode)
	}
}

func TestServeWHIPPublishRejectsOversizedBody(t *testing.T) {
	h, _ := setupTestHandlers()
	req := httptest.NewRequest("POST", "/api/whip/publish/test-room", strings.NewReader(strings.Repeat("a", int(maxSDPBodyBytes)+1)))
	w := httptest.NewRecorder()
	h.ServeWHIPPublish(w, req, "test-room")
	if w.Result().StatusCode != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected 413, got %d", w.Result().StatusCode)
	}
}

func TestServeWHEPPlayRejectsOversizedBody(t *testing.T) {
	h, _ := setupTestHandlers()
	req := httptest.NewRequest("POST", "/api/whep/play/test-room", strings.NewReader(strings.Repeat("a", int(maxSDPBodyBytes)+1)))
	w := httptest.NewRecorder()
	h.ServeWHEPPlay(w, req, "test-room")
	if w.Result().StatusCode != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected 413, got %d", w.Result().StatusCode)
	}
}

func TestServeRecordsListSuccess(t *testing.T) {
	cfg := testutil.TestConfig()
	cfg.RecordDir = t.TempDir()
	mgr := sfu.NewManager(cfg)
	h := NewHTTPHandlers(mgr, cfg)
	if err := os.WriteFile(cfg.RecordDir+"/test.ivf", []byte("test"), 0o644); err != nil {
		t.Fatalf("write record: %v", err)
	}
	req := httptest.NewRequest("GET", "/api/records", nil)
	w := httptest.NewRecorder()
	h.ServeRecordsList(w, req)
	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Result().StatusCode)
	}
}

func TestServeRecordsListMissingDirReturnsEmptyList(t *testing.T) {
	cfg := testutil.TestConfig()
	cfg.RecordDir = t.TempDir() + "/missing"
	mgr := sfu.NewManager(cfg)
	h := NewHTTPHandlers(mgr, cfg)
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
		t.Fatalf("expected empty list, got %d", len(records))
	}
}

func TestServeAdminCloseRoomSuccess(t *testing.T) {
	cfg := testutil.TestConfig()
	cfg.AdminToken = "admin-token"
	mgr := sfu.NewManager(cfg)
	mgr.EnsureRoom("test-room")
	h := NewHTTPHandlers(mgr, cfg)
	req := httptest.NewRequest("POST", "/api/admin/rooms/test-room/close", nil)
	req.Header.Set("Authorization", "Bearer admin-token")
	w := httptest.NewRecorder()
	h.ServeAdminCloseRoom(w, req, "test-room")
	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Result().StatusCode)
	}
}

func TestTokenMatch(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	if !tokenMatch(req, "test-token") {
		t.Fatal("expected token match")
	}
}

func TestAllowCORSWildcardDoesNotSetCredentials(t *testing.T) {
	h, _ := setupTestHandlers()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	w := httptest.NewRecorder()
	h.allowCORS(w, req)
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Fatalf("expected wildcard origin, got %q", got)
	}
	if got := w.Header().Get("Access-Control-Allow-Credentials"); got != "" {
		t.Fatalf("expected no credentials header with wildcard origin, got %q", got)
	}
}

func TestJWTOKRoomRequiresValidExpAndAudience(t *testing.T) {
	secret := "jwt-secret"
	aud := "live-webrtc"
	validToken := signedJWT(t, secret, jwt.MapClaims{
		"room": "test-room",
		"exp":  time.Now().Add(time.Hour).Unix(),
		"iat":  time.Now().Add(-time.Minute).Unix(),
		"aud":  aud,
	})
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+validToken)
	if !jwtOKRoom(req, "test-room", secret, aud) {
		t.Fatal("expected valid jwt to authorize room")
	}

	expiredToken := signedJWT(t, secret, jwt.MapClaims{
		"room": "test-room",
		"exp":  time.Now().Add(-time.Hour).Unix(),
		"iat":  time.Now().Add(-2 * time.Hour).Unix(),
		"aud":  aud,
	})
	req2 := httptest.NewRequest("GET", "/test", nil)
	req2.Header.Set("Authorization", "Bearer "+expiredToken)
	if jwtOKRoom(req2, "test-room", secret, aud) {
		t.Fatal("expected expired jwt to be rejected")
	}

	wrongAudToken := signedJWT(t, secret, jwt.MapClaims{
		"room": "test-room",
		"exp":  time.Now().Add(time.Hour).Unix(),
		"iat":  time.Now().Add(-time.Minute).Unix(),
		"aud":  "other-audience",
	})
	req3 := httptest.NewRequest("GET", "/test", nil)
	req3.Header.Set("Authorization", "Bearer "+wrongAudToken)
	if jwtOKRoom(req3, "test-room", secret, aud) {
		t.Fatal("expected wrong audience jwt to be rejected")
	}
}
