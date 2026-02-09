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
        <v-data-table-server
          :headers="dynamicTableHeaders"
          :items="transactions"
          :items-length="total"
          :items-per-page="itemsPerPage"
          :page="currentPage"
          class="elevation-0"
          hover
          :loading="loading"
          @update:page="onPageChange"
          @update:items-per-page="onItemsPerPageChange"
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
                    :aria-label="$t('views.transactions.viewDetails')"
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
                    :aria-label="$t('views.transactions.approve')"
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
                    :aria-label="$t('views.transactions.reject')"
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
                  {{ $t('views.transactions.showing') }} {{ ((currentPage - 1) * itemsPerPage) + 1 }} {{ $t('views.transactions.to') }} {{ Math.min(currentPage * itemsPerPage, total) }} {{ $t('views.transactions.of') }} {{ total }} {{ $t('views.transactions.transactions') }}
                </template>
                <template v-else>
                  {{ $t('common.noData') }}
                </template>
              </div>
              <div class="d-flex align-center gap-4">
                <div class="d-flex align-center gap-2">
                  <span class="text-body-2">{{ $t('components.table.rowsPerPage') }}:</span>
                  <v-select
                    :model-value="itemsPerPage"
                    :items="itemsPerPageOptions"
                    density="compact"
                    variant="outlined"
                    hide-details
                    style="width: 100px"
                    @update:model-value="onItemsPerPageChange"
                  />
                </div>
                <v-pagination
                  v-if="total > 0"
                  :model-value="currentPage"
                  :length="Math.max(1, Math.ceil(total / itemsPerPage))"
                  :total-visible="7"
                  size="small"
                  @update:model-value="onPageChange"
                />
              </div>
            </div>
          </template>
        </v-data-table-server>
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
          :disabled="generating"
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
          :disabled="generating"
          :placeholder="$t('views.transactions.eventCountPlaceholder')"
          prepend-inner-icon="fa-hashtag"
          class="mb-4"
        />

        <!-- Progress bar -->
        <div v-if="generating && generateTotal > 0" class="mt-4">
          <div class="d-flex justify-space-between mb-1">
            <span class="text-body-2 font-weight-medium">
              {{ $t('views.transactions.generating') }}
            </span>
            <span class="text-body-2 text-grey-darken-1">
              {{ generateProgress }} / {{ generateTotal }}
            </span>
          </div>
          <v-progress-linear
            :model-value="(generateProgress / generateTotal) * 100"
            color="primary"
            height="8"
            rounded
          />
        </div>

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
          :disabled="generating"
          prepend-icon="fa-xmark"
        >
          {{ $t('components.modal.cancel') }}
        </BaseButton>
        <BaseButton
          @click="handleGenerateEvents"
          :loading="generating"
          :disabled="generating || !generateSchemaId || !generateCount || generateCount < 1 || generateCount > 10000"
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
import { useI18n } from 'vue-i18n'
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
const { t } = useI18n()
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
const generateProgress = ref(0)
const generateTotal = ref(0)
const generateBatchSize = 100

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
  const headers: Array<{ title: string; key: string; sortable?: boolean; align?: 'start' | 'center' | 'end' }> = []

  // If schema selected, show schema fields first
  if (selectedSchema.value && schemaFields.value.length > 0) {
    schemaFields.value.forEach((field: ExtractedField) => {
      const safeKey = field.path.replace(/\./g, '_')
      const fieldName = field.path.replace(/\./g, '_').toUpperCase()
      headers.push({
        title: fieldName,
        key: `schema_field_${safeKey}`,
        sortable: false,
        align: 'start' as const
      })
    })
  }

  // Add transaction columns after schema fields
  headers.push(
    { title: t('views.transactions.tableStatus').toUpperCase(), key: 'status', sortable: true, align: 'start' as const },
    { title: t('views.transactions.tableCreated').toUpperCase(), key: 'created_at', sortable: true, align: 'start' as const }
  )

  // Actions column always last
  headers.push({ title: t('views.transactions.tableActions').toUpperCase(), key: 'actions', sortable: false, align: 'center' as const })
  return headers
})

// Pagination event handlers
function onPageChange(page: number) {
  if (page !== currentPage.value) {
    currentPage.value = page
    loadTransactions()
  }
}

function onItemsPerPageChange(perPage: number) {
  if (perPage !== itemsPerPage.value) {
    itemsPerPage.value = perPage
    currentPage.value = 1
    loadTransactions()
  }
}

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

// Watch for synthetic mode changes — reload data and close generate modal
watch(() => systemModeStore.syntheticMode, () => {
  if (!systemModeStore.syntheticMode && showGenerateModal.value) {
    showGenerateModal.value = false
  }
  loadTransactions()
})

onUnmounted(() => {
  stopLiveUpdates()
})

