<template>
  <div class="file-upload-page">
    <div class="header">
      <h1>文件上传</h1>
    </div>
    <div class="content">
      <div class="upload-form-container" @dragover.prevent @drop.prevent="handleDrop">
        <div class="upload-form" @click="triggerFileInput">
          <label for="file-input" class="file-input-label">
            <span v-if="!file">点击或拖拽图片文件到这里</span>
            <span v-else>{{ file.name }}</span>
          </label>
          <input id="file-input" type="file" @change="handleFileChange" class="file-input" />
        </div>
        <button @click="uploadFile" class="button">上传图片文件</button>
        <div class="upload-status" v-if="uploadStatus">
          <p>{{ uploadStatus }}</p>
        </div>
        <div class="file-path" v-if="filePath">
          <p>图片文件路径: {{ filePath }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import request from '@/utils/request';

const file = ref(null);
const uploadStatus = ref('');
const filePath = ref('');

const handleFileChange = (event) => {
  file.value = event.target.files[0];
};

const handleDrop = (event) => {
  file.value = event.dataTransfer.files[0];
};

const triggerFileInput = () => {
  document.getElementById('file-input').click();
};

const uploadFile = async () => {
  if (!file.value) {
    uploadStatus.value = '请选择一个文件';
    return;
  }

  const formData = new FormData();
  formData.append('file', file.value);

  try {
    const response = await request({
      url: '/upload_file_safe',
      method: 'post',
      data: formData,
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    });

    if (response.status === 1) {
      uploadStatus.value = '文件上传成功';
      filePath.value = response.filePath;
    } else {
      uploadStatus.value = response.message;
    }
  } catch (error) {
    uploadStatus.value = '上传失败，请重试';
  }
};
</script>

<style scoped>
.file-upload-page {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 40px;
  font-family: 'Helvetica Neue', Arial, sans-serif;
  background-color: #f0f2f5;
  min-height: 100vh;
}

.header {
  text-align: center;
  margin-bottom: 40px;
  color: #333;
  font-size: 24px;
}

.content {
  width: 100%;
  max-width: 600px;
}

.upload-form-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  background-color: #fff;
  padding: 40px;
  border-radius: 10px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
  transition: box-shadow 0.3s;
  border: 2px dashed #ccc;
  cursor: pointer;
}

.upload-form-container:hover {
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.1);
  border-color: #007BFF;
}

.upload-form {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 100%;
  margin-bottom: 20px;
}

.file-input-label {
  padding: 20px;
  border: none;
  border-radius: 6px;
  background: #f9f9f9;
  color: #333;
  font-size: 16px;
  cursor: pointer;
  transition: background-color 0.3s, color 0.3s;
  width: 100%;
  text-align: center;
}

.file-input-label:hover {
  background-color: #e9e9e9;
}

.file-input {
  display: none; /* Hide the default file input */
}

.button {
  padding: 12px 24px;
  border: none;
  border-radius: 6px;
  background: linear-gradient(135deg, #007BFF, #0056b3);
  color: #fff;
  font-size: 16px;
  cursor: pointer;
  transition: transform 0.3s, box-shadow 0.3s;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.button:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 12px rgba(0, 0, 0, 0.2);
}

.button:active {
  transform: translateY(1px);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.upload-status, .file-path {
  margin-top: 20px;
  text-align: center;
  color: #555;
  font-size: 16px;
}
</style>