# Integration Testing Guide (Go & Vue + TypeScript)

**Audience:** AI agents responsible for generating, maintaining, or refactoring integration tests.

**Goal:** Ensure integration tests are **fast, deterministic, reliable, and maintainable**, validating component interactions and system behavior with real dependencies.

---

## 1. Core Principles (Non‑Negotiable)

1. **Integration tests must be deterministic**

   * No dependency on external services beyond control (production APIs, third-party services)
   * Use controlled test infrastructure (Docker containers, in-memory databases, test servers)
   * Any external dependency must be containerized or mocked at the boundary

2. **Tests must be reasonably fast**

   * Target execution time:

     * Individual test: **< 500ms**
     * Entire integration suite: **< 30s** locally
   * Slow tests (> 1s) require justification and optimization
   * Use parallel execution where possible

3. **Complete isolation between tests**

   * Each test must have independent data and state
   * Database cleanup after each test (transactions, truncate, or recreate)
   * No shared mutable state between tests
   * Tests must pass in any order

4. **One integration scenario per test**

   * Each test validates one end-to-end flow or integration point
   * Complex scenarios should be broken down into focused tests

5. **Tests are code**

   * Apply the same standards as production code: readability, naming, structure, and refactoring

6. **Follow AAA pattern implicitly**

   * Tests must follow **Arrange → Act → Assert** structure
   * **NEVER add comments** marking each section (`// Arrange`, `// Act`, `// Assert`)
   * The structure should be self‑evident through blank lines and code organization

7. **Test real integrations**

   * Integration tests validate interactions between components
   * Use real databases, message queues, HTTP servers (in test mode)
   * Mock only external services beyond your control (third-party APIs, payment gateways)

---

## 2. Coverage Requirements

**Integration tests complement unit tests, not replace them.**

* Focus on critical paths and integration points
* Cover all API endpoints and their error scenarios
* Test database transactions and rollback behavior
* Validate authentication and authorization flows
* Test message queue producers/consumers
* Coverage target: **Critical paths must have integration tests**

Verification:

```bash
# Go - run integration tests
go test -tags=integration ./...

# Go - with coverage
go test -tags=integration -coverprofile=integration-coverage.txt ./...

# TypeScript/Vue - run integration tests
npm run test:integration
```

---

## 3. Go – Integration Testing Guidelines

### 3.1 Mandatory Libraries

Use **only** the following libraries unless explicitly instructed otherwise:

* **Testing framework**

  * `testing` (standard library)
  * Build tags: `//go:build integration`

* **Assertions**

  * `github.com/stretchr/testify/assert`
  * `github.com/stretchr/testify/require`
  * `github.com/stretchr/testify/suite` (for complex setup/teardown)

* **Test containers** (for real dependencies)

  * `github.com/testcontainers/testcontainers-go`
  * Supports PostgreSQL, Redis, MySQL, MongoDB, Kafka, etc.

* **HTTP testing**

  * `net/http/httptest` (standard library)
  * For testing HTTP handlers and clients

* **Database helpers**

  * Project-specific DB helpers for migrations and seeding
  * Transaction-based isolation when possible

---

### 3.2 Test Structure (Go)

**Build tags for separation:**

All integration test files must start with:

```go
//go:build integration

package mypackage_test
```

This allows running integration tests separately:

```bash
# Run only integration tests
go test -tags=integration ./...

# Run only unit tests (default, no tags)
go test ./...
```

**Test structure:**

Follow **Arrange → Act → Assert** pattern **without comments**.

```go
//go:build integration

package api_test

import (
    "context"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/postgres"
)

func TestIntegration_UserAPI_CreateUser_StoresInDatabase(t *testing.T) {
    ctx := context.Background()
    db := setupTestDatabase(t, ctx)
    defer db.Close()
    api := NewAPI(db)

    user := &User{Email: "test@example.com", Name: "Test User"}
    err := api.CreateUser(ctx, user)

    require.NoError(t, err)
    assert.NotEmpty(t, user.ID)

    stored, err := db.FindUserByEmail(ctx, "test@example.com")
    require.NoError(t, err)
    assert.Equal(t, user.Email, stored.Email)
    assert.Equal(t, user.Name, stored.Name)
}
```

