package sfu

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/pion/interceptor"
	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media/ivfwriter"
	"github.com/pion/webrtc/v3/pkg/media/oggwriter"
	"live-webrtc-go/internal/metrics"
	"live-webrtc-go/internal/uploader"
)

var uploadRecordingFile = uploader.Upload
var uploaderEnabled = uploader.Enabled

// Room 表示一个 SFU 房间，维护发布者、订阅者与轨道 fanout。
type Room struct {
	name       string
	mu         sync.RWMutex
	publisher  *webrtc.PeerConnection
	trackFeeds map[string]*trackFanout // key: track ID
	subs       map[*webrtc.PeerConnection]struct{}
	mgr        *Manager
}

// NewRoom 初始化房间默认状态。
func NewRoom(name string, m *Manager) *Room {
	return &Room{
		name:       name,
		trackFeeds: make(map[string]*trackFanout),
		subs:       make(map[*webrtc.PeerConnection]struct{}),
		mgr:        m,
	}
}

func (r *Room) empty() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.publisher == nil && len(r.trackFeeds) == 0 && len(r.subs) == 0
}

func (r *Room) syncSubscriberMetricsLocked() {
	metrics.SetSubscribers(r.name, len(r.subs))
}

func (r *Room) trackFeedCountForTest() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.trackFeeds)
}

func (r *Room) subscriberCountForTest() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.subs)
}

func (r *Room) attachTrackFeed(feed *trackFanout, remote *webrtc.TrackRemote) {
	r.mu.Lock()
	r.trackFeeds[remote.ID()] = feed
	subs := make([]*webrtc.PeerConnection, 0, len(r.subs))
	for sub := range r.subs {
		subs = append(subs, sub)
	}
	pub := r.publisher
	r.mu.Unlock()

	for _, sub := range subs {
		feed.attachToSubscriber(sub)
	}

	go feed.readLoop()
	go r.startPLI(remote)

	if pub != nil && r.mgr != nil && r.mgr.cfg != nil && r.mgr.cfg.RecordEnabled {
		r.setupRecording(feed, remote)
	}
}

func (r *Room) startPLI(remote *webrtc.TrackRemote) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		r.mu.RLock()
		pub := r.publisher
		r.mu.RUnlock()
		if pub == nil {
			return
		}
		_ = pub.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: uint32(remote.SSRC())}})
	}
}

func (r *Room) pruneIfEmpty() {
	if r.mgr != nil {
		r.mgr.deleteRoomIfEmpty(r)
	}
}

func (r *Room) closeFeeds(feeds []*trackFanout) {
	for _, feed := range feeds {
		if path := feed.close(); path != "" && uploaderEnabled() {
			go r.uploadRecording(path)
		}
	}
}

func (r *Room) closeSubscriberTracks(feeds []*trackFanout, pc *webrtc.PeerConnection) {
	for _, feed := range feeds {
		feed.detachFromSubscriber(pc)
	}
}

// stats 返回房间的状态快照。
func (r *Room) stats() RoomInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return RoomInfo{
		Name:         r.name,
		HasPublisher: r.publisher != nil,
		Tracks:       len(r.trackFeeds),
		Subscribers:  len(r.subs),
	}
}

// iceConfig 生成 ICE 配置，优先使用配置中的 STUN/TURN。
func (r *Room) iceConfig() webrtc.Configuration {
	if r.mgr != nil && r.mgr.cfg != nil {
		return r.mgr.cfg.ICEConfig()
	}
	return webrtc.Configuration{ICEServers: []webrtc.ICEServer{{URLs: []string{"stun:stun.l.google.com:19302"}}}}
}

