<template>
  <v-container fluid class="pa-8">
    <div class="d-flex justify-space-between align-center mb-10">
      <div>
        <div class="d-flex align-center gap-3 mb-2">
          <v-icon icon="fa-tasks" size="large" color="primary" />
          <h2 class="text-h4 font-weight-bold">{{ $t('views.rules.title') }}</h2>
        </div>
        <p class="text-body-1 text-grey-darken-1">{{ $t('views.rules.subtitle') }}</p>
      </div>
      <BaseButton @click="openCreateModal" prepend-icon="fa-plus">
        {{ $t('views.rules.createRule') }}
      </BaseButton>
    </div>

    <LoadingSpinner v-if="loading" :text="$t('views.rules.loading')" :centered="false" />

    <ErrorMessage
      v-else-if="error"
      :title="$t('views.rules.errorTitle')"
      :message="error"
      retryable
      @retry="loadRules"
    />

    <div v-else>
      <BaseTable
        :columns="tableColumns"
        :data="paginatedRules"
        :empty-text="$t('views.rules.emptyText')"
      >
      <template #cell-schema="{ row }">
        <span class="text-body-2 font-weight-medium text-grey-darken-2">
          {{ getSchemaName(row.schema_id) || 'N/A' }}
        </span>
      </template>

      <template #cell-name="{ row }">
        <div class="font-weight-semibold text-grey-darken-3">{{ row.name }}</div>
        <div class="text-body-2 text-grey-darken-1">{{ row.description }}</div>
      </template>

      <template #cell-action="{ row }">
        <BaseBadge :variant="getActionBadgeVariant(row.action)" rounded>
          {{ getActionLabel(row.action) }}
        </BaseBadge>
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
          {{ row.enabled ? $t('components.ruleTable.enabled') : $t('components.ruleTable.disabled') }}
        </BaseBadge>
      </template>

      <template #cell-actions="{ row }">
        <div class="d-flex gap-2">
          <BaseButton size="sm" @click="openEditModal(row)" prepend-icon="fa-pencil">
            {{ $t('components.ruleTable.edit') }}
          </BaseButton>
          <BaseButton size="sm" variant="danger" @click="deleteRule(row.id)" prepend-icon="fa-trash">
            {{ $t('components.ruleTable.delete') }}
          </BaseButton>
        </div>
      </template>
      </BaseTable>

      <!-- Pagination -->
      <div v-if="rules.length > pageSize" class="mt-4 d-flex justify-space-between align-center pa-4 bg-grey-lighten-5 rounded-lg">
        <div class="text-body-2 text-grey-darken-1">
          {{ $t('components.ruleTable.showing') }} {{ startIndex + 1 }} {{ $t('components.ruleTable.to') }} {{ endIndex }} {{ $t('components.ruleTable.of') }} {{ rules.length }} {{ $t('components.ruleTable.rules') }}
        </div>
        <div class="d-flex align-center gap-2">
          <v-btn
            :disabled="currentPage === 1"
            @click="prevPage"
            icon
            variant="text"
            size="small"
          >
            <v-icon icon="fa-chevron-left" />
          </v-btn>

          <div class="d-flex align-center gap-1">
            <v-btn
              v-for="page in visiblePages"
              :key="page"
              @click="goToPage(page)"
              :variant="currentPage === page ? 'flat' : 'text'"
              :color="currentPage === page ? 'primary' : undefined"
              size="small"
              class="min-width-40"
            >
              {{ page }}
            </v-btn>
          </div>

          <v-btn
            :disabled="currentPage === totalPages"
            @click="nextPage"
            icon
            variant="text"
            size="small"
          >
            <v-icon icon="fa-chevron-right" />
          </v-btn>
        </div>
      </div>
    </div>

    <RuleFormModal
      v-model="showModal"
      :is-editing="isEditing"
      :editing-rule="editingRule"
      :schema-options="schemaOptions"
      :current-schema="currentSchema"
      :rule-presets="rulePresets"
      :rule-actions="ruleActions"
      :saving="saving"
      :expression-mode="expressionMode"
      :builder-type="builderType"
      :polygon-config="polygonConfig"
      :velocity-config="velocityConfig"
      :condition-rows="conditionRows"
      :polygon-expression="polygonExpression"
      :velocity-expression="velocityExpression"
      :generated-expression="generatedExpression"
      @apply-preset="applyPreset"
      @submit="handleSubmit"
      @cancel="closeModal"
    >
      <template #expression-builder>
        <!-- Rule Expression Section -->
        <div v-if="currentSchema" class="mb-6 pa-4 bg-grey-lighten-5 rounded-lg">
          <div class="d-flex justify-space-between align-center mb-4">
            <h4 class="text-subtitle-1 font-weight-medium d-flex align-center">
              <v-icon icon="fa-code" size="small" class="mr-2" />
              {{ $t('views.rules.modal.ruleConditions') }} <span class="text-red">*</span>
            </h4>
            <div class="d-flex gap-2">
              <v-btn-toggle
                v-model="expressionMode"
                mandatory
                color="primary"
                variant="outlined"
                density="compact"
              >
              <v-btn value="manual" size="small">
                  <v-icon icon="fa-code" size="small" class="mr-1" />
                  {{ $t('views.rules.modal.manual') }}
                </v-btn>
                <v-btn value="builder" size="small">
                  <v-icon icon="fa-wrench" size="small" class="mr-1" />
                  {{ $t('views.rules.modal.builder') }}
                </v-btn>
              </v-btn-toggle>
              <BaseButton v-if="expressionMode === 'manual'" size="sm" @click="addCondition" prepend-icon="fa-plus">
                {{ $t('views.rules.modal.addCondition') }}
              </BaseButton>
            </div>
          </div>

          <!-- Builder Mode: Polygon Builder -->
          <div v-if="expressionMode === 'builder' && builderType === 'polygon'" class="mb-4">
            <div class="d-flex justify-space-between align-center mb-3">
              <h5 class="text-body-1 font-weight-medium">{{ $t('views.rules.modal.polygonBuilder') }}</h5>
              <v-btn-toggle
                v-model="builderType"
                mandatory
                color="primary"
                variant="outlined"
                density="compact"
              >
                <v-btn value="polygon" size="small">{{ $t('views.rules.modal.polygon') }}</v-btn>
                <v-btn value="velocity" size="small">{{ $t('views.rules.modal.velocity') }}</v-btn>
              </v-btn-toggle>
            </div>
            <PolygonBuilder
              :config="polygonConfig"
              :field-options="fieldOptions"
              :expression="polygonExpression"
              @add-point="addPolygonPoint"
              @remove-point="removePolygonPoint"
            />
          </div>

          <!-- Builder Mode: Velocity Builder -->
          <div v-if="expressionMode === 'builder' && builderType === 'velocity'" class="mb-4">
            <div class="d-flex justify-space-between align-center mb-3">
              <h5 class="text-body-1 font-weight-medium">{{ $t('views.rules.modal.velocityBuilder') }}</h5>
              <v-btn-toggle
                v-model="builderType"
                mandatory
                color="primary"
                variant="outlined"
                density="compact"
              >
                <v-btn value="polygon" size="small">{{ $t('views.rules.modal.polygon') }}</v-btn>
                <v-btn value="velocity" size="small">{{ $t('views.rules.modal.velocity') }}</v-btn>
              </v-btn-toggle>
            </div>
            <VelocityBuilder
              :config="velocityConfig"
              :field-options="fieldOptions"
              :expression="velocityExpression"
            />
          </div>

          <!-- Manual Mode: Condition Builder -->
          <ConditionBuilder
            v-if="expressionMode === 'manual'"
            :condition-rows="conditionRows"
            :field-options="fieldOptions"
            :current-schema="currentSchema"
            :generated-expression="generatedExpression"
            @remove-condition="removeCondition"
          />
        </div>
      </template>
    </RuleFormModal>
  </v-container>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { api } from '@/lib/api'
