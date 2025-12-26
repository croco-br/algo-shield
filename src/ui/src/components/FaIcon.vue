<template>
  <i :class="iconClass" :style="iconStyle"></i>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  icon: string
  size?: string | number
  color?: string
}

const props = withDefaults(defineProps<Props>(), {
  size: '1em',
  color: undefined,
})

const iconClass = computed(() => {
  // Support both 'fa-icon-name' and 'fas fa-icon-name' formats
  if (props.icon.startsWith('fa-')) {
    return `fas ${props.icon}`
  }
  if (props.icon.startsWith('fas ') || props.icon.startsWith('far ') || props.icon.startsWith('fal ')) {
    return props.icon
  }
  return `fas fa-${props.icon}`
})

const iconStyle = computed(() => {
  const style: Record<string, string> = {}
  if (props.size) {
    style.fontSize = typeof props.size === 'number' ? `${props.size}px` : props.size
  }
  if (props.color) {
    style.color = props.color
  }
  return style
})
</script>

