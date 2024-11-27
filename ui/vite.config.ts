import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
// https://github.com/vitejs/vite/discussions/9440
export default defineConfig({
  plugins: [react()],
  build: {
    chunkSizeWarningLimit: 900
  }
})
