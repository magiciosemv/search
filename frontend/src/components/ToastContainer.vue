<template>
  <div class="toast-container">
    <TransitionGroup name="toast">
      <div v-for="toast in toasts" :key="toast.id" class="toast" :class="toast.type" @click="remove(toast.id)">
        <div class="toast-icon">
          <svg v-if="toast.type === 'success'" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"><polyline points="20 6 9 17 4 12"/></svg>
          <svg v-else-if="toast.type === 'error'" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"><circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/></svg>
          <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"><circle cx="12" cy="12" r="10"/><line x1="12" y1="16" x2="12" y2="12"/><line x1="12" y1="8" x2="12.01" y2="8"/></svg>
        </div>
        <span class="toast-msg">{{ toast.message }}</span>
      </div>
    </TransitionGroup>
  </div>
</template>

<script setup>
import { useToast } from '../utils/toast.js'
const { toasts, remove } = useToast()
</script>

<style scoped>
.toast-container {
  position: fixed; top: 16px; right: 16px; z-index: 200;
  display: flex; flex-direction: column; gap: 8px; pointer-events: none;
}
.toast {
  display: flex; align-items: center; gap: 10px;
  padding: 12px 18px; border-radius: 10px;
  font-size: 0.8125rem; font-weight: 500;
  pointer-events: auto; cursor: pointer;
  backdrop-filter: blur(12px);
  border: 1px solid var(--border-default);
  max-width: 380px;
  box-shadow: 0 8px 24px rgba(0,0,0,0.3);
}
.toast.success { background: rgba(0, 214, 143, 0.12); border-color: rgba(0, 214, 143, 0.2); color: var(--green-bright); }
.toast.error { background: rgba(239, 68, 68, 0.12); border-color: rgba(239, 68, 68, 0.2); color: var(--red-bright); }
.toast.info { background: rgba(59, 130, 246, 0.12); border-color: rgba(59, 130, 246, 0.2); color: var(--blue-bright); }
.toast-icon { display: flex; flex-shrink: 0; }
.toast-msg { flex: 1; }

.toast-enter-active { animation: toastIn 0.3s cubic-bezier(0.16, 1, 0.3, 1); }
.toast-leave-active { animation: toastOut 0.2s ease forwards; }
.toast-move { transition: transform 0.25s cubic-bezier(0.16, 1, 0.3, 1); }

@keyframes toastIn {
  from { opacity: 0; transform: translateX(40px) scale(0.95); }
  to { opacity: 1; transform: translateX(0) scale(1); }
}
@keyframes toastOut {
  to { opacity: 0; transform: translateX(40px) scale(0.95); }
}
</style>
