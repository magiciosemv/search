<template>
  <Transition name="modal">
    <div v-if="modelValue" class="modal-backdrop" @click.self="$emit('update:modelValue', false)">
      <div class="modal-panel">
        <div class="modal-header">
          <h2 class="modal-title">{{ title }}</h2>
          <button class="modal-close" @click="$emit('update:modelValue', false)">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
          </button>
        </div>

        <div class="modal-body">
          <slot />
        </div>

        <div class="modal-footer">
          <button class="btn btn-ghost" @click="$emit('update:modelValue', false)">Cancel</button>
          <button class="btn btn-primary" @click="$emit('submit')">{{ submitLabel }}</button>
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
  position: fixed;
  inset: 0;
  z-index: 100;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
}

.modal-panel {
  width: 100%;
  max-width: 460px;
  background: var(--bg-surface);
  border: 1px solid var(--border-default);
  border-radius: 16px;
  box-shadow: 0 0 40px rgba(6, 182, 212, 0.08), 0 24px 48px rgba(0,0,0,0.4);
  position: relative;
  overflow: hidden;
}

.modal-panel::before {
  content: '';
  position: absolute;
  top: 0; left: 0; right: 0;
  height: 1px;
  background: linear-gradient(90deg, transparent, var(--cyan-base), transparent);
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid var(--border-subtle);
}

.modal-title {
  font-size: 1.125rem;
  font-weight: 700;
}

.modal-close {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 6px;
  border: none;
  background: transparent;
  color: var(--text-muted);
  cursor: pointer;
  transition: all 0.15s ease;
}

.modal-close:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.modal-body {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  padding: 20px 24px;
  border-top: 1px solid var(--border-subtle);
}

.modal-enter-active .modal-panel { animation: modalIn 0.3s var(--ease-out); }
.modal-leave-active .modal-panel { animation: modalIn 0.2s ease reverse; }
.modal-enter-active { animation: fadeIn 0.2s ease; }
.modal-leave-active { animation: fadeIn 0.15s ease reverse; }

@keyframes modalIn {
  from { opacity: 0; transform: translateY(16px) scale(0.97); }
  to { opacity: 1; transform: translateY(0) scale(1); }
}
</style>
