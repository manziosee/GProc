import { defineStore } from 'pinia'
import { ref } from 'vue'
import axios from 'axios'

export interface ScheduledTask {
  id: string
  name: string
  command: string
  args: string[]
  cron: string
  enabled: boolean
  description?: string
  timeout: string
  env: Record<string, string>
  workingDir?: string
  nextRun: string
  lastRun?: string
  lastResult?: 'success' | 'failure' | 'timeout'
  runCount: number
  created: string
}

export interface TaskExecution {
  id: string
  taskId: string
  startTime: string
  endTime?: string
  status: 'running' | 'completed' | 'failed' | 'timeout'
  exitCode?: number
  output: string
  error?: string
}

export const useSchedulerStore = defineStore('scheduler', () => {
  const tasks = ref<ScheduledTask[]>([])
  const executions = ref<TaskExecution[]>([])
  const loading = ref(false)

  const fetchTasks = async () => {
    loading.value = true
    try {
      const response = await axios.get('/api/v1/scheduler/tasks')
      tasks.value = response.data
    } catch (error) {
      console.error('Failed to fetch scheduled tasks:', error)
    } finally {
      loading.value = false
    }
  }

  const createTask = async (taskData: {
    name: string
    command: string
    args: string[]
    cron: string
    description?: string
    timeout?: string
    env?: Record<string, string>
    workingDir?: string
  }) => {
    try {
      await axios.post('/api/v1/scheduler/tasks', taskData)
      await fetchTasks()
      return { success: true }
    } catch (error: any) {
      return { success: false, error: error.response?.data?.message }
    }
  }

  const updateTask = async (taskId: string, taskData: Partial<ScheduledTask>) => {
    try {
      await axios.put(`/api/v1/scheduler/tasks/${taskId}`, taskData)
      await fetchTasks()
      return { success: true }
    } catch (error: any) {
      return { success: false, error: error.response?.data?.message }
    }
  }

  const deleteTask = async (taskId: string) => {
    try {
      await axios.delete(`/api/v1/scheduler/tasks/${taskId}`)
      await fetchTasks()
      return { success: true }
    } catch (error: any) {
      return { success: false, error: error.response?.data?.message }
    }
  }

  const runTaskNow = async (taskId: string) => {
    try {
      const response = await axios.post(`/api/v1/scheduler/tasks/${taskId}/run`)
      return { success: true, executionId: response.data.executionId }
    } catch (error: any) {
      return { success: false, error: error.response?.data?.message }
    }
  }

  const fetchExecutions = async (taskId?: string) => {
    try {
      const url = taskId 
        ? `/api/v1/scheduler/executions?taskId=${taskId}`
        : '/api/v1/scheduler/executions'
      const response = await axios.get(url)
      executions.value = response.data
    } catch (error) {
      console.error('Failed to fetch task executions:', error)
    }
  }

  const getExecutionLogs = async (executionId: string) => {
    try {
      const response = await axios.get(`/api/v1/scheduler/executions/${executionId}/logs`)
      return { success: true, logs: response.data.logs }
    } catch (error: any) {
      return { success: false, error: error.response?.data?.message }
    }
  }

  return {
    tasks,
    executions,
    loading,
    fetchTasks,
    createTask,
    updateTask,
    deleteTask,
    runTaskNow,
    fetchExecutions,
    getExecutionLogs
  }
})