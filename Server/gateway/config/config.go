package config

import (
	"encoding/json"
	"os"
	"strconv"
)

// Config 网关配置结构体
type Config struct {
	Server      ServerConfig      `json:"server"`       // 服务器配置
	JWT         JWTConfig         `json:"jwt"`          // JWT配置
	Redis       RedisConfig       `json:"redis"`        // Redis配置
	BusinessAPI BusinessAPIConfig `json:"business_api"` // 业务服务API配置
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string `json:"port"` // 服务器端口
	Mode string `json:"mode"` // 运行模式
}

// JWTConfig JWT配置
type JWTConfig struct {
	AccessSecretKey  string `json:"access_secret_key"`  // Access Token密钥
	RefreshSecretKey string `json:"refresh_secret_key"` // Refresh Token密钥
	AccessExpire     int    `json:"access_expire"`      // Access Token过期时间（秒）
	RefreshExpire    int    `json:"refresh_expire"`     // Refresh Token过期时间（秒）
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `json:"host"`     // Redis主机
	Port     string `json:"port"`     // Redis端口
	Password string `json:"password"` // Redis密码
	DB       int    `json:"db"`       // Redis数据库编号
}

// BusinessAPIConfig 业务服务API配置
type BusinessAPIConfig struct {
	BaseURL string `json:"base_url"` // 业务服务基础URL
	Timeout int    `json:"timeout"`  // 请求超时时间（秒）
}

// LoadConfig 加载应用配置
func LoadConfig() *Config {
	// 默认配置
	config := &Config{
		Server: ServerConfig{
			Port: "8080",
			Mode: "debug",
		},
		JWT: JWTConfig{
			AccessSecretKey:  "your-access-secret-key",
			RefreshSecretKey: "your-refresh-secret-key",
			AccessExpire:     900,   // 15分钟
			RefreshExpire:    86400, // 24小时
		},
		Redis: RedisConfig{
			Host:     "localhost",
			Port:     "6379",
			Password: "",
			DB:       0,
		},
		BusinessAPI: BusinessAPIConfig{
			BaseURL: "http://localhost:8081",
			Timeout: 30,
		},
	}

	// 尝试从配置文件加载
	if err := loadFromFile(config); err != nil {
		// 如果配置文件不存在或读取失败，使用默认配置
	}

	// 从环境变量覆盖配置
	loadFromEnv(config)

	return config
}

// loadFromFile 从配置文件加载配置
func loadFromFile(config *Config) error {
	configFile := getConfigFile()
	if configFile == "" {
		return nil
	}

	file, err := os.Open(configFile)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(config)
}

// loadFromEnv 从环境变量加载配置
func loadFromEnv(config *Config) {
	// 服务器配置
	if port := os.Getenv("GATEWAY_PORT"); port != "" {
		config.Server.Port = port
	}
	if mode := os.Getenv("GATEWAY_MODE"); mode != "" {
		config.Server.Mode = mode
	}

	// JWT配置
	if accessKey := os.Getenv("JWT_ACCESS_SECRET_KEY"); accessKey != "" {
		config.JWT.AccessSecretKey = accessKey
	}
	if refreshKey := os.Getenv("JWT_REFRESH_SECRET_KEY"); refreshKey != "" {
		config.JWT.RefreshSecretKey = refreshKey
	}
	if accessExpireStr := os.Getenv("JWT_ACCESS_EXPIRE"); accessExpireStr != "" {
		if expire, err := strconv.Atoi(accessExpireStr); err == nil {
			config.JWT.AccessExpire = expire
		}
	}
	if refreshExpireStr := os.Getenv("JWT_REFRESH_EXPIRE"); refreshExpireStr != "" {
		if expire, err := strconv.Atoi(refreshExpireStr); err == nil {
			config.JWT.RefreshExpire = expire
		}
	}

	// Redis配置
	if host := os.Getenv("REDIS_HOST"); host != "" {
		config.Redis.Host = host
	}
	if port := os.Getenv("REDIS_PORT"); port != "" {
		config.Redis.Port = port
	}
	if password := os.Getenv("REDIS_PASSWORD"); password != "" {
		config.Redis.Password = password
	}
	if dbStr := os.Getenv("REDIS_DB"); dbStr != "" {
		if db, err := strconv.Atoi(dbStr); err == nil {
			config.Redis.DB = db
		}
	}

	// 业务服务API配置
	if baseURL := os.Getenv("BUSINESS_API_BASE_URL"); baseURL != "" {
		config.BusinessAPI.BaseURL = baseURL
	}
	if timeoutStr := os.Getenv("BUSINESS_API_TIMEOUT"); timeoutStr != "" {
		if timeout, err := strconv.Atoi(timeoutStr); err == nil {
			config.BusinessAPI.Timeout = timeout
		}
	}
}

// getConfigFile 获取配置文件路径
func getConfigFile() string {
	// 按优先级查找配置文件
	configFiles := []string{
		"gateway-config.json",
		"config/gateway-config.json",
		"./gateway-config.json",
		"./config/gateway-config.json",
	}

	for _, file := range configFiles {
		if _, err := os.Stat(file); err == nil {
			return file
		}
	}

	return ""
}
