<template>
  <v-container fluid class="pa-8">
    <div class="d-flex justify-space-between align-center mb-10">
      <div>
        <div class="d-flex align-center gap-3 mb-2">
          <v-icon icon="fa-code" size="large" color="primary" />
          <h2 class="text-h4 font-weight-bold">Event Schemas</h2>
        </div>
        <p class="text-body-1 text-grey-darken-1">Define event structures by pasting sample JSON</p>
      </div>
      <BaseButton @click="openCreateModal" prepend-icon="fa-plus">
        Create Schema
      </BaseButton>
    </div>

    <LoadingSpinner v-if="loading" text="Loading schemas..." :centered="false" />

    <ErrorMessage
      v-else-if="error"
      title="Error loading schemas"
      :message="error"
      retryable
      @retry="loadSchemas"
    />

    <BaseTable
      v-else
      :columns="tableColumns"
      :data="schemas"
      empty-text="No schemas defined. Create your first schema to enable flexible rule evaluation."
    >
      <template #cell-name="{ row }">
        <div class="font-weight-semibold text-grey-darken-3">{{ row.name }}</div>
        <div class="text-body-2 text-grey-darken-1">{{ row.description || 'No description' }}</div>
      </template>

      <template #cell-fields="{ row }">
        <span class="text-body-2 font-weight-medium text-grey-darken-2">
          {{ row.extracted_fields?.length || 0 }} fields
        </span>
      </template>

      <template #cell-created="{ row }">
        <span class="text-body-2 text-grey-darken-1">
          {{ formatDate(row.created_at) }}
        </span>
      </template>

      <template #cell-actions="{ row }">
        <div class="d-flex gap-2">
          <BaseButton size="sm" @click="openViewModal(row)" prepend-icon="fa-eye">
            View
          </BaseButton>
          <BaseButton size="sm" @click="openEditModal(row)" prepend-icon="fa-pencil">
            Edit
          </BaseButton>
          <BaseButton size="sm" variant="danger" @click="deleteSchema(row.id)" prepend-icon="fa-trash">
            Delete
          </BaseButton>
        </div>
      </template>
    </BaseTable>

    <!-- Create/Edit Modal -->
    <BaseModal
      v-model="showModal"
      :title="isEditing ? 'Edit Schema' : 'Create New Schema'"
      size="lg"
    >
      <v-form ref="formRef" @submit.prevent="handleSubmit" class="mt-4">
        <BaseInput
          v-model="editingSchema.name"
          label="Name"
          placeholder="e.g., Payment Event, Kafka Order Event"
          required
          :rules="[(v: string) => !!v || 'Name is required']"
          prepend-inner-icon="fa-text"
          class="mb-4"
        />

        <BaseInput
          v-model="editingSchema.description"
          label="Description"
          placeholder="Brief description of this event schema"
          prepend-inner-icon="fa-align-left"
          class="mb-4"
        />

        <div class="mb-4">
          <label class="text-body-2 text-grey-darken-1 d-block mb-2">
            Sample JSON <span class="text-red">*</span>
          </label>
          <div class="json-editor-wrapper">
            <textarea
              v-model="sampleJsonText"
              class="json-editor"
              placeholder='{
  "event_id": "evt_123",
  "amount": 150.00,
  "customer": {
    "id": "cust_456",
    "country": "US"
  }
}'
              rows="12"
              @input="parseSampleJson"
            />
          </div>
          <p v-if="jsonError" class="text-caption text-red mt-1">{{ jsonError }}</p>
          <p class="text-caption text-grey-darken-1 mt-1">
            Paste a sample JSON event. Fields will be automatically extracted.
          </p>
        </div>

        <!-- Extracted Fields Preview -->
        <div v-if="extractedFields.length > 0" class="mb-4 pa-4 bg-grey-lighten-5 rounded-lg">
          <h4 class="text-subtitle-1 font-weight-medium mb-3 d-flex align-center">
            <v-icon icon="fa-list" size="small" class="mr-2" />
            Extracted Fields ({{ extractedFields.length }})
          </h4>
          <v-table density="compact">
            <thead>
              <tr>
                <th>Path</th>
                <th>Type</th>
                <th>Sample Value</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="field in extractedFields" :key="field.path">
                <td><code>{{ field.path }}</code></td>
                <td>
                  <BaseBadge :variant="getTypeBadgeVariant(field.type)" size="sm">
                    {{ field.type }}
                  </BaseBadge>
                </td>
                <td class="text-grey-darken-1">{{ formatSampleValue(field.sample_value) }}</td>
              </tr>
            </tbody>
          </v-table>
        </div>
      </v-form>

      <template #footer>
        <BaseButton variant="ghost" @click="closeModal" prepend-icon="fa-xmark">Cancel</BaseButton>
        <BaseButton 
          @click="handleSubmit" 
          :loading="saving" 
          :disabled="!!jsonError || extractedFields.length === 0"
          prepend-icon="fa-save"
        >
          Save
        </BaseButton>
      </template>
    </BaseModal>

    <!-- View Schema Modal -->
    <BaseModal
      v-model="showViewModal"
      :title="viewingSchema?.name || 'Schema Details'"
      size="lg"
    >
      <div v-if="viewingSchema" class="mt-4">
        <div class="mb-4">
          <label class="text-body-2 text-grey-darken-1 d-block mb-1">Description</label>
          <p class="text-body-1">{{ viewingSchema.description || 'No description' }}</p>
        </div>

        <div class="mb-4">
          <label class="text-body-2 text-grey-darken-1 d-block mb-2">Sample JSON</label>
          <pre class="json-preview">{{ JSON.stringify(viewingSchema.sample_json, null, 2) }}</pre>
        </div>

        <div class="mb-4 pa-4 bg-grey-lighten-5 rounded-lg">
          <h4 class="text-subtitle-1 font-weight-medium mb-3">
            Extracted Fields ({{ viewingSchema.extracted_fields?.length || 0 }})
          </h4>
          <v-table density="compact">
            <thead>
              <tr>
                <th>Path</th>
                <th>Type</th>
                <th>Sample Value</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="field in viewingSchema.extracted_fields" :key="field.path">
                <td><code>{{ field.path }}</code></td>
                <td>
                  <BaseBadge :variant="getTypeBadgeVariant(field.type)" size="sm">
                    {{ field.type }}
                  </BaseBadge>
                </td>
                <td class="text-grey-darken-1">{{ formatSampleValue(field.sample_value) }}</td>
              </tr>
            </tbody>
          </v-table>
        </div>
      </div>

      <template #footer>
        <BaseButton variant="ghost" @click="showViewModal = false" prepend-icon="fa-xmark">Close</BaseButton>
      </template>
    </BaseModal>
  </v-container>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { api } from '@/lib/api'
