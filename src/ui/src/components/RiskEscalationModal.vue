<template>
  <div
    v-if="modelValue"
    class="modal-overlay fixed inset-0 bg-black bg-opacity-50 z-modal-backdrop flex items-center justify-center p-4"
    @click.self="close"
  >
    <div
      class="modal-content bg-white rounded-lg shadow-2xl w-full max-w-[400px] animate-fade-in z-modal"
    >
      <!-- Header -->
      <div class="modal-header px-6 py-4 border-b border-neutral-200 flex items-center justify-between">
        <h2 class="text-lg font-bold text-neutral-900">Escalate Transaction</h2>
        <v-btn
          @click="close"
          icon
          variant="text"
          size="small"
        >
          <v-icon icon="fa-xmark" />
        </v-btn>
      </div>

      <!-- Content -->
      <form @submit.prevent="handleSubmit" class="modal-body px-8 py-8 space-y-8">
        <!-- Transaction ID -->
        <div>
          <label class="block text-sm font-semibold text-neutral-700 mb-2">
            Transaction ID
          </label>
          <input
            type="text"
            :value="transaction?.external_id || transaction?.id"
            disabled
            class="w-full px-4 py-2 bg-neutral-100 border border-neutral-200 rounded-lg text-sm text-neutral-600 cursor-not-allowed"
          />
        </div>

        <!-- Escalation Level -->
        <div>
          <label class="block text-sm font-semibold text-neutral-700 mb-2">
            Escalation Level
            <span class="text-red-500">*</span>
          </label>
          <select
            v-model="form.level"
            required
            class="w-full px-4 py-2 border border-neutral-200 rounded-lg text-sm text-neutral-900 focus:outline-none focus:ring-2 focus:ring-teal-500 focus:border-transparent transition-all"
          >
            <option value="">Select level</option>
            <option value="tier1">Tier 1 - Senior Analyst</option>
            <option value="tier2">Tier 2 - Manager</option>
            <option value="tier3">Tier 3 - Director</option>
            <option value="tier4">Tier 4 - Compliance Officer</option>
          </select>
        </div>

        <!-- Priority -->
        <div>
          <label class="block text-sm font-semibold text-neutral-700 mb-2">
            Priority
            <span class="text-red-500">*</span>
          </label>
          <div class="flex gap-2">
            <button
              v-for="priority in priorities"
              :key="priority.value"
              type="button"
              @click="form.priority = priority.value"
              :class="[
                'flex-1 px-3 py-2 rounded-lg text-sm font-semibold transition-all',
                form.priority === priority.value
                  ? priority.activeClass
                  : 'bg-neutral-100 text-neutral-600 hover:bg-neutral-200'
              ]"
            >
              {{ priority.label }}
            </button>
          </div>
        </div>

        <!-- Comments -->
        <div>
          <label class="block text-sm font-semibold text-neutral-700 mb-2">
            Comments
            <span class="text-red-500">*</span>
          </label>
          <textarea
            v-model="form.comments"
            required
            rows="4"
            placeholder="Provide detailed reason for escalation..."
            class="w-full px-4 py-2 border border-neutral-200 rounded-lg text-sm text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-teal-500 focus:border-transparent transition-all resize-none"
          ></textarea>
          <p class="text-xs text-neutral-500 mt-1">Minimum 20 characters required</p>
        </div>

        <!-- Actions -->
        <div class="flex items-center gap-3 pt-2">
          <button
            type="button"
            @click="close"
            class="flex-1 px-4 py-2.5 text-sm font-semibold text-neutral-700 hover:bg-neutral-100 rounded-lg transition-colors"
          >
            Cancel
          </button>
          <button
            type="submit"
            :disabled="!isFormValid"
            class="flex-1 px-4 py-2.5 bg-teal-600 text-white text-sm font-semibold rounded-lg hover:bg-teal-700 disabled:opacity-50 disabled:cursor-not-allowed transition-all"
          >
            <v-icon icon="fa-paper-plane" class="mr-2" />
            Submit
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import type { Transaction } from '@/types/transaction'

interface Props {
  modelValue: boolean
  transaction: Transaction | null
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  'submit': [data: any]
}>()

const form = ref({
  level: '',
  priority: 'medium',
  comments: '',
})

const priorities = [
  { value: 'low', label: 'Low', activeClass: 'bg-blue-100 text-blue-700' },
  { value: 'medium', label: 'Medium', activeClass: 'bg-orange-100 text-orange-700' },
  { value: 'high', label: 'High', activeClass: 'bg-red-100 text-red-700' },
]

const isFormValid = computed(() => {
  return form.value.level && form.value.comments.length >= 20
})

const close = () => {
  emit('update:modelValue', false)
  // Reset form
  form.value = {
    level: '',
    priority: 'medium',
    comments: '',
  }
}

const handleSubmit = () => {
  if (!isFormValid.value) return

  emit('submit', {
    transactionId: props.transaction?.id,
    ...form.value,
  })

  close()
}
</script>

<style scoped>
.modal-overlay {
  animation: fadeIn 200ms ease-out;
}

.animate-fade-in {
  animation: slideDown 200ms ease-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-20px) scale(0.95);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}
</style>
