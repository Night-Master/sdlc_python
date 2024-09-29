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
          <ElButton type="danger" style="float: right; margin-right: 10px" @click="clearComments()"> Clear Comments </ElButton>
        </div>
        <div class="edit-container h-60vh" v-if="value === 0">
          <pre><code >{{ content }} </code></pre>
        </div>
        <div class="view-container h-60vh" v-else>
          <PreviewComponent  :initialContent="initialContent" />
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
          <ElButton type="danger" style="float: right; margin-right: 10px" @click="clearComments()"> Clear Comments </ElButton>
        </div>
        <div class="edit-container h-60vh" v-if="value2 === 0">
          <pre><code >{{ content2 }} </code></pre>
        </div>
        <div class="view-container h-60vh" v-else>
          <PreviewComponentII :initialContent="initialContent2"/>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue';
import { ElSelect, ElOption, ElButton } from 'element-plus';
import PreviewComponent from './components/PreviewComponent/index.vue';
import PreviewComponentII from './components/PreviewComponent2/index.vue';
import { useRouter } from 'vue-router';
import request from '@/utils/request';

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
const initialContent = ref('');
const initialContent2 = ref('');

// 使用封装的 request 方法来获取 API 代码
const fetchCode = async (apiName, contentRef) => {
  try {
    const response = await request({
      url: '/get_api_code',
      method: 'post',
      data: { api_name: apiName }
    });

    if (response.status === 1) {
      contentRef.value = response.data.code;
    } else {
      alert(`获取代码失败: ${response.message}`);
    }
  } catch (error) {
    alert(`请求失败: ${error.message}`);
  }
};

// 调用后端 API 获取对应的代码
onMounted(() => {
  fetchCode('create_comments', content);
  fetchCode('create_comments_safe', content2);
});

const run1 = () => {
  initialContent.value = `<svg onmouseover="alert('你被xss攻击了')"/>`;
  value.value = 1;
};

const run2 = () => {
  initialContent2.value = `<svg onmouseover="alert('你被xss攻击了')"/>`;
  value2.value = 1;
};

const clearComments = async () => {
  try {
    const response = await request({
      url: '/clear_comments',
      method: 'post'
    });
    if (response.status === 1) {
      alert(response.message);
      location.reload(); // 刷新页面
    } else {
      alert(response.message);
    }
  } catch (error) {
    alert('An error occurred');
  }
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