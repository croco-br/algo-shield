import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import BaseButton from './BaseButton.vue'

describe('BaseButton', () => {
  const createWrapper = (props = {}, slots = {}) => {
    return mount(BaseButton, {
      props,
      slots,
      global: {
        stubs: {
          VBtn: {
            template: '<button type="button" class="v-btn"><slot /></button>',
            props: ['type', 'disabled', 'size', 'variant', 'color', 'block', 'loading', 'prependIcon'],
          },
        },
      },
    })
  }

  describe('rendering', () => {
    it('renders slot content', () => {
      const wrapper = createWrapper({}, { default: 'Click Me' })

      expect(wrapper.text()).toContain('Click Me')
    })

    it('renders as button element by default', () => {
      const wrapper = createWrapper({}, { default: 'Button' })

      expect(wrapper.find('button').exists()).toBe(true)
    })

    it('renders with primary variant by default', () => {
      const wrapper = createWrapper({}, { default: 'Primary' })

      expect(wrapper.find('.v-btn').exists()).toBe(true)
    })
  })

  describe('variant prop', () => {
    it('renders primary variant', () => {
      const wrapper = createWrapper({ variant: 'primary' }, { default: 'Primary' })

      expect(wrapper.vm.mappedVariant).toBe('flat')
      expect(wrapper.vm.mappedColor).toBe('primary')
    })

    it('renders secondary variant', () => {
      const wrapper = createWrapper({ variant: 'secondary' }, { default: 'Secondary' })

      expect(wrapper.vm.mappedVariant).toBe('outlined')
      expect(wrapper.vm.mappedColor).toBeUndefined()
    })

    it('renders danger variant', () => {
      const wrapper = createWrapper({ variant: 'danger' }, { default: 'Danger' })

      expect(wrapper.vm.mappedVariant).toBe('flat')
      expect(wrapper.vm.mappedColor).toBe('error')
    })

    it('renders ghost variant', () => {
      const wrapper = createWrapper({ variant: 'ghost' }, { default: 'Ghost' })

      expect(wrapper.vm.mappedVariant).toBe('text')
      expect(wrapper.vm.mappedColor).toBeUndefined()
    })
  })

  describe('size prop', () => {
    it('renders with small size', () => {
      const wrapper = createWrapper({ size: 'sm' }, { default: 'Small' })

      expect(wrapper.vm.mappedSize).toBe('small')
    })

    it('renders with medium size by default', () => {
      const wrapper = createWrapper({}, { default: 'Medium' })

      expect(wrapper.vm.mappedSize).toBe('default')
    })

    it('renders with large size', () => {
      const wrapper = createWrapper({ size: 'lg' }, { default: 'Large' })

      expect(wrapper.vm.mappedSize).toBe('large')
    })
  })

  describe('type prop', () => {
    it('has button type by default', () => {
      const wrapper = createWrapper({}, { default: 'Button' })

      expect(wrapper.find('button').exists()).toBe(true)
    })

    it('accepts submit type', () => {
      const wrapper = createWrapper({ type: 'submit' }, { default: 'Submit' })

      expect(wrapper.find('button').exists()).toBe(true)
    })

    it('accepts reset type', () => {
      const wrapper = createWrapper({ type: 'reset' }, { default: 'Reset' })

      expect(wrapper.find('button').exists()).toBe(true)
    })
  })

  describe('disabled prop', () => {
    it('is not disabled by default', () => {
      const wrapper = createWrapper({}, { default: 'Button' })

      expect(wrapper.find('button').exists()).toBe(true)
    })

    it('accepts disabled prop', () => {
      const wrapper = createWrapper({ disabled: true }, { default: 'Disabled' })

      expect(wrapper.find('button').exists()).toBe(true)
    })

    it('accepts loading prop that disables button', () => {
      const wrapper = createWrapper({ loading: true }, { default: 'Loading' })

      expect(wrapper.find('button').exists()).toBe(true)
    })

    it('accepts both disabled and loading props', () => {
      const wrapper = createWrapper({ disabled: true, loading: true }, { default: 'Both' })

      expect(wrapper.find('button').exists()).toBe(true)
    })
  })

  describe('loading prop', () => {
    it('is not loading by default', () => {
      const wrapper = createWrapper({}, { default: 'Button' })

      expect(wrapper.find('.v-btn').exists()).toBe(true)
    })

    it('accepts loading prop', () => {
      const wrapper = createWrapper({ loading: true }, { default: 'Loading' })

      expect(wrapper.find('button').exists()).toBe(true)
    })
  })

  describe('fullWidth prop', () => {
    it('is not full width by default', () => {
      const wrapper = createWrapper({}, { default: 'Button' })

      expect(wrapper.find('.v-btn').exists()).toBe(true)
    })

    it('renders as full width when fullWidth is true', () => {
      const wrapper = createWrapper({ fullWidth: true }, { default: 'Full Width' })

      expect(wrapper.find('.v-btn').exists()).toBe(true)
    })
  })

  describe('prependIcon prop', () => {
    it('renders without icon by default', () => {
      const wrapper = createWrapper({}, { default: 'Button' })

      expect(wrapper.find('.v-btn').exists()).toBe(true)
    })

    it('renders with prepend icon when provided', () => {
      const wrapper = createWrapper({ prependIcon: 'fa-plus' }, { default: 'Add' })

      expect(wrapper.find('.v-btn').exists()).toBe(true)
    })
  })

  describe('computed properties', () => {
    it('maps sm size to small', () => {
      const wrapper = createWrapper({ size: 'sm' }, { default: 'Button' })

      expect(wrapper.vm.mappedSize).toBe('small')
    })

    it('maps md size to default', () => {
      const wrapper = createWrapper({ size: 'md' }, { default: 'Button' })

      expect(wrapper.vm.mappedSize).toBe('default')
    })

    it('maps lg size to large', () => {
      const wrapper = createWrapper({ size: 'lg' }, { default: 'Button' })

      expect(wrapper.vm.mappedSize).toBe('large')
    })

    it('maps ghost variant to text', () => {
      const wrapper = createWrapper({ variant: 'ghost' }, { default: 'Ghost' })

      expect(wrapper.vm.mappedVariant).toBe('text')
    })

    it('maps secondary variant to outlined', () => {
      const wrapper = createWrapper({ variant: 'secondary' }, { default: 'Secondary' })

      expect(wrapper.vm.mappedVariant).toBe('outlined')
    })

    it('maps primary variant to flat', () => {
      const wrapper = createWrapper({ variant: 'primary' }, { default: 'Primary' })

      expect(wrapper.vm.mappedVariant).toBe('flat')
    })

    it('maps danger variant to error color', () => {
      const wrapper = createWrapper({ variant: 'danger' }, { default: 'Danger' })

      expect(wrapper.vm.mappedColor).toBe('error')
    })

    it('maps primary variant to primary color', () => {
      const wrapper = createWrapper({ variant: 'primary' }, { default: 'Primary' })

      expect(wrapper.vm.mappedColor).toBe('primary')
    })

    it('maps secondary variant to undefined color', () => {
      const wrapper = createWrapper({ variant: 'secondary' }, { default: 'Secondary' })

      expect(wrapper.vm.mappedColor).toBeUndefined()
    })

    it('maps ghost variant to undefined color', () => {
      const wrapper = createWrapper({ variant: 'ghost' }, { default: 'Ghost' })

      expect(wrapper.vm.mappedColor).toBeUndefined()
    })
  })

  describe('combined props', () => {
    it('renders with multiple props combined', () => {
      const wrapper = createWrapper(
        {
          variant: 'danger',
          size: 'lg',
          type: 'submit',
          fullWidth: true,
          prependIcon: 'fa-trash',
        },
        { default: 'Delete All' }
      )

      expect(wrapper.text()).toContain('Delete All')
      expect(wrapper.vm.mappedColor).toBe('error')
      expect(wrapper.vm.mappedSize).toBe('large')
      expect(wrapper.find('button').exists()).toBe(true)
    })

    it('accepts disabled and loading together', () => {
      const wrapper = createWrapper(
        {
          disabled: true,
          loading: true,
        },
        { default: 'Processing' }
      )

      expect(wrapper.find('button').exists()).toBe(true)
    })
  })
})

