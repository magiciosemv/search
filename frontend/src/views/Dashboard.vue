<template>
  <div class="page-container">
    <div class="stats-grid stagger">
      <div v-for="stat in statCards" :key="stat.label" class="stat-card glass animate-fade-up">
        <div class="stat-top">
          <span class="stat-label">{{ stat.label }}</span>
          <div class="stat-icon" :class="stat.iconClass" v-html="stat.icon"></div>
        </div>
        <div class="stat-value font-mono">{{ stats[stat.key] ?? 0 }}</div>
        <div class="stat-footer text-dim font-mono">{{ stat.sub }}</div>
      </div>
    </div>

    <div class="glass monitor-bar animate-fade-up" style="animation-delay: 180ms">
      <div class="monitor-left">
        <span class="status-dot"></span>
        <span class="monitor-label">Monitor Service</span>
        <span class="monitor-interval font-mono">every 30s</span>
      </div>
      <span class="badge badge-success">{{ monitorRunning ? 'Running' : 'Stopped' }}</span>
      <button class="btn btn-ghost btn-sm" @click="backupDB">
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
        Backup
      </button>
    </div>

    <div class="glass alerts-section animate-fade-up" style="animation-delay: 260ms">
      <div class="alerts-header">
        <h2 class="section-title">Recent Alerts</h2>
        <router-link to="/alerts" class="view-link">View all</router-link>
      </div>
      <EmptyState v-if="!recentAlerts || recentAlerts.length === 0" message="No alerts yet" />
      <div v-else class="alert-list">
        <div v-for="alert in recentAlerts" :key="alert.id" class="alert-row">
          <span class="alert-addr font-mono text-blue">{{ truncateAddress(alert.address || alert.addressStr) }}</span>
          <span class="badge badge-blue">{{ alert.alert_type }}</span>
          <span class="alert-diff font-mono">
            <span class="text-muted">{{ formatSOL(alert.old_value) }}</span>
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" style="color:var(--text-dim)"><polyline points="9 18 15 12 9 6"/></svg>
            <span>{{ formatSOL(alert.new_value) }}</span>
          </span>
          <span class="alert-time text-muted font-mono">{{ formatDate(alert.sent_at) }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, onUnmounted } from 'vue'
import { useFetch, useSSE } from '../utils/api.js'
import { truncateAddress, formatDate, formatSOL } from '../utils/format.js'
import { useToast } from '../utils/toast.js'
import EmptyState from '../components/EmptyState.vue'

const toast = useToast()

const { data: stats, refetch: refetchStats } = useFetch('/api/stats', {})
const { data: recentAlerts, refetch: refetchAlerts } = useFetch('/api/alerts?limit=5', [])

const monitorRunning = true

const backupDB = async () => {
  try {
    const res = await fetch('/api/backup', { headers: { Authorization: 'Bearer solana-monitor-secret-key-2024' } })
    if (!res.ok) throw new Error('Backup failed')
    const blob = await res.blob()
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `solana-monitor-backup-${new Date().toISOString().slice(0, 10)}.db`
    a.click()
    URL.revokeObjectURL(url)
    toast.success('Database backup downloaded')
  } catch (e) { toast.error(e.message) }
}

const statCards = [
  {
    label: 'Wallets', key: 'total_addresses', iconClass: 'icon-blue', sub: 'monitored',
    icon: '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 12V7H5a2 2 0 010-4h14v4"/><path d="M3 5v14a2 2 0 002 2h16v-5"/><path d="M18 12a2 2 0 000 4h4v-4z"/></svg>'
  },
  {
    label: 'Channels', key: 'total_notifications', iconClass: 'icon-purple', sub: 'active',
    icon: '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M18 8A6 6 0 006 8c0 7-3 9-3 9h18s-3-2-3-9"/><path d="M13.73 21a2 2 0 01-3.46 0"/></svg>'
  },
  {
    label: 'Total Alerts', key: 'total_alerts', iconClass: 'icon-amber', sub: 'all time',
    icon: '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>'
  },
  {
    label: 'Today', key: 'today_alerts', iconClass: 'icon-green', sub: '24h',
    icon: '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>'
  }
]

onMounted(() => {
  const { disconnect } = useSSE('/api/events', (type) => {
    if (type === 'balance_update') refetchStats()
    if (type === 'new_alert') { refetchAlerts(); refetchStats() }
  })
  onUnmounted(disconnect)
})
</script>

<style scoped>
.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 14px;
  margin-bottom: 14px;
}
@media (max-width: 900px) { .stats-grid { grid-template-columns: repeat(2, 1fr); } }

.stat-card { padding: 18px 20px; }
.stat-top { display: flex; align-items: center; justify-content: space-between; margin-bottom: 10px; }
.stat-label {
  font-size: 0.6875rem; font-weight: 600; text-transform: uppercase;
  letter-spacing: 0.06em; color: var(--text-muted);
}
.stat-icon {
  width: 32px; height: 32px; display: flex; align-items: center; justify-content: center; border-radius: 8px;
}
.icon-blue { background: var(--blue-dim); color: var(--blue-bright); }
.icon-purple { background: var(--purple-dim); color: var(--purple-base); }
.icon-amber { background: var(--amber-dim); color: var(--amber-bright); }
.icon-green { background: var(--green-dim); color: var(--green-bright); }
.stat-value { font-size: 1.75rem; font-weight: 700; letter-spacing: -0.02em; }
.stat-footer { font-size: 0.625rem; margin-top: 4px; text-transform: uppercase; letter-spacing: 0.06em; }

/* Monitor */
.monitor-bar {
  display: flex; align-items: center; justify-content: space-between;
  padding: 12px 18px; margin-bottom: 14px;
}
.monitor-left { display: flex; align-items: center; gap: 10px; }
.monitor-label { font-size: 0.8125rem; font-weight: 500; }
.monitor-interval {
  font-size: 0.625rem; background: var(--bg-elevated);
  padding: 2px 7px; border-radius: 4px; color: var(--text-dim);
}

/* Alerts */
.alerts-section { padding: 18px 20px; }
.alerts-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 14px; }
.view-link {
  font-size: 0.75rem; font-weight: 600; color: var(--blue-base);
  transition: color var(--duration-fast) ease;
}
.view-link:hover { color: var(--blue-bright); }

.alert-list { display: flex; flex-direction: column; gap: 2px; }
.alert-row {
  display: flex; align-items: center; gap: 14px;
  padding: 10px 12px; border-radius: 8px;
  transition: background var(--duration-fast) ease;
  animation: fadeIn 0.3s var(--ease-out) both;
}
.alert-row:hover { background: rgba(139, 151, 173, 0.04); }
.alert-addr { font-size: 0.75rem; min-width: 110px; }
.alert-diff { display: flex; align-items: center; gap: 6px; font-size: 0.75rem; }
.alert-time { font-size: 0.6875rem; margin-left: auto; white-space: nowrap; }
</style>
