import { beforeEach, describe, expect, it, vi } from 'vitest'

const uniMock = {
  getStorageSync: vi.fn(() => 'tok-1'),
  setStorageSync: vi.fn(),
  removeStorageSync: vi.fn(),
  reLaunch: vi.fn(),
  showToast: vi.fn(),
  request: vi.fn(),
}
vi.stubGlobal('uni', uniMock)

import { RequestError, configureClient, request } from '@/api/client'

describe('request', () => {
  beforeEach(() => {
    uniMock.request.mockReset()
    configureClient({ getToken: () => 'tok-1', onUnauthorized: vi.fn() })
  })

  it('HTTP 200 时直接返回 data', async () => {
    uniMock.request.mockImplementation((opts: any) => {
      opts.success({ statusCode: 200, data: { Code: 0, status: 'success', msg: '登录成功', data: { token: 't' } } })
    })
    await expect(request<{ token: string }>({ url: '/login', method: 'POST', data: { username: 'a' } })).resolves.toEqual({
      token: 't',
    })
  })

  it('自动带上 Authorization: Bearer <token>，并拼上 API_BASE', async () => {
    uniMock.request.mockImplementation((opts: any) => {
      opts.success({ statusCode: 200, data: { status: 'success', data: null } })
    })
    await request({ url: '/order?page=1' })
    const opts = uniMock.request.mock.calls[0][0]
    expect(opts.header.Authorization).toBe('Bearer tok-1')
    expect(opts.url.endsWith('/api/v1/order?page=1')).toBe(true)
  })

  it('auth:false 时不带 Authorization', async () => {
    uniMock.request.mockImplementation((opts: any) => {
      opts.success({ statusCode: 200, data: { status: 'success', data: [] } })
    })
    await request({ url: '/status', auth: false })
    expect(uniMock.request.mock.calls[0][0].header.Authorization).toBeUndefined()
  })

  it('业务失败时抛出 RequestError，message 取后端 msg', async () => {
    uniMock.request.mockImplementation((opts: any) => {
      opts.success({ statusCode: 400, data: { errCode: 102001, status: 'fail', msg: '签到失败', detail: '存在尚未完成的订单' } })
    })
    await expect(request({ url: '/order', method: 'POST' })).rejects.toMatchObject({
      statusCode: 400,
      errCode: 102001,
      message: '签到失败',
    })
  })

  it('401 无 errCode（Auth 中间件的认证失败）触发 onUnauthorized', async () => {
    const onUnauthorized = vi.fn()
    configureClient({ getToken: () => 'tok-1', onUnauthorized })
    uniMock.request.mockImplementation((opts: any) => {
      opts.success({ statusCode: 401, data: { code: 0, msg: '认证失败', data: null } })
    })
    await expect(request({ url: '/order?user_id=1' })).rejects.toBeInstanceOf(RequestError)
    expect(onUnauthorized).toHaveBeenCalledTimes(1)
  })

  it('401 带 errCode（业务错误，如 101002）不触发 onUnauthorized', async () => {
    const onUnauthorized = vi.fn()
    configureClient({ getToken: () => 'tok-1', onUnauthorized })
    uniMock.request.mockImplementation((opts: any) => {
      opts.success({ statusCode: 401, data: { errCode: 101002, status: 'fail', msg: '登录状态有误' } })
    })
    await expect(request({ url: '/good/1', method: 'POST' })).rejects.toMatchObject({
      statusCode: 401,
      errCode: 101002,
      message: '登录状态有误',
    })
    expect(onUnauthorized).not.toHaveBeenCalled()
  })

  it('网络失败时给出可定位的提示', async () => {
    uniMock.request.mockImplementation((opts: any) => {
      opts.fail({ errMsg: 'request:fail' })
    })
    await expect(request({ url: '/status', auth: false })).rejects.toMatchObject({ statusCode: 0 })
  })
})
