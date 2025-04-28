import type {
    CreateTaskRequest,
    CreateTaskResponse,
    DeleteTaskRequest,
    GetTaskLists,
    GetTaskResponse,
    GetUploadUrlResponse,
    UpdateTaskRequest,
} from '@/types'
import { GET, POST } from '@/apis/common.ts'

export const getTask = (id: string) => {
    return GET<GetTaskResponse>(`/tasks/${id}`)
}

export const createTask = (data: CreateTaskRequest) => {
    return POST<CreateTaskResponse>('/tasks/create', JSON.stringify(data))
}

export const updateTask = (data: UpdateTaskRequest) => {
    return POST<string>('/tasks/update', JSON.stringify(data))
}

export const prioritizeTask = (data: UpdateTaskRequest) => {
    return POST<string>('/tasks/prioritize', JSON.stringify(data))
}

export const deleteTask = (data: DeleteTaskRequest) => {
    return POST<string>('/tasks/delete', JSON.stringify(data))
}

export const getUploadPostUrl = (params: { series: string }) => {
    const query = new URLSearchParams(params).toString()
    return GET<GetUploadUrlResponse>(`/tasks/uploadPost?${query}`)
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

    return GET<GetTaskLists>(`/tasks/list?${new URLSearchParams(params)}`)
}

export const setWorkerDevice = (params: { device: string }) => {
    const query = new URLSearchParams(params).toString()
    return GET<string>(`/tasks/device?${query}`)
}
