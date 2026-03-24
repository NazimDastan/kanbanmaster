import { createI18n } from 'vue-i18n'
import en from './en'
import tr from './tr'
import de from './de'
import ar from './ar'

const savedLocale = localStorage.getItem('locale') ?? 'en'

const i18n = createI18n({
  legacy: false,
  locale: savedLocale,
  fallbackLocale: 'en',
  messages: { en, tr, de, ar },
})

export function setLocale(locale: string) {
  i18n.global.locale.value = locale
  localStorage.setItem('locale', locale)
  document.documentElement.setAttribute('dir', locale === 'ar' ? 'rtl' : 'ltr')
  document.documentElement.setAttribute('lang', locale)
}

export default i18n
