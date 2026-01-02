<template>
  <v-container fluid class="pa-8">
    <div class="d-flex justify-space-between align-center mb-10">
      <div>
        <div class="d-flex align-center gap-3 mb-2">
          <v-icon icon="fa-tasks" size="large" color="primary" />
          <h2 class="text-h4 font-weight-bold">Rules Management</h2>
        </div>
        <p class="text-body-1 text-grey-darken-1">Create schema-based rules using expressions to evaluate events</p>
      </div>
      <BaseButton @click="openCreateModal" prepend-icon="fa-plus">
        Create Rule
      </BaseButton>
    </div>

    <LoadingSpinner v-if="loading" text="Loading rules..." :centered="false" />

    <ErrorMessage
      v-else-if="error"
      title="Error loading rules"
      :message="error"
      retryable
      @retry="loadRules"
    />

    <BaseTable
      v-else
      :columns="tableColumns"
      :data="rules"
      empty-text="No rules configured. Create your first rule to get started."
    >
      <template #cell-name="{ row }">
        <div class="font-weight-semibold text-grey-darken-3">{{ row.name }}</div>
        <div class="text-body-2 text-grey-darken-1">{{ row.description }}</div>
      </template>

      <template #cell-type="{ value }">
        <span class="text-body-2 font-weight-medium text-grey-darken-2">{{ value }}</span>
      </template>

      <template #cell-action="{ row }">
        <BaseBadge :variant="getActionBadgeVariant(row.action)" rounded>
          {{ row.action }}
        </BaseBadge>
      </template>

      <template #cell-score="{ value }">
        <span class="text-body-2 font-weight-medium text-grey-darken-2">{{ value }}</span>
      </template>

      <template #cell-priority="{ value }">
        <span class="text-body-2 font-weight-medium text-grey-darken-2">{{ value }}</span>
      </template>

      <template #cell-status="{ row }">
        <BaseBadge
          :variant="row.enabled ? 'success' : 'default'"
          :outline="!row.enabled"
          @click="toggleRule(row)"
          class="cursor-pointer"
        >
          {{ row.enabled ? 'Enabled' : 'Disabled' }}
        </BaseBadge>
      </template>

      <template #cell-actions="{ row }">
        <div class="d-flex gap-2">
          <BaseButton size="sm" @click="openEditModal(row)" prepend-icon="fa-pencil">
            Edit
          </BaseButton>
          <BaseButton size="sm" variant="danger" @click="deleteRule(row.id)" prepend-icon="fa-trash">
            Delete
          </BaseButton>
        </div>
      </template>
    </BaseTable>

    <BaseModal
      v-model="showModal"
      :title="isEditing ? 'Edit Rule' : 'Create New Rule'"
      size="lg"
    >
      <v-form ref="formRef" @submit.prevent="handleSubmit" class="mt-4">
        <!-- Presets Section (only for new rules) -->
        <div v-if="!isEditing" class="mb-6">
          <label class="text-body-2 text-grey-darken-1 d-block mb-2">Quick Start with Preset</label>
          <div class="d-flex flex-wrap gap-2">
            <v-chip
              v-for="preset in rulePresets"
              :key="preset.id"
              @click="applyPreset(preset)"
              variant="outlined"
              color="primary"
              class="cursor-pointer"
            >
              <v-icon :icon="preset.icon" size="small" class="mr-1" />
              {{ preset.label }}
            </v-chip>
          </div>
        </div>

        <BaseInput
          v-model="editingRule.name"
          label="Rule Name"
          placeholder="e.g., High Value Transaction, Suspicious Activity"
          required
          :rules="[(v: string) => !!v || 'Name is required']"
          prepend-inner-icon="fa-text"
          hint="A descriptive name for this rule"
          persistent-hint
          class="mb-4"
        />

        <BaseInput
          v-model="editingRule.description"
          label="Description"
          placeholder="Describe what this rule checks for"
          required
          :rules="[(v: string) => !!v || 'Description is required']"
          prepend-inner-icon="fa-align-left"
          hint="Explain when this rule should trigger"
          persistent-hint
          class="mb-4"
        />

        <BaseSelect
          v-model="editingRule.schema_id"
          label="Event Schema"
          :options="schemaOptions"
          :rules="[(v: string) => !!v || 'Schema is required']"
          hint="Choose the event schema that defines the structure of events this rule will evaluate"
          persistent-hint
          class="mb-4"
        />

        <!-- Available Fields from Selected Schema -->
        <div v-if="selectedSchema && selectedSchema.extracted_fields?.length > 0" class="mb-4 pa-3 bg-blue-lighten-5 rounded-lg">
          <div class="d-flex justify-space-between align-center mb-2">
            <label class="text-caption text-grey-darken-1 d-flex align-center">
              <v-icon icon="fa-info-circle" size="x-small" class="mr-1" />
              Available fields from "{{ selectedSchema.name }}" ({{ selectedSchema.extracted_fields.length }}):
            </label>
            <v-btn
              v-if="selectedSchema.extracted_fields.length > 10"
              size="x-small"
              variant="text"
              @click="showAllFields = !showAllFields"
            >
              {{ showAllFields ? 'Show Less' : `Show All (${selectedSchema.extracted_fields.length})` }}
            </v-btn>
          </div>
          <div class="d-flex flex-wrap gap-1">
            <v-chip
              v-for="field in (showAllFields ? selectedSchema.extracted_fields : selectedSchema.extracted_fields.slice(0, 10))"
              :key="field.path"
              size="x-small"
              variant="outlined"
              color="primary"
            >
              {{ field.path }} <span class="text-grey ml-1">({{ field.type }})</span>
            </v-chip>
          </div>
          <p class="text-caption text-grey-darken-1 mt-2">
            Available fields for use in conditions below.
          </p>
        </div>

        <BaseSelect
          v-model="editingRule.action"
          label="Action"
          :options="ruleActions"
          required
          :rules="[(v: string) => !!v || 'Action is required']"
          class="mb-4"
        />

        <BaseInput
          v-model.number="editingRule.score"
          type="number"
          label="Score"
          :min="0"
          :max="100"
          required
          :rules="[
            (v: any) => (v !== null && v !== undefined && v !== '') || 'Score is required',
            (v: any) => (typeof v === 'number' && v >= 0 && v <= 100) || 'Score must be between 0 and 100'
          ]"
          prepend-inner-icon="fa-hashtag"
          class="mb-4"
        />

        <BaseInput
          v-model.number="editingRule.priority"
          type="number"
          label="Priority"
          :min="0"
          :max="1000"
          required
          :rules="[
            (v: any) => (v !== null && v !== undefined && v !== '') || 'Priority is required',
            (v: any) => (typeof v === 'number' && v >= 0 && v <= 1000) || 'Priority must be between 0 and 1000'
          ]"
          hint="0 = highest priority, 1000 = lowest priority"
          persistent-hint
          prepend-inner-icon="fa-sort"
          class="mb-4"
        />

        <!-- Rule Expression Section -->
        <div v-if="selectedSchema" class="mb-6 pa-4 bg-grey-lighten-5 rounded-lg">
          <div class="d-flex justify-space-between align-center mb-4">
            <h4 class="text-subtitle-1 font-weight-medium d-flex align-center">
              <v-icon icon="fa-code" size="small" class="mr-2" />
              Rule Conditions <span class="text-red">*</span>
            </h4>
            <BaseButton size="sm" @click="addCondition" prepend-icon="fa-plus">
              Add Condition
            </BaseButton>
          </div>

          <div v-if="conditionRows.length === 0" class="text-center pa-6 text-grey-darken-1">
            <v-icon icon="fa-info-circle" size="large" class="mb-2" />
            <p class="text-body-2">No conditions added yet. Click "Add Condition" to start building your rule.</p>
          </div>

          <div v-for="(row, index) in conditionRows" :key="row.id" class="mb-4 pa-3 bg-white rounded-lg border">
            <div class="d-flex align-center gap-2">
              <div class="flex-grow-1">
                <div class="d-flex align-center gap-2 flex-wrap">
                  <!-- Field Dropdown -->
                  <BaseSelect
                    v-model="row.field"
                    label="Field"
                    :options="fieldOptions"
                    :rules="[(v: string) => !!v || 'Field is required']"
                    class="flex-grow-1"
                    style="min-width: 200px;"
                  />

                  <!-- Comparison Operator -->
                  <BaseSelect
                    v-model="row.operator"
                    label="Operator"
                    :options="getOperatorOptions(row.field)"
                    :rules="[(v: string) => !!v || 'Operator is required']"
                    style="min-width: 150px;"
                  />

                  <!-- Value Input -->
                  <div class="flex-grow-1" style="min-width: 200px;">
                    <BaseInput
                      v-if="getValueInputType(row.field) === 'number'"
                      v-model.number="row.value"
                      type="number"
                      label="Value"
                      :rules="[(v: any) => v !== null && v !== undefined && v !== '' || 'Value is required']"
                    />
                    <div v-else-if="getValueInputType(row.field) === 'boolean'">
                      <label class="text-body-2 text-grey-darken-1 d-block mb-1">Value <span class="text-red">*</span></label>
                      <v-select
                        v-model="row.value"
                        :items="[{ title: 'true', value: true }, { title: 'false', value: false }]"
                        variant="outlined"
                        density="compact"
                        :rules="[(v: any) => v !== null && v !== undefined && v !== '' || 'Value is required']"
                      />
                    </div>
                    <BaseInput
                      v-else-if="getValueInputType(row.field) === 'array'"
                      v-model="row.value"
                      label="Value (comma-separated)"
                      placeholder='e.g., "US", "CA", "GB" or 1, 2, 3'
                      :rules="[(v: string) => !!v || 'Value is required']"
                      hint="Enter comma-separated values. Strings should be quoted."
                      persistent-hint
                    />
                    <BaseInput
                      v-else
                      v-model="row.value"
                      label="Value"
                      :rules="[(v: string) => !!v || 'Value is required']"
                      :placeholder="getValuePlaceholder(row.field)"
                    />
                  </div>
                </div>
              </div>

              <!-- Remove Button -->
              <v-btn
                icon="fa-trash"
                size="small"
                variant="text"
                color="error"
                @click="removeCondition(index)"
                class="mt-6"
              />
            </div>

            <!-- Logic Connector (AND/OR) -->
            <div v-if="index < conditionRows.length - 1" class="mt-3 d-flex align-center">
              <v-divider class="flex-grow-1" />
              <v-btn-toggle
                v-model="row.logicOperator"
                mandatory
                color="primary"
                variant="outlined"
                density="compact"
                class="mx-2"
              >
                <v-btn value="and" size="small">
                  AND
                </v-btn>
                <v-btn value="or" size="small">
                  OR
                </v-btn>
              </v-btn-toggle>
              <v-divider class="flex-grow-1" />
            </div>
          </div>

          <!-- Generated Expression Preview -->
          <div v-if="conditionRows.length > 0" class="mt-4 pa-3 bg-grey-darken-1 rounded">
            <div class="text-caption text-white mb-1">Generated Expression:</div>
            <code class="text-white">{{ generatedExpression }}</code>
          </div>
        </div>

        <div class="mb-6">
          <label class="text-body-2 text-grey-darken-1 d-block mb-2">Status</label>
          <v-btn-toggle
            v-model="editingRule.enabled"
            mandatory
            color="primary"
            variant="outlined"
          >
            <v-btn :value="true" class="px-6">
              <v-icon icon="fa-check" size="small" class="mr-2" />
              Enabled
            </v-btn>
            <v-btn :value="false" class="px-6">
              <v-icon icon="fa-ban" size="small" class="mr-2" />
              Disabled
            </v-btn>
          </v-btn-toggle>
        </div>
      </v-form>

      <template #footer>
        <BaseButton variant="ghost" @click="closeModal" prepend-icon="fa-xmark">Cancel</BaseButton>
        <BaseButton @click="handleSubmit" :loading="saving" prepend-icon="fa-save">Save</BaseButton>
      </template>
    </BaseModal>
  </v-container>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch, computed, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { api } from '@/lib/api'
