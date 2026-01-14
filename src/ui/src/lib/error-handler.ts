import { useI18n } from 'vue-i18n';

/**
 * API Error Response structure from backend
 */
export interface APIError {
	code?: string;
	message?: string;
	error?: string;
}

/**
 * Extracts a user-friendly error message from an error object
 *
 * This function handles errors from the backend API that follow the structure:
 * { code: "ERROR_CODE", message: "Error message" }
 *
 * It translates error codes to localized messages using i18n.
 * Falls back to the raw message if no translation is available.
 *
 * @param error - The error object (can be Error, string, or APIError)
 * @returns A user-friendly error message
 */
export function getErrorMessage(error: unknown): string {
	const { t } = useI18n();

	// Handle string errors
	if (typeof error === 'string') {
		return error;
	}

	// Handle Error objects
	if (error instanceof Error) {
		try {
			// Try to parse the error message as JSON (API errors)
			const parsed = JSON.parse(error.message) as APIError;
			if (parsed.code) {
				// Check if we have a translation for this error code
				const translationKey = `errors.${parsed.code}`;
				const translated = t(translationKey);

				// If translation exists and is not the key itself, use it
				if (translated !== translationKey) {
					return translated;
				}
			}

			// Fall back to the message field
			return parsed.message || parsed.error || error.message;
		} catch {
			// Not JSON, return the raw error message
			return error.message;
		}
	}

	// Handle APIError objects directly
	if (error && typeof error === 'object') {
		const apiError = error as APIError;

		if (apiError.code) {
			// Check if we have a translation for this error code
			const translationKey = `errors.${apiError.code}`;
			const translated = t(translationKey);

			// If translation exists and is not the key itself, use it
			if (translated !== translationKey) {
				return translated;
			}
		}

		// Fall back to message or error field
		return apiError.message || apiError.error || t('errors.INTERNAL_ERROR');
	}

	// Unknown error type
	return t('errors.INTERNAL_ERROR');
}

/**
 * Checks if an error is a 403 Forbidden error
 *
 * @param error - The error object
 * @returns true if the error is a 403 Forbidden error
 */
export function isForbiddenError(error: unknown): boolean {
	if (error instanceof Error) {
		return error.message.includes('403') || error.message.includes('Forbidden');
	}
	return false;
}

/**
 * Checks if an error is a 401 Unauthorized error
 *
 * @param error - The error object
 * @returns true if the error is a 401 Unauthorized error
 */
export function isUnauthorizedError(error: unknown): boolean {
	if (error instanceof Error) {
		return error.message.includes('401') || error.message.includes('Unauthorized');
	}
	return false;
}
