<template>
  <div class="pageView">
    <div class="login-container">
      <h1>{{title}} 漏洞平台</h1>
      <form @submit.prevent="handleLogin">
        <div class="input-group">
          <label for="username">Username</label>
          <input
            type="text"
            id="username"
            v-model="username"
            placeholder="Enter your username"
            required
          />
        </div>
        <div class="input-group">
          <label for="password">Password</label>
          <input
            type="password"
            id="password"
            v-model="password"
            placeholder="Enter your password"
            required
          />
        </div>
        <button type="submit" class="btn-login">Login</button>
      </form>
      <div v-if="message" :class="['message', messageType]">{{ message }}</div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, onMounted,watchEffect } from 'vue'
import request from '@/utils/request'
import { useAppStore } from '@/store/modules/app'
import { useCache } from '@/hooks/web/useCache'
import { usePermissionStore } from '@/store/modules/permission'
import { useRouter } from 'vue-router'
import JSEncrypt from 'jsencrypt'
const appStore = useAppStore()
const permissionStore = usePermissionStore()
const { currentRoute, addRoute, push } = useRouter()
const props = defineProps({
  initialUsername: String,
  initialPassword: String
})

const username = ref(props.initialUsername || '')
const password = ref(props.initialPassword || '')
const message = ref('')
const messageType = ref('')
const redirect = ref('')
const publicKey = ref('')
const { wsCache } = useCache()
const encrypt = new JSEncrypt()
const title = ref('')
watchEffect(()=>{
   title.value = appStore.getTitle
console.log(title.value,'获取信息')
})

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

const handleLogin = async () => {
  if (!publicKey.value) {
    message.value = 'Public key not loaded'
    messageType.value = 'error'
    return
  }

  try {
    const encryptedUsername = encrypt.encrypt(username.value)
    const encryptedPassword = encrypt.encrypt(password.value)

    const response = await request({
      url: '/sql_injection_sqlite3_safe',
      method: 'post',
      data: {
        username: encryptedUsername,
        password: encryptedPassword
      }
    })

    if (response.status == 1) {
      message.value = response.message + "\n当前登录的用户：" + response.username
      username.value = response.username
      messageType.value = 'success'

      const res = {
        password: '',
        role: 'admin',
        roleId: '1',
        username: response.username,
        permissions: ['*.*.*']
      }

      wsCache.set(appStore.getUserInfo, res)
      window.location.href = '/#/sql_dashboard'
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

watch(
  () => currentRoute.value,
  (route) => {
    redirect.value = route?.query?.redirect
  },
  {
    immediate: true
  }
)

watch(
  () => [props.initialUsername, props.initialPassword],
  ([newUsername, newPassword]) => {
    username.value = newUsername || ''
    password.value = newPassword || ''
  }
)

onMounted(() => {
  fetchPublicKey()
})
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

.login-container {
  background-color: #fff;
  padding: 40px;
  border-radius: 10px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.05);
  width: 400px;
  text-align: center;
  transition: transform 0.3s ease;
}

.login-container:hover {
  transform: translateY(-5px);
}

.login-container h1 {
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

.btn-login {
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

.btn-login:hover {
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