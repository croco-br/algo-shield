<template>
  <div class="max-w-7xl mx-auto px-8">
    <div class="mb-8">
      <h2 class="text-3xl font-semibold mb-2">Synthetic Data Testing</h2>
      <p class="text-gray-500">Generate and test synthetic events with custom field schemas</p>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
      <!-- Field Configuration Panel -->
      <div class="bg-white rounded-lg border border-gray-200 p-6">
        <h3 class="text-xl font-semibold mb-2">Event Schema</h3>
        <p class="text-sm text-gray-500 mb-6">Define the fields for your synthetic events</p>

        <div class="space-y-3">
          <div
            v-for="(field, index) in fields"
            :key="index"
            class="flex items-center justify-between p-3 bg-gray-50 border border-gray-200 rounded-md"
          >
            <div class="flex-1">
              <div class="font-medium text-gray-900 flex items-center gap-2">
                {{ field.name }}
                <span v-if="field.required" class="px-2 py-0.5 bg-yellow-500 text-white rounded text-xs">Required</span>
              </div>
              <div class="text-sm text-gray-500">{{ field.type }}</div>
            </div>
            <button
              @click="removeField(index)"
              class="text-red-600 text-xl leading-none hover:opacity-80"
              title="Remove field"
            >
              ✕
            </button>
          </div>

          <div v-if="showAddField" class="p-4 bg-gray-50 border-2 border-dashed border-gray-300 rounded-md space-y-3">
            <input
              type="text"
              v-model="newField.name"
              placeholder="Field name"
              class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
            />
            <select
              v-model="newField.type"
              class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
            >
              <option v-for="type in fieldTypes" :key="type.value" :value="type.value">
                {{ type.label }}
              </option>
            </select>
            <label class="flex items-center gap-2 text-sm">
              <input type="checkbox" v-model="newField.required" />
              Required
            </label>
            <div class="flex gap-2">
              <button
                @click="addField"
                class="px-4 py-2 bg-indigo-600 text-white rounded text-sm hover:bg-indigo-700 transition-colors"
              >
                Add
              </button>
              <button
                @click="showAddField = false"
                class="px-4 py-2 border border-gray-300 rounded text-sm hover:bg-gray-50 transition-colors"
              >
                Cancel
              </button>
            </div>
          </div>

          <button
            v-else
            @click="showAddField = true"
            class="w-full px-4 py-2 border border-dashed border-gray-300 rounded text-sm text-gray-500 hover:border-indigo-500 hover:text-indigo-600 transition-colors"
          >
            + Add Field
          </button>
        </div>
      </div>

      <!-- Test Configuration Panel -->
      <div class="bg-white rounded-lg border border-gray-200 p-6">
        <h3 class="text-xl font-semibold mb-2">Test Configuration</h3>
        <p class="text-sm text-gray-500 mb-6">Configure and run your tests</p>

        <div class="mb-6">
          <label for="numEvents" class="block text-sm font-medium text-gray-700 mb-2">
            Number of Events
          </label>
          <input
            id="numEvents"
            type="number"
            min="1"
            max="1000"
            v-model.number="numberOfEvents"
            :disabled="loading"
            class="w-full px-4 py-3 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 disabled:opacity-60"
          />
        </div>

        <div v-if="loading" class="space-y-4">
          <div class="mb-4 p-4 bg-gray-50 border border-gray-200 rounded-md">
            <div class="flex justify-between items-center mb-3 text-sm">
              <span>Processing events...</span>
              <span class="font-semibold text-indigo-600">
                {{ currentProgress }} / {{ numberOfEvents }}
              </span>
            </div>
            <div class="w-full h-2 bg-gray-200 rounded-full overflow-hidden">
              <div
                class="h-full bg-gradient-to-r from-indigo-600 to-purple-600 transition-all duration-300"
                :style="{ width: `${(currentProgress / numberOfEvents) * 100}%` }"
              ></div>
            </div>
          </div>

          <button
            @click="stopTest"
            class="w-full py-3 px-4 bg-red-600 text-white rounded-md font-medium hover:bg-red-700 transition-colors"
          >
            ⏹ Stop Test
          </button>
        </div>

        <button
          v-else
          @click="runTest"
          :disabled="fields.length === 0"
          class="w-full py-3 px-4 bg-indigo-600 text-white rounded-md font-medium hover:bg-indigo-700 transition-colors disabled:opacity-60 disabled:cursor-not-allowed"
        >
          ▶ Run Test
        </button>

        <div v-if="testResults.length > 0" class="mt-6">
          <div class="grid grid-cols-2 gap-4 p-6 bg-gray-50 rounded-md">
            <div class="text-center">
              <div class="text-3xl font-bold text-indigo-600">{{ testResults.length }}</div>
              <div class="text-sm text-gray-500 mt-1">Events Tested</div>
            </div>
            <div class="text-center">
              <div class="text-3xl font-bold text-indigo-600">
                {{ testResults.filter((r) => r.action === 'block').length }}
              </div>
              <div class="text-sm text-gray-500 mt-1">Blocked</div>
            </div>
            <div class="text-center">
              <div class="text-3xl font-bold text-indigo-600">
                {{ testResults.filter((r) => r.action === 'review').length }}
              </div>
              <div class="text-sm text-gray-500 mt-1">Review</div>
            </div>
            <div class="text-center">
              <div class="text-3xl font-bold text-indigo-600">
                {{ testResults.filter((r) => r.action === 'allow').length }}
              </div>
              <div class="text-sm text-gray-500 mt-1">Allowed</div>
            </div>
          </div>

          <div class="flex gap-2 mt-4">
            <button
              @click="exportResults"
              class="flex-1 px-4 py-2 border border-gray-300 rounded text-sm hover:bg-gray-50 transition-colors"
            >
              Export Results
            </button>
            <button
              @click="clearResults"
              class="flex-1 px-4 py-2 border border-gray-300 rounded text-sm hover:bg-gray-50 transition-colors"
            >
              Clear
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Results Panel -->
    <div v-if="testResults.length > 0" class="bg-white rounded-lg border border-gray-200 p-6">
      <h3 class="text-xl font-semibold mb-6">Test Results</h3>

      <div class="space-y-4">
        <div
          v-for="(result, index) in testResults"
          :key="index"
          class="border border-gray-200 rounded-md p-4"
        >
          <div class="flex items-center gap-4 mb-2">
            <span class="font-bold text-gray-500">#{{ index + 1 }}</span>
            <span class="text-sm text-gray-500">
              {{ new Date(result.timestamp).toLocaleTimeString() }}
            </span>
            <span
              :class="[
                'px-3 py-1 rounded-full text-xs font-medium',
                result.action === 'block'
                  ? 'bg-red-500 text-white'
                  : result.action === 'review'
                    ? 'bg-yellow-500 text-white'
                    : 'bg-green-500 text-white'
              ]"
            >
              {{ result.action }}
            </span>
            <span class="ml-auto font-medium">Score: {{ result.score }}</span>
          </div>
          <details class="mt-3">
            <summary class="cursor-pointer text-sm text-indigo-600 hover:text-indigo-700 select-none">
              View Event Data
            </summary>
            <pre class="mt-3 p-4 bg-gray-50 rounded border border-gray-200 overflow-x-auto text-sm">
              {{ JSON.stringify(result.data, null, 2) }}
            </pre>
          </details>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

