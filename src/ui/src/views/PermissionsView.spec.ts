import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'
import { createPinia, setActivePinia } from 'pinia'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import { createI18n } from 'vue-i18n'
import PermissionsView from './PermissionsView.vue'
import { useAuthStore } from '@/stores/auth'
import { api } from '@/lib/api'

vi.mock('@/lib/api', () => ({
  api: {
    get: vi.fn(),
    post: vi.fn(),
    delete: vi.fn(),
    put: vi.fn(),
  },
  setTokenGetter: vi.fn(),
  setCsrfTokenGetter: vi.fn(),
  setSyntheticModeStorage: vi.fn(),
}))

vi.mock('@/components/BaseButton.vue', () => ({
  default: {
    name: 'BaseButton',
    template: '<button type="button"><slot /></button>',
    props: ['loading', 'disabled', 'size', 'variant', 'prependIcon'],
  },
}))

vi.mock('@/components/BaseBadge.vue', () => ({
  default: {
    name: 'BaseBadge',
    template: '<span class="badge"><slot /></span>',
    props: ['variant', 'size', 'closable'],
    emits: ['close'],
  },
}))

vi.mock('@/components/BaseModal.vue', () => ({
  default: {
    name: 'BaseModal',
    template: '<div v-if="modelValue"><slot /></div>',
    props: ['modelValue', 'title', 'size'],
    emits: ['update:modelValue'],
  },
}))

vi.mock('@/components/BaseTable.vue', () => ({
  default: {
    name: 'BaseTable',
    template: '<div><slot /></div>',
    props: ['columns', 'data', 'emptyText'],
  },
}))

vi.mock('@/components/LoadingSpinner.vue', () => ({
  default: {
    name: 'LoadingSpinner',
    template: '<div v-if="text">{{ text }}</div>',
    props: ['text', 'centered'],
  },
}))

vi.mock('@/components/ErrorMessage.vue', () => ({
  default: {
    name: 'ErrorMessage',
    template: '<div v-if="message">{{ message }}</div>',
    props: ['message'],
    emits: ['dismiss'],
  },
}))

