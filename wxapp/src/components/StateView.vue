<template>
  <view class="state-view">
    <text class="state-view__text">{{ text || defaultText }}</text>
    <slot />
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'
const props = defineProps<{ state: 'loading' | 'empty' | 'error'; text?: string }>()
const defaultText = computed(() => {
  if (props.state === 'loading') return '加载中...'
  if (props.state === 'empty') return '暂无数据'
  return '加载失败'
})
</script>

<style scoped>
.state-view {
  padding: 80rpx 0;
  display: flex;
  flex-direction: column;
  align-items: center;
}
.state-view__text {
  color: var(--c-text-weak);
  font-size: 28rpx;
  margin-bottom: 24rpx;
}
</style>
