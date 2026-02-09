import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'
import { createPinia, setActivePinia } from 'pinia'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import { createI18n } from 'vue-i18n'
import RulesView from './RulesView.vue'
import { useAuthStore } from '@/stores/auth'
import { api } from '@/lib/api'

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
    props: ['variant', 'outline', 'rounded'],
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
    props: ['title', 'message', 'retryable'],
    emits: ['retry'],
  },
}))

vi.mock('@/components/RuleFormModal.vue', () => ({
  default: {
    name: 'RuleFormModal',
    template: '<div v-if="modelValue">Rule Form Modal</div>',
    props: [
      'modelValue',
      'isEditing',
      'editingRule',
      'schemaOptions',
      'currentSchema',
      'rulePresets',
      'ruleActions',
      'saving',
      'expressionMode',
      'builderType',
      'polygonConfig',
      'velocityConfig',
      'conditionRows',
      'polygonExpression',
      'velocityExpression',
      'generatedExpression',
    ],
    emits: ['update:modelValue', 'applyPreset', 'submit', 'cancel'],
    slots: {
      'expression-builder': {},
    },
  },
}))

vi.mock('@/components/PolygonBuilder.vue', () => ({
  default: {
    name: 'PolygonBuilder',
    template: '<div>Polygon Builder</div>',
    props: ['config', 'fieldOptions', 'expression'],
    emits: ['addPoint', 'removePoint'],
  },
}))

vi.mock('@/components/VelocityBuilder.vue', () => ({
  default: {
    name: 'VelocityBuilder',
    template: '<div>Velocity Builder</div>',
    props: ['config', 'fieldOptions', 'expression'],
  },
}))

vi.mock('@/components/ConditionBuilder.vue', () => ({
  default: {
    name: 'ConditionBuilder',
    template: '<div>Condition Builder</div>',
    props: ['conditionRows', 'fieldOptions', 'currentSchema', 'generatedExpression'],
    emits: ['removeCondition'],
  },
}))

