import { uiConfig } from "./config";

export interface ApiError {
  error: string;
}

// SECURITY NOTE: Token getter functions are injected to avoid circular dependencies
// The auth store will provide the tokens via getToken() and getCsrfToken() functions
// Using objects to store the getters helps with module isolation in tests
const tokenGetterHolder: { getter: (() => string | null) | null } = {
  getter: null,
};

const csrfTokenGetterHolder: { getter: (() => string | null) | null } = {
  getter: null,
};

// Refresh token function holder — injected by auth store to avoid circular dependencies
const refreshTokenHolder: { refreshFn: (() => Promise<void>) | null } = {
  refreshFn: null,
};

// Force logout function holder — injected by auth store to clear tokens and redirect
const forceLogoutHolder: { fn: (() => void) | null } = {
  fn: null,
};

// Module-level promise to deduplicate concurrent refresh attempts
let activeRefreshPromise: Promise<void> | null = null;

// Flag to prevent multiple force-logout redirects from causing an infinite reload loop
let isLoggingOut = false;

// Flag to indicate a refresh is in progress (used for request queuing)
let isRefreshing = false;

// Queue for requests that arrived while a token refresh is in progress
interface QueuedRequest {
  resolve: (value: unknown) => void;
  reject: (reason: unknown) => void;
  endpoint: string;
  options: RequestInit;
}
let requestQueue: QueuedRequest[] = [];

function replayQueuedRequests(): void {
  const queue = [...requestQueue];
  requestQueue = [];
  for (const item of queue) {
    request(item.endpoint, item.options, true)
      .then(item.resolve)
      .catch(item.reject);
  }
}

function rejectQueuedRequests(reason: unknown): void {
  const queue = [...requestQueue];
  requestQueue = [];
  for (const item of queue) {
    item.reject(reason);
  }
}

export function setTokenGetter(getter: () => string | null): void {
  tokenGetterHolder.getter = getter;
}

export function setCsrfTokenGetter(getter: () => string | null): void {
  csrfTokenGetterHolder.getter = getter;
}

export function setRefreshTokenFn(fn: () => Promise<void>): void {
  refreshTokenHolder.refreshFn = fn;
}

export function setForceLogoutFn(fn: () => void): void {
  forceLogoutHolder.fn = fn;
}

/** Reset logout state after successful login so future 401s can trigger redirect again. */
export function resetLogoutState(): void {
  isLoggingOut = false;
  isRefreshing = false;
  requestQueue = [];
}

// Get synthetic mode from localStorage (synced by systemMode store)
// NOTE: Synthetic mode is safe to store in localStorage as it's not sensitive data
function getSyntheticMode(): boolean {
  try {
    const mode = localStorage.getItem("synthetic_mode");
    return mode === "true";
  } catch {
    // If localStorage is not available, default to false
    return false;
  }
}

// Set synthetic mode in localStorage (called by systemMode store)
export function setSyntheticModeStorage(enabled: boolean): void {
  try {
    localStorage.setItem("synthetic_mode", String(enabled));
  } catch {
    // If localStorage is not available, silently fail
    // This can happen in some environments but shouldn't break the app
  }
}

