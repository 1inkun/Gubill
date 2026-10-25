import type { JwtPayload } from '@/types/api'

const B64_CHARS = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/'
const B64_LOOKUP: Record<string, number> = {}
for (let i = 0; i < B64_CHARS.length; i++) B64_LOOKUP[B64_CHARS[i]] = i

/** base64url -> 字节数组（小程序里没有 atob/Buffer，必须自己解） */
export function base64UrlToBytes(input: string): number[] {
  const clean = input.replace(/-/g, '+').replace(/_/g, '/').replace(/=+$/, '')
  const bytes: number[] = []
  let buffer = 0
  let bits = 0
  for (let i = 0; i < clean.length; i++) {
    const v = B64_LOOKUP[clean[i]]
    if (v === undefined) throw new Error(`base64url 非法字符: ${clean[i]}`)
    buffer = (buffer << 6) | v
    bits += 6
    if (bits >= 8) {
      bits -= 8
      bytes.push((buffer >> bits) & 0xff)
    }
  }
  return bytes
}

/** 字节数组 -> 字符串（含 UTF-8 多字节） */
export function utf8Decode(bytes: number[]): string {
  let out = ''
  let i = 0
  while (i < bytes.length) {
    const b = bytes[i++]
    if (b < 0x80) {
      out += String.fromCharCode(b)
    } else if (b < 0xe0) {
      out += String.fromCharCode(((b & 0x1f) << 6) | (bytes[i++] & 0x3f))
    } else if (b < 0xf0) {
      out += String.fromCharCode(((b & 0x0f) << 12) | ((bytes[i++] & 0x3f) << 6) | (bytes[i++] & 0x3f))
    } else {
      const cp =
        ((b & 0x07) << 18) | ((bytes[i++] & 0x3f) << 12) | ((bytes[i++] & 0x3f) << 6) | (bytes[i++] & 0x3f)
      const off = cp - 0x10000
      out += String.fromCharCode(0xd800 + (off >> 10), 0xdc00 + (off & 0x3ff))
    }
  }
  return out
}

export function decodeJwtPayload(token: string): JwtPayload | null {
  try {
    const parts = token.split('.')
    if (parts.length !== 3) return null
    const payload = JSON.parse(utf8Decode(base64UrlToBytes(parts[1]))) as JwtPayload
    if (!payload || typeof payload !== 'object') return null
    return payload
  } catch {
    return null
  }
}

export function isTokenExpired(payload: JwtPayload | null, nowSeconds = Math.floor(Date.now() / 1000)): boolean {
  if (!payload || typeof payload.exp !== 'number') return true
  return payload.exp <= nowSeconds
}
