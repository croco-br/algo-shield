<template>
  <v-text-field
    :id="id"
    :model-value="modelValue"
    :type="type"
    :label="label"
    :placeholder="placeholder"
    :disabled="disabled"
    :required="required"
    :error="!!error"
    :error-messages="error ? [error] : []"
    :hint="hint"
    :persistent-hint="persistentHint"
    :min="min"
    :max="max"
    :minlength="minlength"
    :maxlength="maxlength"
    :variant="variant"
    :density="density"
    :prepend-inner-icon="prependInnerIcon"
    :hide-details="hideDetails"
    :rules="rules"
    :pattern="pattern"
    @update:model-value="handleUpdate"
  />
</template>

<script setup lang="ts">
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
  persistentHint?: boolean
  min?: number | string
  max?: number | string
  minlength?: number
  maxlength?: number
  variant?: 'outlined' | 'filled' | 'underlined' | 'plain'
  density?: 'default' | 'comfortable' | 'compact'
  prependInnerIcon?: string
  hideDetails?: boolean | 'auto'
  rules?: ((value: any) => boolean | string)[]
  pattern?: string
}

const props = withDefaults(defineProps<Props>(), {
  type: 'text',
  disabled: false,
  required: false,
  variant: 'outlined',
  density: 'comfortable',
  persistentHint: false,
  hideDetails: false,
  rules: () => [],
})

const emit = defineEmits<{
  'update:modelValue': [value: string | number]
}>()

function handleUpdate(value: string | number) {
  emit('update:modelValue', value)
}

// Expose for testing
defineExpose({
  handleUpdate,
})
</script>
