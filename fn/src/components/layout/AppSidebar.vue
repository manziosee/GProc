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
        :class="{ active: currentPage === item.key }"
        @click="$emit('navigate', item.key)"
      >
        <span class="nav-icon">{{ item.icon }}</span>
        <span v-if="!collapsed" class="nav-label">{{ item.label }}</span>
      </div>
    </nav>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  collapsed: boolean
  currentPage: string
}>()

defineEmits<{
  toggle: []
  navigate: [page: string]
}>()

const menuItems = [
  { key: 'dashboard', label: 'Dashboard', icon: '📊' },
  { key: 'processes', label: 'Processes', icon: '⚡' },
  { key: 'logs', label: 'Logs', icon: '📋' },
  { key: 'monitoring', label: 'Monitoring', icon: '📈' },
  { key: 'config', label: 'Config', icon: '🔧' },
  { key: 'scheduler', label: 'Scheduler', icon: '⏰' },
  { key: 'cli', label: 'CLI Reference', icon: '💻' },
  { key: 'users', label: 'Users', icon: '👥' },
  { key: 'settings', label: 'Settings', icon: '⚙️' }
]
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