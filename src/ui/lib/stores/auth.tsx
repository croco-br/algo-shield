'use client';

import { createContext, useContext, useEffect, useState, ReactNode } from 'react';
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

interface AuthContextType {
	user: User | null;
	loading: boolean;
	setToken: (token: string) => Promise<void>;
	logout: () => Promise<void>;
	refresh: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
	const [user, setUser] = useState<User | null>(null);
	const [loading, setLoading] = useState(true);

	async function loadUserFromToken() {
		try {
			const token = typeof window !== 'undefined' ? localStorage.getItem('auth_token') : null;
			if (!token) {
				setUser(null);
				setLoading(false);
				return;
			}

			const userData = await api.get<User>('/api/v1/auth/me');
			setUser(userData);
		} catch (e) {
			// Token invalid, clear it
			if (typeof window !== 'undefined') {
				localStorage.removeItem('auth_token');
			}
			setUser(null);
		} finally {
			setLoading(false);
		}
	}

	useEffect(() => {
		loadUserFromToken();
	}, []);

	const setToken = async (token: string) => {
		if (typeof window !== 'undefined') {
			localStorage.setItem('auth_token', token);
			await loadUserFromToken();
		}
	};

	const logout = async () => {
		try {
			await api.post('/api/v1/auth/logout');
		} catch (e) {
			// Ignore errors on logout
		}
		setUser(null);
		if (typeof window !== 'undefined') {
			localStorage.removeItem('auth_token');
		}
	};

	const refresh = loadUserFromToken;

	return (
		<AuthContext.Provider value={{ user, loading, setToken, logout, refresh }}>
			{children}
		</AuthContext.Provider>
	);
}

export function useAuth() {
	const context = useContext(AuthContext);
	if (context === undefined) {
		throw new Error('useAuth must be used within an AuthProvider');
	}
	return context;
}
