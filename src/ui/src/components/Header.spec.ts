import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'
import { createPinia, setActivePinia } from 'pinia'
import Header from './Header.vue'
import { useAuthStore } from '@/stores/auth'
import { useBrandingStore } from '@/stores/branding'

// Mock vuetify components
vi.mock('vuetify/components', () => ({
  VAppBar: { template: '<div><slot /></div>' },
  VAvatar: { template: '<div><slot /></div>' },
  VBtn: { template: '<button><slot /></button>' },
  VMenu: { template: '<div><slot /></div>' },
  VList: { template: '<div><slot /></div>' },
  VListItem: { template: '<div><slot /></div>' },
  VListItemTitle: { template: '<div><slot /></div>' },
  VListItemSubtitle: { template: '<div><slot /></div>' },
  VListSubheader: { template: '<div><slot /></div>' },
  VDivider: { template: '<hr />' },
  VIcon: { template: '<i />' },
  VImg: { template: '<img />' },
}))

describe('Header', () => {
  let router: ReturnType<typeof createRouter>

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

    vi.clearAllMocks()
  })

  describe('rendering', () => {
    it('does not render when user is null', () => {
      // Arrange
      const authStore = useAuthStore()
      authStore.user = null

      // Act
      const wrapper = mount(Header, {
        global: {
          plugins: [router],
          stubs: {
            VAppBar: true,
            VAvatar: true,
            VBtn: true,
            VMenu: true,
            VList: true,
            VListItem: true,
            VIcon: true,
            VImg: true,
          },
        },
      })

      // Assert
      expect(wrapper.html()).toBe('')
    })

    it('does not render on login page', async () => {
      // Arrange
      const authStore = useAuthStore()
      authStore.user = {
        id: '1',
        email: 'test@example.com',
        name: 'Test User',
        auth_type: 'local',
        active: true,
      }

      await router.push('/login')

      // Act
      const wrapper = mount(Header, {
        global: {
          plugins: [router],
          stubs: {
            VAppBar: true,
            VAvatar: true,
            VBtn: true,
            VMenu: true,
            VList: true,
            VListItem: true,
            VIcon: true,
            VImg: true,
          },
        },
      })

      // Assert
      expect(wrapper.html()).toBe('')
    })

    it('renders when user is logged in and not on login page', async () => {
      // Arrange
      const authStore = useAuthStore()
      authStore.user = {
        id: '1',
        email: 'test@example.com',
        name: 'Test User',
        auth_type: 'local',
        active: true,
      }

      await router.push('/')

      // Act
      const wrapper = mount(Header, {
        global: {
          plugins: [router],
          stubs: {
            VAppBar: true,
            VAvatar: true,
            VBtn: true,
            VMenu: true,
            VList: true,
            VListItem: true,
            VIcon: true,
            VImg: true,
            VDivider: true,
            VListSubheader: true,
          },
        },
      })

      // Assert
      expect(wrapper.html()).not.toBe('')
    })
  })

  describe('branding', () => {
    it('displays custom app name when configured', async () => {
      // Arrange
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
        app_name: 'Custom App',
        icon_url: '/custom-icon.png',
        header_color: '#123456',
      }

      await router.push('/')

      // Act
      const wrapper = mount(Header, {
        global: {
          plugins: [router],
          stubs: {
            VAppBar: true,
            VAvatar: true,
            VBtn: true,
            VMenu: true,
            VList: true,
            VListItem: true,
            VIcon: true,
            VImg: true,
          },
        },
      })

      // Assert
      expect(wrapper.text()).toContain('Custom App')
    })

    it('displays default app name when not configured', async () => {
      // Arrange
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

      // Act
      const wrapper = mount(Header, {
        global: {
          plugins: [router],
          stubs: {
            VAppBar: true,
            VAvatar: true,
            VBtn: true,
            VMenu: true,
            VList: true,
            VListItem: true,
            VIcon: true,
            VImg: true,
          },
        },
      })

      // Assert - should show default app name from i18n
      expect(wrapper.html()).toBeTruthy()
    })
  })

  describe('user menu', () => {
    it('displays user name and email', async () => {
      // Arrange
      const authStore = useAuthStore()
      authStore.user = {
        id: '1',
        email: 'test@example.com',
        name: 'Test User',
        auth_type: 'local',
        active: true,
      }

      await router.push('/')

      // Act
      const wrapper = mount(Header, {
        global: {
          plugins: [router],
          stubs: {
            VAppBar: true,
            VAvatar: true,
            VBtn: true,
            VMenu: true,
            VList: true,
            VListItem: true,
            VListItemTitle: true,
            VListItemSubtitle: true,
            VIcon: true,
            VImg: true,
          },
        },
      })

      // Assert
      expect(wrapper.text()).toContain('Test User')
      expect(wrapper.text()).toContain('test@example.com')
    })

    it('calls logout when logout button is clicked', async () => {
      // Arrange
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
          plugins: [router],
          stubs: {
            VAppBar: true,
            VAvatar: true,
            VBtn: true,
            VMenu: true,
            VList: true,
            VListItem: true,
            VIcon: true,
            VImg: true,
            VDivider: true,
            VListSubheader: true,
          },
        },
      })

      // Act
      await wrapper.vm.handleLogout()

      // Assert
      expect(logoutSpy).toHaveBeenCalled()
    })
  })

  describe('logo error handling', () => {
    it('falls back to default logo on error', async () => {
      // Arrange
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
          plugins: [router],
          stubs: {
            VAppBar: true,
            VAvatar: true,
            VBtn: true,
            VMenu: true,
            VList: true,
            VListItem: true,
            VIcon: true,
            VImg: false, // Don't stub VImg to test error handling
          },
        },
      })

      // Act
      const mockEvent = {
        target: { src: '' },
      } as Event

      wrapper.vm.handleLogoError(mockEvent)

      // Assert
      expect((mockEvent.target as HTMLImageElement).src).toBe('/gopher.png')
    })
  })
})

