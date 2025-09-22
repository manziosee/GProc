import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'Landing',
      component: () => import('../pages/Landing.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/login',
      name: 'Login',
      component: () => import('../pages/Login.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/dashboard',
      name: 'Dashboard',
      component: () => import('../pages/Dashboard.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/processes',
      name: 'ProcessManagement',
      component: () => import('../pages/ProcessManagement.vue'),
      meta: { requiresAuth: true, permission: 'process:read' }
    },
    {
      path: '/cluster',
      name: 'ClusterManagement',
      component: () => import('../pages/ClusterManagement.vue'),
      meta: { requiresAuth: true, permission: 'cluster:read' }
    },
    {
      path: '/monitoring',
      name: 'Monitoring',
      component: () => import('../pages/Monitoring.vue'),
      meta: { requiresAuth: true, permission: 'metrics:read' }
    },
    {
      path: '/deployments',
      name: 'Deployments',
      component: () => import('../pages/Deployments.vue'),
      meta: { requiresAuth: true, permission: 'deployment:read' }
    },
    {
      path: '/scheduler',
      name: 'Scheduler',
      component: () => import('../pages/Scheduler.vue'),
      meta: { requiresAuth: true, permission: 'scheduler:read' }
    },
    {
      path: '/security',
      name: 'Security',
      component: () => import('../pages/Security.vue'),
      meta: { requiresAuth: true, permission: 'user:read' }
    },
    {
      path: '/logs',
      name: 'LogsViewer',
      component: () => import('../pages/LogsViewer.vue'),
      meta: { requiresAuth: true, permission: 'process:read' }
    },
    {
      path: '/settings',
      name: 'Settings',
      component: () => import('../pages/Settings.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/probes',
      name: 'LanguageProbes',
      component: () => import('../pages/LanguageProbes.vue'),
      meta: { requiresAuth: true, permission: 'process:read' }
    },
    {
      path: '/templates',
      name: 'Templates',
      component: () => import('../pages/Templates.vue'),
      meta: { requiresAuth: true, permission: 'template:read' }
    },
    {
      path: '/audit',
      name: 'AuditLogs',
      component: () => import('../pages/AuditLogs.vue'),
      meta: { requiresAuth: true, permission: 'audit:read' }
    },
    {
      path: '/backup',
      name: 'BackupRestore',
      component: () => import('../pages/BackupRestore.vue'),
      meta: { requiresAuth: true, permission: 'backup:read' }
    },
    {
      path: '/secrets',
      name: 'SecretsManagement',
      component: () => import('../pages/SecretsManagement.vue'),
      meta: { requiresAuth: true, permission: 'secrets:read' }
    }
  ]
})

router.beforeEach((to, _from, next) => {
  const authStore = useAuthStore()
  
  // Redirect authenticated users from landing to dashboard
  if (to.path === '/' && authStore.isAuthenticated) {
    next('/dashboard')
    return
  }

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/login')
    return
  }
  
  if (to.meta.permission && typeof to.meta.permission === 'string' && !authStore.hasPermission(
    to.meta.permission.split(':')[0],
    to.meta.permission.split(':')[1]
  )) {
    next('/')
    return
  }
  
  next()
})

export default router