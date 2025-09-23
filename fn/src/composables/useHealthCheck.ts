import { ref, onMounted } from 'vue'
import api from '../config/api'

export function useHealthCheck() {
  const isBackendHealthy = ref(false)
  const backendVersion = ref('')
  const lastCheck = ref<Date | null>(null)

  const checkHealth = async () => {
    try {
      const response = await api.get('/api/v1/health')
      isBackendHealthy.value = response.status === 200
      backendVersion.value = response.data?.version || 'unknown'
      lastCheck.value = new Date()
    } catch (error) {
      isBackendHealthy.value = false
      console.warn('Backend health check failed:', error)
    }
  }

  onMounted(() => {
    checkHealth()
    // Check every 30 seconds
    setInterval(checkHealth, 30000)
  })

  return {
    isBackendHealthy,
    backendVersion,
    lastCheck,
    checkHealth
  }
}