<template>
  <div class="w-full">
    <label
      v-if="label"
      :for="id"
      class="block text-sm font-medium text-gray-700 mb-2"
    >
      {{ label }}
      <span v-if="required" class="text-red-500">*</span>
    </label>
    <select
      :id="id"
      :value="modelValue"
      :disabled="disabled"
      :required="required"
      :class="selectClasses"
      @change="handleChange"
    >
      <option v-if="placeholder" value="" disabled>{{ placeholder }}</option>
      <option
        v-for="option in options"
        :key="getOptionValue(option)"
        :value="getOptionValue(option)"
      >
        {{ getOptionLabel(option) }}
      </option>
    </select>
    <small v-if="hint" class="block text-xs text-gray-500 mt-1">
      {{ hint }}
    </small>
    <small v-if="error" class="block text-xs text-red-600 mt-1">
      {{ error }}
    </small>
  </div>
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
  error?: string
  hint?: string
  valueKey?: string
  labelKey?: string
}

const props = withDefaults(defineProps<Props>(), {
  disabled: false,
  required: false,
  valueKey: 'value',
  labelKey: 'label',
})

const emit = defineEmits<{
  'update:modelValue': [value: string | number]
}>()

const selectClasses = computed(() => {
  const classes = [
    'w-full px-4 py-3',
    'border rounded-md',
    'text-gray-900 text-sm',
    'transition-colors',
    'focus:outline-none focus:ring-2 focus:border-transparent',
    'disabled:opacity-60 disabled:cursor-not-allowed disabled:bg-gray-50',
    'bg-white',
    'appearance-none',
    'bg-no-repeat',
    'bg-[right_0.75rem_center]',
    'bg-[length:1.25rem]',
    'pr-10',
  ]

  if (props.error) {
    classes.push('border-red-300 focus:ring-red-500')
  } else {
    classes.push('border-gray-300 focus:ring-indigo-500')
  }

  return classes.join(' ')
})

const getOptionValue = (option: Option | string | number): string | number => {
  if (typeof option === 'object') {
    return option[props.valueKey]
  }
  return option
}

const getOptionLabel = (option: Option | string | number): string => {
  if (typeof option === 'object') {
    return option[props.labelKey]
  }
  return String(option)
}

const handleChange = (event: Event) => {
  const target = event.target as HTMLSelectElement
  emit('update:modelValue', target.value)
}
</script>

<style scoped>
select {
  background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 20 20'%3e%3cpath stroke='%236b7280' stroke-linecap='round' stroke-linejoin='round' stroke-width='1.5' d='M6 8l4 4 4-4'/%3e%3c/svg%3e");
}
</style>
