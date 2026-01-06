<template>
  <div class="pa-3 bg-white rounded-lg">
    <BaseSelect
      v-model="config.latField"
      label="Latitude Field"
      :options="fieldOptions"
      hint="Select the field containing latitude (e.g., location.lat, latitude)"
      persistent-hint
      class="mb-3"
    />
    <BaseSelect
      v-model="config.lonField"
      label="Longitude Field"
      :options="fieldOptions"
      hint="Select the field containing longitude (e.g., location.lon, longitude)"
      persistent-hint
      class="mb-3"
    />
    
    <div class="mb-3">
      <label class="text-body-2 text-grey-darken-1 d-block mb-2">
        Polygon Coordinates <span class="text-red">*</span>
        <span class="text-caption text-grey">(at least 3 points required)</span>
      </label>
      <div v-for="(point, index) in config.points" :key="index" class="d-flex gap-2 mb-2">
        <BaseInput
          v-model.number="point[0]"
          type="number"
          label="Latitude"
          step="0.000001"
          placeholder="e.g., 37.7749"
          :min="-90"
          :max="90"
          :rules="[
            (v: any) => v !== null && v !== undefined && v !== '' || 'Latitude is required',
            (v: any) => (typeof v === 'number' && v >= -90 && v <= 90) || 'Latitude must be between -90 and 90'
          ]"
          hint="Range: -90 to 90"
          persistent-hint
          style="flex: 1;"
        />
        <BaseInput
          v-model.number="point[1]"
          type="number"
          label="Longitude"
          step="0.000001"
          placeholder="e.g., -122.4194"
          :min="-180"
          :max="180"
          :rules="[
            (v: any) => v !== null && v !== undefined && v !== '' || 'Longitude is required',
            (v: any) => (typeof v === 'number' && v >= -180 && v <= 180) || 'Longitude must be between -180 and 180'
          ]"
          hint="Range: -180 to 180"
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
        Add Point
      </BaseButton>
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

