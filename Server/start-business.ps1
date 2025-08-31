# 启动业务服务的PowerShell脚本

Write-Host "正在启动业务服务..." -ForegroundColor Green

# 检查依赖服务
Write-Host "检查MySQL和Redis服务..." -ForegroundColor Blue
docker-compose up -d mysql redis

# 等待数据库启动
Write-Host "等待数据库启动..." -ForegroundColor Blue
Start-Sleep -Seconds 10

# 进入业务服务目录
Set-Location business

# 检查Go环境
$goVersion = go version 2>$null
if ($LASTEXITCODE -ne 0) {
    Write-Host "错误: 未找到Go环境，请先安装Go" -ForegroundColor Red
    exit 1
}

Write-Host "Go版本: $goVersion" -ForegroundColor Cyan

# 下载依赖
Write-Host "正在下载依赖..." -ForegroundColor Blue
go mod tidy

# 设置环境变量
$env:BUSINESS_PORT = "8081"
$env:BUSINESS_MODE = "debug"
$env:DB_HOST = "localhost"
$env:DB_PORT = "3306"
$env:DB_USERNAME = "root"
$env:DB_PASSWORD = "password"
$env:DB_NAME = "login_db"
$env:REDIS_HOST = "localhost"
$env:REDIS_PORT = "6379"

# 启动服务
Write-Host "正在启动业务服务..." -ForegroundColor Blue
go run main.go

Write-Host "业务服务已启动: http://localhost:8081" -ForegroundColor Green
