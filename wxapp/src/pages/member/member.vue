<template>
  <AppShell>
    <view class="member__list">
      <text class="member__section">选择套餐</text>
      <StateView v-if="loading" state="loading" />
      <StateView v-else-if="errorText" state="error" :text="errorText">
        <AppButton type="outline" @click="load">重新加载</AppButton>
      </StateView>
      <StateView v-else-if="!goods.length" state="empty" text="暂无在售套餐" />
      <template v-else>
        <view v-for="g in goods" :key="g.id" class="app-card member__plan">
          <view class="app-row">
            <text class="member__name">{{ g.name }}</text>
            <text class="app-tag">{{ typeLabel(g.type) }}</text>
          </view>
          <text class="app-muted member__desc">{{ g.description }}</text>
          <view class="member__price">
            <text class="member__price-num">￥{{ formatCents(g.price) }}</text>
            <text class="app-muted member__price-unit">/ {{ typeUnit(g.type) }}</text>
          </view>
          <view class="app-gap" />
          <AppButton :loading="buyingId === g.id" @click="onBuy(g)">立即开通</AppButton>
        </view>
      </template>
    </view>

    <view class="app-gap" />
    <AppCard title="会员说明">
      <text class="app-muted member__note">
        1. 会员有效期自支付完成时刻起算；
        2. 会员折扣与计时计费不叠加，以结算页展示为准；
        3. 下单后在「我的-我的支付订单」中可查看并完成支付。
      </text>
    </AppCard>
  </AppShell>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import AppShell from '@/components/AppShell.vue'
import AppCard from '@/components/AppCard.vue'
import AppButton from '@/components/AppButton.vue'
import StateView from '@/components/StateView.vue'
import { buyGood, getGoods } from '@/api/good'
import { ensureLogin } from '@/store/guard'
import { errorMessage, toast } from '@/utils/toast'
import { formatCents } from '@/utils/price'
import type { Good } from '@/types/api'

const loading = ref(true)
const errorText = ref('')
const goods = ref<Good[]>([])
const buyingId = ref(0)

function typeLabel(type: string): string {
  const map: Record<string, string> = { weekly: '周卡', monthly: '月卡', mothly: '月卡', yearly: '年卡' }
  return map[type] ?? type
}

function typeUnit(type: string): string {
  const map: Record<string, string> = { weekly: '周', monthly: '月', mothly: '月', yearly: '年' }
  return map[type] ?? '期'
}

async function load() {
  if (!ensureLogin()) return
  loading.value = true
  errorText.value = ''
  try {
    goods.value = (await getGoods()) ?? []
  } catch (e) {
    errorText.value = errorMessage(e)
  } finally {
    loading.value = false
  }
}

async function onBuy(g: Good) {
  if (!ensureLogin()) return
  buyingId.value = g.id
  try {
    const res = await buyGood(g.id)
    toast('支付订单已创建', 'success')
    // 统一跳转支付订单详情页，在那里完成支付
    const payOrderId = res?.pay_order_id
    if (payOrderId) {
      uni.navigateTo({ url: `/pages/pay-order-detail/pay-order-detail?id=${payOrderId}` })
    }
  } catch (e) {
    toast(errorMessage(e))
  } finally {
    buyingId.value = 0
  }
}

onShow(load)
</script>

<style scoped>
.member__status {
  font-size: 40rpx;
  font-weight: 700;
  color: var(--c-primary-strong);
}
.member__range {
  margin-top: 12rpx;
  font-size: 26rpx;
}
.member__section {
  font-size: 30rpx;
  font-weight: 600;
  margin-bottom: 20rpx;
}
.member__list {
  display: flex;
  flex-direction: column;
}
.member__plan {
  margin-bottom: 24rpx;
}
.member__name {
  font-size: 32rpx;
  font-weight: 600;
}
.member__desc {
  margin-top: 12rpx;
  font-size: 26rpx;
}
.member__price {
  margin-top: 20rpx;
  display: flex;
  align-items: baseline;
}
.member__price-num {
  font-size: 52rpx;
  font-weight: 700;
  color: var(--c-primary-strong);
}
.member__price-unit {
  margin-left: 8rpx;
  font-size: 24rpx;
}
.member__note {
  font-size: 26rpx;
  line-height: 1.8;
}
</style>
