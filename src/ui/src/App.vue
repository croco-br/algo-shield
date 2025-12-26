<template>
  <div class="min-h-screen bg-background relative">
    <!-- Header -->
    <Header v-if="showHeader" ref="headerRef" />

    <!-- Sidebar -->
    <Sidebar v-if="showSidebar" ref="sidebarRef" @collapse-change="handleSidebarCollapse" />

    <!-- Main Content -->
    <main
      class="transition-all duration-300"
      :style="mainStyles"
    >
      <div :class="[showHeader ? 'max-w-[1920px] mx-auto' : 'min-h-screen', contentPadding]">
        <ProtectedRoute>
          <router-view />
        </ProtectedRoute>
      </div>
    </main>
  </div>
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

const mainStyles = computed(() => {
  const styles: Record<string, string> = {}

  // Set min-height
  if (showHeader.value) {
    styles.minHeight = 'calc(100vh - var(--header-height))'
    styles.marginTop = 'var(--header-height)'
  } else {
    styles.minHeight = '100vh'
  }

  // Set margin-left for sidebar with additional spacing
  if (showSidebar.value) {
    const sidebarWidth = isSidebarCollapsed.value ? 'var(--sidebar-width-collapsed)' : 'var(--sidebar-width)'
    // Add extra spacing (24px = 1.5rem) to prevent overlap with collapse button
    styles.marginLeft = `calc(${sidebarWidth} + 1.5rem)`
    // Add matching spacing on the right side for visual balance
    styles.marginRight = '1.5rem'
  }

  // Ensure main content stays below fixed header/sidebar
  styles.position = 'relative'
  styles.zIndex = '0'

  return styles
})

const contentPadding = computed(() => {
  return showHeader.value ? 'py-16 px-12' : ''
})
</script>

<style scoped>
.bg-background {
  background-color: var(--color-background);
}
</style>
