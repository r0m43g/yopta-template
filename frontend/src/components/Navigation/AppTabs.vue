<!-- frontend/src/components/AppTabs.vue -->
<template>
  <div class="tabs">
    <router-link
      v-for="tab in tabs"
      :key="tab.name"
      :to="tab.route"
      class="tab tab-bordered"
      :class="{ 'tab-active': isActive(tab.route) }"
    >
      {{ tab.name }}
    </router-link>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const tabs = computed(() => {
  let t = [
    {
      name: 'Вход',
      route: '/login',
      icon: 'user',
      auth: false,
    },
    {
      name: 'Регистрация',
      route: '/register',
      icon: 'user-plus',
      auth: false,
    },
    {
      name: 'Профиль',
      route: '/profile',
      icon: 'user',
      auth: true,
    },
    {
      name: 'Первая',
      route: '/first',
      icon: 'document',
      auth: true,
    },
    {
      name: 'Вторая',
      route: '/second',
      icon: 'collection',
      auth: true,
    }
  ]

  return t.filter(tab => tab.auth === !!useAuthStore().token)
})

const isActive = (tabRoute) => {
  return route.path.startsWith(tabRoute)
}
</script>

<style scoped>
.tabs {
  @apply navbar bg-base-100;
}

.tab {
  @apply btn btn-ghost text-xl;
  @apply border-transparent text-gray-500 hover:text-gray-100;
  @apply transition-colors duration-200;
}

.tab-active {
  @apply border-primary text-primary;
  @apply hover:text-primary;
}

.tab:focus {
  @apply outline-none ring-2 ring-primary ring-opacity-50;
}
</style>
