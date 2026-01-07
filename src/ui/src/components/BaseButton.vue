<template>
  <v-btn
    :type="type"
    :disabled="disabled || loading"
    :size="mappedSize"
    :variant="mappedVariant"
    :color="mappedColor"
    :block="fullWidth"
    :loading="loading"
    :prepend-icon="prependIcon"
    class="base-button"
  >
    <slot />
  </v-btn>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  variant?: 'primary' | 'secondary' | 'danger' | 'ghost'
  size?: 'sm' | 'md' | 'lg'
  type?: 'button' | 'submit' | 'reset'
  disabled?: boolean
  loading?: boolean
  fullWidth?: boolean
  prependIcon?: string
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'primary',
  size: 'md',
  type: 'button',
  disabled: false,
  loading: false,
  fullWidth: false,
})

// Map sizes consistently
const mappedSize = computed(() => {
  switch (props.size) {
    case 'sm':
      return 'small'
    case 'lg':
      return 'large'
    default:
      return 'default'
  }
})

// Map our variants to Vuetify variants and colors
const mappedVariant = computed(() => {
  if (props.variant === 'ghost') {
    return 'text'
  }
  if (props.variant === 'secondary') {
    return 'outlined'
  }
  return 'flat' // primary and danger use flat with color
})

const mappedColor = computed(() => {
  if (props.variant === 'danger') {
    return 'error'
  }
  if (props.variant === 'secondary') {
    return undefined // outlined buttons use default color
  }
  if (props.variant === 'ghost') {
    return undefined // text buttons use default color
  }
  return 'primary' // primary variant
})

// Expose for testing
defineExpose({
  mappedSize,
  mappedVariant,
  mappedColor,
})
</script>

<style scoped>
.base-button {
  text-transform: none;
  font-family: var(--font-family-sans);
  font-weight: 500;
  letter-spacing: 0.01em;
}
</style>
