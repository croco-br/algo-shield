import { describe, it, expect, vi, beforeEach } from 'vitest'
import { api, setTokenGetter, setSyntheticModeStorage } from './api'

vi.mock('./config', () => ({
	uiConfig: {
		api: {
			baseUrl: 'http://localhost:8080',
			timeout: 30000,
		},
	},
}))

describe('api', () => {
	beforeEach(() => {
		localStorage.clear()
		setTokenGetter(() => null)
		globalThis.fetch = vi.fn()
	})

	describe('get', () => {
		it('makes GET request with correct URL and headers', async () => {
			const mockResponse = { data: 'test' }
			globalThis.fetch = vi.fn().mockResolvedValue({
				ok: true,
				status: 200,
				headers: new Headers({ 'content-type': 'application/json' }),
				json: async () => mockResponse,
			})

			const result = await api.get('/test')

			expect(globalThis.fetch).toHaveBeenCalledWith(
				'http://localhost:8080/test',
				expect.objectContaining({
					method: 'GET',
					headers: expect.objectContaining({
						'Content-Type': 'application/json',
						'X-Synthetic-Mode': 'false',
					}),
					mode: 'cors',
					credentials: 'omit',
				})
			)
			expect(result).toEqual(mockResponse)
		})

		it('does not include Authorization header when token getter returns null', async () => {
			setTokenGetter(() => null)

			globalThis.fetch = vi.fn().mockResolvedValue({
				ok: true,
				status: 200,
				headers: new Headers({ 'content-type': 'application/json' }),
				json: async () => ({}),
			})

			await api.get('/test')

			const fetchCalls = (globalThis.fetch as any).mock.calls
			const [, options] = fetchCalls[0]

			expect(options.headers['Authorization']).toBeUndefined()
		})

		it('includes X-Synthetic-Mode header based on localStorage', async () => {
			setSyntheticModeStorage(true)

			globalThis.fetch = vi.fn().mockResolvedValue({
				ok: true,
				status: 200,
				headers: new Headers({ 'content-type': 'application/json' }),
				json: async () => ({}),
			})

			await api.get('/test')

			expect(globalThis.fetch).toHaveBeenCalledWith(
				expect.any(String),
				expect.objectContaining({
					headers: expect.objectContaining({
						'X-Synthetic-Mode': 'true',
					}),
				})
			)
		})
	})

	describe('post', () => {
		it('makes POST request with JSON body', async () => {
			const requestData = { email: 'test@example.com', password: 'password123' }
			const mockResponse = { token: 'jwt-token' }

			globalThis.fetch = vi.fn().mockResolvedValue({
				ok: true,
				status: 200,
				headers: new Headers({ 'content-type': 'application/json' }),
				json: async () => mockResponse,
			})

			const result = await api.post('/auth/login', requestData)

			expect(globalThis.fetch).toHaveBeenCalledWith(
				'http://localhost:8080/auth/login',
				expect.objectContaining({
					method: 'POST',
					body: JSON.stringify(requestData),
				})
			)
			expect(result).toEqual(mockResponse)
		})

		it('makes POST request without body when data is undefined', async () => {
			globalThis.fetch = vi.fn().mockResolvedValue({
				ok: true,
				status: 200,
				headers: new Headers({ 'content-type': 'application/json' }),
				json: async () => ({}),
			})

			await api.post('/test')

			expect(globalThis.fetch).toHaveBeenCalledWith(
				expect.any(String),
				expect.objectContaining({
					method: 'POST',
				})
			)
		})
	})

	describe('put', () => {
		it('makes PUT request with JSON body', async () => {
			const requestData = { name: 'Updated Name' }
			const mockResponse = { id: '123', name: 'Updated Name' }

			globalThis.fetch = vi.fn().mockResolvedValue({
				ok: true,
				status: 200,
				headers: new Headers({ 'content-type': 'application/json' }),
				json: async () => mockResponse,
			})

			const result = await api.put('/users/123', requestData)

			expect(globalThis.fetch).toHaveBeenCalledWith(
				'http://localhost:8080/users/123',
				expect.objectContaining({
					method: 'PUT',
					body: JSON.stringify(requestData),
				})
			)
			expect(result).toEqual(mockResponse)
		})
	})

	describe('patch', () => {
		it('makes PATCH request with JSON body', async () => {
			const requestData = { active: true }
			const mockResponse = { id: '123', active: true }

			globalThis.fetch = vi.fn().mockResolvedValue({
				ok: true,
				status: 200,
				headers: new Headers({ 'content-type': 'application/json' }),
				json: async () => mockResponse,
			})

			const result = await api.patch('/users/123', requestData)

			expect(globalThis.fetch).toHaveBeenCalledWith(
				'http://localhost:8080/users/123',
				expect.objectContaining({
					method: 'PATCH',
					body: JSON.stringify(requestData),
				})
			)
			expect(result).toEqual(mockResponse)
		})
	})

	describe('delete', () => {
		it('makes DELETE request', async () => {
			globalThis.fetch = vi.fn().mockResolvedValue({
				ok: true,
				status: 200,
				headers: new Headers({ 'content-type': 'application/json' }),
				json: async () => ({}),
			})

			await api.delete('/users/123')

			expect(globalThis.fetch).toHaveBeenCalledWith(
				'http://localhost:8080/users/123',
				expect.objectContaining({
					method: 'DELETE',
				})
			)
		})
	})

	describe('error handling', () => {
		it('throws error with message from JSON error response', async () => {
			const errorResponse = { error: 'Invalid credentials' }

			globalThis.fetch = vi.fn().mockResolvedValue({
				ok: false,
				status: 401,
				headers: new Headers({ 'content-type': 'application/json' }),
				json: async () => errorResponse,
			})

			await expect(api.get('/test')).rejects.toThrow('Invalid credentials')
		})

		it('throws HTTP status error when JSON error response has no error field', async () => {
			globalThis.fetch = vi.fn().mockResolvedValue({
				ok: false,
				status: 500,
				headers: new Headers({ 'content-type': 'application/json' }),
				json: async () => ({}),
			})

			await expect(api.get('/test')).rejects.toThrow('HTTP 500')
		})

		it('throws formatted error for HTML 404 response', async () => {
			globalThis.fetch = vi.fn().mockResolvedValue({
				ok: false,
				status: 404,
				headers: new Headers({ 'content-type': 'text/html' }),
				clone: () => ({
					text: async () => '<!DOCTYPE html><html>Not Found</html>',
				}),
			})

			await expect(api.get('/test')).rejects.toThrow('API endpoint not found')
		})

		it('throws formatted error for HTML 502 response', async () => {
			globalThis.fetch = vi.fn().mockResolvedValue({
				ok: false,
				status: 502,
				headers: new Headers({ 'content-type': 'text/html' }),
				clone: () => ({
					text: async () => '<!DOCTYPE html><html>Bad Gateway</html>',
				}),
			})

			await expect(api.get('/test')).rejects.toThrow('API server is not available')
		})

		it('throws formatted error for HTML 503 response', async () => {
			globalThis.fetch = vi.fn().mockResolvedValue({
				ok: false,
				status: 503,
				headers: new Headers({ 'content-type': 'text/html' }),
				clone: () => ({
					text: async () => '<!DOCTYPE html><html>Service Unavailable</html>',
				}),
			})

			await expect(api.get('/test')).rejects.toThrow('API server is not available')
		})

		it('throws error with text content for non-HTML non-JSON error response', async () => {
			globalThis.fetch = vi.fn().mockResolvedValue({
				ok: false,
				status: 500,
				headers: new Headers({ 'content-type': 'text/plain' }),
				clone: () => ({
					text: async () => 'Internal Server Error',
				}),
			})

			await expect(api.get('/test')).rejects.toThrow('Internal Server Error')
		})

		it('throws timeout error when request is aborted', async () => {
			globalThis.fetch = vi.fn().mockImplementation(() =>
				new Promise((_, reject) => {
					setTimeout(() => reject(new DOMException('Aborted', 'AbortError')), 100)
				})
			)

			await expect(api.get('/test')).rejects.toThrow('Request timeout')
		})

		it('throws connection error on network failure', async () => {
			globalThis.fetch = vi.fn().mockRejectedValue(new Error('Failed to fetch'))

			await expect(api.get('/test')).rejects.toThrow('Unable to connect to server')
		})

		it('throws connection error on NetworkError', async () => {
			globalThis.fetch = vi.fn().mockRejectedValue(new Error('NetworkError'))

			await expect(api.get('/test')).rejects.toThrow('Unable to connect to server')
		})

		it('re-throws formatted error messages', async () => {
			globalThis.fetch = vi.fn().mockRejectedValue(new Error('Custom error message'))

			await expect(api.get('/test')).rejects.toThrow('Custom error message')
		})

		it('throws generic error for unexpected error types', async () => {
			globalThis.fetch = vi.fn().mockRejectedValue('string error')

			await expect(api.get('/test')).rejects.toThrow('An unexpected error occurred')
		})
	})

	describe('204 No Content', () => {
		it('returns undefined for 204 responses', async () => {
			globalThis.fetch = vi.fn().mockResolvedValue({
				ok: true,
				status: 204,
				headers: new Headers(),
			})

			const result = await api.delete('/test')

			expect(result).toBeUndefined()
		})
	})

	describe('response parsing', () => {
		it('parses JSON response with application/json content-type', async () => {
			const mockResponse = { data: 'test' }

			globalThis.fetch = vi.fn().mockResolvedValue({
				ok: true,
				status: 200,
				headers: new Headers({ 'content-type': 'application/json' }),
				json: async () => mockResponse,
			})

			const result = await api.get('/test')

			expect(result).toEqual(mockResponse)
		})

		it('parses JSON response without content-type header', async () => {
			const mockResponse = { data: 'test' }

			globalThis.fetch = vi.fn().mockResolvedValue({
				ok: true,
				status: 200,
				headers: new Headers(),
				text: async () => JSON.stringify(mockResponse),
			})

			const result = await api.get('/test')

			expect(result).toEqual(mockResponse)
		})

		it('returns undefined for empty text response', async () => {
			globalThis.fetch = vi.fn().mockResolvedValue({
				ok: true,
				status: 200,
				headers: new Headers({ 'content-type': 'text/plain' }),
				text: async () => '',
			})

			const result = await api.get('/test')

			expect(result).toBeUndefined()
		})

		it('throws error when non-JSON response cannot be parsed as JSON', async () => {
			globalThis.fetch = vi.fn().mockResolvedValue({
				ok: true,
				status: 200,
				headers: new Headers({ 'content-type': 'text/plain' }),
				text: async () => 'not json',
			})

			await expect(api.get('/test')).rejects.toThrow('Expected JSON response')
		})
	})

	describe('request configuration', () => {
		it('sets CORS mode and credentials', async () => {
			globalThis.fetch = vi.fn().mockResolvedValue({
				ok: true,
				status: 200,
				headers: new Headers({ 'content-type': 'application/json' }),
				json: async () => ({}),
			})

			await api.get('/test')

			expect(globalThis.fetch).toHaveBeenCalledWith(
				expect.any(String),
				expect.objectContaining({
					mode: 'cors',
					credentials: 'omit',
				})
			)
		})

		it('sets AbortSignal for timeout', async () => {
			globalThis.fetch = vi.fn().mockResolvedValue({
				ok: true,
				status: 200,
				headers: new Headers({ 'content-type': 'application/json' }),
				json: async () => ({}),
			})

			await api.get('/test')

			const fetchCalls = (globalThis.fetch as any).mock.calls
			const [, options] = fetchCalls[0]

			expect(options.signal).toBeInstanceOf(AbortSignal)
		})

		it('merges custom headers from options', async () => {
			globalThis.fetch = vi.fn().mockResolvedValue({
				ok: true,
				status: 200,
				headers: new Headers({ 'content-type': 'application/json' }),
				json: async () => ({}),
			})

			await api.get('/test')

			const fetchCalls = (globalThis.fetch as any).mock.calls
			const [, options] = fetchCalls[0]

			expect(options.headers['Content-Type']).toBe('application/json')
			expect(options.headers['X-Synthetic-Mode']).toBe('false')
		})
	})
})
