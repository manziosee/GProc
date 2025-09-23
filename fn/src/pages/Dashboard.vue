<template>
  <div class="dashboard">
    <div class="dashboard-header">
      <div class="header-content">
        <h1 class="dashboard-title">Overview</h1>
        <p class="dashboard-subtitle">Fleet health and activity in real time</p>
      </div>
      <div class="header-actions">
        <n-button type="primary" size="large" @click="refreshAll">
          <template #icon><span>↻</span></template>
          Refresh
        </n-button>
      </div>
    </div>

    <div class="metrics-grid">
      <div class="metric-card running">
        <div class="metric-header">
          <div class="metric-icon">🟢</div>
          <div class="metric-trend up">live</div>
        </div>
        <div class="metric-value">{{ kpis.processes_running }}</div>
        <div class="metric-label">Running Processes</div>
        <div class="metric-footer">
          <span class="metric-detail">{{ kpis.processes_total }} total</span>
        </div>
      </div>
      
      <div class="metric-card performance">
        <div class="metric-header">
          <div class="metric-icon">📊</div>
          <div class="metric-trend" :class="kpis.cpu_usage > 70 ? 'down' : 'stable'">{{ kpis.cpu_usage > 70 ? 'high' : 'ok' }}</div>
        </div>
        <div class="metric-value">{{ kpis.cpu_usage.toFixed(1) }}%</div>
        <div class="metric-label">CPU Usage</div>
        <div class="metric-footer">
          <span class="metric-detail">cores auto-detected</span>
        </div>
      </div>
      
      <div class="metric-card memory">
        <div class="metric-header">
          <div class="metric-icon">💾</div>
          <div class="metric-trend" :class="kpis.memory_usage > 80 ? 'down' : 'up'">{{ kpis.memory_usage.toFixed(1) }}%</div>
        </div>
        <div class="metric-value">{{ kpis.memory_usage.toFixed(1) }}%</div>
        <div class="metric-label">Memory Usage</div>
        <div class="metric-footer">
          <span class="metric-detail">tracked via Prometheus</span>
        </div>
      </div>
      
      <div class="metric-card uptime">
        <div class="metric-header">
          <div class="metric-icon">⏱️</div>
          <div class="metric-trend stable">{{ kpis.uptime }}</div>
        </div>
        <div class="metric-value">SLA</div>
        <div class="metric-label">System Uptime</div>
        <div class="metric-footer">
          <span class="metric-detail">updated live</span>
        </div>
      </div>
    </div>

    <div class="dashboard-content">
      <div class="content-section wide">
        <div class="section-header">
          <h2 class="section-title">Processes</h2>
          <n-button text @click="refreshProcesses">Reload</n-button>
        </div>
        <ProcessTable
          :items="processes"
          :onStart="startProcess"
          :onStop="stopProcess"
          :onRestart="restartProcess"
          @refresh="refreshProcesses"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { NButton } from 'naive-ui'
import api from '../config/api'
import ProcessTable from '../components/ProcessTable.vue'
import { useWebSocketStore } from '../stores/websocket'

type Proc = {
  id: string
  name: string
  pid: number
  status: string
  restarts: number
  start_time?: string
}

const ws = useWebSocketStore()

const processes = ref<Proc[]>([])
const kpis = ref({ processes_running: 0, processes_total: 0, cpu_usage: 0, memory_usage: 0, uptime: '-' })

async function refreshMetrics() {
  try {
    const { data } = await api.get('/api/v1/metrics')
    kpis.value = {
      processes_running: data.processes_running ?? 0,
      processes_total: data.processes_total ?? 0,
      cpu_usage: data.cpu_usage ?? 0,
      memory_usage: data.memory_usage ?? 0,
      uptime: data.uptime ?? '-'
    }
  } catch (error) {
    console.error('Failed to fetch metrics:', error)
  }
}

async function refreshProcesses() {
  try {
    const { data } = await api.get('/api/v1/processes')
    processes.value = data || []
  } catch (error) {
    console.error('Failed to fetch processes:', error)
    processes.value = []
  }
}

