<template>
  <div class="page-container">
    <div class="stats-grid stagger">
      <div v-for="stat in statCards" :key="stat.label" class="stat-card holo-border animate-fade-up">
        <div class="stat-top">
          <span class="stat-label font-display">{{ t(stat.label) }}</span>
          <div class="stat-icon" :class="stat.iconClass" v-html="stat.icon"></div>
        </div>
        <div class="stat-value">{{ stats[stat.key] ?? 0 }}</div>
        <div class="stat-footer text-dim font-mono">{{ t(stat.sub) }}</div>
      </div>
    </div>

    <div class="glass monitor-bar animate-fade-up" style="animation-delay: 180ms">
      <div class="monitor-left">
        <span class="status-dot"></span>
        <span class="monitor-label">{{ t('dashboard.monitorService') }}</span>
        <span class="monitor-interval font-mono">{{ t('dashboard.every30s') }}</span>
      </div>
      <span class="badge badge-success">{{ monitorRunning ? t('dashboard.running') : t('dashboard.stopped') }}</span>
      <button class="btn btn-ghost btn-sm" @click="backupDB">
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
        {{ t('dashboard.backup') }}
      </button>
    </div>

    <div class="glass alerts-section animate-fade-up" style="animation-delay: 260ms">
      <div class="alerts-header">
        <h2 class="section-title">{{ t('dashboard.recentAlerts') }}</h2>
        <router-link to="/alerts" class="view-link">{{ t('dashboard.viewAll') }}</router-link>
      </div>
      <EmptyState v-if="!recentAlerts || recentAlerts.length === 0" :message="t('dashboard.noAlerts')" />
      <div v-else class="alert-list">
        <div v-for="alert in recentAlerts" :key="alert.id" class="alert-row">
          <span class="alert-addr font-mono neon-cyan">{{ truncateAddress(alert.address || alert.addressStr) }}</span>
          <span class="badge badge-cyan">{{ alert.alert_type }}</span>
          <span class="alert-diff font-mono">
            <span class="text-muted">{{ formatSOL(alert.old_value) }}</span>
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" style="color:var(--cyan-base)"><polyline points="9 18 15 12 9 6"/></svg>
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
import { useI18n } from '../utils/i18n.js'
import EmptyState from '../components/EmptyState.vue'

const toast = useToast()
const { t } = useI18n()

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
    toast.success(t('dashboard.backupDownloaded'))
  } catch (e) { toast.error(e.message) }
}

const statCards = [
  {
    label: 'dashboard.wallets', key: 'total_addresses', iconClass: 'icon-cyan', sub: 'dashboard.monitored',
    icon: '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 12V7H5a2 2 0 010-4h14v4"/><path d="M3 5v14a2 2 0 002 2h16v-5"/><path d="M18 12a2 2 0 000 4h4v-4z"/></svg>'
  },
  {
    label: 'dashboard.channels', key: 'total_notifications', iconClass: 'icon-purple', sub: 'dashboard.active',
    icon: '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M18 8A6 6 0 006 8c0 7-3 9-3 9h18s-3-2-3-9"/><path d="M13.73 21a2 2 0 01-3.46 0"/></svg>'
  },
  {
    label: 'dashboard.totalAlerts', key: 'total_alerts', iconClass: 'icon-amber', sub: 'dashboard.allTime',
    icon: '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>'
  },
  {
    label: 'dashboard.today', key: 'today_alerts', iconClass: 'icon-green', sub: 'dashboard.last24h',
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
  gap: 12px;
  margin-bottom: 12px;
}
@media (max-width: 900px) { .stats-grid { grid-template-columns: repeat(2, 1fr); } }

.stat-card { padding: 16px 18px; }
.stat-top { display: flex; align-items: center; justify-content: space-between; margin-bottom: 10px; }
.stat-label {
  font-size: 0.5625rem; font-weight: 600; text-transform: uppercase;
  letter-spacing: 0.1em; color: var(--cyan-base);
}
.stat-icon {
  width: 30px; height: 30px; display: flex; align-items: center; justify-content: center; border-radius: 8px;
}
.icon-cyan { background: var(--cyan-dim); color: var(--cyan-base); }
.icon-purple { background: var(--purple-dim); color: var(--purple-base); }
.icon-amber { background: var(--amber-dim); color: var(--amber-base); }
.icon-green { background: var(--green-dim); color: var(--green-base); }
.stat-value {
  font-family: var(--font-mono);
  font-size: 1.5rem; font-weight: 700; letter-spacing: -0.02em;
}
.stat-footer { font-size: 0.625rem; margin-top: 4px; text-transform: uppercase; letter-spacing: 0.06em; }

/* Monitor */
.monitor-bar {
  display: flex; align-items: center; justify-content: space-between;
  padding: 10px 16px; margin-bottom: 12px;
}
.monitor-left { display: flex; align-items: center; gap: 10px; }
.monitor-label { font-size: 0.8125rem; font-weight: 500; }
.monitor-interval {
  font-family: var(--font-mono);
  font-size: 0.625rem; background: var(--bg-elevated);
  padding: 2px 7px; border-radius: 4px; color: var(--text-dim);
}
.status-dot {
  width: 8px; height: 8px; border-radius: 50%;
  background: var(--green-base);
  box-shadow: 0 0 6px var(--green-base), 0 0 12px var(--green-base);
}

/* Alerts */
.alerts-section { padding: 16px 18px; }
.alerts-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 14px; }
.view-link {
  font-size: 0.6875rem; font-weight: 600; color: var(--cyan-base);
  transition: color var(--duration-fast) ease;
}
.view-link:hover { opacity: 0.85; }

.alert-list { display: flex; flex-direction: column; gap: 2px; }
.alert-row {
  display: flex; align-items: center; gap: 14px;
  padding: 8px 12px; border-radius: 6px;
  transition: background var(--duration-fast) ease;
  animation: fadeIn 0.3s var(--ease-out) both;
}
.alert-row:hover { background: rgba(139, 151, 173, 0.04); }
.alert-addr { font-size: 0.75rem; min-width: 110px; }
.alert-diff { display: flex; align-items: center; gap: 6px; font-size: 0.75rem; }
.alert-time { font-size: 0.6875rem; margin-left: auto; white-space: nowrap; }
</style>
