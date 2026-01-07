<template>
  <div class="pa-3 bg-white rounded-lg">
    <BaseSelect
      v-model="config.metric"
      :label="$t('views.rules.modal.velocityBuilder.metric')"
      :options="[
        { value: 'count', label: $t('views.rules.modal.velocityBuilder.transactionCount') },
        { value: 'sum', label: $t('views.rules.modal.velocityBuilder.amountSum') }
      ]"
      :hint="$t('views.rules.modal.velocityBuilder.metricHint')"
      persistent-hint
      class="mb-3"
    />
    <BaseSelect
      v-model="config.groupField"
      :label="$t('views.rules.modal.velocityBuilder.groupByField')"
      :options="fieldOptions"
      :hint="$t('views.rules.modal.velocityBuilder.groupByFieldHint')"
      persistent-hint
      class="mb-3"
    />
    <BaseInput
      v-model.number="config.threshold"
      type="number"
      :label="$t('views.rules.modal.velocityBuilder.threshold')"
      :hint="$t('views.rules.modal.velocityBuilder.thresholdHint')"
      persistent-hint
      class="mb-3"
    />
    <div class="d-flex gap-2 mb-3">
      <BaseInput
        v-model.number="config.timeValue"
        type="number"
        :label="$t('views.rules.modal.velocityBuilder.timeValue')"
        class="flex-grow-1"
      />
      <BaseSelect
        v-model="config.timeUnit"
        :label="$t('views.rules.modal.velocityBuilder.timeUnit')"
        :options="[
          { value: 'seconds', label: $t('views.rules.modal.velocityBuilder.seconds') },
          { value: 'minutes', label: $t('views.rules.modal.velocityBuilder.minutes') },
          { value: 'hours', label: $t('views.rules.modal.velocityBuilder.hours') },
          { value: 'days', label: $t('views.rules.modal.velocityBuilder.days') }
        ]"
        style="min-width: 150px;"
      />
    </div>

    <div v-if="expression" class="mt-3 pa-3 bg-grey-darken-1 rounded">
      <div class="text-caption text-white mb-1">{{ $t('views.rules.generatedExpression') }}:</div>
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

