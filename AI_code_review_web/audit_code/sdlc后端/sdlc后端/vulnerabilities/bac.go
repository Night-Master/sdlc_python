package vulnerabilities

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func Get_profile_unauthorized(c *gin.Context) {
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "message": "Internal server error"})
		return
	}
	defer db.Close()

	type ProfileRequest struct {
		Username string `json:"username" binding:"required"`
	}

	var profileReq ProfileRequest
	if err := c.ShouldBindJSON(&profileReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 0, "message": "Invalid request"})
		return
	}

	rows, err := db.Query("SELECT id, username, email, signature, avatar, birthdate, balance FROM users WHERE username = ?", profileReq.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "message": "Internal server error"})
		return
	}
	defer rows.Close()

	var user struct {
		ID        int     `json:"id"`
		Username  string  `json:"username"`
		Email     string  `json:"email"`
		Signature string  `json:"signature"`
		Avatar    string  `json:"avatar"`
		Birthdate string  `json:"birthdate"`
		Balance   float64 `json:"balance"`
	}

	if rows.Next() {
		err = rows.Scan(&user.ID, &user.Username, &user.Email, &user.Signature, &user.Avatar, &user.Birthdate, &user.Balance)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "message": "Internal server error"})
			return
		}
	} else {
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "User not found"})
		return
	}

	orderRows, err := db.Query("SELECT o.order_id, o.amount, o.quantity, p.name FROM orders o JOIN products p ON o.product_id = p.id WHERE o.user_id = ?", user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "message": "Internal server error"})
		return
	}
	defer orderRows.Close()

	var orders []struct {
		OrderID  string  `json:"order_id"`
		Amount   float64 `json:"amount"`
		Quantity int     `json:"quantity"`
		Name     string  `json:"name"`
	}

	for orderRows.Next() {
		var order struct {
			OrderID  string  `json:"order_id"`
			Amount   float64 `json:"amount"`
			Quantity int     `json:"quantity"`
			Name     string  `json:"name"`
		}
		err = orderRows.Scan(&order.OrderID, &order.Amount, &order.Quantity, &order.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "message": "Internal server error"})
			return
		}
		orders = append(orders, order)
	}

	c.JSON(http.StatusOK, gin.H{"status": 1, "message": "Profile fetched successfully", "data": gin.H{"user": user, "orders": orders}})
}

// 在获取用户信息时未检测当前用户名是否和token中的用户名一致，导致每个用户都可以查看数据库中所有人的用户信息，发生越权
func Get_profile(c *gin.Context) {
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "message": "Internal server error"})
		return
	}
	defer db.Close()

	type ProfileRequest struct {
		Username string `json:"username" binding:"required"`
	}

	var profileReq ProfileRequest
	if err := c.ShouldBindJSON(&profileReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 0, "message": "Invalid request"})
		return
	}

	rows, err := db.Query("SELECT id, username, email, signature, avatar, birthdate, balance FROM users WHERE username = ?", profileReq.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "message": "Internal server error"})
		return
	}
	defer rows.Close()

	var user struct {
		ID        int     `json:"id"`
		Username  string  `json:"username"`
		Email     string  `json:"email"`
		Signature string  `json:"signature"`
		Avatar    string  `json:"avatar"`
		Birthdate string  `json:"birthdate"`
		Balance   float64 `json:"balance"`
	}

	if rows.Next() {
		err = rows.Scan(&user.ID, &user.Username, &user.Email, &user.Signature, &user.Avatar, &user.Birthdate, &user.Balance)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "message": "Internal server error"})
			return
		}
	} else {
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "User not found"})
		return
	}

	orderRows, err := db.Query("SELECT o.order_id, o.amount, o.quantity, p.name FROM orders o JOIN products p ON o.product_id = p.id WHERE o.user_id = ?", user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "message": "Internal server error"})
		return
	}
	defer orderRows.Close()

	var orders []struct {
		OrderID  string  `json:"order_id"`
		Amount   float64 `json:"amount"`
		Quantity int     `json:"quantity"`
		Name     string  `json:"name"`
	}

	for orderRows.Next() {
		var order struct {
			OrderID  string  `json:"order_id"`
			Amount   float64 `json:"amount"`
			Quantity int     `json:"quantity"`
			Name     string  `json:"name"`
		}
		err = orderRows.Scan(&order.OrderID, &order.Amount, &order.Quantity, &order.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "message": "Internal server error"})
			return
		}
		orders = append(orders, order)
	}

	c.JSON(http.StatusOK, gin.H{"status": 1, "message": "Profile fetched successfully", "data": gin.H{"user": user, "orders": orders}})
}

// 在获取用户信息时检测当前用户名是否和token中的用户名一致，如果不一致则提示越权
func Get_profile_safe(c *gin.Context) {
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "message": "Internal server error"})
		return
	}
	defer db.Close()

	type ProfileRequest struct {
		Username string `json:"username" binding:"required"`
	}

	var profileReq ProfileRequest
	if err := c.ShouldBindJSON(&profileReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 0, "message": "Invalid request"})
		return
	}

	// 从上下文中获取 username
	tokenUsername, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "错误的登录凭证"})
		return
	}

	// 校验 token 中的用户名和请求中的用户名是否一致
	if tokenUsername != profileReq.Username {
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "您只能查看自己的账号信息"})
		return
	}

	rows, err := db.Query("SELECT id, username, email, signature, avatar, birthdate, balance FROM users WHERE username = ?", profileReq.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "message": "Internal server error"})
		return
	}
	defer rows.Close()

	var user struct {
		ID        int     `json:"id"`
		Username  string  `json:"username"`
		Email     string  `json:"email"`
		Signature string  `json:"signature"`
		Avatar    string  `json:"avatar"`
		Birthdate string  `json:"birthdate"`
		Balance   float64 `json:"balance"`
	}

	if rows.Next() {
		err = rows.Scan(&user.ID, &user.Username, &user.Email, &user.Signature, &user.Avatar, &user.Birthdate, &user.Balance)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "message": "Internal server error"})
			return
		}
	} else {
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "User not found"})
		return
	}

	orderRows, err := db.Query("SELECT o.order_id, o.amount, o.quantity, p.name FROM orders o JOIN products p ON o.product_id = p.id WHERE o.user_id = ?", user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "message": "Internal server error"})
		return
	}
	defer orderRows.Close()

	var orders []struct {
		OrderID  string  `json:"order_id"`
		Amount   float64 `json:"amount"`
		Quantity int     `json:"quantity"`
		Name     string  `json:"name"`
	}

	for orderRows.Next() {
		var order struct {
			OrderID  string  `json:"order_id"`
			Amount   float64 `json:"amount"`
			Quantity int     `json:"quantity"`
			Name     string  `json:"name"`
		}
		err = orderRows.Scan(&order.OrderID, &order.Amount, &order.Quantity, &order.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "message": "Internal server error"})
			return
		}
		orders = append(orders, order)
	}

	c.JSON(http.StatusOK, gin.H{"status": 1, "message": "Profile fetched successfully", "data": gin.H{"user": user, "orders": orders}})
}
