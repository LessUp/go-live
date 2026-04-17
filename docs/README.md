# Documentation

This directory contains user guides, tutorials, and developer documentation.

---

## Directory Structure

```
docs/
├── index.md              # Documentation landing page
├── api.md                # API reference (root)
├── design.md             # Design documentation (root)
├── usage.md              # Usage guide (root)
├── en/                   # English documentation
│   ├── index.md
│   ├── api.md
│   ├── design.md
│   └── usage.md
├── zh/                   # Chinese documentation (中文文档)
│   ├── index.md
│   ├── api.md
│   ├── design.md
│   └── usage.md
├── changelog/            # Changelog management
│   ├── README.md
│   ├── CHANGELOG_GUIDE.md
│   ├── RELEASE_WORKFLOW.md
│   ├── templates/
│   └── scripts/
└── assets/               # Static assets (CSS, images, etc.)
```

---

## Documentation Types

### User Guides
- Installation instructions
- Configuration reference
- Deployment guides
- Troubleshooting

### Developer Documentation
- Architecture overview
- Contributing guidelines
- Code conventions
- Testing guidelines

### API Reference
- REST API documentation
- OpenAPI/Swagger specs
- Authentication guide

---

## Viewing Documentation

Documentation is rendered via GitHub Pages at:
- **English**: https://lessup.github.io/go-live/en/
- **中文**: https://lessup.github.io/go-live/zh/

---

## Contributing to Documentation

1. Update both English (`/docs/en/`) and Chinese (`/docs/zh/`) versions
2. Follow the same structure for consistency
3. Test changes locally with Jekyll
4. Submit PR with documentation updates

---

## Specifications

Technical specifications are maintained in the `/specs` directory:
- [Product Specs](../specs/product/) - Feature requirements
- [RFCs](../specs/rfc/) - Technical designs
- [API Specs](../specs/api/) - OpenAPI definitions
- [Test Specs](../specs/testing/) - BDD test specs

See [specs/README.md](../specs/README.md) for details.
