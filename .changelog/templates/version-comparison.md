# 版本对比: v{OLD_VERSION} vs v{NEW_VERSION}

## 📊 版本概览

| 项目 | v{OLD_VERSION} | v{NEW_VERSION} |
|------|----------------|----------------|
| 发布日期 | {OLD_DATE} | {NEW_DATE} |
| 版本类型 | {OLD_TYPE} | {NEW_TYPE} |

---

## 🔍 详细对比

### 新增功能

| 功能 | v{OLD_VERSION} | v{NEW_VERSION} |
|------|----------------|----------------|
| Feature 1 | ❌ | ✅ |
| Feature 2 | ❌ | ✅ |

### API 变更

| 端点 | v{OLD_VERSION} | v{NEW_VERSION} | 说明 |
|------|----------------|----------------|------|
| `/api/xxx` | ✅ | ✅ | 无变更 |
| `/api/yyy` | ❌ | ✅ | 新增 |

### 配置变更

| 配置项 | v{OLD_VERSION} | v{NEW_VERSION} | 说明 |
|--------|----------------|----------------|------|
| `VAR_1` | 默认值 A | 默认值 B | 变更 |
| `VAR_2` | - | 新配置 | 新增 |

### 依赖变更

| 依赖 | v{OLD_VERSION} | v{NEW_VERSION} |
|------|----------------|----------------|
| package-a | v1.0.0 | v1.1.0 |
| package-b | - | v2.0.0 |

---

## 🚀 升级建议

### 自动升级

```bash
# Docker
docker pull ghcr.io/lessup/go-live:v{NEW_VERSION}
```

### 手动升级

1. 检查破坏性变更
2. 更新配置文件
3. 执行数据库迁移（如需要）
4. 部署新版本

---

## ⚠️ 注意事项

- [ ] 检查破坏性变更
- [ ] 测试关键功能
- [ ] 备份数据

---

<!--
模板变量:
- {OLD_VERSION} - 旧版本号
- {NEW_VERSION} - 新版本号
- {OLD_DATE} - 旧版本发布日期
- {NEW_DATE} - 新版本发布日期
- {OLD_TYPE} - 旧版本类型 (Major/Minor/Patch)
- {NEW_TYPE} - 新版本类型 (Major/Minor/Patch)
-->
