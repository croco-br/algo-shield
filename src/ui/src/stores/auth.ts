import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { api } from '@/lib/api';

export interface Role {
	id: string;
	name: string;
	description: string;
}

export interface User {
	id: string;
	email: string;
	name: string;
	authType: 'local' | 'sso';
	picture_url?: string;
	roles?: Role[];
	active: boolean;
}

export const useAuthStore = defineStore('auth', () => {
	const user = ref<User | null>(null);
	const loading = ref(true);

	async function loadUserFromToken() {
		try {
			const token = localStorage.getItem('auth_token');
			if (!token) {
				user.value = null;
				loading.value = false;
				return;
			}

			const userData = await api.get<User>('/api/v1/auth/me');
			user.value = userData;
		} catch (e) {
			// Token invalid, clear it
			localStorage.removeItem('auth_token');
			user.value = null;
		} finally {
			loading.value = false;
		}
	}

	async function setToken(token: string) {
		localStorage.setItem('auth_token', token);
		await loadUserFromToken();
	}

	async function logout() {
		try {
			await api.post('/api/v1/auth/logout');
		} catch (e) {
			// Ignore errors on logout
		}
		user.value = null;
		localStorage.removeItem('auth_token');
	}

	async function refresh() {
		await loadUserFromToken();
	}

	const isAdmin = computed(() => {
		return user.value?.roles?.some((role) => role.name === 'admin') || false;
	});

	// Initialize on store creation
	loadUserFromToken();

	return {
		user,
		loading,
		setToken,
		logout,
		refresh,
		isAdmin,
	};
});
