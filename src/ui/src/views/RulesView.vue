<template>
  <v-container fluid class="pa-8">
    <div class="d-flex justify-space-between align-center mb-10">
      <div>
        <div class="d-flex align-center gap-3 mb-2">
          <v-icon icon="fa-tasks" size="large" color="primary" />
          <h2 class="text-h4 font-weight-bold">Rules Management</h2>
        </div>
        <p class="text-body-1 text-grey-darken-1">Configure custom rules for fraud detection and AML</p>
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
          label="Name"
          placeholder="Rule name"
          required
          :rules="[(v: string) => !!v || 'Name is required']"
          prepend-inner-icon="fa-text"
          class="mb-4"
        />

        <BaseInput
          v-model="editingRule.description"
          label="Description"
          placeholder="Description"
          required
          :rules="[(v: string) => !!v || 'Description is required']"
          prepend-inner-icon="fa-align-left"
          class="mb-4"
        />

        <BaseSelect
          v-model="editingRule.type"
          label="Type"
          :options="ruleTypes"
          required
          :rules="[(v: string) => !!v || 'Type is required']"
          class="mb-4"
        />

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

        <!-- Rule Conditions Section -->
        <div v-if="editingRule.type" class="mb-6 pa-4 bg-grey-lighten-5 rounded-lg">
          <h4 class="text-subtitle-1 font-weight-medium mb-4 d-flex align-center">
            <v-icon icon="fa-cog" size="small" class="mr-2" />
            Rule Conditions
          </h4>

          <!-- Amount Conditions -->
          <template v-if="editingRule.type === 'amount'">
            <BaseInput
              v-model.number="editingRule.conditions.amount_threshold"
              type="number"
              label="Amount Threshold"
              :min="0"
              required
              :rules="[
                (v: any) => (v !== null && v !== undefined && v !== '') || 'Amount threshold is required',
                (v: any) => (typeof v === 'number' && v > 0) || 'Amount must be greater than 0'
              ]"
              hint="Transactions above this amount will match this rule"
              persistent-hint
              prepend-inner-icon="fa-dollar-sign"
            />
          </template>

          <!-- Velocity Conditions -->
          <template v-else-if="editingRule.type === 'velocity'">
            <BaseInput
              v-model.number="editingRule.conditions.transaction_count"
              type="number"
              label="Transaction Count"
              :min="1"
              required
              :rules="[
                (v: any) => (v !== null && v !== undefined && v !== '') || 'Transaction count is required',
                (v: any) => (typeof v === 'number' && v >= 1) || 'Count must be at least 1'
              ]"
              hint="Maximum number of transactions allowed in the time window"
              persistent-hint
              prepend-inner-icon="fa-hashtag"
              class="mb-4"
            />
            <BaseInput
              v-model.number="editingRule.conditions.time_window_seconds"
              type="number"
              label="Time Window (seconds)"
              :min="1"
              required
              :rules="[
                (v: any) => (v !== null && v !== undefined && v !== '') || 'Time window is required',
                (v: any) => (typeof v === 'number' && v >= 1) || 'Time window must be at least 1 second'
              ]"
              hint="Time period in seconds (e.g., 3600 = 1 hour)"
              persistent-hint
              prepend-inner-icon="fa-clock"
            />
          </template>

          <!-- Blocklist Conditions -->
          <template v-else-if="editingRule.type === 'blocklist'">
            <v-combobox
              v-model="editingRule.conditions.blocklisted_accounts"
              label="Blocked Accounts"
              multiple
              chips
              closable-chips
              clearable
              variant="outlined"
              density="comfortable"
              hint="Type account ID and press Enter to add. Transactions from/to these accounts will match."
              persistent-hint
              prepend-inner-icon="fa-ban"
              :rules="[
                (v: any) => (Array.isArray(v) && v.length > 0) || 'At least one account is required'
              ]"
            />
          </template>

          <!-- Geography Conditions -->
          <template v-else-if="editingRule.type === 'geography'">
            <div class="geography-polygon-container">
              <label class="text-body-2 text-grey-darken-1 d-block mb-2">Restricted Zone Polygon (Lat/Lon Coordinates)</label>
              <v-textarea
                v-model="editingRule.conditions.polygon_coordinates"
                variant="outlined"
                density="comfortable"
                rows="4"
                placeholder='[
  [40.7128, -74.0060],
  [40.7580, -73.9855],
  [40.7282, -73.7949],
  [40.7128, -74.0060]
]'
                prepend-inner-icon="fa-map-marker-alt"
                :rules="[
                  (v: any) => !!v || 'Polygon coordinates are required',
                  (v: any) => isValidPolygon(v) || 'Invalid polygon format. Use JSON array of [lat, lon] pairs.'
                ]"
              />
              <p class="text-caption text-grey-darken-1 mt-2">
                Define a closed polygon using lat/lon coordinate pairs in JSON format.<br>
                Format: <code>[[lat1, lon1], [lat2, lon2], ..., [lat1, lon1]]</code><br>
                The first and last coordinates should be the same to close the polygon.
              </p>
            </div>
          </template>

          <!-- Custom Conditions -->
          <template v-else-if="editingRule.type === 'custom'">
            <div class="custom-expression-container">
              <label class="text-body-2 text-grey-darken-1 d-block mb-2">Custom Expression</label>
              <div class="code-editor-wrapper">
                <textarea
                  ref="customExpressionRef"
                  v-model="editingRule.conditions.custom_expression"
                  class="code-editor"
                  :placeholder="customExpressionPlaceholder"
                  rows="6"
                  @input="updateHighlighting"
                  @scroll="syncScroll"
                />
                <pre class="code-highlight" ref="highlightRef"><code v-html="highlightedExpression"></code></pre>
              </div>
              <p class="text-caption text-grey-darken-1 mt-2">
                <strong>Available fields:</strong> <code>amount</code>, <code>currency</code>, <code>type</code>, 
                <code>from_account</code>, <code>to_account</code>, <code>region</code>, <code>country</code>, 
                <code>metadata.*</code><br>
                <strong>Operators:</strong> <code>==</code>, <code>!=</code>, <code>&gt;</code>, <code>&lt;</code>, 
                <code>&gt;=</code>, <code>&lt;=</code>, <code>and</code>, <code>or</code>, <code>not</code>, 
                <code>in</code>, <code>contains</code><br>
                <strong>Examples:</strong><br>
                • <code>amount &gt; 10000</code><br>
                • <code>currency == "USD" and amount &gt; 5000</code><br>
                • <code>type in ["international", "wire"]</code><br>
                • <code>metadata.risk_score &gt; 50</code>
              </p>
            </div>
          </template>
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
})

