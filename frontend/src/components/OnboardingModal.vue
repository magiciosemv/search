<template>
  <Transition name="onboarding">
    <div v-if="modelValue" class="onboarding-backdrop">
      <div class="onboarding-panel glass">
        <div class="onboarding-header">
          <div class="onboarding-brand">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
              <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5"/>
            </svg>
          </div>
          <h2 class="onboarding-title">Welcome to Solana Monitor</h2>
          <div class="step-dots">
            <span v-for="i in 3" :key="i" class="dot" :class="{ active: step >= i }"></span>
          </div>
        </div>

        <div class="onboarding-body">
          <Transition name="step" mode="out-in">
            <div v-if="step === 1" key="1" class="step-content">
              <div class="step-icon">
                <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"><path d="M21 12V7H5a2 2 0 010-4h14v4"/><path d="M3 5v14a2 2 0 002 2h16v-5"/><path d="M18 12a2 2 0 000 4h4v-4z"/></svg>
              </div>
              <h3>Add Your Wallet</h3>
              <p>Start by adding a Solana wallet address to monitor its balance in real-time.</p>
            </div>
            <div v-else-if="step === 2" key="2" class="step-content">
              <div class="step-icon">
                <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"><path d="M18 8A6 6 0 006 8c0 7-3 9-3 9h18s-3-2-3-9"/><path d="M13.73 21a2 2 0 01-3.46 0"/></svg>
              </div>
              <h3>Set Up Notifications</h3>
              <p>Configure Telegram or email channels to receive instant balance change alerts.</p>
            </div>
            <div v-else key="3" class="step-content">
              <div class="step-icon success-icon">
                <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"><path d="M22 11.08V12a10 10 0 11-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
              </div>
              <h3>You're All Set!</h3>
              <p>The monitor will check balances every 30 seconds and notify you of changes.</p>
            </div>
          </Transition>
        </div>

        <div class="onboarding-footer">
          <button v-if="step < 3" class="btn btn-ghost btn-sm" @click="skip">Skip</button>
          <span class="step-counter font-mono">Step {{ step }} / 3</span>
          <button v-if="step < 3" class="btn btn-primary btn-sm" @click="step++">Next</button>
          <button v-else class="btn btn-primary btn-sm" @click="finish">Get Started</button>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup>
import { ref } from 'vue'

defineProps({ modelValue: Boolean })
const emit = defineEmits(['update:modelValue'])

const step = ref(1)

const skip = () => finish()

const finish = () => {
  localStorage.setItem('onboarded', 'true')
  emit('update:modelValue', false)
}
</script>

<style scoped>
.onboarding-backdrop {
  position: fixed; inset: 0; z-index: 300;
  display: flex; align-items: center; justify-content: center;
  background: rgba(6, 8, 15, 0.7); backdrop-filter: blur(6px);
}
.onboarding-panel {
  width: 100%; max-width: 440px; border-radius: 16px;
  box-shadow: 0 32px 64px rgba(0,0,0,0.5); overflow: hidden;
}
.onboarding-header {
  text-align: center; padding: 28px 28px 0;
}
.onboarding-brand {
  color: var(--blue-base); margin-bottom: 12px;
  display: flex; justify-content: center;
}
.onboarding-title {
  font-size: 1.125rem; font-weight: 700; margin-bottom: 16px;
}
.step-dots {
  display: flex; justify-content: center; gap: 8px;
}
.dot {
  width: 8px; height: 8px; border-radius: 50%;
  background: var(--border-default); transition: all 0.3s ease;
}
.dot.active { background: var(--blue-base); box-shadow: 0 0 6px rgba(59,130,246,0.3); }

.onboarding-body {
  padding: 28px; min-height: 200px;
  display: flex; align-items: center; justify-content: center;
}
.step-content { text-align: center; }
.step-icon {
  color: var(--blue-bright); margin-bottom: 16px;
  display: flex; justify-content: center;
}
.success-icon { color: var(--green-bright); }
.step-content h3 { font-size: 1.0625rem; font-weight: 700; margin-bottom: 8px; }
.step-content p { font-size: 0.8125rem; color: var(--text-secondary); line-height: 1.6; max-width: 320px; margin: 0 auto; }

.onboarding-footer {
  display: flex; align-items: center; justify-content: space-between;
  padding: 18px 28px; border-top: 1px solid var(--border-subtle);
}
.step-counter { font-size: 0.6875rem; color: var(--text-dim); }

/* Transitions */
.onboarding-enter-active { animation: fadeIn 0.25s ease; }
.onboarding-enter-active .onboarding-panel { animation: modalIn 0.35s cubic-bezier(0.16, 1, 0.3, 1); }
.onboarding-leave-active { animation: fadeIn 0.15s ease reverse; }
.step-enter-active { animation: stepIn 0.3s cubic-bezier(0.16, 1, 0.3, 1); }
.step-leave-active { animation: stepOut 0.15s ease forwards; }

@keyframes fadeIn { from { opacity: 0; } to { opacity: 1; } }
@keyframes modalIn {
  from { opacity: 0; transform: translateY(16px) scale(0.96); }
  to { opacity: 1; transform: translateY(0) scale(1); }
}
@keyframes stepIn {
  from { opacity: 0; transform: translateX(20px); }
  to { opacity: 1; transform: translateX(0); }
}
@keyframes stepOut {
  to { opacity: 0; transform: translateX(-20px); }
}
</style>
