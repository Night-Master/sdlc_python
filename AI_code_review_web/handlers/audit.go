package handlers

import (
	"ai_code_review/utils"
	"archive/zip"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	auditMutex  sync.Mutex
	progressMap = make(map[string]int)
	totalFiles  = make(map[string]int)
	statusMap   = make(map[string]string)
	config      *utils.Config
	db          *sql.DB
)

// 初始化审计相关变量
func InitAudit(c *gin.Context, cfg *utils.Config, database *sql.DB) {
	config = cfg
	db = database
}

// GetAIPlatforms 获取AI平台列表
func GetAIPlatforms(c *gin.Context) {
	type AIPlatformResponse struct {
		Name string `json:"name"`
		// 你可以根据需要添加更多字段
	}

	var platforms []AIPlatformResponse
	for _, p := range config.AIPlatforms {
		platforms = append(platforms, AIPlatformResponse{
			Name: p.Name,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":       1,
		"ai_platforms": platforms,
	})
}

// StartAuditDir 启动目录审计任务
func StartAuditDir(c *gin.Context) {
	type AuditDirRequest struct {
		Path         string `json:"path" binding:"required"`
		Language     string `json:"language" binding:"required"` // 例如：go, java, py, js
		Platform     string `json:"platform" binding:"required"`
		ThreadNum    int    `json:"thread_num"`              // 可选，默认值
		WebFramework string `json:"web_framework,omitempty"` // 新增的 Web 框架字段
	}

	var req AuditDirRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 0, "message": "请求参数无效"})
		return
	}

	// 转义路径
	dir := html.EscapeString(req.Path)
	if dir == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": 0, "message": "路径不能为空"})
		return
	}

	// 检查路径是否存在且为目录
	info, err := os.Stat(dir)
	if os.IsNotExist(err) || !info.IsDir() {
		c.JSON(http.StatusBadRequest, gin.H{"status": 0, "message": "提供的路径不存在或不是目录"})
		return
	}

	// 查找AI平台
	var selectedPlatform *utils.AIPlatform
	for _, p := range config.AIPlatforms {
		if p.Name == req.Platform {
			selectedPlatform = &p
			break
		}
	}
	if selectedPlatform == nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 0, "message": "未找到指定的AI平台"})
		return
	}

	// 生成任务ID
	taskID := uuid.New().String()

	// 插入任务记录到数据库，包含扫描参数和 Web 框架
	startTime := time.Now().Format(time.RFC3339)
	_, err = db.Exec(`INSERT INTO tasks (task_id, status, start_time, language, thread_num, ai_platform_name, target_path, web_framework) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		taskID, "scanning", startTime, req.Language, getThreadNum(req.ThreadNum), selectedPlatform.Name, dir, req.WebFramework)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "message": "无法创建任务"})
		return
	}

	// 收集代码文件
	files, err := collectCodeFiles(dir, req.Language)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "message": "收集代码文件出错"})
		return
	}

	if len(files) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "未找到代码文件",
			"taskID":  taskID,
		})
		return
	}

	// 初始化任务进度
	auditMutex.Lock()
	progressMap[taskID] = 0
	totalFiles[taskID] = len(files)
	statusMap[taskID] = "scanning"
	auditMutex.Unlock()

	// 启动审计任务，传递 webFramework 参数
	go auditTask(taskID, selectedPlatform, files, getThreadNum(req.ThreadNum), req.WebFramework)

	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "审计任务已启动",
		"taskID":  taskID,
	})
}

func getThreadNum(threadNum int) int {
	if threadNum <= 0 {
		return config.MaxThreads // 假设在 config 中定义了 MaxThreads
	}
	return threadNum
}

// StartAuditZip 启动压缩包审计任务
func StartAuditZip(c *gin.Context) {
	type AuditZipRequest struct {
		Language     string `form:"language" binding:"required"` // 例如：go, java, py, js
		Platform     string `form:"platform" binding:"required"`
		ThreadNum    int    `form:"threads"`                 // 可选，默认值
		WebFramework string `form:"web_framework,omitempty"` // 新增的 Web 框架字段
	}

	var req AuditZipRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 0, "message": "请求参数无效"})
		return
	}

	// 处理文件上传
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 0, "message": "上传文件失败"})
		return
	}

	// 确保audit_code目录存在
	if err := os.MkdirAll("audit_code", 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "message": "无法创建保存目录"})
		return
	}

	// 保存zip文件
	zipFilePath := filepath.Join("audit_code", file.Filename)
	if err := c.SaveUploadedFile(file, zipFilePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "message": "文件保存失败"})
		return
	}

	// 解压文件
	extractedDir := strings.TrimSuffix(zipFilePath, ".zip")
	if err := unzipFile(zipFilePath, extractedDir); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "message": "解压失败"})
		return
	}

	// 删除zip文件
	if err := os.Remove(zipFilePath); err != nil {
		// 记录日志但不影响流程
		fmt.Printf("删除zip文件失败: %v\n", err)
	}

	// 检查解压后的目录
	info, err := os.Stat(extractedDir)
	if os.IsNotExist(err) || !info.IsDir() {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "message": "解压后的目录不存在或不是目录"})
		return
	}

	// 查找AI平台
	var selectedPlatform *utils.AIPlatform
	for _, p := range config.AIPlatforms {
		if p.Name == req.Platform {
			selectedPlatform = &p
			break
		}
	}
	if selectedPlatform == nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 0, "message": "未找到指定的AI平台"})
		return
	}

	// 生成任务ID
	taskID := uuid.New().String()

	// 插入任务记录到数据库，包含扫描参数和 Web 框架名称
	startTime := time.Now().Format(time.RFC3339)
	_, err = db.Exec(`INSERT INTO tasks (task_id, status, start_time, language, thread_num, ai_platform_name, target_path, web_framework) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		taskID, "scanning", startTime, req.Language, getThreadNum(req.ThreadNum), selectedPlatform.Name, zipFilePath, req.WebFramework)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "message": "无法创建任务"})
		return
	}

	// 收集代码文件
	files, err := collectCodeFiles(extractedDir, req.Language)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "message": "收集代码文件出错"})
		return
	}

	if len(files) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "未找到代码文件",
			"taskID":  taskID,
		})
		return
	}

	// 初始化任务进度
	auditMutex.Lock()
	progressMap[taskID] = 0
	totalFiles[taskID] = len(files)
	statusMap[taskID] = "scanning"
	auditMutex.Unlock()

	// 启动审计任务，传递 webFramework 参数
	go auditTask(taskID, selectedPlatform, files, getThreadNum(req.ThreadNum), req.WebFramework)

	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "审计任务已启动",
		"taskID":  taskID,
	})
}

