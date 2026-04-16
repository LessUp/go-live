package sfu

import (
	"sync"

	"github.com/pion/rtp"
	"github.com/pion/webrtc/v3"
	"live-webrtc-go/internal/metrics"
)

const mtuSize = 1500 // Maximum Transmission Unit size for RTP packets

// rtpWriter 抽象录制写入器接口。
type rtpWriter interface {
	WriteRTP(*rtp.Packet) error
	Close() error
}

type subscriberTrack struct {
	sender *webrtc.RTPSender
	local  *webrtc.TrackLocalStaticRTP
}

type subscriberBinding struct {
	pc     *webrtc.PeerConnection
	sender *webrtc.RTPSender
}

// trackFanout 负责把单个远端 Track 分发给多个订阅者，并可选写盘上传。
type trackFanout struct {
	remote  *webrtc.TrackRemote
	mu      sync.RWMutex
	locals  map[*webrtc.PeerConnection]*subscriberTrack
	closed  chan struct{}
	room    string
	rec     rtpWriter
	recPath string
}

func newTrackFanout(remote *webrtc.TrackRemote, room string) *trackFanout {
	return &trackFanout{
		remote: remote,
		locals: make(map[*webrtc.PeerConnection]*subscriberTrack),
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
		buf := make([]byte, mtuSize)
		for {
			if _, _, err := sender.Read(buf); err != nil {
				return
			}
		}
	}()

	f.mu.Lock()
	f.locals[pc] = &subscriberTrack{sender: sender, local: local}
	f.mu.Unlock()
}

func (f *trackFanout) detachFromSubscriber(pc *webrtc.PeerConnection) {
	f.mu.Lock()
	track, ok := f.locals[pc]
	if ok {
		delete(f.locals, pc)
	}
	f.mu.Unlock()
	if ok && track.sender != nil {
		_ = pc.RemoveTrack(track.sender)
	}
}

func (f *trackFanout) snapshotLocals() []*webrtc.TrackLocalStaticRTP {
	f.mu.RLock()
	defer f.mu.RUnlock()
	locals := make([]*webrtc.TrackLocalStaticRTP, 0, len(f.locals))
	for _, track := range f.locals {
		locals = append(locals, track.local)
	}
	return locals
}

// close 关闭 fanout，移除订阅者发送轨道并返回已关闭录制文件路径。
func (f *trackFanout) close() string {
	select {
	case <-f.closed:
		return ""
	default:
		close(f.closed)
	}

	f.mu.Lock()
	bindings := make([]subscriberBinding, 0, len(f.locals))
	for pc, track := range f.locals {
		bindings = append(bindings, subscriberBinding{pc: pc, sender: track.sender})
	}
	f.locals = make(map[*webrtc.PeerConnection]*subscriberTrack)
	path := ""
	if f.rec != nil {
		path = f.recPath
		_ = f.rec.Close()
		f.rec = nil
		f.recPath = ""
	}
	f.mu.Unlock()

	for _, binding := range bindings {
		if binding.pc != nil && binding.sender != nil {
			_ = binding.pc.RemoveTrack(binding.sender)
		}
	}
	return path
}

// readLoop 持续从远端 Track 读取 RTP，并同步写入录制和所有订阅者。
func (f *trackFanout) readLoop() {
	buf := make([]byte, mtuSize)
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
		locals := f.snapshotLocals()
		if rec != nil {
			_ = rec.WriteRTP(pkt)
		}
		for _, local := range locals {
			clone := *pkt
			if pkt.Payload != nil {
				clone.Payload = append([]byte(nil), pkt.Payload...)
			}
			_ = local.WriteRTP(&clone)
		}
	}
}
