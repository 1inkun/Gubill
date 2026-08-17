<template>
	<div class="login">login</div>
	<view v-if="mode">
		<Login @register="() => { mode = false }" />
	</view>
	<view v-else>
		<Register @success="() => { mode = true }" @login="() => { mode = !mode }" />
	</view>
</template>

<script setup lang="ts">
import Login from '@/components/login.vue';
import Register from '@/components/register.vue';
import { ref } from 'vue';
const mode = ref(true)
</script>

<style scoped>
/* @import url("src/static/styles/login.css"); */
</style>
