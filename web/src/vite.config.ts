import { readFileSync } from 'node:fs';
import { fileURLToPath } from 'node:url';

import tailwindcss from '@tailwindcss/vite';
import react from '@vitejs/plugin-react';
import license, { Dependency } from 'rollup-plugin-license';
import { defineConfig } from 'vite';
import { VitePWA } from 'vite-plugin-pwa';

const markdownBlock = (text: string) => {
    let fence = '~~~';
    while (text.includes(fence)) {
        fence += '~';
    }
    return `${fence}\n${text}\n${fence}`;
};

const formatLicenseReport = (dependencies: Dependency[]) => {
    const reports = dependencies.map((dependency) => {
        const metadata = [`- **License:** ${dependency.license || 'Unknown'}`];
        if (dependency.homepage) {
            metadata.push(`- **Homepage:** ${dependency.homepage}`);
        }
        const sections = [`## ${dependency.name} ${dependency.version}`, metadata.join('\n')];
        if (dependency.licenseText) {
            sections.push('### License text', markdownBlock(dependency.licenseText));
        }
        if (dependency.noticeText) {
            sections.push('### Notice', markdownBlock(dependency.noticeText));
        }
        return {
            name: dependency.name ?? '',
            version: dependency.version ?? '',
            text: sections.join('\n\n')
        };
    });

    // These bundled FDSN assets are outside the npm dependency tree.
    for (const [name, file] of [
        ['Alpine.js', 'alpine-LICENSE.txt'],
        ['Tailwind CSS', 'tailwindcss-LICENSE.txt']
    ]) {
        const text = readFileSync(
            new URL(`../../internal/server/router/fdsnws/static/assets/${file}`, import.meta.url),
            'utf-8'
        );
        reports.push({
            name,
            version: '',
            text: [
                `## ${name} (FDSN web interface)`,
                '- **License:** MIT',
                markdownBlock(text)
            ].join('\n\n')
        });
    }

    return (
        reports
            .sort(
                (left, right) =>
                    left.name.localeCompare(right.name, 'en', { sensitivity: 'base' }) ||
                    left.version.localeCompare(right.version, 'en')
            )
            .map(({ text }) => text)
            .join('\n\n---\n\n') + '\n'
    );
};

// https://vite.dev/config/
export default defineConfig({
    base: './',
    plugins: [
        react(),
        tailwindcss(),
        license({
            thirdParty: {
                includePrivate: true,
                multipleVersions: true,
                output: [
                    {
                        file: fileURLToPath(
                            new URL('../dist/third-party-web-licenses.md', import.meta.url)
                        ),
                        template: formatLicenseReport
                    },
                    // Keep the latest report available to the Vite development server too.
                    {
                        file: fileURLToPath(
                            new URL('./public/third-party-web-licenses.md', import.meta.url)
                        ),
                        template: formatLicenseReport
                    }
                ]
            }
        }),
        VitePWA({
            injectRegister: 'auto',
            registerType: 'prompt',
            workbox: {
                globPatterns: ['**/*.{js,css,html,ico,png,svg,webp,woff2}'],
                navigateFallbackDenylist: [/^\/api/, /^\/fdsnws(?:\/|\?|$)/],
                maximumFileSizeToCacheInBytes: 3000000
            },
            manifest: {
                theme_color: '#8d59ef'
            }
        })
    ],
    build: {
        chunkSizeWarningLimit: 4096,
        sourcemap: false,
        outDir: '../dist',
        rollupOptions: {
            output: {
                manualChunks(id) {
                    if (id.includes('/src/locales/') && id.endsWith('.json')) {
                        return 'locales';
                    }
                }
            }
        }
    },
    server: {
        host: '0.0.0.0',
        port: 3000
    },
    define: {
        'globalThis.__DEV__': JSON.stringify(false) // Disable DevTools message
    }
});
