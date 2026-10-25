import { describe, expect, it } from 'vitest'
import { estimatePriceCents, formatCents } from '@/utils/price'

const HOUR = 3600 * 1000

describe('estimatePriceCents（必须与后端 calcPrice.go 一致）', () => {
  it('刚签到 0 秒为 0 分', () => {
    expect(estimatePriceCents(0, 0)).toBe(0)
  })
  it('10 分钟按 1 小时 2 档计费 = 1000 分（10.00 元）', () => {
    expect(estimatePriceCents(0, 10 * 60 * 1000)).toBe(1000)
  })
  it('正好 1 小时 = 1000 分', () => {
    expect(estimatePriceCents(0, HOUR)).toBe(1000)
  })
  it('1 小时零 1 秒进到第 3 档 = 2000 分', () => {
    expect(estimatePriceCents(0, HOUR + 1000)).toBe(2000)
  })
  it('2 小时 = 2000 分', () => {
    expect(estimatePriceCents(0, 2 * HOUR)).toBe(2000)
  })
})

describe('formatCents', () => {
  it('1000 -> 10.00', () => expect(formatCents(1000)).toBe('10.00'))
  it('0 -> 0.00', () => expect(formatCents(0)).toBe('0.00'))
  it('505 -> 5.05', () => expect(formatCents(505)).toBe('5.05'))
})