import { i18n } from '@/plugins/i18n'
import BaseButton from '@/components/BaseButton.vue'
import BaseBadge from '@/components/BaseBadge.vue'
import BaseTable from '@/components/BaseTable.vue'
import LoadingSpinner from '@/components/LoadingSpinner.vue'
import ErrorMessage from '@/components/ErrorMessage.vue'
import RuleFormModal from '@/components/RuleFormModal.vue'
import PolygonBuilder from '@/components/PolygonBuilder.vue'
import VelocityBuilder from '@/components/VelocityBuilder.vue'
import ConditionBuilder from '@/components/ConditionBuilder.vue'

const router = useRouter()
const authStore = useAuthStore()
const t = i18n.global.t

const tableColumns = [
  { key: 'schema', label: 'components.ruleTable.schema' },
  { key: 'name', label: 'components.ruleTable.name' },
  { key: 'action', label: 'components.ruleTable.action' },
  { key: 'priority', label: 'components.ruleTable.priority' },
  { key: 'status', label: 'components.ruleTable.status' },
  { key: 'actions', label: 'components.ruleTable.actions' },
]

const rules = ref<any[]>([])
const loading = ref(true)
const error = ref('')
const showModal = ref(false)
const isEditing = ref(false)
const saving = ref(false)

// Pagination
const currentPage = ref(1)
const pageSize = ref(10)

const totalPages = computed(() => Math.ceil(rules.value.length / pageSize.value))
const startIndex = computed(() => (currentPage.value - 1) * pageSize.value)
const endIndex = computed(() => Math.min(startIndex.value + pageSize.value, rules.value.length))
const paginatedRules = computed(() => {
  return rules.value.slice(startIndex.value, endIndex.value)
})
const visiblePages = computed(() => {
  const pages: number[] = []
  const maxVisible = 5
  let start = Math.max(1, currentPage.value - Math.floor(maxVisible / 2))
  let end = Math.min(totalPages.value, start + maxVisible - 1)

  if (end - start < maxVisible - 1) {
    start = Math.max(1, end - maxVisible + 1)
  }

  for (let i = start; i <= end; i++) {
    pages.push(i)
  }

  return pages
})

