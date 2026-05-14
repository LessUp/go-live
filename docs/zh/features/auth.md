# 认证系统

多层认证系统，支持 Token 和 JWT。

## 认证流程

```mermaid
flowchart TB
    REQ[请求] --> HEADER{检查认证头}
    
    HEADER -->|Authorization: Bearer| TOKEN[提取 Token]
    HEADER -->|X-Auth-Token| TOKEN
    HEADER -->|无| NOAUTH[无认证]
    
    TOKEN --> ROOM{房间 Token?}
    ROOM -->|找到| ROOMCHECK[比较 ROOM_TOKENS]
    ROOM -->|未找到| GLOBAL
    
    ROOMCHECK -->|匹配| ALLOW[✓ 允许]
    ROOMCHECK -->|不匹配| GLOBAL
    
    GLOBAL[全局 Token 检查] --> GLOBALFOUND{AUTH_TOKEN?}
    GLOBALFOUND -->|是| GLOBALCHECK[比较 Token]
    GLOBALFOUND -->|否| JWT
    
    GLOBALCHECK -->|匹配| ALLOW
    GLOBALCHECK -->|不匹配| JWT
    
    JWT[JWT 检查] --> JWTFOUND{JWT_SECRET?}
    JWTFOUND -->|是| JWTCHECK[验证并解析 JWT]
    JWTFOUND -->|否| NOAUTH
    
    JWTCHECK -->|无效| DENY[✗ 401 未授权]
    JWTCHECK -->|有效| JWTAUTH{检查声明}
    
    JWTAUTH -->|房间匹配| ALLOW
    JWTAUTH -->|房间不匹配| DENY
    
    NOAUTH -->|已配置认证| DENY
    NOAUTH -->|未配置认证| ALLOW
```

## Token 认证

### 全局 Token

所有房间使用同一 Token：

```bash
AUTH_TOKEN=my-secret-token
```

```http
Authorization: Bearer my-secret-token
# 或
X-Auth-Token: my-secret-token
```

### 房间级 Token

不同房间使用不同 Token：

```bash
ROOM_TOKENS=room1:token1;room2:token2;room3:token3
```

| 房间 | Token |
|------|-------|
| `room1` | `token1` |
| `room2` | `token2` |
| `room3` | `token3` |
| `room4` | 回退到 `AUTH_TOKEN` |

## JWT 认证

JWT 通过声明提供更多灵活性。

### JWT 结构

```json
{
  "sub": "user123",
  "room": "demo",
  "role": "admin",
  "admin": true,
  "exp": 1710123456,
  "iat": 1710120000
}
```

### JWT 声明

| 声明 | 类型 | 说明 |
|------|------|------|
| `sub` | string | 主题（用户标识符） |
| `room` | string | 限制到特定房间（可选） |
| `role` | string | 角色（`admin` 用于管理员访问） |
| `admin` | boolean | 管理员标志（role 的替代） |
| `exp` | number | 过期时间戳 |
| `iat` | number | 签发时间戳 |

### JWT 配置

```bash
JWT_SECRET=your-signing-secret
JWT_AUDIENCE=your-app-name  # 可选：要求特定 audience
```

## 管理员认证

管理员端点需要 `ADMIN_TOKEN`：

```bash
ADMIN_TOKEN=admin-secret-token
```

```http
POST /api/admin/rooms/{room}/close
Authorization: Bearer admin-secret-token
```

## 认证优先级

```mermaid
flowchart LR
    A[1. 房间 Token<br/>ROOM_TOKENS] --> B[2. 全局 Token<br/>AUTH_TOKEN]
    B --> C[3. JWT<br/>JWT_SECRET]
    C --> D[4. 无认证]
```

## 安全最佳实践

### Token 安全

- 使用强随机 Token（32+ 字符）
- 定期轮换 Token
- 不同环境使用不同 Token
- 永远不要将 Token 提交到版本控制

### JWT 安全

- 使用强签名密钥（256+ 位）
- 设置适当的过期时间
- 如使用 `JWT_AUDIENCE`，验证 audience 声明
- 生产环境使用 HTTPS

### Token 生成示例

```bash
# 生成随机 Token
openssl rand -hex 32

# 生成 JWT 密钥
openssl rand -base64 32
```

## 认证头格式

支持两种格式：

```http
# Bearer token（推荐）
Authorization: Bearer <token>

# 自定义头
X-Auth-Token: <token>
```

## 错误响应

| 状态码 | 信息 | 原因 |
|--------|------|------|
| 401 | `unauthorized` | 无效/缺失 Token |
| 403 | `forbidden` | JWT 房间不匹配 |
