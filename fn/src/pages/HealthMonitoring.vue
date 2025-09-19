<template>
  <div class="health-monitoring">
    <div class="page-header">
      <h1>Health Monitoring</h1>
      <div class="header-actions">
        <n-button @click="refreshData">
          <RotateCcw :size="16" />
          Refresh
        </n-button>
        <n-button @click="exportReport">
          <Download :size="16" />
          Export Report
        </n-button>
      </div>
    </div>

    <!-- System Overview -->
    <div class="system-overview">
      <n-card title="System Health" :bordered="false">
        <div class="health-grid">
          <div class="health-item" :class="systemHealth.overall.status">
            <div class="health-icon">
              <component :is="getHealthIcon(systemHealth.overall.status)" :size="24" />
            </div>
            <div class="health-content">
              <div class="health-title">Overall Health</div>
              <div class="health-status">{{ systemHealth.overall.text }}</div>
              <div class="health-score">{{ systemHealth.overall.score }}/100</div>
            </div>
          </div>
          
          <div class="health-item" :class="systemHealth.processes.status">
            <div class="health-icon">
              <Activity :size="24" />
            </div>
            <div class="health-content">
              <div class="health-title">Processes</div>
              <div class="health-status">{{ systemHealth.processes.running }}/{{ systemHealth.processes.total }} Running</div>
              <div class="health-details">{{ systemHealth.processes.failed }} Failed</div>
            </div>
          </div>
          
          <div class="health-item" :class="systemHealth.resources.status">
            <div class="health-icon">
              <Cpu :size="24" />
            </div>
            <div class="health-content">
              <div class="health-title">Resources</div>
              <div class="health-status">CPU: {{ systemHealth.resources.cpu }}%</div>
              <div class="health-details">Memory: {{ systemHealth.resources.memory }}%</div>
            </div>
          </div>
          
          <div class="health-item" :class="systemHealth.network.status">
            <div class="health-icon">
              <Wifi :size="24" />
            </div>
            <div class="health-content">
              <div class="health-title">Network</div>
              <div class="health-status">{{ systemHealth.network.status }}</div>
              <div class="health-details">{{ systemHealth.network.latency }}ms avg</div>
            </div>
          </div>
        </div>
      </n-card>
    </div>

    <!-- Real-time Metrics -->
    <div class="metrics-section">
      <n-grid :cols="2" :x-gap="20" :y-gap="20">
        <n-grid-item>
          <n-card title="CPU & Memory Usage" :bordered="false">
            <div class="chart-container">
              <Line :data="resourceChartData" :options="chartOptions" />
            </div>
          </n-card>
        </n-grid-item>
        
        <n-grid-item>
          <n-card title="Network Traffic" :bordered="false">
            <div class="chart-container">
              <Line :data="networkChartData" :options="chartOptions" />
            </div>
          </n-card>
        </n-grid-item>
      </n-grid>
    </div>

    <!-- Process Health Status -->
    <n-card title="Process Health Status" :bordered="false">
      <n-data-table
        :columns="processColumns"
        :data="processHealthData"
        :pagination="false"
        :row-class-name="getProcessRowClass"
      />
    </n-card>

    <!-- Health Checks Configuration -->
    <n-card title="Health Check Configuration" :bordered="false">
      <div class="health-checks-grid">
        <div
          v-for="check in healthChecks"
          :key="check.id"
          class="health-check-card"
          :class="check.status"
        >
          <div class="check-header">
            <div class="check-title">
              <component :is="getHealthIcon(check.status)" :size="16" />
              {{ check.name }}
            </div>
            <div class="check-actions">
              <n-button size="tiny" quaternary @click="editHealthCheck(check)">
                <Edit :size="14" />
              </n-button>
              <n-button size="tiny" quaternary @click="runHealthCheck(check)">
                <Play :size="14" />
              </n-button>
            </div>
          </div>
          
          <div class="check-details">
            <div class="check-url">{{ check.url }}</div>
            <div class="check-info">
              <span>Interval: {{ check.interval }}s</span>
              <span>Timeout: {{ check.timeout }}s</span>
              <span>Last check: {{ formatTime(check.lastCheck) }}</span>
            </div>
            <div class="check-metrics">
              <div class="metric">
                <span class="metric-label">Response Time:</span>
                <span class="metric-value">{{ check.responseTime }}ms</span>
              </div>
              <div class="metric">
                <span class="metric-label">Success Rate:</span>
                <span class="metric-value">{{ check.successRate }}%</span>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <div class="add-health-check">
        <n-button dashed @click="showAddHealthCheck = true" style="width: 100%">
          <Plus :size="16" />
          Add Health Check
        </n-button>
      </div>
    </n-card>

    <!-- Alerts & Incidents -->
    <n-card title="Recent Alerts & Incidents" :bordered="false">
      <div class="alerts-list">
        <div
          v-for="alert in recentAlerts"
          :key="alert.id"
          class="alert-item"
          :class="alert.severity"
        >
          <div class="alert-icon">
            <component :is="getAlertIcon(alert.severity)" :size="16" />
          </div>
          <div class="alert-content">
            <div class="alert-title">{{ alert.title }}</div>
            <div class="alert-message">{{ alert.message }}</div>
            <div class="alert-time">{{ formatTime(alert.timestamp) }}</div>
          </div>
          <div class="alert-actions">
            <n-button size="small" @click="acknowledgeAlert(alert)">
              Acknowledge
            </n-button>
          </div>
        </div>
      </div>
    </n-card>

    <!-- Add Health Check Modal -->
    <n-modal v-model:show="showAddHealthCheck" preset="dialog" title="Add Health Check">
      <template #default>
        <div class="add-health-check-form">
          <n-form :model="newHealthCheck" label-placement="top">
            <n-form-item label="Name">
              <n-input v-model:value="newHealthCheck.name" placeholder="e.g., API Health Check" />
            </n-form-item>
            
            <n-form-item label="URL">
              <n-input v-model:value="newHealthCheck.url" placeholder="http://localhost:8080/health" />
            </n-form-item>
            
            <n-grid :cols="2" :x-gap="16">
              <n-form-item-gi label="Interval (seconds)">
                <n-input-number v-model:value="newHealthCheck.interval" :min="5" :max="3600" />
              </n-form-item-gi>
              
              <n-form-item-gi label="Timeout (seconds)">
                <n-input-number v-model:value="newHealthCheck.timeout" :min="1" :max="60" />
              </n-form-item-gi>
            </n-grid>
            
            <n-form-item label="Expected Status Code">
              <n-input-number v-model:value="newHealthCheck.expectedStatus" :min="100" :max="599" />
            </n-form-item>
          </n-form>
        </div>
      </template>
      
      <template #action>
        <div style="display: flex; gap: 12px;">
          <n-button @click="showAddHealthCheck = false">Cancel</n-button>
          <n-button type="primary" @click="addHealthCheck">Add Health Check</n-button>
        </div>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted, h } from 'vue'
