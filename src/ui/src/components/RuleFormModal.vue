<template>
  <BaseModal
    :model-value="modelValue"
    :title="isEditing ? $t('views.rules.modal.editTitle') : $t('views.rules.modal.createTitle')"
    size="lg"
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <v-form ref="formRef" @submit.prevent="handleSubmit" class="mt-4">
      <!-- Presets Section (only for new rules) -->
      <div v-if="!isEditing" class="mb-6">
        <label class="text-body-2 text-grey-darken-1 d-block mb-2">{{ $t('views.rules.modal.quickStart') }}</label>
        <div class="d-flex flex-wrap gap-2">
          <v-chip
            v-for="preset in rulePresets"
            :key="preset.id"
            @click="$emit('apply-preset', preset)"
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
        :label="$t('views.rules.modal.ruleName')"
        :placeholder="$t('views.rules.modal.ruleNamePlaceholder')"
        required
        :rules="[(v: string) => !!v || $t('views.rules.modal.nameRequired')]"
        prepend-inner-icon="fa-text"
        :hint="$t('views.rules.modal.ruleNameHint')"
        persistent-hint
        class="mb-4"
      />

      <BaseInput
        v-model="editingRule.description"
        :label="$t('views.rules.modal.ruleDescription')"
        :placeholder="$t('views.rules.modal.ruleDescriptionPlaceholder')"
        required
        :rules="[(v: string) => !!v || $t('views.rules.modal.descriptionRequired')]"
        prepend-inner-icon="fa-align-left"
        :hint="$t('views.rules.modal.ruleDescriptionHint')"
        persistent-hint
        class="mb-4"
      />

      <BaseSelect
        v-model="editingRule.schema_id"
        :label="$t('views.rules.modal.eventSchema')"
        :options="schemaOptions"
        :rules="[(v: string) => !!v || $t('views.rules.modal.schemaRequired')]"
        :hint="$t('views.rules.modal.eventSchemaHint')"
        persistent-hint
        class="mb-4"
      />

      <!-- Available Fields from Selected Schema -->
      <div v-if="currentSchema && currentSchema.extracted_fields?.length > 0" class="mb-4 pa-3 bg-blue-lighten-5 rounded-lg">
        <div class="d-flex justify-space-between align-center mb-2">
          <label class="text-caption text-grey-darken-1 d-flex align-center">
            <v-icon icon="fa-info-circle" size="x-small" class="mr-1" />
            {{ $t('views.rules.modal.availableFields') }} "{{ currentSchema.name }}" ({{ currentSchema.extracted_fields.length }}):
          </label>
          <v-btn
            v-if="currentSchema.extracted_fields.length > 10"
            size="x-small"
            variant="text"
            @click="showAllFields = !showAllFields"
          >
            {{ showAllFields ? $t('views.rules.modal.showLess') : `${$t('views.rules.showAll')} (${currentSchema.extracted_fields.length})` }}
          </v-btn>
        </div>
        <div class="d-flex flex-wrap gap-1">
          <v-chip
            v-for="field in (showAllFields ? currentSchema.extracted_fields : currentSchema.extracted_fields.slice(0, 10))"
            :key="field.path"
            size="x-small"
            variant="outlined"
            color="primary"
          >
            {{ field.path }} <span class="text-grey ml-1">({{ field.type }})</span>
          </v-chip>
        </div>
        <p class="text-caption text-grey-darken-1 mt-2">
          {{ $t('views.rules.modal.availableFieldsHint') }}
        </p>
      </div>

      <BaseSelect
        v-model="editingRule.action"
        :label="$t('views.rules.modal.action')"
        :options="ruleActions"
        required
        :rules="[(v: string) => !!v || $t('views.rules.modal.actionRequired')]"
        class="mb-4"
      />

      <BaseInput
        v-model.number="editingRule.priority"
        type="number"
        :label="$t('views.rules.modal.priority')"
        :min="0"
        :max="100"
        required
        :rules="[
          (v: any) => (v !== null && v !== undefined && v !== '') || $t('views.rules.modal.priorityRequired'),
          (v: any) => (typeof v === 'number' && v >= 0 && v <= 100) || $t('views.rules.modal.priorityRange')
        ]"
        :hint="$t('views.rules.modal.priorityHint')"
        persistent-hint
        prepend-inner-icon="fa-sort"
        class="mb-4"
      />

      <!-- Rule Expression Section -->
      <slot name="expression-builder" />

      <!-- Validation Error Display -->
      <v-alert
        v-if="validationError"
        type="error"
        variant="tonal"
        dismissible
        @click:close="validationError = ''"
        class="mb-4"
      >
        {{ validationError }}
      </v-alert>

      <div class="mb-6">
        <label class="text-body-2 text-grey-darken-1 d-block mb-2">{{ $t('views.rules.modal.status') }}</label>
        <v-btn-toggle
          v-model="editingRule.enabled"
          mandatory
          color="primary"
          variant="outlined"
        >
          <v-btn :value="true" class="px-6">
            <v-icon icon="fa-check" size="small" class="mr-2" />
            {{ $t('views.rules.modal.enabled') }}
          </v-btn>
          <v-btn :value="false" class="px-6">
            <v-icon icon="fa-ban" size="small" class="mr-2" />
            {{ $t('views.rules.modal.disabled') }}
          </v-btn>
        </v-btn-toggle>
      </div>
    </v-form>

    <template #footer>
      <BaseButton variant="ghost" @click="$emit('cancel')" prepend-icon="fa-xmark">{{ $t('views.rules.modal.cancel') }}</BaseButton>
      <BaseButton @click="handleSubmit" :loading="saving" prepend-icon="fa-save">{{ $t('views.rules.modal.save') }}</BaseButton>
    </template>
  </BaseModal>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { i18n } from '@/plugins/i18n'
