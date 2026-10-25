<template>
  <AppShell>
    <view class="tabs">
      <view
        v-for="tab in tabs"
        :key="tab.value"
        class="tabs__item"
        :class="{ 'tabs__item--active': tab.value === activeTab }"
        @click="switchTab(tab.value)"
      >
        {{ tab.label }}
      </view>
    </view>

    <StateView v-if="loading && !orders.length" state="loading" />
    <StateView v-else-if="errorText" state="error" :text="errorText">
      <AppButton type="outline" @click="load(true)">重新加载</AppButton>
    </StateView>
    <StateView v-else-if="!orders.length" state="empty" text="暂无支付订单" />
    <template v-else>
      <view
        v-for="order in orders"
        :key="order.id"
        class="app-card pay__item"
        @click="goDetail(order.id)"
      >
        <view class="app-row">
          <text class="pay__desc">{{ order.description || '会员购买' }}</text>
          <text class="app-tag" :class="{ 'app-tag--done': order.status !== 0 }">
            {{ statusText(order.status) }}
          </text>
        </view>
        <view class="app-row pay__line">
          <text class="app-muted pay__no">单号：{{ order.order_no }}</text>
          <text class="pay__fee">￥{{ formatCents(order.total_fee) }}</text>
        </view>
      </view>
      <view class="pay__more app-muted" v-if="loading">加载中...</view>
      <view class="pay__more app-muted" v-else-if="!hasMore">没有更多了</view>
    </template>
  </AppShell>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { onPullDownRefresh, onReachBottom, onShow } from '@dcloudio/uni-app'
import AppShell from '@/components/AppShell.vue'
import AppButton from '@/components/AppButton.vue'
import StateView from '@/components/StateView.vue'
import { listPayOrders, PAY_ORDER_STATUS } from '@/api/pay'
import { PAGE_SIZE } from '@/config'
import { ensureLogin } from '@/store/guard'
import { getUserId } from '@/store/session'
import { errorMessage } from '@/utils/toast'
import { formatCents } from '@/utils/price'
import type { PayOrderData } from '@/types/api'

const tabs = [
  { label: '全部', value: PAY_ORDER_STATUS.ALL },
  { label: '待支付', value: PAY_ORDER_STATUS.PENDING },
  { label: '已支付', value: PAY_ORDER_STATUS.PAID },
]

const activeTab = ref<number>(PAY_ORDER_STATUS.ALL)
const orders = ref<PayOrderData[]>([])
const page = ref(1)
const hasMore = ref(true)
const loading = ref(false)
const errorText = ref('')

function statusText(status: number): string {
  if (status === PAY_ORDER_STATUS.PENDING) return '待支付'
  if (status === PAY_ORDER_STATUS.PAID) return '已支付'
  return '已关闭'
}

async function load(reset = false) {
  if (!ensureLogin()) return
  if (reset) {
    page.value = 1
    hasMore.value = true
    orders.value = []
  }
  if (!hasMore.value || loading.value) return
  loading.value = true
  errorText.value = ''
  try {
    const res = await listPayOrders({ userId: getUserId(), status: activeTab.value, page: page.value })
    const rows = res?.data ?? []
    orders.value = reset ? rows : orders.value.concat(rows)
    hasMore.value = rows.length >= PAGE_SIZE
    if (hasMore.value) page.value += 1
  } catch (e) {
    errorText.value = errorMessage(e)
  } finally {
    loading.value = false
  }
}

function switchTab(value: number) {
  if (activeTab.value === value) return
  activeTab.value = value
  load(true)
}

function goDetail(id: number) {
  uni.navigateTo({ url: `/pages/pay-order-detail/pay-order-detail?id=${id}` })
}

onShow(() => load(true))
onReachBottom(() => load())
onPullDownRefresh(async () => {
  await load(true)
  uni.stopPullDownRefresh()
})
</script>

<style scoped>
.tabs {
  display: flex;
  margin-bottom: 24rpx;
  background-color: var(--c-surface);
  border-radius: 16rpx;
  padding: 8rpx;
}
.tabs__item {
  flex: 1;
  text-align: center;
  padding: 16rpx 0;
  font-size: 28rpx;
  color: var(--c-text-weak);
  border-radius: 12rpx;
}
.tabs__item--active {
  background-color: var(--c-primary-weak);
  color: var(--c-primary-strong);
  font-weight: 600;
}
.pay__item {
  margin-bottom: 24rpx;
}
.pay__desc {
  font-size: 30rpx;
  font-weight: 600;
}
.pay__line {
  margin-top: 16rpx;
}
.pay__no {
  font-size: 24rpx;
}
.pay__fee {
  font-size: 32rpx;
  font-weight: 700;
  color: var(--c-primary-strong);
}
.pay__more {
  text-align: center;
  padding: 24rpx 0;
  font-size: 24rpx;
}
</style>
