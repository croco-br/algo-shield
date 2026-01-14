<template>
  <v-container fluid class="pa-6">
    <!-- Header Section -->
    <div class="mb-6">
      <div class="d-flex align-center justify-space-between flex-wrap gap-3 mb-2">
        <div class="d-flex align-center gap-3">
          <v-icon icon="fa-exchange-alt" size="x-large" color="primary" />
          <div>
            <div class="d-flex align-center gap-2 mb-1">
              <h1 class="text-h4 font-weight-bold ma-0">{{ $t('views.transactions.title') }}</h1>
              <v-chip v-if="isLive" color="success" size="small">
                <v-icon icon="fa-circle" size="x-small" class="mr-1 pulse" />
                {{ $t('views.transactions.live') }}
              </v-chip>
            </div>
            <p class="text-body-2 text-grey-darken-1 ma-0">{{ $t('views.transactions.subtitle') }}</p>
          </div>
        </div>

        <!-- Action Buttons -->
        <div class="d-flex gap-2">
          <BaseButton
            v-if="syntheticMode && filters.schemaId"
            variant="secondary"
            size="sm"
            @click="openGenerateModal"
            prepend-icon="fa-bolt"
          >
            {{ $t('views.transactions.generateSynthetic') }}
          </BaseButton>
          <BaseButton
            :variant="isLive ? 'secondary' : 'primary'"
            size="sm"
            @click="toggleLiveUpdates"
            :prepend-icon="isLive ? 'fa-pause' : 'fa-play'"
            :disabled="!filters.schemaId"
          >
            {{ isLive ? $t('views.transactions.pause') : $t('views.transactions.live') }}
          </BaseButton>
          <BaseButton
            size="sm"
            variant="ghost"
            @click="loadTransactions"
            prepend-icon="fa-refresh"
            :disabled="!filters.schemaId"
          >
            {{ $t('views.transactions.refresh') }}
          </BaseButton>
        </div>
      </div>
    </div>

    <!-- No Schema Selected State -->
    <v-card v-if="!filters.schemaId" class="pa-12 text-center" variant="outlined">
      <v-icon icon="fa-database" size="64" color="grey" class="mb-4" />
      <h3 class="text-h5 font-weight-bold mb-3">{{ $t('views.transactions.noSchemaSelectedTitle') }}</h3>
      <p class="text-body-1 text-grey-darken-1 mb-0" style="max-width: 600px; margin: 0 auto;">
        {{ $t('views.transactions.noSchemaSelectedMessage') }}
      </p>
    </v-card>

    <!-- Main Content (Only shown when schema is selected) -->
    <template v-else>
      <!-- Schema Info Bar -->
      <v-card class="mb-4 pa-4" variant="flat" color="grey-lighten-4">
        <div class="d-flex align-center gap-3">
          <v-icon icon="fa-database" color="primary" />
          <div class="flex-grow-1">
            <div class="text-caption text-grey-darken-1">{{ $t('views.transactions.filterSchema') }}</div>
            <div class="text-subtitle-1 font-weight-medium">{{ selectedSchema?.name || $t('common.notAvailable') }}</div>
          </div>
          <div v-if="activeFiltersCount > 0" class="d-flex align-center gap-2">
            <v-icon icon="fa-filter" size="small" color="primary" />
            <span class="text-body-2 font-weight-medium">
              {{ $t('views.transactions.filtersActive', { count: activeFiltersCount }) }}
            </span>
          </div>
        </div>
      </v-card>

      <!-- Filters Section (Collapsible) -->
      <v-expansion-panels class="mb-4" variant="accordion">
        <v-expansion-panel>
          <v-expansion-panel-title>
            <div class="d-flex align-center gap-2">
              <v-icon icon="fa-filter" size="small" />
              <span class="font-weight-medium">{{ $t('views.transactions.filters') }}</span>
              <v-chip
                v-if="activeFiltersCount > 0"
                size="x-small"
                color="primary"
                class="ml-2"
              >
                {{ activeFiltersCount }}
              </v-chip>
            </div>
          </v-expansion-panel-title>
          <v-expansion-panel-text>
            <v-row class="mt-2">
              <v-col cols="12" md="4">
                <v-select
                  v-model="filters.status"
                  :items="statusOptions"
                  :label="$t('views.transactions.filterStatus')"
                  clearable
                  density="compact"
                  variant="outlined"
                  prepend-inner-icon="fa-circle-check"
                />
              </v-col>
              <v-col cols="12" md="4">
                <v-text-field
                  v-model="filters.startDate"
                  :label="$t('views.transactions.filterStartDate')"
                  type="date"
                  density="compact"
                  variant="outlined"
                  clearable
                  prepend-inner-icon="fa-calendar"
                />
              </v-col>
              <v-col cols="12" md="4">
                <v-text-field
                  v-model="filters.endDate"
                  :label="$t('views.transactions.filterEndDate')"
                  type="date"
                  density="compact"
                  variant="outlined"
                  clearable
                  prepend-inner-icon="fa-calendar"
                />
              </v-col>
            </v-row>
            <v-row>
              <v-col cols="12" class="d-flex gap-2 justify-end pt-0">
                <BaseButton
                  size="sm"
                  @click="applyFilters"
                  prepend-icon="fa-filter"
                >
                  {{ $t('views.transactions.applyFilters') }}
                </BaseButton>
                <BaseButton
                  size="sm"
                  variant="ghost"
                  @click="clearFilters"
                  prepend-icon="fa-times"
                >
                  {{ $t('views.transactions.clearFilters') }}
                </BaseButton>
              </v-col>
            </v-row>
          </v-expansion-panel-text>
        </v-expansion-panel>
      </v-expansion-panels>

      <!-- Loading State -->
      <LoadingSpinner
        v-if="loading"
        :text="$t('views.transactions.loading')"
        :centered="false"
      />

      <!-- Error State -->
      <ErrorMessage
        v-else-if="error"
        :title="$t('views.transactions.errorTitle')"
        :message="error"
        retryable
        @retry="loadTransactions"
      />

      <!-- Transaction Table -->
      <v-card v-else variant="outlined" class="overflow-hidden">
        <!-- Debug Info -->
        <div class="pa-2 bg-grey-lighten-4 text-caption">
          Debug: total={{ total }}, page={{ currentPage }}, itemsPerPage={{ itemsPerPage }},
          transactions.length={{ transactions.length }}
        </div>
        <v-data-table
          :headers="dynamicTableHeaders"
          :items="transactions"
          v-model:items-per-page="itemsPerPage"
          v-model:page="currentPage"
          :items-length="total"
          :items-per-page-options="itemsPerPageOptions"
          class="elevation-0"
          hover
          :loading="loading"
          show-current-page
        >
          <!-- Status Column -->
          <template #item.status="{ item }">
            <BaseBadge :variant="getStatusVariant(item.status)">
              {{ item.status }}
            </BaseBadge>
          </template>

          <!-- Created At Column -->
          <template #item.created_at="{ item }">
            <span class="text-caption">{{ formatDate(item.created_at) }}</span>
          </template>

          <!-- Actions Column -->
          <template #item.actions="{ item }">
            <div class="d-flex gap-1">
              <v-tooltip location="top">
                <template #activator="{ props }">
                  <v-btn
                    v-bind="props"
                    size="small"
                    variant="text"
                    icon="fa-eye"
                    @click="openDetailModal(item)"
                  />
                </template>
                <span>{{ $t('views.transactions.viewDetails') }}</span>
              </v-tooltip>

              <v-tooltip v-if="item.status === 'pending'" location="top">
                <template #activator="{ props }">
                  <v-btn
                    v-bind="props"
                    size="small"
                    variant="text"
                    icon="fa-check"
                    color="success"
                    @click="approveTransaction(item.id)"
                  />
                </template>
                <span>{{ $t('views.transactions.approve') }}</span>
              </v-tooltip>

              <v-tooltip v-if="item.status === 'pending'" location="top">
                <template #activator="{ props }">
                  <v-btn
                    v-bind="props"
                    size="small"
                    variant="text"
                    icon="fa-times"
                    color="error"
                    @click="rejectTransaction(item.id)"
                  />
                </template>
                <span>{{ $t('views.transactions.reject') }}</span>
              </v-tooltip>
            </div>
          </template>

          <!-- No Data Slot -->
          <template #no-data>
            <div class="text-center pa-8">
              <v-icon icon="fa-inbox" size="48" color="grey-lighten-1" class="mb-3" />
              <p class="text-body-1 text-grey-darken-1 mb-0">{{ $t('common.noData') }}</p>
            </div>
          </template>

          <!-- Custom Footer -->
          <template #bottom>
            <div class="v-data-table-footer pa-4 d-flex align-center justify-space-between">
              <div class="text-body-2">
                <template v-if="total > 0">
                  Showing {{ ((currentPage - 1) * itemsPerPage) + 1 }} to {{ Math.min(currentPage * itemsPerPage, total) }} of {{ total }} entries
                </template>
                <template v-else>
                  No entries
                </template>
              </div>
              <div class="d-flex align-center gap-4">
                <div class="d-flex align-center gap-2">
                  <span class="text-body-2">Items per page:</span>
                  <v-select
                    v-model="itemsPerPage"
                    :items="itemsPerPageOptions"
                    density="compact"
                    variant="outlined"
                    hide-details
                    style="width: 100px"
                  />
                </div>
                <v-pagination
                  v-if="total > 0"
                  v-model="currentPage"
                  :length="Math.max(1, Math.ceil(total / itemsPerPage))"
                  :total-visible="7"
                  size="small"
                />
              </div>
            </div>
          </template>
        </v-data-table>
      </v-card>
    </template>

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
          prepend-inner-icon="fa-database"
        />

        <BaseInput
          v-model.number="generateCount"
          :label="$t('views.transactions.eventCount')"
          type="number"
          min="1"
          max="10000"
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
        <BaseButton
          variant="ghost"
          @click="showGenerateModal = false"
          prepend-icon="fa-xmark"
        >
          {{ $t('components.modal.cancel') }}
        </BaseButton>
        <BaseButton
          @click="handleGenerateEvents"
          :loading="generating"
          :disabled="!generateSchemaId || !generateCount || generateCount < 1 || generateCount > 10000"
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
import { useSystemModeStore } from '@/stores/systemMode'
import { useAuthStore } from '@/stores/auth'
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
const authStore = useAuthStore()

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
const itemsPerPageOptions = [
  { value: 10, title: '10' },
  { value: 25, title: '25' },
  { value: 50, title: '50' },
  { value: 100, title: '100' },
  { value: 200, title: '200' }
]

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
  ).slice(0, 10)
})

