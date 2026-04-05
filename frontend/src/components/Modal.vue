<template>
  <div v-if="modelValue" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" @click.self="$emit('update:modelValue', false)">
    <div class="bg-white dark:bg-slate-800 p-6 rounded-lg w-full max-w-md">
      <h2 class="text-xl font-bold mb-4">{{ title }}</h2>
      <div class="space-y-4">
        <slot />
      </div>
      <div class="flex justify-end gap-2 mt-6">
        <button @click="$emit('update:modelValue', false)"
          class="px-4 py-2 border rounded-lg hover:bg-gray-100 dark:hover:bg-slate-700">
          Cancel
        </button>
        <button @click="$emit('submit')"
          class="px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700">
          {{ submitLabel }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  modelValue: Boolean,
  title: String,
  submitLabel: { type: String, default: 'Add' }
})

defineEmits(['update:modelValue', 'submit'])
</script>