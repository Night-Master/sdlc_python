<template>
  <div class="pageView">
    <div class="xml-container">
      <h1>XML Input</h1>
      <form @submit.prevent="handleSubmit">
        <div class="input-group">
          <label for="xmlInput">Enter XML</label>
          <textarea
            id="xmlInput"
            v-model="xmlInput"
            placeholder="Enter your XML here"
            required
          ></textarea>
        </div>
        <button type="submit" class="btn-submit">Submit</button>
      </form>
      <div v-if="message" :class="['message', messageType]">{{ message }}</div>
      <div v-if="parsedData" class="parsed-data">
        <h2>Parsed Data</h2>
        <pre>{{ parsedData }}</pre>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import request from '@/utils/request';

const xmlInput = ref('');
const message = ref('');
const messageType = ref('');
const parsedData = ref('');

const handleSubmit = async () => {
  try {
    const response = await request({
      url: '/xxe',
      method: 'post',
      data: {
        xml: xmlInput.value
      }
    });
    if (response.status === 1) {
      message.value = response.message;
      messageType.value = 'success';
      parsedData.value = response.data;
    } else {
      message.value = response.message;
      messageType.value = 'error';
      parsedData.value = '';
    }
  } catch (error) {
    message.value = 'An error occurred';
    messageType.value = 'error';
    parsedData.value = '';
  }
};
</script>

<style scoped lang="less">
.pageView {
  width: 100%;
  height: 100%;
  position: relative;
  font-family: 'Arial', sans-serif;
  background-color: #f4f4f4;
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  margin: 0;
}
.xml-container {
  background-color: #fff;
  padding: 40px;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  width: 600px;
  text-align: center;
}
.xml-container h1 {
  margin-bottom: 20px;
  color: #333;
}
.input-group {
  margin-bottom: 20px;
  text-align: left;
}
.input-group label {
  display: block;
  margin-bottom: 5px;
  color: #555;
}
.input-group textarea {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  box-sizing: border-box;
  resize: vertical;
  min-height: 200px;
}
.input-group textarea:focus {
  border-color: #007bff;
  outline: none;
}
.btn-submit {
  background-color: #007bff;
  color: #fff;
  padding: 10px 20px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  width: 100%;
  font-size: 16px;
}
.btn-submit:hover {
  background-color: #0056b3;
}
.message {
  margin-top: 20px;
  padding: 10px;
  border-radius: 4px;
  font-size: 14px;
}
.success {
  background-color: #d4edda;
  color: #155724;
  border: 1px solid #c3e6cb;
}
.error {
  background-color: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
}
.parsed-data {
  margin-top: 20px;
  text-align: left;
  background-color: #f9f9f9;
  padding: 10px;
  border-radius: 4px;
  border: 1px solid #ddd;
}
.parsed-data h2 {
  margin-bottom: 10px;
  color: #333;
}
.parsed-data pre {
  white-space: pre-wrap;
  word-wrap: break-word;
  font-size: 14px;
  color: #555;
}
</style>