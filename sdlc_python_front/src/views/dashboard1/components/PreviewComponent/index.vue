<template>
  <div class="pageView">
    <div class="login-container">
      <h1>Login</h1>
      <form action="#" method="post">
        <div class="input-group">
          <label for="username">Username</label>
          <input
            type="text"
            id="username"
            name="username"
            placeholder="Enter your username"
            required
          />
        </div>
        <div class="input-group">
          <label for="password">Password</label>
          <input
            type="password"
            id="password"
            name="password"
            placeholder="Enter your password"
            required
          />
        </div>
        <button type="submit" class="btn-login">Login</button>
      </form>
    </div>
  </div>
</template>
<script setup></script>
<style scoped lang="less">
.pageView {
  width: 100%;
  height: 100%;
  position: relative;
  //   background-color: #f7f7f7;
  font-family: 'Arial', sans-serif;
  background-color: #f4f4f4;
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  margin: 0;
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
