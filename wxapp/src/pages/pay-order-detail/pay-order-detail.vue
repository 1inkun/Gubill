<template>
  <AppShell>
    <StateView v-if="loading" state="loading" />
    <StateView v-else-if="errorText" state="error" :text="errorText">
      <AppButton type="outline" @click="load">重新加载</AppButton>
    </StateView>
    <template v-else-if="order">
      <AppCard title="支付订单">
        <view class="detail__center">
          <text class="detail__fee">￥{{ formatCents(order.total_fee) }}</text>
          <text class="app-tag" :class="{ 'app-tag--done': order.status !== 0 }">
            {{ statusText(order.status) }}
          </text>
        </view>
        <view class="detail__row"><text class="app-muted">商品描述</text><text>{{ order.description || '--' }}</text></view>
        <view class="detail__row"><text class="app-muted">订单编号</text><text>{{ order.order_no }}</text></view>
        <view class="detail__row"><text class="app-muted">订单类型</text><text>{{ order.type || '--' }}</text></view>
        <view v-if="order.transaction_id" class="detail__row">
          <text class="app-muted">支付流水号</text><text>{{ order.transaction_id }}</text>
        </view>
      </AppCard>

      <view class="app-gap" />
      <template v-if="order.status === 0">
        <AppButton :loading="paying" @click="showSheet = true">去支付</AppButton>
      </template>
      <text v-else class="app-muted detail__done">订单已完成支付，感谢您的支持</text>
    </template>
  </AppShell>

  <CheckoutSheet
    :visible="showSheet"
    :amount-cents="order?.total_fee ?? 0"
    :paying="paying"
    @close="showSheet = false"
    @pay="onPay"
  />
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import AppShell from '@/components/AppShell.vue'
import AppCard from '@/components/AppCard.vue'
import AppButton from '@/components/AppButton.vue'
import StateView from '@/components/StateView.vue'
import CheckoutSheet from '@/components/CheckoutSheet.vue'
import { getPayOrderById, payOrder, PAY_ORDER_STATUS, type PayChannel } from '@/api/pay'
import { ensureLogin } from '@/store/guard'
import { errorMessage, toast } from '@/utils/toast'
import { formatCents } from '@/utils/price'
import type { PayOrderData } from '@/types/api'

const payOrderId = ref(0)
const order = ref<PayOrderData | null>(null)
const loading = ref(true)
const paying = ref(false)
const showSheet = ref(false)
const errorText = ref('')

function statusText(status: number): string {
  if (status === PAY_ORDER_STATUS.PENDING) return '待支付'
  if (status === PAY_ORDER_STATUS.PAID) return '已支付'
  return '已关闭'
}

async function load() {
  if (!ensureLogin()) return
  loading.value = true
  errorText.value = ''
  try {
    order.value = await getPayOrderById(payOrderId.value)
  } catch (e) {
    errorText.value = errorMessage(e)
  } finally {
    loading.value = false
  }
}

async function onPay(channel: PayChannel) {
  const o = order.value
  if (!o) return
  paying.value = true
  try {
    await payOrder(o.id, channel)
    showSheet.value = false
    toast('支付成功', 'success')
    await load()
  } catch (e) {
    toast(errorMessage(e))
  } finally {
    paying.value = false
  }
}

onLoad((query) => {
  payOrderId.value = Number((query as Record<string, string>)?.id ?? 0)
  load()
})
</script>

<style scoped>
.detail__center {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 12rpx 0 24rpx;
}
.detail__fee {
  font-size: 72rpx;
  font-weight: 700;
  color: var(--c-primary-strong);
  margin-bottom: 12rpx;
}
.detail__row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16rpx 0;
  border-bottom: 1rpx solid var(--c-border);
  font-size: 28rpx;
}
.detail__done {
  display: block;
  text-align: center;
  font-size: 26rpx;
  padding: 20rpx 0;
}
</style>
