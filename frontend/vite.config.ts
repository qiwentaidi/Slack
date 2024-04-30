import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'

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
        // prismjsPlugin({
	    //     languages: 'all', // 语言
	    //     plugins: ['line-numbers','show-language','copy-to-clipboard','inline-color'],
	    //     theme: 'base16-ateliersulphurpool.light',// 主题
	    //     css: true,
	    // })
    ],
})
