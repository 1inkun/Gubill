<template>
	<view class="card-container">
		<view class="tool-item" @click="goToPayOrders()">
			<!-- 图标 -->
			<text>全部订单</text>
		</view>
	</view>
</template>

<script setup lang="ts">
const goToPayOrders = function () {
	uni.navigateTo({ url: "payOrders" })
}
</script>

<style scoped>
.card-container {
	flex-direction: row;
}

.tool-item {
	display: flex;
	flex-direction: column;
	justify-content: center;
	align-items: center;
}
</style>