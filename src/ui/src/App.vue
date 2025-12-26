<template>
  <div class="min-h-screen bg-background relative">
    <!-- Header -->
    <Header v-if="showHeader" ref="headerRef" />

    <!-- Sidebar -->
    <Sidebar v-if="showSidebar" ref="sidebarRef" @collapse-change="handleSidebarCollapse" />

    <!-- Main Content -->
    <main
      :class="[
        'min-h-screen transition-all duration-300',
        mainMarginClass
      ]"
    >
      <div :class="showHeader ? 'max-w-[1920px] mx-auto px-12 py-16' : 'min-h-screen'">
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

const route = useRoute()
const headerRef = ref()
const sidebarRef = ref()

const showHeader = computed(() => !route.path.startsWith('/login'))
const showSidebar = computed(() => !route.path.startsWith('/login'))

const isSidebarCollapsed = ref(false)

const handleSidebarCollapse = (collapsed: boolean) => {
  isSidebarCollapsed.value = collapsed
}

const mainMarginClass = computed(() => {
  if (showHeader.value && showSidebar.value) {
    return isSidebarCollapsed.value ? 'ml-20 mt-[60px]' : 'ml-60 mt-[60px]'
  }
  if (showHeader.value && !showSidebar.value) {
    return 'mt-[60px]'
  }
  return ''
})
</script>

<style scoped>
.bg-background {
  background-color: var(--color-background);
}
</style>
