import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    vueDevTools(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },
  optimizeDeps: {
    include: ['@mdi/font']
  },
  build: {
    rollupOptions: {
      output: {
        assetFileNames: (assetInfo) => {
          // Preserve font files with their original names and don't process them
          if (assetInfo.name && /\.(woff2?|ttf|eot|otf)$/i.test(assetInfo.name)) {
            return `assets/${assetInfo.name}`
          }
          return 'assets/[name]-[hash][extname]'
        }
      }
    },
    // Ensure font files are copied as-is without processing
    assetsInlineLimit: 0
  },
  server: {
    fs: {
      // Allow serving files from node_modules for @mdi/font
      allow: ['..']
    },
    // Ensure font files are served with correct MIME types
    headers: {
      'Access-Control-Allow-Origin': '*'
    }
  },
  // Ensure CSS imports are handled correctly
  css: {
    preprocessorOptions: {
      scss: {
        additionalData: ''
      }
    }
  }
})
