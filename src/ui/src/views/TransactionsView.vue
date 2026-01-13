<template>
  <v-container fluid class="pa-8">
    <v-row>
      <v-col cols="12">
        <div class="d-flex justify-space-between align-center mb-8">
          <div>
            <div class="d-flex align-center gap-3 mb-2">
              <v-icon icon="fa-exchange-alt" size="large" color="primary" />
              <h2 class="text-h4 font-weight-bold">{{ $t('views.transactions.title') }}</h2>
              <v-chip v-if="isLive" color="success" size="small" class="ml-2">
                <v-icon icon="fa-circle" size="x-small" class="mr-1 pulse" />
                {{ $t('views.transactions.live') }}
              </v-chip>
            </div>
            <p class="text-body-1 text-grey-darken-1">{{ $t('views.transactions.subtitle') }}</p>
          </div>
          <div class="d-flex gap-2">
            <BaseButton 
              v-if="syntheticMode"
              variant="secondary"
              @click="openGenerateModal"
              prepend-icon="fa-bolt"
            >
              {{ $t('views.transactions.generateSynthetic') }}
            </BaseButton>
            <BaseButton 
              :variant="isLive ? 'secondary' : 'primary'"
              @click="toggleLiveUpdates"
              :prepend-icon="isLive ? 'fa-pause' : 'fa-play'"
            >
              {{ isLive ? $t('views.transactions.pause') : $t('views.transactions.live') }}
            </BaseButton>
            <BaseButton @click="loadTransactions" prepend-icon="fa-refresh">
              {{ $t('views.transactions.refresh') }}
            </BaseButton>
          </div>
        </div>

        <!-- Filters -->
        <v-card class="mb-6 pa-4" variant="outlined">
          <v-row>
            <v-col cols="12" md="3">
              <v-select
                v-model="filters.status"
                :items="statusOptions"
                :label="$t('views.transactions.filterStatus')"
                clearable
                density="compact"
                variant="outlined"
              />
            </v-col>
            <v-col cols="12" md="3">
              <v-select
                v-model="filters.schemaId"
                :items="schemaOptions"
                :label="$t('views.transactions.filterSchema')"
                clearable
                density="compact"
                variant="outlined"
              />
            </v-col>
            <v-col cols="12" md="3">
              <v-text-field
                v-model="filters.startDate"
                :label="$t('views.transactions.filterStartDate')"
                type="date"
                density="compact"
                variant="outlined"
                clearable
              />
            </v-col>
            <v-col cols="12" md="3">
              <v-text-field
                v-model="filters.endDate"
                :label="$t('views.transactions.filterEndDate')"
                type="date"
                density="compact"
                variant="outlined"
                clearable
              />
            </v-col>
          </v-row>
          <v-row align="center">
            <v-col cols="12" md="3">
              <v-text-field
                v-model.number="filters.minAmount"
                :label="$t('views.transactions.filterMinAmount')"
                type="number"
                density="compact"
                variant="outlined"
                clearable
                hide-details
              />
            </v-col>
            <v-col cols="12" md="3">
              <v-text-field
                v-model.number="filters.maxAmount"
                :label="$t('views.transactions.filterMaxAmount')"
                type="number"
                density="compact"
                variant="outlined"
                clearable
                hide-details
              />
            </v-col>
            <v-col cols="12" md="6" class="d-flex gap-2">
              <BaseButton @click="applyFilters" prepend-icon="fa-filter">
                {{ $t('views.transactions.applyFilters') }}
              </BaseButton>
              <BaseButton variant="ghost" @click="clearFilters" prepend-icon="fa-times">
                {{ $t('views.transactions.clearFilters') }}
              </BaseButton>
            </v-col>
          </v-row>
        </v-card>

        <LoadingSpinner v-if="loading" :text="$t('views.transactions.loading')" :centered="false" />

        <ErrorMessage
          v-else-if="error"
          :title="$t('views.transactions.errorTitle')"
          :message="error"
          retryable
          @retry="loadTransactions"
        />

        <!-- Transaction Table -->
        <v-card v-else variant="outlined">
          <v-data-table
            :headers="tableHeaders"
            :items="transactions"
            :items-per-page="itemsPerPage"
            :page="currentPage"
            :items-length="total"
            :items-per-page-options="itemsPerPageOptions"
            @update:page="handlePageChange"
            @update:items-per-page="handleItemsPerPageChange"
            class="elevation-0"
            server-items-length
          >
            <template #item.status="{ item }">
              <BaseBadge :variant="getStatusVariant(item.status)">
                {{ item.status }}
              </BaseBadge>
            </template>
            <template #item.amount="{ item }">
              {{ formatCurrency(item.amount, item.currency) }}
            </template>
            <template #item.created_at="{ item }">
              {{ formatDate(item.created_at) }}
            </template>
            <template #item.actions="{ item }">
              <div v-if="item.status === 'in_review'" class="d-flex gap-2">
                <BaseButton 
                  size="sm"
                  variant="primary"
                  @click="approveTransaction(item.id)"
                  prepend-icon="fa-check"
                >
                  {{ $t('views.transactions.approve') }}
                </BaseButton>
                <BaseButton 
                  size="sm"
                  variant="danger"
                  @click="rejectTransaction(item.id)"
                  prepend-icon="fa-times"
                >
                  {{ $t('views.transactions.reject') }}
                </BaseButton>
              </div>
            </template>
          </v-data-table>
        </v-card>
      </v-col>
    </v-row>

    <!-- Generate Synthetic Events Modal -->
    <BaseModal
      v-model="showGenerateModal"
      :title="$t('views.transactions.generateSyntheticTitle')"
      size="md"
    >
      <div class="mt-4">
        <p class="text-body-1 mb-4">
          {{ $t('views.transactions.generateSyntheticDescription') }}
        </p>

        <v-select
          v-model="generateSchemaId"
          :items="generateSchemaOptions"
          :label="$t('views.transactions.selectSchema')"
          density="compact"
          variant="outlined"
          class="mb-4"
        />

        <BaseInput
          v-model.number="generateCount"
          :label="$t('views.transactions.eventCount')"
          type="number"
          min="1"
          max="1000"
          :placeholder="$t('views.transactions.eventCountPlaceholder')"
          prepend-inner-icon="fa-hashtag"
          class="mb-4"
        />

        <v-alert
          v-if="generateResult"
          :type="generateResult.success ? 'success' : 'error'"
          variant="tonal"
          class="mt-4"
        >
          {{ generateResult.message }}
        </v-alert>
      </div>

      <template #footer>
        <BaseButton variant="ghost" @click="showGenerateModal = false" prepend-icon="fa-xmark">
          {{ $t('components.modal.cancel') }}
        </BaseButton>
        <BaseButton 
          @click="handleGenerateEvents" 
          :loading="generating" 
          :disabled="!generateSchemaId || !generateCount || generateCount < 1 || generateCount > 1000"
          prepend-icon="fa-bolt"
        >
          {{ $t('views.transactions.generate') }}
        </BaseButton>
      </template>
    </BaseModal>
  </v-container>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted, computed, watch } from 'vue'
