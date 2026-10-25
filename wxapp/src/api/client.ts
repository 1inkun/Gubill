import { API_BASE, REQUEST_TIMEOUT } from '@/config'

export interface RequestOptions {
  url: string
  method?: 'GET' | 'POST' | 'PUT'
  data?: Record<string, unknown>
  /** 默认 true；登录/注册/在场名单传 false */
  auth?: boolean
}

export class RequestError extends Error {
  statusCode: number
  errCode: number
  constructor(statusCode: number, message: string, errCode = 0) {
    super(message)
    this.name = 'RequestError'
    this.statusCode = statusCode
    this.errCode = errCode
  }
}

let tokenGetter: () => string = () => ''
let unauthorizedHandler: () => void = () => {}

export function configureClient(opts: { getToken: () => string; onUnauthorized?: () => void }): void {
  tokenGetter = opts.getToken
  if (opts.onUnauthorized) unauthorizedHandler = opts.onUnauthorized
}

/**
 * 统一出口。约定：
 * - 成功只看 HTTP 2xx（后端成功响应里 code 字段的 json tag 是错的，实际输出为 Code，不可依赖）；
 * - 失败抛 RequestError，message 直接取后端 msg，可直接 toast；
 * - 401 统一回调登出。
 */
export function request<T>(options: RequestOptions): Promise<T> {
  const withAuth = options.auth !== false
  const token = withAuth ? tokenGetter() : ''
  return new Promise<T>((resolve, reject) => {
    uni.request({
      url: API_BASE + options.url,
      method: (options.method ?? 'GET') as 'GET' | 'POST' | 'PUT',
      data: options.data,
      timeout: REQUEST_TIMEOUT,
      header: {
        'content-type': 'application/json',
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
      },
      success: (res) => {
        const body = (res.data ?? {}) as { msg?: string; errCode?: number; data?: unknown }
        if (res.statusCode >= 200 && res.statusCode < 300) {
          resolve((body.data ?? null) as T)
          return
        }
        if (res.statusCode === 401) {
          // 区分两类 401：
          // - Auth 中间件返回的 {code:0,msg:"认证失败"} 没有 errCode → 会话真失效，需要登出；
          // - 业务 handler 返回的 401 带 errCode（如 101002/102002）→ 业务错误，仅提示，不登出。
          const isSessionInvalid = body.errCode === undefined || body.errCode === null
          if (isSessionInvalid) {
            try {
              unauthorizedHandler()
            } catch {
              /* 忽略登出回调自身异常 */
            }
          }
        }
        reject(new RequestError(res.statusCode, body.msg || `请求失败（${res.statusCode}）`, body.errCode ?? 0))
      },
      fail: () => {
        reject(new RequestError(0, `网络异常，请确认后端已启动（${API_BASE}）`))
      },
    })
  })
}
