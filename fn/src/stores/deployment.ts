import { defineStore } from 'pinia'
import { ref } from 'vue'
import axios from 'axios'

export interface Deployment {
  id: string
  processName: string
  strategy: 'blue-green' | 'rolling' | 'canary'
  version: string
  status: 'pending' | 'in-progress' | 'completed' | 'failed' | 'rolled-back'
  progress: number
  startTime: string
  endTime?: string
  rollbackOnFail: boolean
  healthChecks: boolean
  logs: string[]
}

export interface DeploymentTemplate {
  name: string
  language: string
  command: string
  args: string[]
  env: Record<string, string>
  healthCheck?: {
    url: string
    interval: string
  }
}

export const useDeploymentStore = defineStore('deployment', () => {
  const deployments = ref<Deployment[]>([])
  const templates = ref<DeploymentTemplate[]>([])
  const loading = ref(false)

  const fetchDeployments = async () => {
    loading.value = true
    try {
      const response = await axios.get('/api/v1/deployments')
      deployments.value = response.data
    } catch (error) {
      console.error('Failed to fetch deployments:', error)
    } finally {
      loading.value = false
    }
  }

  const createDeployment = async (deploymentData: {
    processName: string
    strategy: string
    version: string
    rollbackOnFail: boolean
  }) => {
    try {
      const response = await axios.post('/api/v1/deployments', deploymentData)
      await fetchDeployments()
      return { success: true, deployment: response.data }
    } catch (error: any) {
      return { success: false, error: error.response?.data?.message }
    }
  }

  const rollbackDeployment = async (deploymentId: string) => {
    try {
      await axios.post(`/api/v1/deployments/${deploymentId}/rollback`)
      await fetchDeployments()
      return { success: true }
    } catch (error: any) {
      return { success: false, error: error.response?.data?.message }
    }
  }

  const fetchTemplates = async () => {
    try {
      const response = await axios.get('/api/v1/templates')
      templates.value = response.data
    } catch (error) {
      console.error('Failed to fetch templates:', error)
    }
  }

  const createTemplate = async (templateData: DeploymentTemplate) => {
    try {
      await axios.post('/api/v1/templates', templateData)
      await fetchTemplates()
      return { success: true }
    } catch (error: any) {
      return { success: false, error: error.response?.data?.message }
    }
  }

  return {
    deployments,
    templates,
    loading,
    fetchDeployments,
    createDeployment,
    rollbackDeployment,
    fetchTemplates,
    createTemplate
  }
})