package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"live-webrtc-go/internal/sfu"
	"live-webrtc-go/internal/testutil"
)

type recordResponse struct {
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	ModTime string `json:"modTime"`
	URL     string `json:"url"`
}

func setupTestHandlers() (*HTTPHandlers, *sfu.Manager) {
	cfg := testutil.TestConfig()
	mgr := sfu.NewManager(cfg)
	h := NewHTTPHandlers(mgr, cfg)
	return h, mgr
}

func signedJWT(t *testing.T, secret string, claims jwt.MapClaims) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("sign jwt: %v", err)
	}
	return signed
}

func decodeRecords(t *testing.T, resp *http.Response) []recordResponse {
	t.Helper()
	defer func() { _ = resp.Body.Close() }()
	var records []recordResponse
	if err := json.NewDecoder(resp.Body).Decode(&records); err != nil {
		t.Fatalf("decode records: %v", err)
	}
	return records
}

func TestServeRoomsSuccess(t *testing.T) {
	h, _ := setupTestHandlers()
	req := httptest.NewRequest(http.MethodGet, "/api/rooms", nil)
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
		t.Fatal("expected array response")
	}
}

func TestServeRoomsMethodAndRateLimit(t *testing.T) {
	cfg := testutil.TestConfig()
	cfg.RateLimitRPS = 1
	cfg.RateLimitBurst = 1
	h := NewHTTPHandlers(sfu.NewManager(cfg), cfg)

	postReq := httptest.NewRequest(http.MethodPost, "/api/rooms", nil)
	postW := httptest.NewRecorder()
	h.ServeRooms(postW, postReq)
	if postW.Result().StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", postW.Result().StatusCode)
	}

	req1 := httptest.NewRequest(http.MethodGet, "/api/rooms", nil)
	req1.RemoteAddr = "127.0.0.1:9001"
	w1 := httptest.NewRecorder()
	h.ServeRooms(w1, req1)
	if w1.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected first request 200, got %d", w1.Result().StatusCode)
	}

	req2 := httptest.NewRequest(http.MethodGet, "/api/rooms", nil)
	req2.RemoteAddr = "127.0.0.1:9001"
	w2 := httptest.NewRecorder()
	h.ServeRooms(w2, req2)
	if w2.Result().StatusCode != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d", w2.Result().StatusCode)
	}
}

func TestServeBootstrapSuccess(t *testing.T) {
	cfg := testutil.TestConfig()
	cfg.RecordEnabled = true
	cfg.JWTSecret = "jwt-secret"
	cfg.TURN = []string{"turn:turn.example.com:3478"}
	cfg.TURNUsername = "turn-user"
	cfg.TURNPassword = "turn-pass"
	h := NewHTTPHandlers(sfu.NewManager(cfg), cfg)
	req := httptest.NewRequest(http.MethodGet, "/api/bootstrap", nil)
	w := httptest.NewRecorder()

	h.ServeBootstrap(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Result().StatusCode)
	}
	var resp bootstrapResponse
	if err := json.NewDecoder(w.Result().Body).Decode(&resp); err != nil {
		t.Fatalf("decode bootstrap: %v", err)
	}
	if !resp.AuthEnabled {
		t.Fatal("expected authEnabled=true")
	}
	if !resp.RecordEnabled {
		t.Fatal("expected recordEnabled=true")
	}
	if !resp.Features["rooms"] || !resp.Features["records"] {
		t.Fatalf("expected rooms and records features, got %+v", resp.Features)
	}
	if len(resp.ICEServers) != 2 {
		t.Fatalf("expected STUN + TURN servers, got %d", len(resp.ICEServers))
	}
}

func TestServeBootstrapMethodNotAllowed(t *testing.T) {
	h, _ := setupTestHandlers()
	req := httptest.NewRequest(http.MethodPost, "/api/bootstrap", nil)
	w := httptest.NewRecorder()

	h.ServeBootstrap(w, req)

	if w.Result().StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", w.Result().StatusCode)
	}
}

