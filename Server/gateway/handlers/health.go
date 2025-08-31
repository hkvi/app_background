package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthCheck 健康检查处理器
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"service":   "gateway",
		"timestamp": time.Now().Unix(),
		"message":   "网关服务运行正常",
	})
}
