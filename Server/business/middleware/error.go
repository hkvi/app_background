package middleware

import (
	"net/http"

	"business/utils/hkvilog"

	"github.com/gin-gonic/gin"
)

// ErrorHandler 错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录panic错误
				hkvilog.Errorf("Panic recovered: %v", err)

				// 返回内部服务器错误
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "内部服务器错误",
				})
				c.Abort()
			}
		}()

		// 继续处理请求
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			// 获取最后一个错误
			err := c.Errors.Last()

			// 记录错误
			hkvilog.Errorf("Request error: %v", err)

			// 如果还没有响应，返回错误信息
			if !c.Writer.Written() {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "请求处理失败",
				})
			}
		}
	}
}
