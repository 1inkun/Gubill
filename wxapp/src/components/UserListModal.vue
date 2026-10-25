<template>
  <view v-if="visible" class="ulist__mask" @click="$emit('close')">
    <view class="ulist__card app-card" @click.stop>
      <text class="ulist__title">当前在场用户（{{ users.length }}）</text>
      <scroll-view :scroll-y="true" class="ulist__list">
        <view v-if="!users.length" class="ulist__empty app-muted">当前没有人在场</view>
        <view v-for="u in users" :key="u.username" class="ulist__item">
          <view class="ulist__avatar">{{ (u.nickname || u.username || '客').slice(0, 1) }}</view>
          <view class="ulist__meta">
            <text class="ulist__name">{{ u.nickname || u.username }}</text>
            <text class="app-muted ulist__sub">账号：{{ u.username }}</text>
          </view>
        </view>
      </scroll-view>
      <AppButton type="outline" @click="$emit('close')">关 闭</AppButton>
    </view>
  </view>
</template>

<script setup lang="ts">
import AppButton from '@/components/AppButton.vue'
import type { CurrentUser } from '@/types/api'

defineProps<{ visible: boolean; users: CurrentUser[] }>()
defineEmits<{ (e: 'close'): void }>()
</script>

<style scoped>
.ulist__mask {
  position: fixed;
  left: 0;
  top: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
}
.ulist__card {
  width: 620rpx;
  max-height: 70vh;
  display: flex;
  flex-direction: column;
}
.ulist__title {
  font-size: 32rpx;
  font-weight: 600;
  margin-bottom: 20rpx;
}
.ulist__list {
  max-height: 46vh;
  margin-bottom: 24rpx;
}
.ulist__empty {
  padding: 40rpx 0;
  text-align: center;
  font-size: 26rpx;
}
.ulist__item {
  display: flex;
  align-items: center;
  padding: 16rpx 0;
  border-bottom: 1rpx solid var(--c-border);
}
.ulist__avatar {
  width: 80rpx;
  height: 80rpx;
  border-radius: 50%;
  background-color: var(--c-primary-weak);
  color: var(--c-primary-strong);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 34rpx;
  margin-right: 20rpx;
}
.ulist__meta {
  display: flex;
  flex-direction: column;
}
.ulist__name {
  font-size: 28rpx;
  font-weight: 600;
}
.ulist__sub {
  font-size: 22rpx;
  margin-top: 4rpx;
}
</style>
