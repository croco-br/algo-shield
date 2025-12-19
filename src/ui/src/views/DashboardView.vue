<template>
  <div class="dashboard-container">
    <!-- KPI Cards Row -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
      <KPICard
        label="Total Transactions"
        :value="kpiData.totalTransactions"
        icon="fas fa-exchange-alt"
        variant="info"
        :trend="5.2"
        subtitle="vs last month"
        :loading="loading"
      />
      <KPICard
        label="High-Risk Alerts"
        :value="kpiData.highRiskAlerts"
        icon="fas fa-exclamation-triangle"
        variant="danger"
        :trend="-12.5"
        subtitle="vs last month"
        :loading="loading"
      />
      <KPICard
        label="Pending Reviews"
        :value="kpiData.pendingReviews"
        icon="fas fa-clock"
        variant="warning"
        :trend="8.3"
        subtitle="requires attention"
        :loading="loading"
      />
      <KPICard
        label="Compliance Score"
        :value="kpiData.complianceScore + '%'"
        icon="fas fa-shield-alt"
        variant="success"
        :trend="2.1"
        subtitle="vs last month"
        :loading="loading"
        :formatter="(v) => v + '%'"
      />
    </div>

    <!-- Heatmap Row -->
    <div class="mb-8">
      <Heatmap
        :height="400"
        :rows="10"
        :cols="24"
        @cell-click="handleHeatmapCellClick"
      />
    </div>

    <!-- Transaction Table -->
    <TransactionTable
      :data="transactions"
      :pageSize="10"
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
import KPICard from '@/components/KPICard.vue'
import Heatmap from '@/components/Heatmap.vue'
import TransactionTable from '@/components/TransactionTable.vue'
import TransactionDetailModal from '@/components/TransactionDetailModal.vue'
import RiskEscalationModal from '@/components/RiskEscalationModal.vue'

interface Transaction {
  id: string
  date: string
  amount: number
  customer: string
  riskLevel: 'Low' | 'Medium' | 'High'
  status: 'Pending' | 'Approved' | 'Rejected' | 'Under Review'
}

interface KPIData {
  totalTransactions: number
  highRiskAlerts: number
  pendingReviews: number
  complianceScore: number
}

const loading = ref(true)
const showDetailModal = ref(false)
const showEscalationModal = ref(false)
const selectedTransaction = ref<Transaction | null>(null)

const kpiData = ref<KPIData>({
  totalTransactions: 12547,
  highRiskAlerts: 248,
  pendingReviews: 63,
  complianceScore: 94.7,
})

// Generate sample transaction data
const generateTransactions = (): Transaction[] => {
  const customers = [
    'Acme Corporation',
    'TechStart Inc.',
    'Global Ventures LLC',
    'Metro Bank',
    'Digital Solutions',
    'Enterprise Holdings',
    'Quantum Industries',
    'Nexus Group',
    'Stellar Financial',
    'Atlas Trading Co.',
  ]

  const statuses: Transaction['status'][] = ['Pending', 'Approved', 'Rejected', 'Under Review']
  const riskLevels: Transaction['riskLevel'][] = ['Low', 'Low', 'Low', 'Medium', 'Medium', 'High']

  const transactions: Transaction[] = []

  for (let i = 0; i < 50; i++) {
    const date = new Date()
    date.setDate(date.getDate() - Math.floor(Math.random() * 30))

    transactions.push({
      id: `TXN${String(10000 + i).padStart(6, '0')}`,
      date: date.toISOString(),
      amount: Math.floor(Math.random() * 100000) + 1000,
      customer: customers[Math.floor(Math.random() * customers.length)]!,
      riskLevel: riskLevels[Math.floor(Math.random() * riskLevels.length)]!,
      status: statuses[Math.floor(Math.random() * statuses.length)]!,
    })
  }

  return transactions.sort((a, b) => new Date(b.date).getTime() - new Date(a.date).getTime())
}

const transactions = ref<Transaction[]>([])

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
  console.log('Escalation submitted:', data)
  // Here you would typically make an API call
  // For now, just show a success message
  alert('Transaction escalated successfully!')
}

const handleHeatmapCellClick = (cell: any) => {
  console.log('Heatmap cell clicked:', cell)
}

onMounted(async () => {
  loading.value = true
  // Simulate API call
  await new Promise(resolve => setTimeout(resolve, 1000))

  transactions.value = generateTransactions()
  loading.value = false
})
</script>

<style scoped>
.dashboard-container {
  padding: 0;
}

/* Grid system with 24px gutter */
.grid {
  display: grid;
  gap: 24px;
}
</style>
