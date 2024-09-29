<template>
  <div class="pageView">
    <div class="search-container">
      <h1>Search</h1>
      <form @submit.prevent="handleSearch">
        <div class="input-group">
          <label for="searchQuery">Search Query</label>
          <input
            type="text"
            id="searchQuery"
            v-model="searchQuery"
            placeholder="Enter your search query"
            required
          />
        </div>
        <button type="submit" class="btn-search">Search</button>
      </form>
      <div v-if="message" v-html="message" :class="['message', messageType]"></div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue';
import request from '@/utils/request';

const props = defineProps({
  initialSearchQuery: {
    type: String,
    default: ''
  }
});

const searchQuery = ref(props.initialSearchQuery);
const message = ref('');
const messageType = ref('');

const handleSearch = async () => {
  try {
    const response = await request({
      url: '/reflect_xss_safe',
      method: 'post',
      data: {
        input: searchQuery.value
      }
    });
    if (response.status === 1) {
      message.value = "你搜索的内容:"+response.message;
      messageType.value = 'success';
    } else {
      message.value = response.message;
      messageType.value = 'error';
    }
  } catch (error) {
    message.value = 'An error occurred';
    messageType.value = 'error';
  }
};

watch(
  () => props.initialSearchQuery,
  (newQuery) => {
    searchQuery.value = newQuery;
  }
);
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
.search-container {
  background-color: #fff;
  padding: 40px;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  width: 400px;
  text-align: center;
}
.search-container h1 {
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
.btn-search {
  background-color: #007bff;
  color: #fff;
  padding: 10px 20px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  width: 100%;
  font-size: 16px;
}
.btn-search:hover {
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
</style>