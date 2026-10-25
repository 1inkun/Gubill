<script setup lang="ts">
import { onLaunch, onShow } from '@dcloudio/uni-app'
import { setupApiClient } from '@/api/setup'

onLaunch(() => {
  setupApiClient()
  console.log('[app] launched, theme =', (uni.getSystemInfoSync() as any).theme)
})

onShow(() => {
  console.log('[app] show')
})
</script>

<style>
@import './styles/theme.css';
</style>
