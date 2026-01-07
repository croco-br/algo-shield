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
        <h2 class="text-lg font-bold text-neutral-900">{{ $t('components.riskEscalation.title') }}</h2>
        <BaseButton variant="ghost" size="sm" @click="close" prepend-icon="fa-xmark" />
      </div>

      <!-- Content -->
      <form @submit.prevent="handleSubmit" class="modal-body px-8 py-8 space-y-8">
        <!-- Transaction ID -->
        <div>
          <label class="block text-sm font-semibold text-neutral-700 mb-2">
            {{ $t('components.riskEscalation.transactionId') }}
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
            {{ $t('components.riskEscalation.escalationLevel') }}
            <span class="text-red-500">*</span>
          </label>
          <select
            v-model="form.level"
            required
            class="w-full px-4 py-2 border border-neutral-200 rounded-lg text-sm text-neutral-900 focus:outline-none focus:ring-2 focus:ring-teal-500 focus:border-transparent transition-all"
          >
            <option value="">{{ $t('components.riskEscalation.selectLevel') }}</option>
            <option value="tier1">{{ $t('components.riskEscalation.tier1') }}</option>
            <option value="tier2">{{ $t('components.riskEscalation.tier2') }}</option>
            <option value="tier3">{{ $t('components.riskEscalation.tier3') }}</option>
            <option value="tier4">{{ $t('components.riskEscalation.tier4') }}</option>
          </select>
        </div>

        <!-- Priority -->
        <div>
          <label class="block text-sm font-semibold text-neutral-700 mb-2">
            {{ $t('components.riskEscalation.priority') }}
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
            {{ $t('components.riskEscalation.comments') }}
            <span class="text-red-500">*</span>
          </label>
          <textarea
            v-model="form.comments"
            required
            rows="4"
            :placeholder="$t('components.riskEscalation.commentsPlaceholder')"
            class="w-full px-4 py-2 border border-neutral-200 rounded-lg text-sm text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-teal-500 focus:border-transparent transition-all resize-none"
          ></textarea>
          <p class="text-xs text-neutral-500 mt-1">{{ $t('components.riskEscalation.commentsHint') }}</p>
        </div>

        <!-- Actions -->
        <div class="flex items-center gap-3 pt-2">
          <BaseButton variant="ghost" @click="close" class="flex-1">
            {{ $t('components.riskEscalation.cancel') }}
          </BaseButton>
          <BaseButton
            type="submit"
            :disabled="!isFormValid"
            prepend-icon="fa-paper-plane"
            class="flex-1"
          >
            {{ $t('components.riskEscalation.submit') }}
          </BaseButton>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { i18n } from '@/plugins/i18n'
import type { Transaction } from '@/types/transaction'
import BaseButton from '@/components/BaseButton.vue'

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

const priorities = computed(() => [
  { value: 'low', label: i18n.global.t('components.riskEscalation.low'), activeClass: 'bg-blue-100 text-blue-700' },
  { value: 'medium', label: i18n.global.t('components.riskEscalation.medium'), activeClass: 'bg-orange-100 text-orange-700' },
  { value: 'high', label: i18n.global.t('components.riskEscalation.high'), activeClass: 'bg-red-100 text-red-700' },
])

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
