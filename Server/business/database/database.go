package database

import (
	"business/config"
	"business/utils/hkvilog"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// InitDatabase 初始化数据库连接
func InitDatabase(cfg *config.DatabaseConfig) error {
	// 首先连接到MySQL服务器（不指定数据库）
	dsnWithoutDB := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	)

	// 连接MySQL服务器
	tempDB, err := sql.Open("mysql", dsnWithoutDB)
	if err != nil {
		return fmt.Errorf("连接MySQL服务器失败: %v", err)
	}
	defer tempDB.Close()

	// 测试连接
	if err := tempDB.Ping(); err != nil {
		return fmt.Errorf("MySQL服务器ping失败: %v", err)
	}

	// 检查数据库是否存在，如果不存在则创建
	if err := ensureDatabaseExists(tempDB, cfg.DBName); err != nil {
		return fmt.Errorf("确保数据库存在失败: %v", err)
	}

	// 构建包含数据库名称的连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	// 连接到指定数据库
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("数据库连接失败: %v", err)
	}

	// 测试数据库连接
	if err := db.Ping(); err != nil {
		return fmt.Errorf("数据库ping失败: %v", err)
	}

	// 设置连接池参数
	db.SetMaxOpenConns(25) // 最大连接数
	db.SetMaxIdleConns(5)  // 最大空闲连接数

	DB = db
	hkvilog.Info("数据库连接成功")
	return nil
}

// ensureDatabaseExists 确保数据库存在，如果不存在则创建
func ensureDatabaseExists(db *sql.DB, dbName string) error {
	// 检查数据库是否存在
	var exists int
	query := "SELECT COUNT(*) FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?"
	err := db.QueryRow(query, dbName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("检查数据库是否存在失败: %v", err)
	}

	if exists == 0 {
		// 数据库不存在，创建它
		createQuery := fmt.Sprintf("CREATE DATABASE `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", dbName)
		_, err := db.Exec(createQuery)
		if err != nil {
			return fmt.Errorf("创建数据库失败: %v", err)
		}
		hkvilog.Infof("数据库 '%s' 创建成功", dbName)
	} else {
		hkvilog.Infof("数据库 '%s' 已存在", dbName)
	}

	return nil
}

// CloseDatabase 关闭数据库连接
func CloseDatabase() {
	if DB != nil {
		DB.Close()
		hkvilog.Info("数据库连接已关闭")
	}
}

// CreateTables 创建数据库表
func CreateTables() error {
	// 创建用户表
	createUserTable := `
	CREATE TABLE IF NOT EXISTS users (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(50)  NULL,
		password VARCHAR(255) NULL,
		phone VARCHAR(20) UNIQUE  NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		INDEX idx_username (username),
		INDEX idx_phone (phone)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`

	_, err := DB.Exec(createUserTable)
	if err != nil {
		return fmt.Errorf("创建用户表失败: %v", err)
	}

	// 创建验证码表（用于记录历史，实际验证码存储在Redis中）
	createSMSTable := `
	CREATE TABLE IF NOT EXISTS sms_codes (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		phone VARCHAR(20) NOT NULL,
		code VARCHAR(10) NOT NULL,
		type ENUM('login', 'register', 'reset') NOT NULL,
		used BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		expired_at TIMESTAMP NOT NULL,
		INDEX idx_phone_type (phone, type),
		INDEX idx_expired_at (expired_at)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`

	_, err = DB.Exec(createSMSTable)
	if err != nil {
		return fmt.Errorf("创建短信验证码表失败: %v", err)
	}

	hkvilog.Info("数据库表创建成功")
	return nil
}
