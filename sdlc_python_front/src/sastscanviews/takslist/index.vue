<template>
  <div class="task-list-page">
    <div class="header">
      <h1>所有扫描任务</h1>
    </div>
    <div class="content">
      <table class="task-table">
        <thead>
          <tr>
            <th>任务ID</th>
            <th>状态</th>
            <th>开始时间</th>
            <th>完成时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="task in tasks" :key="task.taskID">
            <td>{{ task.taskID }}</td>
            <td>
              <span v-if="task.status === 'completed'" class="status-completed">已完成</span>
              <span v-else class="status-scanning">扫描中</span>
            </td>
            <td>{{ formatTime(task.startTime) }}</td>
            <td>{{ task.endTime ? formatTime(task.endTime) : '-' }}</td>
            <td>
              <button
                v-if="task.status === 'completed'"
                @click="viewResult(task.taskID)"
                class="view-button"
              >
                查看结果
              </button>
              <span v-else class="in-progress">等待完成</span>
            </td>
          </tr>
        </tbody>
      </table>
      <div v-if="loading" class="loading">加载中...</div>
      <div v-if="error" class="error">{{ error }}</div>
      <div v-if="tasks.length === 0 && !loading" class="no-tasks">
        <p>暂无扫描任务。</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import request from '@/utils/request'; // 使用你封装的 request 工具
import { useRouter } from 'vue-router';

const tasks = ref([]);
const loading = ref(false);
const error = ref('');
const router = useRouter();

// 获取所有任务
const fetchTasks = async () => {
  loading.value = true;
  error.value = '';
  try {
    const response = await request({
      url: '/tasks',
      method: 'get',
    });

    if (response.status === 1 && response.tasks) {
      tasks.value = response.tasks;
    } else {
      error.value = response.message || '无法获取任务列表';
    }
  } catch (err) {
    error.value = '无法获取任务列表，请重试';
  } finally {
    loading.value = false;
  }
};

// 跳转到任务结果页面
const viewResult = (taskID) => {
  router.push({ name: 'TaskResult', params: { taskID } });
};

// 格式化时间
const formatTime = (time) => {
  if (!time || time === '0001-01-01T00:00:00Z') return '-';
  return new Date(time).toLocaleString();
};

onMounted(() => {
  fetchTasks();
});
</script>

<style scoped>
.task-list-page {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
  font-family: 'Helvetica Neue', Arial, sans-serif;
  background-color: #f5f7fa;
  min-height: 100vh;
}

.header {
  text-align: center;
  margin-bottom: 20px;
  color: #333;
  font-size: 28px;
  font-weight: bold;
}

.content {
  width: 100%;
  max-width: 1000px;
  background: #fff;
  padding: 30px;
  border-radius: 12px;
  box-shadow: 0 6px 25px rgba(0, 0, 0, 0.1);
}

.task-table {
  width: 100%;
  border-collapse: collapse;
}

.task-table th,
.task-table td {
  padding: 16px 20px;
  border-bottom: 1px solid #e0e0e0;
  text-align: left;
}

.task-table th {
  background-color: #fafafa;
  font-weight: 600;
}

.task-table tr:hover {
  background-color: #f1f1f1;
}

.status-completed {
  color: #28a745;
  font-weight: bold;
}

.status-scanning {
  color: #ffc107;
  font-weight: bold;
}

.view-button {
  padding: 8px 16px;
  border: none;
  border-radius: 6px;
  background-color: #007bff;
  color: white;
  cursor: pointer;
  transition: background-color 0.3s;
}

.view-button:hover {
  background-color: #0056b3;
}

.in-progress {
  color: #6c757d;
  font-style: italic;
}

.loading,
.error,
.no-tasks {
  margin-top: 20px;
  text-align: center;
  color: #555;
  font-size: 16px;
}
</style>
