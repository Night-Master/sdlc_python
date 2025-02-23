package vulnerabilities

import (
	"net/http"
	"os/exec"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

// 对用户输入的命令参数检测策略不严格，只要包含了指定命令即可，可以使用连接符绕过
func ExecuteCommand(c *gin.Context) {
	type CommandRequest struct {
		Command string `json:"command" binding:"required"`
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

// 白名单策略：使用 validCommands 映射来定义允许的命令及其参数。例如，dir、ls、ipconfig 和 ifconfig 命令没有额外的参数。

// 命令参数验证：检查命令参数是否在白名单中。如果命令参数不在白名单中，则返回错误信息。

// 通过这种方式，我们可以确保只有白名单中的命令和参数可以被执行，从而大大降低命令注入的风险。
func ExecuteCommand_safe(c *gin.Context) {
	type CommandRequest struct {
		Command string `json:"command" binding:"required"`
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
