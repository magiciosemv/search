<template>
  <div class="page-container">
    <div class="page-header">
      <div>
        <h1 class="page-title">Notifications</h1>
        <p class="page-subtitle">Alert delivery channels</p>
      </div>
      <button class="btn btn-primary" @click="showAddModal = true">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
        Add Channel
      </button>
    </div>

    <div class="notif-grid stagger">
      <div v-for="notif in notifications" :key="notif.id"
        class="notif-card card animate-fade-up"
        :class="notif.type">
        <div class="notif-top">
          <div class="notif-identity">
            <div class="notif-type-icon" :class="notif.type">
              <svg v-if="notif.type === 'telegram'" width="20" height="20" viewBox="0 0 24 24" fill="currentColor"><path d="M11.944 0A12 12 0 000 12a12 12 0 0012 12 12 12 0 0012-12A12 12 0 0012 0a12 12 0 00-.056 0zm4.962 7.224c.1-.002.321.023.465.14a.506.506 0 01.171.325c.016.093.036.306.02.472-.18 1.898-.962 6.502-1.36 8.627-.168.9-.499 1.201-.82 1.23-.696.065-1.225-.46-1.9-.902-1.056-.693-1.653-1.124-2.678-1.8-1.185-.78-.417-1.21.258-1.91.177-.184 3.247-2.977 3.307-3.23.007-.032.014-.15-.056-.212s-.174-.041-.249-.024c-.106.024-1.793 1.14-5.061 3.345-.479.33-.913.49-1.302.48-.428-.008-1.252-.241-1.865-.44-.752-.245-1.349-.374-1.297-.789.027-.216.325-.437.893-.663 3.498-1.524 5.83-2.529 6.998-3.014 3.332-1.386 4.025-1.627 4.476-1.635z"/></svg>
              <svg v-else width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><rect x="2" y="4" width="20" height="16" rx="2"/><path d="M22 4L12 13 2 4"/></svg>
            </div>
            <div>
              <div class="notif-name">{{ notif.name }}</div>
              <div class="notif-type-label">{{ notif.type === 'telegram' ? 'Telegram Bot' : 'Email SMTP' }}</div>
            </div>
          </div>
          <button class="toggle-btn" :class="{ on: notif.enabled }" @click="toggleNotification(notif)">
            <div class="toggle-track">
              <div class="toggle-thumb"></div>
            </div>
            <span class="toggle-text">{{ notif.enabled ? 'Active' : 'Off' }}</span>
          </button>
        </div>

        <div class="notif-config">
          <span class="config-key">{{ notif.type === 'telegram' ? 'Chat ID' : 'Address' }}</span>
          <span class="config-value font-mono">{{ getConfigValue(notif) }}</span>
        </div>

        <div class="notif-actions">
          <button @click="testNotification(notif.id)" class="btn btn-ghost btn-sm">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M18 8A6 6 0 006 8c0 7-3 9-3 9h18s-3-2-3-9"/><path d="M13.73 21a2 2 0 01-3.46 0"/></svg>
            Test
          </button>
          <button @click="deleteNotification(notif.id)" class="btn btn-danger btn-sm">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2"/></svg>
            Delete
          </button>
        </div>
      </div>
    </div>

    <EmptyState v-if="!notifications || notifications.length === 0"
      title="No channels configured"
      message="Add a Telegram bot or email to receive alerts" />

    <Modal v-model="showAddModal" title="Add Channel" submit-label="Add" @submit="addNotification">
      <div>
        <label class="field-label">Channel Name</label>
        <input v-model="newNotif.name" type="text" placeholder="e.g. My Telegram"
          class="input-field" />
      </div>
      <div>
        <label class="field-label">Type</label>
        <div class="type-selector">
          <button @click="newNotif.type = 'telegram'" class="type-option"
            :class="{ active: newNotif.type === 'telegram' }">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor"><path d="M11.944 0A12 12 0 000 12a12 12 0 0012 12 12 12 0 0012-12A12 12 0 0012 0a12 12 0 00-.056 0zm4.962 7.224c.1-.002.321.023.465.14a.506.506 0 01.171.325c.016.093.036.306.02.472-.18 1.898-.962 6.502-1.36 8.627-.168.9-.499 1.201-.82 1.23-.696.065-1.225-.46-1.9-.902-1.056-.693-1.653-1.124-2.678-1.8-1.185-.78-.417-1.21.258-1.91.177-.184 3.247-2.977 3.307-3.23.007-.032.014-.15-.056-.212s-.174-.041-.249-.024c-.106.024-1.793 1.14-5.061 3.345-.479.33-.913.49-1.302.48-.428-.008-1.252-.241-1.865-.44-.752-.245-1.349-.374-1.297-.789.027-.216.325-.437.893-.663 3.498-1.524 5.83-2.529 6.998-3.014 3.332-1.386 4.025-1.627 4.476-1.635z"/></svg>
            Telegram
          </button>
          <button @click="newNotif.type = 'email'" class="type-option"
            :class="{ active: newNotif.type === 'email' }">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><rect x="2" y="4" width="20" height="16" rx="2"/><path d="M22 4L12 13 2 4"/></svg>
            Email
          </button>
        </div>
      </div>
      <div v-if="newNotif.type === 'telegram'">
        <label class="field-label">Telegram Chat ID</label>
        <input v-model="newNotif.config.chat_id" type="text" placeholder="Enter Chat ID..."
          class="input-field input-mono" />
      </div>
      <div v-if="newNotif.type === 'email'">
        <label class="field-label">Email Address</label>
        <input v-model="newNotif.config.email" type="email" placeholder="your@email.com"
          class="input-field input-mono" />
      </div>
    </Modal>
  </div>
