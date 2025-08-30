package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthResponse 健康检查响应结构体
type HealthResponse struct {
	Status  string `json:"status"`  // 状态
	Message string `json:"message"` // 消息
}

// HealthCheck 健康检查处理器
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status:  "ok",
		Message: "Server is running",
	})
}