**Naming pattern:**

```
TestIntegration_<Component>_<Scenario>_<ExpectedOutcome>
```

Prefix with `TestIntegration_` to distinguish from unit tests.

---

### 3.3 Test Containers for Dependencies

**Use testcontainers-go for real dependencies:**

```go
func setupTestDatabase(t *testing.T, ctx context.Context) *sql.DB {
    t.Helper()

    postgresContainer, err := postgres.Run(ctx,
        "postgres:16-alpine",
        postgres.WithDatabase("testdb"),
        postgres.WithUsername("testuser"),
        postgres.WithPassword("testpass"),
        testcontainers.WithWaitStrategy(
            wait.ForLog("database system is ready to accept connections").
                WithOccurrence(2).
                WithStartupTimeout(30*time.Second)),
    )
    require.NoError(t, err)

    t.Cleanup(func() {
        if err := postgresContainer.Terminate(ctx); err != nil {
            t.Logf("failed to terminate container: %s", err)
        }
    })

    connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
    require.NoError(t, err)

    db, err := sql.Open("postgres", connStr)
    require.NoError(t, err)

    runMigrations(t, db)

    return db
}
```

**Key principles:**

* Create containers in test setup
* Use `t.Cleanup()` for automatic teardown
* Set reasonable timeouts (30s for container startup)
* Reuse containers across tests in same package when possible (via `TestMain`)
* Always wait for service to be ready before returning

---

### 3.4 Database Isolation Strategies

**Option 1: Transactions (Fastest)**

```go
func TestIntegration_WithTransaction(t *testing.T) {
    ctx := context.Background()
    db := getSharedTestDB(t) // Reused connection

    tx, err := db.BeginTx(ctx, nil)
    require.NoError(t, err)
    defer tx.Rollback() // Always rollback, even if test passes

    repo := NewRepository(tx)

    err = repo.CreateUser(ctx, &User{Email: "test@example.com"})

    require.NoError(t, err)
    // Transaction is rolled back, database stays clean
}
```

**Benefits:**
* Extremely fast (no physical cleanup)
* Complete isolation
* Automatically reverted

**Limitations:**
* Cannot test transaction behavior itself
* Doesn't work with multiple connections
* Some ORMs may not support testing with transactions

**Option 2: Truncate Tables (Moderate)**

```go
func TestIntegration_WithTruncate(t *testing.T) {
    ctx := context.Background()
    db := getSharedTestDB(t)

    t.Cleanup(func() {
        truncateAllTables(t, db)
    })

    repo := NewRepository(db)

    err := repo.CreateUser(ctx, &User{Email: "test@example.com"})

    require.NoError(t, err)
}

func truncateAllTables(t *testing.T, db *sql.DB) {
    t.Helper()
    tables := []string{"users", "orders", "products"}

    for _, table := range tables {
        _, err := db.Exec("TRUNCATE TABLE " + table + " CASCADE")
        if err != nil {
            t.Logf("failed to truncate %s: %v", table, err)
        }
    }
}
```

**Option 3: Recreate Database (Slowest, Most Thorough)**

```go
func TestIntegration_WithRecreate(t *testing.T) {
    ctx := context.Background()
    db := setupTestDatabase(t, ctx) // Fresh database per test
    defer db.Close()

    // Test with completely fresh database
}
```

**Recommendation:**
* Use **transactions** for most tests (fastest)
* Use **truncate** when testing transaction behavior
* Use **recreate** only for tests that modify schema or require complete isolation

---

### 3.5 Test Suites for Complex Setup

Use testify suites when setup/teardown is complex:

