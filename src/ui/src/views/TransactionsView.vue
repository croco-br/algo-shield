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

        <!-- Transaction Table - Dynamic based on schema -->
        <v-card v-if="filters.schemaId" variant="outlined">
          <v-data-table
            :headers="dynamicTableHeaders"
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
            <!-- Core columns -->
            <template #item.status="{ item }">
              <BaseBadge :variant="getStatusVariant(item.status)">
                {{ item.status }}
              </BaseBadge>
            </template>
            <template #item.created_at="{ item }">
              {{ formatDate(item.created_at) }}
            </template>


            <!-- Actions -->
            <template #item.actions="{ item }">
              <div class="d-flex gap-2">
                <BaseButton 
                  size="sm"
                  variant="ghost"
                  @click="openDetailModal(item)"
                  prepend-icon="fa-eye"
                >
                  {{ $t('views.transactions.viewDetails') }}
                </BaseButton>
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

    <!-- Transaction Detail Modal -->
    <TransactionDetailModal
      v-model="showDetailModal"
      :transaction="selectedTransaction"
      @close="closeDetailModal"
    />
  </v-container>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted, computed, watch } from 'vue'
import { useRoute } from 'vue-router'
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
import TransactionDetailModal from '@/components/TransactionDetailModal.vue'

interface Schema {
  id: string
  name: string
}

interface ExtractedField {
  path: string
  type: string
  nullable: boolean
  sample_value?: any
}

interface EventSchema {
  id: string
  name: string
  description?: string
  extracted_fields: ExtractedField[]
  sample_json: Record<string, any>
}

const route = useRoute()
const { t } = useLocale()
const systemModeStore = useSystemModeStore()

const loading = ref(true)
const error = ref('')
const transactions = ref<Transaction[]>([])
const total = ref(0)
const schemas = ref<Schema[]>([])
const selectedSchemaData = ref<EventSchema | null>(null)
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

// Transaction detail modal state
const showDetailModal = ref(false)
const selectedTransaction = ref<Transaction | null>(null)

const filters = reactive({
  status: null as string | null,
  schemaId: null as string | null,
  startDate: null as string | null,
  endDate: null as string | null,
})

const statusOptions = computed(() => [
  { title: t('views.transactions.statusAll'), value: null },
  { title: t('views.transactions.statusApproved'), value: 'approved' },
  { title: t('views.transactions.statusRejected'), value: 'rejected' },
  { title: t('views.transactions.statusInReview'), value: 'in_review' },
  { title: t('views.transactions.statusPending'), value: 'pending' },
])

const generateSchemaOptions = computed(() => 
  schemas.value.map((s: Schema) => ({ title: s.name, value: s.id }))
)

const selectedSchema = computed(() => {
  if (!filters.schemaId) return null
  return schemas.value.find(s => s.id === filters.schemaId) || null
})

const schemaFields = computed<ExtractedField[]>(() => {
  if (!selectedSchemaData.value) return []
  return selectedSchemaData.value.extracted_fields.filter(f => 
    f.type !== 'object' && f.type !== 'array'
  ).slice(0, 10) // Limit to 10 most relevant fields
})

// Dynamic table headers based on selected schema
const dynamicTableHeaders = computed(() => {
  const headers: Array<{ title: string; key: string; sortable: boolean; value?: (item: Transaction) => any }> = [
    { title: t('views.transactions.tableStatus'), key: 'status', sortable: true },
    { title: t('views.transactions.tableCreated'), key: 'created_at', sortable: true },
  ]

  // If schema selected, show schema fields
  if (selectedSchema.value && schemaFields.value.length > 0) {
    // Add schema fields as columns with custom value function
    schemaFields.value.forEach((field: ExtractedField) => {
      const safeKey = field.path.replace(/\./g, '_')
      const fieldName = field.path.split('.').pop() || field.path
      headers.push({
        title: fieldName.toUpperCase(),
        key: `schema_field_${safeKey}`,
        sortable: false,
        value: (item: Transaction) => {
          const val = getSchemaFieldValue(item, field.path)
          return val !== undefined ? formatSchemaFieldValue(val, field.type) : '-'
        }
      })
    })
  }

  headers.push({ title: t('views.transactions.tableActions'), key: 'actions', sortable: false })
  return headers
})

onMounted(async () => {
  await systemModeStore.loadMode()
  await loadSchemas()
  
  // Read schemaId from query parameter
  const schemaIdFromQuery = route.query.schemaId as string | undefined
  if (schemaIdFromQuery) {
    filters.schemaId = schemaIdFromQuery
    await loadSchemaData(schemaIdFromQuery)
    await loadTransactions()
  }
})

// Watch for query parameter changes
watch(() => route.query.schemaId, async (newSchemaId) => {
  const schemaId = typeof newSchemaId === 'string' ? newSchemaId : undefined
  if (schemaId && schemaId !== filters.schemaId) {
    filters.schemaId = schemaId
    await loadSchemaData(schemaId)
    await loadTransactions()
  } else if (!schemaId) {
    filters.schemaId = null
    selectedSchemaData.value = null
    transactions.value = []
    total.value = 0
  }
})

// Watch for synthetic mode changes
watch(() => systemModeStore.syntheticMode, () => {
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


async function loadSchemaData(schemaId: string) {
  try {
    const response = await api.get<EventSchema>(`/api/v1/schemas/${schemaId}`)
    selectedSchemaData.value = response
  } catch (e) {
    console.error('Error loading schema data:', e)
    selectedSchemaData.value = null
  }
}

function getSchemaFieldValue(transaction: Transaction, fieldPath: string): any {
  const metadata = transaction.metadata || {}
  const parts = fieldPath.split('.')
  let current: any = metadata
  
  for (const part of parts) {
    if (current === null || current === undefined) return undefined
    current = current[part]
  }
  
  return current
}

function formatSchemaFieldValue(value: any, type: string): string {
  if (value === null || value === undefined) return '-'
  
  switch (type) {
    case 'number':
      return typeof value === 'number' ? value.toLocaleString() : String(value)
    case 'boolean':
      return value ? 'Yes' : 'No'
    case 'datetime':
      if (typeof value === 'string') {
        try {
          const date = new Date(value)
          if (!isNaN(date.getTime())) {
            return date.toLocaleString()
          }
        } catch {
          // Fall through to return string value
        }
      }
      return String(value)
    case 'string':
      return String(value)
    default:
      return String(value)
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
  filters.startDate = null
  filters.endDate = null
  currentPage.value = 1
  if (filters.schemaId) {
    loadTransactions()
  }
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
      if (eventSource) {
        eventSource.close()
        eventSource = null
      }
      startPolling()
    }
  } catch (e) {
    startPolling()
  }
}

function startPolling() {
  if (pollingInterval) return
  
  pollingInterval = setInterval(() => {
    if (isLive.value) {
      loadTransactions()
    }
  }, 5000)
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

function openDetailModal(transaction: Transaction) {
  selectedTransaction.value = transaction
  showDetailModal.value = true
}

function closeDetailModal() {
  showDetailModal.value = false
  selectedTransaction.value = null
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
