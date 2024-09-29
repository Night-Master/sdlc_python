<!-- src/components/Test.vue -->
<template>
  <div class="test-container">
    <h1>{{ projectName }}</h1>
    <div v-if="loading">加载中...</div>
    <div v-else-if="error">{{ error }}</div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import request from '@/utils/request';

const projectName = ref('');
const loading = ref(true);
const error = ref('');

const fetchProjectName = async () => {
  try {
    const response = await request({
      url: '/get_project_name',
      method: 'get',
    });
    projectName.value = response.project_name;
  } catch (err) {
    console.error('获取项目名称失败:', err);
    error.value = '获取项目名称失败';
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  fetchProjectName();
});
</script>

<style scoped>
.test-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
  font-family: Arial, sans-serif;
}

h1 {
  font-size: 2em;
  margin-bottom: 20px;
}

div {
  font-size: 1.2em;
}
</style>
