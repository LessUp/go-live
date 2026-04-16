# Changelog 管理

本目录包含项目 Changelog 和发布相关的所有资源和模板。

---

## 目录结构

```
.changelog/
├── README.md                      # 本文件
├── CHANGELOG_GUIDE.md             # Changelog 编写规范
├── RELEASE_WORKFLOW.md            # 发布流程文档
├── templates/                     # 发布模板
│   ├── release-notes.md           # 双语发布说明模板
│   ├── release-notes-zh.md        # 中文发布说明模板
│   ├── unreleased-entry.md        # 未发布变更模板
│   └── version-comparison.md      # 版本对比模板
└── scripts/                       # 自动化脚本
    └── generate-release-notes.sh  # 生成发布说明脚本
```

---

## 工作流程

### 开发期间

1. **添加变更条目**
   - 在 `CHANGELOG.md` 的 `[Unreleased]` 部分添加条目
   - 使用适当的分类：Added、Changed、Fixed、Security 等
   - 遵循 [CHANGELOG_GUIDE.md](./CHANGELOG_GUIDE.md) 的规范

2. **使用模板**
   - 新功能：[unreleased-entry.md](./templates/unreleased-entry.md)

### 发布期间

详见 [RELEASE_WORKFLOW.md](./RELEASE_WORKFLOW.md)

**简要步骤**：
1. 更新 `CHANGELOG.md`，将 `[Unreleased]` 内容移至新版本
2. 创建发布分支并提交
3. 合并到 `main` 分支
4. 打标签：`git tag -a vX.Y.Z -m "Release vX.Y.Z"`
5. 推送标签：`git push origin vX.Y.Z`
6. 使用模板创建 GitHub Release

---

## 快速参考

### 分类定义

| 分类 | 用途 | 示例 |
|------|------|------|
| **Added** | 新增功能 | 添加 WebSocket 支持 |
| **Changed** | 功能变更 | 优化性能 |
| **Deprecated** | 即将移除的功能 | 标记废弃 API |
| **Removed** | 已移除功能 | 删除旧配置 |
| **Fixed** | Bug 修复 | 修复内存泄漏 |
| **Security** | 安全修复 | 修复漏洞 |

### 命名规范

```
v{MAJOR}.{MINOR}.{PATCH}[-{PRERELEASE}]

示例：
v1.0.0           # 稳定版本
v1.1.0-beta.1    # Beta 版本
v2.0.0-rc.1      # RC 版本
```

### 提交格式

```markdown
✅ Add WebSocket support (#100)
✅ Fix memory leak in track fanout (#99)

❌ Added WebSocket support
❌ fixed memory leak
```

---

## 相关文件

- `/CHANGELOG.md` - 主 Changelog 文件
- `/.github/workflows/ci.yml` - CI 流水线
- `/docs/` - 文档站点

---

## 参考规范

- [Keep a Changelog](https://keepachangelog.com/)
- [Semantic Versioning](https://semver.org/)
