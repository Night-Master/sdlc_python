package vulnerabilities

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/exp/rand"
)

func GetProducts(c *gin.Context) {
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "message": "Internal server error"})
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, name, price FROM products")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "message": "Internal server error"})
		return
	}
	defer rows.Close()

	type Product struct {
		ID    int     `json:"id"`
		Name  string  `json:"name"`
		Price float64 `json:"price"`
	}
	var products []Product
	for rows.Next() {
		var id int
		var name string
		var price float64
		if err := rows.Scan(&id, &name, &price); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "message": "Internal server error"})
			return
		}
		products = append(products, Product{ID: id, Name: name, Price: price})
	}

	c.JSON(http.StatusOK, gin.H{"status": 1, "message": "Products fetched successfully", "data": products})
}

func PurchaseProduct(c *gin.Context) {
	var req struct {
		ProductID int `json:"productId"`
		Quantity  int `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "Invalid request"})
		return
	}

	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "Internal server error"})
		return
	}
	defer db.Close()

	// 从上下文中获取 username
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "Unauthorized"})
		return
	}

	// 确保 username 是 string 类型
	usernameStr, ok := username.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "Internal server error"})
		return
	}

	// 开始事务
	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "Internal server error"})
		return
	}

	// 获取用户ID
	var userID int
	err = tx.QueryRow("SELECT id FROM users WHERE username = ?", usernameStr).Scan(&userID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "User not found"})
		return
	}

	// 获取商品价格
	var productPrice float64
	err = tx.QueryRow("SELECT price FROM products WHERE id = ?", req.ProductID).Scan(&productPrice)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "Product not found"})
		return
	}

	// 计算总金额
	totalAmount := productPrice * float64(req.Quantity)

	// 检查用户余额是否足够
	var userBalance float64
	err = tx.QueryRow("SELECT balance FROM users WHERE id = ?", userID).Scan(&userBalance)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "User not found"})
		return
	}
	if userBalance < totalAmount {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "Insufficient balance", "balance": userBalance})
		return
	}

	// 更新用户余额
	_, err = tx.Exec("UPDATE users SET balance = balance - ? WHERE id = ?", totalAmount, userID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "Failed to update user balance"})
		return
	}

	// 插入订单记录
	orderID := generateOrderID()
	_, err = tx.Exec("INSERT INTO orders (user_id, product_id, order_id, amount, quantity) VALUES (?, ?, ?, ?, ?)", userID, req.ProductID, orderID, totalAmount, req.Quantity)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "Failed to insert order"})
		return
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "Failed to commit transaction"})
		return
	}

	// 获取更新后的用户余额
	err = db.QueryRow("SELECT balance FROM users WHERE id = ?", userID).Scan(&userBalance)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 1, "message": fmt.Sprintf("Successfully purchased %d of product ID %d", req.Quantity, req.ProductID), "balance": userBalance})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": 1, "message": fmt.Sprintf("Successfully purchased %d of product ID %d", req.Quantity, req.ProductID), "balance": userBalance})
}

func PurchaseProduct_safe(c *gin.Context) {
	var req struct {
		ProductID int `json:"productId"`
		Quantity  int `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "Invalid request"})
		return
	}

	// 检查数量是否大于零
	if req.Quantity <= 0 {
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "购买数量必须大于0"})
		return
	}

	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "Internal server error"})
		return
	}
	defer db.Close()

	// 从上下文中获取 username
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "Unauthorized"})
		return
	}

	// 确保 username 是 string 类型
	usernameStr, ok := username.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "Internal server error"})
		return
	}

	// 开始事务
	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "Internal server error"})
		return
	}

	// 获取用户ID
	var userID int
	err = tx.QueryRow("SELECT id FROM users WHERE username = ?", usernameStr).Scan(&userID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "User not found"})
		return
	}

	// 获取商品价格
	var productPrice float64
	err = tx.QueryRow("SELECT price FROM products WHERE id = ?", req.ProductID).Scan(&productPrice)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "Product not found"})
		return
	}

	// 计算总金额
	totalAmount := productPrice * float64(req.Quantity)

	// 检查用户余额是否足够
	var userBalance float64
	err = tx.QueryRow("SELECT balance FROM users WHERE id = ?", userID).Scan(&userBalance)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "User not found"})
		return
	}
	if userBalance < totalAmount {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "Insufficient balance", "balance": userBalance})
		return
	}

	// 更新用户余额
	_, err = tx.Exec("UPDATE users SET balance = balance - ? WHERE id = ?", totalAmount, userID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "Failed to update user balance"})
		return
	}

	// 插入订单记录
	orderID := generateOrderID()
	_, err = tx.Exec("INSERT INTO orders (user_id, product_id, order_id, amount, quantity) VALUES (?, ?, ?, ?, ?)", userID, req.ProductID, orderID, totalAmount, req.Quantity)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "Failed to insert order"})
		return
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 0, "message": "Failed to commit transaction"})
		return
	}

	// 获取更新后的用户余额
	err = db.QueryRow("SELECT balance FROM users WHERE id = ?", userID).Scan(&userBalance)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 1, "message": fmt.Sprintf("Successfully purchased %d of product ID %d", req.Quantity, req.ProductID), "balance": userBalance})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": 1, "message": fmt.Sprintf("Successfully purchased %d of product ID %d", req.Quantity, req.ProductID), "balance": userBalance})
}
func generateOrderID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", b)
}
