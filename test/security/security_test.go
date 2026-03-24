//go:build security
// +build security

package security

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"live-webrtc-go/internal/api"
	"live-webrtc-go/internal/sfu"
	"live-webrtc-go/internal/testutil"
)

func setupSecurityTest() (*api.HTTPHandlers, string, string) {
	cfg := testutil.TestConfig()
	cfg.AllowedOrigin = "https://example.com"
	cfg.AuthToken = "secure-token"
	cfg.RoomTokens = map[string]string{"secure-room": "room-token"}
	cfg.AdminToken = "admin-token"
	cfg.RateLimitRPS = 5
	cfg.RateLimitBurst = 10
	cfg.JWTSecret = "jwt-secret-key"
	cfg.JWTAudience = "live-webrtc"
	mgr := sfu.NewManager(cfg)
	h := api.NewHTTPHandlers(mgr, cfg)
	return h, cfg.JWTSecret, cfg.JWTAudience
}

func signJWT(t *testing.T, secret string, claims jwt.MapClaims) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("sign jwt: %v", err)
	}
	return s
}

func TestSecurityAuthenticationBypass(t *testing.T) {
	h, _, _ := setupSecurityTest()
	sdpOffer := "v=0\r\no=- 1234567890 1234567890 IN IP4 127.0.0.1\r\ns=-\r\nt=0 0\r\n"

	for _, header := range []struct {
		name  string
		key   string
		value string
	}{{"none", "", ""}, {"wrong token", "X-Auth-Token", "wrong-token"}, {"wrong bearer", "Authorization", "Bearer wrong-token"}} {
		req := httptest.NewRequest("POST", "/api/whip/publish/secure-room", bytes.NewReader([]byte(sdpOffer)))
		if header.key != "" {
			req.Header.Set(header.key, header.value)
		}
		w := httptest.NewRecorder()
		h.ServeWHIPPublish(w, req, "secure-room")
		if w.Result().StatusCode != http.StatusUnauthorized {
			t.Fatalf("%s: expected 401, got %d", header.name, w.Result().StatusCode)
		}
	}
}

func TestSecurityRoomTokenAuthentication(t *testing.T) {
	h, _, _ := setupSecurityTest()
	sdpOffer := "v=0\r\no=- 1234567890 1234567890 IN IP4 127.0.0.1\r\ns=-\r\nt=0 0\r\n"

	req1 := httptest.NewRequest("POST", "/api/whip/publish/secure-room", bytes.NewReader([]byte(sdpOffer)))
	req1.Header.Set("X-Auth-Token", "secure-token")
	w1 := httptest.NewRecorder()
	h.ServeWHIPPublish(w1, req1, "secure-room")
	if w1.Result().StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected room token to override global token")
	}

	req2 := httptest.NewRequest("POST", "/api/whip/publish/secure-room", bytes.NewReader([]byte(sdpOffer)))
	req2.Header.Set("X-Auth-Token", "room-token")
	w2 := httptest.NewRecorder()
	h.ServeWHIPPublish(w2, req2, "secure-room")
	if w2.Result().StatusCode == http.StatusUnauthorized {
		t.Fatal("expected room token to authorize request")
	}
}

func TestSecurityJWTAuthentication(t *testing.T) {
	h, secret, aud := setupSecurityTest()
	sdpOffer := "v=0\r\no=- 1234567890 1234567890 IN IP4 127.0.0.1\r\ns=-\r\nt=0 0\r\n"

	req1 := httptest.NewRequest("POST", "/api/whip/publish/test-room", bytes.NewReader([]byte(sdpOffer)))
	req1.Header.Set("Authorization", "Bearer invalid.jwt.token")
	w1 := httptest.NewRecorder()
	h.ServeWHIPPublish(w1, req1, "test-room")
	if w1.Result().StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected invalid jwt to be rejected")
	}

	validToken := signJWT(t, secret, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Add(-time.Minute).Unix(),
		"aud": aud,
	})
	req2 := httptest.NewRequest("POST", "/api/whip/publish/test-room", bytes.NewReader([]byte(sdpOffer)))
	req2.Header.Set("Authorization", "Bearer "+validToken)
	w2 := httptest.NewRecorder()
	h.ServeWHIPPublish(w2, req2, "test-room")
	if w2.Result().StatusCode == http.StatusUnauthorized {
		t.Fatal("expected valid jwt to authorize request")
	}
}

func TestSecurityAdminAuthentication(t *testing.T) {
	h, _, _ := setupSecurityTest()
	req1 := httptest.NewRequest("POST", "/api/admin/rooms/test-room/close", nil)
	w1 := httptest.NewRecorder()
	h.ServeAdminCloseRoom(w1, req1, "test-room")
	if w1.Result().StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 without admin auth")
	}
}

func TestSecurityRateLimiting(t *testing.T) {
	cfg := testutil.TestConfig()
	cfg.RateLimitRPS = 1
	cfg.RateLimitBurst = 2
	mgr := sfu.NewManager(cfg)
	h := api.NewHTTPHandlers(mgr, cfg)
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("GET", "/api/rooms", nil)
		req.RemoteAddr = "192.168.1.100:12345"
		w := httptest.NewRecorder()
		h.ServeRooms(w, req)
		if i < 2 && w.Result().StatusCode == http.StatusTooManyRequests {
			t.Fatalf("request %d should not be rate limited", i+1)
		}
	}
}

func TestSecurityCORSProtection(t *testing.T) {
	h, _, _ := setupSecurityTest()
	req1 := httptest.NewRequest("OPTIONS", "/api/rooms", nil)
	req1.Header.Set("Origin", "https://example.com")
	w1 := httptest.NewRecorder()
	h.ServeRooms(w1, req1)
	if got := w1.Result().Header.Get("Access-Control-Allow-Origin"); got != "https://example.com" {
		t.Fatalf("expected allowed origin, got %q", got)
	}

	req2 := httptest.NewRequest("OPTIONS", "/api/rooms", nil)
	req2.Header.Set("Origin", "https://malicious.com")
	w2 := httptest.NewRecorder()
	h.ServeRooms(w2, req2)
	if got := w2.Result().Header.Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("expected no origin for disallowed caller, got %q", got)
	}
}

func TestSecurityInputValidation(t *testing.T) {
	for _, room := range []string{"../../../etc/passwd", "room/../../config", "<script>alert(1)</script>", "bad room"} {
		if api.ValidRoomNameForTest(room) {
			t.Fatalf("expected invalid room name: %q", room)
		}
	}
}

func TestSecurityLargePayload(t *testing.T) {
	h, _, _ := setupSecurityTest()
	largeSDP := strings.Repeat("a=large-payload-line\r\n", 100000)
	req := httptest.NewRequest("POST", "/api/whip/publish/secure-room", bytes.NewReader([]byte(largeSDP)))
	req.Header.Set("X-Auth-Token", "room-token")
	w := httptest.NewRecorder()
	h.ServeWHIPPublish(w, req, "secure-room")
	if w.Result().StatusCode != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected 413, got %d", w.Result().StatusCode)
	}
}

func TestSecuritySensitiveDataExposure(t *testing.T) {
	h, _, _ := setupSecurityTest()
	req := httptest.NewRequest("POST", "/api/whip/publish/test-room", bytes.NewReader([]byte("invalid-sdp")))
	w := httptest.NewRecorder()
	h.ServeWHIPPublish(w, req, "test-room")
	body := strings.ToLower(w.Body.String())
	for _, pattern := range []string{"password", "secret", "database"} {
		if strings.Contains(body, pattern) {
			t.Fatalf("response leaked sensitive term %q", pattern)
		}
	}
}
