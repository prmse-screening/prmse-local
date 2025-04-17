import type {
    CreateTaskRequest,
    CreateTaskResponse,
    DeleteTaskRequest,
    GetTaskLists,
    GetUploadUrlResponse,
    UpdateTaskRequest,
} from '@/types'
import { http } from '@/apis/common.ts'

export const createTask = (data: CreateTaskRequest) => {
    return http<CreateTaskResponse>('/tasks/create', {
        method: 'POST',
        body: JSON.stringify(data),
        headers: { 'Content-Type': 'application/json' },
    })
}

export const updateTask = (data: UpdateTaskRequest) => {
    return http<string>('/tasks/update', {
        method: 'POST',
        body: JSON.stringify(data),
        headers: { 'Content-Type': 'application/json' },
    })
}

export const prioritizeTask = (data: UpdateTaskRequest) => {
    return http<string>('/tasks/prioritize', {
        method: 'POST',
        body: JSON.stringify(data),
        headers: { 'Content-Type': 'application/json' },
    })
}

export const deleteTask = (data: DeleteTaskRequest) => {
    return http<string>('/tasks/delete', {
        method: 'POST',
        body: JSON.stringify(data),
        headers: { 'Content-Type': 'application/json' },
    })
}

export const getUploadUrl = (params: { series: string }) => {
    const query = new URLSearchParams(params).toString()
    return http<GetUploadUrlResponse>(`/tasks/upload?${query}`)
}

export const getTaskList = (params: { page: number; pageSize: number; status?: number; series?: string }) => {
    const query = new URLSearchParams(params as any).toString()
    return http<GetTaskLists>(`/tasks/list?${query}`)
}

export const setWorkerDevice = (params: { device: string }) => {
    const query = new URLSearchParams(params).toString()
    return http<string>(`/tasks/device?${query}`)
}
