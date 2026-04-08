<template>
  <Transition name="modal">
    <div v-if="modelValue" class="modal-backdrop" @click.self="$emit('update:modelValue', false)">
      <div class="modal-panel glass">
        <div class="modal-header">
          <h2 class="modal-title">{{ title }}</h2>
          <button class="modal-close" @click="$emit('update:modelValue', false)">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
          </button>
        </div>
        <div class="modal-body"><slot /></div>
        <div class="modal-footer">
          <button class="btn btn-ghost btn-sm" @click="$emit('update:modelValue', false)">Cancel</button>
          <button class="btn btn-primary btn-sm" @click="$emit('submit')">{{ submitLabel }}</button>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup>
defineProps({
  modelValue: Boolean,
  title: String,
  submitLabel: { type: String, default: 'Add' }
})
defineEmits(['update:modelValue', 'submit'])
</script>

<style scoped>
.modal-backdrop {
  position: fixed; inset: 0; z-index: 100;
  display: flex; align-items: center; justify-content: center;
  background: rgba(6, 8, 15, 0.6);
  backdrop-filter: blur(6px);
}
.modal-panel {
  width: 100%; max-width: 440px;
  border-radius: 14px;
  border: 1px solid var(--border-default);
  box-shadow: 0 32px 64px rgba(0,0,0,0.5);
}
.modal-header {
  display: flex; align-items: center; justify-content: space-between;
  padding: 18px 22px;
  border-bottom: 1px solid var(--border-subtle);
}
.modal-title { font-size: 1rem; font-weight: 700; }
.modal-close {
  display: flex; align-items: center; justify-content: center;
  width: 30px; height: 30px; border-radius: 8px;
  border: none; background: transparent;
  color: var(--text-muted); cursor: pointer;
  transition: all var(--duration-fast) ease;
}
.modal-close:hover { background: var(--bg-hover); color: var(--text-primary); }
.modal-body { padding: 22px; display: flex; flex-direction: column; gap: 16px; }
.modal-footer {
  display: flex; justify-content: flex-end; gap: 8px;
  padding: 16px 22px;
  border-top: 1px solid var(--border-subtle);
}
.modal-enter-active .modal-panel { animation: modalIn 0.3s var(--ease-out); }
.modal-leave-active .modal-panel { animation: modalIn 0.18s ease reverse; }
.modal-enter-active { animation: fadeIn 0.15s ease; }
.modal-leave-active { animation: fadeIn 0.12s ease reverse; }
@keyframes modalIn {
  from { opacity: 0; transform: translateY(12px) scale(0.97); }
  to { opacity: 1; transform: translateY(0) scale(1); }
}
</style>
