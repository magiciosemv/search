<template>
  <div class="layout">
    <!-- Particle system -->
    <div class="particles">
      <div v-for="p in particles" :key="p.id" class="particle" :style="p.style"></div>
    </div>

    <!-- Sidebar -->
    <aside class="sidebar" :class="{ collapsed: sidebarCollapsed, open: sidebarOpen }">
      <div class="sidebar-brand">
        <div class="brand-icon neon-cyan">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5"/></svg>
        </div>
        <div class="brand-text" v-show="!sidebarCollapsed">
          <span class="brand-name font-display">NEXUS</span>
          <span class="brand-ver font-mono">v2.0</span>
        </div>
      </div>

      <nav class="sidebar-nav">
        <div class="nav-divider" v-show="!sidebarCollapsed"></div>
        <router-link v-for="item in navItems" :key="item.path" :to="item.path"
          class="nav-item" :class="{ active: $route.path === item.path }"
          @click="sidebarOpen = false">
          <div class="nav-icon" :class="{ 'neon-cyan': $route.path === item.path }" v-html="item.icon"></div>
          <span class="nav-label" v-show="!sidebarCollapsed">{{ t('nav.' + item.key) }}</span>
        </router-link>
      </nav>

      <div class="sidebar-footer" v-show="!sidebarCollapsed">
        <div class="footer-row">
          <button class="icon-btn" @click="toggleTheme" :title="theme === 'dark' ? 'Light mode' : 'Dark mode'">
            <svg v-if="theme === 'dark'" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="5"/><line x1="12" y1="1" x2="12" y2="3"/><line x1="12" y1="21" x2="12" y2="23"/><line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/><line x1="1" y1="12" x2="3" y2="12"/><line x1="21" y1="12" x2="23" y2="12"/><line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/></svg>
            <svg v-else width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 12.79A9 9 0 1111.21 3 7 7 0 0021 12.79z"/></svg>
          </button>
          <button class="icon-btn" @click="setLocale(locale === 'en' ? 'zh' : 'en')">
            <span class="font-mono lang-text">{{ locale === 'en' ? 'EN' : '中' }}</span>
          </button>
          <div class="footer-status">
            <span class="status-dot"></span>
            <span class="status-label font-mono">{{ monitorRunning ? 'LIVE' : 'OFF' }}</span>
          </div>
        </div>
      </div>
    </aside>

    <!-- Mobile backdrop -->
    <Transition name="fade">
      <div v-if="sidebarOpen" class="sidebar-backdrop" @click="sidebarOpen = false"></div>
    </Transition>

    <!-- Main area -->
    <div class="main-area">
      <!-- Top navbar -->
      <header class="topbar">
        <div class="topbar-left">
          <button class="icon-btn sidebar-toggle" @click="sidebarCollapsed = !sidebarCollapsed" aria-label="Toggle sidebar">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><line x1="3" y1="6" x2="21" y2="6"/><line x1="3" y1="12" x2="21" y2="12"/><line x1="3" y1="18" x2="21" y2="18"/></svg>
          </button>
          <button class="icon-btn mobile-menu" @click="sidebarOpen = !sidebarOpen">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="3" y1="6" x2="21" y2="6"/><line x1="3" y1="12" x2="21" y2="12"/><line x1="3" y1="18" x2="21" y2="18"/></svg>
          </button>
          <span class="topbar-page font-display">{{ currentPage }}</span>
        </div>
        <div class="topbar-right">
          <div class="topbar-clock font-mono">{{ currentTime }}</div>
          <div class="topbar-status">
            <span class="status-dot"></span>
          </div>
        </div>
      </header>

      <!-- Content -->
      <main class="main-content">
        <router-view v-slot="{ Component }">
          <transition name="page" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </main>
    </div>

    <ToastContainer />
    <ConfirmDialog />
    <OnboardingModal v-model="showOnboarding" />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import ToastContainer from './components/ToastContainer.vue'
