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
            {{ isCollapsed ? 'fa-chevron-right' : 'fa-chevron-left' }}
          </v-icon>
        </v-btn>
      </div>
    </template>

    <!-- Navigation Links -->
    <v-list nav density="comfortable">
      <template v-for="item in navItems" :key="item.path">
        <!-- Transactions with schema submenu -->
        <v-list-group
          v-if="item.path === '/transactions' && !isCollapsed && transactionSchemas.length > 0"
          :value="isTransactionsActive"
          color="primary"
        >
          <template #activator="{ props }">
            <v-list-item
              v-bind="props"
              :prepend-icon="getIcon(item.icon)"
              :title="$t(item.label)"
              :active="isTransactionsActive"
              class="mx-2 mb-1"
            />
          </template>
          <v-list-item
            v-for="schema in transactionSchemas"
            :key="schema.id"
            :to="`/transactions?schemaId=${schema.id}`"
            :active="isSchemaActive(schema.id)"
            :title="`${schema.name} - ${schema.fieldCount} colunas`"
            class="mx-2"
            @click="isMobile && closeMobile()"
          />
        </v-list-group>
        <!-- Regular menu items -->
        <v-list-item
          v-else
          :to="item.path"
          :active="isActive(item.path)"
          :prepend-icon="getIcon(item.icon)"
          :title="isCollapsed ? '' : $t(item.label)"
          :value="item.path"
          @click="isMobile && closeMobile()"
          class="mx-2 mb-1"
        >
          <template v-if="isCollapsed" #prepend>
            <v-icon :icon="getIcon(item.icon)" />
          </template>
        </v-list-item>
      </template>
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
import { api } from '@/lib/api'

const emit = defineEmits<{
  collapseChange: [isCollapsed: boolean]
}>()

interface NavItem {
  label: string
  path: string
  icon: string
  adminOnly?: boolean
}

interface Schema {
  id: string
  name: string
  extracted_fields: Array<{ path: string; type: string; nullable: boolean }>
}

interface TransactionSchema {
  id: string
  name: string
  fieldCount: number
}

const route = useRoute()
const authStore = useAuthStore()

const allNavItems: NavItem[] = [
  { label: 'sidebar.dashboard', path: '/dashboard', icon: 'fa-chart-line' },
  { label: 'sidebar.transactions', path: '/transactions', icon: 'fa-exchange-alt' },
  { label: 'sidebar.schemas', path: '/schemas', icon: 'fa-code' },
  { label: 'sidebar.rules', path: '/rules', icon: 'fa-tasks' },
  { label: 'sidebar.permissions', path: '/permissions', icon: 'fa-users-cog', adminOnly: true },
  { label: 'sidebar.branding', path: '/branding', icon: 'fa-palette', adminOnly: true },
]

const navItems = computed(() => {
  return allNavItems.filter(item => !item.adminOnly || authStore.isAdmin)
})

const isCollapsed = ref(false)
const isMobile = ref(false)
const isOpen = ref(false)
const transactionSchemas = ref<TransactionSchema[]>([])

const isTransactionsActive = computed(() => {
  return route.path === '/transactions'
})

const isSchemaActive = (schemaId: string) => {
  return route.path === '/transactions' && route.query.schemaId === schemaId
}

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

// Return Font Awesome icon name as-is
const getIcon = (icon: string): string => {
  return icon
}

const checkMobile = () => {
  isMobile.value = window.innerWidth < 960
  if (!isMobile.value) {
    isOpen.value = false
  }
}

async function loadTransactionSchemas() {
  try {
    const response = await api.get<{ schemas: Schema[] }>('/api/v1/schemas')
    transactionSchemas.value = (response?.schemas || []).map((schema: Schema) => ({
      id: schema.id,
      name: schema.name,
      fieldCount: schema.extracted_fields?.length || 0
    }))
  } catch (e) {
    console.error('Error loading schemas for sidebar:', e)
  }
}

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
  emit('collapseChange', isCollapsed.value)
  loadTransactionSchemas()
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
