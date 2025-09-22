<template>
  <div class="language-probes">
    <n-card title="Language-Specific Monitoring Probes">
      <n-grid :cols="1" :y-gap="16">
        <!-- Probe Selection -->
        <n-grid-item>
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
            <n-button type="primary" @click="runProbes" :loading="loading">
              <template #icon><Zap /></template>
              Run Probes
            </n-button>
          </n-space>
        </n-grid-item>
        
        <!-- Language-Specific Probe Results -->
        <n-grid-item v-if="probeResults.length">
          <n-card :title="`${selectedLanguage} Probe Results`">
            <n-grid :cols="3" :x-gap="16" :y-gap="16">
              <n-grid-item v-for="result in probeResults" :key="result.name">
                <n-card>
                  <template #header>
                    <n-space align="center">
                      <component :is="getProbeIcon(result.name)" :size="20" />
                      {{ result.name }}
                    </n-space>
                  </template>
                  
                  <n-statistic
                    :value="`${result.value} ${result.unit}`"
                    :value-style="{ color: result.healthy ? '#18a058' : '#d03050' }"
                  />
                  
                  <template #footer>
                    <n-space justify="space-between">
                      <n-tag :type="result.healthy ? 'success' : 'error'">
                        {{ result.healthy ? 'Healthy' : 'Unhealthy' }}
                      </n-tag>
                      <span class="timestamp">{{ formatTime(result.timestamp) }}</span>
                    </n-space>
                  </template>
                </n-card>
              </n-grid-item>
            </n-grid>
          </n-card>
        </n-grid-item>
        
        <!-- Language-Specific Information -->
        <n-grid-item>
          <n-tabs type="line">
            <n-tab-pane name="nodejs" tab="Node.js">
              <n-card title="Node.js Monitoring">
                <n-descriptions :column="2" bordered>
                  <n-descriptions-item label="Event Loop Lag">
                    Measures the delay between scheduled and actual execution of callbacks
                  </n-descriptions-item>
                  <n-descriptions-item label="Heap Usage">
                    V8 heap memory utilization and garbage collection metrics
                  </n-descriptions-item>
                  <n-descriptions-item label="Active Handles">
                    Number of active handles (timers, sockets, etc.)
                  </n-descriptions-item>
                  <n-descriptions-item label="Active Requests">
                    Number of active asynchronous requests
                  </n-descriptions-item>
                </n-descriptions>
              </n-card>
            </n-tab-pane>
            
            <n-tab-pane name="python" tab="Python">
              <n-card title="Python Monitoring">
                <n-descriptions :column="2" bordered>
                  <n-descriptions-item label="GIL Wait Time">
                    Time spent waiting for Global Interpreter Lock
                  </n-descriptions-item>
                  <n-descriptions-item label="Memory Leaks">
                    Detection of unreferenced objects and memory growth
                  </n-descriptions-item>
                  <n-descriptions-item label="Thread Count">
                    Number of active Python threads
                  </n-descriptions-item>
                  <n-descriptions-item label="Import Time">
                    Module import performance metrics
                  </n-descriptions-item>
                </n-descriptions>
              </n-card>
            </n-tab-pane>
            
            <n-tab-pane name="java" tab="Java">
              <n-card title="Java Monitoring">
                <n-descriptions :column="2" bordered>
                  <n-descriptions-item label="Heap Memory">
                    JVM heap usage and garbage collection statistics
                  </n-descriptions-item>
                  <n-descriptions-item label="Thread Pool">
                    Thread pool utilization and queue sizes
                  </n-descriptions-item>
                  <n-descriptions-item label="Class Loading">
                    Number of loaded classes and class loader metrics
                  </n-descriptions-item>
                  <n-descriptions-item label="JIT Compilation">
                    Just-in-time compilation statistics
                  </n-descriptions-item>
                </n-descriptions>
              </n-card>
            </n-tab-pane>
            
            <n-tab-pane name="go" tab="Go">
              <n-card title="Go Monitoring">
                <n-descriptions :column="2" bordered>
                  <n-descriptions-item label="Goroutines">
                    Number of active goroutines and scheduler metrics
                  </n-descriptions-item>
                  <n-descriptions-item label="Memory Stats">
                    Go runtime memory allocation and GC statistics
                  </n-descriptions-item>
                  <n-descriptions-item label="Channel Stats">
                    Channel buffer utilization and blocking operations
                  </n-descriptions-item>
                  <n-descriptions-item label="pprof Integration">
                    CPU and memory profiling data collection
                  </n-descriptions-item>
                </n-descriptions>
              </n-card>
            </n-tab-pane>
            
            <n-tab-pane name="rust" tab="Rust">
              <n-card title="Rust Monitoring">
                <n-descriptions :column="2" bordered>
                  <n-descriptions-item label="Thread Count">
                    Number of OS threads and async task metrics
                  </n-descriptions-item>
                  <n-descriptions-item label="Allocator Stats">
                    Memory allocator performance and fragmentation
                  </n-descriptions-item>
                  <n-descriptions-item label="Panic Recovery">
                    Panic frequency and recovery statistics
                  </n-descriptions-item>
                  <n-descriptions-item label="Async Runtime">
                    Tokio/async-std runtime performance metrics
                  </n-descriptions-item>
                </n-descriptions>
              </n-card>
            </n-tab-pane>
            
            <n-tab-pane name="php" tab="PHP">
              <n-card title="PHP Monitoring">
                <n-descriptions :column="2" bordered>
                  <n-descriptions-item label="OPcache Stats">
                    Opcode cache hit ratio and memory usage
                  </n-descriptions-item>
                  <n-descriptions-item label="FPM Status">
                    PHP-FPM pool status and worker utilization
                  </n-descriptions-item>
                  <n-descriptions-item label="Memory Usage">
                    PHP memory limit utilization and peak usage
                  </n-descriptions-item>
                  <n-descriptions-item label="Extension Stats">
                    Loaded extensions and their performance impact
                  </n-descriptions-item>
                </n-descriptions>
              </n-card>
            </n-tab-pane>
          </n-tabs>
        </n-grid-item>
        
        <!-- Historical Probe Data -->
        <n-grid-item v-if="historicalData.length">
          <n-card title="Historical Probe Data">
            <div class="chart-container">
              <Line :data="chartData" :options="chartOptions" />
            </div>
          </n-card>
        </n-grid-item>
      </n-grid>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Line } from 'vue-chartjs'
