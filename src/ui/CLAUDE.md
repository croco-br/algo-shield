# Vue.js / TypeScript Conventions

## Code Style

- `<script setup lang="ts">` with strict mode for all components
- PascalCase for component files, base components: `BaseButton`, `BaseInput`, `BaseTable`, `BaseBadge`
- Pinia for state management, Tailwind utility classes for styling
- Organize by feature, not by type
- Remove unused code and components regularly

## API Client & CSRF

- **Always use `api.post/put/patch/delete`** from `src/ui/src/lib/api.ts` — never raw `fetch()` or `axios`
- API client auto-adds `Authorization` and `X-CSRF-Token` headers
- `api.ts` uses a holder pattern (`csrfTokenGetterHolder`) to avoid circular deps with auth store
- CSRF token stored in Pinia memory, NOT localStorage
- On 403 "CSRF token missing": user must logout and login again

## i18n (Internationalization)

**No hardcoded user-facing strings.** Use:
- Templates: `{{ $t('key') }}`
- Scripts: `const { t } = useI18n(); t('key')`

Key convention: `common.*`, `auth.*`, `header.*`, `sidebar.*`, `views.<name>.*`, `components.<name>.*`, `errors.*`

Add to BOTH `src/ui/src/locales/pt-BR.json` and `src/ui/src/locales/en-US.json`. Full guide: `src/ui/src/locales/README.md`

## Auth & Tokens

- Access token: Pinia store (memory only) — NEVER localStorage
- Refresh token: localStorage (used only for getting new access tokens)
- CSRF token: Pinia store (memory only)
- Clear all tokens on logout
- Router guards protect routes: `requiresAuth`, `requiresAdmin`

## Test Mocks

- Every `vi.mock('@/lib/api', ...)` block across test files must be kept in sync when new exports are added to `api.ts`
- Also update `vitest.setup.ts` (global mock) when API exports change
- Use `flushPromises()` for async operations in tests
- Target: <300ms per test

## Key References

- Full testing guide: `docs/agents/unit-test.md` (Vue/TS section)
- i18n guide: `src/ui/src/locales/README.md`
- API client: `src/ui/src/lib/api.ts`
- Auth store: `src/ui/src/stores/auth.ts`
