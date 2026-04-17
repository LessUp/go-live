package sfu

import (
	"testing"

	"github.com/pion/rtp"
	"github.com/pion/webrtc/v3"
)

func TestTrackFanout_NewFanout(t *testing.T) {
	// Test that trackFanout can be properly initialized
	f := &trackFanout{
		locals: make(map[*webrtc.PeerConnection]*subscriberTrack),
		closed: make(chan struct{}),
		room:   "test-room",
	}

	if f.room != "test-room" {
		t.Errorf("expected room 'test-room', got %q", f.room)
	}
	if f.locals == nil {
		t.Error("expected locals map to be initialized")
	}
	if f.closed == nil {
		t.Error("expected closed channel to be initialized")
	}
}

func TestTrackFanout_CloseIdempotent(t *testing.T) {
	// Test that close is idempotent (can be called multiple times safely)
	f := &trackFanout{
		locals: make(map[*webrtc.PeerConnection]*subscriberTrack),
		closed: make(chan struct{}),
		room:   "test-room",
	}

	// First close should return empty string (no recorder path)
	path := f.close()
	if path != "" {
		t.Errorf("expected empty path, got %q", path)
	}

	// Second close should also return empty string and not panic
	path = f.close()
	if path != "" {
		t.Errorf("expected empty path on second close, got %q", path)
	}
}

func TestTrackFanout_SnapshotLocalsEmpty(t *testing.T) {
	// Test snapshotLocals with empty fanout
	f := &trackFanout{
		locals: make(map[*webrtc.PeerConnection]*subscriberTrack),
		closed: make(chan struct{}),
		room:   "test-room",
	}

	locals := f.snapshotLocals()
	if len(locals) != 0 {
		t.Errorf("expected 0 locals, got %d", len(locals))
	}
}

func TestTrackFanout_SetRecorder(t *testing.T) {
	// Test that setRecorder properly sets recorder and recPath
	f := &trackFanout{
		locals: make(map[*webrtc.PeerConnection]*subscriberTrack),
		closed: make(chan struct{}),
		room:   "test-room",
	}

	mockRec := &mockRTPWriter{}
	f.setRecorder(mockRec, "/path/to/recording.ogg")

	f.mu.RLock()
	rec := f.rec
	path := f.recPath
	f.mu.RUnlock()

	if rec == nil {
		t.Error("expected recorder to be set")
	}
	if path != "/path/to/recording.ogg" {
		t.Errorf("expected path '/path/to/recording.ogg', got %q", path)
	}
}

// mockRTPWriter implements rtpWriter for testing
type mockRTPWriter struct {
	closed bool
}

func (m *mockRTPWriter) WriteRTP(_ *rtp.Packet) error { return nil }
func (m *mockRTPWriter) Close() error {
	m.closed = true
	return nil
}
