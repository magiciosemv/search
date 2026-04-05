import { ref, onMounted } from 'vue'

export async function apiPost(url, body) {
  const res = await fetch(url, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body)
  })
  if (!res.ok) {
    const err = await res.json()
    throw new Error(err.error || 'Request failed')
  }
  return res.json()
}

export async function apiPut(url, body) {
  const res = await fetch(url, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body)
  })
  if (!res.ok) {
    const err = await res.json()
    throw new Error(err.error || 'Request failed')
  }
  return res.json()
}

export async function apiDelete(url) {
  const res = await fetch(url, { method: 'DELETE' })
  if (!res.ok) throw new Error('Request failed')
}

export function useFetch(url, defaultValue = []) {
  const data = ref(defaultValue)
  const loading = ref(false)

  const execute = async (overrideUrl) => {
    loading.value = true
    try {
      const res = await fetch(overrideUrl || url)
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