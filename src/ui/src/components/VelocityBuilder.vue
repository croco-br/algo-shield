<template>
  <div class="pa-3 bg-white rounded-lg">
    <BaseSelect
      v-model="config.metric"
      label="Metric"
      :options="[
        { value: 'count', label: 'Transaction Count' },
        { value: 'sum', label: 'Amount Sum' }
      ]"
      hint="Choose what to measure"
      persistent-hint
      class="mb-3"
    />
    <BaseSelect
      v-model="config.groupField"
      label="Group By Field"
      :options="fieldOptions"
      hint="Field to group velocity checks by (e.g., origin)"
      persistent-hint
      class="mb-3"
    />
    <BaseInput
      v-model.number="config.threshold"
      type="number"
      label="Threshold"
      hint="Maximum allowed value"
      persistent-hint
      class="mb-3"
    />
    <div class="d-flex gap-2 mb-3">
      <BaseInput
        v-model.number="config.timeValue"
        type="number"
        label="Time Value"
        class="flex-grow-1"
      />
      <BaseSelect
        v-model="config.timeUnit"
        label="Time Unit"
        :options="[
          { value: 'seconds', label: 'Seconds' },
          { value: 'minutes', label: 'Minutes' },
          { value: 'hours', label: 'Hours' },
          { value: 'days', label: 'Days' }
        ]"
        style="min-width: 150px;"
      />
    </div>

    <div v-if="expression" class="mt-3 pa-3 bg-grey-darken-1 rounded">
      <div class="text-caption text-white mb-1">Generated Expression:</div>
      <code class="text-white">{{ expression }}</code>
    </div>
  </div>
</template>

<script setup lang="ts">
import BaseSelect from '@/components/BaseSelect.vue'
import BaseInput from '@/components/BaseInput.vue'

interface Props {
  config: {
    metric: 'count' | 'sum'
    groupField: string
    threshold: number
    timeValue: number
    timeUnit: 'seconds' | 'minutes' | 'hours' | 'days'
  }
  fieldOptions: Array<{ value: string; label: string }>
  expression: string
}

defineProps<Props>()
</script>

