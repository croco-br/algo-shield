<template>
  <div class="w-full">
    <label
      v-if="label"
      :for="id"
      class="block text-sm font-semibold text-slate-700 mb-2.5"
    >
      {{ label }}
      <span v-if="required" class="text-red-500">*</span>
    </label>
    <div class="relative">
      <input
        :id="id"
        :type="type"
        :value="modelValue"
        :placeholder="placeholder"
        :disabled="disabled"
        :required="required"
        :min="min"
        :max="max"
        :minlength="minlength"
        :maxlength="maxlength"
        :class="inputClasses"
        @input="handleInput"
        @focus="isFocused = true"
        @blur="isFocused = false"
      />
      <div v-if="type === 'password'" class="absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none">
        <svg class="w-4 h-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/>
        </svg>
      </div>
    </div>
    <small v-if="hint" class="block text-xs text-slate-500 mt-2">
      {{ hint }}
    </small>
    <small v-if="error" class="block text-xs text-red-600 mt-2">
      {{ error }}
    </small>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

interface Props {
  id?: string
  type?: 'text' | 'email' | 'password' | 'number' | 'tel' | 'url' | 'date'
  modelValue?: string | number
  label?: string
  placeholder?: string
  disabled?: boolean
  required?: boolean
  error?: string
  hint?: string
  min?: number | string
  max?: number | string
  minlength?: number
  maxlength?: number
}

const props = withDefaults(defineProps<Props>(), {
  type: 'text',
  disabled: false,
  required: false,
})

const emit = defineEmits<{
  'update:modelValue': [value: string | number]
}>()

const isFocused = ref(false)

const inputClasses = computed(() => {
  const classes = [
    'w-full px-4 py-3.5',
    'border-2 rounded-lg',
    'text-slate-900 text-sm font-medium',
    'placeholder:text-slate-400 placeholder:font-normal',
    'transition-all duration-200',
    'focus:outline-none focus:ring-4',
    'disabled:opacity-50 disabled:cursor-not-allowed disabled:bg-slate-50',
  ]

  if (props.error) {
    classes.push('border-red-300 focus:border-red-500 focus:ring-red-100')
  } else {
    classes.push('border-slate-200 focus:border-blue-600 focus:ring-blue-50')
  }

  return classes.join(' ')
})

const handleInput = (event: Event) => {
  const target = event.target as HTMLInputElement
  emit('update:modelValue', target.value)
}
</script>