interface Field {
	name: string;
	type: 'string' | 'number' | 'boolean' | 'date' | 'uuid' | 'email' | 'ip';
	required: boolean;
	format?: string;
}

interface TestResult {
	timestamp: string;
	data: any;
	score: number;
	action: string;
	status: 'success' | 'error';
}

const fieldTypes = [
	{ value: 'string', label: 'String' },
	{ value: 'number', label: 'Number' },
	{ value: 'boolean', label: 'Boolean' },
	{ value: 'date', label: 'Date/Time' },
	{ value: 'uuid', label: 'UUID' },
	{ value: 'email', label: 'Email' },
	{ value: 'ip', label: 'IP Address' },
];

const fields = ref<Field[]>([
	{ name: 'transaction_id', type: 'uuid', required: true },
	{ name: 'amount', type: 'number', required: true },
	{ name: 'user_email', type: 'email', required: true },
])
const testResults = ref<TestResult[]>([])
const loading = ref(false)
const shouldStop = ref(false)
const currentProgress = ref(0)
const numberOfEvents = ref(1)
const showAddField = ref(false)
const newField = ref<Partial<Field>>({ name: '', type: 'string', required: true })

function addField() {
	if (newField.value.name && newField.value.type) {
		fields.value.push(newField.value as Field)
		newField.value = { name: '', type: 'string', required: true }
		showAddField.value = false
	}
}

function removeField(index: number) {
	fields.value = fields.value.filter((_, i) => i !== index)
}

function generateValue(field: Field): any {
	switch (field.type) {
		case 'uuid':
			return crypto.randomUUID()
		case 'email':
			return `user${Math.floor(Math.random() * 10000)}@example.com`
		case 'number':
			return Math.floor(Math.random() * 10000) + 100
		case 'boolean':
			return Math.random() > 0.5
		case 'date':
			return new Date().toISOString()
		case 'ip':
			return `${Math.floor(Math.random() * 256)}.${Math.floor(Math.random() * 256)}.${Math.floor(Math.random() * 256)}.${Math.floor(Math.random() * 256)}`
		case 'string':
		default:
			return `value_${Math.random().toString(36).substring(7)}`
	}
}

function generateSyntheticData(): any {
	const data: any = {}
	fields.value.forEach((field) => {
		data[field.name] = generateValue(field)
	})
	return data
}

async function runTest() {
	loading.value = true
	shouldStop.value = false
	currentProgress.value = 0
	testResults.value = []

	try {
		for (let i = 0; i < numberOfEvents.value; i++) {
			if (shouldStop.value) {
				console.log('Test stopped by user')
				break
			}

			const syntheticData = generateSyntheticData()

			// Simulate API call
			await new Promise((resolve) => setTimeout(resolve, 300))

			const mockScore = Math.floor(Math.random() * 100)
			const mockAction = mockScore > 80 ? 'block' : mockScore > 50 ? 'review' : 'allow'

			testResults.value.push({
				timestamp: new Date().toISOString(),
				data: syntheticData,
				score: mockScore,
				action: mockAction,
				status: 'success',
			})

			currentProgress.value = i + 1
		}
	} catch (error) {
		console.error('Test failed:', error)
	}

	loading.value = false
	currentProgress.value = 0
}

function stopTest() {
	shouldStop.value = true
}

function exportResults() {
	const dataStr = JSON.stringify(testResults.value, null, 2)
	const blob = new Blob([dataStr], { type: 'application/json' })
	const url = URL.createObjectURL(blob)
	const link = document.createElement('a')
	link.href = url
	link.download = `synthetic-test-${Date.now()}.json`
	link.click()
	URL.revokeObjectURL(url)
}

function clearResults() {
	testResults.value = []
}
</script>