import BaseModal from '@/components/BaseModal.vue'
import BaseInput from '@/components/BaseInput.vue'
import BaseSelect from '@/components/BaseSelect.vue'
import BaseButton from '@/components/BaseButton.vue'

const t = i18n.global.t

interface Props {
  modelValue: boolean
  isEditing: boolean
  editingRule: any
  schemaOptions: Array<{ value: string; label: string }>
  currentSchema: any
  rulePresets: any[]
  ruleActions: Array<{ value: string; label: string }>
  saving: boolean
  expressionMode: 'builder' | 'manual'
  builderType: 'polygon' | 'velocity'
  polygonConfig: {
    latField: string
    lonField: string
    points: Array<[number, number]>
  }
  velocityConfig: {
    metric: 'count' | 'sum'
    groupField: string
    threshold: number
    timeValue: number
    timeUnit: 'seconds' | 'minutes' | 'hours' | 'days'
  }
  conditionRows: Array<{
    id: string
    field: string
    operator: string
    value: any
    logicOperator: 'and' | 'or'
  }>
  polygonExpression: string
  velocityExpression: string
  generatedExpression: string
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  'apply-preset': [preset: any]
  'submit': []
  'cancel': []
}>()

const formRef = ref<any>(null)
const showAllFields = ref(false)
const validationError = ref('')

// Watch for modal open/close to clear validation error
watch(() => props.modelValue, (isOpen) => {
  if (!isOpen) {
    validationError.value = ''
  }
})

async function handleSubmit() {
  // Clear previous validation error
  validationError.value = ''
  
  // Validate form fields first
  if (formRef.value) {
    const { valid } = await formRef.value.validate()
    if (!valid) {
      return
    }
  }
  
  // Validate conditions
  const conditionError = validateConditions()
  if (conditionError) {
    validationError.value = conditionError
    return
  }
  
  // If all validations pass, emit submit
  emit('submit')
}

function validateConditions(): string | null {
  // Validate that a schema is selected
  if (!props.editingRule.schema_id || !props.currentSchema) {
    return t('views.rules.modal.validation.selectSchema')
  }

  const schemaFields = props.currentSchema.extracted_fields || []
  
  if (schemaFields.length === 0) {
    return t('views.rules.modal.validation.noFields')
  }

  // Validate based on expression mode
  if (props.expressionMode === 'builder') {
    if (props.builderType === 'polygon') {
      if (!props.polygonConfig.latField || !props.polygonConfig.lonField) {
        return t('views.rules.modal.validation.selectLatLon')
      }
      const validPoints = props.polygonConfig.points.filter(p => p[0] !== 0 || p[1] !== 0)
      if (validPoints.length < 3) {
        return 'Polygon must have at least 3 points'
      }
      // Validate latitude and longitude ranges
      for (let i = 0; i < validPoints.length; i++) {
        const point = validPoints[i]
        if (!point || point.length < 2) {
          return `Point ${i + 1}: Invalid point format`
        }
        if (typeof point[0] !== 'number' || point[0] < -90 || point[0] > 90) {
          return `Point ${i + 1}: Latitude must be between -90 and 90 (got ${point[0]})`
        }
        if (typeof point[1] !== 'number' || point[1] < -180 || point[1] > 180) {
          return `Point ${i + 1}: Longitude must be between -180 and 180 (got ${point[1]})`
        }
      }
      const expression = props.polygonExpression
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
    } else if (props.builderType === 'velocity') {
      if (!props.velocityConfig.groupField) {
        return t('views.rules.modal.validation.selectGroupField')
      }
      if (!props.velocityConfig.threshold || props.velocityConfig.threshold <= 0) {
        return t('views.rules.modal.validation.validThreshold')
      }
      if (!props.velocityConfig.timeValue || props.velocityConfig.timeValue <= 0) {
        return t('views.rules.modal.validation.validTimeValue')
      }
      const expression = props.velocityExpression
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
    } else {
      // Invalid builder type
      return t('views.rules.modal.validation.invalidBuilderType')
    }
  } else {
    // Manual mode: Validate condition rows
    if (props.conditionRows.length === 0) {
      return t('views.rules.modal.validation.atLeastOneCondition')
    }

    // Validate each condition row
    for (let i = 0; i < props.conditionRows.length; i++) {
      const row = props.conditionRows[i]
      
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
      const fieldExists = schemaFields.some((f: any) => f.path === row.field)
      if (!fieldExists) {
        return `Condition ${i + 1}: Field "${row.field}" does not exist in the selected schema`
      }

      // Validate operator is appropriate for field type
      const field = schemaFields.find((f: any) => f.path === row.field)
      if (field) {
        const validOperators = getOperatorOptions(row.field).map(o => o.value)
        if (!validOperators.includes(row.operator)) {
          return `Condition ${i + 1}: Operator "${row.operator}" is not valid for field type "${field.type}"`
        }
      }
    }

    // Validate generated expression for manual mode
    const expression = props.generatedExpression
    if (!expression || expression.trim() === '') {
      return 'Generated expression is empty. Please check your conditions.'
    }
    
    // Sanitize the expression
    const sanitizeResult = sanitizeExpression(expression)
    if (!sanitizeResult.valid) {
      return sanitizeResult.error || 'Invalid expression'
    }
  }

  return null
}

function getOperatorOptions(fieldPath: string): Array<{ value: string; label: string }> {
  if (!props.currentSchema || !fieldPath) {
    return []
  }
  
  const field = props.currentSchema.extracted_fields?.find((f: any) => f.path === fieldPath)
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
</script>

