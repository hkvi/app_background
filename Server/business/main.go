package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"business/cache"
	"business/config"
	"business/database"
	"business/routes"
	"business/utils/hkvilog"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 初始化数据库
	if err := database.InitDatabase(&cfg.Database); err != nil {
		hkvilog.Error("数据库初始化失败:", err)
		os.Exit(1)
	}
	defer database.CloseDatabase()

	// 创建数据库表
	if err := database.CreateTables(); err != nil {
		hkvilog.Error("创建数据库表失败:", err)
		os.Exit(1)
	}

	// 初始化Redis
	if err := cache.InitRedis(&cfg.Redis); err != nil {
		hkvilog.Error("Redis初始化失败:", err)
		os.Exit(1)
	}
	defer cache.CloseRedis()

	// 设置Gin运行模式
	gin.SetMode(cfg.Server.Mode)

	// 创建Gin引擎
	r := gin.Default()

	// 设置路由
	routes.SetupRoutes(r, cfg)

	// 构建服务器地址
	addr := fmt.Sprintf(":%s", cfg.Server.Port)

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// 启动服务器（在goroutine中）
	go func() {
		hkvilog.Infof("业务服务启动中，端口: %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			hkvilog.Error("业务服务启动失败:", err)
			os.Exit(1)
		}
	}()

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	hkvilog.Info("正在关闭业务服务...")

	// 设置5秒的超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		hkvilog.Error("业务服务强制关闭:", err)
		os.Exit(1)
	}

	hkvilog.Info("业务服务已退出")
}
