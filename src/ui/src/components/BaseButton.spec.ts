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

  describe('basic functionality', () => {
    it('renders slot content', () => {
      const wrapper = createWrapper({}, { default: 'Click Me' })

      expect(wrapper.text()).toContain('Click Me')
    })

    it('renders as button element', () => {
      const wrapper = createWrapper({}, { default: 'Button' })

      expect(wrapper.find('button').exists()).toBe(true)
    })

    it('accepts various props without errors', () => {
      const wrapper = createWrapper(
        {
          variant: 'primary',
          size: 'lg',
          type: 'submit',
          disabled: false,
          loading: false,
          fullWidth: true,
          prependIcon: 'fa-plus',
        },
        { default: 'Submit' }
      )

      expect(wrapper.find('button').exists()).toBe(true)
      expect(wrapper.text()).toContain('Submit')
    })
  })
})
