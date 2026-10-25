<template>
  <AppShell>
    <StateView v-if="loading" state="loading" />
    <StateView v-else-if="errorText" state="error" :text="errorText">
      <AppButton type="outline" @click="load">重新加载</AppButton>
    </StateView>
    <template v-else-if="order">
      <AppCard title="订单信息">
        <view class="detail__row"><text class="app-muted">订单号</text><text>#{{ order.id }}</text></view>
        <view class="detail__row">
          <text class="app-muted">状态</text>
          <text class="app-tag" :class="{ 'app-tag--done': order.status !== 0 }">
            {{ order.status === 0 ? '进行中' : '已结算' }}
          </text>
        </view>
        <view class="detail__row">
          <text class="app-muted">开始时间</text>
          <text>{{ formatDateTime(parseBackendTime(order.start_at)) }}</text>
        </view>
        <view class="detail__row">
          <text class="app-muted">结束时间</text>
          <text>{{ order.status === 0 ? '进行中' : formatDateTime(parseBackendTime(order.end_at)) }}</text>
        </view>
        <view class="detail__row"><text class="app-muted">时长</text><text>{{ elapsedText }}</text></view>
        <view class="detail__row">
          <text class="app-muted">金额</text>
          <text class="detail__price">
            ￥{{ order.status === 0 ? formatCents(estimated) + '（预估）' : formatCents(order.price) }}
          </text>
        </view>
      </AppCard>

      <view class="app-gap" />
      <AppButton v-if="order.status === 0" :loading="submitting" @click="onFinish">结算并离场</AppButton>
      <CheckinTimer v-if="order.status === 0" :start-at="order.start_at" />
    </template>
  </AppShell>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import AppShell from '@/components/AppShell.vue'
import AppCard from '@/components/AppCard.vue'
import AppButton from '@/components/AppButton.vue'
import StateView from '@/components/StateView.vue'
import CheckinTimer from '@/components/CheckinTimer.vue'
import { finishOrder, getOrderById } from '@/api/order'
import { ensureLogin } from '@/store/guard'
import { durationSeconds, formatDateTime, formatDuration, parseBackendTime } from '@/utils/format'
import { estimatePriceCents, formatCents } from '@/utils/price'
import { errorMessage, toast } from '@/utils/toast'
import type { OrderResponse } from '@/types/api'

const orderId = ref(0)
const order = ref<OrderResponse | null>(null)
const loading = ref(true)
const submitting = ref(false)
const errorText = ref('')

const elapsedText = computed(() => {
  const o = order.value
  if (!o) return '--'
  const start = parseBackendTime(o.start_at)
  const end = parseBackendTime(o.end_at) || new Date()
  return formatDuration(durationSeconds(start, end))
})

const estimated = computed(() => {
  const o = order.value
  if (!o) return 0
  const start = parseBackendTime(o.start_at)
  return start ? estimatePriceCents(start.getTime(), Date.now()) : 0
})

async function load() {
  if (!ensureLogin()) return
  loading.value = true
  errorText.value = ''
  try {
    order.value = await getOrderById(orderId.value)
  } catch (e) {
    errorText.value = errorMessage(e)
  } finally {
    loading.value = false
  }
}

async function onFinish() {
  const o = order.value
  if (!o) return
  submitting.value = true
  try {
    await finishOrder(o.id)
    await load()
    uni.showModal({ title: '结算成功', content: `订单 #${o.id} 已结算`, showCancel: false })
  } catch (e) {
    toast(errorMessage(e))
  } finally {
    submitting.value = false
  }
}

onLoad((query) => {
  orderId.value = Number((query as Record<string, string>)?.id ?? 0)
  load()
})
</script>

<style scoped>
.detail__row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16rpx 0;
  border-bottom: 1rpx solid var(--c-border);
}
.detail__price {
  color: var(--c-primary-strong);
  font-weight: 600;
}
</style>
