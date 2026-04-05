<template>
  <div class="page-container">
    <div class="page-header">
      <div>
        <h1 class="page-title">Alert History</h1>
        <p class="page-subtitle">All triggered notifications</p>
      </div>
    </div>

    <div class="card alerts-table-card animate-fade-up">
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
            <tr v-for="alert in alerts" :key="alert.id">
              <td class="time-cell">
                <div class="time-primary">{{ formatTime(alert.sent_at) }}</div>
                <div class="time-date text-muted">{{ formatDateOnly(alert.sent_at) }}</div>
              </td>
              <td>
                <span class="addr-badge font-mono" :title="alert.address">{{ truncateAddress(alert.address) }}</span>
              </td>
              <td>
                <span class="badge badge-cyan">{{ alert.alert_type }}</span>
              </td>
              <td class="font-mono val-cell">{{ formatSOL(alert.old_value) }}</td>
              <td class="font-mono val-cell text-primary">{{ formatSOL(alert.new_value) }}</td>
              <td>
                <span class="change-val font-mono" :class="getChange(alert) >= 0 ? 'positive' : 'negative'">
                  <svg v-if="getChange(alert) >= 0" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round"><polyline points="18 15 12 9 6 15"/></svg>
                  <svg v-else width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round"><polyline points="6 9 12 15 18 9"/></svg>
                  {{ Math.abs(getChange(alert)).toFixed(4) }}
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <EmptyState v-if="alerts.length === 0" title="" message="No alerts recorded yet" />
    </div>

    <div v-if="alerts.length >= limit" class="load-more-wrap">
      <button @click="loadMore" class="btn btn-ghost">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><polyline points="23 4 23 10 17 10"/><path d="M20.49 15a9 9 0 11-2.12-9.36L23 10"/></svg>
        Load More
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { truncateAddress, formatSOL } from '../utils/format.js'
import EmptyState from '../components/EmptyState.vue'

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

onMounted(fetchAlerts)
</script>

<style scoped>
.alerts-table-card {
  overflow: hidden;
}

.table-wrap {
  overflow-x: auto;
}

.addr-badge {
  display: inline-block;
  padding: 3px 8px;
  background: rgba(6, 182, 212, 0.06);
  border: 1px solid var(--border-subtle);
  border-radius: 6px;
  font-size: 0.8125rem;
  color: var(--cyan-bright);
}

.time-cell { white-space: nowrap; }

.time-primary {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--text-primary);
}

.time-date {
  font-size: 0.6875rem;
  font-family: var(--font-mono);
}

.val-cell {
  font-size: 0.8125rem;
  color: var(--text-secondary);
}

.change-val {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 0.8125rem;
  font-weight: 600;
  padding: 3px 8px;
  border-radius: 6px;
}

.change-val.positive {
  color: var(--green-bright);
  background: rgba(16, 185, 129, 0.08);
}

.change-val.negative {
  color: var(--red-bright);
  background: rgba(239, 68, 68, 0.08);
}

.load-more-wrap {
  display: flex;
  justify-content: center;
  padding: 20px;
}
</style>