function refreshAll() {
  refreshMetrics()
  refreshProcesses()
}

async function startProcess(id: string) {
  try {
    await api.post(`/api/v1/processes/${id}/start`)
    refreshProcesses()
  } catch (error) {
    console.error('Failed to start process:', error)
  }
}
async function stopProcess(id: string) {
  try {
    await api.post(`/api/v1/processes/${id}/stop`)
    refreshProcesses()
  } catch (error) {
    console.error('Failed to stop process:', error)
  }
}
async function restartProcess(id: string) {
  try {
    await api.post(`/api/v1/processes/${id}/restart`)
    refreshProcesses()
  } catch (error) {
    console.error('Failed to restart process:', error)
  }
}

onMounted(() => {
  refreshAll()
  if (!ws.connected) ws.connect()
})
</script>

<style scoped>
.dashboard {
  padding: 32px;
  max-width: 1400px;
  margin: 0 auto;
}

.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  margin-bottom: 32px;
  padding-bottom: 24px;
  border-bottom: 1px solid var(--n-border-color);
}

.dashboard-title {
  font-size: 2.5rem;
  font-weight: 700;
  color: var(--n-text-color);
  margin: 0 0 8px 0;
  letter-spacing: -0.025em;
}

.dashboard-subtitle {
  font-size: 1.125rem;
  color: var(--n-text-color-2);
  margin: 0;
  font-weight: 400;
}

.metrics-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 24px;
  margin-bottom: 40px;
}

.metric-card {
  background: var(--n-color-embedded);
  border: 1px solid var(--n-border-color);
  border-radius: 16px;
  padding: 24px;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
}

.metric-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, var(--accent-color), var(--accent-light));
}

.metric-card.running {
  --accent-color: var(--gproc-success);
  --accent-light: #34d399;
}

.metric-card.performance {
  --accent-color: var(--gproc-secondary);
  --accent-light: #60a5fa;
}

.metric-card.memory {
  --accent-color: var(--gproc-accent);
  --accent-light: #a78bfa;
}

.metric-card.uptime {
  --accent-color: var(--gproc-warning);
  --accent-light: #fbbf24;
}

.metric-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-xl);
}

.metric-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.metric-icon {
  font-size: 1.5rem;
}

.metric-trend {
  font-size: 0.875rem;
  font-weight: 600;
  padding: 4px 8px;
  border-radius: 12px;
}

.metric-trend.up {
  background: rgba(16, 185, 129, 0.1);
  color: var(--gproc-success);
}

.metric-trend.down {
  background: rgba(239, 68, 68, 0.1);
  color: var(--gproc-error);
}

.metric-trend.stable {
  background: rgba(59, 130, 246, 0.1);
  color: var(--gproc-secondary);
}

.metric-value {
  font-size: 2.5rem;
  font-weight: 700;
  color: var(--n-text-color);
  line-height: 1;
  margin-bottom: 8px;
}

.metric-label {
  font-size: 0.875rem;
  color: var(--n-text-color-2);
  font-weight: 500;
  margin-bottom: 12px;
}

.metric-footer {
  padding-top: 12px;
  border-top: 1px solid var(--n-border-color);
}

.metric-detail {
  font-size: 0.75rem;
  color: var(--n-text-color-3);
}

.dashboard-content {
  display: grid;
  grid-template-columns: 1fr;
  gap: 32px;
}

.content-section.wide {
  grid-column: 1 / -1;
}

.content-section {
  background: var(--n-color-embedded);
  border: 1px solid var(--n-border-color);
  border-radius: 16px;
  padding: 24px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.section-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--n-text-color);
  margin: 0;
}

@media (max-width: 768px) {
  .dashboard {
    padding: 16px;
  }
  
  .dashboard-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }
  
  .metrics-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 16px;
  }
  
  .dashboard-content {
    grid-template-columns: 1fr;
    gap: 24px;
  }
}
</style>