<template>
  <div
    class="kpi-card bg-white rounded shadow-card hover:shadow-card-hover transition-all duration-200 hover:-translate-y-0.5"
    :class="{ 'animate-pulse': loading }"
  >
    <div class="p-8">
      <div class="flex items-center justify-between mb-4">
        <div class="text-sm font-medium text-neutral-600">{{ label }}</div>
        <div
          v-if="icon"
          class="w-10 h-10 rounded-full flex items-center justify-center"
          :class="iconBgClass"
        >
          <i :class="[icon, iconColorClass, 'text-lg']"></i>
        </div>
      </div>

      <div class="mb-2">
        <div class="text-3xl font-bold text-neutral-900">
          {{ loading ? '—' : formattedValue }}
        </div>
      </div>

      <div class="flex items-center gap-2">
        <div
          v-if="trend"
          class="flex items-center gap-1 text-sm font-medium"
          :class="trendColorClass"
        >
          <i :class="trendIconClass"></i>
          <span>{{ Math.abs(trend) }}%</span>
        </div>
        <div v-if="subtitle" class="text-sm text-neutral-500">
          {{ subtitle }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  label: string
  value: number | string
  icon?: string
  trend?: number
  subtitle?: string
  variant?: 'default' | 'success' | 'warning' | 'danger' | 'info'
  loading?: boolean
  formatter?: (value: number | string) => string
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'default',
  loading: false,
})

const formattedValue = computed(() => {
  if (props.formatter && typeof props.value === 'number') {
    return props.formatter(props.value)
  }
  return props.value.toLocaleString()
})

const iconBgClass = computed(() => {
  const classes = {
    default: 'bg-neutral-100',
    success: 'bg-green-100',
    warning: 'bg-orange-100',
    danger: 'bg-red-100',
    info: 'bg-teal-100',
  }
  return classes[props.variant]
})

const iconColorClass = computed(() => {
  const classes = {
    default: 'text-neutral-600',
    success: 'text-green-600',
    warning: 'text-orange-600',
    danger: 'text-red-600',
    info: 'text-teal-600',
  }
  return classes[props.variant]
})

const trendColorClass = computed(() => {
  if (!props.trend) return ''
  return props.trend > 0 ? 'text-green-600' : 'text-red-600'
})

const trendIconClass = computed(() => {
  if (!props.trend) return ''
  return props.trend > 0 ? 'fas fa-arrow-up' : 'fas fa-arrow-down'
})
</script>

<style scoped>
.kpi-card {
  box-shadow: 0px 4px 12px rgba(0, 0, 0, 0.1);
}

.kpi-card:hover {
  box-shadow: 0px 6px 16px rgba(0, 0, 0, 0.15);
}

.shadow-card {
  box-shadow: 0px 4px 12px rgba(0, 0, 0, 0.1);
}

.shadow-card-hover {
  box-shadow: 0px 6px 16px rgba(0, 0, 0, 0.15);
}
</style>
