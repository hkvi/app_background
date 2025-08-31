# API接口文档

## 项目概述

本项目采用微服务架构，包含网关服务和业务服务：
- **网关服务** (Gateway): 端口 8080，负责认证、JWT令牌管理、请求代理
- **业务服务** (Business): 端口 8081，负责核心业务逻辑处理

## 基础信息

- **网关服务基础URL**: `http://localhost:8080`
- **业务服务基础URL**: `http://localhost:8081` (内部访问)
- **API版本**: v1
- **数据格式**: JSON
- **字符编码**: UTF-8

## 认证机制

本系统采用JWT双令牌认证机制：
- **Access Token**: 短期令牌，用于API访问认证，默认有效期15分钟
- **Refresh Token**: 长期令牌，用于刷新Access Token，默认有效期24小时

### 认证头格式
```
Authorization: Bearer <access_token>
```

---

## 接口列表

### 1. 健康检查

#### 1.1 网关健康检查
- **URL**: `GET /api/health`
- **描述**: 检查网关服务状态
- **认证**: 无需认证

**响应示例**:
```json
{
  "status": "ok",
  "service": "gateway",
  "timestamp": 1703123456,
  "message": "网关服务运行正常"
}
```

#### 1.2 业务服务健康检查
- **URL**: `GET /api/health` (业务服务)
- **描述**: 检查业务服务状态
- **认证**: 无需认证

**响应示例**:
```json
{
  "status": "ok",
  "service": "business",
  "timestamp": 1703123456,
  "message": "业务服务运行正常"
}
```

---

### 2. 用户认证接口

#### 2.1 用户注册
- **URL**: `POST /api/auth/register`
- **描述**: 用户账号注册
- **认证**: 无需认证
- **限流**: 受登录限流中间件保护

**请求参数**:
```json
{
  "username": "testuser",
  "password": "123456",
  "phone": "13800138000"
}
```

**参数说明**:
| 参数 | 类型 | 必填 | 说明 | 验证规则 |
|------|------|------|------|----------|
| username | string | 是 | 用户名 | 长度3-50字符 |
| password | string | 是 | 密码 | 长度6-100字符 |
| phone | string | 否 | 手机号 | 可选 |

**成功响应** (201):
```json
{
  "message": "注册成功",
  "data": {
    "user": {
      "id": 1,
      "username": "testuser",
      "phone": "13800138000",
      "created_at": "2023-12-01T10:00:00Z"
    }
  }
}
```

**错误响应**:
```json
{
  "error": "用户名已存在"
}
```

#### 2.2 用户登录
- **URL**: `POST /api/auth/login`
- **描述**: 用户账号密码登录
- **认证**: 无需认证
- **限流**: 受登录限流中间件保护

**请求参数**:
```json
{
  "username": "testuser",
  "password": "123456"
}
```

**参数说明**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名 |
| password | string | 是 | 密码 |

**成功响应** (200):
```json
{
  "message": "登录成功",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 900,
    "user": {
      "id": 1,
      "username": "testuser",
      "phone": "13800138000",
      "created_at": "2023-12-01T10:00:00Z"
    }
  }
}
```

**错误响应**:
```json
{
  "error": "用户名或密码错误"
}
```

#### 2.3 刷新令牌
- **URL**: `POST /api/auth/refresh`
- **描述**: 使用refresh token获取新的access token
- **认证**: 无需认证

**请求参数**:
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**成功响应** (200):
```json
{
  "message": "令牌刷新成功",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 900
  }
}
```

#### 2.4 用户退出
- **URL**: `POST /api/auth/logout`
- **描述**: 用户退出登录
- **认证**: 需要认证

**请求头**:
```
Authorization: Bearer <access_token>
```

**成功响应** (200):
```json
{
  "message": "退出成功"
}
```

---

### 3. 短信验证码接口

#### 3.1 发送短信验证码
- **URL**: `POST /api/auth/sms/send`
- **描述**: 发送短信验证码
- **认证**: 无需认证
- **限流**: 受登录限流中间件保护

**请求参数**:
```json
{
  "phone": "13800138000"
}
```

**参数说明**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| phone | string | 是 | 手机号 |

**成功响应** (200):
```json
{
  "message": "验证码发送成功"
}
```

**开发环境响应** (包含验证码):
```json
{
  "message": "验证码发送成功",
  "code": "123456"
}
```

**错误响应**:
```json
{
  "error": "手机号格式错误"
}
```

#### 3.2 短信验证码登录
- **URL**: `POST /api/auth/sms/login`
- **描述**: 使用短信验证码登录
- **认证**: 无需认证
- **限流**: 受登录限流中间件保护

