<template>
  <v-container fluid class="pa-8">
    <v-row>
      <v-col cols="12">
        <div class="d-flex justify-space-between align-center mb-8">
          <div>
            <div class="d-flex align-center gap-3 mb-2">
              <v-icon icon="fa-chart-line" size="large" color="primary" />
              <h2 class="text-h4 font-weight-bold">{{ $t('views.dashboard.title') }}</h2>
            </div>
            <p class="text-body-1 text-grey-darken-1">{{ $t('views.dashboard.subtitle') }}</p>
          </div>
          <div class="d-flex align-center gap-4">
            <span v-if="responseTime" class="text-caption text-grey-darken-1">
              <v-icon icon="fa-clock" size="x-small" class="mr-1" />
              {{ responseTime }}ms
            </span>
            <span v-if="lastUpdated" class="text-caption text-grey-darken-1">
              Updated {{ formatTimeAgo(lastUpdated) }}
            </span>
            <BaseButton @click="loadMetrics" :loading="loading" prepend-icon="fa-refresh">
              {{ $t('views.dashboard.refresh') }}
            </BaseButton>
          </div>
        </div>

        <LoadingSpinner v-if="loading && !metrics" :text="$t('views.dashboard.loading')" :centered="false" />

        <ErrorMessage
          v-else-if="error"
          :title="$t('views.dashboard.errorTitle')"
          :message="error"
          retryable
          @retry="loadMetrics"
        />

        <template v-else-if="metrics">
          <!-- Summary Cards -->
          <v-row class="mb-6">
            <v-col cols="12" md="3">
              <v-card class="pa-4 h-100" variant="outlined">
                <div class="d-flex align-center gap-3">
                  <v-avatar color="primary" size="48">
                    <v-icon icon="fa-exchange-alt" />
                  </v-avatar>
                  <div>
                    <div class="text-h4 font-weight-bold">{{ formatNumber(metrics.total_count) }}</div>
                    <div class="text-body-2 text-grey-darken-1">{{ $t('views.dashboard.totalTransactions') }}</div>
                  </div>
                </div>
              </v-card>
            </v-col>
            <v-col v-for="status in statusCards" :key="status.name" cols="12" md="3">
              <v-card class="pa-4 h-100" variant="outlined">
                <div class="d-flex align-center gap-3">
                  <v-avatar :color="status.color" size="48">
                    <v-icon :icon="status.icon" />
                  </v-avatar>
                  <div>
                    <div class="text-h4 font-weight-bold">{{ formatNumber(status.count) }}</div>
                    <div class="text-body-2 text-grey-darken-1">{{ status.label }}</div>
                  </div>
                </div>
              </v-card>
            </v-col>
          </v-row>

          <!-- Charts Row -->
          <v-row>
            <!-- Status Distribution -->
            <v-col cols="12" md="6">
              <v-card class="pa-4" variant="outlined">
                <h3 class="text-h6 font-weight-medium mb-4">{{ $t('views.dashboard.statusDistribution') }}</h3>
                <div class="chart-container">
                  <div v-for="item in metrics.status_distribution" :key="item.status" class="status-bar-item mb-3">
                    <div class="d-flex justify-space-between mb-1">
                      <span class="text-body-2 font-weight-medium">{{ item.status }}</span>
                      <span class="text-body-2 text-grey-darken-1">{{ item.count }} ({{ getPercentage(item.count) }}%)</span>
                    </div>
                    <v-progress-linear
                      :model-value="getPercentage(item.count)"
                      :color="getStatusColor(item.status)"
                      height="24"
                      rounded
                    />
                  </div>
                </div>
              </v-card>
            </v-col>

            <!-- Temporal View -->
            <v-col cols="12" md="6">
              <v-card class="pa-4" variant="outlined">
                <div class="d-flex justify-space-between align-center mb-4">
                  <h3 class="text-h6 font-weight-medium">{{ $t('views.dashboard.transactionVolume') }}</h3>
                  <v-btn-toggle v-model="selectedPeriod" mandatory density="compact" variant="outlined">
                    <v-btn value="24h" size="small">24h</v-btn>
                    <v-btn value="7d" size="small">7d</v-btn>
                    <v-btn value="30d" size="small">30d</v-btn>
                  </v-btn-toggle>
                </div>
                <div class="chart-container temporal-chart">
                  <div v-for="(item, index) in temporalData" :key="index" class="temporal-bar">
                    <div class="bar" :style="{ height: getBarHeight(item.count) + '%' }">
                      <span class="bar-value">{{ item.count }}</span>
                    </div>
                    <span class="bar-label">{{ formatBucketLabel(item.bucket) }}</span>
                  </div>
                </div>
              </v-card>
            </v-col>
          </v-row>
        </template>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useLocale } from '@/composables/useLocale'
import { api } from '@/lib/api'
import BaseButton from '@/components/BaseButton.vue'
import LoadingSpinner from '@/components/LoadingSpinner.vue'
import ErrorMessage from '@/components/ErrorMessage.vue'

const { t } = useLocale()

interface StatusCount {
  status: string
  count: number
}

