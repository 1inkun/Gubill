<template>
	<div class="login">login</div>
    <view >
        <Login/>
    </view>
</template>

<script setup lang="ts">
import Login from '@/components/login.vue';
import CheckLoginStatus from '@/utils/checkLoginStatus';
import { onMounted, ref } from 'vue';
const loginStatus = ref(false)

onMounted(()=>{
    loginStatus.value = CheckLoginStatus(uni.getStorageSync("tokenStr"))
})
</script>

<style scoped>
/* @import url("src/static/styles/login.css"); */
</style>
