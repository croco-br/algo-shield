import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'
import { createPinia, setActivePinia } from 'pinia'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import { createI18n } from 'vue-i18n'
import DashboardView from './DashboardView.vue'
import { api } from '@/lib/api'

vi.mock('@/lib/api', () => ({
  api: {
    get: vi.fn(),
  },
  setTokenGetter: vi.fn(),
  setCsrfTokenGetter: vi.fn(),
  setSyntheticModeStorage: vi.fn(),
}))

vi.mock('@/composables/useLocale', () => ({
  useLocale: () => ({
    t: (key: string) => key,
  }),
}))

vi.mock('@/components/BaseButton.vue', () => ({
  default: {
    name: 'BaseButton',
    template: '<button type="button"><slot /></button>',
    props: ['loading', 'disabled', 'prependIcon', 'variant', 'size'],
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

describe('DashboardView', () => {
  let router: ReturnType<typeof createRouter>
  let vuetify: ReturnType<typeof createVuetify>
  let i18n: ReturnType<typeof createI18n>
  let pinia: ReturnType<typeof createPinia>

  beforeEach(() => {
    pinia = createPinia()
    setActivePinia(pinia)

    router = createRouter({
      history: createMemoryHistory(),
      routes: [{ path: '/dashboard', component: DashboardView }],
    })

    vuetify = createVuetify({
      components,
      directives,
    })

    i18n = createI18n({
      legacy: false,
      locale: 'en-US',
      messages: {
        'en-US': {
          views: {
            dashboard: {
              title: 'Dashboard',
              subtitle: 'Overview of transaction metrics',
              refresh: 'Refresh',
              loading: 'Loading metrics...',
              errorTitle: 'Error loading dashboard',
              errorLoad: 'Failed to load metrics',
              updated: 'Updated',
              totalTransactions: 'Total Transactions',
              approved: 'Approved',
              rejected: 'Rejected',
              inReview: 'In Review',
              pending: 'Pending',
              statusDistribution: 'Status Distribution',
              transactionVolume: 'Transaction Volume',
              justNow: 'just now',
              minutesAgo: '{count}m ago',
              hoursAgo: '{count}h ago',
            },
          },
        },
      },
    })

    vi.clearAllMocks()
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('renders dashboard view', () => {
    const wrapper = mount(DashboardView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    expect(wrapper.text()).toContain('Dashboard')
  })

  it('loads metrics on mount', async () => {
    const mockMetrics = {
      data: {
        status_distribution: [
          { status: 'approved', count: 100 },
          { status: 'rejected', count: 50 },
        ],
        temporal_24h: [{ bucket: '2024-01-01T00:00:00Z', count: 10 }],
        temporal_7d: [],
        temporal_30d: [],
        total_count: 150,
        cached_at: '2024-01-01T00:00:00Z',
      },
      response_time_ms: 50,
    }

    vi.mocked(api.get).mockResolvedValue(mockMetrics)

    const wrapper = mount(DashboardView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    await flushPromises()

    expect(api.get).toHaveBeenCalledWith('/api/v1/dashboard/metrics')
  })

  it('displays loading spinner while loading', () => {
    vi.mocked(api.get).mockImplementation(() => new Promise(() => {}))

    const wrapper = mount(DashboardView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    expect((wrapper.vm as any).loading).toBe(true)
  })

  it('displays error message when loading fails', async () => {
    vi.mocked(api.get).mockRejectedValue(new Error('Failed to load'))

    const wrapper = mount(DashboardView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    await flushPromises()

    expect((wrapper.vm as any).error).toBeTruthy()
    expect((wrapper.vm as any).loading).toBe(false)
  })

  it('displays metrics when loaded successfully', async () => {
    const mockMetrics = {
      data: {
        status_distribution: [
          { status: 'approved', count: 100 },
          { status: 'rejected', count: 50 },
        ],
        temporal_24h: [{ bucket: '2024-01-01T00:00:00Z', count: 10 }],
        temporal_7d: [],
        temporal_30d: [],
        total_count: 150,
        cached_at: '2024-01-01T00:00:00Z',
      },
      response_time_ms: 50,
    }

    vi.mocked(api.get).mockResolvedValue(mockMetrics)

    const wrapper = mount(DashboardView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    await flushPromises()

    expect((wrapper.vm as any).metrics).toBeTruthy()
    expect((wrapper.vm as any).metrics?.total_count).toBe(150)
  })

  it('refreshes metrics when refresh button is clicked', async () => {
    const mockMetrics = {
      data: {
        status_distribution: [],
        temporal_24h: [],
        temporal_7d: [],
        temporal_30d: [],
        total_count: 0,
        cached_at: '2024-01-01T00:00:00Z',
      },
      response_time_ms: 50,
    }

    vi.mocked(api.get).mockResolvedValue(mockMetrics)

    const wrapper = mount(DashboardView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    await flushPromises()

    vi.mocked(api.get).mockClear()

    await (wrapper.vm as any).loadMetrics()
    await flushPromises()

    expect(api.get).toHaveBeenCalledWith('/api/v1/dashboard/metrics')
  })

  it('auto-refreshes metrics every 30 seconds', async () => {
    const mockMetrics = {
      data: {
        status_distribution: [],
        temporal_24h: [],
        temporal_7d: [],
        temporal_30d: [],
        total_count: 0,
        cached_at: '2024-01-01T00:00:00Z',
      },
      response_time_ms: 50,
    }

    vi.mocked(api.get).mockResolvedValue(mockMetrics)

    const wrapper = mount(DashboardView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    await flushPromises()

    vi.mocked(api.get).mockClear()

    vi.advanceTimersByTime(30000)
    await flushPromises()

    expect(api.get).toHaveBeenCalled()
  })

  it('cleans up interval on unmount', async () => {
    const mockMetrics = {
      data: {
        status_distribution: [],
        temporal_24h: [],
        temporal_7d: [],
        temporal_30d: [],
        total_count: 0,
        cached_at: '2024-01-01T00:00:00Z',
      },
      response_time_ms: 50,
    }

    vi.mocked(api.get).mockResolvedValue(mockMetrics)

    const wrapper = mount(DashboardView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    await flushPromises()

    const callCountBefore = vi.mocked(api.get).mock.calls.length

    wrapper.unmount()

    vi.advanceTimersByTime(30000)
    await flushPromises()

    const callCountAfter = vi.mocked(api.get).mock.calls.length

    expect(callCountAfter).toBe(callCountBefore)
  })

  it('calculates percentage correctly', async () => {
    const mockMetrics = {
      data: {
        status_distribution: [
          { status: 'approved', count: 75 },
          { status: 'rejected', count: 25 },
        ],
        temporal_24h: [],
        temporal_7d: [],
        temporal_30d: [],
        total_count: 100,
        cached_at: '2024-01-01T00:00:00Z',
      },
      response_time_ms: 50,
    }

    vi.mocked(api.get).mockResolvedValue(mockMetrics)

    const wrapper = mount(DashboardView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    await flushPromises()

    expect((wrapper.vm as any).getPercentage(75)).toBe(75)
    expect((wrapper.vm as any).getPercentage(25)).toBe(25)
  })

  it('returns 0 percentage when total count is 0', async () => {
    const mockMetrics = {
      data: {
        status_distribution: [],
        temporal_24h: [],
        temporal_7d: [],
        temporal_30d: [],
        total_count: 0,
        cached_at: '2024-01-01T00:00:00Z',
      },
      response_time_ms: 50,
    }

    vi.mocked(api.get).mockResolvedValue(mockMetrics)

    const wrapper = mount(DashboardView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    await flushPromises()

    expect((wrapper.vm as any).getPercentage(10)).toBe(0)
  })

  it('formats numbers correctly', () => {
    const wrapper = mount(DashboardView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    expect((wrapper.vm as any).formatNumber(1000)).toBe('1.0K')
    expect((wrapper.vm as any).formatNumber(1000000)).toBe('1.0M')
    expect((wrapper.vm as any).formatNumber(500)).toBe('500')
  })

  it('formats time ago correctly', () => {
    const wrapper = mount(DashboardView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    const now = new Date()
    const oneMinuteAgo = new Date(now.getTime() - 60000)
    const oneHourAgo = new Date(now.getTime() - 3600000)

    expect((wrapper.vm as any).formatTimeAgo(now)).toBe('just now')
    expect((wrapper.vm as any).formatTimeAgo(oneMinuteAgo)).toContain('m ago')
    expect((wrapper.vm as any).formatTimeAgo(oneHourAgo)).toContain('h ago')
  })

  it('switches temporal period correctly', async () => {
    const mockMetrics = {
      data: {
        status_distribution: [],
        temporal_24h: [{ bucket: '2024-01-01T00:00:00Z', count: 10 }],
        temporal_7d: [{ bucket: '2024-01-01T00:00:00Z', count: 20 }],
        temporal_30d: [{ bucket: '2024-01-01T00:00:00Z', count: 30 }],
        total_count: 0,
        cached_at: '2024-01-01T00:00:00Z',
      },
      response_time_ms: 50,
    }

    vi.mocked(api.get).mockResolvedValue(mockMetrics)

    const wrapper = mount(DashboardView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    await flushPromises()

    ;(wrapper.vm as any).selectedPeriod = '7d'
    await flushPromises()

    expect((wrapper.vm as any).temporalData.length).toBe(1)
    expect((wrapper.vm as any).temporalData[0].count).toBe(20)
  })
})
