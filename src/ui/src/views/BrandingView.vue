<template>
  <v-container fluid class="pa-4" style="max-width: 100%; overflow-x: hidden;">
    <div class="mb-3">
      <div class="d-flex align-center gap-2 mb-1">
        <v-icon icon="fa-palette" size="default" color="primary" />
        <h1 class="text-h5 font-weight-bold">Branding Configuration</h1>
      </div>
      <p class="text-body-2 text-grey-darken-1">Customize your application branding</p>
    </div>

    <v-alert
      v-if="error"
      type="error"
      :text="error"
      class="mb-3"
      closable
      density="compact"
      @click:close="error = ''"
    />

    <v-alert
      v-if="success"
      type="success"
      :text="success"
      class="mb-3"
      closable
      density="compact"
      @click:close="success = ''"
    />

    <v-row class="g-0" style="margin: 0;">
      <!-- Configuration Form -->
      <v-col cols="12" lg="6" class="pr-lg-2" style="padding-left: 0;">
        <v-card class="pa-3">
          <v-card-title class="text-subtitle-1 mb-3" style="font-weight: 600;">Configuration</v-card-title>

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
              prepend-inner-icon="fa-window-maximize"
              class="mb-3"
              variant="outlined"
              density="compact"
            />

            <!-- Icon URL -->
            <v-text-field
              v-model="form.icon_url"
              label="Logo Icon URL"
              placeholder="/assets/logo.svg"
              hint="URL or path to logo image"
              persistent-hint
              prepend-inner-icon="fa-image"
              class="mb-3"
              variant="outlined"
              density="compact"
            />

            <!-- Favicon URL -->
            <v-text-field
              v-model="form.favicon_url"
              label="Favicon URL"
              placeholder="/favicon.ico"
              hint="URL or path to favicon"
              persistent-hint
              prepend-inner-icon="fa-star"
              class="mb-3"
              variant="outlined"
              density="compact"
            />

            <!-- Primary Color -->
            <div class="mb-3">
              <label class="text-body-2 text-grey-darken-1 mb-1 d-flex align-center gap-2">
                <v-icon icon="fa-palette" size="x-small" />
                Primary Color
              </label>
              <div class="d-flex align-center gap-2">
                <input
                  v-model="form.primary_color"
                  type="color"
                  class="color-picker-input"
                  style="width: 50px; height: 40px; border: 1px solid rgba(0,0,0,0.12); border-radius: 4px; cursor: pointer; flex-shrink: 0;"
                />
                <v-text-field
                  v-model="form.primary_color"
                  placeholder="#3B82F6"
                  pattern="^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$"
                  required
                  density="compact"
                  variant="outlined"
                  prepend-inner-icon="fa-hashtag"
                  class="font-mono flex-grow-1"
                  hide-details
                />
              </div>
            </div>

            <!-- Secondary Color -->
            <div class="mb-3">
              <label class="text-body-2 text-grey-darken-1 mb-1 d-flex align-center gap-2">
                <v-icon icon="fa-palette" size="x-small" />
                Secondary Color
              </label>
              <div class="d-flex align-center gap-2">
                <input
                  v-model="form.secondary_color"
                  type="color"
                  class="color-picker-input"
                  style="width: 50px; height: 40px; border: 1px solid rgba(0,0,0,0.12); border-radius: 4px; cursor: pointer; flex-shrink: 0;"
                />
                <v-text-field
                  v-model="form.secondary_color"
                  placeholder="#10B981"
                  pattern="^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$"
                  required
                  density="compact"
                  variant="outlined"
                  prepend-inner-icon="fa-hashtag"
                  class="font-mono flex-grow-1"
                  hide-details
                />
              </div>
            </div>

            <!-- Header Color -->
            <div class="mb-4">
              <label class="text-body-2 text-grey-darken-1 mb-1 d-flex align-center gap-2">
                <v-icon icon="fa-heading" size="x-small" />
                Header Background Color
              </label>
              <div class="d-flex align-center gap-2">
                <input
                  v-model="form.header_color"
                  type="color"
                  class="color-picker-input"
                  style="width: 50px; height: 40px; border: 1px solid rgba(0,0,0,0.12); border-radius: 4px; cursor: pointer; flex-shrink: 0;"
                />
                <v-text-field
                  v-model="form.header_color"
                  placeholder="#1e1e1e"
                  pattern="^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$"
                  required
                  density="compact"
                  variant="outlined"
                  prepend-inner-icon="fa-hashtag"
                  class="font-mono flex-grow-1"
                  hide-details
                />
              </div>
            </div>

            <!-- Action Buttons -->
            <div class="d-flex flex-column gap-2">
              <BaseButton
                type="submit"
                :loading="loading"
                :disabled="loading"
                full-width
                prepend-icon="fa-save"
              >
                {{ loading ? 'Saving...' : 'Save Configuration' }}
              </BaseButton>
              <BaseButton
                variant="secondary"
                @click="resetToDefaults"
                full-width
                prepend-icon="fa-rotate-left"
              >
                Reset to Defaults
              </BaseButton>
            </div>
          </v-form>
        </v-card>
      </v-col>

      <!-- Live Preview -->
      <v-col cols="12" lg="6" class="pl-lg-2" style="padding-right: 0;">
        <v-card class="pa-3">
          <v-card-title class="text-subtitle-1 mb-3" style="font-weight: 600;">Live Preview</v-card-title>

          <div class="d-flex flex-column gap-3">
          <div class="d-flex flex-column gap-2">
            <!-- Logo Preview -->
            <div>
              <h3 class="text-caption font-weight-medium text-grey-darken-1 mb-1">Logo</h3>
              <v-card variant="tonal" class="pa-2">
                <div class="d-flex align-center gap-2">
                  <img
                    :src="form.icon_url || '/gopher.png'"
                    :alt="form.app_name"
                    style="width: 24px; height: 24px; object-fit: contain;"
                    @error="handleImageError"
                  />
                  <span class="text-body-2 font-weight-bold" :style="{ color: form.primary_color }">
                    {{ form.app_name }}
                  </span>
                </div>
              </v-card>
            </div>

            <!-- Color Swatches -->
            <div>
              <h3 class="text-caption font-weight-medium text-grey-darken-1 mb-1">Colors</h3>
              <div class="d-flex flex-column gap-1">
                <div class="d-flex align-center gap-2">
                  <div
                    style="width: 32px; height: 32px; border-radius: 4px; border: 1px solid rgba(0,0,0,0.12); flex-shrink: 0;"
                    :style="{ backgroundColor: form.primary_color }"
                  />
                  <div>
                    <p class="text-caption font-weight-medium mb-0">Primary</p>
                    <p class="text-caption text-grey font-mono mb-0" style="font-size: 10px;">{{ form.primary_color }}</p>
                  </div>
                </div>
                <div class="d-flex align-center gap-2">
                  <div
                    style="width: 32px; height: 32px; border-radius: 4px; border: 1px solid rgba(0,0,0,0.12); flex-shrink: 0;"
                    :style="{ backgroundColor: form.secondary_color }"
                  />
                  <div>
                    <p class="text-caption font-weight-medium mb-0">Secondary</p>
                    <p class="text-caption text-grey font-mono mb-0" style="font-size: 10px;">{{ form.secondary_color }}</p>
                  </div>
                </div>
                <div class="d-flex align-center gap-2">
                  <div
                    style="width: 32px; height: 32px; border-radius: 4px; border: 1px solid rgba(0,0,0,0.12); flex-shrink: 0;"
                    :style="{ backgroundColor: form.header_color }"
                  />
                  <div>
                    <p class="text-caption font-weight-medium mb-0">Header</p>
                    <p class="text-caption text-grey font-mono mb-0" style="font-size: 10px;">{{ form.header_color }}</p>
                  </div>
                </div>
              </div>
            </div>

            <!-- Header Preview -->
            <div>
              <h3 class="text-caption font-weight-medium text-grey-darken-1 mb-1">Header</h3>
              <v-card
                :style="{ backgroundColor: form.header_color }"
                class="pa-2"
              >
                <div class="d-flex align-center gap-2">
                  <img
                    :src="form.icon_url || '/gopher.png'"
                    :alt="form.app_name"
                    style="width: 20px; height: 20px; object-fit: contain;"
                    @error="handleImageError"
                  />
                  <span class="text-white font-weight-bold text-caption">{{ form.app_name }}</span>
                </div>
              </v-card>
            </div>

            <!-- Button Preview -->
            <div>
              <h3 class="text-caption font-weight-medium text-grey-darken-1 mb-1">Buttons</h3>
              <div class="d-flex gap-2">
                <BaseButton
                  size="sm"
                  :style="{ backgroundColor: form.primary_color }"
                >
                  Primary
                </BaseButton>
                <BaseButton
                  size="sm"
                  :style="{ backgroundColor: form.secondary_color }"
                >
                  Secondary
                </BaseButton>
              </div>
            </div>

            <!-- Browser Title Preview -->
            <div>
              <h3 class="text-caption font-weight-medium text-grey-darken-1 mb-1">Browser Tab</h3>
              <v-card variant="tonal" class="pa-1">
                <div class="d-flex align-center gap-1 text-caption">
                  <div style="width: 12px; height: 12px; background-color: #grey; border-radius: 2px;" />
                  <span class="text-grey-darken-1" style="font-size: 11px;">{{ form.app_name }}</span>
                </div>
              </v-card>
            </div>
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
import BaseButton from '@/components/BaseButton.vue'

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
