import { describe, expect, it } from 'vitest'
import { durationSeconds, formatDateTime, formatDuration, parseBackendTime } from '@/utils/format'

describe('parseBackendTime（按真实 UTC 时刻解析，见 scripts/e2e.sh 实测）', () => {
  it('空值与零值时间返回 null', () => {
    expect(parseBackendTime(undefined)).toBeNull()
    expect(parseBackendTime(null)).toBeNull()
    expect(parseBackendTime('0001-01-01T00:00:00Z')).toBeNull()
  })
  it('解析 Z 结尾的真实 UTC 时刻（与 Date 原生解析一致，跨时区确定）', () => {
    const d = parseBackendTime('2026-10-12T18:49:25Z')!
    expect(d.getTime()).toBe(new Date('2026-10-12T18:49:25Z').getTime())
  })
  it('带 +08:00 偏移的字符串同样按标准解析', () => {
    const d = parseBackendTime('2026-10-12T18:49:25+08:00')!
    expect(d.getTime()).toBe(new Date('2026-10-12T18:49:25+08:00').getTime())
  })
})

describe('formatDateTime / formatDuration / durationSeconds', () => {
  it('格式化', () => {
    expect(formatDateTime(new Date(2026, 9, 13, 2, 39, 5))).toBe('2026-10-13 02:39:05')
    expect(formatDateTime(null)).toBe('--')
  })
  it('时长', () => {
    expect(formatDuration(0)).toBe('00:00:00')
    expect(formatDuration(3725)).toBe('01:02:05')
    expect(durationSeconds(new Date(0), new Date(3725 * 1000))).toBe(3725)
    expect(durationSeconds(null, new Date())).toBe(0)
  })
})