func TestServeWHIPPublishRequiresAuth(t *testing.T) {
	cfg := testutil.TestConfig()
	cfg.AuthToken = "required-token"
	h := NewHTTPHandlers(sfu.NewManager(cfg), cfg)
	req := httptest.NewRequest(http.MethodPost, "/api/whip/publish/test-room", strings.NewReader("v=0\r\n"))
	w := httptest.NewRecorder()

	h.ServeWHIPPublish(w, req, "test-room")

	if w.Result().StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Result().StatusCode)
	}
}

func TestServeWHIPPublishRejectsOversizedBody(t *testing.T) {
	h, _ := setupTestHandlers()
	req := httptest.NewRequest(http.MethodPost, "/api/whip/publish/test-room", strings.NewReader(strings.Repeat("a", int(maxSDPBodyBytes)+1)))
	w := httptest.NewRecorder()

	h.ServeWHIPPublish(w, req, "test-room")

	if w.Result().StatusCode != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected 413, got %d", w.Result().StatusCode)
	}
}

func TestServeWHIPPublishRoomJWTAuth(t *testing.T) {
	cfg := testutil.TestConfig()
	cfg.JWTSecret = "jwt-secret"
	cfg.JWTAudience = "live-webrtc"
	h := NewHTTPHandlers(sfu.NewManager(cfg), cfg)
	token := signedJWT(t, cfg.JWTSecret, jwt.MapClaims{
		"room": "test-room",
		"aud":  cfg.JWTAudience,
		"iat":  time.Now().Add(-time.Minute).Unix(),
		"exp":  time.Now().Add(time.Hour).Unix(),
	})
	req := httptest.NewRequest(http.MethodPost, "/api/whip/publish/test-room", strings.NewReader("v=0\r\n"))
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	h.ServeWHIPPublish(w, req, "test-room")

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("expected invalid SDP to reach handler and return 400, got %d", w.Result().StatusCode)
	}
}

func TestServeWHEPPlayRequiresAuthAndRejectsOversizedBody(t *testing.T) {
	cfg := testutil.TestConfig()
	cfg.RoomTokens = map[string]string{"test-room": "room-token"}
	h := NewHTTPHandlers(sfu.NewManager(cfg), cfg)

	unauthorizedReq := httptest.NewRequest(http.MethodPost, "/api/whep/play/test-room", strings.NewReader("v=0\r\n"))
	unauthorizedW := httptest.NewRecorder()
	h.ServeWHEPPlay(unauthorizedW, unauthorizedReq, "test-room")
	if unauthorizedW.Result().StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", unauthorizedW.Result().StatusCode)
	}

	oversizedReq := httptest.NewRequest(http.MethodPost, "/api/whep/play/test-room", strings.NewReader(strings.Repeat("a", int(maxSDPBodyBytes)+1)))
	oversizedW := httptest.NewRecorder()
	h = NewHTTPHandlers(sfu.NewManager(testutil.TestConfig()), testutil.TestConfig())
	h.ServeWHEPPlay(oversizedW, oversizedReq, "test-room")
	if oversizedW.Result().StatusCode != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected 413, got %d", oversizedW.Result().StatusCode)
	}
}

func TestServeRecordsListSortedAndLocalOnly(t *testing.T) {
	cfg := testutil.TestConfig()
	cfg.RecordDir = t.TempDir()
	h := NewHTTPHandlers(sfu.NewManager(cfg), cfg)

	olderPath := filepath.Join(cfg.RecordDir, "older.ivf")
	if err := os.WriteFile(olderPath, []byte("old"), 0o644); err != nil {
		t.Fatalf("write older record: %v", err)
	}
	olderTime := time.Now().Add(-time.Minute)
	if err := os.Chtimes(olderPath, olderTime, olderTime); err != nil {
		t.Fatalf("chtimes older record: %v", err)
	}

	newerPath := filepath.Join(cfg.RecordDir, "newer.ogg")
	if err := os.WriteFile(newerPath, []byte("newer"), 0o644); err != nil {
		t.Fatalf("write newer record: %v", err)
	}
	newerTime := time.Now().Add(time.Minute)
	if err := os.Chtimes(newerPath, newerTime, newerTime); err != nil {
		t.Fatalf("chtimes newer record: %v", err)
	}

	if err := os.WriteFile(filepath.Join(cfg.RecordDir, "ignore.txt"), []byte("x"), 0o644); err != nil {
		t.Fatalf("write ignored file: %v", err)
	}
	if err := os.Mkdir(filepath.Join(cfg.RecordDir, "nested"), 0o755); err != nil {
		t.Fatalf("mkdir nested: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/records", nil)
	w := httptest.NewRecorder()
	h.ServeRecordsList(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Result().StatusCode)
	}
	records := decodeRecords(t, w.Result())
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
	if records[0].Name != "newer.ogg" || records[1].Name != "older.ivf" {
		t.Fatalf("expected newest-first ordering, got %+v", records)
	}
	if records[0].URL != "/records/newer.ogg" {
		t.Fatalf("expected URL /records/newer.ogg, got %q", records[0].URL)
	}
}

