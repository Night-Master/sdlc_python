<template>
  <div class="header-t">
    <ul>
      <li>安全编码指南</li>
      <li>Github</li>
      <el-popover
        v-if="userInfo.username"
        trigger="hover"
        :width="100"
        popper-style="text-align: center;width: 100px;min-width:100px;"
      >
        <div>
          <span style="width: 100px; text-align: center; cursor: pointer" @click="show()"
            >修改密码</span
          >
        </div>
        <template #reference>
          <li v-if="userInfo.username">{{ userInfo.username }}</li>
        </template>
      </el-popover>
      <li v-if="!userInfo.username" @click="Login">登录</li>
      <li @click="goLogin">注销</li>
    </ul>

    <el-dialog v-model="dialogVisible" @close="close" title="修改密码" width="500">
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
    </el-dialog>
  </div>
</template>
<script setup>
import { useRouter } from 'vue-router'
const { push } = useRouter()
import { ElPopover, ElDialog, ElMessage } from 'element-plus'
import { ref } from 'vue'
import { useAppStoreWithOut } from '@/store/modules/app'
import { useCache } from '@/hooks/web/useCache'
const { wsCache } = useCache()
const appStore = useAppStoreWithOut()
const userInfo = wsCache.get(appStore.getUserInfo) || { username: '' }
console.log(userInfo, 33333)
const dialogVisible = ref(false)
import request from '@/utils/request'
import JSEncrypt from 'jsencrypt'

const currentPassword = ref('')
const newPassword = ref('')
const confirmNewPassword = ref('')
const message = ref('')
const messageType = ref('')
const publicKey = ref('')

const encrypt = new JSEncrypt()
const show = () => {
  dialogVisible.value = true
}
const close = () => {
  currentPassword.value = ''
  newPassword.value = ''
  confirmNewPassword.value = ''
  message.value = ''
  messageType.value = ''
}
const fetchPublicKey = async () => {
  try {
    const response = await request({
      url: '/getPublicKey',
      method: 'get'
    })
    publicKey.value = response
    encrypt.setPublicKey(publicKey.value)
  } catch (error) {
    message.value = 'Failed to fetch public key'
    messageType.value = 'error'
  }
}

const handleChangePassword = async () => {
  if (!publicKey.value) {
    message.value = 'Public key not loaded'
    messageType.value = 'error'
    return
  }

  if (newPassword.value !== confirmNewPassword.value) {
    message.value = '新密码和确认密码不匹配'
    messageType.value = 'error'
    return
  }

  try {
    const encryptedCurrentPassword = encrypt.encrypt(currentPassword.value)
    const encryptedNewPassword = encrypt.encrypt(newPassword.value)

    const response = await request({
      url: '/change_password_safe',
      method: 'post',
      data: {
        currentPassword: encryptedCurrentPassword,
        newPassword: encryptedNewPassword
      }
    })

    if (response.status == 1) {
      // message.value = response.message
      // messageType.value = 'success'
      currentPassword.value = ''
      newPassword.value = ''
      confirmNewPassword.value = ''
      dialogVisible.value = false
      ElMessage({
        message: 'success',
        type: 'success'
      })
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

fetchPublicKey()

const goLogin = () => {
  if (!userInfo.username) return false
  wsCache.delete(appStore.getUserInfo)
  localStorage.setItem('token', '')
  localStorage.setItem('username', '')
  push('/login')
}
const Login = () => {
  push('/login')
}
</script>
<style lang="less" scoped>
ul {
  display: flex;
  li {
    margin: 0 20px;
    color: #007bff;
    text-decoration: none;
    transition: color 0.3s;
    cursor: pointer;
  }
}
.header-t {
  width: 100%;
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
