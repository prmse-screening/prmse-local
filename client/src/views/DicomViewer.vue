<template>
    <v-app>
        <!-- Top App Bar -->
        <v-app-bar>
            <v-btn icon @click="router.replace({ name: 'Home' })">
                <v-icon>mdi-home</v-icon>
            </v-btn>
            <v-spacer></v-spacer>
            <v-btn icon @click="resetCamera">
                <v-icon>mdi-refresh</v-icon>
            </v-btn>
        </v-app-bar>

        <!-- Main Content Area -->
        <v-main>
            <v-container fluid height="100%" @click.right.prevent @click.middle.prevent>
                <!-- Splitpanes for CT Visualization -->
                <splitpanes class="default-theme" vertical @resized="resize">
                    <!-- Left Side (CT Visualization) -->
                    <pane min-size="20">
                        <div id="axial" class="bg-primary fill-height rounded"></div>
                    </pane>

                    <!-- Right Side (Further Split) -->
                    <pane min-size="20">
                        <splitpanes class="default-theme" horizontal @resized="resize">
                            <!-- Upper Part (Right Top) -->
                            <pane min-size="20">
                                <div id="sagittal" class="bg-red fill-height rounded"></div>
                            </pane>

                            <!-- Lower Part (Right Bottom) -->
                            <pane min-size="20">
                                <div id="coronal" class="bg-green fill-height rounded"></div>
                            </pane>
                        </splitpanes>
                    </pane>
                </splitpanes>
            </v-container>

            <!-- Floating Action Button -->
            <v-fab icon="mdi-lightbulb-on" color="primary" location="right bottom" app></v-fab>
        </v-main>
    </v-app>
</template>

<script setup lang="ts">
import { Splitpanes, Pane } from 'splitpanes'
import 'splitpanes/dist/splitpanes.css'

import { onBeforeUnmount, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useDebounceFn } from '@vueuse/core'
import {
    RenderingEngine,
    Enums as csEnums,
    volumeLoader,
    setVolumesForViewports,
    cache,
    getRenderingEngine,
} from '@cornerstonejs/core'
import {
    addTool,
    ToolGroupManager,
    Enums as csToolsEnums,
    StackScrollTool,
    WindowLevelTool,
    ZoomTool,
    PanTool,
} from '@cornerstonejs/tools'
import { generateImageIds, prefetchMetadataInformation, processS3ZipFile } from '@/utils'
import { getTask } from '@/apis'
import { wadouri } from '@cornerstonejs/dicom-image-loader'

const router = useRouter()

const renderingEngineId = 'renderingEngine:1'
const viewportIds = ['CT_AXIAL', 'CT_SAGITTAL', 'CT_CORONAL']
const volumeId = 'volume:1'
const toolGroupId = 'group:1'

const render = async (imageIds: string[]) => {
    const toolGroup = ToolGroupManager.createToolGroup(toolGroupId)
    addTool(StackScrollTool)
    addTool(WindowLevelTool)
    addTool(ZoomTool)
    addTool(PanTool)
    toolGroup?.addTool(StackScrollTool.toolName)
    toolGroup?.addTool(WindowLevelTool.toolName)
    toolGroup?.addTool(ZoomTool.toolName)
    toolGroup?.addTool(PanTool.toolName)
    toolGroup?.setToolActive(StackScrollTool.toolName, {
        bindings: [
            {
                mouseButton: csToolsEnums.MouseBindings.Wheel,
            },
        ],
    })
    toolGroup?.setToolActive(WindowLevelTool.toolName, {
        bindings: [
            {
                mouseButton: csToolsEnums.MouseBindings.Primary,
            },
        ],
    })
    toolGroup?.setToolActive(ZoomTool.toolName, {
        bindings: [
            {
                mouseButton: csToolsEnums.MouseBindings.Secondary,
            },
        ],
    })
    toolGroup?.setToolActive(PanTool.toolName, {
        bindings: [
            {
                mouseButton: csToolsEnums.MouseBindings.Auxiliary,
            },
        ],
    })

    await prefetchMetadataInformation(imageIds)
    if (cache.getVolume(volumeId)) {
        cache.removeVolumeLoadObject(volumeId)
    }

    const renderingEngine = new RenderingEngine(renderingEngineId)
    viewportIds.forEach((viewportId) => toolGroup?.addViewport(viewportId, renderingEngineId))

    const viewportInputArray = [
        {
            viewportId: viewportIds[0],
            type: csEnums.ViewportType.ORTHOGRAPHIC,
            element: document.getElementById('axial') as HTMLDivElement,
            defaultOptions: {
                orientation: csEnums.OrientationAxis.AXIAL,
            },
        },
        {
            viewportId: viewportIds[1],
            type: csEnums.ViewportType.ORTHOGRAPHIC,
            element: document.getElementById('sagittal') as HTMLDivElement,
            defaultOptions: {
                orientation: csEnums.OrientationAxis.SAGITTAL,
            },
        },
        {
            viewportId: viewportIds[2],
            type: csEnums.ViewportType.ORTHOGRAPHIC,
            element: document.getElementById('coronal') as HTMLDivElement,
            defaultOptions: {
                orientation: csEnums.OrientationAxis.CORONAL,
            },
        },
    ]
    renderingEngine.setViewports(viewportInputArray)
    const volume = await volumeLoader.createAndCacheVolume(volumeId, { imageIds })
    volume.load()

    await setVolumesForViewports(renderingEngine, [{ volumeId }], viewportIds)
}

const resetCamera = () =>
    getRenderingEngine(renderingEngineId)
        ?.getViewports()
        ?.forEach((viewport) => {
            viewport?.resetCamera({
                resetPan: true,
                resetZoom: true,
                resetToCenter: true,
                storeAsInitialCamera: false,
            })
            viewport?.render()
        })

const resize = useDebounceFn(() => {
    const renderingEngine = getRenderingEngine(renderingEngineId)
    const viewports = renderingEngine?.getViewports()
    if (!renderingEngine || !viewports) {
        return
    }
    const presentations = viewports.map((viewport) => viewport.getViewPresentation())
    renderingEngine.resize(true, false)
    viewports.forEach((viewport, index) => {
        viewport.setViewPresentation(presentations[index])
    })
}, 350)

onMounted(async () => {
    const id = router.currentRoute.value.params.id
    const res = await getTask(id as string)
    if (res) {
        const files = await processS3ZipFile(res.id)
        console.log(files)
        if (!files) return
        const imageIds = generateImageIds(files)
        await render(imageIds)
    }
})

onBeforeUnmount(() => {
    wadouri.fileManager.purge()
    ToolGroupManager.destroyToolGroup(toolGroupId)
    cache.purgeCache()
    getRenderingEngine(renderingEngineId)?.destroy()
})
</script>
