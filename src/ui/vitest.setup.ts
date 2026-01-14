import { vi } from 'vitest'

// Stub CSS imports globally
// This prevents errors when Vuetify or other libraries try to import CSS files
vi.mock('vuetify/lib/components/VCode/VCode.css', () => ({}))
vi.mock('vuetify/lib/directives/ripple/VRipple.css', () => ({}))
vi.mock('vuetify/styles', () => ({}))
vi.mock('@fortawesome/fontawesome-free/css/all.css', () => ({}))
vi.mock('./assets/main.css', () => ({}))