import {
  NCard,
  NButton,
  NDataTable,
  NGrid,
  NGridItem,
  NModal,
  NForm,
  NFormItem,
  NFormItemGi,
  NInput,
  NInputNumber
} from 'naive-ui'
import { Line } from 'vue-chartjs'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
} from 'chart.js'
import {
  RotateCcw,
  Download,
  Activity,
  Cpu,
  Wifi,
  CheckCircle,
  AlertTriangle,
  XCircle,
  Edit,
  Play,
  Plus
} from 'lucide-vue-next'

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend
)

const showAddHealthCheck = ref(false)

const systemHealth = reactive({
  overall: {
    status: 'healthy',
    text: 'All Systems Operational',
    score: 95
  },
  processes: {
    status: 'healthy',
    running: 18,
    total: 20,
    failed: 2
  },
  resources: {
    status: 'warning',
    cpu: 67,
    memory: 43
  },
  network: {
    status: 'healthy',
    latency: 12
  }
})

const resourceChartData = reactive({
  labels: Array.from({ length: 20 }, (_, i) => {
    const time = new Date(Date.now() - (19 - i) * 30000)
    return time.toLocaleTimeString('en-US', { hour12: false, timeStyle: 'short' })
  }),
  datasets: [
    {
      label: 'CPU Usage (%)',
      data: Array.from({ length: 20 }, () => Math.floor(Math.random() * 30) + 40),
      borderColor: '#3b82f6',
      backgroundColor: 'rgba(59, 130, 246, 0.1)',
      tension: 0.4
    },
    {
      label: 'Memory Usage (%)',
      data: Array.from({ length: 20 }, () => Math.floor(Math.random() * 20) + 35),
      borderColor: '#10b981',
      backgroundColor: 'rgba(16, 185, 129, 0.1)',
      tension: 0.4
    }
  ]
})

