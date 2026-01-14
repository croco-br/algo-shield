<template>
  <BaseModal
    v-model="isOpen"
    :title="$t('components.transactionDetailModal.title')"
    size="lg"
    @close="handleClose"
  >
    <div v-if="loading" class="d-flex justify-center py-8">
      <LoadingSpinner :text="$t('components.transactionDetailModal.loading')" />
    </div>

    <div v-else-if="error" class="pa-4">
      <v-alert type="error" variant="tonal">
        {{ error }}
      </v-alert>
    </div>

    <div v-else-if="transaction" class="transaction-detail">
      <!-- Basic Information Section -->
      <div class="mb-6">
        <h3 class="text-h6 mb-4">{{ $t('components.transactionDetailModal.basicInfo') }}</h3>
        <v-row>
          <v-col cols="12" md="6">
            <div class="detail-item mb-3">
              <span class="detail-label">{{ $t('components.transactionDetailModal.status') }}:</span>
              <BaseBadge :variant="getStatusVariant(transaction.status)">
                {{ transaction.status }}
              </BaseBadge>
            </div>
          </v-col>
          <v-col cols="12" md="6">
            <div class="detail-item mb-3">
              <span class="detail-label">{{ $t('components.transactionDetailModal.createdAt') }}:</span>
              <span class="detail-value">{{ formatDate(transaction.created_at) }}</span>
            </div>
          </v-col>
          <v-col cols="12" md="6" v-if="transaction.processed_at">
            <div class="detail-item mb-3">
              <span class="detail-label">{{ $t('components.transactionDetailModal.processedAt') }}:</span>
              <span class="detail-value">{{ formatDate(transaction.processed_at) }}</span>
            </div>
          </v-col>
          <v-col cols="12" md="6" v-if="transaction.processing_time_ms">
            <div class="detail-item mb-3">
              <span class="detail-label">{{ $t('components.transactionDetailModal.processingTime') }}:</span>
              <span class="detail-value">{{ transaction.processing_time_ms }}ms</span>
            </div>
          </v-col>
        </v-row>
      </div>

      <!-- Schema Information -->
      <div v-if="schema" class="mb-6">
        <h3 class="text-h6 mb-4">
          {{ $t('components.transactionDetailModal.schemaInfo') }}
          <BaseBadge variant="info" size="sm" class="ml-2">{{ schema.name }}</BaseBadge>
        </h3>
        <div v-if="schema.description" class="mb-4">
          <p class="text-body-2 text-grey-darken-1">{{ schema.description }}</p>
        </div>

        <!-- Schema Fields -->
        <div v-if="schemaFields.length > 0">
          <h4 class="text-subtitle-1 mb-3">{{ $t('components.transactionDetailModal.schemaFields') }}</h4>
          <v-table density="compact" class="schema-fields-table">
            <thead>
              <tr>
                <th>{{ $t('components.transactionDetailModal.fieldPath') }}</th>
                <th>{{ $t('components.transactionDetailModal.fieldType') }}</th>
                <th>{{ $t('components.transactionDetailModal.fieldValue') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="field in schemaFields" :key="field.path">
                <td class="font-mono text-caption">{{ field.path }}</td>
                <td>
                  <BaseBadge variant="default" size="sm">{{ field.type }}</BaseBadge>
                </td>
                <td class="field-value">
                  <pre v-if="isComplexValue(field.value)">{{ formatValue(field.value) }}</pre>
                  <span v-else>{{ formatValue(field.value) }}</span>
                </td>
              </tr>
            </tbody>
          </v-table>
        </div>
      </div>

      <!-- Matched Rules -->
      <div v-if="transaction.matched_rules && transaction.matched_rules.length > 0" class="mb-6">
        <h3 class="text-h6 mb-4">{{ $t('components.transactionDetailModal.matchedRules') }}</h3>
        <div class="d-flex flex-wrap gap-2">
          <BaseBadge
            v-for="ruleId in transaction.matched_rules"
            :key="ruleId"
            variant="warning"
            size="sm"
          >
            {{ ruleId }}
          </BaseBadge>
        </div>
      </div>

      <!-- Additional Metadata -->
      <div v-if="additionalMetadata && Object.keys(additionalMetadata).length > 0" class="mb-6">
        <h3 class="text-h6 mb-4">{{ $t('components.transactionDetailModal.additionalMetadata') }}</h3>
        <v-expansion-panels variant="accordion">
          <v-expansion-panel>
            <v-expansion-panel-title>
              {{ $t('components.transactionDetailModal.viewMetadata') }}
            </v-expansion-panel-title>
            <v-expansion-panel-text>
              <pre class="metadata-json">{{ formatJSON(additionalMetadata) }}</pre>
            </v-expansion-panel-text>
          </v-expansion-panel>
        </v-expansion-panels>
      </div>

      <!-- Full Metadata (if no schema) -->
      <div v-else-if="transaction.metadata && Object.keys(transaction.metadata).length > 0" class="mb-6">
        <h3 class="text-h6 mb-4">{{ $t('components.transactionDetailModal.metadata') }}</h3>
        <v-expansion-panels variant="accordion">
          <v-expansion-panel>
            <v-expansion-panel-title>
              {{ $t('components.transactionDetailModal.viewMetadata') }}
            </v-expansion-panel-title>
            <v-expansion-panel-text>
              <pre class="metadata-json">{{ formatJSON(transaction.metadata) }}</pre>
            </v-expansion-panel-text>
          </v-expansion-panel>
        </v-expansion-panels>
      </div>
    </div>

    <template #footer>
      <BaseButton variant="ghost" @click="handleClose" prepend-icon="fa-xmark">
        {{ $t('components.modal.close') }}
      </BaseButton>
    </template>
  </BaseModal>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { api } from '@/lib/api'
import type { Transaction } from '@/types/transaction'
import BaseModal from '@/components/BaseModal.vue'
import BaseBadge from '@/components/BaseBadge.vue'
import BaseButton from '@/components/BaseButton.vue'
import LoadingSpinner from '@/components/LoadingSpinner.vue'
import { useCurrency } from '@/composables/useCurrency'

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

interface Props {
  modelValue: boolean
  transaction: Transaction | null
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  'close': []
}>()


const isOpen = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const loading = ref(false)
const error = ref('')
const schema = ref<EventSchema | null>(null)

watch(() => props.transaction, async (newTransaction) => {
  if (newTransaction?.schema_id) {
    await loadSchema(newTransaction.schema_id)
  } else {
    schema.value = null
  }
}, { immediate: true })

async function loadSchema(schemaId: string) {
  loading.value = true
  error.value = ''
  try {
    const response = await api.get<EventSchema>(`/api/v1/schemas/${schemaId}`)
    schema.value = response
  } catch (e: any) {
    error.value = e.message || 'Failed to load schema'
    schema.value = null
  } finally {
    loading.value = false
  }
}

const schemaFields = computed(() => {
  if (!schema.value || !props.transaction) return []
  
  const fields: Array<{ path: string; type: string; value: any }> = []
  const metadata = props.transaction.metadata || {}
  
  for (const field of schema.value.extracted_fields) {
    const value = getNestedValue(metadata, field.path)
    if (value !== undefined) {
      fields.push({
        path: field.path,
        type: field.type,
        value: value
      })
    }
  }
  
  return fields
})

const additionalMetadata = computed(() => {
  if (!schema.value || !props.transaction) return props.transaction?.metadata || {}
  
  const schemaFieldPaths = new Set(schema.value.extracted_fields.map(f => f.path))
  const metadata = props.transaction.metadata || {}
  const additional: Record<string, any> = {}
  
  // Helper function to check if a key is part of any schema field path
  const isSchemaField = (key: string): boolean => {
    // Check exact match
    if (schemaFieldPaths.has(key)) {
      return true
    }
    // Check if this key is a prefix of any schema field path (for nested fields)
    // e.g., "user" matches "user.created_at" or "user.id"
    for (const path of schemaFieldPaths) {
      if (path.startsWith(key + '.')) {
        return true
      }
    }
    return false
  }
  
  // Recursively filter out schema fields from metadata
  const filterSchemaFields = (obj: any, currentPath: string = ''): any => {
    if (obj === null || obj === undefined) {
      return obj
    }
    
    // If current path matches a schema field exactly, exclude it
    if (currentPath && schemaFieldPaths.has(currentPath)) {
      return undefined
    }
    
    // If it's an object, recursively check nested properties
    if (typeof obj === 'object' && !Array.isArray(obj)) {
      const filtered: Record<string, any> = {}
      for (const [key, value] of Object.entries(obj)) {
        const fullPath = currentPath ? `${currentPath}.${key}` : key
        
        // Skip if this path is a schema field
        if (isSchemaField(fullPath)) {
          continue
        }
        
        // Recursively filter nested objects
        const filteredValue = filterSchemaFields(value, fullPath)
        if (filteredValue !== undefined) {
          filtered[key] = filteredValue
        }
      }
      return Object.keys(filtered).length > 0 ? filtered : undefined
    }
    
    // For non-object values, check if current path is a schema field
    if (currentPath && isSchemaField(currentPath)) {
      return undefined
    }
    
    return obj
  }
  
  const filtered = filterSchemaFields(metadata)
  return filtered && typeof filtered === 'object' ? filtered : {}
})

function getNestedValue(obj: any, path: string): any {
  const parts = path.split('.')
  let current = obj
  for (const part of parts) {
    if (current === null || current === undefined) return undefined
    current = current[part]
  }
  return current
}

function formatValue(value: any): string {
  if (value === null) return 'null'
  if (value === undefined) return 'undefined'
  if (typeof value === 'object') return JSON.stringify(value, null, 2)
  return String(value)
}

function isComplexValue(value: any): boolean {
  return typeof value === 'object' && value !== null && !Array.isArray(value)
}

function formatJSON(obj: any): string {
  return JSON.stringify(obj, null, 2)
}

function formatDate(dateStr: string): string {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleString()
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

function handleClose() {
  isOpen.value = false
  emit('close')
}
</script>

<style scoped>
.transaction-detail {
  padding: 8px 0;
}

.detail-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.detail-label {
  font-size: 0.75rem;
  font-weight: 600;
  color: rgba(0, 0, 0, 0.6);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.detail-value {
  font-size: 0.875rem;
  color: rgba(0, 0, 0, 0.87);
}

.schema-fields-table {
  border: 1px solid rgba(0, 0, 0, 0.12);
  border-radius: 4px;
}

.field-value {
  max-width: 300px;
  overflow-x: auto;
}

.field-value pre {
  margin: 0;
  font-size: 0.75rem;
  white-space: pre-wrap;
  word-break: break-word;
}

.metadata-json {
  background: rgba(0, 0, 0, 0.05);
  padding: 16px;
  border-radius: 4px;
  font-size: 0.875rem;
  overflow-x: auto;
  max-height: 400px;
  overflow-y: auto;
}
</style>
