<template>
  <v-container fluid class="fill-height pa-0">
    <v-row align="start" justify="center" class="fill-height">
      <v-col cols="12" sm="8" md="5" lg="4" class="pt-12">
        <v-card class="pa-8 elevation-2">
          <!-- Logo Section -->
          <div class="text-center mb-8">
            <div class="d-flex align-center justify-center gap-3 mb-3">
              <img src="/gopher.png" alt="AlgoShield" class="w-10 h-10 object-contain" />
              <h1 class="text-h5 font-weight-bold">AlgoShield</h1>
            </div>
            <p class="text-body-2 text-grey-darken-1">
              Enterprise AML Platform
            </p>
          </div>

          <!-- Tabs -->
          <v-tabs v-model="activeTab" class="mb-6">
            <v-tab value="login" @click="error = ''">Login</v-tab>
            <v-tab value="register" @click="error = ''">Register</v-tab>
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
                <v-text-field
                  v-model="email"
                  type="email"
                  label="Email"
                  placeholder="user@example.com"
                  :disabled="loading"
                  required
                  class="mb-4"
                />

                <v-text-field
                  v-model="password"
                  type="password"
                  label="Password"
                  placeholder="••••••••"
                  :disabled="loading"
                  :minlength="8"
                  required
                  class="mb-6"
                />

                <v-btn
                  type="submit"
                  :loading="loading"
                  :disabled="loading"
                  block
                  color="primary"
                  size="large"
                  class="mb-2"
                >
                  {{ loading ? 'Signing in...' : 'Sign In' }}
                </v-btn>
              </v-form>
            </v-window-item>

            <!-- Register Form -->
            <v-window-item value="register">
              <v-form @submit.prevent="handleRegister" class="mt-4">
                <v-text-field
                  v-model="name"
                  type="text"
                  label="Name"
                  placeholder="Your Name"
                  :disabled="loading"
                  required
                  class="mb-4"
                />

                <v-text-field
                  v-model="email"
                  type="email"
                  label="Email"
                  placeholder="user@example.com"
                  :disabled="loading"
                  required
                  class="mb-4"
                />

                <v-text-field
                  v-model="password"
                  type="password"
                  label="Password"
                  placeholder="••••••••"
                  :disabled="loading"
                  :minlength="8"
                  required
                  hint="Minimum 8 characters"
                  persistent-hint
                  class="mb-6"
                />

                <v-btn
                  type="submit"
                  :loading="loading"
                  :disabled="loading"
                  block
                  color="primary"
                  size="large"
                  class="mb-2"
                >
                  {{ loading ? 'Creating account...' : 'Create Account' }}
                </v-btn>
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
import { useAuthStore } from '@/stores/auth'
import { api } from '@/lib/api'
import ErrorMessage from '@/components/ErrorMessage.vue'

const router = useRouter()
const authStore = useAuthStore()

const email = ref('')
const password = ref('')
const name = ref('')
const activeTab = ref<'login' | 'register'>('login')
const loading = ref(false)
const error = ref('')

async function handleLogin() {
  if (!email.value || !password.value) {
    error.value = 'Please enter email and password'
    return
  }

  if (password.value.length < 8) {
    error.value = 'Password must be at least 8 characters'
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
      error.value = 'Invalid response from server'
      return
    }

    await authStore.setToken(response.token)
    router.push('/')
  } catch (e: any) {
    console.error('Login error:', e)
    const errorMsg = e.message || 'Login failed. Please try again.'
    
    if (errorMsg.includes('Invalid email or password') || errorMsg.includes('invalid email or password')) {
      error.value = 'Email ou senha inválidos. Verifique suas credenciais e tente novamente.'
    } else if (errorMsg.includes('not available') || errorMsg.includes('not found')) {
      error.value = 'Servidor da API não está disponível. Verifique se o servidor está rodando.'
    } else if (errorMsg.includes('timeout')) {
      error.value = 'A requisição demorou muito para responder. Tente novamente.'
    } else if (errorMsg.includes('Unable to connect')) {
      error.value = 'Não foi possível conectar ao servidor. Verifique sua conexão e se o servidor está rodando.'
    } else {
      error.value = errorMsg
    }
  } finally {
    loading.value = false
  }
}

async function handleRegister() {
  if (!email.value || !password.value || !name.value) {
    error.value = 'Please fill in all fields'
    return
  }

  if (password.value.length < 8) {
    error.value = 'Password must be at least 8 characters'
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
      error.value = 'Invalid response from server'
      return
    }

    await authStore.setToken(response.token)
    router.push('/')
  } catch (e: any) {
    console.error('Registration error:', e)
    const errorMsg = e.message || 'Registration failed. Please try again.'
    
    if (errorMsg.includes('already exists') || errorMsg.includes('já existe')) {
      error.value = 'Este email já está cadastrado. Use outro email ou faça login.'
    } else if (errorMsg.includes('validation failed') || errorMsg.includes('Validation failed')) {
      error.value = 'Dados inválidos. Verifique os campos e tente novamente.'
    } else if (errorMsg.includes('not available') || errorMsg.includes('not found')) {
      error.value = 'Servidor da API não está disponível. Verifique se o servidor está rodando.'
    } else if (errorMsg.includes('timeout')) {
      error.value = 'A requisição demorou muito para responder. Tente novamente.'
    } else if (errorMsg.includes('Unable to connect')) {
      error.value = 'Não foi possível conectar ao servidor. Verifique sua conexão e se o servidor está rodando.'
    } else {
      error.value = errorMsg
    }
  } finally {
    loading.value = false
  }
}
</script>
