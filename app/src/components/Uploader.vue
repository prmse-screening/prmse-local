<template>
    <v-card :loading="isUploading">
        <v-card-title class="text-h5">Select CT series</v-card-title>
        <v-divider />
        <v-card-item>
            <div
                v-if="isUploading"
                class="text-center d-flex flex-column align-center justify-center"
                style="min-height: 35vh"
            >
                <v-progress-circular indeterminate size="66" color="primary" class="mb-4" />
                <span class="text-subtitle-1">Processing cases {{ curr }} / {{ totalTasks }}...</span>
            </div>
            <div v-else class="text-center d-flex flex-column align-center justify-center" style="min-height: 35vh">
                <v-icon size="90" color="primary" class="mb-4" icon="mdi-cloud-upload" />
                <v-btn color="primary" size="large" variant="elevated" @click="openDialog">Browse files</v-btn>
            </div>
        </v-card-item>
    </v-card>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { listLeafFolders, processDirectory } from '@/utils'
import { createTask, getUploadPostUrl, updateTask } from '@/apis'
import { uploadToS3 } from '@/apis/common.ts'
import { open } from '@tauri-apps/plugin-dialog'
import { TaskStatus } from '@/types'

const curr = ref(1)
const isUploading = ref(false)
const emit = defineEmits<{ (e: 'uploaded'): void }>()
const totalTasks = ref(0)

const openDialog = async () => {
    try {
        const dir = await open({
            multiple: false,
            directory: true,
            openLabel: 'Select a folder',
        })
        if (!dir) return
        const entities = await listLeafFolders(dir)
        if (!entities) return
        totalTasks.value = entities.length
        await upload(entities)
    } catch (error) {
        console.error(error)
    }
}

const upload = async (entities: string[]) => {
    isUploading.value = true
    for (const entity of entities) {
        try {
            console.log(entity)
            const series = await processDirectory(entity)
            if (!series) continue

            const task = await createTask({ series })
            if (!task) continue

            const uploadInfo = await getUploadPostUrl({ series })
            if (!uploadInfo) continue

            const uploadRes = await uploadToS3(uploadInfo.url, uploadInfo.form, entity)
            if (!uploadRes) continue

            task.status = TaskStatus.Pending
            await updateTask(task)
        } catch (err) {
            console.warn(`Entity processing failed:`, err)
        } finally {
            curr.value++
        }
    }
    emit('uploaded')
    isUploading.value = false
    curr.value = 1
}
</script>
