<template>
  <div class="task-creation-page">
    <div class="header">
      <h1>创建扫描任务</h1>
    </div>
    <div class="content">
      <div class="mode-selection">
        <label class="radio-option">
          <input type="radio" value="local" v-model="scanMode" /> 本地扫描
        </label>
        <label class="radio-option">
          <input type="radio" value="cloud" v-model="scanMode" /> 云端扫描
        </label>
      </div>

      <!-- 本地模式下输入路径和拖拽上传 -->
      <div v-if="scanMode === 'local'" class="local-upload-section">
        <input
          v-model="scanPath"
          placeholder="输入本地路径"
          class="input-field"
        />
        <p class="or-text">或</p>
        <div
          class="upload-form-container"
          @dragover.prevent
          @drop.prevent="handleDrop"
        >
          <div class="upload-form" @click="triggerFileInput">
            <label for="file-input" class="file-input-label">
              <span v-if="!file">点击或拖拽代码压缩文件到这里</span>
              <span v-else>{{ file.name }}</span>
            </label>
            <input
              ref="fileInput"
              id="file-input"
              type="file"
              @change="handleFileChange"
              class="file-input"
              accept=".zip"
            />
          </div>
        </div>
      </div>

      <!-- 云端模式下仅上传 -->
      <div v-if="scanMode === 'cloud'" class="cloud-upload-section">
        <div
          class="upload-form-container"
          @dragover.prevent
          @drop.prevent="handleDrop"
        >
          <div class="upload-form" @click="triggerFileInput">
            <label for="file-input" class="file-input-label">
              <span v-if="!file">点击或拖拽代码压缩文件到这里</span>
              <span v-else>{{ file.name }}</span>
            </label>
            <input
              ref="fileInput"
              id="file-input"
              type="file"
              @change="handleFileChange"
              class="file-input"
              accept=".zip"
            />
          </div>
        </div>
      </div>

      <!-- 上传进度条 -->
      <div v-if="uploading" class="progress-bar-container">
        <div class="progress-bar" :style="{ width: uploadProgress + '%' }"></div>
        <span class="progress-text">{{ uploadProgress }}%</span>
      </div>

      <!-- 扫描语言选择 -->
      <div class="input-group">
        <label for="language" class="input-label">选择扫描语言</label>
        <select v-model="selectedLanguage" class="select-field" id="language">
          <option disabled value="">请选择语言</option>
          <option value="go">Go</option>
          <option value="java">Java</option>
          <option value="python">Python</option>
          <option value="javascript">JavaScript</option>
        </select>
      </div>

      <!-- 线程数选择 -->
      <div class="input-group">
        <label for="threads" class="input-label">线程数</label>
        <input
          v-model.number="threadCount"
          type="number"
          min="1"
          class="input-field"
          id="threads"
          placeholder="输入线程数"
        />
      </div>

      <!-- AI平台选择 -->
      <div class="input-group">
        <label for="platform" class="input-label">选择AI平台</label>
        <select v-model="selectedAIPlatform" class="select-field" id="platform">
          <option disabled value="">请选择平台</option>
          <option value="deepseek">DeepSeek</option>
          <option value="tongyi_qianwen">同义千问</option>
        </select>
      </div>

      <!-- 新增的 Web 框架选择 -->
      <div class="input-group">
        <label for="webFramework" class="input-label">
          Web 框架（可选）
          <span class="helper-text">(填写可获得更高的准确性)</span>
        </label>
        <select v-model="selectedWebFramework" class="select-field" id="webFramework">
          <option disabled value="">请选择框架（可选）</option>
          <option value="gin">Gin</option>
          <option value="beego">Beego</option>
          <option value="spring">Spring</option>
          <option value="django">Django</option>
          <option value="flask">Flask</option>
          <option value="express">Express</option>
          <option value="其他">其他</option>
        </select>
        <!-- 当选择“其他”时，显示文本输入框 -->
        <input
          v-if="selectedWebFramework === '其他'"
          v-model="customWebFramework"
          type="text"
          placeholder="请输入自定义框架名称"
          class="input-field"
          style="margin-top: 10px;"
        />
      </div>

      <button @click="createTask" class="button">创建任务</button>
      <div class="task-status" v-if="taskStatus">
        <p>{{ taskStatus }}</p>
      </div>
    </div>
  </div>
</template>
<script setup>
import { ref, watch } from 'vue';
import request from '@/utils/request_ai_code'; // 使用你封装的 request 工具

