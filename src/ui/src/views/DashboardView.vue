<template>
  <div class="dashboard-container">
    <div class="mb-8">
      <h2 class="text-3xl font-bold text-neutral-900 mb-2">Transactions</h2>
      <p class="text-neutral-600">View and manage transaction records</p>
    </div>

    <LoadingSpinner v-if="loading" text="Loading transactions..." :centered="false" />

    <ErrorMessage
      v-else-if="error"
      title="Error loading transactions"
      :message="error"
      retryable
      @retry="loadTransactions"
    />

    <!-- Transaction Table -->
    <TransactionTable
      v-else
      :data="transactions"
      :pageSize="50"
      @row-click="openTransactionDetail"
    />

    <!-- Transaction Detail Modal -->
    <TransactionDetailModal
      v-model="showDetailModal"
      :transaction="selectedTransaction"
      @open-escalation="openEscalationModal"
    />

    <!-- Risk Escalation Modal -->
    <RiskEscalationModal
      v-model="showEscalationModal"
      :transaction="selectedTransaction"
      @submit="handleEscalationSubmit"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '@/lib/api'
import type { Transaction } from '@/types/transaction'
import TransactionTable from '@/components/TransactionTable.vue'
import TransactionDetailModal from '@/components/TransactionDetailModal.vue'
import RiskEscalationModal from '@/components/RiskEscalationModal.vue'
import LoadingSpinner from '@/components/LoadingSpinner.vue'
import ErrorMessage from '@/components/ErrorMessage.vue'

const loading = ref(true)
const error = ref('')
const showDetailModal = ref(false)
const showEscalationModal = ref(false)
const selectedTransaction = ref<Transaction | null>(null)
const transactions = ref<Transaction[]>([])

const loadTransactions = async () => {
  loading.value = true
  error.value = ''
  try {
    const response = await api.get<{ transactions: Transaction[] }>('/api/v1/transactions?limit=50&offset=0')
    transactions.value = response.transactions || []
  } catch (e: any) {
    error.value = e.message || 'Failed to load transactions'
    console.error('Error loading transactions:', e)
  } finally {
    loading.value = false
  }
}

const openTransactionDetail = (transaction: Transaction) => {
  selectedTransaction.value = transaction
  showDetailModal.value = true
}

const openEscalationModal = () => {
  showDetailModal.value = false
  setTimeout(() => {
    showEscalationModal.value = true
  }, 200)
}

const handleEscalationSubmit = (data: any) => {
  // Note: Escalation endpoint doesn't exist in backend yet
  alert('Transaction escalation feature not yet implemented in backend')
}

onMounted(() => {
  loadTransactions()
})
</script>

<style scoped>
/* Dashboard uses Tailwind's grid system with harmonious spacing */
</style>