const networkChartData = reactive({
  labels: Array.from({ length: 20 }, (_, i) => {
    const time = new Date(Date.now() - (19 - i) * 30000)
    return time.toLocaleTimeString('en-US', { hour12: false, timeStyle: 'short' })
  }),
  datasets: [
    {
      label: 'Inbound (MB/s)',
      data: Array.from({ length: 20 }, () => Math.floor(Math.random() * 50) + 10),
      borderColor: '#8b5cf6',
      backgroundColor: 'rgba(139, 92, 246, 0.1)',
      tension: 0.4
    },
    {
      label: 'Outbound (MB/s)',
      data: Array.from({ length: 20 }, () => Math.floor(Math.random() * 30) + 5),
      borderColor: '#f59e0b',
      backgroundColor: 'rgba(245, 158, 11, 0.1)',
      tension: 0.4
    }
  ]
})

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      position: 'top' as const,
    }
  },
  scales: {
    y: {
      beginAtZero: true
    }
  }
}

const processHealthData = ref([
  {
    name: 'web-server',
    status: 'healthy',
    uptime: '2d 14h 32m',
    cpu: 15.3,
    memory: 128.5,
    healthScore: 98,
    lastHealthCheck: new Date(Date.now() - 30000)
  },
  {
    name: 'api-worker',
    status: 'warning',
    uptime: '1d 8h 15m',
    cpu: 45.7,
    memory: 256.8,
    healthScore: 75,
    lastHealthCheck: new Date(Date.now() - 60000)
  },
  {
    name: 'background-job',
    status: 'critical',
    uptime: '0d 0h 5m',
    cpu: 0,
    memory: 0,
    healthScore: 25,
    lastHealthCheck: new Date(Date.now() - 300000)
  }
])

const processColumns = [
  {
    title: 'Process',
    key: 'name',
    render: (row: any) => h('div', { class: 'process-name-cell' }, [
      h('span', { class: 'process-name' }, row.name),
      h('div', { class: 'process-uptime' }, `Uptime: ${row.uptime}`)
    ])
  },
  {
    title: 'Health Score',
    key: 'healthScore',
    render: (row: any) => h('div', { class: 'health-score-cell' }, [
      h('div', { class: 'score-bar' }, [
        h('div', { 
          class: 'score-fill', 
          style: { 
            width: `${row.healthScore}%`,
            backgroundColor: row.healthScore > 80 ? '#10b981' : row.healthScore > 50 ? '#f59e0b' : '#ef4444'
          }
        })
      ]),
      h('span', { class: 'score-text' }, `${row.healthScore}%`)
    ])
  },
  {
    title: 'Resources',
    key: 'resources',
    render: (row: any) => h('div', { class: 'resources-cell' }, [
      h('div', `CPU: ${row.cpu}%`),
      h('div', `Memory: ${row.memory} MB`)
    ])
  },
  {
    title: 'Last Check',
    key: 'lastHealthCheck',
    render: (row: any) => formatTime(row.lastHealthCheck)
  }
]

const healthChecks = ref([
  {
    id: 1,
    name: 'Web Server Health',
    url: 'http://localhost:8080/health',
    interval: 30,
    timeout: 5,
    status: 'healthy',
    responseTime: 45,
    successRate: 99.8,
    lastCheck: new Date(Date.now() - 30000)
  },
  {
    id: 2,
    name: 'API Endpoint',
    url: 'http://localhost:3000/api/health',
    interval: 60,
    timeout: 10,
    status: 'warning',
    responseTime: 1200,
    successRate: 95.2,
    lastCheck: new Date(Date.now() - 60000)
  },
  {
    id: 3,
    name: 'Database Connection',
    url: 'http://localhost:5432/health',
    interval: 120,
    timeout: 15,
    status: 'critical',
    responseTime: 0,
    successRate: 0,
    lastCheck: new Date(Date.now() - 300000)
  }
])

