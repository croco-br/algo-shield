import { vi } from 'vitest'

// Global mock for setTokenGetter and setCsrfTokenGetter to avoid errors in tests
// This is needed because the auth store calls these functions on initialization
vi.mock('@/lib/api', async (importOriginal) => {
	const actual = await importOriginal<typeof import('@/lib/api')>()
	return {
		...actual,
		setTokenGetter: vi.fn(),
		setCsrfTokenGetter: vi.fn(),
	}
})