```go
//go:build integration

package api_test

import (
    "context"
    "database/sql"
    "testing"

    "github.com/stretchr/testify/suite"
    "github.com/testcontainers/testcontainers-go/postgres"
)

type APIIntegrationSuite struct {
    suite.Suite
    ctx       context.Context
    db        *sql.DB
    container *postgres.PostgresContainer
    api       *API
}

func (s *APIIntegrationSuite) SetupSuite() {
    s.ctx = context.Background()

    container, err := postgres.Run(s.ctx, "postgres:16-alpine")
    s.Require().NoError(err)
    s.container = container

    connStr, err := container.ConnectionString(s.ctx)
    s.Require().NoError(err)

    db, err := sql.Open("postgres", connStr)
    s.Require().NoError(err)
    s.db = db

    runMigrations(s.T(), db)

    s.api = NewAPI(db)
}

func (s *APIIntegrationSuite) TearDownSuite() {
    s.db.Close()
    s.container.Terminate(s.ctx)
}

func (s *APIIntegrationSuite) SetupTest() {
    truncateAllTables(s.T(), s.db)
}

func (s *APIIntegrationSuite) TestCreateUser_StoresInDatabase() {
    user := &User{Email: "test@example.com"}

    err := s.api.CreateUser(s.ctx, user)

    s.NoError(err)
    s.NotEmpty(user.ID)
}

func (s *APIIntegrationSuite) TestCreateUser_DuplicateEmail_ReturnsError() {
    user1 := &User{Email: "test@example.com"}
    s.api.CreateUser(s.ctx, user1)

    user2 := &User{Email: "test@example.com"}
    err := s.api.CreateUser(s.ctx, user2)

    s.Error(err)
    s.ErrorIs(err, ErrDuplicateEmail)
}

func TestAPIIntegrationSuite(t *testing.T) {
    suite.Run(t, new(APIIntegrationSuite))
}
```

**When to use suites:**
* Container setup is expensive (> 5s)
* Multiple tests share same infrastructure
* Complex setup/teardown logic
* Need lifecycle hooks (SetupSuite, SetupTest, TearDownTest, TearDownSuite)

**When NOT to use suites:**
* Simple tests with minimal setup
* Tests require complete isolation (different configurations)
* Setup is fast (< 100ms)

---

### 3.6 HTTP Integration Testing

**Testing HTTP handlers:**

```go
func TestIntegration_UserHandler_CreateUser_Returns201(t *testing.T) {
    ctx := context.Background()
    db := setupTestDatabase(t, ctx)
    handler := NewUserHandler(db)

    reqBody := `{"email":"test@example.com","name":"Test User"}`
    req := httptest.NewRequest("POST", "/users", strings.NewReader(reqBody))
    req.Header.Set("Content-Type", "application/json")
    rec := httptest.NewRecorder()

    handler.ServeHTTP(rec, req)

    assert.Equal(t, http.StatusCreated, rec.Code)

    var response UserResponse
    err := json.Unmarshal(rec.Body.Bytes(), &response)
    require.NoError(t, err)
    assert.Equal(t, "test@example.com", response.Email)
}
```

**Testing full HTTP server:**

```go
func TestIntegration_Server_EndToEnd(t *testing.T) {
    ctx := context.Background()
    db := setupTestDatabase(t, ctx)
    server := NewServer(db)

    ts := httptest.NewServer(server.Handler())
    defer ts.Close()

    client := &http.Client{Timeout: 5 * time.Second}

    reqBody := `{"email":"test@example.com"}`
    resp, err := client.Post(ts.URL+"/users", "application/json", strings.NewReader(reqBody))

    require.NoError(t, err)
    defer resp.Body.Close()
    assert.Equal(t, http.StatusCreated, resp.StatusCode)
}
```

---

### 3.7 Performance Optimization (Go)

**1. Reuse containers across tests**

```go
var (
    testDB     *sql.DB
    testDBOnce sync.Once
)

func getSharedTestDB(t *testing.T) *sql.DB {
    t.Helper()

    testDBOnce.Do(func() {
        ctx := context.Background()
        container, err := postgres.Run(ctx, "postgres:16-alpine")
        if err != nil {
            t.Fatalf("failed to start container: %v", err)
        }

        connStr, _ := container.ConnectionString(ctx)
        db, _ := sql.Open("postgres", connStr)
        runMigrations(t, db)

        testDB = db
    })

    return testDB
}
```

