import axios from 'axios';

const service = axios.create({
  // baseURL: 'http://127.0.0.1:2334',
  baseURL: 'http://123.249.114.216:2334',
  // baseURL: 'http://192.168.5.176:2333',

  timeout: 50000
});

// 请求拦截器
service.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers['Authorization'] = token;
    }
    return config;
  },
  error => {
    return Promise.reject(error);
  }
);

// 响应拦截器
service.interceptors.response.use(
  response => {
    if (response.data.token) {
      localStorage.setItem('token', response.data.token);
    }
    if (response.data.username) {
      localStorage.setItem('username', response.data.username);
    }
    return response.data;
  },
  error => {
    if (error.response && error.response.status === 200 && error.response.data.status === 0) {
      // 处理未登录的情况
      console.error(error.response.data.message);
      // 你可以在这里添加更多的处理逻辑，例如跳转到登录页面
    }
    return Promise.reject(error);
  }
);

export default service;