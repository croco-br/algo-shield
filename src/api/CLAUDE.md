# Go API Conventions

## Code Style

- Standard Go conventions (gofmt, golangci-lint), `internal/` for private code
- Error wrapping: `fmt.Errorf("context: %w", err)` — always add context
- **All handlers MUST validate requests** using `go-playground/validator`
- `context.Context` required for all DB/Redis operations
- Interfaces defined in the consumer package, not the provider
- Structured logging with levels (debug, info, warn, error)

## Mocks & Testing

- Mock files (`mock_*_test.go`) are **manually maintained** — not auto-generated
- Mocks use `Mock*` prefix (e.g., `MockTransactionService`) in `_test.go` files
- When interfaces change, update mocks by hand
- Run with `-race` flag always: `go test -race ./...`
- Flaky detection: `go test -count=50 ./...`

## Database Patterns

- `pgxpool` for connection pooling, `pgx.Tx` for atomic operations
- **ALWAYS** parameterized queries (`$1, $2`) — never string concatenation
- Use query timeouts via context
- Connection pool limits must be respected

## Auth Architecture

- JWT access tokens (24h) + refresh tokens (168h)
- Token revocation via Redis sets (individual + user-level blacklists)
- CSRF tokens: Redis-backed, 24h TTL, validated on POST/PUT/PATCH/DELETE
- RBAC roles: `admin`, `rule_editor`, `viewer`
- Password hashing: bcrypt with cost 12+

## Dependency Flow

```
Handler → Service (interfaces) → Repository → DB
```

Route wiring: `src/api/internal/routes/routes.go`

## Key References

- Full testing guide: `docs/agents/unit-test.md`, `docs/agents/integration-test.md`
- Security guidelines: `openspec/project.md` § "Code Security Guidelines"
- Environment config: `.env.example`