**2. Parallel test execution**

```go
func TestIntegration_ParallelSafe(t *testing.T) {
    t.Parallel() // Only if tests are truly isolated

    ctx := context.Background()
    db := getSharedTestDB(t)

    tx, err := db.BeginTx(ctx, nil)
    require.NoError(t, err)
    defer tx.Rollback()

    // Test with transaction isolation
}
```

**3. Use connection pooling**

```go
func setupTestDatabase(t *testing.T, ctx context.Context) *sql.DB {
    // ... container setup ...

    db, err := sql.Open("postgres", connStr)
    require.NoError(t, err)

    db.SetMaxOpenConns(10)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(5 * time.Minute)

    return db
}
```

**4. Minimize container startup**

```go
func TestMain(m *testing.M) {
    // Start shared containers once for entire package
    ctx := context.Background()
    container, _ := postgres.Run(ctx, "postgres:16-alpine")
    defer container.Terminate(ctx)

    connStr, _ := container.ConnectionString(ctx)
    db, _ := sql.Open("postgres", connStr)
    defer db.Close()

    runMigrations(nil, db)
    testDB = db

    code := m.Run()
    os.Exit(code)
}
```

---

### 3.8 Testing Message Queues and Async Operations

**Redis/queue integration:**

```go
func TestIntegration_EventPublisher_PublishesEvent(t *testing.T) {
    ctx := context.Background()

    redisContainer, err := redis.Run(ctx, "redis:7-alpine")
    require.NoError(t, err)
    defer redisContainer.Terminate(ctx)

    connStr, err := redisContainer.ConnectionString(ctx)
    require.NoError(t, err)

    client := redis.NewClient(&redis.Options{Addr: connStr})
    publisher := NewEventPublisher(client)

    event := &OrderCreatedEvent{OrderID: "123"}
    err = publisher.Publish(ctx, "orders", event)

    require.NoError(t, err)

    result, err := client.LPop(ctx, "orders").Result()
    require.NoError(t, err)
    assert.Contains(t, result, "123")
}
```

**Testing async operations with timeouts:**

```go
func TestIntegration_AsyncProcessor_ProcessesEvents(t *testing.T) {
    ctx := context.Background()
    processor := NewEventProcessor()

    done := make(chan bool, 1)
    processor.OnEvent(func(e Event) {
        done <- true
    })

    processor.Start(ctx)
    defer processor.Stop()

    processor.Publish(Event{Type: "test"})

    select {
    case <-done:
        // Success
    case <-time.After(2 * time.Second):
        t.Fatal("timeout waiting for event processing")
    }
}
```

---

### 3.9 Context and Timeout Management

**Always use timeouts in integration tests:**

```go
func TestIntegration_WithTimeout(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    db := getSharedTestDB(t)
    repo := NewRepository(db)

    err := repo.SlowOperation(ctx)

    require.NoError(t, err)
}
```

**Test context cancellation:**

```go
func TestIntegration_Service_ContextCancellation_StopsOperation(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    db := getSharedTestDB(t)
    svc := NewService(db)

    errCh := make(chan error, 1)
    go func() {
        errCh <- svc.LongRunningOperation(ctx)
    }()

    time.Sleep(100 * time.Millisecond)
    cancel()

    err := <-errCh
    assert.ErrorIs(t, err, context.Canceled)
}
```

---

### 3.10 Environment Variables and Configuration

**Test-specific configuration:**

```go
func TestIntegration_WithConfig(t *testing.T) {
    t.Setenv("APP_ENV", "test")
    t.Setenv("LOG_LEVEL", "error")

    config := LoadConfig()

    assert.Equal(t, "test", config.Environment)
    assert.Equal(t, "error", config.LogLevel)
}
```

