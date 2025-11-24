import tailwindcss from '@tailwindcss/vite';
import react from '@vitejs/plugin-react';
import { defineConfig } from 'vite';
import { VitePWA } from 'vite-plugin-pwa';

// https://vite.dev/config/
export default defineConfig({
    base: './',
    plugins: [
        react(),
        tailwindcss(),
        VitePWA({
            injectRegister: 'auto',
            registerType: 'prompt',
            workbox: {
                globPatterns: ['**/*.{js,css,html,ico,png,svg,woff2}'],
                navigateFallbackDenylist: [/^\/api/],
                maximumFileSizeToCacheInBytes: 3000000
            }
        })
    ],
    build: {
        chunkSizeWarningLimit: 1600,
        sourcemap: false,
        outDir: '../dist'
    },
    server: {
        host: '0.0.0.0',
        port: 3000
    },
    define: {
        'globalThis.__DEV__': JSON.stringify(false) // Disable DevTools message
    }
});
