<template>
  <div class="layout">
    <!-- Sidebar -->
    <aside class="sidebar" :class="{ open: sidebarOpen }">
      <div class="sidebar-brand">
        <div class="brand-icon">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
            <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5"/>
          </svg>
        </div>
        <span class="brand-name">SOLANA</span>
        <span class="brand-sub">MONITOR</span>
      </div>

      <nav class="sidebar-nav">
        <div class="nav-section-label">Main</div>
        <router-link v-for="item in navItems" :key="item.path" :to="item.path"
          class="nav-item" :class="{ active: $route.path === item.path }"
          @click="sidebarOpen = false">
          <div class="nav-icon" v-html="item.icon"></div>
          <span class="nav-label">{{ t('nav.' + item.key) }}</span>
        </router-link>
      </nav>

      <div class="sidebar-footer">
        <div class="footer-controls">
          <button class="theme-toggle" @click="toggleTheme" :title="theme === 'dark' ? 'Switch to light' : 'Switch to dark'">
            <svg v-if="theme === 'dark'" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><circle cx="12" cy="12" r="5"/><line x1="12" y1="1" x2="12" y2="3"/><line x1="12" y1="21" x2="12" y2="23"/><line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/><line x1="1" y1="12" x2="3" y2="12"/><line x1="21" y1="12" x2="23" y2="12"/><line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/></svg>
            <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M21 12.79A9 9 0 1111.21 3 7 7 0 0021 12.79z"/></svg>
          </button>
          <button class="lang-toggle" @click="setLocale(locale === 'en' ? 'zh' : 'en')" :title="locale === 'en' ? '切换中文' : 'Switch to English'">
            <span class="lang-label font-mono">{{ locale === 'en' ? 'EN' : '中' }}</span>
          </button>
          <div class="sidebar-status">
            <span class="status-dot"></span>
            <span class="status-text">{{ monitorRunning ? t('dashboard.running') : t('dashboard.stopped') }}</span>
          </div>
        </div>
      </div>
    </aside>

    <!-- Sidebar backdrop (mobile) -->
    <Transition name="fade">
      <div v-if="sidebarOpen" class="sidebar-backdrop" @click="sidebarOpen = false"></div>
    </Transition>

    <!-- Main area -->
    <div class="main-area">
      <!-- Top bar -->
      <header class="topbar">
        <button class="btn btn-icon mobile-menu-btn" @click="sidebarOpen = !sidebarOpen" aria-label="Toggle menu">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
            <line x1="3" y1="6" x2="21" y2="6"/><line x1="3" y1="12" x2="21" y2="12"/><line x1="3" y1="18" x2="21" y2="18"/>
          </svg>
        </button>
        <div class="topbar-breadcrumb">
          <span class="breadcrumb-page">{{ currentPage }}</span>
        </div>
        <div class="topbar-right">
          <div class="topbar-time font-mono">{{ currentTime }}</div>
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
const sidebarOpen = ref(false)
const monitorRunning = ref(true)
const currentTime = ref('')

const navItems = [
  {
    path: '/', key: 'dashboard',
    icon: '<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="7" height="7" rx="1.5"/><rect x="14" y="3" width="7" height="7" rx="1.5"/><rect x="3" y="14" width="7" height="7" rx="1.5"/><rect x="14" y="14" width="7" height="7" rx="1.5"/></svg>'
  },
  {
    path: '/addresses', key: 'wallets',
    icon: '<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M21 12V7H5a2 2 0 010-4h14v4"/><path d="M3 5v14a2 2 0 002 2h16v-5"/><path d="M18 12a2 2 0 000 4h4v-4z"/></svg>'
  },
  {
    path: '/notifications', key: 'channels',
    icon: '<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M18 8A6 6 0 006 8c0 7-3 9-3 9h18s-3-2-3-9"/><path d="M13.73 21a2 2 0 01-3.46 0"/></svg>'
  },
  {
    path: '/alerts', key: 'alerts',
    icon: '<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>'
  },
  {
    path: '/rules', key: 'rules',
    icon: '<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/><polyline points="10 9 9 9 8 9"/></svg>'
  }
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
  overflow: hidden;
}

.sidebar-brand {
  display: flex;
  align-items: baseline;
  gap: 0;
  padding: 20px 20px 20px 20px;
  border-bottom: 1px solid var(--border-subtle);
  flex-shrink: 0;
}

.brand-icon {
  color: var(--blue-base);
  margin-right: 10px;
  margin-top: 2px;
  display: flex;
  align-items: center;
}

