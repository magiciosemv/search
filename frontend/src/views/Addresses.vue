<template>
  <div class="p-6">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold">Addresses</h1>
      <button @click="showAddModal = true"
        class="px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700">
        Add Address
      </button>
    </div>

    <!-- Address List -->
    <div class="bg-white dark:bg-slate-800 rounded-lg shadow overflow-hidden">
      <table class="w-full">
        <thead class="bg-gray-50 dark:bg-slate-700">
          <tr>
            <th class="px-4 py-3 text-left text-sm font-medium">Address</th>
            <th class="px-4 py-3 text-left text-sm font-medium">Label</th>
            <th class="px-4 py-3 text-left text-sm font-medium">Balance</th>
            <th class="px-4 py-3 text-left text-sm font-medium">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200 dark:divide-slate-700">
          <tr v-for="addr in addresses" :key="addr.id" class="hover:bg-gray-50 dark:hover:bg-slate-700">
            <td class="px-4 py-3 font-mono text-sm">{{ addr.address }}</td>
            <td class="px-4 py-3">{{ addr.label || '-' }}</td>
            <td class="px-4 py-3 font-semibold">{{ addr.balance?.toFixed(4) || '0.0000' }} SOL</td>
            <td class="px-4 py-3">
              <button @click="refreshBalance(addr.id)" class="text-purple-600 hover:text-purple-800 mr-3">
                Refresh
              </button>
              <button @click="deleteAddress(addr.id)" class="text-red-600 hover:text-red-800">
                Delete
              </button>
            </td>
          </tr>
          <tr v-if="addresses.length === 0">
            <td colspan="4" class="px-4 py-8 text-center text-gray-500">
              No addresses added yet
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Add Address Modal -->
    <div v-if="showAddModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white dark:bg-slate-800 p-6 rounded-lg w-full max-w-md">
        <h2 class="text-xl font-bold mb-4">Add Address</h2>
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
const newAddress = ref({ address: '', label: '' })

const fetchAddresses = async () => {
  try {
    const res = await fetch('/api/addresses')
    addresses.value = await res.json()
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
    newAddress.value = { address: '', label: '' }
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

onMounted(() => {
  fetchAddresses()
})
</script>