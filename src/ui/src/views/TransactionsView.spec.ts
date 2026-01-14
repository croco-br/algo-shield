import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'
import { createPinia, setActivePinia } from 'pinia'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import { createI18n } from 'vue-i18n'
import TransactionsView from './TransactionsView.vue'
import { useSystemModeStore } from '@/stores/systemMode'
import { api } from '@/lib/api'

vi.mock('@/lib/api', () => ({
  api: {
    get: vi.fn(),
    post: vi.fn(),
    patch: vi.fn(),
  },
}))

vi.mock('@/stores/systemMode', () => ({
  useSystemModeStore: vi.fn(),
}))

vi.mock('@/composables/useLocale', () => ({
  useLocale: () => ({
    t: (key: string) => key,
  }),
}))

vi.mock('@/composables/useCurrency', () => ({
  useCurrency: () => ({
    formatCurrency: (amount: number) => `$${amount}`,
    formatCompactCurrency: (amount: number) => `$${amount}`,
    parseCurrency: (value: string) => parseFloat(value),
  }),
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
    props: ['variant'],
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

vi.mock('@/components/BaseInput.vue', () => ({
  default: {
    name: 'BaseInput',
    template: '<input :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />',
    props: ['modelValue', 'type', 'label', 'placeholder', 'min', 'max', 'prependInnerIcon', 'class'],
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

vi.mock('@/components/TransactionDetailModal.vue', () => ({
  default: {
    name: 'TransactionDetailModal',
    template: '<div v-if="modelValue">Detail Modal</div>',
    props: ['modelValue', 'transaction'],
    emits: ['update:modelValue', 'close'],
  },
}))

describe('TransactionsView', () => {
  let router: ReturnType<typeof createRouter>
  let vuetify: ReturnType<typeof createVuetify>
  let i18n: ReturnType<typeof createI18n>
  let pinia: ReturnType<typeof createPinia>
  let systemModeStore: any

  beforeEach(() => {
    pinia = createPinia()
    setActivePinia(pinia)

    router = createRouter({
      history: createMemoryHistory(),
      routes: [{ path: '/transactions', component: TransactionsView }],
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
            transactions: {
              title: 'Transactions',
              subtitle: 'View and manage transactions',
              live: 'Live',
              pause: 'Pause',
              refresh: 'Refresh',
              loading: 'Loading...',
              errorTitle: 'Error loading transactions',
              generateSynthetic: 'Generate Synthetic',
              generateSyntheticTitle: 'Generate Synthetic Events',
              generateSyntheticDescription: 'Generate synthetic events for testing',
              selectSchema: 'Select Schema',
              eventCount: 'Event Count',
              eventCountPlaceholder: 'Enter count',
              generate: 'Generate',
              filterStatus: 'Status',
              filterStartDate: 'Start Date',
              filterEndDate: 'End Date',
              applyFilters: 'Apply Filters',
              clearFilters: 'Clear Filters',
              viewDetails: 'View Details',
              approve: 'Approve',
              reject: 'Reject',
              statusAll: 'All',
              statusApproved: 'Approved',
              statusRejected: 'Rejected',
              statusInReview: 'In Review',
              statusPending: 'Pending',
              tableStatus: 'Status',
              tableCreated: 'Created',
              tableActions: 'Actions',
            },
          },
          components: {
            modal: {
              cancel: 'Cancel',
            },
          },
        },
      },
    })

    systemModeStore = {
      syntheticMode: false,
      loadMode: vi.fn().mockResolvedValue(undefined),
    }

    vi.mocked(useSystemModeStore).mockReturnValue(systemModeStore)

    vi.clearAllMocks()
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('renders transactions view', () => {
    const wrapper = mount(TransactionsView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    expect(wrapper.text()).toContain('Transactions')
  })

  it('loads schemas on mount', async () => {
    const mockSchemas = { schemas: [] }

    vi.mocked(api.get).mockImplementation((url: string) => {
      if (url === '/api/v1/schemas') {
        return Promise.resolve(mockSchemas)
      }
      return Promise.reject(new Error('Unknown endpoint'))
    })

    const wrapper = mount(TransactionsView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    await flushPromises()

    expect(api.get).toHaveBeenCalledWith('/api/v1/schemas')
  })

  it('loads transactions when schema is selected', async () => {
    const mockSchemas = {
      schemas: [{ id: '1', name: 'Transaction' }],
    }
    const mockTransactions = {
      transactions: [],
      total: 0,
    }

    vi.mocked(api.get).mockImplementation((url: string) => {
      if (url === '/api/v1/schemas') {
        return Promise.resolve(mockSchemas)
      }
      if (url.includes('/api/v1/transactions')) {
        return Promise.resolve(mockTransactions)
      }
      return Promise.reject(new Error('Unknown endpoint'))
    })

    const wrapper = mount(TransactionsView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    await flushPromises()

    ;(wrapper.vm as any).filters.schemaId = '1'
    await (wrapper.vm as any).loadTransactions()
    await flushPromises()

    expect(api.get).toHaveBeenCalled()
  })

  it('displays loading state while loading transactions', () => {
    vi.mocked(api.get).mockImplementation(() => new Promise(() => {}))

    const wrapper = mount(TransactionsView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    expect((wrapper.vm as any).loading).toBe(true)
  })

  it('displays error when loading fails', async () => {
    const mockSchemas = {
      schemas: [
        { id: '1', name: 'Test Schema', event_type: 'test' }
      ]
    }

    vi.mocked(api.get).mockImplementation((url: string) => {
      if (url.includes('/api/v1/schemas') && !url.includes('/api/v1/schemas/')) {
        return Promise.resolve(mockSchemas)
      }
      if (url.includes('/api/v1/schemas/1')) {
        return Promise.resolve({
          id: '1',
          name: 'Test Schema',
          event_type: 'test',
          fields: [],
          extracted_fields: [],
          sample_event: {}
        })
      }
      if (url.includes('/api/v1/transactions')) {
        return Promise.reject(new Error('Failed to load'))
      }
      return Promise.reject(new Error('Unexpected API call'))
    })

    await router.push({ path: '/transactions', query: { schemaId: '1' } })

    const wrapper = mount(TransactionsView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    await flushPromises()

    expect((wrapper.vm as any).error).toBeTruthy()
    expect((wrapper.vm as any).loading).toBe(false)
  })

  it('applies filters correctly', async () => {
    const mockSchemas = { schemas: [] }
    const mockTransactions = {
      transactions: [],
      total: 0,
    }

    vi.mocked(api.get).mockImplementation((url: string) => {
      if (url === '/api/v1/schemas') {
        return Promise.resolve(mockSchemas)
      }
      if (url.includes('/api/v1/transactions')) {
        return Promise.resolve(mockTransactions)
      }
      return Promise.reject(new Error('Unknown endpoint'))
    })

    const wrapper = mount(TransactionsView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    await flushPromises()

    ;(wrapper.vm as any).filters.status = 'approved'
    ;(wrapper.vm as any).filters.schemaId = '1'
    await (wrapper.vm as any).applyFilters()
    await flushPromises()

    expect((wrapper.vm as any).currentPage).toBe(1)
  })

  it('clears filters correctly', async () => {
    const mockSchemas = { schemas: [] }

    vi.mocked(api.get).mockResolvedValue(mockSchemas)

    const wrapper = mount(TransactionsView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    await flushPromises()

    ;(wrapper.vm as any).filters.status = 'approved'
    ;(wrapper.vm as any).filters.startDate = '2024-01-01'
    ;(wrapper.vm as any).filters.endDate = '2024-01-31'
    ;(wrapper.vm as any).clearFilters()

    expect((wrapper.vm as any).filters.status).toBeNull()
    expect((wrapper.vm as any).filters.startDate).toBeNull()
    expect((wrapper.vm as any).filters.endDate).toBeNull()
  })

  it('approves transaction', async () => {
    const mockSchemas = { schemas: [] }
    const mockTransactions = {
      transactions: [{ id: '1', status: 'pending' }],
      total: 1,
    }

    vi.mocked(api.get).mockImplementation((url: string) => {
      if (url === '/api/v1/schemas') {
        return Promise.resolve(mockSchemas)
      }
      if (url.includes('/api/v1/transactions')) {
        return Promise.resolve(mockTransactions)
      }
      return Promise.reject(new Error('Unknown endpoint'))
    })

    vi.mocked(api.patch).mockResolvedValue({})

    const wrapper = mount(TransactionsView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    await flushPromises()

    ;(wrapper.vm as any).filters.schemaId = '1'
    await (wrapper.vm as any).loadTransactions()
    await flushPromises()

    await (wrapper.vm as any).approveTransaction('1')
    await flushPromises()

    expect(api.patch).toHaveBeenCalledWith('/api/v1/transactions/1/approve')
  })

  it('rejects transaction', async () => {
    const mockSchemas = { schemas: [] }
    const mockTransactions = {
      transactions: [{ id: '1', status: 'pending' }],
      total: 1,
    }

    vi.mocked(api.get).mockImplementation((url: string) => {
      if (url === '/api/v1/schemas') {
        return Promise.resolve(mockSchemas)
      }
      if (url.includes('/api/v1/transactions')) {
        return Promise.resolve(mockTransactions)
      }
      return Promise.reject(new Error('Unknown endpoint'))
    })

    vi.mocked(api.patch).mockResolvedValue({})

    const wrapper = mount(TransactionsView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    await flushPromises()

    ;(wrapper.vm as any).filters.schemaId = '1'
    await (wrapper.vm as any).loadTransactions()
    await flushPromises()

    await (wrapper.vm as any).rejectTransaction('1')
    await flushPromises()

    expect(api.patch).toHaveBeenCalledWith('/api/v1/transactions/1/reject')
  })

  it('opens detail modal when view button is clicked', async () => {
    const mockSchemas = { schemas: [] }
    const transaction = {
      id: '1',
      status: 'approved',
      created_at: '2024-01-01T00:00:00Z',
      metadata: {},
    }

    vi.mocked(api.get).mockResolvedValue(mockSchemas)

    const wrapper = mount(TransactionsView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    await flushPromises()

    ;(wrapper.vm as any).openDetailModal(transaction)

    expect((wrapper.vm as any).showDetailModal).toBe(true)
    expect((wrapper.vm as any).selectedTransaction).toEqual(transaction)
  })

  it('formats date correctly', () => {
    const wrapper = mount(TransactionsView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    const formatted = (wrapper.vm as any).formatDate('2024-01-15T12:30:00Z')
    expect(formatted).toBeTruthy()
  })

  it('gets status variant correctly', () => {
    const wrapper = mount(TransactionsView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    expect((wrapper.vm as any).getStatusVariant('approved')).toBe('success')
    expect((wrapper.vm as any).getStatusVariant('rejected')).toBe('danger')
    expect((wrapper.vm as any).getStatusVariant('in_review')).toBe('warning')
    expect((wrapper.vm as any).getStatusVariant('pending')).toBe('info')
  })

  it('calls stopLiveUpdates on unmount', async () => {
    const mockSchemas = { schemas: [] }
    const mockTransactions = {
      transactions: [],
      total: 0,
    }

    vi.mocked(api.get).mockImplementation((url: string) => {
      if (url.includes('/api/v1/schemas')) {
        return Promise.resolve(mockSchemas)
      }
      return Promise.resolve(mockTransactions)
    })

    const wrapper = mount(TransactionsView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    await flushPromises()

    const apiCallCountBefore = vi.mocked(api.get).mock.calls.length

    wrapper.unmount()
    await flushPromises()

    vi.advanceTimersByTime(5000)
    await flushPromises()

    const apiCallCountAfter = vi.mocked(api.get).mock.calls.length

    expect(apiCallCountAfter).toBe(apiCallCountBefore)
  })
})
