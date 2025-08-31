package middleware

import (
	"gateway/config"
	"gateway/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "缺少认证token",
			})
			c.Abort()
			return
		}

		// 检查Authorization格式
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "认证格式错误",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 加载配置
		cfg := config.LoadConfig()

		// 解析访问令牌
		claims, err := utils.ValidateAccessToken(tokenString, cfg.JWT.AccessSecretKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		// 检查令牌是否在黑名单中 (暂时注释，需要token ID)
		// isBlacklisted, err := cache.IsTokenBlacklisted(claims.ID)
		// if err != nil {
		//	c.JSON(http.StatusInternalServerError, gin.H{
		//		"error": "令牌验证失败",
		//	})
		//	c.Abort()
		//	return
		// }

		// if isBlacklisted {
		//	c.JSON(http.StatusUnauthorized, gin.H{
		//		"error": "令牌已失效",
		//	})
		//	c.Abort()
		//	return
		// }

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		// c.Set("token_id", claims.ID)

		c.Next()
	}
}
