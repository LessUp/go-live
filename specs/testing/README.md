# Testing Specifications

This directory contains Behavior-Driven Development (BDD) test specifications in Gherkin format.

---

## Test Spec Files

| File | Feature | Type |
|------|---------|------|
| [whip-publish.feature](./whip-publish.feature) | WHIP Publishing | Integration |
| [whep-playback.feature](./whep-playback.feature) | WHEP Playback | Integration |
| [authentication.feature](./authentication.feature) | Authentication | Security |
| [room-management.feature](./room-management.feature) | Room Management | Integration |

---

## Test Categories

### Unit Tests
- Located in `internal/*/` packages alongside code
- Table-driven test patterns
- Race detection enabled

### Integration Tests
- Located in `test/integration/`
- Require `-tags=integration` build tag
- Test HTTP API endpoints

### E2E Tests
- Located in `test/e2e/`
- Require `-tags=e2e` build tag
- Full WebRTC flow tests

### Security Tests
- Located in `test/security/`
- Auth bypass, rate limiting, input validation

### Performance Tests
- Located in `test/performance/`
- Benchmarks and load testing

---

## Gherkin Feature File Template

```gherkin
Feature: Feature Name
  Brief description of the feature

  Scenario: Scenario description
    Given some precondition
    When some action is taken
    Then some expected outcome occurs
```

---

## Running Tests

```bash
# All tests
make test

# Unit tests only
make test-unit

# Integration tests
make test-integration

# E2E tests
make test-e2e

# All tests (including e2e)
make test-all
```
