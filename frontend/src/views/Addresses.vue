<template>
  <div class="p-6">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold">Wallets</h1>
      <button @click="showAddModal = true"
        class="px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700">
        Add Wallet
      </button>
    </div>

    <!-- Wallet List -->
    <div class="grid gap-4">
      <div v-for="addr in addresses" :key="addr.id"
        class="bg-white dark:bg-slate-800 p-4 rounded-lg shadow">
        <div class="flex items-center justify-between">
          <div>
            <div class="font-mono text-sm">{{ addr.address }}</div>
            <div class="text-sm text-gray-500">{{ addr.label || 'No label' }}</div>
          </div>
          <div class="text-right">
            <div class="text-2xl font-bold">{{ addr.balance?.toFixed(4) || '0.0000' }} <span class="text-sm font-normal">SOL</span></div>
            <div class="text-xs text-gray-500">{{ addr.updated_at ? new Date(addr.updated_at).toLocaleString() : '' }}</div>
          </div>
        </div>

        <!-- Monitor Settings -->
        <div class="mt-4 pt-4 border-t dark:border-slate-700">
          <div class="flex items-center gap-4">
            <label class="text-sm">Alert when balance changes by:</label>
            <input v-model.number="addr.threshold" @change="saveThreshold(addr)" type="number" step="0.1" min="0"
              class="w-24 px-2 py-1 border rounded dark:bg-slate-700" placeholder="SOL" />
            <span class="text-sm text-gray-500">SOL</span>
            <button @click="refreshBalance(addr.id)" class="text-purple-600 hover:text-purple-800 text-sm ml-auto">
              Refresh
            </button>
            <button @click="deleteAddress(addr.id)" class="text-red-600 hover:text-red-800 text-sm">
              Delete
            </button>
          </div>
        </div>
      </div>

      <div v-if="addresses.length === 0" class="bg-white dark:bg-slate-800 p-8 rounded-lg shadow text-center text-gray-500">
        No wallets added yet
      </div>
    </div>

    <!-- Add Wallet Modal -->
    <div v-if="showAddModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white dark:bg-slate-800 p-6 rounded-lg w-full max-w-md">
        <h2 class="text-xl font-bold mb-4">Add Wallet</h2>
        <div class="space-y-4">
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
            <label class="block text-sm font-medium mb-1">Alert when balance changes by (SOL)</label>
            <input v-model.number="newAddress.threshold" type="number" step="0.1" min="0" value="1"
              class="w-full px-3 py-2 border rounded-lg dark:bg-slate-700 dark:border-slate-600" />
          </div>
        </div>
        <div class="flex justify-end gap-2 mt-6">
          <button @click="showAddModal = false" class="px-4 py-2 border rounded-lg hover:bg-gray-100 dark:hover:bg-slate-700">
            Cancel
          </button>
          <button @click="addAddress" class="px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700">
            Add
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const addresses = ref([])
const showAddModal = ref(false)
const newAddress = ref({ address: '', label: '', threshold: 1 })

const fetchAddresses = async () => {
  try {
    const res = await fetch('/api/addresses')
    addresses.value = await res.json()
    // Initialize threshold if not set
    addresses.value.forEach(addr => {
      if (addr.threshold === undefined) addr.threshold = 1
    })
  } catch (e) {
    console.error('Failed to fetch addresses:', e)
  }
}

const addAddress = async () => {
  try {
    const res = await fetch('/api/addresses', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(newAddress.value)
    })
    if (!res.ok) {
      const err = await res.json()
      alert(err.error || 'Failed to add address')
      return
    }
    showAddModal.value = false
    newAddress.value = { address: '', label: '', threshold: 1 }
    fetchAddresses()
  } catch (e) {
    alert('Failed to add address')
  }
}

const refreshBalance = async (id) => {
  try {
    await fetch(`/api/addresses/${id}/refresh`, { method: 'POST' })
    fetchAddresses()
  } catch (e) {
    alert('Failed to refresh balance')
  }
}

const deleteAddress = async (id) => {
  if (!confirm('Are you sure?')) return
  try {
    await fetch(`/api/addresses/${id}`, { method: 'DELETE' })
    fetchAddresses()
  } catch (e) {
    alert('Failed to delete address')
  }
}

const saveThreshold = async (addr) => {
  try {
    // Get or create rule for this address
    const rulesRes = await fetch('/api/rules')
    const rules = await rulesRes.json()
    const existingRule = rules.find(r => r.address_id === addr.id && r.rule_type === 'balance_change')

    if (existingRule) {
      await fetch(`/api/rules/${existingRule.id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ threshold: addr.threshold, enabled: addr.threshold > 0 })
      })
    } else if (addr.threshold > 0) {
      await fetch('/api/rules', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ address_id: addr.id, rule_type: 'balance_change', threshold: addr.threshold })
      })
    }
  } catch (e) {
    console.error('Failed to save threshold:', e)
  }
}

onMounted(() => {
  fetchAddresses()
})
</script>