// Package sfu 提供轻量级房间与轨道分发的教学实现。
package sfu

import (
	"context"
	"sync"

	"live-webrtc-go/internal/config"
	"live-webrtc-go/internal/metrics"
)

// Manager 负责跟踪所有房间的生命周期，提供 Publish/Subscribe 入口。
type Manager struct {
	mu    sync.RWMutex
	rooms map[string]*Room
	cfg   *config.Config
}

func (m *Manager) roomCountLocked() int {
	return len(m.rooms)
}

func (m *Manager) getRoom(name string) *Room {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.rooms[name]
}

func (m *Manager) ensureRoom(name string) *Room {
	m.mu.Lock()
	defer m.mu.Unlock()
	if r, ok := m.rooms[name]; ok {
		return r
	}
	r := NewRoom(name, m)
	m.rooms[name] = r
	metrics.SetRooms(float64(m.roomCountLocked()))
	return r
}

func (m *Manager) deleteRoomIfEmpty(r *Room) bool {
	if r == nil {
		return false
	}
	if !r.empty() {
		return false
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	current, ok := m.rooms[r.name]
	if !ok || current != r || !r.empty() {
		return false
	}
	delete(m.rooms, r.name)
	metrics.SetRooms(float64(m.roomCountLocked()))
	return true
}

func (m *Manager) pruneRoom(name string) {
	m.mu.RLock()
	r := m.rooms[name]
	m.mu.RUnlock()
	m.deleteRoomIfEmpty(r)
}

// NewManager 创建一个房间管理器。
func NewManager(c *config.Config) *Manager {
	return &Manager{rooms: make(map[string]*Room), cfg: c}
}

// getOrCreateRoom 获取或创建房间，首次创建时更新房间计数指标。
func (m *Manager) getOrCreateRoom(name string) *Room {
	return m.ensureRoom(name)
}

// EnsureRoom 显式创建一个空房间，供管理接口或测试场景使用。
func (m *Manager) EnsureRoom(name string) {
	m.ensureRoom(name)
}

// Publish 根据房间名将 SDP Offer 交给对应 Room 处理，返回 SDP Answer。
func (m *Manager) Publish(ctx context.Context, roomName, offerSDP string) (string, error) {
	r := m.ensureRoom(roomName)
	answer, err := r.Publish(ctx, offerSDP)
	if err != nil {
		m.deleteRoomIfEmpty(r)
		return "", err
	}
	return answer, nil
}

// Subscribe 根据房间名将 SDP Offer 交给对应 Room 处理，返回 SDP Answer。
func (m *Manager) Subscribe(ctx context.Context, roomName, offerSDP string) (string, error) {
	r := m.ensureRoom(roomName)
	answer, err := r.Subscribe(ctx, offerSDP)
	if err != nil {
		m.deleteRoomIfEmpty(r)
		return "", err
	}
	return answer, nil
}

// RoomInfo 房间状态快照，用于对外暴露。
type RoomInfo struct {
	Name         string
	HasPublisher bool
	Tracks       int
	Subscribers  int
}

// ListRooms 返回所有房间的状态快照。
func (m *Manager) ListRooms() []RoomInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]RoomInfo, 0, len(m.rooms))
	for _, r := range m.rooms {
		out = append(out, r.stats())
	}
	return out
}

// CloseRoom 主动关闭指定房间并更新房间数量指标。
func (m *Manager) CloseRoom(name string) bool {
	m.mu.Lock()
	r, ok := m.rooms[name]
	if ok {
		delete(m.rooms, name)
	}
	n := len(m.rooms)
	m.mu.Unlock()
	if ok {
		r.Close()
		metrics.SetRooms(float64(n))
	}
	return ok
}

// CloseAll 在服务退出时关闭所有房间，避免 WebRTC 连接泄漏。
func (m *Manager) CloseAll() {
	m.mu.Lock()
	rooms := make([]*Room, 0, len(m.rooms))
	for _, r := range m.rooms {
		rooms = append(rooms, r)
	}
	m.rooms = make(map[string]*Room)
	m.mu.Unlock()
	for _, r := range rooms {
		r.Close()
	}
	metrics.SetRooms(0)
}
