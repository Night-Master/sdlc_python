<template>
  <div class="task-list-page">
    <div class="container">
      <h1>扫描任务列表</h1>
      <table v-if="tasks.length > 0">
        <thead>
          <tr>
            <th>任务ID</th>
            <th>状态</th>
            <th>进度</th>
            <th>扫描参数</th>
            <th>开始时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="task in tasks" :key="task.taskID">
            <td>{{ task.taskID }}</td>
            <td>{{ task.status }}</td>
            <td>
              <div v-if="task.status === 'scanning'" class="progress-bar">
                <div
                  class="progress"
                  :style="{ width: calculateProgress(task.taskID) }"
                ></div>
                <span>{{ progressMap[task.taskID] || 0 }}/{{ totalFilesMap[task.taskID] || 1 }}</span>
              </div>
              <span v-else>{{ task.status === 'completed' ? '已完成' : '出错' }}</span>
            </td>
            <td>
              <div class="scan-params">
                <p><strong>语言：</strong>{{ task.language }}</p>
                <p><strong>线程数：</strong>{{ task.threadNum }}</p>
                <p><strong>AI平台：</strong>{{ task.aiPlatform }}</p>
                <p><strong>Web框架：</strong>{{ task.webFramework || '无' }}</p>
              </div>
            </td>
            <td>{{ formatDate(task.startTime) }}</td>
            <td>
              <router-link :to="`/result/${task.taskID}`" class="view-btn">查看结果</router-link>
            </td>
          </tr>
        </tbody>
      </table>

      <div v-else class="no-tasks">
        <p>暂无扫描任务。</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import request from '@/utils/request_ai_code'; // 使用您封装的 request 工具
import { useIntervalFn } from '@vueuse/core';

// 定义响应式变量
const tasks = ref([]);
const progressMap = ref({});
const totalFilesMap = ref({});

// 格式化日期的方法
const formatDate = (dateString) => {
  if (!dateString) return '-';
  const options = {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  };
  return new Date(dateString).toLocaleDateString(undefined, options);
};

// 计算进度百分比的方法
const calculateProgress = (taskID) => {
  const scanned = progressMap.value[taskID] || 0;
  const total = totalFilesMap.value[taskID] || 1; // 避免除以0
  const percentage = Math.min((scanned / total) * 100, 100).toFixed(2) + '%';
  return percentage;
};

// 获取任务列表的函数
const fetchTasks = async () => {
  try {
    const response = await request({
      url: '/api/audit/tasks',
      method: 'get',
    });

    if (response.status === 1 && response.tasks) {
      tasks.value = response.tasks;
      // 初始化进度和总文件数
      response.tasks.forEach(task => {
        if (task.status === 'scanning') {
          if (!(task.taskID in progressMap.value)) {
            progressMap.value[task.taskID] = 0;
          }
          if (!(task.taskID in totalFilesMap.value)) {
            totalFilesMap.value[task.taskID] = 1; // 默认总数为1，避免除以0
          }
        }
      });
    } else {
      console.error(response.message || '无法获取任务列表');
    }
  } catch (error) {
    console.error('获取任务列表失败:', error);
  }
};

// 获取单个任务进度的函数
const fetchTaskProgress = async (taskID) => {
  try {
    const response = await request({
      url: `/api/audit/progress/${taskID}`,
      method: 'get',
    });

    if (response.status === 1 && response.scanned !== undefined && response.total !== undefined) {
      progressMap.value[taskID] = response.scanned;
      totalFilesMap.value[taskID] = response.total;
    } else {
      console.error(`任务 ${taskID} 的进度获取失败: ${response.message || '未知错误'}`);
    }
  } catch (error) {
    console.error(`获取任务 ${taskID} 进度失败:`, error);
  }
};

// 获取所有扫描任务的进度
const fetchAllProgress = async () => {
  const scanningTasks = tasks.value.filter(task => task.status === 'scanning');

  // 并行请求所有扫描任务的进度
  const progressPromises = scanningTasks.map(task =>
    fetchTaskProgress(task.taskID)
  );

  await Promise.all(progressPromises);
};

// 获取任务列表和进度的函数
const fetchTasksAndProgress = async () => {
  await fetchTasks();
  await fetchAllProgress();
};

// 页面加载时获取任务列表和进度
onMounted(async () => {
  await fetchTasksAndProgress();
});

// 每隔5秒刷新一次任务列表和进度
useIntervalFn(async () => {
  await fetchTasksAndProgress();
}, 5000);
</script>

<style scoped>
.task-list-page {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 20px;
  font-family: 'Helvetica Neue', Arial, sans-serif;
  background: linear-gradient(135deg, #f5f7fa, #c3cfe2);
  min-height: 100vh;
  transition: background 0.5s ease;
}

.container {
  width: 100%;
  max-width: 1400px; /* 增加宽度以适应更多列 */
  background-color: #ffffffcc;
  padding: 40px;
  border-radius: 15px;
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.12);
  backdrop-filter: blur(10px);
  transition: box-shadow 0.3s ease;
  animation: fadeInUp 1s ease;
}

.container:hover {
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.2);
}

h1 {
  text-align: center;
  font-size: 36px;
  margin-bottom: 30px;
  color: #333;
  animation: fadeInDown 1s ease;
}

.export-btn {
  background-color: #28a745;
  color: white;
  padding: 10px 15px;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  margin-bottom: 20px;
  position: absolute;
  top: 30px;
  right: 30px;
  transition: background-color 0.3s;
}

.export-btn:hover {
  background-color: #218838;
}

table {
  width: 100%;
  border-collapse: collapse;
  animation: fadeIn 1s ease;
}

table, th, td {
  border: 1px solid #ddd;
}

th, td {
  padding: 12px;
  text-align: left;
  vertical-align: middle;
}

th {
  background-color: #007bff;
  color: white;
  border: none;
  position: sticky;
  top: 0;
  z-index: 1;
}

tr:nth-child(even) {
  background-color: #f9f9f9;
}

tr:hover {
  background-color: #f1f1f1;
}

.scan-params {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.scan-params p {
  margin: 0;
  font-size: 14px;
  color: #555;
}

.progress-bar {
  position: relative;
  width: 100%;
  height: 20px;
  background-color: #e9ecef;
  border-radius: 10px;
  overflow: hidden;
}

.progress {
  height: 100%;
  background-color: #007bff;
  transition: width 0.5s ease;
}

.progress-bar span {
  position: absolute;
  top: 0;
  left: 50%;
  transform: translateX(-50%);
  font-size: 12px;
  color: white;
}

.view-btn {
  background-color: #007bff;
  color: white;
  padding: 6px 12px;
  border: none;
  border-radius: 5px;
  text-decoration: none;
  cursor: pointer;
  transition: background-color 0.3s;
}

.view-btn:hover {
  background-color: #0056b3;
}

.no-tasks {
  text-align: center;
  color: #555;
  font-size: 18px;
  animation: fadeIn 1s ease;
}

/* 响应式设计 */
@media (max-width: 1400px) {
  .container {
    max-width: 1000px;
  }
}

@media (max-width: 1200px) {
  .container {
    padding: 30px;
  }

  th, td {
    padding: 10px;
    font-size: 14px;
  }

  h1 {
    font-size: 28px;
    margin-bottom: 20px;
  }
}

@media (max-width: 768px) {
  .container {
    padding: 20px;
  }

  table, th, td {
    font-size: 12px;
  }

  .view-btn {
    padding: 4px 8px;
    font-size: 14px;
  }

  .scan-params p {
    font-size: 12px;
  }
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes fadeInDown {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
