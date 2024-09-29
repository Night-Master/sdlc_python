<template>
  <div class="homeBox"> 
    <div class="_title_t">{{ title }}</div>
    <div class="tips" v-if="currentRoute.meta.desc">
      {{ currentRoute.meta.desc }}
    </div>
    <div class="list">
      <div class="item">
        <div class="_top">
          <ElSelect v-model="value" placeholder="" style="width: 140px">
            <ElOption
              v-for="item in options"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </ElSelect>
          <ElButton type="primary" style="float: right" @click="run1()"> Run </ElButton>
        </div>
        <div class="edit-container h-60vh" v-if="value === 0">
          <pre><code >{{ content }} </code></pre>
        </div>
        <div class="view-container h-60vh" v-else>
          <PreviewComponent :initialUsername="username" :initialPassword="password" />
        </div>
      </div>
      <div class="item">
        <div class="_top">
          <ElSelect v-model="value2" placeholder="" style="width: 140px">
            <ElOption
              v-for="item in options2"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </ElSelect>
          <ElButton type="primary" style="float: right" @click="run2()"> Run </ElButton>
        </div>
        <div class="edit-container h-60vh" v-if="value2 === 0">
          <pre><code >{{ content2 }} </code></pre>
        </div>
        <div class="view-container h-60vh" v-else>
          <PreviewComponentII :initialUsername="username2" :initialPassword="password2"/>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue';
import { ElSelect, ElOption, ElButton } from 'element-plus';
import PreviewComponent from './components/PreviewComponent/index.vue';
import PreviewComponentII from './components/PreviewComponent2/index.vue';
import { useRouter } from 'vue-router'
import { useAppStore } from '@/store/modules/app'
const appStore = useAppStore()
const { currentRoute } = useRouter()
const title = ref(currentRoute.value.meta.title)
watch(
  () => currentRoute.value.meta.title,
  (newVal, oldVal) => {
    title.value = newVal
  }
)
const value = ref(0);
const options = ref([
  { value: 0, label: '漏洞代码' },
  { value: 1, label: '应用界面' }
]);
const value2 = ref(0);
const options2 = ref([
  { value: 0, label: '安全代码' },
  { value: 1, label: '应用界面' }
]);
let content ,content2;
if(appStore.getTitle=="Golang")
{
   content = ref(`
  //获取用户信息的接口放在登录态保护的路由之外
  func SetupRoutes(router *gin.Engine) {
    config := cors.DefaultConfig()
    config.AllowAllOrigins = true
    config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
    router.Use(cors.New(config))
    router.POST("/get_profile_unauthorized", vulnerabilities.Get_profile_unauthorized)
    // 保护需要验证的API
    auth := router.Group("/")
    auth.Use(middleware.AuthMiddleware())
    {
      //这里面放需要登录态访问的接口
    }
  }
  `);
   content2 = ref(`
  //获取用户信息的接口放在登录态保护的路由内
  func SetupRoutes(router *gin.Engine) {
    config := cors.DefaultConfig()
    config.AllowAllOrigins = true
    config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
    router.Use(cors.New(config))
    // 保护需要验证的API
    auth := router.Group("/")
    auth.Use(middleware.AuthMiddleware())
    {
      auth.POST("/get_profile_safe", vulnerabilities.Get_profile_safe)
    }
  }
  `);
}
else{
 content = ref(`
//获取用户信息的接口放在登录态保护的路由之外
def setup_routes(app: Flask):
    CORS(app, resources={r"/*": {"origins": "*", "headers": ["Origin", "Content-Length", "Content-Type", "Authorization"]}})
    app.add_url_rule('/get_profile_unauthorized', 'get_profile_unauthorized', get_profile_unauthorized, methods=['POST'])

    # 保护需要验证的API
    #auth_group = []
    #for route in auth_group:
    #    app.add_url_rule(route, route[1:], auth_middleware(globals()[route[1:]]), methods=['GET', 'POST'])
`);
 content2 = ref(`
//获取用户信息的接口放在登录态保护的路由内
  def setup_routes(app: Flask):
    CORS(app, resources={r"/*": {"origins": "*", "headers": ["Origin", "Content-Length", "Content-Type", "Authorization"]}})
    # 保护需要验证的API
    #auth_group = ['get_profile_safe']
    #for route in auth_group:
    #    app.add_url_rule(route, route[1:], auth_middleware(globals()[route[1:]]), methods=['GET', 'POST'])

`);
}

let username = ref('admin');
let password = ref('');
let username2 = ref('admin');
let password2 = ref('');

const run1 = () => {
   username.value = `admin' and '1'='1' --`;
   password.value = '12';
   value.value = 1;
};

const run2 = () => {
  username2.value = `admin' and '1'='1' --`;
  password2.value = '12';
  value2.value = 1;
};

onMounted(() => {
  console.log('mounted');
});
</script>

<style scoped lang="less">
.homeBox {
  display: flex;
  flex-direction: column;
  height: 100vh;
  padding: 20px;
  box-sizing: border-box;
}
.list {
  display: flex;
  flex: 1;
  .item {
    width: 50%;
    display: flex;
    flex-direction: column;
    .edit-container, .view-container {
      flex: 1;
      margin: 10px;
    }
  }
}
._top {
  margin: 20px;
}
.tips {
  margin: 20px 10px;
  padding: 20px;
  background-color: #fff;
  border-radius: 10px;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
  line-height: 1.6;
  color: #555;
  cursor: pointer;
  transition: background-color 0.3s, color 0.3s;
}
.tips:hover {
  background-color: #e9ecef;
}
pre {
  background-color: #fff;
  color: #333;
  padding: 20px;
  border-radius: 10px;
  overflow: auto;
  flex: 1;
  height: 100%;
  width: 100%;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
}
</style>