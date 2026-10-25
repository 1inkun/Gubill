import { decodeJwtPayload, isTokenExpired } from '@/utils/jwt'
import type { JwtPayload } from '@/types/api'

const TOKEN_KEY = 'token'

export function getToken(): string {
  return (uni.getStorageSync(TOKEN_KEY) as string) || ''
}

export function setToken(token: string): void {
  uni.setStorageSync(TOKEN_KEY, token)
}

export function clearSession(): void {
  uni.removeStorageSync(TOKEN_KEY)
}

export function getPayload(): JwtPayload | null {
  return decodeJwtPayload(getToken())
}

export function isLoggedIn(): boolean {
  const payload = getPayload()
  return !!payload && !isTokenExpired(payload)
}

/** 后端 /login 不返回 userId，从 JWT payload 里取；取不到返回 0 */
export function getUserId(): number {
  const payload = getPayload()
  return payload && typeof payload.user_id === 'number' ? payload.user_id : 0
}

export function getUsername(): string {
  return getPayload()?.username ?? ''
}

export function logout(): void {
  clearSession()
}
