---
layout: default
title: 设计说明
nav_order: 3
lang: zh
---

# 设计说明

本文档详细描述 live-webrtc-go 的系统架构、模块拆分、数据流向与扩展点，便于二次开发或架构评审。

{: .no_toc }

## 目录

{: .no_toc .text-delta }

1. TOC
{:toc}

---

## 系统架构

### 架构总览

```
                           ┌─────────────────────────────────────┐
                           │          HTTP Server :8080          │
                           │                                     │
    ┌──────────┐           │  ┌─────────┐    ┌─────────────┐    │
    │  推流端   │ ──WHIP──▶ │  │  认证    │───▶│   SFU       │    │
    │(OBS/网页) │           │  │ 中间件   │    │  管理器      │    │
    └──────────┘           │  └─────────┘    └──────┬──────┘    │
                           │                         │           │
    ┌──────────┐           │         ┌───────────────┤           │
    │  观看端   │ ◀──WHEP── │         │               │           │
    │ (浏览器)  │ ───────▶  │         ▼               ▼           │
    └──────────┘           │  ┌─────────────┐ ┌───────────┐      │
                           │  │    房间      │ │   录制    │      │
                           │  │  (转发)      │ │  与上传   │      │
                           │  └──────┬──────┘ └─────┬─────┘      │
                           └─────────┼──────────────┼────────────┘
                                     │              │
                           ┌─────────▼──────────────▼───────────┐
                           │          对象存储                   │
                           │         (S3/MinIO)                 │
                           └────────────────────────────────────┘
```

### 请求处理链

```
HTTP 请求
    │
    ▼
┌─────────────┐
│    CORS     │ ← ALLOWED_ORIGIN
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   限流器     │ ← RATE_LIMIT_RPS, RATE_LIMIT_BURST
└──────┬──────┘
       │
       ▼
┌─────────────┐
│    认证      │ ← AUTH_TOKEN / ROOM_TOKENS / JWT_SECRET
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   处理器     │ → 业务逻辑
└─────────────┘
```

---

## 核心概念

### 房间（Room）

房间是 SFU 的核心抽象，每个房间：
- 最多一个发布者（Publisher）
- 可有多个订阅者（Subscriber）
- 拥有独立的轨道转发逻辑
- 可配置独立的认证令牌

### 轨道转发（Track Fanout）

当发布者推送媒体轨道时，系统创建轨道转发：
- 从发布者 PeerConnection 读取 RTP 包
- 复制并分发给所有订阅者
- 可选写入录制文件

### 对等连接（PeerConnection）

每个 WebRTC 连接：
- 发布者：接收媒体轨道
- 订阅者：发送媒体轨道
- ICE 协商通过 WHIP/WHEP 协议完成

---

## 模块详解

### cmd/server

**职责**：程序入口，服务初始化

```go
// main.go 主要流程
1. config.Load()           // 加载配置
2. uploader.Init()         // 初始化上传器
3. sfu.NewManager()        // 创建房间管理器
4. api.NewHTTPHandlers()   // 创建 HTTP 处理器
5. RegisterRoutes()        // 注册路由
6. otel.InitTracer()       // 初始化追踪
7. http.Server.Listen()    // 启动服务
8. 优雅退出                 // 优雅关闭
```

### internal/config

**职责**：环境变量解析与默认值

```
┌─────────────────────────────────────┐
│             Config                  │
├─────────────────────────────────────┤
│ HTTPAddr        string              │
│ AllowedOrigin   string              │
│ AuthToken       string              │
│ RoomTokens      map[string]string   │
│ JWTSecret       string              │
│ RecordEnabled   bool                │
│ RecordDir       string              │
│ S3Endpoint      string              │
│ RateLimitRPS    float64             │
│ STUN/TURN       []string            │
│ ...                                 │
└─────────────────────────────────────┘
```

### internal/api

**职责**：HTTP 请求处理

| 文件 | 功能 |
|------|------|
| `handlers.go` | WHIP/WHEP/房间/录制/管理端点处理 |
| `middleware.go` | CORS、限流、Token/JWT 认证 |
| `routes.go` | URL 路由、参数提取、房间名校验 |

**认证优先级**：
```
1. 房间级 Token (ROOM_TOKENS)
    ↓ (未找到或失败)
2. 全局 Token (AUTH_TOKEN)
    ↓ (未找到或失败)
3. JWT (JWT_SECRET)
    ↓ (未找到或失败)
4. 允许访问 (未配置认证)
```

