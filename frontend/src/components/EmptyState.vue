<template>
  <div class="bg-white dark:bg-slate-800 p-8 rounded-lg shadow text-center text-gray-500">
    <div v-if="title" class="text-gray-400 mb-2">{{ title }}</div>
    <div class="text-gray-400 text-sm">{{ message }}</div>
  </div>
</template>

<script setup>
defineProps({
  title: { type: String, default: '' },
  message: { type: String, default: 'No items yet' }
})
</script>