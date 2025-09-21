<template>
  <div class="settings">
    <div class="page-header">
      <h1 class="page-title">Settings</h1>
      <p class="page-subtitle">Configure your GProc environment</p>
    </div>

    <div class="settings-grid">
      <div class="settings-section">
        <div class="section-header">
          <h2 class="section-title">🎨 Appearance</h2>
          <p class="section-description">Customize the look and feel</p>
        </div>
        <div class="setting-item">
          <div class="setting-info">
            <label class="setting-label">Theme</label>
            <span class="setting-description">Choose your preferred theme</span>
          </div>
          <n-switch v-model:value="isDarkMode" @update:value="toggleTheme">
            <template #checked>🌙 Dark</template>
            <template #unchecked>☀️ Light</template>
          </n-switch>
        </div>
        <div class="setting-item">
          <div class="setting-info">
            <label class="setting-label">Compact Mode</label>
            <span class="setting-description">Reduce spacing for more content</span>
          </div>
          <n-switch v-model:value="compactMode" />
        </div>
      </div>

      <div class="settings-section">
        <div class="section-header">
          <h2 class="section-title">🔔 Notifications</h2>
          <p class="section-description">Manage your notification preferences</p>
        </div>
        <div class="setting-item">
          <div class="setting-info">
            <label class="setting-label">Email Notifications</label>
            <span class="setting-description">Receive alerts via email</span>
          </div>
          <n-switch v-model:value="emailNotifications" />
        </div>
        <div class="setting-item">
          <div class="setting-info">
            <label class="setting-label">Desktop Notifications</label>
            <span class="setting-description">Show browser notifications</span>
          </div>
          <n-switch v-model:value="desktopNotifications" />
        </div>
        <div class="setting-item">
          <div class="setting-info">
            <label class="setting-label">Notification Email</label>
            <span class="setting-description">Email address for notifications</span>
          </div>
          <n-input v-model:value="notificationEmail" placeholder="admin@company.com" style="width: 250px" />
        </div>
      </div>

      <div class="settings-section">
        <div class="section-header">
          <h2 class="section-title">⚙️ System</h2>
          <p class="section-description">Configure system behavior</p>
        </div>
        <div class="setting-item">
          <div class="setting-info">
            <label class="setting-label">Auto Refresh</label>
            <span class="setting-description">Automatically refresh data</span>
          </div>
          <n-select v-model:value="autoRefresh" :options="refreshOptions" style="width: 150px" />
        </div>
        <div class="setting-item">
          <div class="setting-info">
            <label class="setting-label">Log Retention</label>
            <span class="setting-description">How long to keep logs</span>
          </div>
          <n-select v-model:value="logRetention" :options="retentionOptions" style="width: 150px" />
        </div>
        <div class="setting-item">
          <div class="setting-info">
            <label class="setting-label">Max Log Size</label>
            <span class="setting-description">Maximum size per log file</span>
          </div>
          <n-input-number v-model:value="maxLogSize" :min="1" :max="1000" style="width: 150px">
            <template #suffix>MB</template>
          </n-input-number>
        </div>
      </div>

      <div class="settings-section">
        <div class="section-header">
          <h2 class="section-title">🔒 Security</h2>
          <p class="section-description">Security and authentication settings</p>
        </div>
        <div class="setting-item">
          <div class="setting-info">
            <label class="setting-label">Session Timeout</label>
            <span class="setting-description">Auto logout after inactivity</span>
          </div>
          <n-select v-model:value="sessionTimeout" :options="timeoutOptions" style="width: 150px" />
        </div>
        <div class="setting-item">
          <div class="setting-info">
            <label class="setting-label">Two-Factor Auth</label>
            <span class="setting-description">Enable 2FA for extra security</span>
          </div>
          <n-switch v-model:value="twoFactorAuth" />
        </div>
        <div class="setting-item">
          <div class="setting-info">
            <label class="setting-label">API Access</label>
            <span class="setting-description">Allow API access</span>
          </div>
          <n-switch v-model:value="apiAccess" />
        </div>
      </div>
    </div>

    <div class="settings-actions">
      <n-space>
        <n-button @click="resetSettings">Reset to Defaults</n-button>
        <n-button type="primary" @click="saveSettings">Save Changes</n-button>
      </n-space>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, inject } from 'vue'