async function loadSchemas() {
  try {
    const response = await api.get<{ schemas: Schema[] }>('/api/v1/schemas')
    schemas.value = response?.schemas || []
  } catch (e) {
    // Schema load failure is non-critical
  }
}

async function loadSchemaData(schemaId: string) {
  try {
    const response = await api.get<EventSchema>(`/api/v1/schemas/${schemaId}`)
    selectedSchemaData.value = response
  } catch (e) {
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

    if (filters.startDate && filters.startDate.trim() !== '') {
      try {
        const startDate = new Date(filters.startDate + 'T00:00:00Z')
        if (!isNaN(startDate.getTime())) {
          params.append('start_date', startDate.toISOString())
        }
      } catch {
        // Invalid date - skip
      }
    }
    if (filters.endDate && filters.endDate.trim() !== '') {
      try {
        const endDate = new Date(filters.endDate + 'T23:59:59.999Z')
        if (!isNaN(endDate.getTime())) {
          params.append('end_date', endDate.toISOString())
        }
      } catch {
        // Invalid date - skip
      }
    }

    const response = await api.get<{ transactions: Transaction[]; total?: number }>(`/api/v1/transactions?${params.toString()}`)

    const rawTransactions = response.transactions || []
    total.value = response.total ?? 0

    // Process transactions to add schema field values
    transactions.value = rawTransactions.map(transaction => {
      const processed: Record<string, any> = { ...transaction }

      if (schemaFields.value.length > 0) {
        schemaFields.value.forEach((field: ExtractedField) => {
          const safeKey = field.path.replace(/\./g, '_')
          const fieldKey = `schema_field_${safeKey}`
          const val = getSchemaFieldValue(transaction, field.path)
          processed[fieldKey] = val !== undefined ? formatSchemaFieldValue(val, field.type) : '-'
        })
      }

      return processed as Transaction
    })

    // If currentPage is beyond available pages, adjust it
    const maxPage = Math.ceil(total.value / itemsPerPage.value)
    if (maxPage > 0 && currentPage.value > maxPage) {
      currentPage.value = maxPage
    }
  } catch (e: any) {
    error.value = e.message || t('views.transactions.errorLoad')
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
  generateProgress.value = 0
  generateTotal.value = 0
  showGenerateModal.value = true
}

async function handleGenerateEvents() {
  if (!generateSchemaId.value || !generateCount.value) return

  try {
    generating.value = true
    generateResult.value = null
    generateProgress.value = 0
    generateTotal.value = generateCount.value

    const total = generateCount.value
    let totalGenerated = 0
    let totalFailed = 0

    // Split into batches for progress tracking
    let remaining = total
    while (remaining > 0) {
      const batchCount = Math.min(remaining, generateBatchSize)

      const response = await api.post<{ generated_count: number; failed_count: number; message: string }>(
        `/api/v1/schemas/${generateSchemaId.value}/generate-events`,
        { count: batchCount }
      )

      totalGenerated += response?.generated_count || 0
      totalFailed += response?.failed_count || 0
      remaining -= batchCount
      generateProgress.value = totalGenerated + totalFailed
    }

    let message: string
    if (totalFailed === 0) {
      message = t('views.transactions.generateSuccess', { count: totalGenerated })
    } else {
      message = t('views.transactions.generatePartial', {
        total: totalGenerated + totalFailed,
        success: totalGenerated,
        failed: totalFailed
      })
    }

    generateResult.value = {
      success: totalFailed === 0,
      message: message,
    }

    await new Promise(resolve => setTimeout(resolve, 100))
    await loadTransactions()

    if (transactions.value.length === 0) {
      generateResult.value = {
        success: false,
        message: t('views.transactions.generateFailed'),
      }
    }
  } catch (e: any) {
    const errorMessage = e.message || t('views.transactions.generateFailed')
    const isCSRFError = errorMessage.includes('CSRF') || errorMessage.includes('403') || errorMessage.includes('Forbidden')

    generateResult.value = {
      success: false,
      message: isCSRFError
        ? t('errors.FORBIDDEN')
        : errorMessage,
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
  startPolling()
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

  if (pollingInterval) {
    clearInterval(pollingInterval)
    pollingInterval = null
  }
}

async function approveTransaction(id: string) {
  try {
    error.value = ''
    await api.patch(`/api/v1/transactions/${id}/approve`)
    const transaction = transactions.value.find(t => t.id === id)
    if (transaction) {
      transaction.status = 'approved'
    }
  } catch (e: any) {
    error.value = e.message || t('views.transactions.errorLoad')
  }
}

async function rejectTransaction(id: string) {
  try {
    error.value = ''
    await api.patch(`/api/v1/transactions/${id}/reject`)
    const transaction = transactions.value.find(t => t.id === id)
    if (transaction) {
      transaction.status = 'rejected'
    }
  } catch (e: any) {
    error.value = e.message || t('views.transactions.errorLoad')
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