**请求参数**:
```json
{
  "phone": "13800138000",
  "code": "123456"
}
```

**参数说明**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| phone | string | 是 | 手机号 |
| code | string | 是 | 验证码 |

**成功响应** (200):
```json
{
  "message": "登录成功",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 900,
    "user": {
      "id": 1,
      "username": "",
      "phone": "13800138000",
      "created_at": "2023-12-01T10:00:00Z"
    }
  }
}
```

---

### 4. 业务代理接口

#### 4.1 业务服务代理
- **URL**: `ANY /api/business/*path`
- **描述**: 代理到业务服务的所有请求
- **认证**: 需要认证
- **方法**: 支持所有HTTP方法

**请求头**:
```
Authorization: Bearer <access_token>
```

**说明**: 该接口会将请求代理到业务服务，并在请求头中添加用户信息。

---

## 错误码说明

### HTTP状态码

| 状态码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 201 | 创建成功 |
| 400 | 请求参数错误 |
| 401 | 未授权/认证失败 |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 429 | 请求过于频繁 |
| 500 | 服务器内部错误 |

### 业务错误码

| 错误信息 | 说明 |
|----------|------|
| "请求参数错误" | 请求体格式不正确或缺少必填参数 |
| "用户名已存在" | 注册时用户名重复 |
| "用户名或密码错误" | 登录凭据不正确 |
| "手机号格式错误" | 手机号格式不符合要求 |
| "验证码错误或已过期" | 短信验证码不正确或已过期 |
| "请求过于频繁" | 触发限流保护 |
| "无效的刷新令牌" | refresh token无效或过期 |
| "令牌已失效" | access token无效或过期 |

---

## 中间件说明

### 1. CORS中间件
- **作用**: 处理跨域请求
- **应用范围**: 网关服务所有接口

### 2. 日志中间件
- **作用**: 记录请求日志
- **应用范围**: 所有服务的所有接口

### 3. 错误处理中间件
- **作用**: 统一错误响应格式
- **应用范围**: 所有服务的所有接口

### 4. 认证中间件
- **作用**: 验证JWT访问令牌
- **应用范围**: 需要认证的接口

### 5. 登录限流中间件
- **作用**: 限制登录相关接口的请求频率
- **应用范围**: 认证相关接口

---

## 部署信息

### 服务端口
- **网关服务**: 8080
- **业务服务**: 8081
- **MySQL**: 3306
- **Redis**: 6379

### Docker部署
使用docker-compose部署：
```bash
docker-compose up -d
```

### 环境变量
主要环境变量配置：
- `JWT_ACCESS_SECRET_KEY`: JWT访问令牌密钥
- `JWT_REFRESH_SECRET_KEY`: JWT刷新令牌密钥
- `SMS_ACCESS_KEY_ID`: 短信服务AccessKey ID
- `SMS_ACCESS_KEY_SECRET`: 短信服务AccessKey Secret

---

## 开发调试

### 测试命令示例

1. **健康检查**:
```bash
curl -X GET http://localhost:8080/api/health
```

2. **用户注册**:
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"123456","phone":"13800138000"}'
```

3. **用户登录**:
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"123456"}'
```

4. **发送短信验证码**:
```bash
curl -X POST http://localhost:8080/api/auth/sms/send \
  -H "Content-Type: application/json" \
  -d '{"phone":"13800138000"}'
```

5. **短信验证码登录**:
```bash
curl -X POST http://localhost:8080/api/auth/sms/login \
  -H "Content-Type: application/json" \
  -d '{"phone":"13800138000","code":"123456"}'
```

6. **刷新令牌**:
```bash
curl -X POST http://localhost:8080/api/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"your_refresh_token_here"}'
```

7. **用户退出**:
```bash
curl -X POST http://localhost:8080/api/auth/logout \
  -H "Authorization: Bearer your_access_token_here"
```

---

## 版本信息

- **文档版本**: v1.0
- **最后更新**: 2023-12-01
- **Go版本**: 1.21+
- **框架**: Gin Web Framework
- **数据库**: MySQL 8.0+
- **缓存**: Redis 7.0+

---

## 注意事项

1. **安全性**: 
   - 生产环境请使用HTTPS
   - 妥善保管JWT密钥
   - 定期轮换密钥

2. **性能优化**:
   - 合理设置令牌过期时间
   - 使用Redis缓存提升性能
   - 配置合适的限流策略

3. **监控告警**:
   - 监控API响应时间
   - 监控错误率
   - 监控服务健康状态

4. **开发环境**:
   - 短信验证码在开发环境下会返回，生产环境不返回
   - 可通过环境变量控制功能开关
