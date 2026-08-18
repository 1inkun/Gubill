<template>
    <view>
        <view>{{ props.data }}</view>
        <view> {{ now }}</view>
        <view>{{ duration }}</view>
        <view v-if="totalValue != undefined">{{ `${Math.ceil(totalValue / 100)} ¥` }}</view>
        <view v-else> {{ `0 ¥` }} </view>
        <button @click="submit">即刻结算</button>
    </view>
</template>

<script setup lang="ts">
import NewInstance from '@/api/instance';
import { FinishSignRes, SignData } from '@/types/components';
import CONFIG from '@/static/config.json'
import { ResponseData } from '@/types/global';
import { computed, onMounted, ref } from 'vue';
import CalcTotalValue from '@/utils/calcTotalValue';

const props = defineProps<{
    data ?: SignData 
}>()
const emit = defineEmits(['finish'])
let intervalID:number = 0
const now = ref(Math.ceil(Date.now() / 1000))

const duration = computed(()=>{
    const data = props.data
    if(data != undefined){
       let d = now.value - data.start_at
       let mins:number = Math.floor(d / 60)
       let hours:number = 0
       if (mins >= 60){
            hours = Math.floor(mins / 60)
            mins -= hours * 60
       }
       return `${hours}小时 ${mins}分钟`
    }
})

const totalValue = computed(()=>{
    const data = props.data
    if(data != undefined){
        const d = now.value - data.start_at
        return CalcTotalValue(d,500)
    }
})

const submit = async function(){
    uni.showLoading({
        title: '加载中',
        mask: true
    })
    if(props.data?.signId == undefined){
        uni.hideLoading()
        return
    }
    const instance = NewInstance(CONFIG.Server.baseUrl,uni.getStorageSync('tokenStr'))
    try {
        const resp = await instance<ResponseData<FinishSignRes>>({
            url: `/sign/${props.data.signId}`,
            method: "PUT"
        })
        const res = resp.data
        if(res?.status != "success"){
            throw new Error(res?.msg)
        }
        const data = res.data
        if(!Array.isArray(data)){
            const payId = data.payId
            emit('finish',payId)
            return
        }
        throw new Error("未知错误")
    } catch (error:any) {
        console.error(error)
    } finally {
        uni.hideLoading()
    }
}

</script>

<style scoped></style>