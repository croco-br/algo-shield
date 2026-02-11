# Generate Synthetic Data and Token/Session Issues

## Why token/session can break after using "Generate synthetic data"

"Generate synthetic data" sends **many POST requests in a row** (batches of 100 events per request). That can surface or trigger several behaviors:

### 1. **401 during generate → refresh → force logout**

- The **access token** (JWT) has a limited lifetime (e.g. 24h). If it is close to expiry and you run a long generate (many batches), one of the batch requests can get **401 Unauthorized**.
- The API client then tries to **refresh** the token. If refresh fails for any reason (expired refresh token, network error, or backend 5xx), the app **force-logs you out** and redirects to `/login`.
- So the "session bug" is often: **token expired mid-flow → refresh failed → you’re logged out**.

### 2. **`isLoggingOut` not reset after re-login (fixed)**

- After a force logout, a flag `isLoggingOut` was set and **never reset**. If you logged in again in the same tab, the next 401 would **not** trigger a redirect to login (the app thought it was already "logging out"), so you could end up in a broken state (e.g. requests fail with 401 but you stay on the page).
- **Fix:** The app now calls `resetLogoutState()` when you set a new token (login or successful refresh), so after re-login the session and 401→redirect behavior work again.

### 3. **Client vs backend timeout**

- Default **client** request timeout is **30 seconds** (`VITE_API_TIMEOUT`). The **backend** batch timeout for generate-events is **30 seconds** (`API_TIMEOUT_BATCH_OPERATION`). A single large batch can hit the limit; the client may see "Request timeout" and the request fails. That does **not** clear the token, but the generate run looks failed.

### 4. **CSRF**

- Each POST (including each generate-events batch) sends the **CSRF token**. The backend does **not** rotate CSRF on each request; it only validates. So CSRF is not the cause of "session bugs" unless Redis (where CSRF is stored) has issues or the token expired (24h TTL).

---

## What was changed in code

- **`src/ui/src/lib/api.ts`**  
  - Exported **`resetLogoutState()`** so the auth store can clear the logout flag after a successful login.  
  - **Only force logout when refresh failed with 401.** If generate gets 401 and the refresh call fails for another reason (network, timeout, 5xx), we no longer redirect to login; we throw so you see the error and can retry.  
  - **Skip the "Unrecoverable 401" force logout when we already tried refresh and it failed for non-401.** Previously we always fell through to that block after a failed refresh, so we were still redirecting to login even when refresh failed due to network/5xx. Now we set `attemptedRefreshAndFailedWithNon401` in that case and skip the second force logout.  
  - Thrown errors now include **`statusCode`** so the auth layer can tell 401 from other failures.

- **`src/ui/src/stores/auth.ts`**  
  - In **`setToken()`**, the store calls **`resetLogoutState()`** before `loadUserFromToken()`.  
  - **Only clear tokens when refresh returns 401.** If refresh fails with network/5xx, we no longer clear tokens; we just throw. You keep your session and can retry (e.g. run generate again).

---

## What you can do as a user

1. **Avoid very long runs**  
   Prefer smaller "Event count" values or run generate multiple times so each run finishes well within token and timeout limits.

2. **If you get logged out right after generate**  
   Log in again; the app should now behave normally (including redirecting to login on 401 when needed).

3. **If you see "Request timeout"**  
   Reduce batch size (event count per run) or increase `VITE_API_TIMEOUT` and/or `API_TIMEOUT_BATCH_OPERATION` if you control the environment.

4. **If you see CSRF/403**  
   Check that you’re not mixing tabs (e.g. one tab logged in, another logging out) and that Redis is healthy; CSRF is stored in Redis with a 24h TTL.
