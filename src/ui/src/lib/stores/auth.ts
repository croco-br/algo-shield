import { writable } from 'svelte/store';
import { browser } from '$app/environment';

export interface User {
	id: string;
	email: string;
	name: string;
	authType: 'local' | 'sso';
}

function createAuthStore() {
	const { subscribe, set, update } = writable<User | null>(null);

	// Load user from localStorage on init
	if (browser) {
		const stored = localStorage.getItem('auth_user');
		if (stored) {
			try {
				set(JSON.parse(stored));
			} catch (e) {
				localStorage.removeItem('auth_user');
			}
		}
	}

	return {
		subscribe,
		login: (email: string, password: string) => {
			// Mock local login
			const user: User = {
				id: Math.random().toString(36).substring(7),
				email,
				name: email.split('@')[0],
				authType: 'local'
			};
			set(user);
			if (browser) {
				localStorage.setItem('auth_user', JSON.stringify(user));
			}
			return user;
		},
		loginSSO: (provider: string) => {
			// Mock SSO login
			const user: User = {
				id: Math.random().toString(36).substring(7),
				email: `user@${provider}.com`,
				name: `${provider} User`,
				authType: 'sso'
			};
			set(user);
			if (browser) {
				localStorage.setItem('auth_user', JSON.stringify(user));
			}
			return user;
		},
		logout: () => {
			set(null);
			if (browser) {
				localStorage.removeItem('auth_user');
			}
		}
	};
}

export const authStore = createAuthStore();
