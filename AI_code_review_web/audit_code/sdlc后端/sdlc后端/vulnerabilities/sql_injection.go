package vulnerabilities

import (
	"backend/utils"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func Sql_injection_sqlite3(c *gin.Context) {
	type LoginRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	var loginReq LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 0, "message": "Invalid request"})
		return
	}

	username := loginReq.Username
	password := loginReq.Password
	log.Println("username:", username)
	log.Println("password:", password)

	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "message": "Internal server error"})
		return
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT * FROM users WHERE username='%s' AND password='%s'", username, password)
	log.Println("query:", query)
	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "message": "Internal server error"})
		return
	}
	defer rows.Close()

	if rows.Next() {
		c.JSON(http.StatusOK, gin.H{"status": 1, "message": "Login successful!"})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "Login failed!"})
	}
}

// 参数化查询，它通过将用户输入的数据作为参数传递给查询，而不是直接将用户输入嵌入到查询字符串中，从而有效地防止SQL注入攻击
func Sql_injection_sqlite3_safe(c *gin.Context) {
	type LoginRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	var loginReq LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 0, "message": "Invalid request"})
		return
	}

	username := loginReq.Username
	password := loginReq.Password
	log.Println("username:", username)
	log.Println("password:", password)

	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "message": "Internal server error"})
		return
	}
	defer db.Close()

	query := "SELECT * FROM users WHERE username=? AND password=?"
	log.Println("query:", query)
	rows, err := db.Query(query, username, password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "message": "Internal server error"})
		return
	}
	defer rows.Close()

	if rows.Next() {
		token, err := utils.GenerateToken(username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "message": "Internal server error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": 1, "message": "Login successful!", "token": token, "username": username})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "Login failed!"})
	}
}
