<template>
  <div class="page-container">
    <div class="page-header">
      <div>
        <h1 class="page-title">{{ t('rules.title') }}</h1>
        <p class="page-subtitle">{{ t('rules.subtitle') }}</p>
      </div>
      <button class="btn btn-primary" @click="showAddModal = true">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
        {{ t('rules.addRule') }}
      </button>
    </div>

    <div class="rules-list stagger">
      <div v-for="rule in rules" :key="rule.id" class="rule-card glass animate-fade-up">
        <div class="rule-top">
          <div class="rule-identity">
            <div class="rule-icon">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"><path d="M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>
            </div>
            <div>
              <div class="rule-type badge badge-blue">{{ rule.rule_type === 'balance_change' ? t('rules.balanceChange') : rule.rule_type }}</div>
              <div class="rule-wallet text-muted font-mono" v-if="getAddress(rule.address_id)">{{ truncateAddress(getAddress(rule.address_id).address) }} <span class="chain-tag">{{ getChainSymbol(getAddress(rule.address_id).chain) }}</span></div>
            </div>
          </div>
          <div class="rule-right">
            <div class="rule-threshold font-mono">{{ rule.threshold }} <span class="text-dim">SOL</span></div>
            <button class="toggle-btn" :class="{ on: rule.enabled }" @click="toggleRule(rule)">
              <div class="toggle-track"><div class="toggle-thumb"></div></div>
            </button>
          </div>
        </div>
        <div class="rule-actions">
          <button @click="editRule(rule)" class="btn btn-ghost btn-sm">{{ t('common.edit') }}</button>
          <button @click="deleteRule(rule.id)" class="btn btn-danger btn-sm">{{ t('common.delete') }}</button>
        </div>
      </div>
    </div>

    <EmptyState v-if="!rules || rules.length === 0" :title="t('rules.noRules')" :message="t('rules.noRulesMsg')" />

    <Modal v-model="showAddModal" :title="editingRule ? t('common.edit') : t('rules.addRule')" :submit-label="editingRule ? t('common.save') : t('common.add')" @submit="submitRule">
      <div>
        <label class="field-label">{{ t('rules.wallet') }}</label>
        <select v-model="form.address_id" class="input-field">
          <option value="">Select wallet...</option>
          <option v-for="addr in addresses" :key="addr.id" :value="addr.id">{{ addr.label || truncateAddress(addr.address) }}</option>
        </select>
      </div>
      <div>
        <label class="field-label">{{ t('rules.ruleType') }}</label>
        <select v-model="form.rule_type" class="input-field">
          <option value="balance_change">{{ t('rules.balanceChange') }}</option>
        </select>
      </div>
      <div>
        <label class="field-label">{{ t('rules.threshold') }}</label>
        <input v-model.number="form.threshold" type="number" step="0.1" min="0" class="input-field input-mono" />
      </div>
    </Modal>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { apiPost, apiPut, apiDelete, useFetch } from '../utils/api.js'
import { useToast } from '../utils/toast.js'
import { useConfirm } from '../utils/confirm.js'
import { truncateAddress } from '../utils/format.js'
import { useI18n } from '../utils/i18n.js'
import Modal from '../components/Modal.vue'
import EmptyState from '../components/EmptyState.vue'

const { t } = useI18n()
const toast = useToast()
const confirm = useConfirm()

const { data: rules, refetch: refetchRules } = useFetch('/api/rules', [])
const { data: addresses } = useFetch('/api/addresses', [])

const showAddModal = ref(false)
const editingRule = ref(null)
const form = ref({ address_id: '', rule_type: 'balance_change', threshold: 1 })

const getAddress = (id) => addresses.value.find(a => a.id === id)

const chainSymbols = {
  solana: 'SOL', ethereum: 'ETH', bitcoin: 'BTC',
  usdt_erc20: 'USDT', usdc_erc20: 'USDC', bsc: 'BNB', polygon: 'MATIC', arbitrum: 'ETH'
}
const getChainSymbol = (chain) => chainSymbols[chain] || chain?.toUpperCase() || 'SOL'

const editRule = (rule) => {
  editingRule.value = rule
  form.value = { address_id: rule.address_id, rule_type: rule.rule_type, threshold: rule.threshold }
  showAddModal.value = true
}

const submitRule = async () => {
  try {
    if (editingRule.value) {
      await apiPut(`/api/rules/${editingRule.value.id}`, { rule_type: form.value.rule_type, threshold: form.value.threshold, enabled: editingRule.value.enabled })
      toast.success(t('rules.ruleUpdated'))
    } else {
      await apiPost('/api/rules', { address_id: Number(form.value.address_id), rule_type: form.value.rule_type, threshold: form.value.threshold })
      toast.success(t('rules.ruleAdded'))
    }
    showAddModal.value = false
    editingRule.value = null
    form.value = { address_id: '', rule_type: 'balance_change', threshold: 1 }
    refetchRules()
  } catch (e) { toast.error(e.message) }
}

const toggleRule = async (rule) => {
  try {
    await apiPut(`/api/rules/${rule.id}`, { rule_type: rule.rule_type, threshold: rule.threshold, enabled: !rule.enabled })
    refetchRules()
  } catch (e) { toast.error(e.message) }
}

const deleteRule = async (id) => {
  if (!await confirm.show(t('rules.removeRule'), t('rules.deleteRule'))) return
  try { await apiDelete(`/api/rules/${id}`); toast.success(t('rules.ruleRemoved')); refetchRules() } catch (e) { toast.error(e.message) }
}
</script>

<style scoped>
.rules-list { display: flex; flex-direction: column; gap: 10px; }
.rule-card { padding: 18px 20px; }
.rule-top { display: flex; align-items: center; justify-content: space-between; margin-bottom: 14px; }
.rule-identity { display: flex; align-items: center; gap: 14px; }
.rule-icon {
  width: 38px; height: 38px; display: flex; align-items: center; justify-content: center;
  border-radius: 10px; background: var(--amber-dim); color: var(--amber-bright);
}
.rule-wallet { font-size: 0.75rem; margin-top: 4px; }
.rule-right { display: flex; align-items: center; gap: 14px; }
.rule-threshold { font-size: 1.125rem; font-weight: 700; }
.rule-actions { display: flex; gap: 6px; padding-top: 14px; border-top: 1px solid var(--border-subtle); }
.chain-tag { color: var(--blue-bright); font-size: 0.625rem; margin-left: 4px; }

/* Toggle */
.toggle-btn { display: flex; align-items: center; background: none; border: none; cursor: pointer; }
.toggle-track {
  width: 38px; height: 22px; border-radius: 11px;
  background: var(--bg-elevated); border: 1px solid var(--border-subtle);
  position: relative; transition: all var(--duration-normal) ease;
}
.toggle-thumb {
  width: 16px; height: 16px; border-radius: 50%;
  background: var(--text-dim); position: absolute; top: 2px; left: 2px;
  transition: all var(--duration-normal) cubic-bezier(0.16, 1, 0.3, 1);
}
.toggle-btn.on .toggle-track { background: rgba(0, 214, 143, 0.12); border-color: rgba(0, 214, 143, 0.3); }
.toggle-btn.on .toggle-thumb { left: 18px; background: var(--green-base); box-shadow: 0 0 8px rgba(0, 214, 143, 0.4); }
</style>
