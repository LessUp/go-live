---
layout: default
title: 使用指南
nav_order: 2
lang: zh
---

# 使用指南

本文档详细介绍 live-webrtc-go 的本地开发、容器部署、API 使用和故障排除。

{: .no_toc }

## 目录

{: .no_toc .text-delta }

1. TOC
{:toc}

---

## 环境要求

| 依赖项 | 版本 | 说明 |
|--------|------|------|
| Go | 1.22+ | 编译运行必需 |
| Docker | 20.10+ | 可选，容器化部署 |
| Docker Compose | 2.0+ | 可选，多容器编排 |
| 浏览器 | Chrome 90+/Firefox 88+ | 需要 WebRTC 支持 |

### WebRTC 端口要求

如果在 NAT 环境下使用，需要：
- 配置 STUN/TURN 服务器
- 确保 UDP 端口未被防火墙阻止

---

## 本地开发

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

### 环境变量文件

创建 `.env.local` 文件：

```bash
# 服务配置
HTTP_ADDR=:8080
ALLOWED_ORIGIN=*

# 认证配置
AUTH_TOKEN=your-secret-token
# ROOM_TOKENS=room1:token1;room2:token2
# JWT_SECRET=jwt-signing-secret

# WebRTC 配置
# STUN_URLS=stun:stun.l.google.com:19302
# TURN_URLS=turn:turn.example.com:3478
# TURN_USERNAME=username
# TURN_PASSWORD=password

# 录制配置
RECORD_ENABLED=1
RECORD_DIR=records

# 限流配置
RATE_LIMIT_RPS=10
RATE_LIMIT_BURST=20
```

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

---

## 配置说明

### 核心配置

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `HTTP_ADDR` | `:8080` | HTTP 监听地址，格式 `host:port` |
| `ALLOWED_ORIGIN` | `*` | CORS 允许来源，`*` 表示任意 |

### 认证配置

| 变量 | 格式 | 说明 |
|------|------|------|
| `AUTH_TOKEN` | 字符串 | 全局访问令牌 |
| `ROOM_TOKENS` | `room1:tok1;room2:tok2` | 房间级令牌映射 |
| `JWT_SECRET` | 字符串 | JWT HMAC 签名密钥 |
| `JWT_AUDIENCE` | 字符串 | 要求的 JWT audience 声明 |
| `ADMIN_TOKEN` | 字符串 | 管理接口访问令牌 |

**认证优先级**：
1. 房间级 Token → 2. 全局 Token → 3. JWT → 4. 无认证

### WebRTC/ICE 配置

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `STUN_URLS` | `stun:stun.l.google.com:19302` | STUN 服务器列表 |
| `TURN_URLS` | - | TURN 服务器列表 |
| `TURN_USERNAME` | - | TURN 用户名 |
| `TURN_PASSWORD` | - | TURN 密码 |

**示例**：
```bash
STUN_URLS=stun:stun1.l.google.com:19302,stun:stun2.l.google.com:19302
TURN_URLS=turn:turn.example.com:3478,turns:turn.example.com:5349
TURN_USERNAME=myuser
TURN_PASSWORD=mypassword
```

### 录制配置

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `RECORD_ENABLED` | `0` | 启用录制，设为 `1` 启用 |
| `RECORD_DIR` | `records` | 录制文件存储目录 |

### S3 上传配置

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `UPLOAD_RECORDINGS` | `0` | 启用上传，设为 `1` 启用 |
| `DELETE_RECORDING_AFTER_UPLOAD` | `0` | 上传后删除本地文件 |
| `S3_ENDPOINT` | - | S3/MinIO 端点地址 |
| `S3_REGION` | - | S3 区域 |
| `S3_BUCKET` | - | 目标存储桶名称 |
| `S3_ACCESS_KEY` | - | Access Key ID |
| `S3_SECRET_KEY` | - | Secret Access Key |
| `S3_USE_SSL` | `1` | 使用 HTTPS 连接 |
| `S3_PATH_STYLE` | `0` | 使用 path-style 寻址 |
| `S3_PREFIX` | - | 对象键前缀 |

