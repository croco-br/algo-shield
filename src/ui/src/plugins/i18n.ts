import { createI18n } from 'vue-i18n'
import ptBR from '../locales/pt-BR.json'
import enUS from '../locales/en-US.json'

// Detect browser locale
function getBrowserLocale(): string {
  const browserLocale = navigator.language || (navigator as any).userLanguage
  
  // Match exact locales
  if (browserLocale === 'pt-BR' || browserLocale === 'pt') {
    return 'pt-BR'
  }
  if (browserLocale === 'en-US' || browserLocale === 'en') {
    return 'en-US'
  }
  
  // Default to English
  return 'en-US'
}

// Get initial locale from localStorage or browser
function getInitialLocale(): string {
  const savedLocale = localStorage.getItem('locale')
  
  if (savedLocale && (savedLocale === 'pt-BR' || savedLocale === 'en-US')) {
    return savedLocale
  }
  
  return getBrowserLocale()
}

export const i18n = createI18n({
  legacy: false, // Use Composition API mode (legacy mode is deprecated in v11 and will be removed in v12)
  locale: getInitialLocale(),
  fallbackLocale: 'en-US',
  messages: {
    'pt-BR': ptBR,
    'en-US': enUS,
  },
  silentTranslationWarn: true,
  silentFallbackWarn: true,
  warnHtmlMessage: false,
})

