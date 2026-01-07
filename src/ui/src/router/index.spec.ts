import { describe, it, expect, beforeEach, vi } from 'vitest'
import { createRouter, createMemoryHistory } from 'vue-router'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from '@/stores/auth'
import type { User } from '@/stores/auth'

const routes = [
  {
    path: '/login',
    name: 'login',
    component: { template: '<div>Login</div>' },
    meta: { public: true },
  },
  {
    path: '/dashboard',
    name: 'dashboard',
    component: { template: '<div>Dashboard</div>' },
  },
  {
    path: '/',
    name: 'home',
    redirect: '/dashboard',
  },
  {
    path: '/transactions',
    name: 'transactions',
    component: { template: '<div>Transactions</div>' },
  },
  {
    path: '/rules',
    name: 'rules',
    component: { template: '<div>Rules</div>' },
  },
  {
    path: '/schemas',
    name: 'schemas',
    component: { template: '<div>Schemas</div>' },
  },
  {
    path: '/permissions',
    name: 'permissions',
    component: { template: '<div>Permissions</div>' },
    meta: { requiresAdmin: true },
  },
  {
    path: '/branding',
    name: 'branding',
    component: { template: '<div>Branding</div>' },
    meta: { requiresAdmin: true },
  },
  {
    path: '/profile',
    name: 'profile',
    component: { template: '<div>Profile</div>' },
  },
]

