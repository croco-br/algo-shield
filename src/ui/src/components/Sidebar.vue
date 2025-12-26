<template>
  <aside
    :class="[
      'sidebar fixed left-0 bg-white border-r border-neutral-200 transition-all duration-300 z-fixed',
      isMobile && !isOpen ? '-translate-x-full' : 'translate-x-0'
    ]"
    :style="{
      top: 'var(--header-height)',
      height: 'calc(100vh - var(--header-height))',
      width: isCollapsed ? 'var(--sidebar-width-collapsed)' : 'var(--sidebar-width)'
    }"
  >
    <!-- Collapse Toggle (Desktop) -->
    <button
      v-if="!isMobile"
      @click="toggleCollapse"
      class="absolute -right-3 top-6 w-6 h-6 bg-white border border-neutral-200 rounded-full flex items-center justify-center shadow-sm hover:shadow-md transition-all z-10"
    >
      <i :class="isCollapsed ? 'fas fa-chevron-right' : 'fas fa-chevron-left'" class="text-xs text-neutral-600"></i>
    </button>

    <!-- Navigation Links -->
    <nav class="py-6">
      <router-link
        v-for="item in navItems"
        :key="item.path"
        :to="item.path"
        :class="[
          'flex items-center gap-4 px-6 py-3 transition-all duration-200',
          isActive(item.path) ? 'text-teal-600 bg-teal-50 border-r-3 border-teal-600' : 'text-neutral-500 hover:text-neutral-900 hover:bg-neutral-50',
          isCollapsed ? 'justify-center' : ''
        ]"
        @click="isMobile && closeMobile()"
      >
        <i :class="[item.icon, 'text-lg', isCollapsed ? 'mx-0' : '']"></i>
        <span
          v-if="!isCollapsed"
          class="font-medium text-sm"
        >
          {{ item.label }}
        </span>
      </router-link>
    </nav>

  </aside>

  <!-- Mobile Overlay -->
  <div
    v-if="isMobile && isOpen"
    class="fixed inset-0 bg-black bg-opacity-50 z-modal-backdrop"
    @click="closeMobile"
  ></div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const emit = defineEmits<{
  collapseChange: [isCollapsed: boolean]
}>()

interface NavItem {
  label: string
  path: string
  icon: string
  adminOnly?: boolean
}

const route = useRoute()
const authStore = useAuthStore()

const allNavItems: NavItem[] = [
  { label: 'Dashboard', path: '/dashboard', icon: 'fas fa-chart-line' },
  { label: 'Transactions', path: '/transactions', icon: 'fas fa-exchange-alt' },
  { label: 'Rules', path: '/rules', icon: 'fas fa-tasks' },
  { label: 'Permissions', path: '/permissions', icon: 'fas fa-users-cog', adminOnly: true },
  { label: 'Branding', path: '/branding', icon: 'fas fa-palette', adminOnly: true },
]

const navItems = computed(() => {
  return allNavItems.filter(item => !item.adminOnly || authStore.isAdmin)
})

const isCollapsed = ref(false)
const isMobile = ref(false)
const isOpen = ref(false)

const toggleCollapse = () => {
  isCollapsed.value = !isCollapsed.value
  emit('collapseChange', isCollapsed.value)
}

const closeMobile = () => {
  isOpen.value = false
}

const isActive = (path: string) => {
  return route.path === path || (path !== '/' && route.path.startsWith(path))
}

const checkMobile = () => {
  isMobile.value = window.innerWidth < 960
  if (!isMobile.value) {
    isOpen.value = false
  }
}

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
  emit('collapseChange', isCollapsed.value)
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})

defineExpose({
  toggleMobile: () => {
    isOpen.value = !isOpen.value
  },
  isCollapsed
})
</script>

<style scoped>
.sidebar {
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.05);
}

.border-r-3 {
  border-right-width: 3px;
}
</style>
