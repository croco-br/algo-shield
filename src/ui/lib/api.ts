import { uiConfig } from './config';

export interface ApiError {
	error: string;
}

async function request<T>(
	endpoint: string,
	options: RequestInit = {}
): Promise<T> {
	const token = typeof window !== 'undefined' ? localStorage.getItem('auth_token') : null;
	
	const headers: Record<string, string> = {
		'Content-Type': 'application/json',
		...(options.headers as Record<string, string> || {}),
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
		if (error instanceof Error) {
			if (error.name === 'AbortError') {
				throw new Error('Request timeout');
			}
			// Handle network errors (Failed to fetch, CORS, etc.)
			if (error.message.includes('Failed to fetch') || error.message.includes('NetworkError')) {
				throw new Error(`Unable to connect to API at ${uiConfig.api.baseUrl}. Please ensure the API server is running.`);
			}
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
