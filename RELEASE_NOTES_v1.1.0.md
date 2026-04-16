# Release v1.1.0 - 2025-04-16

## 🎉 Release Highlights

### English

This release focuses on comprehensive documentation improvements, introducing bilingual (English/Chinese) documentation, professional changelog management, and an enhanced GitHub Pages site.

### 中文

本次发布专注于全面的文档改进，推出了双语（英文/中文）文档、专业的 Changelog 管理系统以及增强的 GitHub Pages 站点。

---

## ✨ What's New / 新特性

### Documentation / 文档

#### English
- **Bilingual Documentation Site** - Complete documentation in English and Chinese
  - English docs at `/docs/en/` with usage guide, design docs, and API reference
  - Chinese docs at `/docs/zh/` with 使用指南, 设计说明, and API 参考
  - Language switcher for seamless navigation
- **Professional Changelog System** - Enhanced changelog management
  - Detailed changelog writing guidelines (CHANGELOG_GUIDE.md)
  - Complete release process documentation (RELEASE_WORKFLOW.md)
  - Bilingual release note templates

#### 中文
- **双语文档站点** - 完整的中英文文档
  - 英文文档位于 `/docs/en/`，包含使用指南、设计说明、API 参考
  - 中文文档位于 `/docs/zh/`，包含使用指南、设计说明、API 参考
  - 语言切换器实现无缝导航
- **专业 Changelog 系统** - 增强的变更日志管理
  - 详细的 Changelog 编写规范 (CHANGELOG_GUIDE.md)
  - 完整的发布流程文档 (RELEASE_WORKFLOW.md)
  - 双语发布说明模板

### Observability / 可观测性

#### English
- OpenTelemetry tracing support for distributed observability
  - Configurable via `OTEL_EXPORTER_OTLP_ENDPOINT` and `OTEL_SERVICE_NAME`
  - Supports both stdout and OTLP (grpc/http) exporters

#### 中文
- OpenTelemetry 追踪支持，用于分布式可观测性
  - 通过 `OTEL_EXPORTER_OTLP_ENDPOINT` 和 `OTEL_SERVICE_NAME` 配置
  - 支持 stdout 和 OTLP (grpc/http) 导出器

### Code Quality / 代码质量

#### English
- `HTTPHandlers.Close()` method for graceful shutdown
- Unit tests for `internal/uploader` package (23.5% coverage)
- Unit tests for `internal/otel` package (47.4% coverage)

#### 中文
- `HTTPHandlers.Close()` 方法实现优雅关闭
- `internal/uploader` 包的单元测试（覆盖率 23.5%）
- `internal/otel` 包的单元测试（覆盖率 47.4%）

---

## 🔧 Changes / 变更

### English
- **Restructured README.md** - More professional structure with quick navigation
- **Enhanced Quick Start** - Clearer installation and setup instructions
- JWT parser refactored with consolidated options pattern
- Enhanced server shutdown with proper error logging

### 中文
- **重构 README.md** - 更专业的结构，快速导航
- **增强快速开始** - 更清晰的安装和配置说明
- JWT 解析器重构，使用统一的选项模式
- 增强服务器关闭，正确记录错误

---

## 🐛 Bug Fixes / 问题修复

### English
- **Resource leak**: Rate limiter GC goroutine now properly stopped on shutdown
- Silent error ignore on server shutdown - now logs errors properly

### 中文
- **资源泄漏**：限流器 GC goroutine 现在在关闭时正确停止
- 服务器关闭时的错误静默忽略 - 现在正确记录错误

---

## 🚀 Upgrade Guide / 升级指南

### Docker / Docker 部署

```bash
# Pull the latest image / 拉取最新镜像
docker pull ghcr.io/lessup/go-live:v1.1.0
```

### Binary / 二进制

```bash
# Download from GitHub Releases / 从 GitHub Releases 下载
curl -LO https://github.com/LessUp/go-live/releases/download/v1.1.0/live-webrtc-go-linux-amd64

# Or build from source / 或从源码构建
git clone https://github.com/LessUp/go-live.git
cd go-live
git checkout v1.1.0
go build -o live-webrtc-go ./cmd/server
```

### Kubernetes / Kubernetes

```yaml
# Update image tag / 更新镜像标签
image: ghcr.io/lessup/go-live:v1.1.0
```

---

## 📋 Full Changelog / 完整变更日志

[View CHANGELOG.md](https://github.com/LessUp/go-live/blob/master/CHANGELOG.md)

---

## 🔗 Links / 相关链接

- [Documentation (EN)](https://lessup.github.io/go-live/en/)
- [文档（中文）](https://lessup.github.io/go-live/zh/)
- [GitHub Repository](https://github.com/LessUp/go-live)
- [Issue Tracker / 问题反馈](https://github.com/LessUp/go-live/issues)

---

**Contributors / 贡献者**: @LessUp

Special thanks to everyone who contributed to this release!
感谢所有为本次发布做出贡献的人！