// GetAuditProgress 获取审计任务的进度
func GetAuditProgress(c *gin.Context) {
	taskID := c.Param("taskID")

	auditMutex.Lock()
	scanned, ok1 := progressMap[taskID]
	total, ok2 := totalFiles[taskID]
	status, ok3 := statusMap[taskID]
	auditMutex.Unlock()

	if !ok1 || !ok2 || !ok3 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  0,
			"message": "任务未找到",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     1,
		"taskID":     taskID,
		"scanned":    scanned,
		"total":      total,
		"taskStatus": status,
	})
}

// GetAuditResult 获取审计任务的结果，新增返回 webFramework 信息
func GetAuditResult(c *gin.Context) {
	taskID := c.Param("taskID")

	// 查询任务状态
	var status string
	err := db.QueryRow("SELECT status FROM tasks WHERE task_id = ?", taskID).Scan(&status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "无法查询数据库或任务未完成",
		})
		return
	}

	if status != "completed" && status != "error" {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"taskID":  taskID,
			"result":  "任务尚未完成",
			"message": "请稍后再试",
		})
		return
	}

	// 查询结果
	safeTaskID := strings.ReplaceAll(taskID, "-", "")
	tableName := fmt.Sprintf("task_%s", safeTaskID)
	query := fmt.Sprintf("SELECT file_path, func_name, line, message, severity, cwe, code_snippet, scan_time FROM %s", tableName)

	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "无法查询数据库或任务未完成",
		})
		return
	}
	defer rows.Close()

	var results []utils.Vulnerability
	for rows.Next() {
		var vuln utils.Vulnerability
		var scanTime string
		err := rows.Scan(&vuln.FilePath, &vuln.FuncName, &vuln.Line, &vuln.Message, &vuln.Severity, &vuln.CWE, &vuln.CodeSnippet, &scanTime)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  0,
				"message": "数据库结果解析失败",
			})
			return
		}
		results = append(results, vuln)
	}

	if len(results) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  0,
			"message": "任务未找到或没有扫描结果",
		})
		return
	}

	// 获取 Web 框架信息
	var webFramework string
	err = db.QueryRow("SELECT web_framework FROM tasks WHERE task_id = ?", taskID).Scan(&webFramework)
	if err != nil {
		webFramework = "未知"
	}

	c.JSON(http.StatusOK, gin.H{
		"status":       1,
		"taskID":       taskID,
		"webFramework": webFramework, // 返回 Web 框架信息
		"result":       results,
	})
}