### internal/sfu

**职责**：WebRTC SFU 核心逻辑

```
┌─────────────────────────────────────────────────────┐
│                     Manager                          │
│  - 管理所有 Room 实例                                │
│  - 创建/删除 Room                                    │
│  - 统计房间数量                                      │
└──────────────────────┬──────────────────────────────┘
                       │ 1:N
                       ▼
┌─────────────────────────────────────────────────────┐
│                      Room                            │
│  - Publisher PeerConnection                          │
│  - Subscriber PeerConnections                        │
│  - TrackFeeds (TrackFanout map)                      │
└──────────────────────┬──────────────────────────────┘
                       │ 1:N
                       ▼
┌─────────────────────────────────────────────────────┐
│                  TrackFanout                         │
│  - Remote Track (from publisher)                     │
│  - Local Tracks (to subscribers)                     │
│  - readLoop: RTP distribution                        │
│  - Optional: Recorder (IVF/OGG writer)               │
└─────────────────────────────────────────────────────┘
```

**关键方法**：

| 方法 | 作用 |
|------|------|
| `Manager.Publish()` | 创建房间，建立发布者连接 |
| `Manager.Subscribe()` | 创建订阅者连接，绑定现有轨道 |
| `Room.attachTrackFeed()` | 新轨道分发到所有订阅者 |
| `TrackFanout.readLoop()` | RTP 包读取与分发循环 |

### internal/metrics

**职责**：Prometheus 指标暴露

| 指标 | 类型 | 说明 |
|------|------|------|
| `live_rooms` | Gauge | 活跃房间数 |
| `live_subscribers` | GaugeVec | 每房间订阅者数 |
| `rtp_bytes_total` | CounterVec | RTP 字节数 |
| `rtp_packets_total` | CounterVec | RTP 包数 |

### internal/uploader

**职责**：S3/MinIO 文件上传

```
上传流程:
1. 检查 Enabled() → client != nil
2. 打开本地文件
3. 构建对象键 (prefix + filename)
4. client.PutObject()
5. (可选) 删除本地文件
```

---

## 数据流

### 推流流程 (WHIP)

```
1. 发布者 → POST /api/whip/publish/{room}
   │
2. HTTPHandlers.ServeWHIPPublish()
   │  ├─ CORS 检查
   │  ├─ 限流检查
   │  └─ 认证检查
   │
3. Manager.Publish(roomName, sdpOffer)
   │  ├─ getOrCreateRoom()
   │  └─ Room.Publish(sdpOffer)
   │
4. Room.Publish()
   │  ├─ 创建 MediaEngine + Interceptors
   │  ├─ NewPeerConnection(ICEConfig)
   │  ├─ SetRemoteDescription(offer)
   │  ├─ CreateAnswer()
   │  ├─ SetLocalDescription(answer)
   │  └─ OnTrack: attachTrackFeed()
   │
5. 返回 SDP Answer
   │
6. TrackFanout.readLoop() 持续运行
   │  ├─ 从 Remote Track 读取 RTP
   │  ├─ 写入录制器（如启用）
   │  └─ 分发到所有 Local Tracks
```

### 播放流程 (WHEP)

```
1. 观看者 → POST /api/whep/play/{room}
   │
2. HTTPHandlers.ServeWHEPPlay()
   │  ├─ CORS/限流/认证检查
   │  └─ Manager.Subscribe()
   │
3. Manager.Subscribe(roomName, sdpOffer)
   │  └─ Room.Subscribe(sdpOffer)
   │
4. Room.Subscribe()
   │  ├─ 检查订阅者上限
   │  ├─ NewPeerConnection()
   │  ├─ 遍历现有 TrackFeeds
   │  │   └─ TrackFanout.attachToSubscriber()
   │  ├─ SetRemoteDescription/CreateAnswer
   │  └─ OnICEStateChange: removeSubscriber()
   │
5. 返回 SDP Answer
```

### 断开连接

```
ICE 状态变更 (Failed/Disconnected/Closed)
    │
    ▼
┌─────────────────────────────────────┐
│ 发布者断开                           │
├─────────────────────────────────────┤
│ 1. closePublisher()                  │
│ 2. 关闭所有 TrackFanout              │
│ 3. 上传录制文件                      │
│ 4. 清空订阅者列表                    │
│ 5. pruneIfEmpty()                    │
└─────────────────────────────────────┘

┌─────────────────────────────────────┐
│ 订阅者断开                           │
├─────────────────────────────────────┤
│ 1. removeSubscriber()                │
│ 2. 从 TrackFanouts 移除绑定          │
│ 3. 关闭 PeerConnection               │
│ 4. pruneIfEmpty()                    │
└─────────────────────────────────────┘
```

