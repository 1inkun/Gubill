<template>
  <view
    class="app-btn"
    :class="[`app-btn--${type}`, { 'app-btn--disabled': disabled || loading }]"
    @click="onClick"
  >
    <text v-if="loading">处理中...</text>
    <slot v-else />
  </view>
</template>

<script setup lang="ts">
const props = withDefaults(defineProps<{ type?: 'primary' | 'outline'; loading?: boolean; disabled?: boolean }>(), {
  type: 'primary',
  loading: false,
  disabled: false,
})
const emit = defineEmits<{ (e: 'click'): void }>()
function onClick() {
  if (props.disabled || props.loading) return
  emit('click')
}
</script>

<style scoped>
.app-btn {
  height: 88rpx;
  line-height: 88rpx;
  text-align: center;
  border-radius: 16rpx;
  font-size: 30rpx;
}
.app-btn--primary {
  background-color: var(--c-primary);
  color: var(--c-on-primary);
}
.app-btn--primary:active {
  background-color: var(--c-primary-hover);
}
.app-btn--outline {
  background-color: transparent;
  color: var(--c-primary-strong);
  border: 1rpx solid var(--c-primary);
}
.app-btn--disabled {
  opacity: 0.5;
}
</style>
