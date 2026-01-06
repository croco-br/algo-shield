<template>
  <BaseModal
    :model-value="modelValue"
    :title="isEditing ? 'Edit Rule' : 'Create New Rule'"
    size="lg"
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <v-form ref="formRef" @submit.prevent="handleSubmit" class="mt-4">
      <!-- Presets Section (only for new rules) -->
      <div v-if="!isEditing" class="mb-6">
        <label class="text-body-2 text-grey-darken-1 d-block mb-2">Quick Start with Preset</label>
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
      <div v-if="currentSchema && currentSchema.extracted_fields?.length > 0" class="mb-4 pa-3 bg-blue-lighten-5 rounded-lg">
        <div class="d-flex justify-space-between align-center mb-2">
          <label class="text-caption text-grey-darken-1 d-flex align-center">
            <v-icon icon="fa-info-circle" size="x-small" class="mr-1" />
            Available fields from "{{ currentSchema.name }}" ({{ currentSchema.extracted_fields.length }}):
          </label>
          <v-btn
            v-if="currentSchema.extracted_fields.length > 10"
            size="x-small"
            variant="text"
            @click="showAllFields = !showAllFields"
          >
            {{ showAllFields ? 'Show Less' : `Show All (${currentSchema.extracted_fields.length})` }}
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
        v-model.number="editingRule.priority"
        type="number"
        label="Priority"
        :min="0"
        :max="100"
        required
        :rules="[
          (v: any) => (v !== null && v !== undefined && v !== '') || 'Priority is required',
          (v: any) => (typeof v === 'number' && v >= 0 && v <= 100) || 'Priority must be between 0 and 100'
        ]"
        hint="0 = highest priority, 100 = lowest priority"
        persistent-hint
        prepend-inner-icon="fa-sort"
        class="mb-4"
      />

      <!-- Rule Expression Section -->
      <slot name="expression-builder" />

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
      <BaseButton variant="ghost" @click="$emit('cancel')" prepend-icon="fa-xmark">Cancel</BaseButton>
      <BaseButton @click="handleSubmit" :loading="saving" prepend-icon="fa-save">Save</BaseButton>
    </template>
  </BaseModal>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import BaseModal from '@/components/BaseModal.vue'
import BaseInput from '@/components/BaseInput.vue'
import BaseSelect from '@/components/BaseSelect.vue'
import BaseButton from '@/components/BaseButton.vue'

interface Props {
  modelValue: boolean
  isEditing: boolean
  editingRule: any
  schemaOptions: Array<{ value: string; label: string }>
  currentSchema: any
  rulePresets: any[]
  ruleActions: Array<{ value: string; label: string }>
  saving: boolean
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

async function handleSubmit() {
  if (formRef.value) {
    const { valid } = await formRef.value.validate()
    if (valid) {
      emit('submit')
    }
  }
}
</script>

