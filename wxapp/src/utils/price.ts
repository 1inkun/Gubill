import { SINGLE_PRICE_CENTS } from '@/config'

/**
 * 与后端 internal/utils/calcPrice.go 逐字对齐的预估算法：
 *   seconds   = ceil(now - start_at)
 *   halfHours = ceil(seconds / 3600) * 2
 *   price(分) = halfHours * 500
 * 禁止使用 billAppFront 里那套 12/24 小时封顶算法（与后端不一致）。
 */
export function estimatePriceCents(startAtMs: number, nowMs: number): number {
  const seconds = Math.ceil(Math.max(0, nowMs - startAtMs) / 1000)
  const halfHours = Math.ceil(seconds / 3600) * 2
  return halfHours * SINGLE_PRICE_CENTS
}

/** 分 -> "10.00" */
export function formatCents(cents: number): string {
  const n = Number.isFinite(cents) ? Math.trunc(cents) : 0
  const abs = Math.abs(n)
  return `${n < 0 ? '-' : ''}${Math.trunc(abs / 100)}.${String(abs % 100).padStart(2, '0')}`
}
