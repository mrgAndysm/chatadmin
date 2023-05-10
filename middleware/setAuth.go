package middleware

import (
	"chatgpt-go/db"
	"chatgpt-go/global"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func Cross() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("xxxcxcxcxc")
		// 设置跨域请求头
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// 处理 OPTIONS 请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
func SetAuthorizationHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := global.OpenAIKey // 使用从环境变量或输入中获取的API密钥
		c.Request.Header.Set("Authorization", "Bearer "+token)
		c.Next()
	}
}

func isNotEmptyString(s string) bool {
	return len(strings.TrimSpace(s)) > 0
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.Request.Header.Get("Authorization")
		token := strings.TrimPrefix(authorization, "Bearer ")
		if global.Config.System.AuthSecretKey == token {
			c.Next()
			return
		}

		response := struct {
			Status  string      `json:"status"`
			Message string      `json:"message"`
			Data    interface{} `json:"data"`
		}{
			Status:  "Unauthorized",
			Message: "Error: 无访问权限 | No access rights",
			Data:    nil,
		}
		if token == "" {
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}
		err := db.RedisDb.Get(token).Err()
		if err != nil {
			fmt.Println("redis-err:", err)
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		c.Next()
	}
}
