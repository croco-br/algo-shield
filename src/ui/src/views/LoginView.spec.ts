import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'
import { createPinia, setActivePinia } from 'pinia'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import { createI18n } from 'vue-i18n'
import LoginView from './LoginView.vue'
import ErrorMessage from '@/components/ErrorMessage.vue'

// Mock the api module
vi.mock('@/lib/api', () => ({
  api: {
    post: vi.fn(),
  },
}))

// Mock child components
vi.mock('@/components/BaseButton.vue', () => ({
  default: {
    name: 'BaseButton',
    template: '<button type="button"><slot /></button>',
    props: ['loading', 'disabled', 'fullWidth', 'size', 'type'],
  },
}))

vi.mock('@/components/BaseInput.vue', () => ({
  default: {
    name: 'BaseInput',
    template: '<input :type="type" :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />',
    props: ['modelValue', 'type', 'label', 'placeholder', 'disabled', 'required', 'minlength', 'hint', 'persistentHint'],
    emits: ['update:modelValue'],
  },
}))

vi.mock('@/components/ErrorMessage.vue', () => ({
  default: {
    name: 'ErrorMessage',
    template: '<v-alert v-if="message" type="error" class="mb-4">{{ message }}</v-alert>',
    props: ['message', 'variant'],
    emits: ['dismiss'],
  },
}))