import BaseButton from '@/components/BaseButton.vue'
import BaseBadge from '@/components/BaseBadge.vue'
import BaseModal from '@/components/BaseModal.vue'
import BaseTable from '@/components/BaseTable.vue'
import BaseSelect from '@/components/BaseSelect.vue'
import BaseInput from '@/components/BaseInput.vue'
import LoadingSpinner from '@/components/LoadingSpinner.vue'
import ErrorMessage from '@/components/ErrorMessage.vue'
import Prism from 'prismjs'
import 'prismjs/themes/prism.css'

const router = useRouter()
const authStore = useAuthStore()

const tableColumns = [
  { key: 'name', label: 'Name' },
  { key: 'type', label: 'Type' },
  { key: 'action', label: 'Action' },
  { key: 'score', label: 'Score' },
  { key: 'priority', label: 'Priority' },
  { key: 'status', label: 'Status' },
  { key: 'actions', label: 'Actions' },
]

const rules = ref<any[]>([])
const loading = ref(true)
const error = ref('')
const showModal = ref(false)
const isEditing = ref(false)
const saving = ref(false)
const formRef = ref<any>(null)
const customExpressionRef = ref<HTMLTextAreaElement | null>(null)
const highlightRef = ref<HTMLPreElement | null>(null)

interface RuleForm {
  id?: string
  name: string
  description: string
  type: string
  action: string
  score: number
  priority: number
  enabled: boolean
  conditions: Record<string, any>
  schema_id?: string
}

