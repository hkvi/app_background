package middleware

import (
	"time"

	"business/utils/hkvilog"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware 日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 计算处理时间
		duration := time.Since(startTime)

		// 获取请求信息
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		statusCode := c.Writer.Status()
		userAgent := c.Request.UserAgent()

		// 记录日志
		hkvilog.Infof("[%s] %s %s %d %v \"%s\"",
			clientIP,
			method,
			path,
			statusCode,
			duration,
			userAgent,
		)
	}
}
