#!/bin/bash
#
# 生成发布说明脚本
# Usage: ./generate-release-notes.sh <version> <date>
# Example: ./generate-release-notes.sh 1.1.0 2025-04-16
#

set -e

VERSION="${1:-}"
DATE="${2:-$(date +%Y-%m-%d)}"

if [ -z "$VERSION" ]; then
    echo "Usage: $0 <version> [date]"
    echo "Example: $0 1.1.0"
    echo "Example: $0 1.1.0 2025-04-16"
    exit 1
fi

# 移除版本号前的 v（如果有）
VERSION=$(echo "$VERSION" | sed 's/^v//')

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEMPLATE_DIR="$(dirname "$SCRIPT_DIR")/templates"
OUTPUT_FILE="release-notes-v${VERSION}.md"

echo "Generating release notes for v${VERSION} (${DATE})..."

# 检查模板文件
if [ ! -f "$TEMPLATE_DIR/release-notes.md" ]; then
    echo "Error: Template file not found: $TEMPLATE_DIR/release-notes.md"
    exit 1
fi

# 读取 CHANGELOG 内容
CHANGELOG_FILE="$(dirname "$(dirname "$SCRIPT_DIR")")/CHANGELOG.md"
if [ -f "$CHANGELOG_FILE" ]; then
    # 提取当前版本的变更内容
    VERSION_CONTENT=$(awk "/^## \\[${VERSION}\]/,/^## \\[Unreleased\]/" "$CHANGELOG_FILE" | head -n -1)
    
    # 提取 Added 部分
    ADDED=$(echo "$VERSION_CONTENT" | awk '/### Added/,/### (Changed|Deprecated|Removed|Fixed|Security|## )/' | head -n -1 | tail -n +2)
    
    # 提取 Changed 部分  
    CHANGED=$(echo "$VERSION_CONTENT" | awk '/### Changed/,/### (Deprecated|Removed|Fixed|Security|## )/' | head -n -1 | tail -n +2)
    
    # 提取 Fixed 部分
    FIXED=$(echo "$VERSION_CONTENT" | awk '/### Fixed/,/### (Security|## )/' | head -n -1 | tail -n +2)
    
    # 提取 Security 部分
    SECURITY=$(echo "$VERSION_CONTENT" | awk '/### Security/,/### ## /' | head -n -1 | tail -n +2)
fi

# 生成发布说明
cat > "$OUTPUT_FILE" << EOF
# Release v${VERSION} - ${DATE}

## 🎉 Release Highlights

### English

This is the v${VERSION} release of live-webrtc-go.

### 中文

这是 live-webrtc-go 的 v${VERSION} 版本发布。

---

## ✨ What's New / 新特性

### English
$(if [ -n "$ADDED" ]; then echo "$ADDED"; else echo "_No new features in this release._"; fi)

### 中文
$(if [ -n "$ADDED" ]; then echo "$ADDED" | sed 's/Add/添加/g' | sed 's/Support/支持/g'; else echo "_本次发布无新特性。_"; fi)

---

## 🔧 Changes / 变更

### English
$(if [ -n "$CHANGED" ]; then echo "$CHANGED"; else echo "_No changes in this release._"; fi)

### 中文
$(if [ -n "$CHANGED" ]; then echo "$CHANGED" | sed 's/Change/变更/g' | sed 's/Improve/改进/g'; else echo "_本次发布无变更。_"; fi)

---

## 🐛 Bug Fixes / 问题修复

### English
$(if [ -n "$FIXED" ]; then echo "$FIXED"; else echo "_No bug fixes in this release._"; fi)

### 中文
$(if [ -n "$FIXED" ]; then echo "$FIXED" | sed 's/Fix/修复/g'; else echo "_本次发布无问题修复。_"; fi)

---

## 🔒 Security / 安全
$(if [ -n "$SECURITY" ]; then echo "$SECURITY"; else echo "_No security fixes in this release._ / _本次发布无安全修复。_"; fi)

---

## 🚀 Upgrade Guide / 升级指南

### Docker / Docker 部署

\`\`\`bash
# Pull the latest image / 拉取最新镜像
docker pull ghcr.io/lessup/go-live:v${VERSION}
\`\`\`

### Binary / 二进制

\`\`\`bash
# Download from GitHub Releases / 从 GitHub Releases 下载
curl -LO https://github.com/LessUp/go-live/releases/download/v${VERSION}/live-webrtc-go-linux-amd64

# Or build from source / 或从源码构建
git clone https://github.com/LessUp/go-live.git
cd go-live
git checkout v${VERSION}
go build -o live-webrtc-go ./cmd/server
\`\`\`

### Kubernetes / Kubernetes

\`\`\`yaml
# Update image tag / 更新镜像标签
image: ghcr.io/lessup/go-live:v${VERSION}
\`\`\`

---

## 📋 Full Changelog / 完整变更日志

[View CHANGELOG.md](https://github.com/LessUp/go-live/blob/master/CHANGELOG.md)

---

## 🔗 Links / 相关链接

- [Documentation (EN)](https://lessup.github.io/go-live/en/)
- [文档（中文）](https://lessup.github.io/go-live/zh/)
- [GitHub Repository](https://github.com/LessUp/go-live)
- [Issue Tracker / 问题反馈](https://github.com/LessUp/go-live/issues)
EOF

echo "Release notes generated: $OUTPUT_FILE"
echo ""
echo "Next steps:"
echo "1. Review and edit $OUTPUT_FILE"
echo "2. Create GitHub release:"
echo "   gh release create v${VERSION} --title \"v${VERSION}\" --notes-file $OUTPUT_FILE"