import { useMonitoringStore } from '../stores/monitoring'
import { useProcessStore } from '../stores/processes'
import { useMessage } from 'naive-ui'
import { Zap, Cpu, Clock, Activity, Database, Code } from 'lucide-vue-next'

const monitoringStore = useMonitoringStore()
const processStore = useProcessStore()
const message = useMessage()

const selectedProcess = ref('')
const selectedLanguage = ref('')
const loading = ref(false)
const historicalData = ref([])

const { probeResults } = monitoringStore
const { processes } = processStore

const processOptions = computed(() =>
  processes.value.map(p => ({ label: p.name, value: p.id }))
)

const languageOptions = [
  { label: 'Node.js', value: 'nodejs' },
  { label: 'Python', value: 'python' },
  { label: 'Java', value: 'java' },
  { label: 'Go', value: 'go' },
  { label: 'Rust', value: 'rust' },
  { label: 'PHP', value: 'php' }
]

const chartData = computed(() => ({
  labels: historicalData.value.map((_, index) => `Run ${index + 1}`),
  datasets: probeResults.value.map(result => ({
    label: result.name,
    data: historicalData.value.map(data => data[result.name] || 0),
    borderColor: getRandomColor(),
    backgroundColor: 'transparent',
    tension: 0.4
  }))
}))

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    title: {
      display: true,
      text: 'Probe Results Over Time'
    }
  },
  scales: {
    y: {
      beginAtZero: true
    }
  }
}

const getProbeIcon = (probeName: string) => {
  const iconMap: Record<string, any> = {
    'Event Loop Lag': Clock,
    'Heap Usage': Memory,
    'CPU Usage': Cpu,
    'Memory Usage': Database,
    'Thread Count': Activity,
    'GIL Wait Time': Clock,
    'Goroutines': Activity,
    'OPcache Stats': Database,
    'default': Code
  }
  return iconMap[probeName] || iconMap.default
}

const getRandomColor = () => {
  const colors = ['#18a058', '#2080f0', '#f0a020', '#d03050', '#722ed1', '#eb2f96']
  return colors[Math.floor(Math.random() * colors.length)]
}

const formatTime = (timestamp: string) => {
  return new Date(timestamp).toLocaleTimeString()
}

const runProbes = async () => {
  if (!selectedProcess.value || !selectedLanguage.value) {
    message.warning('Please select both process and language')
    return
  }
  
  loading.value = true
  const result = await monitoringStore.runProbes(selectedProcess.value, selectedLanguage.value)
  
  if (result.success) {
    message.success('Probes completed successfully')
    
    // Add to historical data
    const dataPoint = result.results.reduce((acc: any, probe: any) => {
      acc[probe.name] = probe.value
      return acc
    }, {})
    historicalData.value.push(dataPoint)
    
    // Keep only last 20 data points
    if (historicalData.value.length > 20) {
      historicalData.value = historicalData.value.slice(-20)
    }
  } else {
    message.error(result.error)
  }
  
  loading.value = false
}

onMounted(() => {
  processStore.fetchProcesses()
})
</script>

<style scoped>
.language-probes {
  padding: 1rem;
}

.timestamp {
  font-size: 0.8rem;
  color: #666;
}

.chart-container {
  height: 400px;
  margin-top: 1rem;
}
</style>