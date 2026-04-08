<template>
  <div class="page-container">
    <div class="page-header">
      <div>
        <h1 class="page-title">Alerts</h1>
        <p class="page-subtitle">Triggered notification history</p>
      </div>
      <span v-if="alerts.length > 0" class="alert-count badge badge-blue">{{ alerts.length }} total</span>
    </div>

    <div class="glass alerts-card animate-fade-up">
      <div v-if="alerts.length > 0" class="table-wrap">
        <table class="data-table">
          <thead>
            <tr>
              <th>Time</th>
              <th>Address</th>
              <th>Type</th>
              <th>Before</th>
              <th>After</th>
              <th>Change</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="alert in alerts" :key="alert.id" class="alert-row">
              <td class="time-cell">
                <div class="time-primary">{{ formatTime(alert.sent_at) }}</div>
                <div class="time-date text-dim font-mono">{{ formatDateOnly(alert.sent_at) }}</div>
              </td>
              <td>
                <span class="addr-chip font-mono" :title="alert.address || alert.addressStr">{{ truncateAddress(alert.address || alert.addressStr) }}</span>
              </td>
              <td>
                <span class="badge badge-blue">{{ alert.alert_type }}</span>
              </td>
              <td class="font-mono val-cell text-muted">{{ formatSOL(alert.old_value) }}</td>
              <td class="font-mono val-cell">{{ formatSOL(alert.new_value) }}</td>
              <td>
                <span class="change-chip font-mono" :class="getChange(alert) >= 0 ? 'positive' : 'negative'">
                  <svg v-if="getChange(alert) >= 0" width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round"><polyline points="18 15 12 9 6 15"/></svg>
                  <svg v-else width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round"><polyline points="6 9 12 15 18 9"/></svg>
                  {{ Math.abs(getChange(alert)).toFixed(4) }}
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <EmptyState v-if="alerts.length === 0" title="No alerts" message="Alerts will appear here when thresholds are triggered" />
    </div>

    <div v-if="alerts.length >= limit" class="load-more">
      <button @click="loadMore" class="btn btn-ghost btn-sm">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><polyline points="23 4 23 10 17 10"/><path d="M20.49 15a9 9 0 11-2.12-9.36L23 10"/></svg>
        Load More
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { truncateAddress, formatSOL } from '../utils/format.js'
import { useSSE } from '../utils/api.js'
import EmptyState from '../components/EmptyState.vue'

const alerts = ref([])
const limit = ref(50)
const offset = ref(0)

const fetchAlerts = async () => {
  try {
    const res = await fetch(`/api/alerts?limit=${limit.value}&offset=${offset.value}`, {
      headers: { Authorization: 'Bearer solana-monitor-secret-key-2024' }
    })
    const newAlerts = await res.json()
    if (offset.value === 0) {
      alerts.value = newAlerts
    } else {
      alerts.value = [...alerts.value, ...newAlerts]
    }
  } catch (e) { console.error('Failed to fetch alerts:', e) }
}

const loadMore = () => {
  offset.value += limit.value
  fetchAlerts()
}

const getChange = (alert) => {
  return (alert.new_value || 0) - (alert.old_value || 0)
}

const formatTime = (date) => {
  if (!date) return ''
  return new Date(date).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' })
}

const formatDateOnly = (date) => {
  if (!date) return ''
  return new Date(date).toLocaleDateString([], { month: 'short', day: 'numeric', year: 'numeric' })
}

onMounted(() => {
  fetchAlerts()
  const { disconnect } = useSSE('/api/events', (type) => {
    if (type === 'new_alert') fetchAlerts()
  })
  onUnmounted(disconnect)
})
</script>

<style scoped>
.alerts-card { overflow: hidden; padding: 0; }
.alert-count { font-size: 0.75rem; font-weight: 600; }

.table-wrap { overflow-x: auto; }

.data-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.8125rem;
}
.data-table th {
  padding: 12px 18px;
  text-align: left;
  font-size: 0.6875rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: var(--text-dim);
  border-bottom: 1px solid var(--border-subtle);
  background: var(--bg-void);
  position: sticky;
  top: 0;
}
.data-table td {
  padding: 11px 18px;
  border-bottom: 1px solid var(--border-subtle);
  transition: background var(--duration-fast) ease;
}
.data-table tr:last-child td { border-bottom: none; }
.data-table tbody tr:hover td { background: rgba(139, 151, 173, 0.04); }

.alert-row { animation: fadeIn 0.3s var(--ease-out) both; }

.time-cell { white-space: nowrap; }
.time-primary { font-size: 0.8125rem; font-weight: 500; color: var(--text-primary); }
.time-date { font-size: 0.625rem; }

.addr-chip {
  display: inline-block;
  padding: 3px 10px;
  background: var(--blue-dim);
  border: 1px solid rgba(59, 130, 246, 0.15);
  border-radius: 6px;
  font-size: 0.75rem;
  color: var(--blue-bright);
}

.val-cell { color: var(--text-secondary); }

.change-chip {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 0.75rem;
  font-weight: 600;
  padding: 3px 10px;
  border-radius: 6px;
}
.change-chip.positive { color: var(--green-bright); background: var(--green-dim); }
.change-chip.negative { color: var(--red-bright); background: var(--red-dim); }

.load-more {
  display: flex;
  justify-content: center;
  padding: 18px;
}
</style>
