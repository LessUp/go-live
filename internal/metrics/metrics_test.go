package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestMetricsInitialValues(t *testing.T) {
	SetRooms(0)
	SetSubscribers("test-room", 0)
	if got := testutil.ToFloat64(Rooms); got != 0 {
		t.Fatalf("expected initial rooms 0, got %f", got)
	}
	if got := testutil.ToFloat64(RTPBytes.WithLabelValues("test-room")); got != 0 {
		t.Fatalf("expected initial RTP bytes 0, got %f", got)
	}
	if got := testutil.ToFloat64(RTPPackets.WithLabelValues("test-room")); got != 0 {
		t.Fatalf("expected initial RTP packets 0, got %f", got)
	}
	if got := testutil.ToFloat64(Subscribers.WithLabelValues("test-room")); got != 0 {
		t.Fatalf("expected initial subscribers 0, got %f", got)
	}
}

func TestSetRooms(t *testing.T) {
	SetRooms(5)
	if got := testutil.ToFloat64(Rooms); got != 5 {
		t.Fatalf("expected rooms 5, got %f", got)
	}
	SetRooms(0)
}

func TestSetSubscribers(t *testing.T) {
	SetSubscribers("test-room", 3)
	if got := testutil.ToFloat64(Subscribers.WithLabelValues("test-room")); got != 3 {
		t.Fatalf("expected subscribers 3, got %f", got)
	}
	SetSubscribers("test-room", 0)
}

func TestIncAndDecSubscribers(t *testing.T) {
	room := "test-room"
	SetSubscribers(room, 0)
	IncSubscribers(room)
	IncSubscribers(room)
	if got := testutil.ToFloat64(Subscribers.WithLabelValues(room)); got != 2 {
		t.Fatalf("expected subscribers 2, got %f", got)
	}
	DecSubscribers(room)
	if got := testutil.ToFloat64(Subscribers.WithLabelValues(room)); got != 1 {
		t.Fatalf("expected subscribers 1, got %f", got)
	}
	SetSubscribers(room, 0)
}

func TestAddBytes(t *testing.T) {
	room := "bytes-room"
	before := testutil.ToFloat64(RTPBytes.WithLabelValues(room))
	AddBytes(room, 1000)
	AddBytes(room, 500)
	if got := testutil.ToFloat64(RTPBytes.WithLabelValues(room)); got != before+1500 {
		t.Fatalf("expected RTP bytes increase by 1500, got %f (before=%f)", got, before)
	}
}

func TestIncPackets(t *testing.T) {
	room := "packets-room"
	before := testutil.ToFloat64(RTPPackets.WithLabelValues(room))
	IncPackets(room)
	IncPackets(room)
	if got := testutil.ToFloat64(RTPPackets.WithLabelValues(room)); got != before+2 {
		t.Fatalf("expected RTP packets increase by 2, got %f (before=%f)", got, before)
	}
}

func TestMetricsConcurrentAccess(t *testing.T) {
	room := "concurrent-room"
	SetSubscribers(room, 0)
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				IncSubscribers(room)
				AddBytes(room, 100)
				IncPackets(room)
				if j%10 == 0 {
					DecSubscribers(room)
				}
			}
			done <- true
		}()
	}
	for i := 0; i < 10; i++ {
		<-done
	}
	if got := testutil.ToFloat64(Subscribers.WithLabelValues(room)); got < 0 {
		t.Fatalf("subscribers should not be negative: %f", got)
	}
	if got := testutil.ToFloat64(RTPBytes.WithLabelValues(room)); got <= 0 {
		t.Fatalf("bytes should be positive: %f", got)
	}
	if got := testutil.ToFloat64(RTPPackets.WithLabelValues(room)); got <= 0 {
		t.Fatalf("packets should be positive: %f", got)
	}
	SetSubscribers(room, 0)
}