function prevPage() {
  if (currentPage.value > 1) {
    currentPage.value--
  }
}

function nextPage() {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
  }
}

function goToPage(page: number) {
  currentPage.value = page
}

interface RuleForm {
  id?: string
  name: string
  description: string
  action: string
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
  action: '',
  priority: 50,
  enabled: true,
  conditions: {},
  schema_id: undefined,
})

// Structured condition rows for the form
const conditionRows = ref<ConditionRow[]>([])
let conditionRowIdCounter = 0

// Expression mode: 'builder' or 'manual'
const expressionMode = ref<'builder' | 'manual'>('builder')
const builderType = ref<'polygon' | 'velocity'>('polygon')

// Polygon builder configuration
const polygonConfig = reactive({
  latField: '',
  lonField: '',
  points: [[0, 0], [0, 0], [0, 0]] as Array<[number, number]>
})

// Velocity builder configuration
const velocityConfig = reactive<{
  metric: 'count' | 'sum'
  groupField: string
  threshold: number
  timeValue: number
  timeUnit: 'seconds' | 'minutes' | 'hours' | 'days'
}>({
  metric: 'count',
  groupField: '',
  threshold: 10,
  timeValue: 1,
  timeUnit: 'hours'
})

// Schemas for selection
const schemas = ref<EventSchema[]>([])
const selectedSchema = ref<EventSchema | null>(null)

// Computed property that resolves schema directly from editingRule.schema_id
// This avoids race conditions with the async watcher that updates selectedSchema
const currentSchema = computed(() => {
  if (!editingRule.schema_id) {
    return null
  }
  return schemas.value.find(s => s.id === editingRule.schema_id) || null
})

// Computed list of schema options for the dropdown
const schemaOptions = computed(() => {
  return schemas.value.map(s => ({
    value: s.id,
    label: s.name,
  }))
})

// Computed list of field options from current schema (resolved directly from editingRule.schema_id)
const fieldOptions = computed(() => {
  if (!currentSchema.value || !currentSchema.value.extracted_fields) {
    return []
  }
  return currentSchema.value.extracted_fields.map(f => ({
    value: f.path,
    label: `${f.path} (${f.type})`,
  }))
})