// Publish 接收主播的 SDP Offer，创建 PeerConnection 并拉起 track fanout。
func (r *Room) Publish(ctx context.Context, offerSDP string) (string, error) {
	_ = ctx

	m := &webrtc.MediaEngine{}
	if err := m.RegisterDefaultCodecs(); err != nil {
		return "", fmt.Errorf("register codecs: %w", err)
	}
	i := &interceptor.Registry{}
	if err := webrtc.RegisterDefaultInterceptors(m, i); err != nil {
		return "", fmt.Errorf("register interceptors: %w", err)
	}

	api := webrtc.NewAPI(webrtc.WithMediaEngine(m), webrtc.WithInterceptorRegistry(i))
	pc, err := api.NewPeerConnection(r.iceConfig())
	if err != nil {
		return "", err
	}

	r.mu.Lock()
	if r.publisher != nil {
		r.mu.Unlock()
		_ = pc.Close()
		return "", errors.New("publisher already exists in this room")
	}
	r.publisher = pc
	r.mu.Unlock()

	pc.OnICEConnectionStateChange(func(s webrtc.ICEConnectionState) {
		if s == webrtc.ICEConnectionStateFailed || s == webrtc.ICEConnectionStateDisconnected || s == webrtc.ICEConnectionStateClosed {
			go r.closePublisher(pc)
		}
	})

	pc.OnTrack(func(remote *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		_ = receiver
		feed := newTrackFanout(remote, r.name)
		r.attachTrackFeed(feed, remote)
	})

	if err := pc.SetRemoteDescription(webrtc.SessionDescription{Type: webrtc.SDPTypeOffer, SDP: offerSDP}); err != nil {
		r.closePublisher(pc)
		return "", err
	}
	answer, err := pc.CreateAnswer(nil)
	if err != nil {
		r.closePublisher(pc)
		return "", err
	}
	g := webrtc.GatheringCompletePromise(pc)
	if err := pc.SetLocalDescription(answer); err != nil {
		r.closePublisher(pc)
		return "", err
	}
	<-g

	if pc.LocalDescription() == nil {
		r.closePublisher(pc)
		return "", errors.New("missing local description")
	}
	return pc.LocalDescription().SDP, nil
}

// Subscribe 为观众创建 PeerConnection，并把已存在的 track fanout 到新订阅者。
func (r *Room) Subscribe(ctx context.Context, offerSDP string) (string, error) {
	_ = ctx

	if r.mgr != nil && r.mgr.cfg != nil && r.mgr.cfg.MaxSubsPerRoom > 0 {
		r.mu.RLock()
		if len(r.subs) >= r.mgr.cfg.MaxSubsPerRoom {
			r.mu.RUnlock()
			return "", fmt.Errorf("subscriber limit reached")
		}
		r.mu.RUnlock()
	}
	m := &webrtc.MediaEngine{}
	if err := m.RegisterDefaultCodecs(); err != nil {
		return "", fmt.Errorf("register codecs: %w", err)
	}
	i := &interceptor.Registry{}
	if err := webrtc.RegisterDefaultInterceptors(m, i); err != nil {
		return "", fmt.Errorf("register interceptors: %w", err)
	}
	api := webrtc.NewAPI(webrtc.WithMediaEngine(m), webrtc.WithInterceptorRegistry(i))

	pc, err := api.NewPeerConnection(r.iceConfig())
	if err != nil {
		return "", err
	}

	r.mu.Lock()
	r.subs[pc] = struct{}{}
	r.syncSubscriberMetricsLocked()
	feeds := make([]*trackFanout, 0, len(r.trackFeeds))
	for _, feed := range r.trackFeeds {
		feeds = append(feeds, feed)
	}
	r.mu.Unlock()

	for _, feed := range feeds {
		feed.attachToSubscriber(pc)
	}

	pc.OnICEConnectionStateChange(func(s webrtc.ICEConnectionState) {
		if s == webrtc.ICEConnectionStateFailed || s == webrtc.ICEConnectionStateDisconnected || s == webrtc.ICEConnectionStateClosed {
			go r.removeSubscriber(pc)
		}
	})

	if err := pc.SetRemoteDescription(webrtc.SessionDescription{Type: webrtc.SDPTypeOffer, SDP: offerSDP}); err != nil {
		r.removeSubscriber(pc)
		return "", err
	}

	answer, err := pc.CreateAnswer(nil)
	if err != nil {
		r.removeSubscriber(pc)
		return "", err
	}
	g := webrtc.GatheringCompletePromise(pc)
	if err := pc.SetLocalDescription(answer); err != nil {
		r.removeSubscriber(pc)
		return "", err
	}
	<-g

	if pc.LocalDescription() == nil {
		r.removeSubscriber(pc)
		return "", errors.New("missing local description")
	}
	return pc.LocalDescription().SDP, nil
}

