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
    include: []
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
        },
        manualChunks: (id) => {
          // Split node_modules into vendor chunks
          if (id.includes('node_modules')) {
            // Vuetify is large, split it into its own chunk
            if (id.includes('vuetify')) {
              return 'vendor-vuetify'
            }
            // Font Awesome icons
            if (id.includes('@fortawesome')) {
              return 'vendor-icons'
            }
            // Vue core libraries
            if (id.includes('vue') || id.includes('pinia') || id.includes('vue-router')) {
              return 'vendor-vue'
            }
            // Other utilities (prismjs, etc.)
            if (id.includes('prismjs')) {
              return 'vendor-utils'
            }
            // All other node_modules
            return 'vendor'
          }
        }
      }
    },
    // Ensure font files are copied as-is without processing
    assetsInlineLimit: 0,
    // Increase chunk size warning limit since we're splitting properly
    chunkSizeWarningLimit: 1000
  },
  server: {
    fs: {
      // Allow serving files from node_modules
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
