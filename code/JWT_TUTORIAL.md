# JWT（JSON Web Token）原理详解

## 目录
1. [JWT 基础概念](#1-jwt-基础概念)
2. [JWT 结构详解](#2-jwt-结构详解)
3. [JWT 工作原理](#3-jwt-工作原理)
4. [JWT 在项目中的应用](#4-jwt-在项目中的应用)
5. [JWT 安全性考虑](#5-jwt-安全性考虑)
6. [实际代码演示](#6-实际代码演示)

---

## 1. JWT 基础概念

### 什么是 JWT？

JWT（JSON Web Token）是一种开放标准（RFC 7519），用于在各方之间安全地传输信息作为JSON对象。它通常用于：

- **身份验证（Authentication）**：验证用户身份
- **授权（Authorization）**：控制用户访问权限
- **信息交换**：在各方之间安全地传输信息

### 为什么使用 JWT？

**传统 Session 方式的问题：**
```
用户登录 → 服务器创建 Session → 返回 Session ID → 客户端存储 Session ID
```

**JWT 方式的优势：**
```
用户登录 → 服务器生成 JWT → 返回 JWT → 客户端存储 JWT
```

**主要优势：**
1. **无状态**：服务器不需要存储会话信息
2. **可扩展性**：适合分布式系统
3. **跨域支持**：可以在不同域名间使用
4. **自包含**：Token 本身包含所有必要信息

---

## 2. JWT 结构详解

JWT 由三部分组成，用点（.）分隔：
```
Header.Payload.Signature
```

### 2.1 Header（头部）

包含 Token 的类型和签名算法：

```json
{
  "alg": "HS256",  // 签名算法
  "typ": "JWT"     // Token 类型
}
```

**Base64 编码后：**
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9
```

### 2.2 Payload（载荷）

包含声明（Claims），分为三种类型：

#### 注册声明（Registered Claims）
```json
{
  "iss": "login-system",     // 签发者
  "sub": "user123",          // 主题
  "aud": "web-app",          // 受众
  "exp": 1640995200,         // 过期时间
  "nbf": 1640991600,         // 生效时间
  "iat": 1640991600,         // 签发时间
  "jti": "unique-id"         // JWT ID
}
```

#### 公共声明（Public Claims）
```json
{
  "user_id": 123,
  "username": "john_doe",
  "role": "admin"
}
```

#### 私有声明（Private Claims）
```json
{
  "custom_field": "custom_value"
}
```

**完整的 Payload 示例：**
```json
{
  "user_id": 123,
  "username": "john_doe",
  "exp": 1640995200,
  "iat": 1640991600
}
```

**Base64 编码后：**
```
eyJ1c2VyX2lkIjoxMjMsInVzZXJuYW1lIjoiam9obl9kb2UiLCJleHAiOjE2NDA5OTUyMDAsImlhdCI6MTY0MDk5MTYwMH0
```

### 2.3 Signature（签名）

使用指定的算法对 Header 和 Payload 进行签名：

```
HMACSHA256(
  base64UrlEncode(header) + "." + base64UrlEncode(payload),
  secret_key
)
```

**签名过程：**
1. 将 Header 和 Payload 进行 Base64 编码
2. 用点（.）连接
3. 使用密钥和指定算法进行签名
4. 对签名结果进行 Base64 编码

**最终 JWT 示例：**
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMjMsInVzZXJuYW1lIjoiam9obl9kb2UiLCJleHAiOjE2NDA5OTUyMDAsImlhdCI6MTY0MDk5MTYwMH0.SIGNATURE_HERE
```

---

## 3. JWT 工作原理

### 3.1 生成 JWT 的流程

```mermaid
graph TD
    A[用户登录] --> B[验证用户名密码]
    B --> C[创建 JWT Claims]
    C --> D[使用密钥签名]
    D --> E[返回 JWT Token]
    E --> F[客户端存储 Token]
```

### 3.2 验证 JWT 的流程

```mermaid
graph TD
    A[客户端发送请求] --> B[携带 JWT Token]
    B --> C[服务器接收 Token]
    C --> D[解析 Token 结构]
    D --> E[验证签名]
    E --> F[检查过期时间]
    F --> G[提取用户信息]
    G --> H[处理业务逻辑]
```

### 3.3 详细步骤

#### 生成 JWT：
1. **创建 Claims**：包含用户信息和标准声明
2. **选择算法**：如 HS256、RS256 等
3. **生成签名**：使用密钥对 Header + Payload 签名
4. **组合 Token**：Header.Payload.Signature

#### 验证 JWT：
1. **解析结构**：分离 Header、Payload、Signature
2. **验证签名**：使用相同密钥重新计算签名
3. **检查时间**：验证 exp、nbf、iat
4. **提取信息**：从 Payload 中获取用户信息

---

## 4. JWT 在项目中的应用

### 4.1 登录流程

```go
// 用户登录成功后生成 JWT
func LoginHandler(c *gin.Context) {
    // 1. 验证用户名密码
    // 2. 生成 JWT Token
    token, err := utils.GenerateToken(user.ID, user.Username, secretKey, expireTime)
    // 3. 返回 Token
    c.JSON(200, gin.H{"token": token})
}
```

### 4.2 中间件验证

```go
// JWT 中间件
func JWTAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 从请求头获取 Token
        token := c.GetHeader("Authorization")
        // 2. 解析和验证 Token
        claims, err := utils.ParseToken(token, secretKey)
        // 3. 将用户信息存储到上下文
        c.Set("user_id", claims.UserID)
        c.Set("username", claims.Username)
        c.Next()
    }
}
```

### 4.3 路由保护

```go
// 需要认证的路由
router.GET("/profile", JWTAuthMiddleware(), ProfileHandler)
router.POST("/logout", JWTAuthMiddleware(), LogoutHandler)
```

---

## 5. JWT 安全性考虑

### 5.1 安全最佳实践

#### 密钥管理
```go
// ❌ 错误：硬编码密钥
secretKey := "my-secret-key"

// ✅ 正确：从环境变量获取
secretKey := os.Getenv("JWT_SECRET_KEY")
```

#### Token 过期时间
```go
// 设置合理的过期时间
expireTime := 3600 // 1小时
```

#### 敏感信息
```go
// ❌ 错误：在 JWT 中存储敏感信息
{
  "password": "hashed_password",
  "credit_card": "1234-5678-9012-3456"
}

// ✅ 正确：只存储必要信息
{
  "user_id": 123,
  "username": "john_doe",
  "role": "user"
}
```

### 5.2 常见攻击防护

#### 1. 重放攻击
- 使用 `jti`（JWT ID）确保唯一性
- 设置合理的过期时间
- 使用 HTTPS 传输

#### 2. 暴力破解
- 使用强密钥（至少 256 位）
- 定期轮换密钥
- 限制登录尝试次数

#### 3. XSS 攻击
- 不在客户端存储敏感信息
- 使用 HttpOnly Cookie（如果适用）
- 实施 CSP（Content Security Policy）

---

## 6. 实际代码演示

### 6.1 生成 JWT 示例

```go
package main

import (
    "fmt"
    "time"
    "github.com/golang-jwt/jwt/v5"
)

func main() {
    // 创建 Claims
    claims := jwt.MapClaims{
        "user_id":  123,
        "username": "john_doe",
        "exp":      time.Now().Add(time.Hour * 1).Unix(),
        "iat":      time.Now().Unix(),
    }

    // 创建 Token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    
    // 签名
    secretKey := "your-secret-key"
    tokenString, err := token.SignedString([]byte(secretKey))
    if err != nil {
        panic(err)
    }

    fmt.Printf("JWT Token: %s\n", tokenString)
}
```

### 6.2 验证 JWT 示例

```go
func validateToken(tokenString string) {
    // 解析 Token
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        // 验证签名方法
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte("your-secret-key"), nil
    })

    if err != nil {
        fmt.Printf("Token 解析失败: %v\n", err)
        return
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        fmt.Printf("用户ID: %v\n", claims["user_id"])
        fmt.Printf("用户名: %v\n", claims["username"])
        fmt.Printf("过期时间: %v\n", claims["exp"])
    } else {
        fmt.Println("无效的 Token")
    }
}
```

### 6.3 完整的 JWT 工具类

```go
package utils

import (
    "errors"
    "time"
    "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UserID   int    `json:"user_id"`
    Username string `json:"username"`
    jwt.RegisteredClaims
}

// 生成 Token
func GenerateToken(userID int, username, secretKey string, expireTime int) (string, error) {
    claims := Claims{
        UserID:   userID,
        Username: username,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireTime) * time.Second)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secretKey))
}

// 解析 Token
func ParseToken(tokenString, secretKey string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("无效的签名方法")
        }
        return []byte(secretKey), nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, errors.New("无效的 token")
}
```

---

## 总结

JWT 是一种强大的身份验证和授权机制，具有以下特点：

### 优势：
- ✅ 无状态，适合分布式系统
- ✅ 自包含，减少服务器存储
- ✅ 跨域支持
- ✅ 标准化，易于实现

### 注意事项：
- ⚠️ Token 大小限制
- ⚠️ 无法撤销（需要额外机制）
- ⚠️ 安全性依赖密钥管理
- ⚠️ 需要合理设置过期时间

### 适用场景：
- 微服务架构
- 单页应用（SPA）
- API 认证
- 移动应用

通过合理使用 JWT，可以构建安全、可扩展的身份验证系统。
