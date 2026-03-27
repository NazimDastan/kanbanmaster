import { createI18n } from 'vue-i18n'
import en from './en'
import tr from './tr'
import de from './de'
import ar from './ar'

type Locale = 'en' | 'tr' | 'de' | 'ar'

const validLocales: Locale[] = ['en', 'tr', 'de', 'ar']
const saved = localStorage.getItem('locale') as Locale | null
const savedLocale: Locale = saved && validLocales.includes(saved) ? saved : 'en'

const i18n = createI18n({
  legacy: false,
  locale: savedLocale,
  fallbackLocale: 'en',
  messages: { en, tr, de, ar },
})

export function setLocale(locale: string) {
  const l = locale as Locale
  i18n.global.locale.value = l
  localStorage.setItem('locale', l)
  document.documentElement.setAttribute('dir', l === 'ar' ? 'rtl' : 'ltr')
  document.documentElement.setAttribute('lang', l)
}

export default i18n
