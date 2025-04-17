import { TaskState } from '@/types/enums'
import type { Task } from '@/types'

export type CreateTaskRequest = {
    series: string
}

export type UpdateTaskRequest = {
    id: number
    series: string
    status: TaskState
    result: string
    model: string
    order: number
    updated: string
}

export type DeleteTaskRequest = UpdateTaskRequest

export type CreateTaskResponse = Task

export type GetTaskLists = {
    tasks: Task[]
    total: number
}

export type GetUploadUrlResponse = {
    url: string
    form: Record<string, string>
}
