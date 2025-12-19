/**
 * AlgoShield Design Tokens
 * Enterprise-grade design system for fraud detection and AML software
 */

export const tokens = {
  colors: {
    // Brand Colors - Transmit trust, stability, and enterprise maturity
    brand: {
      primary: '#1E3A8A',      // Blue-900: trust, banking, stability
      secondary: '#3730A3',    // Indigo-800: technology, precision
      accent: '#10B981',       // Emerald-500: approval, security, success
    },

    // Neutral Scale - Professional and clean
    neutral: {
      50: '#F9FAFB',           // Subtle backgrounds
      100: '#F3F4F6',          // Secondary cards
      200: '#E5E7EB',          // Borders
      300: '#D1D5DB',          // Disabled borders
      400: '#9CA3AF',          // Secondary text
      500: '#6B7280',          // Icons
      600: '#4B5563',          // Primary text
      700: '#374151',          // Headings
      800: '#1F2937',          // Strong emphasis
      900: '#111827',          // Maximum contrast
    },

    // Slate Scale - For backgrounds (softer than neutral)
    slate: {
      50: '#F8FAFC',
      100: '#F1F5F9',
      200: '#E2E8F0',
      400: '#94A3B8',
      500: '#64748B',
      600: '#475569',
      700: '#334155',
      900: '#0F172A',
    },

    // Semantic States
    state: {
      focus: '#2563EB',        // Blue-600: active focus
      error: '#DC2626',        // Red-600: critical errors
      success: '#059669',      // Emerald-600: successful actions
      warning: '#D97706',      // Amber-600: attention required
      info: '#0891B2',         // Cyan-600: informational
    },

    // Action Colors
    action: {
      allow: '#10B981',        // Emerald-500
      block: '#EF4444',        // Red-500
      review: '#F59E0B',       // Amber-500
      score: '#3B82F6',        // Blue-500
    },
  },

  // Typography Scale
  typography: {
    // Font Families
    fontFamily: {
      sans: 'ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif',
      mono: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace',
    },

    // Font Sizes (rem / px)
    fontSize: {
      xs: '0.75rem',         // 12px
      sm: '0.875rem',        // 14px
      base: '1rem',          // 16px
      lg: '1.125rem',        // 18px
      xl: '1.25rem',         // 20px
      '2xl': '1.5rem',       // 24px
      '3xl': '1.875rem',     // 30px
      '4xl': '2.25rem',      // 36px
    },

    // Font Weights
    fontWeight: {
      normal: 400,
      medium: 500,
      semibold: 600,
      bold: 700,
    },

    // Line Heights
    lineHeight: {
      tight: 1.25,
      normal: 1.5,
      relaxed: 1.625,
    },
  },

  // Spacing Scale (px)
  spacing: {
    0: '0',
    1: '0.25rem',    // 4px
    2: '0.5rem',     // 8px
    3: '0.75rem',    // 12px
    4: '1rem',       // 16px
    5: '1.25rem',    // 20px
    6: '1.5rem',     // 24px
    8: '2rem',       // 32px
    10: '2.5rem',    // 40px
    12: '3rem',      // 48px
    14: '3.5rem',    // 56px
    16: '4rem',      // 64px
  },

  // Border Radius
  borderRadius: {
    sm: '0.25rem',    // 4px
    md: '0.375rem',   // 6px
    lg: '0.5rem',     // 8px
    xl: '0.75rem',    // 12px
    '2xl': '1rem',    // 16px
    full: '9999px',   // Fully rounded
  },

  // Shadows - Enterprise-grade depth
  shadows: {
    sm: '0 1px 2px 0 rgba(0, 0, 0, 0.05)',
    md: '0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -2px rgba(0, 0, 0, 0.1)',
    lg: '0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -4px rgba(0, 0, 0, 0.1)',
    xl: '0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 8px 10px -6px rgba(0, 0, 0, 0.1)',
    card: '0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px -1px rgba(0, 0, 0, 0.1)',

    // Colored shadows for emphasis
    primaryButton: '0 10px 25px -5px rgba(37, 99, 235, 0.3)',
    primaryButtonHover: '0 15px 30px -5px rgba(37, 99, 235, 0.4)',

    // Focus rings
    focusRing: '0 0 0 4px',
  },

  // Transitions
  transitions: {
    fast: '100ms cubic-bezier(0.4, 0, 0.2, 1)',
    default: '200ms cubic-bezier(0.4, 0, 0.2, 1)',
    slow: '300ms cubic-bezier(0.4, 0, 0.2, 1)',
  },

  // Component-specific tokens
  components: {
    // Inputs
    input: {
      padding: '14px 16px',           // py-3.5 px-4
      borderWidth: '2px',
      borderRadius: '0.5rem',         // rounded-lg
      fontSize: '0.875rem',           // text-sm
      fontWeight: 500,                // font-medium
    },

    // Buttons
    button: {
      padding: {
        sm: '8px 16px',               // py-2 px-4
        md: '14px 24px',              // py-3.5 px-6
        lg: '16px 32px',              // py-4 px-8
      },
      borderRadius: '0.5rem',         // rounded-lg
      fontSize: {
        sm: '0.875rem',               // text-sm
        md: '0.875rem',               // text-sm
        lg: '1rem',                   // text-base
      },
      fontWeight: 600,                // font-semibold
    },

    // Cards
    card: {
      padding: '56px 48px',           // py-14 px-12
      maxWidth: '480px',
      borderRadius: '0.75rem',        // rounded-xl
      shadow: '0 8px 30px rgba(0, 0, 0, 0.08)',
    },

    // Modal
    modal: {
      overlay: 'rgba(0, 0, 0, 0.5)',
      maxWidth: {
        sm: '400px',
        md: '600px',
        lg: '800px',
      },
    },

    // Badge
    badge: {
      padding: '4px 12px',            // py-1 px-3
      borderRadius: '9999px',         // rounded-full
      fontSize: '0.75rem',            // text-xs
      fontWeight: 600,                // font-semibold
    },
  },
}

export default tokens