// Get operator options based on field type
function getOperatorOptions(fieldPath: string): Array<{ value: string; label: string }> {
  if (!currentSchema.value || !fieldPath) {
    return []
  }
  
  const field = currentSchema.value.extracted_fields.find(f => f.path === fieldPath)
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
  if (!currentSchema.value || !fieldPath) {
    return 'string'
  }
  
  const field = currentSchema.value.extracted_fields.find(f => f.path === fieldPath)
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
  if (!currentSchema.value || !fieldPath) {
    return 'Enter value'
  }
  
  const field = currentSchema.value.extracted_fields.find(f => f.path === fieldPath)
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

  // Resolve schema directly from editingRule.schema_id to avoid async watcher timing issues
  const currentSchema = editingRule.schema_id 
    ? schemas.value.find(s => s.id === editingRule.schema_id) 
    : null

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
    const field = currentSchema?.extracted_fields?.find(f => f.path === row.field)
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

// Polygon expression generator
const polygonExpression = computed(() => {
  if (!polygonConfig.latField || !polygonConfig.lonField || polygonConfig.points.length < 3) {
    return ''
  }
  const pointsStr = polygonConfig.points
    .filter(p => p[0] !== 0 || p[1] !== 0) // Filter out empty points
    .map(p => `[${p[0]}, ${p[1]}]`)
    .join(', ')
  return `pointInPolygon(${polygonConfig.latField}, ${polygonConfig.lonField}, [${pointsStr}])`
})

// Velocity expression generator
const velocityExpression = computed(() => {
  if (!velocityConfig.groupField || !velocityConfig.threshold || !velocityConfig.timeValue) {
    return ''
  }
  const timeWindowSeconds = convertTimeToSeconds(velocityConfig.timeValue, velocityConfig.timeUnit)
  const functionName = velocityConfig.metric === 'count' ? 'velocityCount' : 'velocitySum'
  return `${functionName}(${velocityConfig.groupField}, ${timeWindowSeconds}) > ${velocityConfig.threshold}`
})

// Convert time value and unit to seconds
function convertTimeToSeconds(value: number, unit: 'seconds' | 'minutes' | 'hours' | 'days'): number {
  switch (unit) {
    case 'seconds':
      return value
    case 'minutes':
      return value * 60
    case 'hours':
      return value * 3600
    case 'days':
      return value * 86400
    default:
      return value
  }
}

// Polygon builder functions
function addPolygonPoint() {
  polygonConfig.points.push([0, 0])
}

function removePolygonPoint(index: number) {
  if (polygonConfig.points.length > 3) {
    polygonConfig.points.splice(index, 1)
  }
}

// Rule actions must match backend validation: allow|block|review
const ruleActions = [
  { value: 'allow', label: 'Allow' },
  { value: 'block', label: 'Block' },
  { value: 'review', label: 'Review' },
]

// Preset configurations for common fraud detection scenarios
interface RulePreset {
  id: string
  label: string
  icon: string
  name: string
  description: string
  action: string
  conditions: Record<string, any>
}

const rulePresets: RulePreset[] = [
  {
    id: 'high-value',
    label: 'High Value',
    icon: 'fa-dollar-sign',
    name: 'High Value Transaction',
    description: 'Flag events with high value amounts',
    action: 'review',
    conditions: { custom_expression: 'amount > 10000' },
  },
  {
    id: 'suspicious-flag',
    label: 'Suspicious Flag',
    icon: 'fa-exclamation-triangle',
    name: 'Suspicious Flag Check',
    description: 'Check for suspicious flag in event metadata',
    action: 'review',
    conditions: { custom_expression: 'metadata.is_suspicious == true' },
  },
  {
    id: 'polygon-restriction',
    label: 'Polygon Restriction',
    icon: 'fa-globe',
    name: 'Polygon Restriction',
    description: 'Flag events within a geographic polygon area',
    action: 'review',
    conditions: { custom_expression: 'pointInPolygon(location.lat, location.lon, [[37.7749, -122.4194], [37.7849, -122.4094], [37.7649, -122.4294]])' },
  },
  {
    id: 'velocity-check',
    label: 'High Transaction Frequency',
    icon: 'fa-tachometer-alt',
    name: 'High Transaction Frequency',
    description: 'Flag events with high transaction frequency',
    action: 'review',
    conditions: { custom_expression: 'velocityCount(field, 3600) > 10' },
  },
  {
    id: 'amount-velocity',
    label: 'Amount Velocity',
    icon: 'fa-chart-line',
    name: 'High Amount Velocity',
    description: 'Flag fields with high cumulative amounts in time window',
    action: 'review',
    conditions: { custom_expression: 'velocitySum(field, 3600) > 10000' },
  },
]


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
function getDefaultConditions(): Record<string, any> {
  return { custom_expression: '' }
}

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


// Parse velocity expression: velocityCount(field, timeWindow) > threshold or velocitySum(field, timeWindow) > threshold
function parseVelocityExpression(expression: string): { metric: 'count' | 'sum', groupField: string, timeWindowSeconds: number, threshold: number } | null {
  const velocityCountMatch = expression.match(/velocityCount\s*\(\s*([^,]+)\s*,\s*(\d+)\s*\)\s*>\s*(\d+)/)
  if (velocityCountMatch && velocityCountMatch[1] && velocityCountMatch[2] && velocityCountMatch[3]) {
    return {
      metric: 'count',
      groupField: velocityCountMatch[1].trim(),
      timeWindowSeconds: parseInt(velocityCountMatch[2], 10),
      threshold: parseInt(velocityCountMatch[3], 10),
    }
  }
  
  const velocitySumMatch = expression.match(/velocitySum\s*\(\s*([^,]+)\s*,\s*(\d+)\s*\)\s*>\s*(\d+)/)
  if (velocitySumMatch && velocitySumMatch[1] && velocitySumMatch[2] && velocitySumMatch[3]) {
    return {
      metric: 'sum',
      groupField: velocitySumMatch[1].trim(),
      timeWindowSeconds: parseInt(velocitySumMatch[2], 10),
      threshold: parseInt(velocitySumMatch[3], 10),
    }
  }
  
  return null
}

// Parse manual expression into condition rows
function parseManualExpression(expression: string): ConditionRow[] {
  const rows: ConditionRow[] = []
  
  if (!expression || !expression.trim()) {
    return [{
      id: `condition-${++conditionRowIdCounter}`,
      field: '',
      operator: '',
      value: null,
      logicOperator: 'and',
    }]
  }
  
  // Helper function to parse a single condition
  function parseCondition(conditionStr: string, logicOp: 'and' | 'or' = 'and'): ConditionRow | null {
    const part = conditionStr.trim()
    if (!part) return null
    
    // Pattern: field operator value
    // Examples: amount > 10000, currency == "USD", metadata.is_suspicious == true
    const simpleMatch = part.match(/^([a-zA-Z_][a-zA-Z0-9_.]*)\s*(==|!=|>|<|>=|<=)\s*(.+)$/)
    if (simpleMatch && simpleMatch[1] && simpleMatch[2] && simpleMatch[3]) {
      const field = simpleMatch[1].trim()
      const operator = simpleMatch[2].trim()
      let value = simpleMatch[3].trim()
      
      // Remove quotes from string values
      if ((value.startsWith('"') && value.endsWith('"')) || (value.startsWith("'") && value.endsWith("'"))) {
        value = value.slice(1, -1)
      }
      
      // Handle boolean values - keep as string for the form
      // The form will handle the conversion when needed
      
      return {
        id: `condition-${++conditionRowIdCounter}`,
        field: field,
        operator: operator,
        value: value,
        logicOperator: logicOp,
      }
    }
    
    // Pattern: field in [value1, value2, ...]
    // Example: origin in ["BLOCKED001", "BLOCKED002"]
    const inMatch = part.match(/^([a-zA-Z_][a-zA-Z0-9_.]*)\s+in\s+\[(.+)\]$/)
    if (inMatch && inMatch[1] && inMatch[2]) {
      const field = inMatch[1].trim()
      const valuesStr = inMatch[2].trim()
      
      // Parse array values
      const values = valuesStr.split(',').map(v => {
        let val = v.trim()
        // Remove quotes
        if ((val.startsWith('"') && val.endsWith('"')) || (val.startsWith("'") && val.endsWith("'"))) {
          val = val.slice(1, -1)
        }
        return val
      }).filter(v => v)
      
      return {
        id: `condition-${++conditionRowIdCounter}`,
        field: field,
        operator: 'in',
        value: values.join(', '),
        logicOperator: logicOp,
      }
    }
    
    // Pattern: "value" in field (for array contains)
    // Example: "high-value" in tags
    const containsMatch = part.match(/^"([^"]+)"\s+in\s+([a-zA-Z_][a-zA-Z0-9_.]*)$/)
    if (containsMatch && containsMatch[1] && containsMatch[2]) {
      const value = containsMatch[1]
      const field = containsMatch[2].trim()
      
      return {
        id: `condition-${++conditionRowIdCounter}`,
        field: field,
        operator: 'contains',
        value: value,
        logicOperator: logicOp,
      }
    }
    
    // Pattern: not condition
    // Example: not is_verified
    const notMatch = part.match(/^not\s+(.+)$/i)
    if (notMatch && notMatch[1]) {
      const innerCondition = parseCondition(notMatch[1].trim(), logicOp)
      if (innerCondition) {
        // For 'not', we'll use != operator with boolean false or handle it differently
        // This is a simplified approach
        if (innerCondition.operator === '==') {
          innerCondition.operator = '!='
        }
        return innerCondition
      }
    }
    
    return null
  }
  
  // Split expression by 'and' and 'or', handling parentheses
  // This is a simplified parser - for very complex nested expressions, it may not parse perfectly
  let expr = expression.trim()
  
  // Handle expressions with parentheses by splitting on top-level operators
  // We'll use a simple approach: split on ' and ' and ' or ' that are not inside parentheses
  const tokens: Array<{ text: string, operator?: 'and' | 'or' }> = []
  let current = ''
  let depth = 0
  
  for (let i = 0; i < expr.length; i++) {
    const char = expr[i]
    
    if (char === '(') {
      depth++
      current += char
    } else if (char === ')') {
      depth--
      current += char
    } else if (depth === 0 && i < expr.length - 1) {
      // Check for ' and ' or ' or ' at top level
      const remaining = expr.substring(i)
      if (remaining.toLowerCase().startsWith(' and ')) {
        if (current.trim()) {
          tokens.push({ text: current.trim() })
        }
        tokens.push({ text: 'and', operator: 'and' })
        current = ''
        i += 4 // Skip ' and '
        continue
      } else if (remaining.toLowerCase().startsWith(' or ')) {
        if (current.trim()) {
          tokens.push({ text: current.trim() })
        }
        tokens.push({ text: 'or', operator: 'or' })
        current = ''
        i += 3 // Skip ' or '
        continue
      }
      current += char
    } else {
      current += char
    }
  }
  
  if (current.trim()) {
    tokens.push({ text: current.trim() })
  }
  
  // Process tokens
  let currentLogicOp: 'and' | 'or' = 'and'
  
  for (const token of tokens) {
    if (token.operator) {
      currentLogicOp = token.operator
      continue
    }
    
    // Remove outer parentheses if present
    let conditionText = token.text.trim()
    while (conditionText.startsWith('(') && conditionText.endsWith(')')) {
      conditionText = conditionText.slice(1, -1).trim()
    }
    
    const condition = parseCondition(conditionText, currentLogicOp)
    if (condition) {
      rows.push(condition)
    } else {
      // If we can't parse it, create a row with the raw expression
      rows.push({
        id: `condition-${++conditionRowIdCounter}`,
        field: '',
        operator: '',
        value: conditionText,
        logicOperator: currentLogicOp,
      })
    }
  }
  
  // Set logic operator for first row (default to 'and')
  if (rows.length > 0 && rows[0]) {
    rows[0].logicOperator = 'and'
  }
  
  return rows.length > 0 ? rows : [{
    id: `condition-${++conditionRowIdCounter}`,
    field: '',
    operator: '',
    value: null,
    logicOperator: 'and',
  }]
}

