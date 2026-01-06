<template>
  <div
    v-if="modelValue"
    class="modal-overlay fixed inset-0 bg-black bg-opacity-50 z-modal-backdrop flex items-center justify-center p-4"
    @click.self="close"
  >
    <div
      class="modal-content bg-white rounded-lg shadow-2xl w-full max-w-[800px] max-h-[90vh] overflow-hidden animate-fade-in z-modal"
    >
      <!-- Header -->
      <div class="modal-header px-6 py-4 border-b border-neutral-200 flex items-center justify-between bg-neutral-50">
        <div>
          <h2 class="text-xl font-bold text-neutral-900">Transaction Details</h2>
          <p class="text-sm text-neutral-600 mt-1">External ID: {{ transaction?.external_id }}</p>
        </div>
        <BaseButton variant="ghost" size="sm" @click="close" prepend-icon="fa-xmark" />
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
            <v-icon :icon="getTabIcon(tab.icon)" size="small" class="mr-2" />
            {{ tab.label }}
          </button>
        </div>
      </div>

      <!-- Content -->
      <div class="modal-body px-8 py-8 overflow-y-auto max-h-[calc(90vh-200px)]">
        <!-- Overview Tab -->
        <div v-if="activeTab === 'overview'" class="space-y-8">
          <div class="grid grid-cols-2 gap-8">
            <div>
              <label class="text-sm font-semibold text-neutral-700">Origin</label>
              <p class="text-base text-neutral-900 mt-1">{{ transaction?.origin }}</p>
            </div>
            <div>
              <label class="text-sm font-semibold text-neutral-700">Destination</label>
              <p class="text-base text-neutral-900 mt-1">{{ transaction?.destination }}</p>
            </div>
            <div>
              <label class="text-sm font-semibold text-neutral-700">Amount</label>
              <p class="text-base font-bold text-neutral-900 mt-1">{{ formatCurrency(transaction?.amount || 0, transaction?.currency) }}</p>
            </div>
            <div>
              <label class="text-sm font-semibold text-neutral-700">Date</label>
              <p class="text-base text-neutral-900 mt-1">{{ formatDate(transaction?.created_at || '') }}</p>
            </div>
            <div>
              <label class="text-sm font-semibold text-neutral-700">Status</label>
              <span
                :class="[
                  'inline-block mt-1 px-3 py-1 rounded-full text-sm font-semibold',
                  getStatusBadgeClass(transaction?.status || 'pending')
                ]"
              >
                {{ formatStatus(transaction?.status || 'pending') }}
              </span>
            </div>
            <div>
              <label class="text-sm font-semibold text-neutral-700">Type</label>
              <p class="text-base text-neutral-900 mt-1">{{ transaction?.type }}</p>
            </div>
          </div>

          <div class="border-t border-neutral-200 pt-8">
            <h4 class="text-sm font-semibold text-neutral-700 mb-4">Matched Rules</h4>
            <div v-if="transaction?.matched_rules && transaction.matched_rules.length > 0" class="space-y-2">
              <div
                v-for="(rule, index) in transaction.matched_rules"
                :key="index"
                class="flex items-center justify-between py-2 px-3 bg-neutral-50 rounded"
              >
                <span class="text-sm text-neutral-900">{{ rule }}</span>
              </div>
            </div>
            <div v-else class="text-sm text-neutral-500 py-2">
              No rules matched
            </div>
          </div>
        </div>

        <!-- Escalation Tab -->
        <div v-if="activeTab === 'escalation'" class="space-y-6">
          <div class="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
            <div class="flex items-start gap-3">
              <v-icon icon="fa-triangle-exclamation" color="warning" size="small" class="mt-0.5" />
              <div>
                <h4 class="text-sm font-semibold text-yellow-900">Escalation Required</h4>
                <p class="text-sm text-yellow-700 mt-1">
                  This transaction has been flagged for review and requires escalation to a senior analyst.
                </p>
              </div>
            </div>
          </div>

          <BaseButton
            @click="$emit('open-escalation')"
            full-width
            prepend-icon="fa-arrow-up"
          >
            Escalate Transaction
          </BaseButton>

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
        <BaseButton variant="ghost" @click="close">
          Close
        </BaseButton>
        <BaseButton prepend-icon="fa-check">
          Mark as Reviewed
        </BaseButton>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { Transaction } from '@/types/transaction'
import BaseButton from '@/components/BaseButton.vue'

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
  { key: 'overview', label: 'Overview', icon: 'fa-circle-info' },
  { key: 'escalation', label: 'Escalation', icon: 'fa-flag' },
]

const getTabIcon = (icon: string): string => {
  return icon
}

const activeTab = ref('overview')

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

const capitalize = (str: string) => {
  return str.charAt(0).toUpperCase() + str.slice(1)
}

const formatCurrency = (amount: number, currency: string = 'USD') => {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: currency,
  }).format(amount)
}

const getStatusBadgeClass = (status: string) => {
  const classes = {
    pending: 'bg-yellow-100 text-yellow-800',
    approved: 'bg-green-100 text-green-800',
    rejected: 'bg-red-100 text-red-800',
    in_review: 'bg-blue-100 text-blue-800',
  }
  return classes[status.toLowerCase() as keyof typeof classes] || classes.pending
}

const formatStatus = (status: string) => {
  if (status === 'in_review') {
    return 'In Review'
  }
  return capitalize(status)
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