import ConfirmDialog from './components/ConfirmDialog.vue'
import { useTheme } from './utils/theme.js'
import { useI18n } from './utils/i18n.js'
import OnboardingModal from './components/OnboardingModal.vue'
const { theme, toggleTheme } = useTheme()
const { locale, setLocale, t } = useI18n()
const showOnboarding = ref(!localStorage.getItem('onboarded'))

const route = useRoute()
const sidebarCollapsed = ref(false)
const sidebarOpen = ref(false)
const monitorRunning = ref(true)
const currentTime = ref('')

// Particle system
const particles = ref([])
for (let i = 0; i < 20; i++) {
  particles.value.push({
    id: i,
    style: {
      left: Math.random() * 100 + '%',
      animationDuration: (8 + Math.random() * 15) + 's',
      animationDelay: (Math.random() * 10) + 's',
      width: (1 + Math.random() * 2) + 'px',
      height: (1 + Math.random() * 2) + 'px',
    }
  })
}

const navItems = [
  { path: '/', key: 'dashboard', icon: '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="7" height="7" rx="1"/><rect x="14" y="3" width="7" height="7" rx="1"/><rect x="3" y="14" width="7" height="7" rx="1"/><rect x="14" y="14" width="7" height="7" rx="1"/></svg>' },
  { path: '/addresses', key: 'wallets', icon: '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 12V7H5a2 2 0 010-4h14v4"/><path d="M3 5v14a2 2 0 002 2h16v-5"/><path d="M18 12a2 2 0 000 4h4v-4z"/></svg>' },
  { path: '/notifications', key: 'channels', icon: '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M18 8A6 6 0 006 8c0 7-3 9-3 9h18s-3-2-3-9"/><path d="M13.73 21a2 2 0 01-3.46 0"/></svg>' },
  { path: '/alerts', key: 'alerts', icon: '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>' },
  { path: '/rules', key: 'rules', icon: '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/></svg>' }
]

const pageNames = { '/': 'nav.dashboard', '/addresses': 'nav.wallets', '/notifications': 'nav.channels', '/alerts': 'nav.alerts', '/rules': 'nav.rules' }
const currentPage = computed(() => t(pageNames[route.path] || ''))

let timer
onMounted(() => {
  const update = () => { currentTime.value = new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' }) }
  update()
  timer = setInterval(update, 1000)
})
onUnmounted(() => clearInterval(timer))
</script>

<style scoped>
.layout {
  display: flex;
  height: 100vh;
  overflow: hidden;
}

/* ---- Particles ---- */
.particles {
  position: fixed;
  inset: 0;
  pointer-events: none;
  z-index: 0;
  overflow: hidden;
}
.particle {
  position: absolute;
  bottom: -10px;
  border-radius: 50%;
  background: var(--cyan-base);
  opacity: 0;
  animation: particleFloat linear infinite;
  box-shadow: 0 0 6px var(--cyan-glow);
}

/* ---- Sidebar ---- */
.sidebar {
  width: var(--sidebar-width);
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: var(--bg-base);
  border-right: 1px solid var(--border-subtle);
  flex-shrink: 0;
  z-index: 20;
  transition: width 0.3s var(--ease-out);
}
.sidebar.collapsed {
  width: 56px;
}

.sidebar-brand {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 14px 16px;
  border-bottom: 1px solid var(--border-subtle);
  flex-shrink: 0;
}
.brand-text { display: flex; align-items: baseline; gap: 6px; }
.brand-name {
  font-size: 0.75rem;
  font-weight: 800;
  letter-spacing: 0.2em;
}
.brand-ver {
  font-size: 0.5rem;
  color: var(--text-dim);
}

/* ---- Nav ---- */
.sidebar-nav {
  flex: 1;
  padding: 8px;
  display: flex;
  flex-direction: column;
  gap: 2px;
  overflow-y: auto;
}
.nav-divider {
  height: 1px;
  background: linear-gradient(90deg, transparent, var(--border-default), transparent);
  margin: 4px 8px 8px;
}
.nav-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 9px 12px;
  border-radius: 8px;
  color: var(--text-muted);
  transition: all var(--duration-fast) ease;
  cursor: pointer;
  text-decoration: none;
  position: relative;
  font-weight: 500;
  font-size: 0.8125rem;
}
.nav-item:hover { color: var(--text-secondary); background: rgba(0, 240, 255, 0.03); }
.nav-item.active {
  color: var(--cyan-base);
  background: rgba(0, 240, 255, 0.06);
}
.nav-item.active::before {
  content: '';
  position: absolute;
  left: 0; top: 50%;
  transform: translateY(-50%);
  width: 2px; height: 18px;
  background: var(--cyan-base);
  border-radius: 0 2px 2px 0;
  box-shadow: 0 0 8px var(--cyan-glow);
}
.nav-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 18px; height: 18px;
  flex-shrink: 0;
  transition: all var(--duration-fast) ease;
}

