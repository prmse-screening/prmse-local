import type {
    CreateTaskRequest,
    CreateTaskResponse,
    DeleteTaskRequest,
    GetTaskLists,
    GetTaskResponse,
    GetUploadUrlResponse,
    UpdateTaskRequest,
} from '@/types'
import { http } from '@/apis/common.ts'

export const getTask = (id: string) => {
    return http<GetTaskResponse>(`/tasks/${id}`)
}

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

export const getTaskList = ({
    page,
    pageSize,
    status,
    series,
    sortKey,
    sortOrder,
}: {
    page: number
    pageSize: number
    status?: number
    series?: string
    sortKey?: string
    sortOrder?: string
}) => {
    const params: Record<string, string> = {
        page: page.toString(),
        pageSize: pageSize.toString(),
    }

    if (status != null) params.status = status.toString()
    if (series && series != '') params.series = series.trim()
    if (sortKey) params.sortKey = sortKey.trim()
    if (sortOrder) params.sortOrder = sortOrder.trim()

    return http<GetTaskLists>(`/tasks/list?${new URLSearchParams(params)}`)
}

export const setWorkerDevice = (params: { device: string }) => {
    const query = new URLSearchParams(params).toString()
    return http<string>(`/tasks/device?${query}`)
}