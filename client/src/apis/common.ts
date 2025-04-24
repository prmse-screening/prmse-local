import { toast } from 'vue-sonner'
import type { BaseResponse } from '@/types'

export const BASE_URL = import.meta.env.DEV ? 'http://localhost:8080' : '/api'

export const http = async <T>(url: string, options: RequestInit = {}): Promise<T | null> => {
    try {
        const res = await fetch(BASE_URL + url, {
            headers: {
                ...(options.headers || {}),
            },
            ...options,
        })

        const json = (await res.json()) as BaseResponse<T>

        if (json.status === 'error') {
            toast.error(json.data as string)
            return null
        }

        return json.data
    } catch (error: any) {
        console.error(error)
        toast.error(error.message || 'Network error')
        return null
    }
}

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
