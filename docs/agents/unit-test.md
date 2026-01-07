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

6. **Follow AAA pattern implicitly**

   * Tests must follow **Arrange → Act → Assert** structure.
   * **NEVER add comments** marking each section (`// Arrange`, `// Act`, `// Assert`).
   * The structure should be self‑evident through blank lines and code organization.

7. **Idiomatic separation of concerns**

   * Tests should validate units that have clear, single responsibilities.
   * If a unit is hard to test, it likely violates separation of concerns.
   * Production code must be designed for testability (dependency injection, interfaces, pure functions).

---

## 2. Coverage Requirements

**Minimum coverage threshold: 80%**

* All new code must maintain or improve overall coverage
* Coverage is measured automatically in CI via `go test -coverprofile`
* Coverage reports are uploaded to Codecov for tracking
* AI agents must verify coverage AND lint compliance before completing tasks:

```bash
# Go - verify linting first
go vet ./...
golangci-lint run ./...

# Go - then check coverage
go test -coverprofile=coverage.txt -covermode=atomic ./...
go tool cover -func=coverage.txt | grep total

# TypeScript/Vue - verify linting first
npm run lint

# TypeScript/Vue - then check coverage
npm run test:coverage
```

* If coverage drops below 80%, agent must add tests or ask for clarification
* Focus on meaningful coverage: test logic, not just lines
* Fix all lint errors before checking coverage

---

## 3. Go – Unit Testing Guidelines

### 3.1 Mandatory Libraries

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

### 3.2 Test Structure (Go)

Follow **Arrange → Act → Assert** pattern **without comments**.

Use blank lines to visually separate the three sections.

```go
func Test_Service_DoSomething_WhenCondition_ThenExpectedResult(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    repo := NewMockRepository(ctrl)
    repo.EXPECT().FindByID("123").Return(Entity{ID: "123"}, nil)
    svc := NewService(repo)

    result, err := svc.DoSomething("123")

    require.NoError(t, err)
    assert.Equal(t, "expected", result.Value)
}
```

**Critical**: Never add comments like `// Arrange`, `// Act`, `// Assert`. The structure must be self‑evident.

Naming pattern:

```
Test_<Unit>_<Scenario>_<ExpectedOutcome>
```

**Idiomatic Go testing**:

* Keep setup minimal and focused
* Use table‑driven tests for multiple scenarios
* Avoid test helpers that obscure test logic

---

### 3.3 Idiomatic Separation of Concerns (Go)

Well‑designed Go code should follow these principles:

1. **Single Responsibility**

   * Each type/function does one thing well
   * Services orchestrate, repositories persist, handlers adapt

2. **Dependency Injection via constructors**

   ```go
   // Good: dependencies are explicit
   type Service struct {
       repo Repository
       clock Clock
   }

   func NewService(repo Repository, clock Clock) *Service {
       return &Service{repo: repo, clock: clock}
   }
   ```

   ```go
   // Bad: hidden dependencies, hard to test
   func NewService() *Service {
       return &Service{
           repo: NewPostgresRepo(), // ❌ concrete dependency
           clock: time.Now,         // ❌ global state
       }
   }
   ```

3. **Interface segregation**

   * Define small interfaces at the point of use
   * Consumers define interfaces, not providers

   ```go
   // In service package
   type Repository interface {
       FindByID(id string) (Entity, error)
   }
   ```

4. **Pure functions when possible**

   * Prefer stateless functions for business logic
   * Easier to test, easier to reason about

   ```go
   func CalculateDiscount(price float64, tier string) float64 {
       // No dependencies, no state, pure logic
   }
   ```

---

### 3.4 Mocking Rules (Go)

1. **Mock only interfaces**, never concrete implementations
2. **Mock at architectural boundaries**

   * Repositories (data access)
   * External services (HTTP, gRPC, message queues)
   * Non‑deterministic sources (time, UUID, randomness)

3. **Do not mock value objects or pure functions**

   * If it's deterministic and has no I/O, test it directly

4. **One mock per architectural layer**

   * Don't mock both the repository AND the database client
   * Mock at the highest reasonable abstraction

Correct:

