# Unreleased Entry Template / 未发布变更模板

Use this template to draft changelog entries before adding to `CHANGELOG.md`.

使用此模板起草 Changelog 条目，然后再添加到 `CHANGELOG.md`。

---

## Entry / 条目

### Category / 分类 (choose one / 选择一项)

- [ ] **Added** - 新功能 - New features, APIs, configuration options
- [ ] **Changed** - 变更 - Changes to existing functionality, improvements
- [ ] **Deprecated** - 弃用 - Features to be removed in future releases
- [ ] **Removed** - 移除 - Features removed in this release
- [ ] **Fixed** - 修复 - Bug fixes, error handling improvements
- [ ] **Security** - 安全 - Security fixes, vulnerability patches

### Description / 描述

<!-- Brief, imperative mood description / 简短的祈使句描述 -->
<!-- Example / 示例: "Add WebSocket support for real-time room events" -->

**English:**
Brief description of the change.

**中文：**
变更的简要描述。

### Details / 详情 (optional / 可选)

<!-- Additional context, rationale, or examples -->
<!-- 额外的上下文、原因或示例 -->

**English:**
- Why this change was made
- How it affects users
- Any migration steps needed

**中文：**
- 为何进行此变更
- 对用户的影响
- 所需的迁移步骤

### Grouping / 分组 (optional / 可选)

<!-- Suggest a subheading if this belongs to a group -->
<!-- 如果属于某个组，建议子标题 -->
<!-- Example / 示例: "Observability", "API", "Security", "Performance" -->

Suggested group / 建议分组: 

### Related / 相关

- Issue: #N
- PR: #N
- Discussion: #N

### Breaking Change / 破坏性变更

- [ ] This is a breaking change / 这是一个破坏性变更

If yes, describe / 如果是，请描述：
1. What breaks / 什么会被破坏
2. Migration path / 迁移路径
3. Deprecation timeline / 弃用时间表

### Configuration Changes / 配置变更

<!-- If this adds or changes environment variables -->
<!-- 如果此变更添加或修改了环境变量 -->

| Variable | Default | Description (EN) | Description (ZH) |
|----------|---------|------------------|------------------|
| `NEW_VAR` | `value` | Description | 描述 |

### Testing / 测试

- [ ] Unit tests added / 单元测试已添加
- [ ] Integration tests added / 集成测试已添加
- [ ] Manual testing performed / 已进行手动测试

---

## Examples / 示例

### Example 1: New Feature / 新功能示例

```markdown
### Added

#### Observability
- Add OpenTelemetry tracing support for distributed tracing (#50)
```

### Example 2: Bug Fix / Bug 修复示例

```markdown
### Fixed

- Fix memory leak in track fanout when subscriber disconnects (#44)
```

### Example 3: Breaking Change / 破坏性变更示例

```markdown
### Changed

- **Breaking:** Change API response format for `/api/rooms` endpoint (#60)
  - Old / 旧: `{"rooms": [...]}`
  - New / 新: `[...]`
  - Migration / 迁移: Update client code to handle array directly
```

### Example 4: Security Fix / 安全修复示例

```markdown
### Security

- Fix timing attack vulnerability in token comparison (#43)
  - Use `crypto/subtle.ConstantTimeCompare` for all token checks
  - CVE: Pending assignment
```

---

## How to Submit / 如何提交

1. Fill out this template / 填写此模板
2. Copy the formatted entry / 复制格式化后的条目
3. Add to the appropriate section in `CHANGELOG.md` under `[Unreleased]`
   在 `CHANGELOG.md` 的 `[Unreleased]` 下添加到适当部分
4. Submit your PR / 提交 PR

## After Merge / 合并后

The entry will be included in the next release. On release:
此条目将包含在下一次发布中。发布时：

1. Entries move from `[Unreleased]` to a new version section
   条目从 `[Unreleased]` 移动到新版本部分
2. GitHub release notes are generated
   生成 GitHub 发布说明
3. Version links are updated
   更新版本链接
