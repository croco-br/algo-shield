import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue'),
      meta: { public: true },
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: () => import('@/views/DashboardView.vue'),
    },
    {
      path: '/',
      name: 'home',
      redirect: '/dashboard',
    },
    {
      path: '/transactions',
      name: 'transactions',
      component: () => import('@/views/DashboardView.vue'),
    },
    {
      path: '/risk-analysis',
      name: 'risk-analysis',
      component: () => import('@/views/DashboardView.vue'),
    },
    {
      path: '/reports',
      name: 'reports',
      component: () => import('@/views/DashboardView.vue'),
    },
    {
      path: '/compliance',
      name: 'compliance',
      component: () => import('@/views/DashboardView.vue'),
    },
    {
      path: '/rules',
      name: 'rules',
      component: () => import('@/views/RulesView.vue'),
    },
    {
      path: '/permissions',
      name: 'permissions',
      component: () => import('@/views/PermissionsView.vue'),
      meta: { requiresAdmin: true },
    },
    {
      path: '/synthetic-test',
      name: 'synthetic-test',
      component: () => import('@/views/SyntheticTestView.vue'),
    },
    {
      path: '/health',
      name: 'health',
      component: () => import('@/views/HealthView.vue'),
      meta: { public: true },
    },
  ],
})

router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()
  const isPublicRoute = to.meta.public === true
  const requiresAdmin = to.meta.requiresAdmin === true

  // Wait for auth to load if still loading
  if (authStore.loading) {
    // Wait a bit and check again
    await new Promise(resolve => setTimeout(resolve, 100))
    // If still loading, wait for the store to finish loading
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

export default router
