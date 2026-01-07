import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'
import { createPinia, setActivePinia } from 'pinia'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import { createI18n } from 'vue-i18n'
import Header from './Header.vue'
import { useAuthStore } from '@/stores/auth'
import { useBrandingStore } from '@/stores/branding'

describe('Header', () => {
  let router: ReturnType<typeof createRouter>
  let vuetify: ReturnType<typeof createVuetify>
  let i18n: ReturnType<typeof createI18n>

  beforeEach(() => {
    setActivePinia(createPinia())
    
    router = createRouter({
      history: createMemoryHistory(),
      routes: [
        { path: '/', component: { template: '<div>Home</div>' } },
        { path: '/login', component: { template: '<div>Login</div>' } },
        { path: '/profile', component: { template: '<div>Profile</div>' } },
      ],
    })

    vuetify = createVuetify({
      components,
      directives,
    })

    i18n = createI18n({
      legacy: true,
      locale: 'en-US',
      messages: {
        'en-US': {
          common: { appName: 'AlgoShield' },
        },
      },
    })

    vi.clearAllMocks()
  })

  describe('rendering', () => {
    it('does not render when user is null', () => {
      const authStore = useAuthStore()
      authStore.user = null

      const wrapper = mount(Header, {
        global: {
          plugins: [router, vuetify, i18n],
        },
      })

      expect(wrapper.html()).toBe('')
    })

    it('does not render on login page', async () => {
      const authStore = useAuthStore()
      authStore.user = {
        id: '1',
        email: 'test@example.com',
        name: 'Test User',
        auth_type: 'local',
        active: true,
      }

      await router.push('/login')

      const wrapper = mount(Header, {
        global: {
          plugins: [router, vuetify, i18n],
        },
      })

      expect(wrapper.html()).toBe('')
    })

    it('renders when user is logged in and not on login page', async () => {
      const authStore = useAuthStore()
      authStore.user = {
        id: '1',
        email: 'test@example.com',
        name: 'Test User',
        auth_type: 'local',
        active: true,
      }

      await router.push('/')

      const wrapper = mount(Header, {
        global: {
          plugins: [router, vuetify, i18n],
        },
      })

      expect(wrapper.html()).not.toBe('')
    })
  })

  describe('branding', () => {
    it('displays custom app name when configured', async () => {
      const authStore = useAuthStore()
      authStore.user = {
        id: '1',
        email: 'test@example.com',
        name: 'Test User',
        auth_type: 'local',
        active: true,
      }

      const brandingStore = useBrandingStore()
      brandingStore.config = {
        id: 1,
        app_name: 'Custom App',
        icon_url: '/custom-icon.png',
        primary_color: '#3B82F6',
        secondary_color: '#10B981',
        header_color: '#123456',
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z',
      }

      await router.push('/')

      const wrapper = mount(Header, {
        global: {
          plugins: [router, vuetify, i18n],
        },
      })

      expect(wrapper.text()).toContain('Custom App')
    })

    it('displays default app name when not configured', async () => {
      const authStore = useAuthStore()
      authStore.user = {
        id: '1',
        email: 'test@example.com',
        name: 'Test User',
        auth_type: 'local',
        active: true,
      }

      const brandingStore = useBrandingStore()
      brandingStore.config = null

      await router.push('/')

      const wrapper = mount(Header, {
        global: {
          plugins: [router, vuetify, i18n],
        },
      })

      expect(wrapper.html()).toBeTruthy()
    })
  })

  describe('user menu', () => {
    it('displays user name and email', async () => {
      const authStore = useAuthStore()
      authStore.user = {
        id: '1',
        email: 'test@example.com',
        name: 'Test User',
        auth_type: 'local',
        active: true,
      }

      await router.push('/')

      const wrapper = mount(Header, {
        global: {
          plugins: [router, vuetify, i18n],
        },
      })

      expect(wrapper.text()).toContain('Test User')
      expect(wrapper.text()).toContain('test@example.com')
    })

    it('calls logout when logout button is clicked', async () => {
      const authStore = useAuthStore()
      authStore.user = {
        id: '1',
        email: 'test@example.com',
        name: 'Test User',
        auth_type: 'local',
        active: true,
      }

      const logoutSpy = vi.spyOn(authStore, 'logout').mockResolvedValue()

      await router.push('/')

      const wrapper = mount(Header, {
        global: {
          plugins: [router, vuetify, i18n],
        },
      })

      await wrapper.vm.handleLogout()

      expect(logoutSpy).toHaveBeenCalled()
    })
  })

  describe('logo error handling', () => {
    it('falls back to default logo on error', async () => {
      const authStore = useAuthStore()
      authStore.user = {
        id: '1',
        email: 'test@example.com',
        name: 'Test User',
        auth_type: 'local',
        active: true,
      }

      await router.push('/')

      const wrapper = mount(Header, {
        global: {
          plugins: [router, vuetify, i18n],
        },
      })

      const mockEvent = {
        target: { src: '' },
      } as unknown as Event

      wrapper.vm.handleLogoError(mockEvent)

      expect((mockEvent.target as HTMLImageElement).src).toBe('/gopher.png')
    })
  })
})


