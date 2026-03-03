import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useSystemModeStore } from './systemMode'

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

describe('useSystemModeStore', () => {
	beforeEach(() => {
		setActivePinia(createPinia())
		vi.clearAllMocks()
		localStorage.clear()
	})

	describe('initial state', () => {
		it('initializes syntheticMode from localStorage', () => {
			localStorage.setItem('synthetic_mode', 'true')
			setActivePinia(createPinia())

			const store = useSystemModeStore()

			expect(store.syntheticMode).toBe(true)
		})

		it('defaults to false when localStorage has no value', () => {
			const store = useSystemModeStore()

			expect(store.syntheticMode).toBe(false)
			expect(store.loading).toBe(false)
			expect(store.error).toBeNull()
			expect(store.updatedAt).toBeNull()
		})

		it('defaults to false when localStorage throws', () => {
			const getItemSpy = vi.spyOn(Storage.prototype, 'getItem').mockImplementation(() => {
				throw new Error('Storage error')
			})

			setActivePinia(createPinia())
			const store = useSystemModeStore()

			expect(store.syntheticMode).toBe(false)
			getItemSpy.mockRestore()
		})
	})

	describe('loadMode', () => {
		it('loads mode from API and updates state', async () => {
			const { api } = await import('@/lib/api')
			const { setSyntheticModeStorage } = await import('@/lib/api')
			vi.mocked(api.get).mockResolvedValue({
				enabled: true,
				updated_at: '2026-01-01T00:00:00Z',
			})

			const store = useSystemModeStore()
			await store.loadMode()

			expect(api.get).toHaveBeenCalledWith('/api/v1/system/mode')
			expect(store.syntheticMode).toBe(true)
			expect(store.loading).toBe(false)
			expect(store.error).toBeNull()
			expect(store.updatedAt).toBeInstanceOf(Date)
			expect(setSyntheticModeStorage).toHaveBeenCalledWith(true)
		})

		it('sets error when API fails', async () => {
			const { api } = await import('@/lib/api')
			vi.mocked(api.get).mockRejectedValue(new Error('Network error'))

			const store = useSystemModeStore()
			await store.loadMode()

			expect(store.error).toBe('Network error')
			expect(store.loading).toBe(false)
		})

		it('uses fallback error message when error has no message', async () => {
			const { api } = await import('@/lib/api')
			vi.mocked(api.get).mockRejectedValue({})

			const store = useSystemModeStore()
			await store.loadMode()

			expect(store.error).toBe('Failed to load synthetic mode')
			expect(store.loading).toBe(false)
		})
	})

	describe('toggleMode', () => {
		it('toggles mode via API and returns new value', async () => {
			const { api } = await import('@/lib/api')
			vi.mocked(api.put).mockResolvedValue({
				enabled: true,
				updated_at: '2026-01-01T12:00:00Z',
			})

			const store = useSystemModeStore()
			expect(store.syntheticMode).toBe(false)

			const result = await store.toggleMode()

			expect(api.put).toHaveBeenCalledWith('/api/v1/system/mode', { enabled: true })
			expect(result).toBe(true)
			expect(store.syntheticMode).toBe(true)
			expect(store.loading).toBe(false)
		})

		it('throws and sets error when API fails', async () => {
			const { api } = await import('@/lib/api')
			vi.mocked(api.put).mockRejectedValue(new Error('Toggle failed'))

			const store = useSystemModeStore()

			await expect(store.toggleMode()).rejects.toThrow('Toggle failed')
			expect(store.error).toBe('Toggle failed')
			expect(store.loading).toBe(false)
		})

		it('uses fallback error message when error has no message', async () => {
			const { api } = await import('@/lib/api')
			vi.mocked(api.put).mockRejectedValue({})

			const store = useSystemModeStore()

			await expect(store.toggleMode()).rejects.toBeDefined()
			expect(store.error).toBe('Failed to toggle synthetic mode')
		})
	})

	describe('setMode', () => {
		it('sets mode to specific value via API', async () => {
			const { api } = await import('@/lib/api')
			const { setSyntheticModeStorage } = await import('@/lib/api')
			vi.mocked(api.put).mockResolvedValue({
				enabled: false,
				updated_at: '2026-01-02T00:00:00Z',
			})

			const store = useSystemModeStore()
			const result = await store.setMode(false)

			expect(api.put).toHaveBeenCalledWith('/api/v1/system/mode', { enabled: false })
			expect(result).toBe(false)
			expect(store.syntheticMode).toBe(false)
			expect(store.updatedAt).toBeInstanceOf(Date)
			expect(setSyntheticModeStorage).toHaveBeenCalledWith(false)
		})

		it('throws and sets error when API fails', async () => {
			const { api } = await import('@/lib/api')
			vi.mocked(api.put).mockRejectedValue(new Error('Set failed'))

			const store = useSystemModeStore()

			await expect(store.setMode(true)).rejects.toThrow('Set failed')
			expect(store.error).toBe('Set failed')
			expect(store.loading).toBe(false)
		})

		it('uses fallback error message when error has no message', async () => {
			const { api } = await import('@/lib/api')
			vi.mocked(api.put).mockRejectedValue({})

			const store = useSystemModeStore()

			await expect(store.setMode(true)).rejects.toBeDefined()
			expect(store.error).toBe('Failed to set synthetic mode')
		})
	})
})
