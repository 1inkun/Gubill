export function toast(title: string, icon: 'none' | 'success' = 'none'): void {
  uni.showToast({ title, icon, duration: 2000 })
}

export function errorMessage(e: unknown): string {
  if (e instanceof Error && e.message) return e.message
  if (typeof e === 'string' && e) return e
  return '操作失败，请稍后重试'
}
