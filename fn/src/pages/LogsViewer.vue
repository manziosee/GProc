<template>
  <div class="logs-viewer">
    <div class="page-header">
      <h1>Logs Viewer</h1>
      <div class="header-actions">
        <n-select
          v-model:value="selectedProcess"
          :options="processOptions"
          placeholder="Select process"
          style="width: 200px"
          clearable
        />
        <n-button @click="refreshLogs">
          <RotateCcw :size="16" />
          Refresh
        </n-button>
        <n-button @click="clearLogs" type="warning">
          <Trash2 :size="16" />
          Clear
        </n-button>
        <n-button @click="downloadLogs">
          <Download :size="16" />
          Download
        </n-button>
      </div>
    </div>

    <!-- Log Controls -->
    <n-card :bordered="false" class="log-controls">
      <div class="controls-row">
        <div class="control-group">
          <label>Log Level:</label>
          <n-select
            v-model:value="logLevel"
            :options="logLevelOptions"
            style="width: 120px"
            size="small"
          />
        </div>
        
        <div class="control-group">
          <label>Lines:</label>
          <n-input-number
            v-model:value="maxLines"
            :min="10"
            :max="10000"
            :step="10"
            style="width: 100px"
            size="small"
          />
        </div>
        
        <div class="control-group">
          <label>Auto-scroll:</label>
          <n-switch v-model:value="autoScroll" size="small" />
        </div>
        
        <div class="control-group">
          <label>Live tail:</label>
          <n-switch v-model:value="liveTail" size="small" />
        </div>
        
        <div class="control-group">
          <n-input
            v-model:value="searchQuery"
            placeholder="Search logs..."
            clearable
            size="small"
            style="width: 200px"
          >
            <template #prefix>
              <Search :size="14" />
            </template>
          </n-input>
        </div>
      </div>
    </n-card>

    <!-- Log Terminal -->
    <n-card :bordered="false" class="log-terminal-card">
      <div class="terminal-header">
        <div class="terminal-title">
          <Terminal :size="16" />
          <span>{{ selectedProcess || 'All Processes' }} - Live Logs</span>
        </div>
        <div class="terminal-controls">
          <n-button-group size="small">
            <n-button @click="fontSize = Math.max(10, fontSize - 1)">
              <Minus :size="14" />
            </n-button>
            <n-button @click="fontSize = Math.min(20, fontSize + 1)">
              <Plus :size="14" />
            </n-button>
          </n-button-group>
          <n-button size="small" @click="toggleFullscreen">
            <Maximize :size="14" />
          </n-button>
        </div>
      </div>
      
      <div 
        ref="terminalContainer"
        class="terminal-container"
        :class="{ fullscreen: isFullscreen }"
        :style="{ fontSize: fontSize + 'px' }"
      >
        <div class="log-lines">
          <div
            v-for="(log, index) in filteredLogs"
            :key="index"
            class="log-line"
            :class="getLogLevelClass(log.level)"
          >
            <span class="log-timestamp">{{ formatTimestamp(log.timestamp) }}</span>
            <span class="log-level" :class="log.level">{{ log.level.toUpperCase() }}</span>
            <span class="log-process">{{ log.process }}</span>
            <span class="log-message" v-html="highlightSearch(log.message)"></span>
          </div>
        </div>
        
        <!-- Loading indicator -->
        <div v-if="loading" class="loading-indicator">
          <n-spin size="small" />
          <span>Loading logs...</span>
        </div>
        
        <!-- Empty state -->
        <div v-if="!loading && filteredLogs.length === 0" class="empty-state">
          <FileText :size="48" />
          <h3>No logs found</h3>
          <p>{{ selectedProcess ? 'This process has no logs yet' : 'Select a process to view logs' }}</p>
        </div>
      </div>
    </n-card>

    <!-- Log Statistics -->
    <div class="log-stats">
      <n-card title="Log Statistics" :bordered="false" size="small">
        <div class="stats-grid">
          <div class="stat-item">
            <div class="stat-value">{{ logStats.total }}</div>
            <div class="stat-label">Total Lines</div>
          </div>
          <div class="stat-item">
            <div class="stat-value error">{{ logStats.errors }}</div>
            <div class="stat-label">Errors</div>
          </div>
          <div class="stat-item">
            <div class="stat-value warning">{{ logStats.warnings }}</div>
            <div class="stat-label">Warnings</div>
          </div>
          <div class="stat-item">
            <div class="stat-value">{{ formatFileSize(logStats.size) }}</div>
            <div class="stat-label">Log Size</div>
          </div>
        </div>
      </n-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import {
  NCard,
  NButton,
  NSelect,
  NInputNumber,
  NSwitch,
  NInput,
  NButtonGroup,
  NSpin
} from 'naive-ui'
import {
  RotateCcw,
  Trash2,
  Download,
  Search,
  Terminal,
  Minus,
  Plus,
  Maximize,
  FileText
} from 'lucide-vue-next'

