<template>
  <n-config-provider>
    <n-global-style />
    <n-message-provider>
      <n-dialog-provider>
        <n-notification-provider>
          <div id="app">
            <AppLayout v-if="authStore.isAuthenticated" />
            <router-view v-else />
          </div>
        </n-notification-provider>
      </n-dialog-provider>
    </n-message-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { NConfigProvider, NGlobalStyle, NMessageProvider, NDialogProvider, NNotificationProvider } from 'naive-ui'
import { useRouter } from 'vue-router'
import { useAuthStore } from './stores/auth'
import { useWebSocketStore } from './stores/websocket'
import AppLayout from './components/layout/AppLayout.vue'

const router = useRouter()
const authStore = useAuthStore()
const wsStore = useWebSocketStore()

onMounted(() => {
  // Initialize WebSocket connection if authenticated
  if (authStore.isAuthenticated && authStore.token) {
    wsStore.connect(authStore.token)
  }
})
</script>

<style scoped>
#app {
  height: 100vh;
  width: 100vw;
}
</style>