package vulnerabilities

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

// 未对文件路径和文件类型进行检查
func UploadFile(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "文件获取失败: " + err.Error()})
		return
	}

	// 保存文件到本地
	filename := filepath.Base(file.Filename)
	filePath := "./uploads/" + filename
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "文件保存失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": 1, "message": "文件上传成功", "filePath": filePath})
}

// 对上传的文件进行上传路径和文件类型进行检查
func UploadFile_safe(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "文件获取失败: " + err.Error()})
		return
	}

	// 检查文件类型是否为图片
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".bmp":  true,
		".webp": true,
		".tiff": true,
		".svg":  true,
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExtensions[ext] {
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "不允许的文件类型: " + ext})
		return
	}

	// 生成安全的文件名
	filename := filepath.Base(file.Filename)
	filename = strings.Replace(filename, "..", "", -1) // 移除所有 ".."
	filename = strings.Replace(filename, "/", "", -1)  // 移除所有 "/"
	filename = strings.Replace(filename, "\\", "", -1) // 移除所有 "\"

	// 确保文件保存路径在预期的目录下
	filePath := filepath.Join("./uploads", filename)

	// 保存文件到本地
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "文件保存失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": 1, "message": "文件上传成功", "filePath": filePath})
}

func ListImages(c *gin.Context) {
	files, err := ioutil.ReadDir("./uploads")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "读取文件列表失败: " + err.Error()})
		return
	}

	var images []string
	for _, file := range files {
		if !file.IsDir() {
			images = append(images, file.Name())
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": 1, "message": "获取图片列表成功", "images": images})
}

// 未对下载的文件的路径和文件类型进行检查
func DownloadFile(c *gin.Context) {
	var request struct {
		FileName string `json:"fileName"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 0, "message": "无效的请求: " + err.Error()})
		return
	}

	filePath := filepath.Join("./uploads", request.FileName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"status": 0, "message": "文件不存在"})
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+request.FileName)
	c.Header("Content-Type", "application/octet-stream")
	c.File(filePath)
}

// 对下载的文件的路径和文件类型进行检查
func DownloadFile_safe(c *gin.Context) {
	// 检查文件类型是否为图片
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".bmp":  true,
		".webp": true,
		".tiff": true,
		".svg":  true,
	}
	var request struct {
		FileName string `json:"fileName"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 0, "message": "无效的请求: " + err.Error()})
		return
	}

	// 验证文件名是否包含非法字符
	if strings.ContainsAny(request.FileName, "/\\") {
		c.JSON(http.StatusBadRequest, gin.H{"status": 0, "message": "文件名包含非法字符"})
		return
	}

	// 检查文件扩展名是否在允许的范围内
	ext := strings.ToLower(filepath.Ext(request.FileName))
	if !allowedExtensions[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"status": 0, "message": "不允许的文件类型"})
		return
	}

	filePath := filepath.Join("./uploads", request.FileName)

	// 确保文件路径在预期的目录下
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "message": "获取文件路径失败: " + err.Error()})
		return
	}

	expectedDir, err := filepath.Abs("./uploads")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "message": "获取预期目录失败: " + err.Error()})
		return
	}

	if !strings.HasPrefix(absPath, expectedDir) {
		c.JSON(http.StatusBadRequest, gin.H{"status": 0, "message": "文件路径非法"})
		return
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"status": 0, "message": "文件不存在"})
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+request.FileName)
	c.Header("Content-Type", "application/octet-stream")
	c.File(absPath)
}
