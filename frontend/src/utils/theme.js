import { ref } from 'vue'

const theme = ref(localStorage.getItem('theme') || 'dark')

function applyTheme(t) {
  document.documentElement.setAttribute('data-theme', t)
  localStorage.setItem('theme', t)
}

applyTheme(theme.value)

function toggleTheme() {
  theme.value = theme.value === 'dark' ? 'light' : 'dark'
  applyTheme(theme.value)
}

export function useTheme() {
  return { theme, toggleTheme }
}
