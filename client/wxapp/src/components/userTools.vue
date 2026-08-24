<template>
	<view class="card-container">
		<view class="tool-item" @click="goToPayOrders()">
			<!-- 图标 -->
			<text>全部订单</text>
		</view>
        <view class="tool-item" @click="goToSignOrders()">
			<!-- 图标 -->
			<text>签到记录</text>
		</view>
	</view>
</template>

<script setup lang="ts">
const goToPayOrders = function () {
	uni.navigateTo({ url: "payOrders" })
}
const goToSignOrders = function () {
	uni.navigateTo({ url: "signOrders" })
}
</script>

<style scoped>
.card-container {
	flex-direction: row;
    justify-content: space-around;
}

.tool-item {
	display: flex;
	flex-direction: column;
	justify-content: center;
	align-items: center;
}
</style>