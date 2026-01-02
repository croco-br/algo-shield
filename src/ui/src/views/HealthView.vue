<template>
  <v-container fluid class="pa-8">
    <div class="mb-10">
      <h2 class="text-h4 font-weight-bold mb-2">System Health</h2>
      <p class="text-body-1 text-grey-darken-1">Monitor system status and availability</p>
    </div>

    <LoadingSpinner v-if="loading" text="Loading health status..." :centered="false" />

    <ErrorMessage
      v-else-if="error"
      title="Error loading health status"
      :message="error"
      retryable
      @retry="loadHealth"
    />

    <v-card v-else class="pa-8">
      <div class="d-flex align-center gap-3 mb-6">
        <v-avatar
          :color="healthData.status === 'ok' ? 'success' : 'warning'"
          size="48"
        >
          <v-icon
            :icon="healthData.status === 'ok' ? 'fa-check-circle' : 'fa-circle-exclamation'"
            color="white"
          />
        </v-avatar>
        <div>
          <h3 class="text-h6 font-weight-bold">System Status</h3>
          <p class="text-body-2 text-grey-darken-1">
            {{ healthData.status === 'ok' ? 'All systems operational' : 'System degraded' }}
          </p>
        </div>
      </div>

      <v-list class="bg-transparent">
        <v-list-item class="mb-2">
          <template #prepend>
            <v-avatar
              :color="healthData.status === 'ok' ? 'success' : 'warning'"
              size="8"
            />
          </template>
          <v-list-item-title class="font-weight-semibold">Status</v-list-item-title>
          <template #append>
            <BaseBadge
              :variant="healthData.status === 'ok' ? 'success' : 'warning'"
              size="sm"
            >
              {{ healthData.status.toUpperCase() }}
            </BaseBadge>
          </template>
        </v-list-item>

        <v-list-item class="mb-2">
          <template #prepend>
            <v-avatar
              :color="healthData.postgres === 'healthy' ? 'success' : 'error'"
              size="8"
            />
          </template>
          <v-list-item-title class="font-weight-semibold">PostgreSQL</v-list-item-title>
          <template #append>
            <BaseBadge
              :variant="healthData.postgres === 'healthy' ? 'success' : 'danger'"
              size="sm"
            >
              {{ healthData.postgres.toUpperCase() }}
            </BaseBadge>
          </template>
        </v-list-item>

        <v-list-item class="mb-2">
          <template #prepend>
            <v-avatar
              :color="healthData.redis === 'healthy' ? 'success' : 'error'"
              size="8"
            />
          </template>
          <v-list-item-title class="font-weight-semibold">Redis</v-list-item-title>
          <template #append>
            <BaseBadge
              :variant="healthData.redis === 'healthy' ? 'success' : 'danger'"
              size="sm"
            >
              {{ healthData.redis.toUpperCase() }}
            </BaseBadge>
          </template>
        </v-list-item>

        <v-list-item>
          <template #prepend>
            <v-avatar color="grey" size="8" />
          </template>
          <v-list-item-title class="font-weight-semibold">Timestamp</v-list-item-title>
          <template #append>
            <span class="text-body-2 text-grey-darken-1">
              {{ formatTimestamp(healthData.timestamp) }}
            </span>
          </template>
        </v-list-item>
      </v-list>
    </v-card>
  </v-container>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '@/lib/api'
import BaseBadge from '@/components/BaseBadge.vue'
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
