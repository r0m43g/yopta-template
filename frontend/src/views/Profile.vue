<!-- frontend/src/views/Profile.vue -->
<template>
  <div class="max-w-2xl h-full mx-auto py-8">
    <h1 class="text-3xl font-bold mb-8">Профиль пользователя</h1>
    <div class="card bg-base-100 shadow-xl">
      <div class="card-body">
        <p class="text-lg">Имя: {{ username }}</p>
        <p class="text-lg">Email: {{ email }}</p>
        <p class="mt-4">...</p>
        <div class="form-control">
          <button @click="logout" class="btn btn-primary">Выйти</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { useAuthStore } from '@/stores/auth'
import api from '@/services/api'
export default {
  name: "Profile",
  computed: {
    username() {
      const authStore = useAuthStore()
      return authStore.username
    },
    email() {
      const authStore = useAuthStore()
      return authStore.email
    }
  },
  methods: {
    logout() {
      const authStore = useAuthStore()
      authStore.clearToken()
      this.$router.push({ name: 'login' })
    }
  },
  async mounted() {
    const authStore = useAuthStore()
    const response = await api.get('/profile')
    authStore.setUser(response.data)
  }
}
</script>

<style scoped>
/* стили для профиля */
</style>

