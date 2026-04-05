import { ref, onMounted } from 'vue'

const API_KEY = 'solana-monitor-secret-key-2024'

function authHeaders() {
  const h = {}
  if (API_KEY) {
    h['Authorization'] = `Bearer ${API_KEY}`
  }
  return h
}

export async function apiPost(url, body) {
  const res = await fetch(url, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...authHeaders() },
    body: JSON.stringify(body)
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({}))
    throw new Error(err.error || 'Request failed')
  }
  return res.json()
}

export async function apiPut(url, body) {
  const res = await fetch(url, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', ...authHeaders() },
    body: JSON.stringify(body)
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({}))
    throw new Error(err.error || 'Request failed')
  }
  return res.json()
}

export async function apiDelete(url) {
  const res = await fetch(url, {
    method: 'DELETE',
    headers: authHeaders()
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({}))
    throw new Error(err.error || 'Request failed')
  }
}

export function useFetch(url, defaultValue = []) {
  const data = ref(defaultValue)
  const loading = ref(false)

  const execute = async (overrideUrl) => {
    loading.value = true
    try {
      const res = await fetch(overrideUrl || url, {
        headers: authHeaders()
      })
      const json = await res.json()
      data.value = json || defaultValue
    } catch (e) {
      console.error('Fetch failed:', e)
    } finally {
      loading.value = false
    }
  }

  onMounted(execute)

  return { data, loading, refetch: execute }
}

// useSSE subscribes to Server-Sent Events and calls eventHandler(eventType, data).
// Returns a reconnect function.
export function useSSE(url, eventHandler) {
  let source = null
  let reconnectTimer = null
  let active = true

  function connect() {
    if (!active) return

    source = new EventSource(url)

    source.addEventListener('connected', () => {
      console.log('SSE connected')
    })

    source.addEventListener('balance_update', (e) => {
      try {
        const data = JSON.parse(e.data)
        eventHandler('balance_update', data)
      } catch (err) {
        console.error('SSE parse error:', err)
      }
    })

    source.addEventListener('new_alert', (e) => {
      try {
        const data = JSON.parse(e.data)
        eventHandler('new_alert', data)
      } catch (err) {
        console.error('SSE parse error:', err)
      }
    })

    source.onerror = () => {
      source.close()
      source = null
      if (active) {
        reconnectTimer = setTimeout(connect, 3000)
      }
    }
  }

  connect()

  function disconnect() {
    active = false
    if (reconnectTimer) clearTimeout(reconnectTimer)
    if (source) source.close()
  }

  return { disconnect }
}
