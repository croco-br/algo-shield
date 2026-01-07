import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import LoadingSpinner from './LoadingSpinner.vue'

describe('LoadingSpinner', () => {
  const createWrapper = (props = {}) => {
    return mount(LoadingSpinner, {
      props,
      global: {
        stubs: {
          VProgressCircular: {
            template: '<div class="v-progress-circular"></div>',
            props: ['size', 'width', 'color', 'indeterminate'],
          },
        },
      },
    })
  }

  describe('rendering', () => {
    it('renders progress circular component', () => {
      const wrapper = createWrapper()

      expect(wrapper.find('.v-progress-circular').exists()).toBe(true)
    })

    it('renders without text by default', () => {
      const wrapper = createWrapper()

      expect(wrapper.find('p').exists()).toBe(false)
    })

    it('renders with text when provided', () => {
      const wrapper = createWrapper({ text: 'Loading...' })

      expect(wrapper.find('p').exists()).toBe(true)
      expect(wrapper.text()).toContain('Loading...')
    })

    it('renders centered by default', () => {
      const wrapper = createWrapper()

      expect(wrapper.vm.containerClasses).toContain('justify-center')
    })
  })

  describe('size prop', () => {
    it('uses medium size by default', () => {
      const wrapper = createWrapper()

      expect(wrapper.vm.mappedSize).toBe(40)
      expect(wrapper.vm.mappedWidth).toBe(4)
    })

    it('renders with small size', () => {
      const wrapper = createWrapper({ size: 'sm' })

      expect(wrapper.vm.mappedSize).toBe(24)
      expect(wrapper.vm.mappedWidth).toBe(2)
    })

    it('renders with medium size', () => {
      const wrapper = createWrapper({ size: 'md' })

      expect(wrapper.vm.mappedSize).toBe(40)
      expect(wrapper.vm.mappedWidth).toBe(4)
    })

    it('renders with large size', () => {
      const wrapper = createWrapper({ size: 'lg' })

      expect(wrapper.vm.mappedSize).toBe(48)
      expect(wrapper.vm.mappedWidth).toBe(4)
    })
  })

  describe('centered prop', () => {
    it('is centered by default', () => {
      const wrapper = createWrapper()

      expect(wrapper.vm.containerClasses).toContain('min-h-[50vh]')
      expect(wrapper.vm.containerClasses).toContain('justify-center')
    })

    it('is centered when centered is true', () => {
      const wrapper = createWrapper({ centered: true })

      expect(wrapper.vm.containerClasses).toContain('min-h-[50vh]')
      expect(wrapper.vm.containerClasses).toContain('justify-center')
    })

    it('is not centered when centered is false', () => {
      const wrapper = createWrapper({ centered: false })

      expect(wrapper.vm.containerClasses).not.toContain('min-h-[50vh]')
      expect(wrapper.vm.containerClasses).not.toContain('justify-center')
    })
  })

  describe('fullscreen prop', () => {
    it('is not fullscreen by default', () => {
      const wrapper = createWrapper()

      expect(wrapper.vm.containerClasses).not.toContain('min-h-screen')
    })

    it('is fullscreen when fullscreen is true', () => {
      const wrapper = createWrapper({ fullscreen: true })

      expect(wrapper.vm.containerClasses).toContain('min-h-screen')
      expect(wrapper.vm.containerClasses).toContain('justify-center')
    })

    it('overrides centered when fullscreen is true', () => {
      const wrapper = createWrapper({ fullscreen: true, centered: true })

      expect(wrapper.vm.containerClasses).toContain('min-h-screen')
      expect(wrapper.vm.containerClasses).not.toContain('min-h-[50vh]')
    })

    it('is not fullscreen when fullscreen is false', () => {
      const wrapper = createWrapper({ fullscreen: false })

      expect(wrapper.vm.containerClasses).not.toContain('min-h-screen')
    })
  })

  describe('text prop', () => {
    it('does not render text element when text is not provided', () => {
      const wrapper = createWrapper()

      expect(wrapper.find('p').exists()).toBe(false)
    })

    it('renders text when text prop is provided', () => {
      const wrapper = createWrapper({ text: 'Please wait...' })

      expect(wrapper.find('p').exists()).toBe(true)
      expect(wrapper.text()).toContain('Please wait...')
    })

    it('renders different text messages', () => {
      const wrapper1 = createWrapper({ text: 'Loading data...' })
      expect(wrapper1.text()).toContain('Loading data...')

      const wrapper2 = createWrapper({ text: 'Processing...' })
      expect(wrapper2.text()).toContain('Processing...')
    })
  })

  describe('computed properties', () => {
    describe('containerClasses', () => {
      it('includes base classes', () => {
        const wrapper = createWrapper()

        expect(wrapper.vm.containerClasses).toContain('d-flex')
        expect(wrapper.vm.containerClasses).toContain('flex-column')
        expect(wrapper.vm.containerClasses).toContain('align-center')
        expect(wrapper.vm.containerClasses).toContain('gap-3')
      })

      it('adds fullscreen classes when fullscreen is true', () => {
        const wrapper = createWrapper({ fullscreen: true })

        expect(wrapper.vm.containerClasses).toContain('min-h-screen')
        expect(wrapper.vm.containerClasses).toContain('justify-center')
      })

      it('adds centered classes when centered is true and fullscreen is false', () => {
        const wrapper = createWrapper({ centered: true, fullscreen: false })

        expect(wrapper.vm.containerClasses).toContain('min-h-[50vh]')
        expect(wrapper.vm.containerClasses).toContain('justify-center')
      })

      it('does not add height classes when both centered and fullscreen are false', () => {
        const wrapper = createWrapper({ centered: false, fullscreen: false })

        expect(wrapper.vm.containerClasses).not.toContain('min-h-screen')
        expect(wrapper.vm.containerClasses).not.toContain('min-h-[50vh]')
      })
    })

    describe('mappedSize', () => {
      it('maps sm size to 24', () => {
        const wrapper = createWrapper({ size: 'sm' })

        expect(wrapper.vm.mappedSize).toBe(24)
      })

      it('maps md size to 40', () => {
        const wrapper = createWrapper({ size: 'md' })

        expect(wrapper.vm.mappedSize).toBe(40)
      })

      it('maps lg size to 48', () => {
        const wrapper = createWrapper({ size: 'lg' })

        expect(wrapper.vm.mappedSize).toBe(48)
      })

      it('defaults to 40 for md', () => {
        const wrapper = createWrapper()

        expect(wrapper.vm.mappedSize).toBe(40)
      })
    })

    describe('mappedWidth', () => {
      it('maps sm size to width 2', () => {
        const wrapper = createWrapper({ size: 'sm' })

        expect(wrapper.vm.mappedWidth).toBe(2)
      })

      it('maps md size to width 4', () => {
        const wrapper = createWrapper({ size: 'md' })

        expect(wrapper.vm.mappedWidth).toBe(4)
      })

      it('maps lg size to width 4', () => {
        const wrapper = createWrapper({ size: 'lg' })

        expect(wrapper.vm.mappedWidth).toBe(4)
      })

      it('defaults to width 4 for md', () => {
        const wrapper = createWrapper()

        expect(wrapper.vm.mappedWidth).toBe(4)
      })
    })

    describe('textClasses', () => {
      it('includes base text color class', () => {
        const wrapper = createWrapper({ text: 'Loading' })

        expect(wrapper.vm.textClasses).toContain('text-grey-darken-1')
      })

      it('adds text-sm class for small size', () => {
        const wrapper = createWrapper({ size: 'sm', text: 'Loading' })

        expect(wrapper.vm.textClasses).toContain('text-sm')
      })

      it('adds text-base class for medium size', () => {
        const wrapper = createWrapper({ size: 'md', text: 'Loading' })

        expect(wrapper.vm.textClasses).toContain('text-base')
      })

      it('adds text-lg class for large size', () => {
        const wrapper = createWrapper({ size: 'lg', text: 'Loading' })

        expect(wrapper.vm.textClasses).toContain('text-lg')
      })
    })
  })

  describe('combined props', () => {
    it('renders with all props combined', () => {
      const wrapper = createWrapper({
        size: 'lg',
        text: 'Loading data...',
        centered: true,
        fullscreen: false,
      })

      expect(wrapper.text()).toContain('Loading data...')
      expect(wrapper.vm.mappedSize).toBe(48)
      expect(wrapper.vm.containerClasses).toContain('min-h-[50vh]')
      expect(wrapper.vm.textClasses).toContain('text-lg')
    })

    it('renders fullscreen with large size and text', () => {
      const wrapper = createWrapper({
        size: 'lg',
        text: 'Processing...',
        fullscreen: true,
      })

      expect(wrapper.text()).toContain('Processing...')
      expect(wrapper.vm.mappedSize).toBe(48)
      expect(wrapper.vm.containerClasses).toContain('min-h-screen')
      expect(wrapper.vm.textClasses).toContain('text-lg')
    })

    it('renders small spinner without centering', () => {
      const wrapper = createWrapper({
        size: 'sm',
        centered: false,
        fullscreen: false,
      })

      expect(wrapper.vm.mappedSize).toBe(24)
      expect(wrapper.vm.containerClasses).not.toContain('min-h-screen')
      expect(wrapper.vm.containerClasses).not.toContain('min-h-[50vh]')
    })
  })

  describe('progress circular attributes', () => {
    it('sets indeterminate attribute', () => {
      const wrapper = createWrapper()

      expect(wrapper.find('.v-progress-circular').exists()).toBe(true)
    })

    it('sets primary color', () => {
      const wrapper = createWrapper()

      expect(wrapper.find('.v-progress-circular').exists()).toBe(true)
    })
  })
})

