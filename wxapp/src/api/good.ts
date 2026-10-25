import { request } from './client'
import type { BuyGoodResponse, Good } from '@/types/api'

/** 在售会员商品列表（公开接口，无需鉴权） */
export function getGoods(): Promise<Good[]> {
  return request<Good[]>({ url: '/good', auth: false })
}

/**
 * 购买会员商品：后端创建一条 status=0（待支付）的支付订单。
 * 已知后端问题：/good 路由组未挂 Auth 中间件，导致 handler 取不到 userId
 * 而恒定返回 101002「登录状态有误」，待后端修复（见 README）。
 */
export function buyGood(goodId: number): Promise<BuyGoodResponse> {
  return request<BuyGoodResponse>({ url: `/good/${goodId}`, method: 'POST' })
}