interface ExtractedField {
  path: string
  type: string
  nullable?: boolean
  sample_value?: any
}

interface EventSchema {
  id: string
  name: string
  description?: string
  extracted_fields: ExtractedField[]
}

interface ConditionRow {
  id: string
  field: string
  operator: string
  value: any
  logicOperator: 'and' | 'or'
}

const editingRule = reactive<RuleForm>({
  name: '',
  description: '',
  type: '',
  action: '',
  score: 0,
  priority: 500,
  enabled: true,
  conditions: {},
  schema_id: undefined,
})

// Structured condition rows for the form
const conditionRows = ref<ConditionRow[]>([])
let conditionRowIdCounter = 0

// Schemas for selection
const schemas = ref<EventSchema[]>([])
const selectedSchema = ref<EventSchema | null>(null)
const showAllFields = ref(false)

// Computed list of schema options for the dropdown
const schemaOptions = computed(() => {
  return schemas.value.map(s => ({
    value: s.id,
    label: s.name,
  }))
})

// Computed list of field options from selected schema
const fieldOptions = computed(() => {
  if (!selectedSchema.value || !selectedSchema.value.extracted_fields) {
    return []
  }
  return selectedSchema.value.extracted_fields.map(f => ({
    value: f.path,
    label: `${f.path} (${f.type})`,
  }))
})

// Get operator options based on field type
function getOperatorOptions(fieldPath: string): Array<{ value: string; label: string }> {
  if (!selectedSchema.value || !fieldPath) {
    return []
  }
  
  const field = selectedSchema.value.extracted_fields.find(f => f.path === fieldPath)
  if (!field) {
    return []
  }

  const type = field.type
  
  if (type === 'number') {
    return [
      { value: '==', label: 'Equals (==)' },
      { value: '!=', label: 'Not Equals (!=)' },
      { value: '>', label: 'Greater Than (>)' },
      { value: '<', label: 'Less Than (<)' },
      { value: '>=', label: 'Greater or Equal (>=)' },
      { value: '<=', label: 'Less or Equal (<=)' },
    ]
  } else if (type === 'string') {
    return [
      { value: '==', label: 'Equals (==)' },
      { value: '!=', label: 'Not Equals (!=)' },
      { value: 'in', label: 'In List (in)' },
      { value: 'contains', label: 'Contains (contains)' },
    ]
  } else if (type === 'boolean') {
    return [
      { value: '==', label: 'Equals (==)' },
      { value: '!=', label: 'Not Equals (!=)' },
    ]
  } else if (type === 'array') {
    return [
      { value: 'in', label: 'Item In Array (in)' },
      { value: 'contains', label: 'Contains Item (contains)' },
    ]
  }
  
  // Default operators for unknown types
  return [
    { value: '==', label: 'Equals (==)' },
    { value: '!=', label: 'Not Equals (!=)' },
  ]
}

