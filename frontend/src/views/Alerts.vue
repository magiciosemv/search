<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold mb-6">Alert History</h1>

    <div class="card overflow-hidden">
      <table class="w-full">
        <thead class="bg-gray-50 dark:bg-slate-700">
          <tr>
            <th class="px-4 py-3 text-left text-sm font-medium">Time</th>
            <th class="px-4 py-3 text-left text-sm font-medium">Address</th>
            <th class="px-4 py-3 text-left text-sm font-medium">Type</th>
            <th class="px-4 py-3 text-left text-sm font-medium">Old Value</th>
            <th class="px-4 py-3 text-left text-sm font-medium">New Value</th>
            <th class="px-4 py-3 text-left text-sm font-medium">Change</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200 dark:divide-slate-700">
          <tr v-for="alert in alerts" :key="alert.id" class="hover:bg-gray-50 dark:hover:bg-slate-700">
            <td class="px-4 py-3 text-sm">{{ formatDate(alert.sent_at) }}</td>
            <td class="px-4 py-3 font-mono text-sm">{{ truncateAddress(alert.address) }}</td>
            <td class="px-4 py-3">
              <span class="badge badge-type">{{ alert.alert_type }}</span>
            </td>
            <td class="px-4 py-3 text-sm">{{ formatSOL(alert.old_value) }}</td>
            <td class="px-4 py-3 text-sm">{{ formatSOL(alert.new_value) }}</td>
            <td class="px-4 py-3 text-sm">
              <span :class="getChange(alert) >= 0 ? 'text-green-600' : 'text-red-600'">
                {{ getChange(alert) >= 0 ? '+' : '' }}{{ getChange(alert).toFixed(4) }}
              </span>
            </td>
          </tr>
          <tr v-if="alerts.length === 0">
            <td colspan="6" class="px-4 py-8 text-center text-gray-500">
              No alerts yet
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div class="flex justify-center gap-2 mt-4">
      <button @click="loadMore" v-if="alerts.length >= limit"
        class="px-4 py-2 border rounded-lg hover:bg-gray-100 dark:hover:bg-slate-700">
        Load More
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { truncateAddress, formatDate, formatSOL } from '../utils/format.js'

const alerts = ref([])
const limit = ref(50)
const offset = ref(0)

const fetchAlerts = async () => {
  try {
    const res = await fetch(`/api/alerts?limit=${limit.value}&offset=${offset.value}`)
    const newAlerts = await res.json()
    if (offset.value === 0) {
      alerts.value = newAlerts
    } else {
      alerts.value = [...alerts.value, ...newAlerts]
    }
  } catch (e) {
    console.error('Failed to fetch alerts:', e)
  }
}

const loadMore = () => {
  offset.value += limit.value
  fetchAlerts()
}

const getChange = (alert) => {
  return (alert.new_value || 0) - (alert.old_value || 0)
}

onMounted(fetchAlerts)
</script>
