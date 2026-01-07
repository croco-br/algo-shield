<template>
  <v-container fluid class="fill-height pa-0">
    <v-row align="start" justify="center" class="fill-height">
      <v-col cols="12" sm="8" md="5" lg="4" class="pt-12">
        <v-card class="pa-8 elevation-2">
          <!-- Logo Section -->
          <div class="text-center mb-8">
            <div class="d-flex align-center justify-center gap-3 mb-3">
              <img src="/gopher.png" :alt="$t('common.appName')" class="w-10 h-10 object-contain" />
              <h1 class="text-h5 font-weight-bold">{{ $t('common.appName') }}</h1>
            </div>
            <p class="text-body-2 text-grey-darken-1">
              {{ $t('common.appTagline') }}
            </p>
          </div>

          <!-- Tabs -->
          <v-tabs v-model="activeTab" class="mb-6">
            <v-tab value="login" @click="error = ''">{{ $t('auth.login') }}</v-tab>
            <v-tab value="register" @click="error = ''">{{ $t('auth.register') }}</v-tab>
          </v-tabs>

          <ErrorMessage
            v-if="error"
            :message="error"
            variant="error"
            class="mb-6"
            @dismiss="error = ''"
          />

          <v-window v-model="activeTab">
            <!-- Login Form -->
            <v-window-item value="login">
              <v-form @submit.prevent="handleLogin" class="mt-4">
                <BaseInput
                  v-model="email"
                  type="email"
                  :label="$t('auth.email')"
                  :placeholder="$t('auth.emailPlaceholder')"
                  :disabled="loading"
                  required
                  class="mb-4"
                />

                <BaseInput
                  v-model="password"
                  type="password"
                  :label="$t('auth.password')"
                  :placeholder="$t('auth.passwordPlaceholder')"
                  :disabled="loading"
                  :minlength="8"
                  required
                  class="mb-6"
                />

                <BaseButton
                  type="submit"
                  :loading="loading"
                  :disabled="loading"
                  full-width
                  size="lg"
                  class="mb-2"
                >
                  {{ loading ? $t('auth.signingIn') : $t('auth.signIn') }}
                </BaseButton>
              </v-form>
            </v-window-item>

            <!-- Register Form -->
            <v-window-item value="register">
              <v-form @submit.prevent="handleRegister" class="mt-4">
                <BaseInput
                  v-model="name"
                  type="text"
                  :label="$t('auth.name')"
                  :placeholder="$t('auth.namePlaceholder')"
                  :disabled="loading"
                  required
                  class="mb-4"
                />

                <BaseInput
                  v-model="email"
                  type="email"
                  :label="$t('auth.email')"
                  :placeholder="$t('auth.emailPlaceholder')"
                  :disabled="loading"
                  required
                  class="mb-4"
                />

                <BaseInput
                  v-model="password"
                  type="password"
                  :label="$t('auth.password')"
                  :placeholder="$t('auth.passwordPlaceholder')"
                  :disabled="loading"
                  :minlength="8"
                  required
                  :hint="$t('auth.passwordHint')"
                  persistent-hint
                  class="mb-6"
                />

                <BaseButton
                  type="submit"
                  :loading="loading"
                  :disabled="loading"
                  full-width
                  size="lg"
                  class="mb-2"
                >
                  {{ loading ? $t('auth.creatingAccount') : $t('auth.createAccount') }}
                </BaseButton>
              </v-form>
            </v-window-item>
          </v-window>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { i18n } from '@/plugins/i18n'
import { useAuthStore } from '@/stores/auth'
import { api } from '@/lib/api'
import BaseButton from '@/components/BaseButton.vue'
import BaseInput from '@/components/BaseInput.vue'
import ErrorMessage from '@/components/ErrorMessage.vue'

const router = useRouter()
const authStore = useAuthStore()
const t = i18n.global.t

const email = ref('')
const password = ref('')
const name = ref('')
const activeTab = ref<'login' | 'register'>('login')
const loading = ref(false)
const error = ref('')

async function handleLogin() {
  if (!email.value || !password.value) {
    error.value = t('auth.errors.emailPassword')
    return
  }

  if (password.value.length < 8) {
    error.value = t('auth.errors.passwordLength')
    return
  }

  loading.value = true
  error.value = ''

  try {
    const response = await api.post<{ token: string; user: any }>('/api/v1/auth/login', {
      email: email.value,
      password: password.value,
    })

    if (!response || !response.token) {
      error.value = t('auth.errors.invalidResponse')
      return
    }

    await authStore.setToken(response.token)
    router.push('/')
  } catch (e: any) {
    console.error('Login error:', e)
    const errorMsg = e.message || t('auth.errors.loginFailed')
    
    if (errorMsg.includes('Invalid email or password') || errorMsg.includes('invalid email or password')) {
      error.value = t('auth.errors.invalidCredentials')
    } else if (errorMsg.includes('not available') || errorMsg.includes('not found')) {
      error.value = t('auth.errors.serverUnavailable')
    } else if (errorMsg.includes('timeout')) {
      error.value = t('auth.errors.timeout')
    } else if (errorMsg.includes('Unable to connect')) {
      error.value = t('auth.errors.connectionError')
    } else {
      error.value = errorMsg
    }
  } finally {
    loading.value = false
  }
}

async function handleRegister() {
  if (!email.value || !password.value || !name.value) {
    error.value = t('auth.errors.allFields')
    return
  }

  if (password.value.length < 8) {
    error.value = t('auth.errors.passwordLength')
    return
  }

  loading.value = true
  error.value = ''

  try {
    const response = await api.post<{ token: string; user: any }>('/api/v1/auth/register', {
      email: email.value,
      password: password.value,
      name: name.value,
    })

    if (!response || !response.token) {
      error.value = t('auth.errors.invalidResponse')
      return
    }

    await authStore.setToken(response.token)
    router.push('/')
  } catch (e: any) {
    console.error('Registration error:', e)
    const errorMsg = e.message || t('auth.errors.registrationFailed')
    
    if (errorMsg.includes('already exists') || errorMsg.includes('já existe')) {
      error.value = t('auth.errors.emailExists')
    } else if (errorMsg.includes('validation failed') || errorMsg.includes('Validation failed')) {
      error.value = t('auth.errors.validationFailed')
    } else if (errorMsg.includes('not available') || errorMsg.includes('not found')) {
      error.value = t('auth.errors.serverUnavailable')
    } else if (errorMsg.includes('timeout')) {
      error.value = t('auth.errors.timeout')
    } else if (errorMsg.includes('Unable to connect')) {
      error.value = t('auth.errors.connectionError')
    } else {
      error.value = errorMsg
    }
  } finally {
    loading.value = false
  }
}
</script>
