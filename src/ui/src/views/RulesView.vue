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
      <BaseButton @click="openCreateModal" prepend-icon="fa-plus" style="text-transform: none; font-family: var(--font-family-sans);">
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
      <v-form @submit.prevent="handleSubmit" class="mt-4">
        <v-text-field
          v-model="editingRule.name"
          label="Name"
          placeholder="Rule name"
          required
          prepend-inner-icon="fa-text"
          variant="outlined"
          class="mb-4"
        />

        <v-text-field
          v-model="editingRule.description"
          label="Description"
          placeholder="Description"
          prepend-inner-icon="fa-align-left"
          variant="outlined"
          class="mb-4"
        />

        <BaseSelect
          v-model="editingRule.type"
          label="Type"
          :options="ruleTypes"
          required
          class="mb-4"
        />

        <BaseSelect
          v-model="editingRule.action"
          label="Action"
          :options="ruleActions"
          required
          class="mb-4"
        />

        <v-text-field
          v-model.number="editingRule.score"
          type="number"
          label="Score"
          min="0"
          max="100"
          required
          prepend-inner-icon="fa-hashtag"
          variant="outlined"
          class="mb-4"
        />

        <BaseSelect
          v-model="editingRule.priority"
          label="Priority"
          :options="rulePriorities"
          required
          class="mb-4"
        />

        <v-switch
          v-model="editingRule.enabled"
          label="Enabled"
          class="mb-6"
        />
      </v-form>

      <template #footer>
        <v-btn variant="text" @click="closeModal" prepend-icon="fa-xmark">Cancel</v-btn>
        <v-btn @click="handleSubmit" color="primary" :loading="saving" prepend-icon="fa-save">Save</v-btn>
      </template>
    </BaseModal>
  </v-container>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { api } from '@/lib/api'
import BaseButton from '@/components/BaseButton.vue'
import BaseBadge from '@/components/BaseBadge.vue'
import BaseModal from '@/components/BaseModal.vue'
import BaseTable from '@/components/BaseTable.vue'
import BaseSelect from '@/components/BaseSelect.vue'
import LoadingSpinner from '@/components/LoadingSpinner.vue'
import ErrorMessage from '@/components/ErrorMessage.vue'

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

const editingRule = reactive<{
  id?: string
  name: string
  description: string
  type: string
  action: string
  score: number
  priority: string
  enabled: boolean
}>({
  name: '',
  description: '',
  type: '',
  action: '',
  score: 0,
  priority: '',
  enabled: true,
})

const ruleTypes = [
  { value: 'fraud', label: 'Fraud' },
  { value: 'aml', label: 'AML' },
  { value: 'risk', label: 'Risk' },
]

const ruleActions = [
  { value: 'block', label: 'Block' },
  { value: 'flag', label: 'Flag' },
  { value: 'monitor', label: 'Monitor' },
]

const rulePriorities = [
  { value: 'low', label: 'Low' },
  { value: 'medium', label: 'Medium' },
  { value: 'high', label: 'High' },
]

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
    const response = await api.get<any[]>('/api/v1/rules')
    rules.value = response || []
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
  editingRule.type = 'fraud'
  editingRule.action = 'flag'
  editingRule.score = 0
  editingRule.priority = 'medium'
  editingRule.enabled = true
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
  showModal.value = true
}

function closeModal() {
  showModal.value = false
}

async function handleSubmit() {
  try {
    saving.value = true
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

function getActionBadgeVariant(action: string): 'success' | 'warning' | 'danger' | 'info' | 'default' {
  switch (action.toLowerCase()) {
    case 'block':
      return 'danger'
    case 'flag':
      return 'warning'
    case 'monitor':
      return 'info'
    default:
      return 'default'
  }
}
</script>
