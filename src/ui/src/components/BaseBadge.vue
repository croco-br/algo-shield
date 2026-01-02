<template>
  <v-chip
    :color="mappedColor"
    :variant="mappedVariant"
    :size="mappedSize"
    :closable="closable"
    :rounded="rounded ? 'pill' : undefined"
    class="base-badge"
    @click:close="$emit('close')"
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
  closable?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'default',
  size: 'md',
  outline: false,
  rounded: false,
  closable: false,
})

defineEmits<{
  close: []
}>()

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

// Map sizes to match BaseButton
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
</script>

<style scoped>
.base-badge {
  font-family: var(--font-family-sans);
  font-weight: 500;
  letter-spacing: 0.01em;
}
</style>
