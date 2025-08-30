package routes

import (
	"login/config"
	"login/handlers"
	"login/middleware"
	"login/utils/hkvilog"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置路由
func SetupRoutes(r *gin.Engine, cfg *config.Config) {
	// 使用中间件
	r.Use(middleware.CORSMiddleware())    // CORS中间件
	r.Use(middleware.LoggerMiddleware())  // 日志中间件
	r.Use(middleware.ErrorHandler())      // 错误处理中间件

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
		
		// 用户认证相关接口（无需认证）
		api.POST("/register", userHandler.Register) // 用户注册
		api.POST("/login", userHandler.Login)       // 用户登录
		
		// 短信相关接口（无需认证，但有限流）
		sms := api.Group("/sms")
		sms.Use(middleware.SMSRateLimitMiddleware()) // 短信接口限流
		{
			sms.POST("/send", smsHandler.SendSMS)     // 发送短信验证码
			sms.POST("/login", smsHandler.SMSLogin)   // 短信验证码登录
		}
		
		// 需要认证的接口
		auth := api.Group("")
		auth.Use(middleware.AuthMiddleware()) // 使用认证中间件
		{
			auth.GET("/hello", handlers.Hello)           // Hello World接口
			auth.POST("/logout", userHandler.Logout)     // 用户退出
		}
	}
}