```go
// Service depends on interface
type Service struct {
    repo Repository
}

func NewService(repo Repository) *Service {
    return &Service{repo: repo}
}
```

Incorrect:

```go
// Service creates its own dependencies
func NewService() *Service {
    return &Service{
        repo: postgres.NewRepository(), // ❌ hard-coded
    }
}
```

---

### 3.5 Performance Constraints (Go)

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

### 3.6 Flaky Test Detection (Go)

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

### 3.7 Context and Async (Go)

**Context handling:**

* Always pass `context.Context` as first parameter in tests
* Use `context.Background()` or `context.TODO()` for test contexts
* Test context cancellation and timeout behavior explicitly

```go
func Test_Service_WithContext_WhenCancelled_ThenReturnsError(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    ctx, cancel := context.WithCancel(context.Background())
    cancel() // Cancel immediately

    svc := NewService()

    err := svc.DoWork(ctx)

    assert.ErrorIs(t, err, context.Canceled)
}
```

**Goroutines and concurrency:**

* Use `sync.WaitGroup` or channels to synchronize goroutines in tests
* Test race conditions with `go test -race`
* Avoid `time.Sleep` for synchronization; use explicit signals

```go
func Test_Worker_ProcessesConcurrently(t *testing.T) {
    var wg sync.WaitGroup
    results := make(chan int, 10)

    wg.Add(10)
    for i := 0; i < 10; i++ {
        go func(n int) {
            defer wg.Done()
            results <- Process(n)
        }(i)
    }

    wg.Wait()
    close(results)

    assert.Len(t, results, 10)
}
```

---

### 3.8 Test Files and Organization (Go)

**File location:**

* Test files live in the same directory as the code under test
* Name pattern: `<filename>_test.go`
* Example: `service.go` → `service_test.go`

**Package naming:**

* Use same package name (`package foo`) to test unexported functions
* Use external package name (`package foo_test`) to test public API only
* Prefer internal package for unit tests (access to internals)
* Use external package for integration-style tests

**Mock files:**

* Generated mocks: `mock_<interface>_test.go` (gitignored if regenerated in CI)
* Manual mocks: `mock_<interface>.go` in test helpers package
* Keep mocks close to tests that use them

**Example structure:**

```
internal/auth/
├── service.go
├── service_test.go
├── mock_repository_test.go      # Generated by mockgen
└── mock_user_service_test.go    # Generated by mockgen
```

---

### 3.9 Test Data and Fixtures (Go)

**Inline test data:**

* Prefer inline data for small, simple cases
* Use table-driven tests for multiple scenarios

**Golden files (avoid unless necessary):**

* Use only for large, complex output (JSON, XML, etc.)
* Store in `testdata/` directory (Go convention)
* Update with `-update` flag pattern

**Builder/factory pattern:**

* Create test builders for complex domain objects
* Keep builders in test files or separate `testing` package

```go
func TestUserBuilder() *User {
    return &User{
        ID:    "test-123",
        Email: "test@example.com",
        Role:  "user",
    }
}

func (u *User) WithRole(role string) *User {
    u.Role = role
    return u
}

// Usage
user := TestUserBuilder().WithRole("admin")
```

---

### 3.10 Cleanup and Resource Management (Go)

**Use `t.Cleanup()` for automatic cleanup:**

```go
func Test_WithCleanup(t *testing.T) {
    file, err := os.CreateTemp("", "test")
    require.NoError(t, err)

    t.Cleanup(func() {
        os.Remove(file.Name())
    })

    // Test code here
}
```

**Benefits:**

* Runs even if test fails or panics
* Cleaner than `defer` in table-driven tests
* Stacks multiple cleanup functions in LIFO order

---

### 3.11 Assertions and Error Messages (Go)

**Descriptive assertions:**

```go
// Good: context helps debugging
assert.Equal(t, expected, actual, "user ID should match after creation")

// Bad: no context
assert.Equal(t, expected, actual)
```

**When to use `assert` vs `require`:**

* Use `assert` for non-critical checks (test continues)
* Use `require` when failure makes rest of test invalid (test stops)

