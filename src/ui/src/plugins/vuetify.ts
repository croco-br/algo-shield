import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import { h } from 'vue'
// Import Font Awesome CSS instead of MDI
import '@fortawesome/fontawesome-free/css/all.css'
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

// Custom Font Awesome iconset for Vuetify
const fa = {
  component: (props: any) => {
    let iconClass = props.icon
    // Convert icon name to Font Awesome class
    if (!iconClass.startsWith('fa-') && !iconClass.startsWith('fas ') && !iconClass.startsWith('far ') && !iconClass.startsWith('fal ')) {
      iconClass = `fa-${iconClass}`
    }
    if (!iconClass.startsWith('fas ') && !iconClass.startsWith('far ') && !iconClass.startsWith('fal ')) {
      iconClass = `fas ${iconClass}`
    }
    return h('i', { 
      class: iconClass,
      style: { 
        fontSize: props.size || '1em',
        color: props.color 
      } 
    })
  }
}

export const vuetify = createVuetify({
  components,
  directives,
  icons: {
    defaultSet: 'fa',
    sets: {
      fa,
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

