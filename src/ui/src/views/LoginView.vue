<template>
  <div class="min-h-screen flex items-center justify-center bg-neutral-100 p-4">
    <div class="bg-white rounded shadow-card w-full max-w-[400px] px-10 py-12">
      <!-- Logo Section -->
      <div class="mb-10 text-center">
        <div class="flex items-center justify-center gap-3 mb-3">
          <img src="/gopher.png" alt="AlgoShield" class="w-10 h-10 object-contain" />
          <h1 class="text-2xl font-bold text-neutral-900">AlgoShield</h1>
        </div>
        <p class="text-sm text-neutral-600">
          Enterprise AML Platform
        </p>
      </div>

      <!-- Tabs -->
      <div class="flex gap-0 mb-8 border-b border-neutral-200">
        <button
          @click="activeTab = 'login'; error = ''"
          :class="[
            'flex-1 py-3 text-sm font-semibold transition-all duration-200 border-b-2',
            activeTab === 'login'
              ? 'text-teal-600 border-teal-600'
              : 'text-neutral-500 border-transparent hover:text-neutral-900'
          ]"
        >
          Login
        </button>
        <button
          @click="activeTab = 'register'; error = ''"
          :class="[
            'flex-1 py-3 text-sm font-semibold transition-all duration-200 border-b-2',
            activeTab === 'register'
              ? 'text-teal-600 border-teal-600'
              : 'text-neutral-500 border-transparent hover:text-neutral-900'
          ]"
        >
          Register
        </button>
      </div>

      <ErrorMessage
        v-if="error"
        :message="error"
        variant="error"
        class="mb-6"
        @dismiss="error = ''"
      />

      <form v-if="activeTab === 'login'" @submit.prevent="handleLogin" class="space-y-8">
        <div>
          <label class="block text-sm font-semibold text-neutral-700 mb-3">Email</label>
          <input
            type="email"
            v-model="email"
            placeholder="user@example.com"
            :disabled="loading"
            required
            class="w-full px-4 py-3 border border-neutral-200 rounded text-sm text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-teal-500 focus:border-transparent transition-all disabled:opacity-50 disabled:bg-neutral-50"
          />
        </div>

        <div>
          <label class="block text-sm font-semibold text-neutral-700 mb-3">Password</label>
          <input
            type="password"
            v-model="password"
            placeholder="••••••••"
            :disabled="loading"
            :minlength="8"
            required
            class="w-full px-4 py-3 border border-neutral-200 rounded text-sm text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-teal-500 focus:border-transparent transition-all disabled:opacity-50 disabled:bg-neutral-50"
          />
        </div>

        <div class="pt-4">
          <button
            type="submit"
            :disabled="loading"
            class="w-full px-6 py-3 bg-teal-600 text-white text-sm font-semibold rounded hover:bg-teal-700 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200 ripple"
          >
            {{ loading ? 'Signing in...' : 'Sign In' }}
          </button>
        </div>
      </form>

      <form v-else @submit.prevent="handleRegister" class="space-y-8">
        <div>
          <label class="block text-sm font-semibold text-neutral-700 mb-3">Name</label>
          <input
            type="text"
            v-model="name"
            placeholder="Your Name"
            :disabled="loading"
            required
            class="w-full px-4 py-3 border border-neutral-200 rounded text-sm text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-teal-500 focus:border-transparent transition-all disabled:opacity-50 disabled:bg-neutral-50"
          />
        </div>

        <div>
          <label class="block text-sm font-semibold text-neutral-700 mb-3">Email</label>
          <input
            type="email"
            v-model="email"
            placeholder="user@example.com"
            :disabled="loading"
            required
            class="w-full px-4 py-3 border border-neutral-200 rounded text-sm text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-teal-500 focus:border-transparent transition-all disabled:opacity-50 disabled:bg-neutral-50"
          />
        </div>

        <div>
          <label class="block text-sm font-semibold text-neutral-700 mb-3">Password</label>
          <input
            type="password"
            v-model="password"
            placeholder="••••••••"
            :disabled="loading"
            :minlength="8"
            required
            class="w-full px-4 py-3 border border-neutral-200 rounded text-sm text-neutral-900 placeholder-neutral-400 focus:outline-none focus:ring-2 focus:ring-teal-500 focus:border-transparent transition-all disabled:opacity-50 disabled:bg-neutral-50"
          />
          <p class="text-xs text-neutral-500 mt-2">Minimum 8 characters</p>
        </div>

        <div class="pt-4">
          <button
            type="submit"
            :disabled="loading"
            class="w-full px-6 py-3 bg-teal-600 text-white text-sm font-semibold rounded hover:bg-teal-700 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200 ripple"
          >
            {{ loading ? 'Creating account...' : 'Create Account' }}
          </button>
        </div>
      </form>
    </div>
  </div>
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
    // Use the error message from the API, or provide a user-friendly default
    const errorMsg = e.message || 'Login failed. Please try again.'
    
    // Improve error messages for common cases
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
    // Use the error message from the API, or provide a user-friendly default
    const errorMsg = e.message || 'Registration failed. Please try again.'
    
    // Improve error messages for common cases
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