```go
func Test_Example(t *testing.T) {
    user, err := repo.FindByID("123")
    require.NoError(t, err, "user must exist for test to proceed")

    assert.Equal(t, "John", user.Name, "name should match")
    assert.True(t, user.Active, "user should be active")
}
```

---

### 3.12 Mocking Time, Randomness, and Non-Determinism (Go)

**Time:**

* Define a `Clock` interface and inject it
* Mock in tests

```go
type Clock interface {
    Now() time.Time
}

type realClock struct{}
func (realClock) Now() time.Time { return time.Now() }

type fixedClock struct{ t time.Time }
func (c fixedClock) Now() time.Time { return c.t }

// In tests
fixedTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
clock := fixedClock{t: fixedTime}
```

**UUIDs and random sources:**

* Use interface for ID generation
* Inject predictable generator in tests

```go
type IDGenerator interface {
    Generate() string
}

// In tests
type fakeIDGen struct{ id string }
func (f fakeIDGen) Generate() string { return f.id }
```

---

## 4. Vue + TypeScript – Unit Testing Guidelines

### 4.1 Mandatory Libraries

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

### 4.2 Test Structure (Vue)

Follow **Arrange → Act → Assert** pattern **without comments**.

Use blank lines to visually separate the three sections.

```ts
import { mount } from '@vue/test-utils'
import { describe, it, expect, vi } from 'vitest'
import MyComponent from './MyComponent.vue'

describe('MyComponent', () => {
  it('renders the title when prop is provided', () => {
    const wrapper = mount(MyComponent, {
      props: { title: 'Hello' }
    })

    const text = wrapper.text()

    expect(text).toContain('Hello')
  })
})
```

**Critical**: Never add comments like `// Arrange`, `// Act`, `// Assert`. The structure must be self‑evident.

Naming pattern:

```
<unit>.spec.ts
```

Test names must describe **behavior**, not implementation.

**Idiomatic TypeScript testing**:

* Keep test setup minimal and declarative
* Extract complex setup into factory functions (not in beforeEach)
* Avoid obscuring what is being tested

---

### 4.3 Idiomatic Separation of Concerns (Vue + TypeScript)

Well‑designed Vue/TypeScript code should follow these principles:

1. **Component Single Responsibility**

   * Components handle presentation and user interaction
   * Business logic lives in composables or services
   * State management in stores (Pinia/Vuex)

   ```ts
   // Good: component delegates to composable
   <script setup lang="ts">
   import { useUserAuthentication } from '@/composables/useUserAuthentication'

   const { login, isLoading, error } = useUserAuthentication()
   </script>
   ```

   ```ts
   // Bad: business logic in component
   <script setup lang="ts">
   const login = async () => {
     const token = await fetch('/api/auth') // ❌ direct HTTP
     localStorage.setItem('token', token)   // ❌ side effects
     validateToken(token)                   // ❌ business logic
   }
   </script>
   ```

2. **Composables for reusable logic**

   * Extract stateful logic into composables
   * Composables should be testable independently

   ```ts
   // composables/useUserAuthentication.ts
   export function useUserAuthentication() {
     const { login: apiLogin } = useAuthApi()
     const isLoading = ref(false)

     const login = async (credentials: Credentials) => {
       isLoading.value = true
       try {
         return await apiLogin(credentials)
       } finally {
         isLoading.value = false
       }
     }

     return { login, isLoading }
   }
   ```

3. **Services for API/external communication**

   * Isolate HTTP clients and external dependencies
   * Services return plain data, no Vue reactivity

   ```ts
   // services/authApi.ts
   export class AuthApi {
     constructor(private client: HttpClient) {}

     async login(credentials: Credentials): Promise<AuthResponse> {
       return this.client.post('/auth/login', credentials)
     }
   }
   ```

4. **Pure utility functions**

   * Stateless transformations and calculations
   * No dependencies, fully testable

   ```ts
   // utils/formatters.ts
   export function formatCurrency(amount: number, locale: string): string {
     // Pure function, easy to test
   }
   ```

---

### 4.4 Mocking Rules (Vue + TS)

1. **Mock external dependencies only**

   * HTTP clients
   * Stores (when testing components)
   * Browser APIs (localStorage, fetch, etc.)

