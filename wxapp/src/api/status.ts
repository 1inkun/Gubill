import { request } from './client'
import type { CurrentUser } from '@/types/api'

/** 当前签到中（未结算订单）的用户名单，无需鉴权 */
export function getCheckinStatus(): Promise<CurrentUser[]> {
  return request<CurrentUser[]>({ url: '/status', auth: false })
}
