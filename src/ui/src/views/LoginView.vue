<template>
  <div class="min-h-screen flex items-center justify-center bg-purple-600 p-4 sm:p-8">
    <div class="bg-white rounded-xl shadow-xl w-full max-w-md p-8 sm:p-10">
      <div class="mb-8">
        <div class="flex items-center gap-3 mb-2">
          <img src="/gopher.png" alt="AlgoShield" class="w-12 h-12 object-contain" />
          <h1 class="text-3xl font-semibold text-gray-900">AlgoShield</h1>
        </div>
        <p class="text-sm text-gray-500">Fraud Detection & Anti-Money Laundering</p>
      </div>

      <div class="flex gap-2 mb-6 border-b border-gray-200">
        <button
          @click="activeTab = 'login'; error = ''"
          :class="[
            'flex-1 py-3 text-sm font-medium transition-all',
            activeTab === 'login'
              ? 'text-indigo-600 font-semibold'
              : 'text-gray-500'
          ]"
        >
          <span
            :class="[
              'inline-block pb-3',
              activeTab === 'login' ? 'border-b-2 border-indigo-600' : ''
            ]"
          >
            Login
          </span>
        </button>
        <button
          @click="activeTab = 'register'; error = ''"
          :class="[
            'flex-1 py-3 text-sm font-medium transition-all',
            activeTab === 'register'
              ? 'text-indigo-600 font-semibold'
              : 'text-gray-500'
          ]"
        >
          <span
            :class="[
              'inline-block pb-3',
              activeTab === 'register' ? 'border-b-2 border-indigo-600' : ''
            ]"
          >
            Register
          </span>
        </button>
      </div>

      <div v-if="error" class="bg-red-50 border border-red-200 rounded-md p-3 mb-6 text-sm text-red-600">
        {{ error }}
      </div>

      <form v-if="activeTab === 'login'" @submit.prevent="handleLogin" class="space-y-6">
        <div>
          <label for="email" class="block text-sm font-medium text-gray-700 mb-2">
            Email
          </label>
          <input
            id="email"
            type="email"
            v-model="email"
            placeholder="user@example.com"
            required
            :disabled="loading"
            class="w-full px-4 py-3 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent disabled:opacity-60 disabled:cursor-not-allowed"
          />
        </div>

        <div>
          <label for="password" class="block text-sm font-medium text-gray-700 mb-2">
            Password
          </label>
          <input
            id="password"
            type="password"
            v-model="password"
            placeholder="••••••••"
            required
            minlength="8"
            :disabled="loading"
            class="w-full px-4 py-3 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent disabled:opacity-60 disabled:cursor-not-allowed"
          />
        </div>

        <button
          type="submit"
          :disabled="loading"
          class="w-full py-3 px-4 bg-indigo-600 text-white rounded-md font-medium hover:bg-indigo-700 transition-colors disabled:opacity-60 disabled:cursor-not-allowed"
        >
          {{ loading ? 'Signing in...' : 'Sign In' }}
        </button>
      </form>

      <form v-else @submit.prevent="handleRegister" class="space-y-6">
        <div>
          <label for="reg-name" class="block text-sm font-medium text-gray-700 mb-2">
            Name
          </label>
          <input
            id="reg-name"
            type="text"
            v-model="name"
            placeholder="Your Name"
            required
            :disabled="loading"
            class="w-full px-4 py-3 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent disabled:opacity-60 disabled:cursor-not-allowed"
          />
        </div>

        <div>
          <label for="reg-email" class="block text-sm font-medium text-gray-700 mb-2">
            Email
          </label>
          <input
            id="reg-email"
            type="email"
            v-model="email"
            placeholder="user@example.com"
            required
            :disabled="loading"
            class="w-full px-4 py-3 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent disabled:opacity-60 disabled:cursor-not-allowed"
          />
        </div>

        <div>
          <label for="reg-password" class="block text-sm font-medium text-gray-700 mb-2">
            Password
          </label>
          <input
            id="reg-password"
            type="password"
            v-model="password"
            placeholder="••••••••"
            required
            minlength="8"
            :disabled="loading"
            class="w-full px-4 py-3 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent disabled:opacity-60 disabled:cursor-not-allowed"
          />
          <small class="text-xs text-gray-500 mt-1 block">Minimum 8 characters</small>
        </div>

        <button
          type="submit"
          :disabled="loading"
          class="w-full py-3 px-4 bg-indigo-600 text-white rounded-md font-medium hover:bg-indigo-700 transition-colors disabled:opacity-60 disabled:cursor-not-allowed"
        >
          {{ loading ? 'Creating account...' : 'Create Account' }}
        </button>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { api } from '@/lib/api'

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
