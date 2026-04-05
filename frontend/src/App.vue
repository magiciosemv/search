<template>
  <div class="min-h-screen bg-gray-100 dark:bg-slate-900">
    <div class="flex">
      <aside class="w-64 bg-white dark:bg-slate-800 min-h-screen border-r dark:border-slate-700">
        <div class="p-4 border-b dark:border-slate-700">
          <h1 class="text-xl font-bold text-purple-600">Solana Monitor</h1>
        </div>
        <nav class="p-4">
          <router-link v-for="item in navItems" :key="item.path" :to="item.path"
            class="flex items-center gap-2 px-4 py-2.5 rounded-lg mb-1 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-slate-700 transition"
            :class="{ 'bg-purple-100 dark:bg-purple-900/50 text-purple-700 dark:text-purple-300': $route.path === item.path }">
            <span class="text-sm opacity-60">{{ item.icon }}</span>
            <span class="text-sm">{{ item.name }}</span>
          </router-link>
        </nav>
      </aside>
      <main class="flex-1 overflow-auto">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup>
const navItems = [
  { path: '/', name: 'Dashboard', icon: '◉' },
  { path: '/addresses', name: 'Wallets', icon: '◎' },
  { path: '/notifications', name: 'Notifications', icon: '◈' },
  { path: '/alerts', name: 'Alerts', icon: '◇' },
]
</script>