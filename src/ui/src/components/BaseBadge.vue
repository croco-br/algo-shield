<template>
  <v-chip
    :color="mappedColor"
    :variant="mappedVariant"
    :size="mappedSize"
    :rounded="rounded ? 'pill' : undefined"
  >
    <slot />
  </v-chip>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  variant?: 'success' | 'warning' | 'danger' | 'info' | 'default'
  size?: 'sm' | 'md' | 'lg'
  outline?: boolean
  rounded?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'default',
  size: 'md',
  outline: false,
  rounded: false,
})

// Map variants to Vuetify colors
const mappedColor = computed(() => {
  switch (props.variant) {
    case 'success':
      return 'success'
    case 'warning':
      return 'warning'
    case 'danger':
      return 'error'
    case 'info':
      return 'info'
    default:
      return 'default'
  }
})

// Map outline to Vuetify variant
const mappedVariant = computed(() => {
  return props.outline ? 'outlined' : 'flat'
})

// Map sizes
const mappedSize = computed(() => {
  switch (props.size) {
    case 'sm':
      return 'x-small'
    case 'lg':
      return 'large'
    default:
      return 'small'
  }
})
</script>
