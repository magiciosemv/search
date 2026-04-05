<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold mb-6">Dashboard</h1>

    <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
      <div class="card p-4">
        <div class="text-gray-500 dark:text-gray-400 text-sm">Total Wallets</div>
        <div class="text-2xl font-bold">{{ stats.total_addresses || 0 }}</div>
      </div>
      <div class="card p-4">
        <div class="text-gray-500 dark:text-gray-400 text-sm">Notifications</div>
        <div class="text-2xl font-bold">{{ stats.total_notifications || 0 }}</div>
      </div>
      <div class="card p-4">
        <div class="text-gray-500 dark:text-gray-400 text-sm">Total Alerts</div>
        <div class="text-2xl font-bold">{{ stats.total_alerts || 0 }}</div>
      </div>
      <div class="card p-4">
        <div class="text-gray-500 dark:text-gray-400 text-sm">Today's Alerts</div>
        <div class="text-2xl font-bold">{{ stats.today_alerts || 0 }}</div>
      </div>
    </div>

    <div class="card p-4 mb-8">
      <div class="flex items-center justify-between">
        <div>
          <h2 class="text-lg font-semibold">Monitor Service</h2>
          <p class="text-gray-500 dark:text-gray-400 text-sm">Background balance monitoring</p>
        </div>
        <span :class="stats.monitor_running ? 'badge badge-success' : 'badge badge-muted'">
          {{ stats.monitor_running ? 'Running' : 'Stopped' }}
        </span>
      </div>
    </div>

    <div class="card p-4">
      <h2 class="text-lg font-semibold mb-4">Recent Alerts</h2>
      <EmptyState v-if="!recentAlerts || recentAlerts.length === 0" title="" message="No alerts yet" />
      <div v-else class="space-y-3">
        <div v-for="alert in recentAlerts" :key="alert.id"
          class="flex items-center justify-between p-3 bg-gray-50 dark:bg-slate-700 rounded">
          <div>
            <div class="font-mono text-sm">{{ truncateAddress(alert.address || alert.addressStr) }}</div>
            <div class="text-xs text-gray-500">{{ formatDate(alert.sent_at) }}</div>
          </div>
          <div class="text-right">
            <div class="font-semibold">{{ alert.alert_type }}</div>
            <div class="text-xs text-gray-500">{{ formatSOL(alert.old_value) }} &rarr; {{ formatSOL(alert.new_value) }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useFetch } from '../utils/api.js'
import { truncateAddress, formatDate, formatSOL } from '../utils/format.js'

const { data: stats } = useFetch('/api/stats', {})
const { data: recentAlerts } = useFetch('/api/alerts?limit=5', [])
</script>