const scanMode = ref('local'); // 模式选择，本地或云端
const scanPath = ref(''); // 本地路径
const file = ref(null); // 上传的文件
const taskStatus = ref(''); // 任务状态
const fileInput = ref(null); // 文件输入引用
const selectedLanguage = ref(''); // 选择的语言
const threadCount = ref(1); // 线程数
const selectedAIPlatform = ref(''); // 选择的AI平台

const selectedWebFramework = ref(''); // 选择的Web框架
const customWebFramework = ref(''); // 自定义Web框架

const uploading = ref(false); // 是否正在上传
const uploadProgress = ref(0); // 上传进度百分比

// 处理文件选择
const handleFileChange = (event) => {
  const selectedFile = event.target.files[0];
  if (selectedFile) {
    file.value = selectedFile;
  }
};

// 处理拖拽文件
const handleDrop = (event) => {
  const droppedFile = event.dataTransfer.files[0];
  if (droppedFile) {
    file.value = droppedFile;
  }
};

// 触发文件输入点击
const triggerFileInput = () => {
  if (fileInput.value) {
    fileInput.value.click();
  }
};

// 监听 selectedWebFramework 的变化，清空自定义框架名称
watch(selectedWebFramework, (newVal) => {
  if (newVal !== '其他') {
    customWebFramework.value = '';
  }
});

// 创建扫描任务
const createTask = async () => {
  taskStatus.value = ''; // 重置任务状态
  uploadProgress.value = 0; // 重置上传进度

  // 验证必填字段
  if (!selectedLanguage.value || !selectedAIPlatform.value || threadCount.value < 1) {
    taskStatus.value = '请填写所有必填字段';
    return;
  }

  // 获取最终的 Web 框架值
  let finalWebFramework = '';
  if (selectedWebFramework.value) {
    if (selectedWebFramework.value === '其他') {
      if (!customWebFramework.value) {
        taskStatus.value = '请填写自定义的 Web 框架名称';
        return;
      }
      finalWebFramework = customWebFramework.value;
    } else {
      finalWebFramework = selectedWebFramework.value;
    }
  }

  // 本地扫描模式
  if (scanMode.value === 'local') {
    if (!scanPath.value && !file.value) {
      taskStatus.value = '请输入本地路径或上传压缩文件';
      return;
    }

    if (scanPath.value) {
      // 通过路径扫描
      try {
        const response = await request({
          url: '/api/audit/scan_start_dir',
          method: 'post',
          data: {
            path: scanPath.value,
            language: selectedLanguage.value,
            thread_num: threadCount.value,
            platform: selectedAIPlatform.value,
            web_framework: finalWebFramework, // 传递 Web 框架参数
          },
        });

        handleResponse(response);
      } catch (error) {
        taskStatus.value = '任务创建失败，请重试';
      }
    } else if (file.value) {
      // 通过上传压缩包扫描
      const formData = new FormData();
      formData.append('file', file.value);
      formData.append('language', selectedLanguage.value);
      formData.append('thread_num', threadCount.value);
      formData.append('platform', selectedAIPlatform.value);
      formData.append('web_framework', finalWebFramework); // 传递 Web 框架参数

      try {
        uploading.value = true; // 开始上传
        const response = await request({
          url: '/api/audit/scan_start_zip',
          method: 'post',
          data: formData,
          headers: {
            'Content-Type': 'multipart/form-data',
          },
          onUploadProgress: (progressEvent) => {
            const percentCompleted = Math.round(
              (progressEvent.loaded * 100) / progressEvent.total
            );
            uploadProgress.value = percentCompleted;
          },
        });

        handleResponse(response);
      } catch (error) {
        taskStatus.value = '任务创建失败，请重试';
      } finally {
        uploading.value = false; // 上传完成或失败
      }
    }
  }

  // 云端扫描模式
  else if (scanMode.value === 'cloud') {
    if (!file.value) {
      taskStatus.value = '请上传压缩文件';
      return;
    }

    const formData = new FormData();
    formData.append('file', file.value);
    formData.append('language', selectedLanguage.value);
    formData.append('thread_num', threadCount.value);
    formData.append('platform', selectedAIPlatform.value);
    formData.append('web_framework', finalWebFramework); // 传递 Web 框架参数

    try {
      uploading.value = true; // 开始上传
      const response = await request({
        url: '/api/audit/scan_start_zip',
        method: 'post',
        data: formData,
        headers: {
          'Content-Type': 'multipart/form-data',
        },
        onUploadProgress: (progressEvent) => {
          const percentCompleted = Math.round(
            (progressEvent.loaded * 100) / progressEvent.total
          );
          uploadProgress.value = percentCompleted;
        },
      });

      handleResponse(response);
    } catch (error) {
      taskStatus.value = '任务创建失败，请重试';
    } finally {
      uploading.value = false; // 上传完成或失败
    }
  }
};

