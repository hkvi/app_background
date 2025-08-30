package models

import (
	"time"
)

// User 用户模型
type User struct {
	ID        int       `json:"id" db:"id"`                 // 用户ID
	Username  string    `json:"username" db:"username"`     // 用户名
	Password  string    `json:"-" db:"password"`            // 密码（不返回给前端）
	Phone     string    `json:"phone" db:"phone"`           // 手机号
	CreatedAt time.Time `json:"created_at" db:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"` // 更新时间
}

// UserRegisterRequest 用户注册请求
type UserRegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`  // 用户名
	Password string `json:"password" binding:"required,min=6,max=100"` // 密码
}

// UserLoginRequest 用户登录请求
type UserLoginRequest struct {
	Username string `json:"username" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"` // 密码
}

// SMSLoginRequest 短信验证码登录请求
type SMSLoginRequest struct {
	Phone string `json:"phone" binding:"required"` // 手机号
	Code  string `json:"code" binding:"required"`  // 验证码
}

// SendSMSRequest 发送短信验证码请求
type SendSMSRequest struct {
	Phone string `json:"phone" binding:"required"` // 手机号
}

// UserResponse 用户响应
type UserResponse struct {
	ID        int       `json:"id"`         // 用户ID
	Username  string    `json:"username"`   // 用户名
	Phone     string    `json:"phone"`      // 手机号
	CreatedAt time.Time `json:"created_at"` // 创建时间
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string       `json:"token"` // JWT token
	User  UserResponse `json:"user"`  // 用户信息
}

// SMSResponse 短信响应
type SMSResponse struct {
	Message string `json:"message"` // 响应消息
}
