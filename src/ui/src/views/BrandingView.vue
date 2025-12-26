<template>
  <div class="max-w-7xl mx-auto px-12">
    <div class="mb-10">
      <h1 class="text-3xl font-bold text-slate-900 mb-2">Branding Configuration</h1>
      <p class="text-slate-600 font-medium">Customize your application branding</p>
    </div>

    <div v-if="error" class="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg text-red-700">
      {{ error }}
    </div>

    <div v-if="success" class="mb-6 p-4 bg-green-50 border border-green-200 rounded-lg text-green-700">
      {{ success }}
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
      <!-- Configuration Form -->
      <div class="bg-white rounded-lg shadow-card p-8">
        <h2 class="text-xl font-semibold text-slate-900 mb-6">Configuration</h2>

        <form @submit.prevent="handleSubmit" class="space-y-6">
          <!-- App Name -->
          <div>
            <label for="appName" class="block text-sm font-medium text-slate-700 mb-2">
              Application Name
            </label>
            <input
              id="appName"
              v-model="form.app_name"
              type="text"
              required
              maxlength="100"
              class="w-full px-4 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              placeholder="AlgoShield"
            />
            <p class="text-xs text-slate-500 mt-1">Max 100 characters</p>
          </div>

          <!-- Icon URL -->
          <div>
            <label for="iconUrl" class="block text-sm font-medium text-slate-700 mb-2">
              Logo Icon URL
            </label>
            <input
              id="iconUrl"
              v-model="form.icon_url"
              type="text"
              class="w-full px-4 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              placeholder="/assets/logo.svg"
            />
            <p class="text-xs text-slate-500 mt-1">URL or path to logo image</p>
          </div>

          <!-- Favicon URL -->
          <div>
            <label for="faviconUrl" class="block text-sm font-medium text-slate-700 mb-2">
              Favicon URL
            </label>
            <input
              id="faviconUrl"
              v-model="form.favicon_url"
              type="text"
              class="w-full px-4 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              placeholder="/favicon.ico"
            />
            <p class="text-xs text-slate-500 mt-1">URL or path to favicon</p>
          </div>

          <!-- Primary Color -->
          <div>
            <label for="primaryColor" class="block text-sm font-medium text-slate-700 mb-2">
              Primary Color
            </label>
            <div class="flex gap-3 items-center">
              <input
                id="primaryColor"
                v-model="form.primary_color"
                type="color"
                class="h-12 w-20 border border-slate-300 rounded-lg cursor-pointer"
              />
              <input
                v-model="form.primary_color"
                type="text"
                pattern="^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$"
                required
                class="flex-1 px-4 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 font-mono"
                placeholder="#3B82F6"
              />
            </div>
            <p class="text-xs text-slate-500 mt-1">Hex format: #RGB or #RRGGBB</p>
          </div>

          <!-- Secondary Color -->
          <div>
            <label for="secondaryColor" class="block text-sm font-medium text-slate-700 mb-2">
              Secondary Color
            </label>
            <div class="flex gap-3 items-center">
              <input
                id="secondaryColor"
                v-model="form.secondary_color"
                type="color"
                class="h-12 w-20 border border-slate-300 rounded-lg cursor-pointer"
              />
              <input
                v-model="form.secondary_color"
                type="text"
                pattern="^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$"
                required
                class="flex-1 px-4 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 font-mono"
                placeholder="#10B981"
              />
            </div>
            <p class="text-xs text-slate-500 mt-1">Hex format: #RGB or #RRGGBB</p>
          </div>

          <!-- Action Buttons -->
          <div class="flex gap-3 pt-4">
            <button
              type="submit"
              :disabled="loading"
              class="flex-1 px-6 py-3 bg-blue-600 text-white font-medium rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              {{ loading ? 'Saving...' : 'Save Configuration' }}
            </button>
            <button
              type="button"
              @click="resetToDefaults"
              class="px-6 py-3 border-2 border-slate-300 text-slate-700 font-medium rounded-lg hover:bg-slate-50 transition-colors"
            >
              Reset to Defaults
            </button>
          </div>
        </form>
      </div>

      <!-- Live Preview -->
      <div class="bg-white rounded-lg shadow-card p-8">
        <h2 class="text-xl font-semibold text-slate-900 mb-6">Live Preview</h2>

        <div class="space-y-6">
          <!-- Logo Preview -->
          <div>
            <h3 class="text-sm font-medium text-slate-700 mb-3">Logo</h3>
            <div class="p-6 bg-slate-50 rounded-lg flex items-center gap-3">
              <img
                :src="form.icon_url || '/gopher.png'"
                :alt="form.app_name"
                class="w-10 h-10 object-contain"
                @error="handleImageError"
              />
              <span class="text-lg font-bold" :style="{ color: form.primary_color }">
                {{ form.app_name }}
              </span>
            </div>
          </div>

          <!-- Color Swatches -->
          <div>
            <h3 class="text-sm font-medium text-slate-700 mb-3">Colors</h3>
            <div class="space-y-3">
              <div class="flex items-center gap-3">
                <div
                  class="w-16 h-16 rounded-lg border-2 border-slate-200"
                  :style="{ backgroundColor: form.primary_color }"
                ></div>
                <div>
                  <p class="text-sm font-medium text-slate-900">Primary Color</p>
                  <p class="text-xs text-slate-500 font-mono">{{ form.primary_color }}</p>
                </div>
              </div>
              <div class="flex items-center gap-3">
                <div
                  class="w-16 h-16 rounded-lg border-2 border-slate-200"
                  :style="{ backgroundColor: form.secondary_color }"
                ></div>
                <div>
                  <p class="text-sm font-medium text-slate-900">Secondary Color</p>
                  <p class="text-xs text-slate-500 font-mono">{{ form.secondary_color }}</p>
                </div>
              </div>
            </div>
          </div>

          <!-- Button Preview -->
          <div>
            <h3 class="text-sm font-medium text-slate-700 mb-3">Buttons</h3>
            <div class="flex gap-3">
              <button
                type="button"
                class="px-6 py-2 text-white font-medium rounded-lg"
                :style="{ backgroundColor: form.primary_color }"
              >
                Primary
              </button>
              <button
                type="button"
                class="px-6 py-2 text-white font-medium rounded-lg"
                :style="{ backgroundColor: form.secondary_color }"
              >
                Secondary
              </button>
            </div>
          </div>

          <!-- Browser Title Preview -->
          <div>
            <h3 class="text-sm font-medium text-slate-700 mb-3">Browser Tab</h3>
            <div class="p-4 bg-slate-50 rounded-lg">
              <div class="flex items-center gap-2 text-sm">
                <div class="w-4 h-4 bg-slate-300 rounded"></div>
                <span class="text-slate-700">{{ form.app_name }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useBrandingStore } from '@/stores/branding'