// GetAllAuditTasks 获取所有审计任务
func GetAllAuditTasks(c *gin.Context) {
	auditMutex.Lock()
	defer auditMutex.Unlock()

	var tasks []utils.Task
	dbTaskIDs := make(map[string]bool) // 存储数据库中已存在的 taskID

	// 查询数据库中已完成或错误的任务，包含 Web 框架字段
	rows, err := db.Query("SELECT task_id, status, start_time, end_time, language, thread_num, ai_platform_name, target_path, web_framework FROM tasks WHERE status = 'completed' OR status = 'error'")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  0,
			"message": "无法查询数据库",
		})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var task utils.Task
		err := rows.Scan(&task.TaskID, &task.Status, &task.StartTime, &task.EndTime, &task.Language, &task.ThreadNum, &task.AIPlatform, &task.TargetPath, &task.WebFramework)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  0,
				"message": "数据库结果解析失败",
			})
			return
		}
		tasks = append(tasks, task)
		dbTaskIDs[task.TaskID] = true
	}

	// 获取内存中的任务，并排除已在数据库中的任务
	for taskID, _ := range progressMap {
		if _, exists := dbTaskIDs[taskID]; exists {
			continue // 如果任务ID已在数据库中，则跳过
		}
		taskStatus, ok := statusMap[taskID]
		if !ok {
			taskStatus = "未知状态"
		}

		// 假设内存中任务的其他参数可以从其他地方获取，或者设置为默认值
		task := utils.Task{
			TaskID:     taskID,
			Status:     taskStatus,
			StartTime:  "", // 如果需要，可以添加内存中任务的开始时间
			Language:   "未知",
			ThreadNum:  0,
			AIPlatform: "未知",
			TargetPath: "未知",
			// WebFramework 默认为空
		}
		tasks = append(tasks, task)
	}

	// 返回所有任务信息
	c.JSON(http.StatusOK, gin.H{
		"status": 1,
		"tasks":  tasks,
	})
}

// auditTask 审计任务的核心逻辑
func auditTask(taskID string, platform *utils.AIPlatform, files []string, threadNum int, webFramework string) {
	// 设置默认线程数
	if threadNum <= 0 {
		threadNum = 5
	}

	// 使用通道和WaitGroup实现并发
	fileChan := make(chan string, len(files))
	resultChan := make(chan *utils.Vulnerability, len(files)*10) // 假设每个文件最多10个漏洞

	var wg sync.WaitGroup

	// 启动工作线程
	for i := 0; i < threadNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range fileChan {
				code, err := ioutil.ReadFile(file)
				if err != nil {
					fmt.Printf("读取文件出错：%v\n", err)
					continue
				}

				result, err := auditCode(*platform, string(code), webFramework)
				if err != nil {
					fmt.Printf("审计文件 %s 时出错：%v\n", file, err)
					continue
				}

				// 将文件路径添加到结果中
				for i := range result.Result {
					result.Result[i].FilePath = file
				}

				// 发送结果到结果通道
				for _, vuln := range result.Result {
					resultChan <- &vuln
				}

				// 更新进度
				auditMutex.Lock()
				progressMap[taskID]++
				auditMutex.Unlock()
			}
		}()
	}

	// 发送文件到文件通道
	for _, file := range files {
		fileChan <- file
	}
	close(fileChan)

	// 等待所有工作线程完成
	wg.Wait()
	close(resultChan)

	// 收集所有漏洞
	var allVulns []*utils.Vulnerability
	for vuln := range resultChan {
		allVulns = append(allVulns, vuln)
	}

	// 保存漏洞到数据库
	err := saveToDB(taskID, allVulns)
	if err != nil {
		fmt.Printf("无法保存漏洞报告到数据库: %v\n", err)
		updateTaskStatus(taskID, "error", time.Now().Format(time.RFC3339))
		return
	}

	// 更新任务状态为完成，并记录完成时间
	updateTaskStatus(taskID, "completed", time.Now().Format(time.RFC3339))
}

