<template>
  <div class="page-container">
    <div class="page-header">
      <div>
        <h1 class="page-title">{{ t('wallets.title') }}</h1>
        <p class="page-subtitle">{{ t('wallets.subtitle') }}</p>
      </div>
      <button class="btn btn-primary" @click="showAddModal = true">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
        {{ t('wallets.addWallet') }}
      </button>
    </div>

    <div class="wallet-list stagger">
      <div v-for="addr in addresses" :key="addr.id" class="wallet-card glass animate-fade-up">
        <div class="wallet-header">
          <div class="wallet-identity">
            <div class="wallet-avatar">{{ (addr.label || 'W')[0].toUpperCase() }}</div>
            <div>
              <div class="wallet-label">{{ addr.label || 'Unnamed' }} <span class="chain-badge font-mono">{{ getChainSymbol(addr.chain) }}</span></div>
              <div class="wallet-addr font-mono" :title="addr.address">{{ truncateAddress(addr.address) }}</div>
            </div>
          </div>
          <div class="wallet-balance">
            <span class="balance-value font-mono">{{ formatSOL(addr.balance) }}</span>
            <span class="balance-unit text-dim">{{ getChainSymbol(addr.chain) }}</span>
          </div>
        </div>
        <div class="wallet-footer">
          <div class="threshold-group">
            <span class="threshold-label text-muted">{{ t('wallets.threshold') }}</span>
            <input v-model.number="addr.threshold" @change="saveThreshold(addr)"
              type="number" step="0.1" min="0" class="input-field input-mono threshold-input" />
            <span class="threshold-unit text-dim font-mono">SOL</span>
          </div>
          <div class="wallet-actions">
            <button @click="refreshBalance(addr.id)" class="btn btn-ghost btn-sm">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><polyline points="23 4 23 10 17 10"/><path d="M20.49 15a9 9 0 11-2.12-9.36L23 10"/></svg>
              {{ t('wallets.refresh') }}
            </button>
            <button @click="deleteAddress(addr.id)" class="btn btn-danger btn-sm">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2"/></svg>
            </button>
          </div>
        </div>
        <div class="wallet-updated text-dim font-mono" v-if="addr.updated_at">{{ t('wallets.updated') }} {{ formatDate(addr.updated_at) }}</div>
      </div>
    </div>

    <EmptyState v-if="addresses.length === 0" :title="t('wallets.noWallets')" :message="t('wallets.noWalletsMsg')" />

    <Modal v-model="showAddModal" :title="t('wallets.addWallet')" :submit-label="t('common.add')" @submit="addAddress">
      <div>
        <label class="field-label">{{ t('wallets.chain') }}</label>
        <select v-model="newAddress.chain" class="input-field">
          <option value="solana">Solana (SOL)</option>
          <option value="ethereum">Ethereum (ETH)</option>
          <option value="bitcoin">Bitcoin (BTC)</option>
          <option value="usdt_erc20">USDT (ERC-20)</option>
          <option value="usdc_erc20">USDC (ERC-20)</option>
          <option value="bsc">BNB Chain (BNB)</option>
          <option value="polygon">Polygon (MATIC)</option>
          <option value="arbitrum">Arbitrum (ETH)</option>
        </select>
      </div>
      <div>
        <label class="field-label">{{ addressLabel }}</label>
        <input v-model="newAddress.address" type="text" :placeholder="addressPlaceholder" class="input-field input-mono" />
      </div>
      <div>
        <label class="field-label">{{ t('wallets.label') }}</label>
        <input v-model="newAddress.label" type="text" placeholder="e.g. Main Wallet" class="input-field" />
      </div>
      <div>
        <label class="field-label">{{ t('wallets.alertThreshold') }}</label>
        <input v-model.number="newAddress.threshold" type="number" step="0.1" min="0" class="input-field input-mono" />
      </div>
    </Modal>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { apiPost, apiDelete, useFetch, useSSE } from '../utils/api.js'
import { useToast } from '../utils/toast.js'
import { useConfirm } from '../utils/confirm.js'
import { useI18n } from '../utils/i18n.js'
import { truncateAddress, formatDate, formatSOL } from '../utils/format.js'
import Modal from '../components/Modal.vue'
import EmptyState from '../components/EmptyState.vue'

const { t } = useI18n()
const toast = useToast()
const confirm = useConfirm()
const { data: addresses, refetch } = useFetch('/api/addresses', [])
const showAddModal = ref(false)
const newAddress = ref({ address: '', chain: 'solana', label: '', threshold: 1 })

const chainSymbols = {
  solana: 'SOL', ethereum: 'ETH', bitcoin: 'BTC',
  usdt_erc20: 'USDT', usdc_erc20: 'USDC', usdt_trc20: 'USDT',
  bsc: 'BNB', polygon: 'MATIC', arbitrum: 'ETH'
}

const getChainSymbol = (chain) => chainSymbols[chain] || chain?.toUpperCase() || 'SOL'

