import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from './auth'
import type { User } from './auth'

// Mock the api module
vi.mock('@/lib/api', () => ({
  api: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  },
}))

describe('useAuthStore', () => {
  beforeEach(() => {
    // Create a fresh pinia instance for each test
    setActivePinia(createPinia())
    
    // Clear localStorage
    localStorage.clear()
    
    // Clear all mocks
    vi.clearAllMocks()
  })

  afterEach(() => {
    vi.clearAllMocks()
  })

  describe('loadUserFromToken', () => {
    it('sets user to null when no token in localStorage', async () => {
      // Arrange
      const store = useAuthStore()

      // Act
      await store.refresh()

      // Assert
      expect(store.user).toBeNull()
      expect(store.loading).toBe(false)
    })

    it('loads user data when valid token exists', async () => {
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

      localStorage.setItem('auth_token', 'valid-token')
      vi.mocked(api.get).mockResolvedValue(mockUser)

      const store = useAuthStore()

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
      
      localStorage.setItem('auth_token', 'invalid-token')
      vi.mocked(api.get).mockRejectedValue(new Error('Unauthorized'))

      const store = useAuthStore()

      // Act
      await store.refresh()

      // Assert
      expect(localStorage.getItem('auth_token')).toBeNull()
      expect(store.user).toBeNull()
      expect(store.loading).toBe(false)
    })
  })

  describe('setToken', () => {
    it('stores token and loads user data', async () => {
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
      expect(localStorage.getItem('auth_token')).toBe(token)
      expect(store.user).toEqual(mockUser)
    })
  })

  describe('logout', () => {
    it('calls logout API and clears user data', async () => {
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

      const store = useAuthStore()
      store.user = mockUser
      localStorage.setItem('auth_token', 'some-token')

      // Act
      await store.logout()

      // Assert
      expect(api.post).toHaveBeenCalledWith('/api/v1/auth/logout')
      expect(store.user).toBeNull()
      expect(localStorage.getItem('auth_token')).toBeNull()
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

      vi.mocked(api.post).mockRejectedValue(new Error('Network error'))

      const store = useAuthStore()
      store.user = mockUser
      localStorage.setItem('auth_token', 'some-token')

      // Act
      await store.logout()

      // Assert
      expect(store.user).toBeNull()
      expect(localStorage.getItem('auth_token')).toBeNull()
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
    it('reloads user data from token', async () => {
      // Arrange
      const { api } = await import('@/lib/api')
      const mockUser: User = {
        id: 'user-123',
        email: 'test@example.com',
        name: 'Test User',
        auth_type: 'local',
        active: true,
      }

      localStorage.setItem('auth_token', 'valid-token')
      vi.mocked(api.get).mockResolvedValue(mockUser)

      const store = useAuthStore()

      // Act
      await store.refresh()

      // Assert
      expect(api.get).toHaveBeenCalledWith('/api/v1/auth/me')
      expect(store.user).toEqual(mockUser)
    })
  })
})

