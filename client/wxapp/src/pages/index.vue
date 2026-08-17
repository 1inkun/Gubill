<template>
	<view class="content">
		<view class="text-area">
			<text class="title">{{ title }}</text>
		</view>
	</view>
</template>

<script setup lang="ts">
import { ref } from "vue";
const title = ref("Hello");
</script>

<style>
/* @import url("src/static/styles/index.css"); */
</style>
