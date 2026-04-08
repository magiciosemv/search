<template>
  <div class="page-container">
    <div class="page-header">
      <div>
        <h1 class="page-title">{{ t('channels.title') }}</h1>
        <p class="page-subtitle">{{ t('channels.subtitle') }}</p>
      </div>
      <button class="btn btn-primary" @click="showAddModal = true">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
        {{ t('channels.addChannel') }}
      </button>
    </div>

    <div class="notif-list stagger">
      <div v-for="notif in notifications" :key="notif.id"
        class="notif-card glass animate-fade-up" :class="notif.type">
        <div class="notif-top">
          <div class="notif-identity">
            <div class="notif-icon" :class="notif.type">
              <svg v-if="notif.type === 'telegram'" width="18" height="18" viewBox="0 0 24 24" fill="currentColor"><path d="M11.944 0A12 12 0 000 12a12 12 0 0012 12 12 12 0 0012-12A12 12 0 0012 0a12 12 0 00-.056 0zm4.962 7.224c.1-.002.321.023.465.14a.506.506 0 01.171.325c.016.093.036.306.02.472-.18 1.898-.962 6.502-1.36 8.627-.168.9-.499 1.201-.82 1.23-.696.065-1.225-.46-1.9-.902-1.056-.693-1.653-1.124-2.678-1.8-1.185-.78-.417-1.21.258-1.91.177-.184 3.247-2.977 3.307-3.23.007-.032.014-.15-.056-.212s-.174-.041-.249-.024c-.106.024-1.793 1.14-5.061 3.345-.479.33-.913.49-1.302.48-.428-.008-1.252-.241-1.865-.44-.752-.245-1.349-.374-1.297-.789.027-.216.325-.437.893-.663 3.498-1.524 5.83-2.529 6.998-3.014 3.332-1.386 4.025-1.627 4.476-1.635z"/></svg>
              <svg v-else width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><rect x="2" y="4" width="20" height="16" rx="2"/><path d="M22 4L12 13 2 4"/></svg>
            </div>
            <div>
              <div class="notif-name">{{ notif.name }}</div>
              <div class="notif-type text-muted">{{ notif.type === 'telegram' ? 'Telegram' : 'Email' }}</div>
            </div>
          </div>
          <button class="toggle-btn" :class="{ on: notif.enabled }" @click="toggleNotification(notif)">
            <div class="toggle-track"><div class="toggle-thumb"></div></div>
            <span class="toggle-text">{{ notif.enabled ? 'On' : 'Off' }}</span>
          </button>
        </div>

        <div class="notif-config">
          <span class="config-key text-dim">{{ notif.type === 'telegram' ? t('channels.chatIdLabel') : t('channels.addressLabel') }}</span>
          <span class="config-value font-mono text-secondary">{{ getConfigValue(notif) }}</span>
        </div>

        <div class="notif-actions">
          <button @click="testNotification(notif.id)" class="btn btn-ghost btn-sm">{{ t('channels.test') }}</button>
          <button @click="deleteNotification(notif.id)" class="btn btn-danger btn-sm">{{ t('channels.remove') }}</button>
        </div>
      </div>
    </div>

    <EmptyState v-if="!notifications || notifications.length === 0" :title="t('channels.noChannels')" :message="t('channels.noChannelsMsg')" />

    <Modal v-model="showAddModal" :title="t('channels.addChannel')" :submit-label="t('common.add')" @submit="addNotification">
      <div>
        <label class="field-label">{{ t('channels.channelName') }}</label>
        <input v-model="newNotif.name" type="text" placeholder="e.g. My Telegram" class="input-field" />
      </div>
      <div>
        <label class="field-label">{{ t('channels.type') }}</label>
        <div class="type-selector">
          <button @click="newNotif.type = 'telegram'" class="type-btn" :class="{ active: newNotif.type === 'telegram' }">Telegram</button>
          <button @click="newNotif.type = 'email'" class="type-btn" :class="{ active: newNotif.type === 'email' }">Email</button>
        </div>
      </div>
      <div v-if="newNotif.type === 'telegram'">
        <label class="field-label">{{ t('channels.chatId') }}</label>
        <input v-model="newNotif.config.chat_id" type="text" placeholder="Enter Chat ID..." class="input-field input-mono" />
      </div>
      <div v-if="newNotif.type === 'email'">
        <label class="field-label">{{ t('channels.emailAddress') }}</label>
        <input v-model="newNotif.config.email" type="email" placeholder="your@email.com" class="input-field input-mono" />
      </div>
    </Modal>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { apiPost, apiPut, apiDelete, useFetch } from '../utils/api.js'
import { useToast } from '../utils/toast.js'
import { useConfirm } from '../utils/confirm.js'
import { useI18n } from '../utils/i18n.js'
import Modal from '../components/Modal.vue'
import EmptyState from '../components/EmptyState.vue'

const { t } = useI18n()
const toast = useToast()
const confirm = useConfirm()