// Get input type for value field based on field type
function getValueInputType(fieldPath: string): 'number' | 'string' | 'boolean' | 'array' {
  if (!selectedSchema.value || !fieldPath) {
    return 'string'
  }
  
  const field = selectedSchema.value.extracted_fields.find(f => f.path === fieldPath)
  if (!field) {
    return 'string'
  }

  const type = field.type
  if (type === 'number') {
    return 'number'
  } else if (type === 'boolean') {
    return 'boolean'
  } else if (type === 'array') {
    return 'array'
  }
  
  return 'string'
}

// Get placeholder for value input based on field type and operator
function getValuePlaceholder(fieldPath: string): string {
  if (!selectedSchema.value || !fieldPath) {
    return 'Enter value'
  }
  
  const field = selectedSchema.value.extracted_fields.find(f => f.path === fieldPath)
  if (!field) {
    return 'Enter value'
  }

  const type = field.type
  if (type === 'string') {
    return 'e.g., "USD" or "active"'
  } else if (type === 'number') {
    return 'e.g., 1000 or 5000'
  } else if (type === 'array') {
    return 'e.g., "US", "CA", "GB" or 1, 2, 3'
  }
  
  return 'Enter value'
}

// Add a new condition row
function addCondition() {
  conditionRows.value.push({
    id: `condition-${++conditionRowIdCounter}`,
    field: '',
    operator: '',
    value: null,
    logicOperator: 'and',
  })
}

// Remove a condition row
function removeCondition(index: number) {
  conditionRows.value.splice(index, 1)
}

// Generate expression from condition rows
const generatedExpression = computed(() => {
  if (conditionRows.value.length === 0) {
    return ''
  }

  const parts: string[] = []
  
  for (let i = 0; i < conditionRows.value.length; i++) {
    const row = conditionRows.value[i]
    
    if (!row) {
      continue
    }
    
    if (!row.field || !row.operator || row.value === null || row.value === undefined || row.value === '') {
      continue
    }

    // Build the condition part
    let conditionPart = `${row.field} ${row.operator} `
    
    // Format value based on type and operator
    const field = selectedSchema.value?.extracted_fields.find(f => f.path === row.field)
    const fieldType = field?.type || 'string'
    
    if (row.operator === 'in' || row.operator === 'contains') {
      // For 'in' and 'contains', value should be an array
      if (fieldType === 'string' && typeof row.value === 'string') {
        // Parse comma-separated string values
        const values = row.value.split(',').map(v => v.trim()).filter(v => v)
        const formattedValues = values.map(v => {
          // If it's already quoted, keep it; otherwise quote it
          if ((v.startsWith('"') && v.endsWith('"')) || (v.startsWith("'") && v.endsWith("'"))) {
            return v
          }
          return `"${v}"`
        })
        conditionPart += `[${formattedValues.join(', ')}]`
      } else if (fieldType === 'number' && typeof row.value === 'string') {
        // Parse comma-separated numbers
        const values = row.value.split(',').map(v => v.trim()).filter(v => v)
        conditionPart += `[${values.join(', ')}]`
      } else {
        conditionPart += row.value
      }
    } else if (fieldType === 'string' && typeof row.value === 'string') {
      // String values need quotes unless already quoted
      if (!row.value.startsWith('"') && !row.value.startsWith("'")) {
        conditionPart += `"${row.value}"`
      } else {
        conditionPart += row.value
      }
    } else if (fieldType === 'boolean') {
      // Boolean values don't need quotes
      conditionPart += row.value
    } else {
      // Numbers and other types
      conditionPart += row.value
    }
    
    parts.push(conditionPart)
    
    // Add logic operator if not the last condition
    if (i < conditionRows.value.length - 1 && row.logicOperator) {
      parts.push(row.logicOperator)
    }
  }
  
  return parts.join(' ')
})

// Only custom rule type is supported (schema-based expressions)
const ruleTypes = [
  { value: 'custom', label: 'Custom Expression' },
]

// Rule actions must match backend validation: allow|block|review|score
const ruleActions = [
  { value: 'allow', label: 'Allow' },
  { value: 'block', label: 'Block' },
  { value: 'review', label: 'Review' },
  { value: 'score', label: 'Score Only' },
]

// Preset configurations for common fraud detection scenarios
interface RulePreset {
  id: string
  label: string
  icon: string
  name: string
  description: string
  type: string
  action: string
  score: number
  conditions: Record<string, any>
}

