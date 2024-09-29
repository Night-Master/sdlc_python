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
import { ElSelect, ElOption, ElButton,ElMessageBox } from 'element-plus';
// import { ElSelect, ElOption, ElButton, ElMessageBox } from 'element-plus';
import PreviewComponent from './components/PreviewComponent/index.vue';
import PreviewComponentII from './components/PreviewComponent2/index.vue';
import { useRouter } from 'vue-router'
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
const content = ref(`
//对用户输入的命令参数检测策略不严格，只要包含了指定命令即可，可以使用连接符绕过
func ExecuteCommand(c *gin.Context) {
	type CommandRequest struct {
		Command string \`json:"command" binding:"required"\`
	}

	var cmdReq CommandRequest
	if err := c.ShouldBindJSON(&cmdReq); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "无效请求"})
		return
	}

	command := cmdReq.Command

	// 检查命令是否符合指定命令
	validCommands := []string{"dir", "ls", "ipconfig", "ifconfig"}
	inputCommand := strings.Split(command, " ")[0]
	if !contains(validCommands, inputCommand) {
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "你输入的不是指定的命令"})
		return
	}

	// 根据操作系统执行命令
	var output []byte
	var err error
	if runtime.GOOS == "windows" {
		output, err = exec.Command("cmd", "/C", command).CombinedOutput()
	} else {
		output, err = exec.Command("sh", "-c", command).CombinedOutput()
	}

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "命令执行失败: " + err.Error()})
		return
	}

	// 返回命令执行结果
	c.JSON(http.StatusOK, gin.H{"status": 1, "message": strings.TrimSpace(string(output))})
}
`);
const content2 = ref(`
// 白名单策略：使用 validCommands 映射来定义允许的命令及其参数。例如，dir、ls、ipconfig 和 ifconfig 命令没有额外的参数。
// 命令参数验证：检查命令参数是否在白名单中。如果命令参数不在白名单中，则返回错误信息。
// 通过这种方式，我们可以确保只有白名单中的命令和参数可以被执行，从而大大降低命令注入的风险。
func ExecuteCommand_safe(c *gin.Context) {
	type CommandRequest struct {
		Command string \`json:"command" binding:"required"\`
	}

	var cmdReq CommandRequest
	if err := c.ShouldBindJSON(&cmdReq); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "无效请求"})
		return
	}

	command := cmdReq.Command

	// 检查命令是否符合指定命令
	validCommands := map[string][]string{
		"dir":      {},
		"ls":       {},
		"ipconfig": {},
		"ifconfig": {},
	}

	inputCommand := strings.Split(command, " ")[0]
	if _, exists := validCommands[inputCommand]; !exists {
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "你输入的不是指定的命令"})
		return
	}

	// 检查命令参数是否在白名单中
	cmdParts := strings.Split(command, " ")
	if len(cmdParts) > 1 {
		for _, arg := range cmdParts[1:] {
			if !contains(validCommands[inputCommand], arg) {
				c.JSON(http.StatusOK, gin.H{"status": 0, "message": "命令参数不在白名单中"})
				return
			}
		}
	}

	// 根据操作系统执行命令
	var output []byte
	var err error
	if runtime.GOOS == "windows" {
		output, err = exec.Command("cmd", "/C", command).CombinedOutput()
	} else {
		output, err = exec.Command("sh", "-c", command).CombinedOutput()
	}

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "命令执行失败: " + err.Error()})
		return
	}

	// 返回命令执行结果
	c.JSON(http.StatusOK, gin.H{"status": 1, "message": strings.TrimSpace(string(output))})
}
`);
let username = ref('admin');
let password = ref('');
let username2 = ref('admin');
let password2 = ref('');
const showReminder = () => {
  ElMessageBox.alert('购买负数商品试试', '关于如何利用', {
    confirmButtonText: '确定',
    type: 'warning',
    dangerouslyUseHTMLString: true,
    center: true,
    customClass: 'custom-message-box'
  });
  
};
const run1 = () => {

   value.value = 1;
  //  showReminder();
};

const run2 = () => {
  value2.value = 1;
  // showReminder();
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