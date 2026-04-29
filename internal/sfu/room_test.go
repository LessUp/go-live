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

// TestRoomPublishReturnsSentinelOnDuplicatePublisher is a RED test: it will
// fail to compile because ErrPublisherExists does not exist yet. After the
// sentinel is added, it will fail because Publish returns a plain string error.
func TestRoomPublishReturnsSentinelOnDuplicatePublisher(t *testing.T) {
	mgr := setupTestManager()
	room := mgr.getOrCreateRoom("test-room")

	pub, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		t.Fatalf("new peer connection: %v", err)
	}
	defer func() { _ = pub.Close() }()
	room.publisher = pub // pre-set publisher to simulate duplicate

	_, err = room.Publish(context.Background(), "v=0\r\n")
	if !errors.Is(err, ErrPublisherExists) {
		t.Fatalf("expected ErrPublisherExists for duplicate publisher, got %v", err)
	}
}

// TestRoomSubscribeReturnsSentinelOnSubscriberLimit is a RED test: it will
// fail to compile because ErrSubscriberLimitReached does not exist yet. After
// the sentinel is added, it will fail because Subscribe returns a plain fmt.Errorf.
func TestRoomSubscribeReturnsSentinelOnSubscriberLimit(t *testing.T) {
	cfg := testutil.TestConfig()
	cfg.MaxSubsPerRoom = 1
	mgr := NewManager(cfg)
	room := mgr.getOrCreateRoom("test-room")

	sub, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		t.Fatalf("new peer connection: %v", err)
	}
	defer func() { _ = sub.Close() }()
	// Direct state seeding is intentional: the subscriber-limit guard fires
	// before SDP parsing, so there is no way to reach the limit via a second
	// real Subscribe call in this RED phase (it would fail on SDP first).
	// Using a real PeerConnection (not a nil pointer) keeps the map entry
	// production-representative while still exercising the guard path.
	room.subs[sub] = struct{}{}

	_, err = room.Subscribe(context.Background(), "v=0\r\n")
	if !errors.Is(err, ErrSubscriberLimitReached) {
		t.Fatalf("expected ErrSubscriberLimitReached when limit reached, got %v", err)
	}
}

// TestRoomSubscribeReturnsErrNoPublisherWhenNoPublisher is a RED test:
// it will fail to compile because ErrNoPublisher does not exist yet.
// This covers the "room exists but publisher is nil" WHEP 404 path.
func TestRoomSubscribeReturnsErrNoPublisherWhenNoPublisher(t *testing.T) {
	mgr := setupTestManager()
	room := mgr.getOrCreateRoom("test-room")
	// publisher is nil – no one is streaming

	_, err := room.Subscribe(context.Background(), "v=0\r\n")
	if !errors.Is(err, ErrNoPublisher) {
		t.Fatalf("expected ErrNoPublisher when no publisher, got %v", err)
	}
}

// TestSubscribeLimitRaceCondition is a deterministic regression test for the
// stale-check window bug in Subscribe.
//
// Root cause: Subscribe checked MaxSubsPerRoom under r.mu.RLock(), released the
// lock, created a PeerConnection, then unconditionally inserted it under
// r.mu.Lock(). Concurrent callers could all pass the optimistic check and then
// all insert, exceeding the cap.
//
// Why this test is reliable: testHookBeforeInsert fires inside the write lock,
// before the insertion point. The hook directly injects an external subscriber
// into r.subs (no separate lock needed — the write lock is already held),
// bringing the count up to MaxSubsPerRoom. Subscribe must then detect this in
// the re-check and return ErrSubscriberLimitReached rather than inserting.
func TestSubscribeLimitRaceCondition(t *testing.T) {
	cfg := testutil.TestConfig()
	cfg.MaxSubsPerRoom = 1
	mgr := NewManager(cfg)
	room := mgr.getOrCreateRoom("test-room")

	stubPub, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		t.Fatalf("new peer connection: %v", err)
	}
	defer func() { _ = stubPub.Close() }()
	room.publisher = stubPub

	// externalSub simulates a concurrent Subscribe that already completed
	// insertion before this call acquires the write lock.
	externalSub, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		t.Fatalf("new peer connection: %v", err)
	}
	defer func() { _ = externalSub.Close() }()

	// The hook fires inside r.mu.Lock(), so it can write to r.subs directly.
	room.testHookBeforeInsert = func() {
		room.subs[externalSub] = struct{}{}
	}

	_, err = room.Subscribe(context.Background(), "invalid sdp")

	// With the bug: goroutine inserts unconditionally → no ErrSubscriberLimitReached.
	// With the fix: goroutine rechecks under Lock, finds limit reached, returns error.
	if !errors.Is(err, ErrSubscriberLimitReached) {
		t.Errorf("expected ErrSubscriberLimitReached when limit is reached "+
			"between optimistic check and insertion, got: %v", err)
	}

	// Verify no dangling subscriber was added beyond the injected external stub.
	room.mu.Lock()
	subsAfter := len(room.subs)
	delete(room.subs, externalSub) // clean up external stub
	room.mu.Unlock()
	if subsAfter != 1 {
		t.Errorf("expected exactly 1 subscriber (the external stub) after test, got %d", subsAfter)
	}
}

// TestSubscribePublisherDisconnectRace is a deterministic regression test for
// the publisher-disconnect window bug in Subscribe.
//
// Root cause: Subscribe checked r.publisher != nil under r.mu.RLock(), released
// the lock, then inserted the subscriber unconditionally under r.mu.Lock().
// If the publisher disconnected between the two lock sections, a subscriber
// could be added to a room with no active publisher.
//
// Why this test is reliable: testHookBeforeInsert fires inside the write lock,
// before the insertion point. The hook clears r.publisher directly (the write
// lock is already held), so the re-check immediately after the hook sees nil
// and returns ErrNoPublisher instead of inserting the subscriber.
//
// Expected behaviour after the fix: Subscribe rechecks r.publisher under the
// write lock and returns ErrNoPublisher instead of inserting the subscriber.
func TestSubscribePublisherDisconnectRace(t *testing.T) {
	cfg := testutil.TestConfig()
	mgr := NewManager(cfg)
	room := mgr.getOrCreateRoom("test-room")

	// Seed a non-nil publisher so the optimistic RLock check passes.
	stubPub, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		t.Fatalf("new peer connection: %v", err)
	}
	defer func() { _ = stubPub.Close() }()
	room.publisher = stubPub

	// The hook fires inside r.mu.Lock(), so it can clear the publisher directly.
	room.testHookBeforeInsert = func() {
		room.publisher = nil
	}

	_, err = room.Subscribe(context.Background(), "invalid sdp")
	if !errors.Is(err, ErrNoPublisher) {
		t.Errorf("expected ErrNoPublisher when publisher disconnects mid-Subscribe, got: %v", err)
	}

	// Room must not have a dangling subscriber entry.
	room.mu.RLock()
	subCount := len(room.subs)
	room.mu.RUnlock()
	if subCount != 0 {
		t.Errorf("expected 0 subscribers after publisher-disconnect race, got %d", subCount)
	}
}
