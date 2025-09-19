<template>
  <div class="dashboard">
    <div class="dashboard-header">
      <div class="header-content">
        <h1 class="dashboard-title">Dashboard</h1>
        <p class="dashboard-subtitle">Real-time overview of your process ecosystem</p>
      </div>
      <div class="header-actions">
        <n-button type="primary" size="large">
          <template #icon><span>⚡</span></template>
          Quick Deploy
        </n-button>
      </div>
    </div>

    <div class="metrics-grid">
      <div class="metric-card running">
        <div class="metric-header">
          <div class="metric-icon">🟢</div>
          <div class="metric-trend up">+12%</div>
        </div>
        <div class="metric-value">{{ runningProcesses }}</div>
        <div class="metric-label">Running Processes</div>
        <div class="metric-footer">
          <span class="metric-detail">{{ totalProcesses }} total</span>
        </div>
      </div>
      
      <div class="metric-card performance">
        <div class="metric-header">
          <div class="metric-icon">📊</div>
          <div class="metric-trend down">-3%</div>
        </div>
        <div class="metric-value">{{ cpuUsage }}%</div>
        <div class="metric-label">CPU Usage</div>
        <div class="metric-footer">
          <span class="metric-detail">8 cores available</span>
        </div>
      </div>
      
      <div class="metric-card memory">
        <div class="metric-header">
          <div class="metric-icon">💾</div>
          <div class="metric-trend up">+5%</div>
        </div>
        <div class="metric-value">{{ memoryUsage }}</div>
        <div class="metric-label">Memory Usage</div>
        <div class="metric-footer">
          <span class="metric-detail">16GB total</span>
        </div>
      </div>
      
      <div class="metric-card uptime">
        <div class="metric-header">
          <div class="metric-icon">⏱️</div>
          <div class="metric-trend stable">99.9%</div>
        </div>
        <div class="metric-value">{{ uptime }}</div>
        <div class="metric-label">System Uptime</div>
        <div class="metric-footer">
          <span class="metric-detail">Last restart: 2d ago</span>
        </div>
      </div>
    </div>

    <div class="dashboard-content">
      <div class="content-section">
        <div class="section-header">
          <h2 class="section-title">Recent Activity</h2>
          <n-button text>View All</n-button>
        </div>
        <div class="activity-list">
          <div v-for="activity in recentActivity" :key="activity.id" class="activity-item">
            <div class="activity-icon" :class="activity.type">{{ activity.icon }}</div>
            <div class="activity-content">
              <div class="activity-title">{{ activity.title }}</div>
              <div class="activity-time">{{ activity.time }}</div>
            </div>
          </div>
        </div>
      </div>
      
      <div class="content-section">
        <div class="section-header">
          <h2 class="section-title">Process Health</h2>
          <n-button text>Manage</n-button>
        </div>
        <div class="health-grid">
          <div v-for="process in processHealth" :key="process.name" class="health-item">
            <div class="health-status" :class="process.status"></div>
            <div class="health-info">
              <div class="health-name">{{ process.name }}</div>
              <div class="health-details">{{ process.details }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { NButton } from 'naive-ui'

const runningProcesses = ref(8)
const totalProcesses = ref(12)
const cpuUsage = ref(34)
const memoryUsage = ref('2.4GB')
const uptime = ref('15d 4h')

const recentActivity = ref([
  { id: 1, type: 'success', icon: '✅', title: 'web-server started successfully', time: '2 minutes ago' },
  { id: 2, type: 'warning', icon: '⚠️', title: 'api-service restarted (high memory)', time: '5 minutes ago' },
  { id: 3, type: 'info', icon: '📄', title: 'Configuration updated', time: '10 minutes ago' },
  { id: 4, type: 'success', icon: '🚀', title: 'worker-pool scaled to 4 instances', time: '15 minutes ago' }
])

const processHealth = ref([
  { name: 'web-server', status: 'healthy', details: 'Response time: 45ms' },
  { name: 'api-service', status: 'healthy', details: 'Response time: 32ms' },
  { name: 'worker-pool', status: 'warning', details: 'High CPU usage' },
  { name: 'scheduler', status: 'healthy', details: 'Next run: 2h 15m' }
])
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
  grid-template-columns: 1fr 1fr;
  gap: 32px;
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

.activity-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.activity-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: var(--n-color);
  border-radius: 8px;
  border: 1px solid var(--n-border-color);
}

.activity-icon {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.875rem;
}

.activity-icon.success {
  background: rgba(16, 185, 129, 0.1);
}

.activity-icon.warning {
  background: rgba(245, 158, 11, 0.1);
}

.activity-icon.info {
  background: rgba(59, 130, 246, 0.1);
}

.activity-title {
  font-weight: 500;
  color: var(--n-text-color);
  margin-bottom: 2px;
}

.activity-time {
  font-size: 0.75rem;
  color: var(--n-text-color-3);
}

.health-grid {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.health-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: var(--n-color);
  border-radius: 8px;
  border: 1px solid var(--n-border-color);
}

.health-status {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  flex-shrink: 0;
}

.health-status.healthy {
  background: var(--gproc-success);
  box-shadow: 0 0 0 2px rgba(16, 185, 129, 0.2);
}

.health-status.warning {
  background: var(--gproc-warning);
  box-shadow: 0 0 0 2px rgba(245, 158, 11, 0.2);
}

.health-name {
  font-weight: 500;
  color: var(--n-text-color);
  margin-bottom: 2px;
}

.health-details {
  font-size: 0.75rem;
  color: var(--n-text-color-3);
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