</template>

<script setup>
import { ref } from 'vue'
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
    notifications.value = await fetch('/api/notifications').then(r => r.json())
  } catch (e) { alert(e.message) }
}

const toggleNotification = async (notif) => {
  try {
    await apiPut(`/api/notifications/${notif.id}`, {
      name: notif.name, type: notif.type,
      config: notif.config_map || {}, enabled: !notif.enabled
    })
    notifications.value = await fetch('/api/notifications').then(r => r.json())
  } catch (e) { alert(e.message) }
}

const testNotification = async (id) => {
  try {
    await apiPost(`/api/notifications/${id}/test`, {})
    alert('Test notification sent!')
  } catch (e) { alert(e.message) }
}

const deleteNotification = async (id) => {
  if (!confirm('Remove this channel?')) return
  try {
    await apiDelete(`/api/notifications/${id}`)
    notifications.value = await fetch('/api/notifications').then(r => r.json())
  } catch (e) { alert(e.message) }
}
</script>

<style scoped>
.notif-grid {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.notif-card {
  padding: 20px 24px;
}

.notif-card.telegram { border-left: 3px solid #229ed9; }
.notif-card.email { border-left: 3px solid var(--amber-base); }

.notif-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.notif-identity {
  display: flex;
  align-items: center;
  gap: 14px;
}

.notif-type-icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 10px;
}

.notif-type-icon.telegram {
  background: rgba(34, 158, 217, 0.12);
  color: #229ed9;
}

.notif-type-icon.email {
  background: rgba(245, 158, 11, 0.12);
  color: var(--amber-bright);
}

.notif-name {
  font-size: 1rem;
  font-weight: 600;
}

.notif-type-label {
  font-size: 0.8125rem;
  color: var(--text-muted);
  margin-top: 1px;
}

/* Toggle */
.toggle-btn {
  display: flex;
  align-items: center;
  gap: 10px;
  background: none;
  border: none;
  cursor: pointer;
  padding: 0;
}

.toggle-track {
  width: 40px;
  height: 22px;
  border-radius: 11px;
  background: var(--bg-elevated);
  border: 1px solid var(--border-subtle);
  position: relative;
  transition: all 0.25s ease;
}

.toggle-thumb {
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: var(--text-muted);
  position: absolute;
  top: 2px;
  left: 2px;
  transition: all 0.25s var(--ease-out);
}

.toggle-btn.on .toggle-track {
  background: rgba(16, 185, 129, 0.15);
  border-color: rgba(16, 185, 129, 0.3);
}

.toggle-btn.on .toggle-thumb {
  left: 20px;
  background: var(--green-bright);
  box-shadow: 0 0 8px rgba(16, 185, 129, 0.5);
}

.toggle-text {
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--text-muted);
  font-family: var(--font-mono);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.toggle-btn.on .toggle-text { color: var(--green-bright); }

/* Config */
.notif-config {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 14px;
  background: var(--bg-elevated);
  border-radius: 8px;
  margin-bottom: 16px;
}

.config-key {
  font-size: 0.75rem;
  color: var(--text-muted);
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.config-value {
  font-size: 0.8125rem;
  color: var(--text-secondary);
}

.notif-actions {
  display: flex;
  gap: 8px;
}

/* Type selector */
.type-selector {
  display: flex;
  gap: 8px;
}

.type-option {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 10px;
  border-radius: 8px;
  border: 1px solid var(--border-subtle);
  background: var(--bg-elevated);
  color: var(--text-muted);
  font-family: var(--font-display);
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.type-option:hover {
  cursor: pointer;
  border-color: var(--border-default);
  color: var(--text-secondary);
}

.type-option.active {
  border-color: var(--cyan-base);
  background: rgba(6, 182, 212, 0.08);
  color: var(--cyan-bright);
}

</style>
