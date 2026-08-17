<template>
	<view class="user">user</view>
    <button @click="Login">登录</button>
</template>

<script setup lang="ts">

const Login = function(){
    uni.navigateTo({url: "login"})
}
</script>

<style scoped>
/* @import url("src/static/styles/login.css"); */
</style>
