<template>
  <div class="min-h-screen bg-gray-50">
    <Header v-if="showHeader" />
    <main :class="showHeader ? 'py-8' : ''">
      <ProtectedRoute>
        <router-view />
      </ProtectedRoute>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import Header from './components/Header.vue'
import ProtectedRoute from './components/ProtectedRoute.vue'

const route = useRoute()
const showHeader = computed(() => !route.path.startsWith('/login'))
</script>
