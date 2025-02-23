package vulnerabilities

import (
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
)

func RemoteCodeExecution(c *gin.Context) {
	command := c.Query("command")
	if command == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Command is required"})
		return
	}

	// 不安全的代码：直接执行用户提供的命令
	output, err := exec.Command(command).CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"output": string(output)})
}
