package routes

import (
	"backend/initproject"
	"backend/middleware"
	"backend/vulnerabilities"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	router.Use(cors.New(config))

	router.GET("/init_sqllite3", initproject.Init_sqllite3)
	router.POST("/sql_injection_sqlite3_safe", vulnerabilities.Sql_injection_sqlite3_safe)
	router.POST("/sql_injection_sqlite3", vulnerabilities.Sql_injection_sqlite3)
	router.POST("/reflect_xss", vulnerabilities.ReflectXss)
	router.POST("/reflect_xss_safe", vulnerabilities.ReflectXssSafe)
	router.POST("/get_profile_unauthorized", vulnerabilities.Get_profile_unauthorized)
	// 保护需要验证的API
	auth := router.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		// 获取商品信息接口
		auth.GET("/products", vulnerabilities.GetProducts)

		// 购买商品接口
		auth.POST("/purchase", vulnerabilities.PurchaseProduct)
		auth.POST("/purchase_safe", vulnerabilities.PurchaseProduct_safe)
		// Get_profile_unauthorized
		auth.POST("/get_comments", vulnerabilities.Get_comments)
		auth.POST("/create_comments", vulnerabilities.Create_comments)
		auth.POST("/create_comments_safe", vulnerabilities.Create_comments_safe)
		auth.POST("/clear_comments", vulnerabilities.Clear_comments)
		auth.POST("/execute_command", vulnerabilities.ExecuteCommand)
		auth.POST("/execute_command_safe", vulnerabilities.ExecuteCommand_safe)
		auth.POST("/get_profile", vulnerabilities.Get_profile)
		auth.POST("/get_profile_safe", vulnerabilities.Get_profile_safe)
		auth.POST("/upload_file", vulnerabilities.UploadFile)
		auth.POST("/upload_file_safe", vulnerabilities.UploadFile_safe)
		auth.GET("/list_images", vulnerabilities.ListImages)
		auth.POST("/download", vulnerabilities.DownloadFile)

		auth.POST("/download_safe", vulnerabilities.DownloadFile_safe)

		auth.POST("/xxe", vulnerabilities.XML_parse)

	}
}
