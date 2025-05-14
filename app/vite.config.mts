// Plugins
import Components from 'unplugin-vue-components/vite'
import Vue from '@vitejs/plugin-vue'
import Vuetify, { transformAssetUrls } from 'vite-plugin-vuetify'
import Fonts from 'unplugin-fonts/vite'
import { viteCommonjs } from '@originjs/vite-plugin-commonjs'
import removeConsole from 'vite-plugin-remove-console'

// Utilities
import { defineConfig } from 'vite'
import { fileURLToPath, URL } from 'node:url'

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [
        Vue({
            template: { transformAssetUrls },
        }),
        viteCommonjs(),
        // https://github.com/vuetifyjs/vuetify-loader/tree/master/packages/vite-plugin#readme
        Vuetify(),
        Components(),
        Fonts({
            fontsource: {
                families: [
                    {
                        name: 'Roboto',
                        weights: [100, 300, 400, 500, 700, 900],
                        styles: ['normal', 'italic'],
                    },
                ],
            },
        }),
        removeConsole({
            external: ['log', 'error'],
        }),
    ],
    optimizeDeps: {
        exclude: ['vuetify', '@cornerstonejs/dicom-image-loader'],
        include: ['dicom-parser'],
    },
    worker: {
        format: 'es',
    },
    define: { 'process.env': {} },
    resolve: {
        alias: {
            '@': fileURLToPath(new URL('./src', import.meta.url)),
        },
        extensions: ['.js', '.json', '.jsx', '.mjs', '.ts', '.tsx', '.vue'],
    },
    server: {
        port: 3000,
    },
    css: {
        preprocessorOptions: {
            sass: {
                api: 'modern-compiler',
            },
            scss: {
                api: 'modern-compiler',
            },
        },
    },
})
