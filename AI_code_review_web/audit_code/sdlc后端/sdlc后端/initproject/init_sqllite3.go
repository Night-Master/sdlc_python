package initproject

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func generateOrderID() string {
	rand.Seed(time.Now().UnixNano())
	const digits = "0123456789"
	b := make([]byte, 16)
	for i := range b {
		b[i] = digits[rand.Intn(len(digits))]
	}
	return string(b)
}

func Init_sqllite3(c *gin.Context) {
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		fmt.Println("Failed to connect to the database:", err)
		return
	}
	defer db.Close()

	// 创建 users 表
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		email TEXT NOT NULL,
		signature TEXT,
		avatar TEXT,
		birthdate TEXT,
		balance REAL
	)
	`)
	if err != nil {
		fmt.Println("Failed to create users table:", err)
		return
	}

	// 创建 products 表
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		price REAL NOT NULL
	)
	`)
	if err != nil {
		fmt.Println("Failed to create products table:", err)
		return
	}

	// 创建 orders 表
	_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS orders (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL,
	product_id INTEGER NOT NULL,
	order_id TEXT NOT NULL,
	amount REAL NOT NULL,
	quantity INTEGER NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users (id),
	FOREIGN KEY (product_id) REFERENCES products (id)
)
`)
	if err != nil {
		fmt.Println("Failed to create orders table:", err)
		return
	}

	// 创建 comments 表
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,
		content TEXT NOT NULL
	)
	`)
	if err != nil {
		fmt.Println("Failed to create comments table:", err)
		return
	}

	// 插入一些示例数据到 users 表
	_, err = db.Exec("INSERT INTO users (username, password,email, signature, avatar, birthdate, balance) VALUES (?, ?, ?, ?,?, ?, ?)", "user1", "hello", "123@qq.com", "Hello, I'm user1!", "avatar1.jpg", "1990-01-01", 100.0)
	if err != nil {
		fmt.Println("Failed to insert data into users table:", err)
		return
	}
	_, err = db.Exec("INSERT INTO users (username, password,email, signature, avatar, birthdate, balance) VALUES (?, ?, ?, ?,?, ?, ?)", "user2", "password1", "1223@qq.com", "Hello, I'm user2!", "avatar2.jpg", "1995-05-05", 200.0)
	if err != nil {
		fmt.Println("Failed to insert data into users table:", err)
		return
	}

	// 插入一些示例数据到 products 表
	_, err = db.Exec("INSERT INTO products (name, price) VALUES (?, ?)", "Product A", 50.0)
	if err != nil {
		fmt.Println("Failed to insert data into products table:", err)
		return
	}
	_, err = db.Exec("INSERT INTO products (name, price) VALUES (?, ?)", "Product B", 75.0)
	if err != nil {
		fmt.Println("Failed to insert data into products table:", err)
		return
	}

	// 插入一些示例数据到 orders 表
	user1ID := 1    // 假设 user1 的 ID 为 1
	user2ID := 2    // 假设 user2 的 ID 为 2
	product1ID := 1 // 假设 Product A 的 ID 为 1
	product2ID := 2 // 假设 Product B 的 ID 为 2
	_, err = db.Exec("INSERT INTO orders (user_id, product_id, order_id, amount, quantity) VALUES (?, ?, ?, ?, ?)", user1ID, product1ID, generateOrderID(), 50.0, 1)
	if err != nil {
		fmt.Println("Failed to insert data into orders table:", err)
		return
	}
	_, err = db.Exec("INSERT INTO orders (user_id, product_id, order_id, amount, quantity) VALUES (?, ?, ?, ?, ?)", user2ID, product2ID, generateOrderID(), 75.0, 1)
	if err != nil {
		fmt.Println("Failed to insert data into orders table:", err)
		return
	}

	// 插入一些示例数据到 comments 表
	_, err = db.Exec("INSERT INTO comments (username, content) VALUES (?, ?)", "user1", "hello")
	if err != nil {
		fmt.Println("Failed to insert data into comments table:", err)
		return
	}
	_, err = db.Exec("INSERT INTO comments (username, content) VALUES (?, ?)", "user2", "我要学渗透")
	if err != nil {
		fmt.Println("Failed to insert data into comments table:", err)
		return
	}

	fmt.Println("SQLite 数据库和表已成功创建，并插入了示例数据。")
}
