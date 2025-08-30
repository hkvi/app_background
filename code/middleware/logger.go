package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware 日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// 自定义日志格式
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,                    // 客户端IP
			param.TimeStamp.Format(time.RFC1123), // 时间戳
			param.Method,                      // 请求方法
			param.Path,                        // 请求路径
			param.Request.Proto,               // 协议版本
			param.StatusCode,                  // 状态码
			param.Latency,                     // 响应时间
			param.Request.UserAgent(),         // 用户代理
			param.ErrorMessage,                // 错误信息
		)
	})
}
