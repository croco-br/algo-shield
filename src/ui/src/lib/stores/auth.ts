import { writable } from 'svelte/store';
import { browser } from '$app/environment';
import { api } from '../api';

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

function createAuthStore() {
	const { subscribe, set, update } = writable<User | null>(null);

	async function loadUserFromToken() {
		try {
			const user = await api.get<User>('/api/v1/auth/me');
			set(user);
		} catch (e) {
			// Token invalid, clear it
			if (browser) {
				localStorage.removeItem('auth_token');
			}
			set(null);
		}
	}

	// Load user from token on init
	if (browser) {
		const token = localStorage.getItem('auth_token');
		if (token) {
			// Verify token and load user
			loadUserFromToken();
		}
	}

	return {
		subscribe,
		setToken: async (token: string) => {
			if (browser) {
				localStorage.setItem('auth_token', token);
				await loadUserFromToken();
			}
		},
		logout: async () => {
			try {
				await api.post('/api/v1/auth/logout');
			} catch (e) {
				// Ignore errors on logout
			}
			set(null);
			if (browser) {
				localStorage.removeItem('auth_token');
			}
		},
		refresh: loadUserFromToken
	};
}

export const authStore = createAuthStore();