func TestServeRecordsListSortsByNameWhenTimesMatch(t *testing.T) {
	cfg := testutil.TestConfig()
	cfg.RecordDir = t.TempDir()
	h := NewHTTPHandlers(sfu.NewManager(cfg), cfg)

	alphaPath := filepath.Join(cfg.RecordDir, "alpha.ivf")
	betaPath := filepath.Join(cfg.RecordDir, "beta.ivf")
	if err := os.WriteFile(alphaPath, []byte("a"), 0o644); err != nil {
		t.Fatalf("write alpha: %v", err)
	}
	if err := os.WriteFile(betaPath, []byte("b"), 0o644); err != nil {
		t.Fatalf("write beta: %v", err)
	}
	stamp := time.Now()
	if err := os.Chtimes(alphaPath, stamp, stamp); err != nil {
		t.Fatalf("chtimes alpha: %v", err)
	}
	if err := os.Chtimes(betaPath, stamp, stamp); err != nil {
		t.Fatalf("chtimes beta: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/records", nil)
	w := httptest.NewRecorder()
	h.ServeRecordsList(w, req)
	records := decodeRecords(t, w.Result())

	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
	if records[0].Name != "alpha.ivf" || records[1].Name != "beta.ivf" {
		t.Fatalf("expected alphabetical tie-break, got %+v", records)
	}
}

func TestServeRecordsListMissingDirAndLocalDeletion(t *testing.T) {
	cfg := testutil.TestConfig()
	cfg.RecordDir = filepath.Join(t.TempDir(), "records")
	h := NewHTTPHandlers(sfu.NewManager(cfg), cfg)

	missingReq := httptest.NewRequest(http.MethodGet, "/api/records", nil)
	missingW := httptest.NewRecorder()
	h.ServeRecordsList(missingW, missingReq)
	if missingW.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected 200 for missing dir, got %d", missingW.Result().StatusCode)
	}
	if got := missingW.Result().Header.Get("Content-Type"); !strings.Contains(got, "application/json") {
		t.Fatalf("expected application/json content type, got %q", got)
	}
	if records := decodeRecords(t, missingW.Result()); len(records) != 0 {
		t.Fatalf("expected empty list for missing dir, got %d", len(records))
	}

	if err := os.MkdirAll(cfg.RecordDir, 0o755); err != nil {
		t.Fatalf("mkdir record dir: %v", err)
	}
	path := filepath.Join(cfg.RecordDir, "local.ivf")
	if err := os.WriteFile(path, []byte("x"), 0o644); err != nil {
		t.Fatalf("write local file: %v", err)
	}
	if err := os.Remove(path); err != nil {
		t.Fatalf("remove local file: %v", err)
	}

	deletedReq := httptest.NewRequest(http.MethodGet, "/api/records", nil)
	deletedW := httptest.NewRecorder()
	h.ServeRecordsList(deletedW, deletedReq)
	if records := decodeRecords(t, deletedW.Result()); len(records) != 0 {
		t.Fatalf("expected deleted local file to disappear from list, got %d", len(records))
	}
}

