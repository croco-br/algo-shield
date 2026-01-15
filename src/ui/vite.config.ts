import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    // Only enable devtools in development mode
    ...(process.env.NODE_ENV === 'development' ? [vueDevTools()] : []),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },
  optimizeDeps: {
    // Pre-bundle these dependencies to speed up builds
    include: [
      'vue',
      'vue-router',
      'pinia',
      'vuetify',
      '@fortawesome/fontawesome-svg-core',
      '@fortawesome/vue-fontawesome',
      'vue-i18n'
    ],
    // Enable preloading for critical dependencies
    entries: ['./src/main.ts']
  },
  build: {
    // Enable minification and source maps only in production
    minify: 'esbuild',
    sourcemap: false,
    // Reduce chunk size for better caching
    target: 'esnext',
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
    // Inline small assets (SVG icons, fonts) to reduce HTTP requests
    // Assets smaller than 8KB will be inlined as base64
    assetsInlineLimit: 8192,
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
  // CSS optimization for development and production
  css: {
    devSourcemap: true,
    preprocessorOptions: {
      scss: {
        additionalData: ''
      }
    },
    // PostCSS config in postcss.config.js includes:
    // - @tailwindcss/postcss (Tailwind v4 with automatic purging)
    // - autoprefixer (browser compatibility)
  }
})
