# 版本发布 v{VERSION} - {DATE}

## 🎉 发布亮点

本次发布的简要描述（1-2 句话）。

---

## ✨ 新特性

### 新功能
- ✨ 特性 1 - 简要说明
- ✨ 特性 2 - 简要说明

### 改进
- 🚀 改进 1 - 简要说明
- 🚀 改进 2 - 简要说明

---

## 🐛 问题修复

- 🐛 修复 1 - 描述 (#issue)
- 🐛 修复 2 - 描述 (#issue)

---

## 🔒 安全修复

<!-- 如无安全修复，删除此节 -->

- 🔒 安全修复 1 - 描述（如有 CVE 编号请标注）

---

## ⚠️ 破坏性变更

<!-- 如无破坏性变更，删除此节 -->

⚠️ **重要提示**：本次发布包含破坏性变更。

- **变更 1**：描述和迁移指南
- **变更 2**：描述和迁移指南

### 迁移指南

```bash
# 从前版本迁移的步骤
```

---

## 📦 依赖更新

| 依赖包 | 旧版本 | 新版本 | 原因 |
|--------|--------|--------|------|
| package | v1.0.0 | v1.1.0 | Bug 修复 |

---

## 🙏 贡献者

感谢所有为本次发布做出贡献的人！

- @username - 贡献描述

---

## 🚀 升级指南

### Docker 部署

```bash
# 拉取最新镜像
docker pull ghcr.io/lessup/go-live:v{VERSION}
```

### 二进制部署

```bash
# 从 GitHub Releases 下载
curl -LO https://github.com/LessUp/go-live/releases/download/v{VERSION}/live-webrtc-go-linux-amd64

# 或从源码构建
git clone https://github.com/LessUp/go-live.git
cd go-live
git checkout v{VERSION}
go build -o live-webrtc-go ./cmd/server
```

### Kubernetes

```yaml
# 更新镜像标签
image: ghcr.io/lessup/go-live:v{VERSION}
```

---

## 📝 配置变更

### 新增环境变量

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `NEW_VAR` | `value` | 描述 |

### 弃用环境变量

| 变量名 | 替代方案 | 删除版本 |
|--------|----------|----------|
| `OLD_VAR` | `NEW_VAR` | v{X}.{Y}.{Z} |

---

## 📋 完整变更日志

[查看 CHANGELOG.md](https://github.com/LessUp/go-live/blob/master/CHANGELOG.md)

---

## 🔗 相关链接

- [中文文档](https://lessup.github.io/go-live/zh/)
- [GitHub 仓库](https://github.com/LessUp/go-live)
- [问题反馈](https://github.com/LessUp/go-live/issues)

---

<!-- 
模板变量替换:
- {VERSION} - 例如：1.1.0
- {DATE} - 例如：2025-04-16
-->