interface LogEntry {
  timestamp: Date
  level: 'info' | 'warn' | 'error' | 'debug'
  process: string
  message: string
}

const selectedProcess = ref<string | null>(null)
const logLevel = ref('all')
const maxLines = ref(1000)
const autoScroll = ref(true)
const liveTail = ref(true)
const searchQuery = ref('')
const fontSize = ref(12)
const isFullscreen = ref(false)
const loading = ref(false)
const terminalContainer = ref<HTMLElement>()

const processOptions = ref([
  { label: 'web-server', value: 'web-server' },
  { label: 'api-worker', value: 'api-worker' },
  { label: 'background-job', value: 'background-job' },
  { label: 'database-sync', value: 'database-sync' }
])

const logLevelOptions = [
  { label: 'All', value: 'all' },
  { label: 'Debug', value: 'debug' },
  { label: 'Info', value: 'info' },
  { label: 'Warning', value: 'warn' },
  { label: 'Error', value: 'error' }
]

const logs = ref<LogEntry[]>([
  {
    timestamp: new Date(),
    level: 'info',
    process: 'web-server',
    message: 'Server started on port 8080'
  },
  {
    timestamp: new Date(Date.now() - 1000),
    level: 'warn',
    process: 'web-server',
    message: 'High memory usage detected: 85%'
  },
  {
    timestamp: new Date(Date.now() - 2000),
    level: 'error',
    process: 'api-worker',
    message: 'Database connection failed: timeout after 30s'
  },
  {
    timestamp: new Date(Date.now() - 3000),
    level: 'info',
    process: 'background-job',
    message: 'Processing batch job #1234'
  },
  {
    timestamp: new Date(Date.now() - 4000),
    level: 'debug',
    process: 'database-sync',
    message: 'Syncing table users: 1000 records'
  }
])

const filteredLogs = computed(() => {
  let filtered = logs.value

  // Filter by process
  if (selectedProcess.value) {
    filtered = filtered.filter(log => log.process === selectedProcess.value)
  }

  // Filter by log level
  if (logLevel.value !== 'all') {
    filtered = filtered.filter(log => log.level === logLevel.value)
  }

  // Filter by search query
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(log =>
      log.message.toLowerCase().includes(query) ||
      log.process.toLowerCase().includes(query)
    )
  }

  // Limit lines
  return filtered.slice(-maxLines.value)
})

const logStats = computed(() => {
  const stats = {
    total: logs.value.length,
    errors: logs.value.filter(log => log.level === 'error').length,
    warnings: logs.value.filter(log => log.level === 'warn').length,
    size: logs.value.length * 100 // Approximate size
  }
  return stats
})

const formatTimestamp = (timestamp: Date) => {
  return timestamp.toLocaleTimeString('en-US', { 
    hour12: false,
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    fractionalSecondDigits: 3
  })
}

const getLogLevelClass = (level: string) => {
  return `log-level-${level}`
}

const highlightSearch = (message: string) => {
  if (!searchQuery.value) return message
  
  const regex = new RegExp(`(${searchQuery.value})`, 'gi')
  return message.replace(regex, '<mark>$1</mark>')
}

const formatFileSize = (bytes: number) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const refreshLogs = () => {
  loading.value = true
  setTimeout(() => {
    // Simulate API call
    loading.value = false
  }, 1000)
}

const clearLogs = () => {
  logs.value = []
}

