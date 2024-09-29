<script setup lang="ts">
import { useRouter } from 'vue-router'
import { ref } from 'vue'
const isHasMessage = ref(true)
const { push } = useRouter()
const goMessgae = () => {
  console.log(1233)
  // 去消息页面
  push('/system-setting/message')
}
</script>

<template>
  <div class="hover-trigger Icon-messgae" @click="goMessgae()">
    <Icon
      :style="{ color: 'var(--top-header-text-color)' }"
      icon="ant-design:message-outlined"
    />
    <span v-if="isHasMessage" class="message"></span>
  </div>
</template>
<style lang="less" scoped>
.Icon-messgae {
  height: 100%;
  position: relative;
}
.message {
  display: inline-block;
  width: 5px;
  height: 5px;
  background: red;
  border-radius: 50%;
  position: absolute;
  top: 50%;
  margin-top: -10px;
  right: 8px;
}
</style>