**Benefits of `t.Setenv()`:**
* Automatically restored after test
* Safe for parallel tests
* Prevents env pollution between tests

---

## 4. Vue + TypeScript – Integration Testing Guidelines

### 4.1 Mandatory Libraries

Use the following stack:

* **Test runner**

  * `vitest`

* **Component integration testing**

  * `@vue/test-utils`
  * Mount components with real dependencies (stores, router)

* **API mocking**

  * `msw` (Mock Service Worker) - intercepts HTTP requests
  * Realistic network-level mocking

* **DOM environment**

  * `happy-dom` or `jsdom`

---

### 4.2 Test Structure (Vue)

**File naming:**

```
<component>.integration.spec.ts
```

Separate from unit tests (`.spec.ts`).

**Example integration test:**

```ts
import { describe, it, expect, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { createRouter, createMemoryHistory } from 'vue-router'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'

import UserProfile from './UserProfile.vue'
import { useAuthStore } from '@/stores/auth'

const server = setupServer()

beforeAll(() => server.listen())
afterEach(() => server.resetHandlers())
afterAll(() => server.close())

describe('UserProfile Integration', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('loads and displays user data from API', async () => {
    server.use(
      http.get('/api/users/123', () => {
        return HttpResponse.json({
          id: '123',
          name: 'John Doe',
          email: 'john@example.com'
        })
      })
    )

    const wrapper = mount(UserProfile, {
      props: { userId: '123' }
    })

    await flushPromises()

    expect(wrapper.text()).toContain('John Doe')
    expect(wrapper.text()).toContain('john@example.com')
  })
})
```

---

### 4.3 API Mocking with MSW

**Setup MSW server:**

```ts
// tests/mocks/server.ts
import { setupServer } from 'msw/node'
import { handlers } from './handlers'

export const server = setupServer(...handlers)
```

```ts
// tests/mocks/handlers.ts
import { http, HttpResponse } from 'msw'

export const handlers = [
  http.get('/api/users/:id', ({ params }) => {
    return HttpResponse.json({
      id: params.id,
      name: 'Test User'
    })
  }),

  http.post('/api/users', async ({ request }) => {
    const body = await request.json()
    return HttpResponse.json(
      { id: '123', ...body },
      { status: 201 }
    )
  })
]
```

**Use in tests:**

```ts
import { server } from '@/tests/mocks/server'
import { http, HttpResponse } from 'msw'

beforeAll(() => server.listen({ onUnhandledRequest: 'error' }))
afterEach(() => server.resetHandlers())
afterAll(() => server.close())

it('handles API errors gracefully', async () => {
  server.use(
    http.get('/api/users/123', () => {
      return HttpResponse.json(
        { error: 'Not found' },
        { status: 404 }
      )
    })
  )

  const wrapper = mount(UserProfile, {
    props: { userId: '123' }
  })

  await flushPromises()

  expect(wrapper.find('.error').text()).toContain('Not found')
})
```

---

### 4.4 Store Integration Testing

**Test components with real stores:**

```ts
import { createPinia, setActivePinia } from 'pinia'
import { useAuthStore } from '@/stores/auth'

describe('Dashboard Integration', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('displays user info from auth store', async () => {
    const authStore = useAuthStore()
    authStore.user = {
      id: '123',
      name: 'John Doe',
      role: 'admin'
    }

    const wrapper = mount(Dashboard)

    expect(wrapper.text()).toContain('John Doe')
    expect(wrapper.find('[data-testid="admin-panel"]').exists()).toBe(true)
  })

  it('redirects to login when not authenticated', async () => {
    const router = createRouter({
      history: createMemoryHistory(),
      routes: [
        { path: '/', component: Dashboard },
        { path: '/login', component: Login }
      ]
    })

    const authStore = useAuthStore()
    authStore.user = null

    const wrapper = mount(Dashboard, {
      global: { plugins: [router] }
    })

    await router.isReady()

    expect(router.currentRoute.value.path).toBe('/login')
  })
})
```

---

### 4.5 Router Integration Testing

**Test navigation flows:**

