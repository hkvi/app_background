# 微服务架构 - 网关 + 业务服务

这是一个基于Go语言的微服务架构实现，包含网关服务和业务服务两个独立的服务。

## 架构概览

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   客户端请求    │───▶│   网关服务      │───▶│   业务服务      │
│                 │    │   (Gateway)     │    │   (Business)    │
│   - 认证        │    │   - 路由转发    │    │   - 用户管理    │
│   - 授权        │    │   - JWT鉴权     │    │   - 短信服务    │
│   - 限流        │    │   - 双Token     │    │   - 数据库操作  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │                        │
                              ▼                        ▼
                       ┌─────────────┐         ┌─────────────┐
                       │    Redis    │         │    MySQL    │
                       │   (缓存)    │         │   (数据库)  │
                       └─────────────┘         └─────────────┘
```

## 服务说明

### 网关服务 (Gateway) - 端口 8080
- **职责**: 统一入口、路由转发、JWT认证、限流控制
- **功能**:
  - 双Token认证 (Access Token + Refresh Token)
  - 请求路由和代理
  - 用户认证和授权
  - 限流和安全控制
  - CORS处理

### 业务服务 (Business) - 端口 8081
- **职责**: 核心业务逻辑处理
- **功能**:
  - 用户注册和登录
  - 短信验证码服务
  - 数据库操作
  - 业务数据处理

## 快速启动

### 方式一：Docker Compose (推荐)

1. **启动所有服务**:
   ```powershell
   ./start-all.ps1
   ```

2. **停止所有服务**:
   ```powershell
   ./stop-all.ps1
   ```

### 方式二：独立启动

1. **启动基础服务** (MySQL + Redis):
   ```bash
   docker-compose up -d mysql redis
   ```

2. **启动业务服务**:
   ```powershell
   ./start-business.ps1
   ```

3. **启动网关服务**:
   ```powershell
   ./start-gateway.ps1
   ```

## API 文档

### 基础接口

- **健康检查**: `GET /api/health`
  ```json
  {
    "status": "ok",
    "service": "gateway",
    "timestamp": 1640995200,
    "message": "网关服务运行正常"
  }
  ```

### 认证接口

#### 1. 用户注册
- **URL**: `POST /api/auth/register`
- **请求体**:
  ```json
  {
    "username": "testuser",
    "password": "password123"
  }
  ```

#### 2. 用户登录
- **URL**: `POST /api/auth/login`
- **请求体**:
  ```json
  {
    "username": "testuser",
    "password": "password123"
  }
  ```
- **响应**:
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
        "created_at": "2023-12-31T12:00:00Z"
      }
    }
  }
  ```

#### 3. 发送短信验证码
- **URL**: `POST /api/auth/sms/send`
- **请求体**:
  ```json
  {
    "phone": "13800138000"
  }
  ```

#### 4. 短信验证码登录
- **URL**: `POST /api/auth/sms/login`
- **请求体**:
  ```json
  {
    "phone": "13800138000",
    "code": "201707"
  }
  ```

#### 5. 刷新令牌
- **URL**: `POST /api/auth/refresh`
- **请求体**:
  ```json
  {
    "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
  }
  ```

#### 6. 用户退出
- **URL**: `POST /api/auth/logout`
- **请求头**: `Authorization: Bearer {access_token}`

### 需要认证的接口

所有需要认证的接口都需要在请求头中包含访问令牌：
```
Authorization: Bearer {access_token}
```

## 配置说明

### 网关服务配置 (`gateway/gateway-config.json`)
```json
{
  "server": {
    "port": "8080",
    "mode": "debug"
  },
  "jwt": {
    "access_secret_key": "your-access-secret-key",
    "refresh_secret_key": "your-refresh-secret-key",
    "access_expire": 900,
    "refresh_expire": 86400
  },
  "redis": {
    "host": "localhost",
    "port": "6379",
    "password": "",
    "db": 0
  },
  "business_api": {
    "base_url": "http://localhost:8081",
    "timeout": 30
  }
}
```

