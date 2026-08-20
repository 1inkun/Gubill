<template>
	<view class="content">
		<view v-if="loginStatus">
			<UserInfo />
			<UserTools />
			<RecentPayOrder />
			<button class="btn btn-error" @click="Logout">登出</button>
		</view>
		<view v-else>
			<NeedLogin />
		</view>
	</view>
</template>

<script setup lang="ts">
import NeedLogin from '@/components/needLogin.vue';
import RecentPayOrder from '@/components/recentPayOrder.vue';
import UserInfo from '@/components/userInfo.vue';
import UserTools from '@/components/userTools.vue';
import CheckLoginStatus from '@/utils/checkLoginStatus';
import { onMounted, ref } from 'vue';

const loginStatus = ref(false)
const Login = function () {
	uni.navigateTo({ url: "login" })
}
const Logout = function () {
	// 删除本地存储中用户相关数据
	loginStatus.value = false
	uni.removeStorageSync("tokenStr")
	uni.removeStorageSync("userInfo")
	// 登出后跳转回主页
	uni.switchTab({ url: "index" })
	// 全局登出事件
	uni.$emit('userLogout')
}
onMounted(() => {
	uni.$on('userLogin', () => { loginStatus.value = true })
	loginStatus.value = CheckLoginStatus(uni.getStorageSync("tokenStr"))
})
</script>

<style scoped></style>
