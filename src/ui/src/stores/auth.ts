import { defineStore } from "pinia";
import { ref, computed } from "vue";
import {
  api,
  setTokenGetter,
  setCsrfTokenGetter,
  setRefreshTokenFn,
  setForceLogoutFn,
  resetLogoutState,
} from "@/lib/api";

export interface Role {
  id: string;
  name: string;
  description: string;
}

export interface User {
  id: string;
  email: string;
  name: string;
  auth_type: "local" | "sso";
  picture_url?: string;
  roles?: Role[];
  active: boolean;
}

export const useAuthStore = defineStore("auth", () => {
  const user = ref<User | null>(null);
  const loading = ref(true);
  // SECURITY: Store access token in memory only (XSS protection)
  // Refresh token in localStorage is safer as it's only used for getting new access tokens
  const token = ref<string | null>(null);
  const csrfToken = ref<string | null>(null);
  const refreshToken = ref<string | null>(
    localStorage.getItem("refresh_token"),
  );

  async function loadUserFromToken() {
    try {
      if (!token.value) {
        // Try to use refresh token if available
        if (refreshToken.value) {
          await refreshAccessToken();
          return;
        }
        user.value = null;
        loading.value = false;
        return;
      }

      const userData = await api.get<User>("/api/v1/auth/me");
      user.value = userData;
    } catch (e) {
      // Token invalid, try refresh
      if (refreshToken.value) {
        try {
          await refreshAccessToken();
        } catch (refreshError) {
          // Refresh failed, clear everything
          clearTokens();
        }
      } else {
        clearTokens();
      }
    } finally {
      loading.value = false;
    }
  }

  async function refreshAccessToken() {
    // Guard: don't hit the API if there's no refresh token (prevents 429 rate-limit loops)
    if (!refreshToken.value) {
      clearTokens();
      throw new Error("No refresh token available");
    }

    try {
      const response = await api.post<{
        token: string;
        refresh_token?: string;
        csrf_token?: string;
      }>("/api/v1/auth/refresh", {
        refresh_token: refreshToken.value,
      });
      token.value = response.token;
      csrfToken.value = response.csrf_token || null;

      // Store new refresh token if returned (rotation)
      if (response.refresh_token) {
        refreshToken.value = response.refresh_token;
        localStorage.setItem("refresh_token", response.refresh_token);
      }

      // Load user data with new token
      const userData = await api.get<User>("/api/v1/auth/me");
      user.value = userData;
    } catch (e: unknown) {
      // Only clear tokens when refresh failed with 401 (invalid/expired refresh token).
      // On network/5xx leave tokens so user can retry (e.g. generate again).
      const statusCode =
        e != null && typeof (e as { statusCode?: number }).statusCode === "number"
          ? (e as { statusCode: number }).statusCode
          : undefined;
      if (statusCode === 401) {
        clearTokens();
      }
      throw e;
    }
  }

  function clearTokens() {
    token.value = null;
    csrfToken.value = null;
    refreshToken.value = null;
    user.value = null;
    localStorage.removeItem("refresh_token");
  }

  async function setToken(
    newToken: string,
    newCsrfToken?: string,
    newRefreshToken?: string,
  ) {
    token.value = newToken;
    csrfToken.value = newCsrfToken || null;

    // Only persist refresh token (safer than access token)
    if (newRefreshToken) {
      refreshToken.value = newRefreshToken;
      localStorage.setItem("refresh_token", newRefreshToken);
    }

    resetLogoutState();
    await loadUserFromToken();
  }

  async function logout() {
    try {
      await api.post("/api/v1/auth/logout");
    } catch (e) {
      // Ignore errors on logout
    }
    clearTokens();
  }

  async function refresh() {
    await loadUserFromToken();
  }

  function getToken(): string | null {
    return token.value;
  }

  function getCsrfToken(): string | null {
    return csrfToken.value;
  }

  const isAdmin = computed(() => {
    return user.value?.roles?.some((role) => role.name === "admin") || false;
  });

  // SECURITY: Register token getters and refresh function with API module to avoid circular dependencies
  // This allows api.ts to get the tokens from memory without importing the store
  setTokenGetter(getToken);
  setCsrfTokenGetter(getCsrfToken);
  setRefreshTokenFn(refreshAccessToken);
  setForceLogoutFn(() => {
    clearTokens();
    window.location.href = "/login";
  });

  // Initialize on store creation (token will be null on page load)
  loadUserFromToken();

  return {
    user,
    loading,
    setToken,
    logout,
    refresh,
    getToken,
    getCsrfToken,
    isAdmin,
  };
});
