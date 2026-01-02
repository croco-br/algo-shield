import { uiConfig } from './config';

export interface ApiError {
	error: string;
}

async function request<T>(
	endpoint: string,
	options: RequestInit = {}
): Promise<T> {
	if (!uiConfig.api.baseUrl) {
		throw new Error('API base URL is not configured. Please rebuild the UI container with VITE_API_URL set.');
	}
	if (uiConfig.api.baseUrl.includes('://api:') || uiConfig.api.baseUrl.includes('://postgres:') || uiConfig.api.baseUrl.includes('://redis:')) {
		throw new Error(`API baseUrl uses container hostname (${uiConfig.api.baseUrl}). Browser cannot resolve container hostnames. Use http://localhost:8080 instead.`);
	}
	const token = localStorage.getItem('auth_token');
	
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
		// Explicitly set mode to 'cors' to ensure CORS is handled correctly
		const response = await fetch(`${uiConfig.api.baseUrl}${endpoint}`, {
			...options,
			headers,
			signal: controller.signal,
			mode: 'cors',
			credentials: 'omit', // Explicitly set to avoid credential issues
		});

		clearTimeout(timeoutId);

		// Check if response is JSON
		const contentType = response.headers.get('content-type');
		const isJson = contentType && contentType.includes('application/json');

		if (!response.ok) {
			let errorMessage = `HTTP ${response.status}`;
			
			if (isJson) {
				try {
					const error: ApiError = await response.json();
					errorMessage = error.error || errorMessage;
				} catch {
					// If JSON parse fails, use default message
					errorMessage = `Server error (${response.status})`;
				}
			} else {
				// Response is HTML or other non-JSON format
				// Clone the response to read text without consuming the original body
				const text = await response.clone().text();
				if (text.includes('<!DOCTYPE') || text.includes('<html')) {
					if (response.status === 404) {
						errorMessage = 'API endpoint not found. Please check if the API server is running and configured correctly.';
					} else if (response.status === 502 || response.status === 503) {
						errorMessage = 'API server is not available. Please ensure the API server is running.';
					} else {
						errorMessage = `Server error (${response.status}). The API may not be available.`;
					}
				} else {
					errorMessage = text.substring(0, 200) || errorMessage;
				}
			}
			
			throw new Error(errorMessage);
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
			if (!text || text.trim() === '') {
				return undefined as T;
			}
			// Try to parse as JSON anyway (some servers don't set content-type correctly)
			try {
				return JSON.parse(text);
			} catch {
				throw new Error(`Expected JSON response but received: ${contentType || 'unknown format'}`);
			}
		}
	} catch (error) {
		clearTimeout(timeoutId);
		if (error instanceof Error) {
			if (error.name === 'AbortError') {
				throw new Error('Request timeout. The server took too long to respond.');
			}
			// Handle network errors (Failed to fetch, CORS, etc.)
			if (error.message.includes('Failed to fetch') || error.message.includes('NetworkError')) {
				const apiUrl = uiConfig.api.baseUrl || 'the API server';
				throw new Error(`Unable to connect to server.`);
			}
			// Re-throw if it's already our formatted error
			throw error;
		}
		throw new Error('An unexpected error occurred');
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
