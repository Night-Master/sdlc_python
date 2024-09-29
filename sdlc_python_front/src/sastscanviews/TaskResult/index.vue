<template>
  <div class="task-result-page">
    <div class="container">
      <h1>扫描结果 - 任务ID: {{ taskID }}</h1>
      <button class="export-btn" @click="exportData">导出结果</button>

      <table v-if="result.length > 0">
        <thead>
          <tr>
            <th>文件路径</th>
            <th>函数名</th>
            <th>行号</th>
            <th>消息</th>
            <th>严重性</th>
            <th>CWE</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(vuln, index) in result" :key="index">
            <td>{{ vuln.file_path }}</td>
            <td>{{ vuln.func_name }}</td>
            <td>{{ vuln.line }}</td>
            <td>{{ vuln.message }}</td>
            <td>
              <span :class="`severity-${vuln.severity.toLowerCase()}`">
                {{ vuln.severity }}
              </span>
            </td>
            <td>{{ vuln.cwe }}</td>
          </tr>
        </tbody>
      </table>

      <div v-else-if="!loading && !error" class="no-results">
        <p>没有找到扫描结果。</p>
      </div>

      <div v-if="loading" class="loading">加载中...</div>
      <div v-if="error" class="error">{{ error }}</div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import request from '@/utils/request'; // 使用你封装的 request 工具
import { useRoute } from 'vue-router';

const route = useRoute();
const taskID = route.params.taskID;

const result = ref([]);
const loading = ref(false);
const error = ref('');

// 获取扫描结果
const fetchResult = async () => {
  loading.value = true;
  error.value = '';
  try {
    const response = await request({
      url: `/result/${taskID}`,
      method: 'get',
    });

    if (response.status === 1 && response.result) {
      result.value = response.result;
    } else {
      error.value = response.message || '无法获取扫描结果';
    }
  } catch (err) {
    error.value = '无法获取扫描结果，请重试';
  } finally {
    loading.value = false;
  }
};

// 导出扫描结果为 CSV
const exportData = () => {
  if (result.value.length === 0) {
    alert('没有数据可以导出');
    return;
  }

  const headers = ['文件路径', '函数名', '行号', '消息', '严重性', 'CWE'];
  const rows = result.value.map(vuln => [
    vuln.file_path,
    vuln.func_name,
    vuln.line,
    vuln.message,
    vuln.severity,
    vuln.cwe,
  ]);

  let csvContent = '';
  csvContent += headers.join(',') + '\n';
  rows.forEach(row => {
    // Escape double quotes by doubling them
    const escapedRow = row.map(item => `"${String(item).replace(/"/g, '""')}"`);
    csvContent += escapedRow.join(',') + '\n';
  });

  const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
  const url = URL.createObjectURL(blob);

  const link = document.createElement('a');
  link.setAttribute('href', url);
  link.setAttribute('download', `scan_result_${taskID}.csv`);
  link.style.visibility = 'hidden';
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
};

onMounted(() => {
  fetchResult();
});
</script>

<style scoped>
body {
  margin: 0;
  padding: 0;
}

.task-result-page {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 20px;
  font-family: 'Helvetica Neue', Arial, sans-serif;
  background-color: #f5f7fa;
  min-height: 100vh;
}

.container {
  max-width: 1000px;
  width: 100%;
  background-color: white;
  padding: 30px;
  border-radius: 12px;
  box-shadow: 0 6px 25px rgba(0, 0, 0, 0.1);
  position: relative;
}

h1 {
  text-align: center;
  font-size: 24px;
  margin-bottom: 20px;
  color: #333;
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
}

table, th, td {
  border: 1px solid #ddd;
}

th, td {
  padding: 12px;
  text-align: left;
}

th {
  background-color: #007bff;
  color: white;
  border: none;
}

tr:nth-child(even) {
  background-color: #f9f9f9;
}

tr:hover {
  background-color: #f1f1f1;
}

.severity-high {
  color: #dc3545;
  font-weight: bold;
}

.severity-medium {
  color: #ffc107;
  font-weight: bold;
}

.severity-low {
  color: #28a745;
  font-weight: bold;
}

.loading,
.error,
.no-results {
  margin-top: 20px;
  text-align: center;
  color: #555;
  font-size: 16px;
}

.error {
  color: #dc3545;
}
</style>
