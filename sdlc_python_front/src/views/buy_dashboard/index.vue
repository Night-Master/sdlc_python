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
          <pre><code>{{ content }}</code></pre>
        </div>
        <div class="view-container h-60vh" v-else>
          <PreviewComponent :initialUsername="username" :initialPassword="password" />
        </div>
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
          <pre><code>{{ content2 }}</code></pre>
        </div>
        <div class="view-container h-60vh" v-else>
          <PreviewComponentII :initialUsername="username2" :initialPassword="password2"/>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue';
import { ElSelect, ElOption, ElButton, ElMessageBox } from 'element-plus';
import PreviewComponent from './components/PreviewComponent/index.vue';
import PreviewComponentII from './components/PreviewComponent2/index.vue';
import { useRouter } from 'vue-router';
import request from '@/utils/request'; // 使用封装的 request 方法

const { currentRoute } = useRouter();
const title = ref(currentRoute.value.meta.title);

watch(
  () => currentRoute.value.meta.title,
  (newVal) => {
    title.value = newVal;
  }
);

const value = ref(0);
const options = ref([
  { value: 0, label: '漏洞代码' },
  { value: 1, label: '应用界面' }
]);

const value2 = ref(0);
const options2 = ref([
  { value: 0, label: '安全代码' },
  { value: 1, label: '应用界面' }
]);

// 动态加载代码
const content = ref('');
const content2 = ref('');
const username = ref('admin');
const password = ref('');
const username2 = ref('admin');
const password2 = ref('');

// 使用封装的 request 方法来获取 API 代码
const fetchCode = async (apiName, contentRef) => {
  try {
    const response = await request({
      url: '/get_api_code', // 根据你的API端点调整
      method: 'post',
      data: { api_name: apiName }
    });

    if (response.status === 1) {
      contentRef.value = response.data.code;
    } else {
      ElMessageBox.alert(`获取代码失败: ${response.message}`, '错误', {
        confirmButtonText: '确定',
        type: 'error',
        center: true
      });
    }
  } catch (error) {
    ElMessageBox.alert(`请求失败: ${error.message}`, '错误', {
      confirmButtonText: '确定',
      type: 'error',
      center: true
    });
  }
};

// 调用后端 API 获取对应的代码
onMounted(() => {
  fetchCode('purchase', content);
  fetchCode('purchase_safe', content2);
});

// 运行按钮的行为
const showReminder = () => {
  ElMessageBox.alert('购买负数商品试试', '关于如何利用', {
    confirmButtonText: '确定',
    type: 'warning',
    dangerouslyUseHTMLString: true,
    center: true,
    customClass: 'custom-message-box'
  });
};

const run1 = () => {
  value.value = 1;
  showReminder();
};

const run2 = () => {
  value2.value = 1;
  showReminder();
};
</script>

<style scoped lang="less">
.homeBox {
  display: flex;
  flex-direction: column;
  height: 100vh;
  padding: 20px;
  box-sizing: border-box;
}
.list {
  display: flex;
  flex: 1;
  .item {
    width: 50%;
    display: flex;
    flex-direction: column;
    .edit-container, .view-container {
      flex: 1;
      margin: 10px;
    }
  }
}
._top {
  margin: 20px;
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
