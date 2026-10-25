export interface OrderResponse {
  id: number
  user_id: number
  /** Go time.Time，形如 "2026-10-12T18:49:25Z"（真实 UTC）；未结算时为 "0001-01-01T00:00:00Z" */
  start_at: string
  end_at: string
  /** 0 = 进行中（未结算）, 1 = 已结算 */
  status: number
  /** 单位：分 */
  price: number
}

export interface CurrentUser {
  username: string
  nickname: string
  avatar_file: string
}

export interface LoginResult {
  token: string
}

export interface PageResult<T> {
  data: T[]
  page: number
  pageSize: number
}

export interface JwtPayload {
  user_id: number
  username: string
  role: string
  iat?: number
  exp?: number
}

/** GET /good：在售会员商品（price 单位：分，展示时转元） */
export interface Good {
  id: number
  name: string
  type: string
  description: string
  status: number
  price: number
}

/** POST /good/:good_id：购买商品后创建的支付订单 */
export interface BuyGoodResponse {
  pay_order_id: number
}

/** PUT /order/:order_id：结算签到订单后创建的支付订单 */
export interface FinishOrderResponse {
  pay_order_id: number
}

/** GET /member：用户会员信息（status: 0=生效中, 其他=已过期；无记录时后端返回 500） */
export interface MemberData {
  id: number
  user_id: number
  status: number
  start_at: string
  end_at: string
}

/** GET /pay：支付订单（status: 0=待支付/未结算） */
export interface PayOrderData {
  id: number
  user_id: number
  type: string
  order_no: number
  transaction_id: string
  privider: string
  client_ip: string
  status: number
  total_fee: number
  description: string
}
