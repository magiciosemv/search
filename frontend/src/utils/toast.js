import { ref } from 'vue'

const toasts = ref([])
let nextId = 0

function addToast(type, message, duration = 3500) {
  const id = nextId++
  toasts.value.push({ id, type, message })
  if (duration > 0) {
    setTimeout(() => removeToast(id), duration)
  }
}

function removeToast(id) {
  const idx = toasts.value.findIndex(t => t.id === id)
  if (idx !== -1) toasts.value.splice(idx, 1)
}

export function useToast() {
  return {
    toasts,
    success: (msg) => addToast('success', msg),
    error: (msg) => addToast('error', msg, 5000),
    info: (msg) => addToast('info', msg),
    remove: removeToast
  }
}
