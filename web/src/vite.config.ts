import tailwindcss from '@tailwindcss/vite';
import react from '@vitejs/plugin-react';
import { readFileSync } from 'fs';
import { dirname, resolve } from 'path';
import { defineConfig } from 'vite';
import { VitePWA } from 'vite-plugin-pwa';

const __dirname = dirname(import.meta.filename);

const getBuildTag = () => {
    let commit = 'unknown';
    try {
        const filePath = resolve(__dirname, '..', '..', '.git', 'logs', 'HEAD');
        commit = readFileSync(filePath, 'utf-8')
            .trim()
            .split('\n')
            .pop()!
            .split(' ')[1]
            .slice(0, 8);
    } catch {
        commit = 'unknown';
    }

    let version = 'custombuild';
    try {
        const filePath = resolve(__dirname, '..', '..', 'VERSION');
        version = readFileSync(filePath, 'utf-8').trim();
    } catch {
        version = 'custombuild';
    }

    return `${version}-${commit}-${Math.floor(Date.now() / 1000)}`;
};

// https://vite.dev/config/
export default defineConfig({
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
        'import.meta.env.BUILD_TAG': JSON.stringify(getBuildTag())
    }
});