```ts
import { createRouter, createMemoryHistory } from 'vue-router'
import { routes } from '@/router'

describe('Navigation Integration', () => {
  it('navigates through checkout flow', async () => {
    const router = createRouter({
      history: createMemoryHistory(),
      routes
    })

    await router.push('/cart')
    await router.isReady()

    const wrapper = mount(App, {
      global: { plugins: [router] }
    })

    await wrapper.find('[data-testid="checkout-btn"]').trigger('click')
    await router.isReady()

    expect(router.currentRoute.value.path).toBe('/checkout')
    expect(wrapper.find('h1').text()).toBe('Checkout')
  })
})
```

---

### 4.6 Performance Optimization (Vue)

**1. Reuse MSW server setup**

```ts
// vitest.setup.ts
import { beforeAll, afterEach, afterAll } from 'vitest'
import { server } from './mocks/server'

beforeAll(() => server.listen({ onUnhandledRequest: 'error' }))
afterEach(() => server.resetHandlers())
afterAll(() => server.close())
```

**2. Shallow mount when possible**

```ts
import { shallowMount } from '@vue/test-utils'

it('renders with minimal overhead', () => {
  const wrapper = shallowMount(ComplexComponent, {
    props: { data: testData }
  })

  expect(wrapper.exists()).toBe(true)
})
```

**3. Parallel test execution**

```ts
// vitest.config.ts
export default defineConfig({
  test: {
    pool: 'threads',
    poolOptions: {
      threads: {
        maxThreads: 4,
        minThreads: 2
      }
    }
  }
})
```

**4. Use fake timers**

```ts
import { vi } from 'vitest'

it('handles delayed operations', async () => {
  vi.useFakeTimers()

  const wrapper = mount(AutoSaveComponent)

  wrapper.find('input').setValue('test')

  vi.advanceTimersByTime(1000)
  await flushPromises()

  expect(mockSave).toHaveBeenCalled()

  vi.useRealTimers()
})
```

---

### 4.7 Testing Composables Integration

**Test composables with real dependencies:**

```ts
import { useUserData } from '@/composables/useUserData'
import { server } from '@/tests/mocks/server'

describe('useUserData Integration', () => {
  it('fetches and caches user data', async () => {
    let callCount = 0
    server.use(
      http.get('/api/users/123', () => {
        callCount++
        return HttpResponse.json({ id: '123', name: 'John' })
      })
    )

    const { data, fetch } = useUserData('123')

    await fetch()
    expect(data.value?.name).toBe('John')
    expect(callCount).toBe(1)

    await fetch()
    expect(callCount).toBe(1) // Cached, no second call
  })
})
```

---

## 5. Error Scenarios and Edge Cases

**Integration tests must validate error handling end-to-end.**

### 5.1 Database Errors (Go)

```go
func TestIntegration_UserService_DatabaseDown_ReturnsError(t *testing.T) {
    ctx := context.Background()
    db := setupTestDatabase(t, ctx)

    db.Close() // Simulate database failure

    svc := NewUserService(db)

    _, err := svc.GetUser(ctx, "123")

    assert.Error(t, err)
    assert.Contains(t, err.Error(), "connection")
}
```

### 5.2 Network Errors (Vue)

```ts
it('retries failed API requests', async () => {
  let attempts = 0
  server.use(
    http.get('/api/data', () => {
      attempts++
      if (attempts < 3) {
        return HttpResponse.error()
      }
      return HttpResponse.json({ data: 'success' })
    })
  )

  const wrapper = mount(DataComponent)

  await flushPromises()

  expect(attempts).toBe(3)
  expect(wrapper.text()).toContain('success')
})
```

### 5.3 Timeout Scenarios

```go
func TestIntegration_SlowQuery_Timeout_ReturnsError(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
    defer cancel()

    db := getSharedTestDB(t)
    _, err := db.ExecContext(ctx, "SELECT pg_sleep(1)")

    assert.Error(t, err)
    assert.ErrorIs(t, err, context.DeadlineExceeded)
}
```

