import { createApp } from 'vue'
import { createPinia } from 'pinia'
import axios from 'axios'
import App from './App.vue'
import router from './router'
import './style.css'

// Configure axios defaults
axios.defaults.baseURL = 'http://localhost:8080'
axios.defaults.timeout = 10000

// Add request interceptor for authentication
axios.interceptors.request.use((config) => {
  const token = localStorage.getItem('gproc_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Add response interceptor for error handling
axios.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('gproc_token')
      router.push('/login')
    }
    return Promise.reject(error)
  }
)

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)

app.mount('#app')