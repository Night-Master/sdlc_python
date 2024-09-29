<template>
  <div class="image-gallery-page">
    <div class="header">
      <h1>图片下载(下载文件上传中上传的图片)</h1>
    </div>
    <div class="content">
      <div class="image-list">
        <div v-for="image in images" :key="image.name" class="image-item">
          <div class="image-name">{{ image.name }}</div>
          <div class="image-actions">
            <button @click="downloadFile(image.name)" class="button">下载</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import request from '@/utils/request';

const images = ref([]);

const fetchImages = async () => {
  try {
    const response = await request({
      url: '/list_images',
      method: 'get'
    });

    if (response.status === 1) {
      images.value = response.images.map(image => ({
        name: image,
        url: `http://127.0.0.1:2333/${image}`
      }));
    } else {
      console.error(response.message);
    }
  } catch (error) {
    console.error('获取图片列表失败，请重试', error);
  }
};

const downloadFile = async (fileName) => {
  try {
    const response = await request({
      url: '/download_safe',
      method: 'post',
      data: { fileName },
      responseType: 'blob'
    });

    const url = window.URL.createObjectURL(new Blob([response]));
    const link = document.createElement('a');
    link.href = url;
    link.setAttribute('download', fileName);
    document.body.appendChild(link);
    link.click();
    link.remove();
  } catch (error) {
    console.error('文件下载失败，请重试', error);
  }
};

onMounted(() => {
  fetchImages();
});
</script>

<style scoped>
.image-gallery-page {
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
  max-width: 800px;
}

.image-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 20px;
}

.image-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  background-color: #fff;
  padding: 20px;
  border-radius: 10px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
  transition: box-shadow 0.3s;
}

.image-item:hover {
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.1);
}

.image-name {
  font-size: 18px;
  margin-bottom: 10px;
}

.image-actions {
  display: flex;
  justify-content: center;
  width: 100%;
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
  text-decoration: none;
}

.button:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 12px rgba(0, 0, 0, 0.2);
}

.button:active {
  transform: translateY(1px);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}
</style>