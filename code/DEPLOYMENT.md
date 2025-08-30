# 部署说明

## 数据库自动创建功能

本项目已经实现了数据库自动创建功能，您不再需要手动创建数据库。系统会在启动时自动检查并创建所需的数据库。

## 配置方式

### 1. 配置文件方式（推荐）

在项目根目录创建 `config.json` 文件：

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
    "password": "your_password",
    "dbname": "login_db"
  },
  "jwt": {
    "secret_key": "your-super-secret-jwt-key-change-in-production",
    "expire": 3600
  }
}
```

### 2. 环境变量方式

您也可以通过环境变量来配置：

```bash
# 服务器配置
export SERVER_PORT=8080
export SERVER_MODE=debug

# 数据库配置
export DB_HOST=localhost
export DB_PORT=3306
export DB_USERNAME=root
export DB_PASSWORD=your_password
export DB_NAME=login_db

# JWT配置
export JWT_SECRET_KEY=your-super-secret-jwt-key-change-in-production
export JWT_EXPIRE=3600
```

### 3. 默认配置

如果没有配置文件和环境变量，系统将使用以下默认配置：

- 服务器端口：8080
- 数据库主机：localhost
- 数据库端口：3306
- 数据库用户名：root
- 数据库密码：password
- 数据库名称：login_db

## 部署步骤

### 1. 确保MySQL服务运行

```bash
# Ubuntu/Debian
sudo systemctl start mysql

# CentOS/RHEL
sudo systemctl start mysqld

# Windows
net start mysql
```

### 2. 创建MySQL用户（可选）

如果您不想使用root用户，可以创建专用用户：

```sql
CREATE USER 'login_user'@'localhost' IDENTIFIED BY 'your_password';
GRANT ALL PRIVILEGES ON login_db.* TO 'login_user'@'localhost';
FLUSH PRIVILEGES;
```

### 3. 运行应用

```bash
# 编译
go build -o login main.go

# 运行
./login
```

或者在Windows上：

```cmd
go build -o login.exe main.go
login.exe
```

## 注意事项

1. **数据库权限**：确保MySQL用户有创建数据库的权限
2. **字符集**：数据库将使用 `utf8mb4` 字符集和 `utf8mb4_unicode_ci` 排序规则
3. **安全性**：生产环境中请修改默认密码和JWT密钥
4. **端口**：确保指定的端口没有被其他服务占用

## 故障排除

### 数据库连接失败

1. 检查MySQL服务是否运行
2. 验证用户名和密码是否正确
3. 确认MySQL用户有足够权限

### 端口被占用

修改配置文件中的端口号，或者停止占用端口的服务。

### 权限不足

确保MySQL用户有 `CREATE` 权限：

```sql
GRANT CREATE ON *.* TO 'your_user'@'localhost';
FLUSH PRIVILEGES;
```
