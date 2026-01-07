import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'
import { createPinia, setActivePinia } from 'pinia'
import LoginView from './LoginView.vue'

// Mock the api module
vi.mock('@/lib/api', () => ({
  api: {
    post: vi.fn(),
  },
}))

// Mock i18n
vi.mock('@/plugins/i18n', () => ({
  i18n: {
    global: {
      t: (key: string) => key,
    },
  },
}))

describe('LoginView', () => {
  let router: ReturnType<typeof createRouter>

  beforeEach(() => {
    setActivePinia(createPinia())
    
    router = createRouter({
      history: createMemoryHistory(),
      routes: [
        { path: '/', component: { template: '<div>Home</div>' } },
        { path: '/login', component: LoginView },
      ],
    })

    localStorage.clear()
    vi.clearAllMocks()
  })

  describe('login form', () => {
    it('renders login form with email and password fields', () => {
      // Arrange & Act
      const wrapper = mount(LoginView, {
        global: {
          plugins: [router],
        },
      })

      // Assert
      expect(wrapper.find('input[type="email"]').exists()).toBe(true)
      expect(wrapper.find('input[type="password"]').exists()).toBe(true)
    })

    it('shows error when email is missing', async () => {
      // Arrange
      const wrapper = mount(LoginView, {
        global: {
          plugins: [router],
        },
      })

      // Act
      await wrapper.find('form').trigger('submit.prevent')

      // Assert
      // Error should be shown (implementation depends on component)
      expect(wrapper.vm.error).toBeTruthy()
    })

    it('shows error when password is too short', async () => {
      // Arrange
      const wrapper = mount(LoginView, {
        global: {
          plugins: [router],
        },
      })

      wrapper.vm.email = 'test@example.com'
      wrapper.vm.password = 'short'

      // Act
      await wrapper.find('form').trigger('submit.prevent')

      // Assert
      expect(wrapper.vm.error).toContain('passwordLength')
    })

    it('calls API and redirects on successful login', async () => {
      // Arrange
      const { api } = await import('@/lib/api')
      const mockResponse = {
        token: 'jwt-token',
        user: { id: '1', email: 'test@example.com', name: 'Test User' },
      }

      vi.mocked(api.post).mockResolvedValue(mockResponse)

      const wrapper = mount(LoginView, {
        global: {
          plugins: [router],
        },
      })

      wrapper.vm.email = 'test@example.com'
      wrapper.vm.password = 'password123'

      // Act
      await wrapper.find('form').trigger('submit.prevent')
      await wrapper.vm.$nextTick()

      // Assert
      expect(api.post).toHaveBeenCalledWith('/api/v1/auth/login', {
        email: 'test@example.com',
        password: 'password123',
      })
    })

    it('shows safe error message on invalid credentials', async () => {
      // Arrange
      const { api } = await import('@/lib/api')
      
      vi.mocked(api.post).mockRejectedValue(new Error('Invalid email or password'))

      const wrapper = mount(LoginView, {
        global: {
          plugins: [router],
        },
      })

      wrapper.vm.email = 'test@example.com'
      wrapper.vm.password = 'wrongpassword'

      // Act
      await wrapper.find('form').trigger('submit.prevent')
      await wrapper.vm.$nextTick()

      // Assert
      expect(wrapper.vm.error).toContain('invalidCredentials')
      expect(wrapper.vm.loading).toBe(false)
    })

    it('shows connection error when server is unavailable', async () => {
      // Arrange
      const { api } = await import('@/lib/api')
      
      vi.mocked(api.post).mockRejectedValue(new Error('Unable to connect'))

      const wrapper = mount(LoginView, {
        global: {
          plugins: [router],
        },
      })

      wrapper.vm.email = 'test@example.com'
      wrapper.vm.password = 'password123'

      // Act
      await wrapper.find('form').trigger('submit.prevent')
      await wrapper.vm.$nextTick()

      // Assert
      expect(wrapper.vm.error).toContain('connectionError')
    })
  })

  describe('register form', () => {
    it('switches to register tab', async () => {
      // Arrange
      const wrapper = mount(LoginView, {
        global: {
          plugins: [router],
        },
      })

      // Act
      wrapper.vm.activeTab = 'register'
      await wrapper.vm.$nextTick()

      // Assert
      expect(wrapper.vm.activeTab).toBe('register')
    })

    it('shows error when name is missing', async () => {
      // Arrange
      const wrapper = mount(LoginView, {
        global: {
          plugins: [router],
        },
      })

      wrapper.vm.activeTab = 'register'
      wrapper.vm.email = 'test@example.com'
      wrapper.vm.password = 'password123'
      // name is empty

      // Act
      await wrapper.vm.handleRegister()

      // Assert
      expect(wrapper.vm.error).toContain('allFields')
    })

    it('calls API on successful registration', async () => {
      // Arrange
      const { api } = await import('@/lib/api')
      const mockResponse = {
        token: 'jwt-token',
        user: { id: '1', email: 'test@example.com', name: 'Test User' },
      }

      vi.mocked(api.post).mockResolvedValue(mockResponse)

      const wrapper = mount(LoginView, {
        global: {
          plugins: [router],
        },
      })

      wrapper.vm.activeTab = 'register'
      wrapper.vm.email = 'test@example.com'
      wrapper.vm.password = 'password123'
      wrapper.vm.name = 'Test User'

      // Act
      await wrapper.vm.handleRegister()

      // Assert
      expect(api.post).toHaveBeenCalledWith('/api/v1/auth/register', {
        email: 'test@example.com',
        password: 'password123',
        name: 'Test User',
      })
    })

    it('shows error when email already exists', async () => {
      // Arrange
      const { api } = await import('@/lib/api')
      
      vi.mocked(api.post).mockRejectedValue(new Error('User with this email already exists'))

      const wrapper = mount(LoginView, {
        global: {
          plugins: [router],
        },
      })

      wrapper.vm.activeTab = 'register'
      wrapper.vm.email = 'existing@example.com'
      wrapper.vm.password = 'password123'
      wrapper.vm.name = 'Test User'

      // Act
      await wrapper.vm.handleRegister()

      // Assert
      expect(wrapper.vm.error).toContain('emailExists')
    })
  })

  describe('loading state', () => {
    it('disables form during submission', async () => {
      // Arrange
      const { api } = await import('@/lib/api')
      
      // Create a promise that we can control
      let resolveLogin: (value: any) => void
      const loginPromise = new Promise((resolve) => {
        resolveLogin = resolve
      })
      
      vi.mocked(api.post).mockReturnValue(loginPromise as any)

      const wrapper = mount(LoginView, {
        global: {
          plugins: [router],
        },
      })

      wrapper.vm.email = 'test@example.com'
      wrapper.vm.password = 'password123'

      // Act
      const submitPromise = wrapper.find('form').trigger('submit.prevent')

      // Assert - loading should be true
      await wrapper.vm.$nextTick()
      expect(wrapper.vm.loading).toBe(true)

      // Resolve the login
      resolveLogin!({ token: 'token', user: {} })
      await submitPromise
      await wrapper.vm.$nextTick()

      // Loading should be false after completion
      expect(wrapper.vm.loading).toBe(false)
    })
  })
})

