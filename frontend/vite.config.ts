import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite';
import {ArcoResolver} from 'unplugin-vue-components/resolvers';
import { vitePluginForArco } from '@arco-plugins/vite-vue'

export default defineConfig(({mode}) => {
    console.log({mode});
    const isProduction = mode === 'production';
    return {
        plugins: [
            vue(),
            vitePluginForArco({
                style: 'css'
            }),
            AutoImport({
                resolvers: [ArcoResolver()],
            }),
            Components({
                resolvers: [
                    ArcoResolver({
                        sideEffect: true
                    })
                ]
            })
            // compress({
            //     verbose: true,
            //     disable: false,
            //     threshold: 10240, // 只有大小大于此阈值的文件会被压缩
            //     algorithm: 'gzip',
            //     ext: '.gz',
            //     deleteOriginFile: false, // 是否删除原始文件
            // }),
        ],
        build: {
            rollupOptions: {
                output: {
                    entryFileNames: '[name].js',
                    chunkFileNames: '[name].js',
                    assetFileNames: '[name].[ext]',
                    // manualChunks(id) {
                    //     // 如果模块在 node_modules 中，则将其拆分为单独的块
                    //     if (id.includes('node_modules')) {
                    //         return id.split('node_modules/')[1].split('/')[0].toString(); // 以模块名作为块名
                    //     }
                    // },
                    manualChunks(id) {
                        if (id.includes('@arco-design')) {
                            return 'arco'; // 将 Arco Design 组件分到一个单独的块中
                        }
                    },
                }
            }
        },
    }
})
