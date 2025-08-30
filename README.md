# Go语言Web后端项目

一个基于Go语言的轻量级Web后端系统，实现用户认证和授权功能，支持多种登录方式。

## 功能特性

### 用户管理
- ✅ 用户注册（用户名+密码）
- ✅ 用户登录（用户名+密码）
- ✅ 用户登录（手机号+验证码）
- ✅ 用户退出
- ✅ 受保护的Hello World接口

### 验证码服务
- ✅ 短信验证码发送
- ✅ 验证码验证和过期管理
- ✅ 防刷机制（IP限流和手机号限流）

### 安全特性
- ✅ JWT无状态认证
- ✅ bcrypt密码哈希
- ✅ 输入验证和SQL注入防护
- ✅ 限流中间件
- ✅ CORS跨域支持

## 技术栈

- **Web框架**: Gin
- **数据库**: MySQL
- **缓存**: Redis
- **认证**: JWT
- **密码安全**: bcrypt
- **短信服务**: 阿里云短信服务

## 快速开始

### 环境要求

- Go 1.25.0+
- MySQL 8.0+
- Redis 6.0+

### 安装和运行

1. **克隆项目**
```bash
git clone <项目地址>
cd app_background
```

2. **安装依赖**
```bash
cd code
go mod download
```

3. **配置数据库和Redis**
```bash
# 启动MySQL和Redis服务
sudo systemctl start mysqld
sudo systemctl start redis
```

4. **修改配置文件**
编辑 `code/config.json`，配置数据库、Redis和短信服务参数。

5. **运行项目**
```bash
go run main.go
```

服务将在 `http://localhost:8080` 启动。

## API接口

### 基础接口
- `GET /api/health` - 健康检查

### 用户认证
- `POST /api/register` - 用户注册
- `POST /api/login` - 用户名密码登录
- `POST /api/sms/send` - 发送短信验证码
- `POST /api/sms/login` - 短信验证码登录
- `POST /api/logout` - 用户退出

### 受保护接口
- `GET /api/hello` - Hello World接口（需要认证）

详细的API文档请参考 [API文档](document/API文档.md)

## 项目结构

```
app_background/
├── code/                    # 源代码
│   ├── config/             # 配置管理
│   ├── database/           # 数据库相关
│   ├── handlers/           # HTTP处理器
│   ├── middleware/         # 中间件
│   ├── models/             # 数据模型
│   ├── routes/             # 路由配置
│   ├── services/           # 业务服务层
│   ├── utils/              # 工具函数
│   ├── cache/              # 缓存相关
│   ├── config.json         # 配置文件
│   ├── go.mod              # Go模块文件
│   └── main.go             # 应用入口
├── document/               # 文档
│   ├── API文档.md          # API接口文档
│   ├── 部署文档.md         # 部署指南
│   ├── 技术选型文档.md     # 技术选型说明
│   └── JWT_TUTORIAL.md     # JWT教程
└── README.md               # 项目说明
```

## 使用示例

### 用户名密码登录

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

### 短信验证码登录

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

## 配置说明

### 环境变量

支持通过环境变量覆盖配置：

```bash
export SERVER_PORT=8080
export DB_HOST=localhost
export DB_USERNAME=root
export DB_PASSWORD=password
export DB_NAME=login_db
export REDIS_HOST=localhost
export REDIS_PORT=6379
export JWT_SECRET_KEY=your-secret-key
export SMS_ACCESS_KEY_ID=your-access-key-id
export SMS_ACCESS_KEY_SECRET=your-access-key-secret
```

### 配置文件

`config.json` 配置示例：

```json
{
  "server": {
    "port": "8080",
    "mode": "debug"
  },
  "database": {
    "host": "localhost",
    "port": "3306",
    "username": "root",
    "password": "password",
    "dbname": "login_db"
  },
  "jwt": {
    "secret_key": "your-secret-key",
    "expire": 3600
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

## 部署

详细的部署指南请参考 [部署文档](document/部署文档.md)

### Docker部署

```bash
# 使用docker-compose快速部署
docker-compose up -d
```

### 生产环境部署

```bash
# 编译
go build -o app main.go

# 使用systemd管理服务
sudo systemctl start login-app
sudo systemctl enable login-app
```

## 开发

### 代码规范

- 遵循Go官方代码规范
- 使用gofmt格式化代码
- 添加必要的注释和文档

### 测试

```bash
# 运行测试
go test ./...

# 运行特定测试
go test ./handlers
```

### 日志

项目使用内置的日志系统，支持不同级别的日志输出：

- INFO: 一般信息
- ERROR: 错误信息
- DEBUG: 调试信息（仅在debug模式下输出）

## 安全考虑

- 使用bcrypt进行密码哈希
- JWT token设置合理的过期时间
- 验证码有效期限制（5分钟）
- 防刷机制（IP限流和手机号限流）
- 输入验证防止SQL注入和XSS攻击

## 性能优化

- 数据库连接池配置
- Redis缓存验证码
- 合理的索引设计
- 限流保护

## 贡献

欢迎提交Issue和Pull Request！

## 许可证

MIT License

## 联系方式

如有问题，请提交Issue或联系开发团队。
