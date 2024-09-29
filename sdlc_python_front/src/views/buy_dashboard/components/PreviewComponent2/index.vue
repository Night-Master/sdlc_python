<template>
  <div class="pageView">
    <div class="product-container">
      <h1>Product List</h1>
      <div v-if="products.length === 0" class="no-products">No products available</div>
      <div v-else>
        <div v-for="product in products" :key="product.id" class="product-item">
          <h2>{{ product.name }}</h2>
          <p>Price: ${{ product.price }}</p>
          <button @click="selectProduct(product)" class="btn-select">Select</button>
        </div>
      </div>
      <div v-if="selectedProduct" class="purchase-form">
        <h2>Purchase {{ selectedProduct.name }}</h2>
        <div class="input-group">
          <label for="quantity">Quantity</label>
          <input
            type="number"
            id="quantity"
            v-model="quantity"
            placeholder="Enter quantity"
            required
          />
        </div>
        <button @click="handlePurchase" class="btn-purchase">Purchase</button>
      </div>
      <div v-if="message" v-html="message" :class="['message', messageType]"></div>
      <div v-if="balance !== null" class="balance">Balance: ${{ balance }}</div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import request from '@/utils/request';
const products = ref([]);
const selectedProduct = ref(null);
const quantity = ref(1);
const message = ref('');
const messageType = ref('');
const balance = ref(null);

const fetchProducts = async () => {
  try {
    const response = await request({
      url: '/products',
      method: 'get'
    });
    if (response.status === 1) {
      products.value = response.data;
    } else {
      message.value = response.message;
      messageType.value = 'error';
    }
  } catch (error) {
    message.value = 'An error occurred';
    messageType.value = 'error';
  }
};

const selectProduct = (product) => {
  selectedProduct.value = product;
};

const handlePurchase = async () => {
  if (!selectedProduct.value) return;

  try {
    const response = await request({
      url: '/purchase_safe',
      method: 'post',
      data: {
        productId: selectedProduct.value.id,
        quantity: quantity.value
      }
    });
    if (response.status === 1) {
      message.value = `购买成功: 购买了 ${quantity.value} 件 ${selectedProduct.value.name}`;
      messageType.value = 'success';
      balance.value = response.balance;
    } else {
      message.value = response.message;
      messageType.value = 'error';
      balance.value = response.balance;
    }
  } catch (error) {
    message.value = 'An error occurred';
    messageType.value = 'error';
  }
};

onMounted(() => {
  fetchProducts();
});
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
  align-items: flex-start; /* 修改为 flex-start 以确保顶部对齐 */
  padding-top: 50px; /* 添加 padding-top 以确保内容与顶部保持一定距离 */
  margin: 0;
}
.product-container {
  background-color: #fff;
  padding: 40px;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  width: 400px;
  text-align: center;
  position: absolute; /* 修改为 absolute */
  top: 50px; /* 设置与顶部的固定距离 */
}
.product-container h1 {
  margin-bottom: 20px;
  color: #333;
}
.no-products {
  color: #777;
  font-style: italic;
}
.product-item {
  margin-bottom: 20px;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
}
.product-item h2 {
  margin-bottom: 10px;
  color: #333;
}
.product-item p {
  color: #555;
}
.btn-select {
  background-color: #007bff;
  color: #fff;
  padding: 10px 20px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
}
.btn-select:hover {
  background-color: #0056b3;
}
.purchase-form {
  margin-top: 20px;
}
.purchase-form h2 {
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
.btn-purchase {
  background-color: #28a745;
  color: #fff;
  padding: 10px 20px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  width: 100%;
  font-size: 16px;
}
.btn-purchase:hover {
  background-color: #218838;
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
.balance {
  margin-top: 20px;
  font-size: 16px;
  color: #333;
}
</style>