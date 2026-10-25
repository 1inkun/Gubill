import { isLoggedIn } from './session'

/** 返回 true 表示已登录；未登录时跳登录页 */
export function ensureLogin(redirect = true): boolean {
  if (isLoggedIn()) return true
  if (redirect) uni.reLaunch({ url: '/pages/login/login' })
  return false
}