// Parse polygon expression: pointInPolygon(latField, lonField, [[lat1, lon1], [lat2, lon2], ...])
function parsePolygonExpression(expression: string): { latField: string, lonField: string, points: Array<[number, number]> } | null {
  const match = expression.match(/pointInPolygon\s*\(\s*([^,]+)\s*,\s*([^,]+)\s*,\s*\[(.*?)\]\s*\)/)
  if (!match || !match[1] || !match[2] || !match[3]) return null
  
  const latField = match[1].trim()
  const lonField = match[2].trim()
  const pointsStr = match[3]
  
  if (!pointsStr) return null
  
  // Parse points array: [[lat1, lon1], [lat2, lon2], ...]
  const pointMatches = Array.from(pointsStr.matchAll(/\[\s*([-\d.]+)\s*,\s*([-\d.]+)\s*\]/g))
  const points: Array<[number, number]> = []
  
  for (const pointMatch of pointMatches) {
    if (pointMatch[1] && pointMatch[2]) {
      const lat = parseFloat(pointMatch[1])
      const lon = parseFloat(pointMatch[2])
      if (!isNaN(lat) && !isNaN(lon)) {
        points.push([lat, lon])
      }
    }
  }
  
  if (points.length >= 3) {
    return { latField, lonField, points }
  }
  
  return null
}

