package cache

import (
	"context"
	"fmt"
	"login/config"
	"login/utils/hkvilog"
	"time"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

// InitRedis 初始化Redis连接
func InitRedis(cfg *config.RedisConfig) error {
	// 创建Redis客户端
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: 10, // 连接池大小
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := RedisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("Redis连接失败: %v", err)
	}

	hkvilog.Info("Redis连接成功")
	return nil
}

// CloseRedis 关闭Redis连接
func CloseRedis() {
	if RedisClient != nil {
		RedisClient.Close()
		hkvilog.Info("Redis连接已关闭")
	}
}

// SetSMSCode 设置短信验证码
func SetSMSCode(phone, code string, expiration time.Duration) error {
	ctx := context.Background()
	key := fmt.Sprintf("sms_code:%s", phone)
	return RedisClient.Set(ctx, key, code, expiration).Err()
}

// GetSMSCode 获取短信验证码
func GetSMSCode(phone string) (string, error) {
	ctx := context.Background()
	key := fmt.Sprintf("sms_code:%s", phone)
	return RedisClient.Get(ctx, key).Result()
}

// DeleteSMSCode 删除短信验证码
func DeleteSMSCode(phone string) error {
	ctx := context.Background()
	key := fmt.Sprintf("sms_code:%s", phone)
	return RedisClient.Del(ctx, key).Err()
}

// SetRateLimit 设置限流
func SetRateLimit(key string, expiration time.Duration) error {
	ctx := context.Background()
	return RedisClient.Set(ctx, key, "1", expiration).Err()
}

// GetRateLimit 获取限流状态
func GetRateLimit(key string) (int64, error) {
	ctx := context.Background()
	return RedisClient.Get(ctx, key).Int64()
}

// IncrementRateLimit 增加限流计数
func IncrementRateLimit(key string, expiration time.Duration) (int64, error) {
	ctx := context.Background()
	pipe := RedisClient.Pipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, expiration)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}
	return incr.Val(), nil
}
