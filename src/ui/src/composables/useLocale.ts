import { computed } from 'vue'
import { i18n } from '@/plugins/i18n'

export type Locale = 'pt-BR' | 'en-US'

export interface LocaleOption {
  value: Locale
  label: string
  flag: string
}

export function useLocale() {
  const availableLocales: LocaleOption[] = [
    {
      value: 'pt-BR',
      label: 'Português (Brasil)',
      flag: '🇧🇷',
    },
    {
      value: 'en-US',
      label: 'English (US)',
      flag: '🇺🇸',
    },
  ]

  const currentLocale = computed<Locale>({
    get: () => i18n.global.locale as Locale,
    set: (value: Locale) => {
      i18n.global.locale = value
      localStorage.setItem('locale', value)
    },
  })

  const currentLocaleOption = computed(() => {
    return availableLocales.find((l) => l.value === currentLocale.value) || availableLocales[1]
  })

  function setLocale(newLocale: Locale) {
    currentLocale.value = newLocale
  }

  function toggleLocale() {
    currentLocale.value = currentLocale.value === 'pt-BR' ? 'en-US' : 'pt-BR'
  }

  return {
    locale: currentLocale,
    currentLocaleOption,
    availableLocales,
    setLocale,
    toggleLocale,
    t: i18n.global.t,
  }
}

