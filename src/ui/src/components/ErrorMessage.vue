<template>
  <div v-if="message" :class="containerClasses">
    <div class="flex items-start gap-3">
      <!-- Icon -->
      <div class="flex-shrink-0">
        <svg v-if="variant === 'error'" class="w-5 h-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
        </svg>
        <svg v-else-if="variant === 'warning'" class="w-5 h-5 text-yellow-400" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
        </svg>
        <svg v-else class="w-5 h-5 text-blue-400" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd" />
        </svg>
      </div>

      <!-- Content -->
      <div class="flex-1 min-w-0">
        <h3 v-if="title" :class="titleClasses">{{ title }}</h3>
        <div :class="messageClasses">
          <p>{{ message }}</p>
        </div>
        <div v-if="$slots.actions || retryable" class="mt-4">
          <slot name="actions">
            <button
              v-if="retryable"
              @click="$emit('retry')"
              :class="retryButtonClasses"
            >
              Try again
            </button>
          </slot>
        </div>
      </div>

      <!-- Close button -->
      <button
        v-if="dismissible"
        @click="$emit('dismiss')"
        :class="closeButtonClasses"
        aria-label="Dismiss"
      >
        ×
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

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

const containerClasses = computed(() => {
  const classes = ['rounded-md p-4 border']

  if (props.variant === 'error') {
    classes.push('bg-red-50 border-red-200')
  } else if (props.variant === 'warning') {
    classes.push('bg-yellow-50 border-yellow-200')
  } else {
    classes.push('bg-blue-50 border-blue-200')
  }

  return classes.join(' ')
})

const titleClasses = computed(() => {
  const classes = ['text-sm font-medium mb-1']

  if (props.variant === 'error') {
    classes.push('text-red-800')
  } else if (props.variant === 'warning') {
    classes.push('text-yellow-800')
  } else {
    classes.push('text-blue-800')
  }

  return classes.join(' ')
})

const messageClasses = computed(() => {
  const classes = ['text-sm']

  if (props.variant === 'error') {
    classes.push('text-red-700')
  } else if (props.variant === 'warning') {
    classes.push('text-yellow-700')
  } else {
    classes.push('text-blue-700')
  }

  return classes.join(' ')
})

const retryButtonClasses = computed(() => {
  const classes = ['text-sm font-medium underline hover:no-underline']

  if (props.variant === 'error') {
    classes.push('text-red-800 hover:text-red-900')
  } else if (props.variant === 'warning') {
    classes.push('text-yellow-800 hover:text-yellow-900')
  } else {
    classes.push('text-blue-800 hover:text-blue-900')
  }

  return classes.join(' ')
})

const closeButtonClasses = computed(() => {
  const classes = ['flex-shrink-0 text-2xl leading-none transition-colors']

  if (props.variant === 'error') {
    classes.push('text-red-400 hover:text-red-600')
  } else if (props.variant === 'warning') {
    classes.push('text-yellow-400 hover:text-yellow-600')
  } else {
    classes.push('text-blue-400 hover:text-blue-600')
  }

  return classes.join(' ')
})
</script>
