<!-- frontend/src/components/ValidationError.vue -->
<template>
  <div v-if="errors && errors.length" class="validation-errors">
    <div
      v-for="(error, index) in errors"
      :key="index"
      class="alert alert-error shadow-lg mb-2"
    >
      <div class="flex items-center gap-2">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="stroke-current flex-shrink-0 h-6 w-6"
          fill="none"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
          />
        </svg>
        <span>{{ error }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  errors: {
    type: Array,
    default: () => []
  }
})
</script>

<style scoped>
.validation-errors {
  @apply mb-4;
}

.alert-error {
  @apply bg-error text-error-content;
}

.alert-error svg {
  @apply text-error-content;
}
</style>