const activeFiltersCount = computed(() => {
  let count = 0
  if (filters.status) count++
  if (filters.startDate) count++
  if (filters.endDate) count++
  return count
})

// Dynamic table headers based on selected schema
const dynamicTableHeaders = computed(() => {
  const headers: Array<{ title: string; key: string; sortable: boolean; value?: (item: Transaction) => any }> = []

  // If schema selected, show schema fields first
  if (selectedSchema.value && schemaFields.value.length > 0) {
    schemaFields.value.forEach((field: ExtractedField) => {
      const safeKey = field.path.replace(/\./g, '_')
      const fieldName = field.path.replace(/\./g, '_').toUpperCase()
      headers.push({
        title: fieldName,
        key: `schema_field_${safeKey}`,
        sortable: false,
        value: (item: Transaction) => {
          const val = getSchemaFieldValue(item, field.path)
          return val !== undefined ? formatSchemaFieldValue(val, field.type) : '-'
        }
      })
    })
  }

  // Add transaction columns after schema fields
  headers.push(
    { title: t('views.transactions.tableStatus').toUpperCase(), key: 'status', sortable: true },
    { title: t('views.transactions.tableCreated').toUpperCase(), key: 'created_at', sortable: true }
  )

  // Actions column always last
  headers.push({ title: t('views.transactions.tableActions').toUpperCase(), key: 'actions', sortable: false })
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
  } else {
    loading.value = false
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
    loading.value = false
  }
})

