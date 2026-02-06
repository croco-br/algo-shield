<template>
  <LoadingSpinner v-if="loading && !isPublicRoute" />
  <slot v-else />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import LoadingSpinner from '@/components/LoadingSpinner.vue'

const route = useRoute()
const authStore = useAuthStore()

const loading = computed(() => authStore.loading)
const isPublicRoute = computed(() => route.meta.public === true)
</script>
