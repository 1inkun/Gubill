import { configureClient } from './client'
import { clearSession, getToken } from '@/store/session'

let installed = false

/** 在 App.vue onLaunch 里调用一次：把会话接到请求层，并统一处理 401 */
export function setupApiClient(): void {
  if (installed) return
  installed = true
  configureClient({
    getToken,
    onUnauthorized: () => {
      clearSession()
      const pages = (globalThis as any).getCurrentPages?.() ?? []
      const current: string = pages.length ? pages[pages.length - 1].route : ''
      if (current.indexOf('pages/login/login') !== -1) return
      uni.reLaunch({ url: '/pages/login/login' })
      uni.showToast({ title: '登录已过期，请重新登录', icon: 'none' })
    },
  })
}
