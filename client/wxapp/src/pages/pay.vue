<template>
	<view class="content">
		<view class="card-container" v-if="payData == undefined">
			<view class="header">出现错误</view>
		</view>
		<view class="card-container" v-else>
			<view class="header">订单详情</view>
			<view class="price-container">
				<text>{{ `${payData.value / 100}¥` }}</text>
			</view>
			<view class="item-container">
				<view class="item">
					<text>订单编号:</text>
					<text class="payId">{{ payData.payId }}</text>
				</view>
				<view class="item">
					<text>订单类型:</text>
					<text v-if="payData.businessType == 'sign'">签到订单</text>
					<text v-if="payData.businessType == 'member'">会员订单</text>
				</view>
				<view class="item">
					<text>过期时间:</text>
					<text>{{ expireTime }}</text>
				</view>
			</view>
		</view>
		<button @click="pay" class="btn btn-primary"> 即刻支付 </button>
		<!-- <view>
			{{ payId }}
			{{ payData }}
		</view> -->
	</view>
</template>

<script setup lang="ts">
import NewInstance from "@/api/instance";
import { onLoad } from "@dcloudio/uni-app";
import { computed, onMounted, ref } from "vue";
import CONFIG from "@/static/config.json"
import { ResponseData } from "@/types/global";
import { PayRes } from "@/types/components";
import FormatDate from "@/utils/formatDate";

const instance = NewInstance(CONFIG.Server.baseUrl, uni.getStorageSync('tokenStr'))
const payId = ref<string>();
const payData = ref<PayRes>()
const expireTime = computed(() => {
	if (payData.value != undefined) {
		return FormatDate(payData.value.expire_time)
	}
	return ''
})

const getOrderData = async function () {
	uni.showLoading({
		title: '少女祈祷中',
		mask: true
	})
	try {
		const resp = await instance<ResponseData<PayRes>>({
			url: `/pay/${payId.value}`,
			method: 'GET'
		})
		const res = resp.data
		if (res?.status != "success") {
			throw new Error(res?.msg)
		}
		const data = res.data
		if (!Array.isArray(data)) {
			payData.value = data
			return
		}
		throw new Error("未知错误")
	} catch (error: any) {
		const res = error.response
		console.error(res)
	} finally {
		uni.hideLoading()
	}
}
const pay = async function () {
};
onMounted(() => {
	getOrderData()
})
onLoad((option) => {
	if (option != undefined) payId.value = option.payId;
});
</script>

<style scoped>
.price-container {
	display: flex;
	justify-content: center;
	align-items: center;
	font-size: 4rem;
	font-weight: 600;
	color: var(--error-color-hover);
}

.payId {
	max-width: 200px;
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
}
</style>
