<template>
  <AppShell>
    <AppCard title="接口">
      <view class="setting__row"><text class="app-muted">后端地址</text><text>{{ backendOrigin }}</text></view>
      <view class="setting__row"><text class="app-muted">接口前缀</text><text>{{ apiPath }}</text></view>
      <view class="setting__row"><text class="app-muted">当前平台基址</text><text>{{ apiBase }}</text></view>
    </AppCard>

    <view class="app-gap" />
    <AppCard title="主题">
      <view class="setting__row"><text class="app-muted">黑夜模式</text><text>跟随系统</text></view>
      <text class="app-muted setting__note">修改系统外观后刷新页面即可生效</text>
    </AppCard>

    <view class="app-gap" />
    <AppCard title="本地数据">
      <view class="setting__row" @click="onClear"><text>清除本地缓存（含登录状态）</text><text class="app-muted">›</text></view>
    </AppCard>

    <view class="app-gap" />
    <AppCard title="关于">
      <view class="setting__row"><text class="app-muted">前端版本</text><text>{{ appVersion }}</text></view>
      <view class="setting__row"><text class="app-muted">计费规则</text><text>5 元 / 半小时</text></view>
    </AppCard>

    <view class="app-gap" />
    <AppButton type="outline" @click="onLogout">退出登录</AppButton>
  </AppShell>
</template>

<script setup lang="ts">
import AppShell from '@/components/AppShell.vue'
import AppCard from '@/components/AppCard.vue'
import AppButton from '@/components/AppButton.vue'
import { API_BASE, API_PATH, APP_VERSION, BACKEND_ORIGIN } from '@/config'
import { clearSession } from '@/store/session'
import { toast } from '@/utils/toast'

const backendOrigin = BACKEND_ORIGIN
const apiPath = API_PATH
const apiBase = API_BASE
const appVersion = APP_VERSION

function onClear() {
  uni.showModal({
    title: '清除缓存',
    content: '将清除本地登录状态，需要重新登录',
    success: (res) => {
      if (!res.confirm) return
      clearSession()
      toast('已清除')
      uni.reLaunch({ url: '/pages/login/login' })
    },
  })
}

function onLogout() {
  clearSession()
  uni.reLaunch({ url: '/pages/login/login' })
}
</script>

<style scoped>
.setting__row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20rpx 0;
  border-bottom: 1rpx solid var(--c-border);
  font-size: 28rpx;
}
.setting__note {
  font-size: 24rpx;
}
</style>
