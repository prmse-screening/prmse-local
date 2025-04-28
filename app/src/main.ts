/**
 * main.ts
 *
 * Bootstraps Vuetify and other plugins then mounts the App`
 */

// Plugins
import { registerPlugins } from '@/plugins'

import App from './App.vue'

import { createApp } from 'vue'
import router from '@/router'
import { init as coreInit } from '@cornerstonejs/core'
import { init as csToolsInit } from '@cornerstonejs/tools'
import { init as dicomImageLoaderInit } from '@cornerstonejs/dicom-image-loader'

coreInit()
csToolsInit()
dicomImageLoaderInit()

const app = createApp(App)

registerPlugins(app)

app.use(router).mount('#app')
