import { vi } from 'vitest'

// Polyfill for Vuetify v-dialog/overlay (uses visualViewport; missing in happy-dom)
if (typeof window !== 'undefined' && !window.visualViewport) {
  Object.defineProperty(window, 'visualViewport', {
    value: {
      width: window.innerWidth,
      height: window.innerHeight,
      offsetLeft: 0,
      offsetTop: 0,
      pageLeft: 0,
      pageTop: 0,
      scale: 1,
      addEventListener: vi.fn(),
      removeEventListener: vi.fn(),
    },
    writable: true,
  })
}

// Stub CSS imports globally
// This prevents errors when Vuetify or other libraries try to import CSS files
vi.mock('vuetify/lib/components/VCode/VCode.css', () => ({}))
vi.mock('vuetify/lib/directives/ripple/VRipple.css', () => ({}))
vi.mock('vuetify/styles', () => ({}))
vi.mock('@fortawesome/fontawesome-free/css/all.css', () => ({}))
vi.mock('./assets/main.css', () => ({}))

// SECURITY: Mock setTokenGetter globally for tests
// The auth store calls setTokenGetter on initialization, so we need to mock it
// Individual test files can override this mock if needed
vi.mock('@/lib/api', async (importOriginal) => {
	const actual = await importOriginal<typeof import('@/lib/api')>()
	return {
		...actual,
		setTokenGetter: vi.fn(),
	}
})

