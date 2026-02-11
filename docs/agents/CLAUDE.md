# Testing Conventions

## Requirements

- **Coverage minimum: 80%** for new/modified code
- **Race detection: mandatory** — always use `-race` flag for Go tests
- **Flaky detection:** `go test -count=50` / `npm run test -- --repeat=20`
- **Deterministic:** no dependency on real time, network, filesystem, or randomness

## Patterns

- **AAA pattern** (Arrange-Act-Assert) — NEVER add `// Arrange`, `// Act`, `// Assert` comments
- **Table-driven tests** — consolidate 3+ similar tests into table-driven
- **Unit tests:** mock dependencies, one behavior per test
- **Integration tests:** `//go:build integration`, testcontainers-go, <500ms per test
- **UI tests:** Vitest + Vue Test Utils, `flushPromises()` for async, <300ms per test

## Anti-Patterns (Do NOT Generate)

- Trivial constructor tests (no logic = no test)
- Interface implementation tests (compiler checks this)
- Duplicate validation across layers (test once at service level)
- `exists()` without verifying actual values
- Tests without meaningful assertions

## Go-Specific

- Mock files (`mock_*_test.go`) are manually maintained, not auto-generated
- Use `Mock*` prefix for mock structs
- When interfaces change, update mocks by hand

## Vue/TS-Specific

- Keep `vi.mock('@/lib/api', ...)` blocks in sync across all test files when `api.ts` changes
- Update `vitest.setup.ts` (global mock) when API exports change

## Full Guides

- **Unit testing:** `docs/agents/unit-test.md` (1,945 lines — Go & Vue/TS)
- **Integration testing:** `docs/agents/integration-test.md` (1,375 lines — testcontainers, real DB)
