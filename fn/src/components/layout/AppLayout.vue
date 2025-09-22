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
            />
            <div class="main-content" :class="{ 'sidebar-collapsed': sidebarCollapsed }">
              <AppHeader 
                @toggle-sidebar="sidebarCollapsed = !sidebarCollapsed"
                @toggle-theme="toggleTheme"
                :is-dark="isDark"
              />
              <div class="page-content">
                <router-view />
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
import AppSidebar from './AppSidebar.vue'
import AppHeader from './AppHeader.vue'

const isDark = useLocalStorage('gproc-theme', false)
const sidebarCollapsed = ref(false)

const toggleTheme = () => {
  isDark.value = !isDark.value
  if (isDark.value) {
    document.documentElement.setAttribute('data-theme', 'dark')
  } else {
    document.documentElement.removeAttribute('data-theme')
  }
}

provide('isDark', isDark)

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

@media (max-width: 768px) {
  .main-content {
    margin-left: 0;
  }
  
  .main-content.sidebar-collapsed {
    margin-left: 0;
  }
}
</style>