// frontend/src/stores/auth.js
import { defineStore } from 'pinia'
import { jwtDecode } from 'jwt-decode'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('token') || null,
    user: null,
    username: null,
    email: null
  }),
  actions: {
    setToken(token) {
      this.token = token
      localStorage.setItem('token', token)
      try {
        this.user = jwtDecode(token)
      } catch (error) {
        this.user = null
      }
    },
    clearToken() {
      this.token = null
      this.user = null
      this.username = null
      this.email = null
      localStorage.removeItem('token')
    },
    isTokenValid() {
      if (!this.token) return false
      try {
        const decoded = jwtDecode(this.token)
        return decoded.exp * 1000 > Date.now()
      } catch (error) {
        return false
      }
    },
    setUser(user) {
      this.username = user.username
      this.email = user.email
    }
  },
})

