<template>
  <v-app>
    <!-- Header -->
    <Header v-if="showHeader" ref="headerRef" />

    <!-- Sidebar -->
    <Sidebar v-if="showSidebar" ref="sidebarRef" @collapse-change="handleSidebarCollapse" />

    <!-- Main Content -->
    <v-main v-if="showHeader">
      <div class="max-w-[1920px] mx-auto py-16 px-12">
        <ProtectedRoute>
          <router-view />
        </ProtectedRoute>
      </div>
    </v-main>

    <!-- Login (no header/sidebar) -->
    <v-main v-else>
      <ProtectedRoute>
        <router-view />
      </ProtectedRoute>
    </v-main>
  </v-app>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute } from 'vue-router'
import Header from './components/Header.vue'
import Sidebar from './components/Sidebar.vue'
import ProtectedRoute from './components/ProtectedRoute.vue'
import { useBrandingStore } from './stores/branding'

// Initialize branding store (loads and applies branding automatically)
useBrandingStore()

const route = useRoute()
const headerRef = ref()
const sidebarRef = ref()

const showHeader = computed(() => !route.path.startsWith('/login'))
const showSidebar = computed(() => !route.path.startsWith('/login'))

const isSidebarCollapsed = ref(false)

const handleSidebarCollapse = (collapsed: boolean) => {
  isSidebarCollapsed.value = collapsed
}
</script>
