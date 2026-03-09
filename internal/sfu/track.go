package sfu

import (
	"context"
	"sync"

	"github.com/pion/rtp"
	"github.com/pion/webrtc/v3"
	"live-webrtc-go/internal/metrics"
	"live-webrtc-go/internal/uploader"
)

// rtpWriter 抽象录制写入器接口。
type rtpWriter interface {
	WriteRTP(*rtp.Packet) error
	Close() error
}

// trackFanout 负责把单个远端 Track 分发给多个订阅者，并可选写盘上传。
type trackFanout struct {
	remote *webrtc.TrackRemote
	mu     sync.RWMutex
	// per-subscriber local tracks
	locals  map[*webrtc.PeerConnection]*webrtc.TrackLocalStaticRTP
	closed  chan struct{}
	room    string
	rec     rtpWriter
	recPath string
}

func newTrackFanout(remote *webrtc.TrackRemote, room string) *trackFanout {
	return &trackFanout{
		remote: remote,
		locals: make(map[*webrtc.PeerConnection]*webrtc.TrackLocalStaticRTP),
		closed: make(chan struct{}),
		room:   room,
	}
}

// setRecorder 设置录制写入器与目标文件路径。
func (f *trackFanout) setRecorder(w rtpWriter, path string) {
	f.mu.Lock()
	f.rec = w
	f.recPath = path
	f.mu.Unlock()
}

// attachToSubscriber 为订阅者创建本地 Track，并启动读取循环以清理发送缓冲。
func (f *trackFanout) attachToSubscriber(pc *webrtc.PeerConnection) {
	codec := f.remote.Codec().RTPCodecCapability
	local, err := webrtc.NewTrackLocalStaticRTP(codec, f.remote.ID(), f.remote.StreamID())
	if err != nil {
		return
	}
	sender, err := pc.AddTrack(local)
	if err != nil {
		return
	}
	go func() {
		buf := make([]byte, 1500)
		for {
			if _, _, err := sender.Read(buf); err != nil {
				return
			}
		}
	}()

	f.mu.Lock()
	f.locals[pc] = local
	f.mu.Unlock()
}

func (f *trackFanout) detachFromSubscriber(pc *webrtc.PeerConnection) {
	f.mu.Lock()
	delete(f.locals, pc)
	f.mu.Unlock()
}

// close 关闭录制文件并触发异步上传。
func (f *trackFanout) close() {
	select {
	case <-f.closed:
		return
	default:
		close(f.closed)
	}
	f.mu.Lock()
	if f.rec != nil {
		_ = f.rec.Close()
		if f.recPath != "" {
			go func(p string) { _ = uploader.Upload(context.Background(), p) }(f.recPath)
		}
		f.rec = nil
		f.recPath = ""
	}
	f.mu.Unlock()
}

// readLoop 持续从远端 Track 读取 RTP，并同步写入录制和所有订阅者。
func (f *trackFanout) readLoop() {
	buf := make([]byte, 1500)
	for {
		select {
		case <-f.closed:
			return
		default:
		}
		n, _, err := f.remote.Read(buf)
		if err != nil {
			return
		}
		metrics.AddBytes(f.room, n)
		metrics.IncPackets(f.room)
		pkt := &rtp.Packet{}
		if err := pkt.Unmarshal(buf[:n]); err != nil {
			continue
		}
		f.mu.RLock()
		rec := f.rec
		f.mu.RUnlock()
		if rec != nil {
			_ = rec.WriteRTP(pkt)
		}
		f.mu.RLock()
		for _, local := range f.locals {
			// clone packet for each subscriber to avoid mutation issues
			clone := *pkt
			if pkt.Payload != nil {
				clone.Payload = append([]byte(nil), pkt.Payload...)
			}
			_ = local.WriteRTP(&clone)
		}
		f.mu.RUnlock()
	}
}