const recentAlerts = ref([
  {
    id: 1,
    severity: 'critical',
    title: 'Database Connection Failed',
    message: 'Unable to connect to primary database server',
    timestamp: new Date(Date.now() - 120000)
  },
  {
    id: 2,
    severity: 'warning',
    title: 'High Response Time',
    message: 'API endpoint response time exceeded 1000ms threshold',
    timestamp: new Date(Date.now() - 300000)
  },
  {
    id: 3,
    severity: 'info',
    title: 'Health Check Recovered',
    message: 'Web server health check is now responding normally',
    timestamp: new Date(Date.now() - 600000)
  }
])

const newHealthCheck = reactive({
  name: '',
  url: '',
  interval: 30,
  timeout: 5,
  expectedStatus: 200
})

const getHealthIcon = (status: string) => {
  const icons = {
    healthy: CheckCircle,
    warning: AlertTriangle,
    critical: XCircle
  }
  return icons[status as keyof typeof icons] || CheckCircle
}

const getAlertIcon = (severity: string) => {
  const icons = {
    critical: XCircle,
    warning: AlertTriangle,
    info: CheckCircle
  }
  return icons[severity as keyof typeof icons] || CheckCircle
}

const getProcessRowClass = (row: any) => {
  return `process-row-${row.status}`
}

const formatTime = (time: Date) => {
  const now = new Date()
  const diff = Math.floor((now.getTime() - time.getTime()) / 60000)
  if (diff < 1) return 'Just now'
  if (diff < 60) return `${diff}m ago`
  const hours = Math.floor(diff / 60)
  return `${hours}h ago`
}

const refreshData = () => {
  // Simulate data refresh
  console.log('Refreshing health monitoring data...')
}

const exportReport = () => {
  console.log('Exporting health report...')
}

const editHealthCheck = (check: any) => {
  console.log('Edit health check:', check)
}

const runHealthCheck = (check: any) => {
  console.log('Run health check:', check)
}

const addHealthCheck = () => {
  const newCheck = {
    id: Date.now(),
    ...newHealthCheck,
    status: 'healthy',
    responseTime: 0,
    successRate: 100,
    lastCheck: new Date()
  }
  
  healthChecks.value.push(newCheck)
  showAddHealthCheck.value = false
  
  // Reset form
  Object.assign(newHealthCheck, {
    name: '',
    url: '',
    interval: 30,
    timeout: 5,
    expectedStatus: 200
  })
}

const acknowledgeAlert = (alert: any) => {
  console.log('Acknowledge alert:', alert)
}

// Real-time updates
let updateInterval: NodeJS.Timeout | null = null

onMounted(() => {
  updateInterval = setInterval(() => {
    // Update resource charts
    const newTime = new Date().toLocaleTimeString('en-US', { hour12: false, timeStyle: 'short' })
    
    resourceChartData.labels.shift()
    resourceChartData.labels.push(newTime)
    
    resourceChartData.datasets[0].data.shift()
    resourceChartData.datasets[0].data.push(Math.floor(Math.random() * 30) + 40)
    
    resourceChartData.datasets[1].data.shift()
    resourceChartData.datasets[1].data.push(Math.floor(Math.random() * 20) + 35)
    
    networkChartData.labels.shift()
    networkChartData.labels.push(newTime)
    
    networkChartData.datasets[0].data.shift()
    networkChartData.datasets[0].data.push(Math.floor(Math.random() * 50) + 10)
    
    networkChartData.datasets[1].data.shift()
    networkChartData.datasets[1].data.push(Math.floor(Math.random() * 30) + 5)
  }, 5000)
})

onUnmounted(() => {
  if (updateInterval) {
    clearInterval(updateInterval)
  }
})
</script>

<style scoped>
.health-monitoring {
  display: flex;
  flex-direction: column;
  gap: 24px;
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
}

.health-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
}

.health-item {
  display: flex;
  gap: 16px;
  padding: 20px;
  border-radius: 8px;
  border: 1px solid var(--n-border-color);
}