const brandingStore = useBrandingStore()

const loading = ref(false)
const error = ref('')
const success = ref('')

const form = reactive({
  app_name: '',
  icon_url: '',
  favicon_url: '',
  primary_color: '#3B82F6',
  secondary_color: '#10B981',
})

onMounted(() => {
  // Load current branding configuration
  if (brandingStore.config) {
    form.app_name = brandingStore.config.app_name
    form.icon_url = brandingStore.config.icon_url || ''
    form.favicon_url = brandingStore.config.favicon_url || ''
    form.primary_color = brandingStore.config.primary_color
    form.secondary_color = brandingStore.config.secondary_color
  }
})

async function handleSubmit() {
  try {
    loading.value = true
    error.value = ''
    success.value = ''

    await brandingStore.updateBranding({
      app_name: form.app_name,
      icon_url: form.icon_url || null,
      favicon_url: form.favicon_url || null,
      primary_color: form.primary_color,
      secondary_color: form.secondary_color,
    })

    success.value = 'Branding configuration updated successfully!'
    setTimeout(() => {
      success.value = ''
    }, 5000)
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Failed to update branding configuration'
  } finally {
    loading.value = false
  }
}

function resetToDefaults() {
  form.app_name = 'AlgoShield'
  form.icon_url = '/assets/logo.svg'
  form.favicon_url = '/favicon.ico'
  form.primary_color = '#3B82F6'
  form.secondary_color = '#10B981'
}

function handleImageError(event: Event) {
  const img = event.target as HTMLImageElement
  img.src = '/gopher.png'
}
</script>