// Rule types must match backend validation: amount|velocity|blocklist|geography|custom
const ruleTypes = [
  { value: 'amount', label: 'Amount Threshold' },
  { value: 'velocity', label: 'Velocity Check' },
  { value: 'blocklist', label: 'Blocklist' },
  { value: 'geography', label: 'Geography' },
  { value: 'custom', label: 'Custom' },
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
    label: 'High Value Transaction',
    icon: 'fa-dollar-sign',
    name: 'High Value Transaction',
    description: 'Flag transactions above $10,000 for review',
    type: 'amount',
    action: 'review',
    score: 30,
    conditions: { amount_threshold: 10000 },
  },
  {
    id: 'suspicious-velocity',
    label: 'Suspicious Velocity',
    icon: 'fa-bolt',
    name: 'Suspicious Velocity',
    description: 'Flag accounts with more than 10 transactions per hour',
    type: 'velocity',
    action: 'review',
    score: 50,
    conditions: { transaction_count: 10, time_window_seconds: 3600 },
  },
  {
    id: 'suspicious-metadata',
    label: 'Suspicious Metadata',
    icon: 'fa-exclamation-triangle',
    name: 'Suspicious Metadata Flag',
    description: 'Check for suspicious flag in transaction metadata',
    type: 'custom',
    action: 'review',
    score: 60,
    conditions: { custom_expression: 'metadata.is_suspicious == true' },
  },
  {
    id: 'high-amount-foreign',
    label: 'High Amount Foreign Currency',
    icon: 'fa-coins',
    name: 'High Amount Foreign Currency',
    description: 'Flag high-value transactions in non-USD currency',
    type: 'custom',
    action: 'review',
    score: 35,
    conditions: { custom_expression: 'amount > 5000 and currency != "USD"' },
  },
  {
    id: 'blocked-account',
    label: 'Blocked Account',
    icon: 'fa-ban',
    name: 'Blocked Account',
    description: 'Block transactions from/to specific accounts',
    type: 'blocklist',
    action: 'block',
    score: 100,
    conditions: { blocklisted_accounts: [] },
  },
  {
    id: 'high-risk-region',
    label: 'High Risk Zone',
    icon: 'fa-globe',
    name: 'High Risk Geographic Zone',
    description: 'Flag transactions originating from a restricted geographic polygon',
    type: 'geography',
    action: 'review',
    score: 40,
    conditions: { polygon_coordinates: '[\n  [55.7558, 37.6173],\n  [59.9343, 30.3351],\n  [55.0084, 82.9357],\n  [55.7558, 37.6173]\n]' },
  },
]

