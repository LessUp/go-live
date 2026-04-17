# Changelog Management

This directory contains project changelog and release resources and templates.

---

## Directory Structure

```
docs/changelog/
├── README.md                      # This file
├── CHANGELOG_GUIDE.md             # Changelog writing guidelines
├── RELEASE_WORKFLOW.md            # Release process documentation
├── templates/                     # Release templates
│   ├── release-notes.md           # Bilingual release notes template
│   ├── release-notes-zh.md        # Chinese release notes template
│   ├── unreleased-entry.md        # Unreleased changelog template
│   └── version-comparison.md      # Version comparison template
└── scripts/                       # Automation scripts
    └── generate-release-notes.sh  # Release notes generation script
```

---

## Workflow

### During Development

1. **Add Change Entries**
   - Add entries under the `[Unreleased]` section in `CHANGELOG.md`
   - Use appropriate categories: Added, Changed, Fixed, Security, etc.
   - Follow the guidelines in [CHANGELOG_GUIDE.md](./CHANGELOG_GUIDE.md)

2. **Use Templates**
   - For new features: [unreleased-entry.md](./templates/unreleased-entry.md)

### During Release

See [RELEASE_WORKFLOW.md](./RELEASE_WORKFLOW.md)

**Quick Steps**:
1. Update `CHANGELOG.md`, move `[Unreleased]` content to new version
2. Create release branch and commit
3. Merge to `main` branch
4. Tag: `git tag -a vX.Y.Z -m "Release vX.Y.Z"`
5. Push tag: `git push origin vX.Y.Z`
6. Create GitHub Release using templates

---

## Quick Reference

### Category Definitions

| Category | Purpose | Example |
|----------|---------|---------|
| **Added** | New features | Add WebSocket support |
| **Changed** | Feature changes | Optimize performance |
| **Deprecated** | Soon-to-be-removed features | Mark deprecated API |
| **Removed** | Removed features | Delete old config |
| **Fixed** | Bug fixes | Fix memory leak |
| **Security** | Security fixes | Fix vulnerability |

### Version Naming

```
v{MAJOR}.{MINOR}.{PATCH}[-{PRERELEASE}]

Examples:
v1.0.0           # Stable release
v1.1.0-beta.1    # Beta release
v2.0.0-rc.1      # RC release
```

### Commit Format

```markdown
✅ Add WebSocket support (#100)
✅ Fix memory leak in track fanout (#99)

❌ Added WebSocket support
❌ fixed memory leak
```

---

## Related Files

- `/CHANGELOG.md` - Main changelog file
- `/.github/workflows/ci.yml` - CI pipeline
- `/docs/` - Documentation site

---

## References

- [Keep a Changelog](https://keepachangelog.com/)
- [Semantic Versioning](https://semver.org/)
