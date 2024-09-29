<template>
  <div class="personal-page">
    <div class="header">
      <input v-model="inputUsername" placeholder="输入用户名" @keyup.enter="fetchUserProfile" class="input-field" />
      <button @click="fetchUserProfile" class="button">查看用户</button>
    </div>
    <div class="content">
      <div class="profile">
        <h1>你当前登录的账号是 {{ currentUsername }}，你查看的是 {{ user.username }} 的用户界面</h1>
        <div class="avatar-container">
          <img :src="avatarSrc" alt="Avatar" class="avatar" @error="setDefaultAvatar">
        </div>
        <div class="info">
          <p><strong>Signature:</strong> {{ user.signature }}</p>
          <p><strong>Email:</strong> {{ user.email }}</p>
          <p><strong>Birthdate:</strong> {{ user.birthdate }}</p>
          <p><strong>Balance:</strong> {{ user.balance }}</p>
        </div>
      </div>
      <div class="orders">
        <h2>Orders</h2>
        <table>
          <thead>
            <tr>
              <th>Order ID</th>
              <th>Product Name</th>
              <th>Amount</th>
              <th>Quantity</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="order in orders" :key="order.order_id">
              <td>{{ order.order_id }}</td>
              <td>{{ order.name }}</td>
              <td>{{ order.amount }}</td>
              <td>{{ order.quantity }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import request from '@/utils/request';
import defaultAvatar from '@/assets/imgs/avatar.jpg'; // 设置默认头像路径

const user = ref({
  username: '',
  signature: '',
  email: '',
  birthdate: '',
  avatar: '',
  balance: 0
});

const orders = ref([]);
const currentUsername = ref(localStorage.getItem('username'));
const inputUsername = ref('user2'); // 默认输入值为 user2

const avatarSrc = computed(() => {
  return user.value.avatar || defaultAvatar;
});

const fetchUserProfile = async () => {
  try {
    const response = await request({
      url: '/get_profile_unauthorized',
      method: 'post',
      data: { 'username': inputUsername.value } // 使用用户输入的用户名
    });
    if (response.status === 1) {
      user.value = response.data.user;
      orders.value = response.data.orders;
    } else {
      alert(response.message);
    }
  } catch (error) {
    alert('An error occurred');
  }
};

const setDefaultAvatar = (event) => {
  if (event.target.src !== defaultAvatar) {
    event.target.src = defaultAvatar;
  } else {
    event.target.onerror = null; // 移除错误处理程序，防止无限加载
  }
};

onMounted(() => {
  // 清空 localStorage 并弹窗提示
  localStorage.removeItem('token');
  localStorage.removeItem('username');
  alert('由于本漏洞演示需要未登录状态，已经清空登录凭据');

  fetchUserProfile();
});
</script>

<style scoped>
.personal-page {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
  font-family: 'Helvetica Neue', Arial, sans-serif;
  background-color: #fafafa;
}

.header {
  text-align: center;
  margin-bottom: 40px; /* 增加间隔 */
  color: #333;
}

.input-field {
  margin-right: 10px;
  padding: 10px 15px;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-size: 16px;
  transition: border-color 0.3s, box-shadow 0.3s;
}

.input-field:focus {
  border-color: #007BFF;
  box-shadow: 0 0 5px rgba(0, 123, 255, 0.5);
  outline: none;
}

.button {
  padding: 10px 20px;
  border: none;
  border-radius: 4px;
  background-color: #007BFF;
  color: #fff;
  font-size: 16px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.button:hover {
  background-color: #0056b3;
}

.content {
  width: 100%;
  max-width: 600px;
}

.profile {
  display: flex;
  flex-direction: column;
  align-items: center;
  background-color: #fff;
  padding: 30px;
  border-radius: 10px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
}

.avatar-container {
  width: 150px;
  height: 150px;
  border-radius: 50%;
  overflow: hidden;
  margin-bottom: 20px;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.1);
}

.avatar {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.info {
  text-align: center;
  color: #555;
}

.info p {
  margin: 10px 0;
}

.orders {
  margin-top: 20px;
  width: 100%;
  background-color: #fff;
  padding: 20px;
  border-radius: 10px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
}

.orders table {
  width: 100%;
  border-collapse: collapse;
}

.orders th, .orders td {
  border: 1px solid #ddd;
  padding: 8px;
  text-align: center;
}

.orders th {
  background-color: #f2f2f2;
  color: #333;
}

.orders tr:nth-child(even) {
  background-color: #f9f9f9;
}

.orders tr:hover {
  background-color: #ddd;
}
</style>