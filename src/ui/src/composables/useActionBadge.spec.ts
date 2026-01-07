import { describe, it, expect } from 'vitest'
import { useActionBadge } from './useActionBadge'

describe('useActionBadge', () => {
  describe('variant computation', () => {
    it('returns success variant for allow action', () => {
      const { variant } = useActionBadge('allow')

      expect(variant.value).toBe('success')
    })

    it('returns danger variant for block action', () => {
      const { variant } = useActionBadge('block')

      expect(variant.value).toBe('danger')
    })

    it('returns warning variant for review action', () => {
      const { variant } = useActionBadge('review')

      expect(variant.value).toBe('warning')
    })

    it('returns info variant for unknown action', () => {
      const { variant } = useActionBadge('unknown')

      expect(variant.value).toBe('info')
    })

    it('returns info variant for empty string', () => {
      const { variant } = useActionBadge('')

      expect(variant.value).toBe('info')
    })

    it('returns info variant for custom action type', () => {
      const { variant } = useActionBadge('custom-action')

      expect(variant.value).toBe('info')
    })
  })

  describe('edge cases', () => {
    it('handles mixed case input correctly', () => {
      const { variant } = useActionBadge('ALLOW')

      expect(variant.value).toBe('info')
    })

    it('handles action with special characters', () => {
      const { variant } = useActionBadge('allow-special')

      expect(variant.value).toBe('info')
    })

    it('handles action with numbers', () => {
      const { variant } = useActionBadge('action123')

      expect(variant.value).toBe('info')
    })
  })
})

