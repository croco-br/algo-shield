<template>
  <div class="max-w-7xl mx-auto px-8">
    <div class="mb-8">
      <h2 class="text-3xl font-bold text-slate-900 mb-2">System Health</h2>
      <p class="text-slate-600 font-medium">Monitor system status and availability</p>
    </div>

    <div class="bg-white rounded-xl border-2 border-slate-200 shadow-sm p-8">
      <div class="flex items-center gap-3 mb-6">
        <div class="w-12 h-12 rounded-full bg-emerald-100 flex items-center justify-center">
          <svg class="w-6 h-6 text-emerald-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
          </svg>
        </div>
        <div>
          <h3 class="text-xl font-bold text-slate-900">System Status</h3>
          <p class="text-sm text-slate-600">All systems operational</p>
        </div>
      </div>

      <div class="space-y-4">
        <div class="flex items-center justify-between py-3 px-4 bg-slate-50 rounded-lg">
          <div class="flex items-center gap-3">
            <div class="w-2 h-2 rounded-full bg-emerald-500"></div>
            <span class="font-semibold text-slate-700">Status</span>
          </div>
          <span class="text-sm font-medium text-emerald-600 uppercase">{{ healthData.status }}</span>
        </div>

        <div class="flex items-center justify-between py-3 px-4 bg-slate-50 rounded-lg">
          <div class="flex items-center gap-3">
            <div class="w-2 h-2 rounded-full bg-blue-500"></div>
            <span class="font-semibold text-slate-700">Service</span>
          </div>
          <span class="text-sm font-medium text-slate-900">{{ healthData.service }}</span>
        </div>

        <div class="flex items-center justify-between py-3 px-4 bg-slate-50 rounded-lg">
          <div class="flex items-center gap-3">
            <div class="w-2 h-2 rounded-full bg-slate-400"></div>
            <span class="font-semibold text-slate-700">Timestamp</span>
          </div>
          <span class="text-sm font-medium text-slate-600">{{ healthData.timestamp }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

interface HealthData {
  status: string
  timestamp: string
  service: string
}

const healthData = ref<HealthData>({
  status: 'ok',
  timestamp: new Date().toISOString(),  
  service: "algo-shield-ui"
})

onMounted(() => {
  // Health endpoint can be called here if needed
  // For now, we just show static data
})
</script>
