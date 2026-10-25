<template>
  <view class="timer">
    <text class="timer__label app-muted">累计签到时间</text>
    <text class="timer__clock">{{ elapsedText }}</text>
    <text class="timer__price">预估 ￥{{ formatCents(estimated) }}</text>
  </view>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { parseBackendTime, formatDuration } from '@/utils/format'
import { estimatePriceCents, formatCents } from '@/utils/price'

const props = defineProps<{ startAt: string }>()
const now = ref(Date.now())
let timer: ReturnType<typeof setInterval> | null = null

const startMs = computed(() => parseBackendTime(props.startAt)?.getTime() ?? now.value)
const elapsedText = computed(() => formatDuration((now.value - startMs.value) / 1000))
const estimated = computed(() => estimatePriceCents(startMs.value, now.value))

onMounted(() => {
  timer = setInterval(() => {
    now.value = Date.now()
  }, 1000)
})
onUnmounted(() => {
  if (timer) clearInterval(timer)
  timer = null
})
</script>

<style scoped>
.timer {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20rpx 0;
}
.timer__clock {
  font-size: 64rpx;
  font-weight: 700;
  color: var(--c-primary-strong);
  letter-spacing: 4rpx;
  margin-top: 8rpx;
}
.timer__hint {
  color: var(--c-text-weak);
  font-size: 26rpx;
  margin-top: 8rpx;
}
.timer__price {
  margin-top: 12rpx;
  font-size: 30rpx;
  color: var(--c-text);
}
</style>
