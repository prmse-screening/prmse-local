<template>
    <v-card class="mx-auto">
        <v-card-title>
            <v-icon icon="mdi-list-box-outline" class="mr-2" />
            <span class="text-h6">Tasks</span>
        </v-card-title>
        <v-card-item>
            <v-row class="align-center justify-space-between" no-gutters>
                <v-col cols="12" md="2" class="d-flex align-center mb-2 mb-md-0">
                    <v-btn
                        prepend-icon="mdi-export-variant"
                        style="width: 100%"
                        color="primary"
                        variant="tonal"
                        @click="exportData"
                    >
                        Export
                    </v-btn>
                </v-col>
                <v-col cols="12" md="6" class="mb-2 mb-md-0">
                    <v-text-field
                        v-model="series"
                        density="compact"
                        label="Search"
                        prepend-inner-icon="mdi-magnify"
                        variant="solo-filled"
                        flat
                        hide-details
                        single-line
                    />
                </v-col>
                <v-col cols="12" md="3">
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
                        style="width: 100%"
                    />
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

            <template v-slot:no-data>
                <v-btn
                    prepend-icon="mdi-backup-restore"
                    color="primary"
                    variant="tonal"
                    border
                    @click="loadItems({ page: 1, itemsPerPage: 10 })"
                >
                    Reset data
                </v-btn>
            </template>
        </v-data-table-server>
    </v-card>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
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
import { watchDebounced } from '@vueuse/core'
import { useRouter } from 'vue-router'
import { parseExportUrl } from '@/apis/common.ts'

const router = useRouter()
const itemsPerPage = ref(10)
const totalItems = ref(0)
const loading = ref(false)
const headers: DataTableHeader[] = [
    { title: 'ID', align: 'start', sortable: false, key: 'id' },
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

watchDebounced(series, () => (search.value = Date.now().toString()), { debounce: 500, maxWait: 1000 })
watch(status, () => (search.value = Date.now().toString()))

const viewItem = async (id: number) => {
    await router.replace({ name: 'Viewer', params: { id } })
}

const deleteItem = async (item: DeleteTaskRequest) => {
    const res = await deleteTask(item)
    if (res) {
        search.value = Date.now().toString()
    }
}

const prioritizeItem = async (item: UpdateTaskRequest) => {
    const res = await prioritizeTask(item)
    if (res) {
        search.value = Date.now().toString()
    }
}

const exportData = async () => {
    const url = parseExportUrl(status.value, series.value)
    const link = document.createElement('a')
    link.href = url
    link.download = `tasks_${new Date().toLocaleString()}.csv`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
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
}
</script>

<style scoped></style>
