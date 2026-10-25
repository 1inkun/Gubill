import { request } from './client'
import type { MemberData } from '@/types/api'

/**
 * 查询用户会员信息。
 * 注意：后端 models.SearchMemberParams 缺少 form tag，gin 绑定 query 时用
 * 字段名（UserId）而不是 user_id——传小写会得到 102002「非法操作」（实测）。
 */
export function getMemberData(userId: number): Promise<MemberData> {
  return request<MemberData>({ url: `/member?UserId=${userId}` })
}
