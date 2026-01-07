import { vi } from 'vitest'

// Stub CSS imports globally
vi.mock('vuetify/lib/components/VCode/VCode.css', () => ({}))
vi.mock('vuetify/lib/directives/ripple/VRipple.css', () => ({}))

