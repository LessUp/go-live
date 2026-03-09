package sfu

import (
	"context"
	"errors"
	"fmt"
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
)

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
	var servers []webrtc.ICEServer
	if r.mgr != nil && r.mgr.cfg != nil {
		if len(r.mgr.cfg.STUN) > 0 {
			servers = append(servers, webrtc.ICEServer{URLs: r.mgr.cfg.STUN})
		}
		if len(r.mgr.cfg.TURN) > 0 {
			s := webrtc.ICEServer{URLs: r.mgr.cfg.TURN}
			if r.mgr.cfg.TURNUsername != "" || r.mgr.cfg.TURNPassword != "" {
				s.Username = r.mgr.cfg.TURNUsername
				s.Credential = r.mgr.cfg.TURNPassword
				s.CredentialType = webrtc.ICECredentialTypePassword
			}
			servers = append(servers, s)
		}
	}
	if len(servers) == 0 {
		servers = []webrtc.ICEServer{{URLs: []string{"stun:stun.l.google.com:19302"}}}
	}
	return webrtc.Configuration{ICEServers: servers}
}

// Publish 接收主播的 SDP Offer，创建 PeerConnection 并拉起 track fanout。
func (r *Room) Publish(ctx context.Context, offerSDP string) (string, error) {
	r.mu.Lock()
	if r.publisher != nil {
		r.mu.Unlock()
		return "", errors.New("publisher already exists in this room")
	}
	r.mu.Unlock()

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

	pc.OnICEConnectionStateChange(func(s webrtc.ICEConnectionState) {
		if s == webrtc.ICEConnectionStateFailed || s == webrtc.ICEConnectionStateDisconnected || s == webrtc.ICEConnectionStateClosed {
			go r.closePublisher(pc)
		}
	})

	pc.OnTrack(func(remote *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		feed := newTrackFanout(remote, r.name)
		r.mu.Lock()
		r.trackFeeds[remote.ID()] = feed
		// attach existing subscribers
		for sub := range r.subs {
			feed.attachToSubscriber(sub)
		}
		r.mu.Unlock()

		go feed.readLoop()

		go func() {
			// 周期性发送 PLI，提醒发布端刷新关键帧，减轻画面马赛克
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
		}()

		if r.mgr != nil && r.mgr.cfg != nil && r.mgr.cfg.RecordEnabled {
			r.setupRecording(feed, remote)
		}
	})

	if err := pc.SetRemoteDescription(webrtc.SessionDescription{Type: webrtc.SDPTypeOffer, SDP: offerSDP}); err != nil {
		_ = pc.Close()
		return "", err
	}
	answer, err := pc.CreateAnswer(nil)
	if err != nil {
		_ = pc.Close()
		return "", err
	}
	g := webrtc.GatheringCompletePromise(pc)
	if err := pc.SetLocalDescription(answer); err != nil {
		_ = pc.Close()
		return "", err
	}
	<-g

	r.mu.Lock()
	r.publisher = pc
	r.mu.Unlock()

	return pc.LocalDescription().SDP, nil
}

// Subscribe 为观众创建 PeerConnection，并把已存在的 track fanout 到新订阅者。
func (r *Room) Subscribe(ctx context.Context, offerSDP string) (string, error) {
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

	pc.OnICEConnectionStateChange(func(s webrtc.ICEConnectionState) {
		if s == webrtc.ICEConnectionStateFailed || s == webrtc.ICEConnectionStateDisconnected || s == webrtc.ICEConnectionStateClosed {
			go r.removeSubscriber(pc)
		}
	})

	r.mu.RLock()
	for _, feed := range r.trackFeeds {
		feed.attachToSubscriber(pc)
	}
	r.mu.RUnlock()

	if err := pc.SetRemoteDescription(webrtc.SessionDescription{Type: webrtc.SDPTypeOffer, SDP: offerSDP}); err != nil {
		_ = pc.Close()
		return "", err
	}

	answer, err := pc.CreateAnswer(nil)
	if err != nil {
		_ = pc.Close()
		return "", err
	}
	g := webrtc.GatheringCompletePromise(pc)
	if err := pc.SetLocalDescription(answer); err != nil {
		_ = pc.Close()
		return "", err
	}
	<-g

	r.mu.Lock()
	r.subs[pc] = struct{}{}
	r.mu.Unlock()
	metrics.IncSubscribers(r.name)

	return pc.LocalDescription().SDP, nil
}

// setupRecording 针对音频/视频分别创建 OGG/IVF 写入器做简单录制。
func (r *Room) setupRecording(feed *trackFanout, remote *webrtc.TrackRemote) {
	_ = os.MkdirAll(r.mgr.cfg.RecordDir, 0o755)
	base := fmt.Sprintf("%s_%s_%d", r.name, remote.ID(), time.Now().Unix())
	mime := remote.Codec().MimeType
	switch {
	case mime == webrtc.MimeTypeOpus:
		p := filepath.Join(r.mgr.cfg.RecordDir, base+".ogg")
		if w, err := oggwriter.New(p, 48000, 2); err == nil {
			feed.setRecorder(w, p)
		}
	case mime == webrtc.MimeTypeVP8 || mime == webrtc.MimeTypeVP9:
		p := filepath.Join(r.mgr.cfg.RecordDir, base+".ivf")
		if w, err := ivfwriter.New(p); err == nil {
			feed.setRecorder(w, p)
		}
	}
}

// closePublisher 在发布者掉线时清理资源，并断开所有 fanout。
func (r *Room) closePublisher(pc *webrtc.PeerConnection) {
	r.mu.Lock()
	if r.publisher == pc {
		for _, f := range r.trackFeeds {
			f.close()
		}
		r.trackFeeds = make(map[string]*trackFanout)
		r.publisher = nil
	}
	r.mu.Unlock()
	_ = pc.Close()
}

// removeSubscriber 在订阅者离线时解除与 track fanout 的绑定。
func (r *Room) removeSubscriber(pc *webrtc.PeerConnection) {
	r.mu.Lock()
	if _, ok := r.subs[pc]; ok {
		for _, f := range r.trackFeeds {
			f.detachFromSubscriber(pc)
		}
		delete(r.subs, pc)
	}
	r.mu.Unlock()
	_ = pc.Close()
	metrics.DecSubscribers(r.name)
}

// Close 主动关闭房间内所有连接。
func (r *Room) Close() {
	r.mu.Lock()
	pub := r.publisher
	feeds := r.trackFeeds
	subs := r.subs
	r.publisher = nil
	r.trackFeeds = make(map[string]*trackFanout)
	r.subs = make(map[*webrtc.PeerConnection]struct{})
	r.mu.Unlock()

	if pub != nil {
		_ = pub.Close()
	}
	for _, f := range feeds {
		f.close()
	}
	for s := range subs {
		_ = s.Close()
	}
}