describe('LoginView', () => {
  let router: ReturnType<typeof createRouter>
  let vuetify: ReturnType<typeof createVuetify>
  let i18n: ReturnType<typeof createI18n>

  beforeEach(() => {
    setActivePinia(createPinia())
    
    router = createRouter({
      history: createMemoryHistory(),
      routes: [
        { path: '/', component: { template: '<div>Home</div>' } },
        { path: '/login', component: LoginView },
      ],
    })

    vuetify = createVuetify({
      components,
      directives,
    })

    i18n = createI18n({
      legacy: true,
      locale: 'en-US',
      messages: {
        'en-US': {
          common: { appName: 'AlgoShield', appTagline: 'Fraud Detection' },
          auth: {
            login: 'Login',
            register: 'Register',
            email: 'Email',
            emailPlaceholder: 'Enter your email',
            password: 'Password',
            passwordPlaceholder: 'Enter your password',
            name: 'Name',
            namePlaceholder: 'Enter your name',
            passwordHint: 'Minimum 8 characters',
            signIn: 'Sign In',
            signingIn: 'Signing In...',
            createAccount: 'Create Account',
            creatingAccount: 'Creating Account...',
            errors: {
              emailPassword: 'Email and password are required',
              passwordLength: 'Password must be at least 8 characters',
              invalidResponse: 'Invalid response from server',
              loginFailed: 'Login failed',
              invalidCredentials: 'Invalid email or password',
              serverUnavailable: 'Server is unavailable',
              timeout: 'Request timeout',
              connectionError: 'Unable to connect to server',
              allFields: 'All fields are required',
              registrationFailed: 'Registration failed',
              emailExists: 'Email already exists',
              validationFailed: 'Validation failed',
            },
          },
        },
      },
    })

    localStorage.clear()
    vi.clearAllMocks()
  })

  describe('login form', () => {
    it('renders login form with email and password fields', () => {
      const wrapper = mount(LoginView, {
        global: {
          plugins: [router, vuetify, i18n],
        },
      })

      expect(wrapper.find('input[type="email"]').exists()).toBe(true)
      expect(wrapper.find('input[type="password"]').exists()).toBe(true)
    })

    it('shows error when email is missing', async () => {
      const wrapper = mount(LoginView, {
        global: {
          plugins: [router, vuetify, i18n],
        },
      })

      await wrapper.find('form').trigger('submit.prevent')
      await wrapper.vm.$nextTick()

      // Check if error message component is displayed
      expect(wrapper.findComponent(ErrorMessage).exists()).toBe(true)
    })

    it('shows error when password is too short', async () => {
      const wrapper = mount(LoginView, {
        global: {
          plugins: [router, vuetify, i18n],
        },
      })

      wrapper.vm.email = 'test@example.com'
      wrapper.vm.password = 'short'

      await wrapper.find('form').trigger('submit.prevent')
      await wrapper.vm.$nextTick()

      expect(wrapper.vm.error).toBeTruthy()
    })

    it('calls API and redirects on successful login', async () => {
      const { api } = await import('@/lib/api')
      const mockResponse = {
        token: 'jwt-token',
        user: { id: '1', email: 'test@example.com', name: 'Test User' },
      }

      vi.mocked(api.post).mockResolvedValue(mockResponse)

      const wrapper = mount(LoginView, {
        global: {
          plugins: [router, vuetify, i18n],
        },
      })

      wrapper.vm.email = 'test@example.com'
      wrapper.vm.password = 'password123'

      await wrapper.find('form').trigger('submit.prevent')
      await wrapper.vm.$nextTick()

      expect(api.post).toHaveBeenCalledWith('/api/v1/auth/login', {
        email: 'test@example.com',
        password: 'password123',
      })
    })

    it('shows safe error message on invalid credentials', async () => {
      const { api } = await import('@/lib/api')
      
      vi.mocked(api.post).mockRejectedValue(new Error('Invalid email or password'))

      const wrapper = mount(LoginView, {
        global: {
          plugins: [router, vuetify, i18n],
        },
      })

      wrapper.vm.email = 'test@example.com'
      wrapper.vm.password = 'wrongpassword'

      await wrapper.find('form').trigger('submit.prevent')
      await wrapper.vm.$nextTick()

      expect(wrapper.vm.error).toBeTruthy()
      expect(wrapper.vm.loading).toBe(false)
    })

    it('shows connection error when server is unavailable', async () => {
      const { api } = await import('@/lib/api')
      
      vi.mocked(api.post).mockRejectedValue(new Error('Unable to connect'))

      const wrapper = mount(LoginView, {
        global: {
          plugins: [router, vuetify, i18n],
        },
      })

      wrapper.vm.email = 'test@example.com'
      wrapper.vm.password = 'password123'

      await wrapper.find('form').trigger('submit.prevent')
      await wrapper.vm.$nextTick()

      expect(wrapper.vm.error).toBeTruthy()
    })
  })

  describe('register form', () => {
    it('switches to register tab', async () => {
      const wrapper = mount(LoginView, {
        global: {
          plugins: [router, vuetify, i18n],
        },
      })

      wrapper.vm.activeTab = 'register'
      await wrapper.vm.$nextTick()

      expect(wrapper.vm.activeTab).toBe('register')
    })

    it('shows error when name is missing', async () => {
      const wrapper = mount(LoginView, {
        global: {
          plugins: [router, vuetify, i18n],
        },
      })

      wrapper.vm.activeTab = 'register'
      wrapper.vm.email = 'test@example.com'
      wrapper.vm.password = 'password123'

      await wrapper.vm.handleRegister()

      expect(wrapper.vm.error).toBeTruthy()
    })

    it('calls API on successful registration', async () => {
      const { api } = await import('@/lib/api')
      const mockResponse = {
        token: 'jwt-token',
        user: { id: '1', email: 'test@example.com', name: 'Test User' },
      }

      vi.mocked(api.post).mockResolvedValue(mockResponse)

      const wrapper = mount(LoginView, {
        global: {
          plugins: [router, vuetify, i18n],
        },
      })

      wrapper.vm.activeTab = 'register'
      wrapper.vm.email = 'test@example.com'
      wrapper.vm.password = 'password123'
      wrapper.vm.name = 'Test User'

      await wrapper.vm.handleRegister()

      expect(api.post).toHaveBeenCalledWith('/api/v1/auth/register', {
        email: 'test@example.com',
        password: 'password123',
        name: 'Test User',
      })
    })

    it('shows error when email already exists', async () => {
      const { api } = await import('@/lib/api')
      
      vi.mocked(api.post).mockRejectedValue(new Error('User with this email already exists'))

      const wrapper = mount(LoginView, {
        global: {
          plugins: [router, vuetify, i18n],
        },
      })

      wrapper.vm.activeTab = 'register'
      wrapper.vm.email = 'existing@example.com'
      wrapper.vm.password = 'password123'
      wrapper.vm.name = 'Test User'

      await wrapper.vm.handleRegister()

      expect(wrapper.vm.error).toBeTruthy()
    })
  })

  describe('loading state', () => {
    it('disables form during submission', async () => {
      const { api } = await import('@/lib/api')
      
      let resolveLogin: (value: any) => void
      const loginPromise = new Promise((resolve) => {
        resolveLogin = resolve
      })
      
      vi.mocked(api.post).mockReturnValue(loginPromise as any)

      const wrapper = mount(LoginView, {
        global: {
          plugins: [router, vuetify, i18n],
        },
      })

      wrapper.vm.email = 'test@example.com'
      wrapper.vm.password = 'password123'

      const submitPromise = wrapper.find('form').trigger('submit.prevent')

      await wrapper.vm.$nextTick()
      expect(wrapper.vm.loading).toBe(true)

      resolveLogin!({ token: 'token', user: {} })
      await submitPromise
      await wrapper.vm.$nextTick()

      expect(wrapper.vm.loading).toBe(false)
    })
  })
})