---

## 认证体系

### Token 认证

```
优先级 1: 房间 Token (ROOM_TOKENS)
┌─────────────────────────────────────┐
│ ROOM_TOKENS="room1:abc;room2:def"   │
│                                      │
│ 访问 room1 → 检查 token == "abc"    │
│ 访问 room2 → 检查 token == "def"    │
│ 访问 room3 → 回退到全局 Token       │
└─────────────────────────────────────┘

优先级 2: 全局 Token (AUTH_TOKEN)
┌─────────────────────────────────────┐
│ AUTH_TOKEN="secret123"              │
│                                      │
│ 所有房间使用相同 Token               │
└─────────────────────────────────────┘
```

### JWT 认证

```go
// JWT Claims 结构
type roomClaims struct {
    Room  string `json:"room,omitempty"`   // 限制房间
    Role  string `json:"role,omitempty"`   // "admin" 角色
    Admin any    `json:"admin,omitempty"`  // true/1 管理员
    jwt.RegisteredClaims
}

// 使用场景
1. 房间访问: claims.Room == room or claims.Room == ""
2. 管理接口: claims.Role == "admin" or claims.Admin == true
```

---

## 录制与上传

### 录制格式

| 编解码器 | 文件格式 | 写入器 |
|----------|----------|--------|
| Opus | .ogg | oggwriter (48kHz, 双声道) |
| VP8 | .ivf | ivfwriter |
| VP9 | .ivf | ivfwriter |

### 文件命名

```
{room}_{trackID}_{unixTimestamp}.{ext}

示例: demo_video0_1710123456.ivf
```

### 上传流程

```
Room.closePublisher()
    │
    ▼
TrackFanout.close() → 返回录制文件路径
    │
    ▼
uploader.Enabled()?
    │ Yes
    ▼
go uploadRecording(path)
    │
    ▼
Upload(ctx, path)
    ├─ PutObject(S3Bucket, objectKey, file)
    └─ (可选) os.Remove(localFile)
```

---

## 可观测性

### Prometheus 指标

```
# 活跃房间数
live_rooms

# 每房间订阅者数
live_subscribers{room="demo"}

# RTP 字节数（累计）
live_rtp_bytes_total{room="demo"}

# RTP 包数（累计）
live_rtp_packets_total{room="demo"}
```

### OpenTelemetry 追踪

```
环境变量:
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317
OTEL_EXPORTER_OTLP_PROTOCOL=grpc
OTEL_SERVICE_NAME=live-webrtc-go

追踪的 span:
- HTTP Handler: {method} {path}
```

### 健康检查

```
GET /healthz → "ok" (200 OK)
```

---

## 扩展点

### 1. 多实例部署

当前房间状态在内存中，多实例需要：
- 外部存储（Redis/数据库）存储房间映射
- 会话亲和（Sticky Session）
- 或客户端重定向

### 2. 媒体处理

可在 `TrackFanout.readLoop()` 前插入：
- 转码（FFmpeg 集成）
- 多码率
- 截图/水印

### 3. 认证扩展

在 `middleware.go` 中扩展：
- OAuth2 集成
- Webhook 回调验证
- IP 白名单

### 4. 存储扩展

实现 `rtpWriter` 接口：
```go
type rtpWriter interface {
    WriteRTP(*rtp.Packet) error
    Close() error
}
```

可支持：
- 实时转封装（MP4）
- 流式上传（不落地）
- CDN 推送

---

## 性能考量

### 内存使用

- 每个 TrackFanout: 约 1-2 MB（RTP 缓冲）
- 每个订阅者: 约 1500 bytes (MTU buffer)
- 录制缓冲: 取决于写入频率

### CPU 使用

- RTP 包处理: 主循环在 `readLoop()`
- 编解码协商: 仅连接建立时
- 指标更新: 每个 RTP 包

### 优化建议

1. 零拷贝 RTP 转发（需修改 TrackFanout）
2. 批量指标更新
3. 连接池（多房间场景）
4. SIMD 优化（大量订阅者）
