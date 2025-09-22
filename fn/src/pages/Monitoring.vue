<template>
  <div class="monitoring">
    <n-grid :cols="1" :y-gap="16">
      <!-- Real-time Metrics -->
      <n-grid-item>
        <n-card title="System Metrics">
          <n-grid :cols="4" :x-gap="16">
            <n-grid-item>
              <n-statistic label="Running Processes" :value="currentMetrics?.processesRunning || 0" />
            </n-grid-item>
            <n-grid-item>
              <n-statistic label="CPU Usage" :value="`${currentMetrics?.cpuUsage || 0}%`" />
            </n-grid-item>
            <n-grid-item>
              <n-statistic label="Memory Usage" :value="`${currentMetrics?.memoryUsage || 0}%`" />
            </n-grid-item>
            <n-grid-item>
              <n-statistic label="Uptime" :value="currentMetrics?.uptime || '0s'" />
            </n-grid-item>
          </n-grid>
          
          <!-- Charts -->
          <div class="charts-container">
            <div class="chart">
              <h3>CPU Usage</h3>
              <Line :data="cpuChartData" :options="chartOptions" />
            </div>
            <div class="chart">
              <h3>Memory Usage</h3>
              <Line :data="memoryChartData" :options="chartOptions" />
            </div>
          </div>
        </n-card>
      </n-grid-item>
      
      <!-- Alerts -->
      <n-grid-item>
        <n-card title="Active Alerts">
          <template #header-extra>
            <n-button @click="fetchAlerts">
              <template #icon><RefreshCw /></template>
              Refresh
            </n-button>
          </template>
          
          <n-data-table
            :columns="alertColumns"
            :data="alerts"
            :pagination="{ pageSize: 10 }"
          />
        </n-card>
      </n-grid-item>
      
      <!-- Language Probes -->
      <n-grid-item>
        <n-card title="Language-Specific Probes">
          <n-space>
            <n-select
              v-model:value="selectedProcess"
              placeholder="Select Process"
              :options="processOptions"
              style="width: 200px"
            />
            <n-select
              v-model:value="selectedLanguage"
              placeholder="Select Language"
              :options="languageOptions"
              style="width: 150px"
            />
            <n-button type="primary" @click="runProbes" :loading="probesLoading">
              Run Probes
            </n-button>
          </n-space>
          
          <div v-if="probeResults.length" class="probe-results">
            <n-grid :cols="3" :x-gap="16" :y-gap="16">
              <n-grid-item v-for="result in probeResults" :key="result.name">
                <n-card :title="result.name">
                  <n-statistic
                    :value="`${result.value} ${result.unit}`"
                    :value-style="{ color: result.healthy ? '#18a058' : '#d03050' }"
                  />
                  <template #footer>
                    <n-tag :type="result.healthy ? 'success' : 'error'">
                      {{ result.healthy ? 'Healthy' : 'Unhealthy' }}
                    </n-tag>
                  </template>
                </n-card>
              </n-grid-item>
            </n-grid>
          </div>
        </n-card>
      </n-grid-item>
      
      <!-- Anomaly Detection -->
      <n-grid-item>
        <n-card title="ML-Based Anomaly Detection">
          <template #header-extra>
            <n-button @click="fetchAnomalies">
              <template #icon><Brain /></template>
              Detect Anomalies
            </n-button>
          </template>
          
          <div v-if="anomalies.length">
            <n-timeline>
              <n-timeline-item
                v-for="anomaly in anomalies"
                :key="anomaly.id"
                :type="anomaly.severity === 'critical' ? 'error' : 'warning'"
              >
                <template #header>{{ anomaly.title }}</template>
                <p>{{ anomaly.description }}</p>
                <n-tag :type="anomaly.severity === 'critical' ? 'error' : 'warning'">
                  {{ anomaly.severity }}
                </n-tag>
              </n-timeline-item>
            </n-timeline>
          </div>
          <n-empty v-else description="No anomalies detected" />
        </n-card>
      </n-grid-item>
    </n-grid>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { Line } from 'vue-chartjs'
import { Chart as ChartJS, CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend } from 'chart.js'
import { useMonitoringStore } from '../stores/monitoring'
import { useProcessStore } from '../stores/processes'
import { useMessage } from 'naive-ui'
import { RefreshCw, Brain, AlertTriangle } from 'lucide-vue-next'

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend)

