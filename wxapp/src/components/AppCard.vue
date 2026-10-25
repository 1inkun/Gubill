<template>
  <view class="app-card">
    <view v-if="title" class="app-card__title">{{ title }}</view>
    <slot />
  </view>
</template>

<script setup lang="ts">
defineProps<{ title?: string }>()
</script>

<style scoped>
.app-card__title {
  font-size: 30rpx;
  font-weight: 600;
  margin-bottom: 20rpx;
}
</style>
