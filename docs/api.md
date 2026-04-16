---
layout: default
title: API 参考
---

# API 参考

本文档详细介绍 live-webrtc-go 的所有 HTTP API 端点。

## 目录

- [认证方式](#认证方式)
- [流媒体接口](#流媒体接口)
- [状态查询接口](#状态查询接口)
- [管理接口](#管理接口)
- [健康与监控](#健康与监控)
- [错误响应](#错误响应)

## 认证方式

系统支持三种认证方式，按优先级依次尝试：

### 1. Bearer Token

```http
Authorization: Bearer <token>
```

### 2. X-Auth-Token Header

```http
X-Auth-Token: <token>
```

### 3. URL Query Parameter

```http
GET /api/rooms?token=<token>
```

### 认证优先级

1. 房间级 Token (`ROOM_TOKENS`)
2. 全局 Token (`AUTH_TOKEN`)
3. JWT (`JWT_SECRET`)
4. 无认证（未配置认证时）

## 流媒体接口

### WHIP 推流

发布媒体流到指定房间。

```http
POST /api/whip/publish/{room}
Content-Type: application/sdp
Authorization: Bearer <token>
```

**参数**

| 参数 | 位置 | 类型 | 必需 | 说明 |
|------|------|------|------|------|
| `room` | path | string | 是 | 房间名，匹配 `^[A-Za-z0-9_-]{1,64}$` |

**请求体**

SDP Offer (text/plain)

**响应**

| 状态码 | 说明 |
|--------|------|
| 200 | 成功，返回 SDP Answer |
| 400 | 无效的 SDP 或房间名 |
| 401 | 认证失败 |
| 409 | 房间已有发布者 |
| 429 | 请求频率超限 |

**示例**

```bash
curl -X POST "http://localhost:8080/api/whip/publish/demo" \
  -H "Content-Type: application/sdp" \
  -H "Authorization: Bearer mytoken" \
  --data-binary @offer.sdp
```

---

### WHEP 播放

从指定房间订阅媒体流。

```http
POST /api/whep/play/{room}
Content-Type: application/sdp
Authorization: Bearer <token>
```

**参数**

| 参数 | 位置 | 类型 | 必需 | 说明 |
|------|------|------|------|------|
| `room` | path | string | 是 | 房间名 |

**请求体**

SDP Offer (text/plain)

**响应**

| 状态码 | 说明 |
|--------|------|
| 200 | 成功，返回 SDP Answer |
| 400 | 无效的 SDP 或房间名 |
| 401 | 认证失败 |
| 404 | 房间不存在 |
| 429 | 请求频率超限 |

**示例**

```bash
curl -X POST "http://localhost:8080/api/whep/play/demo" \
  -H "Content-Type: application/sdp" \
  -H "Authorization: Bearer mytoken" \
  --data-binary @offer.sdp
```

---

## 状态查询接口

### 获取前端配置

返回前端应用所需的运行时配置。

```http
GET /api/bootstrap
```

**响应**

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

**字段说明**

| 字段 | 类型 | 说明 |
|------|------|------|
| `authEnabled` | boolean | 是否启用认证 |
| `recordEnabled` | boolean | 是否启用录制 |
| `iceServers` | array | ICE 服务器配置 |
| `features` | object | 功能开关 |

---

### 获取房间列表

返回所有活跃房间及其状态。

```http
GET /api/rooms
```

**响应**

```json
[
  {
    "name": "demo",
    "hasPublisher": true,
    "subscriberCount": 5
  },
  {
    "name": "test",
    "hasPublisher": false,
    "subscriberCount": 0
  }
]
```

**字段说明**

| 字段 | 类型 | 说明 |
|------|------|------|
| `name` | string | 房间名 |
| `hasPublisher` | boolean | 是否有发布者 |
| `subscriberCount` | number | 订阅者数量 |

---

### 获取录制列表

返回所有录制文件的元数据。

```http
GET /api/records
```

**响应**

```json
[
  {
    "name": "demo_video0_1710123456.ivf",
    "size": 1048576,
    "modTime": "2024-03-10T12:34:56Z"
  }
]
```

**字段说明**

| 字段 | 类型 | 说明 |
|------|------|------|
| `name` | string | 文件名 |
| `size` | number | 文件大小（字节） |
| `modTime` | string | 修改时间（ISO 8601） |

---

## 管理接口

### 关闭房间

强制关闭指定房间，断开所有连接。

```http
POST /api/admin/rooms/{room}/close
Authorization: Bearer <admin-token>
```

**参数**

| 参数 | 位置 | 类型 | 必需 | 说明 |
|------|------|------|------|------|
| `room` | path | string | 是 | 要关闭的房间名 |

**响应**

| 状态码 | 说明 |
|--------|------|
| 200 | 成功关闭 |
| 401 | 认证失败（需要 Admin Token） |
| 404 | 房间不存在 |

**示例**

```bash
curl -X POST \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  http://localhost:8080/api/admin/rooms/demo/close
```

---

## 健康与监控

### 健康检查

```http
GET /healthz
```

**响应**

```
ok
```

状态码：200

---

### Prometheus 指标

```http
GET /metrics
```

**可用指标**

| 指标名 | 类型 | 标签 | 说明 |
|--------|------|------|------|
| `live_rooms` | Gauge | - | 活跃房间数 |
| `live_subscribers` | GaugeVec | `room` | 每房间订阅者数 |
| `live_rtp_bytes_total` | CounterVec | `room` | RTP 字节总数 |
| `live_rtp_packets_total` | CounterVec | `room` | RTP 包总数 |

**示例**

```bash
curl http://localhost:8080/metrics
```

---

## 错误响应

所有错误响应格式：

```json
{
  "error": "错误描述"
}
```

### 常见错误码

| 状态码 | 错误信息 | 原因 |
|--------|----------|------|
| 400 | `invalid room name` | 房间名格式错误 |
| 400 | `invalid SDP` | SDP 格式错误 |
| 401 | `unauthorized` | 认证失败或未提供 |
| 404 | `room not found` | 房间不存在（WHEP 播放时） |
| 409 | `publisher already exists` | 房间已有发布者 |
| 429 | `too many requests` | 触发限流 |
| 500 | `internal server error` | 服务器内部错误 |

---

## CORS 配置

所有 API 响应包含 CORS 头：

```http
Access-Control-Allow-Origin: <ALLOWED_ORIGIN>
Access-Control-Allow-Methods: GET, POST, OPTIONS
Access-Control-Allow-Headers: Content-Type, Authorization, X-Auth-Token
```

预检请求 (OPTIONS) 自动返回 204。

---

## 请求限制

| 限制 | 值 | 配置 |
|------|-----|------|
| SDP 请求体大小 | 1 MB | 硬编码 |
| 房间名长度 | 1-64 字符 | 正则 `^[A-Za-z0-9_-]{1,64}$` |
| 请求频率 | 可配置 | `RATE_LIMIT_RPS`, `RATE_LIMIT_BURST` |
| 每房间订阅者 | 可配置 | `MAX_SUBS_PER_ROOM` |
