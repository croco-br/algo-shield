import { uiConfig } from './config';

export interface ApiError {
	error: string;
}

async function request<T>(
	endpoint: string,
	options: RequestInit = {}
): Promise<T> {
	const token = typeof window !== 'undefined' ? localStorage.getItem('auth_token') : null;
	
	const headers: HeadersInit = {
		'Content-Type': 'application/json',
		...options.headers,
	};

	if (token) {
		headers['Authorization'] = `Bearer ${token}`;
	}

	// Create AbortController for timeout
	const controller = new AbortController();
	const timeoutId = setTimeout(() => controller.abort(), uiConfig.api.timeout);

	try {
		const response = await fetch(`${uiConfig.api.baseUrl}${endpoint}`, {
			...options,
			headers,
			signal: controller.signal,
		});

		clearTimeout(timeoutId);

		if (!response.ok) {
			const error: ApiError = await response.json().catch(() => ({ error: 'Unknown error' }));
			throw new Error(error.error || `HTTP ${response.status}`);
		}

		return response.json();
	} catch (error) {
		clearTimeout(timeoutId);
		if (error instanceof Error && error.name === 'AbortError') {
			throw new Error('Request timeout');
		}
		throw error;
	}
}

export const api = {
	get: <T>(endpoint: string) => request<T>(endpoint, { method: 'GET' }),
	post: <T>(endpoint: string, data?: unknown) =>
		request<T>(endpoint, {
			method: 'POST',
			body: data ? JSON.stringify(data) : undefined,
		}),
	put: <T>(endpoint: string, data?: unknown) =>
		request<T>(endpoint, {
			method: 'PUT',
			body: data ? JSON.stringify(data) : undefined,
		}),
	delete: <T>(endpoint: string) => request<T>(endpoint, { method: 'DELETE' }),
};
