# 快速开始

本指南介绍本地开发、容器部署和基本使用方法。

## 环境要求

| 依赖项 | 版本 | 说明 |
|--------|------|------|
| Go | 1.22+ | 编译运行必需 |
| Docker | 20.10+ | 可选，容器化部署 |
| 浏览器 | Chrome 90+/Firefox 88+ | 需要 WebRTC 支持 |

## 快速启动

### 方式一：直接运行

```bash
# 克隆项目
git clone https://github.com/LessUp/go-live.git
cd go-live

# 下载依赖
go mod tidy

# 运行服务
go run ./cmd/server
```

### 方式二：使用开发脚本（推荐）

```bash
# 基础启动（加载 .env.local 如存在）
./scripts/start.sh

# 启动前执行 go mod tidy
RUN_TIDY=1 ./scripts/start.sh
```

脚本功能：
- 创建 `records`、`.gocache`、`.gomodcache` 目录
- 设置 `GOCACHE` 和 `GOMODCACHE` 环境变量
- 加载 `.env.local` 文件（如存在）
- 启动服务

### 启动后验证

```bash
# 健康检查
curl http://localhost:8080/healthz
# 输出: ok

# 查看房间列表
curl http://localhost:8080/api/rooms
# 输出: []

# 获取前端配置
curl http://localhost:8080/api/bootstrap | jq
```

## Docker 部署

### 构建镜像

```bash
docker build -t live-webrtc-go:latest .
```

### 基础运行

```bash
docker run --rm -p 8080:8080 live-webrtc-go:latest
```

### 启用录制

```bash
docker run --rm -p 8080:8080 \
  -e RECORD_ENABLED=1 \
  -e RECORD_DIR=/records \
  -v $(pwd)/records:/records \
  live-webrtc-go:latest
```

### Docker Compose

```yaml
version: '3.8'

services:
  live-webrtc:
    build: .
    ports:
      - "8080:8080"
    environment:
      - HTTP_ADDR=:8080
      - AUTH_TOKEN=${AUTH_TOKEN}
      - RECORD_ENABLED=1
      - RECORD_DIR=/records
    volumes:
      - ./records:/records
    restart: unless-stopped
```

## 使用 OBS 推流

1. 打开 OBS → 设置 → 流
2. 服务选择 "WHIP"
3. 服务器输入: `http://localhost:8080/api/whip/publish/myroom`
4. 令牌输入: `your-token`（如配置了认证）
5. 开始推流

## 使用 JavaScript 播放

```javascript
// 获取配置
const config = await fetch('/api/bootstrap').then(r => r.json());

// 创建 PeerConnection
const pc = new RTCPeerConnection({
  iceServers: config.iceServers
});

// 创建收发器（只接收）
pc.addTransceiver('video', { direction: 'recvonly' });
pc.addTransceiver('audio', { direction: 'recvonly' });

// 处理远程轨道
pc.ontrack = (event) => {
  const video = document.getElementById('video');
  video.srcObject = event.streams[0];
};

// 创建 Offer
const offer = await pc.createOffer();
await pc.setLocalDescription(offer);

// 发送 WHEP 请求
const response = await fetch('/api/whep/play/myroom', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/sdp',
    'Authorization': 'Bearer your-token'
  },
  body: offer.sdp
});

const answer = await response.text();
await pc.setRemoteDescription({ type: 'answer', sdp: answer });
```

## 下一步

- [系统架构](/zh/architecture/overview) - 了解系统设计
- [API 端点](/zh/api/endpoints) - 完整 API 参考
- [配置项](/zh/api/configuration) - 所有配置选项