import { useLocale } from '@/composables/useLocale'
import { useCurrency } from '@/composables/useCurrency'
import { useSystemModeStore } from '@/stores/systemMode'
import { api } from '@/lib/api'
import type { Transaction } from '@/types/transaction'
import BaseButton from '@/components/BaseButton.vue'
import BaseBadge from '@/components/BaseBadge.vue'
import BaseModal from '@/components/BaseModal.vue'
import BaseInput from '@/components/BaseInput.vue'
import LoadingSpinner from '@/components/LoadingSpinner.vue'
import ErrorMessage from '@/components/ErrorMessage.vue'

const { t } = useLocale()
const { formatCurrency } = useCurrency()
const systemModeStore = useSystemModeStore()

interface Schema {
  id: string
  name: string
}

const loading = ref(true)
const error = ref('')
const transactions = ref<Transaction[]>([])
const total = ref(0)
const schemas = ref<Schema[]>([])
const isLive = ref(false)
const syntheticMode = computed(() => systemModeStore.syntheticMode)
let pollingInterval: ReturnType<typeof setInterval> | null = null
let eventSource: EventSource | null = null

// Pagination state
const currentPage = ref(1)
const itemsPerPage = ref(50)
const itemsPerPageOptions = [10, 25, 50, 100, 200]

