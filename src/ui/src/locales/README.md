# Internationalization (i18n) Documentation

This project uses `vue-i18n` v9 for internationalization support.

## Supported Languages

- **Portuguese (Brazil)**: `pt-BR`
- **English (US)**: `en-US`

## File Structure

```
src/ui/src/
├── locales/
│   ├── pt-BR.json    # Portuguese translations
│   ├── en-US.json    # English translations
│   └── README.md     # This file
├── plugins/
│   └── i18n.ts       # i18n configuration
└── composables/
    └── useLocale.ts  # Language switching composable
```

## How to Add New Translation Keys

### 1. Add the key to both language files

Add the same key to both `pt-BR.json` and `en-US.json` with the appropriate translation:

**pt-BR.json:**
```json
{
  "common": {
    "newKey": "Novo Texto"
  }
}
```

**en-US.json:**
```json
{
  "common": {
    "newKey": "New Text"
  }
}
```

### 2. Use the translation key in your component

**In template:**
```vue
<template>
  <div>{{ $t('common.newKey') }}</div>
</template>
```

**In script:**
```vue
<script setup lang="ts">
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const message = t('common.newKey')
</script>
```

## Translation Key Naming Convention

Use dot-notation to organize keys hierarchically:

```
common.*              # Common/shared translations (buttons, actions, etc.)
auth.*                # Authentication-related translations
header.*              # Header component translations
sidebar.*             # Sidebar component translations
views.*               # View-specific translations
  ├── dashboard.*
  ├── profile.*
  └── ...
components.*          # Component-specific translations
  ├── table.*
  ├── modal.*
  └── ...
errors.*              # Error messages
languages.*           # Language names
```

### Examples:
- `common.save` - Simple key for common text
- `views.dashboard.title` - View-specific title
- `auth.errors.invalidCredentials` - Error message
- `components.table.noData` - Component-specific text

## How Users Change Language

Users can change the language by:
1. Clicking on their avatar in the top-right corner of the header
2. Selecting a language from the dropdown menu

The selected language is automatically saved to `localStorage` and persists across browser sessions.

## Browser Locale Detection

On first visit (when no language preference is saved):
- The application detects the user's browser locale
- If the browser locale is `pt` or `pt-BR`, it defaults to Portuguese
- If the browser locale is `en` or `en-US`, it defaults to English
- For any other locale, it defaults to English

## Language Switching in Code

To programmatically change the language, use the `useLocale` composable:

```vue
<script setup lang="ts">
import { useLocale } from '@/composables/useLocale'

const { locale, setLocale, availableLocales } = useLocale()

// Get current locale
console.log(locale.value) // 'pt-BR' or 'en-US'

// Change locale
setLocale('pt-BR')

// Get all available locales
console.log(availableLocales)
// [
//   { value: 'pt-BR', label: 'Português (Brasil)', flag: '🇧🇷' },
//   { value: 'en-US', label: 'English (US)', flag: '🇺🇸' }
// ]
</script>
```

## Adding a New Language

To add support for a new language:

1. Create a new JSON file: `src/ui/src/locales/[locale-code].json`
2. Copy the structure from `en-US.json` and translate all values
3. Update `src/ui/src/plugins/i18n.ts`:
   ```typescript
   import newLocale from '../locales/new-locale.json'
   
   export const i18n = createI18n({
     // ...
     messages: {
       'pt-BR': ptBR,
       'en-US': enUS,
       'new-locale': newLocale, // Add here
     },
   })
   ```
4. Update `src/ui/src/composables/useLocale.ts`:
   ```typescript
   export type Locale = 'pt-BR' | 'en-US' | 'new-locale'
   
   const availableLocales: LocaleOption[] = [
     // ... existing locales
     {
       value: 'new-locale',
       label: 'New Language',
       flag: '🏴',
     },
   ]
   ```

## Best Practices

1. **Always add translations for all languages** - Don't leave any language file incomplete
2. **Use descriptive keys** - `auth.errors.invalidCredentials` is better than `err1`
3. **Keep translations organized** - Use the hierarchical structure consistently
4. **Test translations** - Verify that translated text fits in the UI across all languages
5. **Use parameters for dynamic text**:
   ```vue
   <!-- Translation file -->
   "welcome": "Welcome, {name}!"
   
   <!-- Component -->
   {{ $t('welcome', { name: user.name }) }}
   ```
6. **Avoid hardcoded strings** - Always use translation keys, even if the text is the same in all languages

## Troubleshooting

### Translation key not found
If you see a key like `common.myKey` displayed instead of the translation:
- Check that the key exists in both language files
- Verify the key path is correct (check for typos)
- Restart the dev server after adding new translations

### TypeScript errors with $t
If you see TypeScript errors about `$t` not existing:
- Ensure `src/ui/src/vue-i18n.d.ts` exists
- Restart your IDE/TypeScript server

### Language not persisting
- Check browser console for localStorage errors
- Verify that `localStorage` is enabled in the browser
- Clear browser cache and try again

