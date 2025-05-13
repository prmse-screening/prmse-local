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
            <v-btn icon>
                <v-icon>mdi-file-chart-outline</v-icon>
                <v-dialog activator="parent" persistent max-width="66vw">
                    <template v-slot:default="{ isActive }">
                        <ResultCard @close="isActive.value = false" :id="taskId" />
                    </template>
                </v-dialog>
            </v-btn>
        </v-app-bar>

        <!-- Main Content Area -->
        <v-main>
            <v-container
                fluid
                class="d-flex justify-center align-center"
                height="100%"
                @click.right.prevent
                @click.middle.prevent
            >
                <splitpanes class="default-theme" vertical @resized="resize">
                    <!-- Left Side (CT Visualization) -->
                    <pane min-size="20">
                        <div :id="`CT_AXIAL_${taskId}`" class="d-flex align-center justify-center fill-height rounded">
                            <v-progress-circular v-if="loading" indeterminate />
                        </div>
                    </pane>

                    <!-- Right Side (Further Split) -->
                    <pane min-size="20">
                        <splitpanes class="default-theme" horizontal @resized="resize">
                            <!-- Upper Part (Right Top) -->
                            <pane min-size="20">
                                <div
                                    :id="`CT_SAGITTAL_${taskId}`"
                                    class="d-flex align-center justify-center fill-height rounded"
                                >
                                    <v-progress-circular v-if="loading" indeterminate />
                                </div>
                            </pane>

                            <!-- Lower Part (Right Bottom) -->
                            <pane min-size="20">
                                <div
                                    :id="`CT_CORONAL_${taskId}`"
                                    class="d-flex align-center justify-center fill-height rounded"
                                >
                                    <v-progress-circular v-if="loading" indeterminate />
                                </div>
                            </pane>
                        </splitpanes>
                    </pane>
                </splitpanes>
            </v-container>
        </v-main>
    </v-app>
</template>

<script setup lang="ts">
import { Splitpanes, Pane } from 'splitpanes'
import 'splitpanes/dist/splitpanes.css'

import { onMounted, onUnmounted, ref } from 'vue'
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
import { deleteImageIds, generateImageIds, prefetchMetadataInformation, processS3ZipFile } from '@/utils'
import { getTask } from '@/apis'
import { wadouri } from '@cornerstonejs/dicom-image-loader'
import ResultCard from '@/components/ResultCard.vue'
import type IToolGroup from '@cornerstonejs/tools/types/IToolGroup'

const loading = ref(false)
const router = useRouter()

const taskId = router.currentRoute.value.params.id as string
const renderingEngineId = `renderingEngine:${taskId}`
const viewportIds = [`CT_AXIAL_${taskId}`, `CT_SAGITTAL_${taskId}`, `CT_CORONAL_${taskId}`]
const volumeId = `volume:${taskId}`
const toolGroupId = `group:${taskId}`
let imageIds = []

const initTools = (): IToolGroup | undefined => {
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
    return toolGroup
}
const render = async (imageIds: string[]) => {
    const toolGroup = initTools()
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
            element: document.getElementById(`CT_AXIAL_${taskId}`) as HTMLDivElement,
            defaultOptions: { orientation: csEnums.OrientationAxis.AXIAL },
        },
        {
            viewportId: viewportIds[1],
            type: csEnums.ViewportType.ORTHOGRAPHIC,
            element: document.getElementById(`CT_SAGITTAL_${taskId}`) as HTMLDivElement,
            defaultOptions: { orientation: csEnums.OrientationAxis.SAGITTAL },
        },
        {
            viewportId: viewportIds[2],
            type: csEnums.ViewportType.ORTHOGRAPHIC,
            element: document.getElementById(`CT_CORONAL_${taskId}`) as HTMLDivElement,
            defaultOptions: { orientation: csEnums.OrientationAxis.CORONAL },
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
    loading.value = true
    const res = await getTask(router.currentRoute.value.params.id as string)
    if (res) {
        const files = await processS3ZipFile(res.id)
        loading.value = false
        if (!files) return
        imageIds = generateImageIds(files)
        await render(imageIds)
    }
    loading.value = false
})

onUnmounted(() => {
    wadouri.fileManager.purge()
    ToolGroupManager.destroyToolGroup(toolGroupId)
    getRenderingEngine(renderingEngineId)?.destroy()
    cache.purgeCache()
    cache.purgeVolumeCache()
})
</script>