2. **Do not mock**

   * Vue reactivity system
   * Component internal methods unless unavoidable
   * Computed properties (test them through outputs)

3. **Preferred mocking patterns**

```ts
vi.mock('@/services/api', () => ({
  fetchUser: vi.fn()
}))
```

```ts
const spy = vi.spyOn(console, 'error').mockImplementation(() => {})
```

4. **Reset mocks between tests**

```ts
afterEach(() => {
  vi.clearAllMocks()
})
```

5. **One mock per architectural layer**

   * Don't mock both the composable AND the service it uses
   * Mock at the boundary (API client, not the composable wrapper)

---

### 4.5 Performance Constraints (Vue)

* Prefer `mount` only when needed
* Use `shallowMount` for logic-only tests
* Avoid full DOM rendering when testing computed values or methods

Target:

* Individual test: **< 5ms**
* Full unit suite: **< 1s**

---

### 4.6 Flaky Test Detection (Vue)

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

### 4.7 Async and Promises (Vue + TS)

**Testing async code:**

* Always use `async/await` in test functions
* Use `flushPromises()` from `@vue/test-utils` to wait for pending promises

```ts
import { flushPromises } from '@vue/test-utils'

it('loads user data', async () => {
  const wrapper = mount(UserProfile, {
    props: { userId: '123' }
  })

  await flushPromises()

  expect(wrapper.text()).toContain('John Doe')
})
```

**Testing composables with async:**

```ts
it('handles async state', async () => {
  const { result, isLoading } = useAsyncData()

  expect(isLoading.value).toBe(true)

  await flushPromises()

  expect(isLoading.value).toBe(false)
  expect(result.value).toBeDefined()
})
```

**Vue nextTick:**

* Use `await nextTick()` after reactive state changes
* Ensures DOM updates before assertions

```ts
it('updates DOM after state change', async () => {
  const wrapper = mount(Counter)

  await wrapper.find('button').trigger('click')
  await nextTick()

  expect(wrapper.find('.count').text()).toBe('1')
})
```

---

### 4.8 Test Files and Organization (Vue + TS)

**File location:**

* Test files live next to components: `Component.vue` → `Component.spec.ts`
* Composables: `useAuth.ts` → `useAuth.spec.ts`
* Services: `api.ts` → `api.spec.ts`

**Example structure:**

```
src/
├── components/
│   ├── Header.vue
│   └── Header.spec.ts
├── composables/
│   ├── useAuth.ts
│   └── useAuth.spec.ts
└── services/
    ├── api.ts
    └── api.spec.ts
```

**Naming convention:**

* Use `.spec.ts` suffix (vitest convention)
* Mirror source file name exactly

---

### 4.9 Mocking Stores, Router, and Plugins (Vue + TS)

**Mocking Pinia stores:**

```ts
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from '@/stores/auth'

beforeEach(() => {
  setActivePinia(createPinia())
})

it('uses auth store', () => {
  const authStore = useAuthStore()
  authStore.user = { id: '123', name: 'Test' }

  const wrapper = mount(UserProfile)

  expect(wrapper.text()).toContain('Test')
})
```

**Mocking vue-router:**

```ts
import { mount } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'

const router = createRouter({
  history: createMemoryHistory(),
  routes: [{ path: '/user/:id', component: UserProfile }]
})

it('renders user from route params', async () => {
  await router.push('/user/123')
  await router.isReady()

  const wrapper = mount(UserProfile, {
    global: {
      plugins: [router]
    }
  })

  expect(wrapper.text()).toContain('User 123')
})
```

**Mocking global plugins:**

```ts
const wrapper = mount(Component, {
  global: {
    mocks: {
      $t: (key: string) => key // Mock i18n
    }
  }
})
```

---

### 4.10 Assertions and DOM Queries (Vue + TS)

**Prefer semantic queries:**

```ts
// Good: semantic
wrapper.find('[data-testid="submit-button"]')
wrapper.find('button[type="submit"]')

// Bad: implementation details
wrapper.find('.btn-primary')
wrapper.findAll('div')[2]
```

