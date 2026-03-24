package sfu

import (
	"context"
	"testing"

	"live-webrtc-go/internal/testutil"
)

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
