import { TaskStatus } from '@/types/enums'

export type Task = {
    id: number
    series: string
    status: TaskStatus
    result: string
    model: string
    order: number
    updated: string
}