// setupRecording 针对音频/视频分别创建 OGG/IVF 写入器做简单录制。
func (r *Room) setupRecording(feed *trackFanout, remote *webrtc.TrackRemote) {
	if r.mgr == nil || r.mgr.cfg == nil {
		return
	}
	// G301: 0o755 is required for WebRTC recording directory access
	// #nosec G301
	if err := os.MkdirAll(r.mgr.cfg.RecordDir, 0o755); err != nil {
		slog.Error("create record dir", "room", r.name, "error", err)
		return
	}
	base := fmt.Sprintf("%s_%s_%d", r.name, remote.ID(), time.Now().Unix())
	mime := remote.Codec().MimeType
	switch mime {
	case webrtc.MimeTypeOpus:
		p := filepath.Join(r.mgr.cfg.RecordDir, base+".ogg")
		if w, err := oggwriter.New(p, 48000, 2); err == nil {
			feed.setRecorder(w, p)
		} else {
			slog.Error("create ogg recorder", "room", r.name, "error", err)
		}
	case webrtc.MimeTypeVP8, webrtc.MimeTypeVP9:
		p := filepath.Join(r.mgr.cfg.RecordDir, base+".ivf")
		if w, err := ivfwriter.New(p); err == nil {
			feed.setRecorder(w, p)
		} else {
			slog.Error("create ivf recorder", "room", r.name, "error", err)
		}
	}
}

// closePublisher 在发布者掉线时清理资源，并断开所有 fanout。
func (r *Room) closePublisher(pc *webrtc.PeerConnection) {
	r.mu.Lock()
	if r.publisher != pc {
		r.mu.Unlock()
		_ = pc.Close()
		return
	}
	feeds := make([]*trackFanout, 0, len(r.trackFeeds))
	for _, f := range r.trackFeeds {
		feeds = append(feeds, f)
	}
	r.trackFeeds = make(map[string]*trackFanout)
	r.publisher = nil
	r.mu.Unlock()

	r.closeFeeds(feeds)
	_ = pc.Close()
	r.pruneIfEmpty()
}

// removeSubscriber 在订阅者离线时解除与 track fanout 的绑定。
func (r *Room) removeSubscriber(pc *webrtc.PeerConnection) {
	r.mu.Lock()
	feeds := make([]*trackFanout, 0, len(r.trackFeeds))
	if _, ok := r.subs[pc]; ok {
		for _, f := range r.trackFeeds {
			feeds = append(feeds, f)
		}
		delete(r.subs, pc)
		r.syncSubscriberMetricsLocked()
		r.mu.Unlock()
		r.closeSubscriberTracks(feeds, pc)
		_ = pc.Close()
		r.pruneIfEmpty()
		return
	}
	r.mu.Unlock()
	_ = pc.Close()
}

// Close 主动关闭房间内所有连接。
func (r *Room) Close() {
	r.mu.Lock()
	pub := r.publisher
	feeds := make([]*trackFanout, 0, len(r.trackFeeds))
	for _, f := range r.trackFeeds {
		feeds = append(feeds, f)
	}
	subs := make([]*webrtc.PeerConnection, 0, len(r.subs))
	for s := range r.subs {
		subs = append(subs, s)
	}
	r.publisher = nil
	r.trackFeeds = make(map[string]*trackFanout)
	r.subs = make(map[*webrtc.PeerConnection]struct{})
	metrics.SetSubscribers(r.name, 0)
	r.mu.Unlock()

	if pub != nil {
		_ = pub.Close()
	}
	r.closeFeeds(feeds)
	for _, s := range subs {
		_ = s.Close()
	}
}

func (r *Room) uploadRecording(path string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := uploadRecordingFile(ctx, path); err != nil {
		slog.Error("upload recording", "room", r.name, "path", path, "error", err)
	}
}
