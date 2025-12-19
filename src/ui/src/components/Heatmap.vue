<template>
  <div class="heatmap-container bg-white rounded shadow-card p-6">
    <div class="flex items-center justify-between mb-4">
      <h3 class="text-lg font-semibold text-neutral-900">Risk Density Heatmap</h3>
      <div class="flex items-center gap-4">
        <div class="flex items-center gap-2 text-sm">
          <span class="text-neutral-600">Legend:</span>
          <div class="flex items-center gap-1">
            <div class="w-4 h-4 rounded bg-teal-100"></div>
            <span class="text-xs text-neutral-500">Low</span>
          </div>
          <div class="flex items-center gap-1">
            <div class="w-4 h-4 rounded bg-teal-400"></div>
            <span class="text-xs text-neutral-500">Medium</span>
          </div>
          <div class="flex items-center gap-1">
            <div class="w-4 h-4 rounded bg-teal-600"></div>
            <span class="text-xs text-neutral-500">High</span>
          </div>
        </div>
      </div>
    </div>

    <div class="heatmap-grid" :style="{ height: height + 'px' }">
      <div
        v-for="(cell, index) in cells"
        :key="index"
        class="heatmap-cell group cursor-pointer transition-all duration-200"
        :class="getCellClass(cell.value)"
        :style="getCellStyle(index)"
        @mouseenter="hoveredCell = cell"
        @mouseleave="hoveredCell = null"
        @click="$emit('cell-click', cell)"
      >
        <div class="cell-tooltip opacity-0 group-hover:opacity-100 transition-opacity">
          <div class="text-xs font-semibold">{{ cell.label }}</div>
          <div class="text-xs">Risk: {{ cell.value }}</div>
        </div>
      </div>
    </div>

    <!-- Tooltip -->
    <div
      v-if="hoveredCell"
      class="fixed bg-neutral-900 text-white px-3 py-2 rounded shadow-lg text-xs z-50 pointer-events-none"
      :style="tooltipStyle"
    >
      <div class="font-semibold">{{ hoveredCell.label }}</div>
      <div>Risk Level: {{ hoveredCell.value }}</div>
      <div>Transactions: {{ hoveredCell.count || 0 }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'

interface HeatmapCell {
  label: string
  value: number
  count?: number
  timestamp?: string
}

interface Props {
  height?: number
  rows?: number
  cols?: number
  data?: HeatmapCell[]
}

const props = withDefaults(defineProps<Props>(), {
  height: 400,
  rows: 10,
  cols: 24,
})

const emit = defineEmits<{
  'cell-click': [cell: HeatmapCell]
}>()

const hoveredCell = ref<HeatmapCell | null>(null)
const cells = ref<HeatmapCell[]>([])

// Generate sample data if not provided
const generateSampleData = () => {
  if (props.data) {
    cells.value = props.data
    return
  }

  const generated: HeatmapCell[] = []
  const hours = ['00:00', '01:00', '02:00', '03:00', '04:00', '05:00', '06:00', '07:00', '08:00', '09:00', '10:00', '11:00', '12:00', '13:00', '14:00', '15:00', '16:00', '17:00', '18:00', '19:00', '20:00', '21:00', '22:00', '23:00']
  const days = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun', 'Week 1', 'Week 2', 'Week 3']

  for (let row = 0; row < props.rows; row++) {
    for (let col = 0; col < props.cols; col++) {
      const value = Math.floor(Math.random() * 100)
      const label = `${days[row % days.length]} ${hours[col % hours.length]}`
      const count = Math.floor(Math.random() * 500)
      generated.push({ label, value, count })
    }
  }

  cells.value = generated
}

const getCellClass = (value: number) => {
  if (value < 30) return 'bg-teal-100 hover:bg-teal-200'
  if (value < 60) return 'bg-teal-400 hover:bg-teal-500'
  return 'bg-teal-600 hover:bg-teal-700'
}

const getCellStyle = (index: number) => {
  const row = Math.floor(index / props.cols)
  const col = index % props.cols
  const cellWidth = `calc(100% / ${props.cols})`
  const cellHeight = `calc(100% / ${props.rows})`

  return {
    gridColumn: col + 1,
    gridRow: row + 1,
    width: cellWidth,
    height: cellHeight,
  }
}

const tooltipStyle = computed(() => {
  // This would be calculated based on mouse position in a real implementation
  return {
    top: '50%',
    left: '50%',
  }
})

onMounted(() => {
  generateSampleData()
})
</script>

<style scoped>
.heatmap-grid {
  display: grid;
  gap: 2px;
  background: #f0f0f0;
  padding: 2px;
  border-radius: 4px;
  position: relative;
}

.heatmap-cell {
  position: relative;
  border-radius: 2px;
}

.cell-tooltip {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  pointer-events: none;
  z-index: 10;
}

.shadow-card {
  box-shadow: 0px 4px 12px rgba(0, 0, 0, 0.1);
}
</style>
