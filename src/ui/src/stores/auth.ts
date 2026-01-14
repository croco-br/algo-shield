import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { api, setTokenGetter, setCsrfTokenGetter } from '@/lib/api';

export interface Role {
	id: string;
	name: string;
	description: string;
}

export interface User {
	id: string;
	email: string;
	name: string;
	auth_type: 'local' | 'sso';
	picture_url?: string;
	roles?: Role[];
	active: boolean;
}

export const useAuthStore = defineStore('auth', () => {
	const user = ref<User | null>(null);
	const loading = ref(true);
	// SECURITY: Store token in memory only (never localStorage to prevent XSS attacks)
	// This means token will be lost on page reload, requiring re-authentication
	const token = ref<string | null>(null);
	// SECURITY: CSRF token stored in memory alongside JWT token
	const csrfToken = ref<string | null>(null);

	async function loadUserFromToken() {
		try {
			if (!token.value) {
				user.value = null;
				loading.value = false;
				return;
			}

			const userData = await api.get<User>('/api/v1/auth/me');
			user.value = userData;
		} catch (e) {
			// Token invalid, clear it
			token.value = null;
			user.value = null;
		} finally {
			loading.value = false;
		}
	}

	async function setToken(newToken: string, newCsrfToken?: string) {
		token.value = newToken;
		csrfToken.value = newCsrfToken || null;
		await loadUserFromToken();
	}

	async function logout() {
		try {
			await api.post('/api/v1/auth/logout');
		} catch (e) {
			// Ignore errors on logout
		}
		user.value = null;
		token.value = null;
		csrfToken.value = null;
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
		return user.value?.roles?.some((role) => role.name === 'admin') || false;
	});

	// SECURITY: Register token getters with API module to avoid circular dependencies
	// This allows api.ts to get the tokens from memory without importing the store
	setTokenGetter(getToken);
	setCsrfTokenGetter(getCsrfToken);

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
