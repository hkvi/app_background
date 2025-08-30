package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HelloResponse Hello World响应结构体
type HelloResponse struct {
	Message string `json:"message"` // 消息
	UserID  int    `json:"user_id"` // 用户ID
	Username string `json:"username"` // 用户名
}

// Hello 需要认证的Hello World处理器
func Hello(c *gin.Context) {
	// 从上下文中获取用户信息（由认证中间件设置）
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "无法获取用户信息",
		})
		return
	}

	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "无法获取用户信息",
		})
		return
	}

	c.JSON(http.StatusOK, HelloResponse{
		Message:  "Hello World! 欢迎使用认证系统",
		UserID:   userID.(int),
		Username: username.(string),
	})
}
