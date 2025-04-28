export enum TaskStatus {
    Preparing = 0,
    Pending = 1,
    Running = 2,
    Success = 3,
    Failed = 4,
}

export const STATUS_COLOR: Record<TaskStatus, string> = {
    [TaskStatus.Preparing]: '',
    [TaskStatus.Pending]: 'secondary',
    [TaskStatus.Running]: 'primary',
    [TaskStatus.Success]: 'success',
    [TaskStatus.Failed]: 'error',
}

export const STATUS_NAME: Record<TaskStatus, string> = {
    [TaskStatus.Preparing]: 'Preparing',
    [TaskStatus.Pending]: 'Pending',
    [TaskStatus.Running]: 'Running',
    [TaskStatus.Success]: 'Success',
    [TaskStatus.Failed]: 'Failed',
}

export const STATUS_LIST = [
    { label: 'Preparing', value: TaskStatus.Preparing },
    { label: 'Pending', value: TaskStatus.Pending },
    { label: 'Running', value: TaskStatus.Running },
    { label: 'Success', value: TaskStatus.Success },
    { label: 'Failed', value: TaskStatus.Failed },
]
