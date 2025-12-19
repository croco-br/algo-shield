<template>
  <div v-if="loading && !isPublicRoute" class="flex items-center justify-center min-h-[50vh]">
    <div class="w-10 h-10 border-4 border-gray-200 border-t-indigo-600 rounded-full animate-spin"></div>
  </div>
  <slot v-else />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const authStore = useAuthStore()

const loading = computed(() => authStore.loading)
const isPublicRoute = computed(() => route.meta.public === true)
</script>
