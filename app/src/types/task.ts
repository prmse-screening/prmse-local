import { TaskStatus } from '@/types/enums'

export type Result = { prediction: number[]; threshold: number }

export type Task = {
    id: number
    series: string
    status: TaskStatus
    result: Result
    model: string
    order: number
    updated: string
}
