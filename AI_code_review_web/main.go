package main

import (
	"ai_code_review/handlers"
	"ai_code_review/utils"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	db, err := utils.InitDB()
	if err != nil {
		log.Fatalf("无法初始化数据库: %v", err)
	}
	defer db.Close()

	// 加载配置文件
	config, err := utils.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("无法加载配置文件: %v", err)
	}

	// 初始化审计处理器
	handlers.Initialize(config, db)

	// 初始化Gin
	router := gin.Default()

	// 配置CORS中间件
	corsConfig := cors.New(cors.Config{
		AllowOrigins: []string{"*"}, // 允许的来源
		// AllowOrigins:     []string{"http://localhost:4000", "http://127.0.0.1:4000", "http://127.0.0.1:80", "http://localhost:80"}, // 允许的来源
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // 允许的HTTP方法
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // 允许的请求头
		ExposeHeaders:    []string{"Content-Length"},                          // 暴露的响应头
		AllowCredentials: true,                                                // 允许发送凭据（如cookies）
		MaxAge:           12 * time.Hour,                                      // 预检请求的缓存时间
	})

	router.Use(corsConfig)

	// 设置路由
	api := router.Group("/api")
	{
		api.POST("/audit/scan_start_dir", handlers.StartAuditDir)
		api.POST("/audit/scan_start_zip", handlers.StartAuditZip)
		api.GET("/audit/progress/:taskID", handlers.GetAuditProgress)
		api.GET("/audit/result/:taskID", handlers.GetAuditResult)
		api.GET("/audit/tasks", handlers.GetAllAuditTasks)
		api.GET("/ai_platforms", handlers.GetAIPlatforms) // 新增的API端点
	}

	// 启动服务器
	if err := router.Run(":2334"); err != nil {
		log.Fatalf("无法启动服务器: %v", err)
	}
}
