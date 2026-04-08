import { ref } from 'vue'

const visible = ref(false)
const title = ref('')
const message = ref('')
let resolvePromise = null

function showConfirm(msg, ttl = 'Confirm') {
  title.value = ttl
  message.value = msg
  visible.value = true
  return new Promise((resolve) => {
    resolvePromise = resolve
  })
}

function ok() {
  visible.value = false
  if (resolvePromise) resolvePromise(true)
  resolvePromise = null
}

function cancel() {
  visible.value = false
  if (resolvePromise) resolvePromise(false)
  resolvePromise = null
}

export function useConfirm() {
  return { visible, title, message, show: showConfirm, ok, cancel }
}
