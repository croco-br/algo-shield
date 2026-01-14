<template>
  <v-app-bar
    v-if="user && !isLoginPage"
    :height="'var(--header-height)'"
    :color="brandingConfig?.header_color || 'var(--color-header-background)'"
    fixed
    elevation="0"
    :class="['border-b border-neutral-800', { 'synthetic-mode': syntheticMode }]"
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

      <!-- Right: Synthetic Mode + User -->
      <div class="d-flex align-center gap-4">
        <!-- Synthetic Mode Toggle -->
        <div class="d-flex align-center gap-2">
          <v-chip 
            v-if="syntheticMode" 
            color="warning" 
            size="small" 
            class="font-weight-bold"
          >
            {{ $t('header.syntheticMode') }}
          </v-chip>
          <v-tooltip :text="syntheticMode ? $t('header.disableSynthetic') : $t('header.enableSynthetic')" location="bottom">
            <template #activator="{ props: tooltipProps }">
              <v-switch
                v-bind="tooltipProps"
                v-model="syntheticMode"
                :loading="syntheticModeLoading"
                :disabled="syntheticModeLoading || !isAdmin"
                :color="syntheticMode ? 'success' : 'grey'"
                hide-details
                density="compact"
                @update:model-value="handleSyntheticModeChange"
              />
            </template>
          </v-tooltip>
        </div>
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
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useBrandingStore } from '@/stores/branding'
import { useSystemModeStore } from '@/stores/systemMode'
import { useLocale } from '@/composables/useLocale'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const brandingStore = useBrandingStore()
const systemModeStore = useSystemModeStore()
const { locale, availableLocales, setLocale } = useLocale()

const showUserMenu = ref(false)

const user = computed(() => authStore.user)
const isLoginPage = computed(() => route.path.startsWith('/login'))
const brandingConfig = computed(() => brandingStore.config)
const isAdmin = computed(() => authStore.user?.roles?.some(r => r.name === 'admin') ?? false)

// Synthetic mode state - sync with store
const syntheticMode = ref(systemModeStore.syntheticMode)
const syntheticModeLoading = computed(() => systemModeStore.loading)

// Sync with store changes
watch(() => systemModeStore.syntheticMode, (newValue) => {
  syntheticMode.value = newValue
}, { immediate: true })

onMounted(async () => {
  if (user.value) {
    await systemModeStore.loadMode()
  }
})

const handleSyntheticModeChange = async (enabled: boolean | null) => {
  if (enabled === null) return
  try {
    await systemModeStore.setMode(enabled)
    // Reload the page to refresh all data with the new mode
    window.location.reload()
  } catch {
    // Revert on error - sync back with store
    syntheticMode.value = systemModeStore.syntheticMode
  }
}

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
  handleSyntheticModeChange,
})
</script>

<style scoped>
.v-app-bar.synthetic-mode {
  border-bottom: 3px solid rgb(var(--v-theme-warning)) !important;
}
</style>
