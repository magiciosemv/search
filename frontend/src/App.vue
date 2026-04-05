<template>
  <div class="app-layout">
    <aside class="sidebar">
      <div class="sidebar-brand">
        <div class="brand-icon">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
            <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </div>
        <div>
          <div class="brand-name">SOLANA</div>
          <div class="brand-sub">MONITOR</div>
        </div>
      </div>

      <nav class="sidebar-nav">
        <router-link v-for="(item, i) in navItems" :key="item.path" :to="item.path"
          class="nav-item" :class="{ active: $route.path === item.path }"
          :style="{ animationDelay: i * 60 + 'ms' }">
          <div class="nav-icon" v-html="item.icon"></div>
          <span class="nav-label">{{ item.name }}</span>
          <div v-if="$route.path === item.path" class="nav-indicator"></div>
        </router-link>
      </nav>

      <div class="sidebar-footer">
        <div class="status-line">
          <div class="glow-dot"></div>
          <span class="status-text">System Active</span>
        </div>
      </div>
    </aside>

    <main class="main-content">
      <router-view v-slot="{ Component }">
        <transition name="page" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>
  </div>
</template>

<script setup>
const navItems = [
  {
    path: '/', name: 'Dashboard',
    icon: '<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="7" height="7" rx="1"/><rect x="14" y="3" width="7" height="7" rx="1"/><rect x="3" y="14" width="7" height="7" rx="1"/><rect x="14" y="14" width="7" height="7" rx="1"/></svg>'
  },
  {
    path: '/addresses', name: 'Wallets',
    icon: '<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 12V7H5a2 2 0 010-4h14v4"/><path d="M3 5v14a2 2 0 002 2h16v-5"/><path d="M18 12a2 2 0 000 4h4v-4z"/></svg>'
  },
  {
    path: '/notifications', name: 'Notifications',
    icon: '<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M18 8A6 6 0 006 8c0 7-3 9-3 9h18s-3-2-3-9"/><path d="M13.73 21a2 2 0 01-3.46 0"/></svg>'
  },
  {
    path: '/alerts', name: 'Alerts',
    icon: '<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>'
  }
]
</script>

<style scoped>
.app-layout {
  display: flex;
  height: 100vh;
  overflow: hidden;
}

.sidebar {
  width: var(--sidebar-w);
  min-width: var(--sidebar-w);
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: linear-gradient(180deg, var(--bg-surface) 0%, var(--bg-base) 100%);
  border-right: 1px solid var(--border-subtle);
  position: relative;
  overflow: hidden;
}

.sidebar::after {
  content: '';
  position: absolute;
  top: 0; right: 0;
  width: 1px;
  height: 100%;
  background: linear-gradient(180deg, var(--cyan-base), transparent 60%);
  opacity: 0.3;
}

.sidebar-brand {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 24px 20px;
  border-bottom: 1px solid var(--border-subtle);
}

.brand-icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, rgba(6, 182, 212, 0.15), rgba(6, 182, 212, 0.05));
  border: 1px solid var(--border-default);
  border-radius: 10px;
  color: var(--cyan-bright);
}

.brand-name {
  font-size: 1rem;
  font-weight: 800;
  letter-spacing: 0.15em;
  color: var(--text-primary);
  line-height: 1;
}

.brand-sub {
  font-size: 0.625rem;
  font-weight: 600;
  letter-spacing: 0.3em;
  color: var(--cyan-base);
  margin-top: 2px;
}

.sidebar-nav {
  flex: 1;
  padding: 16px 12px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 11px 14px;
  border-radius: 8px;
  color: var(--text-muted);
  font-size: 0.875rem;
  font-weight: 500;
  transition: all 0.2s ease;
  position: relative;
  animation: slideIn 0.4s var(--ease-out) both;
  text-decoration: none;
  cursor: pointer;
  border: 1px solid transparent;
}

.nav-item:hover {
  color: var(--text-secondary);
  background: rgba(6, 182, 212, 0.04);
}

.nav-item.active {
  color: var(--cyan-bright);
  background: rgba(6, 182, 212, 0.08);
  border: 1px solid rgba(6, 182, 212, 0.12);
}

.nav-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
  opacity: 0.7;
}

.nav-item.active .nav-icon { opacity: 1; }

.nav-indicator {
  position: absolute;
  right: -13px;
  width: 3px;
  height: 20px;
  background: var(--cyan-base);
  border-radius: 3px 0 0 3px;
  box-shadow: 0 0 8px var(--cyan-base);
}

.sidebar-footer {
  padding: 16px 20px;
  border-top: 1px solid var(--border-subtle);
}

.status-line {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--text-muted);
  letter-spacing: 0.04em;
}

.status-text {
  text-transform: uppercase;
  font-family: var(--font-mono);
  font-size: 0.6875rem;
}

.main-content {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
}

/* Page transition */
.page-enter-active { animation: fadeInUp 0.35s var(--ease-out); }
.page-leave-active { animation: fadeIn 0.15s ease reverse; }
</style>
