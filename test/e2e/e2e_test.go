//go:build e2e
// +build e2e

package e2e

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

type TestServer struct {
	baseURL string
	client  *http.Client
}

func NewTestServer(baseURL string) *TestServer {
	return &TestServer{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

func (ts *TestServer) request(method, path string, body []byte, headers map[string]string) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}
	req, err := http.NewRequest(method, ts.baseURL+path, bodyReader)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/sdp")
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return ts.client.Do(req)
}

func serverURL() string {
	if url := os.Getenv("TEST_SERVER_URL"); url != "" {
		return url
	}
	return "http://localhost:8080"
}

func TestE2EServerStartup(t *testing.T) {
	ts := NewTestServer(serverURL())
	resp, err := ts.request("GET", "/healthz", nil, nil)
	if err != nil {
		t.Fatalf("connect server: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if string(body) != "ok" {
		t.Fatalf("expected ok body, got %q", string(body))
	}
}

func TestE2ERoomsAPI(t *testing.T) {
	ts := NewTestServer(serverURL())
	resp, err := ts.request("GET", "/api/rooms", nil, nil)
	if err != nil {
		t.Fatalf("request rooms: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var rooms []map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&rooms); err != nil {
		t.Fatalf("decode rooms: %v", err)
	}
}

func TestE2ERecordsAPI(t *testing.T) {
	ts := NewTestServer(serverURL())
	resp, err := ts.request("GET", "/api/records", nil, nil)
	if err != nil {
		t.Fatalf("request records: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var records []map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&records); err != nil {
		t.Fatalf("decode records: %v", err)
	}
}

func TestE2EMetricsEndpoint(t *testing.T) {
	ts := NewTestServer(serverURL())
	resp, err := ts.request("GET", "/metrics", nil, nil)
	if err != nil {
		t.Fatalf("request metrics: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/plain") {
		t.Fatalf("expected text/plain content type, got %q", contentType)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read metrics: %v", err)
	}
	for _, metric := range []string{"webrtc_rooms", "go_gc_duration_seconds", "process_cpu_seconds_total"} {
		if !strings.Contains(string(body), metric) {
			t.Fatalf("expected metric %q in output", metric)
		}
	}
}

func TestE2EStaticFiles(t *testing.T) {
	ts := NewTestServer(serverURL())
	for _, path := range []string{"/web/index.html", "/web/publisher.html", "/web/player.html"} {
		resp, err := ts.request("GET", path, nil, nil)
		if err != nil {
			t.Fatalf("request %s: %v", path, err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200 for %s, got %d", path, resp.StatusCode)
		}
	}
}

func TestE2ECORSHeaders(t *testing.T) {
	ts := NewTestServer(serverURL())
	resp, err := ts.request("OPTIONS", "/api/rooms", nil, map[string]string{"Origin": "http://localhost:3000"})
	if err != nil {
		t.Fatalf("cors request: %v", err)
	}
	defer resp.Body.Close()
	if resp.Header.Get("Access-Control-Allow-Origin") == "" {
		t.Fatal("expected CORS origin header")
	}
}

func TestE2EWebRTCPublishSubscribeRejectsInvalidSDP(t *testing.T) {
	ts := NewTestServer(serverURL())
	sdpOffer := []byte("v=0\r\no=- 1234567890 1234567890 IN IP4 127.0.0.1\r\ns=-\r\nt=0 0\r\n")
	resp, err := ts.request("POST", "/api/whip/publish/test-room", sdpOffer, nil)
	if err != nil {
		t.Fatalf("publish request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid publish SDP, got %d", resp.StatusCode)
	}
	resp2, err := ts.request("POST", "/api/whep/play/test-room", sdpOffer, nil)
	if err != nil {
		t.Fatalf("play request: %v", err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid subscribe SDP, got %d", resp2.StatusCode)
	}
}
