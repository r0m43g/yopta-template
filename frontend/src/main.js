// frontend/src/main.js
import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { createPinia } from 'pinia'
import './index.css' // подключение стилей (Tailwind, DaisyUI и т.д.)

const app = createApp(App)

app.use(createPinia())
app.use(router)

app.mount('#app')

