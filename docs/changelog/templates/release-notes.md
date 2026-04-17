# Release v{VERSION} - {DATE}

## 🎉 Release Highlights

### English

Brief description of this release (1-2 sentences).

### 中文

本次发布的简要描述（1-2 句话）。

---

## ✨ What's New / 新特性

### English

#### Features
- ✨ Feature 1 - brief description
- ✨ Feature 2 - brief description

#### Improvements
- 🚀 Improvement 1 - brief description
- 🚀 Improvement 2 - brief description

### 中文

#### 新特性
- ✨ 特性 1 - 简要说明
- ✨ 特性 2 - 简要说明

#### 改进
- 🚀 改进 1 - 简要说明
- 🚀 改进 2 - 简要说明

---

## 🐛 Bug Fixes / 问题修复

### English
- 🐛 Fix 1 - description (#issue)
- 🐛 Fix 2 - description (#issue)

### 中文
- 🐛 修复 1 - 描述 (#issue)
- 🐛 修复 2 - 描述 (#issue)

---

## 🔒 Security / 安全

<!-- Remove this section if none -->

### English
- 🔒 Security fix 1 - description (CVE-XXXX-XXXXX if applicable)

### 中文
- 🔒 安全修复 1 - 描述（如有 CVE 编号请标注）

---

## ⚠️ Breaking Changes / 破坏性变更

<!-- Remove this section if none -->

### English

⚠️ **Important**: This release contains breaking changes.

- **Change 1**: Description and migration guide
- **Change 2**: Description and migration guide

### 中文

⚠️ **重要提示**：本次发布包含破坏性变更。

- **变更 1**：描述和迁移指南
- **变更 2**：描述和迁移指南

### Migration Guide / 迁移指南

```bash
# Steps to migrate from previous version
# 从前版本迁移的步骤
```

---

## 📦 Dependencies / 依赖更新

| Package | Old Version | New Version | Reason |
|---------|-------------|-------------|--------|
| package | v1.0.0 | v1.1.0 | Bug fix |

---

## 🙏 Contributors / 贡献者

Thanks to all contributors who made this release possible!

感谢所有为本次发布做出贡献的人！

- @username - contribution description / 贡献描述

---

## 🚀 Upgrade Guide / 升级指南

### Docker / Docker 部署

```bash
# Pull the latest image
docker pull ghcr.io/lessup/go-live:v{VERSION}

# 拉取最新镜像
docker pull ghcr.io/lessup/go-live:v{VERSION}
```

### Binary / 二进制

```bash
# Download from GitHub Releases
# 从 GitHub Releases 下载
curl -LO https://github.com/LessUp/go-live/releases/download/v{VERSION}/live-webrtc-go-linux-amd64

# Or build from source
# 或从源码构建
git clone https://github.com/LessUp/go-live.git
cd go-live
git checkout v{VERSION}
go build -o live-webrtc-go ./cmd/server
```

### Kubernetes / Kubernetes

```yaml
# Update image tag
# 更新镜像标签
image: ghcr.io/lessup/go-live:v{VERSION}
```

---

## 📝 Configuration Changes / 配置变更

<!-- Document any new or changed configuration options -->

### New Environment Variables / 新增环境变量

| Variable | Default | Description (EN) | Description (ZH) |
|----------|---------|------------------|------------------|
| `NEW_VAR` | `value` | English description | 中文描述 |

### Deprecated Environment Variables / 弃用环境变量

| Variable | Replacement | Removal Version |
|----------|-------------|-----------------|
| `OLD_VAR` | `NEW_VAR` | v{X}.{Y}.{Z} |

---

## 📋 Full Changelog / 完整变更日志

[View CHANGELOG.md](https://github.com/LessUp/go-live/blob/master/CHANGELOG.md)

---

## 🔗 Links / 相关链接

- [Documentation (EN)](https://lessup.github.io/go-live/en/)
- [文档（中文）](https://lessup.github.io/go-live/zh/)
- [GitHub Repository](https://github.com/LessUp/go-live)
- [Issue Tracker](https://github.com/LessUp/go-live/issues)

---

<!-- 
Template variables to replace / 模板变量替换:
- {VERSION} - e.g., 1.1.0
- {DATE} - e.g., 2025-04-16
- {PREVIOUS_VERSION} - e.g., 1.0.0
-->
