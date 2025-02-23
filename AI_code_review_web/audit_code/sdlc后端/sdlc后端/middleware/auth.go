package middleware

import (
	"net/http"

	"backend/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusOK, gin.H{"status": 0, "message": "您未登录，请勿未授权访问！"})
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"status": 0, "message": "登录凭证过期"})
			c.Abort()
			return
		}

		c.Set("username", claims.Username)
		c.Next()
	}
}