// Convert seconds to time value and unit
function convertSecondsToTimeValue(seconds: number): { value: number, unit: 'seconds' | 'minutes' | 'hours' | 'days' } {
  if (seconds % 86400 === 0) {
    return { value: seconds / 86400, unit: 'days' }
  } else if (seconds % 3600 === 0) {
    return { value: seconds / 3600, unit: 'hours' }
  } else if (seconds % 60 === 0) {
    return { value: seconds / 60, unit: 'minutes' }
  } else {
    return { value: seconds, unit: 'seconds' }
  }
}

// Apply preset configuration
function applyPreset(preset: RulePreset) {
  // Override all current field values
  editingRule.name = preset.name
  editingRule.description = preset.description
  editingRule.action = preset.action
  
  const expression = preset.conditions.custom_expression || ''
  
  // Try to parse as velocity expression
  const velocityParsed = parseVelocityExpression(expression)
  if (velocityParsed) {
    expressionMode.value = 'builder'
    builderType.value = 'velocity'
    velocityConfig.metric = velocityParsed.metric
    velocityConfig.groupField = velocityParsed.groupField
    const timeConfig = convertSecondsToTimeValue(velocityParsed.timeWindowSeconds)
    velocityConfig.timeValue = timeConfig.value
    velocityConfig.timeUnit = timeConfig.unit
    velocityConfig.threshold = velocityParsed.threshold
    // Clear condition rows since we're using builder
    conditionRows.value = []
    return
  }
  
  // Try to parse as polygon expression
  const polygonConfigParsed = parsePolygonExpression(expression)
  if (polygonConfigParsed) {
    expressionMode.value = 'builder'
    builderType.value = 'polygon'
    
    // Find lat/lon fields in schema by substring match
    if (currentSchema.value && currentSchema.value.extracted_fields) {
      const fields = currentSchema.value.extracted_fields
      const latFieldLower = polygonConfigParsed.latField.toLowerCase()
      const lonFieldLower = polygonConfigParsed.lonField.toLowerCase()
      
      // Find latitude field (check for 'lat' substring)
      const latFieldMatch = fields.find(f => 
        f.path.toLowerCase().includes('lat') || latFieldLower.includes(f.path.toLowerCase())
      )
      if (latFieldMatch) {
        polygonConfig.latField = latFieldMatch.path
      } else {
        polygonConfig.latField = polygonConfigParsed.latField
      }
      
      // Find longitude field (check for 'lon' or 'lng' substring)
      const lonFieldMatch = fields.find(f => 
        f.path.toLowerCase().includes('lon') || 
        f.path.toLowerCase().includes('lng') || 
        lonFieldLower.includes(f.path.toLowerCase())
      )
      if (lonFieldMatch) {
        polygonConfig.lonField = lonFieldMatch.path
      } else {
        polygonConfig.lonField = polygonConfigParsed.lonField
      }
    } else {
      // Fallback to parsed values if no schema
      polygonConfig.latField = polygonConfigParsed.latField
      polygonConfig.lonField = polygonConfigParsed.lonField
    }
    
    // Ensure we have at least 3 points, pad with [0,0] if needed
    while (polygonConfigParsed.points.length < 3) {
      polygonConfigParsed.points.push([0, 0])
    }
    polygonConfig.points = polygonConfigParsed.points
    // Clear condition rows since we're using builder
    conditionRows.value = []
    return
  }
  
  // Default to manual mode for other expressions
  expressionMode.value = 'manual'
  editingRule.conditions = { ...preset.conditions }
  // For manual mode, we just set the expression directly
  // The user can see it in the generated expression preview
  // Clear condition rows since we're using the expression directly
  conditionRows.value = [{
    id: `condition-${++conditionRowIdCounter}`,
    field: '',
    operator: '',
    value: null,
    logicOperator: 'and',
  }]
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
    // Sort by priority ascending (0 = highest priority, 100 = lowest)
    const loadedRules = response?.rules || []
    rules.value = loadedRules.sort((a, b) => (a.priority ?? 100) - (b.priority ?? 100))
    // Reset to first page when rules are reloaded
    currentPage.value = 1
  } catch (e: any) {
    error.value = e.message || (window as any).$i18n?.global?.t?.('views.rules.errorLoad') || 'Failed to load rules'
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
  editingRule.action = 'review'
  editingRule.priority = 50
  editingRule.enabled = true
  editingRule.conditions = getDefaultConditions()
  editingRule.schema_id = schemas.value.length > 0 ? schemas.value[0]?.id : undefined
  
  // Reset expression mode and builders (default to manual mode)
  expressionMode.value = 'manual'
  builderType.value = 'polygon'
  polygonConfig.latField = ''
  polygonConfig.lonField = ''
  polygonConfig.points = [[0, 0], [0, 0], [0, 0]]
  velocityConfig.metric = 'count'
  velocityConfig.groupField = ''
  velocityConfig.threshold = 10
  velocityConfig.timeValue = 1
  velocityConfig.timeUnit = 'hours'
  
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
  editingRule.action = rule.action
  editingRule.priority = rule.priority
  editingRule.enabled = rule.enabled
  editingRule.schema_id = rule.schema_id
  
  // Preserve existing conditions, but ensure custom_expression exists
  const expression = rule.conditions?.custom_expression || ''
  
  if (expression) {
    editingRule.conditions = { custom_expression: expression }
    
    // Try to parse as velocity expression first
    const velocityParsed = parseVelocityExpression(expression)
    if (velocityParsed) {
      expressionMode.value = 'builder'
      builderType.value = 'velocity'
      velocityConfig.metric = velocityParsed.metric
      velocityConfig.groupField = velocityParsed.groupField
      const timeConfig = convertSecondsToTimeValue(velocityParsed.timeWindowSeconds)
      velocityConfig.timeValue = timeConfig.value
      velocityConfig.timeUnit = timeConfig.unit
      velocityConfig.threshold = velocityParsed.threshold
      conditionRows.value = []
      showModal.value = true
      return
    }
    
    // Try to parse as polygon expression
    const polygonConfigParsed = parsePolygonExpression(expression)
    if (polygonConfigParsed) {
      expressionMode.value = 'builder'
      builderType.value = 'polygon'
      
      // Find lat/lon fields in schema by substring match
      const currentSchema = editingRule.schema_id 
        ? schemas.value.find(s => s.id === editingRule.schema_id) 
        : null
      
      if (currentSchema && currentSchema.extracted_fields) {
        const fields = currentSchema.extracted_fields
        const latFieldLower = polygonConfigParsed.latField.toLowerCase()
        const lonFieldLower = polygonConfigParsed.lonField.toLowerCase()
        
        // Find latitude field (check for 'lat' substring)
        const latFieldMatch = fields.find(f => 
          f.path.toLowerCase().includes('lat') || latFieldLower.includes(f.path.toLowerCase())
        )
        if (latFieldMatch) {
          polygonConfig.latField = latFieldMatch.path
        } else {
          polygonConfig.latField = polygonConfigParsed.latField
        }
        
        // Find longitude field (check for 'lon' or 'lng' substring)
        const lonFieldMatch = fields.find(f => 
          f.path.toLowerCase().includes('lon') || 
          f.path.toLowerCase().includes('lng') || 
          lonFieldLower.includes(f.path.toLowerCase())
        )
        if (lonFieldMatch) {
          polygonConfig.lonField = lonFieldMatch.path
        } else {
          polygonConfig.lonField = polygonConfigParsed.lonField
        }
      } else {
        // Fallback to parsed values if no schema
        polygonConfig.latField = polygonConfigParsed.latField
        polygonConfig.lonField = polygonConfigParsed.lonField
      }
      
      // Ensure we have at least 3 points, pad with [0,0] if needed
      while (polygonConfigParsed.points.length < 3) {
        polygonConfigParsed.points.push([0, 0])
      }
      polygonConfig.points = polygonConfigParsed.points
      conditionRows.value = []
      showModal.value = true
      return
    }
    
    // For other expressions, parse into condition rows for manual editing
    expressionMode.value = 'manual'
    conditionRows.value = parseManualExpression(expression)
  } else {
    editingRule.conditions = getDefaultConditions()
    expressionMode.value = 'builder'
    builderType.value = 'polygon'
    conditionRows.value = [{
      id: `condition-${++conditionRowIdCounter}`,
      field: '',
      operator: '',
      value: null,
      logicOperator: 'and',
    }]
  }
  
  showModal.value = true
}

