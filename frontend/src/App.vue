<template>
  <div class="min-h-screen bg-gray-100 dark:bg-slate-900">
    <!-- Sidebar -->
    <div class="flex">
      <aside class="w-64 bg-white dark:bg-slate-800 min-h-screen border-r dark:border-slate-700">
        <div class="p-4 border-b dark:border-slate-700">
          <h1 class="text-xl font-bold text-purple-600">Solana Monitor</h1>
        </div>
        <nav class="p-4">
          <router-link v-for="item in navItems" :key="item.path" :to="item.path"
            class="flex items-center gap-2 px-4 py-2 rounded-lg mb-1 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-slate-700"
            :class="{ 'bg-purple-100 dark:bg-purple-900 text-purple-700 dark:text-purple-300': $route.path === item.path }">
            <span>{{ item.icon }}</span>
            {{ item.name }}
          </router-link>
        </nav>
      </aside>

      <!-- Main Content -->
      <main class="flex-1">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup>
const navItems = [
  { path: '/', name: 'Dashboard', icon: '📊' },
  { path: '/addresses', name: 'Addresses', icon: '📍' },
  { path: '/rules', name: 'Rules', icon: '📋' },
  { path: '/notifications', name: 'Notifications', icon: '🔔' },
  { path: '/alerts', name: 'Alerts', icon: '⚠️' },
]
</script>