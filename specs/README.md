# Specifications

This directory contains all project specifications following **Spec-Driven Development (SDD)** methodology.

All code implementation must be based on the specifications defined here. The `/specs` directory is the **Single Source of Truth** for product requirements, technical designs, API definitions, and database schemas.

---

## Directory Structure

```
specs/
├── product/            # Product requirements & feature definitions (PRDs)
├── rfc/                # Technical designs & architecture (RFCs)
├── api/                # API specifications (OpenAPI, GraphQL, etc.)
├── db/                 # Database schema definitions
└── testing/            # BDD test specifications (Gherkin .feature files)
```

---

## Specification Types

### `/specs/product/` - Product Specifications

Product Requirement Documents (PRDs) defining:
- Feature requirements
- User stories and use cases
- Acceptance criteria
- Business logic rules

**Naming**: `{feature-name}.md`

### `/specs/rfc/` - Request for Comments (Technical Designs)

Architecture and technical design documents:
- System architecture decisions
- Technology selection rationale
- Design patterns and conventions
- Performance considerations

**Naming**: `{NNNN}-{short-title}.md` (e.g., `0001-core-architecture.md`)

### `/specs/api/` - API Specifications

Machine-readable and human-readable API definitions:
- REST API (OpenAPI/Swagger)
- WebSocket protocols
- GraphQL schemas

### `/specs/db/` - Database Specifications

Database schema definitions:
- SQL schemas
- DBML files
- Migration specifications

### `/specs/testing/` - Test Specifications

Behavior-Driven Development (BDD) test specifications:
- Gherkin `.feature` files
- Acceptance test scenarios
- Integration test specifications

---

## How to Contribute

1. **New Feature**: Create a product spec in `/specs/product/`
2. **Technical Design**: Create an RFC in `/specs/rfc/`
3. **API Changes**: Update `/specs/api/` first, then implement
4. **Database Changes**: Update `/specs/db/` first, then migrate

See [CONTRIBUTING.md](../CONTRIBUTING.md) for detailed guidelines.

---

## Specification-Driven Workflow

When implementing any feature or making changes:

1. **Review existing specs** - Check if specs already exist
2. **Create/update specs first** - Define the specification before coding
3. **Implement based on specs** - Code must follow specs exactly
4. **Verify against specs** - Tests must validate spec compliance

**Never skip the spec definition step.**