// Generate synthetic events state
const showGenerateModal = ref(false)
const generateSchemaId = ref<string | null>(null)
const generateCount = ref(10)
const generating = ref(false)
const generateResult = ref<{ success: boolean; message: string } | null>(null)

const filters = reactive({
  status: 'pending' as string | null, // Default filter: pending transactions
  schemaId: null as string | null,
  startDate: null as string | null,
  endDate: null as string | null,
  minAmount: null as number | null,
  maxAmount: null as number | null,
})

const statusOptions = computed(() => [
  { title: t('views.transactions.statusAll'), value: null },
  { title: t('views.transactions.statusApproved'), value: 'approved' },
  { title: t('views.transactions.statusRejected'), value: 'rejected' },
  { title: t('views.transactions.statusInReview'), value: 'in_review' },
  { title: t('views.transactions.statusPending'), value: 'pending' },
])

const schemaOptions = computed(() => [
  { title: t('views.transactions.allSchemas'), value: null },
  ...schemas.value.map((s: Schema) => ({ title: s.name, value: s.id }))
])

const generateSchemaOptions = computed(() => 
  schemas.value.map((s: Schema) => ({ title: s.name, value: s.id }))
)



const tableHeaders = computed(() => [
  { title: t('views.transactions.tableExternalId'), key: 'external_id', sortable: true },
  { title: t('views.transactions.tableStatus'), key: 'status', sortable: true },
  { title: t('views.transactions.tableAmount'), key: 'amount', sortable: true },
  { title: t('views.transactions.tableOrigin'), key: 'origin', sortable: true },
  { title: t('views.transactions.tableDestination'), key: 'destination', sortable: true },
  { title: t('views.transactions.tableType'), key: 'type', sortable: true },
  { title: t('views.transactions.tableCreated'), key: 'created_at', sortable: true },
  { title: t('views.transactions.tableActions'), key: 'actions', sortable: false },
])

onMounted(async () => {
  await systemModeStore.loadMode()
  await Promise.all([loadTransactions(), loadSchemas()])
})

// Watch for synthetic mode changes
watch(() => systemModeStore.syntheticMode, () => {
  // Modal will be hidden automatically if synthetic mode is disabled
  if (!systemModeStore.syntheticMode && showGenerateModal.value) {
    showGenerateModal.value = false
  }
})

onUnmounted(() => {
  stopLiveUpdates()
})

async function loadSchemas() {
  try {
    const response = await api.get<{ schemas: Schema[] }>('/api/v1/schemas')
    schemas.value = response?.schemas || []
  } catch (e) {
    console.error('Error loading schemas:', e)
  }
}

async function loadTransactions() {
  loading.value = true
  error.value = ''
  try {
    const params = new URLSearchParams()
    const offset = (currentPage.value - 1) * itemsPerPage.value
    params.append('limit', String(itemsPerPage.value))
    params.append('offset', String(offset))

    if (filters.status) params.append('status', filters.status)
    if (filters.schemaId) params.append('schema_id', filters.schemaId)
    if (filters.startDate) params.append('start_date', new Date(filters.startDate).toISOString())
    if (filters.endDate) params.append('end_date', new Date(filters.endDate).toISOString())
    if (filters.minAmount) params.append('min_amount', String(filters.minAmount))
    if (filters.maxAmount) params.append('max_amount', String(filters.maxAmount))

    const response = await api.get<{ transactions: Transaction[]; total?: number }>(`/api/v1/transactions?${params.toString()}`)
    transactions.value = response.transactions || []
    total.value = response.total || transactions.value.length
  } catch (e: any) {
    error.value = e.message || 'Failed to load transactions'
    console.error('Error loading transactions:', e)
  } finally {
    loading.value = false
  }
}

