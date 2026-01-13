/**
 * Currency formatting composable
 * Provides consistent currency formatting across the application
 */

export interface CurrencyFormatOptions {
  locale?: string
  minimumFractionDigits?: number
  maximumFractionDigits?: number
}

/**
 * Format a monetary amount for display
 * @param amount - The amount to format (can be integer cents or decimal)
 * @param currency - The ISO 4217 currency code (e.g., 'USD', 'BRL', 'EUR')
 * @param options - Formatting options
 * @returns Formatted currency string
 */
export function formatCurrency(
  amount: number | undefined | null,
  currency?: string,
  options: CurrencyFormatOptions = {}
): string {
  if (amount === undefined || amount === null) {
    return '-'
  }

  const {
    locale = 'en-US',
    minimumFractionDigits = 2,
    maximumFractionDigits = 2
  } = options

  try {
    return new Intl.NumberFormat(locale, {
      style: 'currency',
      currency: currency || 'USD',
      minimumFractionDigits,
      maximumFractionDigits
    }).format(amount)
  } catch {
    // Fallback for invalid currency codes
    return new Intl.NumberFormat(locale, {
      style: 'decimal',
      minimumFractionDigits,
      maximumFractionDigits
    }).format(amount) + ' ' + (currency || 'USD')
  }
}

/**
 * Format a number as compact currency (e.g., $1.5M, $10K)
 * @param amount - The amount to format
 * @param currency - The ISO 4217 currency code
 * @param locale - The locale for formatting
 * @returns Compact formatted currency string
 */
export function formatCompactCurrency(
  amount: number | undefined | null,
  currency?: string,
  locale = 'en-US'
): string {
  if (amount === undefined || amount === null) {
    return '-'
  }

  try {
    return new Intl.NumberFormat(locale, {
      style: 'currency',
      currency: currency || 'USD',
      notation: 'compact',
      compactDisplay: 'short',
      maximumFractionDigits: 1
    }).format(amount)
  } catch {
    // Fallback for compact notation
    if (amount >= 1000000) {
      return formatCurrency(amount / 1000000, currency) + 'M'
    }
    if (amount >= 1000) {
      return formatCurrency(amount / 1000, currency) + 'K'
    }
    return formatCurrency(amount, currency)
  }
}

/**
 * Parse a formatted currency string back to a number
 * @param value - The formatted currency string
 * @returns The numeric value or NaN if parsing fails
 */
export function parseCurrency(value: string): number {
  // Remove currency symbols, spaces, and thousands separators
  const cleaned = value
    .replace(/[^\d.,-]/g, '')
    .replace(/,/g, '')
  
  return parseFloat(cleaned)
}

/**
 * Vue composable for currency formatting
 * Provides reactive currency formatting utilities
 */
export function useCurrency() {
  return {
    formatCurrency,
    formatCompactCurrency,
    parseCurrency
  }
}
