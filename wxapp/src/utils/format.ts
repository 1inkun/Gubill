/**
 * 后端时间戳解析：实测（scripts/e2e.sh）后端返回的是真实 UTC 时刻
 * （"2026-10-12T18:49:25Z" 与 date -u 吻合，MySQL 会话时区为 UTC），
 * 所以按标准 RFC3339 解析即可；未结算订单的 end_at 是零值
 * "0001-01-01T00:00:00Z"，视为无值返回 null。
 */
export function parseBackendTime(value?: string | null): Date | null {
  if (!value) return null
  const year = Number(value.slice(0, 4))
  if (!Number.isFinite(year) || year < 1970) return null
  const date = new Date(value)
  return Number.isNaN(date.getTime()) ? null : date
}

const pad = (n: number) => String(n).padStart(2, '0')

export function formatDateTime(date: Date | null): string {
  if (!date) return '--'
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(
    date.getMinutes(),
  )}:${pad(date.getSeconds())}`
}

export function formatDuration(totalSeconds: number): string {
  const s = Math.max(0, Math.floor(totalSeconds))
  return `${pad(Math.floor(s / 3600))}:${pad(Math.floor((s % 3600) / 60))}:${pad(s % 60)}`
}

export function durationSeconds(start: Date | null, end: Date | null): number {
  if (!start || !end) return 0
  return Math.max(0, Math.floor((end.getTime() - start.getTime()) / 1000))
}
