<template>
    <v-card class="mx-auto">
        <v-card-title>
            <v-icon icon="mdi-list-box-outline" class="mr-2" />
            <span class="text-h6">Tasks</span>
        </v-card-title>
        <v-card-item>
            <v-row class="align-center justify-space-between mb-2" no-gutters>
                <v-col cols="12" md="4">
                    <v-text-field
                        v-model="series"
                        density="compact"
                        label="Search series"
                        prepend-inner-icon="mdi-magnify"
                        variant="solo-filled"
                        flat
                        hide-details
                        single-line
                    />
                </v-col>
                <v-col cols="12" md="4">
                    <v-select
                        v-model="status"
                        flat
                        chips
                        single-line
                        hide-details
                        prepend-inner-icon="mdi-filter-outline"
                        :items="STATUS_LIST"
                        item-title="label"
                        item-value="value"
                        variant="solo-filled"
                        density="compact"
                        label="Filter by status"
                        clearable
                    />
                </v-col>
                <v-col cols="12" md="auto">
                    <v-btn
                        prepend-icon="mdi-export-variant"
                        color="primary"
                        variant="tonal"
                        @click="exportData"
                        :loading="isExport"
                    >
                        Export
                    </v-btn>
                </v-col>
                <v-col cols="12" md="auto">
                    <v-btn
                        prepend-icon="mdi-refresh"
                        color="secondary"
                        variant="tonal"
                        @click="refresh"
                    >
                        Refresh
                    </v-btn>
                </v-col>
            </v-row>
        </v-card-item>
        <v-data-table-server
            v-model:items-per-page="itemsPerPage"
            :headers="headers"
            :items="serverItems"
            :items-length="totalItems"
            :loading="loading"
            :search="search"
            show-current-page
            item-value="name"
            @update:options="loadItems"
        >
            <template v-slot:item.status="{ value }">
                <v-chip :color="STATUS_COLOR[value as TaskStatus]" :text="STATUS_NAME[value as TaskStatus]"></v-chip>
            </template>
            <template v-slot:item.updated="{ value }">
                {{ dayjs(value).format('YYYY-MM-DD HH:mm:ss') }}
            </template>
            <template v-slot:item.actions="{ item }">
                <div class="d-flex ga-2">
                    <v-icon
                        color="medium-emphasis"
                        icon="mdi-image-multiple"
                        size="small"
                        @click="viewItem(item.id)"
                        :disabled="item.status === 0"
                    />
                    <v-menu :disabled="item.status === 2">
                        <template v-slot:activator="{ props }">
                            <v-icon v-bind="props" color="medium-emphasis" icon="mdi-dots-vertical" size="small" />
                        </template>
                        <v-list>
                            <v-list-item @click="prioritizeItem(item as DeleteTaskRequest)">
                                <template v-slot:prepend>
                                    <v-icon icon="mdi-order-bool-descending-variant" />
                                </template>
                                <v-list-item-title>Prioritize</v-list-item-title>
                            </v-list-item>
                            <v-list-item @click="deleteItem(item as DeleteTaskRequest)">
                                <template v-slot:prepend>
                                    <v-icon icon="mdi-delete" />
                                </template>
                                <v-list-item-title>Delete</v-list-item-title>
                            </v-list-item>
                        </v-list>
                    </v-menu>
                </div>
            </template>
            <!--            <template v-slot:no-data>-->
            <!--                <v-btn-->
            <!--                    prepend-icon="mdi-backup-restore"-->
            <!--                    color="primary"-->
            <!--                    variant="tonal"-->
            <!--                    border-->
            <!--                    @click="loadItems({ page: 1, itemsPerPage: 10 })"-->
            <!--                >-->
            <!--                    Reset data-->
            <!--                </v-btn>-->
            <!--            </template>-->
        </v-data-table-server>
    </v-card>
</template>

<script setup lang="ts">
import { onUnmounted, ref } from 'vue'
import { deleteTask, getTaskList, prioritizeTask } from '@/apis'
import {
    type DeleteTaskRequest,
    STATUS_COLOR,
    STATUS_LIST,
    STATUS_NAME,
    type Task,
    TaskStatus,
    type UpdateTaskRequest,
} from '@/types'
import { type DataTableHeader } from 'vuetify'
import dayjs from 'dayjs'
import { VDataTableServer } from 'vuetify/components'
import { useIntervalFn, watchDebounced } from '@vueuse/core'
import { useRouter } from 'vue-router'
import { parseExportUrl } from '@/apis/common.ts'
import { save } from '@tauri-apps/plugin-dialog'
import { invoke } from '@tauri-apps/api/core'
import { toast } from 'vue-sonner'

const isExport = ref(false)
const router = useRouter()
const itemsPerPage = ref(10)
const totalItems = ref(0)
const loading = ref(false)
const headers: DataTableHeader[] = [
    { title: 'ID', align: 'start', sortable: true, key: 'id' },
    { title: 'Series', align: 'start', sortable: false, key: 'series' },
    { title: 'Status', align: 'start', sortable: false, key: 'status' },
    { title: 'Model', align: 'start', sortable: false, key: 'model' },
    { title: 'Order', align: 'start', sortable: true, key: 'order' },
    { title: 'Updated', align: 'start', sortable: true, key: 'updated' },
    { title: 'Actions', align: 'start', sortable: false, key: 'actions' },
]
const search = ref('')
const series = ref('')
const status = ref<TaskStatus | undefined>(undefined)
const serverItems = ref<Task[]>([])

const triggerLoad = () => (search.value = Date.now().toString())
const refresh = () => triggerLoad()

watchDebounced([series, status], () => triggerLoad(), { debounce: 500, maxWait: 1000 })

defineExpose({ refresh })

const { pause, resume } = useIntervalFn(async () => refresh(), 3000, { immediate: false })

const viewItem = async (id: number) => {
    await router.replace({ name: 'Viewer', params: { id } })
}

const deleteItem = async (item: DeleteTaskRequest) => {
    const res = await deleteTask(item)
    if (res) {
        triggerLoad()
    }
}

const prioritizeItem = async (item: UpdateTaskRequest) => {
    const res = await prioritizeTask(item)
    if (res) {
        triggerLoad()
    }
}

const exportData = async () => {
    const path = await save({
        title: `Exports tasks`,
        defaultPath: `tasks_${dayjs(new Date()).format('YYYY-MM-DD')}`,
        filters: [
            {
                name: `CSV`,
                extensions: ['csv'],
            },
        ],
    })
    if (!path) return

    isExport.value = true
    const url = parseExportUrl(status.value, series.value)
    const res = await invoke<boolean>('export', { path, url })
    if (res) {
        toast.success(`Tasks exported to ${path}`)
    }
    isExport.value = false
}

const loadItems = async ({
    page,
    itemsPerPage,
    sortBy,
}: {
    page: number
    itemsPerPage: number
    sortBy?: { key: string; order: string }[]
}) => {
    pause()
    loading.value = true
    const sortKey = sortBy && sortBy[0]?.key
    const sortOrder = sortBy && sortBy[0]?.order

    const res = await getTaskList({
        page,
        pageSize: itemsPerPage,
        status: status.value,
        series: series.value,
        sortKey,
        sortOrder,
    })
    if (res) {
        serverItems.value = res.tasks
        totalItems.value = res.total
    }
    loading.value = false
    resume()
}

onUnmounted(() => pause())
</script>