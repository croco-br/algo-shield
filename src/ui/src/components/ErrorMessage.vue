<template>
  <v-alert
    v-if="message"
    :type="mappedType"
    :title="title"
    :dismissible="dismissible"
    @click:close="$emit('dismiss')"
    class="mb-4"
  >
    <template v-if="!title">
      {{ message }}
    </template>
    <template v-else>
      <div>{{ message }}</div>
    </template>

    <template v-if="$slots.actions || retryable" #append>
      <slot name="actions">
        <BaseButton
          v-if="retryable"
          variant="ghost"
          size="sm"
          @click="$emit('retry')"
        >
          Try again
        </BaseButton>
      </slot>
    </template>
  </v-alert>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import BaseButton from '@/components/BaseButton.vue'

interface Props {
  variant?: 'error' | 'warning' | 'info'
  message: string
  title?: string
  dismissible?: boolean
  retryable?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'error',
  dismissible: true,
  retryable: false,
})

defineEmits<{
  'dismiss': []
  'retry': []
}>()

// Map variant to Vuetify alert type
const mappedType = computed(() => {
  switch (props.variant) {
    case 'warning':
      return 'warning'
    case 'info':
      return 'info'
    default:
      return 'error'
  }
})

// Expose for testing
defineExpose({
  mappedType,
})
</script>