---

## 6. Cross‑Cutting Rules for AI Agents

### 6.1 When to Write Integration Tests

**Integration tests are required for:**

* API endpoints (all HTTP handlers)
* Database operations (CRUD with real DB)
* Authentication/authorization flows
* External service integrations (with mocked services)
* Message queue producers/consumers
* Multi-component workflows (e.g., checkout process)
* Complex state management interactions

**Do NOT write integration tests for:**

* Pure functions (use unit tests)
* Simple utility functions
* UI-only components without API calls
* Code already covered by E2E tests

---

### 6.2 Test Isolation Checklist

Before marking an integration test complete, verify:

- [ ] Test passes when run alone
- [ ] Test passes when run with all other tests
- [ ] Test passes in random order (`go test -shuffle=on`)
- [ ] Test cleans up all resources (database, files, connections)
- [ ] Test doesn't depend on external services (use containers)
- [ ] Test has timeout protection (< 10s for Go, < 5s for Vue)
- [ ] Test is parallel-safe (if using `t.Parallel()`)

---

### 6.3 Performance Checklist

Integration tests must be optimized:

- [ ] Individual test: < 500ms (Go), < 300ms (Vue)
- [ ] Reuses containers/infrastructure when possible
- [ ] Uses transactions for database isolation (when applicable)
- [ ] Runs in parallel with other tests (when isolated)
- [ ] Minimizes network round-trips
- [ ] Uses connection pooling
- [ ] Avoids unnecessary sleeps (use channels/signals instead)

---

### 6.4 Coverage Verification

**Run integration tests separately:**

```bash
# Go - integration tests only
go test -tags=integration -v ./...
go test -tags=integration -coverprofile=integration-coverage.txt ./...

# Vue - integration tests
npm run test:integration
npm run test:integration -- --coverage
```

**Verify both unit and integration coverage:**

```bash
# Go - all tests
go test -tags=integration -coverprofile=all-coverage.txt ./...
go tool cover -func=all-coverage.txt | grep total

# Should see > 90% total coverage
```

---

### 6.5 Mandatory Rules

AI agents **must not**:

* Write integration tests that depend on production services
* Create tests that require manual setup (databases, config files)
* Use sleeps for synchronization (use proper waits)
* Mock internal components (use real implementations)
* Leave containers running after tests
* Share mutable state between integration tests
* Skip cleanup in case of test failure
* Write tests that take > 1s without justification

AI agents **must**:

* Use testcontainers for databases and services
* Clean up resources using `t.Cleanup()` or `defer`
* Set reasonable timeouts (< 10s)
* Test both success and failure paths
* Use transactions for database isolation when possible
* Mark tests with appropriate build tags (`//go:build integration`)
* Ensure tests are parallel-safe or explicitly serialize
* Use MSW for HTTP mocking in Vue tests
* Verify tests pass with `-shuffle=on` (Go) or `--sequence.shuffle` (vitest)

---

### 6.6 Lint Compliance

**Integration tests must pass all linters without errors.**

Same rules as unit tests apply:

* No unused variables or imports
* Proper error handling
* No console.log or debugging code
* Type safety (TypeScript)
* Proper cleanup of resources

Run linters before completing:

```bash
# Go
go vet ./...
golangci-lint run ./...

# TypeScript/Vue
npm run lint
```

---

## 7. Definition of a Good Integration Test

An integration test is considered **valid** if:

* It tests real component interactions (no mocked internals)
* It runs in isolation (can run alone or with others)
* It is deterministic across multiple runs
* It completes in < 500ms (Go) or < 300ms (Vue)
* It uses real dependencies via containers (databases, queues)
* It cleans up all resources automatically
* It follows AAA pattern **without comments**
* It has timeout protection (context.WithTimeout or test timeout)
* It tests both success and error scenarios
* It validates data persistence and side effects
* It passes all linters without errors or warnings
* It runs in CI without external dependencies
* It can run in parallel with other tests (or explicitly serialized)

