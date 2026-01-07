import { DefineComponent } from 'vue'

declare module '@vue/runtime-core' {
  export interface ComponentCustomProperties {
    $t: (key: string, params?: Record<string, any>) => string
    $tc: (key: string, choice?: number, params?: Record<string, any>) => string
    $te: (key: string) => boolean
  }
}

export {}