import BaseButton from '@/components/BaseButton.vue'
import BaseBadge from '@/components/BaseBadge.vue'
import BaseModal from '@/components/BaseModal.vue'
import BaseTable from '@/components/BaseTable.vue'
import BaseInput from '@/components/BaseInput.vue'
import LoadingSpinner from '@/components/LoadingSpinner.vue'
import ErrorMessage from '@/components/ErrorMessage.vue'

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
  sample_json: Record<string, any>
  extracted_fields: ExtractedField[]
  created_at: string
  updated_at: string
}

const router = useRouter()
const authStore = useAuthStore()

const tableColumns = [
  { key: 'name', label: 'Name' },
  { key: 'fields', label: 'Fields' },
  { key: 'created', label: 'Created' },
  { key: 'actions', label: 'Actions' },
]

const schemas = ref<EventSchema[]>([])
const loading = ref(true)
const error = ref('')
const showModal = ref(false)
const showViewModal = ref(false)
const isEditing = ref(false)
const saving = ref(false)
const formRef = ref<any>(null)

const editingSchema = reactive({
  id: '',
  name: '',
  description: '',
})

const sampleJsonText = ref('')
const jsonError = ref('')
const extractedFields = ref<ExtractedField[]>([])
const viewingSchema = ref<EventSchema | null>(null)

// Maximum nesting depth for field extraction
const MAX_DEPTH = 5

onMounted(() => {
  if (authStore.user) {
    loadSchemas()
  } else {
    router.push('/login')
  }
})

async function loadSchemas() {
  try {
    loading.value = true
    error.value = ''
    const response = await api.get<{ schemas: EventSchema[] }>('/api/v1/schemas')
    schemas.value = response?.schemas || []
  } catch (e: any) {
    error.value = e.message || 'Failed to load schemas'
    console.error('Error loading schemas:', e)
  } finally {
    loading.value = false
  }
}

function openCreateModal() {
  isEditing.value = false
  editingSchema.id = ''
  editingSchema.name = ''
  editingSchema.description = ''
  sampleJsonText.value = ''
  jsonError.value = ''
  extractedFields.value = []
  showModal.value = true
}

function openEditModal(schema: EventSchema) {
  isEditing.value = true
  editingSchema.id = schema.id
  editingSchema.name = schema.name
  editingSchema.description = schema.description || ''
  sampleJsonText.value = JSON.stringify(schema.sample_json, null, 2)
  extractedFields.value = schema.extracted_fields || []
  jsonError.value = ''
  showModal.value = true
}

function openViewModal(schema: EventSchema) {
  viewingSchema.value = schema
  showViewModal.value = true
}

