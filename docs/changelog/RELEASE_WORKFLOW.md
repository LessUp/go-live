# Release Process

This document describes the version release process for the live-webrtc-go project.

---

## Table of Contents

- [Pre-Release Checklist](#pre-release-checklist)
- [Version Numbering Rules](#version-numbering-rules)
- [Release Steps](#release-steps)
- [Post-Release](#post-release)
- [Hotfix Process](#hotfix-process)

---

## Pre-Release Checklist

### Code Checks

- [ ] All tests pass (`make test-all`)
- [ ] Code coverage meets target (>70%)
- [ ] Static analysis passes (`make lint`)
- [ ] Security scan passes (`make security`)
- [ ] Documentation is updated

### Version Checks

- [ ] Version number follows [SemVer](https://semver.org/)
- [ ] CHANGELOG.md is updated
- [ ] Version links in changelog are updated

### Release Notes

- [ ] English release notes are prepared
- [ ] Chinese release notes are prepared
- [ ] Known issues are documented
- [ ] Upgrade guide is prepared

---

## Version Numbering Rules

This project follows [Semantic Versioning](https://semver.org/):

### Version Format

```
MAJOR.MINOR.PATCH[-PRERELEASE][+BUILD]

Examples:
1.0.0           # Stable release
1.1.0-beta.1    # Beta pre-release
2.0.0-rc.1      # RC pre-release
```

### Version Increment Rules

| Version Type | When to Increment | Example |
|--------------|-------------------|---------|
| **MAJOR** | Breaking API changes | Remove endpoints, change response format |
| **MINOR** | Backward-compatible features | Add new endpoints |
| **PATCH** | Backward-compatible fixes | Bug fixes, security patches |

### Pre-release Versions

| Suffix | Meaning | Purpose |
|--------|---------|---------|
| `-alpha.N` | Internal testing | Early development testing |
| `-beta.N` | Public testing | Feature freeze testing |
| `-rc.N` | Release candidate | Pre-release verification |

---

## Release Steps

### Step 1: Create Release Branch

```bash
# Create release branch from main
git checkout -b release/v1.1.0
```

### Step 2: Update Version Information

#### 2.1 Update CHANGELOG.md

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

#### 2.2 Update Version Links

Add at the bottom of CHANGELOG.md:

```markdown
[Unreleased]: https://github.com/LessUp/go-live/compare/v1.1.0...HEAD
[1.1.0]: https://github.com/LessUp/go-live/releases/tag/v1.1.0
```

### Step 3: Commit Changes

```bash
git add CHANGELOG.md
git commit -m "chore(release): prepare for v1.1.0"
git push origin release/v1.1.0
```

### Step 4: Create Pull Request

Create a PR from `release/v1.1.0` to `main` with title:

```
Release v1.1.0
```

### Step 5: Merge and Tag

```bash
# After PR is merged, switch to main branch
git checkout main
git pull origin main

# Create tag
git tag -a v1.1.0 -m "Release v1.1.0"

# Push tag
git push origin v1.1.0
```

### Step 6: Create GitHub Release

Create bilingual Release using GitHub CLI:

```bash
gh release create v1.1.0 \
  --title "v1.1.0 - Release Title" \
  --notes-file docs/changelog/templates/release-notes-v1.1.0.md
```

Or use interactive mode:

```bash
gh release create v1.1.0
```

### Step 7: Publish Docker Image

```bash
# Build image
docker build -t ghcr.io/lessup/go-live:v1.1.0 .
docker build -t ghcr.io/lessup/go-live:latest .

# Push image
docker push ghcr.io/lessup/go-live:v1.1.0
docker push ghcr.io/lessup/go-live:latest
```

---

## Post-Release

### Verify Release

- [ ] GitHub Release page is accessible
- [ ] Release notes are complete
- [ ] Docker image is pullable
- [ ] Documentation site is updated

### Notification Channels

- [ ] Add new version badge in README
- [ ] Send notifications to communities (if applicable)
- [ ] Update project website

### Cleanup

- [ ] Delete release branch
- [ ] Close related issues
- [ ] Update milestones

---

## Hotfix Process

For critical issues requiring immediate fix:

### Step 1: Create Hotfix Branch

```bash
# Create branch from latest tag
git checkout -b hotfix/v1.1.1 v1.1.0
```

### Step 2: Fix and Test

```bash
# Fix the code
# ...

# Run tests
make test
```

### Step 3: Update CHANGELOG

```markdown
## [1.1.1] - 2025-04-17

### Fixed
- Fix critical security vulnerability (#110)
```

### Step 4: Quick Release

```bash
# Commit and push
git add .
git commit -m "fix: resolve critical security issue"
git push origin hotfix/v1.1.1

# Create PR, quick review and merge
# ...

# Tag and release
git checkout main
git pull origin main
git tag -a v1.1.1 -m "Hotfix v1.1.1"
git push origin v1.1.1

# Create Release
gh release create v1.1.1 --title "v1.1.1 - Security Hotfix" --notes "..."
```

---

## Related Templates

- [Release Notes Template](./templates/release-notes.md)
- [Chinese Release Notes Template](./templates/release-notes-zh.md)
- [Unreleased Entry Template](./templates/unreleased-entry.md)

---

## References

- [Semantic Versioning](https://semver.org/)
- [GitHub Releases](https://docs.github.com/en/repositories/releasing-projects-on-github/about-releases)
- [Keep a Changelog](https://keepachangelog.com/)
