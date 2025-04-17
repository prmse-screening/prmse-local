<template>
    <v-data-table-server
        v-model:items-per-page="itemsPerPage"
        :headers="headers"
        :items="serverItems"
        :items-length="totalItems"
        :loading="loading"
        :search="search"
        item-value="name"
        @update:options="loadItems"
    />
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { getTaskList } from '@/apis'
import type { Task } from '@/types'
import type { DataTableHeader } from 'vuetify'

const itemsPerPage = ref(10)
const loading = ref(false)
const headers: DataTableHeader[] = [
    {
        title: 'ID',
        align: 'start',
        sortable: false,
        key: 'id',
    },
    {
        title: 'Series',
        align: 'start',
        sortable: false,
        key: 'series',
    },
    {
        title: 'Status',
        align: 'start',
        sortable: false,
        key: 'status',
    },
    {
        title: 'Model',
        align: 'start',
        sortable: false,
        key: 'model',
    },
    {
        title: 'Order',
        align: 'start',
        sortable: true,
        key: 'order',
    },
    {
        title: 'Updated',
        align: 'start',
        sortable: true,
        key: 'updated',
    },
]
const search = ref('')
const serverItems = ref<Task[]>([])
const totalItems = ref(0)

const loadItems = async ({ page, itemsPerPage }: { page: number; itemsPerPage: number }) => {
    loading.value = true
    const res = await getTaskList({
        page,
        pageSize: itemsPerPage,
    })
    console.log(res)
    if (res) {
        serverItems.value = res.tasks
        totalItems.value = res.total
    }
    loading.value = false
}
</script>

<style scoped></style>
