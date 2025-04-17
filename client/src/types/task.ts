import { TaskState } from '@/types/enums'

export type Task = {
    id: number
    series: string
    status: TaskState
    result: string
    model: string
    order: number
    updated: string
}
