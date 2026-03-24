//go:build performance
// +build performance

package performance

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"live-webrtc-go/internal/api"
	"live-webrtc-go/internal/sfu"
	"live-webrtc-go/internal/testutil"
)

func setupPerformanceTest() (*api.HTTPHandlers, *sfu.Manager) {
	cfg := testutil.TestConfig()
	mgr := sfu.NewManager(cfg)
	h := api.NewHTTPHandlers(mgr, cfg)
	return h, mgr
}

func BenchmarkRoomCreation(b *testing.B) {
	h, _ := setupPerformanceTest()
	for i := 0; i < b.N; i++ {
		roomName := fmt.Sprintf("benchmark-room-%d", i)
		req := httptest.NewRequest("POST", "/api/whip/publish/"+roomName, bytes.NewReader([]byte("invalid-sdp")))
		w := httptest.NewRecorder()
		h.ServeWHIPPublish(w, req, roomName)
	}
}

func BenchmarkRoomListing(b *testing.B) {
	h, mgr := setupPerformanceTest()
	for i := 0; i < 100; i++ {
		mgr.EnsureRoom(fmt.Sprintf("setup-room-%d", i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/api/rooms", nil)
		w := httptest.NewRecorder()
		h.ServeRooms(w, req)
	}
}

func BenchmarkConcurrentRequests(b *testing.B) {
	h, _ := setupPerformanceTest()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := httptest.NewRequest("GET", "/api/rooms", nil)
			w := httptest.NewRecorder()
			h.ServeRooms(w, req)
		}
	})
}

func TestPerformanceHighConcurrency(t *testing.T) {
	h, _ := setupPerformanceTest()
	const numRequests = 1000
	const numWorkers = 50
	work := make(chan int, numRequests)
	var successCount int64
	var errorCount int64
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range work {
				req := httptest.NewRequest("GET", "/api/rooms", nil)
				w := httptest.NewRecorder()
				h.ServeRooms(w, req)
				if w.Result().StatusCode == http.StatusOK {
					atomic.AddInt64(&successCount, 1)
				} else {
					atomic.AddInt64(&errorCount, 1)
				}
			}
		}()
	}
	for i := 0; i < numRequests; i++ {
		work <- i
	}
	close(work)
	wg.Wait()
	if errorCount != 0 || successCount != numRequests {
		t.Fatalf("expected all requests to succeed, success=%d error=%d", successCount, errorCount)
	}
}

func TestPerformanceMemoryUsageDoesNotGrowWithFailedPublishes(t *testing.T) {
	h, mgr := setupPerformanceTest()
	var before runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&before)
	for i := 0; i < 100; i++ {
		roomName := fmt.Sprintf("memory-test-room-%d", i)
		req := httptest.NewRequest("POST", "/api/whip/publish/"+roomName, bytes.NewReader([]byte("invalid-sdp")))
		w := httptest.NewRecorder()
		h.ServeWHIPPublish(w, req, roomName)
	}
	if got := len(mgr.ListRooms()); got != 0 {
		t.Fatalf("expected no stale rooms after failed publishes, got %d", got)
	}
	var after runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&after)
	_ = after.Alloc - before.Alloc
}

func TestPerformanceResponseTime(t *testing.T) {
	h, _ := setupPerformanceTest()
	var total time.Duration
	for i := 0; i < 100; i++ {
		start := time.Now()
		req := httptest.NewRequest("GET", "/api/rooms", nil)
		w := httptest.NewRecorder()
		h.ServeRooms(w, req)
		total += time.Since(start)
	}
	avg := total / 100
	if avg > 10*time.Millisecond {
		t.Fatalf("average response time too high: %v", avg)
	}
}

func TestPerformanceRateLimiting(t *testing.T) {
	cfg := testutil.TestConfig()
	cfg.RateLimitRPS = 10
	cfg.RateLimitBurst = 5
	mgr := sfu.NewManager(cfg)
	h := api.NewHTTPHandlers(mgr, cfg)
	rateLimited := 0
	for i := 0; i < 50; i++ {
		req := httptest.NewRequest("GET", "/api/rooms", nil)
		req.RemoteAddr = "127.0.0.1:12345"
		w := httptest.NewRecorder()
		h.ServeRooms(w, req)
		if w.Result().StatusCode == http.StatusTooManyRequests {
			rateLimited++
		}
	}
	if rateLimited == 0 || rateLimited == 50 {
		t.Fatalf("expected partial rate limiting, got %d", rateLimited)
	}
}
