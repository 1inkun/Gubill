<template>
  <AppShell>
    <view @click="showUsers = true">
      <AppCard title="在场情况">
        <view class="home__center">
          <text class="home__count">{{ onlineUsers.length }}</text>
          <text class="app-muted home__count-label">人正在签到中</text>
          <text class="app-muted home__more">点击查看在场用户 ›</text>
        </view>
      </AppCard>
    </view>

    <view class="app-gap" />

    <AppCard title="我的签到">
      <StateView v-if="loading" state="loading" />
      <StateView v-else-if="errorText" state="error" :text="errorText">
        <AppButton type="outline" @click="refresh">重新加载</AppButton>
      </StateView>
      <template v-else>
        <CheckinTimer v-if="runningOrder" :start-at="runningOrder.start_at" />
        <view v-if="runningOrder" class="home__actions">
          <AppButton :loading="submitting" @click="onFinish">结算并离场</AppButton>
        </view>
        <template v-else>
          <text class="app-muted home__tip">当前未签到，签到后开始按时计费</text>
          <view class="home__actions">
            <AppButton :loading="submitting" @click="onCheckin">即刻签到</AppButton>
          </view>
        </template>
      </template>
    </AppCard>

    <view class="app-gap" />
    <AppCard title="计费说明">
      <text class="app-muted home__tip">5 元 / 半小时（按半小时向上取整），最小计费 10 元</text>
    </AppCard>

    <UserListModal :visible="showUsers" :users="onlineUsers" @close="showUsers = false" />
  </AppShell>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import AppShell from '@/components/AppShell.vue'
import AppCard from '@/components/AppCard.vue'
import AppButton from '@/components/AppButton.vue'
import CheckinTimer from '@/components/CheckinTimer.vue'
import UserListModal from '@/components/UserListModal.vue'
import StateView from '@/components/StateView.vue'
import { createOrder, finishOrder, listOrders, ORDER_STATUS } from '@/api/order'
import { getCheckinStatus } from '@/api/status'
import { ensureLogin } from '@/store/guard'
import { getUserId } from '@/store/session'
import { errorMessage, toast } from '@/utils/toast'
import type { CurrentUser, OrderResponse } from '@/types/api'

const loading = ref(true)
const submitting = ref(false)
const errorText = ref('')
const onlineUsers = ref<CurrentUser[]>([])
const runningOrder = ref<OrderResponse | null>(null)
const showUsers = ref(false)

async function loadOnline() {
  onlineUsers.value = (await getCheckinStatus()) ?? []
}

async function loadMyOrder() {
  const res = await listOrders({ userId: getUserId(), status: ORDER_STATUS.RUNNING, page: 1 })
  const rows = res?.data ?? []
  runningOrder.value = rows.length ? rows[0] : null
}

async function refresh() {
  if (!ensureLogin()) return
  loading.value = true
  errorText.value = ''
  try {
    await Promise.all([loadOnline(), loadMyOrder()])
  } catch (e) {
    errorText.value = errorMessage(e)
  } finally {
    loading.value = false
  }
}

async function onCheckin() {
  if (!ensureLogin()) return
  submitting.value = true
  try {
    runningOrder.value = await createOrder(getUserId())
    toast('签到成功', 'success')
    await loadOnline()
  } catch (e) {
    toast(errorMessage(e))
  } finally {
    submitting.value = false
  }
}

async function onFinish() {
  const order = runningOrder.value
  if (!order || !ensureLogin()) return
  submitting.value = true
  try {
    const res = await finishOrder(order.id)
    runningOrder.value = null
    await loadOnline()
    toast('结算成功', 'success')
    // 统一跳转支付订单详情页，在那里完成支付
    const payOrderId = res?.pay_order_id
    if (payOrderId) {
      uni.navigateTo({ url: `/pages/pay-order-detail/pay-order-detail?id=${payOrderId}` })
    }
  } catch (e) {
    toast(errorMessage(e))
  } finally {
    submitting.value = false
  }
}

onShow(refresh)
</script>

<style scoped>
.home__center {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20rpx 0 8rpx;
}
.home__count {
  font-size: 120rpx;
  font-weight: 700;
  color: var(--c-primary-strong);
  line-height: 1.1;
}
.home__count-label {
  margin-top: 8rpx;
  font-size: 28rpx;
}
.home__more {
  margin-top: 20rpx;
  font-size: 24rpx;
}
.home__tip {
  font-size: 26rpx;
}
.home__actions {
  margin-top: 28rpx;
}
</style>