// collectCodeFiles 收集指定目录下的代码文件
func collectCodeFiles(path string, lang string) ([]string, error) {
	var files []string
	ext := getFileExtension(lang)
	if ext == "" {
		return files, fmt.Errorf("不支持的编程语言: %s", lang)
	}
	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ext) {
			files = append(files, p)
		}
		return nil
	})
	return files, err
}

// getFileExtension 根据编程语言获取文件扩展名
func getFileExtension(lang string) string {
	extensions := map[string]string{
		"go":     ".go",
		"java":   ".java",
		"py":     ".py",
		"python": ".py",
		"js":     ".js",
		// 可以添加更多语言
	}
	if ext, ok := extensions[strings.ToLower(lang)]; ok {
		return ext
	}
	return ""
}

// auditCode 调用AI平台进行代码审计
func auditCode(platform utils.AIPlatform, code string, webFramework string) (*utils.AuditResult, error) {
	// 构建提示词
	var promptBuilder strings.Builder

	if webFramework != "" {
		promptBuilder.WriteString(fmt.Sprintf("您正在审计使用 %s 框架的代码。\n", webFramework))
	}

	// 修改提示词，要求AI模型返回代码片段
	promptBuilder.WriteString(fmt.Sprintf("请对以下代码进行安全审计，并按照格式返回纯 JSON 结果，包括漏洞的代码片段。不要包含任何额外的字符、代码块或格式：%s\n代码：\n%s，如果无漏洞，则result为空", getResultFormat(), code))
	prompt := promptBuilder.String()

	// 构建请求体
	requestBody := map[string]interface{}{
		"model": platform.Model,
		"messages": []map[string]string{
			{"role": "system", "content": "你是一名高级代码安全审计工程师。"},
			{"role": "user", "content": prompt},
		},
		"stream": false,
	}
	requestData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", platform.Endpoint, bytes.NewBuffer(requestData))
	if err != nil {
		return nil, err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+platform.APIKey)

	// 发送请求
	client := &http.Client{
		Timeout: 6000 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API 返回错误状态码 %d: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("无法解析 API 响应为 JSON: %v", err)
	}

	content := extractContent(response)
	if content == "" {
		return nil, fmt.Errorf("响应中没有内容")
	}

	// 打印响应内容以调试
	fmt.Println("AI 返回的内容：", content)

	// 解析审计结果
	var auditResult utils.AuditResult
	err = json.Unmarshal([]byte(content), &auditResult)
	if err != nil {
		// 尝试使用字符串操作提取 JSON 后再次解析
		cleanedContent, cleanErr := extractJSON(content)
		if cleanErr != nil {
			return nil, fmt.Errorf("无法提取有效的 JSON: %v", cleanErr)
		}

		fmt.Println("提取并清理后的 JSON 内容：", cleanedContent)

		err = json.Unmarshal([]byte(cleanedContent), &auditResult)
		if err != nil {
			return nil, fmt.Errorf("解析提取的 JSON 失败: %v", err)
		}
	}

	return &auditResult, nil
}

