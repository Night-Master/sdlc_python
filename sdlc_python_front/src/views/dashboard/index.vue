<template>
  <div>
    <div class="_title_t">{{ title }}</div>
  <div class="pageView">
    <div class="login-container">
      <h1>Login</h1>
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
    </div>
  </div>
  </div>
</template>

<script setup>
import { ref,watch } from 'vue';
import request from '@/utils/request';
import { useRouter } from 'vue-router'
const { currentRoute } = useRouter()
console.log(currentRoute,3333)
const title = ref(currentRoute.value.meta.title)
// 监听录音变化
watch(
  () => currentRoute.value.meta.title,
  (newVal, oldVal) => {
    title.value = newVal
  }
)
const username = ref('');
const password = ref('');

const handleLogin = async () => {
  try {
    const response = await request({
      url: '/sql_injection_sqlite3',
      method: 'post',
      data: {
        username: username.value,
        password: password.value
      }
    });
    console.log(response);
    // 处理登录成功后的逻辑
  } catch (error) {
    console.error(error);
    // 处理登录失败后的逻辑
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
}
.login-container {
  background-color: #fff;
  padding: 40px;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  width: 400px;
  text-align: center;
}
.login-container h1 {
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
.input-group input {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  box-sizing: border-box;
}
.input-group input:focus {
  border-color: #007bff;
  outline: none;
}
.btn-login {
  background-color: #007bff;
  color: #fff;
  padding: 10px 20px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  width: 100%;
  font-size: 16px;
}
.btn-login:hover {
  background-color: #0056b3;
}
</style>