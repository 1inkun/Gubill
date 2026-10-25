/** 后端真实路由组是 /api/v1；swagger.json 里的 /api/v1vi 是 router.go 注释笔误 */
export const BACKEND_ORIGIN = 'http://localhost:8080'
export const API_PATH = '/api/v1'

/**
 * H5 走 vite 代理（同源，绕开 CORS）；小程序端没有代理，必须绝对地址。
 * 用运行时探测而不是条件编译：这样 vitest（node 环境）也能直接跑这个文件。
 */
export function resolveApiBase(hasDOM: boolean): string {
  return hasDOM ? API_PATH : BACKEND_ORIGIN + API_PATH
}

export const API_BASE = resolveApiBase(typeof window !== 'undefined' && typeof document !== 'undefined')

/** 与后端 os.Setenv("price","500") 一致：500 分（5 元）/半小时 */
export const SINGLE_PRICE_CENTS = 500
/** 后端 SQL 里 LIMIT 写死 10，page_size 无效 */
export const PAGE_SIZE = 10
export const REQUEST_TIMEOUT = 10000
export const APP_VERSION = '1.0.0'
