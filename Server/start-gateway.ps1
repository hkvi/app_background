# 启动网关服务的PowerShell脚本

Write-Host "正在启动网关服务..." -ForegroundColor Green

# 检查依赖服务
Write-Host "检查Redis服务..." -ForegroundColor Blue
docker-compose up -d redis

# 等待Redis启动
Start-Sleep -Seconds 3

# 进入网关目录
Set-Location gateway

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
$env:GATEWAY_PORT = "8080"
$env:GATEWAY_MODE = "debug"
$env:REDIS_HOST = "localhost"
$env:REDIS_PORT = "6379"
$env:BUSINESS_API_BASE_URL = "http://localhost:8081"

# 启动服务
Write-Host "正在启动网关服务..." -ForegroundColor Blue
go run main.go

Write-Host "网关服务已启动: http://localhost:8080" -ForegroundColor Green