func TestServeRecordsListMethodRateLimitAndErrorPath(t *testing.T) {
	cfg := testutil.TestConfig()
	cfg.RecordDir = filepath.Join(t.TempDir(), "records-file")
	cfg.RateLimitRPS = 1
	cfg.RateLimitBurst = 1
	if err := os.WriteFile(cfg.RecordDir, []byte("not-a-directory"), 0o644); err != nil {
		t.Fatalf("write records file: %v", err)
	}
	h := NewHTTPHandlers(sfu.NewManager(cfg), cfg)

	postReq := httptest.NewRequest(http.MethodPost, "/api/records", nil)
	postW := httptest.NewRecorder()
	h.ServeRecordsList(postW, postReq)
	if postW.Result().StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", postW.Result().StatusCode)
	}

	errorReq := httptest.NewRequest(http.MethodGet, "/api/records", nil)
	errorReq.RemoteAddr = "127.0.0.1:9002"
	errorW := httptest.NewRecorder()
	h.ServeRecordsList(errorW, errorReq)
	if errorW.Result().StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", errorW.Result().StatusCode)
	}

	rateReq := httptest.NewRequest(http.MethodGet, "/api/records", nil)
	rateReq.RemoteAddr = "127.0.0.1:9002"
	rateW := httptest.NewRecorder()
	h.ServeRecordsList(rateW, rateReq)
	if rateW.Result().StatusCode != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d", rateW.Result().StatusCode)
	}
}

func TestServeAdminCloseRoomUnauthorizedAndSuccess(t *testing.T) {
	cfg := testutil.TestConfig()
	cfg.AdminToken = "admin-token"
	mgr := sfu.NewManager(cfg)
	mgr.EnsureRoom("test-room")
	h := NewHTTPHandlers(mgr, cfg)

	unauthorizedReq := httptest.NewRequest(http.MethodPost, "/api/admin/rooms/test-room/close", nil)
	unauthorizedW := httptest.NewRecorder()
	h.ServeAdminCloseRoom(unauthorizedW, unauthorizedReq, "test-room")
	if unauthorizedW.Result().StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", unauthorizedW.Result().StatusCode)
	}

	authorizedReq := httptest.NewRequest(http.MethodPost, "/api/admin/rooms/test-room/close", nil)
	authorizedReq.Header.Set("Authorization", "Bearer admin-token")
	authorizedW := httptest.NewRecorder()
	h.ServeAdminCloseRoom(authorizedW, authorizedReq, "test-room")
	if authorizedW.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", authorizedW.Result().StatusCode)
	}

	notFoundReq := httptest.NewRequest(http.MethodPost, "/api/admin/rooms/test-room/close", nil)
	notFoundReq.Header.Set("Authorization", "Bearer admin-token")
	notFoundW := httptest.NewRecorder()
	h.ServeAdminCloseRoom(notFoundW, notFoundReq, "test-room")
	if notFoundW.Result().StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", notFoundW.Result().StatusCode)
	}
}

func TestTokenMatchSupportsBearerAndXAuthToken(t *testing.T) {
	bearerReq := httptest.NewRequest(http.MethodGet, "/test", nil)
	bearerReq.Header.Set("Authorization", "Bearer test-token")
	if !tokenMatch(bearerReq, "test-token") {
		t.Fatal("expected bearer token to match")
	}

	xAuthReq := httptest.NewRequest(http.MethodGet, "/test", nil)
	xAuthReq.Header.Set("X-Auth-Token", "test-token")
	if !tokenMatch(xAuthReq, "test-token") {
		t.Fatal("expected X-Auth-Token to match")
	}

	wrongReq := httptest.NewRequest(http.MethodGet, "/test", nil)
	wrongReq.Header.Set("Authorization", "Bearer wrong-token")
	if tokenMatch(wrongReq, "test-token") {
		t.Fatal("expected mismatched token to fail")
	}
}

