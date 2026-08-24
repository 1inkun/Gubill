<template>
	<view class="content">
		<view class="card-container" v-if="payId == undefined">
			<view class="header">出现错误</view>
		</view>
		<view class="card-container">
			<view class="header">订单详情</view>
		</view>
		<button @click="pay" class="btn btn-primary"> 即刻支付 </button>
		<view>
			{{ payId }}
            {{ payData }}
		</view>
	</view>
</template>

<script setup lang="ts">
import NewInstance from "@/api/instance";
import { onLoad } from "@dcloudio/uni-app";
import { onMounted, ref } from "vue";
import CONFIG from "@/static/config.json"
import { ResponseData } from "@/types/global";
import { PayRes } from "@/types/components";

const instance = NewInstance(CONFIG.Server.baseUrl,uni.getStorageSync('tokenStr'))
const payId = ref<string>();
const payData = ref<PayRes>()

const getOrderData = async function(){
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
        if (res?.status != "success"){
            throw new Error(res?.msg)
        }
        const data = res.data
        if (!Array.isArray(data)){
            payData.value = data
            return
        }
        throw new Error("未知错误")
    } catch (error:any) {
        const res = error.response
        console.error(res)
    } finally {
        uni.hideLoading()
    }
}
const pay = async function () {};
onMounted(()=>{
    getOrderData()
})
onLoad((option) => {
	if (option != undefined) payId.value = option.payId;
});
</script>

<style scoped></style>
