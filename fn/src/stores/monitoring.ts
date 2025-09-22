import { defineStore } from 'pinia'
import { ref } from 'vue'
import axios from 'axios'

export interface SystemMetrics {
  processesRunning: number
  processesTotal: number
  cpuUsage: number
  memoryUsage: number
  uptime: string
  timestamp: string
}

export interface Alert {
  id: string
  name: string
  severity: 'info' | 'warning' | 'critical'
  message: string
  resource: string
  target: string
  timestamp: string
  acknowledged: boolean
}

export interface ProbeResult {
  name: string
  language: string
  processId: string
  value: number
  unit: string
  healthy: boolean
  timestamp: string
  details?: Record<string, any>
}

export const useMonitoringStore = defineStore('monitoring', () => {
  const metrics = ref<SystemMetrics[]>([])
  const alerts = ref<Alert[]>([])
  const probeResults = ref<ProbeResult[]>([])
  const loading = ref(false)

  const fetchMetrics = async () => {
    try {
      const response = await axios.get('/api/v1/metrics')
      const newMetric = {
        ...response.data,
        timestamp: new Date().toISOString()
      }
      metrics.value.push(newMetric)
      
      // Keep only last 100 metrics
      if (metrics.value.length > 100) {
        metrics.value = metrics.value.slice(-100)
      }
    } catch (error) {
      console.error('Failed to fetch metrics:', error)
    }
  }

  const fetchAlerts = async () => {
    try {
      const response = await axios.get('/api/v1/alerts')
      alerts.value = response.data
    } catch (error) {
      console.error('Failed to fetch alerts:', error)
    }
  }

  const acknowledgeAlert = async (alertId: string) => {
    try {
      await axios.post(`/api/v1/alerts/${alertId}/acknowledge`)
      await fetchAlerts()
      return { success: true }
    } catch (error: any) {
      return { success: false, error: error.response?.data?.message }
    }
  }

  const runProbes = async (processId: string, language: string) => {
    loading.value = true
    try {
      const response = await axios.post(`/api/v1/probes/run`, {
        processId,
        language
      })
      probeResults.value = response.data
      return { success: true, results: response.data }
    } catch (error: any) {
      return { success: false, error: error.response?.data?.message }
    } finally {
      loading.value = false
    }
  }

  const getAnomalies = async () => {
    try {
      const response = await axios.get('/api/v1/ml/anomalies')
      return { success: true, anomalies: response.data }
    } catch (error: any) {
      return { success: false, error: error.response?.data?.message }
    }
  }

  return {
    metrics,
    alerts,
    probeResults,
    loading,
    fetchMetrics,
    fetchAlerts,
    acknowledgeAlert,
    runProbes,
    getAnomalies
  }
})