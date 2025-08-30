package middleware

import (
	"net/http"
	"login/cache"
	"login/utils/hkvilog"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware(limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取客户端IP
		clientIP := c.ClientIP()
		key := "rate_limit:" + clientIP

		// 增加计数
		count, err := cache.IncrementRateLimit(key, window)
		if err != nil {
			hkvilog.Errorf("限流检查失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "服务器内部错误",
			})
			c.Abort()
			return
		}

		// 检查是否超过限制
		if count > int64(limit) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// SMSRateLimitMiddleware 短信接口专用限流中间件
func SMSRateLimitMiddleware() gin.HandlerFunc {
	return RateLimitMiddleware(10, time.Minute) // 1分钟内最多10次请求
}
