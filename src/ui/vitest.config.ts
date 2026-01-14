import { fileURLToPath } from 'node:url'
import { defineConfig, configDefaults } from 'vitest/config'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [
    {
      name: 'vitest-css-stub',
      resolveId(id) {
        // Intercept CSS imports before they reach Node.js
        if (
          id.endsWith('.css') ||
          id.endsWith('.scss') ||
          id.endsWith('.sass') ||
          id.endsWith('.less') ||
          id.endsWith('.styl')
        ) {
          // Return the id with a marker that we can recognize in load()
          return id + '?stub'
        }
        return null
      },
      load(id) {
        // Handle stubbed CSS imports
        if (id.includes('?stub')) {
          return 'export default {}'
        }
        return null
      },
      transform(_src, id) {
        // Stub out Vue style blocks
        if (id.includes('?vue&type=style')) {
          return {
            code: 'export default {}',
            map: null
          }
        }
      }
    },
    vue()
  ],
  test: {
    environment: 'happy-dom',
    exclude: [
      ...configDefaults.exclude,
      'e2e/**',
      '**/Header.spec.ts',
      '**/LoginView.spec.ts',
      '**/useLocale.spec.ts',
    ],
    root: fileURLToPath(new URL('./', import.meta.url)),
    globals: true,
    css: false,
    setupFiles: ['./vitest.setup.ts'],
    server: {
      deps: {
        inline: ['vuetify']
      }
    },
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html'],
      exclude: [
        'node_modules/',
        'dist/',
        '**/*.spec.ts',
        '**/*.d.ts',
        'vite.config.ts',
        'vitest.config.ts',
        'vitest.setup.ts',
      ],
    },
  },
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },
  optimizeDeps: {
    exclude: ['vuetify'],
  },
})

