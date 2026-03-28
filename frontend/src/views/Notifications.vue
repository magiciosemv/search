<template>
  <div class="p-6">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold">Notifications</h1>
      <button @click="showAddModal = true"
        class="px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700">
        Add Notification
      </button>
    </div>

    <!-- Notifications List -->
    <div class="grid gap-4">
      <div v-for="notif in notifications" :key="notif.id"
        class="bg-white dark:bg-slate-800 p-4 rounded-lg shadow">
        <div class="flex items-center justify-between">
          <div>
            <div class="font-semibold">{{ notif.name }}</div>
            <div class="text-sm text-gray-500">{{ notif.type }} - {{ getConfigValue(notif) }}</div>
          </div>
          <div class="flex items-center gap-3">
            <span :class="notif.enabled ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'"
              class="px-2 py-1 rounded text-xs">
              {{ notif.enabled ? 'Enabled' : 'Disabled' }}
            </span>
            <button @click="testNotification(notif.id)" class="text-purple-600 hover:text-purple-800 text-sm">
              Test
            </button>
            <button @click="deleteNotification(notif.id)" class="text-red-600 hover:text-red-800 text-sm">
              Delete
            </button>
          </div>
        </div>
      </div>
      <div v-if="notifications.length === 0" class="bg-white dark:bg-slate-800 p-8 rounded-lg shadow text-center text-gray-500">
        No notifications configured yet
      </div>
    </div>

    <!-- Add Notification Modal -->
    <div v-if="showAddModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white dark:bg-slate-800 p-6 rounded-lg w-full max-w-md">
        <h2 class="text-xl font-bold mb-4">Add Notification</h2>
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium mb-1">Name</label>
            <input v-model="newNotif.name" type="text" placeholder="e.g., My Telegram"
              class="w-full px-3 py-2 border rounded-lg dark:bg-slate-700 dark:border-slate-600" />
          </div>
          <div>
            <label class="block text-sm font-medium mb-1">Type</label>
            <select v-model="newNotif.type"
              class="w-full px-3 py-2 border rounded-lg dark:bg-slate-700 dark:border-slate-600">
              <option value="telegram">Telegram</option>
              <option value="email">Email</option>
            </select>
          </div>
          <div v-if="newNotif.type === 'telegram'">
            <label class="block text-sm font-medium mb-1">Chat ID</label>
            <input v-model="newNotif.config.chat_id" type="text" placeholder="Your Telegram Chat ID"
              class="w-full px-3 py-2 border rounded-lg dark:bg-slate-700 dark:border-slate-600" />
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
    notifications.value = await res.json()
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

const testNotification = async (id) => {
  try {
    const res = await fetch(`/api/notifications/${id}/test`, { method: 'POST' })
    if (res.ok) {
      alert('Test notification sent!')
    } else {
      alert('Failed to send test notification')
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