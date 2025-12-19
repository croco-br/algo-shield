<template>
  <header v-if="user && !isLoginPage" class="bg-white border-b border-gray-200 shadow-sm">
    <div class="max-w-7xl mx-auto px-8">
      <nav class="py-6 flex items-center justify-between gap-8">
        <div class="flex items-center gap-3">
          <img src="/gopher.png" alt="AlgoShield" class="w-10 h-10 object-contain" />
          <div>
            <h1 class="text-2xl font-semibold text-gray-900">AlgoShield</h1>
            <p class="text-sm text-gray-500">Fraud Detection & Anti-Money Laundering</p>
          </div>
        </div>

        <div class="flex gap-2 flex-1 justify-center">
          <router-link
            to="/"
            :class="[
              'px-6 py-3 rounded-md text-sm font-medium transition-all',
              $route.path === '/' 
                ? 'text-indigo-600 bg-indigo-50' 
                : 'text-gray-500 hover:text-gray-900 hover:bg-gray-50'
            ]"
          >
            Rules
          </router-link>
          <router-link
            to="/synthetic-test"
            :class="[
              'px-6 py-3 rounded-md text-sm font-medium transition-all',
              $route.path === '/synthetic-test' 
                ? 'text-indigo-600 bg-indigo-50' 
                : 'text-gray-500 hover:text-gray-900 hover:bg-gray-50'
            ]"
          >
            Synthetic Test
          </router-link>
          <router-link
            v-if="isAdmin"
            to="/permissions"
            :class="[
              'px-6 py-3 rounded-md text-sm font-medium transition-all',
              $route.path === '/permissions' 
                ? 'text-indigo-600 bg-indigo-50' 
                : 'text-gray-500 hover:text-gray-900 hover:bg-gray-50'
            ]"
          >
            Permissions
          </router-link>
        </div>

        <div class="relative" ref="menuRef">
          <button
            @click="showUserMenu = !showUserMenu"
            class="flex items-center gap-3 px-2 py-2 border border-gray-200 rounded-lg bg-white hover:bg-gray-50 transition-all"
          >
            <div class="w-10 h-10 rounded-full bg-gradient-to-br from-indigo-500 to-purple-600 text-white flex items-center justify-center font-semibold text-lg overflow-hidden">
              <img v-if="user.picture_url" :src="user.picture_url" :alt="user.name" class="w-full h-full object-cover" />
              <span v-else>{{ user.name.charAt(0).toUpperCase() }}</span>
            </div>
            <div class="text-left">
              <div class="text-sm font-medium text-gray-900">{{ user.name }}</div>
              <div class="text-xs text-gray-500">{{ user.email }}</div>
            </div>
            <svg
              class="w-4 h-4 text-gray-500"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <polyline points="6 9 12 15 18 9"></polyline>
            </svg>
          </button>

          <div v-if="showUserMenu" class="absolute top-full right-0 mt-2 bg-white border border-gray-200 rounded-lg shadow-lg min-w-[240px] z-50">
            <div class="p-4 border-b border-gray-200">
              <div class="font-semibold text-gray-900 mb-1">{{ user.name }}</div>
              <div class="text-sm text-gray-500 mb-2">{{ user.email }}</div>
              <div class="flex flex-wrap gap-2 mt-2">
                <span
                  v-for="role in (user.roles || [])"
                  :key="role.id"
                  class="inline-block px-3 py-1 bg-indigo-600 text-white rounded text-xs font-medium"
                >
                  {{ role.name }}
                </span>
              </div>
            </div>
            <button
              @click="handleLogout"
              class="w-full flex items-center gap-3 px-4 py-3 text-sm text-red-600 hover:bg-gray-50 transition-colors border-t border-gray-200"
            >
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
                <polyline points="16 17 21 12 16 7"></polyline>
                <line x1="21" y1="12" x2="9" y2="12"></line>
              </svg>
              Logout
            </button>
          </div>
        </div>
      </nav>
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
