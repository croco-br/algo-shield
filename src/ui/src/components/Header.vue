<template>
  <header
    v-if="user && !isLoginPage"
    class="fixed top-0 left-0 right-0 h-[60px] bg-dark-slate border-b border-neutral-800 z-50"
    style="background: linear-gradient(180deg, #1f1f1f 0%, #1e1e1e 100%)"
  >
    <div class="flex items-center justify-between h-full px-8">
      <!-- Left: Logo + Search -->
      <div class="flex items-center gap-6">
        <!-- Logo -->
        <div class="flex items-center gap-3">
          <img src="/gopher.png" alt="AlgoShield" class="w-8 h-8 object-contain" />
          <span class="text-white font-bold text-lg">AlgoShield</span>
        </div>

        <!-- Global Search -->
        <div class="relative">
          <input
            type="search"
            placeholder="Search transactions, customers, alerts"
            class="w-[400px] h-[36px] pl-10 pr-4 bg-neutral-800 border border-neutral-700 rounded text-sm text-white placeholder-neutral-500 focus:outline-none focus:ring-2 focus:ring-teal-500 focus:border-transparent transition-all"
          />
          <i class="fas fa-search absolute left-3 top-1/2 -translate-y-1/2 text-neutral-500 text-sm"></i>
        </div>
      </div>

      <!-- Right: Notifications + User -->
      <div class="flex items-center gap-4">
        <!-- Notifications -->
        <button
          class="relative w-10 h-10 flex items-center justify-center rounded-full hover:bg-neutral-800 transition-colors"
        >
          <i class="far fa-bell text-white text-lg"></i>
          <span class="absolute top-1 right-1 w-2 h-2 bg-red-500 rounded-full"></span>
        </button>

        <!-- User Menu -->
        <div class="relative" ref="menuRef">
          <button
            @click="showUserMenu = !showUserMenu"
            class="flex items-center gap-2 px-2 py-1 rounded-lg hover:bg-neutral-800 transition-colors"
          >
            <div class="w-8 h-8 rounded-full bg-gradient-to-br from-teal-500 to-teal-600 text-white flex items-center justify-center font-semibold text-sm overflow-hidden">
              <img v-if="user.picture_url" :src="user.picture_url" :alt="user.name" class="w-full h-full object-cover" />
              <span v-else>{{ user.name.charAt(0).toUpperCase() }}</span>
            </div>
            <i class="fas fa-chevron-down text-neutral-400 text-xs"></i>
          </button>

          <!-- Dropdown -->
          <div
            v-if="showUserMenu"
            class="absolute top-full right-0 mt-2 bg-white border border-neutral-200 rounded-lg shadow-xl min-w-[220px] z-50"
          >
            <div class="p-4 border-b border-neutral-200">
              <div class="font-semibold text-neutral-900 text-sm">{{ user.name }}</div>
              <div class="text-xs text-neutral-500 mt-1">{{ user.email }}</div>
            </div>
            <div class="py-2">
              <router-link
                to="/profile"
                class="flex items-center gap-3 px-4 py-2 text-sm text-neutral-700 hover:bg-neutral-50 transition-colors"
              >
                <i class="fas fa-user w-4"></i>
                <span>Profile</span>
              </router-link>
            </div>
            <div class="border-t border-neutral-200">
              <button
                @click="handleLogout"
                class="w-full flex items-center gap-3 px-4 py-3 text-sm font-medium text-red-600 hover:bg-red-50 transition-colors"
              >
                <i class="fas fa-sign-out-alt w-4"></i>
                <span>Logout</span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const showUserMenu = ref(false)
const menuRef = ref<HTMLElement | null>(null)

const user = computed(() => authStore.user)
const isAdmin = computed(() => authStore.isAdmin)
const isLoginPage = computed(() => route.path.startsWith('/login'))

const handleLogout = async () => {
  await authStore.logout()
  router.push('/login')
}

function handleClickOutside(event: MouseEvent) {
  if (menuRef.value && !menuRef.value.contains(event.target as Node)) {
    showUserMenu.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>
