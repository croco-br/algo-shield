<template>
  <v-dialog
    :model-value="modelValue"
    :max-width="maxWidth"
    :persistent="!closeOnBackdrop"
    @update:model-value="handleUpdate"
  >
    <v-card>
      <!-- Header -->
      <v-card-title v-if="title || $slots.header" class="d-flex align-center justify-space-between">
        <slot name="header">
          <span>{{ title }}</span>
        </slot>
        <v-btn
          v-if="closable"
          icon="fa-xmark"
          variant="text"
          size="small"
          @click="close"
        />
      </v-card-title>

      <!-- Body -->
      <v-card-text :style="{ maxHeight: maxContentHeight, overflowY: 'auto' }">
        <slot />
      </v-card-text>

      <!-- Footer -->
      <v-card-actions v-if="$slots.footer" class="d-flex justify-end gap-2">
        <slot name="footer" />
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  modelValue: boolean
  title?: string
  size?: 'sm' | 'md' | 'lg' | 'xl'
  closable?: boolean
  closeOnBackdrop?: boolean
  maxHeight?: string
}

const props = withDefaults(defineProps<Props>(), {
  size: 'md',
  closable: true,
  closeOnBackdrop: true,
  maxHeight: '85vh',
})

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  'close': []
}>()

// Map size to Vuetify max-width
const maxWidth = computed(() => {
  switch (props.size) {
    case 'sm':
      return '500'
    case 'lg':
      return '900'
    case 'xl':
      return '1200'
    default:
      return '700' // md
  }
})

const maxContentHeight = computed(() => {
  return props.maxHeight
})

const close = () => {
  emit('update:modelValue', false)
  emit('close')
}

const handleUpdate = (value: boolean) => {
  emit('update:modelValue', value)
  if (!value) {
    emit('close')
  }
}
</script>
