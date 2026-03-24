import 'vuetify/styles'
import '@mdi/font/css/materialdesignicons.css'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'

const shared = {
  primary: '#6366f1',
  secondary: '#a855f7',
  accent: '#06b6d4',
  success: '#10b981',
  warning: '#f59e0b',
  error: '#ef4444',
  'on-primary': '#ffffff',
  'on-secondary': '#ffffff',
}

const kanbanDark = {
  dark: true,
  colors: {
    ...shared,
    background: '#0a0a0f',
    surface: '#161625',
    'surface-variant': '#1e1e32',
    'on-background': '#f1f5f9',
    'on-surface': '#f1f5f9',
  },
}

const kanbanLight = {
  dark: false,
  colors: {
    ...shared,
    background: '#f8fafc',
    surface: '#ffffff',
    'surface-variant': '#f1f5f9',
    'on-background': '#1e293b',
    'on-surface': '#1e293b',
  },
}

const vuetify = createVuetify({
  components,
  directives,
  theme: {
    defaultTheme: localStorage.getItem('theme') === 'light' ? 'kanbanLight' : 'kanbanDark',
    themes: { kanbanDark, kanbanLight },
  },
  defaults: {
    VBtn: { rounded: 'lg', variant: 'flat' },
    VCard: { rounded: 'xl', elevation: 0, color: 'surface' },
    VTextField: { variant: 'outlined', density: 'comfortable', rounded: 'lg', color: 'primary' },
    VSelect: { variant: 'outlined', density: 'comfortable', rounded: 'lg', color: 'primary' },
    VChip: { rounded: 'xl', size: 'small' },
    VNavigationDrawer: { color: 'surface' },
    VAppBar: { color: 'surface', elevation: 0 },
  },
})

export default vuetify
