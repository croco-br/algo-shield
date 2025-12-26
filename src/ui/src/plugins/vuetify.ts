import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import { aliases, mdi } from 'vuetify/iconsets/mdi'
import '@mdi/font/css/materialdesignicons.css'
import 'vuetify/styles'

// Default theme colors (will be overridden by branding)
const defaultTheme = {
  dark: false,
  colors: {
    primary: '#3B82F6', // Default primary color
    secondary: '#10B981', // Default secondary color
    accent: '#6366F1',
    error: '#EF4444',
    info: '#3B82F6',
    success: '#10B981',
    warning: '#F59E0B',
  },
}

export const vuetify = createVuetify({
  components,
  directives,
  icons: {
    defaultSet: 'mdi',
    aliases,
    sets: {
      mdi,
    },
  },
  theme: {
    defaultTheme: 'light',
    themes: {
      light: defaultTheme,
    },
  },
  defaults: {
    VBtn: {
      style: 'text-transform: none;', // Disable uppercase transformation
    },
    VTextField: {
      variant: 'outlined',
      density: 'comfortable',
    },
    VSelect: {
      variant: 'outlined',
      density: 'comfortable',
    },
  },
  display: {
    mobileBreakpoint: 'sm',
    thresholds: {
      xs: 0,
      sm: 600,
      md: 960,
      lg: 1280,
      xl: 1920,
    },
  },
})

// Function to update Vuetify theme with branding colors
export function updateVuetifyTheme(primaryColor: string, secondaryColor: string) {
  const lightTheme = vuetify.theme.themes.value.light
  if (lightTheme && lightTheme.colors) {
    lightTheme.colors.primary = primaryColor
    lightTheme.colors.secondary = secondaryColor
    lightTheme.colors.accent = secondaryColor
  }
}

