import { defineStore } from 'pinia'
import { ref, watch } from 'vue'
import { api, setSyntheticModeStorage } from '@/lib/api'

interface SyntheticModeResponse {
  enabled: boolean
  updated_at: string
}

// Helper function to get synthetic mode from localStorage
function getSyntheticModeFromStorage(): boolean {
  try {
    const mode = localStorage.getItem('synthetic_mode')
    return mode === 'true'
  } catch {
    return false
  }
}

export const useSystemModeStore = defineStore('systemMode', () => {
  // Initialize from localStorage if available
  const syntheticMode = ref(getSyntheticModeFromStorage())
  const loading = ref(false)
  const error = ref<string | null>(null)
  const updatedAt = ref<Date | null>(null)

  // Sync synthetic mode with localStorage for API requests
  watch(syntheticMode, (newValue) => {
    setSyntheticModeStorage(newValue)
  })

  async function loadMode() {
    loading.value = true
    error.value = null
    try {
      const response = await api.get<SyntheticModeResponse>('/api/v1/system/mode')
      syntheticMode.value = response.enabled
      setSyntheticModeStorage(response.enabled)
      updatedAt.value = new Date(response.updated_at)
    } catch (e: any) {
      error.value = e.message || 'Failed to load synthetic mode'
      console.error('Error loading synthetic mode:', e)
    } finally {
      loading.value = false
    }
  }

  async function toggleMode() {
    loading.value = true
    error.value = null
    try {
      const response = await api.put<SyntheticModeResponse>('/api/v1/system/mode', {
        enabled: !syntheticMode.value
      })
      syntheticMode.value = response.enabled
      setSyntheticModeStorage(response.enabled)
      updatedAt.value = new Date(response.updated_at)
      return response.enabled
    } catch (e: any) {
      error.value = e.message || 'Failed to toggle synthetic mode'
      console.error('Error toggling synthetic mode:', e)
      throw e
    } finally {
      loading.value = false
    }
  }

  async function setMode(enabled: boolean) {
    loading.value = true
    error.value = null
    try {
      const response = await api.put<SyntheticModeResponse>('/api/v1/system/mode', {
        enabled
      })
      syntheticMode.value = response.enabled
      setSyntheticModeStorage(response.enabled)
      updatedAt.value = new Date(response.updated_at)
      return response.enabled
    } catch (e: any) {
      error.value = e.message || 'Failed to set synthetic mode'
      console.error('Error setting synthetic mode:', e)
      throw e
    } finally {
      loading.value = false
    }
  }

  return {
    syntheticMode,
    loading,
    error,
    updatedAt,
    loadMode,
    toggleMode,
    setMode
  }
})
