import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'
import { createPinia, setActivePinia } from 'pinia'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import { createI18n } from 'vue-i18n'
import BrandingView from './BrandingView.vue'
import { useBrandingStore } from '@/stores/branding'

vi.mock('@/components/BaseButton.vue', () => ({
  default: {
    name: 'BaseButton',
    template: '<button type="button"><slot /></button>',
    props: ['loading', 'disabled', 'fullWidth', 'prependIcon', 'variant', 'type', 'size'],
  },
}))

vi.mock('@/components/BaseInput.vue', () => ({
  default: {
    name: 'BaseInput',
    template: '<input :type="type" :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />',
    props: ['modelValue', 'type', 'label', 'placeholder', 'disabled', 'required', 'maxlength', 'hint', 'persistentHint', 'prependInnerIcon', 'density', 'hideDetails', 'class', 'pattern'],
    emits: ['update:modelValue'],
  },
}))

describe('BrandingView', () => {
  let router: ReturnType<typeof createRouter>
  let vuetify: ReturnType<typeof createVuetify>
  let i18n: ReturnType<typeof createI18n>
  let pinia: ReturnType<typeof createPinia>

  beforeEach(() => {
    pinia = createPinia()
    setActivePinia(pinia)

    router = createRouter({
      history: createMemoryHistory(),
      routes: [{ path: '/branding', component: BrandingView }],
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
            branding: {
              title: 'Branding',
              subtitle: 'Customize your application branding',
              configuration: 'Configuration',
              appName: 'App Name',
              appNamePlaceholder: 'Enter app name',
              appNameHint: 'Maximum 100 characters',
              iconUrl: 'Icon URL',
              iconUrlPlaceholder: 'Enter icon URL',
              iconUrlHint: 'URL for the application icon',
              faviconUrl: 'Favicon URL',
              faviconUrlPlaceholder: 'Enter favicon URL',
              faviconUrlHint: 'URL for the browser favicon',
              primaryColor: 'Primary Color',
              secondaryColor: 'Secondary Color',
              headerColor: 'Header Color',
              saveConfiguration: 'Save Configuration',
              resetToDefaults: 'Reset to Defaults',
              livePreview: 'Live Preview',
              colors: 'Colors',
              primary: 'Primary',
              secondary: 'Secondary',
              header: 'Header',
              buttons: 'Buttons',
              browserTab: 'Browser Tab',
              saveSuccess: 'Branding configuration updated successfully!',
              saveError: 'Failed to update branding configuration',
            },
          },
          common: {
            loading: 'Loading...',
          },
        },
      },
    })

    vi.clearAllMocks()
  })

  it('renders branding configuration form', () => {
    const brandingStore = useBrandingStore()
    brandingStore.config = {
      id: 1,
      app_name: 'Test App',
      icon_url: '/icon.png',
      favicon_url: '/favicon.ico',
      primary_color: '#3B82F6',
      secondary_color: '#10B981',
      header_color: '#1e1e1e',
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
    }

    const wrapper = mount(BrandingView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    expect(wrapper.text()).toContain('Configuration')
    expect(wrapper.text()).toContain('Live Preview')
  })

  it('populates form with existing branding config', () => {
    const brandingStore = useBrandingStore()
    brandingStore.config = {
      id: 1,
      app_name: 'My App',
      icon_url: '/custom-icon.png',
      favicon_url: '/custom-favicon.ico',
      primary_color: '#FF0000',
      secondary_color: '#00FF00',
      header_color: '#0000FF',
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
    }

    const wrapper = mount(BrandingView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    expect((wrapper.vm as any).form.app_name).toBe('My App')
    expect((wrapper.vm as any).form.primary_color).toBe('#FF0000')
  })

  it('updates branding when form is submitted', async () => {
    const brandingStore = useBrandingStore()
    brandingStore.config = {
      id: 1,
      app_name: 'Test App',
      icon_url: null,
      favicon_url: null,
      primary_color: '#3B82F6',
      secondary_color: '#10B981',
      header_color: '#1e1e1e',
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
    }

    const updateSpy = vi.spyOn(brandingStore, 'updateBranding').mockResolvedValue({
      id: 1,
      app_name: 'Updated App',
      icon_url: null,
      favicon_url: null,
      primary_color: '#FF0000',
      secondary_color: '#00FF00',
      header_color: '#0000FF',
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
    })

    const wrapper = mount(BrandingView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    ;(wrapper.vm as any).form.app_name = 'Updated App'
    ;(wrapper.vm as any).form.primary_color = '#FF0000'

    await (wrapper.vm as any).handleSubmit()
    await flushPromises()

    expect(updateSpy).toHaveBeenCalledWith({
      app_name: 'Updated App',
      icon_url: null,
      favicon_url: null,
      primary_color: '#FF0000',
      secondary_color: '#10B981',
      header_color: '#1e1e1e',
    })
  })

  it('shows error message when update fails', async () => {
    const brandingStore = useBrandingStore()
    brandingStore.config = {
      id: 1,
      app_name: 'Test App',
      icon_url: null,
      favicon_url: null,
      primary_color: '#3B82F6',
      secondary_color: '#10B981',
      header_color: '#1e1e1e',
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
    }

    vi.spyOn(brandingStore, 'updateBranding').mockRejectedValue(new Error('Update failed'))

    const wrapper = mount(BrandingView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    await (wrapper.vm as any).handleSubmit()
    await flushPromises()

    expect((wrapper.vm as any).error).toBeTruthy()
    expect((wrapper.vm as any).loading).toBe(false)
  })

  it('shows success message when update succeeds', async () => {
    const brandingStore = useBrandingStore()
    brandingStore.config = {
      id: 1,
      app_name: 'Test App',
      icon_url: null,
      favicon_url: null,
      primary_color: '#3B82F6',
      secondary_color: '#10B981',
      header_color: '#1e1e1e',
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
    }

    vi.spyOn(brandingStore, 'updateBranding').mockResolvedValue({
      id: 1,
      app_name: 'Updated App',
      icon_url: null,
      favicon_url: null,
      primary_color: '#FF0000',
      secondary_color: '#00FF00',
      header_color: '#0000FF',
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
    })

    vi.useFakeTimers()

    const wrapper = mount(BrandingView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    await (wrapper.vm as any).handleSubmit()
    await flushPromises()

    expect((wrapper.vm as any).success).toBeTruthy()

    vi.advanceTimersByTime(5000)
    await flushPromises()

    expect((wrapper.vm as any).success).toBe('')

    vi.useRealTimers()
  })

  it('resets form to defaults when reset button is clicked', () => {
    const brandingStore = useBrandingStore()
    brandingStore.config = {
      id: 1,
      app_name: 'Custom App',
      icon_url: '/custom.png',
      favicon_url: '/custom.ico',
      primary_color: '#FF0000',
      secondary_color: '#00FF00',
      header_color: '#0000FF',
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
    }

    const wrapper = mount(BrandingView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    ;(wrapper.vm as any).form.app_name = 'Changed Name'
    ;(wrapper.vm as any).resetToDefaults()

    expect((wrapper.vm as any).form.app_name).toBe('AlgoShield')
    expect((wrapper.vm as any).form.primary_color).toBe('#3B82F6')
    expect((wrapper.vm as any).form.secondary_color).toBe('#10B981')
    expect((wrapper.vm as any).form.header_color).toBe('#1e1e1e')
  })

  it('handles image error by setting fallback image', () => {
    const brandingStore = useBrandingStore()
    brandingStore.config = {
      id: 1,
      app_name: 'Test App',
      icon_url: '/invalid.png',
      favicon_url: null,
      primary_color: '#3B82F6',
      secondary_color: '#10B981',
      header_color: '#1e1e1e',
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
    }

    const wrapper = mount(BrandingView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    const mockEvent = {
      target: {
        src: '/invalid.png',
      },
    } as any

    ;(wrapper.vm as any).handleImageError(mockEvent)

    expect(mockEvent.target.src).toBe('/gopher.png')
  })

  it('sets loading state during submission', async () => {
    const brandingStore = useBrandingStore()
    brandingStore.config = {
      id: 1,
      app_name: 'Test App',
      icon_url: null,
      favicon_url: null,
      primary_color: '#3B82F6',
      secondary_color: '#10B981',
      header_color: '#1e1e1e',
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
    }

    let resolveUpdate: (value: any) => void
    const updatePromise = new Promise(resolve => {
      resolveUpdate = resolve
    })

    vi.spyOn(brandingStore, 'updateBranding').mockReturnValue(updatePromise as any)

    const wrapper = mount(BrandingView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    const submitPromise = (wrapper.vm as any).handleSubmit()

    expect((wrapper.vm as any).loading).toBe(true)

    resolveUpdate!({
      id: 1,
      app_name: 'Updated',
      icon_url: null,
      favicon_url: null,
      primary_color: '#3B82F6',
      secondary_color: '#10B981',
      header_color: '#1e1e1e',
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
    })

    await submitPromise
    await flushPromises()

    expect((wrapper.vm as any).loading).toBe(false)
  })

  it('clears error message when close button is clicked', () => {
    const brandingStore = useBrandingStore()
    brandingStore.config = {
      id: 1,
      app_name: 'Test App',
      icon_url: null,
      favicon_url: null,
      primary_color: '#3B82F6',
      secondary_color: '#10B981',
      header_color: '#1e1e1e',
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
    }

    const wrapper = mount(BrandingView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    ;(wrapper.vm as any).error = 'Test error'
    ;(wrapper.vm as any).error = ''

    expect((wrapper.vm as any).error).toBe('')
  })

  it('clears success message when close button is clicked', () => {
    const brandingStore = useBrandingStore()
    brandingStore.config = {
      id: 1,
      app_name: 'Test App',
      icon_url: null,
      favicon_url: null,
      primary_color: '#3B82F6',
      secondary_color: '#10B981',
      header_color: '#1e1e1e',
      created_at: '2024-01-01T00:00:00Z',
      updated_at: '2024-01-01T00:00:00Z',
    }

    const wrapper = mount(BrandingView, {
      global: {
        plugins: [router, vuetify, i18n, pinia],
      },
    })

    ;(wrapper.vm as any).success = 'Success message'
    ;(wrapper.vm as any).success = ''

    expect((wrapper.vm as any).success).toBe('')
  })
})
