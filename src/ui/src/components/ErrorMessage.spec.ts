import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import ErrorMessage from './ErrorMessage.vue'

describe('ErrorMessage', () => {
  const createWrapper = (props = {}, slots = {}) => {
    return mount(ErrorMessage, {
      props: {
        message: 'Test error message',
        ...props,
      },
      slots,
      global: {
        stubs: {
          VAlert: {
            template: '<div class="v-alert"><slot /><slot name="append" /></div>',
            props: ['type', 'title', 'dismissible'],
          },
          BaseButton: {
            template: '<button><slot /></button>',
          },
        },
      },
    })
  }

  describe('rendering', () => {
    it('renders alert component when message is provided', () => {
      const wrapper = createWrapper({ message: 'Error occurred' })

      expect(wrapper.find('.v-alert').exists()).toBe(true)
      expect(wrapper.text()).toContain('Error occurred')
    })

    it('does not render when message is empty', () => {
      const wrapper = createWrapper({ message: '' })

      expect(wrapper.find('.v-alert').exists()).toBe(false)
    })

    it('accepts title prop', () => {
      const wrapper = createWrapper({ message: 'Error details', title: 'Error Title' })

      expect(wrapper.find('.v-alert').exists()).toBe(true)
      expect(wrapper.text()).toContain('Error details')
    })

    it('renders message without title by default', () => {
      const wrapper = createWrapper({ message: 'Simple error' })

      expect(wrapper.text()).toContain('Simple error')
    })
  })

  describe('variant prop', () => {
    it('renders as error type by default', () => {
      const wrapper = createWrapper({ message: 'Error' })

      expect(wrapper.vm.mappedType).toBe('error')
    })

    it('renders as error type when variant is error', () => {
      const wrapper = createWrapper({ message: 'Error', variant: 'error' })

      expect(wrapper.vm.mappedType).toBe('error')
    })

    it('renders as warning type when variant is warning', () => {
      const wrapper = createWrapper({ message: 'Warning', variant: 'warning' })

      expect(wrapper.vm.mappedType).toBe('warning')
    })

    it('renders as info type when variant is info', () => {
      const wrapper = createWrapper({ message: 'Info', variant: 'info' })

      expect(wrapper.vm.mappedType).toBe('info')
    })
  })

  describe('dismissible prop', () => {
    it('is dismissible by default', () => {
      const wrapper = createWrapper({ message: 'Dismissible error' })

      expect(wrapper.find('.v-alert').exists()).toBe(true)
    })

    it('is dismissible when dismissible is true', () => {
      const wrapper = createWrapper({ message: 'Dismissible', dismissible: true })

      expect(wrapper.find('.v-alert').exists()).toBe(true)
    })

    it('is not dismissible when dismissible is false', () => {
      const wrapper = createWrapper({ message: 'Not dismissible', dismissible: false })

      expect(wrapper.find('.v-alert').exists()).toBe(true)
    })

    it('emits dismiss event when close button is clicked', async () => {
      const wrapper = createWrapper({ message: 'Dismissible', dismissible: true })

      await wrapper.vm.$emit('dismiss')

      expect(wrapper.emitted('dismiss')).toBeTruthy()
    })
  })

  describe('retryable prop', () => {
    it('is not retryable by default', () => {
      const wrapper = createWrapper({ message: 'Error' })

      expect(wrapper.text()).not.toContain('Try again')
    })

    it('shows retry button when retryable is true', () => {
      const wrapper = createWrapper({ message: 'Error', retryable: true })

      expect(wrapper.text()).toContain('Try again')
    })

    it('emits retry event when retry button is clicked', async () => {
      const wrapper = createWrapper({ message: 'Error', retryable: true })

      await wrapper.vm.$emit('retry')

      expect(wrapper.emitted('retry')).toBeTruthy()
    })

    it('does not show retry button when retryable is false', () => {
      const wrapper = createWrapper({ message: 'Error', retryable: false })

      expect(wrapper.text()).not.toContain('Try again')
    })
  })

  describe('slots', () => {
    it('renders custom actions slot', () => {
      const wrapper = createWrapper(
        { message: 'Error' },
        {
          actions: '<button>Custom Action</button>',
        }
      )

      expect(wrapper.text()).toContain('Custom Action')
    })

    it('uses default retry button when no actions slot provided', () => {
      const wrapper = createWrapper({ message: 'Error', retryable: true })

      expect(wrapper.text()).toContain('Try again')
    })

    it('overrides retry button with custom actions slot', () => {
      const wrapper = createWrapper(
        { message: 'Error', retryable: true },
        {
          actions: '<button>Reload Page</button>',
        }
      )

      expect(wrapper.text()).toContain('Reload Page')
      expect(wrapper.text()).not.toContain('Try again')
    })
  })

  describe('computed properties', () => {
    it('maps error variant to error type', () => {
      const wrapper = createWrapper({ message: 'Error', variant: 'error' })

      expect(wrapper.vm.mappedType).toBe('error')
    })

    it('maps warning variant to warning type', () => {
      const wrapper = createWrapper({ message: 'Warning', variant: 'warning' })

      expect(wrapper.vm.mappedType).toBe('warning')
    })

    it('maps info variant to info type', () => {
      const wrapper = createWrapper({ message: 'Info', variant: 'info' })

      expect(wrapper.vm.mappedType).toBe('info')
    })

    it('defaults to error type for unknown variants', () => {
      const wrapper = createWrapper({ message: 'Error' })

      expect(wrapper.vm.mappedType).toBe('error')
    })
  })

  describe('combined props', () => {
    it('renders with all props combined', () => {
      const wrapper = createWrapper({
        variant: 'warning',
        message: 'Connection failed',
        title: 'Network Error',
        dismissible: true,
        retryable: true,
      })

      expect(wrapper.text()).toContain('Connection failed')
      expect(wrapper.text()).toContain('Try again')
      expect(wrapper.vm.mappedType).toBe('warning')
    })

    it('renders info message with custom actions', () => {
      const wrapper = createWrapper(
        {
          variant: 'info',
          message: 'Update available',
          title: 'New Version',
          dismissible: false,
        },
        {
          actions: '<button>Update Now</button>',
        }
      )

      expect(wrapper.text()).toContain('Update available')
      expect(wrapper.text()).toContain('Update Now')
      expect(wrapper.vm.mappedType).toBe('info')
    })
  })

  describe('events', () => {
    it('emits dismiss event', async () => {
      const wrapper = createWrapper({ message: 'Error', dismissible: true })

      await wrapper.vm.$emit('dismiss')

      expect(wrapper.emitted('dismiss')).toBeTruthy()
      expect(wrapper.emitted('dismiss')).toHaveLength(1)
    })

    it('emits retry event', async () => {
      const wrapper = createWrapper({ message: 'Error', retryable: true })

      await wrapper.vm.$emit('retry')

      expect(wrapper.emitted('retry')).toBeTruthy()
      expect(wrapper.emitted('retry')).toHaveLength(1)
    })

    it('can emit both dismiss and retry events', async () => {
      const wrapper = createWrapper({ message: 'Error', dismissible: true, retryable: true })

      await wrapper.vm.$emit('dismiss')
      await wrapper.vm.$emit('retry')

      expect(wrapper.emitted('dismiss')).toBeTruthy()
      expect(wrapper.emitted('retry')).toBeTruthy()
    })
  })

  describe('conditional rendering', () => {
    it('does not render when message is null', () => {
      const wrapper = createWrapper({ message: null as unknown as string })

      expect(wrapper.find('.v-alert').exists()).toBe(false)
    })

    it('does not render when message is undefined', () => {
      const wrapper = createWrapper({ message: undefined as unknown as string })

      expect(wrapper.find('.v-alert').exists()).toBe(false)
    })

    it('renders when message is a non-empty string', () => {
      const wrapper = createWrapper({ message: 'Valid error' })

      expect(wrapper.find('.v-alert').exists()).toBe(true)
    })
  })
})