function closeModal() {
  showModal.value = false
}

function parseSampleJson() {
  if (!sampleJsonText.value.trim()) {
    jsonError.value = ''
    extractedFields.value = []
    return
  }

  try {
    const parsed = JSON.parse(sampleJsonText.value)
    if (typeof parsed !== 'object' || Array.isArray(parsed) || parsed === null) {
      jsonError.value = 'Sample JSON must be an object'
      extractedFields.value = []
      return
    }
    jsonError.value = ''
    extractedFields.value = extractFieldsFromJson(parsed, '', 0)
  } catch (e) {
    jsonError.value = 'Invalid JSON format'
    extractedFields.value = []
  }
}

function extractFieldsFromJson(data: Record<string, any>, prefix: string, depth: number): ExtractedField[] {
  if (depth >= MAX_DEPTH) return []

  const fields: ExtractedField[] = []

  for (const [key, value] of Object.entries(data)) {
    const path = prefix ? `${prefix}.${key}` : key
    const fieldType = inferType(value)

    if (fieldType === 'object' && value !== null && !Array.isArray(value)) {
      // Recurse into nested objects
      const nestedFields = extractFieldsFromJson(value, path, depth + 1)
      fields.push(...nestedFields)
    } else {
      fields.push({
        path,
        type: fieldType,
        nullable: value === null,
        sample_value: value,
      })
    }
  }

  return fields
}

function inferType(value: any): string {
  if (value === null) return 'null'
  if (typeof value === 'boolean') return 'boolean'
  if (typeof value === 'number') return 'number'
  if (typeof value === 'string') return 'string'
  if (Array.isArray(value)) return 'array'
  if (typeof value === 'object') return 'object'
  return 'string'
}

async function handleSubmit() {
  if (formRef.value) {
    const { valid } = await formRef.value.validate()
    if (!valid) return
  }

  if (!sampleJsonText.value.trim() || jsonError.value) {
    error.value = 'Valid sample JSON is required'
    return
  }

  try {
    saving.value = true
    error.value = ''

    const sampleJson = JSON.parse(sampleJsonText.value)
    const payload = {
      name: editingSchema.name,
      description: editingSchema.description,
      sample_json: sampleJson,
    }

    if (isEditing.value && editingSchema.id) {
      await api.put(`/api/v1/schemas/${editingSchema.id}`, payload)
    } else {
      await api.post('/api/v1/schemas', payload)
    }

    closeModal()
    await loadSchemas()
  } catch (e: any) {
    error.value = e.message || 'Failed to save schema'
  } finally {
    saving.value = false
  }
}

async function deleteSchema(id: string) {
  if (!confirm('Are you sure you want to delete this schema?')) return

  try {
    await api.delete(`/api/v1/schemas/${id}`)
    await loadSchemas()
  } catch (e: any) {
    error.value = e.message || 'Failed to delete schema'
  }
}

function formatDate(dateStr: string): string {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}

function formatSampleValue(value: any): string {
  if (value === null) return 'null'
  if (value === undefined) return ''
  if (typeof value === 'string') return `"${value.substring(0, 30)}${value.length > 30 ? '...' : ''}"`
  if (typeof value === 'object') return JSON.stringify(value).substring(0, 40) + '...'
  return String(value)
}

function getTypeBadgeVariant(type: string): 'success' | 'warning' | 'info' | 'default' {
  switch (type) {
    case 'string': return 'success'
    case 'number': return 'info'
    case 'boolean': return 'warning'
    case 'array': return 'default'
    default: return 'default'
  }
}

// Watch for JSON text changes with debounce
let debounceTimer: ReturnType<typeof setTimeout>
watch(sampleJsonText, () => {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => {
    parseSampleJson()
  }, 300)
})
</script>

<style scoped>
.json-editor-wrapper {
  position: relative;
  font-family: 'Fira Code', 'Monaco', 'Consolas', monospace;
  font-size: 14px;
  line-height: 1.5;
}

.json-editor {
  width: 100%;
  min-height: 200px;
  padding: 12px;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  font-family: inherit;
  font-size: inherit;
  line-height: inherit;
  resize: vertical;
  background: #f5f5f5;
}

.json-editor:focus {
  outline: none;
  border-color: #1976d2;
  box-shadow: 0 0 0 2px rgba(25, 118, 210, 0.2);
}

.json-preview {
  background: #f5f5f5;
  padding: 12px;
  border-radius: 8px;
  overflow-x: auto;
  font-family: 'Fira Code', 'Monaco', 'Consolas', monospace;
  font-size: 13px;
  line-height: 1.5;
  max-height: 300px;
  overflow-y: auto;
}

code {
  background: #e8e8e8;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 13px;
}
</style>

