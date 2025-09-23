import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '../config/api'

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

  const register = async (username: string, password: string, email: string) => {
    try {
      await api.post('/api/v1/auth/register', {
        username,
        password,
        email
      })
      
      return { success: true }
    } catch (error: any) {
      return { 
        success: false, 
        error: error.response?.data?.message || 'Registration failed'
      }
    }
  }

  const login = async (username: string, password: string, mfaCode?: string) => {
    try {
      const response = await api.post('/api/v1/auth/login', {
        username,
        password,
        mfa_code: mfaCode
      })
      
      token.value = response.data.token
      user.value = response.data.user
      localStorage.setItem('gproc_token', token.value!)
      
      // Token is handled by interceptor
      
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
    // Token removal handled by interceptor
  }

  const hasPermission = (resource: string, action: string, scope: string = '*') => {
    if (!user.value) return false
    return user.value.roles.includes('admin') || 
           user.value.permissions.some(p => 
             p.includes(`${resource}:${action}:${scope}`)
           )
  }

  // Token handling moved to API interceptor

  return {
    user,
    token,
    isAuthenticated,
    register,
    login,
    logout,
    hasPermission
  }
})