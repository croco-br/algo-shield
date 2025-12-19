<template>
  <div
    v-if="modelValue"
    class="modal-overlay fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center p-4"
    @click.self="close"
  >
    <div
      class="modal-content bg-white rounded-lg shadow-2xl w-full max-w-[800px] max-h-[90vh] overflow-hidden animate-fade-in"
    >
      <!-- Header -->
      <div class="modal-header px-6 py-4 border-b border-neutral-200 flex items-center justify-between bg-neutral-50">
        <div>
          <h2 class="text-xl font-bold text-neutral-900">Transaction Details</h2>
          <p class="text-sm text-neutral-600 mt-1">ID: {{ transaction?.id }}</p>
        </div>
        <button
          @click="close"
          class="w-8 h-8 flex items-center justify-center rounded-full hover:bg-neutral-200 transition-colors"
        >
          <i class="fas fa-times text-neutral-600"></i>
        </button>
      </div>

      <!-- Tabs -->
      <div class="tabs border-b border-neutral-200 bg-white">
        <div class="flex px-6">
          <button
            v-for="tab in tabs"
            :key="tab.key"
            @click="activeTab = tab.key"
            :class="[
              'px-4 py-3 text-sm font-semibold transition-all duration-200 border-b-2',
              activeTab === tab.key
                ? 'text-teal-600 border-teal-600'
                : 'text-neutral-600 border-transparent hover:text-neutral-900'
            ]"
          >
            <i :class="[tab.icon, 'mr-2']"></i>
            {{ tab.label }}
          </button>
        </div>
      </div>

      <!-- Content -->
      <div class="modal-body px-6 py-6 overflow-y-auto max-h-[calc(90vh-200px)]">
        <!-- Overview Tab -->
        <div v-if="activeTab === 'overview'" class="space-y-6">
          <div class="grid grid-cols-2 gap-6">
            <div>
              <label class="text-sm font-semibold text-neutral-700">Customer</label>
              <p class="text-base text-neutral-900 mt-1">{{ transaction?.customer }}</p>
            </div>
            <div>
              <label class="text-sm font-semibold text-neutral-700">Amount</label>
              <p class="text-base font-bold text-neutral-900 mt-1">{{ formatCurrency(transaction?.amount || 0) }}</p>
            </div>
            <div>
              <label class="text-sm font-semibold text-neutral-700">Date</label>
              <p class="text-base text-neutral-900 mt-1">{{ formatDate(transaction?.date || '') }}</p>
            </div>
            <div>
              <label class="text-sm font-semibold text-neutral-700">Risk Level</label>
              <span
                :class="[
                  'inline-block mt-1 px-3 py-1 rounded-full text-sm font-semibold',
                  getRiskBadgeClass(transaction?.riskLevel || 'Low')
                ]"
              >
                {{ transaction?.riskLevel }}
              </span>
            </div>
          </div>

          <div class="border-t border-neutral-200 pt-6">
            <h4 class="text-sm font-semibold text-neutral-700 mb-3">Risk Factors</h4>
            <div class="space-y-2">
              <div class="flex items-center justify-between py-2 px-3 bg-neutral-50 rounded">
                <span class="text-sm text-neutral-900">High transaction value</span>
                <span class="text-sm font-semibold text-orange-600">+30 pts</span>
              </div>
              <div class="flex items-center justify-between py-2 px-3 bg-neutral-50 rounded">
                <span class="text-sm text-neutral-900">Unusual time of transaction</span>
                <span class="text-sm font-semibold text-orange-600">+20 pts</span>
              </div>
              <div class="flex items-center justify-between py-2 px-3 bg-neutral-50 rounded">
                <span class="text-sm text-neutral-900">New customer</span>
                <span class="text-sm font-semibold text-orange-600">+15 pts</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Documents Tab -->
        <div v-if="activeTab === 'documents'" class="space-y-4">
          <div
            v-for="(doc, index) in mockDocuments"
            :key="index"
            class="flex items-center justify-between p-4 border border-neutral-200 rounded-lg hover:border-teal-300 transition-colors cursor-pointer"
          >
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 bg-neutral-100 rounded flex items-center justify-center">
                <i class="fas fa-file-alt text-neutral-600"></i>
              </div>
              <div>
                <p class="text-sm font-semibold text-neutral-900">{{ doc.name }}</p>
                <p class="text-xs text-neutral-500">{{ doc.size }} • {{ doc.date }}</p>
              </div>
            </div>
            <button class="px-3 py-1.5 text-sm font-medium text-teal-600 hover:bg-teal-50 rounded transition-colors">
              <i class="fas fa-download mr-2"></i>
              Download
            </button>
          </div>
        </div>

        <!-- Escalation Tab -->
        <div v-if="activeTab === 'escalation'" class="space-y-4">
          <div class="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
            <div class="flex items-start gap-3">
              <i class="fas fa-exclamation-triangle text-yellow-600 mt-0.5"></i>
              <div>
                <h4 class="text-sm font-semibold text-yellow-900">Escalation Required</h4>
                <p class="text-sm text-yellow-700 mt-1">
                  This transaction has been flagged for review and requires escalation to a senior analyst.
                </p>
              </div>
            </div>
          </div>

          <button
            @click="$emit('open-escalation')"
            class="w-full px-4 py-3 bg-teal-600 text-white font-semibold rounded-lg hover:bg-teal-700 transition-colors"
          >
            <i class="fas fa-level-up-alt mr-2"></i>
            Escalate Transaction
          </button>

          <div class="border-t border-neutral-200 pt-4">
            <h4 class="text-sm font-semibold text-neutral-700 mb-3">Escalation History</h4>
            <div class="space-y-3">
              <div class="text-sm text-neutral-500 text-center py-4">
                No escalations yet
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Footer -->
      <div class="modal-footer px-6 py-4 border-t border-neutral-200 bg-neutral-50 flex items-center justify-end gap-3">
        <button
          @click="close"
          class="px-4 py-2 text-sm font-semibold text-neutral-700 hover:bg-neutral-200 rounded-lg transition-colors"
        >
          Close
        </button>
        <button
          class="px-4 py-2 bg-teal-600 text-white text-sm font-semibold rounded-lg hover:bg-teal-700 transition-colors"
        >
          <i class="fas fa-check mr-2"></i>
          Mark as Reviewed
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

interface Transaction {
  id: string
  date: string
  amount: number
  customer: string
  riskLevel: 'Low' | 'Medium' | 'High'
  status: string
}

interface Props {
  modelValue: boolean
  transaction: Transaction | null
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  'open-escalation': []
}>()

const tabs = [
  { key: 'overview', label: 'Overview', icon: 'fas fa-info-circle' },
  { key: 'documents', label: 'Documents', icon: 'fas fa-file-alt' },
  { key: 'escalation', label: 'Escalation', icon: 'fas fa-flag' },
]

const activeTab = ref('overview')

const mockDocuments = [
  { name: 'Transaction Receipt.pdf', size: '245 KB', date: 'Dec 15, 2024' },
  { name: 'Customer ID Verification.pdf', size: '1.2 MB', date: 'Dec 15, 2024' },
  { name: 'Bank Statement.pdf', size: '890 KB', date: 'Dec 14, 2024' },
]

const close = () => {
  emit('update:modelValue', false)
}

const formatDate = (date: string) => {
  return new Date(date).toLocaleDateString('en-US', {
    month: 'long',
    day: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
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
</script>

<style scoped>
.modal-overlay {
  animation: fadeIn 200ms ease-out;
}

.animate-fade-in {
  animation: slideDown 200ms ease-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
