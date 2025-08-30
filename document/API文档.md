# API接口文档

## 基础信息

- **基础URL**: `http://localhost:8080/api`
- **请求格式**: JSON
- **响应格式**: JSON
- **字符编码**: UTF-8

## 认证方式

系统支持两种认证方式：
1. **JWT Token认证**: 在请求头中添加 `Authorization: Bearer <token>`
2. **无认证**: 部分接口无需认证

## 接口列表

### 1. 健康检查

**接口地址**: `GET /health`

**描述**: 检查服务是否正常运行

**请求参数**: 无

**响应示例**:
```json
{
  "status": "ok",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

### 2. 用户注册

**接口地址**: `POST /register`

**描述**: 用户注册（用户名+密码）

**请求参数**:
```json
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
    "phone": "",
    "created_at": "2024-01-01T12:00:00Z"
  }
}
```

### 3. 用户登录

**接口地址**: `POST /login`

**描述**: 用户名+密码登录

**请求参数**:
```json
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
      "phone": "",
      "created_at": "2024-01-01T12:00:00Z"
    }
  }
}
```

### 4. 发送短信验证码

**接口地址**: `POST /sms/send`

**描述**: 发送短信验证码

**请求参数**:
```json
{
  "phone": "13800138000"
}
```

**响应示例**:
```json
{
  "message": "验证码发送成功",
  "code": "123456"
}
```

**限流规则**:
- 同一手机号1分钟内只能发送1次
- 同一IP地址1分钟内最多发送10次

### 5. 短信验证码登录

**接口地址**: `POST /sms/login`

**描述**: 手机号+验证码登录

**请求参数**:
```json
{
  "phone": "13800138000",
  "code": "123456"
}
```

**响应示例**:
```json
{
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 2,
      "username": "",
      "phone": "13800138000",
      "created_at": "2024-01-01T12:00:00Z"
    }
  }
}
```

**说明**: 
- 如果用户不存在，会自动创建新用户
- 验证码使用后立即失效

### 6. Hello World接口

**接口地址**: `GET /hello`

**描述**: 需要认证的Hello World接口

**请求头**:
```
Authorization: Bearer <token>
```

**响应示例**:
```json
{
  "message": "Hello World!",
  "user_id": 1,
  "username": "testuser"
}
```

### 7. 用户退出

**接口地址**: `POST /logout`

**描述**: 用户退出登录

**请求头**:
```
Authorization: Bearer <token>
```

**响应示例**:
```json
{
  "message": "退出成功"
}
```

## 错误响应格式

所有接口的错误响应都遵循以下格式：

```json
{
  "error": "错误描述信息"
}
```

## 常见HTTP状态码

- `200 OK`: 请求成功
- `201 Created`: 创建成功
- `400 Bad Request`: 请求参数错误
- `401 Unauthorized`: 未认证或认证失败
- `429 Too Many Requests`: 请求过于频繁
- `500 Internal Server Error`: 服务器内部错误

## 使用示例

### 1. 用户名密码登录流程

```bash
# 1. 注册用户
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"123456"}'

# 2. 登录
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"123456"}'

# 3. 访问受保护接口
curl -X GET http://localhost:8080/api/hello \
  -H "Authorization: Bearer <token>"
```

### 2. 短信验证码登录流程

```bash
# 1. 发送验证码
curl -X POST http://localhost:8080/api/sms/send \
  -H "Content-Type: application/json" \
  -d '{"phone":"13800138000"}'

# 2. 验证码登录
curl -X POST http://localhost:8080/api/sms/login \
  -H "Content-Type: application/json" \
  -d '{"phone":"13800138000","code":"123456"}'

# 3. 访问受保护接口
curl -X GET http://localhost:8080/api/hello \
  -H "Authorization: Bearer <token>"
```

## 注意事项

1. **开发环境**: 短信验证码会在响应中返回，生产环境不会返回
2. **验证码有效期**: 5分钟
3. **JWT Token有效期**: 24小时
4. **手机号格式**: 支持中国大陆手机号（11位数字，以1开头）
5. **密码要求**: 最少6位字符
6. **用户名要求**: 3-50位字符
