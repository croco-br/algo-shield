<template>
  <div class="max-w-7xl mx-auto px-8">
    <div class="bg-white rounded-lg border border-gray-200 p-8">
      <h2 class="text-2xl font-semibold mb-4">Health Check</h2>
      <div class="space-y-2">
        <p><strong>Status:</strong> {{ healthData.status }}</p>
        <p><strong>Service:</strong> {{ healthData.service }}</p>
        <p><strong>Timestamp:</strong> {{ healthData.timestamp }}</p>
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
  timestamp: new Date().toISOString()
})

onMounted(() => {
  // Health endpoint can be called here if needed
  // For now, we just show static data
})
</script>
