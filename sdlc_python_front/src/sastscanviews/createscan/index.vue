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

      <!-- 本地模式下可以输入路径或者拖拽上传 -->
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

      <!-- 云端模式下只能上传 -->
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

      <button @click="createTask" class="button">创建任务</button>
      <div class="task-status" v-if="taskStatus">
        <p>{{ taskStatus }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import request from '@/utils/request'; // 使用你封装的 request 工具

const scanMode = ref('local'); // 模式选择，本地或云端
const scanPath = ref(''); // 本地路径
const file = ref(null); // 上传的文件
const taskStatus = ref(''); // 任务状态
const fileInput = ref(null); // 文件输入引用

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

// 创建扫描任务
const createTask = async () => {
  taskStatus.value = ''; // 重置任务状态

  // 本地扫描逻辑
  if (scanMode.value === 'local') {
    if (!scanPath.value && !file.value) {
      taskStatus.value = '请输入本地路径或上传压缩文件';
      return;
    }

    if (scanPath.value) {
      // 通过路径扫描
      try {
        const response = await request({
          url: '/scan_start_dir',
          method: 'post',
          data: {
            path: scanPath.value,
          },
        });

        if (response.status === 1 && response.taskID) {
          taskStatus.value = `任务创建成功，任务ID: ${response.taskID}`;
        } else {
          taskStatus.value = response.message || '任务创建失败';
        }
      } catch (error) {
        taskStatus.value = '任务创建失败，请重试';
      }
    } else if (file.value) {
      // 通过上传压缩包扫描
      const formData = new FormData();
      formData.append('file', file.value);

      try {
        const response = await request({
          url: '/scan_start_zip',
          method: 'post',
          data: formData,
          headers: {
            'Content-Type': 'multipart/form-data',
          },
        });

        if (response.status === 1 && response.taskID) {
          taskStatus.value = `任务创建成功，任务ID: ${response.taskID}`;
        } else {
          taskStatus.value = response.message || '任务创建失败';
        }
      } catch (error) {
        taskStatus.value = '任务创建失败，请重试';
      }
    }
  }

  // 云端扫描逻辑
  else if (scanMode.value === 'cloud') {
    if (!file.value) {
      taskStatus.value = '请上传压缩文件';
      return;
    }

    const formData = new FormData();
    formData.append('file', file.value);

    try {
      const response = await request({
        url: '/scan_start_zip',
        method: 'post',
        data: formData,
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });

      if (response.status === 1 && response.taskID) {
        taskStatus.value = `任务创建成功，任务ID: ${response.taskID}`;
      } else {
        taskStatus.value = response.message || '任务创建失败';
      }
    } catch (error) {
      taskStatus.value = '任务创建失败，请重试';
    }
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
  background-color: #f5f7fa;
  min-height: 100vh;
}

.header {
  text-align: center;
  margin-bottom: 40px;
  color: #333;
  font-size: 32px;
  font-weight: bold;
}

.content {
  width: 100%;
  max-width: 700px;
  background: #fff;
  padding: 40px;
  border-radius: 15px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

.mode-selection {
  display: flex;
  justify-content: center;
  gap: 40px;
  margin-bottom: 30px;
}

.radio-option {
  font-size: 18px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 10px;
}

.radio-option input {
  transform: scale(1.2);
}

.local-upload-section,
.cloud-upload-section {
  margin-bottom: 30px;
}

.input-field {
  width: 100%;
  padding: 12px 16px;
  border: 2px solid #ccc;
  border-radius: 8px;
  margin-bottom: 20px;
  font-size: 16px;
  transition: border-color 0.3s;
}

.input-field:focus {
  border-color: #007bff;
  outline: none;
}

.or-text {
  text-align: center;
  margin: 10px 0;
  color: #888;
  font-size: 16px;
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
}

.upload-form-container:hover {
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
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
  transition: box-shadow 0.3s, transform 0.3s;
}

.button:hover {
  box-shadow: 0 4px 15px rgba(0, 123, 255, 0.3);
  transform: translateY(-2px);
}

.button:active {
  transform: translateY(1px);
  box-shadow: 0 2px 10px rgba(0, 123, 255, 0.2);
}

.task-status {
  margin-top: 20px;
  text-align: center;
  color: #555;
  font-size: 16px;
}
</style>
