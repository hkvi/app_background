package routes

import (
	"gateway/config"
	"gateway/handlers"
	"gateway/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置路由
func SetupRoutes(r *gin.Engine, cfg *config.Config) {
	// 使用中间件
	r.Use(middleware.CORSMiddleware())   // CORS中间件
	r.Use(middleware.LoggerMiddleware()) // 日志中间件
	r.Use(middleware.ErrorHandler())     // 错误处理中间件

	// 创建处理器实例
	authHandler := handlers.NewAuthHandler(cfg)

	// API路由组
	api := r.Group("/api")
	{
		// 健康检查接口
		api.GET("/health", handlers.HealthCheck)

		// 认证相关接口（无需认证）
		auth := api.Group("/auth")
		auth.Use(middleware.LoginRateLimitMiddleware()) // 登录限流
		{
			auth.POST("/login", authHandler.Login)          // 用户登录
			auth.POST("/register", authHandler.Register)    // 用户注册
			auth.POST("/sms/send", authHandler.SendSMS)     // 发送短信验证码
			auth.POST("/sms/login", authHandler.SMSLogin)   // 短信验证码登录
			auth.POST("/refresh", authHandler.RefreshToken) // 刷新令牌
		}

		// 需要认证的接口
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware()) // 使用认证中间件
		{
			protected.POST("/auth/logout", authHandler.Logout) // 用户退出

			// 代理到业务服务的接口
			business := protected.Group("/business")
			{
				// 这里可以添加需要代理到业务服务的路由
				business.Any("/*path", handlers.ProxyToBusiness)
			}
		}
	}
}
