<template>
  <div class="empty-state">
    <div class="empty-icon">
      <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
        <circle cx="12" cy="12" r="10"/><path d="M8 15h8M9 9h.01M15 9h.01"/>
      </svg>
    </div>
    <div v-if="title" class="empty-title">{{ title }}</div>
    <div class="empty-msg">{{ message }}</div>
  </div>
</template>

<script setup>
defineProps({
  title: { type: String, default: '' },
  message: { type: String, default: 'No data' }
})
</script>

<style scoped>
.empty-state {
  display: flex; flex-direction: column;
  align-items: center; justify-content: center;
  padding: 48px 20px; text-align: center;
}
.empty-icon {
  color: var(--text-dim); opacity: 0.35;
  padding: 16px;
  border-radius: 50%;
  border: 1px dashed var(--border-default);
  margin-bottom: 16px;
}
.empty-title {
  font-size: 0.875rem; font-weight: 700;
  color: var(--text-secondary); margin-bottom: 4px;
}
.empty-msg { font-size: 0.8125rem; color: var(--text-muted); }
</style>
