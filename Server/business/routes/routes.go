package routes

import (
	"business/config"
	"business/handlers"
	"business/middleware"
	"business/utils/hkvilog"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置路由
func SetupRoutes(r *gin.Engine, cfg *config.Config) {
	// 使用中间件
	r.Use(middleware.LoggerMiddleware()) // 日志中间件
	r.Use(middleware.ErrorHandler())     // 错误处理中间件

	// 创建处理器实例
	userHandler := handlers.NewUserHandler()
	smsHandler, err := handlers.NewSMSHandler(&cfg.SMS)
	if err != nil {
		hkvilog.Errorf("创建短信处理器失败: %v", err)
		return
	}

	// API路由组
	api := r.Group("/api")
	{
		// 健康检查接口
		api.GET("/health", handlers.HealthCheck)

		// 认证相关接口
		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.Register) // 用户注册
			auth.POST("/login", userHandler.Login)       // 用户登录
		}

		// 短信相关接口
		sms := api.Group("/sms")
		{
			sms.POST("/send", smsHandler.SendSMS)   // 发送短信验证码
			sms.POST("/login", smsHandler.SMSLogin) // 短信验证码登录
		}

		// 其他业务接口可以在这里添加
		// 例如：用户信息管理、订单管理等
	}
}
