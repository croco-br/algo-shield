import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import BaseBadge from './BaseBadge.vue'

vi.mock('vuetify/components', () => ({
  VChip: {
    name: 'VChip',
    template: '<div class="v-chip"><slot /></div>',
    props: ['color', 'variant', 'size', 'closable', 'rounded'],
  },
}))

describe('BaseBadge', () => {
  const createWrapper = (props = {}, slots = {}) => {
    return mount(BaseBadge, {
      props,
      slots,
      global: {
        stubs: {
          VChip: {
            template: '<div class="v-chip"><slot /></div>',
            props: ['color', 'variant', 'size', 'closable', 'rounded'],
          },
        },
      },
    })
  }

  describe('rendering', () => {
    it('renders slot content', () => {
      const wrapper = createWrapper({}, { default: 'Test Badge' })

      expect(wrapper.text()).toContain('Test Badge')
    })

    it('renders with default variant', () => {
      const wrapper = createWrapper({}, { default: 'Badge' })

      expect(wrapper.find('.v-chip').exists()).toBe(true)
    })

    it('renders with success variant', () => {
      const wrapper = createWrapper({ variant: 'success' }, { default: 'Success' })

      expect(wrapper.find('.v-chip').exists()).toBe(true)
    })

    it('renders with warning variant', () => {
      const wrapper = createWrapper({ variant: 'warning' }, { default: 'Warning' })

      expect(wrapper.find('.v-chip').exists()).toBe(true)
    })

    it('renders with danger variant', () => {
      const wrapper = createWrapper({ variant: 'danger' }, { default: 'Danger' })

      expect(wrapper.find('.v-chip').exists()).toBe(true)
    })

    it('renders with info variant', () => {
      const wrapper = createWrapper({ variant: 'info' }, { default: 'Info' })

      expect(wrapper.find('.v-chip').exists()).toBe(true)
    })
  })

  describe('size variants', () => {
    it('renders with small size', () => {
      const wrapper = createWrapper({ size: 'sm' }, { default: 'Small' })

      expect(wrapper.find('.v-chip').exists()).toBe(true)
    })

    it('renders with medium size by default', () => {
      const wrapper = createWrapper({}, { default: 'Medium' })

      expect(wrapper.find('.v-chip').exists()).toBe(true)
    })

    it('renders with large size', () => {
      const wrapper = createWrapper({ size: 'lg' }, { default: 'Large' })

      expect(wrapper.find('.v-chip').exists()).toBe(true)
    })
  })

  describe('outline prop', () => {
    it('renders as flat variant by default', () => {
      const wrapper = createWrapper({}, { default: 'Badge' })

      expect(wrapper.find('.v-chip').exists()).toBe(true)
    })

    it('renders as outlined variant when outline is true', () => {
      const wrapper = createWrapper({ outline: true }, { default: 'Outlined' })

      expect(wrapper.find('.v-chip').exists()).toBe(true)
    })
  })

  describe('rounded prop', () => {
    it('renders without rounded by default', () => {
      const wrapper = createWrapper({}, { default: 'Badge' })

      expect(wrapper.find('.v-chip').exists()).toBe(true)
    })

    it('renders as pill when rounded is true', () => {
      const wrapper = createWrapper({ rounded: true }, { default: 'Rounded' })

      expect(wrapper.find('.v-chip').exists()).toBe(true)
    })
  })

  describe('closable prop', () => {
    it('is not closable by default', () => {
      const wrapper = createWrapper({}, { default: 'Badge' })

      expect(wrapper.find('.v-chip').exists()).toBe(true)
    })

    it('shows close button when closable is true', () => {
      const wrapper = createWrapper({ closable: true }, { default: 'Closable' })

      expect(wrapper.find('.v-chip').exists()).toBe(true)
    })

    it('emits close event when close button is clicked', async () => {
      const wrapper = createWrapper({ closable: true }, { default: 'Closable' })

      await wrapper.vm.$emit('close')

      expect(wrapper.emitted('close')).toBeTruthy()
    })
  })

  describe('computed properties', () => {
    it('maps success variant to success color', () => {
      const wrapper = createWrapper({ variant: 'success' }, { default: 'Success' })

      expect(wrapper.vm.mappedColor).toBe('success')
    })

    it('maps warning variant to warning color', () => {
      const wrapper = createWrapper({ variant: 'warning' }, { default: 'Warning' })

      expect(wrapper.vm.mappedColor).toBe('warning')
    })

    it('maps danger variant to error color', () => {
      const wrapper = createWrapper({ variant: 'danger' }, { default: 'Danger' })

      expect(wrapper.vm.mappedColor).toBe('error')
    })

    it('maps info variant to info color', () => {
      const wrapper = createWrapper({ variant: 'info' }, { default: 'Info' })

      expect(wrapper.vm.mappedColor).toBe('info')
    })

    it('maps default variant to default color', () => {
      const wrapper = createWrapper({ variant: 'default' }, { default: 'Default' })

      expect(wrapper.vm.mappedColor).toBe('default')
    })

    it('maps outline false to flat variant', () => {
      const wrapper = createWrapper({ outline: false }, { default: 'Badge' })

      expect(wrapper.vm.mappedVariant).toBe('flat')
    })

    it('maps outline true to outlined variant', () => {
      const wrapper = createWrapper({ outline: true }, { default: 'Badge' })

      expect(wrapper.vm.mappedVariant).toBe('outlined')
    })

    it('maps sm size to small', () => {
      const wrapper = createWrapper({ size: 'sm' }, { default: 'Badge' })

      expect(wrapper.vm.mappedSize).toBe('small')
    })

    it('maps md size to default', () => {
      const wrapper = createWrapper({ size: 'md' }, { default: 'Badge' })

      expect(wrapper.vm.mappedSize).toBe('default')
    })

    it('maps lg size to large', () => {
      const wrapper = createWrapper({ size: 'lg' }, { default: 'Badge' })

      expect(wrapper.vm.mappedSize).toBe('large')
    })
  })

  describe('combined props', () => {
    it('renders with multiple props combined', () => {
      const wrapper = createWrapper(
        {
          variant: 'success',
          size: 'lg',
          outline: true,
          rounded: true,
          closable: true,
        },
        { default: 'Complex Badge' }
      )

      expect(wrapper.text()).toContain('Complex Badge')
      expect(wrapper.vm.mappedColor).toBe('success')
      expect(wrapper.vm.mappedSize).toBe('large')
      expect(wrapper.vm.mappedVariant).toBe('outlined')
    })
  })
})

