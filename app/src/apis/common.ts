import { toast } from 'vue-sonner'
import type { BaseResponse, S3UploadForm } from '@/types'
import { fetch } from '@tauri-apps/plugin-http'
import { invoke } from '@tauri-apps/api/core'
import { useRouter } from 'vue-router'
import router from "@/router";

export const BASE_URL = 'http://localhost:8080'
const request = async <T>(url: string, method: 'GET' | 'POST', body?: any): Promise<T | null> => {
    try {
        const headers: Record<string, string> = {}
        if (body) headers['Content-Type'] = 'application/json'
        const response = await fetch(`${localStorage.getItem('base')}${url}`, {
            headers,
            method,
            body: body ?? null,
            connectTimeout: 3000
        })
        const json: BaseResponse = await response.json()
        if (json.status === 'error') {
            toast.error(json.data as string)
            return null
        }
        return json.data as T
    } catch (e: any) {
        console.error('HTTP Error:', e)
        toast.error('Failed to connect server')
        localStorage.removeItem('base')
        await router.replace({ name: 'Config' })
        return null
    }
}

export const GET = <T>(url: string): Promise<T | null> => request(url, 'GET')
export const POST = <T>(url: string, body: any): Promise<T | null> => request(url, 'POST', body)

export const ping = async () => {
    const res = await GET<string>('/ping')
    return res === 'pong'
}

// export const uploadToS3 = async (url: string, form: S3UploadForm, file: File): Promise<boolean> => {
//     const formData = new FormData()
//     Object.entries(form).forEach(([key, value]) => {
//         formData.append(key, value)
//     })
//
//     formData.append('file', file)
//
//     try {
//         const res = await fetch(url, {
//             method: 'POST',
//             body: formData,
//         })
//
//         if (res.ok) {
//             return true
//         } else {
//             console.error('Upload failed:', await res.text())
//             return false
//         }
//     } catch (error) {
//         console.error('Upload error:', error)
//         return false
//     }
// }

export const uploadToS3 = async (url: string, form: S3UploadForm, folder: string): Promise<boolean> => {
    const res = await invoke('upload', { url, form, folder })
    if (res === 'success') {
        return true
    } else {
        console.error('Upload failed:', res)
        return false
    }
}

export const parseDownloadUrl = (id: number) => `${localStorage.getItem('base')}/dicom/${id}`
export const parseExportUrl = (status?: number, series?: string) => {
    const params: Record<string, string> = {}
    if (status != null) params.status = status.toString()
    if (series && series != '') params.series = series.trim()

    if (Object.keys(params).length === 0) return `${BASE_URL}/tasks/export`
    return `${BASE_URL}/tasks/export?${new URLSearchParams(params)}`
}
