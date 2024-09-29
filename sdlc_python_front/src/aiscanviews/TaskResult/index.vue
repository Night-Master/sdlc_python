<template>
  <div class="task-result-page">
    <div class="container">
      <h1>扫描结果 - 任务ID: {{ taskID }}</h1>
      <button class="export-btn" @click="exportData">导出结果</button>

      <table v-if="result.length > 0">
        <thead>
          <tr>
            <th class="file-path">文件路径</th>
            <th>函数名</th>
            <th>行号</th>
            <th>消息</th>
            <th>严重性</th>
            <th>CWE</th>
            <th>片段</th> <!-- 更改列标题 -->
          </tr>
        </thead>
        <tbody>
          <tr v-for="(vuln, index) in result" :key="index">
            <td class="file-path">{{ shortenPath(vuln.file_path) }}</td>
            <td>{{ vuln.func_name }}</td>
            <td>{{ vuln.line }}</td>
            <td>{{ vuln.message }}</td>
            <td>
              <span :class="getSeverityClass(vuln.severity)">
                {{ vuln.severity }}
              </span>
            </td>
            <td>{{ vuln.cwe }}</td>
            <td>
              <button class="view-btn" @click="viewCodeSnippet(vuln)">查看</button> <!-- 更改按钮标签 -->
            </td>
          </tr>
        </tbody>
      </table>

      <div v-else-if="!loading && !error" class="no-results">
        <p>没有找到扫描结果。</p>
      </div>

      <div v-if="loading" class="loading">加载中...</div>
      <div v-if="error" class="error">{{ error }}</div>

      <!-- 弹窗 -->
      <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
        <div class="modal-content">
          <h2>代码片段</h2>
          <pre>{{ codeSnippet }}</pre>
          <button class="close-btn" @click="closeModal">关闭</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import request from '@/utils/request_ai_code'; // 使用你封装的 request 工具
import { useRoute } from 'vue-router';

const route = useRoute();
const taskID = route.params.taskID;

const result = ref([]);
const loading = ref(false);
const error = ref('');
const showModal = ref(false);
const codeSnippet = ref('');

// 获取扫描结果
const fetchResult = async () => {
  loading.value = true;
  error.value = '';
  try {
    const response = await request({
      url: `/api/audit/result/${taskID}`,
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

  let csvContent = '\uFEFF'; // 添加 BOM 解决中文乱码
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

// 查看代码片段
const viewCodeSnippet = (vuln) => {
  if (vuln.code_snippet) {
    codeSnippet.value = vuln.code_snippet;
    showModal.value = true;
  } else {
    alert('没有代码片段可显示');
  }
};

// 关闭弹窗
const closeModal = () => {
  showModal.value = false;
};

// 短化文件路径以适应表格布局
const shortenPath = (path) => {
  const maxLength = 40; // 根据需要调整
  if (path.length > maxLength) {
    return '...' + path.slice(-maxLength);
  }
  return path;
};

// 获取严重性对应的样式类
const getSeverityClass = (severity) => {
  const lower = severity.toLowerCase();
  switch (lower) {
    case 'critical':
      return 'severity-critical';
    case 'high':
      return 'severity-high';
    case 'medium':
      return 'severity-medium';
    case 'low':
      return 'severity-low';
    default:
      return 'severity-unknown';
  }
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
  align-items: flex-start;
  padding: 20px;
  font-family: 'Helvetica Neue', Arial, sans-serif;
  background-color: #f5f7fa;
  min-height: 100vh;
}

.container {
  max-width: 1200px; /* 根据需要调整 */
  width: 100%;
  background-color: white;
  padding: 30px;
  border-radius: 12px;
  box-shadow: 0 6px 25px rgba(0, 0, 0, 0.1);
  position: relative;
}

h1 {
  text-align: center;
  font-size: 28px;
  margin-bottom: 20px;
  color: #333;
}

.export-btn {
  background-color: #28a745;
  color: white;
  padding: 10px 20px;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  margin-bottom: 20px;
  position: absolute;
  top: 30px;
  right: 30px;
  transition: background-color 0.3s, transform 0.2s;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.export-btn:hover {
  background-color: #218838;
  transform: translateY(-2px);
}

table {
  width: 100%;
  border-collapse: collapse;
  table-layout: fixed; /* 固定表格布局 */
}

th, td {
  padding: 12px;
  text-align: left;
  word-wrap: break-word;
}

th {
  background-color: #007bff;
  color: white;
  border: none;
}

.file-path {
  width: 20%; /* 缩窄文件路径列 */
}

th, td {
  border: 1px solid #ddd;
}

tr:nth-child(even) {
  background-color: #f9f9f9;
}

tr:hover {
  background-color: #f1f1f1;
}

/* 严重性样式 */
.severity-critical {
  background-color: #dc3545;
  color: white;
  padding: 4px 8px;
  border-radius: 4px;
  font-weight: bold;
}

.severity-high {
  background-color: #fd7e14;
  color: white;
  padding: 4px 8px;
  border-radius: 4px;
  font-weight: bold;
}

.severity-medium {
  background-color: #ffc107;
  color: #212529;
  padding: 4px 8px;
  border-radius: 4px;
  font-weight: bold;
}

.severity-low {
  background-color: #28a745;
  color: white;
  padding: 4px 8px;
  border-radius: 4px;
  font-weight: bold;
}

.severity-unknown {
  background-color: #6c757d;
  color: white;
  padding: 4px 8px;
  border-radius: 4px;
  font-weight: bold;
}

/* 查看按钮样式 */
.view-btn {
  background-color: #17a2b8;
  color: white;
  padding: 6px 12px;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  transition: background-color 0.3s, transform 0.2s;
  font-size: 14px;
}

.view-btn:hover {
  background-color: #138496;
  transform: translateY(-2px);
}

/* 弹窗样式 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background-color: #fff;
  padding: 20px;
  border-radius: 8px;
  width: 90%;
  max-width: 800px;
  max-height: 80vh;
  overflow-y: auto;
  position: relative;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
  animation: fadeIn 0.3s ease-out;
}

.modal-content h2 {
  margin-top: 0;
  color: #333;
}

.modal-content pre {
  background-color: #f8f9fa;
  padding: 15px;
  border-radius: 5px;
  overflow-x: auto;
  white-space: pre-wrap;
  font-family: 'Courier New', Courier, monospace;
  font-size: 14px;
}

.close-btn {
  background-color: #dc3545;
  color: white;
  padding: 8px 16px;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  margin-top: 15px;
  transition: background-color 0.3s, transform 0.2s;
}

.close-btn:hover {
  background-color: #c82333;
  transform: translateY(-2px);
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

/* 动画 */
@keyframes fadeIn {
  from {
    opacity: 0;
    transform: scale(0.95);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .container {
    padding: 20px;
  }

  th, td {
    padding: 10px;
    font-size: 14px;
  }

  h1 {
    font-size: 22px;
  }

  .export-btn {
    padding: 8px 16px;
    top: 20px;
    right: 20px;
  }

  .view-btn {
    padding: 5px 10px;
    font-size: 13px;
  }

  .modal-content {
    width: 95%;
    max-width: 90%;
  }
}
</style>
