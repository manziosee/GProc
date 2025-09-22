import { defineStore } from 'pinia'
import { ref } from 'vue'
import axios from 'axios'

export interface Role {
  name: string
  permissions: Permission[]
  description?: string
}

export interface Permission {
  resource: string
  actions: string[]
  scope: string
}

export interface AuditEvent {
  timestamp: string
  userId: string
  username: string
  action: string
  resource: string
  target: string
  result: string
  details?: string
  ip?: string
  userAgent?: string
}

export interface SSOConfig {
  enabled: boolean
  provider: 'oauth2' | 'saml'
  clientId?: string
  entityId?: string
  metadataUrl?: string
}

export const useSecurityStore = defineStore('security', () => {
  const users = ref<any[]>([])
  const roles = ref<Role[]>([])
  const auditEvents = ref<AuditEvent[]>([])
  const ssoConfig = ref<SSOConfig | null>(null)
  const loading = ref(false)

  const fetchUsers = async () => {
    loading.value = true
    try {
      const response = await axios.get('/api/v1/users')
      users.value = response.data
    } catch (error) {
      console.error('Failed to fetch users:', error)
    } finally {
      loading.value = false
    }
  }

  const createUser = async (userData: {
    username: string
    email: string
    password: string
    roles: string[]
    enabled: boolean
  }) => {
    try {
      await axios.post('/api/v1/users', userData)
      await fetchUsers()
      return { success: true }
    } catch (error: any) {
      return { success: false, error: error.response?.data?.message }
    }
  }

  const updateUser = async (userId: string, userData: Partial<any>) => {
    try {
      await axios.put(`/api/v1/users/${userId}`, userData)
      await fetchUsers()
      return { success: true }
    } catch (error: any) {
      return { success: false, error: error.response?.data?.message }
    }
  }

  const deleteUser = async (userId: string) => {
    try {
      await axios.delete(`/api/v1/users/${userId}`)
      await fetchUsers()
      return { success: true }
    } catch (error: any) {
      return { success: false, error: error.response?.data?.message }
    }
  }

  const fetchRoles = async () => {
    try {
      const response = await axios.get('/api/v1/roles')
      roles.value = response.data
    } catch (error) {
      console.error('Failed to fetch roles:', error)
    }
  }

  const createRole = async (roleData: Role) => {
    try {
      await axios.post('/api/v1/roles', roleData)
      await fetchRoles()
      return { success: true }
    } catch (error: any) {
      return { success: false, error: error.response?.data?.message }
    }
  }

  const fetchAuditEvents = async (limit: number = 100) => {
    try {
      const response = await axios.get(`/api/v1/audit?limit=${limit}`)
      auditEvents.value = response.data
    } catch (error) {
      console.error('Failed to fetch audit events:', error)
    }
  }

  const enableMFA = async (userId: string) => {
    try {
      const response = await axios.post(`/api/v1/users/${userId}/mfa/enable`)
      return { success: true, qrCode: response.data.qrCode, secret: response.data.secret }
    } catch (error: any) {
      return { success: false, error: error.response?.data?.message }
    }
  }

  const verifyMFA = async (userId: string, code: string) => {
    try {
      await axios.post(`/api/v1/users/${userId}/mfa/verify`, { code })
      return { success: true }
    } catch (error: any) {
      return { success: false, error: error.response?.data?.message }
    }
  }

  const fetchSSOConfig = async () => {
    try {
      const response = await axios.get('/api/v1/sso/config')
      ssoConfig.value = response.data
    } catch (error) {
      console.error('Failed to fetch SSO config:', error)
    }
  }

  const updateSSOConfig = async (config: SSOConfig) => {
    try {
      await axios.put('/api/v1/sso/config', config)
      ssoConfig.value = config
      return { success: true }
    } catch (error: any) {
      return { success: false, error: error.response?.data?.message }
    }
  }

  return {
    users,
    roles,
    auditEvents,
    ssoConfig,
    loading,
    fetchUsers,
    createUser,
    updateUser,
    deleteUser,
    fetchRoles,
    createRole,
    fetchAuditEvents,
    enableMFA,
    verifyMFA,
    fetchSSOConfig,
    updateSSOConfig
  }
})