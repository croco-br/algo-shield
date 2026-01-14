import { describe, it, expect, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import { createI18n } from 'vue-i18n'
import BaseTable from './BaseTable.vue'

describe('BaseTable', () => {
  let vuetify: ReturnType<typeof createVuetify>
  let i18n: ReturnType<typeof createI18n>

  beforeEach(() => {
    vuetify = createVuetify({
      components,
      directives,
    })

    i18n = createI18n({
      legacy: false,
      locale: 'en-US',
      messages: {
        'en-US': {
          common: {
            name: 'Name',
            email: 'Email',
            status: 'Status',
          },
        },
      },
    })
  })

  it('renders table with columns and data', () => {
    const columns = [
      { key: 'name', label: 'Name' },
      { key: 'email', label: 'Email' },
    ]
    const data = [
      { id: 1, name: 'John Doe', email: 'john@example.com' },
      { id: 2, name: 'Jane Smith', email: 'jane@example.com' },
    ]

    const wrapper = mount(BaseTable, {
      props: { columns, data },
      global: {
        plugins: [vuetify, i18n],
      },
    })

    expect(wrapper.find('.v-data-table').exists()).toBe(true)
  })

  it('displays empty text when no data', () => {
    const columns = [{ key: 'name', label: 'Name' }]
    const data: any[] = []

    const wrapper = mount(BaseTable, {
      props: { columns, data },
      global: {
        plugins: [vuetify, i18n],
      },
    })

    expect(wrapper.text()).toContain('No data available')
  })

  it('uses custom empty text', () => {
    const columns = [{ key: 'name', label: 'Name' }]
    const data: any[] = []
    const emptyText = 'Custom empty message'

    const wrapper = mount(BaseTable, {
      props: { columns, data, emptyText },
      global: {
        plugins: [vuetify, i18n],
      },
    })

    expect(wrapper.text()).toContain(emptyText)
  })

  it('hides empty message when hideEmpty is true', () => {
    const columns = [{ key: 'name', label: 'Name' }]
    const data: any[] = []

    const wrapper = mount(BaseTable, {
      props: { columns, data, hideEmpty: true },
      global: {
        plugins: [vuetify, i18n],
      },
    })

    expect(wrapper.text()).not.toContain('No data available')
  })

  it('translates column labels with translation keys', () => {
    const columns = [
      { key: 'name', label: 'common.name' },
      { key: 'email', label: 'common.email' },
    ]
    const data = [{ id: 1, name: 'John', email: 'john@example.com' }]

    const wrapper = mount(BaseTable, {
      props: { columns, data },
      global: {
        plugins: [vuetify, i18n],
      },
    })

    expect(wrapper.html()).toContain('Name')
    expect(wrapper.html()).toContain('Email')
  })

  it('renders nested values with dot notation', () => {
    const columns = [
      { key: 'user.name', label: 'Name' },
      { key: 'user.profile.email', label: 'Email' },
    ]
    const data = [
      {
        id: 1,
        user: {
          name: 'John',
          profile: { email: 'john@example.com' },
        },
      },
    ]

    const wrapper = mount(BaseTable, {
      props: { columns, data },
      global: {
        plugins: [vuetify, i18n],
      },
    })

    const vm = wrapper.vm as any
    expect(vm.getNestedValue(data[0], 'user.name')).toBe('John')
    expect(vm.getNestedValue(data[0], 'user.profile.email')).toBe('john@example.com')
  })

  it('handles missing nested values gracefully', () => {
    const columns = [{ key: 'user.missing.value', label: 'Missing' }]
    const data = [{ id: 1, user: { name: 'John' } }]

    const wrapper = mount(BaseTable, {
      props: { columns, data },
      global: {
        plugins: [vuetify, i18n],
      },
    })

    const vm = wrapper.vm as any
    expect(vm.getNestedValue(data[0], 'user.missing.value')).toBeUndefined()
  })

  it('uses custom rowKey for item identification', () => {
    const columns = [{ key: 'name', label: 'Name' }]
    const data = [
      { uuid: 'abc-123', name: 'John' },
      { uuid: 'def-456', name: 'Jane' },
    ]

    const wrapper = mount(BaseTable, {
      props: { columns, data, rowKey: 'uuid' },
      global: {
        plugins: [vuetify, i18n],
      },
    })

    const vm = wrapper.vm as any
    expect(vm.getItemKey(data[0])).toBe('abc-123')
    expect(vm.getItemKey(data[1])).toBe('def-456')
  })

  it('falls back to default id rowKey', () => {
    const columns = [{ key: 'name', label: 'Name' }]
    const data = [{ id: 1, name: 'John' }]

    const wrapper = mount(BaseTable, {
      props: { columns, data },
      global: {
        plugins: [vuetify, i18n],
      },
    })

    const vm = wrapper.vm as any
    expect(vm.getItemKey(data[0])).toBe(1)
  })

  it('generates random key when rowKey is missing', () => {
    const columns = [{ key: 'name', label: 'Name' }]
    const data = [{ name: 'John' }]

    const wrapper = mount(BaseTable, {
      props: { columns, data },
      global: {
        plugins: [vuetify, i18n],
      },
    })

    const vm = wrapper.vm as any
    const key = vm.getItemKey(data[0])
    expect(typeof key).toBe('string')
    expect(key.length).toBeGreaterThan(0)
  })

  it('renders header slot', () => {
    const columns = [{ key: 'name', label: 'Name' }]
    const data = [{ id: 1, name: 'John' }]

    const wrapper = mount(BaseTable, {
      props: { columns, data },
      slots: {
        header: '<div class="custom-header">Custom Header</div>',
      },
      global: {
        plugins: [vuetify, i18n],
      },
    })

    expect(wrapper.find('.custom-header').exists()).toBe(true)
    expect(wrapper.text()).toContain('Custom Header')
  })

  it('renders footer slot', () => {
    const columns = [{ key: 'name', label: 'Name' }]
    const data = [{ id: 1, name: 'John' }]

    const wrapper = mount(BaseTable, {
      props: { columns, data },
      slots: {
        footer: '<div class="custom-footer">Custom Footer</div>',
      },
      global: {
        plugins: [vuetify, i18n],
      },
    })

    expect(wrapper.find('.custom-footer').exists()).toBe(true)
    expect(wrapper.text()).toContain('Custom Footer')
  })

  it('renders empty slot', () => {
    const columns = [{ key: 'name', label: 'Name' }]
    const data: any[] = []

    const wrapper = mount(BaseTable, {
      props: { columns, data },
      slots: {
        empty: '<div class="custom-empty">Custom Empty</div>',
      },
      global: {
        plugins: [vuetify, i18n],
      },
    })

    expect(wrapper.find('.custom-empty').exists()).toBe(true)
    expect(wrapper.text()).toContain('Custom Empty')
  })

  it('renders custom cell slot', () => {
    const columns = [
      { key: 'name', label: 'Name' },
      { key: 'status', label: 'Status' },
    ]
    const data = [{ id: 1, name: 'John', status: 'active' }]

    const wrapper = mount(BaseTable, {
      props: { columns, data },
      slots: {
        'cell-status': '<span class="status-badge">{{ row.status }}</span>',
      },
      global: {
        plugins: [vuetify, i18n],
      },
    })

    expect(wrapper.find('.status-badge').exists()).toBe(true)
  })

  it('normalizes headers correctly', () => {
    const columns = [
      { key: 'name', label: 'Name', width: '200px' },
      { key: 'email', label: 'Email' },
    ]
    const data = [{ id: 1, name: 'John', email: 'john@example.com' }]

    const wrapper = mount(BaseTable, {
      props: { columns, data },
      global: {
        plugins: [vuetify, i18n],
      },
    })

    const vm = wrapper.vm as any
    expect(vm.normalizedHeaders).toHaveLength(2)
    expect(vm.normalizedHeaders[0].key).toBe('name')
    expect(vm.normalizedHeaders[0].width).toBe('200px')
    expect(vm.normalizedHeaders[0].sortable).toBe(false)
  })

  it('handles empty columns array', () => {
    const columns: any[] = []
    const data = [{ id: 1, name: 'John' }]

    const wrapper = mount(BaseTable, {
      props: { columns, data },
      global: {
        plugins: [vuetify, i18n],
      },
    })

    expect(wrapper.find('.v-data-table').exists()).toBe(true)
  })

  it('handles complex nested data structures', () => {
    const columns = [
      { key: 'metadata.transaction.amount', label: 'Amount' },
      { key: 'metadata.user.details.country', label: 'Country' },
    ]
    const data = [
      {
        id: 1,
        metadata: {
          transaction: { amount: 100.5 },
          user: { details: { country: 'US' } },
        },
      },
    ]

    const wrapper = mount(BaseTable, {
      props: { columns, data },
      global: {
        plugins: [vuetify, i18n],
      },
    })

    const vm = wrapper.vm as any
    expect(vm.getNestedValue(data[0], 'metadata.transaction.amount')).toBe(100.5)
    expect(vm.getNestedValue(data[0], 'metadata.user.details.country')).toBe('US')
  })
})
