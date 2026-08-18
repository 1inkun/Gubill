<template>
	<view class="card-container">
		<!-- 标题 -->
		<view class="header">
			<text>登录</text>
		</view>
		<!-- 表单 -->
		<form class="form" @submit="Login">
			<input class="input" type="text" name="username" placeholder="请输入用户名">
			<input class="input" type="text" name="password" placeholder="请输入密码" :password="true">

			<button class="btn btn-primary" hover-class="btn-primary-hover" form-type="submit">登录</button>
		</form>
		<!-- 其他功能 -->
		<view class="tool-container">
			<view @click="() => { emit('register') }">前往注册</view>
			<view>忘记密码?</view>
		</view>
		<!-- 用户许可协议 -->
		<!-- <view class="foot-note">登录即代表同意《用户协议》和《隐私政策》</view> -->
	</view>
</template>

<script setup lang="ts">
import NewInstance from '@/api/instance';
import { LoginData, LoginRes } from '@/types/components';
import CONFIG from '@/static/config.json'
import { ResponseData, UserInfo, CustomJWTClaims } from '@/types/global';
import { jwtDecode } from "jwt-decode";

const emit = defineEmits(['register'])
const instance = NewInstance(CONFIG.Server.baseUrl, "")
const checkPasswd = /^(?=.*[a-z])(?=.*\d).{6,30}$/
const checkUsername = /^[A-Za-z0-9]{4,20}$/

const checkLoginInfo = function (loginData: LoginData): boolean {
	if (!checkPasswd.test(loginData.password)) {
		console.error("密码有误")
		return false
	}
	if (!checkUsername.test(loginData.username)) {
		console.error("用户名有误")
		return false
	}
	return true
}

const SaveUserData = function (tokenStr: string): boolean {
	if (tokenStr == "") {
		return false
	}
	uni.setStorageSync("tokenStr", tokenStr)
	const decode = jwtDecode<CustomJWTClaims>(tokenStr)
	// console.log(decode)
	if (decode.userId == undefined) {
		return false
	}
	let userInfo: UserInfo = {
		username: decode.username,
		userId: decode.userId,
		nickname: decode.nickname,
		role: decode.role
	}
	uni.setStorageSync("userInfo", userInfo)
	return true
}


const Login = async function (e: any) {
	uni.showLoading({
		title: '加载中',
		mask: true
	})
	const loginData: LoginData = e.detail.value
	let ok = checkLoginInfo(loginData)
	if (!ok) {
		uni.hideLoading()
		return
	}
	try {
		const resp = await instance<ResponseData<LoginRes>, LoginData>({
			url: "/user/login",
			method: "POST",
			data: loginData
		})
		const res = resp.data
		if (res?.status != "success") {
			throw new Error(res?.msg)
		}
		const data = res.data
		if (!Array.isArray(data)) {
			let ok = SaveUserData(data.token)
			if (!ok) {
				throw new Error("登录失败")
			}
			uni.$emit('userLogin')
			uni.switchTab({ url: "index" })
		}
	} catch (error: any) {
		console.error(error)
	} finally {
		uni.hideLoading()
	}
}
</script>

<style scoped>
.header {
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
	margin-top: 12px;
	margin-bottom: 8px;
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