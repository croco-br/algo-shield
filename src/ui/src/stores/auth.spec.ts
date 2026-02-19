import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore, setRouterNavigate } from './auth'
import type { User } from './auth'
import { setTokenGetter } from '@/lib/api'

// Mock the api module
vi.mock('@/lib/api', () => ({
	api: {
		get: vi.fn(),
		post: vi.fn(),
		put: vi.fn(),
		delete: vi.fn(),
	},
	setTokenGetter: vi.fn(),
	setCsrfTokenGetter: vi.fn(),
	setRefreshTokenFn: vi.fn(),
	setForceLogoutFn: vi.fn(),
	setSyntheticModeStorage: vi.fn(),
	resetLogoutState: vi.fn(),
}))

describe('useAuthStore', () => {
	beforeEach(() => {
		// Create a fresh pinia instance for each test
		setActivePinia(createPinia())

		// Clear all mocks
		vi.clearAllMocks()
	})

	afterEach(() => {
		vi.clearAllMocks()
	})

	describe('loadUserFromToken', () => {
		it('sets user to null when no token in memory', async () => {
			// Arrange
			const store = useAuthStore()

			// Act
			await store.refresh()

			// Assert
			expect(store.user).toBeNull()
			expect(store.loading).toBe(false)
		})

		it('loads user data when valid token exists in memory', async () => {
			// Arrange
			const { api } = await import('@/lib/api')
			const mockUser: User = {
				id: 'user-123',
				email: 'test@example.com',
				name: 'Test User',
				auth_type: 'local',
				active: true,
				roles: [{ id: 'role-1', name: 'admin', description: 'Administrator' }],
			}

			vi.mocked(api.get).mockResolvedValue(mockUser)

			const store = useAuthStore()
			// Manually set token in memory (simulating login)
			await store.setToken('valid-token')

			// Act
			await store.refresh()

			// Assert
			expect(api.get).toHaveBeenCalledWith('/api/v1/auth/me')
			expect(store.user).toEqual(mockUser)
			expect(store.loading).toBe(false)
		})

		it('clears token when API returns error', async () => {
			// Arrange
			const { api } = await import('@/lib/api')

			vi.mocked(api.get).mockRejectedValue(new Error('Unauthorized'))

			const store = useAuthStore()
			// Set token first
			await store.setToken('invalid-token')

			// Clear the mock from setToken call
			vi.clearAllMocks()
			vi.mocked(api.get).mockRejectedValue(new Error('Unauthorized'))

			// Act
			await store.refresh()

			// Assert
			expect(store.getToken()).toBeNull()
			expect(store.user).toBeNull()
			expect(store.loading).toBe(false)
		})
	})

	describe('setToken', () => {
		it('stores token in memory and loads user data', async () => {
			// Arrange
			const { api } = await import('@/lib/api')
			const mockUser: User = {
				id: 'user-123',
				email: 'test@example.com',
				name: 'Test User',
				auth_type: 'local',
				active: true,
			}

			vi.mocked(api.get).mockResolvedValue(mockUser)

			const store = useAuthStore()
			const token = 'new-token'

			// Act
			await store.setToken(token)

			// Assert
			// SECURITY: Token is in memory, NOT in localStorage
			expect(store.getToken()).toBe(token)
			expect(store.user).toEqual(mockUser)
		})
	})

	describe('logout', () => {
		it('calls logout API and clears user data from memory', async () => {
			// Arrange
			const { api } = await import('@/lib/api')
			const mockUser: User = {
				id: 'user-123',
				email: 'test@example.com',
				name: 'Test User',
				auth_type: 'local',
				active: true,
			}

			vi.mocked(api.post).mockResolvedValue({ message: 'Logged out successfully' })
			vi.mocked(api.get).mockResolvedValue(mockUser)

			const store = useAuthStore()
			await store.setToken('some-token')

			// Clear mocks from setup
			vi.clearAllMocks()
			vi.mocked(api.post).mockResolvedValue({ message: 'Logged out successfully' })

			// Act
			await store.logout()

			// Assert
			expect(api.post).toHaveBeenCalledWith('/api/v1/auth/logout', {
				refresh_token: undefined,
			})
			expect(store.user).toBeNull()
			expect(store.getToken()).toBeNull()
		})

		it('clears user data even when logout API fails', async () => {
			// Arrange
			const { api } = await import('@/lib/api')
			const mockUser: User = {
				id: 'user-123',
				email: 'test@example.com',
				name: 'Test User',
				auth_type: 'local',
				active: true,
			}

			vi.mocked(api.get).mockResolvedValue(mockUser)

			const store = useAuthStore()
			await store.setToken('some-token')

			// Clear mocks and setup error
			vi.clearAllMocks()
			vi.mocked(api.post).mockRejectedValue(new Error('Network error'))

			// Act
			await store.logout()

			// Assert
			expect(store.user).toBeNull()
			expect(store.getToken()).toBeNull()
		})
	})

	describe('isAdmin', () => {
		it('returns true when user has admin role', () => {
			// Arrange
			const store = useAuthStore()
			store.user = {
				id: 'user-123',
				email: 'admin@example.com',
				name: 'Admin User',
				auth_type: 'local',
				active: true,
				roles: [
					{ id: 'role-1', name: 'admin', description: 'Administrator' },
					{ id: 'role-2', name: 'user', description: 'Regular user' },
				],
			}

			// Act & Assert
			expect(store.isAdmin).toBe(true)
		})

		it('returns false when user has no admin role', () => {
			// Arrange
			const store = useAuthStore()
			store.user = {
				id: 'user-123',
				email: 'user@example.com',
				name: 'Regular User',
				auth_type: 'local',
				active: true,
				roles: [{ id: 'role-2', name: 'user', description: 'Regular user' }],
			}

			// Act & Assert
			expect(store.isAdmin).toBe(false)
		})

		it('returns false when user has no roles', () => {
			// Arrange
			const store = useAuthStore()
			store.user = {
				id: 'user-123',
				email: 'user@example.com',
				name: 'Regular User',
				auth_type: 'local',
				active: true,
			}

			// Act & Assert
			expect(store.isAdmin).toBe(false)
		})

		it('returns false when user is null', () => {
			// Arrange
			const store = useAuthStore()
			store.user = null

			// Act & Assert
			expect(store.isAdmin).toBe(false)
		})
	})

	describe('refresh', () => {
		it('reloads user data from memory token', async () => {
			// Arrange
			const { api } = await import('@/lib/api')
			const mockUser: User = {
				id: 'user-123',
				email: 'test@example.com',
				name: 'Test User',
				auth_type: 'local',
				active: true,
			}

			vi.mocked(api.get).mockResolvedValue(mockUser)

			const store = useAuthStore()
			await store.setToken('valid-token')

			// Clear mocks from setup
			vi.clearAllMocks()
			vi.mocked(api.get).mockResolvedValue(mockUser)

			// Act
			await store.refresh()

			// Assert
			expect(api.get).toHaveBeenCalledWith('/api/v1/auth/me')
			expect(store.user).toEqual(mockUser)
		})
	})

	describe('getToken', () => {
		it('returns null when no token is set', () => {
			// Arrange
			const store = useAuthStore()

			// Act & Assert
			expect(store.getToken()).toBeNull()
		})

		it('returns token after setToken is called', async () => {
			// Arrange
			const { api } = await import('@/lib/api')
			const mockUser: User = {
				id: 'user-123',
				email: 'test@example.com',
				name: 'Test User',
				auth_type: 'local',
				active: true,
			}

			vi.mocked(api.get).mockResolvedValue(mockUser)

			const store = useAuthStore()
			const token = 'test-token'

			// Act
			await store.setToken(token)

			// Assert
			expect(store.getToken()).toBe(token)
		})
	})

	describe('security: memory-only token storage', () => {
		it('token is lost on page reload (no localStorage)', () => {
			// This test documents the security behavior
			// SECURITY: Tokens are stored in memory only, not localStorage
			// This prevents XSS attacks from stealing tokens
			// Trade-off: Users must re-authenticate after page reload
			const store = useAuthStore()

			// Assert: No token in localStorage
			expect(localStorage.getItem('auth_token')).toBeNull()

			// Assert: Token only accessible via getToken()
			expect(store.getToken()).toBeNull()
		})
	})

	describe('whenReady', () => {
		it('resolves after loading completes', async () => {
			const { api } = await import('@/lib/api')
			vi.mocked(api.get).mockResolvedValue(null)

			const store = useAuthStore()

			// whenReady should resolve once loading finishes
			await store.whenReady()

			expect(store.loading).toBe(false)
		})

		it('resolves immediately if already loaded', async () => {
			const store = useAuthStore()

			// Wait for initial load
			await store.whenReady()

			// Calling again should resolve immediately
			const resolved = await Promise.race([
				store.whenReady().then(() => true),
				new Promise<boolean>(resolve => setTimeout(() => resolve(false), 50)),
			])
			expect(resolved).toBe(true)
		})
	})

	describe('CSRF token preservation', () => {
		it('preserves CSRF token when refresh response omits csrf_token', async () => {
			const { api } = await import('@/lib/api')
			const mockUser: User = {
				id: 'user-123',
				email: 'test@example.com',
				name: 'Test User',
				auth_type: 'local',
				active: true,
			}

			vi.mocked(api.get).mockResolvedValue(mockUser)

			const store = useAuthStore()
			// Set initial tokens including CSRF
			await store.setToken('token-1', 'csrf-original', 'refresh-1')

			// Clear mocks and simulate refresh response without csrf_token
			vi.clearAllMocks()
			vi.mocked(api.post).mockResolvedValue({
				token: 'token-2',
				// csrf_token intentionally omitted
			})
			vi.mocked(api.get).mockResolvedValue(mockUser)

			await store.refresh()

			// CSRF token should be preserved, not overwritten with null
			expect(store.getCsrfToken()).toBe('csrf-original')
		})

		it('updates CSRF token when refresh response includes csrf_token', async () => {
			const { api } = await import('@/lib/api')
			const mockUser: User = {
				id: 'user-123',
				email: 'test@example.com',
				name: 'Test User',
				auth_type: 'local',
				active: true,
			}

			vi.mocked(api.get).mockResolvedValue(mockUser)

			const store = useAuthStore()
			await store.setToken('token-1', 'csrf-old', 'refresh-1')

			vi.clearAllMocks()
			// First api.get call fails (token expired) → triggers refreshAccessToken in catch
			// Then api.post (refresh) returns new csrf_token
			// Then second api.get (from refreshAccessToken) returns user
			vi.mocked(api.get)
				.mockRejectedValueOnce(new Error('Unauthorized'))
				.mockResolvedValueOnce(mockUser)
			vi.mocked(api.post).mockResolvedValue({
				token: 'token-2',
				csrf_token: 'csrf-new',
			})

			await store.refresh()

			expect(store.getCsrfToken()).toBe('csrf-new')
		})
	})

	describe('force logout with router navigation', () => {
		it('uses injected router navigate when available', async () => {
			const mockNavigate = vi.fn().mockResolvedValue(undefined)
			setRouterNavigate(mockNavigate)

			const { setForceLogoutFn } = await import('@/lib/api')

			// Create store — this registers the force logout callback
			useAuthStore()

			// Get the force logout function that was registered
			const calls = vi.mocked(setForceLogoutFn).mock.calls
			const forceLogoutFn = calls[calls.length - 1]?.[0]
			expect(forceLogoutFn).toBeDefined()

			// Call force logout
			forceLogoutFn!()

			expect(mockNavigate).toHaveBeenCalledWith('/login')
		})
	})
})
