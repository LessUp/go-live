# SFU 核心实现

SFU（选择性转发单元）核心逻辑的详细文档。

## 组件层级

```mermaid
flowchart TB
    M[管理器 Manager] -->|1:N| R1[房间 1]
    M -->|1:N| R2[房间 2]
    M -->|1:N| RN[房间 N]

    R1 -->|1:N| TF1[轨道分发 1]
    R1 -->|1:N| TF2[轨道分发 2]

    TF1 -->|读取| RT1[远程轨道<br/>来自发布者]
    TF1 -->|写入| LT1[本地轨道<br/>到订阅者]
    TF1 -->|写入| REC1[录制器<br/>IVF/OGG]
```

## 管理器 (Manager)

管理器是管理所有房间的顶层组件。

```go
type Manager struct {
    rooms     map[string]*Room    // 房间名 → 房间实例
    roomsMu   sync.RWMutex        // 保护 rooms map
    iceConfig ICEConfig           // ICE 服务器配置
    config    Config              // 服务器配置
}
```

### 关键方法

| 方法 | 用途 |
|------|------|
| `Publish(roomName, sdpOffer)` | 创建房间，建立发布者连接 |
| `Subscribe(roomName, sdpOffer)` | 创建订阅者连接，绑定现有轨道 |
| `CloseRoom(roomName)` | 强制关闭房间 |
| `ListRooms()` | 返回所有活跃房间 |
| `RoomCount()` | 返回活跃房间数量 |

## 房间 (Room)

每个房间代表一个独立的流媒体会话。

```go
type Room struct {
    name         string
    publisher    *webrtc.PeerConnection
    subscribers  map[string]*webrtc.PeerConnection
    trackFeeds   map[uint32]*TrackFanout  // SSRC → TrackFanout
    mu           sync.RWMutex
    config       Config
}
```

### 房间生命周期

```mermaid
stateDiagram-v2
    [*] --> Empty: 房间创建
    Empty --> Active: 发布者连接
    Active --> Recording: 录制启用
    Recording --> Active: 发布者断开
    Active --> Empty: 发布者断开
    Empty --> [*]: 房间清理
```

## 轨道分发 (TrackFanout)

TrackFanout 处理单个媒体轨道的 RTP 包分发。

```go
type TrackFanout struct {
    remoteTrack  *webrtc.TrackRemote
    localTracks  map[string]*webrtc.TrackLocalStaticRTP
    recorder     rtpWriter  // IVF 或 OGG 写入器
    stopChan     chan struct{}
}
```

### RTP 分发流程

```mermaid
sequenceDiagram
    participant P as 发布者
    participant TF as TrackFanout
    participant S1 as 订阅者 1
    participant S2 as 订阅者 2
    participant R as 录制器

    P->>TF: RTP 包
    TF->>TF: ReadRTP()
    
    par 分发
        TF->>S1: WriteRTP()
        TF->>S2: WriteRTP()
        TF->>R: WriteRTP()
    end
```

## 并发模型

- **Manager.roomsMu**：保护 rooms map
- **Room.mu**：保护 publisher、subscribers、trackFeeds
- **TrackFanout**：每个轨道单个 goroutine (readLoop)

### Goroutine 生命周期

```mermaid
flowchart TB
    subgraph Main["主 Goroutine"]
        HTTP[HTTP 服务器]
    end

    subgraph PerRoom["每房间"]
        PCH[发布者 OnTrack 处理器]
        SUB[订阅者处理器]
    end

    subgraph PerTrack["每轨道"]
        RL[readLoop]
    end

    HTTP -->|创建房间| PCH
    PCH -->|OnTrack| RL
    RL -->|WriteRTP| SUB
```

## 下一步

- [数据流](/zh/architecture/data-flow) - 完整请求流程图
- [WHIP 协议](/zh/protocols/whip) - WHIP 发布详情
- [WHEP 协议](/zh/protocols/whep) - WHEP 播放详情
