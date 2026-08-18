<template>
	<view class="card-container">
		<!-- 标题 -->
		<view class="header">
			<text>请登录</text>
		</view>
		<!-- 正文 -->
		<view class="secondary-card">
			<text class="emoji">
				🤔
			</text>
			<text>
				当前功能需要登录后才能够使用，请先登录
			</text>
		</view>
		<!-- 跳转登录 -->
		<button class="btn btn-primary" hover-class="btn-primary-hover" @click="goToLogin">
			即刻登录
		</button>
	</view>
</template>

<script setup lang="ts">
const goToLogin = function () {
	uni.navigateTo({ url: "login" })
}
</script>

<style>
.header {
	font-size: 1.25rem;
	margin-top: 8px;
	margin-bottom: 16px;
	padding-left: 8px;
	border-left: 8px solid var(--primary-color);
}

.btn {
	display: inline-block;
}

.secondary-card {
	align-items: center;
}

.emoji {
	font-size: 4rem;
}

.btn {
	margin-bottom: 20px;
}
</style>