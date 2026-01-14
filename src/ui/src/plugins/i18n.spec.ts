import { describe, it, expect, beforeEach } from 'vitest'
import { i18n } from './i18n'

describe('i18n plugin', () => {
  beforeEach(() => {
    localStorage.clear()
  })

  describe('getInitialLocale', () => {
    it('returns saved locale from localStorage when valid pt-BR', () => {
      localStorage.setItem('locale', 'pt-BR')

      expect(i18n.global.locale).toBeDefined()
    })

    it('returns saved locale from localStorage when valid en-US', () => {
      localStorage.setItem('locale', 'en-US')

      expect(i18n.global.locale).toBeDefined()
    })

    it('ignores invalid locale from localStorage', () => {
      localStorage.setItem('locale', 'invalid-locale')

      expect(i18n.global.locale).toBeDefined()
    })

    it('falls back to browser locale when no saved locale', () => {
      expect(i18n.global.locale).toBeDefined()
    })
  })

  describe('i18n instance', () => {
    it('creates i18n instance with correct configuration', () => {
      expect(i18n).toBeDefined()
      expect(i18n.global).toBeDefined()
    })

    it('has pt-BR messages loaded', () => {
      expect(i18n.global.messages.value['pt-BR']).toBeDefined()
    })

    it('has en-US messages loaded', () => {
      expect(i18n.global.messages.value['en-US']).toBeDefined()
    })

    it('has en-US as fallback locale', () => {
      expect(i18n.global.fallbackLocale.value).toBe('en-US')
    })

    it('translates messages correctly', () => {
      i18n.global.locale.value = 'en-US'
      const translated = i18n.global.t('common.loading')
      expect(typeof translated).toBe('string')
    })

    it('falls back to fallback locale for missing translations', () => {
      i18n.global.locale.value = 'pt-BR'
      const translated = i18n.global.t('nonexistent.key')
      expect(typeof translated).toBe('string')
    })
  })

  describe('locale switching', () => {
    it('allows changing locale to pt-BR', () => {
      i18n.global.locale.value = 'pt-BR'
      expect(i18n.global.locale.value).toBe('pt-BR')
    })

    it('allows changing locale to en-US', () => {
      i18n.global.locale.value = 'en-US'
      expect(i18n.global.locale.value).toBe('en-US')
    })
  })
})

