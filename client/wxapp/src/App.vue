<script setup lang="ts">
import { onLaunch, onShow, onHide } from "@dcloudio/uni-app";
onLaunch(() => {
	console.log("App Launch");
});
onShow(() => {
	console.log("App Show");
});
onHide(() => {
	console.log("App Hide");
});
</script>
<style>
:root {
	--background-color: #f5f5f5;
	--card-background-color: #ffffff;
	--primary-color: #8acbec;
	--primary-color-hover: #95cdff;
	--error-color: #fce4ec;
	--error-color-hover: #f8bbd0;
	--success-color: #e0f2f1;
	--success-color-hover: #b2dfdb;
}

page {
	background-color: var(--background-color);
}

.content {
	display: flex;
	flex-direction: column;
	justify-content: center;
	margin: 0 16px;
}

.card-container {
	display: flex;
	flex-direction: column;
	justify-content: center;
	padding: 16px;
	border-radius: 24px;
	background-color: var(--card-background-color);
}

.header {
	font-size: 1.25rem;
	padding-left: 8px;
	border-left: 8px solid var(--primary-color);
	margin-top: 8px;
	margin-bottom: 16px;
}

.secondary-card {
	display: flex;
	flex-direction: column;
	justify-content: center;
	padding: 16px;
	border-radius: 16px;
	margin-bottom: 16px;
	background-color: var(--background-color);
}

.item {
	display: inline-flex;
	flex-direction: row;
	justify-content: space-between;
	margin-bottom: 16px;
}

.btn {
	width: 100%;
	padding: 8px 0;
	font-weight: 600;
	font-size: 16px;
	border: none;
	border-radius: 12px;
	transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
	letter-spacing: 0.3px;
	box-shadow: 0 8px 24px rgba(15, 23, 42, 0.15);
}

.btn-primary {
	background: var(--primary-color);
	color: #ffffff;
}

.btn-primary-hover {
	background: var(--primary-color-hover);
	transform: translateY(-2px);
	box-shadow: 0 12px 32px rgba(15, 23, 42, 0.2);
}

.btn-error {
	background-color: var(--error-color);
}

.btn-error-hover {
	background: var(--error-color-hover);
	transform: translateY(-2px);
	box-shadow: 0 12px 32px rgba(15, 23, 42, 0.2);
}
</style>
