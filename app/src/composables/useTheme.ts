import { ref, watchEffect } from 'vue'
import vuetify from '@/plugins/vuetify'

type Theme = 'dark' | 'light'

const theme = ref<Theme>((localStorage.getItem('theme') as Theme) ?? 'dark')

const vars: Record<Theme, Record<string, string>> = {
  dark: {
    '--bg': '#0a0a0f',
    '--bg-card': '#0f0f1a',
    '--bg-input': 'rgba(255,255,255,0.05)',
    '--border': 'rgba(255,255,255,0.04)',
    '--text': '#f1f5f9',
    '--text-secondary': '#94a3b8',
  },
  light: {
    '--bg': '#f8fafc',
    '--bg-card': '#ffffff',
    '--bg-input': 'rgba(0,0,0,0.04)',
    '--border': 'rgba(0,0,0,0.08)',
    '--text': '#1e293b',
    '--text-secondary': '#64748b',
  },
}

function applyTheme(t: Theme) {
  const root = document.documentElement

  // Apply CSS variables
  Object.entries(vars[t]).forEach(([k, v]) => root.style.setProperty(k, v))

  // Toggle class
  root.classList.toggle('dark', t === 'dark')
  root.classList.toggle('light', t === 'light')

  // Body
  document.body.style.backgroundColor = t === 'dark' ? '#0a0a0f' : '#f8fafc'
  document.body.style.color = t === 'dark' ? '#f1f5f9' : '#1e293b'

  // Vuetify
  vuetify.theme.global.name.value = t === 'dark' ? 'kanbanDark' : 'kanbanLight'
}

watchEffect(() => {
  applyTheme(theme.value)
  localStorage.setItem('theme', theme.value)
})

export function useTheme() {
  return {
    theme,
    toggle() { theme.value = theme.value === 'dark' ? 'light' : 'dark' },
  }
}
