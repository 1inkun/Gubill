<template>
  <AppShell>
    <view class="register">
      <AppCard title="注册新账号">
        <AppInput v-model="form.username" label="用户名（4-20 位）" placeholder="必填，登录用" />
        <AppInput v-model="form.nickname" label="昵称" placeholder="选填" />
        <AppInput v-model="form.email" label="邮箱" placeholder="必填，用于找回" />
        <AppInput v-model="form.password" label="密码（6-32 位）" placeholder="必填" password />
        <AppInput v-model="form.confirmPassword" label="确认密码" placeholder="再次输入密码" password />
        <view v-if="errorText" class="register__error">{{ errorText }}</view>
        <view class="app-gap" />
        <AppButton :loading="submitting" @click="onSubmit">注 册</AppButton>
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
import { register } from '@/api/auth'
import { toast } from '@/utils/toast'

const form = reactive({ username: '', nickname: '', email: '', password: '', confirmPassword: '' })
const submitting = ref(false)
const errorText = ref('')

function validate(): string {
  const username = form.username.trim()
  if (username.length < 4 || username.length > 20) return '用户名长度需为 4-20 位'
  if (form.password.length < 6 || form.password.length > 32) return '密码长度需为 6-32 位'
  if (form.confirmPassword !== form.password) return '两次输入的密码不一致'
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email.trim())) return '请输入正确的邮箱'
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
    await register({
      username: form.username.trim(),
      nickname: form.nickname.trim(),
      email: form.email.trim(),
      password: form.password,
    })
    toast('注册成功，请登录', 'success')
    uni.reLaunch({ url: '/pages/login/login' })
  } catch (e) {
    errorText.value = e instanceof Error ? e.message : '注册失败'
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.register {
  padding-top: 60rpx;
}
.register__error {
  margin-top: 20rpx;
  color: var(--c-error-text);
  font-size: 26rpx;
}
</style>