**Use `get()` vs `find()`:**

* Use `get()` when element must exist (throws if not found)
* Use `find()` when testing absence

```ts
// Element must exist
const button = wrapper.get('button')

// Test element doesn't exist
expect(wrapper.find('.error').exists()).toBe(false)
```

**Descriptive assertions:**

```ts
// Good
expect(wrapper.text()).toContain('Welcome')
expect(wrapper.find('input').element.value).toBe('test@example.com')

// With message
expect(userCount).toBe(5, 'should display all 5 users')
```

---

## 5. Error Scenarios and Edge Cases

**AI agents must test both happy path and error scenarios without distinction.**

### 5.1 Error Path Coverage

Every function that can fail must have tests for:

* Expected errors (validation, not found, etc.)
* Unexpected errors (network, database, etc.)
* Edge cases (empty input, null, boundary values)

**Go example:**

```go
func Test_Service_FindUser_WhenNotFound_ThenReturnsError(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    repo := NewMockRepository(ctrl)
    repo.EXPECT().FindByID("999").Return(nil, ErrNotFound)

    svc := NewService(repo)

    user, err := svc.FindUser("999")

    assert.Nil(t, user)
    assert.ErrorIs(t, err, ErrNotFound)
}
```

**Vue example:**

```ts
it('displays error when login fails', async () => {
  vi.mocked(authApi.login).mockRejectedValue(new Error('Invalid credentials'))

  const wrapper = mount(LoginForm)
  await wrapper.find('form').trigger('submit')
  await flushPromises()

  expect(wrapper.find('.error').text()).toBe('Invalid credentials')
})
```

### 5.2 Edge Cases

Test boundary conditions explicitly:

* Empty arrays/objects
* Null/undefined values
* Zero, negative numbers
* Very large inputs
* Special characters in strings

---

## 6. Cross‑Cutting Rules for AI Agents

### 6.1 When Code is Not Testable

**If code is hard to test due to poor design, AI agents MUST refactor it.**

* Do not work around bad design with complex test setup
* Refactor production code to be testable:
  * Extract dependencies into interfaces
  * Use dependency injection
  * Separate concerns (business logic from I/O)
* **Scope of refactoring:**
  * Generally limit changes to the vertical slice being tested
  * If refactoring requires changes beyond the slice, ask user before proceeding
* Document refactoring rationale in commit messages

**Red flags requiring refactoring:**

* Function has 3+ dependencies to mock
* Cannot test without filesystem, network, or database
* Global state or singletons prevent isolation
* Logic mixed with I/O operations

---

### 6.2 Test Prioritization

**AI agents must test both happy path and error scenarios without distinction.**

* Do not prioritize one over the other
* Every function/method needs:
  * Success cases (expected valid inputs)
  * Failure cases (expected errors, validation)
  * Edge cases (boundaries, nulls, empty values)

**Order of implementation:**

1. Write all test cases together (table-driven when possible)
2. Cover success and failure paths equally
3. Add edge cases based on function signature and business logic

---

### 6.3 Coverage Goals

**Minimum threshold: 80% line coverage**

* AI agents must verify coverage after writing tests
* If coverage is below 80%, add more tests or ask for guidance
* Coverage is a metric, not a goal:
  * Focus on meaningful tests, not just lines covered
  * 100% coverage with poor tests is worse than 80% with good tests

**How to check:**

```bash
# Go
go test -coverprofile=coverage.txt -covermode=atomic ./...
go tool cover -func=coverage.txt | grep total

# Vue (vitest)
npm run test:coverage
```

---

### 6.4 Handling Broken Tests

**If AI agents encounter broken tests, they MUST fix them based on current project state.**

* Do not skip or ignore broken tests
* Understand why test is failing:
  * Code changed but test didn't update?
  * Test was flaky?
  * Requirements changed?
* Fix test to match current implementation and requirements
* If test is no longer relevant, remove it (with justification)

**Process:**

1. Run test to understand failure
2. Read current production code
3. Update test to align with current behavior
4. Verify test passes and is deterministic

---

### 6.5 Git Integration and Test Validation

**This project uses git hooks to enforce test quality.**

