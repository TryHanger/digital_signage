import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import { resolve } from 'path'

// read backend port from env at build time (optional)
const BACKEND_PORT = process.env.BACKEND_PORT || '8080'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
    },
  },
  server: {
    proxy: {
      // proxy /api requests to backend server (adjust port if your backend runs elsewhere)
      '/api': {
        target: `http://localhost:${BACKEND_PORT}`,
        changeOrigin: true,
        secure: false,
        rewrite: (p) => p.replace(/^\/api/, '/api'),
      },
    },
  },
})
