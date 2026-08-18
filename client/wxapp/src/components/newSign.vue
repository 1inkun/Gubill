<template>
    <button @click="submit">即刻签到</button>
</template>

<script setup lang="ts">
import NewInstance from '@/api/instance';
import CheckLoginStatus from '@/utils/checkLoginStatus';
import CONFIG from '@/static/config.json'
import { ResponseData } from '@/types/global';
import { GenSignRes } from '@/types/components';
const emit = defineEmits(['sign-in'])

const submit = async function() {
    uni.showLoading({
        title: '加载中',
        mask: true
    })
    const loginStatus = CheckLoginStatus(uni.getStorageSync("tokenStr"))
    if(!loginStatus){
        uni.hideLoading()
        console.error("还未登录")
        return
    }
    const instance = NewInstance(CONFIG.Server.baseUrl,uni.getStorageSync("tokenStr"))
    try {
        const resp = await instance<ResponseData<GenSignRes>>({
            url: "/sign",
            method: "POST"
        })
        const res = resp.data
        if (res?.status != "success") {
            throw new Error(res?.msg)
        }
        const data = res.data
        if (!Array.isArray(data)){
            const signId = data.signId
            emit('sign-in',signId)
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