* Pre-commit hook runs:
  * `go fmt` and linting
  * Unit tests (`go test -short`)
  * Auto-stages formatted files and go.mod/go.sum

* Pre-push hook runs:
  * Full test suite with race detection
  * Coverage report

**AI agents must:**

* Ensure all tests pass locally before creating commits
* Tests are automatically validated by hooks
* Do not bypass hooks (`--no-verify`) unless explicitly instructed
* If hooks fail, fix issues before proceeding

---

### 6.6 CI/CD Integration

**Continuous Integration validates all tests automatically.**

* CI pipeline runs on every push and PR:
  1. Lint
  2. Build
  3. Security scan (non-blocking)
  4. Test with coverage (`go test -race -coverprofile`)

* Tests must pass in CI for PR approval
* Coverage reports uploaded to Codecov automatically

**AI agents should:**

* Write tests that work in CI environment (no local dependencies)
* Avoid tests that depend on specific machine configuration
* Ensure tests are truly deterministic (pass in any environment)

---

### 6.7 Mandatory Rules

AI agents **must not**:

* Add comments marking AAA sections (`// Arrange`, `// Act`, `// Assert`)
* Introduce sleeps, retries, or waits to "fix" flaky tests
* Silence errors without assertion
* Mock more than one architectural layer per test
* Write tests for code that violates single responsibility
* Skip refactoring code that is hard to test
* Prioritize happy path over error scenarios
* Generate tests with lint errors or warnings
* Leave unused variables, imports, or dead code in tests
* Disable linter rules without explicit justification

AI agents **must**:

* Use blank lines to visually separate AAA sections
* Prefer pure functions and dependency injection
* Refactor production code when testability is poor (within scope)
* Test both success and failure paths equally
* Maintain minimum 80% coverage
* Fix broken tests to match current project state
* Ensure tests pass with git hooks and in CI
* Add descriptive error messages to assertions
* Ensure each unit under test has a single, clear responsibility
* Validate that dependencies are injected, not hard‑coded
* Run linters and fix all errors/warnings before completing tasks
* Ensure tests are properly formatted (go fmt, prettier)
* Remove all unused variables, imports, and debugging code

---

### 6.8 Generated Code and Mocks

**Generated mocks (Go):**

* Use `mockgen` to generate mocks from interfaces
* Name pattern: `mock_<interface>_test.go`
* Include `_test.go` suffix so they're only compiled for tests
* Regenerate mocks when interfaces change

**Example mockgen command:**

```bash
mockgen -source=repository.go -destination=mock_repository_test.go -package=auth
```

**Generated code (protobuf, OpenAPI, etc.):**

* Do NOT write unit tests for generated code
* Test the code that USES generated code
* Focus on integration points and business logic
* Exception: If you modify generated code manually, test those modifications

**Commit strategy:**

* Commit generated mocks to version control (included in this project)
* Ensures consistent test environment across developers
* CI regenerates to detect drift

---

### 6.9 Test Helpers - When to Use Them

**Test helpers are appropriate when:**

* Setting up complex domain objects used across multiple tests
* Creating reusable test fixtures (builders, factories)
* Abstracting repetitive setup that doesn't obscure test intent

**Test helpers should be avoided when:**

* They hide what is actually being tested
* They contain assertions (helpers should set up, not assert)
* They make tests harder to understand in isolation
* They introduce their own logic that needs testing

**Good test helper (Go):**

```go
func createTestUser(t *testing.T, opts ...func(*User)) *User {
    t.Helper() // Marks this as a helper for better error reporting

    user := &User{
        ID:    "test-id",
        Email: "test@example.com",
        Role:  "user",
    }

    for _, opt := range opts {
        opt(user)
    }

    return user
}

// Usage
user := createTestUser(t, func(u *User) { u.Role = "admin" })
```

**Good test helper (Vue):**

```ts
function createMockAuthStore(overrides = {}) {
  return {
    user: null,
    isAuthenticated: false,
    login: vi.fn(),
    logout: vi.fn(),
    ...overrides
  }
}

// Usage
const authStore = createMockAuthStore({ isAuthenticated: true })
```

