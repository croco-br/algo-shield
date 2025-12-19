<template>
  <div class="transaction-table-container bg-white rounded shadow-card">
    <!-- Header -->
    <div class="px-6 py-4 border-b border-neutral-200 bg-neutral-50">
      <div class="flex items-center justify-between">
        <h3 class="text-lg font-semibold text-neutral-900">Recent Transactions</h3>
        <div class="flex items-center gap-3">
          <button
            class="px-3 py-1.5 text-sm font-medium text-neutral-600 hover:text-neutral-900 hover:bg-neutral-100 rounded transition-colors"
          >
            <i class="fas fa-filter mr-2"></i>
            Filter
          </button>
          <button
            class="px-3 py-1.5 text-sm font-medium text-neutral-600 hover:text-neutral-900 hover:bg-neutral-100 rounded transition-colors"
          >
            <i class="fas fa-download mr-2"></i>
            Export
          </button>
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
              class="px-6 py-3 text-left text-xs font-semibold text-neutral-600 uppercase tracking-wider"
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
              index % 2 === 0 ? 'bg-white' : 'bg-neutral-25'
            ]"
            @click="$emit('row-click', transaction)"
          >
            <td class="px-6 py-4 whitespace-nowrap text-sm font-mono text-neutral-900">
              {{ transaction.id }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-neutral-600">
              {{ formatDate(transaction.date) }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm font-semibold text-neutral-900">
              {{ formatCurrency(transaction.amount) }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-neutral-900">
              {{ transaction.customer }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <span
                :class="[
                  'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-semibold',
                  getRiskBadgeClass(transaction.riskLevel)
                ]"
              >
                {{ transaction.riskLevel }}
              </span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <span
                :class="[
                  'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-semibold',
                  getStatusBadgeClass(transaction.status)
                ]"
              >
                {{ transaction.status }}
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
          <button
            :disabled="currentPage === 1"
            @click="prevPage"
            class="px-3 py-1.5 text-sm font-medium rounded transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            :class="currentPage === 1 ? 'text-neutral-400' : 'text-neutral-600 hover:bg-neutral-200'"
          >
            <i class="fas fa-chevron-left"></i>
          </button>

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

          <button
            :disabled="currentPage === totalPages"
            @click="nextPage"
            class="px-3 py-1.5 text-sm font-medium rounded transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            :class="currentPage === totalPages ? 'text-neutral-400' : 'text-neutral-600 hover:bg-neutral-200'"
          >
            <i class="fas fa-chevron-right"></i>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

interface Transaction {
  id: string
  date: string
  amount: number
  customer: string
  riskLevel: 'Low' | 'Medium' | 'High'
  status: 'Pending' | 'Approved' | 'Rejected' | 'Under Review'
}

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
  { key: 'id', label: 'ID' },
  { key: 'date', label: 'Date' },
  { key: 'amount', label: 'Amount' },
  { key: 'customer', label: 'Customer' },
  { key: 'riskLevel', label: 'Risk Level' },
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

const formatCurrency = (amount: number) => {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
  }).format(amount)
}

const getRiskBadgeClass = (level: string) => {
  const classes = {
    Low: 'bg-green-100 text-green-800',
    Medium: 'bg-orange-100 text-orange-800',
    High: 'bg-red-100 text-red-800',
  }
  return classes[level as keyof typeof classes] || classes.Low
}

const getStatusBadgeClass = (status: string) => {
  const classes = {
    Pending: 'bg-yellow-100 text-yellow-800',
    Approved: 'bg-green-100 text-green-800',
    Rejected: 'bg-red-100 text-red-800',
    'Under Review': 'bg-blue-100 text-blue-800',
  }
  return classes[status as keyof typeof classes] || classes.Pending
}
</script>

<style scoped>
.bg-neutral-25 {
  background-color: #fafafa;
}

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
