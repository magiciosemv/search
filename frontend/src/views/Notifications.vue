<template>
  <div class="p-6">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold">Notifications</h1>
      <button @click="showAddModal = true" class="btn-primary">Add Channel</button>
    </div>

    <div class="grid gap-4">
      <div v-for="notif in notifications" :key="notif.id" class="card p-5">
        <div class="flex items-center justify-between mb-3">
          <div class="flex items-center gap-3">
            <span class="w-8 h-8 rounded-full flex items-center justify-center text-white text-xs font-bold"
              :class="notif.type === 'telegram' ? 'bg-blue-500' : 'bg-amber-500'">
              {{ notif.type === 'telegram' ? 'TG' : 'EM' }}
            </span>
            <div>
              <div class="font-semibold text-lg">{{ notif.name }}</div>
              <div class="text-sm text-gray-500">{{ notif.type === 'telegram' ? 'Telegram' : 'Email' }}</div>
            </div>
          </div>
          <span :class="notif.enabled ? 'badge badge-success' : 'badge badge-muted'"
            class="cursor-pointer select-none" @click="toggleNotification(notif)">
            {{ notif.enabled ? 'Enabled' : 'Disabled' }}
          </span>
        </div>

        <div class="bg-gray-50 dark:bg-slate-700 rounded-lg p-3 text-sm">
          <span class="text-gray-500">{{ notif.type === 'telegram' ? 'Chat ID:' : 'Email:' }}</span>
          <span class="font-mono ml-2">{{ getConfigValue(notif) }}</span>
        </div>

        <div class="flex justify-end gap-2 mt-4">
          <button @click="testNotification(notif.id)"
            class="px-3 py-1.5 text-sm border border-purple-300 text-purple-600 rounded-lg hover:bg-purple-50">
            Test
          </button>
          <button @click="deleteNotification(notif.id)"
            class="px-3 py-1.5 text-sm border border-red-300 text-red-500 rounded-lg hover:bg-red-50">
            Delete
          </button>
        </div>
      </div>

      <EmptyState v-if="!notifications || notifications.length === 0"
        title="No notification channels configured"
        message="Add a Telegram or Email channel to receive alerts" />
    </div>

    <Modal v-model="showAddModal" title="Add Channel" @submit="addNotification">
      <div>
        <label class="block text-sm font-medium mb-1">Channel Name</label>
        <input v-model="newNotif.name" type="text" placeholder="e.g., My Telegram"
          class="w-full px-3 py-2 border rounded-lg dark:bg-slate-700 dark:border-slate-600" />
      </div>
      <div>
        <label class="block text-sm font-medium mb-2">Type</label>
        <div class="flex gap-2">
          <button @click="newNotif.type = 'telegram'"
            :class="newNotif.type === 'telegram' ? 'bg-purple-600 text-white' : 'bg-gray-100 text-gray-700 dark:bg-slate-700 dark:text-gray-300'"
            class="flex-1 px-4 py-2 rounded-lg text-sm font-medium transition">
            Telegram
          </button>
          <button @click="newNotif.type = 'email'"
            :class="newNotif.type === 'email' ? 'bg-purple-600 text-white' : 'bg-gray-100 text-gray-700 dark:bg-slate-700 dark:text-gray-300'"
            class="flex-1 px-4 py-2 rounded-lg text-sm font-medium transition">
            Email
          </button>
        </div>
      </div>
      <div v-if="newNotif.type === 'telegram'">
        <label class="block text-sm font-medium mb-1">Telegram Chat ID</label>
        <input v-model="newNotif.config.chat_id" type="text" placeholder="Your Telegram Chat ID"
          class="w-full px-3 py-2 border rounded-lg dark:bg-slate-700 dark:border-slate-600" />
      </div>
      <div v-if="newNotif.type === 'email'">
        <label class="block text-sm font-medium mb-1">Email Address</label>
        <input v-model="newNotif.config.email" type="email" placeholder="your@email.com"
          class="w-full px-3 py-2 border rounded-lg dark:bg-slate-700 dark:border-slate-600" />
      </div>
    </Modal>
  </div>
</template>

<script setup>
import { apiPost, apiPut, apiDelete, useFetch } from '../utils/api.js'
import Modal from '../components/Modal.vue'
import EmptyState from '../components/EmptyState.vue'

const { data: notifications } = useFetch('/api/notifications', [])
const showAddModal = ref(false)
const newNotif = ref({ name: '', type: 'telegram', config: { chat_id: '', email: '' } })

const getConfigValue = (notif) => {
  if (notif.type === 'telegram') return notif.config_map?.chat_id || ''
  if (notif.type === 'email') return notif.config_map?.email || ''
  return ''
}

const addNotification = async () => {
  try {
    const config = newNotif.value.type === 'telegram'
      ? { chat_id: newNotif.value.config.chat_id }
      : { email: newNotif.value.config.email }
    await apiPost('/api/notifications', { name: newNotif.value.name, type: newNotif.value.type, config })
    showAddModal.value = false
    newNotif.value = { name: '', type: 'telegram', config: { chat_id: '', email: '' } }
  } catch (e) {
    alert(e.message)
  }
}

const toggleNotification = async (notif) => {
  try {
    await apiPut(`/api/notifications/${notif.id}`, {
      name: notif.name, type: notif.type,
      config: notif.config_map || {}, enabled: !notif.enabled
    })
    notifications.value = await fetch('/api/notifications').then(r => r.json())
  } catch (e) {
    alert(e.message)
  }
}

const testNotification = async (id) => {
  try {
    await apiPost(`/api/notifications/${id}/test`, {})
    alert('Test notification sent!')
  } catch (e) {
    alert(e.message)
  }
}

const deleteNotification = async (id) => {
  if (!confirm('Are you sure?')) return
  try {
    await apiDelete(`/api/notifications/${id}`)
    notifications.value = await fetch('/api/notifications').then(r => r.json())
  } catch (e) {
    alert(e.message)
  }
}
</script>