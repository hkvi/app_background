package cache

import (
	"context"
	"fmt"
	"time"

	"gateway/config"
	"gateway/utils/hkvilog"

	"github.com/go-redis/redis/v8"
)

// RedisClient Redis客户端
var RedisClient *redis.Client

// InitRedis 初始化Redis连接
func InitRedis(cfg *config.RedisConfig) error {
	// 创建Redis客户端
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("Redis连接失败: %v", err)
	}

	hkvilog.Info("Redis连接成功")
	return nil
}

// CloseRedis 关闭Redis连接
func CloseRedis() {
	if RedisClient != nil {
		if err := RedisClient.Close(); err != nil {
			hkvilog.Errorf("关闭Redis连接失败: %v", err)
		} else {
			hkvilog.Info("Redis连接已关闭")
		}
	}
}

// StoreRefreshToken 存储刷新令牌到Redis
func StoreRefreshToken(userID int, refreshToken string, expireTime time.Duration) error {
	ctx := context.Background()
	key := fmt.Sprintf("refresh_token:%d", userID)

	err := RedisClient.Set(ctx, key, refreshToken, expireTime).Err()
	if err != nil {
		return fmt.Errorf("存储刷新令牌失败: %v", err)
	}

	return nil
}

// GetRefreshToken 从Redis获取刷新令牌
func GetRefreshToken(userID int) (string, error) {
	ctx := context.Background()
	key := fmt.Sprintf("refresh_token:%d", userID)

	token, err := RedisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("刷新令牌不存在")
		}
		return "", fmt.Errorf("获取刷新令牌失败: %v", err)
	}

	return token, nil
}

// DeleteRefreshToken 从Redis删除刷新令牌
func DeleteRefreshToken(userID int) error {
	ctx := context.Background()
	key := fmt.Sprintf("refresh_token:%d", userID)

	err := RedisClient.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("删除刷新令牌失败: %v", err)
	}

	return nil
}

// BlacklistToken 将令牌加入黑名单
func BlacklistToken(tokenID string, expireTime time.Duration) error {
	ctx := context.Background()
	key := fmt.Sprintf("blacklist:%s", tokenID)

	err := RedisClient.Set(ctx, key, "1", expireTime).Err()
	if err != nil {
		return fmt.Errorf("添加令牌黑名单失败: %v", err)
	}

	return nil
}

// IsTokenBlacklisted 检查令牌是否在黑名单中
func IsTokenBlacklisted(tokenID string) (bool, error) {
	ctx := context.Background()
	key := fmt.Sprintf("blacklist:%s", tokenID)

	_, err := RedisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil // 不在黑名单中
		}
		return false, fmt.Errorf("检查令牌黑名单失败: %v", err)
	}

	return true, nil // 在黑名单中
}