### 业务服务配置 (`business/business-config.json`)
```json
{
  "server": {
    "port": "8081",
    "mode": "debug"
  },
  "database": {
    "host": "localhost",
    "port": "3306",
    "username": "root",
    "password": "password",
    "dbname": "login_db"
  },
  "redis": {
    "host": "localhost",
    "port": "6379",
    "password": "",
    "db": 0
  },
  "sms": {
    "access_key_id": "your-access-key-id",
    "access_key_secret": "your-access-key-secret",
    "sign_name": "your-sign-name",
    "template_code": "your-template-code",
    "region_id": "cn-hangzhou"
  }
}
```

## 环境变量

复制 `env.example` 为 `.env` 并配置以下环境变量：

```env
# JWT密钥配置
JWT_ACCESS_SECRET_KEY=your-access-secret-key-change-in-production
JWT_REFRESH_SECRET_KEY=your-refresh-secret-key-change-in-production

# 短信服务配置
SMS_ACCESS_KEY_ID=your-access-key-id
SMS_ACCESS_KEY_SECRET=your-access-key-secret
SMS_SIGN_NAME=your-sign-name
SMS_TEMPLATE_CODE=your-template-code
SMS_REGION_ID=cn-hangzhou
```

## 开发指南

### 项目结构
```
Server/
├── gateway/                 # 网关服务
│   ├── main.go
│   ├── config/
│   ├── handlers/
│   ├── middleware/
│   ├── routes/
│   ├── utils/
│   └── cache/
├── business/               # 业务服务
│   ├── main.go
│   ├── config/
│   ├── handlers/
│   ├── middleware/
│   ├── routes/
│   ├── services/
│   ├── models/
│   ├── database/
│   ├── cache/
│   └── utils/
├── docker-compose.yml      # Docker编排文件
└── *.ps1                  # 启动脚本
```

### 添加新的API

1. **在业务服务中添加处理器**:
   - 在 `business/handlers/` 中创建新的处理器
   - 在 `business/routes/routes.go` 中添加路由

2. **在网关中配置代理**:
   - 网关会自动代理 `/api/business/*` 路径到业务服务
   - 需要认证的接口会自动验证JWT令牌

### 数据库迁移

数据库表会在业务服务启动时自动创建。如需添加新表，在 `business/database/database.go` 的 `CreateTables` 函数中添加。

## 监控和日志

- **查看所有服务日志**: `docker-compose logs -f`
- **查看特定服务日志**: `docker-compose logs -f gateway-service`
- **实时监控**: 所有服务都有详细的请求日志输出

## 安全考虑

1. **生产环境**:
   - 修改所有默认密钥和密码
   - 使用HTTPS
   - 配置防火墙规则
   - 启用日志监控

2. **JWT安全**:
   - Access Token短期有效（15分钟）
   - Refresh Token长期有效（24小时）
   - 支持令牌黑名单机制

3. **限流保护**:
   - 短信接口限流
   - 登录接口限流
   - IP级别限流

## 故障排除

### 常见问题

1. **端口冲突**: 确保8080、8081、3306、6379端口未被占用
2. **Docker问题**: 确保Docker Desktop正在运行
3. **数据库连接**: 等待MySQL完全启动（约10-15秒）
4. **依赖问题**: 运行 `go mod tidy` 更新依赖

### 重置环境

```bash
# 停止并删除所有容器
docker-compose down -v

# 清理镜像
docker system prune -f

# 重新启动
./start-all.ps1
```

## 技术栈

- **语言**: Go 1.21
- **Web框架**: Gin
- **数据库**: MySQL 8.0
- **缓存**: Redis 7
- **认证**: JWT (双Token)
- **容器化**: Docker + Docker Compose
- **短信服务**: 阿里云短信服务

## 许可证

MIT License
