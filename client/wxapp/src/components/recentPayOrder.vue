<template>
	<view class="card-container">
		<view class="header">
			<text>最近订单</text>
		</view>
		<view v-if="datas == undefined">
			<text>未找到数据</text>
		</view>
		<view v-else v-for="data in datas" >
            <view class="secondary-card order" @click="goToPay(data.payId)">
                <view class="item-container type">
                    <view v-if="data.businessType == 'sign'">
                        <text>签到订单</text>
                    </view>
                    <view v-else-if="data.businessType == 'member'">
                        <text>会员订单</text>
                    </view>
                    <view v-if="data.status == 0" class="badge">
                        <text>未支付</text>
                    </view>
                    <view v-if="data.status == 1" class="badge">
                        <text>已完成</text>
                    </view>
                    <view v-if="data.status == -1" class="badge">
                        <text>已取消</text>
                    </view>
                </view>
                <view class="item-container"></view>
                <view class="item-container">
                    <text>总价:</text>
                    <view class="price">
                        <text>{{ `${data.value / 100} ¥` }}</text>
                    </view>
                </view>
                <!-- {{ data }} -->
            </view>
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
const GetRecentPayOrder = async function () {
	uni.showLoading({
		title: '加载中',
		mask: true
	})
	const pagin: Pagin = {
		page: 1,
		pageSize: 3,
	}
	try {
		const resp = await instance<ResponseData<PaginData<PayRes>>, Pagin>({
			url: "/pay",
			method: "GET",
			data: pagin
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

const goToPay = function(payId:string) {
    uni.navigateTo({ url: `/pages/pay?payId=${payId}` })
}
onMounted(() => {
	GetRecentPayOrder()
})
</script>

<style scoped>
.order {
    display: flex;
    flex-direction: column;
}
.item-container {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    margin-bottom: 4px;
}
.type {
    font-size: 1.25rem;
}
.price {
    font-weight: 600;
    color: var(--error-color-hover);
}
</style>