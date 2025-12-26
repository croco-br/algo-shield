<template>
  <v-container fluid class="pa-8">
    <div class="mb-10">
      <h1 class="text-h4 font-weight-bold mb-2">Branding Configuration</h1>
      <p class="text-body-1 text-grey-darken-1">Customize your application branding</p>
    </div>

    <v-alert
      v-if="error"
      type="error"
      :text="error"
      class="mb-6"
      closable
      @click:close="error = ''"
    />

    <v-alert
      v-if="success"
      type="success"
      :text="success"
      class="mb-6"
      closable
      @click:close="success = ''"
    />

    <v-row>
      <!-- Configuration Form -->
      <v-col cols="12" lg="6">
        <v-card class="pa-8">
          <v-card-title class="text-h6 mb-6">Configuration</v-card-title>

          <v-form @submit.prevent="handleSubmit">
            <!-- App Name -->
            <v-text-field
              v-model="form.app_name"
              label="Application Name"
              placeholder="AlgoShield"
              required
              maxlength="100"
              hint="Max 100 characters"
              persistent-hint
              class="mb-4"
            />

            <!-- Icon URL -->
            <v-text-field
              v-model="form.icon_url"
              label="Logo Icon URL"
              placeholder="/assets/logo.svg"
              hint="URL or path to logo image"
              persistent-hint
              class="mb-4"
            />

            <!-- Favicon URL -->
            <v-text-field
              v-model="form.favicon_url"
              label="Favicon URL"
              placeholder="/favicon.ico"
              hint="URL or path to favicon"
              persistent-hint
              class="mb-4"
            />

            <!-- Primary Color -->
            <div class="mb-4">
              <label class="text-body-2 text-grey-darken-1 mb-2 d-block">Primary Color</label>
              <div class="d-flex align-center gap-3">
                <input
                  v-model="form.primary_color"
                  type="color"
                  class="h-12 w-20 border border-grey rounded-lg cursor-pointer"
                />
                <v-text-field
                  v-model="form.primary_color"
                  placeholder="#3B82F6"
                  pattern="^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$"
                  required
                  density="compact"
                  variant="outlined"
                  class="font-mono"
                />
              </div>
              <p class="text-caption text-grey mt-1">Hex format: #RGB or #RRGGBB</p>
            </div>

            <!-- Secondary Color -->
            <div class="mb-4">
              <label class="text-body-2 text-grey-darken-1 mb-2 d-block">Secondary Color</label>
              <div class="d-flex align-center gap-3">
                <input
                  v-model="form.secondary_color"
                  type="color"
                  class="h-12 w-20 border border-grey rounded-lg cursor-pointer"
                />
                <v-text-field
                  v-model="form.secondary_color"
                  placeholder="#10B981"
                  pattern="^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$"
                  required
                  density="compact"
                  variant="outlined"
                  class="font-mono"
                />
              </div>
              <p class="text-caption text-grey mt-1">Hex format: #RGB or #RRGGBB</p>
            </div>

            <!-- Header Color -->
            <div class="mb-6">
              <label class="text-body-2 text-grey-darken-1 mb-2 d-block">Header Background Color</label>
              <div class="d-flex align-center gap-3">
                <input
                  v-model="form.header_color"
                  type="color"
                  class="h-12 w-20 border border-grey rounded-lg cursor-pointer"
                />
                <v-text-field
                  v-model="form.header_color"
                  placeholder="#1e1e1e"
                  pattern="^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$"
                  required
                  density="compact"
                  variant="outlined"
                  class="font-mono"
                />
              </div>
              <p class="text-caption text-grey mt-1">Hex format: #RGB or #RRGGBB</p>
            </div>

            <!-- Action Buttons -->
            <div class="d-flex gap-3">
              <v-btn
                type="submit"
                :loading="loading"
                :disabled="loading"
                color="primary"
                block
                size="large"
              >
                {{ loading ? 'Saving...' : 'Save Configuration' }}
              </v-btn>
              <v-btn
                variant="outlined"
                @click="resetToDefaults"
                size="large"
              >
                Reset to Defaults
              </v-btn>
            </div>
          </v-form>
        </v-card>
      </v-col>

      <!-- Live Preview -->
      <v-col cols="12" lg="6">
        <v-card class="pa-8">
          <v-card-title class="text-h6 mb-6">Live Preview</v-card-title>

          <div class="d-flex flex-column gap-6">
            <!-- Logo Preview -->
            <div>
              <h3 class="text-body-2 font-weight-medium text-grey-darken-1 mb-3">Logo</h3>
              <v-card variant="tonal" class="pa-6">
                <div class="d-flex align-center gap-3">
                  <img
                    :src="form.icon_url || '/gopher.png'"
                    :alt="form.app_name"
                    class="w-10 h-10 object-contain"
                    @error="handleImageError"
                  />
                  <span class="text-h6 font-weight-bold" :style="{ color: form.primary_color }">
                    {{ form.app_name }}
                  </span>
                </div>
              </v-card>
            </div>

            <!-- Color Swatches -->
            <div>
              <h3 class="text-body-2 font-weight-medium text-grey-darken-1 mb-3">Colors</h3>
              <div class="d-flex flex-column gap-3">
                <div class="d-flex align-center gap-3">
                  <div
                    class="w-16 h-16 rounded-lg border"
                    :style="{ backgroundColor: form.primary_color }"
                  />
                  <div>
                    <p class="text-body-2 font-weight-medium">Primary Color</p>
                    <p class="text-caption text-grey font-mono">{{ form.primary_color }}</p>
                  </div>
                </div>
                <div class="d-flex align-center gap-3">
                  <div
                    class="w-16 h-16 rounded-lg border"
                    :style="{ backgroundColor: form.secondary_color }"
                  />
                  <div>
                    <p class="text-body-2 font-weight-medium">Secondary Color</p>
                    <p class="text-caption text-grey font-mono">{{ form.secondary_color }}</p>
                  </div>
                </div>
                <div class="d-flex align-center gap-3">
                  <div
                    class="w-16 h-16 rounded-lg border"
                    :style="{ backgroundColor: form.header_color }"
                  />
                  <div>
                    <p class="text-body-2 font-weight-medium">Header Background</p>
                    <p class="text-caption text-grey font-mono">{{ form.header_color }}</p>
                  </div>
                </div>
              </div>
            </div>

            <!-- Header Preview -->
            <div>
              <h3 class="text-body-2 font-weight-medium text-grey-darken-1 mb-3">Header</h3>
              <v-card
                :style="{ backgroundColor: form.header_color }"
                class="pa-4"
              >
                <div class="d-flex align-center gap-3">
                  <img
                    :src="form.icon_url || '/gopher.png'"
                    :alt="form.app_name"
                    class="w-8 h-8 object-contain"
                    @error="handleImageError"
                  />
                  <span class="text-white font-weight-bold text-body-2">{{ form.app_name }}</span>
                </div>
              </v-card>
            </div>

            <!-- Button Preview -->
            <div>
              <h3 class="text-body-2 font-weight-medium text-grey-darken-1 mb-3">Buttons</h3>
              <div class="d-flex gap-3">
                <v-btn
                  :style="{ backgroundColor: form.primary_color }"
                  color="primary"
                >
                  Primary
                </v-btn>
                <v-btn
                  :style="{ backgroundColor: form.secondary_color }"
                  color="secondary"
                >
                  Secondary
                </v-btn>
              </div>
            </div>

            <!-- Browser Title Preview -->
            <div>
              <h3 class="text-body-2 font-weight-medium text-grey-darken-1 mb-3">Browser Tab</h3>
              <v-card variant="tonal" class="pa-4">
                <div class="d-flex align-center gap-2 text-body-2">
                  <div class="w-4 h-4 bg-grey rounded" />
                  <span class="text-grey-darken-1">{{ form.app_name }}</span>
                </div>
              </v-card>
            </div>
          </div>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { useBrandingStore, type BrandingConfig } from '@/stores/branding'

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
  header_color: '#1e1e1e',
})

// Watch for config changes and populate form when data is available
watch(
  () => brandingStore.config,
  (config: BrandingConfig | null) => {
    if (config) {
      form.app_name = config.app_name
      form.icon_url = config.icon_url || ''
      form.favicon_url = config.favicon_url || ''
      form.primary_color = config.primary_color
      form.secondary_color = config.secondary_color
      form.header_color = config.header_color
    }
  },
  { immediate: true }
)

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
      header_color: form.header_color,
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
  form.header_color = '#1e1e1e'
}

function handleImageError(event: Event) {
  const img = event.target as HTMLImageElement
  img.src = '/gopher.png'
}
</script>