**Bad test helper (obscures logic):**

```go
// Bad: hides what's being tested
func testServiceBehavior(t *testing.T, input string) {
    svc := setupService()
    result := svc.Process(input)
    assert.Equal(t, "expected", result) // ❌ Hidden assertion
}
```

**Key principle:** Test helpers should make tests more readable, not less. If a helper requires documentation to understand, it's too complex.

---

### 6.10 Lint Compliance

**All generated tests must pass linting without errors or warnings.**

#### Go Linting Rules

AI agents must ensure tests comply with:

* **No unused variables or imports**
  * Every declared variable must be used
  * Remove unused imports immediately
  * Use `_` for intentionally ignored values

```go
// Bad: unused variable
result, err := svc.Process()
require.NoError(t, err)
// result is never used ❌

// Good: use the result
result, err := svc.Process()
require.NoError(t, err)
assert.Equal(t, expected, result)

// Good: explicitly ignore if not needed
_, err := svc.Process()
require.NoError(t, err)
```

* **No shadowed variables**
  * Avoid redeclaring variables in inner scopes
  * Use different names or proper scoping

```go
// Bad: shadowing err
user, err := repo.FindByID("123")
require.NoError(t, err)
if user.Active {
    err := svc.Activate(user) // ❌ shadows outer err
    require.NoError(t, err)
}

// Good: reuse err or use distinct name
user, err := repo.FindByID("123")
require.NoError(t, err)
if user.Active {
    err = svc.Activate(user) // ✓ reuses err
    require.NoError(t, err)
}
```

* **Proper error handling**
  * Never ignore errors without explicit `_`
  * Use `require.NoError` or `assert.Error` as appropriate

```go
// Bad: ignored error
svc.Process() // ❌ error not handled

// Good: error checked
err := svc.Process()
require.NoError(t, err)

// Good: error explicitly ignored with justification
_ = conn.Close() // Cleanup, error not critical in test
```

* **Correct use of `t.Fatal` and `t.Error`**
  * Never use `t.Fatal` in goroutines (causes panic)
  * Use `t.Error` + channels to report goroutine errors

```go
// Bad: Fatal in goroutine
go func() {
    if err != nil {
        t.Fatal(err) // ❌ will panic
    }
}()

// Good: Error + channel
errCh := make(chan error, 1)
go func() {
    if err != nil {
        errCh <- err
    }
    close(errCh)
}()
require.NoError(t, <-errCh)
```

* **Consistent formatting**
  * Run `go fmt` before committing
  * Follow standard Go conventions (handled by formatter)

* **No hardcoded credentials or sensitive data**
  * Use placeholder values in tests
  * Never commit real credentials

#### TypeScript/Vue Linting Rules

AI agents must ensure tests comply with:

* **No unused variables or imports**
  * Remove unused imports immediately
  * Prefix intentionally unused variables with `_`

```ts
// Bad: unused import
import { mount } from '@vue/test-utils'
import { ref } from 'vue' // ❌ never used

// Good: only needed imports
import { mount } from '@vue/test-utils'
```

* **Explicit types, avoid `any`**
  * Always type function parameters and returns
  * Use proper types instead of `any`
  * If `any` is unavoidable, add `// eslint-disable-next-line @typescript-eslint/no-explicit-any`

```ts
// Bad: implicit any
function createMockUser(data) { // ❌ implicit any
  return { ...data }
}

// Good: explicit types
function createMockUser(data: Partial<User>): User {
  return { id: '1', name: 'Test', ...data }
}

// Acceptable with disable comment (rare cases)
// eslint-disable-next-line @typescript-eslint/no-explicit-any
function genericHandler(data: any) {
  return data
}
```

* **Await all promises**
  * Never create promises without awaiting
  * Use `void` prefix if promise result is intentionally not awaited

```ts
// Bad: promise not awaited
it('tests async operation', () => {
  wrapper.vm.loadData() // ❌ promise not awaited
  expect(wrapper.vm.loading).toBe(false)
})

// Good: promise awaited
it('tests async operation', async () => {
  await wrapper.vm.loadData()
  expect(wrapper.vm.loading).toBe(false)
})

// Good: explicitly void if not awaiting (rare)
void wrapper.vm.loadData() // Fire and forget
```

