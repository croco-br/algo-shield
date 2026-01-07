import { describe, it, expect, beforeEach, vi } from 'vitest'

describe('config', () => {
  beforeEach(() => {
    vi.resetModules()
  })

  describe('default configuration', () => {
    it('uses default values when environment variables are not set', async () => {
      vi.stubEnv('VITE_API_URL', '')
      vi.stubEnv('VITE_API_TIMEOUT', '')
      
      const { uiConfig } = await import('./config')

      expect(uiConfig.api.baseUrl).toBe('')
      expect(uiConfig.api.timeout).toBe(30000)
      expect(uiConfig.api.retry.maxAttempts).toBe(3)
      expect(uiConfig.api.retry.initialDelay).toBe(1000)
      expect(uiConfig.api.retry.maxDelay).toBe(10000)
      expect(uiConfig.api.retry.multiplier).toBe(2.0)
      expect(uiConfig.ui.toast.duration).toBe(5000)
      expect(uiConfig.ui.polling.interval).toBe(10000)
    })
  })

  describe('environment variable overrides', () => {
    it('overrides API base URL from environment', async () => {
      vi.stubEnv('VITE_API_URL', 'http://custom-api.com')
      
      const { uiConfig } = await import('./config')

      expect(uiConfig.api.baseUrl).toBe('http://custom-api.com')
    })

    it('overrides API timeout from environment', async () => {
      vi.stubEnv('VITE_API_TIMEOUT', '60000')
      
      const { uiConfig } = await import('./config')

      expect(uiConfig.api.timeout).toBe(60000)
    })

    it('overrides retry max attempts from environment', async () => {
      vi.stubEnv('VITE_API_RETRY_MAX_ATTEMPTS', '5')
      
      const { uiConfig } = await import('./config')

      expect(uiConfig.api.retry.maxAttempts).toBe(5)
    })

    it('overrides retry initial delay from environment', async () => {
      vi.stubEnv('VITE_API_RETRY_INITIAL_DELAY', '2000')
      
      const { uiConfig } = await import('./config')

      expect(uiConfig.api.retry.initialDelay).toBe(2000)
    })

    it('overrides retry max delay from environment', async () => {
      vi.stubEnv('VITE_API_RETRY_MAX_DELAY', '20000')
      
      const { uiConfig } = await import('./config')

      expect(uiConfig.api.retry.maxDelay).toBe(20000)
    })

    it('overrides retry multiplier from environment', async () => {
      vi.stubEnv('VITE_API_RETRY_MULTIPLIER', '1.5')
      
      const { uiConfig } = await import('./config')

      expect(uiConfig.api.retry.multiplier).toBe(1.5)
    })

    it('overrides toast duration from environment', async () => {
      vi.stubEnv('VITE_UI_TOAST_DURATION', '3000')
      
      const { uiConfig } = await import('./config')

      expect(uiConfig.ui.toast.duration).toBe(3000)
    })

    it('overrides polling interval from environment', async () => {
      vi.stubEnv('VITE_UI_POLLING_INTERVAL', '5000')
      
      const { uiConfig } = await import('./config')

      expect(uiConfig.ui.polling.interval).toBe(5000)
    })
  })

  describe('invalid environment variables', () => {
    it('uses default when integer parse fails', async () => {
      vi.stubEnv('VITE_API_TIMEOUT', 'not-a-number')
      
      const { uiConfig } = await import('./config')

      expect(uiConfig.api.timeout).toBe(30000)
    })

    it('uses default when float parse fails', async () => {
      vi.stubEnv('VITE_API_RETRY_MULTIPLIER', 'invalid-float')
      
      const { uiConfig } = await import('./config')

      expect(uiConfig.api.retry.multiplier).toBe(2.0)
    })

    it('handles zero as valid integer value', async () => {
      vi.stubEnv('VITE_API_TIMEOUT', '0')
      
      const { uiConfig } = await import('./config')

      expect(uiConfig.api.timeout).toBe(0)
    })

    it('handles negative numbers correctly for integers', async () => {
      vi.stubEnv('VITE_API_TIMEOUT', '-1000')
      
      const { uiConfig } = await import('./config')

      expect(uiConfig.api.timeout).toBe(-1000)
    })

    it('handles negative numbers correctly for floats', async () => {
      vi.stubEnv('VITE_API_RETRY_MULTIPLIER', '-1.5')
      
      const { uiConfig } = await import('./config')

      expect(uiConfig.api.retry.multiplier).toBe(-1.5)
    })
  })

  describe('edge cases', () => {
    it('handles empty string environment variables', async () => {
      vi.stubEnv('VITE_API_URL', '')
      vi.stubEnv('VITE_API_TIMEOUT', '')
      
      const { uiConfig } = await import('./config')

      expect(uiConfig.api.baseUrl).toBe('')
      expect(uiConfig.api.timeout).toBe(30000)
    })

    it('handles whitespace in string values', async () => {
      vi.stubEnv('VITE_API_URL', '  http://api.com  ')
      
      const { uiConfig } = await import('./config')

      expect(uiConfig.api.baseUrl).toBe('  http://api.com  ')
    })

    it('handles whitespace in numeric values', async () => {
      vi.stubEnv('VITE_API_TIMEOUT', '  5000  ')
      
      const { uiConfig } = await import('./config')

      expect(uiConfig.api.timeout).toBe(5000)
    })

    it('handles decimal strings for integer values', async () => {
      vi.stubEnv('VITE_API_TIMEOUT', '5000.5')
      
      const { uiConfig } = await import('./config')

      expect(uiConfig.api.timeout).toBe(5000)
    })

    it('handles integer strings for float values', async () => {
      vi.stubEnv('VITE_API_RETRY_MULTIPLIER', '3')
      
      const { uiConfig } = await import('./config')

      expect(uiConfig.api.retry.multiplier).toBe(3.0)
    })
  })

  describe('configuration structure', () => {
    it('has all required API configuration fields', async () => {
      const { uiConfig } = await import('./config')

      expect(uiConfig.api).toBeDefined()
      expect(uiConfig.api.baseUrl).toBeDefined()
      expect(uiConfig.api.timeout).toBeDefined()
      expect(uiConfig.api.retry).toBeDefined()
      expect(uiConfig.api.retry.maxAttempts).toBeDefined()
      expect(uiConfig.api.retry.initialDelay).toBeDefined()
      expect(uiConfig.api.retry.maxDelay).toBeDefined()
      expect(uiConfig.api.retry.multiplier).toBeDefined()
    })

    it('has all required UI configuration fields', async () => {
      const { uiConfig } = await import('./config')

      expect(uiConfig.ui).toBeDefined()
      expect(uiConfig.ui.toast).toBeDefined()
      expect(uiConfig.ui.toast.duration).toBeDefined()
      expect(uiConfig.ui.polling).toBeDefined()
      expect(uiConfig.ui.polling.interval).toBeDefined()
    })
  })
})