**MinIO 示例**：
```bash
S3_ENDPOINT=minio.example.com:9000
S3_ACCESS_KEY=minioadmin
S3_SECRET_KEY=minioadmin
S3_BUCKET=recordings
S3_USE_SSL=0
S3_PATH_STYLE=1
```

### 限流配置

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `RATE_LIMIT_RPS` | `0` | 每 IP 每秒请求数，`0` 禁用 |
| `RATE_LIMIT_BURST` | `0` | 突发容量 |

### TLS 配置

| 变量 | 说明 |
|------|------|
| `TLS_CERT_FILE` | TLS 证书文件路径 |
| `TLS_KEY_FILE` | TLS 私钥文件路径 |

```bash
TLS_CERT_FILE=/etc/ssl/cert.pem
TLS_KEY_FILE=/etc/ssl/key.pem
```

### 调试配置

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `PPROF` | `0` | 启用 pprof 端点 |
| `OTEL_SERVICE_NAME` | `live-webrtc-go` | OpenTelemetry 服务名称 |

---

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

### 完整配置示例

```bash
docker run --rm -p 8080:8080 \
  -e HTTP_ADDR=:8080 \
  -e AUTH_TOKEN=mytoken \
  -e RECORD_ENABLED=1 \
  -e RECORD_DIR=/records \
  -e UPLOAD_RECORDINGS=1 \
  -e S3_ENDPOINT=s3.amazonaws.com \
  -e S3_BUCKET=my-bucket \
  -e S3_ACCESS_KEY=$AWS_ACCESS_KEY_ID \
  -e S3_SECRET_KEY=$AWS_SECRET_ACCESS_KEY \
  -v $(pwd)/records:/records \
  live-webrtc-go:latest
```

### Docker Compose

`docker-compose.yml`:

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
      - RATE_LIMIT_RPS=10
      - RATE_LIMIT_BURST=20
    volumes:
      - ./records:/records
    restart: unless-stopped
```

启动：
```bash
docker compose up -d
```

---

## Kubernetes 部署

### Deployment 示例

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: live-webrtc
spec:
  replicas: 3
  selector:
    matchLabels:
      app: live-webrtc
  template:
    metadata:
      labels:
        app: live-webrtc
    spec:
      containers:
      - name: live-webrtc
        image: live-webrtc:latest
        ports:
        - containerPort: 8080
        env:
        - name: HTTP_ADDR
          value: ":8080"
        - name: AUTH_TOKEN
          valueFrom:
            secretKeyRef:
              name: live-webrtc-secret
              key: auth-token
        volumeMounts:
        - name: records
          mountPath: /records
      volumes:
      - name: records
        persistentVolumeClaim:
          claimName: records-pvc
```

### Service 示例

```yaml
apiVersion: v1
kind: Service
metadata:
  name: live-webrtc
spec:
  selector:
    app: live-webrtc
  ports:
  - port: 80
    targetPort: 8080
  type: LoadBalancer
```

---

## API 使用示例

### 使用 curl

```bash
# 设置变量
HOST=http://localhost:8080
TOKEN=your-token

# 健康检查
curl $HOST/healthz

# 获取前端配置
curl $HOST/api/bootstrap | jq

# 列出房间
curl $HOST/api/rooms | jq

# 列出录制文件
curl $HOST/api/records | jq

# 关闭房间（管理接口）
curl -X POST \
  -H "Authorization: Bearer $TOKEN" \
  $HOST/api/admin/rooms/myroom/close
```

### 使用 OBS 推流

1. 打开 OBS → 设置 → 流
2. 服务选择 "WHIP"
3. 服务器输入: `http://localhost:8080/api/whip/publish/myroom`
4. 令牌输入: `your-token`（如配置了认证）
5. 开始推流