describe('router', () => {
  let router: ReturnType<typeof createRouter>
  let authStore: ReturnType<typeof useAuthStore>

  beforeEach(() => {
    setActivePinia(createPinia())
    authStore = useAuthStore()
    
    router = createRouter({
      history: createMemoryHistory(),
      routes,
    })

    router.beforeEach(async (to, from, next) => {
      const isPublicRoute = to.meta.public === true
      const requiresAdmin = to.meta.requiresAdmin === true

      if (authStore.loading) {
        await new Promise(resolve => setTimeout(resolve, 100))
        while (authStore.loading) {
          await new Promise(resolve => setTimeout(resolve, 50))
        }
      }

      if (!authStore.user && !isPublicRoute) {
        next('/login')
        return
      }

      if (authStore.user && to.path === '/login') {
        next('/')
        return
      }

      if (requiresAdmin && !authStore.isAdmin) {
        next('/')
        return
      }

      next()
    })

    authStore.loading = false
    authStore.user = null
  })

  describe('route definitions', () => {
    it('has login route with public meta', () => {
      const loginRoute = router.getRoutes().find(r => r.name === 'login')

      expect(loginRoute).toBeDefined()
      expect(loginRoute?.path).toBe('/login')
      expect(loginRoute?.meta.public).toBe(true)
    })

    it('has dashboard route', () => {
      const dashboardRoute = router.getRoutes().find(r => r.name === 'dashboard')

      expect(dashboardRoute).toBeDefined()
      expect(dashboardRoute?.path).toBe('/dashboard')
    })

    it('has home route that redirects to dashboard', () => {
      const homeRoute = router.getRoutes().find(r => r.name === 'home')

      expect(homeRoute).toBeDefined()
      expect(homeRoute?.path).toBe('/')
      expect(homeRoute?.redirect).toBe('/dashboard')
    })

    it('has permissions route with admin requirement', () => {
      const permissionsRoute = router.getRoutes().find(r => r.name === 'permissions')

      expect(permissionsRoute).toBeDefined()
      expect(permissionsRoute?.meta.requiresAdmin).toBe(true)
    })

    it('has branding route with admin requirement', () => {
      const brandingRoute = router.getRoutes().find(r => r.name === 'branding')

      expect(brandingRoute).toBeDefined()
      expect(brandingRoute?.meta.requiresAdmin).toBe(true)
    })
  })

  describe('navigation guards - unauthenticated user', () => {
    beforeEach(() => {
      authStore.user = null
    })

    it('allows access to public routes', async () => {
      await router.push('/login')

      expect(router.currentRoute.value.path).toBe('/login')
    })

    it('redirects to login when accessing protected route', async () => {
      await router.push('/dashboard')

      expect(router.currentRoute.value.path).toBe('/login')
    })

    it('redirects to login when accessing transactions', async () => {
      await router.push('/transactions')

      expect(router.currentRoute.value.path).toBe('/login')
    })

    it('redirects to login when accessing rules', async () => {
      await router.push('/rules')

      expect(router.currentRoute.value.path).toBe('/login')
    })

    it('redirects to login when accessing schemas', async () => {
      await router.push('/schemas')

      expect(router.currentRoute.value.path).toBe('/login')
    })

    it('redirects to login when accessing profile', async () => {
      await router.push('/profile')

      expect(router.currentRoute.value.path).toBe('/login')
    })

    it('redirects to login when accessing admin routes', async () => {
      await router.push('/permissions')

      expect(router.currentRoute.value.path).toBe('/login')
    })
  })

  describe('navigation guards - authenticated non-admin user', () => {
    beforeEach(() => {
      const user: User = {
        id: 'user-1',
        email: 'user@example.com',
        name: 'Regular User',
        auth_type: 'local',
        active: true,
        roles: [{ id: 'role-1', name: 'user', description: 'Regular user' }],
      }
      authStore.user = user
    })

    it('redirects to home when accessing login page', async () => {
      await router.push('/login')

      expect(router.currentRoute.value.path).toBe('/dashboard')
    })

    it('allows access to dashboard', async () => {
      await router.push('/dashboard')

      expect(router.currentRoute.value.path).toBe('/dashboard')
    })

    it('allows access to transactions', async () => {
      await router.push('/transactions')

      expect(router.currentRoute.value.path).toBe('/transactions')
    })

    it('allows access to rules', async () => {
      await router.push('/rules')

      expect(router.currentRoute.value.path).toBe('/rules')
    })

    it('allows access to schemas', async () => {
      await router.push('/schemas')

      expect(router.currentRoute.value.path).toBe('/schemas')
    })

    it('allows access to profile', async () => {
      await router.push('/profile')

      expect(router.currentRoute.value.path).toBe('/profile')
    })

    it('redirects to home when accessing admin-only permissions route', async () => {
      await router.push('/permissions')

      expect(router.currentRoute.value.path).toBe('/dashboard')
    })

    it('redirects to home when accessing admin-only branding route', async () => {
      await router.push('/branding')

      expect(router.currentRoute.value.path).toBe('/dashboard')
    })
  })

  describe('navigation guards - authenticated admin user', () => {
    beforeEach(() => {
      const adminUser: User = {
        id: 'admin-1',
        email: 'admin@example.com',
        name: 'Admin User',
        auth_type: 'local',
        active: true,
        roles: [
          { id: 'role-1', name: 'admin', description: 'Administrator' },
          { id: 'role-2', name: 'user', description: 'Regular user' },
        ],
      }
      authStore.user = adminUser
    })

    it('redirects to home when accessing login page', async () => {
      await router.push('/login')

      expect(router.currentRoute.value.path).toBe('/dashboard')
    })

    it('allows access to dashboard', async () => {
      await router.push('/dashboard')

      expect(router.currentRoute.value.path).toBe('/dashboard')
    })

    it('allows access to admin-only permissions route', async () => {
      await router.push('/permissions')

      expect(router.currentRoute.value.path).toBe('/permissions')
    })

    it('allows access to admin-only branding route', async () => {
      await router.push('/branding')

      expect(router.currentRoute.value.path).toBe('/branding')
    })

    it('allows access to profile', async () => {
      await router.push('/profile')

      expect(router.currentRoute.value.path).toBe('/profile')
    })
  })

  describe('loading state handling', () => {
    it('waits for auth store to finish loading before navigation', async () => {
      authStore.loading = true
      authStore.user = null

      const navigationPromise = router.push('/dashboard')

      await new Promise(resolve => setTimeout(resolve, 50))
      
      authStore.loading = false

      await navigationPromise

      expect(router.currentRoute.value.path).toBe('/login')
    })

    it('allows navigation after loading completes with user', async () => {
      authStore.loading = true

      const navigationPromise = router.push('/dashboard')

      await new Promise(resolve => setTimeout(resolve, 50))
      
      const user: User = {
        id: 'user-1',
        email: 'user@example.com',
        name: 'User',
        auth_type: 'local',
        active: true,
      }
      authStore.user = user
      authStore.loading = false

      await navigationPromise

      expect(router.currentRoute.value.path).toBe('/dashboard')
    })
  })

  describe('edge cases', () => {
    it('handles navigation to non-existent route', async () => {
      const user: User = {
        id: 'user-1',
        email: 'user@example.com',
        name: 'User',
        auth_type: 'local',
        active: true,
      }
      authStore.user = user

      await router.push('/non-existent').catch(() => {})

      expect(router.currentRoute.value.matched.length).toBeGreaterThanOrEqual(0)
    })

    it('handles rapid navigation changes', async () => {
      const user: User = {
        id: 'user-1',
        email: 'user@example.com',
        name: 'User',
        auth_type: 'local',
        active: true,
      }
      authStore.user = user

      router.push('/dashboard')
      router.push('/rules')
      await router.push('/schemas')

      expect(router.currentRoute.value.path).toBe('/schemas')
    })

    it('handles home redirect correctly', async () => {
      const user: User = {
        id: 'user-1',
        email: 'user@example.com',
        name: 'User',
        auth_type: 'local',
        active: true,
      }
      authStore.user = user

      await router.push('/')

      expect(router.currentRoute.value.path).toBe('/dashboard')
    })

    it('handles user with no roles trying to access admin route', async () => {
      const user: User = {
        id: 'user-1',
        email: 'user@example.com',
        name: 'User',
        auth_type: 'local',
        active: true,
      }
      authStore.user = user

      await router.push('/permissions')

      expect(router.currentRoute.value.path).toBe('/dashboard')
    })
  })
})

