<template>
  <v-app-bar
    v-if="user && !isLoginPage"
    :height="'var(--header-height)'"
    :color="brandingConfig?.header_color || 'var(--color-header-background)'"
    fixed
    elevation="0"
    class="border-b border-neutral-800"
  >
    <div class="d-flex align-center justify-space-between w-100 px-4">
      <!-- Left: Logo -->
      <div class="d-flex align-center gap-6">
        <!-- Logo -->
        <div class="d-flex align-center gap-2 sm:gap-3">
          <v-avatar size="32" class="flex-shrink-0">
            <img
              :src="brandingConfig?.icon_url || '/gopher.png'"
              :alt="brandingConfig?.app_name || 'AlgoShield'"
              @error="handleLogoError"
              class="w-full h-full object-contain"
              loading="eager"
            />
          </v-avatar>
          <span class="text-white font-bold text-sm sm:text-lg truncate max-w-[200px] sm:max-w-none">
            {{ brandingConfig?.app_name || $t('common.appName') }}
          </span>
        </div>
      </div>

      <!-- Right: User -->
      <div class="d-flex align-center gap-2">
        <!-- User Menu -->
        <v-menu
          v-model="showUserMenu"
          location="bottom end"
          offset="8"
        >
          <template #activator="{ props: menuProps }">
            <v-btn
              v-bind="menuProps"
              variant="text"
              color="white"
              class="d-flex align-center gap-2"
            >
              <v-avatar size="32">
                <v-img
                  v-if="user.picture_url"
                  :src="user.picture_url"
                  :alt="user.name"
                  cover
                />
                <span v-else class="text-white">
                  {{ user.name.charAt(0).toUpperCase() }}
                </span>
              </v-avatar>
              <v-icon icon="fa-chevron-down" size="small" />
            </v-btn>
          </template>

          <v-list>
            <v-list-item>
              <template #prepend>
                <v-icon icon="fa-user" size="small" class="mr-2" />
              </template>
              <v-list-item-title class="font-weight-semibold">
                {{ user.name }}
              </v-list-item-title>
              <v-list-item-subtitle>
                {{ user.email }}
              </v-list-item-subtitle>
            </v-list-item>
            <v-divider />
            <v-list-item
              to="/profile"
              prepend-icon="fa-user"
            >
              {{ $t('header.profile') }}
            </v-list-item>
            <v-divider />
            <v-list-subheader class="text-uppercase text-caption">
              {{ $t('header.language') }}
            </v-list-subheader>
            <v-list-item
              v-for="localeOption in availableLocales"
              :key="localeOption.value"
              @click="setLocale(localeOption.value)"
              :active="locale === localeOption.value"
              class="cursor-pointer"
            >
              <template #prepend>
                <span class="mr-2" style="font-size: 1.2em">{{ localeOption.flag }}</span>
              </template>
              <v-list-item-title>
                {{ localeOption.label }}
              </v-list-item-title>
            </v-list-item>
            <v-divider />
            <v-list-item
              @click="handleLogout"
              prepend-icon="fa-sign-out-alt"
              class="text-error"
            >
              {{ $t('auth.logout') }}
            </v-list-item>
          </v-list>
        </v-menu>
      </div>
    </div>
  </v-app-bar>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useBrandingStore } from '@/stores/branding'
import { useLocale } from '@/composables/useLocale'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const brandingStore = useBrandingStore()
const { locale, availableLocales, setLocale } = useLocale()

const showUserMenu = ref(false)

const user = computed(() => authStore.user)
const isLoginPage = computed(() => route.path.startsWith('/login'))
const brandingConfig = computed(() => brandingStore.config)

const handleLogout = async () => {
  await authStore.logout()
  router.push('/login')
}

const handleLogoError = (event: Event) => {
  // Fallback to default logo if custom logo fails to load
  const img = event.target as HTMLImageElement
  img.src = '/gopher.png'
}

// Expose for testing
defineExpose({
  handleLogout,
  handleLogoError,
})
</script>

<style scoped>
</style>
