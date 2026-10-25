import { describe, expect, it } from 'vitest'
import { decodeJwtPayload, isTokenExpired } from '@/utils/jwt'

/**
 * 夹具用 node 自带的 Buffer(base64url) 构造——它是与被测解码器相互独立的实现，
 * 不构成循环验证。（终端里 JWT 形状的字符串会被脱敏，所以不把夹具硬编码在测试外。）
 */
const makeToken = (payload: object, pad = false) => {
  const b64 = Buffer.from(JSON.stringify(payload)).toString('base64url')
  return `a.${pad ? b64 + '=' : b64}.b`
}

const TOKEN_A = makeToken({ user_id: 100001, username: 'admin', role: 'admin', iat: 1760000000, exp: 1760021600 })
const TOKEN_B = makeToken({ user_id: 7, username: '中文用户', role: 'user', exp: 1760021600 })
// 0xfb,0xff,0xbf 的 base64url 会自然出现 -/_ 字符
const TOKEN_C = makeToken({ user_id: 3, username: Buffer.from([0xfb, 0xff, 0xbf]).toString('hex') })

describe('decodeJwtPayload', () => {
  it('解出 user_id / username / role', () => {
    const p = decodeJwtPayload(TOKEN_A)
    expect(p).not.toBeNull()
    expect(p!.user_id).toBe(100001)
    expect(p!.username).toBe('admin')
    expect(p!.role).toBe('admin')
  })
  it('能处理含中文的 payload（UTF-8 正确）', () => {
    expect(decodeJwtPayload(TOKEN_B)!.username).toBe('中文用户')
  })
  it('能处理带 = padding 与 URL 安全字符 -/_ 的 base64url', () => {
    expect(decodeJwtPayload(TOKEN_A + '==')!.user_id).toBe(100001)
    expect(decodeJwtPayload(TOKEN_C)!.username).toBe(Buffer.from([0xfb, 0xff, 0xbf]).toString('hex'))
  })
  it('非法 token 返回 null', () => {
    expect(decodeJwtPayload('not-a-jwt')).toBeNull()
    expect(decodeJwtPayload('')).toBeNull()
    expect(decodeJwtPayload('a.@@@@.b')).toBeNull()
  })
})

describe('isTokenExpired', () => {
  it('exp 小于当前时间视为过期', () => {
    expect(isTokenExpired({ user_id: 1, username: 'u', role: 'user', exp: 100 }, 200)).toBe(true)
    expect(isTokenExpired({ user_id: 1, username: 'u', role: 'user', exp: 300 }, 200)).toBe(false)
  })
  it('没有 exp 视为过期', () => {
    expect(isTokenExpired(null)).toBe(true)
    expect(isTokenExpired({ user_id: 1, username: 'u', role: 'user' })).toBe(true)
  })
})
