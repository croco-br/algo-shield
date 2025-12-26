<template>
  <div :class="containerClasses">
    <v-progress-circular
      :size="mappedSize"
      :width="mappedWidth"
      color="primary"
      indeterminate
    />
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
  const classes = ['d-flex flex-column align-center gap-3']

  if (props.fullscreen) {
    classes.push('min-h-screen justify-center')
  } else if (props.centered) {
    classes.push('min-h-[50vh] justify-center')
  }

  return classes.join(' ')
})

const mappedSize = computed(() => {
  switch (props.size) {
    case 'sm':
      return 24
    case 'lg':
      return 48
    default:
      return 40
  }
})

const mappedWidth = computed(() => {
  switch (props.size) {
    case 'sm':
      return 2
    case 'lg':
      return 4
    default:
      return 4
  }
})

const textClasses = computed(() => {
  const classes = ['text-grey-darken-1']

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
