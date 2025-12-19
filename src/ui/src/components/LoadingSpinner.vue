<template>
  <div :class="containerClasses">
    <div :class="spinnerClasses"></div>
    <p v-if="text" :class="textClasses">{{ text }}</p>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  size?: 'sm' | 'md' | 'lg'
  text?: string
  centered?: boolean
  fullscreen?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  size: 'md',
  centered: true,
  fullscreen: false,
})

const containerClasses = computed(() => {
  const classes = ['flex flex-col items-center gap-3']

  if (props.fullscreen) {
    classes.push('min-h-screen justify-center')
  } else if (props.centered) {
    classes.push('min-h-[50vh] justify-center')
  }

  return classes.join(' ')
})

const spinnerClasses = computed(() => {
  const classes = [
    'border-gray-200 border-t-indigo-600',
    'rounded-full animate-spin',
  ]

  if (props.size === 'sm') {
    classes.push('w-6 h-6 border-2')
  } else if (props.size === 'lg') {
    classes.push('w-12 h-12 border-4')
  } else {
    classes.push('w-10 h-10 border-4')
  }

  return classes.join(' ')
})

const textClasses = computed(() => {
  const classes = ['text-gray-500']

  if (props.size === 'sm') {
    classes.push('text-sm')
  } else if (props.size === 'lg') {
    classes.push('text-lg')
  } else {
    classes.push('text-base')
  }

  return classes.join(' ')
})
</script>