.brand-name {
  font-family: var(--font-mono);
  font-size: 0.8125rem;
  font-weight: 700;
  letter-spacing: 0.18em;
  color: var(--text-primary);
}

.brand-sub {
  font-family: var(--font-mono);
  font-size: 0.625rem;
  font-weight: 500;
  letter-spacing: 0.18em;
  color: var(--text-dim);
  margin-left: 6px;
}

/* ---- Nav ---- */
.sidebar-nav {
  flex: 1;
  padding: 16px 10px;
  display: flex;
  flex-direction: column;
  gap: 2px;
  overflow-y: auto;
}

.nav-section-label {
  font-size: 0.625rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  color: var(--text-dim);
  padding: 4px 10px 10px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 11px;
  padding: 10px 12px;
  border-radius: 8px;
  color: var(--text-muted);
  transition: all var(--duration-fast) ease;
  cursor: pointer;
  text-decoration: none;
  position: relative;
  font-weight: 500;
  font-size: 0.8125rem;
}

.nav-item:hover {
  color: var(--text-secondary);
  background: rgba(139, 151, 173, 0.05);
}

.nav-item.active {
  color: var(--text-primary);
  background: rgba(59, 130, 246, 0.08);
}

.nav-item.active::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 3px;
  height: 20px;
  background: var(--blue-base);
  border-radius: 0 3px 3px 0;
}

.nav-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
  flex-shrink: 0;
}

.nav-label {
  white-space: nowrap;
}

/* ---- Sidebar Footer ---- */
.sidebar-footer {
  padding: 14px 16px;
  border-top: 1px solid var(--border-subtle);
  flex-shrink: 0;
}
.footer-controls {
  display: flex;
  align-items: center;
  gap: 10px;
}
.theme-toggle {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 30px;
  height: 30px;
  border-radius: 8px;
  border: 1px solid var(--border-subtle);
  background: transparent;
  color: var(--text-muted);
  cursor: pointer;
  transition: all var(--duration-fast) ease;
}
.theme-toggle:hover {
  color: var(--text-primary);
  background: var(--bg-hover);
  border-color: var(--border-default);
}
.lang-toggle {
  display: flex; align-items: center; justify-content: center;
  height: 30px; padding: 0 8px; border-radius: 8px;
  border: 1px solid var(--border-subtle); background: transparent;
  color: var(--text-muted); cursor: pointer;
  transition: all var(--duration-fast) ease;
}
.lang-toggle:hover {
  color: var(--text-primary); background: var(--bg-hover);
  border-color: var(--border-default);
}
.lang-label { font-size: 0.625rem; font-weight: 700; letter-spacing: 0.04em; }
.sidebar-status {
  display: flex;
  align-items: center;
  gap: 8px;
}
.status-text {
  font-family: var(--font-mono);
  font-size: 0.625rem;
  font-weight: 500;
  letter-spacing: 0.04em;
  color: var(--green-base);
  text-transform: uppercase;
}

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
  padding: 0 28px;
  gap: 16px;
  border-bottom: 1px solid var(--border-subtle);
  background: var(--bg-void);
  flex-shrink: 0;
}

.mobile-menu-btn {
  display: none;
  background: transparent;
  color: var(--text-secondary);
  border: none;
}

.topbar-breadcrumb {
  flex: 1;
}

.breadcrumb-page {
  font-weight: 600;
  font-size: 0.875rem;
  color: var(--text-secondary);
}

.topbar-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.topbar-time {
  font-size: 0.6875rem;
  color: var(--text-dim);
  letter-spacing: 0.03em;
}

/* ---- Main Content ---- */
.main-content {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
}

/* ---- Backdrop ---- */
.sidebar-backdrop {
  display: none;
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  z-index: 15;
  backdrop-filter: blur(2px);
}

/* ---- Page transition ---- */
.page-enter-active { animation: fadeInUp 0.3s var(--ease-out); }
.page-leave-active { animation: fadeIn 0.12s ease reverse; }
.fade-enter-active { animation: fadeIn 0.2s ease; }
.fade-leave-active { animation: fadeIn 0.15s ease reverse; }

/* ---- Mobile ---- */
@media (max-width: 768px) {
  .mobile-menu-btn { display: flex; }

  .sidebar {
    position: fixed;
    top: 0; left: 0;
    height: 100vh;
    transform: translateX(-100%);
    transition: transform 0.3s var(--ease-out);
    z-index: 25;
  }
  .sidebar.open { transform: translateX(0); }
  .sidebar-backdrop { display: block; }

  .topbar { padding: 0 16px; }
}
</style>