// 处理后端响应
const handleResponse = (response) => {
  if (response.status === 1 && response.taskID) {
    taskStatus.value = `任务创建成功，任务ID: ${response.taskID}`;
  } else {
    taskStatus.value = response.message || '任务创建失败';
  }
};
</script>

<style scoped>
.task-creation-page {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 60px 20px;
  font-family: 'Helvetica Neue', Arial, sans-serif;
  background: linear-gradient(135deg, #f5f7fa, #c3cfe2);
  min-height: 100vh;
  transition: background 0.5s ease;
}

.header {
  text-align: center;
  margin-bottom: 40px;
  color: #333;
  font-size: 36px;
  font-weight: bold;
  animation: fadeInDown 1s ease;
}

.content {
  width: 100%;
  max-width: 700px;
  background: #ffffffcc;
  padding: 40px;
  border-radius: 15px;
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.12);
  backdrop-filter: blur(10px);
  transition: box-shadow 0.3s ease;
  animation: fadeInUp 1s ease;
}

.content:hover {
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.2);
}

.mode-selection {
  display: flex;
  justify-content: center;
  gap: 40px;
  margin-bottom: 30px;
  animation: fadeIn 1s ease;
}

.radio-option {
  font-size: 18px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 10px;
  transition: color 0.3s;
}

.radio-option:hover {
  color: #007bff;
}

.radio-option input {
  transform: scale(1.2);
  cursor: pointer;
}

.local-upload-section,
.cloud-upload-section {
  margin-bottom: 30px;
  animation: fadeIn 1s ease;
}

.input-group {
  margin-bottom: 20px;
  animation: fadeIn 1s ease;
}

.input-label {
  display: block;
  font-size: 16px;
  margin-bottom: 8px;
  color: #555;
}

.helper-text {
  font-size: 12px;
  color: #888;
  margin-left: 10px;
}

.input-field,
.select-field {
  width: 100%;
  padding: 12px 16px;
  border: 2px solid #ccc;
  border-radius: 8px;
  font-size: 16px;
  transition: border-color 0.3s, box-shadow 0.3s;
}

.input-field:focus,
.select-field:focus {
  border-color: #007bff;
  box-shadow: 0 0 8px rgba(0, 123, 255, 0.2);
  outline: none;
}

.or-text {
  text-align: center;
  margin: 10px 0;
  color: #888;
  font-size: 16px;
  animation: fadeIn 1s ease;
}

.upload-form-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  background-color: #f0f2f5;
  padding: 30px;
  border-radius: 10px;
  border: 2px dashed #ccc;
  transition: box-shadow 0.3s, border-color 0.3s;
  cursor: pointer;
  width: 100%;
  animation: fadeIn 1s ease;
}

.upload-form-container:hover {
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.1);
  border-color: #007bff;
}

.upload-form {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 100%;
}

.file-input-label {
  padding: 20px;
  background: #fff;
  border: 2px solid #007bff;
  border-radius: 10px;
  color: #333;
  font-size: 18px;
  width: 100%;
  text-align: center;
  transition: background-color 0.3s, color 0.3s;
}

.file-input-label:hover {
  background-color: #e6f0ff;
}

.file-input {
  display: none;
}

.button {
  width: 100%;
  padding: 15px;
  border: none;
  border-radius: 10px;
  background-color: #007bff;
  color: white;
  font-size: 18px;
  cursor: pointer;
  transition: box-shadow 0.3s, transform 0.3s, background-color 0.3s;
  animation: fadeInUp 1s ease;
}

.button:hover {
  box-shadow: 0 8px 20px rgba(0, 123, 255, 0.3);
  transform: translateY(-3px);
  background-color: #0056b3;
}

.button:active {
  transform: translateY(1px);
  box-shadow: 0 4px 10px rgba(0, 123, 255, 0.2);
}

.task-status {
  margin-top: 20px;
  text-align: center;
  color: #555;
  font-size: 16px;
  animation: fadeIn 1s ease;
}

/* 上传进度条样式 */
.progress-bar-container {
  position: relative;
  width: 100%;
  height: 25px;
  background-color: #e0e0e0;
  border-radius: 12.5px;
  overflow: hidden;
  margin-bottom: 20px;
  animation: fadeIn 1s ease;
}

.progress-bar {
  height: 100%;
  background-color: #007bff;
  width: 0%;
  transition: width 0.4s ease;
}

.progress-text {
  position: absolute;
  top: 0;
  left: 50%;
  transform: translateX(-50%);
  line-height: 25px;
  color: #fff;
  font-weight: bold;
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
