import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'
import { createPinia, setActivePinia } from 'pinia'
import ProtectedRoute from './ProtectedRoute.vue'
import { useAuthStore } from '@/stores/auth'

// Stub LoadingSpinner so we don't depend on Vuetify's v-progress-circular in tests.
// The real component uses v-progress-circular, which isn't resolved without a full Vuetify mount.
vi.mock('@/components/LoadingSpinner.vue', () => ({
  default: {
    name: 'LoadingSpinner',
    template: '<div class="animate-spin">Loading</div>',
    props: ['text', 'size', 'centered', 'fullscreen'],
  },
}))

describe('ProtectedRoute', () => {
  let router: ReturnType<typeof createRouter>
  let pinia: ReturnType<typeof createPinia>

  beforeEach(() => {
    pinia = createPinia()
    setActivePinia(pinia)

    router = createRouter({
      history: createMemoryHistory(),
      routes: [
        { path: '/public', component: { template: '<div>Public</div>' }, meta: { public: true } },
        { path: '/protected', component: { template: '<div>Protected</div>' } },
      ],
    })

    vi.clearAllMocks()
  })

  it('renders slot content when not loading', async () => {
    await router.push('/protected')
    await router.isReady()

    const authStore = useAuthStore()
    authStore.loading = false

    const wrapper = mount(ProtectedRoute, {
      slots: {
        default: '<div class="test-content">Content</div>',
      },
      global: {
        plugins: [router, pinia],
      },
    })

    expect(wrapper.find('.test-content').exists()).toBe(true)
    expect(wrapper.text()).toContain('Content')
  })

  it('shows loading spinner when loading on protected route', async () => {
    await router.push('/protected')
    await router.isReady()

    const authStore = useAuthStore()
    authStore.loading = true

    const wrapper = mount(ProtectedRoute, {
      slots: {
        default: '<div class="test-content">Content</div>',
      },
      global: {
        plugins: [router, pinia],
      },
    })

    expect(wrapper.find('.animate-spin').exists()).toBe(true)
    expect(wrapper.find('.test-content').exists()).toBe(false)
  })

  it('does not show loading spinner on public route even when loading', async () => {
    await router.push('/public')
    await router.isReady()

    const authStore = useAuthStore()
    authStore.loading = true

    const wrapper = mount(ProtectedRoute, {
      slots: {
        default: '<div class="test-content">Content</div>',
      },
      global: {
        plugins: [router, pinia],
      },
    })

    expect(wrapper.find('.animate-spin').exists()).toBe(false)
    expect(wrapper.find('.test-content').exists()).toBe(true)
  })

  it('renders slot when loading is false on protected route', async () => {
    await router.push('/protected')
    await router.isReady()

    const authStore = useAuthStore()
    authStore.loading = false

    const wrapper = mount(ProtectedRoute, {
      slots: {
        default: '<div class="test-content">Content</div>',
      },
      global: {
        plugins: [router, pinia],
      },
    })

    expect(wrapper.find('.animate-spin').exists()).toBe(false)
    expect(wrapper.find('.test-content').exists()).toBe(true)
  })

  it('renders slot immediately on public route regardless of loading state', async () => {
    await router.push('/public')
    await router.isReady()

    const authStore = useAuthStore()
    authStore.loading = false

    const wrapper = mount(ProtectedRoute, {
      slots: {
        default: '<div class="test-content">Content</div>',
      },
      global: {
        plugins: [router, pinia],
      },
    })

    expect(wrapper.find('.test-content').exists()).toBe(true)
    expect(wrapper.text()).toContain('Content')
  })

  it('correctly identifies public routes', async () => {
    await router.push('/public')
    await router.isReady()

    const wrapper = mount(ProtectedRoute, {
      slots: {
        default: '<div>Content</div>',
      },
      global: {
        plugins: [router, pinia],
      },
    })

    const vm = wrapper.vm as any
    expect(vm.isPublicRoute).toBe(true)
  })

  it('correctly identifies protected routes', async () => {
    await router.push('/protected')
    await router.isReady()

    const wrapper = mount(ProtectedRoute, {
      slots: {
        default: '<div>Content</div>',
      },
      global: {
        plugins: [router, pinia],
      },
    })

    const vm = wrapper.vm as any
    expect(vm.isPublicRoute).toBe(false)
  })

  it('reactively updates when auth loading state changes', async () => {
    await router.push('/protected')
    await router.isReady()

    const authStore = useAuthStore()
    authStore.loading = true

    const wrapper = mount(ProtectedRoute, {
      slots: {
        default: '<div class="test-content">Content</div>',
      },
      global: {
        plugins: [router, pinia],
      },
    })

    expect(wrapper.find('.animate-spin').exists()).toBe(true)

    authStore.loading = false
    await wrapper.vm.$nextTick()

    expect(wrapper.find('.animate-spin').exists()).toBe(false)
    expect(wrapper.find('.test-content').exists()).toBe(true)
  })
})
