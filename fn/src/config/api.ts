import axios from 'axios'

// Configure axios with base URL from environment
// Force production backend for demo
const baseURL = 'https://gproc-backend-demo.fly.dev'

const api = axios.create({
  baseURL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Add auth token to requests
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('gproc_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Handle auth errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('gproc_token')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export default api