func TestJWTOKRoomAndJWTAdmin(t *testing.T) {
	secret := "jwt-secret"
	audience := "live-webrtc"

	roomToken := signedJWT(t, secret, jwt.MapClaims{
		"room": "test-room",
		"aud":  audience,
		"iat":  time.Now().Add(-time.Minute).Unix(),
		"exp":  time.Now().Add(time.Hour).Unix(),
	})
	roomReq := httptest.NewRequest(http.MethodGet, "/test", nil)
	roomReq.Header.Set("Authorization", "Bearer "+roomToken)
	if !jwtOKRoom(roomReq, "test-room", secret, audience) {
		t.Fatal("expected room jwt to authorize")
	}
	if jwtOKRoom(roomReq, "other-room", secret, audience) {
		t.Fatal("expected room mismatch to be rejected")
	}

	expiredToken := signedJWT(t, secret, jwt.MapClaims{
		"room": "test-room",
		"aud":  audience,
		"iat":  time.Now().Add(-2 * time.Hour).Unix(),
		"exp":  time.Now().Add(-time.Hour).Unix(),
	})
	expiredReq := httptest.NewRequest(http.MethodGet, "/test", nil)
	expiredReq.Header.Set("Authorization", "Bearer "+expiredToken)
	if jwtOKRoom(expiredReq, "test-room", secret, audience) {
		t.Fatal("expected expired jwt to fail")
	}

	adminToken := signedJWT(t, secret, jwt.MapClaims{
		"role": "admin",
		"aud":  audience,
		"iat":  time.Now().Add(-time.Minute).Unix(),
		"exp":  time.Now().Add(time.Hour).Unix(),
	})
	adminReq := httptest.NewRequest(http.MethodGet, "/test", nil)
	adminReq.Header.Set("Authorization", "Bearer "+adminToken)
	if !jwtAdmin(adminReq, secret, audience) {
		t.Fatal("expected admin jwt to authorize")
	}
}

func TestAllowCORSAndHostMatch(t *testing.T) {
	h, _ := setupTestHandlers()
	wildcardReq := httptest.NewRequest(http.MethodGet, "/test", nil)
	wildcardReq.Header.Set("Origin", "https://example.com")
	wildcardW := httptest.NewRecorder()
	h.allowCORS(wildcardW, wildcardReq)
	if got := wildcardW.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Fatalf("expected wildcard origin, got %q", got)
	}
	if got := wildcardW.Header().Get("Access-Control-Allow-Credentials"); got != "" {
		t.Fatalf("expected no credentials header for wildcard origin, got %q", got)
	}

	cfg := testutil.TestConfig()
	cfg.AllowedOrigin = "example.com"
	h = NewHTTPHandlers(sfu.NewManager(cfg), cfg)
	strictReq := httptest.NewRequest(http.MethodGet, "/test", nil)
	strictReq.Header.Set("Origin", "https://example.com")
	strictW := httptest.NewRecorder()
	h.allowCORS(strictW, strictReq)
	if got := strictW.Header().Get("Access-Control-Allow-Origin"); got != "https://example.com" {
		t.Fatalf("expected reflected origin, got %q", got)
	}
	if got := strictW.Header().Get("Access-Control-Allow-Credentials"); got != "true" {
		t.Fatalf("expected credentials header, got %q", got)
	}
	if !hostMatch("example.com", "https://example.com") {
		t.Fatal("expected hostMatch to accept exact host")
	}
	if hostMatch("example.com", "https://evil.com") {
		t.Fatal("expected hostMatch to reject different host")
	}
}

// TestRoomListJSONContract is a RED test: it asserts that /api/rooms emits
// canonical lower-camel keys ("name", "subscribers") and that the exported Go
// field name "Name" is absent. It will FAIL until RoomInfo carries JSON tags.
func TestRoomListJSONContract(t *testing.T) {
	h, mgr := setupTestHandlers()
	mgr.EnsureRoom("test-room")

	req := httptest.NewRequest(http.MethodGet, "/api/rooms", nil)
	w := httptest.NewRecorder()
	h.ServeRooms(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Result().StatusCode)
	}

	body := w.Body.Bytes()
	if !bytes.Contains(body, []byte(`"name"`)) {
		t.Errorf("expected canonical key \"name\" in response body, got: %s", body)
	}
	if !bytes.Contains(body, []byte(`"subscribers"`)) {
		t.Errorf("expected canonical key \"subscribers\" in response body, got: %s", body)
	}
	if bytes.Contains(body, []byte(`"Name"`)) {
		t.Errorf("expected exported Go field \"Name\" to be absent, got: %s", body)
	}
}

