<template>
  <div class="max-w-7xl mx-auto px-12">
    <div class="flex justify-between items-center mb-10">
      <div>
        <h2 class="text-3xl font-bold text-slate-900 mb-2">Rules Management</h2>
        <p class="text-slate-600 font-medium">Configure custom rules for fraud detection and AML</p>
      </div>
      <BaseButton @click="openCreateModal">
        <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
        </svg>
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
        <div class="font-semibold text-slate-900">{{ row.name }}</div>
        <div class="text-sm text-slate-500">{{ row.description }}</div>
      </template>

      <template #cell-type="{ value }">
        <span class="text-slate-700 font-medium">{{ value }}</span>
      </template>

      <template #cell-action="{ row }">
        <BaseBadge :variant="getActionBadgeVariant(row.action)" rounded>
          {{ row.action }}
        </BaseBadge>
      </template>

      <template #cell-score="{ value }">
        <span class="text-slate-700 font-medium">{{ value }}</span>
      </template>

      <template #cell-priority="{ value }">
        <span class="text-slate-700 font-medium">{{ value }}</span>
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
        <div class="flex gap-2">
          <BaseButton size="sm" @click="openEditModal(row)">
            Edit
          </BaseButton>
          <BaseButton size="sm" variant="danger" @click="deleteRule(row.id)">
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
      <div class="space-y-8">
        <BaseInput
          id="name"
          v-model="editingRule.name"
          label="Name"
          placeholder="Rule name"
        />

        <BaseInput
          id="description"
          v-model="editingRule.description"
          label="Description"
          placeholder="Description"
        />

        <BaseSelect
          id="type"
          v-model="editingRule.type"
          label="Type"
          :options="ruleTypes"
        />

        <BaseSelect
          id="action"
          v-model="editingRule.action"
          label="Action"
          :options="ruleActions"
        />

        <BaseInput
          id="score"
          v-model.number="editingRule.score"
          type="number"
          label="Risk Score"
          placeholder="0-100"
        />

        <BaseInput
          id="priority"
          v-model.number="editingRule.priority"
          type="number"
          label="Priority"
          placeholder="Lower number = higher priority"
        />
      </div>

      <template #footer>
        <BaseButton variant="secondary" @click="closeModal">
          Cancel
        </BaseButton>
        <BaseButton @click="saveRule">
          {{ isEditing ? 'Update' : 'Create' }}
        </BaseButton>
      </template>
    </BaseModal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '@/lib/api'
import BaseButton from '@/components/BaseButton.vue'
import BaseBadge from '@/components/BaseBadge.vue'
import BaseInput from '@/components/BaseInput.vue'
import BaseSelect from '@/components/BaseSelect.vue'
import BaseModal from '@/components/BaseModal.vue'
import BaseTable from '@/components/BaseTable.vue'
import LoadingSpinner from '@/components/LoadingSpinner.vue'
import ErrorMessage from '@/components/ErrorMessage.vue'
import { useActionBadge } from '@/composables/useActionBadge'

interface Rule {
	id: string;
	name: string;
	description: string;
	type: string;
	action: string;
	priority: number;
	enabled: boolean;
	conditions: any;
	score: number;
	created_at: string;
	updated_at: string;
}

const tableColumns = [
	{ key: 'name', label: 'Name' },
	{ key: 'type', label: 'Type' },
	{ key: 'action', label: 'Action' },
	{ key: 'score', label: 'Score' },
	{ key: 'priority', label: 'Priority' },
	{ key: 'status', label: 'Status' },
	{ key: 'actions', label: 'Actions' },
]

const ruleTypes = [
	{ value: 'amount', label: 'Amount Threshold' },
	{ value: 'velocity', label: 'Transaction Velocity' },
	{ value: 'blocklist', label: 'Blocklist' },
	{ value: 'pattern', label: 'Pattern Match' },
	{ value: 'custom', label: 'Custom' },
];

const ruleActions = [
	{ value: 'allow', label: 'Allow' },
	{ value: 'block', label: 'Block' },
	{ value: 'review', label: 'Review' },
	{ value: 'score', label: 'Add Score' },
];

const rules = ref<Rule[]>([])
const loading = ref(true)
const error = ref<string | null>(null)
const showModal = ref(false)
const editingRule = ref<Partial<Rule>>({})
const isEditing = ref(false)

onMounted(() => {
	loadRules()
})

async function loadRules() {
	loading.value = true
	error.value = null
	try {
		const data = await api.get<{ rules: Rule[] }>('/api/v1/rules')
		rules.value = data.rules || []
	} catch (err) {
		const errorMessage = err instanceof Error ? err.message : 'Failed to load rules'
		console.error('Failed to load rules:', err)
		error.value = errorMessage
	} finally {
		loading.value = false
	}
}

function openCreateModal() {
	editingRule.value = {
		name: '',
		description: '',
		type: 'amount',
		action: 'score',
		priority: 10,
		enabled: true,
		conditions: {},
		score: 0,
	}
	isEditing.value = false
	showModal.value = true
}

function openEditModal(rule: Rule) {
	editingRule.value = { ...rule }
	isEditing.value = true
	showModal.value = true
}

function closeModal() {
	showModal.value = false
	editingRule.value = {}
}

async function saveRule() {
	try {
		if (isEditing.value) {
			await api.put(`/api/v1/rules/${editingRule.value.id}`, editingRule.value)
		} else {
			await api.post('/api/v1/rules', editingRule.value)
		}
		closeModal()
		loadRules()
	} catch (err) {
		const errorMessage = err instanceof Error ? err.message : 'Failed to save rule'
		console.error('Failed to save rule:', errorMessage)
		alert(`Error: ${errorMessage}`)
	}
}

async function deleteRule(id: string) {
	if (!confirm('Are you sure you want to delete this rule?')) return

	try {
		await api.delete(`/api/v1/rules/${id}`)
		loadRules()
	} catch (err) {
		const errorMessage = err instanceof Error ? err.message : 'Failed to delete rule'
		console.error('Failed to delete rule:', errorMessage)
		alert(`Error: ${errorMessage}`)
	}
}

async function toggleRule(rule: Rule) {
	try {
		await api.put(`/api/v1/rules/${rule.id}`, { ...rule, enabled: !rule.enabled })
		loadRules()
	} catch (err) {
		const errorMessage = err instanceof Error ? err.message : 'Failed to toggle rule'
		console.error('Failed to toggle rule:', errorMessage)
		alert(`Error: ${errorMessage}`)
	}
}

function getActionBadgeVariant(action: string): 'success' | 'warning' | 'danger' | 'info' | 'default' {
	const { variant } = useActionBadge(action)
	return variant.value
}
</script>
