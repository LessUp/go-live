---
name: verify
description: Run lint and unit tests to verify changes. Use after making code changes to catch issues early.
---

Run the project's quick verification pipeline:

1. Run `make lint` — formats, vets, and lints all Go code
2. Run `make test-unit` — runs unit tests with race detection

Report any failures clearly. If all pass, confirm briefly.
