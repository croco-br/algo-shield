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
      <div :class="showHeader ? 'max-w-[1920px] mx-auto px-12 py-12' : ''">
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

<style>
/* Minimalist Global Styles */
:root {
  --color-background: #f0f0f0;
  --color-dark-slate: #1e1e1e;
  --color-teal: #00bfa5;
  --color-neutral-900: #333333;
  --color-neutral-600: #a0a0a0;
}

.bg-background {
  background-color: var(--color-background);
}

/* Subtle paper texture overlay (only on dashboard background, not cards) */
body::before {
  content: '';
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-image: url("data:image/svg+xml,%3Csvg width='100' height='100' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noise'%3E%3CfeTurbulence baseFrequency='0.65' numOctaves='3' /%3E%3C/filter%3E%3Crect width='100' height='100' filter='url(%23noise)' opacity='0.03'/%3E%3C/svg%3E");
  pointer-events: none;
  z-index: 0;
}

/* FontAwesome CDN already included in index.html */
</style>
