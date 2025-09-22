<template>
  <div class="sidebar" :class="{ collapsed }">
    <div class="sidebar-header">
      <div class="logo-container">
        <img src="/logo.svg" alt="GProc" class="logo" />
        <div v-if="!collapsed" class="brand-text">
          <span class="brand-name">GProc</span>
          <span class="brand-tagline">Enterprise</span>
        </div>
      </div>
    </div>
    
    <nav class="sidebar-nav">
      <div 
        v-for="item in menuItems" 
        :key="item.key"
        class="nav-item"
        :class="{ active: route.name === item.key }"
        @click="navigateTo(item)"
      >
        <span class="nav-icon">{{ item.icon }}</span>
        <span v-if="!collapsed" class="nav-label">{{ item.label }}</span>
      </div>
    </nav>
  </div>
</template>

<script setup lang="ts">
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '../../stores/auth'

defineProps<{
  collapsed: boolean
}>()

defineEmits<{
  toggle: []
}>()

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const menuItems = [
  { key: 'Dashboard', label: 'Dashboard', icon: '📊', route: '/' },
  { key: 'ProcessManagement', label: 'Processes', icon: '⚡', route: '/processes' },
  { key: 'ClusterManagement', label: 'Cluster', icon: '🔗', route: '/cluster' },
  { key: 'Monitoring', label: 'Monitoring', icon: '📈', route: '/monitoring' },
  { key: 'Deployments', label: 'Deployments', icon: '🚀', route: '/deployments' },
  { key: 'Scheduler', label: 'Scheduler', icon: '⏰', route: '/scheduler' },
  { key: 'Security', label: 'Security', icon: '🔒', route: '/security' },
  { key: 'LogsViewer', label: 'Logs', icon: '📋', route: '/logs' },
  { key: 'LanguageProbes', label: 'Probes', icon: '🔍', route: '/probes' },
  { key: 'Templates', label: 'Templates', icon: '📄', route: '/templates' },
  { key: 'AuditLogs', label: 'Audit', icon: '📝', route: '/audit' },
  { key: 'BackupRestore', label: 'Backup', icon: '💾', route: '/backup' },
  { key: 'SecretsManagement', label: 'Secrets', icon: '🔐', route: '/secrets' },
  { key: 'Settings', label: 'Settings', icon: '⚙️', route: '/settings' }
]

const navigateTo = (item: any) => {
  router.push(item.route)
}
</script>

<style scoped>
.sidebar {
  width: 240px;
  background: var(--n-color-embedded);
  border-right: 1px solid var(--n-border-color);
  transition: width 0.3s ease;
  position: fixed;
  left: 0;
  top: 0;
  height: 100vh;
  z-index: 100;
}

.sidebar.collapsed {
  width: 64px;
}

.sidebar-header {
  padding: 24px 20px;
  border-bottom: 1px solid var(--n-border-color);
}

.logo-container {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo {
  width: 32px;
  height: 32px;
  flex-shrink: 0;
}

.brand-text {
  display: flex;
  flex-direction: column;
  line-height: 1.2;
}

.brand-name {
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--n-text-color);
  letter-spacing: -0.025em;
}

.brand-tagline {
  font-size: 0.75rem;
  color: var(--n-text-color-2);
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.collapsed .logo-container {
  justify-content: center;
}

.sidebar-nav {
  padding: 20px 0;
}

.nav-item {
  display: flex;
  align-items: center;
  padding: 12px 20px;
  cursor: pointer;
  transition: background 0.2s ease;
}

.nav-item:hover {
  background: var(--n-color-hover);
}

.nav-item.active {
  background: var(--n-primary-color-hover);
  color: var(--n-primary-color);
}

.nav-icon {
  margin-right: 12px;
  font-size: 18px;
}

.collapsed .nav-item {
  justify-content: center;
}

.collapsed .nav-icon {
  margin-right: 0;
}
</style>