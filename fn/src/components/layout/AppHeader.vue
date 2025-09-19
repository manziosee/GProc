<template>
  <header class="app-header">
    <div class="header-left">
      <button @click="$emit('toggle-sidebar')" class="sidebar-toggle">
        ☰
      </button>
      <div class="header-brand">
        <img src="/logo.svg" alt="GProc" class="header-logo" />
        <div class="header-text">
          <h1 class="header-title">GProc</h1>
          <span class="header-subtitle">Enterprise Process Manager</span>
        </div>
      </div>
    </div>
    
    <div class="header-right">
      <button @click="$emit('toggle-theme')" class="theme-toggle">
        {{ isDark ? '☀️' : '🌙' }}
      </button>
      <div class="user-menu">
        <span>Admin</span>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
defineProps<{
  isDark: boolean
}>()

defineEmits<{
  'toggle-sidebar': []
  'toggle-theme': []
}>()
</script>

<style scoped>
.app-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 24px;
  height: 64px;
  background: var(--n-color-embedded);
  border-bottom: 1px solid var(--n-border-color);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.header-brand {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-logo {
  width: 28px;
  height: 28px;
}

.header-text {
  display: flex;
  flex-direction: column;
  line-height: 1.1;
}

.header-title {
  font-size: 1.125rem;
  font-weight: 700;
  color: var(--n-text-color);
  margin: 0;
  letter-spacing: -0.025em;
}

.header-subtitle {
  font-size: 0.75rem;
  color: var(--n-text-color-2);
  font-weight: 500;
}

.sidebar-toggle {
  background: none;
  border: none;
  font-size: 18px;
  cursor: pointer;
  padding: 8px;
  border-radius: 4px;
  color: var(--n-text-color);
}

.sidebar-toggle:hover {
  background: var(--n-color-hover);
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.theme-toggle {
  background: none;
  border: none;
  font-size: 18px;
  cursor: pointer;
  padding: 8px;
  border-radius: 4px;
}

.theme-toggle:hover {
  background: var(--n-color-hover);
}

.user-menu {
  padding: 8px 16px;
  background: var(--n-color);
  border-radius: 20px;
  border: 1px solid var(--n-border-color);
}
</style>