async function request<T>(
  endpoint: string,
  options: RequestInit = {},
  _isRetry = false,
): Promise<T> {
  if (!uiConfig.api.baseUrl) {
    throw new Error(
      "API base URL is not configured. Please rebuild the UI container with VITE_API_URL set.",
    );
  }
  if (
    uiConfig.api.baseUrl.includes("://api:") ||
    uiConfig.api.baseUrl.includes("://postgres:") ||
    uiConfig.api.baseUrl.includes("://redis:")
  ) {
    throw new Error(
      `API baseUrl uses container hostname (${uiConfig.api.baseUrl}). Browser cannot resolve container hostnames. Use http://localhost:8080 instead.`,
    );
  }
  // SECURITY: Get tokens from memory-only Pinia store (never localStorage)
  const token = tokenGetterHolder.getter ? tokenGetterHolder.getter() : null;
  const csrfToken = csrfTokenGetterHolder.getter
    ? csrfTokenGetterHolder.getter()
    : null;
  const syntheticMode = getSyntheticMode();

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    "X-Synthetic-Mode": String(syntheticMode),
    ...((options.headers as Record<string, string>) || {}),
  };

  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  // SECURITY: Add CSRF token for state-changing requests (POST, PUT, PATCH, DELETE)
  // The CSRF middleware on backend validates this token for all non-GET requests
  const method = options.method?.toUpperCase();
  if (
    csrfToken &&
    method &&
    ["POST", "PUT", "PATCH", "DELETE"].includes(method)
  ) {
    headers["X-CSRF-Token"] = csrfToken;
  }

  // Create AbortController for timeout
  const controller = new AbortController();
  const timeoutId = setTimeout(() => controller.abort(), uiConfig.api.timeout);

  try {
    // Explicitly set mode to 'cors' to ensure CORS is handled correctly
    const response = await fetch(`${uiConfig.api.baseUrl}${endpoint}`, {
      ...options,
      headers,
      signal: controller.signal,
      mode: "cors",
      credentials: "omit", // Explicitly set to avoid credential issues
    });

    clearTimeout(timeoutId);

    // Check if response is JSON
    const contentType = response.headers.get("content-type");
    const isJson = contentType && contentType.includes("application/json");

    if (!response.ok) {
      const isAuthEndpoint =
        endpoint.includes("/auth/login") ||
        endpoint.includes("/auth/register") ||
        endpoint.includes("/auth/refresh");

      // Detect 403 CSRF errors (distinct from 403 permission errors)
      let isCsrf403 = false;
      if (response.status === 403 && !isAuthEndpoint) {
        try {
          const body = await response.clone().text();
          isCsrf403 = /csrf/i.test(body);
        } catch {
          // If we can't read the body, treat as non-CSRF 403
        }
      }

      // Recoverable errors: 401 (expired access token) or 403 CSRF (expired/invalid CSRF token)
      // Both are resolved by refreshing the session (which returns new access + CSRF tokens)
      const isRecoverable =
        (response.status === 401 || isCsrf403) &&
        token &&
        refreshTokenHolder.refreshFn &&
        !_isRetry &&
        !isAuthEndpoint &&
        !isLoggingOut;

      let attemptedRefreshAndFailedWithNon401 = false;

      if (isRecoverable) {
        // If a refresh is already in progress, queue this request for replay after refresh
        if (isRefreshing && activeRefreshPromise) {
          return new Promise<T>((resolve, reject) => {
            requestQueue.push({ resolve: resolve as (value: unknown) => void, reject, endpoint, options });
          });
        }

        try {
          isRefreshing = true;
          if (!activeRefreshPromise) {
            activeRefreshPromise = refreshTokenHolder
              .refreshFn!()
              .finally(() => {
                activeRefreshPromise = null;
                isRefreshing = false;
              });
          }
          await activeRefreshPromise;
          replayQueuedRequests();
          return request<T>(endpoint, options, true);
        } catch (refreshErr: unknown) {
          rejectQueuedRequests(refreshErr);
          const statusCode =
            refreshErr != null &&
            typeof (refreshErr as { statusCode?: number }).statusCode ===
              "number"
              ? (refreshErr as { statusCode: number }).statusCode
              : undefined;
          if (statusCode === 401 && forceLogoutHolder.fn && !isLoggingOut) {
            isLoggingOut = true;
            forceLogoutHolder.fn();
          } else if (statusCode !== 401) {
            attemptedRefreshAndFailedWithNon401 = true;
          }
        }
      }

      if (
        response.status === 401 &&
        !isAuthEndpoint &&
        forceLogoutHolder.fn &&
        !isLoggingOut &&
        !attemptedRefreshAndFailedWithNon401
      ) {
        isLoggingOut = true;
        forceLogoutHolder.fn();
      }

      let errorMessage = `HTTP ${response.status}`;

      if (isJson) {
        try {
          const error: ApiError = await response.json();
          errorMessage = error.error || errorMessage;
        } catch {
          errorMessage = `Server error (${response.status})`;
        }
      } else {
        try {
          const text = await response.clone().text();
          if (text.includes("<!DOCTYPE") || text.includes("<html")) {
            if (response.status === 404) {
              errorMessage =
                "API endpoint not found. Please check if the API server is running and configured correctly.";
            } else if (response.status === 502 || response.status === 503) {
              errorMessage =
                "API server is not available. Please ensure the API server is running.";
            } else {
              errorMessage = `Server error (${response.status}). The API may not be available.`;
            }
          } else {
            errorMessage = text.substring(0, 200) || errorMessage;
          }
        } catch {
          errorMessage = `Server error (${response.status})`;
        }
      }

      const err = new Error(errorMessage) as Error & { statusCode?: number };
      err.statusCode = response.status;
      throw err;
    }

    // Handle 204 No Content responses (common for DELETE operations)
    if (response.status === 204) {
      return undefined as T;
    }

    // Parse response as JSON
    if (isJson) {
      return response.json();
    } else {
      // For successful responses without JSON content-type, try to parse anyway
      // or return undefined if empty
      const text = await response.text();
      if (!text || text.trim() === "") {
        return undefined as T;
      }
      // Try to parse as JSON anyway (some servers don't set content-type correctly)
      try {
        return JSON.parse(text);
      } catch {
        throw new Error(
          `Expected JSON response but received: ${contentType || "unknown format"}`,
        );
      }
    }
  } catch (error) {
    clearTimeout(timeoutId);
    if (error instanceof Error) {
      if (error.name === "AbortError") {
        throw new Error(
          "Request timeout. The server took too long to respond.",
        );
      }
      // Handle network errors (Failed to fetch, CORS, etc.)
      if (
        error.message.includes("Failed to fetch") ||
        error.message.includes("NetworkError")
      ) {
        throw new Error(`Unable to connect to server.`);
      }
      // Re-throw if it's already our formatted error
      throw error;
    }
    throw new Error("An unexpected error occurred");
  }
}

export const api = {
  get: <T>(endpoint: string) => request<T>(endpoint, { method: "GET" }),
  post: <T>(endpoint: string, data?: unknown) =>
    request<T>(endpoint, {
      method: "POST",
      body: data ? JSON.stringify(data) : undefined,
    }),
  put: <T>(endpoint: string, data?: unknown) =>
    request<T>(endpoint, {
      method: "PUT",
      body: data ? JSON.stringify(data) : undefined,
    }),
  patch: <T>(endpoint: string, data?: unknown) =>
    request<T>(endpoint, {
      method: "PATCH",
      body: data ? JSON.stringify(data) : undefined,
    }),
  delete: <T>(endpoint: string) => request<T>(endpoint, { method: "DELETE" }),
};
