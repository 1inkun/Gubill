<template>
  <view v-if="visible" class="sheet__mask" @click="$emit('close')">
    <view class="sheet" @click.stop>
      <view class="sheet__bar" />
      <text class="sheet__title">收银台</text>
      <text class="sheet__amount">￥{{ formatCents(amountCents) }}</text>

      <view class="sheet__options">
        <view
          v-for="opt in options"
          :key="opt.value"
          class="sheet__option"
          :class="{ 'sheet__option--active': channel === opt.value }"
          @click="channel = opt.value"
        >
          <view class="sheet__icon" :class="`sheet__icon--${opt.value}`">{{ opt.icon }}</view>
          <view class="sheet__opt-meta">
            <text class="sheet__opt-name">{{ opt.label }}</text>
            <text class="app-muted sheet__opt-desc">{{ opt.desc }}</text>
          </view>
          <view class="sheet__radio" :class="{ 'sheet__radio--on': channel === opt.value }" />
        </view>
      </view>

      <view class="app-gap" />
      <AppButton :loading="paying" @click="onConfirm">确认支付</AppButton>
      <view class="app-gap" />
      <AppButton type="outline" @click="$emit('close')">取 消</AppButton>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import AppButton from '@/components/AppButton.vue'
import { formatCents } from '@/utils/price'
import type { PayChannel } from '@/api/pay'

const props = defineProps<{ visible: boolean; amountCents: number; paying?: boolean }>()
const emit = defineEmits<{ (e: 'close'): void; (e: 'pay', channel: PayChannel): void }>()

const options: { value: PayChannel; label: string; desc: string; icon: string }[] = [
  { value: 'alipay', label: '支付宝支付', desc: '推荐使用', icon: '支' },
  { value: 'wechat', label: '微信支付', desc: '快捷支付', icon: '微' },
]

const channel = ref<PayChannel>('alipay')

// 每次拉起收银台时重置为默认支付方式
watch(
  () => props.visible,
  (v) => {
    if (v) channel.value = 'alipay'
  },
)

function onConfirm() {
  emit('pay', channel.value)
}
</script>

<style scoped>
.sheet__mask {
  position: fixed;
  left: 0;
  top: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 9999;
  display: flex;
  align-items: flex-end;
}
.sheet {
  width: 100%;
  background-color: var(--c-surface);
  border-radius: 32rpx 32rpx 0 0;
  padding: 20rpx 32rpx calc(32rpx + 20rpx);
  display: flex;
  flex-direction: column;
}
.sheet__bar {
  align-self: center;
  width: 72rpx;
  height: 8rpx;
  border-radius: 8rpx;
  background-color: var(--c-border);
  margin-bottom: 20rpx;
}
.sheet__title {
  font-size: 32rpx;
  font-weight: 600;
}
.sheet__amount {
  margin-top: 12rpx;
  font-size: 60rpx;
  font-weight: 700;
  color: var(--c-primary-strong);
}
.sheet__options {
  margin-top: 28rpx;
  display: flex;
  flex-direction: column;
}
.sheet__option {
  display: flex;
  align-items: center;
  padding: 24rpx 20rpx;
  border: 1rpx solid var(--c-border);
  border-radius: 16rpx;
  margin-bottom: 20rpx;
}
.sheet__option--active {
  border-color: var(--c-primary);
  background-color: var(--c-primary-weak);
}
.sheet__icon {
  width: 72rpx;
  height: 72rpx;
  border-radius: 16rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32rpx;
  color: var(--c-on-primary);
  background-color: var(--c-primary);
  margin-right: 20rpx;
}
.sheet__opt-meta {
  display: flex;
  flex-direction: column;
  flex: 1;
}
.sheet__opt-name {
  font-size: 30rpx;
  font-weight: 600;
}
.sheet__opt-desc {
  font-size: 22rpx;
  margin-top: 4rpx;
}
.sheet__radio {
  width: 32rpx;
  height: 32rpx;
  border-radius: 50%;
  border: 2rpx solid var(--c-border);
  box-sizing: border-box;
}
.sheet__radio--on {
  border-color: var(--c-primary);
  background-color: var(--c-primary);
  box-shadow: inset 0 0 0 6rpx var(--c-surface);
}
</style>