const rulePresets: RulePreset[] = [
  {
    id: 'high-value',
    label: 'High Value',
    icon: 'fa-dollar-sign',
    name: 'High Value Transaction',
    description: 'Flag events with high value amounts',
    type: 'custom',
    action: 'review',
    score: 30,
    conditions: { custom_expression: 'amount > 10000' },
  },
  {
    id: 'suspicious-flag',
    label: 'Suspicious Flag',
    icon: 'fa-exclamation-triangle',
    name: 'Suspicious Flag Check',
    description: 'Check for suspicious flag in event metadata',
    type: 'custom',
    action: 'review',
    score: 60,
    conditions: { custom_expression: 'metadata.is_suspicious == true' },
  },
  {
    id: 'high-amount-foreign',
    label: 'High Amount Foreign',
    icon: 'fa-coins',
    name: 'High Amount Foreign Currency',
    description: 'Flag high-value events in non-USD currency',
    type: 'custom',
    action: 'review',
    score: 35,
    conditions: { custom_expression: 'amount > 5000 and currency != "USD"' },
  },
  {
    id: 'country-restriction',
    label: 'Country Restriction',
    icon: 'fa-globe',
    name: 'Restricted Country',
    description: 'Flag events from restricted countries',
    type: 'custom',
    action: 'review',
    score: 40,
    conditions: { custom_expression: 'user.country in ["RU", "CN", "KP"]' },
  },
]

// Placeholder for custom expression textarea (computed based on selected schema)
// Syntax follows expr-lang: https://github.com/expr-lang/expr
const customExpressionPlaceholder = computed(() => {
  if (!selectedSchema.value || !selectedSchema.value.extracted_fields?.length) {
    return 'Select a schema and click on fields above to build your expression...\n\nExamples:\nfield_name > 10000\nfield_name == "value"\nfield1 > 5000 and field2 != "USD"'
  }
  const firstField = selectedSchema.value.extracted_fields[0]?.path || 'field_name'
  return `Click on fields above or type directly:\n\nExamples:\n${firstField} > 10000\n${firstField} == "value"\n${firstField} > 5000 and ${selectedSchema.value.extracted_fields[1]?.path || 'other_field'} != "USD"`
})

