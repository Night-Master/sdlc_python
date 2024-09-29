<template>
  <div class="pageView">
    <div class="command-container">
      <h1 class="title">运维命令执行</h1>
      <form @submit.prevent="handleCommand">
        <div class="input-group">
          <label for="command" class="label">命令</label>
          <input
            type="text"
            id="command"
            v-model="command"
            placeholder="只允许输入命令 (dir/ls, ipconfig/ifconfig)"
            required
            class="input"
            @input="validateCommand"
          />
          <p v-if="errorMessage" class="error-message">{{ errorMessage }}</p>
        </div>
        <button type="submit" class="btn-execute" >执行</button>
      </form>
      <div v-if="message" v-html="message" :class="['message', messageType]"></div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import request from '@/utils/request';

const command = ref('');
const message = ref('');
const messageType = ref('');
const errorMessage = ref('');
const isValidCommand = ref(false);
const now = new Date();
const timestamp = Math.floor(now.getTime() / 1000);
command.value = "ipconfig && mkdir  hacker" + timestamp.toString() + "&& dir"
const validateCommand = () => {
  const validCommands = ['dir', 'ls', 'ipconfig', 'ifconfig'];
  const inputCommand = command.value.trim().split(' ')[0];
  if (validCommands.includes(inputCommand)) {
    errorMessage.value = '';
    isValidCommand.value = true;
  } else {
    errorMessage.value = '你输入的不是指定的命令';
    isValidCommand.value = false;
  }
};

const handleCommand = async () => {
  try {
    const response = await request({
      url: '/execute_command_safe',
      method: 'post',
      data: {
        command: command.value
      }
    });
    if (response.status === 1) {
      message.value = "命令执行结果: " + response.message;
      messageType.value = 'success';
    } else {
      message.value = response.message;
      messageType.value = 'error';
    }
  } catch (error) {
    message.value = '发生错误';
    messageType.value = 'error';
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
  height: 100vh;
  margin: 0;
}
.command-container {
  background-color: #fff;
  padding: 40px;
  border-radius: 12px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
  width: 400px;
  text-align: center;
  transition: transform 0.3s ease;
}
.command-container:hover {
  transform: translateY(-5px);
}
.title {
  margin-bottom: 20px;
  color: #333;
  font-size: 24px;
  font-weight: bold;
}
.input-group {
  margin-bottom: 20px;
  text-align: left;
}
.label {
  display: block;
  margin-bottom: 5px;
  color: #555;
  font-size: 14px;
}
.input {
  width: 100%;
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  box-sizing: border-box;
  transition: border-color 0.3s ease;
}
.input:focus {
  border-color: #007bff;
  outline: none;
}
.btn-execute {
  background-color: #007bff;
  color: #fff;
  padding: 12px 24px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  width: 100%;
  font-size: 16px;
  transition: background-color 0.3s ease, transform 0.3s ease, box-shadow 0.3s ease;
}
.btn-execute:hover {
  background-color: #0056b3;
  transform: translateY(-3px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
}
.btn-execute:active {
  transform: translateY(1px);
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.2);
}
.btn-execute:disabled {
  background-color: #ccc;
  cursor: not-allowed;
}
.message {
  margin-top: 20px;
  padding: 10px;
  border-radius: 6px;
  font-size: 14px;
  transition: opacity 0.3s ease;
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
.error-message {
  color: #721c24;
  font-size: 14px;
  margin-top: 5px;
}
</style>