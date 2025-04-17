export type BaseResponse<T = any> = {
    code: number
    status: string
    data: T
}