import { NSwitch, NInput, NSelect, NInputNumber, NButton, NSpace } from 'naive-ui'

const isDark = inject('isDark') as any
const isDarkMode = ref(isDark?.value || false)
const compactMode = ref(false)
const emailNotifications = ref(true)
const desktopNotifications = ref(false)
const notificationEmail = ref('admin@company.com')
const autoRefresh = ref('30s')
const logRetention = ref('30d')
const maxLogSize = ref(100)
const sessionTimeout = ref('1h')
const twoFactorAuth = ref(false)
const apiAccess = ref(true)

const refreshOptions = [
  { label: '10 seconds', value: '10s' },
  { label: '30 seconds', value: '30s' },
  { label: '1 minute', value: '1m' },
  { label: '5 minutes', value: '5m' },
  { label: 'Disabled', value: 'off' }
]

const retentionOptions = [
  { label: '7 days', value: '7d' },
  { label: '30 days', value: '30d' },
  { label: '90 days', value: '90d' },
  { label: '1 year', value: '1y' }
]

const timeoutOptions = [
  { label: '30 minutes', value: '30m' },
  { label: '1 hour', value: '1h' },
  { label: '4 hours', value: '4h' },
  { label: '8 hours', value: '8h' },
  { label: 'Never', value: 'never' }
]

const toggleTheme = () => {
  if (isDark) {
    isDark.value = isDarkMode.value
    // Apply theme to document
    if (isDark.value) {
      document.documentElement.setAttribute('data-theme', 'dark')
    } else {
      document.documentElement.removeAttribute('data-theme')
    }
  }
}

const saveSettings = () => {
  console.log('Saving settings...')
}

const resetSettings = () => {
  console.log('Resetting settings...')
}
</script>

<style scoped>
.settings {
  padding: 32px;
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 32px;
  padding-bottom: 24px;
  border-bottom: 1px solid var(--n-border-color);
}

.page-title {
  font-size: 2.5rem;
  font-weight: 700;
  color: var(--n-text-color);
  margin: 0 0 8px 0;
  letter-spacing: -0.025em;
}

.page-subtitle {
  font-size: 1.125rem;
  color: var(--n-text-color-2);
  margin: 0;
  font-weight: 400;
}

.settings-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 32px;
  margin-bottom: 32px;
}

.settings-section {
  background: var(--n-color-embedded);
  border: 1px solid var(--n-border-color);
  border-radius: 16px;
  padding: 24px;
  transition: all 0.3s ease;
}

.settings-section:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.section-header {
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid var(--n-border-color);
}

.section-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--n-text-color);
  margin: 0 0 4px 0;
}

.section-description {
  font-size: 0.875rem;
  color: var(--n-text-color-2);
  margin: 0;
}

.setting-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 0;
  border-bottom: 1px solid var(--n-border-color);
}

.setting-item:last-child {
  border-bottom: none;
}

.setting-info {
  flex: 1;
  margin-right: 16px;
}

.setting-label {
  display: block;
  font-weight: 500;
  color: var(--n-text-color);
  margin-bottom: 2px;
}

.setting-description {
  font-size: 0.875rem;
  color: var(--n-text-color-2);
}

.settings-actions {
  display: flex;
  justify-content: flex-end;
  padding-top: 24px;
  border-top: 1px solid var(--n-border-color);
}

@media (max-width: 768px) {
  .settings {
    padding: 16px;
  }
  
  .settings-grid {
    grid-template-columns: 1fr;
    gap: 24px;
  }
  
  .setting-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
  
  .setting-info {
    margin-right: 0;
  }
}
</style>