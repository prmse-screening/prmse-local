import { TaskStatus } from '@/types/enums'
import type {Result, Task} from '@/types'

// Request
export type CreateTaskRequest = {
    series: string
}

export type UpdateTaskRequest = {
    id: number
    series: string
    status: TaskStatus
    result: Result
    model: string
    order: number
    updated: string
}

export type DeleteTaskRequest = UpdateTaskRequest

// Response
export type GetTaskResponse = Task

export type CreateTaskResponse = Task

export type GetTaskLists = {
    tasks: Task[]
    total: number
}

export type S3UploadForm = Record<string, string>

export type GetUploadUrlResponse = {
    url: string
    form: S3UploadForm
}
