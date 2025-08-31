package handlers

import (
	"business/config"
	"business/models"
	"business/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SMSHandler 短信处理器
type SMSHandler struct {
	smsService  *services.SMSService
	userService *services.UserService
}

// NewSMSHandler 创建短信处理器实例
func NewSMSHandler(cfg *config.SMSConfig) (*SMSHandler, error) {
	smsService, err := services.NewSMSService(cfg)
	if err != nil {
		return nil, err
	}

	return &SMSHandler{
		smsService:  smsService,
		userService: services.NewUserService(),
	}, nil
}

// SendSMS 发送短信验证码
func (h *SMSHandler) SendSMS(c *gin.Context) {
	var req models.SendSMSRequest

	// 绑定请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 验证手机号格式
	if err := h.smsService.ValidatePhoneNumber(req.Phone); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 获取客户端IP
	clientIP := c.ClientIP()

	// 检查限流
	if err := h.smsService.CheckRateLimit(req.Phone, clientIP); err != nil {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 发送短信验证码
	code, err := h.smsService.SendSMSCode(req.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 在开发环境下返回验证码（生产环境应该注释掉）
	if gin.Mode() == gin.DebugMode {
		c.JSON(http.StatusOK, gin.H{
			"message": "验证码发送成功",
			"code":    code, // 仅开发环境返回验证码
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "验证码发送成功",
		})
	}
}

// SMSLogin 短信验证码登录
func (h *SMSHandler) SMSLogin(c *gin.Context) {
	var req models.SMSLoginRequest

	// 绑定请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 验证手机号格式
	if err := h.smsService.ValidatePhoneNumber(req.Phone); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 执行短信登录
	response, err := h.userService.LoginBySMS(req.Phone, req.Code, h.smsService)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"data":    response,
	})
}
