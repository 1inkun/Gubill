<template>
	<view class="card-container">
		<view class="status">
			<view> 已经游玩: </view>
			<view class="duration">
				<text>{{ duration }}</text>
			</view>
			<view class="price">
				<view>预计消费:</view>
				<view class="num" v-if="totalValue != undefined">{{
					`${Math.ceil(totalValue / 100)} ¥`
				}}</view>
				<view class="num" v-else> {{ `0 ¥` }} </view>
			</view>
		</view>
		<view class="item">
			<view>本次签到开始于:</view>
			<view>{{ startTime }}</view>
		</view>
		<button
			class="btn btn-primary"
			hover-class="btn-primary-hover"
			@click="submit"
		>
			即刻结算
		</button>
	</view>
</template>

<script setup lang="ts">
import NewInstance from "@/api/instance";
import { FinishSignRes, SignData } from "@/types/components";
import CONFIG from "@/static/config.json";
import { ResponseData } from "@/types/global";
import { computed, onBeforeUnmount, onMounted, ref } from "vue";
import CalcTotalValue from "@/utils/calcTotalValue";

const props = defineProps<{
	data?: SignData;
}>();
const emit = defineEmits(["finish"]);
let intervalID: number = 0;
const duration = ref(`0小时0分钟`);

const totalValue = ref(0)

const startTime = computed(() => {
	const data = props.data;
	if (data != undefined) {
		const date = new Date(data.start_at * 1000);
		let year = date.getFullYear();
		let month = date.getMonth() + 1;
		let day = date.getDate();
		return `${year} - ${month} - ${day}`;
	}
});

const submit = async function () {
	uni.showLoading({
		title: "加载中",
		mask: true,
	});
	if (props.data?.signId == undefined) {
		uni.hideLoading();
		return;
	}
	const instance = NewInstance(
		CONFIG.Server.baseUrl,
		uni.getStorageSync("tokenStr"),
	);
	try {
		const resp = await instance<ResponseData<FinishSignRes>>({
			url: `/sign/${props.data.signId}`,
			method: "PUT",
		});
		const res = resp.data;
		if (res?.status != "success") {
			throw new Error(res?.msg);
		}
		const data = res.data;
		if (!Array.isArray(data)) {
			const payId = data.payId;
			emit("finish", payId);
			return;
		}
		throw new Error("未知错误");
	} catch (error: any) {
		console.error(error);
	} finally {
		uni.hideLoading();
	}
};

const calcDuration = function () {
	const now = ref(Math.ceil(Date.now() / 1000));
	const data = props.data;
	if (data != undefined) {
		let d = now.value - data.start_at;
		let mins: number = Math.floor(d / 60);
		let hours: number = 0;
		if (mins >= 60) {
			hours = Math.floor(mins / 60);
			mins -= hours * 60;
		}
		duration.value = `${hours}小时${mins}分钟`;
		totalValue.value = CalcTotalValue(d, 500);
	}
};
onMounted(() => {
    calcDuration()
	intervalID = setInterval(calcDuration, 6000);
});
onBeforeUnmount(()=>{
    clearInterval(intervalID)
})
</script>

<style scoped>
.duration {
	display: flex;
	flex-direction: column;
	align-items: center;
	color: var(--primary-color);
	font-size: 2rem;
	font-weight: 600;
}

.price {
	display: flex;
	flex-direction: row;
	justify-content: flex-end;
	align-items: center;
	margin-bottom: 16px;
}

.price .num {
	color: var(--error-color-hover);
	font-size: 1.25rem;
	padding-left: 8px;
}
</style>
