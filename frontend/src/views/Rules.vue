<template>
  <div class="p-6">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold">Rules</h1>
      <button @click="showAddModal = true"
        class="px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700">
        Add Rule
      </button>
    </div>

    <!-- Rules List -->
    <div class="bg-white dark:bg-slate-800 rounded-lg shadow overflow-hidden">
      <table class="w-full">
        <thead class="bg-gray-50 dark:bg-slate-700">
          <tr>
            <th class="px-4 py-3 text-left text-sm font-medium">Address</th>
            <th class="px-4 py-3 text-left text-sm font-medium">Rule Type</th>
            <th class="px-4 py-3 text-left text-sm font-medium">Threshold</th>
            <th class="px-4 py-3 text-left text-sm font-medium">Status</th>
            <th class="px-4 py-3 text-left text-sm font-medium">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200 dark:divide-slate-700">
          <tr v-for="rule in rules" :key="rule.id" class="hover:bg-gray-50 dark:hover:bg-slate-700">
            <td class="px-4 py-3 font-mono text-sm">{{ getAddressLabel(rule.address_id) }}</td>
            <td class="px-4 py-3">{{ formatRuleType(rule.rule_type) }}</td>
            <td class="px-4 py-3">{{ rule.threshold }}</td>
            <td class="px-4 py-3">
              <span :class="rule.enabled ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'"
                class="px-2 py-1 rounded text-xs">
                {{ rule.enabled ? 'Active' : 'Disabled' }}
              </span>
            </td>
            <td class="px-4 py-3">
              <button @click="toggleRule(rule)" class="text-purple-600 hover:text-purple-800 mr-3">
                {{ rule.enabled ? 'Disable' : 'Enable' }}
              </button>
              <button @click="deleteRule(rule.id)" class="text-red-600 hover:text-red-800">
                Delete
              </button>
            </td>
          </tr>
          <tr v-if="rules.length === 0">
            <td colspan="5" class="px-4 py-8 text-center text-gray-500">
              No rules defined yet
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Add Rule Modal -->
    <div v-if="showAddModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white dark:bg-slate-800 p-6 rounded-lg w-full max-w-md">
        <h2 class="text-xl font-bold mb-4">Add Rule</h2>
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium mb-1">Address</label>
            <select v-model="newRule.address_id"
              class="w-full px-3 py-2 border rounded-lg dark:bg-slate-700 dark:border-slate-600">
              <option value="">Select address</option>
              <option v-for="addr in addresses" :key="addr.id" :value="addr.id">
                {{ addr.address.slice(0, 8) }}...{{ addr.address.slice(-4) }} - {{ addr.label || 'No label' }}
              </option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium mb-1">Rule Type</label>
            <select v-model="newRule.rule_type"
              class="w-full px-3 py-2 border rounded-lg dark:bg-slate-700 dark:border-slate-600">
              <option value="balance_change">Balance Change (SOL)</option>
              <option value="balance_change_percent">Balance Change (%)</option>
              <option value="threshold_above">Threshold Above</option>
              <option value="threshold_below">Threshold Below</option>
              <option value="large_incoming">Large Incoming</option>
              <option value="large_outgoing">Large Outgoing</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium mb-1">Threshold</label>
            <input v-model.number="newRule.threshold" type="number" step="0.1" placeholder="e.g., 1.0"
              class="w-full px-3 py-2 border rounded-lg dark:bg-slate-700 dark:border-slate-600" />
          </div>
        </div>
        <div class="flex justify-end gap-2 mt-6">
          <button @click="showAddModal = false" class="px-4 py-2 border rounded-lg hover:bg-gray-100 dark:hover:bg-slate-700">
            Cancel
          </button>
          <button @click="addRule" class="px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700">
            Add
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const rules = ref([])
const addresses = ref([])
const showAddModal = ref(false)
const newRule = ref({ address_id: '', rule_type: 'balance_change', threshold: 1 })

const fetchRules = async () => {
  try {
    const res = await fetch('/api/rules')
    rules.value = await res.json()
  } catch (e) {
    console.error('Failed to fetch rules:', e)
  }
}

const fetchAddresses = async () => {
  try {
    const res = await fetch('/api/addresses')
    addresses.value = await res.json()
  } catch (e) {
    console.error('Failed to fetch addresses:', e)
  }
}

const getAddressLabel = (id) => {
  const addr = addresses.value.find(a => a.id === id)
  if (!addr) return 'Unknown'
  return addr.address.slice(0, 8) + '...' + addr.address.slice(-4)
}

const formatRuleType = (type) => {
  const types = {
    'balance_change': 'Balance Change (SOL)',
    'balance_change_percent': 'Balance Change (%)',
    'threshold_above': 'Above Threshold',
    'threshold_below': 'Below Threshold',
    'large_incoming': 'Large Incoming',
    'large_outgoing': 'Large Outgoing'
  }
  return types[type] || type
}

const addRule = async () => {
  try {
    const res = await fetch('/api/rules', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(newRule.value)
    })
    if (!res.ok) {
      const err = await res.json()
      alert(err.error || 'Failed to add rule')
      return
    }
    showAddModal.value = false
    newRule.value = { address_id: '', rule_type: 'balance_change', threshold: 1 }
    fetchRules()
  } catch (e) {
    alert('Failed to add rule')
  }
}

const toggleRule = async (rule) => {
  try {
    await fetch(`/api/rules/${rule.id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ ...rule, enabled: !rule.enabled })
    })
    fetchRules()
  } catch (e) {
    alert('Failed to update rule')
  }
}

const deleteRule = async (id) => {
  if (!confirm('Are you sure?')) return
  try {
    await fetch(`/api/rules/${id}`, { method: 'DELETE' })
    fetchRules()
  } catch (e) {
    alert('Failed to delete rule')
  }
}

onMounted(() => {
  fetchRules()
  fetchAddresses()
})
</script>