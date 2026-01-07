<template>
  <div>
    <div v-if="conditionRows.length === 0" class="text-center pa-6 text-grey-darken-1">
      <v-icon icon="fa-info-circle" size="large" class="mb-2" />
      <p class="text-body-2">{{ $t('views.rules.modal.conditionBuilder.noConditions') }}</p>
    </div>

    <div v-for="(row, index) in conditionRows" :key="row.id" class="mb-4 pa-3 bg-white rounded-lg border">
      <div class="d-flex align-center gap-2">
        <div class="flex-grow-1">
          <div class="d-flex align-center gap-2 flex-wrap">
            <!-- Field Dropdown -->
            <BaseSelect
              v-model="row.field"
              :label="$t('views.rules.modal.conditionBuilder.field')"
              :options="fieldOptions"
              :rules="[(v: string) => !!v || $t('views.rules.modal.conditionBuilder.fieldRequired')]"
              class="flex-grow-1"
              style="min-width: 200px;"
            />

            <!-- Comparison Operator -->
            <BaseSelect
              v-model="row.operator"
              :label="$t('views.rules.modal.conditionBuilder.operator')"
              :options="getOperatorOptions(row.field)"
              :rules="[(v: string) => !!v || $t('views.rules.modal.conditionBuilder.operatorRequired')]"
              style="min-width: 150px;"
            />

            <!-- Value Input -->
            <div class="flex-grow-1" style="min-width: 200px;">
              <BaseInput
                v-if="getValueInputType(row.field) === 'number'"
                v-model.number="row.value"
                type="number"
                :label="$t('views.rules.modal.conditionBuilder.value')"
                :rules="[(v: any) => v !== null && v !== undefined && v !== '' || $t('views.rules.modal.conditionBuilder.valueRequired')]"
              />
              <div v-else-if="getValueInputType(row.field) === 'boolean'">
                <label class="text-body-2 text-grey-darken-1 d-block mb-1">{{ $t('views.rules.modal.conditionBuilder.value') }} <span class="text-red">*</span></label>
                <v-select
                  v-model="row.value"
                  :items="[{ title: 'true', value: true }, { title: 'false', value: false }]"
                  variant="outlined"
                  density="compact"
                  :rules="[(v: any) => v !== null && v !== undefined && v !== '' || $t('views.rules.modal.conditionBuilder.valueRequired')]"
                />
              </div>
              <BaseInput
                v-else-if="getValueInputType(row.field) === 'array'"
                v-model="row.value"
                :label="$t('views.rules.modal.conditionBuilder.valueCommaSeparated')"
                :placeholder="$t('views.rules.modal.conditionBuilder.valuePlaceholder')"
                :rules="[(v: string) => !!v || $t('views.rules.modal.conditionBuilder.valueRequired')]"
                :hint="$t('views.rules.modal.conditionBuilder.valueHint')"
                persistent-hint
              />
              <BaseInput
                v-else
                v-model="row.value"
                :label="$t('views.rules.modal.conditionBuilder.value')"
                :rules="[(v: string) => !!v || $t('views.rules.modal.conditionBuilder.valueRequired')]"
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
          @click="$emit('remove-condition', index)"
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
            {{ $t('views.rules.modal.conditionBuilder.and') }}
          </v-btn>
          <v-btn value="or" size="small">
            {{ $t('views.rules.modal.conditionBuilder.or') }}
          </v-btn>
        </v-btn-toggle>
        <v-divider class="flex-grow-1" />
      </div>
    </div>

    <!-- Generated Expression Preview -->
    <div v-if="conditionRows.length > 0" class="mt-4 pa-3 bg-grey-darken-1 rounded">
      <div class="text-caption text-white mb-1">{{ $t('views.rules.generatedExpression') }}:</div>
      <code class="text-white">{{ generatedExpression }}</code>
    </div>
  </div>
</template>

<script setup lang="ts">
import { i18n } from '@/plugins/i18n'
import BaseSelect from '@/components/BaseSelect.vue'
import BaseInput from '@/components/BaseInput.vue'

const t = i18n.global.t

interface ConditionRow {
  id: string
  field: string
  operator: string
  value: any
  logicOperator: 'and' | 'or'
}

interface Props {
  conditionRows: ConditionRow[]
  fieldOptions: Array<{ value: string; label: string }>
  currentSchema: any
  generatedExpression: string
}

const props = defineProps<Props>()

defineEmits<{
  'remove-condition': [index: number]
}>()

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
  
  return [
    { value: '==', label: 'Equals (==)' },
    { value: '!=', label: 'Not Equals (!=)' },
  ]
}

function getValueInputType(fieldPath: string): 'number' | 'string' | 'boolean' | 'array' {
  if (!props.currentSchema || !fieldPath) {
    return 'string'
  }
  
  const field = props.currentSchema.extracted_fields?.find((f: any) => f.path === fieldPath)
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

function getValuePlaceholder(fieldPath: string): string {
  const defaultValue = t('views.rules.modal.validation.enterValue')
  
  if (!props.currentSchema || !fieldPath) {
    return defaultValue
  }
  
  const field = props.currentSchema.extracted_fields?.find((f: any) => f.path === fieldPath)
  if (!field) {
    return defaultValue
  }

  const type = field.type
  if (type === 'string') {
    return 'e.g., "USD" or "active"'
  } else if (type === 'number') {
    return 'e.g., 1000 or 5000'
  } else if (type === 'array') {
    return 'e.g., "US", "CA", "GB" or 1, 2, 3'
  }
  
  return defaultValue
}
</script>