interface TemporalCount {
  bucket: string
  count: number
}

interface DashboardMetrics {
  status_distribution: StatusCount[]
  temporal_24h: TemporalCount[]
  temporal_7d: TemporalCount[]
  temporal_30d: TemporalCount[]
  total_count: number
  cached_at: string
}

const loading = ref(false)
const error = ref('')
const metrics = ref<DashboardMetrics | null>(null)
const responseTime = ref<number | null>(null)
const lastUpdated = ref<Date | null>(null)
const selectedPeriod = ref('24h')

let autoRefreshInterval: ReturnType<typeof setInterval> | null = null

const statusCards = computed(() => {
  if (!metrics.value?.status_distribution) return []
  
  const statusConfig: Record<string, { color: string; icon: string; labelKey: string }> = {
    approved: { color: 'success', icon: 'fa-check', labelKey: 'views.dashboard.approved' },
    rejected: { color: 'error', icon: 'fa-times', labelKey: 'views.dashboard.rejected' },
    in_review: { color: 'warning', icon: 'fa-clock', labelKey: 'views.dashboard.inReview' },
    pending: { color: 'info', icon: 'fa-hourglass', labelKey: 'views.dashboard.pending' },
  }

  return metrics.value.status_distribution.map(item => {
    const config = statusConfig[item.status] || { color: 'grey', icon: 'fa-question', labelKey: '' }
    return {
      name: item.status,
      count: item.count,
      color: config.color,
      icon: config.icon,
      label: config.labelKey ? t(config.labelKey) : item.status
    }
  })
})

const temporalData = computed(() => {
  if (!metrics.value) return []
  
  switch (selectedPeriod.value) {
    case '24h': return metrics.value.temporal_24h || []
    case '7d': return metrics.value.temporal_7d || []
    case '30d': return metrics.value.temporal_30d || []
    default: return []
  }
})

const maxTemporalCount = computed(() => {
  return Math.max(...temporalData.value.map(d => d.count), 1)
})

onMounted(() => {
  loadMetrics()
  // Auto-refresh every 30 seconds
  autoRefreshInterval = setInterval(() => {
    loadMetrics()
  }, 30000)
})

onUnmounted(() => {
  if (autoRefreshInterval) {
    clearInterval(autoRefreshInterval)
  }
})

async function loadMetrics() {
  loading.value = true
  error.value = ''
  
  try {
    const response = await api.get<{ data: DashboardMetrics; response_time_ms: number }>('/api/v1/dashboard/metrics')
    metrics.value = response.data
    responseTime.value = response.response_time_ms
    lastUpdated.value = new Date()
  } catch (e: any) {
    error.value = e.message || 'Failed to load dashboard metrics'
    console.error('Error loading metrics:', e)
  } finally {
    loading.value = false
  }
}

function getPercentage(count: number): number {
  if (!metrics.value?.total_count) return 0
  return Math.round((count / metrics.value.total_count) * 100)
}

function getStatusColor(status: string): string {
  switch (status) {
    case 'approved': return 'success'
    case 'rejected': return 'error'
    case 'in_review': return 'warning'
    case 'pending': return 'info'
    default: return 'grey'
  }
}

function getBarHeight(count: number): number {
  return Math.max((count / maxTemporalCount.value) * 100, 5)
}

function formatBucketLabel(bucket: string): string {
  const date = new Date(bucket)
  if (selectedPeriod.value === '24h') {
    return date.toLocaleTimeString([], { hour: '2-digit' })
  }
  return date.toLocaleDateString([], { month: 'short', day: 'numeric' })
}

function formatNumber(num: number): string {
  if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M'
  if (num >= 1000) return (num / 1000).toFixed(1) + 'K'
  return String(num)
}

function formatTimeAgo(date: Date): string {
  const seconds = Math.floor((new Date().getTime() - date.getTime()) / 1000)
  if (seconds < 60) return 'just now'
  if (seconds < 3600) return Math.floor(seconds / 60) + 'm ago'
  return Math.floor(seconds / 3600) + 'h ago'
}
</script>

<style scoped>
.chart-container {
  min-height: 200px;
}

.temporal-chart {
  display: flex;
  align-items: flex-end;
  gap: 4px;
  height: 200px;
  padding-top: 20px;
}

.temporal-bar {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  height: 100%;
}

.temporal-bar .bar {
  width: 100%;
  max-width: 40px;
  background: linear-gradient(180deg, rgb(var(--v-theme-primary)) 0%, rgba(var(--v-theme-primary), 0.7) 100%);
  border-radius: 4px 4px 0 0;
  display: flex;
  align-items: flex-start;
  justify-content: center;
  position: relative;
  transition: height 0.3s ease;
}

.temporal-bar .bar-value {
  position: absolute;
  top: -20px;
  font-size: 10px;
  color: #666;
}

.temporal-bar .bar-label {
  font-size: 10px;
  color: #999;
  margin-top: 4px;
  white-space: nowrap;
}

.status-bar-item {
  transition: transform 0.2s ease;
}

.status-bar-item:hover {
  transform: translateX(4px);
}
</style>
