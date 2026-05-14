# 配置项

所有环境变量和配置选项的完整参考。

## 核心配置

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `HTTP_ADDR` | `:8080` | HTTP 监听地址，格式 `host:port` |
| `ALLOWED_ORIGIN` | `*` | CORS 允许来源，`*` 表示任意 |

## 认证配置

| 变量 | 格式 | 说明 |
|------|------|------|
| `AUTH_TOKEN` | 字符串 | 全局访问令牌 |
| `ROOM_TOKENS` | `room1:tok1;room2:tok2` | 房间级令牌映射 |
| `JWT_SECRET` | 字符串 | JWT HMAC 签名密钥 |
| `JWT_AUDIENCE` | 字符串 | 要求的 JWT audience 声明 |
| `ADMIN_TOKEN` | 字符串 | 管理接口访问令牌 |

**认证优先级**：房间级 Token → 全局 Token → JWT → 无认证

## WebRTC / ICE 配置

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

## 录制配置

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `RECORD_ENABLED` | `0` | 启用录制，设为 `1` 启用 |
| `RECORD_DIR` | `records` | 录制文件存储目录 |

## S3 上传配置

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

## 限流配置

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `RATE_LIMIT_RPS` | `0` | 每 IP 每秒请求数，`0` 禁用 |
| `RATE_LIMIT_BURST` | `0` | 突发容量 |

## 房间限制

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `MAX_SUBS_PER_ROOM` | `0` | 每房间最大订阅者数，`0` = 无限制 |

## TLS 配置

| 变量 | 说明 |
|------|------|
| `TLS_CERT_FILE` | TLS 证书文件路径 |
| `TLS_KEY_FILE` | TLS 私钥文件路径 |

```bash
TLS_CERT_FILE=/etc/ssl/cert.pem
TLS_KEY_FILE=/etc/ssl/key.pem
```

## 调试配置

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `PPROF` | `0` | 启用 pprof 端点 |
| `OTEL_SERVICE_NAME` | `live-webrtc-go` | OpenTelemetry 服务名称 |

## 环境变量文件

创建 `.env.local` 文件用于开发：

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

# 录制配置
RECORD_ENABLED=1
RECORD_DIR=records

# 限流配置
RATE_LIMIT_RPS=10
RATE_LIMIT_BURST=20
```

## 故障排除

| 问题 | 可能原因 | 解决方案 |
|------|----------|----------|
| `publisher already exists in this room` | 房间已有发布者 | 使用不同房间名或等待发布者断开 |
| `unauthorized` | 认证失败 | 检查 Token 或 JWT 配置 |
| `too many requests` | 触发限流 | 增加 `RATE_LIMIT_BURST` 或等待 |
| `no active publisher in room` | 房间无发布者 | 确保发布者已连接 |
| `subscriber limit reached` | 已达订阅者上限 | 增大限制或等待订阅者断开 |
| `ICE connection failed` | NAT 穿透问题 | 配置 TURN 服务器 |

### 调试步骤

1. **检查服务状态**: `curl http://localhost:8080/healthz`
2. **查看指标**: `curl http://localhost:8080/metrics`
3. **启用 pprof**: `PPROF=1 go run ./cmd/server`
4. **检查日志**: 关注 `ERROR` 级别日志