.health-item.healthy {
  background: rgba(16, 185, 129, 0.1);
  border-color: #10b981;
}

.health-item.warning {
  background: rgba(245, 158, 11, 0.1);
  border-color: #f59e0b;
}

.health-item.critical {
  background: rgba(239, 68, 68, 0.1);
  border-color: #ef4444;
}

.health-icon {
  flex-shrink: 0;
  color: inherit;
}

.health-content {
  flex: 1;
}

.health-title {
  font-weight: 600;
  color: var(--n-text-color);
  margin-bottom: 4px;
}

.health-status {
  color: var(--n-text-color-2);
  margin-bottom: 2px;
}

.health-score {
  font-size: 18px;
  font-weight: 700;
  color: var(--n-text-color);
}

.health-details {
  font-size: 12px;
  color: var(--n-text-color-3);
}

.chart-container {
  height: 300px;
}

.process-name-cell {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.process-name {
  font-weight: 600;
  color: var(--n-text-color);
}

.process-uptime {
  font-size: 12px;
  color: var(--n-text-color-3);
}

.health-score-cell {
  display: flex;
  align-items: center;
  gap: 12px;
}

.score-bar {
  width: 60px;
  height: 8px;
  background: var(--n-color-hover);
  border-radius: 4px;
  overflow: hidden;
}

.score-fill {
  height: 100%;
  transition: width 0.3s ease;
}

.score-text {
  font-weight: 600;
  color: var(--n-text-color);
}

.resources-cell {
  display: flex;
  flex-direction: column;
  gap: 2px;
  font-size: 14px;
  color: var(--n-text-color-2);
}

.health-checks-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
  gap: 20px;
  margin-bottom: 20px;
}

.health-check-card {
  border: 1px solid var(--n-border-color);
  border-radius: 8px;
  padding: 16px;
}

.health-check-card.healthy {
  border-left: 4px solid #10b981;
}

.health-check-card.warning {
  border-left: 4px solid #f59e0b;
}

.health-check-card.critical {
  border-left: 4px solid #ef4444;
}

.check-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.check-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: var(--n-text-color);
}

.check-actions {
  display: flex;
  gap: 4px;
}

.check-details {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.check-url {
  font-family: monospace;
  font-size: 12px;
  color: var(--n-text-color-2);
  background: var(--n-color-hover);
  padding: 4px 8px;
  border-radius: 4px;
}

.check-info {
  display: flex;
  gap: 16px;
  font-size: 12px;
  color: var(--n-text-color-3);
}

.check-metrics {
  display: flex;
  gap: 16px;
}

.metric {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.metric-label {
  font-size: 11px;
  color: var(--n-text-color-3);
  text-transform: uppercase;
}

.metric-value {
  font-weight: 600;
  color: var(--n-text-color);
}

.alerts-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.alert-item {
  display: flex;
  gap: 12px;
  padding: 16px;
  border-radius: 8px;
  border-left: 4px solid transparent;
}

.alert-item.critical {
  background: rgba(239, 68, 68, 0.1);
  border-left-color: #ef4444;
}

.alert-item.warning {
  background: rgba(245, 158, 11, 0.1);
  border-left-color: #f59e0b;
}

.alert-item.info {
  background: rgba(59, 130, 246, 0.1);
  border-left-color: #3b82f6;
}

.alert-icon {
  flex-shrink: 0;
  margin-top: 2px;
}

.alert-content {
  flex: 1;
}

.alert-title {
  font-weight: 600;
  color: var(--n-text-color);
  margin-bottom: 4px;
}

.alert-message {
  color: var(--n-text-color-2);
  margin-bottom: 8px;
}

.alert-time {
  font-size: 12px;
  color: var(--n-text-color-3);
}

.alert-actions {
  flex-shrink: 0;
}

:deep(.process-row-healthy td) {
  border-left: 3px solid #10b981;
}

:deep(.process-row-warning td) {
  border-left: 3px solid #f59e0b;
}

:deep(.process-row-critical td) {
  border-left: 3px solid #ef4444;
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }
  
  .health-checks-grid {
    grid-template-columns: 1fr;
  }
}
</style>