const downloadLogs = () => {
  const logText = filteredLogs.value
    .map(log => `${formatTimestamp(log.timestamp)} [${log.level.toUpperCase()}] ${log.process}: ${log.message}`)
    .join('\n')
  
  const blob = new Blob([logText], { type: 'text/plain' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `gproc-logs-${new Date().toISOString().split('T')[0]}.txt`
  a.click()
  URL.revokeObjectURL(url)
}

const toggleFullscreen = () => {
  isFullscreen.value = !isFullscreen.value
}

// Auto-scroll to bottom when new logs arrive
watch(filteredLogs, async () => {
  if (autoScroll.value) {
    await nextTick()
    if (terminalContainer.value) {
      terminalContainer.value.scrollTop = terminalContainer.value.scrollHeight
    }
  }
})

// Simulate live tail
let liveInterval: NodeJS.Timeout | null = null

const startLiveTail = () => {
  if (liveInterval) return
  
  liveInterval = setInterval(() => {
    if (liveTail.value) {
      const messages = [
        'Processing request #' + Math.floor(Math.random() * 10000),
        'Memory usage: ' + Math.floor(Math.random() * 100) + '%',
        'Database query executed in ' + Math.floor(Math.random() * 1000) + 'ms',
        'User authentication successful',
        'Cache miss for key: user_' + Math.floor(Math.random() * 1000)
      ]
      
      const levels: Array<'info' | 'warn' | 'error' | 'debug'> = ['info', 'warn', 'error', 'debug']
      const processes = ['web-server', 'api-worker', 'background-job', 'database-sync']
      
      logs.value.push({
        timestamp: new Date(),
        level: levels[Math.floor(Math.random() * levels.length)],
        process: processes[Math.floor(Math.random() * processes.length)],
        message: messages[Math.floor(Math.random() * messages.length)]
      })
      
      // Keep only recent logs to prevent memory issues
      if (logs.value.length > 5000) {
        logs.value = logs.value.slice(-3000)
      }
    }
  }, 2000)
}

const stopLiveTail = () => {
  if (liveInterval) {
    clearInterval(liveInterval)
    liveInterval = null
  }
}

onMounted(() => {
  startLiveTail()
})

onUnmounted(() => {
  stopLiveTail()
})

watch(liveTail, (newValue) => {
  if (newValue) {
    startLiveTail()
  } else {
    stopLiveTail()
  }
})
</script>

<style scoped>
.logs-viewer {
  display: flex;
  flex-direction: column;
  gap: 20px;
  height: calc(100vh - 120px);
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-header h1 {
  margin: 0;
  color: var(--n-text-color);
  font-size: 28px;
  font-weight: 700;
}

.header-actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

.log-controls {
  flex-shrink: 0;
}

.controls-row {
  display: flex;
  gap: 24px;
  align-items: center;
  flex-wrap: wrap;
}

.control-group {
  display: flex;
  align-items: center;
  gap: 8px;
}

.control-group label {
  font-size: 14px;
  color: var(--n-text-color-2);
  white-space: nowrap;
}

.log-terminal-card {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.terminal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid var(--n-border-color);
  background: var(--n-color-embedded);
}

.terminal-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: var(--n-text-color);
}

.terminal-controls {
  display: flex;
  gap: 8px;
  align-items: center;
}

.terminal-container {
  flex: 1;
  background: #1a1a1a;
  color: #e5e5e5;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  overflow-y: auto;
  position: relative;
  min-height: 400px;
}

.terminal-container.fullscreen {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 9999;
  min-height: 100vh;
}

.log-lines {
  padding: 16px;
}

.log-line {
  display: flex;
  gap: 12px;
  padding: 2px 0;
  border-left: 3px solid transparent;
  padding-left: 8px;
  margin-left: -8px;
}

.log-line.log-level-error {
  border-left-color: #ef4444;
  background: rgba(239, 68, 68, 0.1);
}

.log-line.log-level-warn {
  border-left-color: #f59e0b;
  background: rgba(245, 158, 11, 0.1);
}

.log-line.log-level-info {
  border-left-color: #3b82f6;
}

.log-line.log-level-debug {
  border-left-color: #6b7280;
  opacity: 0.7;
}

.log-timestamp {
  color: #9ca3af;
  font-size: 0.85em;
  min-width: 100px;
  flex-shrink: 0;
}

.log-level {
  min-width: 60px;
  font-weight: 600;
  font-size: 0.8em;
  flex-shrink: 0;
}

.log-level.error {
  color: #ef4444;
}

.log-level.warn {
  color: #f59e0b;
}

.log-level.info {
  color: #3b82f6;
}

.log-level.debug {
  color: #6b7280;
}

.log-process {
  color: #10b981;
  min-width: 120px;
  font-weight: 500;
  flex-shrink: 0;
}

.log-message {
  flex: 1;
  word-break: break-word;
}

.log-message :deep(mark) {
  background: #fbbf24;
  color: #1f2937;
  padding: 1px 2px;
  border-radius: 2px;
}

.loading-indicator {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 20px;
  color: #9ca3af;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: #6b7280;
  text-align: center;
}

.empty-state h3 {
  margin: 16px 0 8px;
  color: var(--n-text-color-2);
}

.log-stats {
  flex-shrink: 0;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
}

.stat-item {
  text-align: center;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: var(--n-text-color);
  margin-bottom: 4px;
}

.stat-value.error {
  color: #ef4444;
}

.stat-value.warning {
  color: #f59e0b;
}

.stat-label {
  font-size: 12px;
  color: var(--n-text-color-3);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }
  
  .controls-row {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }
  
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>