<template>
  <v-card>
    <v-card-title v-if="$slots.header">
      <slot name="header" />
    </v-card-title>

    <v-data-table
      :headers="normalizedHeaders"
      :items="data"
      :item-value="getItemKey"
      :items-per-page="-1"
      hide-default-footer
      class="elevation-0"
    >
      <template
        v-for="column in columns"
        :key="`cell-${column.key}`"
        #[`item.${column.key}`]="{ item }"
      >
        <slot
          :name="`cell-${column.key}`"
          :row="item"
          :value="getNestedValue(item, column.key)"
          :index="data.indexOf(item)"
        >
          {{ getNestedValue(item, column.key) }}
        </slot>
      </template>

      <template #no-data>
        <div v-if="!hideEmpty" class="pa-12 text-center text-grey">
          <slot name="empty">
            {{ emptyText }}
          </slot>
        </div>
      </template>
    </v-data-table>

    <v-card-actions v-if="$slots.footer">
      <slot name="footer" />
    </v-card-actions>
  </v-card>
</template>

<script setup lang="ts">
import { computed } from 'vue'

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

// Normalize columns to Vuetify headers format
const normalizedHeaders = computed(() => {
  return props.columns.map(column => ({
    title: column.label,
    key: column.key,
    width: column.width,
    align: 'start' as const,
    sortable: false,
  }))
})

const getItemKey = (item: any): string | number => {
  if (props.rowKey && item[props.rowKey] !== undefined) {
    return item[props.rowKey]
  }
  return String(item) || Math.random().toString()
}

const getNestedValue = (obj: any, path: string): any => {
  return path.split('.').reduce((current, key) => current?.[key], obj)
}
</script>
