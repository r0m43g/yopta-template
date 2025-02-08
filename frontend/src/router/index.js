import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { jwtDecode } from 'jwt-decode'

import Login from '../views/Auth/Login.vue'
import Register from '../views/Auth/Register.vue'
import Profile from '../views/Profile.vue'
import First from '../views/First.vue'
import Second from '../views/Second.vue'

const routes = [
  { path: '/', name: 'home', component: Login },
  { path: '/first', name: 'first', component: First, meta: { requiresAuth: true } },
  { path: '/second', name: 'second', component: Second, meta: { requiresAuth: true } },
  { path: '/login', name: 'login', component: Login },
  { path: '/register', name: 'register', component: Register },
  { path: '/profile', name: 'profile', component: Profile, meta: { requiresAuth: true } },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// Функция для проверки валидности JWT
function isTokenValid(token) {
  try {
    const decoded = jwtDecode(token)
    return decoded.exp * 1000 > Date.now()
  } catch (error) {
    return false
  }
}

router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()

  // Если токен не в store, загружаем его из localStorage
  let token = authStore.token
  if (!token) {
    token = localStorage.getItem('token')
    if (token) {
      authStore.setToken(token)
    }
  }

  if (to.meta.requiresAuth) {
    if (!token || !isTokenValid(token)) {
      authStore.clearToken()
      return next({ name: 'login' })
    }
  }
  next()
})

export default router

