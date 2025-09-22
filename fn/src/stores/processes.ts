import { defineStore } from 'pinia'
import { ref } from 'vue'
import axios from 'axios'

export interface Process {
  id: string
  name: string
  command: string
  args: string[]
  status: 'running' | 'stopped' | 'failed' | 'starting' | 'stopping'
  pid: number
  restarts: number
  uptime: string
  cpu: number
  memory: number
  group: string
  autoRestart: boolean
  maxRestarts: number
  healthCheck?: {
    url: string
    interval: string
    timeout: string
    retries: number
    status: 'healthy' | 'unhealthy' | 'unknown'
  }
  resourceLimits?: {
    memoryMB: number
    cpuLimit: number
  }
  notifications?: {
    email?: string
    slack?: string
  }
  env: Record<string, string>
  workingDir: string
  logFile: string
  created: string
  lastRestart?: string
}

export const useProcessStore = defineStore('processes', () => {
  const processes = ref<Process[]>([])
  const loading = ref(false)

  const fetchProcesses = async () => {
    loading.value = true
    try {
      const response = await axios.get('/api/v1/processes')
      processes.value = response.data
    } catch (error) {
      console.error('Failed to fetch processes:', error)
    } finally {
      loading.value = false
    }
  }

  const startProcess = async (processId: string) => {
    try {
      await axios.post(`/api/v1/processes/${processId}/start`)
      await fetchProcesses()
      return { success: true }
    } catch (error: any) {
      return { success: false, error: error.response?.data?.message }
    }
  }

  const stopProcess = async (processId: string) => {
    try {
      await axios.post(`/api/v1/processes/${processId}/stop`)
      await fetchProcesses()
      return { success: true }
    } catch (error: any) {
      return { success: false, error: error.response?.data?.message }
    }
  }

  const restartProcess = async (processId: string) => {
    try {
      await axios.post(`/api/v1/processes/${processId}/restart`)
      await fetchProcesses()
      return { success: true }
    } catch (error: any) {
      return { success: false, error: error.response?.data?.message }
    }
  }

  const createProcess = async (processData: Partial<Process>) => {
    try {
      await axios.post('/api/v1/processes', processData)
      await fetchProcesses()
      return { success: true }
    } catch (error: any) {
      return { success: false, error: error.response?.data?.message }
    }
  }

  const deleteProcess = async (processId: string) => {
    try {
      await axios.delete(`/api/v1/processes/${processId}`)
      await fetchProcesses()
      return { success: true }
    } catch (error: any) {
      return { success: false, error: error.response?.data?.message }
    }
  }

  const getProcessLogs = async (processId: string, lines: number = 100) => {
    try {
      const response = await axios.get(`/api/v1/processes/${processId}/logs?lines=${lines}`)
      return { success: true, logs: response.data.logs }
    } catch (error: any) {
      return { success: false, error: error.response?.data?.message }
    }
  }

  return {
    processes,
    loading,
    fetchProcesses,
    startProcess,
    stopProcess,
    restartProcess,
    createProcess,
    deleteProcess,
    getProcessLogs
  }
})