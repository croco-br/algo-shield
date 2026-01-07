<template>
  <div class="pa-3 bg-white rounded-lg">
    <BaseSelect
      v-model="config.latField"
      :label="$t('views.rules.modal.polygonBuilder.latitudeField')"
      :options="fieldOptions"
      :hint="$t('views.rules.modal.polygonBuilder.latitudeFieldHint')"
      persistent-hint
      class="mb-3"
    />
    <BaseSelect
      v-model="config.lonField"
      :label="$t('views.rules.modal.polygonBuilder.longitudeField')"
      :options="fieldOptions"
      :hint="$t('views.rules.modal.polygonBuilder.longitudeFieldHint')"
      persistent-hint
      class="mb-3"
    />
    
    <div class="mb-3">
      <label class="text-body-2 text-grey-darken-1 d-block mb-2">
        {{ $t('views.rules.modal.polygonBuilder.polygonCoordinates') }} <span class="text-red">*</span>
        <span class="text-caption text-grey">{{ $t('views.rules.modal.polygonBuilder.pointsRequired') }}</span>
      </label>
      <div v-for="(point, index) in config.points" :key="index" class="d-flex gap-2 mb-2">
        <BaseInput
          v-model.number="point[0]"
          type="number"
          :label="$t('views.rules.modal.polygonBuilder.latitude')"
          step="0.000001"
          :placeholder="$t('views.rules.modal.polygonBuilder.latitudePlaceholder')"
          :min="-90"
          :max="90"
          :rules="[
            (v: any) => v !== null && v !== undefined && v !== '' || $t('views.rules.modal.polygonBuilder.latitudeRequired'),
            (v: any) => (typeof v === 'number' && v >= -90 && v <= 90) || $t('views.rules.modal.polygonBuilder.latitudeRange')
          ]"
          :hint="$t('views.rules.modal.polygonBuilder.latitudeRangeHint')"
          persistent-hint
          style="flex: 1;"
        />
        <BaseInput
          v-model.number="point[1]"
          type="number"
          :label="$t('views.rules.modal.polygonBuilder.longitude')"
          step="0.000001"
          :placeholder="$t('views.rules.modal.polygonBuilder.longitudePlaceholder')"
          :min="-180"
          :max="180"
          :rules="[
            (v: any) => v !== null && v !== undefined && v !== '' || $t('views.rules.modal.polygonBuilder.longitudeRequired'),
            (v: any) => (typeof v === 'number' && v >= -180 && v <= 180) || $t('views.rules.modal.polygonBuilder.longitudeRange')
          ]"
          :hint="$t('views.rules.modal.polygonBuilder.longitudeRangeHint')"
          persistent-hint
          style="flex: 1;"
        />
        <v-btn
          icon="fa-trash"
          size="small"
          variant="text"
          color="error"
          @click="$emit('remove-point', index)"
          class="mt-4"
        />
      </div>
      <BaseButton size="sm" @click="$emit('add-point')" prepend-icon="fa-plus" variant="secondary">
        {{ $t('views.rules.modal.polygonBuilder.addPoint') }}
      </BaseButton>
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
import BaseButton from '@/components/BaseButton.vue'

interface Props {
  config: {
    latField: string
    lonField: string
    points: Array<[number, number]>
  }
  fieldOptions: Array<{ value: string; label: string }>
  expression: string
}

defineProps<Props>()

defineEmits<{
  'add-point': []
  'remove-point': [index: number]
}>()
</script>

