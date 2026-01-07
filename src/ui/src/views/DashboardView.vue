<template>
  <v-container fluid class="pa-8">
    <v-row>
      <v-col cols="12">
        <div class="mb-8">
          <div class="d-flex align-center gap-3 mb-2">
            <v-icon icon="fa-exchange-alt" size="large" color="primary" />
            <h2 class="text-h4 font-weight-bold">{{ $t('views.dashboard.title') }}</h2>
          </div>
          <p class="text-body-1 text-grey-darken-1">{{ $t('views.dashboard.subtitle') }}</p>
        </div>

        <LoadingSpinner v-if="loading" :text="$t('views.dashboard.loading')" :centered="false" />

        <ErrorMessage
          v-else-if="error"
          :title="$t('views.dashboard.errorTitle')"
          :message="error"
          retryable
          @retry="loadTransactions"
        />

        <!-- Transaction Table -->
        <TransactionTable
          v-else
          :data="transactions"
          :pageSize="50"
        />
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { i18n } from '@/plugins/i18n'
import { api } from '@/lib/api'
import type { Transaction } from '@/types/transaction'
import TransactionTable from '@/components/TransactionTable.vue'
import LoadingSpinner from '@/components/LoadingSpinner.vue'
import ErrorMessage from '@/components/ErrorMessage.vue'

const t = i18n.global.t

const loading = ref(true)
const error = ref('')
const transactions = ref<Transaction[]>([])

const loadTransactions = async () => {
  loading.value = true
  error.value = ''
  try {
    const response = await api.get<{ transactions: Transaction[] }>('/api/v1/transactions?limit=50&offset=0')
    transactions.value = response.transactions || []
  } catch (e: any) {
    error.value = e.message || t('views.dashboard.errorRetry')
    console.error('Error loading transactions:', e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadTransactions()
})
</script>
