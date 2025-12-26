<template>
  <div class="bg-white rounded-lg border border-gray-200 overflow-hidden">
    <div v-if="$slots.header" class="px-6 py-4 border-b border-gray-200">
      <slot name="header" />
    </div>

    <div class="overflow-x-auto">
      <table class="w-full border-collapse">
        <thead class="bg-gray-50">
          <tr>
            <th
              v-for="column in columns"
              :key="column.key"
              :class="[
                'px-6 py-4 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider border-b border-gray-200',
                column.headerClass
              ]"
              :style="column.width ? { width: column.width } : {}"
            >
              {{ column.label }}
            </th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="(row, index) in data"
            :key="getRowKey(row, index)"
            :class="['hover:bg-gray-50 transition-colors', rowClass]"
          >
            <td
              v-for="column in columns"
              :key="column.key"
              :class="[
                'px-6 py-5 border-b border-gray-200',
                column.cellClass
              ]"
            >
              <slot
                :name="`cell-${column.key}`"
                :row="row"
                :value="getNestedValue(row, column.key)"
                :index="index"
              >
                {{ getNestedValue(row, column.key) }}
              </slot>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="$slots.footer" class="px-6 py-4 border-t border-gray-200 bg-gray-50">
      <slot name="footer" />
    </div>

    <div v-if="data.length === 0 && !hideEmpty" class="px-6 py-12 text-center text-gray-500">
      <slot name="empty">
        {{ emptyText }}
      </slot>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Column {
  key: string
  label: string
  width?: string
  headerClass?: string
  cellClass?: string
}

interface Props {
  columns: Column[]
  data: any[]
  rowKey?: string
  rowClass?: string
  emptyText?: string
  hideEmpty?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  rowKey: 'id',
  emptyText: 'No data available',
  hideEmpty: false,
})

const getRowKey = (row: any, index: number): string | number => {
  if (props.rowKey && row[props.rowKey]) {
    return row[props.rowKey]
  }
  return index
}

const getNestedValue = (obj: any, path: string): any => {
  return path.split('.').reduce((current, key) => current?.[key], obj)
}
</script>
