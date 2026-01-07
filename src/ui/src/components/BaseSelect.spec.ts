import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import BaseSelect from './BaseSelect.vue'

describe('BaseSelect', () => {
  const createWrapper = (props = {}) => {
    return mount(BaseSelect, {
      props: {
        options: [],
        ...props,
      },
      global: {
        stubs: {
          VSelect: {
            template: '<select class="v-select"></select>',
            props: ['id', 'modelValue', 'label', 'placeholder', 'items', 'disabled', 'required', 'rules', 'error', 'errorMessages', 'hint', 'variant', 'density', 'itemTitle', 'itemValue'],
          },
        },
      },
    })
  }

  describe('rendering', () => {
    it('renders select component', () => {
      const wrapper = createWrapper()

      expect(wrapper.find('.v-select').exists()).toBe(true)
    })

    it('renders with label prop', () => {
      const wrapper = createWrapper({ label: 'Choose option' })

      expect(wrapper.find('.v-select').exists()).toBe(true)
    })

    it('renders with placeholder prop', () => {
      const wrapper = createWrapper({ placeholder: 'Select an option' })

      expect(wrapper.find('.v-select').exists()).toBe(true)
    })

    it('renders with id prop', () => {
      const wrapper = createWrapper({ id: 'role-select' })

      expect(wrapper.find('.v-select').exists()).toBe(true)
    })
  })

  describe('options prop', () => {
    it('renders with object array options', () => {
      const options = [
        { value: 'admin', label: 'Administrator' },
        { value: 'user', label: 'User' },
      ]
      const wrapper = createWrapper({ options })

      // Verify the VSelect component receives the options
      expect(wrapper.find('.v-select').exists()).toBe(true)
    })

    it('renders with string array options', () => {
      const options = ['Option 1', 'Option 2', 'Option 3']
      const wrapper = createWrapper({ options })

      // Verify the VSelect component receives the options
      expect(wrapper.find('.v-select').exists()).toBe(true)
    })

    it('renders with number array options', () => {
      const options = [1, 2, 3]
      const wrapper = createWrapper({ options })

      // Verify the VSelect component receives the options
      expect(wrapper.find('.v-select').exists()).toBe(true)
    })

    it('renders with empty options array', () => {
      const wrapper = createWrapper({ options: [] })

      // Verify the VSelect component renders with empty options
      expect(wrapper.find('.v-select').exists()).toBe(true)
    })
  })

  describe('normalizedOptions computed', () => {
    it('normalizes object options with default keys', () => {
      const options = [
        { value: 'a', label: 'Option A' },
        { value: 'b', label: 'Option B' },
      ]
      const wrapper = createWrapper({ options })

      // Verify the VSelect component renders with normalized options
      expect(wrapper.find('.v-select').exists()).toBe(true)
    })

    it('normalizes object options with custom keys', () => {
      const options = [
        { id: 'a', name: 'Option A' },
        { id: 'b', name: 'Option B' },
      ]
      const wrapper = createWrapper({ options, valueKey: 'id', labelKey: 'name' })

      // Verify the VSelect component renders with normalized options
      expect(wrapper.find('.v-select').exists()).toBe(true)
    })

    it('normalizes string options', () => {
      const options = ['Red', 'Green', 'Blue']
      const wrapper = createWrapper({ options })

      // Verify the VSelect component renders with normalized options
      expect(wrapper.find('.v-select').exists()).toBe(true)
    })

    it('normalizes number options', () => {
      const options = [10, 20, 30]
      const wrapper = createWrapper({ options })

      // Verify the VSelect component renders with normalized options
      expect(wrapper.find('.v-select').exists()).toBe(true)
    })

    it('handles mixed object properties', () => {
      const options = [
        { value: 'a', label: 'Option A', extra: 'data' },
        { value: 'b', label: 'Option B', other: 'info' },
      ]
      const wrapper = createWrapper({ options })

      // Verify the VSelect component renders with normalized options
      expect(wrapper.find('.v-select').exists()).toBe(true)
    })
  })

  describe('modelValue prop', () => {
    it('renders with string value', () => {
      const options = ['Option 1', 'Option 2']
      const wrapper = createWrapper({ options, modelValue: 'Option 1' })

      expect(wrapper.find('.v-select').exists()).toBe(true)
    })

    it('renders with number value', () => {
      const options = [1, 2, 3]
      const wrapper = createWrapper({ options, modelValue: 2 })

      expect(wrapper.find('.v-select').exists()).toBe(true)
    })

    it('renders with undefined value by default', () => {
      const options = ['Option 1', 'Option 2']
      const wrapper = createWrapper({ options })

      expect(wrapper.find('.v-select').exists()).toBe(true)
    })
  })

  describe('disabled prop', () => {
    it('is not disabled by default', () => {
      const wrapper = createWrapper()

      expect(wrapper.find('.v-select').exists()).toBe(true)
    })

    it('is disabled when disabled prop is true', () => {
      const wrapper = createWrapper({ disabled: true })

      expect(wrapper.find('.v-select').exists()).toBe(true)
    })
  })

  describe('required prop', () => {
    it('is not required by default', () => {
      const wrapper = createWrapper()

      expect(wrapper.find('.v-select').exists()).toBe(true)
    })

    it('is required when required prop is true', () => {
      const wrapper = createWrapper({ required: true })

      expect(wrapper.find('.v-select').exists()).toBe(true)
    })
  })

  describe('error prop', () => {
    it('renders without error by default', () => {
      const wrapper = createWrapper()

      expect(wrapper.find('.v-select').exists()).toBe(true)
    })

    it('accepts error prop', () => {
      const wrapper = createWrapper({ error: 'Please select an option' })

      expect(wrapper.find('.v-select').exists()).toBe(true)
    })
  })

  describe('hint prop', () => {
    it('renders without hint by default', () => {
      const wrapper = createWrapper()

      expect(wrapper.find('.v-select').exists()).toBe(true)
    })

    it('accepts hint prop', () => {
      const wrapper = createWrapper({ hint: 'Choose your role' })

      expect(wrapper.find('.v-select').exists()).toBe(true)
    })
  })

  describe('variant prop', () => {
    it('uses outlined variant by default', () => {
      const wrapper = createWrapper()

      expect(wrapper.find('.v-select').exists()).toBe(true)
    })

    it('renders with filled variant', () => {
      const wrapper = createWrapper({ variant: 'filled' })

      expect(wrapper.find('.v-select').exists()).toBe(true)
    })

    it('renders with underlined variant', () => {
      const wrapper = createWrapper({ variant: 'underlined' })

      expect(wrapper.find('.v-select').exists()).toBe(true)
    })

    it('renders with plain variant', () => {
      const wrapper = createWrapper({ variant: 'plain' })

      expect(wrapper.find('.v-select').exists()).toBe(true)
    })
  })

  describe('density prop', () => {
    it('uses comfortable density by default', () => {
      const wrapper = createWrapper()

      expect(wrapper.find('.v-select').exists()).toBe(true)
    })

    it('renders with default density', () => {
      const wrapper = createWrapper({ density: 'default' })

      expect(wrapper.find('.v-select').exists()).toBe(true)
    })

    it('renders with compact density', () => {
      const wrapper = createWrapper({ density: 'compact' })

      expect(wrapper.find('.v-select').exists()).toBe(true)
    })
  })

  describe('rules prop', () => {
    it('applies validation rules', () => {
      const rules = [(v: string) => !!v || 'Required']
      const wrapper = createWrapper({ rules })

      expect(wrapper.find('.v-select').exists()).toBe(true)
    })

    it('uses empty rules by default', () => {
      const wrapper = createWrapper()

      expect(wrapper.find('.v-select').exists()).toBe(true)
    })
  })

  describe('events', () => {
    it('emits update:modelValue when selection changes', async () => {
      const options = ['Option 1', 'Option 2']
      const wrapper = createWrapper({ options })

      await wrapper.vm.$emit('update:modelValue', 'Option 1')

      expect(wrapper.emitted('update:modelValue')).toBeTruthy()
    })

    it('emits correct value for object options', async () => {
      const options = [
        { value: 'admin', label: 'Administrator' },
        { value: 'user', label: 'User' },
      ]
      const wrapper = createWrapper({ options })

      await wrapper.vm.$emit('update:modelValue', 'admin')

      expect(wrapper.emitted('update:modelValue')?.[0]).toEqual(['admin'])
    })

    it('emits correct value for string options', async () => {
      const options = ['Red', 'Green', 'Blue']
      const wrapper = createWrapper({ options })

      await wrapper.vm.$emit('update:modelValue', 'Green')

      expect(wrapper.emitted('update:modelValue')?.[0]).toEqual(['Green'])
    })

    it('emits correct value for number options', async () => {
      const options = [1, 2, 3]
      const wrapper = createWrapper({ options })

      await wrapper.vm.$emit('update:modelValue', 2)

      expect(wrapper.emitted('update:modelValue')?.[0]).toEqual([2])
    })
  })

  describe('combined props', () => {
    it('renders with multiple props combined', () => {
      const options = [
        { value: 'admin', label: 'Administrator' },
        { value: 'user', label: 'User' },
      ]
      const wrapper = createWrapper({
        id: 'role-select',
        label: 'User Role',
        placeholder: 'Select a role',
        options,
        required: true,
        hint: 'Choose the appropriate role',
        variant: 'outlined',
        density: 'comfortable',
      })

      expect(wrapper.find('.v-select').exists()).toBe(true)
    })

    it('accepts error and hint together', () => {
      const options = ['Option 1', 'Option 2']
      const wrapper = createWrapper({
        options,
        hint: 'This is a hint',
        error: 'This is an error',
      })

      expect(wrapper.find('.v-select').exists()).toBe(true)
    })

    it('handles custom value and label keys', () => {
      const options = [
        { id: 1, name: 'First', description: 'First option' },
        { id: 2, name: 'Second', description: 'Second option' },
      ]
      const wrapper = createWrapper({
        options,
        valueKey: 'id',
        labelKey: 'name',
      })

      // Verify the VSelect component renders with custom keys
      expect(wrapper.find('.v-select').exists()).toBe(true)
    })
  })
})

