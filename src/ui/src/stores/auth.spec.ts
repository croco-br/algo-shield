import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from './auth'
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
	setSyntheticModeStorage: vi.fn(),
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
			expect(api.post).toHaveBeenCalledWith('/api/v1/auth/logout')
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
})