const chainPlaceholders = {
  solana: { label: 'Solana Address', placeholder: 'e.g. 7xKXtg2CW87d97TXJSDpbD5jBkheTqA83TZRuJosgAsU' },
  ethereum: { label: 'Ethereum Address', placeholder: 'e.g. 0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045' },
  bitcoin: { label: 'Bitcoin Address', placeholder: 'e.g. bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh' },
  usdt_erc20: { label: 'USDT Wallet Address (ERC-20)', placeholder: 'e.g. 0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045' },
  usdc_erc20: { label: 'USDC Wallet Address (ERC-20)', placeholder: 'e.g. 0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045' },
  usdt_trc20: { label: 'USDT Wallet Address (TRC-20)', placeholder: 'e.g. TNPeeaaFBmJKcy9S3Qsf48mFyv7aaYfBee' },
  bsc: { label: 'BNB Chain Address', placeholder: 'e.g. 0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045' },
  polygon: { label: 'Polygon Address', placeholder: 'e.g. 0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045' },
  arbitrum: { label: 'Arbitrum Address', placeholder: 'e.g. 0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045' },
}

const addressLabel = computed(() => chainPlaceholders[newAddress.value.chain]?.label || 'Wallet Address')
const addressPlaceholder = computed(() => chainPlaceholders[newAddress.value.chain]?.placeholder || 'Enter wallet address...')

onMounted(() => {
  const { disconnect } = useSSE('/api/events', (type) => {
    if (type === 'balance_update') refetch()
  })
  onUnmounted(disconnect)
})

const addAddress = async () => {
  try {
    await apiPost('/api/addresses', newAddress.value)
    showAddModal.value = false
    newAddress.value = { address: '', chain: 'solana', label: '', threshold: 1 }
    addresses.value = await fetch('/api/addresses', { headers: { Authorization: 'Bearer solana-monitor-secret-key-2024' } }).then(r => r.json())
    toast.success(t('wallets.walletAdded'))
  } catch (e) { toast.error(e.message) }
}

const refreshBalance = async (id) => {
  try { await apiPost(`/api/addresses/${id}/refresh`, {}); refetch(); toast.success(t('wallets.balanceRefreshed')) } catch (e) { toast.error(e.message) }
}

const deleteAddress = async (id) => {
  if (!await confirm.show(t('wallets.removeWallet'), t('wallets.deleteWallet'))) return
  try { await apiDelete(`/api/addresses/${id}`); refetch(); toast.success(t('wallets.walletRemoved')) } catch (e) { toast.error(e.message) }
}

const saveThreshold = async (addr) => {
  try {
    const rules = await fetch('/api/rules', { headers: { Authorization: 'Bearer solana-monitor-secret-key-2024' } }).then(r => r.json())
    const existing = rules.find(r => r.address_id === addr.id && r.rule_type === 'balance_change')
    if (existing) {
      await apiPut(`/api/rules/${existing.id}`, { threshold: addr.threshold, enabled: addr.threshold > 0 })
    } else if (addr.threshold > 0) {
      await apiPost('/api/rules', { address_id: addr.id, rule_type: 'balance_change', threshold: addr.threshold })
    }
  } catch (e) { toast.error(t('wallets.saveThresholdFailed')) }
}
</script>

<style scoped>
.wallet-list { display: flex; flex-direction: column; gap: 10px; }
.wallet-card { padding: 18px 20px; }
.wallet-card:active { transform: scale(0.998); }

.wallet-header { display: flex; align-items: flex-start; justify-content: space-between; margin-bottom: 16px; }
.wallet-identity { display: flex; align-items: center; gap: 14px; }
.wallet-avatar {
  width: 40px; height: 40px; display: flex; align-items: center; justify-content: center;
  border-radius: 10px; font-weight: 700; font-size: 0.9375rem;
  background: var(--blue-dim); border: 1px solid rgba(59, 130, 246, 0.18);
  color: var(--blue-bright);
}
.wallet-label { font-size: 0.9375rem; font-weight: 600; }
.wallet-addr { font-size: 0.75rem; color: var(--text-dim); margin-top: 2px; }
.wallet-balance { text-align: right; }
.balance-value { font-size: 1.375rem; font-weight: 700; }
.balance-unit { font-size: 0.6875rem; margin-left: 3px; }

.wallet-footer {
  display: flex; align-items: center; justify-content: space-between;
  padding-top: 14px; border-top: 1px solid var(--border-subtle);
}
.threshold-group { display: flex; align-items: center; gap: 7px; }
.threshold-label { font-size: 0.6875rem; font-weight: 600; }
.threshold-input { width: 72px !important; padding: 6px 10px !important; font-size: 0.75rem !important; }
.threshold-unit { font-size: 0.6875rem; }
.wallet-actions { display: flex; gap: 6px; }
.wallet-updated { font-size: 0.625rem; margin-top: 10px; }
.chain-badge {
  font-size: 0.5625rem;
  padding: 1px 6px;
  border-radius: 4px;
  background: var(--blue-dim);
  color: var(--blue-bright);
  margin-left: 6px;
  font-weight: 600;
  letter-spacing: 0.04em;
}
</style>
