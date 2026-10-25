<template>
  <AppShell>
    <view class="login">
      <view class="login__brand">
        <text class="login__title">场馆计时计费</text>
        <text class="login__subtitle">登录后即可签到、查看订单</text>
      </view>
      <AppCard>
        <AppInput v-model="form.username" label="用户名" placeholder="4-20 位" />
        <AppInput v-model="form.password" label="密码" placeholder="6-32 位" password />
        <view v-if="errorText" class="login__error">{{ errorText }}</view>
        <view class="app-gap" />
        <AppButton :loading="submitting" @click="onSubmit">登 录</AppButton>
        <view class="login__link" @click="goRegister">没有账号？去注册</view>
      </AppCard>
    </view>
  </AppShell>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import AppShell from '@/components/AppShell.vue'
import AppCard from '@/components/AppCard.vue'
import AppInput from '@/components/AppInput.vue'
import AppButton from '@/components/AppButton.vue'
import { login } from '@/api/auth'
import { setToken } from '@/store/session'
import { toast } from '@/utils/toast'

const form = reactive({ username: '', password: '' })
const submitting = ref(false)
const errorText = ref('')

function validate(): string {
  const username = form.username.trim()
  if (username.length < 4 || username.length > 20) return '用户名长度需为 4-20 位'
  if (form.password.length < 6 || form.password.length > 32) return '密码长度需为 6-32 位'
  return ''
}

async function onSubmit() {
  const err = validate()
  if (err) {
    errorText.value = err
    return
  }
  errorText.value = ''
  submitting.value = true
  try {
    const res = await login({ username: form.username.trim(), password: form.password })
    setToken(res.token)
    toast('登录成功', 'success')
    uni.reLaunch({ url: '/pages/index/index' })
  } catch (e) {
    errorText.value = e instanceof Error ? e.message : '登录失败'
  } finally {
    submitting.value = false
  }
}

function goRegister() {
  uni.navigateTo({ url: '/pages/register/register' })
}
</script>

<style scoped>
.login {
  padding-top: 120rpx;
}
.login__brand {
  margin-bottom: 60rpx;
  display: flex;
  flex-direction: column;
}
.login__title {
  font-size: 52rpx;
  font-weight: 700;
  color: var(--c-primary-strong);
}
.login__subtitle {
  margin-top: 12rpx;
  color: var(--c-text-weak);
  font-size: 26rpx;
}
.login__error {
  margin-top: 20rpx;
  color: var(--c-error-text);
  font-size: 26rpx;
}
.login__link {
  margin-top: 28rpx;
  text-align: center;
  color: var(--c-primary-strong);
  font-size: 26rpx;
}
</style>
