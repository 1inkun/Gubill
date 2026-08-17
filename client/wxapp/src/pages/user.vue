<template>
	<view class="user">user</view>
	<view v-if="LoginStatus">
		<button @click="Logout">登出</button>
	</view>
	<view v-else>
		<button @click="Login">登录</button>
	</view>
</template>

<script setup lang="ts">
import CheckLoginStatus from '@/utils/checkLoginStatus';
import { onMounted, ref } from 'vue';

const LoginStatus = ref(false)
const Login = function () {
	uni.navigateTo({ url: "login" })
}
const Logout = function () {
	// 删除本地存储中用户相关数据
	LoginStatus.value = false
	uni.removeStorageSync("tokenStr")
	uni.removeStorageSync("userInfo")
	// 登出后跳转回主页
	uni.switchTab({ url: "index" })
}
onMounted(() => {
	LoginStatus.value = CheckLoginStatus(uni.getStorageSync("tokenStr"))
})
</script>

<style scoped>
/* @import url("src/static/styles/login.css"); */
</style>
