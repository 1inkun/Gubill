<template>
  <view class="app-card order-item" @click="$emit('select', order.id)">
    <view class="app-row">
      <text class="order-item__id">订单 #{{ order.id }}</text>
      <text class="app-tag" :class="{ 'app-tag--done': order.status !== 0 }">
        {{ order.status === 0 ? '进行中' : '已结算' }}
      </text>
    </view>
    <view class="app-row order-item__line">
      <text class="app-muted">开始</text>
      <text>{{ formatDateTime(parseBackendTime(order.start_at)) }}</text>
    </view>
    <view class="app-row order-item__line">
      <text class="app-muted">时长</text>
      <text>{{ formatDuration(durationSeconds(parseBackendTime(order.start_at), parseBackendTime(order.end_at) || new Date())) }}</text>
    </view>
    <view class="app-row order-item__line">
      <text class="app-muted">金额</text>
      <text class="order-item__price">￥{{ formatCents(order.price) }}</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { durationSeconds, formatDateTime, formatDuration, parseBackendTime } from '@/utils/format'
import { formatCents } from '@/utils/price'
import type { OrderResponse } from '@/types/api'

defineProps<{ order: OrderResponse }>()
// 注意：不要用 tap 作为自定义事件名——uni-app 会把 @tap 当原生 tap 事件（click 的映射）监听，
// 事件对象会覆盖 emit 的 payload（实测踩坑：URL 里出现 id=[object Object]）
defineEmits<{ (e: 'select', id: number): void }>()
</script>

<style scoped>
.order-item {
  margin-bottom: 24rpx;
}
.order-item__id {
  font-size: 30rpx;
  font-weight: 600;
}
.order-item__line {
  margin-top: 12rpx;
  font-size: 26rpx;
}
.order-item__price {
  color: var(--c-primary-strong);
  font-weight: 600;
}
</style>
