package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT声明结构体
type Claims struct {
	UserID    int    `json:"user_id"`    // 用户ID
	Username  string `json:"username"`   // 用户名
	TokenType string `json:"token_type"` // token类型：access 或 refresh
	jwt.RegisteredClaims
}

// GenerateTokenPair 生成访问令牌和刷新令牌对
func GenerateTokenPair(userID int, username, accessSecretKey, refreshSecretKey string, accessExpire, refreshExpire int) (string, string, error) {
	// 生成访问令牌
	accessToken, err := GenerateToken(userID, username, "access", accessSecretKey, accessExpire)
	if err != nil {
		return "", "", err
	}

	// 生成刷新令牌
	refreshToken, err := GenerateToken(userID, username, "refresh", refreshSecretKey, refreshExpire)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// GenerateToken 生成JWT token
func GenerateToken(userID int, username, tokenType, secretKey string, expireTime int) (string, error) {
	// 创建声明
	claims := Claims{
		UserID:    userID,
		Username:  username,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireTime) * time.Second)), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                                              // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                                              // 生效时间
		},
	}

	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名token
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken 解析JWT token
func ParseToken(tokenString, secretKey string) (*Claims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("无效的签名方法")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// 验证token并获取声明
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的token")
}

// ValidateAccessToken 验证访问令牌
func ValidateAccessToken(tokenString, secretKey string) (*Claims, error) {
	claims, err := ParseToken(tokenString, secretKey)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != "access" {
		return nil, errors.New("无效的访问令牌类型")
	}

	return claims, nil
}

// ValidateRefreshToken 验证刷新令牌
func ValidateRefreshToken(tokenString, secretKey string) (*Claims, error) {
	claims, err := ParseToken(tokenString, secretKey)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != "refresh" {
		return nil, errors.New("无效的刷新令牌类型")
	}

	return claims, nil
}
