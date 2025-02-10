<!-- frontend/components/PasswordStrength.vue -->
<template>
  <div class="flex flex-col">
    <label class="input input-bordered flex items-center gap-2 my-2" :class="strengthClass">
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
      <input v-model="password" type="password" @input="checkStrength" placeholder="Пароль" class="grow" required />
    </label>
    <div class="flex flex-row gap-1">
      <div v-for="n in 4" :key="n" class="h-px w-1/4 transition-colors" :class="n <= strength ? strengthColor : 'bg-red-500'"></div>
    </div>
    <span class="bg-green-500"></span><span class="bg-yellow-500"></span><span class="bg-orange-500"></span><span class="bg-red-500"></span>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const password = ref('')
const strength = ref(0)

const checkStrength = () => {
  let score = 0
  if (password.value.length >= 8) score++
  if (/[A-Z]/.test(password.value)) score++
  if (/[0-9]/.test(password.value)) score++
  if (/[^A-Za-z0-9]/.test(password.value)) score++

  strength.value = score
}

const strengthColor = computed(() => {
  const colors = ['red', 'orange', 'yellow', 'green']
return `bg-${colors[strength.value - 1]}-500`
})

const strengthClass = computed(() => {
  return {
  'text-red-500': strength.value <2,
  'text-orange-500': strength.value === 2,
  'text-green-500': strength.value > 2
  }
})
</script>
