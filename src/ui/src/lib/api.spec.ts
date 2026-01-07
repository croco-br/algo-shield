import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { api } from './api'

// Mock uiConfig
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
    // Clear localStorage
    localStorage.clear()
    
    // Reset fetch mock
    globalThis.fetch = vi.fn()
  })

  afterEach(() => {
    vi.clearAllMocks()
  })

  describe('get', () => {
    it('makes GET request with correct headers', async () => {
      // Arrange
      const mockResponse = { data: 'test' }
      globalThis.fetch = vi.fn().mockResolvedValue({
        ok: true,
        status: 200,
        headers: new Headers({ 'content-type': 'application/json' }),
        json: async () => mockResponse,
      })

      // Act
      const result = await api.get('/test')

      // Assert
      expect(globalThis.fetch).toHaveBeenCalledWith(
        'http://localhost:8080/test',
        expect.objectContaining({
          method: 'GET',
          headers: expect.objectContaining({
            'Content-Type': 'application/json',
          }),
        })
      )
      expect(result).toEqual(mockResponse)
    })

    it('includes Authorization header when token exists', async () => {
      // Arrange
      const token = 'test-token'
      localStorage.setItem('auth_token', token)

      const mockResponse = { data: 'test' }
      globalThis.fetch = vi.fn().mockResolvedValue({
        ok: true,
        status: 200,
        headers: new Headers({ 'content-type': 'application/json' }),
        json: async () => mockResponse,
      })

      // Act
      await api.get('/test')

      // Assert
      expect(globalThis.fetch).toHaveBeenCalledWith(
        'http://localhost:8080/test',
        expect.objectContaining({
          headers: expect.objectContaining({
            'Authorization': `Bearer ${token}`,
          }),
        })
      )
    })
  })

  describe('post', () => {
    it('makes POST request with body', async () => {
      // Arrange
      const requestData = { email: 'test@example.com', password: 'password123' }
      const mockResponse = { token: 'jwt-token' }
      
      globalThis.fetch = vi.fn().mockResolvedValue({
        ok: true,
        status: 200,
        headers: new Headers({ 'content-type': 'application/json' }),
        json: async () => mockResponse,
      })

      // Act
      const result = await api.post('/auth/login', requestData)

      // Assert
      expect(globalThis.fetch).toHaveBeenCalledWith(
        'http://localhost:8080/auth/login',
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify(requestData),
        })
      )
      expect(result).toEqual(mockResponse)
    })
  })

  describe('error handling', () => {
    it('throws error with message from JSON response', async () => {
      // Arrange
      const errorResponse = { error: 'Invalid credentials' }
      
      globalThis.fetch = vi.fn().mockResolvedValue({
        ok: false,
        status: 401,
        headers: new Headers({ 'content-type': 'application/json' }),
        json: async () => errorResponse,
      })

      // Act & Assert
      await expect(api.get('/test')).rejects.toThrow('Invalid credentials')
    })

    it('throws generic error for non-JSON error responses', async () => {
      // Arrange
      globalThis.fetch = vi.fn().mockResolvedValue({
        ok: false,
        status: 500,
        headers: new Headers({ 'content-type': 'text/html' }),
        clone: () => ({
          text: async () => '<html>Server Error</html>',
        }),
      })

      // Act & Assert
      await expect(api.get('/test')).rejects.toThrow()
    })

    it('throws timeout error when request takes too long', async () => {
      // Arrange
      globalThis.fetch = vi.fn().mockImplementation(() => 
        new Promise((_, reject) => {
          setTimeout(() => reject(new DOMException('Aborted', 'AbortError')), 100)
        })
      )

      // Act & Assert
      await expect(api.get('/test')).rejects.toThrow('Request timeout')
    })

    it('throws connection error on network failure', async () => {
      // Arrange
      globalThis.fetch = vi.fn().mockRejectedValue(new Error('Failed to fetch'))

      // Act & Assert
      await expect(api.get('/test')).rejects.toThrow('Unable to connect to server')
    })
  })

  describe('204 No Content', () => {
    it('returns undefined for 204 responses', async () => {
      // Arrange
      globalThis.fetch = vi.fn().mockResolvedValue({
        ok: true,
        status: 204,
        headers: new Headers(),
      })

      // Act
      const result = await api.delete('/test')

      // Assert
      expect(result).toBeUndefined()
    })
  })

  describe('CORS handling', () => {
    it('sets correct CORS mode and credentials', async () => {
      // Arrange
      globalThis.fetch = vi.fn().mockResolvedValue({
        ok: true,
        status: 200,
        headers: new Headers({ 'content-type': 'application/json' }),
        json: async () => ({}),
      })

      // Act
      await api.get('/test')

      // Assert
      expect(globalThis.fetch).toHaveBeenCalledWith(
        'http://localhost:8080/test',
        expect.objectContaining({
          mode: 'cors',
          credentials: 'omit',
        })
      )
    })
  })
})

