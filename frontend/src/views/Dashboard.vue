<template>
  <div class="page-container">
    <div class="page-header">
      <div>
        <h1 class="page-title">Dashboard</h1>
        <p class="page-subtitle">Real-time wallet monitoring overview</p>
      </div>
    </div>

    <div class="stats-grid stagger">
      <div class="stat-card card animate-fade-up">
        <div class="stat-header">
          <span class="stat-label">Total Wallets</span>
          <div class="stat-icon wallet-icon">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 12V7H5a2 2 0 010-4h14v4"/><path d="M3 5v14a2 2 0 002 2h16v-5"/><path d="M18 12a2 2 0 000 4h4v-4z"/></svg>
          </div>
        </div>
        <div class="stat-value">{{ stats.total_addresses || 0 }}</div>
      </div>

      <div class="stat-card card animate-fade-up">
        <div class="stat-header">
          <span class="stat-label">Channels</span>
          <div class="stat-icon bell-icon">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M18 8A6 6 0 006 8c0 7-3 9-3 9h18s-3-2-3-9"/><path d="M13.73 21a2 2 0 01-3.46 0"/></svg>
          </div>
        </div>
        <div class="stat-value">{{ stats.total_notifications || 0 }}</div>
      </div>

      <div class="stat-card card animate-fade-up">
        <div class="stat-header">
          <span class="stat-label">Total Alerts</span>
          <div class="stat-icon alert-icon">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>
          </div>
        </div>
        <div class="stat-value">{{ stats.total_alerts || 0 }}</div>
      </div>

      <div class="stat-card card animate-fade-up">
        <div class="stat-header">
          <span class="stat-label">Today</span>
          <div class="stat-icon today-icon">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
          </div>
        </div>
        <div class="stat-value">{{ stats.today_alerts || 0 }}</div>
      </div>
    </div>

    <!-- Monitor Status -->
    <div class="monitor-card card animate-fade-up" style="animation-delay: 200ms">
      <div class="monitor-info">
        <div>
          <div class="monitor-title">Monitor Service</div>
          <div class="monitor-desc">Background balance checking (30s interval)</div>
        </div>
        <span class="badge" :class="stats.monitor_running ? 'badge-success' : 'badge-danger'">
          <span class="badge-dot" :class="stats.monitor_running ? 'dot-green' : 'dot-red'"></span>
          {{ stats.monitor_running ? 'Running' : 'Stopped' }}
        </span>
      </div>
      <div class="monitor-bar">
        <div class="monitor-bar-fill" :class="{ active: stats.monitor_running }"></div>
      </div>
    </div>

    <!-- Recent Alerts -->
    <div class="card recent-card animate-fade-up" style="animation-delay: 280ms">
      <div class="recent-header">
        <h2 class="section-title">Recent Alerts</h2>
        <router-link to="/alerts" class="view-all-link">View All</router-link>
      </div>

      <EmptyState v-if="!recentAlerts || recentAlerts.length === 0"
        title="" message="No alerts yet — add wallets to start monitoring" />

      <div v-else class="alert-list">
        <div v-for="(alert, i) in recentAlerts" :key="alert.id"
          class="alert-row" :style="{ animationDelay: (300 + i * 60) + 'ms' }">
          <div class="alert-addr font-mono">{{ truncateAddress(alert.address || alert.addressStr) }}</div>
          <div class="alert-type-badge badge badge-cyan">{{ alert.alert_type }}</div>
          <div class="alert-values font-mono">
            <span class="text-muted">{{ formatSOL(alert.old_value) }}</span>
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="color: var(--text-muted)"><polyline points="9 18 15 12 9 6"/></svg>
            <span class="text-primary">{{ formatSOL(alert.new_value) }}</span>
          </div>
          <div class="alert-time text-muted">{{ formatDate(alert.sent_at) }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, onUnmounted } from 'vue'
import { useFetch, useSSE } from '../utils/api.js'
import { truncateAddress, formatDate, formatSOL } from '../utils/format.js'
import EmptyState from '../components/EmptyState.vue'

const { data: stats, refetch: refetchStats } = useFetch('/api/stats', {})
const { data: recentAlerts, refetch: refetchAlerts } = useFetch('/api/alerts?limit=5', [])

onMounted(() => {
  const { disconnect } = useSSE('/api/events', (type, data) => {
    if (type === 'balance_update') {
      refetchStats()
    }
    if (type === 'new_alert') {
      refetchAlerts()
      refetchStats()
    }
  })
  onUnmounted(disconnect)
})
</script>

<style scoped>
.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 20px;
}

@media (max-width: 900px) {
  .stats-grid { grid-template-columns: repeat(2, 1fr); }
}

.stat-card {
  padding: 20px;
}

.stat-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.stat-label {
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.stat-icon {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
}

.wallet-icon { background: rgba(6, 182, 212, 0.1); color: var(--cyan-bright); }
.bell-icon { background: rgba(139, 92, 246, 0.1); color: var(--purple-base); }
.alert-icon { background: rgba(245, 158, 11, 0.1); color: var(--amber-bright); }
.today-icon { background: rgba(16, 185, 129, 0.1); color: var(--green-bright); }

.stat-value {
  font-size: 1.75rem;
  font-weight: 800;
  font-family: var(--font-mono);
  letter-spacing: -0.02em;
}

/* Monitor */
.monitor-card {
  padding: 20px 24px;
  margin-bottom: 20px;
}

.monitor-info {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 14px;
}

.monitor-title {
  font-size: 1rem;
  font-weight: 600;
}

.monitor-desc {
  font-size: 0.8125rem;
  color: var(--text-muted);
  margin-top: 2px;
}

.badge-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}

.dot-green {
  background: var(--green-bright);
  box-shadow: 0 0 6px var(--green-base);
}

.dot-red {
  background: var(--red-bright);
  box-shadow: 0 0 6px var(--red-base);
}

.monitor-bar {
  height: 2px;
  background: var(--bg-elevated);
  border-radius: 2px;
  overflow: hidden;
}

.monitor-bar-fill {
  height: 100%;
  width: 0%;
  background: linear-gradient(90deg, var(--cyan-base), var(--green-base));
  border-radius: 2px;
  transition: width 1s ease;
}

.monitor-bar-fill.active {
  width: 100%;
  animation: barPulse 2s ease-in-out infinite;
}

@keyframes barPulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.4; }
}

/* Recent Alerts */
.recent-card {
  padding: 24px;
}

.recent-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.section-title {
  font-size: 1rem;
  font-weight: 600;
}

.view-all-link {
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--cyan-base);
  transition: color 0.15s ease;
}

.view-all-link:hover { color: var(--cyan-bright); }

.alert-list {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.alert-row {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 10px 12px;
  border-radius: 8px;
  transition: background 0.15s ease;
  animation: fadeIn 0.4s var(--ease-out) both;
}

.alert-row:hover { background: var(--bg-hover); }

.alert-addr {
  font-size: 0.8125rem;
  color: var(--cyan-bright);
  min-width: 130px;
}

.alert-type-badge { font-size: 0.6875rem; flex-shrink: 0; }

.alert-values {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.8125rem;
}

.alert-time {
  font-size: 0.75rem;
  margin-left: auto;
  white-space: nowrap;
}
</style>
