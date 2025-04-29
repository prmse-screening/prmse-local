import { TaskStatus } from '@/types/enums'

export type Task = {
    id: number
    series: string
    status: TaskStatus
    result: { prediction: number[]; threshold: number }
    model: string
    order: number
    updated: string
}
