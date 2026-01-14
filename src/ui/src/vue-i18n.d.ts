import { DefineComponent } from 'vue'

// Global properties for templates ($t, $te)
declare module '@vue/runtime-core' {
  export interface ComponentCustomProperties {
    $t: (key: string, params?: Record<string, any>) => string
    $te: (key: string) => boolean
  }
}

// Type definitions for Composition API mode
declare module 'vue-i18n' {
  // Define the structure of locale messages for type safety
  export interface DefineLocaleMessage {
    [key: string]: string | DefineLocaleMessage
  }
}

export {}

