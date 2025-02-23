package vulnerabilities

import (
	"database/sql"
	"html"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

// 未对接受的参数进行处理
func ReflectXss(c *gin.Context) {
	type XssRequest struct {
		Input string `json:"input" binding:"required"`
	}

	var xssReq XssRequest
	if err := c.ShouldBindJSON(&xssReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 0, "message": "Invalid request"})
		return
	}

	input := xssReq.Input
	c.JSON(http.StatusOK, gin.H{"status": 1, "message": input})
}

// 对接受的参数使用 html.EscapeString转义
func ReflectXssSafe(c *gin.Context) {
	type XssRequest struct {
		Input string `json:"input" binding:"required"`
	}

	var xssReq XssRequest
	if err := c.ShouldBindJSON(&xssReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 0, "message": "Invalid request"})
		return
	}

	input := html.EscapeString(xssReq.Input)
	c.JSON(http.StatusOK, gin.H{"status": 1, "message": input})
}

func Get_comments(c *gin.Context) {
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "message": "Internal server error"})
		return
	}
	rows, err := db.Query("SELECT username, content FROM comments")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "message": "Internal server error"})
		return
	}
	defer rows.Close()

	type Comment struct {
		Username string `json:"username"`
		Content  string `json:"content"`
	}
	var comments []Comment
	for rows.Next() {
		var username, content string
		if err := rows.Scan(&username, &content); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "message": "Internal server error"})
			return
		}
		comments = append(comments, Comment{Username: username, Content: content})
	}

	c.JSON(http.StatusOK, gin.H{"status": 1, "message": "Comments fetched successfully", "data": comments})
}

// 未对接收上来的Content字段进行转义处理
func Create_comments(c *gin.Context) {
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "message": "Internal server error"})
		return
	}
	defer db.Close()

	type CommentRequest struct {
		Content string `json:"content" binding:"required"`
	}

	var commentReq CommentRequest
	if err := c.ShouldBindJSON(&commentReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 0, "message": "Invalid request"})
		return
	}

	// 从上下文中获取 username
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": 0, "message": "Unauthorized"})
		return
	}

	// 确保 username 是 string 类型
	usernameStr, ok := username.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "message": "Internal server error"})
		return
	}

	_, err = db.Exec("INSERT INTO comments (username, content) VALUES (?, ?)", usernameStr, commentReq.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "message": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": 1, "message": "Comment added successfully"})
}

// 对接收上来的Content字段进行HTMLEscapeString转义处理
func Create_comments_safe(c *gin.Context) {
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "message": "Internal server error"})
		return
	}
	defer db.Close()

	type CommentRequest struct {
		Content string `json:"content" binding:"required"`
	}

	var commentReq CommentRequest
	if err := c.ShouldBindJSON(&commentReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 0, "message": "Invalid request"})
		return
	}

	// 从上下文中获取 username
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": 0, "message": "Unauthorized"})
		return
	}

	// 确保 username 是 string 类型
	usernameStr, ok := username.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "message": "Internal server error"})
		return
	}

	// 转义用户输入的内容以防止XSS攻击
	escapedContent := html.EscapeString(commentReq.Content)

	_, err = db.Exec("INSERT INTO comments (username, content) VALUES (?, ?)", usernameStr, escapedContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "message": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": 1, "message": "Comment added successfully"})
}

func Clear_comments(c *gin.Context) {
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "Internal server error"})
		return
	}
	defer db.Close()

	// 获取 username
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "检测到未授权访问"})
		return
	}

	// 确保 username 是 string 类型
	usernameStr, ok := username.(string)
	if !ok || usernameStr != "admin" {
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "您不是管理员"})
		return
	}

	// 执行清空操作
	_, err = db.Exec("DELETE FROM comments")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": 1, "message": "All comments cleared successfully"})
}