// Watch for status filter changes to auto-reload
watch(() => filters.status, () => {
  if (filters.schemaId) {
    currentPage.value = 1
    loadTransactions()
  }
})

// Watch for synthetic mode changes
watch(() => systemModeStore.syntheticMode, () => {
  if (!systemModeStore.syntheticMode && showGenerateModal.value) {
    showGenerateModal.value = false
  }
})

// Debug watcher for total
watch(total, (newVal, oldVal) => {
  console.log('[WATCH] total changed from', oldVal, 'to', newVal)
})

// Watch for page changes
watch(currentPage, (newPage, oldPage) => {
  if (filters.schemaId && newPage !== oldPage) {
    console.log('[Watch] currentPage changed from', oldPage, 'to', newPage)
    loadTransactions()
  }
})

// Watch for items per page changes
watch(itemsPerPage, (newValue, oldValue) => {
  if (filters.schemaId && newValue !== oldValue) {
    console.log('[Watch] itemsPerPage changed from', oldValue, 'to', newValue)
    // Reset to first page but don't trigger currentPage watcher if already on page 1
    const wasOnFirstPage = currentPage.value === 1
    currentPage.value = 1

    // Only call loadTransactions if we were already on page 1
    // (otherwise the currentPage watcher will handle it)
    if (wasOnFirstPage) {
      loadTransactions()
    }
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
  if (!filters.schemaId) return

  // Ensure currentPage is valid (at least 1)
  if (currentPage.value < 1) {
    currentPage.value = 1
  }

  loading.value = true
  error.value = ''
  try {
    const params = new URLSearchParams()
    const offset = (currentPage.value - 1) * itemsPerPage.value
    params.append('limit', String(itemsPerPage.value))
    params.append('offset', String(offset))

    if (filters.status) params.append('status', filters.status)
    if (filters.schemaId) params.append('schema_id', filters.schemaId)

    console.log('[Load Transactions] Request params:', {
      schemaId: filters.schemaId,
      status: filters.status,
      syntheticMode: syntheticMode.value,
      limit: itemsPerPage.value,
      offset: offset
    })
    if (filters.startDate && filters.startDate.trim() !== '') {
      try {
        const startDate = new Date(filters.startDate + 'T00:00:00Z')
        if (!isNaN(startDate.getTime())) {
          params.append('start_date', startDate.toISOString())
        }
      } catch (e) {
        console.error('Error parsing start date:', e)
      }
    }
    if (filters.endDate && filters.endDate.trim() !== '') {
      try {
        const endDate = new Date(filters.endDate + 'T23:59:59.999Z')
        if (!isNaN(endDate.getTime())) {
          params.append('end_date', endDate.toISOString())
        }
      } catch (e) {
        console.error('Error parsing end date:', e)
      }
    }

    const response = await api.get<{ transactions: Transaction[]; total?: number }>(`/api/v1/transactions?${params.toString()}`)
    console.log('[Load Transactions] Raw Response:', response)
    console.log('[Load Transactions] Response.total type:', typeof response.total, 'value:', response.total)

    transactions.value = response.transactions || []
    total.value = response.total ?? 0

    console.log('[Load Transactions] After assignment - total.value:', total.value)
    console.log('[Load Transactions] transactions.value.length:', transactions.value.length)
    console.log('[Load Transactions] First transaction:', transactions.value[0])
    console.log('[Load Transactions] Response:', {
      count: transactions.value.length,
      total: total.value,
      transactions: transactions.value
    })
    console.log('[Load Transactions] dynamicTableHeaders:', dynamicTableHeaders.value)

    // If currentPage is beyond available pages, adjust it
    const maxPage = Math.ceil(total.value / itemsPerPage.value)
    if (maxPage > 0 && currentPage.value > maxPage) {
      console.log('[Load Transactions] currentPage', currentPage.value, 'exceeds maxPage', maxPage, '- adjusting')
      currentPage.value = maxPage
      // Don't call loadTransactions here to avoid infinite loop - the watcher will handle it
    }
  } catch (e: any) {
    error.value = e.message || t('views.transactions.errorLoad')
    console.error('[Load Transactions] Error:', e)
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

// Generate synthetic events handlers
function openGenerateModal() {
  generateSchemaId.value = filters.schemaId
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

    console.log('[Generate Events] Sending request:', {
      schemaId: generateSchemaId.value,
      count: generateCount.value,
      syntheticMode: syntheticMode.value,
      hasToken: !!authStore.getToken(),
      hasCsrfToken: !!authStore.getCsrfToken()
    })

    const response = await api.post<{ generated_count: number; failed_count: number; message: string }>(
      `/api/v1/schemas/${generateSchemaId.value}/generate-events`,
      payload
    )

    console.log('[Generate Events] Response:', response)

    // Use i18n messages based on success/failure counts
    let message: string
    if (response?.failed_count === 0) {
      message = t('views.transactions.generateSuccess', { count: response.generated_count })
    } else {
      message = t('views.transactions.generatePartial', {
        total: (response?.generated_count || 0) + (response?.failed_count || 0),
        success: response?.generated_count || 0,
        failed: response?.failed_count || 0
      })
    }

    generateResult.value = {
      success: response?.failed_count === 0,
      message: message,
    }

    // Worker should process in <50ms per project SLA
    // Wait 100ms to account for queue + processing + DB write
    await new Promise(resolve => setTimeout(resolve, 100))
    console.log('[Generate Events] Reloading after 100ms...')
    await loadTransactions()

    if (transactions.value.length === 0) {
      console.error('[Generate Events] CRITICAL: No transactions found')
      console.error('Expected processing time: <50ms per transaction')
      console.error('Possible causes:')
      console.error('  1. Worker container NOT RUNNING - check: docker ps | grep worker')
      console.error('  2. Worker not connected to Redis queue')
      console.error('  3. Synthetic mode not active (current:', syntheticMode.value, ')')
      console.error('  4. Worker crashed or has errors - check: docker logs algo-shield-worker')
      console.error('  5. Schema ID mismatch in worker processing')
      console.error('')
      console.error('Enable Live Updates to see transactions as they arrive')

      generateResult.value = {
        success: false,
        message: 'Events queued but not processed. Worker may not be running. Check console for details.',
      }
    } else {
      console.log(`[Generate Events] Success: ${transactions.value.length} transactions loaded`)
    }
  } catch (e: any) {
    console.error('[Generate Events] Error:', e)

    // Check if it's a CSRF error (403 Forbidden with CSRF in message)
    const errorMessage = e.message || 'Failed to generate events'
    const isCSRFError = errorMessage.includes('CSRF') || errorMessage.includes('403') || errorMessage.includes('Forbidden')

    if (isCSRFError) {
      console.error('[Generate Events] CSRF token error detected')
      console.error('This typically happens when:')
      console.error('  1. User logged in before CSRF system was implemented')
      console.error('  2. CSRF token expired (24h TTL)')
      console.error('  3. Redis connection lost')
      console.error('')
      console.error('Solution: Please logout and login again to refresh CSRF token')

      generateResult.value = {
        success: false,
        message: 'CSRF token error. Please logout and login again.',
      }
    } else {
      generateResult.value = {
        success: false,
        message: errorMessage,
      }
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
    error.value = ''
    await api.patch(`/api/v1/transactions/${id}/approve`)
    // Update transaction status locally instead of reloading entire list
    const transaction = transactions.value.find(t => t.id === id)
    if (transaction) {
      transaction.status = 'approved'
    }
  } catch (e: any) {
    error.value = e.message || 'Failed to approve transaction'
    console.error('Error approving transaction:', e)
  }
}

async function rejectTransaction(id: string) {
  try {
    error.value = ''
    await api.patch(`/api/v1/transactions/${id}/reject`)
    // Update the transaction status locally instead of reloading
    const transaction = transactions.value.find(t => t.id === id)
    if (transaction) {
      transaction.status = 'rejected'
    }
  } catch (e: any) {
    error.value = e.message || 'Failed to reject transaction'
    console.error('Error rejecting transaction:', e)
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

/* Better spacing for action buttons */
.v-data-table :deep(.v-data-table__td) {
  padding: 8px 16px !important;
}

/* Improve expansion panel appearance */
.v-expansion-panels {
  border: 1px solid rgba(0, 0, 0, 0.12);
  border-radius: 4px;
}
</style>
