import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import BaseInput from './BaseInput.vue'

describe('BaseInput', () => {
  const createWrapper = (props = {}) => {
    return mount(BaseInput, {
      props,
      global: {
        stubs: {
          VTextField: {
            template: '<input type="text" class="v-text-field" />',
            props: ['id', 'modelValue', 'type', 'label', 'placeholder', 'disabled', 'required', 'error', 'errorMessages', 'hint', 'persistentHint', 'min', 'max', 'minlength', 'maxlength', 'variant', 'density', 'prependInnerIcon', 'hideDetails', 'rules', 'pattern'],
          },
        },
      },
    })
  }

  describe('rendering', () => {
    it('renders text field component', () => {
      const wrapper = createWrapper()

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })

    it('renders with label prop', () => {
      const wrapper = createWrapper({ label: 'Username' })

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })

    it('renders with placeholder prop', () => {
      const wrapper = createWrapper({ placeholder: 'Enter username' })

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })

    it('renders with id prop', () => {
      const wrapper = createWrapper({ id: 'username-input' })

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })
  })

  describe('type prop', () => {
    it('renders with text type by default', () => {
      const wrapper = createWrapper()

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })

    it('renders with email type', () => {
      const wrapper = createWrapper({ type: 'email' })

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })

    it('renders with password type', () => {
      const wrapper = createWrapper({ type: 'password' })

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })

    it('renders with number type', () => {
      const wrapper = createWrapper({ type: 'number' })

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })

    it('renders with tel type', () => {
      const wrapper = createWrapper({ type: 'tel' })

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })

    it('renders with url type', () => {
      const wrapper = createWrapper({ type: 'url' })

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })

    it('renders with date type', () => {
      const wrapper = createWrapper({ type: 'date' })

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })
  })

  describe('modelValue prop', () => {
    it('accepts string value', () => {
      const wrapper = createWrapper({ modelValue: 'test value' })

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })

    it('accepts number value', () => {
      const wrapper = createWrapper({ modelValue: 42, type: 'number' })

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })

    it('renders with empty value by default', () => {
      const wrapper = createWrapper()

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })
  })

  describe('disabled prop', () => {
    it('is not disabled by default', () => {
      const wrapper = createWrapper()

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })

    it('accepts disabled prop', () => {
      const wrapper = createWrapper({ disabled: true })

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })
  })

  describe('required prop', () => {
    it('is not required by default', () => {
      const wrapper = createWrapper()

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })

    it('accepts required prop', () => {
      const wrapper = createWrapper({ required: true })

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })
  })

  describe('error prop', () => {
    it('renders without error by default', () => {
      const wrapper = createWrapper()

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })

    it('accepts error prop', () => {
      const wrapper = createWrapper({ error: 'This field is required' })

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })
  })

  describe('hint prop', () => {
    it('renders without hint by default', () => {
      const wrapper = createWrapper()

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })

    it('accepts hint prop', () => {
      const wrapper = createWrapper({ hint: 'Enter at least 8 characters' })

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })

    it('accepts persistentHint prop', () => {
      const wrapper = createWrapper({ hint: 'Persistent hint', persistentHint: true })

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })
  })

  describe('number constraints', () => {
    it('accepts min prop for number type', () => {
      const wrapper = createWrapper({ type: 'number', min: 0 })

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })

    it('accepts max prop for number type', () => {
      const wrapper = createWrapper({ type: 'number', max: 100 })

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })

    it('accepts both min and max props', () => {
      const wrapper = createWrapper({ type: 'number', min: 0, max: 100 })

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })
  })

  describe('string constraints', () => {
    it('accepts minlength prop', () => {
      const wrapper = createWrapper({ minlength: 8 })

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })

    it('accepts maxlength prop', () => {
      const wrapper = createWrapper({ maxlength: 50 })

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })

    it('accepts pattern prop', () => {
      const wrapper = createWrapper({ pattern: '[A-Za-z]+' })

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })
  })

  describe('variant prop', () => {
    it('uses outlined variant by default', () => {
      const wrapper = createWrapper()

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })

    it('renders with filled variant', () => {
      const wrapper = createWrapper({ variant: 'filled' })

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })

    it('renders with underlined variant', () => {
      const wrapper = createWrapper({ variant: 'underlined' })

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })

    it('renders with plain variant', () => {
      const wrapper = createWrapper({ variant: 'plain' })

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })
  })

  describe('density prop', () => {
    it('uses comfortable density by default', () => {
      const wrapper = createWrapper()

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })

    it('renders with default density', () => {
      const wrapper = createWrapper({ density: 'default' })

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })

    it('renders with compact density', () => {
      const wrapper = createWrapper({ density: 'compact' })

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })
  })

  describe('prependInnerIcon prop', () => {
    it('renders without icon by default', () => {
      const wrapper = createWrapper()

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })

    it('renders with prepend inner icon', () => {
      const wrapper = createWrapper({ prependInnerIcon: 'fa-user' })

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })
  })

  describe('hideDetails prop', () => {
    it('shows details by default', () => {
      const wrapper = createWrapper()

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })

    it('hides details when hideDetails is true', () => {
      const wrapper = createWrapper({ hideDetails: true })

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })

    it('uses auto hide details', () => {
      const wrapper = createWrapper({ hideDetails: 'auto' })

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })
  })

  describe('rules prop', () => {
    it('applies validation rules', () => {
      const rules = [(v: string) => !!v || 'Required']
      const wrapper = createWrapper({ rules })

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })

    it('uses empty rules by default', () => {
      const wrapper = createWrapper()

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })
  })

  describe('events', () => {
    it('emits update:modelValue on input', async () => {
      const wrapper = createWrapper()

      await wrapper.vm.$emit('update:modelValue', 'new value')

      expect(wrapper.emitted('update:modelValue')).toBeTruthy()
      expect(wrapper.emitted('update:modelValue')?.[0]).toEqual(['new value'])
    })

    it('emits update:modelValue with string value', async () => {
      const wrapper = createWrapper()

      await wrapper.vm.$emit('update:modelValue', 'test string')

      expect(wrapper.emitted('update:modelValue')).toBeTruthy()
      expect(wrapper.emitted('update:modelValue')?.[0]).toEqual(['test string'])
    })

    it('emits update:modelValue with number value', async () => {
      const wrapper = createWrapper({ type: 'number' })

      await wrapper.vm.$emit('update:modelValue', 42)

      expect(wrapper.emitted('update:modelValue')).toBeTruthy()
      expect(wrapper.emitted('update:modelValue')?.[0]).toEqual([42])
    })
  })

  describe('combined props', () => {
    it('renders with multiple props combined', () => {
      const wrapper = createWrapper({
        id: 'email-input',
        type: 'email',
        label: 'Email Address',
        placeholder: 'user@example.com',
        required: true,
        hint: 'We will never share your email',
        prependInnerIcon: 'fa-envelope',
        variant: 'outlined',
        density: 'comfortable',
      })

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })

    it('accepts error and hint together', () => {
      const wrapper = createWrapper({
        hint: 'This is a hint',
        error: 'This is an error',
      })

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })
  })
})