// Sanitize expression to prevent code injection
// Only allows safe expr-lang syntax
function sanitizeExpression(expression: string): { valid: boolean; error?: string } {
  if (!expression || expression.trim() === '') {
    return { valid: false, error: 'Expression is required' }
  }

  // Maximum length check
  if (expression.length > 1000) {
    return { valid: false, error: 'Expression too long (max 1000 characters)' }
  }

  // Disallow dangerous patterns
  const dangerousPatterns = [
    /\bimport\b/i,          // import statements
    /\brequire\b/i,         // require calls
    /\beval\b/i,            // eval
    /\bexec\b/i,            // exec
    /\bFunction\b/,         // Function constructor
    /\b__\w+__\b/,          // dunder methods
    /\$\{/,                 // template literals
    /`/,                    // backticks
    /\bprocess\b/i,         // process access
    /\bwindow\b/i,          // window access
    /\bdocument\b/i,        // document access
    /\bglobal\b/i,          // global access
    /;\s*$/,                // trailing semicolons (multiple statements)
    /\bfor\b/i,             // for loops
    /\bwhile\b/i,           // while loops
    /\bfunction\b/i,        // function definitions
    /=>/,                   // arrow functions
  ]

  for (const pattern of dangerousPatterns) {
    if (pattern.test(expression)) {
      return { valid: false, error: 'Expression contains disallowed syntax' }
    }
  }

  // Only allow safe characters: alphanumeric, operators, parentheses, brackets, quotes, dots, underscores, spaces
  const safePattern = /^[a-zA-Z0-9_\s\.\,\(\)\[\]"'<>=!&|+\-*/%]+$/
  if (!safePattern.test(expression)) {
    return { valid: false, error: 'Expression contains invalid characters' }
  }

  return { valid: true }
}


// Get default conditions (only custom expressions supported)
function getDefaultConditions(type: string): Record<string, any> {
  return { custom_expression: '' }
}

// Insert field into expression at cursor position
function insertFieldIntoExpression(fieldPath: string) {
  if (!customExpressionRef.value) return
  
  const textarea = customExpressionRef.value
  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const currentExpression = editingRule.conditions.custom_expression || ''
  
  // Insert field at cursor position
  const newExpression = currentExpression.substring(0, start) + fieldPath + currentExpression.substring(end)
  editingRule.conditions.custom_expression = newExpression
  
  // Set cursor position after inserted field
  nextTick(() => {
    const newPosition = start + fieldPath.length
    textarea.setSelectionRange(newPosition, newPosition)
    textarea.focus()
    updateHighlighting()
  })
}

// Watch for type changes and reset conditions
watch(
  () => editingRule.type,
  (newType, oldType) => {
    if (newType && newType !== oldType && !isEditing.value) {
      // Only reset conditions for new rules or when type actually changes
      editingRule.conditions = getDefaultConditions(newType)
    }
  }
)

// Watch for schema_id changes and update selectedSchema
watch(
  () => editingRule.schema_id,
  (newSchemaId, oldSchemaId) => {
    if (newSchemaId) {
      selectedSchema.value = schemas.value.find(s => s.id === newSchemaId) || null
      // Reset condition rows when schema changes (clear field selections)
      if (oldSchemaId && newSchemaId !== oldSchemaId) {
        conditionRows.value.forEach(row => {
          row.field = ''
          row.operator = ''
          row.value = null
        })
      }
    } else {
      selectedSchema.value = null
    }
  }
)

// Syntax highlighting for custom expressions
const highlightedExpression = computed(() => {
  const expression = editingRule.conditions.custom_expression || ''
  if (!expression) return ''
  try {
    const grammar = Prism.languages.javascript
    if (grammar) {
      return Prism.highlight(expression, grammar, 'javascript')
    }
    return expression
  } catch {
    return expression
  }
})

// Generate dynamic examples based on selected schema fields
const dynamicExamples = computed(() => {
  if (!selectedSchema.value || !selectedSchema.value.extracted_fields?.length) {
    return 'Select a schema to see examples based on your fields.'
  }

  const fields = selectedSchema.value.extracted_fields
  const examples: string[] = []

  // Group fields by type
  const numberFields = fields.filter(f => f.type === 'number')
  const stringFields = fields.filter(f => f.type === 'string')
  const booleanFields = fields.filter(f => f.type === 'boolean')
  const arrayFields = fields.filter(f => f.type === 'array')

  // Number field examples (show comparison operators)
  if (numberFields.length > 0 && numberFields[0]) {
    const field = numberFields[0].path
    examples.push(`• <code>${field} &gt; 10000</code> (greater than)`)
    examples.push(`• <code>${field} &lt; 1000</code> (less than)`)
    examples.push(`• <code>${field} &gt;= 5000</code> (greater or equal)`)
    examples.push(`• <code>${field} &lt;= 100</code> (less or equal)`)
    examples.push(`• <code>${field} == 0</code> (equals)`)
    examples.push(`• <code>${field} != -1</code> (not equals)`)
    
    // Combined number examples
    if (numberFields.length > 1 && numberFields[1]) {
      examples.push(`• <code>${field} &gt; 1000 and ${numberFields[1].path} &lt; 500</code> (combine with and)`)
      examples.push(`• <code>${field} &gt; 5000 or ${numberFields[1].path} &lt; 100</code> (combine with or)`)
    }
  }

  // String field examples (show string operators)
  if (stringFields.length > 0 && stringFields[0]) {
    const field = stringFields[0].path
    examples.push(`• <code>${field} == "value"</code> (equals)`)
    examples.push(`• <code>${field} != "blocked"</code> (not equals)`)
    examples.push(`• <code>${field} in ["US", "CA", "GB"]</code> (in list)`)
    examples.push(`• <code>"test" in ${field}</code> (contains substring)`)
    
    // Combined string examples
    if (stringFields.length > 1 && stringFields[1]) {
      examples.push(`• <code>${field} == "USD" or ${stringFields[1].path} == "EUR"</code> (combine with or)`)
      examples.push(`• <code>${field} != "blocked" and ${stringFields[1].path} != "suspended"</code> (combine with and)`)
    }
  }

  // Boolean field examples
  if (booleanFields.length > 0 && booleanFields[0]) {
    const field = booleanFields[0].path
    examples.push(`• <code>${field} == true</code> (is true)`)
    examples.push(`• <code>${field} == false</code> (is false)`)
    examples.push(`• <code>not ${field}</code> (negation)`)
    
    // Combined boolean examples
    if (booleanFields.length > 1 && booleanFields[1]) {
      examples.push(`• <code>${field} and ${booleanFields[1].path}</code> (both true)`)
      examples.push(`• <code>${field} or ${booleanFields[1].path}</code> (either true)`)
    }
  }

  // Array field examples
  if (arrayFields.length > 0 && arrayFields[0]) {
    const field = arrayFields[0].path
    examples.push(`• <code>"item" in ${field}</code> (item in array)`)
    examples.push(`• <code>not ("blocked" in ${field})</code> (item not in array)`)
  }

  // Combined examples with multiple field types
  if (numberFields.length > 0 && numberFields[0] && stringFields.length > 0 && stringFields[0]) {
    examples.push(`• <code>${numberFields[0].path} &gt; 5000 and ${stringFields[0].path} == "USD"</code> (number and string)`)
    examples.push(`• <code>${numberFields[0].path} &lt; 100 or ${stringFields[0].path} == "test"</code> (number or string)`)
  }

  if (numberFields.length > 0 && numberFields[0] && booleanFields.length > 0 && booleanFields[0]) {
    examples.push(`• <code>${numberFields[0].path} &gt; 1000 and ${booleanFields[0].path}</code> (number and boolean)`)
  }

  if (stringFields.length > 0 && stringFields[0] && booleanFields.length > 0 && booleanFields[0]) {
    examples.push(`• <code>${stringFields[0].path} == "active" and ${booleanFields[0].path}</code> (string and boolean)`)
  }

  // Nested field examples (fields with dots)
  const nestedFields = fields.filter(f => f.path.includes('.'))
  if (nestedFields.length > 0 && nestedFields[0]) {
    const nestedField = nestedFields[0]
    if (nestedField.type === 'number') {
      examples.push(`• <code>${nestedField.path} &gt; 100</code> (nested number field)`)
    } else if (nestedField.type === 'string') {
      examples.push(`• <code>${nestedField.path} == "value"</code> (nested string field)`)
    } else if (nestedField.type === 'boolean') {
      examples.push(`• <code>${nestedField.path} == true</code> (nested boolean field)`)
    }
  }

  // Complex examples with not operator
  if (numberFields.length > 0 && numberFields[0]) {
    examples.push(`• <code>not (${numberFields[0].path} &lt; 0)</code> (using not)`)
  }
  if (stringFields.length > 0 && stringFields[0]) {
    examples.push(`• <code>not (${stringFields[0].path} in ["blocked", "suspended"])</code> (using not with in)`)
  }

  // Limit to 10 examples to provide good coverage without overwhelming
  return examples.slice(0, 10).join('<br>')
})

function updateHighlighting() {
  nextTick(() => {
    syncScroll()
  })
}

function syncScroll() {
  if (customExpressionRef.value && highlightRef.value) {
    highlightRef.value.scrollTop = customExpressionRef.value.scrollTop
    highlightRef.value.scrollLeft = customExpressionRef.value.scrollLeft
  }
}

// Apply preset configuration
function applyPreset(preset: RulePreset) {
  editingRule.name = preset.name
  editingRule.description = preset.description
  editingRule.type = preset.type
  editingRule.action = preset.action
  editingRule.score = preset.score
  editingRule.conditions = { ...preset.conditions }
}

onMounted(() => {
  if (authStore.user) {
    loadRules()
    loadSchemas()
  } else {
    router.push('/login')
  }
})

async function loadSchemas() {
  try {
    const response = await api.get<{ schemas: EventSchema[] }>('/api/v1/schemas')
    schemas.value = response?.schemas || []
  } catch (e) {
    console.error('Error loading schemas:', e)
  }
}

async function loadRules() {
  try {
    loading.value = true
    error.value = ''
    const response = await api.get<{ rules: any[] }>('/api/v1/rules')
    // Sort by priority ascending (0 = highest priority, 1000 = lowest)
    const loadedRules = response?.rules || []
    rules.value = loadedRules.sort((a, b) => (a.priority ?? 1000) - (b.priority ?? 1000))
  } catch (e: any) {
    error.value = e.message || 'Failed to load rules'
    console.error('Error loading rules:', e)
  } finally {
    loading.value = false
  }
}

function openCreateModal() {
  isEditing.value = false
  delete editingRule.id
  editingRule.name = ''
  editingRule.description = ''
  editingRule.type = 'custom'
  editingRule.action = 'review'
  editingRule.score = 0
  editingRule.priority = 500
  editingRule.enabled = true
  editingRule.conditions = getDefaultConditions('custom')
  editingRule.schema_id = schemas.value.length > 0 ? schemas.value[0]?.id : undefined
  showAllFields.value = false
  
  // Initialize with one empty condition row
  conditionRows.value = [{
    id: `condition-${++conditionRowIdCounter}`,
    field: '',
    operator: '',
    value: null,
    logicOperator: 'and',
  }]
  
  showModal.value = true
}

function openEditModal(rule: any) {
  isEditing.value = true
  editingRule.id = rule.id
  editingRule.name = rule.name
  editingRule.description = rule.description
  editingRule.type = 'custom' // Always custom in schema-based system
  editingRule.action = rule.action
  editingRule.score = rule.score
  editingRule.priority = rule.priority
  editingRule.enabled = rule.enabled
  editingRule.schema_id = rule.schema_id
  
  // Preserve existing conditions, but ensure custom_expression exists
  if (rule.conditions && rule.conditions.custom_expression) {
    editingRule.conditions = { custom_expression: rule.conditions.custom_expression }
    // Try to parse existing expression into condition rows
    // For now, start with empty rows - user can rebuild or we can add parsing later
    conditionRows.value = [{
      id: `condition-${++conditionRowIdCounter}`,
      field: '',
      operator: '',
      value: null,
      logicOperator: 'and',
    }]
  } else {
    editingRule.conditions = getDefaultConditions('custom')
    conditionRows.value = [{
      id: `condition-${++conditionRowIdCounter}`,
      field: '',
      operator: '',
      value: null,
      logicOperator: 'and',
    }]
  }
  
  showAllFields.value = false
  showModal.value = true
}

function closeModal() {
  showModal.value = false
}

async function handleSubmit() {
  // Validate form before submission
  if (formRef.value) {
    const { valid } = await formRef.value.validate()
    if (!valid) return
  }

  // Additional validation for conditions
  const conditionError = validateConditions()
  if (conditionError) {
    error.value = conditionError
    return
  }

  // Build expression from condition rows
  editingRule.conditions.custom_expression = generatedExpression.value

  try {
    saving.value = true
    error.value = ''
    if (isEditing.value && editingRule.id) {
      // Update existing rule
      await api.put(`/api/v1/rules/${editingRule.id}`, editingRule)
    } else {
      // Create new rule
      await api.post('/api/v1/rules', editingRule)
    }
    closeModal()
    await loadRules()
  } catch (e: any) {
    error.value = e.message || 'Failed to save rule'
  } finally {
    saving.value = false
  }
}

function validateConditions(): string | null {
  // Validate that a schema is selected
  if (!editingRule.schema_id || !selectedSchema.value) {
    return 'Please select a schema first'
  }

  const schemaFields = selectedSchema.value.extracted_fields || []
  
  if (schemaFields.length === 0) {
    return 'Selected schema has no fields. Please select a different schema.'
  }

  // Validate condition rows
  if (conditionRows.value.length === 0) {
    return 'At least one condition is required'
  }

  // Validate each condition row
  for (let i = 0; i < conditionRows.value.length; i++) {
    const row = conditionRows.value[i]
    
    if (!row) {
      continue
    }
    
    if (!row.field) {
      return `Condition ${i + 1}: Field is required`
    }
    
    if (!row.operator) {
      return `Condition ${i + 1}: Operator is required`
    }
    
    if (row.value === null || row.value === undefined || row.value === '') {
      return `Condition ${i + 1}: Value is required`
    }

    // Validate that the field exists in the schema
    const fieldExists = schemaFields.some(f => f.path === row.field)
    if (!fieldExists) {
      return `Condition ${i + 1}: Field "${row.field}" does not exist in the selected schema`
    }

    // Validate operator is appropriate for field type
    const field = schemaFields.find(f => f.path === row.field)
    if (field) {
      const validOperators = getOperatorOptions(row.field).map(o => o.value)
      if (!validOperators.includes(row.operator)) {
        return `Condition ${i + 1}: Operator "${row.operator}" is not valid for field type "${field.type}"`
      }
    }
  }

  // Validate generated expression
  const expression = generatedExpression.value
  if (!expression || expression.trim() === '') {
    return 'Generated expression is empty. Please check your conditions.'
  }
  
  // Sanitize the expression
  const sanitizeResult = sanitizeExpression(expression)
  if (!sanitizeResult.valid) {
    return sanitizeResult.error || 'Invalid expression'
  }

  return null
}

// Extract all field references from an expression and validate they exist in the schema
function validateAllFieldsExist(expression: string, validFieldPaths: string[]): string[] {
  // Keywords and operators to exclude from field validation
  const keywords = new Set([
    'and', 'or', 'not', 'in', 'true', 'false', 'null',
    'contains', 'len', 'abs', 'min', 'max', 'sum', 'avg'
  ])

  // First, identify string literals to exclude them from field extraction
  const stringRanges: Array<{ start: number; end: number }> = []
  let inString = false
  let stringStart = -1
  let i = 0
  
  while (i < expression.length) {
    if (expression[i] === '"' && (i === 0 || expression[i - 1] !== '\\')) {
      if (!inString) {
        stringStart = i
        inString = true
      } else {
        stringRanges.push({ start: stringStart, end: i })
        inString = false
      }
    }
    i++
  }

  // Helper to check if a position is inside a string literal
  const isInStringLiteral = (pos: number): boolean => {
    return stringRanges.some(range => pos >= range.start && pos <= range.end)
  }

  // Extract all potential field references from the expression
  // This regex matches identifiers (field names) and nested paths (field.path)
  // Pattern: word characters, dots, and underscores (for field paths)
  const fieldPattern = /\b([a-zA-Z_][a-zA-Z0-9_]*(?:\.[a-zA-Z_][a-zA-Z0-9_]*)*)\b/g
  
  const foundFields = new Set<string>()
  let match
  
  while ((match = fieldPattern.exec(expression)) !== null) {
    const potentialField = match[1]
    
    if (!potentialField) {
      continue
    }
    
    const matchStart = match.index ?? 0
    
    // Skip if it's a keyword
    if (keywords.has(potentialField.toLowerCase())) {
      continue
    }
    
    // Skip if it's a number (starts with digit)
    if (/^\d/.test(potentialField)) {
      continue
    }
    
    // Skip if it's inside a string literal
    if (isInStringLiteral(matchStart)) {
      continue
    }
    
    foundFields.add(potentialField)
  }

  // Check which found fields don't exist in the schema
  // We need to check exact matches - a field path must exactly match a schema field
  const validFieldPathsSet = new Set(validFieldPaths)
  const invalidFields: string[] = []
  
  for (const field of foundFields) {
    // Check for exact match
    if (!validFieldPathsSet.has(field)) {
      invalidFields.push(field)
    }
  }

  return invalidFields
}

async function deleteRule(id: string) {
  if (!confirm('Are you sure you want to delete this rule?')) return

  try {
    await api.delete(`/api/v1/rules/${id}`)
    await loadRules()
  } catch (e: any) {
    error.value = e.message || 'Failed to delete rule'
  }
}

async function toggleRule(rule: any) {
  try {
    await api.put(`/api/v1/rules/${rule.id}`, { ...rule, enabled: !rule.enabled })
    await loadRules()
  } catch (e: any) {
    error.value = e.message || 'Failed to toggle rule'
  }
}

// Badge variants for rule actions (matches backend RuleAction enum in src/pkg/models/rule.go)
function getActionBadgeVariant(action: string): 'success' | 'warning' | 'danger' | 'info' | 'default' {
  switch (action.toLowerCase()) {
    case 'allow':
      return 'success'
    case 'block':
      return 'danger'
    case 'review':
      return 'warning'
    case 'score':
      return 'info'
    default:
      return 'default'
  }
}
</script>

<style scoped>
.cursor-pointer {
  cursor: pointer;
}

.custom-expression-container {
  position: relative;
}

.code-editor-wrapper {
  position: relative;
  font-family: 'Fira Code', 'Monaco', 'Consolas', monospace;
  font-size: 14px;
  line-height: 1.5;
}

.code-editor {
  width: 100%;
  min-height: 100px;
  padding: 12px;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  font-family: inherit;
  font-size: inherit;
  line-height: inherit;
  resize: vertical;
  background: transparent;
  color: transparent;
  caret-color: #333;
  position: relative;
  z-index: 1;
}

.code-editor:focus {
  outline: none;
  border-color: #1976d2;
  box-shadow: 0 0 0 2px rgba(25, 118, 210, 0.2);
}

.code-highlight {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  padding: 12px;
  margin: 0;
  font-family: inherit;
  font-size: inherit;
  line-height: inherit;
  white-space: pre-wrap;
  word-wrap: break-word;
  overflow: hidden;
  pointer-events: none;
  border: 1px solid transparent;
  border-radius: 8px;
  background: #f5f5f5;
  z-index: 0;
}

.code-highlight code {
  font-family: inherit;
  background: transparent;
}
</style>
