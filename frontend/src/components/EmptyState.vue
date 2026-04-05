<template>
  <div class="empty-state">
    <div class="empty-icon">
      <svg width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
        <circle cx="12" cy="12" r="10"/>
        <path d="M8 15h8M9 9h.01M15 9h.01"/>
      </svg>
    </div>
    <div v-if="title" class="empty-title">{{ title }}</div>
    <div class="empty-msg">{{ message }}</div>
  </div>
</template>

<script setup>
defineProps({
  title: { type: String, default: '' },
  message: { type: String, default: 'No data available' }
})
</script>

<style scoped>
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48px 24px;
  text-align: center;
}

.empty-icon {
  color: var(--text-muted);
  opacity: 0.3;
  margin-bottom: 16px;
  padding: 16px;
  border-radius: 50%;
  border: 1px dashed var(--border-default);
}

.empty-title {
  font-size: 0.9375rem;
  font-weight: 600;
  color: var(--text-secondary);
  margin-bottom: 4px;
}

.empty-msg {
  font-size: 0.8125rem;
  color: var(--text-muted);
}
</style>
