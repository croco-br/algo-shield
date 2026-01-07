# Unit Testing Guide (Go & Vue + TypeScript)

**Audience:** AI agents responsible for generating, maintaining, or refactoring unit tests.

**Goal:** Ensure tests are **fast, deterministic, reliable, and maintainable**, with explicit constraints on performance, flakiness, and mocking strategy.

---

## 1. Core Principles (Non‑Negotiable)

1. **Unit tests must be deterministic**

   * No dependency on real time, network, filesystem, randomness, or global state.
   * Any non-determinism must be explicitly mocked or controlled.

2. **Tests must be fast**

   * Target execution time:

     * Individual test: **< 10ms**
     * Entire test suite (unit only): **< 1s** locally
   * Slow tests indicate incorrect test scope (integration or E2E).

3. **One logical behavior per test**

   * Each test should validate a single responsibility or rule.

4. **Tests are code**

   * Apply the same standards as production code: readability, naming, structure, and refactoring.

5. **No shared mutable state between tests**

   * Tests must be order‑independent and parallel‑safe.

---

## 2. Go – Unit Testing Guidelines

### 2.1 Mandatory Libraries

Use **only** the following libraries unless explicitly instructed otherwise:

* **Testing framework**

  * `testing` (standard library)

* **Assertions** (fast and widely adopted)

  * `github.com/stretchr/testify/assert`
  * `github.com/stretchr/testify/require`

* **Mocking**

  * `go.uber.org/mock/gomock`
  * `mockgen` for mock generation

* **Flaky test detection (CI)**

  * `go test -count=100` (repeat execution)

Avoid:

* Reflection-heavy or DSL-based testing frameworks
* Custom assertion libraries

---

### 2.2 Test Structure (Go)

Follow **Arrange → Act → Assert** explicitly.

```go
func Test_Service_DoSomething_WhenCondition_ThenExpectedResult(t *testing.T) {
    // Arrange
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    repo := NewMockRepository(ctrl)
    repo.EXPECT().FindByID("123").Return(Entity{ID: "123"}, nil)

    svc := NewService(repo)

    // Act
    result, err := svc.DoSomething("123")

    // Assert
    require.NoError(t, err)
    assert.Equal(t, "expected", result.Value)
}
```

Naming pattern:

```
Test_<Unit>_<Scenario>_<ExpectedOutcome>
```

---

### 2.3 Mocking Rules (Go)

1. **Mock only interfaces**, never concrete implementations
2. **Mock at architectural boundaries**

   * Repositories
   * External services
   * Time, UUID, randomness
3. **Do not mock value objects or pure functions**
4. Prefer **constructor injection** to enable mocking

Correct:

```go
func NewService(repo Repository) *Service
```

Incorrect:

```go
func NewService() *Service // hides dependencies
```

---

### 2.4 Performance Constraints (Go)

* Avoid `time.Sleep`
* Avoid real crypto, IO, or JSON encoding unless essential
* Use table‑driven tests for multiple cases

```go
func Test_Validator(t *testing.T) {
    tests := []struct {
        name  string
        input string
        valid bool
    }{
        {"valid", "abc", true},
        {"invalid", "", false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            assert.Equal(t, tt.valid, Validate(tt.input))
        })
    }
}
```

---

### 2.5 Flaky Test Detection (Go)

AI agents **must ensure**:

* Tests pass with:

  ```bash
  go test ./... -count=50
  ```
* Tests are safe with:

  ```bash
  go test ./... -race
  ```

Any failure under repetition or race detection indicates a **bug in the test or code**.

---

## 3. Vue + TypeScript – Unit Testing Guidelines

### 3.1 Mandatory Libraries

Use the following stack (optimized for speed and reliability):

* **Test runner**

  * `vitest`

* **Component testing**

  * `@vue/test-utils`

* **DOM environment**

  * `happy-dom` (preferred over jsdom for performance)

* **Mocking / Spies**

  * Built‑in `vi.mock`, `vi.fn`, `vi.spyOn`

Avoid:

* Jest (slower, legacy)
* Enzyme

---

### 3.2 Test Structure (Vue)

```ts
import { mount } from '@vue/test-utils'
import { describe, it, expect, vi } from 'vitest'
import MyComponent from './MyComponent.vue'

describe('MyComponent', () => {
  it('renders the title when prop is provided', () => {
    const wrapper = mount(MyComponent, {
      props: { title: 'Hello' }
    })

    expect(wrapper.text()).toContain('Hello')
  })
})
```

Naming pattern:

```
<unit>.spec.ts
```

Test names must describe **behavior**, not implementation.

---

### 3.3 Mocking Rules (Vue + TS)

1. Mock **external dependencies only**

   * HTTP clients
   * Stores
   * Browser APIs

2. Do **not** mock:

   * Vue reactivity
   * Component internal methods unless unavoidable

3. Preferred mocking patterns:

```ts
vi.mock('@/services/api', () => ({
  fetchUser: vi.fn()
}))
```

```ts
const spy = vi.spyOn(console, 'error').mockImplementation(() => {})
```

4. Reset mocks between tests:

```ts
afterEach(() => {
  vi.clearAllMocks()
})
```

---

### 3.4 Performance Constraints (Vue)

* Prefer `mount` only when needed
* Use `shallowMount` for logic-only tests
* Avoid full DOM rendering when testing computed values or methods

Target:

* Individual test: **< 5ms**
* Full unit suite: **< 1s**

---

### 3.5 Flaky Test Detection (Vue)

AI agents must validate:

```bash
vitest run --repeat 20
```

Additional rules:

* Never rely on `setTimeout` or real timers
* Use fake timers when required:

```ts
vi.useFakeTimers()
vi.runAllTimers()
```

* Do not rely on DOM ordering unless explicitly asserted

---

## 4. Cross‑Cutting Rules for AI Agents

AI agents **must not**:

* Introduce sleeps, retries, or waits to “fix” flaky tests
* Silence errors without assertion
* Mock more than one architectural layer per test

AI agents **must**:

* Prefer pure functions and dependency injection
* Flag tests that are slow, flaky, or overly coupled
* Suggest refactoring production code if testability is poor

---

## 5. Definition of a Good Unit Test

A unit test is considered **valid** if:

* It runs in isolation
* It is deterministic across multiple runs
* It completes in milliseconds
* It clearly communicates intent
* It fails for exactly one reason

If any of these conditions are violated, the test must be rewritten.

---

**End of document.**
