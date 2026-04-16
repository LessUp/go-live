# 发布流程文档

本文档详细说明 live-webrtc-go 项目的版本发布流程。

---

## 目录

- [发布前检查](#发布前检查)
- [版本号规则](#版本号规则)
- [发布步骤](#发布步骤)
- [发布后事项](#发布后事项)
- [紧急修复流程](#紧急修复流程)

---

## 发布前检查

### 代码检查

- [ ] 所有测试通过 (`make test-all`)
- [ ] 代码覆盖率达标 (>70%)
- [ ] 静态分析无错误 (`make lint`)
- [ ] 安全扫描通过 (`make security`)
- [ ] 文档已更新

### 版本检查

- [ ] 版本号符合 [SemVer](https://semver.org/) 规范
- [ ] CHANGELOG.md 已更新
- [ ] `.changelog/README.md` 中的版本链接已更新

### 发布说明准备

- [ ] 英文发布说明已准备
- [ ] 中文发布说明已准备
- [ ] 已知问题已记录
- [ ] 升级指南已准备

---

## 版本号规则

本项目遵循 [Semantic Versioning](https://semver.org/)：

### 版本格式

```
MAJOR.MINOR.PATCH[-PRERELEASE][+BUILD]

Examples:
1.0.0           # 稳定版本
1.1.0-beta.1    # Beta 预发布版本
2.0.0-rc.1      # RC 预发布版本
```

### 版本递增规则

| 版本类型 | 递增时机 | 示例 |
|----------|----------|------|
| **MAJOR** | 破坏性 API 变更 | 移除端点、响应格式变更 |
| **MINOR** | 向后兼容的功能添加 | 新增端点、新功能 |
| **PATCH** | 向后兼容的问题修复 | Bug 修复、安全补丁 |

### 预发布版本

| 后缀 | 含义 | 用途 |
|------|------|------|
| `-alpha.N` | 内部测试版本 | 开发早期测试 |
| `-beta.N` | 公开测试版本 | 功能冻结后的测试 |
| `-rc.N` | 候选发布版本 | 正式发布前的验证 |

---

## 发布步骤

### 步骤 1: 创建发布分支

```bash
# 从 main 分支创建发布分支
git checkout -b release/v1.1.0
```

### 步骤 2: 更新版本信息

#### 2.1 更新 CHANGELOG.md

```markdown
## [1.1.0] - 2025-04-16

### Added
- Add new feature (#100)

### Changed
- Improve performance (#99)

### Fixed
- Fix bug (#98)

## [Unreleased]

[Empty - new entries go here]
```

#### 2.2 更新版本链接

在 CHANGELOG.md 底部添加：

```markdown
[Unreleased]: https://github.com/LessUp/go-live/compare/v1.1.0...HEAD
[1.1.0]: https://github.com/LessUp/go-live/releases/tag/v1.1.0
```

### 步骤 3: 提交变更

```bash
git add CHANGELOG.md
git commit -m "chore(release): prepare for v1.1.0"
git push origin release/v1.1.0
```

### 步骤 4: 创建 Pull Request

创建 PR 从 `release/v1.1.0` 到 `main`，标题：

```
Release v1.1.0
```

### 步骤 5: 合并并打标签

```bash
# PR 合并后，切换到 main 分支
git checkout main
git pull origin main

# 创建标签
git tag -a v1.1.0 -m "Release v1.1.0"

# 推送标签
git push origin v1.1.0
```

### 步骤 6: 创建 GitHub Release

使用 GitHub CLI 创建双语 Release：

```bash
gh release create v1.1.0 \
  --title "v1.1.0 - Release Title" \
  --notes-file .changelog/templates/release-notes-v1.1.0.md
```

或使用交互式方式：

```bash
gh release create v1.1.0
```

### 步骤 7: 发布 Docker 镜像

```bash
# 构建镜像
docker build -t ghcr.io/lessup/go-live:v1.1.0 .
docker build -t ghcr.io/lessup/go-live:latest .

# 推送镜像
docker push ghcr.io/lessup/go-live:v1.1.0
docker push ghcr.io/lessup/go-live:latest
```

---

## 发布后事项

### 验证发布

- [ ] GitHub Release 页面可访问
- [ ] 发布说明完整准确
- [ ] Docker 镜像可拉取
- [ ] 文档站点已更新

### 通知渠道

- [ ] 在 README 中添加新版本徽章
- [ ] 发送通知到相关社区（如适用）
- [ ] 更新项目网站

### 清理工作

- [ ] 删除发布分支
- [ ] 关闭相关 Issue
- [ ] 更新里程碑

---

## 紧急修复流程

对于需要立即修复的严重问题：

### 步骤 1: 创建热修复分支

```bash
# 从最新标签创建分支
git checkout -b hotfix/v1.1.1 v1.1.0
```

### 步骤 2: 修复并测试

```bash
# 修复代码
# ...

# 运行测试
make test
```

### 步骤 3: 更新 CHANGELOG

```markdown
## [1.1.1] - 2025-04-17

### Fixed
- Fix critical security vulnerability (#110)
```

### 步骤 4: 快速发布

```bash
# 提交并推送
git add .
git commit -m "fix: resolve critical security issue"
git push origin hotfix/v1.1.1

# 创建 PR，快速审核合并
# ...

# 打标签并发布
git checkout main
git pull origin main
git tag -a v1.1.1 -m "Hotfix v1.1.1"
git push origin v1.1.1

# 创建 Release
gh release create v1.1.1 --title "v1.1.1 - Security Hotfix" --notes "..."
```

---

## 相关模板

- [发布说明模板](./templates/release-notes.md)
- [中文发布说明模板](./templates/release-notes-zh.md)
- [未发布变更模板](./templates/unreleased-entry.md)

---

## 参考链接

- [Semantic Versioning](https://semver.org/)
- [GitHub Releases](https://docs.github.com/en/repositories/releasing-projects-on-github/about-releases)
- [Keep a Changelog](https://keepachangelog.com/)
