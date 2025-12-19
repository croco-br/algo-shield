<template>
  <div class="max-w-7xl mx-auto px-8">
    <div class="flex justify-between items-center mb-8">
      <div>
        <h2 class="text-3xl font-semibold mb-2">Rules Management</h2>
        <p class="text-gray-500">Configure custom rules for fraud detection and AML</p>
      </div>
      <button
        @click="openCreateModal"
        class="px-6 py-3 bg-indigo-600 text-white rounded-md font-medium hover:bg-indigo-700 transition-colors"
      >
        + Create Rule
      </button>
    </div>

    <div v-if="loading" class="bg-white rounded-lg border border-gray-200 p-8">
      <p>Loading rules...</p>
    </div>

    <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-lg p-8">
      <div class="flex items-start">
        <div class="flex-shrink-0">
          <svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
          </svg>
        </div>
        <div class="ml-3">
          <h3 class="text-sm font-medium text-red-800">Error loading rules</h3>
          <div class="mt-2 text-sm text-red-700">
            <p>{{ error }}</p>
          </div>
          <div class="mt-4">
            <button
              @click="loadRules"
              class="text-sm font-medium text-red-800 hover:text-red-900 underline"
            >
              Try again
            </button>
          </div>
        </div>
      </div>
    </div>

    <div v-else-if="rules.length === 0" class="bg-white rounded-lg border border-gray-200 p-8">
      <p>No rules configured. Create your first rule to get started.</p>
    </div>

    <div v-else class="bg-white rounded-lg border border-gray-200 overflow-hidden">
      <table class="w-full border-collapse">
        <thead>
          <tr class="bg-gray-50">
            <th class="px-4 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider border-b border-gray-200">
              Name
            </th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider border-b border-gray-200">
              Type
            </th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider border-b border-gray-200">
              Action
            </th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider border-b border-gray-200">
              Score
            </th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider border-b border-gray-200">
              Priority
            </th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider border-b border-gray-200">
              Status
            </th>
            <th class="px-4 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider border-b border-gray-200">
              Actions
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="rule in rules" :key="rule.id" class="hover:bg-gray-50">
            <td class="px-4 py-4 border-b border-gray-200">
              <div class="font-semibold text-gray-900">{{ rule.name }}</div>
              <div class="text-sm text-gray-500">{{ rule.description }}</div>
            </td>
            <td class="px-4 py-4 border-b border-gray-200 text-gray-900">{{ rule.type }}</td>
            <td class="px-4 py-4 border-b border-gray-200">
              <span :class="['px-3 py-1 rounded-full text-xs font-medium', getActionBadgeClass(rule.action)]">
                {{ rule.action }}
              </span>
            </td>
            <td class="px-4 py-4 border-b border-gray-200 text-gray-900">{{ rule.score }}</td>
            <td class="px-4 py-4 border-b border-gray-200 text-gray-900">{{ rule.priority }}</td>
            <td class="px-4 py-4 border-b border-gray-200">
              <button
                @click="toggleRule(rule)"
                :class="[
                  'px-3 py-1 rounded text-sm border',
                  rule.enabled
                    ? 'bg-green-500 text-white border-green-500'
                    : 'bg-white text-gray-700 border-gray-300'
                ]"
              >
                {{ rule.enabled ? 'Enabled' : 'Disabled' }}
              </button>
            </td>
            <td class="px-4 py-4 border-b border-gray-200">
              <button
                @click="openEditModal(rule)"
                class="px-4 py-2 mr-2 bg-indigo-600 text-white rounded text-sm hover:bg-indigo-700 transition-colors"
              >
                Edit
              </button>
              <button
                @click="deleteRule(rule.id)"
                class="px-4 py-2 bg-red-600 text-white rounded text-sm hover:bg-red-700 transition-colors"
              >
                Delete
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div
      v-if="showModal"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
      @click="closeModal"
    >
      <div
        class="bg-white rounded-lg p-8 max-w-2xl w-full mx-4 max-h-[90vh] overflow-y-auto"
        @click.stop
      >
        <h3 class="text-2xl font-semibold mb-6">{{ isEditing ? 'Edit Rule' : 'Create New Rule' }}</h3>

        <div class="space-y-6">
          <div>
            <label for="name" class="block text-sm font-medium text-gray-700 mb-2">
              Name
            </label>
            <input
              id="name"
              type="text"
              v-model="editingRule.name"
              placeholder="Rule name"
              class="w-full px-4 py-3 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
            />
          </div>

          <div>
            <label for="description" class="block text-sm font-medium text-gray-700 mb-2">
              Description
            </label>
            <input
              id="description"
              type="text"
              v-model="editingRule.description"
              placeholder="Description"
              class="w-full px-4 py-3 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
            />
          </div>

          <div>
            <label for="type" class="block text-sm font-medium text-gray-700 mb-2">
              Type
            </label>
            <select
              id="type"
              v-model="editingRule.type"
              class="w-full px-4 py-3 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
            >
              <option v-for="type in ruleTypes" :key="type.value" :value="type.value">
                {{ type.label }}
              </option>
            </select>
          </div>

          <div>
            <label for="action" class="block text-sm font-medium text-gray-700 mb-2">
              Action
            </label>
            <select
              id="action"
              v-model="editingRule.action"
              class="w-full px-4 py-3 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
            >
              <option v-for="action in ruleActions" :key="action.value" :value="action.value">
                {{ action.label }}
              </option>
            </select>
          </div>

          <div>
            <label for="score" class="block text-sm font-medium text-gray-700 mb-2">
              Risk Score
            </label>
            <input
              id="score"
              type="number"
              v-model.number="editingRule.score"
              placeholder="0-100"
              class="w-full px-4 py-3 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
            />
          </div>

          <div>
            <label for="priority" class="block text-sm font-medium text-gray-700 mb-2">
              Priority
            </label>
            <input
              id="priority"
              type="number"
              v-model.number="editingRule.priority"
              placeholder="Lower number = higher priority"
              class="w-full px-4 py-3 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
            />
          </div>
        </div>

        <div class="flex justify-end gap-4 mt-8">
          <button
            @click="closeModal"
            class="px-6 py-3 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50 transition-colors"
          >
            Cancel
          </button>
          <button
            @click="saveRule"
            class="px-6 py-3 bg-indigo-600 text-white rounded-md hover:bg-indigo-700 transition-colors"
          >
            {{ isEditing ? 'Update' : 'Create' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '@/lib/api'

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

function getActionBadgeClass(action: string) {
	switch (action) {
		case 'block':
			return 'bg-red-100 text-red-800'
		case 'review':
			return 'bg-yellow-100 text-yellow-800'
		case 'allow':
			return 'bg-green-100 text-green-800'
		default:
			return 'bg-gray-100 text-gray-800'
	}
}
</script>
