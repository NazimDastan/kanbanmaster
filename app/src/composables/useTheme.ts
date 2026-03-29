import { ref, watchEffect } from 'vue'
import vuetify from '@/plugins/vuetify'

type Theme = 'dark' | 'light'
type AccentColor = 'purple' | 'blue' | 'green' | 'orange'

const accentColors: Record<AccentColor, {
  primary: string
  primaryLight: string
  primaryDark: string
  secondary: string
}> = {
  purple: { primary: '#6366f1', primaryLight: '#818cf8', primaryDark: '#4f46e5', secondary: '#a855f7' },
  blue: { primary: '#3b82f6', primaryLight: '#60a5fa', primaryDark: '#2563eb', secondary: '#06b6d4' },
  green: { primary: '#10b981', primaryLight: '#34d399', primaryDark: '#059669', secondary: '#14b8a6' },
  orange: { primary: '#f59e0b', primaryLight: '#fbbf24', primaryDark: '#d97706', secondary: '#f97316' },
}

const lightBgTints: Record<AccentColor, string> = {
  purple: '#f5f3ff',
  blue: '#eff6ff',
  green: '#ecfdf5',
  orange: '#fffbeb',
}

const validAccents: AccentColor[] = ['purple', 'blue', 'green', 'orange']
const savedAccent = localStorage.getItem('accentColor') as AccentColor | null
const accentColor = ref<AccentColor>(
  savedAccent && validAccents.includes(savedAccent) ? savedAccent : 'purple',
)

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

function applyTheme(t: Theme, accent: AccentColor) {
  const root = document.documentElement
  const themeVars = { ...vars[t] }

  // Apply accent-tinted background for light mode
  if (t === 'light') {
    themeVars['--bg'] = lightBgTints[accent]
  }

  Object.entries(themeVars).forEach(([k, v]) => root.style.setProperty(k, v))

  // Apply accent color CSS variables
  const colors = accentColors[accent]
  root.style.setProperty('--color-primary', colors.primary)
  root.style.setProperty('--color-primary-light', colors.primaryLight)
  root.style.setProperty('--color-primary-dark', colors.primaryDark)
  root.style.setProperty('--color-secondary', colors.secondary)

  root.classList.toggle('dark', t === 'dark')
  root.classList.toggle('light', t === 'light')
  document.body.style.backgroundColor = themeVars['--bg']
  document.body.style.color = themeVars['--text']
  // Vuetify 4 uses theme.change() instead of theme.global.name.value
  const themeName = t === 'dark' ? 'kanbanDark' : 'kanbanLight'
  if (typeof vuetify.theme.change === 'function') {
    vuetify.theme.change(themeName)
  } else {
    vuetify.theme.global.name.value = themeName
  }
}

watchEffect(() => {
  applyTheme(theme.value, accentColor.value)
  localStorage.setItem('theme', theme.value)
  localStorage.setItem('accentColor', accentColor.value)
})

export type { AccentColor }

export function useTheme() {
  return {
    theme,
    accentColor,
    accentColors,
    toggle() { theme.value = theme.value === 'dark' ? 'light' : 'dark' },
    setAccent(color: AccentColor) { accentColor.value = color },
  }
}
