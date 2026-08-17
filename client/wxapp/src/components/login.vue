<template>
    <view>
        <form @submit="Login">
            <input type="text" name="username" placeholder="请输入用户名">
            <input type="text" name="password" placeholder="请输入密码" :password="true">
            <button form-type="submit">登录</button>
        </form>
    </view>
</template>

<script setup lang="ts">
import NewInstance from '@/api/instance';
import { LoginData, LoginRes } from '@/types/components';
import CONFIG from '@/static/config.json'
import { ResponseData, UserInfo, CustomJWTClaims } from '@/types/global';
import { jwtDecode } from "jwt-decode";

const instance = NewInstance(CONFIG.Server.baseUrl,"")
const checkPasswd = /^(?=.*[a-z])(?=.*\d).{6,30}$/
const checkUsername = /^[A-Za-z0-9]{4,20}$/

const checkLoginInfo = function(loginData: LoginData):boolean{
    if (!checkPasswd.test(loginData.password)){
        console.error("密码有误")
        return false
    }
    if (!checkUsername.test(loginData.username)){
        console.error("用户名有误")
        return false
    }
    return true
}

const SaveUserData = function(tokenStr: string):boolean {
    if(tokenStr == ""){
        return false
    }
    uni.setStorageSync("tokenStr",tokenStr)
    const decode = jwtDecode<CustomJWTClaims>(tokenStr)
    // console.log(decode)
    if (decode.userId == undefined) {
        return false
    }
    let userInfo:UserInfo = {
        username: decode.username,
        userId: decode.userId,
        nickname: decode.nickname,
        role: decode.role
    }
    uni.setStorageSync("userInfo",userInfo)
    return true
}


const Login = async function(e:any){
    uni.showLoading({
        title: '加载中',
        mask: true
    })
    const loginData:LoginData = e.detail.value
    let ok = checkLoginInfo(loginData)
    if (!ok) {
        uni.hideLoading()
        return
    }
    try {
        const resp = await instance<ResponseData<LoginRes>,LoginData>({
            url: "/user/login",
            method: "POST",
            data: loginData
        })
        const res = resp.data
        if(res?.status != "success") {
            throw new Error(res?.msg)
        }
        const data = res.data
        if(!Array.isArray(data)){
            let ok = SaveUserData(data.token)
            if (!ok) {
                throw new Error("登录失败")
            }
            uni.switchTab({url: "index"})
        } 
    } catch (error:any) {
        console.error(error)
    } finally {
        uni.hideLoading()
    }
}
</script>

<style>
/* @import url("src/static/styles/login.css"); */
</style>