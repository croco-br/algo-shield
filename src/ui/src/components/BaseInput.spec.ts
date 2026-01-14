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

  describe('basic functionality', () => {
    it('renders text field component', () => {
      const wrapper = createWrapper()

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })

    it('accepts text input props without errors', () => {
      const wrapper = createWrapper({
        id: 'test-input',
        label: 'Email',
        placeholder: 'Enter your email',
        type: 'email',
        required: true,
        disabled: false,
      })

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })

    it('accepts number input props without errors', () => {
      const wrapper = createWrapper({
        type: 'number',
        min: 0,
        max: 100,
        label: 'Age',
      })

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })

    it('accepts error state props without errors', () => {
      const wrapper = createWrapper({
        error: true,
        errorMessages: 'This field is required',
        hint: 'Please enter a valid value',
      })

      expect(wrapper.find('.v-text-field').exists()).toBe(true)
    })
  })
})
