package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWT 演示程序
func main() {
	fmt.Println("=== JWT 原理演示 ===\n")

	// 1. 演示 JWT 结构
	demoJWTStructure()

	// 2. 演示 JWT 生成过程
	demoJWTGeneration()

	// 3. 演示 JWT 验证过程
	demoJWTValidation()

	// 4. 演示 JWT 安全性
	demoJWTSecurity()
}

// 演示 JWT 结构
func demoJWTStructure() {
	fmt.Println("1. JWT 结构演示")
	fmt.Println("JWT 由三部分组成：Header.Payload.Signature\n")

	// Header
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}
	headerJSON, _ := json.Marshal(header)
	headerBase64 := base64.RawURLEncoding.EncodeToString(headerJSON)

	// Payload
	payload := map[string]interface{}{
		"user_id":  123,
		"username": "john_doe",
		"exp":      time.Now().Add(time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}
	payloadJSON, _ := json.Marshal(payload)
	payloadBase64 := base64.RawURLEncoding.EncodeToString(payloadJSON)

	fmt.Printf("Header (JSON): %s\n", string(headerJSON))
	fmt.Printf("Header (Base64): %s\n", headerBase64)
	fmt.Printf("Payload (JSON): %s\n", string(payloadJSON))
	fmt.Printf("Payload (Base64): %s\n", payloadBase64)
	fmt.Println("Signature: [使用密钥签名的结果]")
	fmt.Printf("完整 JWT: %s.%s.[SIGNATURE]\n\n", headerBase64, payloadBase64)
}

// 演示 JWT 生成过程
func demoJWTGeneration() {
	fmt.Println("2. JWT 生成过程演示")

	// 创建 Claims
	claims := jwt.MapClaims{
		"user_id":  123,
		"username": "john_doe",
		"role":     "user",
		"exp":      time.Now().Add(time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}

	// 创建 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名
	secretKey := "your-secret-key"
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		fmt.Printf("生成 JWT 失败: %v\n", err)
		return
	}

	fmt.Printf("生成的 JWT: %s\n", tokenString)
	fmt.Printf("JWT 长度: %d 字符\n\n", len(tokenString))

	// 解析 JWT 结构
	parts := strings.Split(tokenString, ".")
	if len(parts) == 3 {
		// 解码 Header
		headerBytes, _ := base64.RawURLEncoding.DecodeString(parts[0])
		var header map[string]interface{}
		json.Unmarshal(headerBytes, &header)
		fmt.Printf("解析的 Header: %+v\n", header)

		// 解码 Payload
		payloadBytes, _ := base64.RawURLEncoding.DecodeString(parts[1])
		var payload map[string]interface{}
		json.Unmarshal(payloadBytes, &payload)
		fmt.Printf("解析的 Payload: %+v\n", payload)
		fmt.Printf("签名: %s\n\n", parts[2])
	}
}

// 演示 JWT 验证过程
func demoJWTValidation() {
	fmt.Println("3. JWT 验证过程演示")

	// 生成一个测试 Token
	claims := jwt.MapClaims{
		"user_id":  456,
		"username": "alice",
		"exp":      time.Now().Add(time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := "your-secret-key"
	tokenString, _ := token.SignedString([]byte(secretKey))

	fmt.Printf("测试 Token: %s\n", tokenString)

	// 验证 Token
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		fmt.Printf("Token 验证失败: %v\n", err)
		return
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		fmt.Printf("✅ Token 验证成功!\n")
		fmt.Printf("用户ID: %v\n", claims["user_id"])
		fmt.Printf("用户名: %v\n", claims["username"])
		fmt.Printf("过期时间: %v\n", time.Unix(int64(claims["exp"].(float64)), 0))
	} else {
		fmt.Println("❌ Token 无效")
	}

	// 演示过期 Token
	fmt.Println("\n--- 过期 Token 演示 ---")
	expiredClaims := jwt.MapClaims{
		"user_id":  789,
		"username": "bob",
		"exp":      time.Now().Add(-time.Hour).Unix(), // 1小时前过期
		"iat":      time.Now().Add(-2 * time.Hour).Unix(),
	}
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	expiredTokenString, _ := expiredToken.SignedString([]byte(secretKey))

	_, expiredErr := jwt.Parse(expiredTokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if expiredErr != nil {
		fmt.Printf("❌ 过期 Token 验证失败: %v\n", expiredErr)
	}

	fmt.Println()
}

// 演示 JWT 安全性
func demoJWTSecurity() {
	fmt.Println("4. JWT 安全性演示")

	// 演示篡改 Payload
	fmt.Println("--- Payload 篡改演示 ---")
	originalClaims := jwt.MapClaims{
		"user_id":  123,
		"username": "john_doe",
		"role":     "user",
		"exp":      time.Now().Add(time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}
	originalToken := jwt.NewWithClaims(jwt.SigningMethodHS256, originalClaims)
	secretKey := "your-secret-key"
	originalTokenString, _ := originalToken.SignedString([]byte(secretKey))

	// 尝试篡改 Payload
	parts := strings.Split(originalTokenString, ".")
	tamperedPayload := map[string]interface{}{
		"user_id":  123,
		"username": "john_doe",
		"role":     "admin", // 尝试提升权限
		"exp":      time.Now().Add(time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}
	tamperedPayloadJSON, _ := json.Marshal(tamperedPayload)
	tamperedPayloadBase64 := base64.RawURLEncoding.EncodeToString(tamperedPayloadJSON)

	tamperedToken := parts[0] + "." + tamperedPayloadBase64 + "." + parts[2]

	// 验证篡改后的 Token
	_, tamperedErr := jwt.Parse(tamperedToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if tamperedErr != nil {
		fmt.Printf("❌ 篡改的 Token 验证失败: %v\n", tamperedErr)
		fmt.Println("✅ 这说明 JWT 的签名机制有效防止了篡改!")
	}

	// 演示密钥重要性
	fmt.Println("\n--- 密钥重要性演示 ---")
	wrongKey := "wrong-secret-key"
	_, wrongKeyErr := jwt.Parse(originalTokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(wrongKey), nil
	})

	if wrongKeyErr != nil {
		fmt.Printf("❌ 使用错误密钥验证失败: %v\n", wrongKeyErr)
		fmt.Println("✅ 这说明密钥的安全性很重要!")
	}

	fmt.Println("\n=== 演示结束 ===")
	fmt.Println("\n安全建议:")
	fmt.Println("1. 使用强密钥（至少256位）")
	fmt.Println("2. 定期轮换密钥")
	fmt.Println("3. 设置合理的过期时间")
	fmt.Println("4. 使用 HTTPS 传输")
	fmt.Println("5. 不要在 JWT 中存储敏感信息")
}
