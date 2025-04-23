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
                <span class="text-subtitle-1">Processing cases {{ curr }} / {{ entities.length }}...</span>
            </div>
            <div v-else class="text-center d-flex flex-column align-center justify-center" style="min-height: 35vh">
                <v-icon size="90" color="primary" class="mb-4" icon="mdi-cloud-upload" />
                <v-btn color="primary" size="large" variant="elevated" @click="open">Browse files</v-btn>
            </div>
        </v-card-item>
    </v-card>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { compressFilesToZip, listLeafFolders, processDirectory } from '@/utils'
import { createTask, getUploadUrl, updateTask } from '@/apis'
import { uploadToS3 } from '@/apis/common.ts'
import { TaskStatus } from '@/types'

const curr = ref(0)
const isUploading = ref(false)
const emit = defineEmits<{ onUploaded: [] }>()
const entities = ref<FileSystemDirectoryHandle[]>([])

const open = async () => {
    try {
        //@ts-ignore FileSystem APIs
        const dirHandle = await window.showDirectoryPicker()
        entities.value = await listLeafFolders(dirHandle)
        isUploading.value = true
        await upload()
    } catch (error) {
        console.error(error)
    }
}

const upload = async () => {
    isUploading.value = true
    for (const entity of entities.value) {
        try {
            const { files, series } = await processDirectory(entity)
            if (!series) continue

            const task = await createTask({ series })
            if (!task) continue

            const uploadInfo = await getUploadUrl({ series })
            if (!uploadInfo) continue

            const file = await compressFilesToZip(files)
            await uploadToS3(uploadInfo.url, uploadInfo.form, file)
            task.status = TaskStatus.Pending
            await updateTask(task)
            curr.value++
        } catch (e) {
            console.error(e)
        }
    }
    emit('onUploaded')
    isUploading.value = false
}
</script>
