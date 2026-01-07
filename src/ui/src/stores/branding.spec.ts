import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useBrandingStore } from './branding'
import type { BrandingConfig } from './branding'

vi.mock('@/lib/api', () => ({
  api: {
    get: vi.fn(),
    put: vi.fn(),
  },
}))

vi.mock('@/plugins/vuetify', () => ({
  updateVuetifyTheme: vi.fn(),
}))

describe('useBrandingStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
    vi.clearAllMocks()
    
    document.title = ''
    document.documentElement.style.cssText = ''
    document.querySelectorAll("link[rel*='icon']").forEach(link => link.remove())
  })

  afterEach(() => {
    vi.clearAllMocks()
  })

  describe('loadBranding', () => {
    it('loads and applies branding config from API', async () => {
      const { api } = await import('@/lib/api')
      const mockConfig: BrandingConfig = {
        id: 1,
        app_name: 'Test App',
        icon_url: '/test-icon.png',
        favicon_url: '/test-favicon.ico',
        primary_color: '#FF0000',
        secondary_color: '#00FF00',
        header_color: '#0000FF',
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      }

      vi.mocked(api.get).mockResolvedValue(mockConfig)

      const store = useBrandingStore()
      await store.loadBranding()

      expect(api.get).toHaveBeenCalledWith('/api/v1/branding')
      expect(store.config).toEqual(mockConfig)
      expect(store.loading).toBe(false)
      expect(store.error).toBeNull()
      expect(document.title).toBe('Test App')
    })

    it('applies default branding when API fails', async () => {
      const { api } = await import('@/lib/api')
      
      vi.mocked(api.get).mockRejectedValue(new Error('Network error'))

      const store = useBrandingStore()
      await store.loadBranding()

      expect(store.error).toBeTruthy()
      expect(store.loading).toBe(false)
      expect(document.title).toBe('AlgoShield')
    })

    it('handles non-Error exceptions', async () => {
      const { api } = await import('@/lib/api')
      
      vi.mocked(api.get).mockRejectedValue('String error')

      const store = useBrandingStore()
      await store.loadBranding()

      expect(store.error).toBe('Failed to load branding')
      expect(store.loading).toBe(false)
    })
  })

  describe('updateBranding', () => {
    it('updates branding config via API', async () => {
      const { api } = await import('@/lib/api')
      const updateData = {
        app_name: 'Updated App',
        icon_url: '/new-icon.png',
        favicon_url: '/new-favicon.ico',
        primary_color: '#111111',
        secondary_color: '#222222',
        header_color: '#333333',
      }

      const mockResponse: BrandingConfig = {
        id: 1,
        ...updateData,
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-02T00:00:00Z',
      }

      vi.mocked(api.get).mockResolvedValue(mockResponse)
      vi.mocked(api.put).mockResolvedValue(mockResponse)

      const store = useBrandingStore()
      await store.loadBranding()
      
      vi.clearAllMocks()
      vi.mocked(api.put).mockResolvedValue(mockResponse)

      const result = await store.updateBranding(updateData)

      expect(api.put).toHaveBeenCalledWith('/api/v1/branding', updateData)
      expect(store.config).toEqual(mockResponse)
      expect(result).toEqual(mockResponse)
      expect(store.loading).toBe(false)
      expect(store.error).toBeNull()
      expect(document.title).toBe('Updated App')
    })

    it('sets error when update fails', async () => {
      const { api } = await import('@/lib/api')
      const updateData = {
        app_name: 'Updated App',
        primary_color: '#111111',
        secondary_color: '#222222',
        header_color: '#333333',
      }

      vi.mocked(api.put).mockRejectedValue(new Error('Update failed'))

      const store = useBrandingStore()

      await expect(store.updateBranding(updateData)).rejects.toThrow('Update failed')
      expect(store.error).toBeTruthy()
      expect(store.loading).toBe(false)
    })

    it('handles non-Error exceptions during update', async () => {
      const { api } = await import('@/lib/api')
      const updateData = {
        app_name: 'Updated App',
        primary_color: '#111111',
        secondary_color: '#222222',
        header_color: '#333333',
      }

      vi.mocked(api.put).mockRejectedValue('String error')

      const store = useBrandingStore()

      await expect(store.updateBranding(updateData)).rejects.toEqual('String error')
      expect(store.error).toBe('Failed to update branding')
      expect(store.loading).toBe(false)
    })
  })

  describe('applyBranding', () => {
    it('updates document title', async () => {
      const { api } = await import('@/lib/api')
      const mockConfig: BrandingConfig = {
        id: 1,
        app_name: 'Custom Title',
        primary_color: '#FF0000',
        secondary_color: '#00FF00',
        header_color: '#0000FF',
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      }

      vi.mocked(api.get).mockResolvedValue(mockConfig)

      const store = useBrandingStore()
      await store.loadBranding()

      expect(document.title).toBe('Custom Title')
    })

    it('updates CSS custom properties', async () => {
      const { api } = await import('@/lib/api')
      const mockConfig: BrandingConfig = {
        id: 1,
        app_name: 'Test App',
        primary_color: '#FF0000',
        secondary_color: '#00FF00',
        header_color: '#0000FF',
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      }

      vi.mocked(api.get).mockResolvedValue(mockConfig)

      const store = useBrandingStore()
      await store.loadBranding()

      const root = document.documentElement
      expect(root.style.getPropertyValue('--color-primary')).toBe('#FF0000')
      expect(root.style.getPropertyValue('--color-secondary')).toBe('#00FF00')
      expect(root.style.getPropertyValue('--color-header-background')).toBe('#0000FF')
    })

    it('updates favicon when provided', async () => {
      const { api } = await import('@/lib/api')
      const mockConfig: BrandingConfig = {
        id: 1,
        app_name: 'Test App',
        favicon_url: '/custom-favicon.ico',
        primary_color: '#FF0000',
        secondary_color: '#00FF00',
        header_color: '#0000FF',
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      }

      vi.mocked(api.get).mockResolvedValue(mockConfig)

      const store = useBrandingStore()
      await store.loadBranding()

      const faviconLink = document.querySelector("link[rel='icon']") as HTMLLinkElement
      expect(faviconLink).toBeTruthy()
      expect(faviconLink?.href).toContain('/custom-favicon.ico')
    })

    it('does not add favicon when not provided', async () => {
      const { api } = await import('@/lib/api')
      const mockConfig: BrandingConfig = {
        id: 1,
        app_name: 'Test App',
        favicon_url: null,
        primary_color: '#FF0000',
        secondary_color: '#00FF00',
        header_color: '#0000FF',
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      }

      vi.mocked(api.get).mockResolvedValue(mockConfig)

      const store = useBrandingStore()
      await store.loadBranding()

      const faviconLink = document.querySelector("link[rel='icon']")
      expect(faviconLink).toBeNull()
    })

    it('removes existing favicon before adding new one', async () => {
      const existingFavicon = document.createElement('link')
      existingFavicon.rel = 'icon'
      existingFavicon.href = '/old-favicon.ico'
      document.head.appendChild(existingFavicon)

      const { api } = await import('@/lib/api')
      const mockConfig: BrandingConfig = {
        id: 1,
        app_name: 'Test App',
        favicon_url: '/new-favicon.ico',
        primary_color: '#FF0000',
        secondary_color: '#00FF00',
        header_color: '#0000FF',
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      }

      vi.mocked(api.get).mockResolvedValue(mockConfig)

      const store = useBrandingStore()
      await store.loadBranding()

      const faviconLinks = document.querySelectorAll("link[rel='icon']")
      expect(faviconLinks.length).toBe(1)
      expect((faviconLinks[0] as HTMLLinkElement).href).toContain('/new-favicon.ico')
    })
  })

  describe('applyDefaultBranding', () => {
    it('applies default values on initialization', () => {
      const store = useBrandingStore()

      expect(document.title).toBe('AlgoShield')
      
      const root = document.documentElement
      expect(root.style.getPropertyValue('--color-primary')).toBe('#3B82F6')
      expect(root.style.getPropertyValue('--color-secondary')).toBe('#10B981')
      expect(root.style.getPropertyValue('--color-header-background')).toBe('#1e1e1e')
    })
  })

  describe('vuetify theme integration', () => {
    it('calls updateVuetifyTheme with correct colors', async () => {
      const { api } = await import('@/lib/api')
      const { updateVuetifyTheme } = await import('@/plugins/vuetify')
      
      const mockConfig: BrandingConfig = {
        id: 1,
        app_name: 'Test App',
        primary_color: '#FF0000',
        secondary_color: '#00FF00',
        header_color: '#0000FF',
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      }

      vi.mocked(api.get).mockResolvedValue(mockConfig)

      const store = useBrandingStore()
      await store.loadBranding()

      expect(updateVuetifyTheme).toHaveBeenCalledWith('#FF0000', '#00FF00')
    })

    it('calls updateVuetifyTheme with default colors on error', async () => {
      const { api } = await import('@/lib/api')
      const { updateVuetifyTheme } = await import('@/plugins/vuetify')
      
      vi.mocked(api.get).mockRejectedValue(new Error('API Error'))

      const store = useBrandingStore()
      await store.loadBranding()

      expect(updateVuetifyTheme).toHaveBeenCalledWith('#3B82F6', '#10B981')
    })
  })
})

