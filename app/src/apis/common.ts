import { toast } from 'vue-sonner'
import type { BaseResponse } from '@/types'
import { fetch } from '@tauri-apps/plugin-http'

export const BASE_URL = 'http://localhost:8080'
const request = async <T = Data>(url: string, method: 'GET' | 'POST', body?: any): Promise<T | null> => {
    try {
        if (body) headers['Content-Type'] = 'application/json'
        const response = await fetch(`${localStorage.getItem('base')}${url}`, {
            headers,
            method,
            body: body ? JSON.stringify(body) : null,
        })
        const json: BaseResponse = await response.json()
        if (json.status === 'error') {
            toast.error(res.data as string)
            return null
        }
        return json.data as T
    } catch (e) {
        console.error('HTTP Error:', error)
        toast.error(error.message || 'Network error')
        return null
    }
}

export const GET = (url: string): Promise<Data | null> => request(url, 'GET')
export const POST = (url: string, body: any): Promise<Data | null> => request(url, 'POST', body)

type S3UploadForm = Record<string, string>
export const uploadToS3 = async (url: string, form: S3UploadForm, file: File): Promise<boolean> => {
    const formData = new FormData()
    Object.entries(form).forEach(([key, value]) => {
        formData.append(key, value)
    })

    formData.append('file', file)

    try {
        const res = await fetch(url, {
            method: 'POST',
            body: formData,
        })

        if (res.ok) {
            return true
        } else {
            console.error('Upload failed:', await res.text())
            return false
        }
    } catch (error) {
        console.error('Upload error:', error)
        return false
    }
}

export const parseDownloadUrl = (id: number) => `${BASE_URL}/dicom/${id}`
export const parseExportUrl = (status?: number, series?: string) => {
    const params: Record<string, string> = {}
    if (status != null) params.status = status.toString()
    if (series && series != '') params.series = series.trim()

    if (Object.keys(params).length === 0) return `${BASE_URL}/tasks/export`
    return `${BASE_URL}/tasks/export?${new URLSearchParams(params)}`
}
