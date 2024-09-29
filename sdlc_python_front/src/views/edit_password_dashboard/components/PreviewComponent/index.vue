<template>
  <div class="pageView">
    <div class="change-password-container">
      <h1>修改密码-明文传输</h1>
      <form @submit.prevent="handleChangePassword">
        <div class="input-group">
          <label for="currentPassword">当前密码</label>
          <input
            type="password"
            id="currentPassword"
            v-model="currentPassword"
            placeholder="Enter your current password"
            required
          />
        </div>
        <div class="input-group">
          <label for="newPassword">新密码</label>
          <input
            type="password"
            id="newPassword"
            v-model="newPassword"
            placeholder="Enter your new password"
            required
          />
        </div>
        <div class="input-group">
          <label for="confirmNewPassword">确认新密码</label>
          <input
            type="password"
            id="confirmNewPassword"
            v-model="confirmNewPassword"
            placeholder="Confirm your new password"
            required
          />
        </div>
        <button type="submit" class="btn-change-password">修改密码</button>
      </form>
      <div v-if="message" :class="['message', messageType]">{{ message }}</div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import request from '@/utils/request'

const currentPassword = ref('')
const newPassword = ref('')
const confirmNewPassword = ref('')
const message = ref('')
const messageType = ref('')

const handleChangePassword = async () => {
  if (newPassword.value !== confirmNewPassword.value) {
    message.value = '新密码和确认密码不匹配'
    messageType.value = 'error'
    return
  }

  try {
    const response = await request({
      url: '/change_password_unsafe',
      method: 'post',
      data: {
        currentPassword: currentPassword.value,
        newPassword: newPassword.value
      }
    })

    if (response.status == 1) {
      message.value = response.message
      messageType.value = 'success'
      currentPassword.value = ''
      newPassword.value = ''
      confirmNewPassword.value = ''
    } else {
      message.value = response.message
      messageType.value = 'error'
    }
  } catch (error) {
    console.log(error, 2222)
    message.value = 'An error occurred'
    messageType.value = 'error'
  }
}
</script>

<style scoped lang="less">
.pageView {
  width: 100%;
  height: 100%;
  position: relative;
  font-family: 'Helvetica Neue', sans-serif;
  background: linear-gradient(135deg, #f0f3f8, #e3e7ed);
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  margin: 0;
}

.change-password-container {
  background-color: #fff;
  padding: 40px;
  border-radius: 10px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.05);
  width: 400px;
  text-align: center;
  transition: transform 0.3s ease;
}

.change-password-container:hover {
  transform: translateY(-5px);
}

.change-password-container h1 {
  margin-bottom: 20px;
  color: #333;
  font-size: 28px;
  font-weight: 500;
}

.input-group {
  margin-bottom: 20px;
  text-align: left;
}

.input-group label {
  display: block;
  margin-bottom: 5px;
  color: #555;
  font-size: 14px;
  font-weight: 500;
}

.input-group input {
  width: 100%;
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  box-sizing: border-box;
  transition: border-color 0.3s ease, box-shadow 0.3s ease;
  background-color: #f9f9f9;
}

.input-group input:focus {
  border-color: #007bff;
  outline: none;
  box-shadow: 0 0 0 3px rgba(0, 123, 255, 0.25);
}

.btn-change-password {
  background-color: #007bff;
  color: #fff;
  padding: 12px 24px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  width: 100%;
  font-size: 16px;
  font-weight: 500;
  transition: background-color 0.3s ease, transform 0.3s ease;
}

.btn-change-password:hover {
  background-color: #0056b3;
  transform: translateY(-2px);
}

.message {
  margin-top: 20px;
  padding: 10px;
  border-radius: 6px;
  font-size: 14px;
  text-align: center;
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
</style>