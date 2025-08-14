import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'
import svgLoader from 'vite-svg-loader'

const rootPath = new URL('.', import.meta.url).pathname

// https://vitejs.dev/config/
export default defineConfig({
    server: {
        hmr: {
            host: 'localhost',
            protocol: 'ws',
        },
    },
    plugins: [
        vue(),
        AutoImport({
            resolvers: [ElementPlusResolver()],
        }),
        Components({
            resolvers: [ElementPlusResolver()],
        }),
        svgLoader(),
    ],
    resolve: {
        alias: {
            '@': rootPath + 'src',
            'wailsjs': rootPath + 'wailsjs',
        },
    },
    build: {
        target: ['chrome112', 'firefox113', 'safari16.5', 'edge112']
    }
})
