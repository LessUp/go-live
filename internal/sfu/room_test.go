package sfu

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/pion/rtp"
	"github.com/pion/webrtc/v3"
	"live-webrtc-go/internal/testutil"
)

type stubRTPWriter struct {
	closed bool
}

func (w *stubRTPWriter) WriteRTP(_ *rtp.Packet) error { return nil }
func (w *stubRTPWriter) Close() error {
	w.closed = true
	return nil
}

func setupTestManager() *Manager {
	return NewManager(testutil.TestConfig())
}

func TestNewManager(t *testing.T) {
	mgr := setupTestManager()
	if mgr == nil {
		t.Fatal("expected manager")
	}
	if mgr.rooms == nil {
		t.Fatal("expected rooms map")
	}
}

func TestManagerGetOrCreateRoomReturnsSameInstance(t *testing.T) {
	mgr := setupTestManager()
	room1 := mgr.getOrCreateRoom("test-room")
	room2 := mgr.getOrCreateRoom("test-room")
	if room1 != room2 {
		t.Fatal("expected same room instance")
	}
}

func TestManagerPublishInvalidSDPDoesNotLeakRoom(t *testing.T) {
	mgr := setupTestManager()
	_, err := mgr.Publish(context.Background(), "test-room", "invalid-sdp")
	if err == nil {
		t.Fatal("expected invalid SDP error")
	}
	if got := len(mgr.ListRooms()); got != 0 {
		t.Fatalf("expected no rooms after failed publish, got %d", got)
	}
}

func TestManagerSubscribeInvalidSDPDoesNotLeakRoom(t *testing.T) {
	mgr := setupTestManager()
	_, err := mgr.Subscribe(context.Background(), "test-room", "invalid-sdp")
	if err == nil {
		t.Fatal("expected invalid SDP error")
	}
	if got := len(mgr.ListRooms()); got != 0 {
		t.Fatalf("expected no rooms after failed subscribe, got %d", got)
	}
}

func TestManagerCloseRoom(t *testing.T) {
	mgr := setupTestManager()
	mgr.getOrCreateRoom("test-room")
	if !mgr.CloseRoom("test-room") {
		t.Fatal("expected room close to succeed")
	}
	if got := len(mgr.ListRooms()); got != 0 {
		t.Fatalf("expected no rooms after close, got %d", got)
	}
}

func TestRoomStatsReflectState(t *testing.T) {
	mgr := setupTestManager()
	room := mgr.getOrCreateRoom("test-room")
	stats := room.stats()
	if stats.Name != "test-room" {
		t.Fatalf("expected room name test-room, got %s", stats.Name)
	}
	if stats.HasPublisher || stats.Tracks != 0 || stats.Subscribers != 0 {
		t.Fatalf("unexpected initial room stats: %+v", stats)
	}
}

func TestRoomCloseClearsState(t *testing.T) {
	mgr := setupTestManager()
	room := mgr.getOrCreateRoom("test-room")
	room.Close()
	stats := room.stats()
	if stats.HasPublisher || stats.Tracks != 0 || stats.Subscribers != 0 {
		t.Fatalf("expected empty room after close, got %+v", stats)
	}
}

func TestClosePublisherClearsFeedsAndTriggersUpload(t *testing.T) {
	mgr := setupTestManager()
	room := mgr.getOrCreateRoom("test-room")
	pub, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		t.Fatalf("new peer connection: %v", err)
	}
	defer func() { _ = pub.Close() }()
	room.publisher = pub

	recordPath := filepath.Join(t.TempDir(), "record.ivf")
	if err := os.WriteFile(recordPath, []byte("recording"), 0o644); err != nil {
		t.Fatalf("write record: %v", err)
	}
	writer := &stubRTPWriter{}
	feed := &trackFanout{
		locals:  make(map[*webrtc.PeerConnection]*subscriberTrack),
		closed:  make(chan struct{}),
		room:    room.name,
		rec:     writer,
		recPath: recordPath,
	}
	room.trackFeeds["track"] = feed

	called := make(chan string, 1)
	prevUpload := uploadRecordingFile
	prevEnabled := uploaderEnabled
	uploadRecordingFile = func(_ context.Context, path string) error {
		called <- path
		return nil
	}
	uploaderEnabled = func() bool { return true }
	defer func() {
		uploadRecordingFile = prevUpload
		uploaderEnabled = prevEnabled
	}()

	room.closePublisher(pub)

	if room.publisher != nil {
		t.Fatal("expected publisher to be cleared")
	}
	if room.trackFeedCountForTest() != 0 {
		t.Fatalf("expected no feeds after close, got %d", room.trackFeedCountForTest())
	}
	if !writer.closed {
		t.Fatal("expected recorder to be closed")
	}
	select {
	case got := <-called:
		if got != recordPath {
			t.Fatalf("expected upload path %q, got %q", recordPath, got)
		}
	case <-time.After(time.Second):
		t.Fatal("expected upload to be triggered")
	}
}

func TestClosePublisherLeavesFileForFailedUpload(t *testing.T) {
	mgr := setupTestManager()
	room := mgr.getOrCreateRoom("test-room")
	pub, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		t.Fatalf("new peer connection: %v", err)
	}
	defer func() { _ = pub.Close() }()
	room.publisher = pub

	recordPath := filepath.Join(t.TempDir(), "record.ivf")
	if err := os.WriteFile(recordPath, []byte("recording"), 0o644); err != nil {
		t.Fatalf("write record: %v", err)
	}
	feed := &trackFanout{
		locals:  make(map[*webrtc.PeerConnection]*subscriberTrack),
		closed:  make(chan struct{}),
		room:    room.name,
		rec:     &stubRTPWriter{},
		recPath: recordPath,
	}
	room.trackFeeds["track"] = feed

	prevUpload := uploadRecordingFile
	prevEnabled := uploaderEnabled
	uploadDone := make(chan struct{})
	uploadRecordingFile = func(_ context.Context, path string) error {
		defer close(uploadDone)
		return errors.New("upload failed")
	}
	uploaderEnabled = func() bool { return true }
	defer func() {
		<-uploadDone
		uploadRecordingFile = prevUpload
		uploaderEnabled = prevEnabled
	}()

	room.closePublisher(pub)

	<-uploadDone
	if _, err := os.Stat(recordPath); err != nil {
		t.Fatalf("expected local file to remain after failed upload: %v", err)
	}
}

func TestRemoveSubscriberDetachesFeedAndUpdatesState(t *testing.T) {
	mgr := setupTestManager()
	room := mgr.getOrCreateRoom("test-room")
	sub, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		t.Fatalf("new peer connection: %v", err)
	}
	defer func() { _ = sub.Close() }()
	room.subs[sub] = struct{}{}

	feed := &trackFanout{
		locals: make(map[*webrtc.PeerConnection]*subscriberTrack),
		closed: make(chan struct{}),
		room:   room.name,
	}
	feed.locals[sub] = &subscriberTrack{}
	room.trackFeeds["track"] = feed

	room.removeSubscriber(sub)

	if room.subscriberCountForTest() != 0 {
		t.Fatalf("expected no subscribers, got %d", room.subscriberCountForTest())
	}
	if len(feed.locals) != 0 {
		t.Fatalf("expected subscriber track detached, got %d", len(feed.locals))
	}
}
