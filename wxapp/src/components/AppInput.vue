<template>
  <view class="app-field">
    <text v-if="label" class="app-field__label">{{ label }}</text>
    <input
      class="app-field__input"
      :value="modelValue"
      :password="password"
      :type="type"
      :placeholder="placeholder"
      placeholder-class="app-field__placeholder"
      @input="onInput"
    />
  </view>
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    modelValue: string
    label?: string
    placeholder?: string
    password?: boolean
    type?: 'text' | 'number' | 'digit'
  }>(),
  { label: '', placeholder: '', password: false, type: 'text' },
)
const emit = defineEmits<{ (e: 'update:modelValue', v: string): void }>()
// uni-app 的 input 事件在 H5/小程序端都带 detail.value；类型上组件库声明为 Event，这里做防御式读取
function onInput(e: unknown) {
  const detail = (e as { detail?: { value?: string } })?.detail
  const target = (e as { target?: { value?: string } })?.target
  emit('update:modelValue', detail?.value ?? target?.value ?? '')
}
</script>

<style scoped>
.app-field {
  padding: 16rpx 0;
  border-bottom: 1rpx solid var(--c-border);
}
.app-field__label {
  font-size: 24rpx;
  color: var(--c-text-weak);
}
.app-field__input {
  height: 64rpx;
  font-size: 30rpx;
  color: var(--c-text);
  width: 100%;
}
.app-field__placeholder {
  color: var(--c-text-weak);
}
</style>
