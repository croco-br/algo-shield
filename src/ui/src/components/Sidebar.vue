<template>
  <v-navigation-drawer
    :model-value="!isMobile || isOpen"
    :width="isCollapsed ? 'var(--sidebar-width-collapsed)' : 'var(--sidebar-width)'"
    :temporary="isMobile"
    :permanent="!isMobile"
    :location="'left'"
    class="border-r"
    :style="{
      top: 'var(--header-height)',
      height: 'calc(100vh - var(--header-height))',
    }"
  >
    <!-- Collapse Toggle (Desktop) -->
    <template v-if="!isMobile" #append>
      <div class="d-flex justify-end pa-2">
        <v-btn
          icon
          variant="text"
          size="small"
          @click="toggleCollapse"
          class="position-absolute"
          style="right: -12px; top: 24px; z-index: 10;"
        >
          <v-icon>
            {{ isCollapsed ? 'mdi-chevron-right' : 'mdi-chevron-left' }}
          </v-icon>
        </v-btn>
      </div>
    </template>

    <!-- Navigation Links -->
    <v-list nav density="comfortable">
      <v-list-item
        v-for="item in navItems"
        :key="item.path"
        :to="item.path"
        :active="isActive(item.path)"
        :prepend-icon="getIcon(item.icon)"
        :title="isCollapsed ? '' : item.label"
        :value="item.path"
        @click="isMobile && closeMobile()"
        class="mx-2 mb-1"
      >
        <template v-if="isCollapsed" #prepend>
          <v-icon :icon="getIcon(item.icon)" />
        </template>
      </v-list-item>
    </v-list>
  </v-navigation-drawer>

  <!-- Mobile Overlay -->
  <v-overlay
    v-if="isMobile && isOpen"
    :model-value="isMobile && isOpen"
    class="align-start justify-start"
    @click="closeMobile"
  />
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
  { label: 'Dashboard', path: '/dashboard', icon: 'mdi-chart-line' },
  { label: 'Transactions', path: '/transactions', icon: 'mdi-swap-horizontal' },
  { label: 'Rules', path: '/rules', icon: 'mdi-format-list-checkbox' },
  { label: 'Permissions', path: '/permissions', icon: 'mdi-account-cog', adminOnly: true },
  { label: 'Branding', path: '/branding', icon: 'mdi-palette', adminOnly: true },
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

// Map FontAwesome icons to Material Design Icons
const getIcon = (faIcon: string): string => {
  const iconMap: Record<string, string> = {
    'fas fa-chart-line': 'mdi-chart-line',
    'fas fa-exchange-alt': 'mdi-swap-horizontal',
    'fas fa-tasks': 'mdi-format-list-checkbox',
    'fas fa-users-cog': 'mdi-account-cog',
    'fas fa-palette': 'mdi-palette',
  }
  return iconMap[faIcon] || 'mdi-circle'
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