const { data: notifications } = useFetch('/api/notifications', [])
const showAddModal = ref(false)
const newNotif = ref({ name: '', type: 'telegram', config: { chat_id: '', email: '' } })

const getConfigValue = (notif) => notif.type === 'telegram' ? (notif.config_map?.chat_id || '') : (notif.config_map?.email || '')

const addNotification = async () => {
  try {
    const config = newNotif.value.type === 'telegram' ? { chat_id: newNotif.value.config.chat_id } : { email: newNotif.value.config.email }
    await apiPost('/api/notifications', { name: newNotif.value.name, type: newNotif.value.type, config })
    showAddModal.value = false
    newNotif.value = { name: '', type: 'telegram', config: { chat_id: '', email: '' } }
    notifications.value = await fetch('/api/notifications', { headers: { Authorization: 'Bearer solana-monitor-secret-key-2024' } }).then(r => r.json())
    toast.success(t('channels.channelAdded'))
  } catch (e) { toast.error(e.message) }
}

const toggleNotification = async (notif) => {
  try {
    await apiPut(`/api/notifications/${notif.id}`, { name: notif.name, type: notif.type, config: notif.config_map || {}, enabled: !notif.enabled })
    notifications.value = await fetch('/api/notifications', { headers: { Authorization: 'Bearer solana-monitor-secret-key-2024' } }).then(r => r.json())
    toast.success(t(notif.enabled ? 'channels.channelDisabled' : 'channels.channelEnabled'))
  } catch (e) { toast.error(e.message) }
}

const testNotification = async (id) => {
  try { await apiPost(`/api/notifications/${id}/test`, {}); toast.success(t('channels.testSent')) } catch (e) { toast.error(e.message) }
}

const deleteNotification = async (id) => {
  if (!await confirm.show(t('channels.removeChannel'), t('channels.deleteChannel'))) return
  try { await apiDelete(`/api/notifications/${id}`); notifications.value = await fetch('/api/notifications', { headers: { Authorization: 'Bearer solana-monitor-secret-key-2024' } }).then(r => r.json()); toast.success(t('channels.channelRemoved')) } catch (e) { toast.error(e.message) }
}
</script>

<style scoped>
.notif-list { display: flex; flex-direction: column; gap: 10px; }
.notif-card { padding: 18px 20px; }
.notif-card.telegram { border-left: 3px solid #229ED9; }
.notif-card.email { border-left: 3px solid var(--amber-base); }

.notif-top { display: flex; align-items: center; justify-content: space-between; margin-bottom: 14px; }
.notif-identity { display: flex; align-items: center; gap: 14px; }
.notif-icon {
  width: 38px; height: 38px; display: flex; align-items: center; justify-content: center; border-radius: 10px;
}
.notif-icon.telegram { background: rgba(34, 158, 217, 0.1); color: #229ED9; }
.notif-icon.email { background: rgba(245, 158, 11, 0.1); color: var(--amber-bright); }
.notif-name { font-size: 0.9375rem; font-weight: 600; }
.notif-type { font-size: 0.75rem; margin-top: 2px; }

/* Toggle */
.toggle-btn { display: flex; align-items: center; gap: 8px; background: none; border: none; cursor: pointer; }
.toggle-track {
  width: 38px; height: 22px; border-radius: 11px;
  background: var(--bg-elevated); border: 1px solid var(--border-subtle);
  position: relative; transition: all var(--duration-normal) ease;
}
.toggle-thumb {
  width: 16px; height: 16px; border-radius: 50%;
  background: var(--text-dim); position: absolute; top: 2px; left: 2px;
  transition: all var(--duration-normal) var(--ease-out);
}
.toggle-btn.on .toggle-track { background: rgba(0, 214, 143, 0.12); border-color: rgba(0, 214, 143, 0.3); }
.toggle-btn.on .toggle-thumb { left: 18px; background: var(--green-base); box-shadow: 0 0 8px rgba(0, 214, 143, 0.4); }
.toggle-text { font-size: 0.6875rem; font-family: var(--font-mono); font-weight: 600; color: var(--text-dim); letter-spacing: 0.04em; }
.toggle-btn.on .toggle-text { color: var(--green-base); }

.notif-config {
  display: flex; align-items: center; gap: 10px;
  padding: 9px 14px; background: var(--bg-void); border-radius: 8px; margin-bottom: 14px;
}
.config-key { font-size: 0.6875rem; font-weight: 600; }
.config-value { font-size: 0.8125rem; }
.notif-actions { display: flex; gap: 6px; }

/* Type selector */
.type-selector { display: flex; gap: 8px; }
.type-btn {
  flex: 1; padding: 10px; border-radius: 8px; font-size: 0.8125rem; font-weight: 600;
  border: 1px solid var(--border-subtle); background: var(--bg-void); color: var(--text-dim);
  font-family: var(--font-sans); cursor: pointer; transition: all var(--duration-fast) ease;
}
.type-btn:hover { border-color: var(--border-default); color: var(--text-secondary); }
.type-btn.active { border-color: var(--blue-base); background: var(--blue-dim); color: var(--blue-bright); }
</style>
