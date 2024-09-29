<template>
  <div class="title">{{ title }}</div>
</template>
<script setup>
import { ref, onMounted } from 'vue'
import { ElCard } from 'element-plus'
import moment from 'moment'
const props = defineProps({
  title: {
    type: String,
    default: ''
  }
})

onMounted(() => {
  console.log('mounted', props.title)
})
</script>
<style scoped lang="less">
.title{
    margin-bottom: 15px;
    padding-bottom: 15px;
    font-size: 14px;
    position: relative;
    // border-bottom: 1px solid #ebeef5;
    padding-left: 20px;
}
.title::before{
    content: '';
    position: absolute;
    width: 5px;
    height: 20px;
    background:  var(--el-color-primary);
    left: 5px;
}
</style>
