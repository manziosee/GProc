<template>
  <div class="daemon-status">
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">Daemon Status</h1>
        <p class="page-subtitle">Monitor and control GProc daemon service</p>
      </div>
      <div class="daemon-controls">
        <n-button :type="daemonStatus.running ? 'warning' : 'primary'" @click="toggleDaemon">
          <template #icon>
            <span>{{ daemonStatus.running ? '⏹️' : '▶️' }}</span>
          </template>
          {{ daemonStatus.running ? 'Stop Daemon' : 'Start Daemon' }}
        </n-button>
        <n-button @click="restartDaemon">
          <template #icon><span>🔄</span></template>
          Restart
        </n-button>
      </div>
    </div>

    <div class="status-overview">
      <div class="status-card" :class="daemonStatus.running ? 'running' : 'stopped'">
        <div class="status-indicator">
          <div class="indicator-dot" :class="daemonStatus.running ? 'active' : 'inactive'"></div>
          <span class="status-text">{{ daemonStatus.running ? 'Running' : 'Stopped' }}</span>
        </div>
        <div class="status-details">
          <div class="detail-item">
            <span class="detail-label">PID:</span>
            <span class="detail-value">{{ daemonStatus.pid || 'N/A' }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">Uptime:</span>
            <span class="detail-value">{{ daemonStatus.uptime || 'N/A' }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">Version:</span>
            <span class="detail-value">{{ daemonStatus.version }}</span>
          </div>
        </div>
      </div>

      <div class="metrics-cards">
        <div class="metric-card">
          <div class="metric-icon">📊</div>
          <div class="metric-info">
            <div class="metric-value">{{ daemonMetrics.processesManaged }}</div>
            <div class="metric-label">Processes Managed</div>
          </div>
        </div>
        
        <div class="metric-card">
          <div class="metric-icon">🔄</div>
          <div class="metric-info">
            <div class="metric-value">{{ daemonMetrics.restarts }}</div>
            <div class="metric-label">Auto Restarts</div>
          </div>
        </div>
        
        <div class="metric-card">
          <div class="metric-icon">💾</div>
          <div class="metric-info">
            <div class="metric-value">{{ daemonMetrics.memoryUsage }}</div>
            <div class="metric-label">Memory Usage</div>
          </div>
        </div>
        
        <div class="metric-card">
          <div class="metric-icon">⚡</div>
          <div class="metric-info">
            <div class="metric-value">{{ daemonMetrics.cpuUsage }}%</div>
            <div class="metric-label">CPU Usage</div>
          </div>
        </div>
      </div>
    </div>

    <div class="daemon-config">
      <h2 class="section-title">Daemon Configuration</h2>
      <div class="config-grid">
        <div class="config-section">
          <h3>Service Settings</h3>
          <div class="config-items">
            <div class="config-item">
              <span class="config-label">Auto Start:</span>
              <n-switch v-model:value="config.autoStart" @update:value="updateConfig" />
            </div>
            <div class="config-item">
              <span class="config-label">Log Level:</span>
              <n-select v-model:value="config.logLevel" :options="logLevelOptions" style="width: 120px" />
            </div>
            <div class="config-item">
              <span class="config-label">Max Processes:</span>
              <n-input-number v-model:value="config.maxProcesses" :min="1" :max="100" style="width: 120px" />
            </div>
          </div>
        </div>

        <div class="config-section">
          <h3>Monitoring</h3>
          <div class="config-items">
            <div class="config-item">
              <span class="config-label">Health Check Interval:</span>
              <n-input-number v-model:value="config.healthCheckInterval" :min="5" :max="300" style="width: 120px">
                <template #suffix>s</template>
              </n-input-number>
            </div>
            <div class="config-item">
              <span class="config-label">Resource Monitoring:</span>
              <n-switch v-model:value="config.resourceMonitoring" />
            </div>
            <div class="config-item">
              <span class="config-label">Web Dashboard Port:</span>
              <n-input-number v-model:value="config.webPort" :min="3000" :max="9999" style="width: 120px" />
            </div>
          </div>
        </div>

        <div class="config-section">
          <h3>Notifications</h3>
          <div class="config-items">
            <div class="config-item">
              <span class="config-label">Email Alerts:</span>
              <n-switch v-model:value="config.emailAlerts" />
            </div>
            <div class="config-item">
              <span class="config-label">Slack Notifications:</span>
              <n-switch v-model:value="config.slackNotifications" />
            </div>
            <div class="config-item">
              <span class="config-label">Alert Threshold:</span>
              <n-select v-model:value="config.alertThreshold" :options="thresholdOptions" style="width: 120px" />
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="daemon-logs">
      <div class="logs-header">
        <h2 class="section-title">Daemon Logs</h2>
        <div class="logs-controls">
          <n-button size="small" @click="clearLogs">Clear</n-button>
          <n-button size="small" @click="downloadLogs">Download</n-button>
          <n-switch v-model:value="autoScroll" size="small">
            <template #checked>Auto Scroll</template>
            <template #unchecked>Manual</template>
          </n-switch>
        </div>
      </div>
      
      <div class="logs-container" ref="logsContainer">
        <div v-for="log in daemonLogs" :key="log.id" class="log-entry" :class="log.level">
          <span class="log-timestamp">{{ formatTimestamp(log.timestamp) }}</span>
          <span class="log-level">{{ log.level.toUpperCase() }}</span>
          <span class="log-message">{{ log.message }}</span>
        </div>
      </div>
    </div>

    <div class="system-commands">
      <h2 class="section-title">System Commands</h2>
      <div class="commands-grid">
        <n-button @click="runCommand('status')" class="command-btn">
          <template #icon><span>📊</span></template>
          Check Status
        </n-button>
        <n-button @click="runCommand('reload')" class="command-btn">
          <template #icon><span>🔄</span></template>
          Reload Config
        </n-button>
        <n-button @click="runCommand('cleanup')" class="command-btn">
          <template #icon><span>🧹</span></template>
          Cleanup Logs
        </n-button>
        <n-button @click="runCommand('backup')" class="command-btn">
          <template #icon><span>💾</span></template>
          Backup State
        </n-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted, nextTick } from 'vue'
import { NButton, NSwitch, NSelect, NInputNumber } from 'naive-ui'

const logsContainer = ref<HTMLElement>()
const autoScroll = ref(true)

const daemonStatus = reactive({
  running: true,
  pid: 12345,
  uptime: '2d 14h 32m',
  version: '1.0.0'
})

const daemonMetrics = reactive({
  processesManaged: 18,
  restarts: 5,
  memoryUsage: '45.2MB',
  cpuUsage: 2.3
})

const config = reactive({
  autoStart: true,
  logLevel: 'info',
  maxProcesses: 50,
  healthCheckInterval: 30,
  resourceMonitoring: true,
  webPort: 3000,
  emailAlerts: false,
  slackNotifications: true,
  alertThreshold: 'medium'
})

const logLevelOptions = [
  { label: 'Debug', value: 'debug' },
  { label: 'Info', value: 'info' },
  { label: 'Warning', value: 'warning' },
  { label: 'Error', value: 'error' }
]

const thresholdOptions = [
  { label: 'Low', value: 'low' },
  { label: 'Medium', value: 'medium' },
  { label: 'High', value: 'high' }
]

const daemonLogs = ref([
  {
    id: 1,
    timestamp: new Date(Date.now() - 300000),
    level: 'info',
    message: 'Daemon started successfully'
  },
  {
    id: 2,
    timestamp: new Date(Date.now() - 240000),
    level: 'info',
    message: 'Process web-server started (PID: 8765)'
  },
  {
    id: 3,
    timestamp: new Date(Date.now() - 180000),
    level: 'warning',
    message: 'Process api-worker high memory usage detected'
  },
  {
    id: 4,
    timestamp: new Date(Date.now() - 120000),
    level: 'info',
    message: 'Health check passed for all processes'
  },
  {
    id: 5,
    timestamp: new Date(Date.now() - 60000),
    level: 'error',
    message: 'Failed to restart process background-job'
  }
])

let logUpdateInterval: number

const formatTimestamp = (timestamp: Date) => {
  return timestamp.toLocaleTimeString('en-US', { 
    hour12: false, 
    hour: '2-digit', 
    minute: '2-digit', 
    second: '2-digit' 
  })
}

const toggleDaemon = () => {
  daemonStatus.running = !daemonStatus.running
  if (daemonStatus.running) {
    daemonStatus.pid = Math.floor(Math.random() * 90000) + 10000
    daemonStatus.uptime = '0d 0h 0m'
    addLog('info', 'Daemon started successfully')
  } else {
    daemonStatus.pid = null
    daemonStatus.uptime = null
    addLog('info', 'Daemon stopped')
  }
}

const restartDaemon = () => {
  addLog('info', 'Restarting daemon...')
  daemonStatus.running = false
  setTimeout(() => {
    daemonStatus.running = true
    daemonStatus.pid = Math.floor(Math.random() * 90000) + 10000
    daemonStatus.uptime = '0d 0h 0m'
    addLog('info', 'Daemon restarted successfully')
  }, 2000)
}

const updateConfig = () => {
  addLog('info', 'Configuration updated')
}

const addLog = (level: string, message: string) => {
  const newLog = {
    id: Date.now(),
    timestamp: new Date(),
    level,
    message
  }
  
  daemonLogs.value.push(newLog)
  
  // Keep only last 100 logs
  if (daemonLogs.value.length > 100) {
    daemonLogs.value.shift()
  }
  
  if (autoScroll.value) {
    nextTick(() => {
      if (logsContainer.value) {
        logsContainer.value.scrollTop = logsContainer.value.scrollHeight
      }
    })
  }
}

const clearLogs = () => {
  daemonLogs.value = []
  addLog('info', 'Logs cleared')
}

const downloadLogs = () => {
  const logText = daemonLogs.value
    .map(log => `[${formatTimestamp(log.timestamp)}] ${log.level.toUpperCase()}: ${log.message}`)
    .join('\n')
  
  const blob = new Blob([logText], { type: 'text/plain' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `gproc-daemon-logs-${new Date().toISOString().split('T')[0]}.txt`
  a.click()
  URL.revokeObjectURL(url)
}

const runCommand = (command: string) => {
  addLog('info', `Executing command: ${command}`)
  
  setTimeout(() => {
    switch (command) {
      case 'status':
        addLog('info', 'Status check completed - All systems operational')
        break
      case 'reload':
        addLog('info', 'Configuration reloaded successfully')
        break
      case 'cleanup':
        addLog('info', 'Log cleanup completed - 15 old files removed')
        break
      case 'backup':
        addLog('info', 'State backup created successfully')
        break
    }
  }, 1000)
}

onMounted(() => {
  // Simulate periodic log updates
  logUpdateInterval = setInterval(() => {
    if (daemonStatus.running) {
      const messages = [
        'Health check completed for all processes',
        'Process metrics updated',
        'Configuration sync completed',
        'Resource monitoring update'
      ]
      const randomMessage = messages[Math.floor(Math.random() * messages.length)]
      addLog('info', randomMessage)
    }
  }, 10000)
})

onUnmounted(() => {
  if (logUpdateInterval) {
    clearInterval(logUpdateInterval)
  }
})
</script>

<style scoped>
.daemon-status {
  padding: 32px;
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
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

.daemon-controls {
  display: flex;
  gap: 12px;
}

.status-overview {
  display: grid;
  grid-template-columns: 1fr 2fr;
  gap: 24px;
  margin-bottom: 40px;
}

.status-card {
  background: var(--n-color-embedded);
  border: 1px solid var(--n-border-color);
  border-radius: 16px;
  padding: 24px;
  position: relative;
  overflow: hidden;
}

.status-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: var(--status-color);
}

.status-card.running {
  --status-color: var(--gproc-success);
}

.status-card.stopped {
  --status-color: var(--gproc-error);
}

.status-indicator {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 20px;
}

.indicator-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
}

.indicator-dot.active {
  background: var(--gproc-success);
  box-shadow: 0 0 0 2px rgba(16, 185, 129, 0.2);
  animation: pulse 2s infinite;
}

.indicator-dot.inactive {
  background: var(--gproc-error);
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.status-text {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--n-text-color);
}

.status-details {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.detail-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.detail-label {
  color: var(--n-text-color-2);
  font-weight: 500;
}

.detail-value {
  color: var(--n-text-color);
  font-family: monospace;
}

.metrics-cards {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.metric-card {
  background: var(--n-color);
  border: 1px solid var(--n-border-color);
  border-radius: 12px;
  padding: 16px;
  display: flex;
  align-items: center;
  gap: 12px;
}

.metric-icon {
  font-size: 1.5rem;
}

.metric-value {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--n-text-color);
  line-height: 1;
}

.metric-label {
  font-size: 0.875rem;
  color: var(--n-text-color-2);
}

.daemon-config {
  margin-bottom: 40px;
}

.section-title {
  font-size: 1.5rem;
  font-weight: 600;
  color: var(--n-text-color);
  margin: 0 0 20px 0;
}

.config-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 24px;
}

.config-section {
  background: var(--n-color-embedded);
  border: 1px solid var(--n-border-color);
  border-radius: 12px;
  padding: 20px;
}

.config-section h3 {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--n-text-color);
  margin: 0 0 16px 0;
}

.config-items {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.config-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.config-label {
  color: var(--n-text-color-2);
  font-weight: 500;
}

.daemon-logs {
  margin-bottom: 40px;
}

.logs-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.logs-controls {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logs-container {
  background: var(--n-color-embedded);
  border: 1px solid var(--n-border-color);
  border-radius: 12px;
  height: 300px;
  overflow-y: auto;
  padding: 16px;
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 0.875rem;
}

.log-entry {
  display: flex;
  gap: 12px;
  padding: 4px 0;
  border-bottom: 1px solid var(--n-border-color);
}

.log-entry:last-child {
  border-bottom: none;
}

.log-timestamp {
  color: var(--n-text-color-3);
  flex-shrink: 0;
  width: 80px;
}

.log-level {
  flex-shrink: 0;
  width: 60px;
  font-weight: 600;
}

.log-entry.info .log-level {
  color: var(--gproc-secondary);
}

.log-entry.warning .log-level {
  color: var(--gproc-warning);
}

.log-entry.error .log-level {
  color: var(--gproc-error);
}

.log-message {
  color: var(--n-text-color);
  flex: 1;
}

.system-commands {
  margin-bottom: 40px;
}

.commands-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.command-btn {
  height: 60px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

@media (max-width: 768px) {
  .daemon-status {
    padding: 16px;
  }
  
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }
  
  .status-overview {
    grid-template-columns: 1fr;
  }
  
  .metrics-cards {
    grid-template-columns: 1fr;
  }
  
  .config-grid {
    grid-template-columns: 1fr;
  }
  
  .commands-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>