import { computed } from 'vue'

export type ActionType = 'allow' | 'block' | 'review'

export function useActionBadge(action: ActionType | string) {
  const variant = computed(() => {
    switch (action) {
      case 'allow':
        return 'success'
      case 'block':
        return 'danger'
      case 'review':
        return 'warning'
      default:
        return 'info'
    }
  })

  return {
    variant,
  }
}
