<template>
  <div class="homeBox">
    <div class="_title_t">{{ title }}</div>
    <div class="tips" v-if="currentRoute.meta.desc">
      {{ currentRoute.meta.desc }}
    </div>
    <div class="list">
      <div class="item">
        <div class="_top">
          <ElSelect v-model="value" placeholder="" style="width: 140px">
            <ElOption
              v-for="item in options"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </ElSelect>
          <ElButton type="primary" style="float: right" @click="run1()"> Run </ElButton>
        </div>
        <div class="edit-container h-60vh" v-if="value === 0">
          <!-- <CodeEditor ref="MonacoEditRef" v-model="content" language="go" theme="vs" /> -->
          <pre><code >{{ content }} </code></pre>
        </div>
        <div class="view-container h-60vh" v-else><PreviewComponent /></div>
      </div>
      <div class="item">
        <div class="_top">
          <ElSelect v-model="value2" placeholder="" style="width: 140px">
            <ElOption
              v-for="item in options2"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </ElSelect>
          <ElButton type="primary" style="float: right" @click="run2()"> Run </ElButton>
        </div>
        <div class="edit-container h-60vh" v-if="value2 === 0">
          <!-- <CodeEditor ref="MonacoEditRef2" v-model="content2" language="go" theme="vs" /> -->
          <pre><code >{{ content2 }} </code></pre>
        </div>
        <div class="view-container h-60vh" v-else><PreviewComponentII /></div>
      </div>
    </div>
  </div>
</template>
<script setup>
import { CodeEditor, } from '@/components/CodeEditor'
import PreviewComponent from './components/PreviewComponent/index.vue'
import PreviewComponentII from './components/PreviewComponent2/index.vue'
import { ref, onMounted,watch } from 'vue'
import { ElSelect, ElOption, ElButton } from 'element-plus'
import { useRouter } from 'vue-router'
const { currentRoute } = useRouter()
console.log(currentRoute,3333)
const title = ref(currentRoute.value.meta.title)
// 监听录音变化
watch(
  () => currentRoute.value.meta.title,
  (newVal, oldVal) => {
    title.value = newVal
  }
)
const value = ref(0)
const options = ref([
  {
    value: 0,
    label: '漏洞代码'
  },
  {
    value: 1,
    label: '漏洞界面'
  }
])
const value2 = ref(0)
const options2 = ref([
  {
    value: 0,
    label: '漏洞代码'
  },
  {
    value: 1,
    label: '漏洞界面'
  }
])
const content = ref(
  'public class HelloWorld {\n  public static void main(String[] args) {\n    System.out.println("Hello, World!");\n  }\n}'
)
const content2 = ref(
  `<div class="sss">
        <div class="ddddd">
          <CodeEditor ref="MonacoEditRef222" v-model="content33" />
        </div>
</div>`
)
const run1 = () => {
  value.value = 1
}
const run2 = () => {
  value2.value = 1
}
onMounted(() => {
  console.log('mounted')
})
</script>
<style scoped lang="less">
.icon {
  font-size: 48px;
  margin-right: 10px;
}
.title {
  display: flex;
  align-items: center;
  font-size: 24px;
}
.box {
  margin: 30px;
  font-size: 18px;
  p {
    margin: 6px 0;
  }
}
.margin-right {
  margin-right: 10px;
}
.list {
  display: flex;
  .item {
    width: 50%;
    .edit-container {
      margin: 10px;
    }
  }
}
._top {
  margin: 20px;
}
.view-container {
  // border: 1px solid #f1f1f1;
  box-shadow: 0 0 5px #f1f1f1;
  margin: 10px;
}
.tips {
  margin: 20px 10px;
  padding: 20px;
  background-color: #fff;
  border-radius: 10px;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
  line-height: 1.6;
  color: #555;
  cursor: pointer;
  transition: background-color 0.3s, color 0.3s;
}
.tips:hover {
  background-color: #e9ecef;
}
pre {
  background-color: #fff;
  color: #333;
  padding: 20px;
  border-radius: 10px;
  overflow: auto;
  flex: 1;
  height: 100%;
  width: 100%;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
}
</style>
<style>
.homeBox {
  .el-input__wrapper {
    background: inherit !important;
    border: none !important;
    box-shadow: none !important;
  }
  .el-input__wrapper:hover {
    background: inherit;
    border: none;
    box-shadow: none;
  }
  .el-select .el-input.is-focus .el-input__wrapper {
    border: none !important;
    box-shadow: none !important;
  }
  .el-select .is-focus.el-input__wrapper {
    border: none !important;
    box-shadow: none !important;
  }
  .el-select .el-input__inner {
    font-weight: bold;
    font-size: 20px;
  }
}
</style>
