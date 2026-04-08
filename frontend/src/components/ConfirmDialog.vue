<template>
  <Transition name="confirm">
    <div v-if="visible" class="confirm-backdrop" @click.self="cancel">
      <div class="confirm-panel glass">
        <div class="confirm-icon-wrap">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"><path d="M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>
        </div>
        <div class="confirm-title">{{ title }}</div>
        <div class="confirm-msg">{{ message }}</div>
        <div class="confirm-actions">
          <button class="btn btn-ghost btn-sm" @click="cancel">Cancel</button>
          <button class="btn btn-danger btn-sm" @click="ok">Confirm</button>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup>
import { useConfirm } from '../utils/confirm.js'
const { visible, title, message, ok, cancel } = useConfirm()
</script>

<style scoped>
.confirm-backdrop {
  position: fixed; inset: 0; z-index: 150;
  display: flex; align-items: center; justify-content: center;
  background: rgba(6, 8, 15, 0.6); backdrop-filter: blur(4px);
}
.confirm-panel {
  width: 100%; max-width: 360px; padding: 24px;
  text-align: center; border-radius: 14px;
  box-shadow: 0 24px 48px rgba(0,0,0,0.4);
}
.confirm-icon-wrap { color: var(--amber-bright); margin-bottom: 14px; }
.confirm-title { font-size: 1rem; font-weight: 700; margin-bottom: 6px; }
.confirm-msg { font-size: 0.8125rem; color: var(--text-secondary); margin-bottom: 20px; }
.confirm-actions { display: flex; justify-content: center; gap: 10px; }

.confirm-enter-active { animation: fadeIn 0.15s ease; }
.confirm-enter-active .confirm-panel { animation: modalIn 0.25s cubic-bezier(0.16, 1, 0.3, 1); }
.confirm-leave-active { animation: fadeIn 0.1s ease reverse; }
@keyframes modalIn {
  from { opacity: 0; transform: translateY(8px) scale(0.97); }
  to { opacity: 1; transform: translateY(0) scale(1); }
}
@keyframes fadeIn {
  from { opacity: 0; } to { opacity: 1; }
}
</style>
