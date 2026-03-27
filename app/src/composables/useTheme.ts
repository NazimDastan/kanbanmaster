import { ref, watchEffect } from 'vue'
import vuetify from '@/plugins/vuetify'

type Theme = 'dark' | 'light'

const theme = ref<Theme>((localStorage.getItem('theme') as Theme) ?? 'dark')

const vars: Record<Theme, Record<string, string>> = {
  dark: {
    '--bg': '#0a0a0f',
    '--bg-card': '#0f0f1a',
    '--bg-elevated': '#161625',
    '--bg-hover': 'rgba(255,255,255,0.02)',
    '--bg-input': 'rgba(255,255,255,0.05)',
    '--bg-active': 'rgba(99,102,241,0.15)',
    '--border': 'rgba(255,255,255,0.04)',
    '--border-hover': 'rgba(255,255,255,0.1)',
    '--text': '#f1f5f9',
    '--text-secondary': '#94a3b8',
    '--text-muted': 'rgba(255,255,255,0.25)',
    '--text-faint': 'rgba(255,255,255,0.15)',
    '--shadow': '0 4px 24px rgba(0,0,0,0.3)',
  },
  light: {
    '--bg': '#f8fafc',
    '--bg-card': '#ffffff',
    '--bg-elevated': '#f1f5f9',
    '--bg-hover': 'rgba(0,0,0,0.02)',
    '--bg-input': 'rgba(0,0,0,0.04)',
    '--bg-active': 'rgba(99,102,241,0.1)',
    '--border': 'rgba(0,0,0,0.08)',
    '--border-hover': 'rgba(0,0,0,0.15)',
    '--text': '#1e293b',
    '--text-secondary': '#64748b',
    '--text-muted': 'rgba(0,0,0,0.3)',
    '--text-faint': 'rgba(0,0,0,0.1)',
    '--shadow': '0 4px 24px rgba(0,0,0,0.08)',
  },
}

function applyTheme(t: Theme) {
  const root = document.documentElement
  Object.entries(vars[t]).forEach(([k, v]) => root.style.setProperty(k, v))
  root.classList.toggle('dark', t === 'dark')
  root.classList.toggle('light', t === 'light')
  document.body.style.backgroundColor = vars[t]['--bg']
  document.body.style.color = vars[t]['--text']
  // Vuetify 4 uses theme.change() instead of theme.global.name.value
  const themeName = t === 'dark' ? 'kanbanDark' : 'kanbanLight'
  if (typeof vuetify.theme.change === 'function') {
    vuetify.theme.change(themeName)
  } else {
    vuetify.theme.global.name.value = themeName
  }
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
