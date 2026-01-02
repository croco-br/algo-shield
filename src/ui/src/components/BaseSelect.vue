<template>
  <v-select
    :id="id"
    :model-value="modelValue"
    :label="label"
    :placeholder="placeholder"
    :items="normalizedOptions"
    :disabled="disabled"
    :required="required"
    :rules="rules"
    :error="!!error"
    :error-messages="error ? [error] : []"
    :hint="hint"
    :variant="variant"
    :density="density"
    item-title="label"
    item-value="value"
    @update:model-value="handleUpdate"
  />
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Option {
  value: string | number
  label: string
  [key: string]: any
}

interface Props {
  id?: string
  modelValue?: string | number
  label?: string
  placeholder?: string
  options: Option[] | string[] | number[]
  disabled?: boolean
  required?: boolean
  rules?: ((value: any) => boolean | string)[]
  error?: string
  hint?: string
  valueKey?: string
  labelKey?: string
  variant?: 'outlined' | 'filled' | 'underlined' | 'plain'
  density?: 'default' | 'comfortable' | 'compact'
}

const props = withDefaults(defineProps<Props>(), {
  disabled: false,
  required: false,
  rules: () => [],
  valueKey: 'value',
  labelKey: 'label',
  variant: 'outlined',
  density: 'comfortable',
})

const emit = defineEmits<{
  'update:modelValue': [value: string | number]
}>()

// Normalize options to Vuetify format
const normalizedOptions = computed(() => {
  return props.options.map(option => {
    if (typeof option === 'object') {
      return {
        value: option[props.valueKey],
        label: option[props.labelKey],
      }
    }
    return {
      value: option,
      label: String(option),
    }
  })
})

function handleUpdate(value: any) {
  // Vuetify select with item-title and item-value returns the value directly
  emit('update:modelValue', value)
}
</script>