// extractContent 提取并清理响应内容
func extractContent(response map[string]interface{}) string {
	choices, ok := response["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return ""
	}
	firstChoice, ok := choices[0].(map[string]interface{})
	if !ok {
		return ""
	}
	message, ok := firstChoice["message"].(map[string]interface{})
	if !ok {
		return ""
	}
	content, ok := message["content"].(string)
	if !ok {
		return ""
	}

	// 清理内容，移除反引号和多余的空白字符
	cleanedContent := strings.TrimSpace(content)
	cleanedContent = strings.TrimPrefix(cleanedContent, "json")
	cleanedContent = strings.TrimPrefix(cleanedContent, "```json")
	cleanedContent = strings.TrimSuffix(cleanedContent, "```")
	cleanedContent = strings.TrimSpace(cleanedContent)

	return cleanedContent
}

// extractJSON 使用字符串操作提取JSON
func extractJSON(content string) (string, error) {
	start := strings.Index(content, "{")
	end := strings.LastIndex(content, "}")
	if start == -1 || end == -1 || start > end {
		return "", fmt.Errorf("未找到有效的 JSON 对象")
	}
	return content[start : end+1], nil
}

// getResultFormat 获取结果格式提示
func getResultFormat() string {
	return `{"result":[{"func_name":"函数名称","line":行号,"message":"漏洞描述","severity":"严重程度(Critical/high/medium/low)","cwe":"CWE-编号","code_snippet":"漏洞代码片段"}]}`
}

// unzipFile 解压zip文件到指定目录
func unzipFile(zipFile, destDir string) error {
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(destDir, f.Name)

		// 防止Zip路径穿越攻击
		if !strings.HasPrefix(fpath, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("非法的文件路径: %s", fpath)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, f.Mode()); err != nil {
				return err
			}
			continue
		}

		// 确保父目录存在
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)

		// 关闭文件句柄
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}

// updateTaskStatus 更新任务状态和完成时间
func updateTaskStatus(taskID, newStatus, endTime string) {
	auditMutex.Lock()
	defer auditMutex.Unlock()
	statusMap[taskID] = newStatus
	closeTaskInDB(taskID, newStatus, endTime)
}

// closeTaskInDB 更新数据库中的任务记录
func closeTaskInDB(taskID, status, endTime string) {
	_, err := db.Exec(`UPDATE tasks SET status = ?, end_time = ? WHERE task_id = ?`, status, endTime, taskID)
	if err != nil {
		fmt.Printf("无法更新任务状态: %v\n", err)
	}
}

// saveToDB 将漏洞结果保存到SQLite数据库，并包含代码片段
func saveToDB(taskID string, vulns []*utils.Vulnerability) error {
	// 移除 taskID 中的破折号，确保表名安全
	safeTaskID := strings.ReplaceAll(taskID, "-", "")

	// 动态创建表，表名使用安全的 taskID
	tableName := fmt.Sprintf("task_%s", safeTaskID)
	createTableSQL := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            file_path TEXT,
            func_name TEXT,
            line INTEGER,
            message TEXT,
            severity TEXT,
            cwe TEXT,
            code_snippet TEXT,
            scan_time DATETIME DEFAULT CURRENT_TIMESTAMP
    );`, tableName)

	_, err := db.Exec(createTableSQL)
	if err != nil {
		return err
	}

	// 插入漏洞数据
	insertSQL := fmt.Sprintf(`INSERT INTO %s (file_path, func_name, line, message, severity, cwe, code_snippet) VALUES (?, ?, ?, ?, ?, ?, ?);`, tableName)

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(insertSQL)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, vuln := range vulns {
		// 直接使用AI模型返回的代码片段
		// 如果AI模型未返回代码片段，可以选择设置为空字符串
		snippet := vuln.CodeSnippet
		if snippet == "" {
			fmt.Printf("漏洞 %s 在文件 %s 的函数 %s 未提供代码片段\n", vuln.CWE, vuln.FilePath, vuln.FuncName)
		}

		_, err = stmt.Exec(vuln.FilePath, vuln.FuncName, vuln.Line, vuln.Message, vuln.Severity, vuln.CWE, snippet)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
