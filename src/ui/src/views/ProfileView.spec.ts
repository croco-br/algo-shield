import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'
import { createPinia, setActivePinia } from 'pinia'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import { createI18n } from 'vue-i18n'
import ProfileView from './ProfileView.vue'
import { useAuthStore } from '@/stores/auth'

vi.mock('@/components/BaseBadge.vue', () => ({
  default: {
    name: 'BaseBadge',
    template: '<span class="badge"><slot /></span>',
    props: ['variant', 'size'],
  },
}))

describe('ProfileView', () => {
  let router: ReturnType<typeof createRouter>
  let vuetify: ReturnType<typeof createVuetify>
  let i18n: ReturnType<typeof createI18n>
  let pinia: ReturnType<typeof createPinia>

  beforeEach(() => {
    pinia = createPinia()
    setActivePinia(pinia)

    router = createRouter({
      history: createMemoryHistory(),
      routes: [{ path: '/profile', component: ProfileView }],
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
          views: {
            profile: {
              title: 'Profile',
              subtitle: 'Your profile information',
              userInfo: 'User Information',
              userInfoSubtitle: 'View your account details',
              name: 'Name',
              email: 'Email',
              authType: 'Authentication Type',
              roles: 'Roles',
              status: 'Status',
              notAvailable: 'Not available',
              noRoles: 'No roles assigned',
            },
          },
          common: {
            active: 'Active',
            inactive: 'Inactive',
          },
        },
      },
    })

    vi.clearAllMocks()
  })

  it('renders profile view with user information', () => {
    const authStore = useAuthStore()
    authStore.user = {
      id: '1',
      email: 'test@example.com',
      name: 'Test User',
      auth_type: 'local',
      active: true,
      roles: [],
    }

    const wrapper = mount(ProfileView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    expect(wrapper.text()).toContain('User Information')
    expect(wrapper.text()).toContain('Test User')
    expect(wrapper.text()).toContain('test@example.com')
  })

  it('displays user name when available', () => {
    const authStore = useAuthStore()
    authStore.user = {
      id: '1',
      email: 'test@example.com',
      name: 'John Doe',
      auth_type: 'local',
      active: true,
      roles: [],
    }

    const wrapper = mount(ProfileView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    expect(wrapper.text()).toContain('John Doe')
  })

  it('displays not available when user name is missing', () => {
    const authStore = useAuthStore()
    authStore.user = {
      id: '1',
      email: 'test@example.com',
      name: '',
      auth_type: 'local',
      active: true,
      roles: [],
    }

    const wrapper = mount(ProfileView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    expect(wrapper.text()).toContain('Not available')
  })

  it('displays user email when available', () => {
    const authStore = useAuthStore()
    authStore.user = {
      id: '1',
      email: 'user@example.com',
      name: 'Test User',
      auth_type: 'local',
      active: true,
      roles: [],
    }

    const wrapper = mount(ProfileView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    expect(wrapper.text()).toContain('user@example.com')
  })

  it('displays authentication type', () => {
    const authStore = useAuthStore()
    authStore.user = {
      id: '1',
      email: 'test@example.com',
      name: 'Test User',
      auth_type: 'sso',
      active: true,
      roles: [],
    }

    const wrapper = mount(ProfileView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    expect(wrapper.text().toLowerCase()).toContain('sso')
  })

  it('displays user roles when available', () => {
    const authStore = useAuthStore()
    authStore.user = {
      id: '1',
      email: 'test@example.com',
      name: 'Test User',
      auth_type: 'local',
      active: true,
      roles: [
        { id: '1', name: 'admin', description: 'Administrator' },
        { id: '2', name: 'user', description: 'User' },
      ],
    }

    const wrapper = mount(ProfileView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    expect(wrapper.text()).toContain('admin')
    expect(wrapper.text()).toContain('user')
  })

  it('displays no roles message when user has no roles', () => {
    const authStore = useAuthStore()
    authStore.user = {
      id: '1',
      email: 'test@example.com',
      name: 'Test User',
      auth_type: 'local',
      active: true,
      roles: [],
    }

    const wrapper = mount(ProfileView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    expect(wrapper.text()).toContain('No roles assigned')
  })

  it('displays active status badge when user is active', () => {
    const authStore = useAuthStore()
    authStore.user = {
      id: '1',
      email: 'test@example.com',
      name: 'Test User',
      auth_type: 'local',
      active: true,
      roles: [],
    }

    const wrapper = mount(ProfileView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    expect(wrapper.text()).toContain('Active')
  })

  it('displays inactive status badge when user is not active', () => {
    const authStore = useAuthStore()
    authStore.user = {
      id: '1',
      email: 'test@example.com',
      name: 'Test User',
      auth_type: 'local',
      active: false,
      roles: [],
    }

    const wrapper = mount(ProfileView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    expect(wrapper.text()).toContain('Inactive')
  })

  it('calls refresh when component mounts', async () => {
    const authStore = useAuthStore()
    authStore.user = {
      id: '1',
      email: 'test@example.com',
      name: 'Test User',
      auth_type: 'local',
      active: true,
      roles: [],
    }

    const refreshSpy = vi.spyOn(authStore, 'refresh')

    mount(ProfileView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    await new Promise(resolve => setTimeout(resolve, 0))

    expect(refreshSpy).toHaveBeenCalled()
  })

  it('handles null user gracefully', () => {
    const authStore = useAuthStore()
    authStore.user = null

    const wrapper = mount(ProfileView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    expect(wrapper.text()).toContain('Not available')
  })
})
