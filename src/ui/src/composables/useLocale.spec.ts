import { describe, it, expect, beforeEach, vi } from 'vitest'
import { createI18n } from 'vue-i18n'
import { useLocale } from './useLocale'

const mockI18n = createI18n({
  legacy: true,
  locale: 'en-US',
  messages: {
    'pt-BR': { test: 'teste' },
    'en-US': { test: 'test' },
  },
})

vi.mock('@/plugins/i18n', () => ({
  i18n: mockI18n,
}))

describe('useLocale', () => {
  beforeEach(() => {
    localStorage.clear()
    mockI18n.global.locale = 'en-US'
  })

  describe('availableLocales', () => {
    it('returns list of available locales', () => {
      const { availableLocales } = useLocale()

      expect(availableLocales).toHaveLength(2)
      expect(availableLocales[0]?.value).toBe('pt-BR')
      expect(availableLocales[0]?.label).toBe('Português (Brasil)')
      expect(availableLocales[0]?.flag).toBe('🇧🇷')
      expect(availableLocales[1]?.value).toBe('en-US')
      expect(availableLocales[1]?.label).toBe('English (US)')
      expect(availableLocales[1]?.flag).toBe('🇺🇸')
    })
  })

  describe('currentLocale', () => {
    it('returns current locale from i18n', () => {
      const { locale } = useLocale()

      expect(locale.value).toBe('en-US')
    })

    it('updates locale when set', () => {
      const { locale } = useLocale()

      locale.value = 'pt-BR'

      expect(locale.value).toBe('pt-BR')
      expect(localStorage.getItem('locale')).toBe('pt-BR')
    })

    it('persists locale to localStorage', () => {
      const { locale } = useLocale()

      locale.value = 'pt-BR'

      expect(localStorage.getItem('locale')).toBe('pt-BR')
    })
  })

  describe('currentLocaleOption', () => {
    it('returns current locale option object', () => {
      const { currentLocaleOption } = useLocale()

      expect(currentLocaleOption.value?.value).toBe('en-US')
      expect(currentLocaleOption.value?.label).toBe('English (US)')
      expect(currentLocaleOption.value?.flag).toBe('🇺🇸')
    })

    it('updates when locale changes', () => {
      const { locale, currentLocaleOption } = useLocale()

      locale.value = 'pt-BR'

      expect(currentLocaleOption.value?.value).toBe('pt-BR')
      expect(currentLocaleOption.value?.label).toBe('Português (Brasil)')
      expect(currentLocaleOption.value?.flag).toBe('🇧🇷')
    })

    it('falls back to en-US when locale is invalid', async () => {
      const { i18n } = await import('@/plugins/i18n')
      // @ts-expect-error - Testing invalid locale fallback
      i18n.global.locale = 'invalid-locale'
      
      const { currentLocaleOption } = useLocale()

      expect(currentLocaleOption.value?.value).toBe('en-US')
    })
  })

  describe('setLocale', () => {
    it('changes current locale to pt-BR', () => {
      const { locale, setLocale } = useLocale()

      setLocale('pt-BR')

      expect(locale.value).toBe('pt-BR')
      expect(localStorage.getItem('locale')).toBe('pt-BR')
    })

    it('changes current locale to en-US', () => {
      const { locale, setLocale } = useLocale()
      locale.value = 'pt-BR'

      setLocale('en-US')

      expect(locale.value).toBe('en-US')
      expect(localStorage.getItem('locale')).toBe('en-US')
    })
  })

  describe('toggleLocale', () => {
    it('toggles from en-US to pt-BR', () => {
      const { locale, toggleLocale } = useLocale()
      locale.value = 'en-US'

      toggleLocale()

      expect(locale.value).toBe('pt-BR')
      expect(localStorage.getItem('locale')).toBe('pt-BR')
    })

    it('toggles from pt-BR to en-US', () => {
      const { locale, toggleLocale } = useLocale()
      locale.value = 'pt-BR'

      toggleLocale()

      expect(locale.value).toBe('en-US')
      expect(localStorage.getItem('locale')).toBe('en-US')
    })

    it('toggles multiple times correctly', () => {
      const { locale, toggleLocale } = useLocale()
      locale.value = 'en-US'

      toggleLocale()
      expect(locale.value).toBe('pt-BR')

      toggleLocale()
      expect(locale.value).toBe('en-US')

      toggleLocale()
      expect(locale.value).toBe('pt-BR')
    })
  })

  describe('translation function', () => {
    it('provides access to i18n translation function', () => {
      const { t } = useLocale()

      expect(typeof t).toBe('function')
      expect(t('test')).toBe('test')
    })

    it('translates based on current locale', () => {
      const { t, setLocale } = useLocale()

      setLocale('en-US')
      expect(t('test')).toBe('test')

      setLocale('pt-BR')
      expect(t('test')).toBe('teste')
    })
  })
})

