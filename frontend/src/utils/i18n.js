import { ref, computed } from 'vue'

const locale = ref(localStorage.getItem('locale') || 'en')

const messages = {}

function loadMessages(lang) {
  if (messages[lang]) return
  import(`./locales/${lang}.js`).then(m => {
    messages[lang] = m.default
  })
}

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

function setLocale(lang) {
  locale.value = lang
  localStorage.setItem('locale', lang)
  if (!messages[lang]) {
    loadMessages(lang)
  }
  document.documentElement.setAttribute('lang', lang)
}

export function useI18n() {
  return { locale, t, setLocale }
}
