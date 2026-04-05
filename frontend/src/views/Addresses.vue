<template>
  <div class="page-container">
    <div class="page-header">
      <div>
        <h1 class="page-title">Wallets</h1>
        <p class="page-subtitle">Monitored Solana addresses</p>
      </div>
      <button class="btn btn-primary" @click="showAddModal = true">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
        Add Wallet
      </button>
    </div>

    <div class="wallet-grid stagger">
      <div v-for="addr in addresses" :key="addr.id"
        class="wallet-card card animate-fade-up">
        <div class="wallet-card-top">
          <div class="wallet-identity">
            <div class="wallet-avatar">{{ (addr.label || 'W')[0].toUpperCase() }}</div>
            <div>
              <div class="wallet-label">{{ addr.label || 'Unnamed' }}</div>
              <div class="wallet-addr font-mono" :title="addr.address">{{ truncateAddress(addr.address) }}</div>
            </div>
          </div>
          <div class="wallet-balance">
            <div class="balance-value font-mono">{{ formatSOL(addr.balance) }}</div>
            <div class="balance-unit">SOL</div>
          </div>
        </div>

        <div class="wallet-card-footer">
          <div class="threshold-group">
            <span class="threshold-label">Threshold</span>
            <input v-model.number="addr.threshold" @change="saveThreshold(addr)"
              type="number" step="0.1" min="0"
              class="input-field input-mono threshold-input" />
            <span class="threshold-unit">SOL</span>
          </div>
          <div class="wallet-actions">
            <button @click="refreshBalance(addr.id)" class="btn btn-ghost btn-sm" title="Refresh balance">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><polyline points="23 4 23 10 17 10"/><path d="M20.49 15a9 9 0 11-2.12-9.36L23 10"/></svg>
              Refresh
            </button>
            <button @click="deleteAddress(addr.id)" class="btn btn-danger btn-sm" title="Remove wallet">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2"/></svg>
              Delete
            </button>
          </div>
        </div>

        <div class="wallet-updated text-muted" v-if="addr.updated_at">
          Last updated: {{ formatDate(addr.updated_at) }}
        </div>
      </div>
    </div>

    <EmptyState v-if="addresses.length === 0" title="No wallets yet"
      message="Add a Solana address to start monitoring" />

    <Modal v-model="showAddModal" title="Add Wallet" submit-label="Add" @submit="addAddress">
      <div>
        <label class="field-label">Solana Address</label>
        <input v-model="newAddress.address" type="text" placeholder="Enter wallet address..."
          class="input-field input-mono" />
      </div>
      <div>
        <label class="field-label">Label (optional)</label>
        <input v-model="newAddress.label" type="text" placeholder="e.g. Main Wallet"
          class="input-field" />
      </div>
      <div>
        <label class="field-label">Alert Threshold (SOL)</label>
        <input v-model.number="newAddress.threshold" type="number" step="0.1" min="0"
          class="input-field input-mono" />
      </div>
    </Modal>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { apiPost, apiPut, apiDelete, useFetch, useSSE } from '../utils/api.js'
import { truncateAddress, formatDate, formatSOL } from '../utils/format.js'
import Modal from '../components/Modal.vue'
import EmptyState from '../components/EmptyState.vue'

const { data: addresses, refetch } = useFetch('/api/addresses', [])
const showAddModal = ref(false)
const newAddress = ref({ address: '', label: '', threshold: 1 })

onMounted(() => {
  const { disconnect } = useSSE('/api/events', (type, data) => {
    if (type === 'balance_update') {
      refetch()
    }
  })
  onUnmounted(disconnect)
})

const addAddress = async () => {
  try {
    await apiPost('/api/addresses', newAddress.value)
    showAddModal.value = false
    newAddress.value = { address: '', label: '', threshold: 1 }
    addresses.value = await fetch('/api/addresses').then(r => r.json())
  } catch (e) { alert(e.message) }
}

const refreshBalance = async (id) => {
  try {
    await apiPost(`/api/addresses/${id}/refresh`, {})
    addresses.value = await fetch('/api/addresses').then(r => r.json())
  } catch (e) { alert(e.message) }
}

const deleteAddress = async (id) => {
  if (!confirm('Remove this wallet?')) return
  try {
    await apiDelete(`/api/addresses/${id}`)
    addresses.value = await fetch('/api/addresses').then(r => r.json())
  } catch (e) { alert(e.message) }
}

const saveThreshold = async (addr) => {
  try {
    const rules = await fetch('/api/rules').then(r => r.json())
    const existing = rules.find(r => r.address_id === addr.id && r.rule_type === 'balance_change')
    if (existing) {
      await apiPut(`/api/rules/${existing.id}`, { threshold: addr.threshold, enabled: addr.threshold > 0 })
    } else if (addr.threshold > 0) {
      await apiPost('/api/rules', { address_id: addr.id, rule_type: 'balance_change', threshold: addr.threshold })
    }
  } catch (e) { console.error('Failed to save threshold:', e) }
}
</script>

<style scoped>
.wallet-grid {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.wallet-card {
  padding: 20px 24px;
  transition: border-color 0.3s ease, box-shadow 0.3s ease, transform 0.2s ease;
}

.wallet-card:hover {
  transform: translateY(-1px);
}

.wallet-card-top {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 16px;
}

.wallet-identity {
  display: flex;
  align-items: center;
  gap: 14px;
}

.wallet-avatar {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 10px;
  background: linear-gradient(135deg, rgba(6, 182, 212, 0.15), rgba(6, 182, 212, 0.05));
  border: 1px solid rgba(6, 182, 212, 0.2);
  font-weight: 700;
  font-size: 1rem;
  color: var(--cyan-bright);
  flex-shrink: 0;
}

.wallet-label {
  font-size: 0.9375rem;
  font-weight: 600;
}

.wallet-addr {
  font-size: 0.8125rem;
  color: var(--text-muted);
  margin-top: 2px;
}

.wallet-balance {
  text-align: right;
}

.balance-value {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--text-primary);
}

.balance-unit {
  font-size: 0.75rem;
  color: var(--text-muted);
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.wallet-card-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 14px;
  border-top: 1px solid var(--border-subtle);
}

.threshold-group {
  display: flex;
  align-items: center;
  gap: 8px;
}

.threshold-label {
  font-size: 0.75rem;
  color: var(--text-muted);
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.threshold-input {
  width: 80px !important;
  padding: 6px 10px !important;
  font-size: 0.8125rem !important;
}

.threshold-unit {
  font-size: 0.75rem;
  color: var(--text-muted);
  font-family: var(--font-mono);
}

.wallet-actions {
  display: flex;
  gap: 8px;
}

.wallet-updated {
  font-size: 0.6875rem;
  margin-top: 10px;
  font-family: var(--font-mono);
}

</style>