### 使用 JavaScript 播放

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

// 发送 WHIP 请求
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

---

## 前端集成

### Bootstrap API 响应

```json
{
  "authEnabled": true,
  "recordEnabled": true,
  "iceServers": [
    {
      "urls": ["stun:stun.l.google.com:19302"]
    }
  ],
  "features": {
    "rooms": true,
    "records": true
  }
}
```

### 认证头格式

支持两种格式：

```
Authorization: Bearer <token>
X-Auth-Token: <token>
```

---

## 故障排除

### 常见问题

| 问题 | 可能原因 | 解决方案 |
|------|----------|----------|
| `publisher already exists` | 房间已有发布者 | 使用不同房间名或等待发布者断开 |
| `unauthorized` | 认证失败 | 检查 Token 或 JWT 配置 |
| `too many requests` | 触发限流 | 增加 `RATE_LIMIT_BURST` 或等待 |
| `room not found` | 房间不存在 | 确保发布者已连接 |
| `ICE connection failed` | NAT 穿透问题 | 配置 TURN 服务器 |
| 无视频/音频 | 编解码器不匹配 | 检查浏览器支持的编解码器 |

### 调试步骤

1. **检查服务状态**
   ```bash
   curl http://localhost:8080/healthz
   curl http://localhost:8080/api/rooms
   ```

2. **查看指标**
   ```bash
   curl http://localhost:8080/metrics
   ```

3. **启用 pprof**
   ```bash
   PPROF=1 go run ./cmd/server
   # 访问 http://localhost:8080/debug/pprof/
   ```

4. **检查日志**
   - 查看控制台输出的 `slog` 日志
   - 关注 `ERROR` 级别日志

### WebRTC 连接问题

**症状**：浏览器显示 "ICE connection failed"

**排查**：
1. 检查 STUN 服务器可达性
2. 配置 TURN 服务器（NAT 环境）
3. 确认 UDP 端口未被阻止

**TURN 配置示例**：
```bash
TURN_URLS=turn:turn.example.com:3478
TURN_USERNAME=username
TURN_PASSWORD=password
```

### 认证问题

**症状**：API 返回 `401 Unauthorized`

**排查**：
1. 确认发送了正确的 Authorization 头
2. 检查 Token 是否正确
3. 如使用 JWT，验证签名和有效期

### 录制问题

**症状**：录制目录为空

**排查**：
1. 确认 `RECORD_ENABLED=1`
2. 检查目录权限：`ls -la records/`
3. 查看日志是否有写入错误
4. 确认发布者推送的是 VP8/VP9/Opus 编码

### 上传问题

**症状**：文件未上传到 S3

**排查**：
1. 确认 `UPLOAD_RECORDINGS=1`
2. 检查所有 S3_* 变量已配置
3. 验证 S3 凭证有效性
4. 检查网络连通性

```bash
# 测试 S3 连接
aws s3 ls --endpoint-url http://minio:9000
```

### 性能问题

**症状**：延迟高、卡顿

**排查**：
1. 检查服务器 CPU/内存使用
2. 查看订阅者数量是否过多
3. 考虑设置 `MAX_SUBS_PER_ROOM` 限制
4. 检查网络带宽

```bash
# 查看指标
curl http://localhost:8080/metrics | grep live_
```

---

## 开发命令参考

```bash
# 构建
make build

# 格式化
make fmt

# 代码检查
make lint

# 安全扫描
make security

# 运行测试
make test          # 单元 + 集成 + 安全
make test-all      # 包含 e2e + 性能
make test-unit     # 仅单元测试

# 覆盖率
make coverage
open coverage.html
```

---

## 相关文档

- [设计说明](design.md) - 架构和模块详解
- [API 参考](api.md) - 完整 API 文档
- [GitHub 仓库](https://github.com/LessUp/go-live) - 源代码
