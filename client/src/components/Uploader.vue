<template>
    <v-card class="mx-auto" :loading="isUploading">
        <div v-if="!isUploading">
            <v-card-title class="text-h5">Select CT series</v-card-title>
            <v-divider />
            <v-card-item>
                <div
                    v-if="isProcessing"
                    class="text-center d-flex flex-column align-center justify-center"
                    style="height: 35vh"
                >
                    <v-progress-circular indeterminate size="70" color="primary" class="mb-4" />
                    <span class="text-subtitle-1">Processing files...</span>
                </div>
                <div v-else class="text-center d-flex flex-column align-center justify-center" style="height: 35vh">
                    <v-icon size="90" color="primary" class="mb-4" icon="mdi-cloud-upload" />
                    <v-btn color="primary" size="large" variant="elevated" @click="open">Browse files</v-btn>
                </div>
            </v-card-item>
        </div>
        <div v-else>
            <v-card-title class="text-h5">Uploading CT series</v-card-title>
            <v-divider />
            <v-card-item>
                <v-virtual-scroll :items="groupedSeries" height="35vh">
                    <template v-slot:default="{ item }">
                        <v-list-item density="compact" :subtitle="`${item.files.length} files`" :title="item.series">
                            <template v-slot:prepend>
                                <v-icon>mdi-file</v-icon>
                            </template>
                            <template v-slot:append></template>
                            <v-progress-linear
                                :model-value="(item.index.value / item.files.length) * 100"
                                color="success"
                            />
                        </v-list-item>
                    </template>
                </v-virtual-scroll>
            </v-card-item>
        </div>
    </v-card>
</template>

<script setup lang="ts">
import { type Ref, ref } from 'vue'
import { getSeries } from '@/utils/dicom.ts'
import { createTask, getUploadUrl, updateTask } from '@/apis'
import { uploadToS3 } from '@/apis/common.ts'
import { TaskState } from '@/types'

const curr = ref('')
const isProcessing = ref(false)
const isUploading = ref(false)
const emit = defineEmits<{ onUploaded: [] }>()
let groupedSeries: { series: string; files: File[]; index: Ref<number> }[]

const open = async () => {
    try {
        //@ts-ignore FileSystem APIs
        const dirHandle = await window.showDirectoryPicker()
        isProcessing.value = true
        const groupedFiles: Record<string, File[]> = {}
        await processDirectory(dirHandle, groupedFiles)
        groupedSeries = Object.entries(groupedFiles).map(([series, files]) => ({
            series,
            files,
            index: ref(0),
        }))
        isProcessing.value = false
        isUploading.value = true
        await upload()
    } catch (error) {
        console.error(error)
        isProcessing.value = false
    }
}

async function processDirectory(dirHandle: FileSystemDirectoryHandle, groupedFiles: Record<string, File[]>) {
    //@ts-ignore FileSystem APIs
    for await (const entry of dirHandle.values()) {
        if (entry.kind === 'file') {
            const file = await entry.getFile()
            const series = await getSeries(file)
            if (series) {
                if (!groupedFiles[series]) {
                    groupedFiles[series] = []
                }
                groupedFiles[series].push(file)
            }
        } else if (entry.kind === 'directory') {
            await processDirectory(entry, groupedFiles)
        }
    }
}

const upload = async () => {
    for (let i = 0; i < groupedSeries.length; i++) {
        const group = groupedSeries[i]
        curr.value = group.series

        const task = await createTask({ series: group.series })
        if (!task) continue

        const uploadInfo = await getUploadUrl({ series: group.series })
        if (!uploadInfo) continue
        const key = uploadInfo.form.key

        let j = 0
        const CHUNK_SIZE = 6
        while (j < group.files.length) {
            const chunk = group.files.slice(j, j + CHUNK_SIZE)
            await Promise.all(
                chunk.map(async (file) => {
                    uploadInfo.form.key = `${key}${file.name}`
                    await uploadToS3(uploadInfo.url, uploadInfo.form, file)
                })
            )
            group.index.value += CHUNK_SIZE
            j += CHUNK_SIZE
        }
        task.status = TaskState.Pending
        await updateTask(task)
        emit('onUploaded')
    }
    await new Promise((resolve) => setTimeout(resolve, 3000))
    isUploading.value = false
}
</script>
