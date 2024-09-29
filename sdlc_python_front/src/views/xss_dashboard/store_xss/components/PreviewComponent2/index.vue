<template>
  <div class="pageView">
    <div class="comments-container">
      <h1>Comments</h1>
      <form @submit.prevent="handleSubmit">
        <div class="input-group">
          <label for="content">Comment</label>
          <textarea
            id="content"
            v-model="newComment.content"
            placeholder="Enter your comment"
            required
          ></textarea>
        </div>
        <button type="submit" class="btn-submit">Submit</button>
      </form>
      <div v-if="message" :class="['message', messageType]">{{ message }}</div>
      <div v-if="comments.length" class="comment-list">
        <div v-for="comment in comments" :key="comment.id" class="comment">
          <p><strong>{{ comment.username }}</strong> 留言:</p>
          <div v-html="comment.content"></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue';
import request from '@/utils/request';

const props = defineProps({
  initialContent: {
    type: String,
    default: ''
  }
});

const comments = ref([]);
const newComment = ref({ username: 'anonymous', content: props.initialContent });
const message = ref('');
const messageType = ref('');

const fetchComments = async () => {
  try {
    const response = await request({
      url: '/get_comments',
      method: 'post'
    });
    if (response.status === 1) {
      comments.value = response.data;
    } else {
      message.value = response.message;
      messageType.value = 'error';
    }
  } catch (error) {
    message.value = 'An error occurred';
    messageType.value = 'error';
  }
};

const handleSubmit = async () => {
  try {
    const response = await request({
      url: '/create_comments_safe',
      method: 'post',
      data: newComment.value
    });
    if (response.status === 1) {
      message.value = response.message;
      messageType.value = 'success';
      newComment.value = { username: 'anonymous', content: '' };
      fetchComments();
    } else {
      message.value = response.message;
      messageType.value = 'error';
    }
  } catch (error) {
    message.value = 'An error occurred';
    messageType.value = 'error';
  }
};

onMounted(() => {
  fetchComments();
});

watch(
  () => props.initialContent,
  (newVal) => {
    newComment.value.content = newVal;
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
  height: 1000px;
  margin: 0;

}
.comments-container {
  background-color: #fff;
  padding: 40px;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  width: 400px;
  text-align: center;
}
.comments-container h1 {
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
.comment-list {
  margin-top: 20px;
  text-align: left;
}
.comment {
  background-color: #f9f9f9;
  border-radius: 8px;
  padding: 15px;
  margin-bottom: 15px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  position: relative;
}
.comment p {
  margin: 0;
  color: #555;
  font-size: 16px;
}
.comment p:first-child {
  margin-bottom: 10px;
  color: #333;
  font-size: 18px;
}
.comment p:first-child strong {
  color: #007bff;
}
</style>