**Red flags indicating poor integration test design:**

* Test requires manual database setup
* Test depends on specific execution order
* Test leaves data/files after completion
* Test takes > 1s without clear reason
* Test uses production API endpoints
* Test has hardcoded waits/sleeps
* Test shares global state
* Test mocks internal application components
* Test requires specific machine configuration

**When you see these red flags:**

1. Use testcontainers for dependencies
2. Add proper cleanup with `t.Cleanup()` or `defer`
3. Use transactions for database isolation
4. Replace sleeps with proper synchronization
5. Ensure complete isolation between tests
6. Ask user if fundamental refactoring is needed

---

## 8. CI/CD Integration

**GitHub Actions example:**

```yaml
name: Integration Tests

on: [push, pull_request]

jobs:
  integration-tests-go:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'

      - name: Run integration tests
        run: |
          go test -tags=integration -v -race -timeout=5m ./...
          go test -tags=integration -coverprofile=integration.txt ./...

      - uses: codecov/codecov-action@v4
        with:
          files: ./integration.txt

  integration-tests-vue:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-node@v4
        with:
          node-version: '20'

      - run: npm ci

      - name: Run integration tests
        run: npm run test:integration -- --coverage

      - uses: codecov/codecov-action@v4
        with:
          files: ./coverage/integration/coverage-final.json
```

**Key CI requirements:**

* Integration tests run in separate job (parallel with unit tests)
* Timeout limits (5m for Go, 3m for Vue)
* Coverage reporting to track integration coverage separately
* Docker-in-Docker support for testcontainers

---

## 9. Common Patterns and Examples

### 9.1 Testing Authentication Flow (Go)

```go
func TestIntegration_AuthFlow_LoginAndAccess(t *testing.T) {
    ctx := context.Background()
    db := setupTestDatabase(t, ctx)
    api := NewAPI(db)

    user := &User{Email: "test@example.com", Password: "secret123"}
    err := api.Register(ctx, user)
    require.NoError(t, err)

    token, err := api.Login(ctx, "test@example.com", "secret123")

    require.NoError(t, err)
    assert.NotEmpty(t, token)

    profile, err := api.GetProfile(ctx, token)
    require.NoError(t, err)
    assert.Equal(t, "test@example.com", profile.Email)
}
```

### 9.2 Testing Form Submission (Vue)

```ts
describe('ContactForm Integration', () => {
  it('submits form and displays success message', async () => {
    server.use(
      http.post('/api/contact', async ({ request }) => {
        const body = await request.json()
        return HttpResponse.json(
          { message: 'Received', id: '123' },
          { status: 201 }
        )
      })
    )

    const wrapper = mount(ContactForm)

    await wrapper.find('input[name="email"]').setValue('test@example.com')
    await wrapper.find('textarea[name="message"]').setValue('Hello')
    await wrapper.find('form').trigger('submit')
    await flushPromises()

    expect(wrapper.find('.success').text()).toContain('Received')
    expect(wrapper.find('form').exists()).toBe(false)
  })
})
```

### 9.3 Testing Pagination (Go)

```go
func TestIntegration_ListUsers_Pagination_ReturnsCorrectPages(t *testing.T) {
    ctx := context.Background()
    db := setupTestDatabase(t, ctx)
    repo := NewUserRepository(db)

    for i := 0; i < 25; i++ {
        repo.Create(ctx, &User{Email: fmt.Sprintf("user%d@test.com", i)})
    }

    page1, err := repo.List(ctx, Pagination{Page: 1, Limit: 10})
    require.NoError(t, err)
    assert.Len(t, page1.Items, 10)
    assert.Equal(t, 25, page1.Total)

    page2, err := repo.List(ctx, Pagination{Page: 2, Limit: 10})
    require.NoError(t, err)
    assert.Len(t, page2.Items, 10)

    page3, err := repo.List(ctx, Pagination{Page: 3, Limit: 10})
    require.NoError(t, err)
    assert.Len(t, page3.Items, 5)
}
```

---

**End of document.**
