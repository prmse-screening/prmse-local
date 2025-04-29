<template>
    <v-container fluid min-height="100vh">
        <v-row class="d-flex flex-column" no-gutters>
            <v-col cols="12" class="mb-4">
                <uploader @uploaded="table?.refresh" />
            </v-col>
            <v-col cols="12">
                <task-table ref="table" />
            </v-col>
        </v-row>
    </v-container>
</template>

<script setup lang="ts">
import Uploader from '@/components/Uploader.vue'
import TaskTable from '@/components/TaskTable.vue'
import { useTemplateRef } from 'vue'

const table = useTemplateRef('table')
</script>
