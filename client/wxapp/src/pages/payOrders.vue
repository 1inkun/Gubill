<template>
	<view class="card-container">
		<view class="header">
			<text>全部订单</text>
		</view>
		<view v-if="datas == undefined">
			<text>未找到订单</text>
		</view>
		<view class="secondary-card" v-else v-for="data in datas">
			{{ data }}
		</view>
	</view>
</template>

<script setup lang="ts">
import NewInstance from '@/api/instance';
import CONFIG from "@/static/config.json"
import { PayRes } from '@/types/components';
import { Pagin, PaginData, ResponseData } from '@/types/global';
import { onMounted, ref } from 'vue';

const instance = NewInstance(CONFIG.Server.baseUrl, uni.getStorageSync("tokenStr"))
const datas = ref<Array<PayRes>>()
const pagin = ref<Pagin>({
	page: 1,
	pageSize: 10,
})
const GetRecentPayOrder = async function () {
	uni.showLoading({
		title: '加载中',
		mask: true
	})
	try {
		const resp = await instance<ResponseData<PaginData<PayRes>>, Pagin>({
			url: "/pay",
			method: "GET",
			data: pagin.value
		})
		const res = resp.data
		if (res?.status != "success") {
			throw new Error(res?.msg)
		}
		const data = res.data
		if (!Array.isArray(data)) {
			datas.value = data.results
			return
		}
		throw new Error("未知错误")
	} catch (error: any) {
		console.error(error)
	} finally {
		uni.hideLoading()
	}
}
onMounted(() => {
	GetRecentPayOrder()
})
</script>

<style scoped></style>