function applyFilters() {
  currentPage.value = 1
  loadTransactions()
}

function clearFilters() {
  filters.status = null
  filters.schemaId = null
  filters.startDate = null
  filters.endDate = null
  filters.minAmount = null
  filters.maxAmount = null
  currentPage.value = 1
  loadTransactions()
}

// Pagination handlers
function handlePageChange(page: number) {
  currentPage.value = page
  loadTransactions()
}

function handleItemsPerPageChange(newValue: number | string) {
  itemsPerPage.value = typeof newValue === 'string' ? parseInt(newValue, 10) : newValue
  currentPage.value = 1
  loadTransactions()
}

// Generate synthetic events handlers
function openGenerateModal() {
  generateSchemaId.value = null
  generateCount.value = 10
  generateResult.value = null
  showGenerateModal.value = true
}

async function handleGenerateEvents() {
  if (!generateSchemaId.value || !generateCount.value) return

  try {
    generating.value = true
    generateResult.value = null

    const payload = {
      count: generateCount.value,
    }

    const response = await api.post<{ generated_count: number; message: string }>(
      `/api/v1/schemas/${generateSchemaId.value}/generate-events`,
      payload
    )

    generateResult.value = {
      success: true,
      message: response?.message || `Successfully generated ${response?.generated_count || generateCount.value} events`,
    }

    // Reload transactions after a short delay to see the new events
    setTimeout(() => {
      loadTransactions()
    }, 1000)
  } catch (e: any) {
    generateResult.value = {
      success: false,
      message: e.message || 'Failed to generate events',
    }
  } finally {
    generating.value = false
  }
}

function toggleLiveUpdates() {
  if (isLive.value) {
    stopLiveUpdates()
  } else {
    startLiveUpdates()
  }
}

function startLiveUpdates() {
  isLive.value = true
  
  // Try SSE first, fall back to polling
  try {
    const token = localStorage.getItem('token')
    eventSource = new EventSource(`/api/v1/transactions/stream?token=${token}`)
    
    eventSource.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)
        if (data.type === 'update' || data.type === 'new') {
          loadTransactions()
        }
      } catch (e) {
        console.error('Error parsing SSE message:', e)
      }
    }

    eventSource.onerror = () => {
      // Fall back to polling if SSE fails
      if (eventSource) {
        eventSource.close()
        eventSource = null
      }
      startPolling()
    }
  } catch (e) {
    // Fall back to polling
    startPolling()
  }
}

function startPolling() {
  if (pollingInterval) return
  
  pollingInterval = setInterval(() => {
    if (isLive.value) {
      loadTransactions()
    }
  }, 5000) // Poll every 5 seconds
}

function stopLiveUpdates() {
  isLive.value = false
  
  if (eventSource) {
    eventSource.close()
    eventSource = null
  }
  
  if (pollingInterval) {
    clearInterval(pollingInterval)
    pollingInterval = null
  }
}

async function approveTransaction(id: string) {
  try {
    await api.patch(`/api/v1/transactions/${id}/approve`)
    await loadTransactions()
  } catch (e: any) {
    error.value = e.message || 'Failed to approve transaction'
  }
}

async function rejectTransaction(id: string) {
  try {
    await api.patch(`/api/v1/transactions/${id}/reject`)
    await loadTransactions()
  } catch (e: any) {
    error.value = e.message || 'Failed to reject transaction'
  }
}

function getStatusVariant(status: string): 'success' | 'danger' | 'warning' | 'info' | 'default' {
  switch (status) {
    case 'approved': return 'success'
    case 'rejected': return 'danger'
    case 'in_review': return 'warning'
    case 'pending': return 'info'
    default: return 'default'
  }
}

function formatDate(dateStr: string): string {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleString()
}
</script>

<style scoped>
.pulse {
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.3; }
}
</style>
