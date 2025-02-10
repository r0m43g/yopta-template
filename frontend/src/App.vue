<!-- frontend/src/App.vue -->
<template>
  <div class="flex flex-col items-center justify-center h-screen w-screen">
    <AppTabs />
    <router-view />
  </div>
</template>

<script>
import AppTabs from './components/Navigation/AppTabs.vue'

export default {
  name: "App",
  components: {
    AppTabs
  }
}
</script>

<style>
/* глобальные стили */
</style>