* **No console.log or debugging code**
  * Remove all `console.log`, `debugger`, `console.dir` statements
  * Use proper assertions instead

```ts
// Bad: debugging code left in
it('calculates total', () => {
  const result = calculateTotal(items)
  console.log(result) // ❌ remove before commit
  expect(result).toBe(100)
})

// Good: clean test
it('calculates total', () => {
  const result = calculateTotal(items)
  expect(result).toBe(100)
})
```

* **Consistent formatting**
  * Follow project ESLint/Prettier configuration
  * Tests are auto-formatted on commit

* **No @ts-ignore without justification**
  * Avoid `@ts-ignore` and `@ts-expect-error`
  * If absolutely necessary, add explanation comment

```ts
// Bad: suppressing errors without reason
// @ts-ignore
wrapper.vm.privateMethod() // ❌ no explanation

// Acceptable with justification (rare)
// @ts-expect-error - testing internal API for regression
wrapper.vm.$_privateMethod()
```

#### Lint Validation Process

**Before completing any test generation task:**

1. **Run linters locally**

```bash
# Go
go vet ./...
golangci-lint run ./...

# TypeScript/Vue
npm run lint
# or
npm run lint:fix  # auto-fix where possible
```

2. **Fix all errors and warnings**
   * Do not commit code with lint errors
   * Fix warnings proactively
   * Never disable linting rules without explicit approval

3. **Verify in CI**
   * Tests must pass lint stage in CI pipeline
   * CI fails on any lint errors

#### Common Lint Pitfalls to Avoid

**Go:**
* Returning `nil` error as naked `nil` - use explicit type: `return nil, ErrNotFound` not `return nil, nil`
* Unchecked type assertions - use comma-ok pattern: `val, ok := x.(Type)`
* Comparing errors with `==` instead of `errors.Is()`
* Not checking `Close()` errors - at minimum: `defer func() { _ = f.Close() }()`

**TypeScript/Vue:**
* Using `{}` instead of `Record<string, unknown>` for object types
* Declaring variables with `var` instead of `const`/`let`
* Not handling promise rejections in async tests
* Using `toBe()` for object comparison instead of `toEqual()`

#### Linter Configuration

This project uses:

* **Go:** `golangci-lint` with project-specific configuration (`.golangci.yml` if present)
* **TypeScript/Vue:** ESLint + Prettier with rules defined in `eslintrc`, `prettier.config.js`

**AI agents must:**
* Follow the existing linter configuration
* Never modify linter rules without explicit user approval
* Fix code to comply with linters, not the other way around
* Ask for guidance if a lint rule seems incorrect

---

## 7. Definition of a Good Unit Test

A unit test is considered **valid** if:

* It runs in isolation (no shared state, no order dependency)
* It is deterministic across multiple runs (`-count=50` for Go, `--repeat=20` for vitest)
* It completes in milliseconds (< 10ms for Go, < 5ms for Vue)
* It clearly communicates intent through naming and structure
* It follows AAA pattern **without comments** (blank lines only)
* It tests a unit with a single, clear responsibility
* It mocks only at architectural boundaries
* It fails for exactly one reason
* It tests both happy path and error scenarios
* It includes descriptive assertion messages
* It contributes to maintaining 80% coverage
* It passes in CI environment (no local dependencies)
* It passes race detection (`go test -race` for Go)
* **It passes all linters without errors or warnings** (see section 6.10)

If any of these conditions are violated, the test must be rewritten.

**Red flags that indicate poor separation of concerns**:

* Test requires mocking 3+ dependencies
* Setup section is longer than the assertion section
* Multiple unrelated behaviors tested in one function
* Cannot test without complex setup or global state
* Test name contains "and", "or", "also"
* Test has complex conditional logic
* Test modifies global state

**When you see these red flags:**

1. Refactor production code (within vertical slice scope)
2. Extract dependencies into interfaces
3. Use dependency injection
4. Separate business logic from I/O
5. Ask user if refactoring requires changes beyond current slice

---

**End of document.**
