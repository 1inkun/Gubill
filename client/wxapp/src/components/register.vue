<template>
	<view class="card-container">
		<!-- 标题 -->
		<view class="header">
			<text>注册</text>
		</view>
		<!-- <button @click="() => { emit('login') }">返回</button> -->
		<!-- 表单 -->
		<form class="form" @submit="Register">
			<input class="input" type="text" name="username" placeholder="输入账户名">
			<input class="input" type="text" name="nickname" placeholder="输入昵称">
			<input class="input" type="text" name="email" placeholder="输入邮箱">
			<input class="input" type="text" name="password" placeholder="输入密码" :password="true" v-model="password">
			<input class="input" type="text" placeholder="确认密码" :password="true" v-model="checkPassword">
			<button class="btn btn-primary" hover-class="btn-primary-hover" form-type="submit"
				:disabled="samepasswd">即刻注册</button>
		</form>
	</view>
</template>

<script setup lang="ts">
import NewInstance from '@/api/instance'
import { RegisterData, RegisterRes } from '@/types/components'
import { computed, ref } from 'vue'
import CONFIG from '@/static/config.json'
import { ResponseData } from '@/types/global'
import { ShowSuccessMsg } from '@/utils/toast'

const emit = defineEmits(['success', 'login'])
const instance = NewInstance(CONFIG.Server.baseUrl, "")

const password = ref("")
const checkPassword = ref("")
const samepasswd = computed(() => {
	return password.value == checkPassword.value ? false : true
})
const checkPasswd = /^(?=.*[a-z])(?=.*\d).{6,30}$/
const checkUsername = /^[A-Za-z0-9]{4,20}$/
const checkEmail = /\w[-\w.+]*@([A-Za-z0-9][-A-Za-z0-9]+\.)+[A-Za-z]{2,14}/

const checkRegisterData = function (data: RegisterData): boolean {
	console.log(data)
	if (!checkUsername.test(data.username)) {
		return false
	}
	if (!checkEmail.test(data.email)) {
		return false
	}
	if (!checkPasswd.test(data.password)) {
		return false
	}
	return true
}
const Register = async function (e: any) {
	uni.showLoading({
		title: '加载中',
		mask: true
	})
	const registerData: RegisterData = e.detail.value
	let ok = checkRegisterData(registerData)
	if (!ok) {
		uni.hideLoading()
		return
	}
	try {
		const resp = await instance<ResponseData<RegisterRes>, RegisterData>({
			url: "/user/register",
			method: "POST",
			data: registerData
		})
		const res = resp.data
		if (res?.status != "success") {
			throw new Error(res?.msg)
		}
		if (!Array.isArray(res.data)) {
			if (res.data.userId != undefined) {
				// 登录成功触发success事件
				emit('success')
                ShowSuccessMsg("注册成功")
				return
			}
		}
		throw new Error("未知错误")
	} catch (error: any) {
		const res = error.response.data
		console.error(res)
	} finally {
		uni.hideLoading()
	}
}
</script>

<style scoped>
.header {
	justify-content: space-between;
	border-left: 8px solid var(--primary-color);
	padding-left: 8px;
	margin-top: 8px;
	margin-bottom: 16px;
	font-size: 1.25rem;
}

.form {
	display: flex;
	flex-direction: column;
}

.input {
	width: auto;
	/* display: inline-block; */
	padding: 12px;
	margin-bottom: 16px;
	border-radius: 12px;
	border: 1px, solid, lightgray;
}

.tool-container {
	display: flex;
	justify-content: space-between;
	margin-top: 16px;
	margin: 4px 0;
	color: darkgray;
}

.foot-note {
	text-align: center;
	font-size: 11px;
	color: #c0c0c0;
	margin-top: 8px;
	line-height: 1.6;
}
</style>