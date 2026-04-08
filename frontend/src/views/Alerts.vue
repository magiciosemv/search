<template>
  <div class="page-container">
    <div class="page-header">
      <div>
        <h1 class="page-title">Alerts</h1>
        <p class="page-subtitle">Triggered notification history</p>
      </div>
      <div class="header-actions">
        <span v-if="filteredAlerts.length !== alerts.length" class="badge badge-blue">{{ filteredAlerts.length }} / {{ alerts.length }}</span>
        <button v-if="alerts.length > 0" class="btn btn-ghost btn-sm" @click="exportCSV">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
          Export CSV
        </button>
      </div>
    </div>

    <div v-if="alerts.length > 0" class="filter-bar glass">
      <div class="filter-group">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><polygon points="22 3 2 3 10 12.46 10 19 14 21 14 12.46 22 3"/></svg>
        <input v-model="filterAddr" type="text" placeholder="Filter by address..." class="input-field input-mono filter-input" />
        <select v-model="filterType" class="input-field filter-select">
          <option value="">All types</option>
          <option v-for="t in alertTypes" :key="t" :value="t">{{ t }}</option>
        </select>
      </div>
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
            <tr v-for="alert in filteredAlerts" :key="alert.id" class="alert-row">
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

    <EmptyState v-if="alerts.length > 0 && filteredAlerts.length === 0" title="No matches" message="Try different filter criteria" />

    <div v-if="alerts.length >= limit" class="load-more">
      <button @click="loadMore" class="btn btn-ghost btn-sm">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><polyline points="23 4 23 10 17 10"/><path d="M20.49 15a9 9 0 11-2.12-9.36L23 10"/></svg>
        Load More
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { truncateAddress, formatSOL } from '../utils/format.js'
import { useSSE } from '../utils/api.js'
import EmptyState from '../components/EmptyState.vue'
import { useToast } from '../utils/toast.js'
const toast = useToast()

const alerts = ref([])
const limit = ref(50)
const offset = ref(0)
const filterAddr = ref('')
const filterType = ref('')

const filteredAlerts = computed(() => {
  return alerts.value.filter(a => {
    if (filterAddr.value && !(a.address || a.addressStr || '').toLowerCase().includes(filterAddr.value.toLowerCase())) return false
    if (filterType.value && a.alert_type !== filterType.value) return false
    return true
  })
})

const alertTypes = computed(() => [...new Set(alerts.value.map(a => a.alert_type))])

const exportCSV = () => {
  const rows = [['Time', 'Address', 'Type', 'Before', 'After', 'Change']]
  filteredAlerts.value.forEach(a => {
    rows.push([
      a.sent_at || '',
      a.address || a.addressStr || '',
      a.alert_type || '',
      String(a.old_value || 0),
      String(a.new_value || 0),
      String(((a.new_value || 0) - (a.old_value || 0)).toFixed(4))
    ])
  })
  const csv = rows.map(r => r.map(c => `"${c}"`).join(',')).join('\n')
  const blob = new Blob([csv], { type: 'text/csv' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = `alerts-${new Date().toISOString().slice(0, 10)}.csv`
  link.click()
  URL.revokeObjectURL(url)
  toast.success('CSV exported')
}

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

.header-actions { display: flex; align-items: center; gap: 10px; }
.filter-bar { padding: 12px 16px; margin-bottom: 14px; }
.filter-group { display: flex; align-items: center; gap: 10px; color: var(--text-dim); }
.filter-input { width: 240px; }
.filter-select {
  width: 140px; appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 24 24' fill='none' stroke='%235A6680' stroke-width='2'%3E%3Cpolyline points='6 9 12 15 18 9'/%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 10px center;
  padding-right: 28px;
}
</style>
