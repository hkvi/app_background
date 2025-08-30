# Go语言Web后端认证系统

基于Go语言和Gin框架开发的轻量级Web后端认证系统，实现了用户注册、登录、退出和受保护的Hello World接口。

## 功能特性

- ✅ 用户注册和登录
- ✅ JWT无状态认证
- ✅ 密码bcrypt加密
- ✅ 受保护的API接口
- ✅ 优雅的服务器关闭
- ✅ 完整的错误处理
- ✅ CORS跨域支持
- ✅ 请求日志记录

## 技术栈

- **Web框架**: Gin v1.10.1
- **数据库**: MySQL
- **认证**: JWT (JSON Web Token)
- **密码加密**: bcrypt
- **数据库驱动**: go-sql-driver/mysql

## 项目结构

```
code/
├── main.go              # 应用入口
├── config/              # 配置文件
│   └── config.go
├── database/            # 数据库相关
│   └── database.go
├── handlers/            # HTTP处理器
│   ├── health.go
│   ├── user.go
│   └── hello.go
├── middleware/          # 中间件
│   ├── cors.go
│   ├── error.go
│   ├── logger.go
│   └── auth.go
├── models/              # 数据模型
│   └── user.go
├── routes/              # 路由配置
│   └── routes.go
├── services/            # 业务逻辑层
│   └── user_service.go
├── utils/               # 工具函数
│   ├── password.go
│   └── jwt.go
└── go.mod               # 依赖管理
```

## 快速开始

### 1. 环境要求

- Go 1.25.0+
- MySQL 5.7+

### 2. 安装依赖

```bash
go mod tidy
```

### 3. 配置数据库

修改 `config/config.go` 中的数据库配置：

```go
Database: DatabaseConfig{
    Host:     "localhost",
    Port:     "3306",
    Username: "root",
    Password: "your_password",
    DBName:   "login_db",
},
```

### 4. 创建数据库

```sql
CREATE DATABASE login_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 5. 运行应用

```bash
go run main.go
```

服务器将在 `http://localhost:8080` 启动。

## API接口文档

### 基础信息

- **基础URL**: `http://localhost:8080/api`
- **内容类型**: `application/json`
- **认证方式**: Bearer Token (JWT)

### 接口列表

#### 1. 健康检查

```http
GET /api/health
```

**响应示例**:
```json
{
    "status": "ok",
    "message": "Server is running"
}
```

#### 2. 用户注册

```http
POST /api/register
Content-Type: application/json

{
    "username": "testuser",
    "password": "123456"
}
```

**响应示例**:
```json
{
    "message": "注册成功",
    "user": {
        "id": 1,
        "username": "testuser",
        "created_at": "2024-01-01T00:00:00Z"
    }
}
```

#### 3. 用户登录

```http
POST /api/login
Content-Type: application/json

{
    "username": "testuser",
    "password": "123456"
}
```

**响应示例**:
```json
{
    "message": "登录成功",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "user": {
            "id": 1,
            "username": "testuser",
            "created_at": "2024-01-01T00:00:00Z"
        }
    }
}
```

#### 4. Hello World (需要认证)

```http
GET /api/hello
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例**:
```json
{
    "message": "Hello World! 欢迎使用认证系统",
    "user_id": 1,
    "username": "testuser"
}
```

#### 5. 用户退出

```http
POST /api/logout
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例**:
```json
{
    "message": "退出成功"
}
```

## 错误处理

所有接口都遵循统一的错误响应格式：

```json
{
    "error": "错误描述信息"
}
```

常见HTTP状态码：
- `200`: 请求成功
- `201`: 创建成功
- `400`: 请求参数错误
- `401`: 未认证或认证失败
- `500`: 服务器内部错误

## 安全特性

1. **密码安全**: 使用bcrypt算法对密码进行哈希存储
2. **JWT安全**: 使用强密钥签名，支持过期时间设置
3. **输入验证**: 对用户输入进行严格验证
4. **SQL注入防护**: 使用参数化查询
5. **CORS配置**: 支持跨域请求

## 开发说明

### 添加新接口

1. 在 `handlers/` 目录下创建处理器
2. 在 `routes/routes.go` 中添加路由
3. 如需认证，使用 `middleware.AuthMiddleware()`

### 数据库操作

1. 在 `services/` 目录下添加业务逻辑
2. 使用 `database.DB` 进行数据库操作
3. 使用参数化查询防止SQL注入

### 配置管理

修改 `config/config.go` 中的配置项，支持：
- 服务器端口和模式
- 数据库连接参数
- JWT密钥和过期时间

## 部署建议

### 生产环境

1. 修改JWT密钥为强密钥
2. 配置HTTPS
3. 设置合适的数据库连接池参数
4. 使用环境变量管理敏感配置
5. 配置日志记录和监控

### Docker部署

```dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

## 许可证

MIT License
