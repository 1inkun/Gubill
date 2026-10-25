import { request } from './client'
import { PAGE_SIZE } from '@/config'
import type { FinishOrderResponse, OrderResponse, PageResult } from '@/types/api'

/**
 * 后端 GET /order 的 status 语义（实测确认）：
 *   0 = 只查进行中；1 = 只查已结算；
 *   -2（任何 < -1 的值）= 不过滤，即「全部订单」；
 *   不传则按 Go 零值 0 处理。
 */
export const ORDER_STATUS = {
  RUNNING: 0,
  FINISHED: 1,
  ALL: -2,
} as const

export function createOrder(userId: number): Promise<OrderResponse> {
  return request<OrderResponse>({ url: '/order', method: 'POST', data: { user_id: userId } })
}

/** 按 id 查询返回的是单个对象，不是分页结构 */
export function getOrderById(orderId: number): Promise<OrderResponse> {
  return request<OrderResponse>({ url: `/order?id=${orderId}` })
}

/** 按用户查询返回 { data: Order[], page, pageSize }，注意是两层 data */
export function listOrders(params: { userId: number; status: number; page: number }): Promise<PageResult<OrderResponse>> {
  return request<PageResult<OrderResponse>>({
    url: `/order?user_id=${params.userId}&status=${params.status}&page=${params.page}&page_size=${PAGE_SIZE}`,
  })
}

/** 结算签到订单：后端会同时创建一条待支付订单，返回其 pay_order_id */
export function finishOrder(orderId: number): Promise<FinishOrderResponse | null> {
  return request<FinishOrderResponse | null>({ url: `/order/${orderId}`, method: 'PUT' })
}