describe('PermissionsView', () => {
  let router: ReturnType<typeof createRouter>
  let vuetify: ReturnType<typeof createVuetify>
  let i18n: ReturnType<typeof createI18n>
  let pinia: ReturnType<typeof createPinia>

  beforeEach(() => {
    pinia = createPinia()
    setActivePinia(pinia)

    router = createRouter({
      history: createMemoryHistory(),
      routes: [
        { path: '/', component: { template: '<div>Home</div>' } },
        { path: '/permissions', component: PermissionsView },
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
          views: {
            permissions: {
              title: 'Permissions',
              subtitle: 'Manage user permissions',
              loading: 'Loading...',
              emptyText: 'No users found',
              errorLoad: 'Failed to load data',
              errorAssignRole: 'Failed to assign role',
              errorRemoveRole: 'Failed to remove role',
              errorToggleStatus: 'Failed to toggle user status',
            },
          },
          components: {
            permissionsTable: {
              user: 'User',
              email: 'Email',
              roles: 'Roles',
              status: 'Status',
              actions: 'Actions',
              addRole: 'Add Role',
              active: 'Active',
              deactivate: 'Deactivate',
              activate: 'Activate',
            },
          },
          common: {
            inactive: 'Inactive',
          },
        },
      },
    })

    vi.clearAllMocks()
  })

  it('renders permissions view', () => {
    const authStore = useAuthStore()
    authStore.user = {
      id: '1',
      email: 'admin@example.com',
      name: 'Admin',
      auth_type: 'local',
      active: true,
      roles: [{ id: '1', name: 'admin', description: 'Administrator' }],
    }

    const wrapper = mount(PermissionsView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    expect(wrapper.text()).toContain('Permissions')
  })

  it('redirects non-admin users to home', async () => {
    const authStore = useAuthStore()
    authStore.user = {
      id: '1',
      email: 'user@example.com',
      name: 'User',
      auth_type: 'local',
      active: true,
      roles: [],
    }

    const pushSpy = vi.spyOn(router, 'push')

    mount(PermissionsView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    await flushPromises()

    expect(pushSpy).toHaveBeenCalledWith('/')
  })

  it('loads users and roles on mount', async () => {
    const authStore = useAuthStore()
    authStore.user = {
      id: '1',
      email: 'admin@example.com',
      name: 'Admin',
      auth_type: 'local',
      active: true,
      roles: [{ id: '1', name: 'admin', description: 'Administrator' }],
    }

    const mockUsers = { users: [] }
    const mockRoles: any[] = []

    vi.mocked(api.get).mockImplementation((url: string) => {
      if (url === '/api/v1/permissions/users') {
        return Promise.resolve(mockUsers)
      }
      if (url === '/api/v1/roles') {
        return Promise.resolve(mockRoles)
      }
      return Promise.reject(new Error('Unknown endpoint'))
    })

    const wrapper = mount(PermissionsView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    await flushPromises()

    expect(api.get).toHaveBeenCalledWith('/api/v1/permissions/users')
    expect(api.get).toHaveBeenCalledWith('/api/v1/roles')
  })

  it('displays error when loading fails', async () => {
    const authStore = useAuthStore()
    authStore.user = {
      id: '1',
      email: 'admin@example.com',
      name: 'Admin',
      auth_type: 'local',
      active: true,
      roles: [{ id: '1', name: 'admin', description: 'Administrator' }],
    }

    vi.mocked(api.get).mockRejectedValue(new Error('Failed to load'))

    const wrapper = mount(PermissionsView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    await flushPromises()

    expect((wrapper.vm as any).error).toBeTruthy()
    expect((wrapper.vm as any).loading).toBe(false)
  })

  it('opens role modal when add role button is clicked', async () => {
    const authStore = useAuthStore()
    authStore.user = {
      id: '1',
      email: 'admin@example.com',
      name: 'Admin',
      auth_type: 'local',
      active: true,
      roles: [{ id: '1', name: 'admin', description: 'Administrator' }],
    }

    const mockUsers = {
      users: [
        {
          id: '2',
          email: 'user@example.com',
          name: 'User',
          active: true,
          roles: [],
        },
      ],
    }
    const mockRoles = [
      { id: '1', name: 'admin', description: 'Administrator' },
      { id: '2', name: 'user', description: 'User' },
    ]

    vi.mocked(api.get).mockImplementation((url: string) => {
      if (url === '/api/v1/permissions/users') {
        return Promise.resolve(mockUsers)
      }
      if (url === '/api/v1/roles') {
        return Promise.resolve(mockRoles)
      }
      return Promise.reject(new Error('Unknown endpoint'))
    })

    const wrapper = mount(PermissionsView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    await flushPromises()

    ;(wrapper.vm as any).openRoleModal(mockUsers.users[0])

    expect((wrapper.vm as any).showRoleModal).toBe(true)
    expect((wrapper.vm as any).selectedUser).toEqual(mockUsers.users[0])
  })

  it('assigns role to user', async () => {
    const authStore = useAuthStore()
    authStore.user = {
      id: '1',
      email: 'admin@example.com',
      name: 'Admin',
      auth_type: 'local',
      active: true,
      roles: [{ id: '1', name: 'admin', description: 'Administrator' }],
    }

    const mockUsers = {
      users: [
        {
          id: '2',
          email: 'user@example.com',
          name: 'User',
          active: true,
          roles: [],
        },
      ],
    }
    const mockRoles = [
      { id: '1', name: 'admin', description: 'Administrator' },
      { id: '2', name: 'user', description: 'User' },
    ]

    vi.mocked(api.get).mockImplementation((url: string) => {
      if (url === '/api/v1/permissions/users') {
        return Promise.resolve(mockUsers)
      }
      if (url === '/api/v1/roles') {
        return Promise.resolve(mockRoles)
      }
      return Promise.reject(new Error('Unknown endpoint'))
    })

    vi.mocked(api.post).mockResolvedValue({})

    const wrapper = mount(PermissionsView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    await flushPromises()

    ;(wrapper.vm as any).selectedUser = mockUsers.users[0]
    await (wrapper.vm as any).assignRole('2')
    await flushPromises()

    expect(api.post).toHaveBeenCalledWith('/api/v1/permissions/users/2/roles', { role_id: '2' })
    expect((wrapper.vm as any).showRoleModal).toBe(false)
  })

  it('removes role from user', async () => {
    const authStore = useAuthStore()
    authStore.user = {
      id: '1',
      email: 'admin@example.com',
      name: 'Admin',
      auth_type: 'local',
      active: true,
      roles: [{ id: '1', name: 'admin', description: 'Administrator' }],
    }

    const mockUsers = {
      users: [
        {
          id: '2',
          email: 'user@example.com',
          name: 'User',
          active: true,
          roles: [{ id: '2', name: 'user', description: 'User' }],
        },
      ],
    }
    const mockRoles = [
      { id: '1', name: 'admin', description: 'Administrator' },
      { id: '2', name: 'user', description: 'User' },
    ]

    vi.mocked(api.get).mockImplementation((url: string) => {
      if (url === '/api/v1/permissions/users') {
        return Promise.resolve(mockUsers)
      }
      if (url === '/api/v1/roles') {
        return Promise.resolve(mockRoles)
      }
      return Promise.reject(new Error('Unknown endpoint'))
    })

    vi.mocked(api.delete).mockResolvedValue({})

    const wrapper = mount(PermissionsView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    await flushPromises()

    await (wrapper.vm as any).removeRole('2', '2')
    await flushPromises()

    expect(api.delete).toHaveBeenCalledWith('/api/v1/permissions/users/2/roles/2')
  })

  it('toggles user active status', async () => {
    const authStore = useAuthStore()
    authStore.user = {
      id: '1',
      email: 'admin@example.com',
      name: 'Admin',
      auth_type: 'local',
      active: true,
      roles: [{ id: '1', name: 'admin', description: 'Administrator' }],
    }

    const mockUsers = {
      users: [
        {
          id: '2',
          email: 'user@example.com',
          name: 'User',
          active: true,
          roles: [],
        },
      ],
    }
    const mockRoles: any[] = []

    vi.mocked(api.get).mockImplementation((url: string) => {
      if (url === '/api/v1/permissions/users') {
        return Promise.resolve(mockUsers)
      }
      if (url === '/api/v1/roles') {
        return Promise.resolve(mockRoles)
      }
      return Promise.reject(new Error('Unknown endpoint'))
    })

    vi.mocked(api.put).mockResolvedValue({})

    const wrapper = mount(PermissionsView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    await flushPromises()

    await (wrapper.vm as any).toggleUserActive(mockUsers.users[0])
    await flushPromises()

    expect(api.put).toHaveBeenCalledWith('/api/v1/permissions/users/2/active', { active: false })
  })

  it('filters available roles correctly', async () => {
    const authStore = useAuthStore()
    authStore.user = {
      id: '1',
      email: 'admin@example.com',
      name: 'Admin',
      auth_type: 'local',
      active: true,
      roles: [{ id: '1', name: 'admin', description: 'Administrator' }],
    }

    const mockUsers = {
      users: [
        {
          id: '2',
          email: 'user@example.com',
          name: 'User',
          active: true,
          roles: [{ id: '2', name: 'user', description: 'User' }],
        },
      ],
    }
    const mockRoles = {
      roles: [
        { id: '1', name: 'admin', description: 'Administrator' },
        { id: '2', name: 'user', description: 'User' },
        { id: '3', name: 'editor', description: 'Editor' },
      ],
    }

    vi.mocked(api.get).mockImplementation((url: string) => {
      if (url === '/api/v1/permissions/users') {
        return Promise.resolve(mockUsers)
      }
      if (url === '/api/v1/roles') {
        return Promise.resolve(mockRoles)
      }
      return Promise.reject(new Error('Unknown endpoint'))
    })

    const wrapper = mount(PermissionsView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    await flushPromises()

    ;(wrapper.vm as any).selectedUser = mockUsers.users[0]
    await flushPromises()

    const availableRoles = (wrapper.vm as any).availableRoles
    expect(availableRoles.length).toBe(2)
    expect(availableRoles.find((r: any) => r.name === 'admin')).toBeTruthy()
    expect(availableRoles.find((r: any) => r.name === 'editor')).toBeTruthy()
  })

  it('shows error when assign role fails', async () => {
    const authStore = useAuthStore()
    authStore.user = {
      id: '1',
      email: 'admin@example.com',
      name: 'Admin',
      auth_type: 'local',
      active: true,
      roles: [{ id: '1', name: 'admin', description: 'Administrator' }],
    }

    const mockUsers = {
      users: [
        {
          id: '2',
          email: 'user@example.com',
          name: 'User',
          active: true,
          roles: [],
        },
      ],
    }
    const mockRoles: any[] = []

    vi.mocked(api.get).mockImplementation((url: string) => {
      if (url === '/api/v1/permissions/users') {
        return Promise.resolve(mockUsers)
      }
      if (url === '/api/v1/roles') {
        return Promise.resolve(mockRoles)
      }
      return Promise.reject(new Error('Unknown endpoint'))
    })

    vi.mocked(api.post).mockRejectedValue(new Error('Failed to assign'))

    const wrapper = mount(PermissionsView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    await flushPromises()

    ;(wrapper.vm as any).selectedUser = mockUsers.users[0]
    await (wrapper.vm as any).assignRole('1')
    await flushPromises()

    expect((wrapper.vm as any).error).toBeTruthy()
  })
})
