<template>
    <v-app>
        <v-main>
            <v-container height="100%" class="d-flex justify-center align-center">
                <v-card class="pa-4" width="35vw">
                    <v-card-title class="text-center text-h5 mb-4">Server configuration</v-card-title>
                    <v-card-text>
                        <v-form @submit.prevent="testConnection" ref="form">
                            <v-text-field
                                v-model="backendUrl"
                                label="Server url"
                                placeholder="For example: https://example.com:8080"
                                :rules="[rules.required, rules.validUrl]"
                                variant="outlined"
                                :error-messages="errorMessage"
                                clearable
                            />
                            <v-btn
                                block
                                color="primary"
                                size="large"
                                type="submit"
                                :loading="loading"
                                class="mt-4"
                            >
                                Try to connect
                            </v-btn>
                        </v-form>
                    </v-card-text>
                </v-card>
            </v-container>
        </v-main>
    </v-app>
</template>

<script setup lang="ts">
import { ref, reactive, useTemplateRef } from 'vue'
import { useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import { ping } from '@/apis'
import type { VForm } from 'vuetify/components'

const router = useRouter()
const form = useTemplateRef<VForm>('form')
const backendUrl = ref('http://localhost:8080')
const errorMessage = ref('')
const loading = ref(false)

const rules = reactive({
    required: (value: string) => !!value || 'Please input a url',
    validUrl: (value: string) => {
        const pattern = /^(https?:\/\/)([\w.-]+)(:[0-9]+)?(\/\S*)?$/i
        return pattern.test(value) || 'Please input a valid url'
    },
})

const testConnection = async () => {
    const validateRes = await form.value?.validate()

    if (!validateRes?.valid) {
        return
    }

    loading.value = true
    errorMessage.value = ''

    try {
        const url = backendUrl.value
        if (!url.startsWith('http://') && !url.startsWith('https://')) {
            return
        }

        localStorage.setItem('base', url)

        const response = await ping()

        if (response) {
            toast.success('Connection successful!')
            await router.replace({ name: 'Home' })
        } else {
            throw new Error(`The server returns an error.`)
        }
    } catch (error: any) {
        localStorage.removeItem('base')
        errorMessage.value = 'Connection failed, please check if the URL is correct'
    } finally {
        loading.value = false
    }
}
</script>
