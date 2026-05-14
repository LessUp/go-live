# WHEP 协议

WHEP (WebRTC-HTTP Egress Protocol) 用于订阅媒体流。

## 概述

```mermaid
sequenceDiagram
    participant V as 观众
    participant S as 服务器

    V->>S: POST /api/whep/play/{room}<br/>Content-Type: application/sdp<br/>SDP Offer
    S->>S: 验证房间名
    S->>S: 检查认证
    S->>S: 检查发布者存在
    S->>S: 创建 PeerConnection
    S->>S: 绑定现有轨道
    S-->>V: 201 Created<br/>Content-Type: application/sdp<br/>SDP Answer
    
    Note over V,S: WebRTC 连接建立
    
    S->>V: RTP 媒体包
```

## 端点

```http
POST /api/whep/play/{room}
Content-Type: application/sdp
Authorization: Bearer <token>
```

### 参数

| 参数 | 位置 | 类型 | 说明 |
|------|------|------|------|
| `room` | path | string | 房间名（1-64 字符，`A-Za-z0-9_-`） |

### 响应码

| 状态码 | 说明 |
|--------|------|
| 201 | 成功 - 返回 SDP Answer |
| 400 | 无效的房间名或 SDP |
| 401 | 认证失败 |
| 403 | 订阅者数量已达上限 |
| 404 | 房间无活跃发布者 |
| 429 | 请求频率超限 |

## 连接流程

```mermaid
flowchart TB
    subgraph Viewer
        A[创建 PeerConnection] --> B[添加 recvonly 收发器]
        B --> C[创建 Offer]
        C --> D[设置本地描述]
        D --> E[HTTP POST 到服务器]
    end

    subgraph Server
        E --> F[验证请求]
        F --> G{发布者存在?}
        G -->|否| ERR1[404 Not Found]
        G -->|是| H{订阅者限制?}
        H -->|已达上限| ERR2[403 Forbidden]
        H -->|OK| I[创建 PeerConnection]
        I --> J[绑定现有轨道]
        J --> K[返回 Answer]
    end

    subgraph Established
        K --> L[ICE 协商]
        L --> M[接收媒体]
    end
```

## 浏览器示例

```javascript
// 获取 ICE 配置
const config = await fetch('/api/bootstrap').then(r => r.json());

// 创建 PeerConnection
const pc = new RTCPeerConnection({
  iceServers: config.iceServers
});

// 创建 recvonly 收发器
pc.addTransceiver('video', { direction: 'recvonly' });
pc.addTransceiver('audio', { direction: 'recvonly' });

// 处理传入轨道
const video = document.getElementById('video');
pc.ontrack = (event) => {
  video.srcObject = event.streams[0];
};

// 创建 offer
const offer = await pc.createOffer();
await pc.setLocalDescription(offer);

// 等待 ICE 收集
await new Promise(resolve => {
  if (pc.iceGatheringState === 'complete') {
    resolve();
  } else {
    pc.onicegatheringstatechange = () => {
      if (pc.iceGatheringState === 'complete') resolve();
    };
  }
});

// 发送 WHEP 请求
const response = await fetch('/api/whep/play/myroom', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/sdp',
    'Authorization': 'Bearer mytoken'
  },
  body: pc.localDescription.sdp
});

if (response.ok) {
  const answer = await response.text();
  await pc.setRemoteDescription({ type: 'answer', sdp: answer });
}
```

## 错误处理

| 错误 | 原因 | 解决方案 |
|------|------|----------|
| `404 Not Found` | 房间无发布者 | 等待发布者 |
| `403 Forbidden` | 订阅者已达上限 | 增加 `MAX_SUBS_PER_ROOM` |
| `401 Unauthorized` | 无效/缺失 Token | 检查认证 |
