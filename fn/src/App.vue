<template>
  <n-config-provider :theme="isDark ? darkTheme : lightTheme">
    <n-global-style />
    <n-message-provider>
      <n-dialog-provider>
        <n-notification-provider>
          <div class="app-container">
            <AppSidebar 
              :collapsed="sidebarCollapsed" 
              @toggle="sidebarCollapsed = !sidebarCollapsed"
              @navigate="currentPage = $event"
              :current-page="currentPage"
            />
            <div class="main-content" :class="{ 'sidebar-collapsed': sidebarCollapsed }">
              <AppHeader 
                @toggle-sidebar="sidebarCollapsed = !sidebarCollapsed"
                @toggle-theme="toggleTheme"
                :is-dark="isDark"
              />
              <div class="page-content">
                <Transition name="page" mode="out-in">
                  <Dashboard v-if="currentPage === 'dashboard'" />
                  <ProcessManagement v-else-if="currentPage === 'processes'" />
                  <LogsViewer v-else-if="currentPage === 'logs'" />
                  <HealthMonitoring v-else-if="currentPage === 'monitoring'" />
                  <LoadBalancer v-else-if="currentPage === 'loadbalancer'" />
                  <ConfigManager v-else-if="currentPage === 'config'" />
                  <ScheduledTasks v-else-if="currentPage === 'scheduler'" />
                  <DaemonStatus v-else-if="currentPage === 'daemon'" />
                  <CLI v-else-if="currentPage === 'cli'" />
                  <UserManagement v-else-if="currentPage === 'users'" />
                  <Settings v-else-if="currentPage === 'settings'" />
                </Transition>
              </div>
            </div>
          </div>
        </n-notification-provider>
      </n-dialog-provider>
    </n-message-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { ref, provide, onMounted } from 'vue'
import { darkTheme, lightTheme, NConfigProvider, NGlobalStyle, NMessageProvider, NDialogProvider, NNotificationProvider } from 'naive-ui'
import { useLocalStorage } from '@vueuse/core'
import AppSidebar from './components/layout/AppSidebar.vue'
import AppHeader from './components/layout/AppHeader.vue'
import Dashboard from './pages/Dashboard.vue'
import ProcessManagement from './pages/ProcessManagement.vue'
import LogsViewer from './pages/LogsViewer.vue'
import HealthMonitoring from './pages/HealthMonitoring.vue'
import LoadBalancer from './pages/LoadBalancer.vue'
import ConfigManager from './pages/ConfigManager.vue'
import ScheduledTasks from './pages/ScheduledTasks.vue'
import DaemonStatus from './pages/DaemonStatus.vue'
import CLI from './pages/CLI.vue'
import UserManagement from './pages/UserManagement.vue'
import Settings from './pages/Settings.vue'

const isDark = useLocalStorage('gproc-theme', false)
const sidebarCollapsed = ref(false)
const currentPage = ref('dashboard')

const toggleTheme = () => {
  isDark.value = !isDark.value
  // Apply theme to document
  if (isDark.value) {
    document.documentElement.setAttribute('data-theme', 'dark')
  } else {
    document.documentElement.removeAttribute('data-theme')
  }
}

// Provide theme state for child components
provide('isDark', isDark)

// Initialize theme on mount
onMounted(() => {
  if (isDark.value) {
    document.documentElement.setAttribute('data-theme', 'dark')
  }
})
</script>

<style scoped>
.app-container {
  display: flex;
  height: 100vh;
  background: var(--n-color);
}

.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  transition: margin-left 0.3s ease;
  margin-left: 240px;
}

.main-content.sidebar-collapsed {
  margin-left: 64px;
}

.page-content {
  flex: 1;
  padding: 24px;
  overflow-y: auto;
  background: var(--n-color);
}

/* Page transitions */
.page-enter-active,
.page-leave-active {
  transition: all 0.3s ease;
}

.page-enter-from {
  opacity: 0;
  transform: translateX(20px);
}

.page-leave-to {
  opacity: 0;
  transform: translateX(-20px);
}

@media (max-width: 768px) {
  .main-content {
    margin-left: 0;
  }
  
  .main-content.sidebar-collapsed {
    margin-left: 0;
  }
}
</style>