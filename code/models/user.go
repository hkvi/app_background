package models

import (
	"time"
)

// User 用户模型
type User struct {
	ID        int       `json:"id" db:"id"`               // 用户ID
	Username  string    `json:"username" db:"username"`   // 用户名
	Password  string    `json:"-" db:"password"`          // 密码（不返回给前端）
	CreatedAt time.Time `json:"created_at" db:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"` // 更新时间
}

// UserRegisterRequest 用户注册请求
type UserRegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"` // 用户名
	Password string `json:"password" binding:"required,min=6,max=100"` // 密码
}

// UserLoginRequest 用户登录请求
type UserLoginRequest struct {
	Username string `json:"username" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"` // 密码
}

// UserResponse 用户响应
type UserResponse struct {
	ID        int       `json:"id"`         // 用户ID
	Username  string    `json:"username"`   // 用户名
	CreatedAt time.Time `json:"created_at"` // 创建时间
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string       `json:"token"` // JWT token
	User  UserResponse `json:"user"`  // 用户信息
}
