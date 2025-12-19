<template>
  <span :class="badgeClasses">
    <slot />
  </span>
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

const badgeClasses = computed(() => {
  const classes = [
    'inline-flex items-center justify-center',
    'font-medium',
    'transition-colors',
  ]

  // Size classes
  if (props.size === 'sm') {
    classes.push('px-2 py-0.5 text-xs')
  } else if (props.size === 'lg') {
    classes.push('px-4 py-2 text-base')
  } else {
    classes.push('px-3 py-1 text-sm')
  }

  // Rounded
  if (props.rounded) {
    classes.push('rounded-full')
  } else {
    classes.push('rounded')
  }

  // Variant classes
  if (props.outline) {
    // Outline variants
    if (props.variant === 'success') {
      classes.push('bg-white text-green-700 border border-green-300')
    } else if (props.variant === 'warning') {
      classes.push('bg-white text-yellow-700 border border-yellow-300')
    } else if (props.variant === 'danger') {
      classes.push('bg-white text-red-700 border border-red-300')
    } else if (props.variant === 'info') {
      classes.push('bg-white text-indigo-700 border border-indigo-300')
    } else {
      classes.push('bg-white text-gray-700 border border-gray-300')
    }
  } else {
    // Solid variants
    if (props.variant === 'success') {
      classes.push('bg-green-100 text-green-800')
    } else if (props.variant === 'warning') {
      classes.push('bg-yellow-100 text-yellow-800')
    } else if (props.variant === 'danger') {
      classes.push('bg-red-100 text-red-800')
    } else if (props.variant === 'info') {
      classes.push('bg-indigo-100 text-indigo-800')
    } else {
      classes.push('bg-gray-100 text-gray-800')
    }
  }

  return classes.join(' ')
})
</script>
