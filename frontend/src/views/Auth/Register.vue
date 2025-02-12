<!-- frontend/src/views/Auth/Register.vue -->
<template>
  <div class="max-w-2xl h-full mx-auto py-8 flex items-center justify-center">
    <div class="card lg:card-side glass  shadow-xl">
      <figure>
        <img src="https://picsum.photos/300/430" alt="Placeholder image" class="rounded-lg shadow-lg" />
      </figure>
      <div class="card-body">
        <h1 class="text-3xl font-bold mb-8">Регистрация</h1>
        <form @submit.prevent="register">
          <label class="input input-bordered flex items-center gap-2 my-2">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 16 16"
              fill="currentColor"
              class="h-4 w-4 opacity-70">
              <path
                d="M8 8a3 3 0 1 0 0-6 3 3 0 0 0 0 6ZM12.735 14c.618 0 1.093-.561.872-1.139a6.002 6.002 0 0 0-11.215 0c-.22.578.254 1.139.872 1.139h9.47Z" />
            </svg>
            <input v-model="username" type="text" class="grow" placeholder="Имя пользователя" required />
          </label>
          <label class="input input-bordered flex items-center gap-2 my-2">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 16 16"
              fill="currentColor"
              class="h-4 w-4 opacity-70">
              <path
                d="M8 8a3 3 0 1 0 0-6 3 3 0 0 0 0 6ZM12.735 14c.618 0 1.093-.561.872-1.139a6.002 6.002 0 0 0-11.215 0c-.22.578.254 1.139.872 1.139h9.47Z" />
            </svg>
            <input v-model="email" type="email" class="grow" placeholder="Email" required />
          </label>
          <label class="input input-bordered flex items-center gap-2 my-2">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 16 16"
              fill="currentColor"
              class="h-4 w-4 opacity-70">
              <path
                fill-rule="evenodd"
                d="M14 6a4 4 0 0 1-4.899 3.899l-1.955 1.955a.5.5 0 0 1-.353.146H5v1.5a.5.5 0 0 1-.5.5h-2a.5.5 0 0 1-.5-.5v-2.293a.5.5 0 0 1 .146-.353l3.955-3.955A4 4 0 1 1 14 6Zm-4-2a.75.75 0 0 0 0 1.5.5.5 0 0 1 .5.5.75.75 0 0 0 1.5 0 2 2 0 0 0-2-2Z"
                clip-rule="evenodd" />
            </svg>
            <input v-model="password" type="password" placeholder="Пароль" class="grow" required />
          </label>
          <password-strength :password="password" />
          <label class="input input-bordered flex items-center gap-2 my-2">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 16 16"
              fill="currentColor"
              class="h-4 w-4 opacity-70">
              <path
                fill-rule="evenodd"
                d="M14 6a4 4 0 0 1-4.899 3.899l-1.955 1.955a.5.5 0 0 1-.353.146H5v1.5a.5.5 0 0 1-.5.5h-2a.5.5 0 0 1-.5-.5v-2.293a.5.5 0 0 1 .146-.353l3.955-3.955A4 4 0 1 1 14 6Zm-4-2a.75.75 0 0 0 0 1.5.5.5 0 0 1 .5.5.75.75 0 0 0 1.5 0 2 2 0 0 0-2-2Z"
                clip-rule="evenodd" />
            </svg>
            <input v-model="confirmPassword" type="password" placeholder="Подтвердите пароль" class="grow" required />
          </label>
          <div class="form-control my-6">
            <button class="btn btn-outline flex-1 mx-2" type="reset">Сбросить</button>
            <button class="btn btn-outline btn-primary flex-1 mx-2" type="submit">Регистрация</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import api from '@/services/api'
import PasswordStrength from '@/components/Auth/PasswordStrength.vue'
export default {
  components: {
    PasswordStrength
  },
  name: "Register",
  data() {
    return {
      username: '',
      email: '',
      password: '',
      confirmPassword: ''
    }
  },
  methods: {
    async register() {
      if (this.password !== this.confirmPassword) {
        alert('Пароли не совпадают')
        return
      }
      try {
        await api.post('/register', { username: this.username, email: this.email, password: this.password })
        this.$router.push({ name: 'login' })
      } catch (error) {
        console.error("Ошибка регистрации", error)
      }
    }
  }
}
</script>

<style scoped>
.form-control {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
}
/* стили для страницы регистрации */
</style>

