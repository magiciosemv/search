<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold mb-6">Dashboard</h1>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
      <div class="bg-white dark:bg-slate-800 p-4 rounded-lg shadow">
        <div class="text-gray-500 dark:text-gray-400 text-sm">Total Wallets</div>
        <div class="text-2xl font-bold">{{ stats.total_addresses || 0 }}</div>
      </div>
      <div class="bg-white dark:bg-slate-800 p-4 rounded-lg shadow">
        <div class="text-gray-500 dark:text-gray-400 text-sm">Notifications</div>
        <div class="text-2xl font-bold">{{ stats.total_notifications || 0 }}</div>
      </div>
      <div class="bg-white dark:bg-slate-800 p-4 rounded-lg shadow">
        <div class="text-gray-500 dark:text-gray-400 text-sm">Total Alerts</div>
        <div class="text-2xl font-bold">{{ stats.total_alerts || 0 }}</div>
      </div>
      <div class="bg-white dark:bg-slate-800 p-4 rounded-lg shadow">
        <div class="text-gray-500 dark:text-gray-400 text-sm">Today's Alerts</div>
        <div class="text-2xl font-bold">{{ stats.today_alerts || 0 }}</div>
      </div>
    </div>

    <!-- Monitor Status -->
    <div class="bg-white dark:bg-slate-800 p-4 rounded-lg shadow mb-8">
      <div class="flex items-center justify-between">
        <div>
          <h2 class="text-lg font-semibold">Monitor Service</h2>
          <p class="text-gray-500 dark:text-gray-400 text-sm">Background balance monitoring</p>
        </div>
        <span :class="stats.monitor_running ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'"
          class="px-3 py-1 rounded-full text-sm font-medium">
          {{ stats.monitor_running ? 'Running' : 'Stopped' }}
        </span>
      </div>
    </div>

    <!-- Recent Alerts -->
    <div class="bg-white dark:bg-slate-800 p-4 rounded-lg shadow">
      <h2 class="text-lg font-semibold mb-4">Recent Alerts</h2>
      <div v-if="!recentAlerts || recentAlerts.length === 0" class="text-gray-500 text-center py-8">
        No alerts yet
      </div>
      <div v-else class="space-y-3">
        <div v-for="alert in recentAlerts" :key="alert.id"
          class="flex items-center justify-between p-3 bg-gray-50 dark:bg-slate-700 rounded">
          <div>
            <div class="font-mono text-sm">{{ truncateAddress(alert.address || alert.addressStr) }}</div>
            <div class="text-xs text-gray-500">{{ formatDate(alert.sent_at) }}</div>
          </div>
          <div class="text-right">
            <div class="font-semibold">{{ alert.alert_type }}</div>
            <div class="text-xs text-gray-500">{{ (alert.old_value || 0).toFixed(4) }} → {{ (alert.new_value || 0).toFixed(4) }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const stats = ref({})
const recentAlerts = ref([])

const fetchStats = async () => {
  try {
    const res = await fetch('/api/stats')
    stats.value = await res.json()
  } catch (e) {
    console.error('Failed to fetch stats:', e)
  }
}

const fetchAlerts = async () => {
  try {
    const res = await fetch('/api/alerts?limit=5')
    recentAlerts.value = await res.json()
  } catch (e) {
    console.error('Failed to fetch alerts:', e)
  }
}

const truncateAddress = (addr) => {
  if (!addr) return ''
  return addr.slice(0, 6) + '...' + addr.slice(-4)
}

const formatDate = (date) => {
  return new Date(date).toLocaleString()
}

onMounted(() => {
  fetchStats()
  fetchAlerts()
})
</script>