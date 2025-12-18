/**
 * UI Configuration
 * 
 * All configuration values can be overridden via environment variables
 * prefixed with NEXT_PUBLIC_ (Next.js convention for client-side env vars)
 */

export interface UIConfig {
	api: {
		baseUrl: string;
		timeout: number; // Request timeout in milliseconds
		retry: {
			maxAttempts: number;
			initialDelay: number; // Initial retry delay in milliseconds
			maxDelay: number; // Maximum retry delay in milliseconds
			multiplier: number; // Exponential backoff multiplier
		};
	};
	ui: {
		toast: {
			duration: number; // Toast notification duration in milliseconds
		};
		polling: {
			interval: number; // Polling interval for real-time updates in milliseconds
		};
	};
}

function getEnvInt(key: string, defaultValue: number): number {
	const value = process.env[key];
	if (value === undefined || value === '') {
		return defaultValue;
	}
	const parsed = parseInt(value, 10);
	return isNaN(parsed) ? defaultValue : parsed;
}

function getEnvFloat(key: string, defaultValue: number): number {
	const value = process.env[key];
	if (value === undefined || value === '') {
		return defaultValue;
	}
	const parsed = parseFloat(value);
	return isNaN(parsed) ? defaultValue : parsed;
}

function getEnvString(key: string, defaultValue: string): string {
	const value = process.env[key];
	return value !== undefined && value !== '' ? value : defaultValue;
}

export const uiConfig: UIConfig = {
	api: {
		// Use relative URL to leverage Next.js rewrites (works in both Docker and local dev)
		// If NEXT_PUBLIC_API_URL is set, use it (for external API or custom setup)
		// Otherwise, use empty string to use relative URLs via Next.js proxy
		baseUrl: getEnvString('NEXT_PUBLIC_API_URL', ''),
		timeout: getEnvInt('NEXT_PUBLIC_API_TIMEOUT', 30000), // 30 seconds default
		retry: {
			maxAttempts: getEnvInt('NEXT_PUBLIC_API_RETRY_MAX_ATTEMPTS', 3),
			initialDelay: getEnvInt('NEXT_PUBLIC_API_RETRY_INITIAL_DELAY', 1000), // 1 second
			maxDelay: getEnvInt('NEXT_PUBLIC_API_RETRY_MAX_DELAY', 10000), // 10 seconds
			multiplier: getEnvFloat('NEXT_PUBLIC_API_RETRY_MULTIPLIER', 2.0),
		},
	},
	ui: {
		toast: {
			duration: getEnvInt('NEXT_PUBLIC_UI_TOAST_DURATION', 5000), // 5 seconds
		},
		polling: {
			interval: getEnvInt('NEXT_PUBLIC_UI_POLLING_INTERVAL', 10000), // 10 seconds
		},
	},
};