/* ---- Sidebar Footer ---- */
.sidebar-footer { padding: 12px; border-top: 1px solid var(--border-subtle); flex-shrink: 0; }
.footer-row { display: flex; align-items: center; gap: 6px; }
.icon-btn {
  display: flex; align-items: center; justify-content: center;
  width: 28px; height: 28px; border-radius: 6px;
  border: 1px solid var(--border-subtle); background: transparent;
  color: var(--text-dim); cursor: pointer;
  transition: all var(--duration-fast) ease;
}
.icon-btn:hover { color: var(--cyan-base); border-color: var(--border-default); background: rgba(0, 240, 255, 0.03); }
.lang-text { font-size: 0.5625rem; font-weight: 700; }
.footer-status { display: flex; align-items: center; gap: 6px; margin-left: auto; }
.status-label { font-size: 0.5rem; font-weight: 600; letter-spacing: 0.08em; color: var(--green-base); }

/* ---- Main Area ---- */
.main-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
  overflow: hidden;
}

/* ---- Topbar ---- */
.topbar {
  height: var(--topbar-height);
  display: flex;
  align-items: center;
  padding: 0 20px;
  gap: 12px;
  border-bottom: 1px solid var(--border-subtle);
  background: var(--bg-void);
  flex-shrink: 0;
}
.topbar-left { display: flex; align-items: center; gap: 8px; }
.sidebar-toggle { display: flex; }
.mobile-menu { display: none; }
.topbar-page {
  font-size: 0.6875rem;
  font-weight: 600;
  letter-spacing: 0.08em;
  color: var(--text-secondary);
}
.topbar-right { display: flex; align-items: center; gap: 12px; margin-left: auto; }
.topbar-clock { font-size: 0.625rem; color: var(--text-dim); letter-spacing: 0.04em; }
.topbar-status { display: flex; align-items: center; }

/* ---- Main Content ---- */
.main-content {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
}

/* ---- Backdrop ---- */
.sidebar-backdrop {
  display: none;
  position: fixed; inset: 0;
  background: rgba(6, 6, 14, 0.7);
  z-index: 15;
  backdrop-filter: blur(4px);
}

/* ---- Transitions ---- */
.page-enter-active { animation: fadeInUp 0.3s var(--ease-out); }
.page-leave-active { animation: fadeIn 0.1s ease reverse; }
.fade-enter-active { animation: fadeIn 0.2s ease; }
.fade-leave-active { animation: fadeIn 0.15s ease reverse; }

/* ---- Mobile ---- */
@media (max-width: 768px) {
  .sidebar-toggle { display: none; }
  .mobile-menu { display: flex; }
  .sidebar {
    position: fixed;
    top: 0; left: 0;
    height: 100vh;
    width: var(--sidebar-width) !important;
    transform: translateX(-100%);
    transition: transform 0.3s var(--ease-out);
    z-index: 25;
  }
  .sidebar.open { transform: translateX(0); }
  .sidebar-backdrop { display: block; }
  .topbar { padding: 0 12px; }
}
</style>
