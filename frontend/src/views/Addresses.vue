<template>
  <div class="p-6">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold">Wallets</h1>
      <button @click="showAddModal = true" class="btn-primary">Add Wallet</button>
    </div>

    <div class="grid gap-4">
      <div v-for="addr in addresses" :key="addr.id" class="card p-4">
        <div class="flex items-center justify-between">
          <div>
            <div class="font-mono text-sm">{{ addr.address }}</div>
            <div class="text-sm text-gray-500">{{ addr.label || 'No label' }}</div>
          </div>
          <div class="text-right">
            <div class="text-2xl font-bold">{{ formatSOL(addr.balance) }} <span class="text-sm font-normal">SOL</span></div>
            <div class="text-xs text-gray-500">{{ formatDate(addr.updated_at) }}</div>
          </div>
        </div>
        <div class="mt-4 pt-4 border-t dark:border-slate-700">
          <div class="flex items-center gap-3">
            <label class="text-sm text-gray-500">Alert threshold:</label>
            <input v-model.number="addr.threshold" @change="saveThreshold(addr)" type="number" step="0.1" min="0"
              class="w-24 px-2 py-1.5 border rounded-lg dark:bg-slate-700 dark:border-slate-600 text-sm" />
            <span class="text-sm text-gray-400">SOL</span>
            <div class="ml-auto flex gap-2">
              <button @click="refreshBalance(addr.id)" class="text-purple-600 hover:text-purple-800 text-sm font-medium">
                Refresh
              </button>
              <button @click="deleteAddress(addr.id)" class="text-red-500 hover:text-red-700 text-sm font-medium">
                Delete
              </button>
            </div>
          </div>
        </div>
      </div>
      <EmptyState v-if="addresses.length === 0" message="No wallets added yet" />
    </div>

    <Modal v-model="showAddModal" title="Add Wallet" @submit="addAddress">
      <div>
        <label class="block text-sm font-medium mb-1">Solana Address</label>
        <input v-model="newAddress.address" type="text" placeholder="Enter Solana address"
          class="w-full px-3 py-2 border rounded-lg dark:bg-slate-700 dark:border-slate-600" />
      </div>
      <div>
        <label class="block text-sm font-medium mb-1">Label (optional)</label>
        <input v-model="newAddress.label" type="text" placeholder="e.g., Main Wallet"
          class="w-full px-3 py-2 border rounded-lg dark:bg-slate-700 dark:border-slate-600" />
      </div>
      <div>
        <label class="block text-sm font-medium mb-1">Alert threshold (SOL)</label>
        <input v-model.number="newAddress.threshold" type="number" step="0.1" min="0"
          class="w-full px-3 py-2 border rounded-lg dark:bg-slate-700 dark:border-slate-600" />
      </div>
    </Modal>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { apiPost, apiDelete, useFetch } from '../utils/api.js'
import { formatDate, formatSOL } from '../utils/format.js'
import Modal from '../components/Modal.vue'
import EmptyState from '../components/EmptyState.vue'

const { data: addresses } = useFetch('/api/addresses', [])
const showAddModal = ref(false)
const newAddress = ref({ address: '', label: '', threshold: 1 })

const addAddress = async () => {
  try {
    await apiPost('/api/addresses', newAddress.value)
    showAddModal.value = false
    newAddress.value = { address: '', label: '', threshold: 1 }
    addresses.value = await fetch('/api/addresses').then(r => r.json())
  } catch (e) {
    alert(e.message)
  }
}

const refreshBalance = async (id) => {
  try {
    await apiPost(`/api/addresses/${id}/refresh`, {})
    addresses.value = await fetch('/api/addresses').then(r => r.json())
  } catch (e) {
    alert(e.message)
  }
}

const deleteAddress = async (id) => {
  if (!confirm('Are you sure?')) return
  try {
    await apiDelete(`/api/addresses/${id}`)
    addresses.value = await fetch('/api/addresses').then(r => r.json())
  } catch (e) {
    alert(e.message)
  }
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
  } catch (e) {
    console.error('Failed to save threshold:', e)
  }
}
</script>