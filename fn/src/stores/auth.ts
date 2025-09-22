import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import axios from 'axios'

export interface User {
  id: string
  username: string
  email: string
  roles: string[]
  permissions: string[]
  enabled: boolean
  mfaEnabled: boolean
  lastSeen: string
  created: string
}

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const token = ref<string | null>(localStorage.getItem('gproc_token'))
  const isAuthenticated = computed(() => !!token.value && !!user.value)

  const login = async (username: string, password: string, mfaCode?: string) => {
    try {
      const response = await axios.post('/api/v1/auth/login', {
        username,
        password,
        mfa_code: mfaCode
      })
      
      token.value = response.data.token
      user.value = response.data.user
      localStorage.setItem('gproc_token', token.value!)
      
      axios.defaults.headers.common['Authorization'] = `Bearer ${token.value}`
      
      return { success: true }
    } catch (error: any) {
      return { 
        success: false, 
        error: error.response?.data?.message || 'Login failed',
        requiresMFA: error.response?.status === 428
      }
    }
  }

  const logout = () => {
    user.value = null
    token.value = null
    localStorage.removeItem('gproc_token')
    delete axios.defaults.headers.common['Authorization']
  }

  const hasPermission = (resource: string, action: string, scope: string = '*') => {
    if (!user.value) return false
    return user.value.roles.includes('admin') || 
           user.value.permissions.some(p => 
             p.includes(`${resource}:${action}:${scope}`)
           )
  }

  if (token.value) {
    axios.defaults.headers.common['Authorization'] = `Bearer ${token.value}`
  }

  return {
    user,
    token,
    isAuthenticated,
    login,
    logout,
    hasPermission
  }
})