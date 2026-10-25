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
    <StateView v-else-if="!orders.length" state="empty" text="还没有订单，去首页签到吧" />
    <template v-else>
      <OrderItem v-for="order in orders" :key="order.id" :order="order" @select="goDetail" />
      <view class="orders__more app-muted" v-if="loading">加载中...</view>
      <view class="orders__more app-muted" v-else-if="!hasMore">没有更多了</view>
    </template>
  </AppShell>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { onPullDownRefresh, onReachBottom, onShow } from '@dcloudio/uni-app'
import AppShell from '@/components/AppShell.vue'
import AppButton from '@/components/AppButton.vue'
import OrderItem from '@/components/OrderItem.vue'
import StateView from '@/components/StateView.vue'
import { listOrders, ORDER_STATUS } from '@/api/order'
import { PAGE_SIZE } from '@/config'
import { ensureLogin } from '@/store/guard'
import { getUserId } from '@/store/session'
import { errorMessage } from '@/utils/toast'
import type { OrderResponse } from '@/types/api'

const tabs = [
  { label: '全部', value: ORDER_STATUS.ALL },
  { label: '进行中', value: ORDER_STATUS.RUNNING },
  { label: '已结算', value: ORDER_STATUS.FINISHED },
]

const activeTab = ref<number>(ORDER_STATUS.ALL)
const orders = ref<OrderResponse[]>([])
const page = ref(1)
const hasMore = ref(true)
const loading = ref(false)
const errorText = ref('')

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
    const res = await listOrders({ userId: getUserId(), status: activeTab.value, page: page.value })
    const rows = res?.data ?? []
    orders.value = reset ? rows : orders.value.concat(rows)
    // 后端 SQL 里 LIMIT 写死 10，page_size 无效，所以用返回条数判断是否还有下一页
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
  uni.navigateTo({ url: `/pages/order-detail/order-detail?id=${id}` })
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
.orders__more {
  text-align: center;
  padding: 24rpx 0;
  font-size: 24rpx;
}
</style>
