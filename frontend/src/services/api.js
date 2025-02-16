// frontend/src/services/api.js
import axios from 'axios'
import { useAuthStore } from '../stores/auth'

// Создаем экземпляр axios с базовым URL
const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:6033/api',
  timeout: 10000,
})

// Request interceptor для добавления токена в заголовки
api.interceptors.request.use(
  config => {
    const authStore = useAuthStore()
    if (authStore.token) {
      config.headers['Authorization'] = `Bearer ${authStore.token}`
    }
    return config
  },
  error => Promise.reject(error)
)

// Response interceptor для обработки 401 и обновления токена
/*
api.interceptors.response.use(
  response => response,
  async error => {
    const originalRequest = error.config

    if (error.response && error.response.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true
      const authStore = useAuthStore()
      try {
        // Предполагается, что сервер предоставляет endpoint /auth/refresh для обновления токена
        const response = await api.post('/auth/refresh', { token: authStore.token })
        const newToken = response.data.token
        authStore.setToken(newToken)
        originalRequest.headers['Authorization'] = `Bearer ${newToken}`
        return api(originalRequest)
      } catch (refreshError) {
        authStore.clearToken()
        return Promise.reject(refreshError)
      }
    }
    return Promise.reject(error)
  }
)
*/
export default api

