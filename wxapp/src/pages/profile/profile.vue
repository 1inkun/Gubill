<template>
  <AppShell>
    <view @click="goUserInfo">
      <AppCard>
        <view class="profile__head">
          <view class="profile__avatar">{{ avatarText }}</view>
          <view class="profile__meta">
            <text class="profile__name">{{ nickname || username || '未登录' }}</text>
            <text class="app-muted profile__sub">账号：{{ username || '--' }}</text>
            <text class="app-muted profile__sub">用户 ID：{{ userId || '--' }}</text>
          </view>
          <view class="profile__edit app-muted">修改 ›</view>
        </view>
      </AppCard>
    </view>

    <view class="app-gap" />
    <AppCard title="我的会员">
      <StateView v-if="loading" state="loading" />
      <template v-else-if="member">
        <view class="app-row">
          <text class="profile__member-status">{{ member.status === 0 ? '生效中' : '已过期' }}</text>
          <text class="app-tag" :class="{ 'app-tag--done': member.status !== 0 }">
            {{ member.status === 0 ? '会员' : '失效' }}
          </text>
        </view>
        <text class="app-muted profile__card-sub">
          有效期至 {{ formatDateTime(parseBackendTime(member.end_at)) }}
        </text>
      </template>
      <template v-else>
        <text class="app-muted profile__card-sub">尚未开通会员，去「会员购买」选择套餐</text>
      </template>
    </AppCard>

    <view class="app-gap" />
    <AppCard title="待支付订单">
      <StateView v-if="loading" state="loading" />
      <template v-else-if="pendingOrders.length">
        <view
          v-for="p in pendingOrders"
          :key="p.id"
          class="profile__pay"
          @click="goPayOrderDetail(p.id)"
        >
          <view class="profile__pay-main">
            <text class="profile__pay-desc">{{ p.description || '会员购买' }}</text>
            <text class="app-muted profile__pay-no">单号：{{ p.order_no }}</text>
          </view>
          <view class="profile__pay-right">
            <text class="profile__pay-fee">￥{{ formatCents(p.total_fee) }}</text>
            <text class="profile__pay-status">待支付 ›</text>
          </view>
        </view>
        <view class="profile__pay-all app-muted" @click="goPayOrders">查看全部支付订单 ›</view>
      </template>
      <template v-else>
        <text class="app-muted profile__card-sub">暂无待支付订单</text>
        <view class="profile__pay-all app-muted" @click="goPayOrders">查看全部支付订单 ›</view>
      </template>
    </AppCard>

    <view class="app-gap" />
    <AppCard title="快捷入口">
      <view class="profile__row" @click="goOrders"><text>历史签到记录</text><text class="app-muted">›</text></view>
      <view class="profile__row" @click="goPayOrders"><text>我的支付订单</text><text class="app-muted">›</text></view>
    </AppCard>

    <view class="app-gap" />
    <AppButton type="outline" @click="onLogout">退出登录</AppButton>
  </AppShell>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import AppShell from '@/components/AppShell.vue'
import AppCard from '@/components/AppCard.vue'
import AppButton from '@/components/AppButton.vue'
import StateView from '@/components/StateView.vue'
import { getCheckinStatus } from '@/api/status'
import { getMemberData } from '@/api/member'
import { listPayOrders, PAY_ORDER_STATUS } from '@/api/pay'
import { clearSession, getUserId, getUsername, isLoggedIn } from '@/store/session'
import { formatCents } from '@/utils/price'
import { formatDateTime, parseBackendTime } from '@/utils/format'
import type { MemberData, PayOrderData } from '@/types/api'

const username = ref('')
const userId = ref(0)
const nickname = ref('')
const member = ref<MemberData | null>(null)
const pendingOrders = ref<PayOrderData[]>([])
const loading = ref(true)

const avatarText = computed(() => (nickname.value || username.value || '未').slice(0, 1))

/** 尽力而为：/status 只返回当前在场的人，昵称拿不到就留空 */
async function loadProfile() {
  username.value = getUsername()
  userId.value = getUserId()
  nickname.value = ''
  member.value = null
  pendingOrders.value = []
  if (!isLoggedIn()) {
    loading.value = false
    return
  }
  loading.value = true
  try {
    const users = (await getCheckinStatus()) ?? []
    const me = users.find((u) => u.username === username.value)
    nickname.value = me?.nickname ?? ''
  } catch {
    /* 拿不到昵称不影响主流程 */
  }
  // 会员信息与待支付订单并行加载；/member 查无记录时后端返回 500，视为未开通
  const [memberRes, payRes] = await Promise.allSettled([
    getMemberData(userId.value),
    listPayOrders({ userId: userId.value, status: PAY_ORDER_STATUS.PENDING, page: 1 }),
  ])
  if (memberRes.status === 'fulfilled' && memberRes.value && memberRes.value.id) {
    member.value = memberRes.value
  }
  if (payRes.status === 'fulfilled') {
    pendingOrders.value = (payRes.value?.data ?? []).slice(0, 2)
  }
  loading.value = false
}

function goUserInfo() {
  // 用户信息设置为静态展示页，接口由后续接入
  uni.navigateTo({ url: '/pages/user-info/user-info' })
}

function goOrders() {
  // 历史签到记录（原「我的订单」），从 tabBar 移入的子页面
  uni.navigateTo({ url: '/pages/orders/orders' })
}

function goPayOrders() {
  uni.navigateTo({ url: '/pages/pay-orders/pay-orders' })
}

function goPayOrderDetail(id: number) {
  uni.navigateTo({ url: `/pages/pay-order-detail/pay-order-detail?id=${id}` })
}

function onLogout() {
  uni.showModal({
    title: '退出登录',
    content: '确定要退出当前账号吗？',
    success: (res) => {
      if (!res.confirm) return
      clearSession()
      uni.reLaunch({ url: '/pages/login/login' })
    },
  })
}

onShow(loadProfile)
</script>

<style scoped>
.profile__head {
  display: flex;
  align-items: center;
}
.profile__avatar {
  width: 120rpx;
  height: 120rpx;
  border-radius: 50%;
  background-color: var(--c-primary-weak);
  color: var(--c-primary-strong);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 48rpx;
  margin-right: 24rpx;
}
.profile__meta {
  display: flex;
  flex-direction: column;
  flex: 1;
}
.profile__edit {
  font-size: 24rpx;
}
.profile__name {
  font-size: 34rpx;
  font-weight: 600;
}
.profile__sub {
  font-size: 24rpx;
  margin-top: 8rpx;
}
.profile__card-sub {
  font-size: 26rpx;
}
.profile__member-status {
  font-size: 36rpx;
  font-weight: 700;
  color: var(--c-primary-strong);
}
.profile__pay {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16rpx 0;
  border-bottom: 1rpx solid var(--c-border);
}
.profile__pay-main {
  display: flex;
  flex-direction: column;
}
.profile__pay-desc {
  font-size: 28rpx;
  font-weight: 600;
}
.profile__pay-no {
  font-size: 22rpx;
  margin-top: 4rpx;
}
.profile__pay-right {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
}
.profile__pay-fee {
  font-size: 28rpx;
  font-weight: 700;
  color: var(--c-primary-strong);
}
.profile__pay-status {
  font-size: 22rpx;
  color: var(--c-error-text);
  margin-top: 4rpx;
}
.profile__pay-all {
  margin-top: 16rpx;
  font-size: 24rpx;
  text-align: center;
}
.profile__row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24rpx 0;
  border-bottom: 1rpx solid var(--c-border);
}
</style>
