<template>
    <view>
        <form @submit="Register">
            <input type="text" placeholder="输入账户名">
            <input type="text" placeholder="输入昵称">
            <input type="text" placeholder="输入邮箱">
            <input type="text" placeholder="输入密码" :password="true">
            <input type="text" placeholder="确认密码" :password="true">
            <button form-type="submit"></button>
        </form>
    </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'


const password = ref("")
const checkPassword = ref("")
const computed = (()=>{
    return password.value == checkPassword.value ? true : false
})
const Register = async function(e:any) {
    uni.showLoading({
        title: '加载中',
        mask: true
    })
    console.log(e.detail)
    try {
        
    } catch (error:any) {
        console.error(error)
    } finally {
        uni.hideLoading()
    }
}
</script>

<style scoped></style>