<template>
  <div>
    <ul>
      <li
        v-for="(item, index) in list"
        :class="{ active: index === current }"
        :key="index"
        @click="tab(item,index)"
        >{{ item.title }}</li
      >
    </ul>
  </div>
</template>
<script setup >
const projectType =  localStorage.getItem('projectType')  || 'bug_platform';
import { ref } from 'vue'
import {currentProjectHomePage} from '@/utils/index'
import { useAppStore } from '@/store/modules/app'
const appStore = useAppStore()
console.log(appStore.getTitle,'项目获取')
const list = ref([
  {
    project: 'bug_platform',
    title: appStore.getTitle+'漏洞平台',
    path:currentProjectHomePage('bug_platform')
  },
  // {
  //   project: 'sast',
  //   title: '代码sast扫描工具',
  //   path:currentProjectHomePage('sast')
  // },
  {
    project: 'aiscan',
    title: 'AI代码扫描工具',
    path:currentProjectHomePage('sast')
  }
])
const currentIndex = list.value.findIndex(item=>item.project === projectType)
const current = ref(currentIndex)

const tab = (item,index) => {
  console.log(item)
  current.value = index
  localStorage.setItem('projectType',item.project)
  if(item.path){
    //  item.path
    window.location.href =`${window.location.origin}${window.location.pathname}#${item.path}` 
    window.location.reload()
  }
}
</script>
<style lang="less" scoped>
ul {
  display: flex;
  color: #999;
  font-size: 20px;
  align-items: center;
  li {
    margin-right: 15px;
    cursor: pointer;
    position: relative;
    padding-right: 15px;
    &:first-child {
      margin-left: 0px;
    }
    &:after {
      content: '';
      width: 2px;
      height: 20px;
      background-color: #999;
      position: absolute;
      right: 0px;
      top: 50%;
      transform: translateY(-50%);
    }
    &:last-child:after{
      display: none;
    }
  }
  .active {
    color: #007bff;
  }
}
</style>