// TestStatusForSFUError is a RED test: it will fail to compile because
// statusForSFUError does not exist yet and neither do sfu.ErrPublisherExists,
// sfu.ErrSubscriberLimitReached, sfu.ErrNoPublisher. Once the sentinel errors
// and the helper are added in the production task this will compile and the
// assertions will drive the correct status mapping.
func TestStatusForSFUError(t *testing.T) {
	tests := []struct {
		err  error
		want int
	}{
		{sfu.ErrPublisherExists, http.StatusConflict},         // 409
		{sfu.ErrSubscriberLimitReached, http.StatusForbidden}, // 403
		{sfu.ErrNoPublisher, http.StatusNotFound},             // 404
		{errors.New("generic error"), http.StatusBadRequest},  // 400
	}
	for _, tt := range tests {
		got := statusForSFUError(tt.err)
		if got != tt.want {
			t.Errorf("statusForSFUError(%v) = %d, want %d", tt.err, got, tt.want)
		}
	}
}

// TestServeWHEPPlay404WhenRoomAbsent is a RED test: asserts that subscribing
// to a room that does not exist in the manager returns 404. Currently the
// manager auto-creates rooms (ensureRoom), so this returns 400 instead.
func TestServeWHEPPlay404WhenRoomAbsent(t *testing.T) {
	h, _ := setupTestHandlers()
	// No room created – manager is empty.
	req := httptest.NewRequest(http.MethodPost, "/api/whep/play/absent-room", strings.NewReader("v=0\r\n"))
	w := httptest.NewRecorder()

	h.ServeWHEPPlay(w, req, "absent-room")

	if w.Result().StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 when room absent from manager, got %d", w.Result().StatusCode)
	}
}

// TestWriteJSONError is a RED test: it will fail to compile because
// writeJSONError does not exist yet. Once the helper is added to handlers.go
// the test asserts the three-point contract:
//   - HTTP status is preserved
//   - Content-Type header includes "application/json"
//   - body decodes to {"error": "<message>"}
func TestWriteJSONError(t *testing.T) {
	w := httptest.NewRecorder()
	writeJSONError(w, http.StatusNotFound, "no active publisher in room")

	resp := w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", resp.StatusCode)
	}
	ct := resp.Header.Get("Content-Type")
	if !strings.Contains(ct, "application/json") {
		t.Fatalf("expected Content-Type to contain application/json, got %q", ct)
	}
	var body struct {
		Error string `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode JSON error body: %v", err)
	}
	if body.Error != "no active publisher in room" {
		t.Fatalf("expected error message %q, got %q", "no active publisher in room", body.Error)
	}
}

// TestServeWHEPPlayJSONErrorWhenRoomAbsent is a RED end-to-end test: it asserts
// that when a subscriber requests a room with no active publisher the handler
// returns all three parts of the approved domain-error contract:
//   - HTTP 404
//   - Content-Type: application/json
//   - body {"error": "no active publisher in room"}
//
// The 404 status code is already correct. This test FAILS because ServeWHEPPlay
// uses http.Error(), which sets Content-Type: text/plain and writes a plain-text
// body instead of application/json with a JSON object.
func TestServeWHEPPlayJSONErrorWhenRoomAbsent(t *testing.T) {
	h, _ := setupTestHandlers()
	// No room created; manager has no active publisher.
	req := httptest.NewRequest(http.MethodPost, "/api/whep/play/absent-room", strings.NewReader("v=0\r\n"))
	w := httptest.NewRecorder()

	h.ServeWHEPPlay(w, req, "absent-room")

	resp := w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
	ct := resp.Header.Get("Content-Type")
	if !strings.Contains(ct, "application/json") {
		t.Fatalf("expected Content-Type application/json, got %q", ct)
	}
	var body struct {
		Error string `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("body is not valid JSON: %v", err)
	}
	const wantMsg = "no active publisher in room"
	if body.Error != wantMsg {
		t.Fatalf("expected JSON error %q, got %q", wantMsg, body.Error)
	}
}
