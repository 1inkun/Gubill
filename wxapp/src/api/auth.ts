import { request } from './client'
import type { LoginResult } from '@/types/api'

export function login(body: { username: string; password: string }): Promise<LoginResult> {
  return request<LoginResult>({ url: '/login', method: 'POST', data: body, auth: false })
}

export function register(body: {
  username: string
  nickname: string
  password: string
  email: string
}): Promise<null> {
  return request<null>({ url: '/register', method: 'POST', data: body, auth: false })
}
