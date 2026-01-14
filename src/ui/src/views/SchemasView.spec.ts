import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'
import { createPinia, setActivePinia } from 'pinia'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import { createI18n } from 'vue-i18n'
import SchemasView from './SchemasView.vue'
import { useAuthStore } from '@/stores/auth'
import { api } from '@/lib/api'

vi.mock('@/lib/api', () => ({
  api: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  },
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
    props: ['variant', 'size'],
  },
}))

vi.mock('@/components/BaseModal.vue', () => ({
  default: {
    name: 'BaseModal',
    template: '<div v-if="modelValue"><slot /></div>',
    props: ['modelValue', 'title', 'size'],
    emits: ['update:modelValue'],
    slots: {
      footer: {},
    },
  },
}))

vi.mock('@/components/BaseTable.vue', () => ({
  default: {
    name: 'BaseTable',
    template: '<div><slot /></div>',
    props: ['columns', 'data', 'emptyText'],
  },
}))

vi.mock('@/components/BaseInput.vue', () => ({
  default: {
    name: 'BaseInput',
    template: '<input :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />',
    props: ['modelValue', 'label', 'placeholder', 'required', 'prependInnerIcon', 'class', 'rules'],
    emits: ['update:modelValue'],
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

describe('SchemasView', () => {
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
        { path: '/schemas', component: SchemasView },
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
            schemas: {
              title: 'Schemas',
              subtitle: 'Manage event schemas',
              createSchema: 'Create Schema',
              loading: 'Loading...',
              errorTitle: 'Error loading schemas',
              emptyText: 'No schemas found',
              errorLoad: 'Failed to load schemas',
              errorSave: 'Failed to save schema',
              errorDelete: 'Failed to delete schema',
              modalCreateTitle: 'Create Schema',
              modalEditTitle: 'Edit Schema',
              name: 'Name',
              namePlaceholder: 'Enter schema name',
              description: 'Description',
              descriptionPlaceholder: 'Enter schema description',
              sampleJson: 'Sample JSON',
              sampleJsonRequired: '*',
              sampleJsonPlaceholder: 'Paste JSON here...',
              pasteJsonHint: 'Paste a sample JSON object to extract fields',
              extractedFields: 'Extracted Fields',
              path: 'Path',
              type: 'Type',
              sampleValue: 'Sample Value',
              noDescription: 'No description',
              fields: 'fields',
            },
          },
          components: {
            schemaTable: {
              name: 'Name',
              fields: 'Fields',
              created: 'Created',
              actions: 'Actions',
              view: 'View',
              edit: 'Edit',
              delete: 'Delete',
            },
            modal: {
              cancel: 'Cancel',
              save: 'Save',
            },
          },
        },
      },
    })

    vi.clearAllMocks()
  })

  it('renders schemas view', () => {
    const authStore = useAuthStore()
    authStore.user = {
      id: '1',
      email: 'user@example.com',
      name: 'User',
      auth_type: 'local',
      active: true,
      roles: [],
    }

    const wrapper = mount(SchemasView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    expect(wrapper.text()).toContain('Schemas')
  })

  it('redirects to login when user is not authenticated', async () => {
    const authStore = useAuthStore()
    authStore.user = null

    const pushSpy = vi.spyOn(router, 'push')

    mount(SchemasView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    await flushPromises()

    expect(pushSpy).toHaveBeenCalledWith('/login')
  })

  it('loads schemas on mount', async () => {
    const authStore = useAuthStore()
    authStore.user = {
      id: '1',
      email: 'user@example.com',
      name: 'User',
      auth_type: 'local',
      active: true,
      roles: [],
    }

    const mockSchemas = {
      schemas: [
        {
          id: '1',
          name: 'Transaction',
          description: 'Transaction schema',
          sample_json: { amount: 100 },
          extracted_fields: [],
          created_at: '2024-01-01T00:00:00Z',
          updated_at: '2024-01-01T00:00:00Z',
        },
      ],
    }

    vi.mocked(api.get).mockResolvedValue(mockSchemas)

    const wrapper = mount(SchemasView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    await flushPromises()

    expect(api.get).toHaveBeenCalledWith('/api/v1/schemas')
    expect((wrapper.vm as any).schemas.length).toBe(1)
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

    const wrapper = mount(SchemasView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    await flushPromises()

    expect((wrapper.vm as any).error).toBeTruthy()
    expect((wrapper.vm as any).loading).toBe(false)
  })

  it('opens create modal when create button is clicked', () => {
    const authStore = useAuthStore()
    authStore.user = {
      id: '1',
      email: 'user@example.com',
      name: 'User',
      auth_type: 'local',
      active: true,
      roles: [],
    }

    const wrapper = mount(SchemasView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    ;(wrapper.vm as any).openCreateModal()

    expect((wrapper.vm as any).showModal).toBe(true)
    expect((wrapper.vm as any).isEditing).toBe(false)
    expect((wrapper.vm as any).editingSchema.name).toBe('')
  })

  it('opens edit modal with schema data', () => {
    const authStore = useAuthStore()
    authStore.user = {
      id: '1',
      email: 'user@example.com',
      name: 'User',
      auth_type: 'local',
      active: true,
      roles: [],
    }

    const schema = {
      id: '1',
      name: 'Transaction',
      description: 'Transaction schema',
      sample_json: { amount: 100 },
      extracted_fields: [],
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
    }

    const wrapper = mount(SchemasView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    ;(wrapper.vm as any).openEditModal(schema)

    expect((wrapper.vm as any).showModal).toBe(true)
    expect((wrapper.vm as any).isEditing).toBe(true)
    expect((wrapper.vm as any).editingSchema.name).toBe('Transaction')
  })

  it('parses valid JSON and extracts fields', () => {
    const authStore = useAuthStore()
    authStore.user = {
      id: '1',
      email: 'user@example.com',
      name: 'User',
      auth_type: 'local',
      active: true,
      roles: [],
    }

    const wrapper = mount(SchemasView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    ;(wrapper.vm as any).sampleJsonText = JSON.stringify({ amount: 100, currency: 'USD' })
    ;(wrapper.vm as any).parseSampleJson()

    expect((wrapper.vm as any).jsonError).toBe('')
    expect((wrapper.vm as any).extractedFields.length).toBeGreaterThan(0)
  })

  it('shows error for invalid JSON', () => {
    const authStore = useAuthStore()
    authStore.user = {
      id: '1',
      email: 'user@example.com',
      name: 'User',
      auth_type: 'local',
      active: true,
      roles: [],
    }

    const wrapper = mount(SchemasView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    ;(wrapper.vm as any).sampleJsonText = 'invalid json'
    ;(wrapper.vm as any).parseSampleJson()

    expect((wrapper.vm as any).jsonError).toBeTruthy()
    expect((wrapper.vm as any).extractedFields.length).toBe(0)
  })

  it('creates new schema', async () => {
    const authStore = useAuthStore()
    authStore.user = {
      id: '1',
      email: 'user@example.com',
      name: 'User',
      auth_type: 'local',
      active: true,
      roles: [],
    }

    const mockSchemas = { schemas: [] }
    const newSchema = {
      id: '1',
      name: 'New Schema',
      description: 'Description',
      sample_json: { amount: 100 },
      extracted_fields: [],
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
    }

    vi.mocked(api.get).mockResolvedValue(mockSchemas)
    vi.mocked(api.post).mockResolvedValue(newSchema)

    const wrapper = mount(SchemasView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    await flushPromises()

    ;(wrapper.vm as any).editingSchema.name = 'New Schema'
    ;(wrapper.vm as any).editingSchema.description = 'Description'
    ;(wrapper.vm as any).sampleJsonText = JSON.stringify({ amount: 100 })
    ;(wrapper.vm as any).parseSampleJson()

    await (wrapper.vm as any).handleSubmit()
    await flushPromises()

    expect(api.post).toHaveBeenCalledWith('/api/v1/schemas', {
      name: 'New Schema',
      description: 'Description',
      sample_json: { amount: 100 },
    })
  })

  it('updates existing schema', async () => {
    const authStore = useAuthStore()
    authStore.user = {
      id: '1',
      email: 'user@example.com',
      name: 'User',
      auth_type: 'local',
      active: true,
      roles: [],
    }

    const mockSchemas = { schemas: [] }
    const updatedSchema = {
      id: '1',
      name: 'Updated Schema',
      description: 'Updated',
      sample_json: { amount: 200 },
      extracted_fields: [],
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
    }

    vi.mocked(api.get).mockResolvedValue(mockSchemas)
    vi.mocked(api.put).mockResolvedValue(updatedSchema)

    const wrapper = mount(SchemasView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    await flushPromises()

    ;(wrapper.vm as any).isEditing = true
    ;(wrapper.vm as any).editingSchema.id = '1'
    ;(wrapper.vm as any).editingSchema.name = 'Updated Schema'
    ;(wrapper.vm as any).editingSchema.description = 'Updated'
    ;(wrapper.vm as any).sampleJsonText = JSON.stringify({ amount: 200 })
    ;(wrapper.vm as any).parseSampleJson()

    await (wrapper.vm as any).handleSubmit()
    await flushPromises()

    expect(api.put).toHaveBeenCalledWith('/api/v1/schemas/1', {
      name: 'Updated Schema',
      description: 'Updated',
      sample_json: { amount: 200 },
    })
  })

  it('deletes schema', async () => {
    const authStore = useAuthStore()
    authStore.user = {
      id: '1',
      email: 'user@example.com',
      name: 'User',
      auth_type: 'local',
      active: true,
      roles: [],
    }

    const mockSchemas = { schemas: [] }

    vi.mocked(api.get).mockResolvedValue(mockSchemas)
    vi.mocked(api.delete).mockResolvedValue({})

    globalThis.confirm = vi.fn(() => true) as typeof confirm

    const wrapper = mount(SchemasView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    await flushPromises()

    await (wrapper.vm as any).deleteSchema('1')
    await flushPromises()

    expect(api.delete).toHaveBeenCalledWith('/api/v1/schemas/1')
  })

  it('formats date correctly', () => {
    const authStore = useAuthStore()
    authStore.user = {
      id: '1',
      email: 'user@example.com',
      name: 'User',
      auth_type: 'local',
      active: true,
      roles: [],
    }

    const wrapper = mount(SchemasView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    const formatted = (wrapper.vm as any).formatDate('2024-01-15T12:00:00Z')
    expect(formatted).toContain('Jan')
    expect(formatted).toContain('15')
    expect(formatted).toContain('2024')
  })

  it('formats sample value correctly', () => {
    const authStore = useAuthStore()
    authStore.user = {
      id: '1',
      email: 'user@example.com',
      name: 'User',
      auth_type: 'local',
      active: true,
      roles: [],
    }

    const wrapper = mount(SchemasView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    expect((wrapper.vm as any).formatSampleValue('test')).toBe('"test"')
    expect((wrapper.vm as any).formatSampleValue(100)).toBe('100')
    expect((wrapper.vm as any).formatSampleValue(null)).toBe('null')
  })
})