// Placeholder for custom expression textarea
// Syntax follows expr-lang: https://github.com/expr-lang/expr
const customExpressionPlaceholder = `Examples:
amount > 10000
currency == "USD"
type == "international"
amount > 5000 and currency != "USD"
from_account == to_account
metadata.is_suspicious == true
region in ["RU", "CN", "KP"]`

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

// Validate polygon JSON format
function isValidPolygon(value: string): boolean {
  if (!value || value.trim() === '') return false
  try {
    const parsed = JSON.parse(value)
    if (!Array.isArray(parsed) || parsed.length < 3) return false
    // Check each point is [lat, lon] pair
    for (const point of parsed) {
      if (!Array.isArray(point) || point.length !== 2) return false
      if (typeof point[0] !== 'number' || typeof point[1] !== 'number') return false
    }
    return true
  } catch {
    return false
  }
}

// Get default conditions based on rule type
function getDefaultConditions(type: string): Record<string, any> {
  switch (type) {
    case 'amount':
      return { amount_threshold: 10000 }
    case 'velocity':
      return { time_window_seconds: 3600, transaction_count: 10 }
    case 'blocklist':
      return { blocklisted_accounts: [] }
    case 'geography':
      return { polygon_coordinates: '[\n  [40.7128, -74.0060],\n  [40.7580, -73.9855],\n  [40.7282, -73.7949],\n  [40.7128, -74.0060]\n]' }
    case 'custom':
      return { custom_expression: '' }
    default:
      return {}
  }
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
  } else {
    router.push('/login')
  }
})

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
  editingRule.type = 'amount'
  editingRule.action = 'review'
  editingRule.score = 0
  editingRule.priority = 500
  editingRule.enabled = true
  editingRule.conditions = getDefaultConditions('amount')
  showModal.value = true
}

function openEditModal(rule: any) {
  isEditing.value = true
  editingRule.id = rule.id
  editingRule.name = rule.name
  editingRule.description = rule.description
  editingRule.type = rule.type
  editingRule.action = rule.action
  editingRule.score = rule.score
  editingRule.priority = rule.priority
  editingRule.enabled = rule.enabled
  // Preserve existing conditions or set defaults
  editingRule.conditions = rule.conditions && Object.keys(rule.conditions).length > 0
    ? { ...rule.conditions }
    : getDefaultConditions(rule.type)
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

  // Additional validation for conditions based on type
  const conditionError = validateConditions()
  if (conditionError) {
    error.value = conditionError
    return
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
    error.value = e.message || 'Failed to save rule'
  } finally {
    saving.value = false
  }
}

function validateConditions(): string | null {
  const { type, conditions } = editingRule

  switch (type) {
    case 'amount':
      if (!conditions.amount_threshold || conditions.amount_threshold <= 0) {
        return 'Amount threshold must be greater than 0'
      }
      break
    case 'velocity':
      if (!conditions.transaction_count || conditions.transaction_count < 1) {
        return 'Transaction count must be at least 1'
      }
      if (!conditions.time_window_seconds || conditions.time_window_seconds < 1) {
        return 'Time window must be at least 1 second'
      }
      break
    case 'blocklist':
      if (!Array.isArray(conditions.blocklisted_accounts) || conditions.blocklisted_accounts.length === 0) {
        return 'At least one blocked account is required'
      }
      break
    case 'geography':
      if (!conditions.polygon_coordinates || !isValidPolygon(conditions.polygon_coordinates)) {
        return 'Valid polygon coordinates are required (JSON array of [lat, lon] pairs)'
      }
      break
    case 'custom':
      if (!conditions.custom_expression || conditions.custom_expression.trim() === '') {
        return 'Custom expression is required'
      }
      // Sanitize the expression
      const sanitizeResult = sanitizeExpression(conditions.custom_expression)
      if (!sanitizeResult.valid) {
        return sanitizeResult.error || 'Invalid expression'
      }
      break
  }

  return null
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
