import { ref, reactive } from 'vue'

const locale = ref(localStorage.getItem('locale') || 'en')
const messages = reactive({})
const loaded = ref(0)

async function loadMessages(lang) {
  if (messages[lang]) return
  const m = await import(`./locales/${lang}.js`)
  messages[lang] = m.default
  loaded.value++
}

// Load initial locale synchronously if possible, or async
loadMessages(locale.value)

function t(key) {
  const keys = key.split('.')
  let result = messages[locale.value]
  if (!result) return key
  for (const k of keys) {
    result = result[k]
    if (result === undefined) return key
  }
  return result
}

async function setLocale(lang) {
  await loadMessages(lang)
  locale.value = lang
  localStorage.setItem('locale', lang)
  document.documentElement.setAttribute('lang', lang)
}

export function useI18n() {
  return { locale, t, setLocale }
}
