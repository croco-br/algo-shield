<template>
  <button
    :type="type"
    :disabled="disabled || loading"
    :class="buttonClasses"
  >
    <span v-if="loading" class="inline-block w-4 h-4 border-2 border-current border-t-transparent rounded-full animate-spin mr-2"></span>
    <slot />
  </button>
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
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'primary',
  size: 'md',
  type: 'button',
  disabled: false,
  loading: false,
  fullWidth: false,
})

const buttonClasses = computed(() => {
  const classes = [
    'inline-flex items-center justify-center',
    'font-semibold rounded-lg',
    'transition-all duration-200',
    'focus:outline-none focus:ring-4 focus:ring-offset-2',
    'disabled:opacity-50 disabled:cursor-not-allowed',
    'active:scale-[0.98]',
  ]

  // Size classes
  if (props.size === 'sm') {
    classes.push('px-5 py-2.5 text-sm')
  } else if (props.size === 'lg') {
    classes.push('px-9 py-5 text-base')
  } else {
    classes.push('px-7 py-4 text-sm')
  }

  // Variant classes
  if (props.variant === 'primary') {
    classes.push(
      'bg-gradient-to-r from-blue-600 to-blue-700',
      'text-white shadow-lg shadow-blue-500/30',
      'hover:from-blue-700 hover:to-blue-800 hover:shadow-xl hover:shadow-blue-500/40',
      'focus:ring-blue-200'
    )
  } else if (props.variant === 'secondary') {
    classes.push(
      'bg-white text-slate-700 shadow-sm',
      'border-2 border-slate-200',
      'hover:bg-slate-50 hover:border-slate-300',
      'focus:ring-slate-200'
    )
  } else if (props.variant === 'danger') {
    classes.push(
      'bg-gradient-to-r from-red-600 to-red-700',
      'text-white shadow-lg shadow-red-500/30',
      'hover:from-red-700 hover:to-red-800 hover:shadow-xl hover:shadow-red-500/40',
      'focus:ring-red-200'
    )
  } else if (props.variant === 'ghost') {
    classes.push(
      'bg-transparent text-slate-700',
      'hover:bg-slate-50',
      'focus:ring-slate-200'
    )
  }

  // Full width
  if (props.fullWidth) {
    classes.push('w-full')
  }

  return classes.join(' ')
})
</script>
