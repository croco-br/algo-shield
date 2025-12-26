<template>
  <Teleport to="body">
    <Transition name="modal">
      <div
        v-if="modelValue"
        class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-modal-backdrop p-4"
        @click="handleBackdropClick"
      >
        <div
          :class="modalClasses"
          @click.stop
          role="dialog"
          aria-modal="true"
        >
          <!-- Header -->
          <div v-if="title || $slots.header" class="flex items-center justify-between p-6 border-b border-gray-200">
            <slot name="header">
              <h2 class="text-2xl font-semibold text-gray-900">{{ title }}</h2>
            </slot>
            <button
              v-if="closable"
              @click="close"
              class="text-gray-400 hover:text-gray-600 transition-colors w-8 h-8 flex items-center justify-center text-2xl leading-none"
              aria-label="Close modal"
            >
              ×
            </button>
          </div>

          <!-- Body -->
          <div class="p-6 overflow-y-auto" :style="{ maxHeight: maxContentHeight }">
            <slot />
          </div>

          <!-- Footer -->
          <div v-if="$slots.footer" class="flex items-center justify-end gap-4 p-6 border-t border-gray-200">
            <slot name="footer" />
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
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

const modalClasses = computed(() => {
  const classes = [
    'bg-white rounded-lg shadow-xl z-modal',
    'w-full mx-4',
    'flex flex-col',
  ]

  // Size classes
  if (props.size === 'sm') {
    classes.push('max-w-md')
  } else if (props.size === 'lg') {
    classes.push('max-w-3xl')
  } else if (props.size === 'xl') {
    classes.push('max-w-5xl')
  } else {
    classes.push('max-w-2xl')
  }

  return classes.join(' ')
})

const maxContentHeight = computed(() => {
  return props.maxHeight
})

const close = () => {
  emit('update:modelValue', false)
  emit('close')
}

const handleBackdropClick = () => {
  if (props.closeOnBackdrop) {
    close()
  }
}
</script>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-active .bg-white,
.modal-leave-active .bg-white {
  transition: transform 0.3s ease;
}

.modal-enter-from .bg-white,
.modal-leave-to .bg-white {
  transform: scale(0.95);
}
</style>
