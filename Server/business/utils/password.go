package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 使用bcrypt对密码进行哈希
func HashPassword(password string) (string, error) {
	// 使用默认成本（10）进行密码哈希
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword 验证密码是否匹配
func CheckPassword(password, hash string) bool {
	// 比较密码和哈希值
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
