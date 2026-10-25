import { describe, expect, it } from 'vitest'
import { API_PATH, PAGE_SIZE, SINGLE_PRICE_CENTS, resolveApiBase } from '@/config'

describe('config', () => {
  it('H5 用相对路径（走 vite 代理）', () => {
    expect(resolveApiBase(true)).toBe('/api/v1')
  })
  it('小程序/无 DOM 环境用绝对地址', () => {
    expect(resolveApiBase(false)).toBe('http://localhost:8080/api/v1')
  })
  it('常量与后端一致', () => {
    expect(API_PATH).toBe('/api/v1')
    expect(SINGLE_PRICE_CENTS).toBe(500)
    expect(PAGE_SIZE).toBe(10)
  })
})
