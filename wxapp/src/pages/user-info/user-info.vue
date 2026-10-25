<template>
  <AppShell>
    <AppCard title="基本信息">
      <AppInput v-model="form.nickname" label="昵称" placeholder="展示用的昵称" />
      <AppInput v-model="form.email" label="邮箱" placeholder="用于找回账号" />
      <AppInput v-model="form.avatar" label="头像地址" placeholder="头像图片 URL（选填）" />
      <view v-if="errorText" class="uinfo__error">{{ errorText }}</view>
      <view class="app-gap" />
      <AppButton :loading="submitting" @click="onSave">保存修改</AppButton>
    </AppCard>

    <view class="app-gap" />
    <AppCard title="账号信息">
      <view class="uinfo__row"><text class="app-muted">用户名</text><text>{{ username || '--' }}</text></view>
      <view class="uinfo__row"><text class="app-muted">用户 ID</text><text>{{ userId || '--' }}</text></view>
      <view class="uinfo__row"><text class="app-muted">角色</text><text>{{ role || '--' }}</text></view>
    </AppCard>

    <view class="app-gap" />
    <text class="app-muted uinfo__note">当前为静态展示页，保存接口由后续接入</text>
  </AppShell>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import AppShell from '@/components/AppShell.vue'
import AppCard from '@/components/AppCard.vue'
import AppInput from '@/components/AppInput.vue'
import AppButton from '@/components/AppButton.vue'
import { getPayload, getUserId, getUsername } from '@/store/session'
import { toast } from '@/utils/toast'

const form = reactive({ nickname: '', email: '', avatar: '' })
const username = ref('')
const userId = ref(0)
const role = ref('')
const submitting = ref(false)
const errorText = ref('')

onShow(() => {
  username.value = getUsername()
  userId.value = getUserId()
  role.value = getPayload()?.role ?? ''
})

function onSave() {
  if (form.nickname.trim().length > 20) {
    errorText.value = '昵称长度不能超过 20 位'
    return
  }
  if (form.email.trim() && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email.trim())) {
    errorText.value = '请输入正确的邮箱'
    return
  }
  errorText.value = ''
  submitting.value = true
  // 静态页占位：用户信息修改接口由后续接入
  setTimeout(() => {
    submitting.value = false
    toast('用户信息接口待接入')
  }, 400)
}
</script>

<style scoped>
.uinfo__error {
  margin-top: 20rpx;
  color: var(--c-error-text);
  font-size: 26rpx;
}
.uinfo__row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20rpx 0;
  border-bottom: 1rpx solid var(--c-border);
  font-size: 28rpx;
}
.uinfo__note {
  font-size: 24rpx;
  text-align: center;
}
</style>
