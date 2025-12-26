<template>
  <div class="max-w-7xl mx-auto px-12">
    <div class="mb-10">
      <h2 class="text-3xl font-bold text-slate-900 mb-2">System Health</h2>
      <p class="text-slate-600 font-medium">Monitor system status and availability</p>
    </div>

    <LoadingSpinner v-if="loading" text="Loading health status..." :centered="false" />

    <ErrorMessage
      v-else-if="error"
      title="Error loading health status"
      :message="error"
      retryable
      @retry="loadHealth"
    />

    <div v-else class="bg-white rounded-xl border-2 border-slate-200 shadow-sm p-8">
      <div class="flex items-center gap-3 mb-6">
        <div
          :class="[
            'w-12 h-12 rounded-full flex items-center justify-center',
            healthData.status === 'ok' ? 'bg-emerald-100' : 'bg-orange-100'
          ]"
        >
          <svg
            :class="[
              'w-6 h-6',
              healthData.status === 'ok' ? 'text-emerald-600' : 'text-orange-600'
            ]"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
        </div>
        <div>
          <h3 class="text-xl font-bold text-slate-900">System Status</h3>
          <p class="text-sm text-slate-600">
            {{ healthData.status === 'ok' ? 'All systems operational' : 'System degraded' }}
          </p>
        </div>
      </div>

      <div class="space-y-5">
        <div class="flex items-center justify-between py-4 px-5 bg-slate-50 rounded-lg">
          <div class="flex items-center gap-3">
            <div
              :class="[
                'w-2 h-2 rounded-full',
                healthData.status === 'ok' ? 'bg-emerald-500' : 'bg-orange-500'
              ]"
            ></div>
            <span class="font-semibold text-slate-700">Status</span>
          </div>
          <span
            :class="[
              'text-sm font-medium uppercase',
              healthData.status === 'ok' ? 'text-emerald-600' : 'text-orange-600'
            ]"
          >
            {{ healthData.status }}
          </span>
        </div>

        <div class="flex items-center justify-between py-4 px-5 bg-slate-50 rounded-lg">
          <div class="flex items-center gap-3">
            <div
              :class="[
                'w-2 h-2 rounded-full',
                healthData.postgres === 'healthy' ? 'bg-emerald-500' : 'bg-red-500'
              ]"
            ></div>
            <span class="font-semibold text-slate-700">PostgreSQL</span>
          </div>
          <span
            :class="[
              'text-sm font-medium uppercase',
              healthData.postgres === 'healthy' ? 'text-emerald-600' : 'text-red-600'
            ]"
          >
            {{ healthData.postgres }}
          </span>
        </div>

        <div class="flex items-center justify-between py-4 px-5 bg-slate-50 rounded-lg">
          <div class="flex items-center gap-3">
            <div
              :class="[
                'w-2 h-2 rounded-full',
                healthData.redis === 'healthy' ? 'bg-emerald-500' : 'bg-red-500'
              ]"
            ></div>
            <span class="font-semibold text-slate-700">Redis</span>
          </div>
          <span
            :class="[
              'text-sm font-medium uppercase',
              healthData.redis === 'healthy' ? 'text-emerald-600' : 'text-red-600'
            ]"
          >
            {{ healthData.redis }}
          </span>
        </div>

        <div class="flex items-center justify-between py-4 px-5 bg-slate-50 rounded-lg">
          <div class="flex items-center gap-3">
            <div class="w-2 h-2 rounded-full bg-slate-400"></div>
            <span class="font-semibold text-slate-700">Timestamp</span>
          </div>
          <span class="text-sm font-medium text-slate-600">{{ formatTimestamp(healthData.timestamp) }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '@/lib/api'
import LoadingSpinner from '@/components/LoadingSpinner.vue'
import ErrorMessage from '@/components/ErrorMessage.vue'

interface HealthData {
  status: string
  timestamp: string
  postgres: string
  redis: string
}

const loading = ref(true)
const error = ref('')
const healthData = ref<HealthData>({
  status: 'ok',
  timestamp: '',
  postgres: 'unknown',
  redis: 'unknown',
})

const loadHealth = async () => {
  loading.value = true
  error.value = ''
  try {
    const response = await api.get<HealthData>('/health')
    healthData.value = {
      status: response.status || 'ok',
      timestamp: response.timestamp || new Date().toISOString(),
      postgres: response.postgres || 'unknown',
      redis: response.redis || 'unknown',
    }
  } catch (e: any) {
    error.value = e.message || 'Failed to load health status'
    console.error('Error loading health:', e)
  } finally {
    loading.value = false
  }
}

const formatTimestamp = (timestamp: string) => {
  if (!timestamp) return 'N/A'
  try {
    return new Date(timestamp).toLocaleString('en-US', {
      year: 'numeric',
      month: 'short',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
    })
  } catch {
    return timestamp
  }
}

onMounted(() => {
  loadHealth()
})
</script>
