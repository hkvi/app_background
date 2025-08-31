package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"gateway/cache"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// RateLimitMiddleware 通用限流中间件
func RateLimitMiddleware(keyPrefix string, maxRequests int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取客户端IP
		clientIP := c.ClientIP()

		// 构建Redis键
		key := fmt.Sprintf("%s:%s", keyPrefix, clientIP)

		// 检查限流
		if !checkRateLimit(key, maxRequests, window) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// SMSRateLimitMiddleware 短信限流中间件
func SMSRateLimitMiddleware() gin.HandlerFunc {
	return RateLimitMiddleware("sms_rate_limit", 5, time.Minute*1) // 1分钟内最多5次
}

// LoginRateLimitMiddleware 登录限流中间件
func LoginRateLimitMiddleware() gin.HandlerFunc {
	return RateLimitMiddleware("login_rate_limit", 10, time.Minute*5) // 5分钟内最多10次
}

// checkRateLimit 检查限流
func checkRateLimit(key string, maxRequests int, window time.Duration) bool {
	ctx := context.Background()

	// 使用滑动窗口算法
	now := time.Now().Unix()
	windowStart := now - int64(window.Seconds())

	// 删除窗口外的记录
	cache.RedisClient.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(windowStart, 10))

	// 获取当前窗口内的请求数
	count, err := cache.RedisClient.ZCard(ctx, key).Result()
	if err != nil {
		// 如果Redis出错，允许请求通过（降级处理）
		return true
	}

	// 检查是否超过限制
	if int(count) >= maxRequests {
		return false
	}

	// 添加当前请求
	cache.RedisClient.ZAdd(ctx, key, &redis.Z{
		Score:  float64(now),
		Member: fmt.Sprintf("%d", now),
	})

	// 设置键的过期时间
	cache.RedisClient.Expire(ctx, key, window)

	return true
}
