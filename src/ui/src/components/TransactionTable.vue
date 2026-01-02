<template>
  <div class="transaction-table-container bg-white rounded shadow-card">
    <!-- Header -->
    <div class="px-6 py-4 border-b border-neutral-200 bg-neutral-50">
      <div class="flex items-center justify-between">
        <h3 class="text-lg font-semibold text-neutral-900">Recent Transactions</h3>
        <div class="flex items-center gap-3">
          <BaseButton variant="ghost" size="sm" prepend-icon="fa-filter">
            Filter
          </BaseButton>
          <BaseButton variant="ghost" size="sm" prepend-icon="fa-download">
            Export
          </BaseButton>
        </div>
      </div>
    </div>

    <!-- Table -->
    <div class="overflow-x-auto">
      <table class="w-full">
        <thead class="bg-neutral-50 border-b border-neutral-200">
          <tr>
            <th
              v-for="column in columns"
              :key="column.key"
              class="px-8 py-4 text-left text-xs font-semibold text-neutral-600 uppercase tracking-wider"
            >
              {{ column.label }}
            </th>
          </tr>
        </thead>
        <tbody class="divide-y divide-neutral-100">
          <tr
            v-for="(transaction, index) in paginatedData"
            :key="transaction.id"
            :class="[
              'transition-colors cursor-pointer hover:bg-neutral-50',
              index % 2 === 0 ? 'bg-white' : 'bg-neutral-50'
            ]"
            @click="$emit('row-click', transaction)"
          >
            <td class="px-8 py-5 whitespace-nowrap text-sm font-mono text-neutral-900">
              {{ transaction.external_id }}
            </td>
            <td class="px-8 py-5 whitespace-nowrap text-sm text-neutral-600">
              {{ formatDate(transaction.created_at) }}
            </td>
            <td class="px-8 py-5 whitespace-nowrap text-sm font-semibold text-neutral-900">
              {{ formatCurrency(transaction.amount, transaction.currency) }}
            </td>
            <td class="px-8 py-5 whitespace-nowrap text-sm text-neutral-900">
              {{ transaction.from_account }}
            </td>
            <td class="px-8 py-5 whitespace-nowrap text-sm text-neutral-900">
              {{ transaction.to_account }}
            </td>
            <td class="px-8 py-5 whitespace-nowrap">
              <span
                :class="[
                  'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-semibold',
                  getRiskBadgeClass(transaction.risk_level)
                ]"
              >
                {{ capitalize(transaction.risk_level) }}
              </span>
            </td>
            <td class="px-8 py-5 whitespace-nowrap">
              <span
                :class="[
                  'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-semibold',
                  getStatusBadgeClass(transaction.status)
                ]"
              >
                {{ capitalize(transaction.status) }}
              </span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Pagination -->
    <div class="px-6 py-4 border-t border-neutral-200 bg-neutral-50">
      <div class="flex items-center justify-between">
        <div class="text-sm text-neutral-600">
          Showing {{ startIndex + 1 }} to {{ endIndex }} of {{ data.length }} transactions
        </div>
        <div class="flex items-center gap-2">
          <v-btn
            :disabled="currentPage === 1"
            @click="prevPage"
            icon
            variant="text"
            size="small"
          >
            <v-icon icon="fa-chevron-left" />
          </v-btn>

          <div class="flex items-center gap-1">
            <button
              v-for="page in visiblePages"
              :key="page"
              @click="goToPage(page)"
              :class="[
                'w-8 h-8 flex items-center justify-center text-sm font-medium rounded transition-colors',
                currentPage === page
                  ? 'bg-teal-600 text-white'
                  : 'text-neutral-600 hover:bg-neutral-200'
              ]"
            >
              {{ page }}
            </button>
          </div>

          <v-btn
            :disabled="currentPage === totalPages"
            @click="nextPage"
            icon
            variant="text"
            size="small"
          >
            <v-icon icon="fa-chevron-right" />
          </v-btn>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import type { Transaction } from '@/types/transaction'
import BaseButton from '@/components/BaseButton.vue'

interface Column {
  key: string
  label: string
}

interface Props {
  data: Transaction[]
  pageSize?: number
}

const props = withDefaults(defineProps<Props>(), {
  pageSize: 10,
})

const emit = defineEmits<{
  'row-click': [transaction: Transaction]
}>()

const columns: Column[] = [
  { key: 'external_id', label: 'External ID' },
  { key: 'created_at', label: 'Date' },
  { key: 'amount', label: 'Amount' },
  { key: 'from_account', label: 'From Account' },
  { key: 'to_account', label: 'To Account' },
  { key: 'risk_level', label: 'Risk Level' },
  { key: 'status', label: 'Status' },
]

const currentPage = ref(1)

const totalPages = computed(() => Math.ceil(props.data.length / props.pageSize))

const startIndex = computed(() => (currentPage.value - 1) * props.pageSize)
const endIndex = computed(() => Math.min(startIndex.value + props.pageSize, props.data.length))

const paginatedData = computed(() => {
  return props.data.slice(startIndex.value, endIndex.value)
})

const visiblePages = computed(() => {
  const pages: number[] = []
  const maxVisible = 5
  let start = Math.max(1, currentPage.value - Math.floor(maxVisible / 2))
  let end = Math.min(totalPages.value, start + maxVisible - 1)

  if (end - start < maxVisible - 1) {
    start = Math.max(1, end - maxVisible + 1)
  }

  for (let i = start; i <= end; i++) {
    pages.push(i)
  }

  return pages
})

const prevPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--
  }
}

const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
  }
}

const goToPage = (page: number) => {
  currentPage.value = page
}

const formatDate = (date: string) => {
  return new Date(date).toLocaleDateString('en-US', {
    month: 'short',
    day: '2-digit',
    year: 'numeric',
  })
}

const formatCurrency = (amount: number, currency: string = 'USD') => {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: currency,
  }).format(amount)
}

const capitalize = (str: string) => {
  return str.charAt(0).toUpperCase() + str.slice(1)
}

const getRiskBadgeClass = (level: string) => {
  const classes = {
    low: 'bg-green-100 text-green-800',
    medium: 'bg-orange-100 text-orange-800',
    high: 'bg-red-100 text-red-800',
  }
  return classes[level.toLowerCase() as keyof typeof classes] || classes.low
}

const getStatusBadgeClass = (status: string) => {
  const classes = {
    pending: 'bg-yellow-100 text-yellow-800',
    approved: 'bg-green-100 text-green-800',
    rejected: 'bg-red-100 text-red-800',
    review: 'bg-blue-100 text-blue-800',
  }
  return classes[status.toLowerCase() as keyof typeof classes] || classes.pending
}
</script>

<style scoped>
.shadow-card {
  box-shadow: 0px 4px 12px rgba(0, 0, 0, 0.1);
}

table {
  min-height: 44px;
}

tbody tr {
  height: 44px;
}
</style>
