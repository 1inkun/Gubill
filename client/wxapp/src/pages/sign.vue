<template>
    <view v-if="currentSignData?.signId != undefined">
        <SignDataDetail :data="currentSignData" @finish="finishSign" />
    </view>
    <view v-else>
        <NewSign @sign-in="getCurrentSignData"/>
    </view>
</template>

<script setup lang="ts">
import NewSign from '@/components/newSign.vue';
import NewInstance from '@/api/instance';
import CheckLoginStatus from '@/utils/checkLoginStatus';
import { onMounted, ref } from 'vue';
import CONFIG from '@/static/config.json'
import { ResponseData } from '@/types/global';
import { SignData } from '@/types/components'
import SignDataDetail from '@/components/signDataDetail.vue';

const currentSignData = ref<SignData>()
const loginStatus = ref(false)
const checkCurrentSignData = async function(tokenStr:string){
    uni.showLoading({
        title: '加载中',
        mask: true
    })
    const instance = NewInstance(CONFIG.Server.baseUrl,tokenStr)
    try {
        const resp = await instance<ResponseData<SignData>>({
            url: '/sign?status=0',
            method: 'GET'
        })
        const res = resp.data
        if(res?.status != "success") {
            throw new Error(res?.msg)
        }
        const data = res.data
        if(Array.isArray(data)){
            currentSignData.value = data[0]
        }
    } catch (error:any) {
        console.error(error)
    } finally {
        uni.hideLoading()
    }
}

const getCurrentSignData = async function(signId:string){
    uni.showLoading({
        title: '加载中',
        mask: true
    })
    const instance = NewInstance(CONFIG.Server.baseUrl,uni.getStorageSync("tokenStr"))
    try {
        const resp = await instance<ResponseData<SignData>>({
            url:`/sign/${signId}`,
            method:"GET"
        })
        const res = resp.data
        if (res?.status != "success"){
            throw new Error(res?.msg)
        }
        const data = res.data
        if (!Array.isArray(data)) {
            currentSignData.value = data
        }
    } catch (error:any) {
        console.error(error)
    } finally {
        uni.hideLoading()
    }
}

const finishSign = function(payId:string){
    currentSignData.value = undefined
}

onMounted(()=>{
    uni.$once('userLogin', () => {
        loginStatus.value = true 
        checkCurrentSignData(uni.getStorageSync("tokenStr"))
    } )
    loginStatus.value = CheckLoginStatus(uni.getStorageSync("tokenStr"))
    if(loginStatus.value) {
        checkCurrentSignData(uni.getStorageSync("tokenStr"))
    }
})
</script>

<style scoped></style>