function closeModal() {
  showModal.value = false
}

async function handleSubmit() {
  // Build expression from current mode
  if (expressionMode.value === 'builder') {
    if (builderType.value === 'polygon') {
      editingRule.conditions.custom_expression = polygonExpression.value
    } else if (builderType.value === 'velocity') {
      editingRule.conditions.custom_expression = velocityExpression.value
    }
  } else {
    // Build expression from condition rows
    editingRule.conditions.custom_expression = generatedExpression.value
  }

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
    error.value = e.message || (window as any).$i18n?.global?.t?.('views.rules.errorSave') || 'Failed to save rule'
  } finally {
    saving.value = false
  }
}

function validateConditions(): string | null {
  // Validate that a schema is selected
  if (!editingRule.schema_id || !currentSchema.value) {
    return t('views.rules.modal.validation.selectSchema')
  }

  const schemaFields = currentSchema.value.extracted_fields || []
  
  if (schemaFields.length === 0) {
    return t('views.rules.modal.validation.noFields')
  }

  // Validate based on expression mode
  if (expressionMode.value === 'builder') {
    if (builderType.value === 'polygon') {
      if (!polygonConfig.latField || !polygonConfig.lonField) {
        return t('views.rules.modal.validation.selectLatLon')
      }
      const validPoints = polygonConfig.points.filter(p => p[0] !== 0 || p[1] !== 0)
      if (validPoints.length < 3) {
        return t('views.rules.modal.validation.polygonPoints')
      }
      // Validate latitude and longitude ranges
      for (let i = 0; i < validPoints.length; i++) {
        const point = validPoints[i]
        if (!point || point.length < 2) {
          return t('views.rules.modal.validation.pointInvalidFormat', { index: i + 1 })
        }
        if (typeof point[0] !== 'number' || point[0] < -90 || point[0] > 90) {
          return t('views.rules.modal.validation.pointLatitudeRange', { index: i + 1, value: point[0] })
        }
        if (typeof point[1] !== 'number' || point[1] < -180 || point[1] > 180) {
          return t('views.rules.modal.validation.pointLongitudeRange', { index: i + 1, value: point[1] })
        }
      }
      const expression = polygonExpression.value
      if (!expression || expression.trim() === '') {
        return t('views.rules.modal.validation.invalidPolygon')
      }
      // Sanitize the polygon expression
      const sanitizeResult = sanitizeExpression(expression)
      if (!sanitizeResult.valid) {
        return sanitizeResult.error || t('views.rules.modal.validation.invalidPolygonExpression')
      }
      // Builder mode validation complete - return early
      return null
    } else if (builderType.value === 'velocity') {
      if (!velocityConfig.groupField) {
        return t('views.rules.modal.validation.selectGroupField')
      }
      if (!velocityConfig.threshold || velocityConfig.threshold <= 0) {
        return t('views.rules.modal.validation.validThreshold')
      }
      if (!velocityConfig.timeValue || velocityConfig.timeValue <= 0) {
        return t('views.rules.modal.validation.validTimeValue')
      }
      const expression = velocityExpression.value
      if (!expression || expression.trim() === '') {
        return t('views.rules.modal.validation.invalidVelocity')
      }
      // Sanitize the velocity expression
      const sanitizeResult = sanitizeExpression(expression)
      if (!sanitizeResult.valid) {
        return sanitizeResult.error || t('views.rules.modal.validation.invalidVelocityExpression')
      }
      // Builder mode validation complete - return early
      return null
    }
  } else {
    // Manual mode: Validate condition rows
    if (conditionRows.value.length === 0) {
      return t('views.rules.modal.validation.atLeastOneCondition')
    }

    // Validate each condition row
    for (let i = 0; i < conditionRows.value.length; i++) {
      const row = conditionRows.value[i]
      
      if (!row) {
        continue
      }
      
      if (!row.field) {
        return t('views.rules.modal.validation.conditionFieldRequired', { index: i + 1 })
      }
      
      if (!row.operator) {
        return t('views.rules.modal.validation.conditionOperatorRequired', { index: i + 1 })
      }
      
      if (row.value === null || row.value === undefined || row.value === '') {
        return t('views.rules.modal.validation.conditionValueRequired', { index: i + 1 })
      }

      // Validate that the field exists in the schema
      const fieldExists = schemaFields.some(f => f.path === row.field)
      if (!fieldExists) {
        return t('views.rules.modal.validation.conditionFieldNotExist', { index: i + 1, field: row.field })
      }

      // Validate operator is appropriate for field type
      const field = schemaFields.find(f => f.path === row.field)
      if (field) {
        const validOperators = getOperatorOptions(row.field).map(o => o.value)
        if (!validOperators.includes(row.operator)) {
          return t('views.rules.modal.validation.conditionOperatorNotValid', { index: i + 1, operator: row.operator, type: field.type })
        }
      }
    }

    // Validate generated expression for manual mode
    const expression = generatedExpression.value
    if (!expression || expression.trim() === '') {
      return t('views.rules.modal.validation.emptyExpression')
    }
    
    // Sanitize the expression
    const sanitizeResult = sanitizeExpression(expression)
    if (!sanitizeResult.valid) {
      return sanitizeResult.error || t('views.rules.modal.validation.invalidExpression')
    }
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
  
  for (const field of Array.from(foundFields)) {
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
    error.value = e.message || (window as any).$i18n?.global?.t?.('views.rules.errorDelete') || 'Failed to delete rule'
  }
}

async function toggleRule(rule: any) {
  try {
    await api.put(`/api/v1/rules/${rule.id}`, { ...rule, enabled: !rule.enabled })
    await loadRules()
  } catch (e: any) {
    error.value = e.message || (window as any).$i18n?.global?.t?.('views.rules.errorToggle') || 'Failed to toggle rule'
  }
}

// Get schema name by ID
function getSchemaName(schemaId: string | undefined): string {
  if (!schemaId) {
    return ''
  }
  const schema = schemas.value.find(s => s.id === schemaId)
  return schema?.name || ''
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
    default:
      return 'default'
  }
}

function getActionLabel(action: string): string {
  const actionLower = action.toLowerCase()
  const key = `views.rules.actions.${actionLower}`
  try {
    return i18n.global.t(key)
  } catch {
    return action
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
