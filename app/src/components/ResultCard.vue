<template>
    <div>
        <v-card rounded="lg">
            <v-card-title class="d-flex justify-space-between align-center">
                <div class="text-h5 text-medium-emphasis ps-2">Result</div>
                <v-btn icon="mdi-close" variant="text" @click="emit('close')"></v-btn>
            </v-card-title>
            <v-divider class="mb-4"></v-divider>
            <v-skeleton-loader v-if="isActive" class="mx-4" type="image"></v-skeleton-loader>
            <v-sheet
                v-else
                class="mx-4"
                :color="health ? 'success' : 'warning'"
                elevation="12"
                max-width="calc(100% - 32px)"
                rounded="lg"
            >
                <v-sparkline :model-value="res?.prediction" line-width="2" padding="16" max="1" min="0" auto-draw>
                    <template #label="{ value }"> {{ (Number(value) * 100).toFixed(3) }}%</template>
                </v-sparkline>
            </v-sheet>

            <v-card-text v-if="isActive" class="pt-10">
                <v-skeleton-loader class="mb-2" type="heading"></v-skeleton-loader>
                <v-skeleton-loader class="mb-2" type="text"></v-skeleton-loader>
            </v-card-text>
            <v-card-text v-else class="pt-10">
                <div class="text-h6 font-weight-light mb-2">Probability of lung cancer in the next six years</div>
                <div class="subheading font-weight-light mb-2" :class="health ? 'text-success' : 'text-warning'">
                    {{ health ? 'No irregularities found in the lung scan.' : 'Manual decision recommended.' }}
                </div>
            </v-card-text>
        </v-card>
    </div>
</template>

<script setup lang="ts">
import { useIntervalFn } from '@vueuse/core'
import { computed, onUnmounted, ref } from 'vue'
import { getTask } from '@/apis'
import type { Result } from '@/types'

const emit = defineEmits<{
    close: []
}>()
const props = defineProps<{
    id: string
}>()

const res = ref<Result>({ prediction: [], threshold: 0 })
const health = computed(() => res.value.prediction[0] < res.value.threshold)

const { pause, isActive } = useIntervalFn(async () => {
    const task = await getTask(props.id)
    if (task && task.result) {
        res.value = task.result
        pause()

    }
}, 1000)

onUnmounted(() => pause())
</script>

<style scoped></style>
