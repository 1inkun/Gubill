import { request } from './client'
import { PAGE_SIZE } from '@/config'
import type { PageResult, PayOrderData } from '@/types/api'

/**
 * 支付订单状态：0=待支付（未结算），与签到订单的 status 含义不同。
 * 按 user_id 查询时返回分页结构 { data: PayOrderData[], page, pageSize }。
 */
export const PAY_ORDER_STATUS = {
  PENDING: 0,
  PAID: 1,
  ALL: -2,
} as const

export function listPayOrders(params: { userId: number; status: number; page: number }): Promise<PageResult<PayOrderData>> {
  return request<PageResult<PayOrderData>>({
    url: `/pay?user_id=${params.userId}&status=${params.status}&page=${params.page}&page_size=${PAGE_SIZE}`,
  })
}

export function getPayOrderById(id: number): Promise<PayOrderData> {
  return request<PayOrderData>({ url: `/pay?id=${id}` })
}

export type PayChannel = 'alipay' | 'wechat'

/**
 * 发起支付请求。
 * 注意：后端支付接口尚未完成（change.md：支付请求暂未完成），
 * 这里约定端点为 POST /pay/:pay_order_id、body 携带支付方式，
 * 后端补齐后若路径/参数有变，只需调整本函数。
 */
export function payOrder(payOrderId: number, channel: PayChannel): Promise<{ msg: string } | null> {
  return request<{ msg: string } | null>({
    url: `/pay/${payOrderId}`,
    method: 'POST',
    data: { type: channel },
  })
}