describe('RulesView', () => {
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
        { path: '/login', component: { template: '<div>Login</div>' } },
        { path: '/rules', component: RulesView },
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
            rules: {
              title: 'Rules',
              subtitle: 'Manage fraud detection rules',
              createRule: 'Create Rule',
              loading: 'Loading...',
              errorTitle: 'Error loading rules',
              emptyText: 'No rules found',
              errorLoad: 'Failed to load rules',
              errorSave: 'Failed to save rule',
              errorDelete: 'Failed to delete rule',
              errorToggle: 'Failed to toggle rule',
              actions: {
                allow: 'Allow',
                block: 'Block',
                review: 'Review',
              },
              presets: {
                highValue: {
                  label: 'High Value',
                  name: 'High Value Transaction',
                  description: 'Flag transactions above threshold',
                },
                suspiciousFlag: {
                  label: 'Suspicious Flag',
                  name: 'Suspicious Flag',
                  description: 'Flag suspicious transactions',
                },
                polygonRestriction: {
                  label: 'Polygon Restriction',
                  name: 'Polygon Restriction',
                  description: 'Restrict by geographic area',
                },
                velocityCheck: {
                  label: 'Velocity Check',
                  name: 'Velocity Check',
                  description: 'Check transaction velocity',
                },
                amountVelocity: {
                  label: 'Amount Velocity',
                  name: 'Amount Velocity',
                  description: 'Check amount velocity',
                },
              },
              modal: {
                ruleConditions: 'Rule Conditions',
                manual: 'Manual',
                builder: 'Builder',
                addCondition: 'Add Condition',
                polygonBuilder: {
                  title: 'Polygon Builder',
                },
                velocityBuilder: {
                  title: 'Velocity Builder',
                },
                polygon: 'Polygon',
                velocity: 'Velocity',
                validation: {
                  selectSchema: 'Please select a schema',
                  noFields: 'Schema has no fields',
                  selectLatLon: 'Please select latitude and longitude fields',
                  polygonPoints: 'Polygon must have at least 3 points',
                  pointInvalidFormat: 'Point {index} has invalid format',
                  pointLatitudeRange: 'Point {index} latitude {value} is out of range',
                  pointLongitudeRange: 'Point {index} longitude {value} is out of range',
                  invalidPolygon: 'Invalid polygon expression',
                  invalidPolygonExpression: 'Invalid polygon expression',
                  selectGroupField: 'Please select a group field',
                  validThreshold: 'Please enter a valid threshold',
                  validTimeValue: 'Please enter a valid time value',
                  invalidVelocity: 'Invalid velocity expression',
                  invalidVelocityExpression: 'Invalid velocity expression',
                  atLeastOneCondition: 'At least one condition is required',
                  conditionFieldRequired: 'Condition {index} field is required',
                  conditionOperatorRequired: 'Condition {index} operator is required',
                  conditionValueRequired: 'Condition {index} value is required',
                  conditionFieldNotExist: 'Condition {index} field {field} does not exist',
                  conditionOperatorNotValid: 'Condition {index} operator {operator} is not valid for type {type}',
                  emptyExpression: 'Expression cannot be empty',
                  invalidExpression: 'Invalid expression',
                },
              },
            },
          },
          components: {
            ruleTable: {
              schema: 'Schema',
              name: 'Name',
              action: 'Action',
              priority: 'Priority',
              status: 'Status',
              actions: 'Actions',
              enabled: 'Enabled',
              disabled: 'Disabled',
              edit: 'Edit',
              delete: 'Delete',
              showing: 'Showing',
              to: 'to',
              of: 'of',
              rules: 'rules',
            },
          },
        },
      },
    })

    vi.clearAllMocks()
  })

  it('renders rules view', () => {
    const authStore = useAuthStore()
    authStore.user = {
      id: '1',
      email: 'user@example.com',
      name: 'User',
      auth_type: 'local',
      active: true,
      roles: [],
    }

    const wrapper = mount(RulesView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    expect(wrapper.text()).toContain('Rules')
  })

  it('redirects to login when user is not authenticated', async () => {
    const authStore = useAuthStore()
    authStore.user = null

    const pushSpy = vi.spyOn(router, 'push')

    mount(RulesView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    await flushPromises()

    expect(pushSpy).toHaveBeenCalledWith('/login')
  })

  it('loads rules and schemas on mount', async () => {
    const authStore = useAuthStore()
    authStore.user = {
      id: '1',
      email: 'user@example.com',
      name: 'User',
      auth_type: 'local',
      active: true,
      roles: [],
    }

    const mockRules = { rules: [] }
    const mockSchemas = { schemas: [] }

    vi.mocked(api.get).mockImplementation((url: string) => {
      if (url === '/api/v1/rules') {
        return Promise.resolve(mockRules)
      }
      if (url === '/api/v1/schemas') {
        return Promise.resolve(mockSchemas)
      }
      return Promise.reject(new Error('Unknown endpoint'))
    })

    const wrapper = mount(RulesView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    await flushPromises()

    expect(api.get).toHaveBeenCalledWith('/api/v1/rules')
    expect(api.get).toHaveBeenCalledWith('/api/v1/schemas')
  })

  it('displays error when loading fails', async () => {
    const authStore = useAuthStore()
    authStore.user = {
      id: '1',
      email: 'user@example.com',
      name: 'User',
      auth_type: 'local',
      active: true,
      roles: [],
    }

    vi.mocked(api.get).mockRejectedValue(new Error('Failed to load'))

    const wrapper = mount(RulesView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    await flushPromises()

    expect((wrapper.vm as any).error).toBeTruthy()
    expect((wrapper.vm as any).loading).toBe(false)
  })

  it('opens create modal when create button is clicked', async () => {
    const authStore = useAuthStore()
    authStore.user = {
      id: '1',
      email: 'user@example.com',
      name: 'User',
      auth_type: 'local',
      active: true,
      roles: [],
    }

    const mockRules = { rules: [] }
    const mockSchemas = { schemas: [] }

    vi.mocked(api.get).mockImplementation((url: string) => {
      if (url === '/api/v1/rules') {
        return Promise.resolve(mockRules)
      }
      if (url === '/api/v1/schemas') {
        return Promise.resolve(mockSchemas)
      }
      return Promise.reject(new Error('Unknown endpoint'))
    })

    const wrapper = mount(RulesView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    await flushPromises()

    ;(wrapper.vm as any).openCreateModal()

    expect((wrapper.vm as any).showModal).toBe(true)
    expect((wrapper.vm as any).isEditing).toBe(false)
  })

  it('deletes rule', async () => {
    const authStore = useAuthStore()
    authStore.user = {
      id: '1',
      email: 'user@example.com',
      name: 'User',
      auth_type: 'local',
      active: true,
      roles: [],
    }

    const mockRules = { rules: [] }
    const mockSchemas = { schemas: [] }

    vi.mocked(api.get).mockImplementation((url: string) => {
      if (url === '/api/v1/rules') {
        return Promise.resolve(mockRules)
      }
      if (url === '/api/v1/schemas') {
        return Promise.resolve(mockSchemas)
      }
      return Promise.reject(new Error('Unknown endpoint'))
    })

    vi.mocked(api.delete).mockResolvedValue({})

    globalThis.confirm = vi.fn(() => true) as typeof confirm

    const wrapper = mount(RulesView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    await flushPromises()

    await (wrapper.vm as any).deleteRule('1')
    await flushPromises()

    expect(api.delete).toHaveBeenCalledWith('/api/v1/rules/1')
  })

  it('toggles rule enabled status', async () => {
    const authStore = useAuthStore()
    authStore.user = {
      id: '1',
      email: 'user@example.com',
      name: 'User',
      auth_type: 'local',
      active: true,
      roles: [],
    }

    const mockRules = { rules: [] }
    const mockSchemas = { schemas: [] }

    vi.mocked(api.get).mockImplementation((url: string) => {
      if (url === '/api/v1/rules') {
        return Promise.resolve(mockRules)
      }
      if (url === '/api/v1/schemas') {
        return Promise.resolve(mockSchemas)
      }
      return Promise.reject(new Error('Unknown endpoint'))
    })

    vi.mocked(api.put).mockResolvedValue({})

    const wrapper = mount(RulesView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    await flushPromises()

    const rule = { id: '1', enabled: true, name: 'Test Rule', action: 'review', priority: 50 }
    await (wrapper.vm as any).toggleRule(rule)
    await flushPromises()

    expect(api.put).toHaveBeenCalledWith('/api/v1/rules/1', { ...rule, enabled: false })
  })

  it('gets action badge variant correctly', () => {
    const authStore = useAuthStore()
    authStore.user = {
      id: '1',
      email: 'user@example.com',
      name: 'User',
      auth_type: 'local',
      active: true,
      roles: [],
    }

    const wrapper = mount(RulesView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    expect((wrapper.vm as any).getActionBadgeVariant('allow')).toBe('success')
    expect((wrapper.vm as any).getActionBadgeVariant('block')).toBe('danger')
    expect((wrapper.vm as any).getActionBadgeVariant('review')).toBe('warning')
    expect((wrapper.vm as any).getActionBadgeVariant('unknown')).toBe('default')
  })

  it('handles pagination correctly', async () => {
    const authStore = useAuthStore()
    authStore.user = {
      id: '1',
      email: 'user@example.com',
      name: 'User',
      auth_type: 'local',
      active: true,
      roles: [],
    }

    const mockRules = {
      rules: Array.from({ length: 25 }, (_, i) => ({
        id: String(i),
        name: `Rule ${i}`,
        enabled: true,
        action: 'review',
        priority: 50,
      })),
    }
    const mockSchemas = { schemas: [] }

    vi.mocked(api.get).mockImplementation((url: string) => {
      if (url === '/api/v1/rules') {
        return Promise.resolve(mockRules)
      }
      if (url === '/api/v1/schemas') {
        return Promise.resolve(mockSchemas)
      }
      return Promise.reject(new Error('Unknown endpoint'))
    })

    const wrapper = mount(RulesView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    await flushPromises()

    expect((wrapper.vm as any).rules.length).toBe(25)
    expect((wrapper.vm as any).paginatedRules.length).toBe(10)

    ;(wrapper.vm as any).nextPage()
    await flushPromises()

    expect((wrapper.vm as any).currentPage).toBe(2)
  })
})