const monitoringStore = useMonitoringStore()
const processStore = useProcessStore()
const message = useMessage()

const selectedProcess = ref('')
const selectedLanguage = ref('')
const probesLoading = ref(false)
const anomalies = ref([])

const { metrics, alerts, probeResults } = monitoringStore
const { processes } = processStore

const currentMetrics = computed(() => metrics.value[metrics.value.length - 1])

const processOptions = computed(() => 
  processes.value.map(p => ({ label: p.name, value: p.id }))
)

const languageOptions = [
  { label: 'Node.js', value: 'node' },
  { label: 'Python', value: 'python' },
  { label: 'Java', value: 'java' },
  { label: 'Go', value: 'go' },
  { label: 'Rust', value: 'rust' },
  { label: 'PHP', value: 'php' }
]

const cpuChartData = computed(() => ({
  labels: metrics.value.map(m => new Date(m.timestamp).toLocaleTimeString()),
  datasets: [{
    label: 'CPU Usage %',
    data: metrics.value.map(m => m.cpuUsage),
    borderColor: '#18a058',
    backgroundColor: 'rgba(24, 160, 88, 0.1)',
    tension: 0.4
  }]
}))

const memoryChartData = computed(() => ({
  labels: metrics.value.map(m => new Date(m.timestamp).toLocaleTimeString()),
  datasets: [{
    label: 'Memory Usage %',
    data: metrics.value.map(m => m.memoryUsage),
    borderColor: '#2080f0',
    backgroundColor: 'rgba(32, 128, 240, 0.1)',
    tension: 0.4
  }]
}))

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  scales: {
    y: {
      beginAtZero: true,
      max: 100
    }
  }
}

const alertColumns = [
  { title: 'Severity', key: 'severity', render: (row: any) => h('n-tag', {
    type: row.severity === 'critical' ? 'error' : row.severity === 'warning' ? 'warning' : 'info'
  }, { default: () => row.severity }) },
  { title: 'Name', key: 'name' },
  { title: 'Message', key: 'message' },
  { title: 'Resource', key: 'resource' },
  { title: 'Target', key: 'target' },
  { title: 'Time', key: 'timestamp', render: (row: any) => new Date(row.timestamp).toLocaleString() },
  {
    title: 'Actions',
    key: 'actions',
    render: (row: any) => h('n-button', {
      size: 'small',
      type: 'primary',
      disabled: row.acknowledged,
      onClick: () => acknowledgeAlert(row.id)
    }, { default: () => row.acknowledged ? 'Acknowledged' : 'Acknowledge' })
  }
]

const runProbes = async () => {
  if (!selectedProcess.value || !selectedLanguage.value) {
    message.warning('Please select both process and language')
    return
  }
  
  probesLoading.value = true
  const result = await monitoringStore.runProbes(selectedProcess.value, selectedLanguage.value)
  
  if (result.success) {
    message.success('Probes completed successfully')
  } else {
    message.error(result.error)
  }
  
  probesLoading.value = false
}

const acknowledgeAlert = async (alertId: string) => {
  const result = await monitoringStore.acknowledgeAlert(alertId)
  if (result.success) {
    message.success('Alert acknowledged')
  } else {
    message.error(result.error)
  }
}

const fetchAlerts = () => {
  monitoringStore.fetchAlerts()
}

const fetchAnomalies = async () => {
  const result = await monitoringStore.getAnomalies()
  if (result.success) {
    anomalies.value = result.anomalies
  } else {
    message.error(result.error)
  }
}

let metricsInterval: NodeJS.Timeout

onMounted(() => {
  processStore.fetchProcesses()
  monitoringStore.fetchAlerts()
  
  // Start real-time metrics collection
  metricsInterval = setInterval(() => {
    monitoringStore.fetchMetrics()
  }, 5000)
})

onUnmounted(() => {
  if (metricsInterval) {
    clearInterval(metricsInterval)
  }
})
</script>

<style scoped>
.monitoring {
  padding: 1rem;
}

.charts-container {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 2rem;
  margin-top: 2rem;
}

.chart {
  height: 300px;
}

.chart h3 {
  margin-bottom: 1rem;
  text-align: center;
}

.probe-results {
  margin-top: 2rem;
}
</style>