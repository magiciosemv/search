<template>
  <div class="p-6">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold">Notifications</h1>
      <button @click="showAddModal = true"
        class="px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700">
        Add Channel
      </button>
    </div>

    <!-- Notification Channels -->
    <div class="grid gap-4">
      <div v-for="notif in notifications" :key="notif.id"
        class="bg-white dark:bg-slate-800 p-5 rounded-lg shadow">
        <div class="flex items-center justify-between mb-3">
          <div class="flex items-center gap-3">
            <span class="text-2xl">{{ notif.type === 'telegram' ? '&#9993;' : '&#9993;' }}</span>
            <div>
              <div class="font-semibold text-lg">{{ notif.name }}</div>
              <div class="text-sm text-gray-500">{{ notif.type === 'telegram' ? 'Telegram' : 'Email' }}</div>
            </div>
          </div>
          <span :class="notif.enabled ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'"
            class="px-3 py-1 rounded-full text-xs font-medium cursor-pointer"
            @click="toggleNotification(notif)">
            {{ notif.enabled ? 'Enabled' : 'Disabled' }}
          </span>
        </div>

        <div class="bg-gray-50 dark:bg-slate-700 rounded-lg p-3 text-sm">
          <div v-if="notif.type === 'telegram'">
            <span class="text-gray-500">Chat ID:</span>
            <span class="font-mono ml-2">{{ getConfigValue(notif) }}</span>
          </div>
          <div v-else>
            <span class="text-gray-500">Email:</span>
            <span class="ml-2">{{ getConfigValue(notif) }}</span>
          </div>
        </div>

        <div class="flex justify-end gap-3 mt-4">
          <button @click="testNotification(notif.id)"
            class="px-3 py-1.5 text-sm border border-purple-600 text-purple-600 rounded-lg hover:bg-purple-50">
            Test
          </button>
          <button @click="deleteNotification(notif.id)"
            class="px-3 py-1.5 text-sm border border-red-400 text-red-600 rounded-lg hover:bg-red-50">
            Delete
          </button>
        </div>
      </div>

      <div v-if="!notifications || notifications.length === 0"
        class="bg-white dark:bg-slate-800 p-12 rounded-lg shadow text-center">
        <div class="text-gray-400 mb-2">No notification channels configured</div>
        <div class="text-gray-400 text-sm">Add a Telegram or Email channel to receive alerts</div>
      </div>
    </div>

    <!-- Add Notification Modal -->
    <div v-if="showAddModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white dark:bg-slate-800 p-6 rounded-lg w-full max-w-md">
        <h2 class="text-xl font-bold mb-4">Add Notification Channel</h2>
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium mb-1">Channel Name</label>
            <input v-model="newNotif.name" type="text" placeholder="e.g., My Telegram"
              class="w-full px-3 py-2 border rounded-lg dark:bg-slate-700 dark:border-slate-600" />
          </div>
          <div>
            <label class="block text-sm font-medium mb-1">Type</label>
            <div class="flex gap-2">
              <button @click="newNotif.type = 'telegram'"
                :class="newNotif.type === 'telegram' ? 'bg-purple-600 text-white' : 'bg-gray-100 text-gray-700 dark:bg-slate-700 dark:text-gray-300'"
                class="flex-1 px-4 py-2 rounded-lg text-sm font-medium">
                Telegram
              </button>
              <button @click="newNotif.type = 'email'"
                :class="newNotif.type === 'email' ? 'bg-purple-600 text-white' : 'bg-gray-100 text-gray-700 dark:bg-slate-700 dark:text-gray-300'"
                class="flex-1 px-4 py-2 rounded-lg text-sm font-medium">
                Email
              </button>
            </div>
          </div>
          <div v-if="newNotif.type === 'telegram'">
            <label class="block text-sm font-medium mb-1">Telegram Chat ID</label>
            <input v-model="newNotif.config.chat_id" type="text" placeholder="Enter your Telegram Chat ID"
              class="w-full px-3 py-2 border rounded-lg dark:bg-slate-700 dark:border-slate-600" />
            <p class="text-xs text-gray-400 mt-1">Send /start to @BotFather, then message your bot to get Chat ID</p>
          </div>
          <div v-if="newNotif.type === 'email'">
            <label class="block text-sm font-medium mb-1">Email Address</label>
            <input v-model="newNotif.config.email" type="email" placeholder="your@email.com"
              class="w-full px-3 py-2 border rounded-lg dark:bg-slate-700 dark:border-slate-600" />
          </div>
        </div>
        <div class="flex justify-end gap-2 mt-6">
          <button @click="showAddModal = false" class="px-4 py-2 border rounded-lg hover:bg-gray-100 dark:hover:bg-slate-700">
            Cancel
          </button>
          <button @click="addNotification" class="px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700">
            Add
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const notifications = ref([])
const showAddModal = ref(false)
const newNotif = ref({
  name: '',
  type: 'telegram',
  config: { chat_id: '', email: '' }
})

const fetchNotifications = async () => {
  try {
    const res = await fetch('/api/notifications')
    const data = await res.json()
    notifications.value = data || []
  } catch (e) {
    console.error('Failed to fetch notifications:', e)
  }
}

const getConfigValue = (notif) => {
  if (notif.type === 'telegram') return notif.config_map?.chat_id || ''
  if (notif.type === 'email') return notif.config_map?.email || ''
  return ''
}

const addNotification = async () => {
  const config = {}
  if (newNotif.value.type === 'telegram') {
    config.chat_id = newNotif.value.config.chat_id
  } else {
    config.email = newNotif.value.config.email
  }

  try {
    const res = await fetch('/api/notifications', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        name: newNotif.value.name,
        type: newNotif.value.type,
        config
      })
    })
    if (!res.ok) {
      const err = await res.json()
      alert(err.error || 'Failed to add notification')
      return
    }
    showAddModal.value = false
    newNotif.value = { name: '', type: 'telegram', config: { chat_id: '', email: '' } }
    fetchNotifications()
  } catch (e) {
    alert('Failed to add notification')
  }
}

const toggleNotification = async (notif) => {
  try {
    await fetch(`/api/notifications/${notif.id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        name: notif.name,
        type: notif.type,
        config: notif.config_map || {},
        enabled: !notif.enabled
      })
    })
    fetchNotifications()
  } catch (e) {
    alert('Failed to toggle notification')
  }
}

const testNotification = async (id) => {
  try {
    const res = await fetch(`/api/notifications/${id}/test`, { method: 'POST' })
    if (res.ok) {
      alert('Test notification sent!')
    } else {
      const err = await res.json()
      alert(err.error || 'Failed to send test notification')
    }
  } catch (e) {
    alert('Failed to send test notification')
  }
}

const deleteNotification = async (id) => {
  if (!confirm('Are you sure?')) return
  try {
    await fetch(`/api/notifications/${id}`, { method: 'DELETE' })
    fetchNotifications()
  } catch (e) {
    alert('Failed to delete notification')
  }
}

onMounted(() => {
  fetchNotifications()
})
</script>