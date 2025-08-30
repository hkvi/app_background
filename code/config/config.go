package config

import (
	"encoding/json"
	"os"
	"strconv"
)

// Config 应用配置结构体
type Config struct {
	Server   ServerConfig   `json:"server"`   // 服务器配置
	Database DatabaseConfig `json:"database"` // 数据库配置
	JWT      JWTConfig      `json:"jwt"`      // JWT配置
	Redis    RedisConfig    `json:"redis"`    // Redis配置
	SMS      SMSConfig      `json:"sms"`      // 短信服务配置
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string `json:"port"` // 服务器端口
	Mode string `json:"mode"` // 运行模式
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `json:"host"`     // 数据库主机
	Port     string `json:"port"`     // 数据库端口
	Username string `json:"username"` // 数据库用户名
	Password string `json:"password"` // 数据库密码
	DBName   string `json:"dbname"`   // 数据库名称
}

// JWTConfig JWT配置
type JWTConfig struct {
	SecretKey string `json:"secret_key"` // JWT密钥
	Expire    int    `json:"expire"`     // 过期时间（秒）
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `json:"host"`     // Redis主机
	Port     string `json:"port"`     // Redis端口
	Password string `json:"password"` // Redis密码
	DB       int    `json:"db"`       // Redis数据库编号
}

// SMSConfig 短信服务配置
type SMSConfig struct {
	AccessKeyID     string `json:"access_key_id"`     // 阿里云AccessKey ID
	AccessKeySecret string `json:"access_key_secret"` // 阿里云AccessKey Secret
	SignName        string `json:"sign_name"`         // 短信签名
	TemplateCode    string `json:"template_code"`     // 短信模板代码
	RegionID        string `json:"region_id"`         // 地域ID
}

// LoadConfig 加载应用配置
func LoadConfig() *Config {
	// 默认配置
	config := &Config{
		Server: ServerConfig{
			Port: "8080",
			Mode: "debug",
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     "3306",
			Username: "root",
			Password: "password",
			DBName:   "login_db",
		},
		JWT: JWTConfig{
			SecretKey: "your-secret-key",
			Expire:    3600,
		},
		Redis: RedisConfig{
			Host:     "localhost",
			Port:     "6379",
			Password: "",
			DB:       0,
		},
		SMS: SMSConfig{
			AccessKeyID:     "your-access-key-id",
			AccessKeySecret: "your-access-key-secret",
			SignName:        "your-sign-name",
			TemplateCode:    "your-template-code",
			RegionID:        "cn-hangzhou",
		},
	}

	// 尝试从配置文件加载
	if err := loadFromFile(config); err != nil {
		// 如果配置文件不存在或读取失败，使用默认配置
		// 这里可以添加日志记录
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
	if port := os.Getenv("SERVER_PORT"); port != "" {
		config.Server.Port = port
	}
	if mode := os.Getenv("SERVER_MODE"); mode != "" {
		config.Server.Mode = mode
	}

	// 数据库配置
	if host := os.Getenv("DB_HOST"); host != "" {
		config.Database.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		config.Database.Port = port
	}
	if username := os.Getenv("DB_USERNAME"); username != "" {
		config.Database.Username = username
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		config.Database.Password = password
	}
	if dbname := os.Getenv("DB_NAME"); dbname != "" {
		config.Database.DBName = dbname
	}

	// JWT配置
	if secretKey := os.Getenv("JWT_SECRET_KEY"); secretKey != "" {
		config.JWT.SecretKey = secretKey
	}
	if expireStr := os.Getenv("JWT_EXPIRE"); expireStr != "" {
		if expire, err := strconv.Atoi(expireStr); err == nil {
			config.JWT.Expire = expire
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

	// SMS配置
	if accessKeyID := os.Getenv("SMS_ACCESS_KEY_ID"); accessKeyID != "" {
		config.SMS.AccessKeyID = accessKeyID
	}
	if accessKeySecret := os.Getenv("SMS_ACCESS_KEY_SECRET"); accessKeySecret != "" {
		config.SMS.AccessKeySecret = accessKeySecret
	}
	if signName := os.Getenv("SMS_SIGN_NAME"); signName != "" {
		config.SMS.SignName = signName
	}
	if templateCode := os.Getenv("SMS_TEMPLATE_CODE"); templateCode != "" {
		config.SMS.TemplateCode = templateCode
	}
	if regionID := os.Getenv("SMS_REGION_ID"); regionID != "" {
		config.SMS.RegionID = regionID
	}

	// SMS配置
	if accessKeyID := os.Getenv("SMS_ACCESS_KEY_ID"); accessKeyID != "" {
		config.SMS.AccessKeyID = accessKeyID
	}
	if accessKeySecret := os.Getenv("SMS_ACCESS_KEY_SECRET"); accessKeySecret != "" {
		config.SMS.AccessKeySecret = accessKeySecret
	}
	if signName := os.Getenv("SMS_SIGN_NAME"); signName != "" {
		config.SMS.SignName = signName
	}
	if templateCode := os.Getenv("SMS_TEMPLATE_CODE"); templateCode != "" {
		config.SMS.TemplateCode = templateCode
	}
	if regionID := os.Getenv("SMS_REGION_ID"); regionID != "" {
		config.SMS.RegionID = regionID
	}
}

// getConfigFile 获取配置文件路径
func getConfigFile() string {
	// 按优先级查找配置文件
	configFiles := []string{
		"config.json",
		"config/config.json",
		"./config.json",
		"./config/config.json",
	}

	for _, file := range configFiles {
		if _, err := os.Stat(file); err == nil {
			return file
		}
	}

	return ""
}
